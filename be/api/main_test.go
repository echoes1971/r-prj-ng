package api

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"rprj/be/dblayer"
	"rprj/be/models"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

var AppConfig models.Config
var testAdminLogin string
var testAdminPwd string
var testUser dblayer.DBEntityInterface

// TestMain is executed before running tests
func TestMain(m *testing.M) {

	err := models.LoadConfig("../config.json", &AppConfig)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	dblayer.InitDBLayer(AppConfig)
	log.Println("DB initialized for tests")

	repo := SetupTestRepo(nil, "-1", []string{"-2"}, AppConfig.TablePrefix)
	// Create test admin user if not exists: test0000 / pass0000 with 4 random digits
	randomDigits := Random4digits()
	testAdminLogin = "test" + randomDigits
	testAdminPwd = "pass" + randomDigits

	// Search for existing user with the same login
	searchUser := repo.GetInstanceByTableName("users")
	searchUser.SetValue("login", testAdminLogin)
	foundUsers, err := repo.Search(searchUser, false, false, "")
	if err != nil {
		log.Fatalf("Failed to search for test admin user: %v", err)
	}
	if len(foundUsers) == 0 {
		log.Printf("Test admin user '%s' does not exist, creating it...\n", testAdminLogin)
		testUser, err = repo.CreateObject("users", map[string]any{
			"login":    testAdminLogin,
			"pwd":      testAdminPwd,
			"fullname": "Test User " + randomDigits,
		}, map[string]any{
			"group_ids": []string{"-1"}, // Admin group
		})
		if err != nil {
			log.Fatalf("Failed to create test admin user: %v", err)
		}
		log.Printf("Created test admin user: login='%s' pwd='%s'\n", testUser.GetValue("login"), "pass"+randomDigits)
	} else {
		log.Printf("Test admin user '%s' already exists, using existing user\n", testAdminLogin)
		testUser = foundUsers[0]
	}

	// Run tests
	code := m.Run()

	// Teardown: delete test admin user
	deletedUser, err := repo.Delete(testUser)
	log.Print("Deleted test admin user=", deletedUser.ToString())

	if err != nil {
		log.Printf("Failed to delete test admin user: %v\n", err)
	} else {
		log.Println("Deleted test admin user")
	}

	// Teardown: close connection
	dblayer.CloseDBConnection()

	// Exit with test code
	os.Exit(code)
}

/* ***** Helper functions for tests ***** */

func RandInt(min, max int) int {
	return min + rand.Intn(max-min)
}

/* Returns a random 4-digit string */
func Random4digits() string {
	const digits = "0123456789"
	result := make([]byte, 4)
	// Generate random number between 0000 and 9999
	for i := 0; i < 4; i++ {
		result[i] = digits[RandInt(0, len(digits))]
	}
	return string(result)
}
func SetupTestRepo(t *testing.T, user_id string, group_ids []string, schema string) *dblayer.DBRepository {
	dbContext := &dblayer.DBContext{
		UserID:   user_id,
		GroupIDs: group_ids,
		Schema:   schema,
	}
	repo := dblayer.NewDBRepository(dbContext, dblayer.Factory, dblayer.DbConnection)
	repo.Verbose = false
	return repo
}

// 2 consecutive logins creates the same token and an error on primary key violation
var tokens = make(map[string]string)

func ApiTestDoLogin(t *testing.T, login, pwd string) string {
	if token, exists := tokens[login+pwd]; exists {
		return token
	}

	creds := Credentials{
		Login: login,
		Pwd:   pwd,
	}
	body, _ := json.Marshal(creds)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LoginHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("DoLogin: wrong status code: got %v, want %v", rr.Code, http.StatusOK)
	}

	var resp TokenResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("DoLogin: error parsing JSON response: %v", err)
	}

	if resp.AccessToken == "" {
		t.Fatalf("DoLogin: access_token missing in response")
	}

	tokens[login+pwd] = resp.AccessToken

	return resp.AccessToken
}

func ApiTestDecodeAccessToken(t *testing.T, tokenString string) jwt.MapClaims {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return JWTKey, nil
	})
	if err != nil || !token.Valid {
		t.Fatalf("DecodeAccessToken: invalid token: %v", err)
	}
	// log.Printf("Decoded claims: %+v\n", claims)

	return claims
}

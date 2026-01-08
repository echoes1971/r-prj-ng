package api

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"rprj/be/dblayer"
	"rprj/be/models"

	"github.com/golang-jwt/jwt/v5"
)

// default JWT key; replace with a secure value loaded from your app configuration at startup
var JWTKey = []byte("change-me-secret")

type PingResponse struct {
	Ping string `json:"ping"`
}

// PingHandler godoc
// @Summary Health check
// @Description Returns pong if server is alive
// @Tags health
// @Produce json
// @Success 200 {object} PingResponse
// @Router /ping [get]
func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(PingResponse{Ping: "Pong"})
}

type Credentials struct {
	Login string `json:"login" example:"admin"`
	Pwd   string `json:"pwd" example:"password123"`
}

type TokenResponse struct {
	AccessToken string   `json:"access_token"`
	ExpiresAt   int64    `json:"expires_at"`
	UserID      string   `json:"user_id"`
	Groups      []string `json:"groups"`
}

/* *** DBFiles *** */
var dbFiles_root_directory string = "."
var dbFiles_dest_directory string = "files"

// OAuth config values (set by InitAPI)
var GoogleClientID string
var GoogleClientSecret string
var GoogleRedirectURL string
var GitHubClientID string
var GitHubClientSecret string
var GitHubRedirectURL string
var TelegramBotToken string
var TelegramBotID string

func InitAPI(config models.Config) {
	JWTKey = []byte(config.JWTSecret)
	dbFiles_root_directory = config.RootDirectory
	dbFiles_dest_directory = config.FilesDirectory
	GoogleClientID = config.GoogleClientID
	GoogleClientSecret = config.GoogleClientSecret
	GoogleRedirectURL = config.GoogleRedirectURL
	GitHubClientID = config.GitHubClientID
	GitHubClientSecret = config.GitHubClientSecret
	GitHubRedirectURL = config.GitHubRedirectURL
	TelegramBotToken = config.TelegramBotToken
	TelegramBotID = config.TelegramBotID
	log.Print("API initialized with JWT key from config")
}

// LoginHandler godoc
// @Summary User login
// @Description Authenticate user and receive JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body Credentials true "Login credentials"
// @Success 200 {object} map[string]interface{} "token and user info"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Router /login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		RespondSimpleError(w, ErrInvalidRequest, "Invalid request format", http.StatusBadRequest)
		return
	}

	dbContext := &dblayer.DBContext{
		UserID:   "-1",           // DANGEROUS!!!! Think of something better here!!!
		GroupIDs: []string{"-2"}, // Same here!!!
		Schema:   dblayer.DbSchema,
	}

	repo := dblayer.NewDBRepository(dbContext, dblayer.Factory, dblayer.DbConnection)
	repo.Verbose = false

	user := repo.GetInstanceByTableName("users")
	if user == nil {
		RespondSimpleError(w, ErrInternalServer, "Failed to create user instance", http.StatusInternalServerError)
		return
	}
	user.SetValue("login", creds.Login)
	foundUsers, err := repo.Search(user, false, false, "")
	if err != nil || len(foundUsers) == 0 {
		RespondSimpleError(w, ErrUnauthorized, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	foundUser, ok := foundUsers[0].(*dblayer.DBUser)
	if !ok {
		RespondSimpleError(w, ErrUnauthorized, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Verify password (supports both encrypted and legacy unencrypted passwords)
	if !foundUser.VerifyPassword(creds.Pwd) {
		RespondSimpleError(w, ErrUnauthorized, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Get User groups
	userGroupsInstance := repo.GetInstanceByTableName("users_groups")
	if userGroupsInstance == nil {
		RespondSimpleError(w, ErrInternalServer, "Failed to create user-groups instance", http.StatusInternalServerError)
		return
	}
	userGroupsInstance.SetValue("user_id", foundUser.GetValue("id"))
	userGroups, err := repo.Search(userGroupsInstance, false, false, "")
	if err != nil {
		RespondSimpleError(w, ErrInternalServer, "Failed to get user groups: "+err.Error(), http.StatusInternalServerError)
		return
	}
	group_list := []string{}
	for _, ug := range userGroups {
		ugEntry := ug
		group_list = append(group_list, ugEntry.GetValue("group_id").(string))
	}

	if foundUser.GetValue("group_id") != "" && slices.Index(group_list, foundUser.GetValue("group_id").(string)) < 0 {
		group_list = append(group_list, foundUser.GetValue("group_id").(string))
	}

	// Genera JWT
	expiration := time.Now().Add(1 * time.Hour)
	claims := &jwt.MapClaims{
		"user_id": foundUser.GetValue("id"),
		"login":   foundUser.GetValue("login"),
		"groups":  strings.Join(group_list, ","),
		"exp":     expiration.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWTKey)
	if err != nil {
		RespondSimpleError(w, ErrInternalServer, "Could not generate token", http.StatusInternalServerError)
		return
	}

	// Salva token in tabella oauth_tokens
	if err := SaveToken(repo, foundUser.GetValue("id").(string), tokenString, expiration.Unix()); err != nil {
		log.Print("Error saving token:", err)
		RespondSimpleError(w, ErrInternalServer, "Could not save token", http.StatusInternalServerError)
		return
	}

	// Risposta al client
	resp := TokenResponse{
		AccessToken: tokenString,
		ExpiresAt:   expiration.Unix(),
		UserID:      foundUser.GetValue("id").(string),
		Groups:      group_list,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// LogoutHandler godoc
// @Summary User logout
// @Description Invalidate the current JWT token. The token must be provided in the Authorization header.
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string "logout message"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /logout [post]
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	tokenString, err := GetTokenFromRequest(r)
	if err != nil {
		RespondSimpleError(w, ErrUnauthorized, "Missing or invalid token", http.StatusUnauthorized)
		return
	}

	claims, err := GetClaimsFromRequest(r)
	if err != nil {
		RespondSimpleError(w, ErrUnauthorized, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("Logging out user %s with token %s", claims["user_id"], tokenString)

	dbContext := &dblayer.DBContext{
		UserID:   claims["user_id"],
		GroupIDs: strings.Split(claims["groups"], ","),
		Schema:   dblayer.DbSchema,
	}

	repo := dblayer.NewDBRepository(dbContext, dblayer.Factory, dblayer.DbConnection)
	repo.Verbose = false

	// Elimina il token dalla tabella oauth_tokens
	if err := DeleteToken(repo, tokenString); err != nil {
		log.Print("Error deleting token:", err)
		RespondSimpleError(w, ErrInternalServer, "Could not delete token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

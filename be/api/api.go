package api

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	"rprj/be/dblayer"

	"github.com/golang-jwt/jwt/v5"
)

// default JWT key; replace with a secure value loaded from your app configuration at startup
var JWTKey = []byte("change-me-secret")

type PingResponse struct {
	Ping string `json:"ping"`
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(PingResponse{Ping: "Pong"})
}

type Credentials struct {
	Login string `json:"login"`
	Pwd   string `json:"pwd"`
}

type TokenResponse struct {
	AccessToken string   `json:"access_token"`
	ExpiresAt   int64    `json:"expires_at"`
	UserID      string   `json:"user_id"`
	Groups      []string `json:"groups"`
}

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
	if !ok || foundUser.GetUnencryptedPwd() != creds.Pwd {
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

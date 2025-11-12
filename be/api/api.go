package api

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	"rprj/be/db"

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
	Groups      []string `json:"groups"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Verifica utente nel DB
	user, err := db.GetUserByLogin(creds.Login)
	if err != nil || user == nil || user.Pwd != creds.Pwd {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Retrieve user groups
	groups, err := db.GetUserGroupsByUserID(user.ID)
	if err != nil {
		http.Error(w, "could not retrieve user groups", http.StatusInternalServerError)
		return
	}
	group_list := []string{}
	for _, g := range groups {
		group_list = append(group_list, g.GroupID)
	}
	if user.GroupID != "" && slices.Index(group_list, user.GroupID) < 0 {
		group_list = append(group_list, user.GroupID)
	}

	// Genera JWT
	expiration := time.Now().Add(1 * time.Hour)
	claims := &jwt.MapClaims{
		"user_id": user.ID,
		"login":   user.Login,
		"groups":  strings.Join(group_list, ","),
		"exp":     expiration.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWTKey)
	if err != nil {
		http.Error(w, "could not generate token", http.StatusInternalServerError)
		return
	}

	// Salva token in tabella oauth_tokens
	if err := db.SaveToken(user.ID, tokenString, expiration.Unix()); err != nil {
		http.Error(w, "could not save token", http.StatusInternalServerError)
		return
	}

	// Risposta al client
	resp := TokenResponse{
		AccessToken: tokenString,
		ExpiresAt:   expiration.Unix(),
		Groups:      group_list,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

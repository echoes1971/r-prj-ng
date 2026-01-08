package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"rprj/be/dblayer"

	"github.com/golang-jwt/jwt/v5"
)

// TelegramOAuthCallback handles Telegram Login Widget callback
// Verifies data hash using HMAC-SHA256 with bot token
func TelegramOAuthCallback(w http.ResponseWriter, r *http.Request) {
	if TelegramBotToken == "" {
		RespondSimpleError(w, ErrInternalServer, "Telegram OAuth not configured", http.StatusInternalServerError)
		return
	}

	// Parse query parameters
	query := r.URL.Query()
	authDate := query.Get("auth_date")
	firstName := query.Get("first_name")
	lastName := query.Get("last_name")
	id := query.Get("id")
	username := query.Get("username")
	photoURL := query.Get("photo_url")
	hash := query.Get("hash")

	if id == "" || hash == "" || authDate == "" {
		RespondSimpleError(w, ErrInvalidRequest, "Missing required Telegram data", http.StatusBadRequest)
		return
	}

	// Verify hash
	if !verifyTelegramHash(query, TelegramBotToken) {
		log.Printf("TelegramOAuthCallback: invalid hash")
		RespondSimpleError(w, ErrUnauthorized, "Invalid Telegram hash", http.StatusUnauthorized)
		return
	}

	// Check auth_date freshness (optional: reject if older than 1 day)
	authDateInt, _ := strconv.ParseInt(authDate, 10, 64)
	if time.Now().Unix()-authDateInt > 86400 {
		RespondSimpleError(w, ErrUnauthorized, "Telegram auth data expired", http.StatusUnauthorized)
		return
	}

	// Build identifier (prefer username, fallback to telegram:id)
	identifier := username
	if identifier == "" {
		identifier = fmt.Sprintf("telegram:%s", id)
	}
	fullname := strings.TrimSpace(firstName + " " + lastName)
	if fullname == "" {
		fullname = identifier
	}

	// Find or create user
	dbContext := &dblayer.DBContext{UserID: "-1", GroupIDs: []string{"-2"}, Schema: dblayer.DbSchema}
	repo := dblayer.NewDBRepository(dbContext, dblayer.Factory, dblayer.DbConnection)
	repo.Verbose = false

	userInst := repo.GetInstanceByTableName("users")
	if userInst == nil {
		RespondSimpleError(w, ErrInternalServer, "Users instance not available", http.StatusInternalServerError)
		return
	}

	userInst.SetValue("login", identifier)
	found, err := repo.Search(userInst, false, false, "")
	if err != nil {
		RespondSimpleError(w, ErrInternalServer, "Search failed", http.StatusInternalServerError)
		return
	}

	var userID string
	if len(found) > 0 {
		userID = found[0].GetValue("id").(string)
	} else {
		userInst.SetValue("login", identifier)
		userInst.SetValue("fullname", fullname)
		userInst.SetValue("pwd", "")
		userInst.SetValue("group_id", "-4")
		userInst.SetValue("permissions", "rwx------")
		created, err := repo.Insert(userInst)
		if err != nil {
			log.Printf("TelegramOAuthCallback: failed to create user: %v", err)
			RespondSimpleError(w, ErrInternalServer, "Failed to create user", http.StatusInternalServerError)
			return
		}
		userID = created.GetValue("id").(string)
	}

	// Collect groups
	userGroupsInstance := repo.GetInstanceByTableName("users_groups")
	group_list := []string{}
	if userGroupsInstance != nil {
		userGroupsInstance.SetValue("user_id", userID)
		ugs, _ := repo.Search(userGroupsInstance, false, false, "")
		for _, ug := range ugs {
			group_list = append(group_list, ug.GetValue("group_id").(string))
		}
	}

	// Ensure primary group included
	u := repo.GetInstanceByTableName("users")
	u.SetValue("id", userID)
	foundUsers, _ := repo.Search(u, false, false, "")
	if len(foundUsers) > 0 {
		if g := foundUsers[0].GetValue("group_id"); g != nil {
			gp := g.(string)
			present := false
			for _, v := range group_list {
				if v == gp {
					present = true
				}
			}
			if !present {
				group_list = append(group_list, gp)
			}
		}
	}

	expiration := time.Now().Add(1 * time.Hour)
	claims := &jwt.MapClaims{
		"user_id": userID,
		"login":   identifier,
		"groups":  stringsJoin(group_list, ","),
		"exp":     expiration.Unix(),
	}
	tokenJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := tokenJwt.SignedString(JWTKey)
	if err != nil {
		RespondSimpleError(w, ErrInternalServer, "Could not generate token", http.StatusInternalServerError)
		return
	}

	_ = SaveToken(repo, userID, tokenString, expiration.Unix())

	// Build payload for frontend
	groupsCSV := ""
	for i, g := range group_list {
		if i > 0 {
			groupsCSV += ","
		}
		groupsCSV += g
	}
	payload := map[string]interface{}{
		"provider":     "telegram",
		"access_token": tokenString,
		"expires_at":   expiration.Unix(),
		"user_id":      userID,
		"login":        identifier,
		"groups":       groupsCSV,
		"photo_url":    photoURL,
	}
	payloadJSON, _ := json.Marshal(payload)
	html := fmt.Sprintf(`<!doctype html><html><head><meta charset="utf-8"></head><body><script>
    (function(){
      try {
        var data = %s;
        if (window.opener && window.opener !== window) {
          window.opener.postMessage(data, "*");
          window.close();
        } else {
          localStorage.setItem('token', data.access_token);
          localStorage.setItem('expires_at', data.expires_at);
          localStorage.setItem('user_id', data.user_id);
          if (data.login) localStorage.setItem('username', data.login);
          if (data.groups) localStorage.setItem('groups', JSON.stringify(data.groups.split(',')));
          window.location.href = '/';
        }
      } catch(e) { console.error(e); window.location.href = '/'; }
    })();
    </script></body></html>`, payloadJSON)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

// verifyTelegramHash verifies Telegram Login Widget data hash
func verifyTelegramHash(params map[string][]string, botToken string) bool {
	receivedHash := ""
	if h, ok := params["hash"]; ok && len(h) > 0 {
		receivedHash = h[0]
	}
	if receivedHash == "" {
		return false
	}

	// Build data-check-string
	keys := []string{}
	for k := range params {
		if k != "hash" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var dataCheckParts []string
	for _, k := range keys {
		if vals, ok := params[k]; ok && len(vals) > 0 {
			dataCheckParts = append(dataCheckParts, fmt.Sprintf("%s=%s", k, vals[0]))
		}
	}
	dataCheckString := strings.Join(dataCheckParts, "\n")

	// secret_key = SHA256(bot_token)
	secretKeyHash := sha256.Sum256([]byte(botToken))
	secretKey := secretKeyHash[:]

	// HMAC-SHA256(data_check_string, secret_key)
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(dataCheckString))
	calculatedHash := hex.EncodeToString(mac.Sum(nil))

	return calculatedHash == receivedHash
}

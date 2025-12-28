package api

import (
	"log"
	"net/http"
	"strings"
	"time"

	"rprj/be/dblayer"

	"github.com/golang-jwt/jwt/v5"
)

func GetTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", http.ErrNoCookie
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", http.ErrNoCookie
	}
	tokenString := parts[1]
	return tokenString, nil
}

func GetClaimsFromRequest(r *http.Request) (map[string]string, error) {

	tokenString, err := GetTokenFromRequest(r)
	if err != nil {
		return nil, err
	}

	dbContext := &dblayer.DBContext{
		UserID:   "-1",           // DANGEROUS!!!! Think of something better here!!!
		GroupIDs: []string{"-2"}, // Same here!!!
		Schema:   dblayer.DbSchema,
	}
	repo := dblayer.NewDBRepository(dbContext, dblayer.Factory, dblayer.DbConnection)
	repo.Verbose = false

	// Validate the token
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return JWTKey, nil
	})
	if err != nil || !token.Valid {
		log.Print("Deleting token from db due to invalidity.")
		DeleteToken(repo, tokenString)
		return nil, http.ErrNoCookie
	}

	// Extract claims as map[string]string
	result := make(map[string]string)
	for key, value := range claims {
		if strVal, ok := value.(string); ok {
			result[key] = strVal
		}
	}

	return result, nil
}

// Middleware check token JWT
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			RespondSimpleError(w, ErrMissingAuthorization, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// it must be in the format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			RespondSimpleError(w, ErrInvalidToken, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]
		// log.Printf("Token ricevuto: %s\n", tokenString)

		dbContext := &dblayer.DBContext{
			UserID:   "-1",           // DANGEROUS!!!! Think of something better here!!!
			GroupIDs: []string{"-2"}, // Same here!!!
			Schema:   dblayer.DbSchema,
		}
		repo := dblayer.NewDBRepository(dbContext, dblayer.Factory, dblayer.DbConnection)
		repo.Verbose = false

		// Valida il token
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return JWTKey, nil
		})
		// log.Printf("Claims: %+v\n", claims)
		// log.Printf("err: %v\n", err)
		if err != nil || !token.Valid {
			RespondSimpleError(w, ErrInvalidToken, "Invalid or expired token", http.StatusUnauthorized)
			log.Print("Deleting token from db due to invalidity.")
			DeleteToken(repo, tokenString)
			return
		}

		// Retrieve user ID from claims
		userID := claims["user_id"].(string)
		// log.Printf("User ID authenticated: %s\n", userID)

		// Retrieve group ids from claims
		// groupIDs := []string{}
		// if g, ok := claims["groups"].(string); ok && g != "" {
		// 	groupIDs = strings.Split(g, ",")
		// }
		// log.Printf("Group IDs: %+v\n", groupIDs)

		// Search the token in the database to ensure it's valid
		if !IsTokenValid(repo, tokenString, userID) {
			RespondSimpleError(w, ErrInvalidToken, "Token not recognized", http.StatusUnauthorized)
			log.Print("Token not found in the database")
			return
		}

		// Passa la richiesta all'handler successivo
		next.ServeHTTP(w, r)
	})
}

func SaveToken(repo *dblayer.DBRepository, userID string, tokenString string, expiry int64) error {

	dbOAuthToken := repo.GetInstanceByTableName("oauth_tokens")
	if dbOAuthToken == nil {
		log.Println("Errore creazione istanza oauth_tokens")
		return nil
	}
	dbOAuthToken.SetValue("user_id", userID)
	dbOAuthToken.SetValue("token_id", tokenString)
	dbOAuthToken.SetValue("access_token", tokenString)
	dbOAuthToken.SetValue("expires_at", time.Unix(expiry, 0))

	_, err := repo.Insert(dbOAuthToken)
	return err
}

func IsTokenValid(repo *dblayer.DBRepository, tokenString string, userID string) bool {

	search := repo.GetInstanceByTableName("oauth_tokens")
	if search == nil {
		log.Println("Errore creazione istanza oauth_tokens")
		return false
	}
	search.SetValue("token_id", tokenString)
	search.SetValue("user_id", userID)

	results, err := repo.Search(search, false, false, "")
	if err != nil {
		log.Println("Errore verifica token:", err)
		return false
	}
	return len(results) > 0
}

func DeleteToken(repo *dblayer.DBRepository, tokenString string) error {
	search := repo.GetInstanceByTableName("oauth_tokens")
	if search == nil {
		log.Println("Errore creazione istanza oauth_tokens")
		return nil
	}
	search.SetValue("token_id", tokenString)
	_, err := repo.Delete(search)
	if err != nil {
		log.Println("Errore cancellazione token:", err)
	}
	return err
}

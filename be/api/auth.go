package api

import (
	"log"
	"net/http"
	"strings"

	"rprj/be/db"

	"github.com/golang-jwt/jwt/v5"
)

// Middleware che controlla il token JWT
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Legge l'header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Deve essere nel formato "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid Authorization header", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]
		log.Printf("Token ricevuto: %s\n", tokenString)

		// Valida il token
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return JWTKey, nil
		})
		log.Printf("Claims: %+v\n", claims)
		log.Printf("err: %v\n", err)
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			log.Print("Deleting token from db due to invalidity.")
			db.DeleteToken(tokenString)
			return
		}

		// Retrieve user ID from claims
		userID := claims["user_id"].(string)
		log.Printf("User ID autenticato: %s\n", userID)

		// Retrieve group ids from claims
		groupIDs := []string{}
		if g, ok := claims["groups"].(string); ok && g != "" {
			groupIDs = strings.Split(g, ",")
		}
		log.Printf("Group IDs: %+v\n", groupIDs)

		// Search the token in the database to ensure it's valid
		if !db.IsTokenValid(tokenString, userID) {
			http.Error(w, "token not recognized", http.StatusUnauthorized)
			log.Print("Token not found in the database")
			return
		}

		// Passa la richiesta all'handler successivo
		next.ServeHTTP(w, r)
	})
}

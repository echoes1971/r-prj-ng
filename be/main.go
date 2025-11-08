package main

/*

Test:

curl -X POST http://localhost:1971/login \
  -H "Content-Type: application/json" \
  -d '{"login":"roberto","pwd":"echoestrade"}'

curl -X GET http://localhost:1971/users/316 \
  -H "Authorization: Bearer <access_token>"

curl -X GET http://localhost:1971/users/ \
  -H "Authorization: Bearer <access_token>"


Test Suites:

# Esegue tutti i test dei vari package
go test -v ./...

go test -v ./api
go test -v ./api -run TestLoginHandler

go clean -testcache
# Se ho funzioni BenchmarkXxx
go test -bench ./api
*/

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"rprj/be/api"
	"rprj/be/db"
	"rprj/be/models"

	"github.com/gorilla/mux"
)

var AppConfig models.Config

// Struttura di esempio per la risposta JSON
type Response struct {
	Message string `json:"message"`
}

func main() {

	configFile := "config.json"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	err := models.LoadConfig(configFile, &AppConfig)
	if err != nil {
		log.Fatalf("Errore caricamento configurazione: %v", err)
	}

	fmt.Printf("Configurazione caricata: %+v\n", AppConfig)

	// Passa la config ai pacchetti
	api.JWTKey = []byte(AppConfig.JWTSecret)
	db.Init(AppConfig.DBUrl, AppConfig.TablePrefix)

	db.TestConnection()

	// Routing
	r := mux.NewRouter()

	// Endpoint pubblico: login
	r.HandleFunc("/login", api.LoginHandler).Methods("POST")

	// Endpoint pubblico: hello
	r.HandleFunc("/hello", api.HelloHandler).Methods("GET")

	// Endpoint protetti: CRUD utenti
	userRoutes := r.PathPrefix("/users").Subrouter()
	userRoutes.Use(api.AuthMiddleware) // applica il middleware

	userRoutes.HandleFunc("/{id}", api.GetUserHandler).Methods("GET")
	userRoutes.HandleFunc("", api.GetAllUsersHandler).Methods("GET")
	userRoutes.HandleFunc("", api.CreateUserHandler).Methods("POST")
	userRoutes.HandleFunc("/{id}", api.UpdateUserHandler).Methods("PUT")
	userRoutes.HandleFunc("/{id}", api.DeleteUserHandler).Methods("DELETE")

	log.Println("Server in ascolto su :", AppConfig.ServerPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", AppConfig.ServerPort), r))
}

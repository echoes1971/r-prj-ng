package main

/*

Test:

curl -X POST http://localhost:1971/login \
  -H "Content-Type: application/json" \
  -d '{"login":"adm","pwd":"mysecretpass"}'

curl -X GET http://localhost:1971/users/316 \
  -H "Authorization: Bearer <access_token>"

curl -X GET http://localhost:1971/users/ \
  -H "Authorization: Bearer <access_token>"


Some curls for testing:

curl -X GET http://localhost:8080/api/content/a996-3e3aed1c-a911
curl -X GET http://localhost:8080/api/nav/breadcrumb/a996-3e3aed1c-a911
curl -X GET http://localhost:8080/api/nav/children/2c53-b677a6c6-74a1

curl -X GET http://localhost:8080/api/nav/2c53-b677a6c6-74a1/indexes


curl -X GET http://localhost:8080/api/nav/children/-10
curl -X GET http://localhost:8080/api/content/-10
curl -X GET http://localhost:8080/api/nav/children/-12
curl -X GET http://localhost:8080/api/content/-22
curl -X GET http://localhost:8080/api/nav/breadcrumb/-22

curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"login":"adm","pwd":"mysecretpass"}'

curl -X GET http://localhost:8080/api/nav/breadcrumb/d626-5380f5d0-019d \
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
	"strings"

	"rprj/be/api"
	"rprj/be/dblayer"
	"rprj/be/models"

	"github.com/gorilla/mux"
)

var AppConfig models.Config

func main() {

	configFile := "config.json"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	err := models.LoadConfig(configFile, &AppConfig)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	fmt.Printf("Loaded config: %+v\n", AppConfig)

	if appName := os.Getenv("APP_NAME"); appName != "" {
		AppConfig.AppName = appName
	}
	if tablePrefix := os.Getenv("TABLE_PREFIX"); tablePrefix != "" {
		AppConfig.TablePrefix = tablePrefix
	}

	// Override Ollama settings from environment variables if present
	if ollamaURL := os.Getenv("OLLAMA_URL"); ollamaURL != "" {
		AppConfig.OllamaURL = ollamaURL
	}
	if ollamaModel := os.Getenv("OLLAMA_MODEL"); ollamaModel != "" {
		AppConfig.OllamaModel = strings.ReplaceAll(ollamaModel, "\"", "")
	}

	// File system directories
	AppConfig.RootDirectory = "."
	if rootDir := os.Getenv("ROOT_DIRECTORY"); rootDir != "" {
		AppConfig.RootDirectory = rootDir
	}
	AppConfig.FilesDirectory = "files"
	if filesDir := os.Getenv("FILES_DIRECTORY"); filesDir != "" {
		AppConfig.FilesDirectory = filesDir
	}

	dblayer.InitDBLayer(AppConfig)
	// dblayer.EnsureDBSchema()
	dblayer.InitDBData()

	api.InitAPI(AppConfig)
	api.OllamaInit(AppConfig.AppName, AppConfig.OllamaURL, AppConfig.OllamaModel)

	// Routing
	r := mux.NewRouter()
	// remove cors
	r.Use(mux.CORSMethodMiddleware(r))

	// Endpoints navigation
	r.HandleFunc("/content/{objectId}", api.GetNavigationHandler).Methods("GET")
	r.HandleFunc("/content/country/{countryId}", api.GetCountryHandler).Methods("GET")
	r.HandleFunc("/nav/children/{folderId}", api.GetChildrenHandler).Methods("GET")
	r.HandleFunc("/nav/breadcrumb/{objectId}", api.GetBreadcrumbHandler).Methods("GET")
	r.HandleFunc("/nav/{objectId}/indexes", api.GetIndexesHandler).Methods("GET")

	// Endpoint pubblico: login
	r.HandleFunc("/login", api.LoginHandler).Methods("POST")

	// Endpoint pubblico: hello
	r.HandleFunc("/ping", api.PingHandler).Methods("GET")

	// Endpoint pubblico: Ollama
	// curl -X POST http://localhost:8080/api/ollama -H "Content-Type: application/json" -d '{"prompt":"Hello Ollama!"}'
	r.HandleFunc("/ollama", api.OllamaHandler).Methods("POST")
	// Endpoint pubblico: Ollama default page response
	// curl -X GET http://localhost:8080/api/ollama/defaultpage
	r.HandleFunc("/ollama/defaultpage", api.DefaultPageOllamaHandler).Methods("GET")

	// Endpoint pubblico: Get all countries
	r.HandleFunc("/countries", api.GetCountriesHandler).Methods("GET")

	// Endpoint protected: CRUD utenti
	userRoutes := r.PathPrefix("/users").Subrouter()
	userRoutes.Use(api.AuthMiddleware) // applica il middleware

	userRoutes.HandleFunc("/{id}", api.GetUserHandler).Methods("GET")
	userRoutes.HandleFunc("/{id}/person", api.GetUserPersonHandler).Methods("GET")
	userRoutes.HandleFunc("", api.GetAllUsersHandler).Methods("GET")
	userRoutes.HandleFunc("", api.CreateUserHandler).Methods("POST")
	userRoutes.HandleFunc("/{id}", api.UpdateUserHandler).Methods("PUT")
	userRoutes.HandleFunc("/{id}", api.DeleteUserHandler).Methods("DELETE")

	// Endpoint protected: CRUD gruppi
	groupRoutes := r.PathPrefix("/groups").Subrouter()
	groupRoutes.Use(api.AuthMiddleware) // applica il middleware

	groupRoutes.HandleFunc("/{id}", api.GetGroupHandler).Methods("GET")
	groupRoutes.HandleFunc("", api.GetAllGroupsHandler).Methods("GET")
	groupRoutes.HandleFunc("", api.CreateGroupHandler).Methods("POST")
	groupRoutes.HandleFunc("/{id}", api.UpdateGroupHandler).Methods("PUT")
	groupRoutes.HandleFunc("/{id}", api.DeleteGroupHandler).Methods("DELETE")

	// Endpoint protected: CRUD objects (generic DBObject operations)
	objectRoutes := r.PathPrefix("/objects").Subrouter()
	objectRoutes.Use(api.AuthMiddleware) // applica il middleware

	objectRoutes.HandleFunc("/search", api.SearchObjectsHandler).Methods("GET")
	objectRoutes.HandleFunc("/creatable-types", api.GetCreatableTypesHandler).Methods("GET")
	objectRoutes.HandleFunc("", api.CreateObjectHandler).Methods("POST")
	objectRoutes.HandleFunc("/{id}", api.UpdateObjectHandler).Methods("PUT")
	objectRoutes.HandleFunc("/{id}", api.DeleteObjectHandler).Methods("DELETE")

	// Endpoint protected: File download
	fileRoutes := r.PathPrefix("/files").Subrouter()
	fileRoutes.Use(api.AuthMiddleware)
	fileRoutes.HandleFunc("/preview-tokens", api.GenerateFileTokensHandler).Methods("POST")

	// File download without auth middleware (uses token or permission check)
	r.HandleFunc("/files/{id}/download", api.DownloadFileHandler).Methods("GET")

	log.Println("Server in ascolto su :", AppConfig.ServerPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", AppConfig.ServerPort), r))
}

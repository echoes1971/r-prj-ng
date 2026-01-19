package main

// @title ρBee (rhobee) API
// @version 1.0
// @description This is the API documentation for the ρBee (rhobee) backend server. It provides endpoints for user authentication, navigation, file management, and integration with Ollama AI services.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url https://opensource.org/licenses/Apache-2.0

// @host localhost:1971
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @tag.name auth
// @tag.description Authentication endpoints for login and token management

// @tag.name health
// @tag.description Health check and server status endpoints

// @tag.name objects
// @tag.description CRUD operations for DBObjects (notes, pages, companies, etc.)

// @tag.name files
// @tag.description File upload, download and preview token generation

// @tag.name users
// @tag.description User management endpoints

// @tag.name groups
// @tag.description Group management endpoints

// @tag.name navigation
// @tag.description Navigation and content retrieval endpoints

// @tag.name ollama
// @tag.description Endpoints for Ollama AI integration

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

	_ "rprj/be/docs" // Import generated docs

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
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

	// Override DB settings from environment variables if present
	if dbEngine := os.Getenv("DB_ENGINE"); dbEngine != "" {
		AppConfig.DBEngine = dbEngine
	}
	if dbURL := os.Getenv("DB_URL"); dbURL != "" {
		AppConfig.DBUrl = dbURL
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

	// OAuth settings from environment variables
	if googleClientId := os.Getenv("GOOGLE_CLIENT_ID"); googleClientId != "" {
		AppConfig.GoogleClientID = googleClientId
	}
	if googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET"); googleClientSecret != "" {
		AppConfig.GoogleClientSecret = googleClientSecret
	}
	if googleRedirectURL := os.Getenv("GOOGLE_REDIRECT_URL"); googleRedirectURL != "" {
		AppConfig.GoogleRedirectURL = googleRedirectURL
	}
	if githubClientId := os.Getenv("GITHUB_CLIENT_ID"); githubClientId != "" {
		AppConfig.GitHubClientID = githubClientId
	}
	if githubClientSecret := os.Getenv("GITHUB_CLIENT_SECRET"); githubClientSecret != "" {
		AppConfig.GitHubClientSecret = githubClientSecret
	}
	if githubRedirectURL := os.Getenv("GITHUB_REDIRECT_URL"); githubRedirectURL != "" {
		AppConfig.GitHubRedirectURL = githubRedirectURL
	}
	if telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN"); telegramBotToken != "" {
		AppConfig.TelegramBotToken = telegramBotToken
	}
	// Extract bot_id from token (format: "123456789:ABCdef...")
	if AppConfig.TelegramBotToken != "" && AppConfig.TelegramBotID == "" {
		parts := strings.Split(AppConfig.TelegramBotToken, ":")
		if len(parts) > 0 {
			AppConfig.TelegramBotID = parts[0]
		}
	}

	dblayer.InitDBLayer(AppConfig)
	dblayer.EnsureDBSchema()
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
	r.HandleFunc("/nav/search", api.NavigationSearchHandler).Methods("GET")

	// Public Endpoints: login, logout
	r.HandleFunc("/login", api.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", api.LogoutHandler).Methods("POST")

	// OAuth endpoints (Google, GitHub, Telegram)
	r.HandleFunc("/oauth/google/start", api.GoogleOAuthStart).Methods("GET")
	r.HandleFunc("/oauth/google/callback", api.GoogleOAuthCallback).Methods("GET")
	r.HandleFunc("/oauth/github/start", api.GitHubOAuthStart).Methods("GET")
	r.HandleFunc("/oauth/github/callback", api.GitHubOAuthCallback).Methods("GET")
	r.HandleFunc("/oauth/telegram/callback", api.TelegramOAuthCallback).Methods("GET")

	// Public Endpoint: hello
	r.HandleFunc("/ping", api.PingHandler).Methods("GET")

	// Public Endpoint: Ollama
	// curl -X POST http://localhost:8080/api/ollama -H "Content-Type: application/json" -d '{"prompt":"Hello Ollama!"}'
	r.HandleFunc("/ollama", api.OllamaHandler).Methods("POST")
	// Public Endpoint: Ollama default page response
	// curl -X GET http://localhost:8080/api/ollama/defaultpage
	r.HandleFunc("/ollama/defaultpage", api.DefaultPageOllamaHandler).Methods("GET")

	// Public Endpoint: Get all countries
	r.HandleFunc("/countries", api.GetCountriesHandler).Methods("GET")

	// Protected Endpoint: CRUD users
	userRoutes := r.PathPrefix("/users").Subrouter()
	userRoutes.Use(api.AuthMiddleware) // applica il middleware

	userRoutes.HandleFunc("/{id}", api.GetUserHandler).Methods("GET")
	userRoutes.HandleFunc("/{id}/person", api.GetUserPersonHandler).Methods("GET")
	userRoutes.HandleFunc("", api.GetAllUsersHandler).Methods("GET")
	userRoutes.HandleFunc("", api.CreateUserHandler).Methods("POST")
	userRoutes.HandleFunc("/{id}", api.UpdateUserHandler).Methods("PUT")
	userRoutes.HandleFunc("/{id}", api.DeleteUserHandler).Methods("DELETE")

	// Protected Endpoint: CRUD groups
	groupRoutes := r.PathPrefix("/groups").Subrouter()
	groupRoutes.Use(api.AuthMiddleware) // applica il middleware

	groupRoutes.HandleFunc("/{id}", api.GetGroupHandler).Methods("GET")
	groupRoutes.HandleFunc("", api.GetAllGroupsHandler).Methods("GET")
	groupRoutes.HandleFunc("", api.CreateGroupHandler).Methods("POST")
	groupRoutes.HandleFunc("/{id}", api.UpdateGroupHandler).Methods("PUT")
	groupRoutes.HandleFunc("/{id}", api.DeleteGroupHandler).Methods("DELETE")

	// Protected Endpoint: CRUD objects (generic DBObject operations)
	objectRoutes := r.PathPrefix("/objects").Subrouter()
	objectRoutes.Use(api.AuthMiddleware) // applica il middleware

	// objectRoutes.HandleFunc("/search", api.SearchObjectsHandler).Methods("GET")
	objectRoutes.HandleFunc("/creatable-types", api.GetCreatableTypesHandler).Methods("GET")
	objectRoutes.HandleFunc("", api.CreateObjectHandler).Methods("POST")
	objectRoutes.HandleFunc("/{id}", api.UpdateObjectHandler).Methods("PUT")
	objectRoutes.HandleFunc("/{id}", api.DeleteObjectHandler).Methods("DELETE")

	// Protected Endpoint: File download
	fileRoutes := r.PathPrefix("/files").Subrouter()
	fileRoutes.Use(api.AuthMiddleware)
	fileRoutes.HandleFunc("/preview-tokens", api.GenerateFileTokensHandler).Methods("POST")

	// File download without auth middleware (uses token or permission check)
	r.HandleFunc("/files/{id}/download", api.DownloadFileHandler).Methods("GET")
	r.HandleFunc("/objects/search", api.SearchObjectsHandler).Methods("GET")

	// Protected Endpoint: Admin
	adminRoutes := r.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(api.AuthMiddleware)
	adminRoutes.HandleFunc("/dashboard", api.DashboardHandler).Methods("GET")

	// Swagger documentation - only in development
	enableSwagger := os.Getenv("ENABLE_SWAGGER")
	if enableSwagger == "true" || enableSwagger == "1" {
		r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
		log.Println("Swagger UI disponibile su: http://localhost:" + fmt.Sprintf("%d", AppConfig.ServerPort) + "/swagger/index.html")
	} else {
		log.Println("Swagger UI disabilitato (set ENABLE_SWAGGER=true per abilitare)")
	}

	log.Println("Server in ascolto su :", AppConfig.ServerPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", AppConfig.ServerPort), r))
}

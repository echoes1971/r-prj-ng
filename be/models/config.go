package models

import (
	"encoding/json"
	"log"
	"os"
)

// Backend Configuration Structure
type Config struct {
	AppName        string `json:"app_name"`
	ServerPort     int    `json:"server_port"`
	DBEngine       string `json:"db_engine"`
	DBUrl          string `json:"db_url"`
	TablePrefix    string `json:"table_prefix"`
	JWTSecret      string `json:"jwt_secret"`
	LogLevel       string `json:"log_level"`
	OllamaModel    string `json:"ollama_model"`
	OllamaURL      string `json:"ollama_url"`
	RootDirectory  string `json:"root_directory"`
	FilesDirectory string `json:"files_directory"`
	// OAuth configuration
	GoogleClientID     string `json:"google_client_id"`
	GoogleClientSecret string `json:"google_client_secret"`
	GoogleRedirectURL  string `json:"google_redirect_url"`
	GitHubClientID     string `json:"github_client_id"`
	GitHubClientSecret string `json:"github_client_secret"`
	GitHubRedirectURL  string `json:"github_redirect_url"`
	TelegramBotToken   string `json:"telegram_bot_token"`
	TelegramBotID      string `json:"telegram_bot_id"`
}

func LoadConfig(filename string, config *Config) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Decodifica JSON in struct
	if err := json.NewDecoder(file).Decode(config); err != nil {
		log.Fatal(err)
	}
	// Implementazione per caricare la configurazione da un file JSON
	// (omessa per brevità)
	return err
}

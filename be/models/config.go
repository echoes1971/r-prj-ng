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

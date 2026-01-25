package dblayer

import (
	_ "embed"
	"encoding/json"
	"log"
)

// TableData represents initial data for a single table
type TableData struct {
	Name    string     `json:"name"`
	Columns []string   `json:"columns"`
	Data    [][]string `json:"data"`
}

// InitialData represents the complete initial data structure
type InitialData struct {
	Comment string      `json:"comment"`
	Tables  []TableData `json:"tables"`
}

//go:embed initialdata.json
var initialDataJSON []byte

// LoadInitialData parses the embedded JSON and returns the initial data structure
func LoadInitialData() (*InitialData, error) {
	var data InitialData
	err := json.Unmarshal(initialDataJSON, &data)
	if err != nil {
		log.Printf("Failed to parse initial data JSON: %v", err)
		return nil, err
	}
	return &data, nil
}

package dblayer

import (
	"log"
	"testing"
)

func TestCreateTableStrings(t *testing.T) {
	factory := NewDBEFactory(true)
	user := NewDBUser()
	factory.Register(user)
	group := NewDBGroup()
	factory.Register(group)
	userGroup := NewUserGroup()
	factory.Register(userGroup)

	// Loop through all registered entities and print their CREATE TABLE strings
	for _, className := range factory.GetAllClassNames() {
		entity := factory.GetInstanceByClassName(className)
		createTableSQL := GetCreateTableSQL(entity, DbSchema)
		log.Printf("CREATE TABLE SQL for %s:\n%s\n", className, createTableSQL)
	}
}

// go test -v ./dblayer -run TestEnsureDBSchema -config ../config_test.json
func TestEnsureDBSchema(t *testing.T) {
	// Call EnsureDBSchema to test table creation logic
	EnsureDBSchema(true)
	InitDBData()

	// IF db engine is sqlite3, use sqlite3 mechanism to know where is the db file via DbConnection.Exec
	if dbEngine == "sqlite3" {
		var dbPath string
		err := DbConnection.QueryRow("PRAGMA database_list;").Scan(new(interface{}), new(interface{}), &dbPath)
		if err != nil {
			t.Fatalf("Error getting SQLite database path: %v", err)
		}
		log.Printf("SQLite database file path: %s", dbPath)
	}
}

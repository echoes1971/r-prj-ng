package dblayer

import (
	"database/sql"
	"log"
	"strings"
)

// The database can be mysql, sqlite, postgres, etc.
var dbEngine string
var dbUrl string
var DbSchema string
var DbConnection *sql.DB
var Factory *DBEFactory

func InitDBLayer(dbEngineName, dbUrlAddress, schema string) {
	dbEngine = dbEngineName
	dbUrl = dbUrlAddress
	DbSchema = strings.ReplaceAll(schema, "_", "")
	log.Print("DB Schema:", DbSchema)

	log.Print("Initializing DBEFactory...")

	Factory = NewDBEFactory(false)
	Factory.Register(NewDBVersion())
	Factory.Register(NewOAuthToken())
	Factory.Register(NewDBUser())
	Factory.Register(NewUserGroup())
	Factory.Register(NewDBGroup())
	log.Print("Initializing DB connection...")
	var err error
	DbConnection, err = sql.Open("mysql", dbUrl)
	if err != nil {
		log.Fatal("Error opening DB connection:", err)
	}
	err = DbConnection.Ping()
	if err != nil {
		log.Fatal("Error pinging DB:", err)
	}

	// TODO: make this configurable
	// Configure connection pool to handle concurrent operations
	DbConnection.SetMaxOpenConns(25)   // Maximum number of open connections to the database
	DbConnection.SetMaxIdleConns(10)   // Maximum number of connections in the idle connection pool
	DbConnection.SetConnMaxLifetime(0) // Maximum amount of time a connection may be reused (0 = unlimited)

}

func InitDBConnection() {
	if DbConnection != nil {
		return
	}
	log.Print("Initializing DB connection...")
	var err error
	DbConnection, err = sql.Open(dbEngine, dbUrl)
	if err != nil {
		log.Fatal("Error opening DB connection:", err)
	}
	err = DbConnection.Ping()
	if err != nil {
		log.Fatal("Error pinging DB:", err)
	}
}

// The database can be mysql, sqlite, postgres, etc.
func ensureTableExistsAndUpdatedForMysql(dbe DBEntityInterface) error {
	// Check if table exists
	tableName := DbSchema + "_" + dbe.GetTableName()
	// This is internal code, so we can build the query directly
	query := "show tables like '" + tableName + "'"
	var existingTable string
	err := DbConnection.QueryRow(query).Scan(&existingTable)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error checking existence of table %s: %v", tableName, err)
		return err
	}

	if existingTable == "" {
		// Table does not exist, create it
		// Compose the create table SQL using the DBEntity's information about columns, types, keys, etc.
		createTableSQL := dbe.GetCreateTableSQL(DbSchema)
		log.Printf("Creating table with SQL: %s", createTableSQL)
		_, err := DbConnection.Exec(createTableSQL)
		if err != nil {
			return err
		}
		log.Printf("Created table %s", tableName)
	} else {
		// Table exists, check for schema updates
		// For simplicity, we will not implement schema migration logic here
		log.Printf("Table %s already exists", tableName)
	}

	return nil
}
func ensureTableExistsAndUpdatedForSqlite(dbe DBEntityInterface) error {
	// TODO: Implement for Sqlite
	return nil
}

// Iterate over all registered DBEntity types and create tables if they do not exist or update their schema
func EnsureDBSchema() {
	var classInstances []DBEntityInterface
	for _, className := range Factory.GetAllClassNames() {
		dbe := Factory.GetInstanceByClassName(className)
		if dbe != nil {
			classInstances = append(classInstances, dbe)
		}
	}
	// Sort classInstances based on dependencies (foreign keys) if needed
	// Sort based on dependencies using a topological sort
	sorted := make([]DBEntityInterface, 0, len(classInstances))
	visited := make(map[string]bool)
	temp := make(map[string]bool)

	var visit func(DBEntityInterface) bool
	visit = func(dbe DBEntityInterface) bool {
		className := dbe.GetTypeName()
		if temp[className] {
			// Circular dependency detected
			log.Printf("Warning: Circular dependency detected for %s", className)
			return false
		}
		if visited[className] {
			return true
		}

		temp[className] = true
		// Get foreign key dependencies
		foreignKeys := dbe.GetForeignKeys()
		for _, fk := range foreignKeys {
			// Find the referenced table's DBEntity
			for _, depDbe := range classInstances {
				if depDbe.GetTableName() == fk.RefTable {
					if !visit(depDbe) {
						return false
					}
					break
				}
			}
		}
		temp[className] = false
		visited[className] = true
		sorted = append(sorted, dbe)
		return true
	}

	for _, dbe := range classInstances {
		if !visited[dbe.GetTypeName()] {
			visit(dbe)
		}
	}
	// reverse sorted to get correct order
	// for i, j := 0, len(sorted)-1; i < j; i, j = i+1, j-1 {
	// 	sorted[i], sorted[j] = sorted[j], sorted[i]
	// }
	// slices.Reverse(sorted)

	// Print sorted class names for debugging
	log.Print("DB Entities creation order:")
	for _, dbe := range sorted {
		log.Printf(" - %s\n", dbe.GetTableName())
		for _, fk := range dbe.GetForeignKeys() {
			log.Printf("    FK: %s -> %s(%s)\n", fk.Column, fk.RefTable, fk.RefColumn)
		}
	}

	classInstances = sorted

	for _, dbe := range classInstances {
		var err error
		className := dbe.GetTypeName()
		switch dbEngine {
		case "mysql":
			err = ensureTableExistsAndUpdatedForMysql(dbe)
		case "sqlite":
			err = ensureTableExistsAndUpdatedForSqlite(dbe)
		default:
			log.Fatal("Unsupported dbEngine:", dbEngine)
		}
		if err != nil {
			log.Fatal("Error ensuring table for ", className, ":", err)
		}
	}
}
func CloseDBConnection() {
	if DbConnection != nil {
		log.Print("Closing DB connection...")
		DbConnection.Close()
	}
}

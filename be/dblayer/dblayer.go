package dblayer

import (
	"database/sql"
	"log"
	"strings"
)

var dbUrl string
var DbSchema string
var DbConnection *sql.DB
var Factory *DBEFactory

func InitDBLayer(dbUrlAddress, schema string) {
	dbUrl = dbUrlAddress
	DbSchema = strings.ReplaceAll(schema, "_", "")
	log.Print("DB Schema:", DbSchema)

	log.Print("Initializing DBEFactory...")

	Factory = NewDBEFactory(false)
	Factory.Register(NewDBVersion())
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
	DbConnection, err = sql.Open("mysql", dbUrl)
	if err != nil {
		log.Fatal("Error opening DB connection:", err)
	}
	err = DbConnection.Ping()
	if err != nil {
		log.Fatal("Error pinging DB:", err)
	}
}

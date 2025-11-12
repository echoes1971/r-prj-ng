package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var dbURL string
var tablePrefix string

var DB *sql.DB

// TestConnection apre la connessione, esegue la query e stampa il risultato
func TestConnection(dbURL string) {
	// URL di connessione
	dsn := dbURL

	// Apri connessione
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Errore apertura connessione:", err)
	}
	defer db.Close()

	// Verifica che la connessione sia attiva
	if err := db.Ping(); err != nil {
		log.Fatal("Errore ping DB:", err)
	}

	fmt.Println("Connessione a MariaDB riuscita!")

	// Esegui la query
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM " + tablePrefix + "users;").Scan(&count)
	if err != nil {
		log.Fatal("Errore esecuzione query:", err)
	}

	fmt.Printf("Numero di utenti in rra_users: %d\n", count)
}

func Init(url string, prefix string) {
	dbURL = url
	tablePrefix = prefix
	fmt.Println("DB inizializzato con URL:", dbURL)
	fmt.Println("Prefisso tabelle:", tablePrefix)

	var err error
	DB, err = sql.Open("mysql", dbURL)
	if err != nil {
		log.Fatal("Errore apertura connessione:", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatal("Errore ping DB:", err)
	}
	log.Println("Connessione a MariaDB riuscita!")
}

/* Generate a random UUID-like string of 16 hex characters */
func uuid16HexGo() (string, error) {
	b := make([]byte, 8) // 8 bytes = 16 hex chars
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

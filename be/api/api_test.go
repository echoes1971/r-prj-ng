package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"rprj/be/dblayer"
	"rprj/be/models"
)

var AppConfig models.Config

// TestMain viene eseguito prima di tutti i test del package
func TestMain(m *testing.M) {

	err := models.LoadConfig("../config.json", &AppConfig)
	if err != nil {
		log.Fatalf("Errore caricamento configurazione: %v", err)
	}

	dblayer.InitDBLayer(AppConfig.DBEngine, AppConfig.DBUrl, AppConfig.TablePrefix)

	log.Println("DB inizializzato per i test")

	// Esegui i test
	code := m.Run()

	// Teardown: chiudi la connessione
	dblayer.CloseDBConnection()

	// Exit con il codice dei test
	os.Exit(code)
}

func TestLoginHandler(t *testing.T) {
	// Prepara il body JSON con credenziali
	creds := Credentials{
		Login: "adm",
		Pwd:   "mysecretpass",
	}
	body, _ := json.Marshal(creds)

	// Crea una request POST
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Recorder per catturare la risposta
	rr := httptest.NewRecorder()

	// Chiama direttamente l'handler
	handler := http.HandlerFunc(LoginHandler)
	handler.ServeHTTP(rr, req)

	// Controlla lo status code
	if rr.Code != http.StatusOK {
		t.Errorf("status code errato: got %v, want %v", rr.Code, http.StatusOK)
	}

	// Controlla che la risposta contenga un access_token
	var resp TokenResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Errorf("errore parsing risposta JSON: %v", err)
	}
	log.Printf("resp: %v\n", resp)

	if resp.AccessToken == "" {
		t.Errorf("access_token mancante nella risposta")
	}
}

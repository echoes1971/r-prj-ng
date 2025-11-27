package api

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// go test -v ./api -run TestObjectHandlerSearchObject
func TestObjectHandlerSearchObject(t *testing.T) {
	token := ApiTestDoLogin(t, testUser.GetValue("login").(string), testUser.GetValue("pwd").(string))

	req := httptest.NewRequest(http.MethodGet, "/object/search?classname=DBFolder&name=Home&limit=10", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(SearchObjectsHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", rr.Code)
	}

	var response map[string]any
	log.Print("Response body:", rr.Body.String())
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	log.Print("Response:", response)
	if err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	if response["success"] != true {
		t.Fatalf("Expected success status, got %v", response["success"])
	}

	objects, ok := response["objects"].([]any)
	if !ok {
		t.Fatalf("Expected objects to be a list, got %T", response["objects"])
	}

	if len(objects) == 0 {
		t.Fatalf("Expected at least one search result, got 0")
	}

	firstResult, ok := objects[0].(map[string]any)
	if !ok {
		t.Fatalf("Expected first result to be a map, got %T", objects[0])
	}
	log.Print("First result:", firstResult)

	if firstResult["name"] != "Home" {
		t.Fatalf("Expected first result name to be 'Home', got '%v'", firstResult["name"])
	}

	log.Printf("TestObjectHandlerSearchObject passed, found object: %v", firstResult)
}

func TestObjectHandlerSearchObjectUser(t *testing.T) {
	token := ApiTestDoLogin(t, testUser.GetValue("login").(string), testUser.GetValue("pwd").(string))
	log.Print("Obtained token:", token)

	req := httptest.NewRequest(http.MethodGet, "/object/search?classname=DBUser&name="+testUser.GetValue("login").(string)+"&limit=10", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(SearchObjectsHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", rr.Code)
	}

	var response map[string]any
	log.Print("Response body:", rr.Body.String())
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	log.Print("Response:", response)
	if err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	if response["success"] != true {
		t.Fatalf("Expected success status, got %v", response["success"])
	}

	objects, ok := response["objects"].([]any)
	if !ok {
		t.Fatalf("Expected objects to be a list, got %T", response["objects"])
	}

	if len(objects) == 0 {
		t.Fatalf("Expected at least one search result, got 0")
	}

	firstResult, ok := objects[0].(map[string]any)
	if !ok {
		t.Fatalf("Expected first result to be a map, got %T", objects[0])
	}
	log.Print("First result:", firstResult)

	if firstResult["name"] != testUser.GetValue("login").(string) {
		t.Fatalf("Expected first result name to be '%v', got '%v'", testUser.GetValue("login").(string), firstResult["login"])
	}

	log.Printf("TestObjectHandlerSearchObject passed, found object: %v", firstResult)
}

func TestObjectHandlerSearchObjectGroup(t *testing.T) {
	token := ApiTestDoLogin(t, testUser.GetValue("login").(string), testUser.GetValue("pwd").(string))
	log.Print("Obtained token:", token)

	req := httptest.NewRequest(http.MethodGet, "/object/search?classname=DBGroup&name="+testUser.GetValue("login").(string)+"&limit=10", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(SearchObjectsHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", rr.Code)
	}

	var response map[string]any
	log.Print("Response body:", rr.Body.String())
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	log.Print("Response:", response)
	if err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	if response["success"] != true {
		t.Fatalf("Expected success status, got %v", response["success"])
	}

	objects, ok := response["objects"].([]any)
	if !ok {
		t.Fatalf("Expected objects to be a list, got %T", response["objects"])
	}

	if len(objects) == 0 {
		t.Fatalf("Expected at least one search result, got 0")
	}

	firstResult, ok := objects[0].(map[string]any)
	if !ok {
		t.Fatalf("Expected first result to be a map, got %T", objects[0])
	}
	log.Print("First result:", firstResult)

	if !strings.Contains(firstResult["name"].(string), testUser.GetValue("login").(string)) {
		t.Fatalf("Expected first result name to contain '%v', got '%v'", testUser.GetValue("login").(string), firstResult["name"])
	}

	log.Printf("TestObjectHandlerSearchObjectGroup passed, found object: %v", firstResult)
}

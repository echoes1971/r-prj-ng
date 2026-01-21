package api

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// go test -v ./api -run TestAdminDashboardHandler
func TestAdminDashboardHandler(t *testing.T) {
	// repo := SetupTestRepo(t, "-1", []string{"-2"}, DbSchema)

	// First, log in to get a token
	token := ApiTestDoLogin(t, testAdminLogin, testAdminPwd)

	// Now test the /admin/dashboard endpoint
	req := httptest.NewRequest("GET", "/admin/dashboard", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	DashboardHandler(w, req)

	if w.Code != http.StatusOK {
		log.Println("Response body:", w.Body.String())
		t.Fatalf("DashboardHandler returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}

	var dashboardResp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &dashboardResp)
	if err != nil {
		t.Fatalf("Failed to parse dashboard response: %v", err)
	}
	log.Printf("Dashboard response: %+v\n", dashboardResp)

	usersCount, ok := dashboardResp["users_count"].(float64)
	log.Printf("Dashboard users_count: %v\n", usersCount)
	if !ok {
		t.Fatalf("Unexpected dashboard response users_count: %v", dashboardResp["users_count"])
	}
}

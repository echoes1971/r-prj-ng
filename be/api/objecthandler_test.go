package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
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

// NOTE: this should be handled by NavigationSearchHandler
// func TestObjectHandlerSearchObjectUser(t *testing.T) {
// 	token := ApiTestDoLogin(t, testUser.GetValue("login").(string), testUser.GetValue("pwd").(string))
// 	log.Print("Obtained token:", token)
// 	req := httptest.NewRequest(http.MethodGet, "/object/search?classname=DBUser&name="+testUser.GetValue("login").(string)+"&limit=10", nil)
// 	req.Header.Set("Authorization", "Bearer "+token)
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(SearchObjectsHandler)
// 	handler.ServeHTTP(rr, req)
// 	if rr.Code != http.StatusOK {
// 		log.Print(" Response body:", rr.Body.String())
// 		t.Fatalf("Expected status OK, got %v", rr.Code)
// 	}
// 	var response map[string]any
// 	log.Print("Response body:", rr.Body.String())
// 	err := json.Unmarshal(rr.Body.Bytes(), &response)
// 	log.Print("Response:", response)
// 	if err != nil {
// 		t.Fatalf("Failed to parse response JSON: %v", err)
// 	}
// 	if response["success"] != true {
// 		t.Fatalf("Expected success status, got %v", response["success"])
// 	}
// 	objects, ok := response["objects"].([]any)
// 	if !ok {
// 		t.Fatalf("Expected objects to be a list, got %T", response["objects"])
// 	}
// 	if len(objects) == 0 {
// 		t.Fatalf("Expected at least one search result, got 0")
// 	}
// 	firstResult, ok := objects[0].(map[string]any)
// 	if !ok {
// 		t.Fatalf("Expected first result to be a map, got %T", objects[0])
// 	}
// 	log.Print("First result:", firstResult)
// 	if firstResult["name"] != testUser.GetValue("login").(string) {
// 		t.Fatalf("Expected first result name to be '%v', got '%v'", testUser.GetValue("login").(string), firstResult["login"])
// 	}
// 	log.Printf("TestObjectHandlerSearchObject passed, found object: %v", firstResult)
// }

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
	log.Print(" Response body:", rr.Body.String())
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

/*
The UI:
1. Create a DBFile without uploading a file
2. Upload a file to the created DBFile object

This test simulates the same steps via API:
1. Create a DBFile object via CreateObjectHandler
2. Upload a file to the created DBFile object via UploadFileHandler
*/
func TestObjectHandlerUploadFile(t *testing.T) {
	token := ApiTestDoLogin(t, testUser.GetValue("login").(string), testUser.GetValue("pwd").(string))

	containerID := "c64dfecb8296b08f" // testFolderHomeID

	//{"classname":"DBFile","father_id":"515","name":"New File","description":""}
	reqBody := `{"classname":"DBFile","permissions":"rwxrw-rw-","father_id":"` + containerID + `","name":"Test Upload File.txt","description":"Test file upload via API"}`
	req := httptest.NewRequest(http.MethodPost, "/object", strings.NewReader(reqBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(CreateObjectHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Expected status Created from CreateObjectHandler, got %v", rr.Code)
	}

	var createResponse map[string]any
	log.Print("CreateObjectHandler response body:", rr.Body.String())
	err := json.Unmarshal(rr.Body.Bytes(), &createResponse)
	log.Print("CreateObjectHandler response:", createResponse)
	if err != nil {
		t.Fatalf("Failed to parse CreateObjectHandler response JSON: %v", err)
	}

	if createResponse["success"] != true {
		t.Fatalf("Expected success status from CreateObjectHandler, got %v", createResponse["success"])
	}

	createdObject, ok := createResponse["data"].(map[string]any)
	if !ok {
		t.Fatalf("Expected created object to be a map, got %T", createResponse["data"])
	}
	createdObjectJSON, _ := json.MarshalIndent(createdObject, "", "  ")
	log.Print("Created object:", string(createdObjectJSON))

	createdObjectID, ok := createdObject["id"].(string)
	if !ok || createdObjectID == "" {
		t.Fatalf("Expected created object to have a valid ID, got '%v'", createdObject["id"])
	}

	// Now upload a file to the created DBFile object calling UpdateObjectHandler,
	// passing the file in a multipart/form-data request
	reqBody = `--boundary
Content-Disposition: form-data; name="file"; filename="Test Upload File.txt"
Content-Type: text/plain

This is a test file uploaded via API.
--boundary--
Content-Disposition: form-data; name="file"; filename="Screenshot_20251005_173213.png"
Content-Type: image/png


--boundary--
Content-Disposition: form-data; name="name"

New File
--boundary--
Content-Disposition: form-data; name="description"


--boundary--
Content-Disposition: form-data; name="alt_link"


--boundary--
Content-Disposition: form-data; name="fk_obj_id"

0
--boundary--
Content-Disposition: form-data; name="permissions"

rwxrw-r--
--boundary--
`
	req.Header.Set("Content-Type", "multipart/form-data; boundary=boundary")

	// format ID like XXXX-XXXXXXXX-XXXX
	formattedID := fmt.Sprintf("%s-%s-%s", createdObjectID[:4], createdObjectID[4:12], createdObjectID[12:16])

	log.Print("Uploading file to object ID:", formattedID)
	req = httptest.NewRequest(http.MethodPut, "/objects/"+formattedID, strings.NewReader(reqBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "multipart/form-data; boundary=boundary")

	// Use mux router to properly set path variables
	rr = httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/objects/{id}", UpdateObjectHandler).Methods("PUT")
	router.ServeHTTP(rr, req)
	log.Print("UpdateObjectHandler response body:", rr.Body.String())

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status OK from UpdateObjectHandler, got %v", rr.Code)
	}

	var uploadResponse map[string]any
	log.Print("UpdateObjectHandler response body:", rr.Body.String())
	err = json.Unmarshal(rr.Body.Bytes(), &uploadResponse)
	log.Print("UpdateObjectHandler response:", uploadResponse)
	if err != nil {
		t.Fatalf("Failed to parse UpdateObjectHandler response JSON: %v", err)
	}

	if uploadResponse["success"] != true {
		t.Fatalf("Expected success status from UpdateObjectHandler, got %v", uploadResponse["success"])
	}

	uploadedObject, ok := uploadResponse["data"].(map[string]any)
	if !ok {
		t.Fatalf("Expected uploaded object to be a map, got %T", uploadResponse["data"])
	}
	log.Print("Uploaded object:", uploadedObject)

	if uploadedObject["id"] != createdObjectID {
		t.Fatalf("Expected uploaded object ID to be '%v', got '%v'", createdObjectID, uploadedObject["id"])
	}

	log.Printf("TestObjectHandlerUploadFile passed, uploaded file to object ID: %v", createdObjectID)

}

// go test -v ./api -run TestGetCreatableTypesHandler
func TestGetCreatableTypesHandler(t *testing.T) {
	token := ApiTestDoLogin(t, testAdminLogin, testAdminPwd)

	// claims := ApiTestDecodeAccessToken(t, token)
	// claimsJSON, _ := json.Marshal(claims)
	// log.Printf("Decoded token claims: %+v", string(claimsJSON))

	// Create a folder
	repo := SetupTestRepo(t,
		testUser.GetValue("id").(string),
		[]string{testUser.GetValue("group_id").(string)},
		AppConfig.TablePrefix)

	// folder, err := repo.CreateObject("folders", map[string]any{"name": "testfolder", "permissions": "rwxrw----"}, map[string]any{})
	folder, err := repo.CreateObject("folders", map[string]any{"name": "testfolder"}, map[string]any{})
	if err != nil {
		t.Fatalf("Failed to create folder: %v", err)
	}
	log.Printf("Created folder: %s %+v", folder.GetValue("id").(string), folder.ToJSON())

	// Test the GetCreatableTypesHandler
	// req := httptest.NewRequest(http.MethodGet, "/object/creatable_types?father_id=515", nil)
	req := httptest.NewRequest(http.MethodGet, "/object/creatable_types?father_id="+folder.GetValue("id").(string), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(GetCreatableTypesHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		log.Print("Response body:", rr.Body.String())
		t.Fatalf("Expected status OK, got %v", rr.Code)
	}

	var response map[string]any
	log.Print("Response body:", rr.Body.String())
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	// log.Print("Response:", response)
	if err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	if response["success"] != true {
		t.Fatalf("Expected success status, got %v", response["success"])
	}

	types, ok := response["types"].([]any)
	if !ok {
		t.Fatalf("Expected types to be a list, got %T", response["types"])
	}

	if len(types) == 0 {
		t.Fatalf("Expected at least one creatable type, got 0")
	}

	// Delete the created folder
	folder, err = repo.Delete(folder)
	if err != nil {
		t.Fatalf("Failed to soft delete folder: %v", err)
	}
	folder, err = repo.Delete(folder)
	if err != nil {
		t.Fatalf("Failed to hard delete folder: %v", err)
	}
	log.Printf("TestGetCreatableTypesHandler passed, found %d creatable types", len(types))
}

package dblayer

import (
	"log"
	"testing"
)

func TestCRUDDBObject(t *testing.T) {
	repo := setupTestRepo(t)

	// Create
	dbObj := repo.factory.GetInstanceByTableNameWithValues("objects", map[string]any{
		"name":        "Test Object",
		"description": "This is a test object",
	})
	if dbObj == nil {
		t.Fatal("Failed to create DBObject instance")
	}

	// Insert
	log.Print("Insert dbObj=", dbObj.ToString())
	_, err := repo.Insert(dbObj)
	if err != nil {
		t.Fatalf("Failed to create DBObject: %v", err)
	}
	if dbObj.GetValue("id") == nil {
		t.Fatal("DBObject ID not set after creation")
	}
	log.Printf("Created DBObject: %v", dbObj.ToString())

	// Read
	readObj := repo.GetInstanceByTableName("objects")
	readObj.SetValue("id", dbObj.GetValue("id"))
	foundObjs, err := repo.Search(readObj, false, false, "")
	if err != nil {
		t.Fatalf("Failed to read DBObject: %v", err)
	}
	if len(foundObjs) != 1 {
		t.Fatalf("Expected to find 1 DBObject, found %d", len(foundObjs))
	}
	foundObj, ok := foundObjs[0].(*DBObject)
	if !ok {
		t.Fatal("Found instance is not of type DBObject")
	}
	if foundObj.GetValue("name") != "Test Object" {
		t.Fatalf("Expected name 'Test Object', got '%v'", foundObj.GetValue("name"))
	}
	log.Printf("Read DBObject: %v", foundObj.ToString())

	// Update
	foundObj.SetValue("description", "Updated description")
	updatedObj, err := repo.Update(foundObj)
	if err != nil {
		t.Fatalf("Failed to update DBObject: %v", err)
	}
	if updatedObj.GetValue("description") != "Updated description" {
		t.Fatalf("Expected description 'Updated description', got '%v'", updatedObj.GetValue("description"))
	}
	log.Printf("Updated DBObject: %v", updatedObj.ToString())

	// Verify Update
	verifyObj := repo.GetInstanceByTableName("objects")
	verifyObj.SetValue("id", dbObj.GetValue("id"))
	foundObjs, err = repo.Search(verifyObj, false, false, "")
	if err != nil {
		t.Fatalf("Failed to read DBObject for verification: %v", err)
	}
	if len(foundObjs) != 1 {
		t.Fatalf("Expected to find 1 DBObject for verification, found %d", len(foundObjs))
	}
	verifiedObj, ok := foundObjs[0].(*DBObject)
	if !ok {
		t.Fatal("Verified instance is not of type DBObject")
	}
	if verifiedObj.GetValue("description") != "Updated description" {
		t.Fatalf("Expected description 'Updated description', got '%v'", verifiedObj.GetValue("description"))
	}
	log.Printf("Verified DBObject: %v", verifiedObj.ToString())

	// Delete
	_, err = repo.Delete(verifiedObj)
	if err != nil {
		t.Fatalf("Failed to delete DBObject: %v", err)
	}

	// Verify Soft Deletion
	deleteObj := repo.GetInstanceByTableName("objects")
	deleteObj.SetValue("id", dbObj.GetValue("id"))
	foundObjs, err = repo.Search(deleteObj, false, false, "")
	if err != nil {
		t.Fatalf("Failed to read DBObject for deletion verification: %v", err)
	}
	if len(foundObjs) != 1 {
		t.Fatalf("Expected to find 1 DBObject after deletion, found %d", len(foundObjs))
	}
	log.Printf("Verified Soft Deletion: %v", verifiedObj.ToString())

	// Verify Hard Deletion
	_, err = repo.Delete(verifiedObj)
	if err != nil {
		t.Fatalf("Failed to hard delete DBObject: %v", err)
	}
	foundObjs, err = repo.Search(deleteObj, false, false, "")
	if err != nil {
		t.Fatalf("Failed to read DBObject for hard deletion verification: %v", err)
	}
	if len(foundObjs) != 0 {
		t.Fatalf("Expected to find 0 DBObjects after hard deletion, found %d", len(foundObjs))
	}
	log.Printf("Verified Hard Deletion: %v", verifiedObj.ToString())

	log.Print("DBObject CRUD operations test completed successfully")
}

func TestObjectById(t *testing.T) {
	// Setup
	dbContext := &DBContext{
		UserID:   "-1",
		GroupIDs: []string{"-2"},
		Schema:   "rprj",
	}

	repo := NewDBRepository(dbContext, Factory, DbConnection)
	repo.Verbose = false
	// Create a test object
	dbObj := repo.GetInstanceByTableName("notes")
	if dbObj == nil {
		t.Fatal("Failed to create DBObject instance")
	}
	dbObj.SetValue("name", "Test Note for GetById")
	dbObj.SetValue("description", "This is a test note object")

	createdObj, err := repo.Insert(dbObj)
	if err != nil {
		t.Fatalf("Failed to create DBObject: %v", err)
	}
	objID := createdObj.GetValue("id").(string)

	// Test GetObjectById
	retrievedObj := repo.ObjectByID(objID, true)
	if retrievedObj == nil {
		t.Fatalf("ObjectByID returned nil for ID %s", objID)
	}
	if retrievedObj.GetValue("name") != "Test Note for GetById" {
		t.Fatalf("Expected name 'Test Note for GetById', got '%v'", retrievedObj.GetValue("name"))
	}
	log.Printf("ObjectByID successful: %v", retrievedObj.ToString())

	// Cleanup
	deletedObject, err := repo.Delete(retrievedObj)
	// deletedObject, err := repo.Delete(createdObj)
	if err != nil {
		t.Fatalf("Failed to delete DBObject during cleanup: %v", err)
	}
	log.Printf("Deleted object: %v", deletedObject.ToString())
	// Second time to force the hard delete
	deletedObject, err = repo.Delete(deletedObject)
	if err != nil {
		t.Fatalf("Failed to delete DBObject during cleanup: %v", err)
	}
	log.Printf("Hard Deleted object: %v", deletedObject.ToString())
}

func TestFullObjectById(t *testing.T) {
	// Setup
	dbContext := &DBContext{
		UserID:   "-1",
		GroupIDs: []string{"-2"},
		Schema:   "rprj",
	}

	repo := NewDBRepository(dbContext, Factory, DbConnection)
	repo.Verbose = false
	// Create a test object
	dbObj := repo.GetInstanceByTableName("pages")
	if dbObj == nil {
		t.Fatal("Failed to create DBObject instance")
	}
	dbObj.SetValue("name", "Test Page for FullObjectById")
	dbObj.SetValue("description", "This is a test page content")
	dbObj.SetValue("html", "<h1>Test Page</h1><p>This is a test page content</p>")

	createdObj, err := repo.Insert(dbObj)
	if err != nil {
		t.Fatalf("Failed to create DBObject: %v", err)
	}
	objID := createdObj.GetValue("id").(string)

	// Test FullObjectByID
	retrievedObj := repo.FullObjectById(objID, true)
	if retrievedObj == nil {
		t.Fatalf("FullObjectByID returned nil for ID %s", objID)
	}
	if retrievedObj.GetValue("name") != "Test Page for FullObjectById" {
		t.Fatalf("Expected name 'Test Page for FullObjectById', got '%v'", retrievedObj.GetValue("name"))
	}
	log.Printf("FullObjectByID successful: %v", retrievedObj.ToString())

	// Cleanup
	deletedObject, err := repo.Delete(retrievedObj)
	if err != nil {
		t.Fatalf("Failed to delete DBObject during cleanup: %v", err)
	}
	log.Printf("Deleted object: %v", deletedObject.ToString())

	// Second time to force the hard delete
	deletedObject, err = repo.Delete(deletedObject)
	if err != nil {
		t.Fatalf("Failed to delete DBObject during cleanup: %v", err)
	}
	log.Printf("Hard Deleted object: %v", deletedObject.ToString())
}

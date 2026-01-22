package dblayer

import (
	"log"
	"os"
	"testing"
)

func TestGetObjectByID(t *testing.T) {
	// Setup
	dbContext := &DBContext{
		UserID:   "-1",
		GroupIDs: []string{"-2"},
		Schema:   DbSchema,
	}

	parent_id_value := "-10"

	repo := NewDBRepository(dbContext, Factory, DbConnection)
	repo.Verbose = false

	// Create parent folder
	testFolder := repo.GetInstanceByTableName("folders")
	if testFolder == nil {
		t.Fatal("Failed to create DBFolder instance")
	}
	testFolder.SetValue("name", "Parent Folder")
	testFolder.SetValue("fk_obj_id", parent_id_value)
	// repo.Verbose = true
	createdFolder, err := repo.Insert(testFolder)
	repo.Verbose = false
	if err != nil {
		t.Fatalf("Failed to create parent folder: %v", err)
	}
	log.Print("Created parent folder: ", createdFolder.ToString())

	// Retrieve folder by ID
	retrievedObj := repo.ObjectByID(createdFolder.GetValue("id").(string), true)
	if retrievedObj == nil {
		t.Fatal("Retrieved object is nil")
	}

	log.Print("Retrieved folder: ", retrievedObj.ToString())

	// Cleanup
	err = hardDeleteForTests(repo, createdFolder.(DBObjectInterface))
	if err != nil {
		t.Fatalf("Failed to hard delete parent folder: %v", err)
	}

}

func TestDBFolderDefaultValues(t *testing.T) {
	// Setup
	dbContext := &DBContext{
		UserID:   "-1",
		GroupIDs: []string{"-2"},
		Schema:   DbSchema,
	}

	fk_obj_id_value := "-10"

	repo := NewDBRepository(dbContext, Factory, DbConnection)
	repo.Verbose = false

	// Create parent folder
	parentFolder := repo.GetInstanceByTableName("folders")
	if parentFolder == nil {
		t.Fatal("Failed to create DBFolder instance")
	}
	parentFolder.SetValue("name", "Parent Folder")
	parentFolder.SetValue("fk_obj_id", fk_obj_id_value)
	repo.Verbose = true
	createdParent, err := repo.Insert(parentFolder)
	repo.Verbose = false
	if err != nil {
		t.Fatalf("Failed to create parent folder: %v", err)
	}
	log.Print("Created parent folder: ", createdParent.ToString())

	if createdParent.GetValue("fk_obj_id") != fk_obj_id_value {
		err := hardDeleteForTests(repo, createdParent.(DBObjectInterface))
		if err != nil {
			t.Fatalf("Failed to hard delete parent folder: %v", err)
		}
		t.Fatalf("Expected fk_obj_id to be '%s', got '%v'", fk_obj_id_value, createdParent.GetValue("fk_obj_id"))
	}

	// Create child folder without setting some fields
	childFolder := repo.GetInstanceByTableName("folders").(DBObjectInterface)
	if childFolder == nil {
		t.Fatal("Failed to create DBFolder instance")
	}
	childFolder.SetValue("name", "Child Folder")
	childFolder.SetValue("father_id", createdParent.GetValue("id"))

	repo.Verbose = false
	childFolder.SetDefaultValues(repo)
	repo.Verbose = false

	log.Print("Child folder after SetDefaultValues: ", childFolder.ToString())

	log.Print("Parent folder.fk_obj_id: ", createdParent.GetValue("fk_obj_id"))
	log.Print(" Child folder.fk_obj_id: ", childFolder.GetValue("fk_obj_id"))

	// Verify default values
	if childFolder.GetValue("fk_obj_id") != createdParent.GetValue("fk_obj_id") {
		t.Fatalf("Expected fk_obj_id to be '%v', got '%v'", createdParent.GetValue("fk_obj_id"), childFolder.GetValue("fk_obj_id"))
	}

	// Optionally, insert the child folder to verify no errors occur
	// createdChild, err := repo.Insert(childFolder)

	// Delete parent folder
	err = hardDeleteForTests(repo, createdParent.(DBObjectInterface))
	if err != nil {
		t.Fatalf("Failed to hard delete parent folder: %v", err)
	}
}

func TestDBFileUpload(t *testing.T) {
	repo := setupTestRepo(t)

	// Create and upload file - simplified!
	file := createTestFile(t, repo, "testdata/images/test_image.jpg", map[string]any{
		"name":        "Test Image",
		"description": "Test image for upload",
	}, nil)

	// Verify file was renamed with prefix
	filename := file.GetValue("filename").(string)
	expectedPrefix := "r_" + file.GetValue("id").(string) + "_"
	if len(filename) < len(expectedPrefix) || filename[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("Expected filename to start with '%s', got '%s'", expectedPrefix, filename)
	}

	// Verify checksum was calculated
	if file.GetValue("checksum") == nil || file.GetValue("checksum").(string) == "" {
		t.Error("Expected checksum to be calculated")
	}

	// Verify MIME type was detected
	mime := file.GetValue("mime").(string)
	if mime != "image/jpeg" && mime[:10] != "image/jpeg" {
		t.Errorf("Expected MIME type to be 'image/jpeg', got '%s'", mime)
	}

	// Verify thumbnail was created
	fullpath := file.GetFullpath(nil)
	thumbPath := fullpath + "_thumb.jpg"
	if _, err := os.Stat(thumbPath); os.IsNotExist(err) {
		t.Errorf("Expected thumbnail to be created at '%s'", thumbPath)
	} else {
		log.Printf("Thumbnail created at: %s", thumbPath)
	}

	// Cleanup: hard delete file (removes physical file and thumbnail)
	err := hardDeleteForTests(repo, file)
	if err != nil {
		t.Fatalf("Failed to hard delete file: %v", err)
	}

	// Verify physical file was deleted
	if _, err := os.Stat(fullpath); !os.IsNotExist(err) {
		t.Errorf("Expected physical file to be deleted at '%s'", fullpath)
	}

	// Verify thumbnail was deleted
	if _, err := os.Stat(thumbPath); !os.IsNotExist(err) {
		t.Errorf("Expected thumbnail to be deleted at '%s'", thumbPath)
	}
}

func TestDBFileUploadPNG(t *testing.T) {
	repo := setupTestRepo(t)

	file := createTestFile(t, repo, "testdata/images/test_image.png", map[string]any{
		"name": "Test PNG Image",
	}, nil)
	mime := file.GetValue("mime").(string)
	if mime != "image/png" && mime[:9] != "image/png" {
		t.Errorf("Expected MIME type to be 'image/png', got '%s'", mime)
	}

	// Verify thumbnail created
	fullpath := file.GetFullpath(nil)
	thumbPath := fullpath + "_thumb.jpg"
	if _, err := os.Stat(thumbPath); os.IsNotExist(err) {
		t.Errorf("Expected thumbnail to be created at '%s'", thumbPath)
	}

	// Cleanup
	err := hardDeleteForTests(repo, file)
	if err != nil {
		t.Fatalf("Failed to hard delete PNG file: %v", err)
	}
}

func TestDBFileUploadGIF(t *testing.T) {
	repo := setupTestRepo(t)

	file := createTestFile(t, repo, "testdata/images/test_image.gif", map[string]any{
		"name": "Test GIF Image",
	}, nil)
	mime := file.GetValue("mime").(string)
	if mime != "image/gif" && mime[:9] != "image/gif" {
		t.Errorf("Expected MIME type to be 'image/gif', got '%s'", mime)
	}

	// Verify thumbnail created
	fullpath := file.GetFullpath(nil)
	thumbPath := fullpath + "_thumb.jpg"
	if _, err := os.Stat(thumbPath); os.IsNotExist(err) {
		t.Errorf("Expected thumbnail to be created at '%s'", thumbPath)
	}

	// Cleanup
	err := hardDeleteForTests(repo, file)
	if err != nil {
		t.Fatalf("Failed to hard delete GIF file: %v", err)
	}
}

func TestDBFileUploadNonImage(t *testing.T) {
	repo := setupTestRepo(t)

	file := createTestFile(t, repo, "testdata/files/test_document.txt", map[string]any{
		"name": "Test Text Document",
	}, nil)
	mime := file.GetValue("mime").(string)
	if mime[:10] != "text/plain" {
		t.Errorf("Expected MIME type to start with 'text/plain', got '%s'", mime)
	}

	// Verify NO thumbnail was created (not an image)
	fullpath := file.GetFullpath(nil)
	thumbPath := fullpath + "_thumb.jpg"
	if _, err := os.Stat(thumbPath); !os.IsNotExist(err) {
		t.Errorf("Expected NO thumbnail for text file, but found one at '%s'", thumbPath)
	}

	// Cleanup
	err := hardDeleteForTests(repo, file)
	if err != nil {
		t.Fatalf("Failed to hard delete text file: %v", err)
	}
}

func TestDBFileUploadPDF(t *testing.T) {
	repo := setupTestRepo(t)

	file := createTestFile(t, repo, "testdata/files/test_document.pdf", map[string]any{
		"name": "Test PDF Document",
	}, nil)

	mime := file.GetValue("mime").(string)
	if mime != "application/pdf" && mime[:15] != "application/pdf" {
		t.Errorf("Expected MIME type to be 'application/pdf', got '%s'", mime)
	}

	// Verify NO thumbnail was created (not an image)
	fullpath := file.GetFullpath(nil)
	thumbPath := fullpath + "_thumb.jpg"
	if _, err := os.Stat(thumbPath); !os.IsNotExist(err) {
		t.Errorf("Expected NO thumbnail for PDF file, but found one at '%s'", thumbPath)
	}

	// Cleanup
	err := hardDeleteForTests(repo, file)
	if err != nil {
		t.Fatalf("Failed to hard delete PDF file: %v", err)
	}
}

func TestDBFileSmallImage(t *testing.T) {
	repo := setupTestRepo(t)

	// Test with small image (50x50) - should still create thumbnail
	file := createTestFile(t, repo, "testdata/images/small_image.jpg", map[string]any{
		"name": "Small Test Image",
	}, nil)

	// Verify thumbnail was still created even for small image
	fullpath := file.GetFullpath(nil)
	thumbPath := fullpath + "_thumb.jpg"
	if _, err := os.Stat(thumbPath); os.IsNotExist(err) {
		t.Errorf("Expected thumbnail to be created even for small image at '%s'", thumbPath)
	}

	// Cleanup
	err := hardDeleteForTests(repo, file)
	if err != nil {
		t.Fatalf("Failed to hard delete small image: %v", err)
	}
}

func TestDBFileUpdate(t *testing.T) {
	repo := setupTestRepo(t)

	// Create initial file
	file := createTestFile(t, repo, "testdata/images/test_image.jpg", map[string]any{
		"name": "Original File",
	}, nil)

	originalFilename := file.GetValue("filename").(string)
	originalChecksum := file.GetValue("checksum").(string)

	log.Printf("Original file: %s, checksum: %s", originalFilename, originalChecksum)

	// Update file metadata (change name and description)
	file.SetValue("name", "Updated File Name")
	file.SetValue("description", "Updated description")

	updated, err := repo.Update(file)
	if err != nil {
		t.Fatalf("Failed to update file: %v", err)
	}

	updatedFile := updated.(*DBFile)

	// Verify metadata updated
	if updatedFile.GetValue("name").(string) != "Updated File Name" {
		t.Errorf("Expected name to be 'Updated File Name', got '%s'", updatedFile.GetValue("name").(string))
	}

	// Verify filename and checksum unchanged (same physical file)
	if updatedFile.GetValue("filename").(string) != originalFilename {
		t.Errorf("Expected filename to remain '%s', got '%s'", originalFilename, updatedFile.GetValue("filename").(string))
	}

	if updatedFile.GetValue("checksum").(string) != originalChecksum {
		t.Errorf("Expected checksum to remain '%s', got '%s'", originalChecksum, updatedFile.GetValue("checksum").(string))
	}

	// Cleanup
	err = hardDeleteForTests(repo, updatedFile)
	if err != nil {
		t.Fatalf("Failed to hard delete file: %v", err)
	}
}

func TestDBFileSoftDelete(t *testing.T) {
	repo := setupTestRepo(t)

	// Create file
	file := createTestFile(t, repo, "testdata/images/test_image.jpg", map[string]any{
		"name": "File to Delete",
	}, nil)

	fullpath := file.GetFullpath(nil)
	thumbPath := fullpath + "_thumb.jpg"

	// Verify file and thumbnail exist
	if _, err := os.Stat(fullpath); os.IsNotExist(err) {
		t.Fatalf("Expected file to exist at '%s'", fullpath)
	}
	if _, err := os.Stat(thumbPath); os.IsNotExist(err) {
		t.Fatalf("Expected thumbnail to exist at '%s'", thumbPath)
	}

	// Soft delete (first delete)
	deleted, err := repo.Delete(file)
	if err != nil {
		t.Fatalf("Failed to soft delete file: %v", err)
	}

	deletedFile := deleted.(*DBFile)

	// Verify file still exists physically after soft delete
	if _, err := os.Stat(fullpath); os.IsNotExist(err) {
		t.Errorf("Expected file to still exist after soft delete at '%s'", fullpath)
	}

	// Verify thumbnail still exists after soft delete
	if _, err := os.Stat(thumbPath); os.IsNotExist(err) {
		t.Errorf("Expected thumbnail to still exist after soft delete at '%s'", thumbPath)
	}

	// Hard delete (second delete)
	err = hardDeleteForTests(repo, deletedFile)
	if err != nil {
		t.Fatalf("Failed to hard delete file: %v", err)
	}

	// Verify file removed physically after hard delete
	if _, err := os.Stat(fullpath); !os.IsNotExist(err) {
		t.Errorf("Expected file to be deleted after hard delete at '%s'", fullpath)
	}

	// Verify thumbnail removed after hard delete
	if _, err := os.Stat(thumbPath); !os.IsNotExist(err) {
		t.Errorf("Expected thumbnail to be deleted after hard delete at '%s'", thumbPath)
	}
}

func TestDBFileWithFolder(t *testing.T) {
	repo := setupTestRepo(t)

	fk_obj_id_value := "-10"

	// Create parent folder
	folder := createTestFolder(t, repo, map[string]any{
		"name":      "Test Folder",
		"fk_obj_id": fk_obj_id_value,
	}, nil)

	// Upload file into folder
	file := createTestFile(t, repo, "testdata/images/test_image.jpg", map[string]any{
		"name":      "File in Folder",
		"father_id": folder.GetValue("id"),
	}, nil)

	// Verify fk_obj_id inherited from folder
	if file.GetValue("fk_obj_id") != fk_obj_id_value {
		t.Errorf("Expected fk_obj_id to be '%s' (inherited from folder), got '%v'", fk_obj_id_value, file.GetValue("fk_obj_id"))
	}

	log.Printf("File fk_obj_id: %v (inherited from folder)", file.GetValue("fk_obj_id"))

	// Cleanup
	err := hardDeleteForTests(repo, file)
	if err != nil {
		t.Fatalf("Failed to hard delete file: %v", err)
	}

	err = hardDeleteForTests(repo, folder.(DBObjectInterface))
	if err != nil {
		t.Fatalf("Failed to hard delete folder: %v", err)
	}
}

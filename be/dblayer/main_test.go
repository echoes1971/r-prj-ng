package dblayer

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"rprj/be/models"
	"testing"
)

// go test -v ./dblayer -run TestXXX -config ../config.json
var (
	configFile = flag.String("config", "../config.json", "Path to configuration file")
)

func TestMain(m *testing.M) {
	// config := models.Config{
	// 	AppName:        "R-Project Test Suite",
	// 	ServerPort:     1971,
	// 	DBEngine:       "mysql",
	// 	DBUrl:          "root:mysecret@tcp(localhost:3306)/rproject",
	// 	TablePrefix:    "rprj",
	// 	RootDirectory:  ".",
	// 	FilesDirectory: "files",
	// }
	flag.Parse()
	var config models.Config
	err := models.LoadConfig(*configFile, &config)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}
	config.RootDirectory, err = filepath.Abs(config.RootDirectory)
	if err != nil {
		log.Fatalf("Error getting absolute path of root directory: %v", err)
	}
	log.Printf("Using root directory: %s", config.RootDirectory)
	config.FilesDirectory = "test_files"
	log.Printf("Using files directory: %s", config.FilesDirectory)

	InitDBLayer(config)
	EnsureDBSchema(false)
	InitDBData()

	// Esegui i test
	exitCode := m.Run()

	// Teardown: chiudi la connessione
	CloseDBConnection()

	os.Exit(exitCode)
}

/* ***** Helper functions for tests ***** */

func RandInt(min, max int) int {
	return min + rand.Intn(max-min)
}

/* Returns a random 4-digit string */
func Random4digits() string {
	const digits = "0123456789"
	result := make([]byte, 4)
	// Generate random number between 0000 and 9999
	for i := 0; i < 4; i++ {
		result[i] = digits[RandInt(0, len(digits))]
	}
	return string(result)
}

func hardDeleteForTests(repo *DBRepository, object DBObjectInterface) error {
	deletedObject, err := repo.Delete(object)
	if err != nil {
		return err
	}
	// Second time to force the hard delete
	deletedObject, err = repo.Delete(deletedObject)
	if err != nil {
		return err
	}
	return nil
}

// setupTestRepo creates a test repository with standard test context
func setupTestRepo(t *testing.T) *DBRepository {
	dbContext := &DBContext{
		UserID:   "-1",
		GroupIDs: []string{"-2"},
		Schema:   DbSchema,
	}
	repo := NewDBRepository(dbContext, Factory, DbConnection)
	repo.Verbose = false
	return repo
}
func SetupTestRepo(t *testing.T, user_id string, group_ids []string, schema string) *DBRepository {
	dbContext := &DBContext{
		UserID:   user_id,
		GroupIDs: group_ids,
		Schema:   schema,
	}
	repo := NewDBRepository(dbContext, Factory, DbConnection)
	repo.Verbose = false
	return repo
}

// prepareTestFile copies a test file from testdata to upload directory
func prepareTestFile(t *testing.T, srcPath, destFilename string) string {
	uploadDir := dbFiles_root_directory + "/" + dbFiles_dest_directory
	log.Print("Preparing test file. Source:", srcPath, " Destination dir:", uploadDir)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		t.Fatalf("Failed to create upload directory: %v", err)
	}

	srcFile, err := os.Open(srcPath)
	if err != nil {
		t.Fatalf("Failed to open test file %s: %v", srcPath, err)
	}
	defer srcFile.Close()

	destPath := filepath.Join(uploadDir, destFilename)
	destFile, err := os.Create(destPath)
	if err != nil {
		t.Fatalf("Failed to create destination file: %v", err)
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(srcFile)
	if err != nil {
		t.Fatalf("Failed to copy test file: %v", err)
	}

	return destFilename
}

// createTestObject creates an entity with the provided values using repo.CreateObject
// Usage: createTestObject(t, repo, "files", map[string]any{"name": "Test", "filename": "test.jpg"})
func createTestObject(t *testing.T, repo *DBRepository, tableName string, values map[string]any, metadata map[string]any) DBEntityInterface {
	created, err := repo.CreateObject(tableName, values, metadata)
	if err != nil {
		t.Fatalf("Failed to create %s: %v", tableName, err)
	}
	return created
}

// createTestFile creates a DBFile with automatic file preparation
// Usage: createTestFile(t, repo, "dblayer/dblayer/testdata/images/test.jpg", map[string]any{"name": "Test Image"})
func createTestFile(t *testing.T, repo *DBRepository, srcPath string, values map[string]any, metadata map[string]any) *DBFile {
	// Generate unique filename
	filename := filepath.Base(srcPath)
	prepareTestFile(t, srcPath, filename)

	// Set filename if not provided
	if _, ok := values["filename"]; !ok {
		values["filename"] = filename
	}

	created := createTestObject(t, repo, "files", values, metadata)
	return created.(*DBFile)
}

// createTestFolder creates a DBFolder
// Usage: createTestFolder(t, repo, map[string]any{"name": "Test Folder", "fk_obj_id": "-10"})
func createTestFolder(t *testing.T, repo *DBRepository, values map[string]any, metadata map[string]any) DBEntityInterface {
	return createTestObject(t, repo, "folders", values, metadata)
}

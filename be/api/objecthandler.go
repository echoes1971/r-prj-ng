package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"rprj/be/dblayer"

	"github.com/gorilla/mux"
)

// CreateObjectHandler creates a new DBObject
// POST /api/objects
// Body: {"classname": "DBNote", "father_id": "123", "name": "My Note", "description": "...", ...}
func CreateObjectHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := GetClaimsFromRequest(r)
	if err != nil {
		RespondSimpleError(w, ErrUnauthorized, "Unauthorized", http.StatusUnauthorized)
		return
	}

	dbContext := &dblayer.DBContext{
		UserID:   claims["user_id"],
		GroupIDs: strings.Split(claims["groups"], ","),
		Schema:   dblayer.DbSchema,
	}

	repo := dblayer.NewDBRepository(dbContext, dblayer.Factory, dblayer.DbConnection)
	repo.Verbose = false

	var requestData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		log.Printf("CreateObjectHandler: Failed to decode request body: %v", err)
		RespondSimpleError(w, ErrInvalidRequest, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Extract classname to determine table
	classname, ok := requestData["classname"].(string)
	if !ok || classname == "" {
		RespondSimpleError(w, ErrMissingField, "Missing or invalid 'classname' field", http.StatusBadRequest)
		return
	}

	// Get instance to determine table name
	instance := repo.GetInstanceByClassName(classname)
	if instance == nil {
		RespondSimpleError(w, ErrInvalidRequest, "Unknown classname: "+classname, http.StatusBadRequest)
		return
	}

	if !instance.IsDBObject() {
		RespondSimpleError(w, ErrInvalidRequest, "Classname is not a DBObject: "+classname, http.StatusBadRequest)
		return
	}

	tableName := instance.GetTableName()

	// Remove classname from values (not a DB field)
	delete(requestData, "classname")

	// Set automatic fields
	requestData["owner"] = dbContext.UserID
	if len(dbContext.GroupIDs) > 0 {
		requestData["group_id"] = dbContext.GroupIDs[0]
	}
	requestData["creator"] = dbContext.UserID
	requestData["creation_date"] = time.Now().Format("2006-01-02 15:04:05")
	requestData["last_modify"] = dbContext.UserID
	requestData["last_modify_date"] = time.Now().Format("2006-01-02 15:04:05")

	fatherID, _ := requestData["father_id"].(string)
	if len(fatherID) == 18 {
		fatherID = strings.ReplaceAll(fatherID, "-", "")
		requestData["father_id"] = fatherID
	}

	// Set default permissions if not provided
	if _, ok := requestData["permissions"]; !ok {
		requestData["permissions"] = "rwxr-x---" // Owner: rwx, Group: r-x, Others: ---
	}

	// Create the object
	// TODO: pass metadata if any
	log.Print("CreateObjectHandler: requestData ", requestData)
	created, err := repo.CreateObject(tableName, requestData, nil)
	if err != nil {
		log.Printf("CreateObjectHandler: Failed to create object: %v", err)
		RespondSimpleError(w, ErrInternalServer, "Failed to create object: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("CreateObjectHandler: Created %s with ID=%s", classname, created.GetValue("id"))

	// Convert entity to map
	resultMap := created.GetAllValues()

	// Return created object
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    resultMap,
		"metadata": map[string]interface{}{
			"classname": classname,
		},
	})
}

// UpdateObjectHandler updates an existing DBObject
// PUT /api/objects/{id}
// Body: {"name": "Updated Name", "description": "Updated description", ...}
func UpdateObjectHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := GetClaimsFromRequest(r)
	if err != nil {
		RespondSimpleError(w, ErrUnauthorized, "Unauthorized", http.StatusUnauthorized)
		return
	}

	dbContext := &dblayer.DBContext{
		UserID:   claims["user_id"],
		GroupIDs: strings.Split(claims["groups"], ","),
		Schema:   dblayer.DbSchema,
	}

	repo := dblayer.NewDBRepository(dbContext, dblayer.Factory, dblayer.DbConnection)
	repo.Verbose = false

	vars := mux.Vars(r)
	objectID := vars["id"]

	if objectID == "" {
		RespondSimpleError(w, ErrInvalidRequest, "Missing object ID", http.StatusBadRequest)
		return
	}
	if len(objectID) == 18 {
		objectID = strings.ReplaceAll(objectID, "-", "")
	}

	// Get existing object to determine classname and check permissions
	existingObj := repo.ObjectByID(objectID, true)
	if existingObj == nil {
		RespondSimpleError(w, ErrObjectNotFound, "Object not found", http.StatusNotFound)
		return
	}

	// Check write permission
	if !repo.CheckWritePermission(existingObj) {
		RespondSimpleError(w, ErrForbidden, "You don't have permission to edit this object", http.StatusForbidden)
		return
	}

	// Get classname from metadata
	classname, ok := existingObj.GetMetadata("classname").(string)
	if !ok {
		RespondSimpleError(w, ErrInternalServer, "Object classname not found", http.StatusInternalServerError)
		return
	}

	// Get full object
	fullObj := repo.FullObjectById(objectID, true)
	if fullObj == nil {
		RespondSimpleError(w, ErrObjectNotFound, "Full object not found", http.StatusNotFound)
		return
	}

	tableName := fullObj.GetTableName()

	// Decode update values based on Content-Type
	var updateValues map[string]interface{}
	var metadataValues map[string]interface{}
	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		// Parse multipart form for file uploads
		err := r.ParseMultipartForm(32 << 20) // 32 MB max
		if err != nil {
			log.Printf("UpdateObjectHandler: Failed to parse multipart form: %v", err)
			RespondSimpleError(w, ErrInvalidRequest, "Invalid multipart form", http.StatusBadRequest)
			return
		}

		updateValues = make(map[string]interface{})
		metadataValues = make(map[string]interface{})

		// Extract form fields
		for key, values := range r.MultipartForm.Value {
			if len(values) > 0 {
				updateValues[key] = values[0]
			}
		}

		// Handle file upload if present
		file, header, err := r.FormFile("file")
		if err == nil {
			defer file.Close()

			// Create files directory if it doesn't exist
			filesDir := filepath.Join(dbFiles_root_directory, dbFiles_dest_directory)
			log.Print("UpdateObjectHandler: filesDir=", filesDir)
			if err := os.MkdirAll(filesDir, 0755); err != nil {
				log.Printf("UpdateObjectHandler: Failed to create files directory: %v", err)
				RespondSimpleError(w, ErrInternalServer, "Failed to create storage directory", http.StatusInternalServerError)
				return
			}

			// Generate filename with r_{id}_ prefix
			baseFilename := filepath.Base(header.Filename)
			savedFilename := "r_" + objectID + "_" + baseFilename
			filePath := filepath.Join(filesDir, savedFilename)

			// Create destination file
			dst, err := os.Create(filePath)
			if err != nil {
				log.Printf("UpdateObjectHandler: Failed to create file %s: %v", filePath, err)
				RespondSimpleError(w, ErrInternalServer, "Failed to save file", http.StatusInternalServerError)
				return
			}
			defer dst.Close()

			// Copy uploaded file to destination
			if _, err := io.Copy(dst, file); err != nil {
				log.Printf("UpdateObjectHandler: Failed to copy file data: %v", err)
				os.Remove(filePath) // Clean up partial file
				RespondSimpleError(w, ErrInternalServer, "Failed to save file data", http.StatusInternalServerError)
				return
			}

			// Update database fields
			updateValues["filename"] = savedFilename
			updateValues["mime"] = header.Header.Get("Content-Type")
			// updateValues["path"] = filePath
			metadataValues["path"] = filePath

			log.Printf("UpdateObjectHandler: File saved successfully: %s (%s)", savedFilename, header.Header.Get("Content-Type"))
		}
	} else {
		// Parse JSON body
		if err := json.NewDecoder(r.Body).Decode(&updateValues); err != nil {
			log.Printf("UpdateObjectHandler: Failed to decode request body: %v", err)
			RespondSimpleError(w, ErrInvalidRequest, "Invalid request body", http.StatusBadRequest)
			return
		}
	}

	// Set automatic update fields
	updateValues["last_modify"] = dbContext.UserID
	updateValues["last_modify_date"] = time.Now().Format("2006-01-02 15:04:05")

	// Remove protected fields that shouldn't be updated via API
	delete(updateValues, "id")
	delete(updateValues, "owner")
	delete(updateValues, "creator")
	delete(updateValues, "creation_date")
	delete(updateValues, "deleted_by")
	delete(updateValues, "deleted_date")

	// Update the object
	updated, err := repo.UpdateObject(tableName, objectID, updateValues, metadataValues)
	if err != nil {
		log.Printf("UpdateObjectHandler: Failed to update object: %v", err)
		RespondSimpleError(w, ErrInternalServer, "Failed to update object: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("UpdateObjectHandler: Updated %s with ID=%s", classname, objectID)

	resultMap := updated.GetAllValues()

	// Return updated object
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    resultMap,
		"metadata": map[string]interface{}{
			"classname": classname,
		},
	})
}

// DeleteObjectHandler soft-deletes a DBObject
// DELETE /api/objects/{id}
func DeleteObjectHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := GetClaimsFromRequest(r)
	if err != nil {
		RespondSimpleError(w, ErrUnauthorized, "Unauthorized", http.StatusUnauthorized)
		return
	}

	dbContext := &dblayer.DBContext{
		UserID:   claims["user_id"],
		GroupIDs: strings.Split(claims["groups"], ","),
		Schema:   dblayer.DbSchema,
	}

	repo := dblayer.NewDBRepository(dbContext, dblayer.Factory, dblayer.DbConnection)
	repo.Verbose = false

	vars := mux.Vars(r)
	objectID := vars["id"]

	if objectID == "" {
		RespondSimpleError(w, ErrInvalidRequest, "Missing object ID", http.StatusBadRequest)
		return
	}
	if len(objectID) == 18 {
		objectID = strings.ReplaceAll(objectID, "-", "")
	}

	// Get existing object to check permissions
	existingObj := repo.ObjectByID(objectID, true)
	if existingObj == nil {
		RespondSimpleError(w, ErrObjectNotFound, "Object not found", http.StatusNotFound)
		return
	}

	// Check write permission (needed to delete)
	if !repo.CheckWritePermission(existingObj) {
		RespondSimpleError(w, ErrForbidden, "You don't have permission to delete this object", http.StatusForbidden)
		return
	}

	// Get classname from metadata
	classname, ok := existingObj.GetMetadata("classname").(string)
	if !ok {
		RespondSimpleError(w, ErrInternalServer, "Object classname not found", http.StatusInternalServerError)
		return
	}

	// Get full object
	fullObj := repo.FullObjectById(objectID, true)
	if fullObj == nil {
		RespondSimpleError(w, ErrObjectNotFound, "Full object not found", http.StatusNotFound)
		return
	}

	// Soft delete (sets deleted_date and deleted_by)
	deleted, err := repo.Delete(fullObj)
	if err != nil {
		log.Printf("DeleteObjectHandler: Failed to delete object: %v", err)
		RespondSimpleError(w, ErrInternalServer, "Failed to delete object: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("DeleteObjectHandler: Soft-deleted %s with ID=%s", classname, objectID)

	// Convert entity to map
	resultMap := deleted.GetAllValues()

	// Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Object deleted successfully",
		"data":    resultMap,
	})
}

// GetCreatableTypesHandler returns the list of DBObject types that can be created as children of a given parent
// GET /api/objects/creatable-types?father_id=123
// GET /api/objects/creatable-types (returns all DBObject types if no father_id)
func GetCreatableTypesHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := GetClaimsFromRequest(r)
	if err != nil {
		RespondSimpleError(w, ErrUnauthorized, "Unauthorized", http.StatusUnauthorized)
		return
	}

	dbContext := &dblayer.DBContext{
		UserID:   claims["user_id"],
		GroupIDs: strings.Split(claims["groups"], ","),
		Schema:   dblayer.DbSchema,
	}

	repo := dblayer.NewDBRepository(dbContext, dblayer.Factory, dblayer.DbConnection)

	fatherID := r.URL.Query().Get("father_id")
	if len(fatherID) == 18 {
		fatherID = strings.ReplaceAll(fatherID, "-", "")
	}

	var creatableTypes []string

	if fatherID == "" || fatherID == "0" {
		// No father specified - return all DBObject types
		for _, className := range dblayer.Factory.GetAllClassNames() {
			instance := dblayer.Factory.GetInstanceByClassName(className)
			if instance != nil && instance.IsDBObject() {
				creatableTypes = append(creatableTypes, className)
			}
		}
	} else {
		// Father specified - get parent object and check TableChildren
		parentObj := repo.ObjectByID(fatherID, true)
		if parentObj == nil {
			RespondSimpleError(w, ErrInternalServer, "Parent object not found", http.StatusNotFound)
			return
		}

		// Check read permission on parent
		if !repo.CheckReadPermission(parentObj) {
			RespondSimpleError(w, ErrForbidden, "No permission to access parent object", http.StatusForbidden)
			return
		}

		parentInstance := dblayer.Factory.GetInstanceByTableName(parentObj.GetTableName())
		if parentInstance == nil {
			RespondSimpleError(w, ErrInternalServer, "Unknown parent type", http.StatusInternalServerError)
			return
		}

		parentTable := parentInstance.GetTableName()

		// Get allowed child tables from TableChildren
		if childTables, exists := dblayer.Factory.TableChildren[parentTable]; exists {
			for _, childTable := range childTables {
				childInstance := dblayer.Factory.GetInstanceByTableName(childTable)
				if childInstance != nil {
					creatableTypes = append(creatableTypes, childInstance.GetTypeName())
				}
			}
		} else {
			// No restrictions defined - allow all DBObject types
			for _, className := range dblayer.Factory.GetAllClassNames() {
				instance := dblayer.Factory.GetInstanceByClassName(className)
				if instance != nil && instance.IsDBObject() {
					creatableTypes = append(creatableTypes, className)
				}
			}
		}
	}

	// Return the list
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"types":   creatableTypes,
	})
}

// SearchObjectsHandler searches for objects by classname and name pattern
// GET /api/objects/search?classname=DBCompany&name=acme&limit=10
func SearchObjectsHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := GetClaimsFromRequest(r)
	if err != nil {
		RespondSimpleError(w, ErrUnauthorized, "Unauthorized", http.StatusUnauthorized)
		return
	}

	dbContext := &dblayer.DBContext{
		UserID:   claims["user_id"],
		GroupIDs: strings.Split(claims["groups"], ","),
		Schema:   dblayer.DbSchema,
	}

	repo := dblayer.NewDBRepository(dbContext, dblayer.Factory, dblayer.DbConnection)
	repo.Verbose = false

	classname := r.URL.Query().Get("classname")
	namePattern := r.URL.Query().Get("name")
	limit := r.URL.Query().Get("limit")

	if classname == "" {
		RespondSimpleError(w, ErrInvalidRequest, "Missing classname parameter", http.StatusBadRequest)
		return
	}

	// Get instance for the requested classname
	searchInstance := repo.GetInstanceByClassName(classname)
	if searchInstance == nil {
		RespondSimpleError(w, ErrInvalidRequest, "Unknown classname: "+classname, http.StatusBadRequest)
		return
	}

	// if !searchInstance.IsDBObject() {
	// 	RespondSimpleError(w, ErrInvalidRequest, "Classname is not a DBObject: "+classname, http.StatusBadRequest)
	// 	return
	// }

	// Set search criteria
	if namePattern != "" {
		switch classname {
		case "DBCountry":
			// For countries, search by Common_Name field
			searchInstance.SetValue("Common_Name", namePattern)
		case "DBUser":
			// For users, search by login field
			searchInstance.SetValue("login", namePattern)
		default:
			searchInstance.SetValue("name", namePattern)
		}
	}
	orderBy := "name"
	switch classname {
	case "DBUser":
		orderBy = "login"
	case "DBCountry":
		orderBy = "Common_Name"
	}

	// Search with LIKE and case-insensitive
	results, err := repo.Search(searchInstance, true, false, orderBy)
	if err != nil {
		log.Printf("SearchObjectsHandler: Search failed: %v", err)
		RespondSimpleError(w, ErrInternalServer, "Search failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Apply limit if specified
	maxResults := len(results)
	if limit != "" {
		var limitInt int
		if _, err := fmt.Sscanf(limit, "%d", &limitInt); err == nil && limitInt > 0 && limitInt < maxResults {
			maxResults = limitInt
		}
	}

	// Convert results to map array
	var resultList []map[string]interface{}
	for i := 0; i < maxResults && i < len(results); i++ {
		entity := results[i]
		// Check read permission
		if !repo.CheckReadPermission(entity) {
			continue
		}

		resultMap := make(map[string]interface{})
		resultMap["id"] = entity.GetValue("id")
		switch classname {
		case "DBCountry":
			resultMap["name"] = entity.GetValue("Common_Name")
		case "DBUser":
			resultMap["name"] = entity.GetValue("login")
		default:
			resultMap["name"] = entity.GetValue("name")
		}
		if desc := entity.GetValue("description"); desc != nil {
			resultMap["description"] = desc
		}
		if desc := entity.GetValue("fullname"); desc != nil {
			resultMap["description"] = desc
		}
		resultList = append(resultList, resultMap)
	}

	log.Printf("SearchObjectsHandler: Found %d %s objects matching '%s'", len(resultList), classname, namePattern)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"objects": resultList,
	})
}

// DownloadFileHandler serves file content for DBFile objects
// GET /api/files/{id}/download
func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := GetClaimsFromRequest(r)

	var dbContext dblayer.DBContext
	if err == nil {
		dbContext = dblayer.DBContext{
			UserID:   claims["user_id"],
			GroupIDs: strings.Split(claims["groups"], ","),
			Schema:   dblayer.DbSchema,
		}
	} else {
		dbContext = dblayer.DBContext{
			UserID:   "-7",           // Anonymous user
			GroupIDs: []string{"-4"}, // Guests group
			Schema:   dblayer.DbSchema,
		}
	}

	repo := dblayer.NewDBRepository(&dbContext, dblayer.Factory, dblayer.DbConnection)

	vars := mux.Vars(r)
	fileID := vars["id"]

	// Load the DBFile object
	tableName := "files"
	entity := repo.GetEntityByID(tableName, fileID)
	if entity == nil {
		log.Printf("DownloadFileHandler: Failed to load file %s", fileID)
		RespondSimpleError(w, ErrObjectNotFound, "File not found", http.StatusNotFound)
		return
	}

	// Cast to DBFile to use getFullpath
	dbFile, ok := entity.(*dblayer.DBFile)
	if !ok {
		log.Printf("DownloadFileHandler: Entity is not a DBFile")
		RespondSimpleError(w, ErrInternalServer, "Invalid file entity", http.StatusInternalServerError)
		return
	}

	// Get file metadata
	filename := dbFile.GetValue("filename")
	if filename == nil {
		log.Printf("DownloadFileHandler: File %s has no filename", fileID)
		RespondSimpleError(w, ErrInternalServer, "File has no filename", http.StatusInternalServerError)
		return
	}

	mime := dbFile.GetValue("mime")
	if mime == nil {
		mime = "application/octet-stream"
	}

	// Construct file path: ./files/{father_id}//{filename}
	baseDir := filepath.Join(dbFiles_root_directory, dbFiles_dest_directory)
	fatherID := dbFile.GetValue("father_id")
	var filePath string
	if fatherID != nil && fatherID != "" && fatherID != "0" {
		// The beforeUpdate creates: files/{father_id}/$/{filename}
		// This structure includes the filename as both a directory and the final file
		filePath = filepath.Join(baseDir, fatherID.(string), filename.(string))
	} else {
		filePath = filepath.Join(baseDir, filename.(string))
	}

	// Open file from disk
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("DownloadFileHandler: Failed to open file %s: %v", filePath, err)
		RespondSimpleError(w, ErrObjectNotFound, "File not found on disk", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Get file info for size
	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("DownloadFileHandler: Failed to stat file %s: %v", filePath, err)
		RespondSimpleError(w, ErrInternalServer, "Failed to read file info", http.StatusInternalServerError)
		return
	}

	// Set headers
	w.Header().Set("Content-Type", mime.(string))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// For images, display inline; for other files, force download
	if strings.HasPrefix(mime.(string), "image/") {
		w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))
	} else {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	}

	// Stream file to response
	if _, err := io.Copy(w, file); err != nil {
		log.Printf("DownloadFileHandler: Failed to stream file %s: %v", filePath, err)
		return
	}

	log.Printf("DownloadFileHandler: Served file %s (%s)", filename, mime)
}

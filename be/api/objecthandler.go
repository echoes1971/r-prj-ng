package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"rprj/be/dblayer"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

// CreatableTypesResponse godoc
// @Description Response structure for creatable types
type CreatableTypesResponse struct {
	Success bool     `json:"success"`
	Types   []string `json:"types"`
}

// FileTokensRequest godoc
// @Description Request structure for file tokens
type FileTokensRequest struct {
	FileIDs []string `json:"file_ids"`
}

// FileTokensResponse godoc
// @Description Response structure for file tokens
type FileTokensResponse struct {
	Success bool              `json:"success"`
	Tokens  map[string]string `json:"tokens"`
}

// ObjectResponse godoc
// @Description Standard response structure for object operations
type ObjectResponse struct {
	Success  bool                   `json:"success"`
	Data     map[string]interface{} `json:"data"`
	Metadata map[string]interface{} `json:"metadata"`
	Message  string                 `json:"message,omitempty"`
}

// ObjectsSearchResponse godoc
// @Description Response structure for object search
type ObjectsSearchResponse struct {
	Success bool                     `json:"success"`
	Objects []map[string]interface{} `json:"objects"`
}

// CreateObjectHandler godoc
// @Summary Create a new DBObject
// @Description Creates a new DBObject based on the provided classname and fields
// @Tags objects
// @Accept json
// @Produce json
// @Param object body map[string]interface{} true "Object data"
// @Success 201 {object} ObjectResponse "Created object data"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /objects [post]
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

	// Decode request based on Content-Type
	var requestData map[string]interface{}
	var metadataValues map[string]interface{}
	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		// Parse multipart form for file uploads
		err := r.ParseMultipartForm(32 << 20) // 32 MB max
		if err != nil {
			log.Printf("CreateObjectHandler: Failed to parse multipart form: %v", err)
			RespondSimpleError(w, ErrInvalidRequest, "Invalid multipart form", http.StatusBadRequest)
			return
		}

		requestData = make(map[string]interface{})
		metadataValues = make(map[string]interface{})

		// Extract form fields
		for key, values := range r.MultipartForm.Value {
			if len(values) > 0 {
				requestData[key] = values[0]
			}
		}

		// Handle file upload if present - this makes it a DBFile
		file, header, err := r.FormFile("file")
		if err == nil {
			defer file.Close()

			// Create files directory if it doesn't exist
			filesDir := filepath.Join(dbFiles_root_directory, dbFiles_dest_directory)
			log.Print("CreateObjectHandler: filesDir=", filesDir)
			if err := os.MkdirAll(filesDir, 0755); err != nil {
				log.Printf("CreateObjectHandler: Failed to create files directory: %v", err)
				RespondSimpleError(w, ErrInternalServer, "Failed to create storage directory", http.StatusInternalServerError)
				return
			}

			// Generate filename with r_{id}_ prefix
			baseFilename := filepath.Base(header.Filename)
			savedFilename := baseFilename //"r_" + newID + "_" + baseFilename
			filePath := filepath.Join(filesDir, savedFilename)

			// Create destination file
			dst, err := os.Create(filePath)
			if err != nil {
				log.Printf("CreateObjectHandler: Failed to create file %s: %v", filePath, err)
				RespondSimpleError(w, ErrInternalServer, "Failed to save file", http.StatusInternalServerError)
				return
			}
			defer dst.Close()

			// Copy uploaded file to destination
			if _, err := io.Copy(dst, file); err != nil {
				log.Printf("CreateObjectHandler: Failed to copy file data: %v", err)
				os.Remove(filePath) // Clean up partial file
				RespondSimpleError(w, ErrInternalServer, "Failed to save file data", http.StatusInternalServerError)
				return
			}

			// Set file fields - this indicates it's a DBFile
			// requestData["id"] = newID
			requestData["filename"] = savedFilename
			requestData["mime"] = header.Header.Get("Content-Type")
			// metadataValues["path"] = filePath

			log.Printf("CreateObjectHandler: File saved successfully: %s (%s)", savedFilename, header.Header.Get("Content-Type"))
		}
	} else {
		// Parse JSON body
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			log.Printf("CreateObjectHandler: Failed to decode request body: %v", err)
			RespondSimpleError(w, ErrInvalidRequest, "Invalid request body", http.StatusBadRequest)
			return
		}
	}

	// Extract classname to determine table
	classname, ok := requestData["classname"].(string)
	if !ok || classname == "" {
		// If no classname but we have a file, it's a DBFile
		if _, hasFile := requestData["filename"]; hasFile {
			classname = "DBFile"
			requestData["classname"] = classname
		} else {
			RespondSimpleError(w, ErrMissingField, "Missing or invalid 'classname' field", http.StatusBadRequest)
			return
		}
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
	created, err := repo.CreateObject(tableName, requestData, metadataValues)
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
	json.NewEncoder(w).Encode(ObjectResponse{
		Success: true,
		Data:    resultMap,
		Metadata: map[string]interface{}{
			"classname": classname,
		},
	})
}

// Body: {"name": "Updated Name", "description": "Updated description", ...}
// UpdateObjectHandler godoc
// @Summary Update an existing DBObject
// @Description Updates an existing DBObject by its ID
// @Tags objects
// @Accept json
// @Produce json
// @Param id path string true "Object ID"
// @Param object body map[string]interface{} true "Object fields to update"
// @Success 200 {object} ObjectResponse "Updated object data"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Object not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /objects/{id} [put]
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
	// Rule: You cannot change group_id unless you are in that group
	if fullObj.GetValue("group_id") != nil && fullObj.GetValue("group_id") != "" && !slices.Contains(dbContext.GroupIDs, fmt.Sprintf("%v", fullObj.GetValue("group_id"))) {
		log.Printf("UpdateObjectHandler: Removing group_id from updateValues, current=%v, user groups=%v", fullObj.GetValue("group_id"), dbContext.GroupIDs)
		delete(updateValues, "group_id")
	}
	// Rule: You cannot change owner unless you are that owner
	if fullObj.GetValue("owner") != nil && fullObj.GetValue("owner") != "" && fmt.Sprintf("%v", fullObj.GetValue("owner")) != dbContext.UserID {
		log.Printf("UpdateObjectHandler: Removing owner from updateValues, current=%v, user=%v", fullObj.GetValue("owner"), dbContext.UserID)
		delete(updateValues, "owner")
	}
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
	json.NewEncoder(w).Encode(ObjectResponse{
		Success: true,
		Data:    resultMap,
		Metadata: map[string]interface{}{
			"classname": classname,
		},
	})
}

// DeleteObjectHandler godoc
// @Summary Delete a DBObject
// @Description Soft-deletes a DBObject by its ID
// @Tags objects
// @Produce json
// @Param id path string true "Object ID"
// @Success 200 {object} ObjectResponse "Deletion success message"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Object not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /objects/{id} [delete]
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
	existingObj := repo.ObjectByID(objectID, false)
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
	fullObj := repo.FullObjectById(objectID, false)
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
	json.NewEncoder(w).Encode(ObjectResponse{
		Success:  true,
		Message:  "Object deleted successfully",
		Data:     resultMap,
		Metadata: deleted.GetAllMetadata(),
	})
}

// GetCreatableTypesHandler godoc
// @Summary Get creatable object types
// @Description Returns the list of DBObject types that can be created as children of a given parent object, returns all DBObject types if no father_id
// @Tags objects
// @Produce json
// @Param father_id query string false "Father object ID"
// @Success 200 {object} CreatableTypesResponse "List of creatable types"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /objects/creatable-types [get]
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
		parentObj := repo.FullObjectById(fatherID, true)
		// parentObj := repo.ObjectByID(fatherID, true)
		// log.Print("GetCreatableTypesHandler: fatherID=", fatherID, " parentObj=", parentObj.ToString())
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
		// log.Printf("GetCreatableTypesHandler: parentTable=%s, TableChildren keys=%v", parentTable, func() []string {
		// 	keys := make([]string, 0, len(dblayer.Factory.TableChildren))
		// 	for k := range dblayer.Factory.TableChildren {
		// 		keys = append(keys, k)
		// 	}
		// 	return keys
		// }())
		if childTables, exists := dblayer.Factory.TableChildren[parentTable]; exists {
			// log.Printf("GetCreatableTypesHandler: Found %d child tables for %s: %v", len(childTables), parentTable, childTables)
			for _, childTable := range childTables {
				childInstance := dblayer.Factory.GetInstanceByTableName(childTable)
				if childInstance != nil {
					creatableTypes = append(creatableTypes, childInstance.GetTypeName())
				} else {
					log.Printf("GetCreatableTypesHandler: Unknown child table: %s", childTable)
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

	// Sort types alphabetically
	sort.Strings(creatableTypes)

	// Return the list
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(CreatableTypesResponse{
		Success: true,
		Types:   creatableTypes,
	})
}

// SearchObjectsHandler godoc
// @Summary Search objects
// @Description Search for objects by classname, name pattern and other filters
// @Tags objects
// @Produce json
// @Param token header string false "Temporary JWT token for access"
// @Param classname query string true "Class name (e.g., DBCompany, DBNote)"
// @Param name query string false "Name pattern for search"
// @Param searchJson query string false "JSON object with additional search parameters"
// @Param orderBy query string false "Field to order by (e.g., name, creation_date)"
// @Param limit query int false "Maximum number of results"
// @Param offset query int false "Offset for pagination"
// @Param type query string false "Filter type (e.g., 'link' for linkable objects)"
// @Param includeDeleted query string false "Include deleted objects"
// @Success 200 {object} ObjectsSearchResponse "List of matching objects"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 500 {object} ErrorResponse "Internal error"
// @Router /objects/search [get]
func SearchObjectsHandler(w http.ResponseWriter, r *http.Request) {
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
	repo.Verbose = false

	classname := r.URL.Query().Get("classname")
	if classname == "" {
		RespondSimpleError(w, ErrInvalidRequest, "Missing classname parameter", http.StatusBadRequest)
		return
	}

	namePattern := strings.TrimSpace(r.URL.Query().Get("name"))
	// includeDeleted
	includeDeletedParam := r.URL.Query().Get("includeDeleted")
	includeDeleted := includeDeletedParam != "" && includeDeletedParam != "0" && strings.ToLower(includeDeletedParam) != "false"
	log.Print("SearchObjectsHandler: includeDeletedParam=", includeDeletedParam, " includeDeleted=", includeDeleted)
	// orderBy
	orderBy := r.URL.Query().Get("orderBy")
	if orderBy != "" {
		orderBy = strings.TrimSpace(orderBy)
	}

	// searchJson
	searchJson := r.URL.Query().Get("searchJson")
	log.Print("SearchObjectsHandler: searchJson=", searchJson)
	var searchParams map[string]any
	if searchJson != "" {
		err := json.Unmarshal([]byte(searchJson), &searchParams)
		if err != nil {
			log.Printf("SearchObjectsHandler: Failed to parse searchJson: %v", err)
		}
	}
	// var searchParams map[string]any
	// if searchJson != "" {
	// 	err := json.Unmarshal([]byte(searchJson), &searchParams)
	// 	if err != nil {
	// 		log.Printf("SearchObjectsHandler: Failed to parse searchJson: %v", err)
	// 	}
	// }
	log.Print("SearchObjectsHandler: searchParams=", searchParams)
	// limit and offset
	limit := r.URL.Query().Get("limit")
	offset := 0
	if r.URL.Query().Get("offset") != "" {
		fmt.Sscanf(r.URL.Query().Get("offset"), "%d", &offset)
	}
	// type
	searchType := r.URL.Query().Get("type") // optional, "link" to filter only linkable objects i.e. objects I can write

	// Get instance for the requested classname
	searchInstance := repo.GetInstanceByClassName(classname)
	if searchInstance == nil {
		RespondSimpleError(w, ErrInvalidRequest, "Unknown classname: "+classname, http.StatusBadRequest)
		return
	}
	// searchParams has key "$or" ?
	if _, ok := searchParams["$or"]; ok {
		log.Print("SearchObjectsHandler: searchParams has $or")
		searchInstance.SetMetadata("or", searchParams["$or"])
	}

	// if !searchInstance.IsDBObject() {
	// 	log.Print("SearchObjectsHandler: Classname is not a DBObject: ", classname)
	// 	RespondSimpleError(w, ErrInvalidRequest, "Classname is not a DBObject: "+classname, http.StatusBadRequest)
	// 	return
	// }
	searchInstanceDescription := repo.GetInstanceByClassName(classname)

	// Set search criteria
	if searchJson == "" {
		searchInstance.SetValue("name", namePattern)
		searchInstanceDescription.SetValue("description", namePattern)
	} else {
		for key, val := range searchParams {
			// Skip if it starts with $ (handled separately)
			if strings.HasPrefix(key, "$") {
				continue
			}
			searchInstance.SetValue(key, val)
		}
	}

	var results []dblayer.DBEntityInterface
	// Search with LIKE and case-insensitive
	if classname != "DBObject" {
		repo.Verbose = true
		results, err = repo.Search(searchInstance, true, false, orderBy)
		repo.Verbose = false
		if err != nil {
			log.Printf("SearchObjectsHandler: Search failed: %v", err)
			RespondSimpleError(w, ErrInternalServer, "Search failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		var resultsDescription []dblayer.DBEntityInterface
		if searchJson == "" {
			resultsDescription, err = repo.Search(searchInstanceDescription, true, false, orderBy)
		}
		if err != nil {
			log.Printf("SearchObjectsHandler: Search failed: %v", err)
			RespondSimpleError(w, ErrInternalServer, "Search failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		// Merge results
		resultsIDs := make(map[string]bool)
		for _, res := range results {
			id := res.GetValue("id").(string)
			resultsIDs[id] = true
		}
		for _, resDesc := range resultsDescription {
			id := resDesc.GetValue("id").(string)
			if _, exists := resultsIDs[id]; !exists {
				results = append(results, resDesc)
			}
		}
		// IF !includeDeleted, filter out deleted objects
		if !includeDeleted {
			var filteredResults []dblayer.DBEntityInterface
			for _, res := range results {
				if res.GetValue("deleted_date") == nil {
					filteredResults = append(filteredResults, res)
				}
			}
			results = filteredResults
		}
	} else if searchInstance.IsDBObject() {
		// Search by name AND description for better results
		results = repo.SearchByNameAndDescription(namePattern, orderBy, !includeDeleted)
	}

	// Apply limit if specified
	maxResults := len(results)
	log.Print("SearchObjectsHandler: limit=", limit)
	if limit != "" {
		var limitInt int
		if _, err := fmt.Sscanf(limit, "%d", &limitInt); err == nil && limitInt > 0 && limitInt < maxResults {
			maxResults = limitInt
		}
	}
	log.Print("SearchObjectsHandler: offset=", offset)
	log.Print("SearchObjectsHandler: maxResults=", maxResults)

	log.Print("SearchObjectsHandler: results=", len(results))
	log.Print("SearchObjectsHandler: classname=", classname)
	// Convert results to map array
	var resultList []map[string]interface{}
	for i := 0; i < len(results); i++ {
		entity := results[i]
		if entity.HasMetadata("classname") && entity.GetMetadata("classname") == "DBFile" {
			// read the full object to get file metadata, so we can display an image preview
			entity = repo.FullObjectById(entity.GetValue("id").(string), !includeDeleted)
			// TODO verify that this should be redundant with the includeDeleted above
			if entity == nil {
				log.Printf("SearchObjectsHandler: It has been soft deleted ID=%s", results[i].GetValue("id").(string))
				continue
			}
		}
		// IF searched classname is != DBObject, then filter other classnames
		// TODO: this should not happen, verify and remove if confirmed
		if classname != "DBObject" && entity.GetMetadata("classname") != classname {
			log.Printf("SearchObjectsHandler: Skipping object ID=%s with classname=%s", entity.GetValue("id").(string), entity.GetMetadata("classname"))
			continue
		}

		// Check read permission
		if !repo.CheckReadPermission(entity) {
			log.Printf("SearchObjectsHandler: No read permission for object ID=%s", entity.GetValue("id").(string))
			continue
		}

		// If type=link, check write permission (I want only objects that I can attach to)
		if searchType == "link" && !repo.CheckWritePermission(entity) {
			log.Printf("SearchObjectsHandler: No write permission for object ID=%s (type=link)", entity.GetValue("id").(string))
			continue
		}

		resultMap := make(map[string]interface{})
		resultMap["id"] = entity.GetValue("id")
		resultMap["name"] = entity.GetValue("name")
		if desc := entity.GetValue("description"); desc != nil {
			resultMap["description"] = desc
		}

		if searchJson != "" {
			// Include all fields in searchJson mode
			for key, val := range entity.GetAllValues() {
				resultMap[key] = val
			}
		}

		resultMap["classname"] = entity.GetMetadata("classname")

		// Include mime type for DBFile objects (useful for filtering images)
		if mime := entity.GetValue("mime"); mime != nil {
			resultMap["mime"] = mime
		}

		resultList = append(resultList, resultMap)
	}

	// Apply offset and limit to resultList
	returnList := []map[string]interface{}{}
	for i := offset; len(returnList) < maxResults && i < len(resultList); i++ {
		returnList = append(returnList, resultList[i])
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ObjectsSearchResponse{
		Success: true,
		Objects: returnList,
	})
}

// DownloadFileHandler godoc
// @Summary Download file
// @Description Downloads the file content for a given DBFile object ID
// @Tags files
// @Produce octet-stream
// @Param token header string false "Temporary JWT token for access"
// @Param id path string true "File ID"
// @Param preview query string false "Set to 'yes' or 'true' to get thumbnail preview if available"
// @Success 200 {file} file "File content"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "File not found"
// @Failure 500 {object} ErrorResponse "Internal error"
// @Router /files/{id}/download [get]
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

	previewParam := r.URL.Query().Get("preview")

	// Load the DBFile object
	// tableName := "files"
	// entity := repo.GetEntityByID(tableName, fileID)
	entity := repo.FullObjectById(fileID, true)
	if entity == nil {
		log.Printf("DownloadFileHandler: Failed to load file %s", fileID)
		RespondSimpleError(w, ErrObjectNotFound, "File not found", http.StatusNotFound)
		return
	}

	// Check token is valid
	tokenString := r.URL.Query().Get("token")
	if tokenString != "" {
		// IF I have a token, I check it
		claimsDownload := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claimsDownload, func(t *jwt.Token) (interface{}, error) {
			return JWTKey, nil
		})
		// log.Printf("Claims: %+v\n", claimsDownload)
		// log.Printf("err: %v\n", err)
		if err != nil || !token.Valid {
			RespondSimpleError(w, ErrInvalidToken, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		downloadID, ok := claimsDownload["id"].(string)
		if !ok || downloadID != fileID {
			RespondSimpleError(w, ErrInvalidToken, "Token does not match file ID", http.StatusUnauthorized)
			return
		}
	} else {
		// No token: check read permissions
		if !repo.CheckReadPermission(entity) {
			RespondSimpleError(w, ErrUnauthorized, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	// Cast to DBFile to use getFullpath
	dbFile, ok := entity.(*dblayer.DBFile)
	if !ok {
		log.Printf("DownloadFileHandler: Entity is not a DBFile")
		RespondSimpleError(w, ErrInternalServer, "Invalid file entity", http.StatusInternalServerError)
		return
	}
	// log.Print("DownloadFileHandler: dbFile=", dbFile.ToJSON())

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

	filePath := dbFile.GetFullpath(nil)
	// In future, a thumbnail could be provided also for non-image files
	// e.g. PDF first page preview, video snapshot, etc.
	if previewParam == "yes" || previewParam == "true" {
		filePath = dbFile.GetThumbnailFullpath(nil)
		// log.Print("DownloadFileHandler: thumbnail filePath=", filePath)
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

	// log.Printf("DownloadFileHandler: Served file %s (%s)", filename, mime)
}

// GenerateFileTokensHandler godoc
// @Summary generates temporary JWT tokens for multiple files
// @Description Generates temporary JWT tokens for multiple files specified by their IDs
// @Tags files
// @Accept json
// @Produce json
// @Param request body map[string][]string true "List of file IDs" `{"file_ids": ["abc123", "def456"]}`
// @Success 200 {object} map[string]interface{} "Generated tokens" {"success": true, "tokens": {"abc123": "JWT_TOKEN", "def456": "JWT_TOKEN"}}
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Security BearerAuth
// @Router /files/preview-tokens [post]
func GenerateFileTokensHandler(w http.ResponseWriter, r *http.Request) {
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

	// Parse request body
	var requestData FileTokensRequest
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		log.Printf("GenerateFileTokensHandler: Failed to decode request body: %v", err)
		RespondSimpleError(w, ErrInvalidRequest, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(requestData.FileIDs) == 0 {
		RespondSimpleError(w, ErrInvalidRequest, "No file IDs provided", http.StatusBadRequest)
		return
	}

	tokens := make(map[string]string)

	// Generate token for each file (only if user has read permission)
	for _, fileID := range requestData.FileIDs {
		// Normalize file ID
		if len(fileID) == 18 {
			fileID = strings.ReplaceAll(fileID, "-", "")
		}

		// Load file and check permissions
		entity := repo.FullObjectById(fileID, true)
		if entity == nil {
			log.Printf("GenerateFileTokensHandler: File %s not found", fileID)
			continue // Skip files that don't exist
		}

		// Check read permission
		if !repo.CheckReadPermission(entity) {
			log.Printf("GenerateFileTokensHandler: User %s has no read permission for file %s", dbContext.UserID, fileID)
			continue // Skip files user can't access
		}

		// Generate JWT token for this file (valid for 15 minutes)
		expirationTime := time.Now().Add(15 * time.Minute)
		tokenClaims := jwt.MapClaims{
			"id":      fileID,
			"user_id": dbContext.UserID,
			"exp":     expirationTime.Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
		tokenString, err := token.SignedString(JWTKey)
		if err != nil {
			log.Printf("GenerateFileTokensHandler: Failed to generate token for file %s: %v", fileID, err)
			continue
		}

		tokens[fileID] = tokenString
		log.Printf("GenerateFileTokensHandler: Generated token for file %s (user: %s, expires: %s)", fileID, dbContext.UserID, expirationTime.Format(time.RFC3339))
	}

	// Return tokens
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(FileTokensResponse{
		Success: true,
		Tokens:  tokens,
	})
}

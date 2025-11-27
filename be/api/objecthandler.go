package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	// Decode update values
	var updateValues map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateValues); err != nil {
		log.Printf("UpdateObjectHandler: Failed to decode request body: %v", err)
		RespondSimpleError(w, ErrInvalidRequest, "Invalid request body", http.StatusBadRequest)
		return
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
	updated, err := repo.UpdateObject(tableName, objectID, updateValues)
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

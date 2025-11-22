package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"rprj/be/dblayer"

	"github.com/gorilla/mux"
)

// GET /content/:objectId
//
//	curl -X GET http://localhost:8080/api/content/xxxx-xxxxxxxx-xxxx \
//	  -H "Authorization: Bearer <access_token>"
func GetNavigationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	objectID := vars["objectId"]
	// The object ID is in the format xxxx-xxxxxxxx-xxxx: remove all the '-' characters
	if len(objectID) == 18 {
		objectID = strings.ReplaceAll(objectID, "-", "")
	}

	claims, err := GetClaimsFromRequest(r)
	// if err != nil {
	// 	RespondSimpleError(w, ErrUnauthorized, "Unauthorized", http.StatusUnauthorized)
	// 	return
	// }

	var dbContext dblayer.DBContext
	if err == nil {
		log.Print("GetNavigationHandler: authenticated user:", claims["user_id"])
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
	repo.Verbose = true

	obj := repo.FullObjectById(objectID, true)
	if obj == nil {
		RespondSimpleError(w, ErrObjectNotFound, "Object not found", http.StatusNotFound)
		return
	}

	// Check read permissions
	if !repo.CheckReadPermission(obj) {
		RespondSimpleError(w, ErrForbidden, "Access denied", http.StatusForbidden)
		return
	}

	if !obj.HasMetadata("classname") {
		obj.SetMetadata("classname", obj.GetTypeName())
	}

	// Returns { data: { ... } , metadata: { ... } }
	response := map[string]interface{}{
		"data":     obj.GetAllValues(),
		"metadata": obj.GetAllMetadata(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /nav/children/:folderId
//
//	curl -X GET http://localhost:8080/api/nav/children/xxxx-xxxxxxxx-xxxx \
//	  -H "Authorization: Bearer <access_token>"
func GetChildrenHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	folderId := vars["folderId"]
	// The folder ID is in the format xxxx-xxxxxxxx-xxxx: remove all the '-' characters
	if len(folderId) == 18 {
		folderId = strings.ReplaceAll(folderId, "-", "")
	}

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
	repo.Verbose = true

	children := repo.GetChildren(folderId, true)

	// Convert to response format
	childrenData := make([]map[string]interface{}, 0, len(children))
	for _, child := range children {
		if !child.HasMetadata("classname") {
			child.SetMetadata("classname", child.GetTypeName())
		}
		childrenData = append(childrenData, map[string]interface{}{
			"data":     child.GetAllValues(),
			"metadata": child.GetAllMetadata(),
		})
	}

	response := map[string]interface{}{
		"children": childrenData,
		"count":    len(childrenData),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /nav/breadcrumb/:objectId
//
//	curl -X GET http://localhost:8080/api/nav/breadcrumb/xxxx-xxxxxxxx-xxxx \
//	  -H "Authorization: Bearer <access_token>"
func GetBreadcrumbHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	objectID := vars["objectId"]
	// The object ID is in the format xxxx-xxxxxxxx-xxxx: remove all the '-' characters
	if len(objectID) == 18 {
		objectID = strings.ReplaceAll(objectID, "-", "")
	}

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
	repo.Verbose = true

	breadcrumb := repo.GetBreadcrumb(objectID)

	// Convert to response format
	breadcrumbData := make([]map[string]interface{}, 0, len(breadcrumb))
	for _, item := range breadcrumb {
		if !item.HasMetadata("classname") {
			item.SetMetadata("classname", item.GetTypeName())
		}
		breadcrumbData = append(breadcrumbData, map[string]interface{}{
			"data":     item.GetAllValues(),
			"metadata": item.GetAllMetadata(),
		})
	}

	response := map[string]interface{}{
		"breadcrumb": breadcrumbData,
		"count":      len(breadcrumbData),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

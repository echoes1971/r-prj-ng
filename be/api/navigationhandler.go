package api

import (
	"encoding/json"
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

	var dbContext dblayer.DBContext
	if err == nil {
		// log.Print("GetNavigationHandler: authenticated user:", claims["user_id"])
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

	// Check permissions
	canEdit := repo.CheckWritePermission(obj)
	obj.SetMetadata("can_edit", canEdit)

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
	repo.Verbose = false

	children := repo.GetChildren(folderId, true)

	// Convert to response format
	childrenData := make([]map[string]interface{}, 0, len(children))
	for _, child := range children {
		if (child.GetTypeName() == "DBPage" || child.GetMetadata("classname") == "DBPage") && child.GetValue("name") == "index" {
			// Skip index pages
			continue
		}
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
	repo.Verbose = false

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

// GET /nav/:objectId/indexes
//
//	curl -X GET http://localhost:8080/api/nav/xxxx-xxxxxxxx-xxxx/indexes \
func GetIndexesHandler(w http.ResponseWriter, r *http.Request) {
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
	repo.Verbose = false

	search := repo.GetInstanceByTableName("pages")
	if search == nil {
		RespondSimpleError(w, ErrInternalServer, "Failed to create page instance", http.StatusInternalServerError)
		return
	}
	search.SetValue("father_id", objectID)
	search.SetValue("name", "index")
	pages, err := repo.Search(search, false, false, "")
	if err != nil {
		RespondSimpleError(w, ErrInternalServer, "Search failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Filter by permissions
	indexes := make([]map[string]interface{}, 0, len(pages))
	for _, p := range pages {
		if repo.CheckReadPermission(p) {
			if !p.HasMetadata("classname") {
				p.SetMetadata("classname", p.GetTypeName())
			}
			indexes = append(indexes, map[string]interface{}{
				"data":     p.GetAllValues(),
				"metadata": p.GetAllMetadata(),
			})
		}
	}

	response := map[string]interface{}{
		"indexes": indexes,
		"count":   len(indexes),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /content/country/:countryId
//
//	curl -X GET http://localhost:8080/api/content/country/xxxx-xxxxxxxx-xxxx \
//	  -H "Authorization: Bearer <access_token>"
func GetCountryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	countryID := vars["countryId"]
	// The country ID is in the format xxxx-xxxxxxxx-xxxx: remove all the '-' characters
	if len(countryID) == 18 {
		countryID = strings.ReplaceAll(countryID, "-", "")
	}

	dbContext := dblayer.DBContext{
		UserID:   "-7",           // Anonymous user - countries are public
		GroupIDs: []string{"-4"}, // Guests group
		Schema:   dblayer.DbSchema,
	}

	repo := dblayer.NewDBRepository(&dbContext, dblayer.Factory, dblayer.DbConnection)
	repo.Verbose = false

	country := repo.GetEntityByID("countrylist", countryID)
	if country == nil {
		http.Error(w, "Country not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(country.GetAllValues())
}

// GetCountriesHandler returns all countries from countrylist table
// GET /api/countries
func GetCountriesHandler(w http.ResponseWriter, r *http.Request) {
	dbContext := dblayer.DBContext{
		UserID:   "-7",           // Anonymous user - countries are public
		GroupIDs: []string{"-4"}, // Guests group
		Schema:   dblayer.DbSchema,
	}

	repo := dblayer.NewDBRepository(&dbContext, dblayer.Factory, dblayer.DbConnection)
	repo.Verbose = false

	// Get all countries - create empty instance for search
	countryInstance := repo.GetInstanceByTableName("countrylist")
	if countryInstance == nil {
		http.Error(w, "Country table not found", http.StatusInternalServerError)
		return
	}

	// Search with empty criteria returns all
	countries, err := repo.Search(countryInstance, false, false, "Common_Name")
	if err != nil {
		http.Error(w, "Error fetching countries: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to simple array of country objects
	result := make([]map[string]interface{}, 0, len(countries))
	for _, country := range countries {
		result = append(result, country.GetAllValues())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"countries": result,
	})
}

// Convert to response format

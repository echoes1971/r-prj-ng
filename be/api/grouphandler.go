package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"rprj/be/dblayer"

	"github.com/gorilla/mux"
)

// GetAllGroupsHandler godoc
// @Summary Get all groups
// @Description Retrieve a list of all groups
// @Tags groups
// @Produce json
// @Param search query string false "Search term to filter groups by name"
// @Param order_by query string false "Field to order the results by"
// @Success 200 {array} map[string]interface{}
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Security BearerAuth
// @Router /groups [get]
func GetAllGroupsHandler(w http.ResponseWriter, r *http.Request) {
	searchBy := r.URL.Query().Get("search")
	orderBy := r.URL.Query().Get("order_by")

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

	search := repo.GetInstanceByTableName("groups")
	if search == nil {
		RespondSimpleError(w, ErrInternalServer, "Failed to create group instance", http.StatusInternalServerError)
		return
	}
	if searchBy != "" {
		search.SetValue("name", "%"+searchBy+"%")
		// search.SetValue("description", "%"+searchBy+"%")
	}
	groups, err := repo.Search(search, true, false, orderBy)
	if err != nil {
		RespondSimpleError(w, ErrInternalServer, "Search failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]map[string]interface{}, len(groups))
	for i, g := range groups {
		response[i] = map[string]interface{}{
			"ID":          g.GetValue("id"),
			"Name":        g.GetValue("name"),
			"Description": g.GetValue("description"),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetGroupHandler godoc
// @Summary Get group by ID
// @Description Retrieve a group by its ID
// @Tags groups
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Group Not Found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Security BearerAuth
// @Router /groups/{id} [get]
func GetGroupHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		RespondError(w, ErrMissingField, "Field is required", map[string]string{"field": "id"}, http.StatusBadRequest)
		return
	}

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

	group := repo.GetInstanceByTableName("groups")
	if group == nil {
		RespondSimpleError(w, ErrInternalServer, "Failed to create group instance", http.StatusInternalServerError)
		return
	}
	group.SetValue("id", id)
	foundGroups, err := repo.Search(group, false, false, "")
	if err != nil {
		RespondSimpleError(w, ErrInternalServer, "Failed to get group: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if len(foundGroups) == 0 {
		RespondError(w, ErrGroupNotFound, "Group not found", map[string]string{"id": id}, http.StatusNotFound)
		return
	}
	group = foundGroups[0]

	// Get User groups
	userGroupsInstance := repo.GetInstanceByTableName("users_groups")
	if userGroupsInstance == nil {
		RespondSimpleError(w, ErrInternalServer, "Failed to create user-groups instance", http.StatusInternalServerError)
		return
	}
	userGroupsInstance.SetValue("group_id", id)
	groupUsers, err := repo.Search(userGroupsInstance, false, false, "")
	if err != nil {
		RespondSimpleError(w, ErrInternalServer, "Failed to get user groups: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build response with user IDs
	userIDs := make([]string, len(groupUsers))
	for i, gu := range groupUsers {
		userIDs[i] = gu.GetValue("user_id").(string)
	}

	response := map[string]interface{}{
		"id":          group.GetValue("id"),
		"name":        group.GetValue("name"),
		"description": group.GetValue("description"),
		"user_ids":    userIDs,
	}

	json.NewEncoder(w).Encode(response)
}

// CreateGroupHandler godoc
// @Summary Create a new group
// @Description Create a new group with the provided details
// @Tags groups
// @Accept json
// @Produce json
// @Param group body object true "Group details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 409 {object} ErrorResponse "Group Already Exists"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Security BearerAuth
// @Router /groups [post]
func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		// UserIDs     []string `json:"user_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondSimpleError(w, ErrInvalidRequest, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" {
		RespondError(w, ErrMissingField, "Field is required", map[string]string{"field": "name"}, http.StatusBadRequest)
		return
	}

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

	// dblayer.InitDBConnection()
	repo := dblayer.NewDBRepository(dbContext, dblayer.Factory, dblayer.DbConnection)
	repo.Verbose = false

	dbGroup := repo.GetInstanceByTableName("groups")
	if dbGroup == nil {
		RespondSimpleError(w, ErrInternalServer, "Failed to create group instance", http.StatusInternalServerError)
		return
	}
	dbGroup.SetValue("name", req.Name)
	dbGroup.SetValue("description", req.Description)
	// dbGroup.SetMetadata("user_ids", req.UserIDs) // Users are added after a group is created

	createdGroup, err := repo.Insert(dbGroup)
	if err != nil {
		// Check if it's a duplicate name error
		if strings.Contains(err.Error(), "already exists") {
			RespondError(w, ErrGroupAlreadyExists, "Group already exists", map[string]string{"name": req.Name}, http.StatusConflict)
		} else {
			RespondSimpleError(w, ErrInternalServer, "Failed to create group: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := map[string]interface{}{
		"id":          createdGroup.GetValue("id"),
		"name":        createdGroup.GetValue("name"),
		"description": createdGroup.GetValue("description"),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateGroupHandler godoc
// @Summary Update group by ID
// @Description Update a group's details by its ID
// @Tags groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param group body object true "Group details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Group Not Found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Security BearerAuth
// @Router /groups/{id} [put]
func UpdateGroupHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		RespondError(w, ErrMissingField, "Field is required", map[string]string{"field": "id"}, http.StatusBadRequest)
		return
	}

	var req struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		UserIDs     []string `json:"user_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondSimpleError(w, ErrInvalidRequest, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" {
		RespondError(w, ErrMissingField, "Field is required", map[string]string{"field": "name"}, http.StatusBadRequest)
		return
	}

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

	group := repo.GetInstanceByTableName("groups")
	if group == nil {
		RespondSimpleError(w, ErrInternalServer, "Failed to create group instance", http.StatusInternalServerError)
		return
	}
	group.SetValue("id", id)
	group.SetValue("name", req.Name)
	group.SetValue("description", req.Description)
	group.SetMetadata("user_ids", req.UserIDs) // Users are updated after a group is updated

	g, err := repo.Update(group)
	if err != nil {
		RespondSimpleError(w, ErrInternalServer, "Failed to update group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ID":          g.GetValue("id"),
		"Name":        g.GetValue("name"),
		"Description": g.GetValue("description"),
	})
}

// DeleteGroupHandler godoc
// @Summary Delete group by ID
// @Description Delete a group by its ID
// @Tags groups
// @Param id path string true "Group ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Security BearerAuth
// @Router /groups/{id} [delete]
func DeleteGroupHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		RespondError(w, ErrMissingField, "Field is required", map[string]string{"field": "id"}, http.StatusBadRequest)
		return
	}

	// Groups with negative ID cannot be deleted (system groups)
	if strings.HasPrefix(id, "-") {
		RespondError(w, ErrForbidden, "Cannot delete system groups", map[string]string{"id": id}, http.StatusForbidden)
		return
	}

	claims, err := GetClaimsFromRequest(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	dbContext := &dblayer.DBContext{
		UserID:   claims["user_id"],
		GroupIDs: strings.Split(claims["groups"], ","),
		Schema:   dblayer.DbSchema,
	}

	repo := dblayer.NewDBRepository(dbContext, dblayer.Factory, dblayer.DbConnection)
	repo.Verbose = false

	group := repo.GetInstanceByTableName("groups")
	if group == nil {
		RespondSimpleError(w, ErrInternalServer, "Failed to create group instance", http.StatusInternalServerError)
		return
	}
	group.SetValue("id", id)

	_, err = repo.Delete(group)
	if err != nil {
		RespondSimpleError(w, ErrInternalServer, "Failed to delete group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

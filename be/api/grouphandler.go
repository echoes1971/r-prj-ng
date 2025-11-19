package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"rprj/be/dblayer"

	"github.com/gorilla/mux"
)

// GET /groups
func GetAllGroupsHandler(w http.ResponseWriter, r *http.Request) {
	searchBy := r.URL.Query().Get("search")
	orderBy := r.URL.Query().Get("order_by")

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
	repo.Verbose = true

	search := repo.GetInstanceByTableName("groups")
	if search == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create user instance"})
		return
	}
	if searchBy != "" {
		search.SetValue("name", "%"+searchBy+"%")
		// search.SetValue("description", "%"+searchBy+"%")
	}
	groups, err := repo.Search(search, true, false, orderBy)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Search failed: " + err.Error()})
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

// GET /groups/{id}
func GetGroupHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	claims, err := GetClaimsFromRequest(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	dbContext := &dblayer.DBContext{
		UserID:   claims["user_id"],
		GroupIDs: strings.Split(claims["groups"], ","),
		Schema:   dblayer.DbSchema,
	}

	repo := dblayer.NewDBRepository(dbContext, dblayer.Factory, dblayer.DbConnection)
	repo.Verbose = true

	group := repo.GetInstanceByTableName("groups")
	if group == nil {
		http.Error(w, "failed to create group instance", http.StatusInternalServerError)
		return
	}
	group.SetValue("id", id)
	foundGroups, err := repo.Search(group, false, false, "")
	if err != nil {
		http.Error(w, "failed to get group: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if len(foundGroups) == 0 {
		http.NotFound(w, r)
		return
	}
	group = foundGroups[0]

	// Get User groups
	userGroupsInstance := repo.GetInstanceByTableName("users_groups")
	if userGroupsInstance == nil {
		http.Error(w, "failed to create user-groups instance", http.StatusInternalServerError)
		return
	}
	userGroupsInstance.SetValue("group_id", id)
	groupUsers, err := repo.Search(userGroupsInstance, false, false, "")
	if err != nil {
		http.Error(w, "failed to get user groups: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build response with user IDs
	userIDs := make([]string, len(groupUsers))
	for i, gu := range groupUsers {
		userIDs[i] = gu.GetValue("user_id")
	}

	response := map[string]interface{}{
		"id":          group.GetValue("id"),
		"name":        group.GetValue("name"),
		"description": group.GetValue("description"),
		"user_ids":    userIDs,
	}

	json.NewEncoder(w).Encode(response)
}

// POST /groups
func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		// UserIDs     []string `json:"user_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request format"})
		return
	}

	// Validate required fields
	if req.Name == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Group name is required"})
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

	// dblayer.InitDBConnection()
	repo := dblayer.NewDBRepository(dbContext, dblayer.Factory, dblayer.DbConnection)
	repo.Verbose = true

	dbGroup := repo.GetInstanceByTableName("groups")
	if dbGroup == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create group instance"})
		return
	}
	dbGroup.SetValue("name", req.Name)
	dbGroup.SetValue("description", req.Description)
	// dbGroup.SetMetadata("user_ids", req.UserIDs) // Users are added after a group is created

	createdGroup, err := repo.Insert(dbGroup)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		// Check if it's a duplicate name error
		if strings.Contains(err.Error(), "already exists") {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create group: " + err.Error()})
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

// PUT /groups/{id}
func UpdateGroupHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing group ID"})
		return
	}

	var req struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		UserIDs     []string `json:"user_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request format"})
		return
	}

	// Validate required fields
	if req.Name == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Group name is required"})
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
	repo.Verbose = true

	group := repo.GetInstanceByTableName("groups")
	if group == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create group instance"})
		return
	}
	group.SetValue("id", id)
	group.SetValue("name", req.Name)
	group.SetValue("description", req.Description)
	group.SetMetadata("user_ids", req.UserIDs) // Users are updated after a group is updated

	g, err := repo.Update(group)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update group: " + err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ID":          g.GetValue("id"),
		"Name":        g.GetValue("name"),
		"Description": g.GetValue("description"),
	})
}

// DELETE /groups/{id}
func DeleteGroupHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing group ID"})
		return
	}

	// Groups with negative ID cannot be deleted (system groups)
	if strings.HasPrefix(id, "-") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cannot delete system groups"})
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
	repo.Verbose = true

	group := repo.GetInstanceByTableName("groups")
	if group == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create group instance"})
		return
	}
	group.SetValue("id", id)

	_, err = repo.Delete(group)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete group: " + err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

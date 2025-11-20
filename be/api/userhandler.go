package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"rprj/be/dblayer"

	"github.com/gorilla/mux"
)

// GET /users
func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
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

	search := repo.GetInstanceByTableName("users")
	if search == nil {
		RespondSimpleError(w, ErrInternalServer, "Failed to create user instance", http.StatusInternalServerError)
		return
	}
	if searchBy != "" {
		search.SetValue("login", "%"+searchBy+"%")
		// search.SetValue("fullname", "%"+searchBy+"%")
	}
	users, err := repo.Search(search, true, false, orderBy)
	if err != nil {
		RespondSimpleError(w, ErrInternalServer, "Failed to search users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]map[string]interface{}, len(users))
	for i, u := range users {
		response[i] = map[string]interface{}{
			"ID":       u.GetValue("id"),
			"Login":    u.GetValue("login"),
			"Fullname": u.GetValue("fullname"),
			"GroupID":  u.GetValue("group_id"),
		}
	}
	json.NewEncoder(w).Encode(response)
}

// GET /users/{id}
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
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
	repo.Verbose = true

	user := repo.GetInstanceByTableName("users")
	if user == nil {
		RespondSimpleError(w, ErrInternalServer, "Failed to create user instance", http.StatusInternalServerError)
		return
	}
	user.SetValue("id", id)
	foundUsers, err := repo.Search(user, false, false, "")
	if err != nil {
		RespondSimpleError(w, ErrInternalServer, "Failed to get user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if len(foundUsers) == 0 {
		RespondError(w, ErrUserNotFound, "User not found", map[string]string{"id": id}, http.StatusNotFound)
		return
	}
	user = foundUsers[0]

	// Get User groups
	userGroupsInstance := repo.GetInstanceByTableName("users_groups")
	if userGroupsInstance == nil {
		RespondSimpleError(w, ErrInternalServer, "Failed to create user-groups instance", http.StatusInternalServerError)
		return
	}
	userGroupsInstance.SetValue("user_id", id)
	userGroups, err := repo.Search(userGroupsInstance, false, false, "")
	if err != nil {
		RespondSimpleError(w, ErrInternalServer, "Failed to get user groups: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build response with group IDs
	groupIDs := make([]string, len(userGroups))
	for i, ug := range userGroups {
		groupIDs[i] = ug.GetValue("group_id").(string)
	}

	response := map[string]interface{}{
		"id":        user.GetValue("id"),
		"login":     user.GetValue("login"),
		"fullname":  user.GetValue("fullname"),
		"group_id":  user.GetValue("group_id"),
		"group_ids": groupIDs,
	}

	json.NewEncoder(w).Encode(response)
}

// POST /users
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login    string   `json:"login"`
		Pwd      string   `json:"pwd"`
		Fullname string   `json:"fullname"`
		GroupID  string   `json:"group_id"`
		GroupIDs []string `json:"group_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondSimpleError(w, ErrInvalidRequest, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Login == "" {
		RespondError(w, ErrMissingField, "Field is required", map[string]string{"field": "login"}, http.StatusBadRequest)
		return
	}
	if req.Pwd == "" {
		RespondError(w, ErrMissingField, "Field is required", map[string]string{"field": "password"}, http.StatusBadRequest)
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
	repo.Verbose = true

	dbUser := repo.GetInstanceByTableName("users")
	if dbUser == nil {
		RespondSimpleError(w, ErrInternalServer, "Failed to create user instance", http.StatusInternalServerError)
		return
	}
	dbUser.SetValue("login", req.Login)
	dbUser.SetValue("pwd", req.Pwd)
	dbUser.SetValue("fullname", req.Fullname)
	dbUser.SetMetadata("group_ids", req.GroupIDs)

	createdUser, err := repo.Insert(dbUser)
	if err != nil {
		// Check if it's a duplicate login error
		if strings.Contains(err.Error(), "already exists") {
			RespondError(w, ErrUserAlreadyExists, "User already exists", map[string]string{"login": req.Login}, http.StatusConflict)
		} else {
			RespondSimpleError(w, ErrInternalServer, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := map[string]interface{}{
		"id":        createdUser.GetValue("id"),
		"login":     createdUser.GetValue("login"),
		"fullname":  createdUser.GetValue("fullname"),
		"group_id":  createdUser.GetValue("group_id"),
		"group_ids": req.GroupIDs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// PUT /users/{id}
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		RespondError(w, ErrMissingField, "Field is required", map[string]string{"field": "id"}, http.StatusBadRequest)
		return
	}

	var req struct {
		Login    string   `json:"login"`
		Pwd      string   `json:"pwd"`
		Fullname string   `json:"fullname"`
		GroupID  string   `json:"group_id"`
		GroupIDs []string `json:"group_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondSimpleError(w, ErrInvalidRequest, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Login == "" {
		RespondError(w, ErrMissingField, "Field is required", map[string]string{"field": "login"}, http.StatusBadRequest)
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
	repo.Verbose = true

	user := repo.GetInstanceByTableName("users")
	if user == nil {
		RespondSimpleError(w, ErrInternalServer, "Failed to create user instance", http.StatusInternalServerError)
		return
	}

	user.SetValue("id", id)
	// user.SetValue("login", req.Login) // login CANNOT be changed
	user.SetValue("fullname", req.Fullname)
	if req.Pwd != "" {
		user.SetValue("pwd", req.Pwd)
	}
	user.SetValue("group_id", req.GroupID)
	user.SetMetadata("group_ids", req.GroupIDs)

	u, err := repo.Update(user)
	if err != nil {
		// Check if it's a duplicate login error
		if strings.Contains(err.Error(), "already exists") {
			RespondError(w, ErrUserAlreadyExists, "User already exists", map[string]string{"login": req.Login}, http.StatusConflict)
		} else {
			RespondSimpleError(w, ErrInternalServer, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

// DELETE /users/{id}
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
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
	repo.Verbose = true

	user := repo.GetInstanceByTableName("users")
	if user == nil {
		RespondSimpleError(w, ErrInternalServer, "Failed to create user instance", http.StatusInternalServerError)
		return
	}
	user.SetValue("id", id)

	_, err = repo.Delete(user)
	if err != nil {
		RespondSimpleError(w, ErrInternalServer, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

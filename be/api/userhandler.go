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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create user instance"})
		return
	}
	if searchBy != "" {
		search.SetValue("login", "%"+searchBy+"%")
		// search.SetValue("fullname", "%"+searchBy+"%")
	}
	users, err := repo.Search(search, true, false, orderBy)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to search users: " + err.Error()})
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

	user := repo.GetInstanceByTableName("users")
	if user == nil {
		http.Error(w, "failed to create user instance", http.StatusInternalServerError)
		return
	}
	user.SetValue("id", id)
	foundUsers, err := repo.Search(user, false, false, "")
	if err != nil {
		http.Error(w, "failed to get user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if len(foundUsers) == 0 {
		http.NotFound(w, r)
		return
	}
	user = foundUsers[0]

	// Get User groups
	userGroupsInstance := repo.GetInstanceByTableName("users_groups")
	if userGroupsInstance == nil {
		http.Error(w, "failed to create user-groups instance", http.StatusInternalServerError)
		return
	}
	userGroupsInstance.SetValue("user_id", id)
	userGroups, err := repo.Search(userGroupsInstance, false, false, "")
	if err != nil {
		http.Error(w, "failed to get user groups: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build response with group IDs
	groupIDs := make([]string, len(userGroups))
	for i, ug := range userGroups {
		groupIDs[i] = ug.GetValue("group_id")
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request format"})
		return
	}

	// Validate required fields
	if req.Login == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Login is required"})
		return
	}
	if req.Pwd == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Password is required"})
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

	dbUser := repo.GetInstanceByTableName("users")
	if dbUser == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create user instance"})
		return
	}
	dbUser.SetValue("login", req.Login)
	dbUser.SetValue("pwd", req.Pwd)
	dbUser.SetValue("fullname", req.Fullname)

	createdUser, err := repo.Insert(dbUser)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		// Check if it's a duplicate login error
		if strings.Contains(err.Error(), "already exists") {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create user: " + err.Error()})
		}
		return
	}

	// Assign the user to the specified groups
	for _, gID := range req.GroupIDs {
		dbUserGroup := repo.GetInstanceByTableName("users_groups")
		if dbUserGroup == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create user-group association instance"})
			return
		}
		dbUserGroup.SetValue("user_id", createdUser.GetValue("id"))
		dbUserGroup.SetValue("group_id", gID)
		_, err := repo.Insert(dbUserGroup)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to assign user to group: " + err.Error()})
			return
		}
	}

	response := map[string]interface{}{
		"id":        createdUser.GetValue("id"),
		"login":     createdUser.GetValue("login"),
		"fullname":  createdUser.GetValue("fullname"),
		"group_id":  createdUser.GetValue("group_id"),
		"group_ids": req.GroupIDs,
	}

	// u := models.DBUser{
	// 	Login:    req.Login,
	// 	Pwd:      req.Pwd,
	// 	Fullname: req.Fullname,
	// }

	// // Create user with transaction (creates group, user, and associations atomically)
	// createdUser, _, err := db.CreateUser(u, req.Login, req.GroupIDs)
	// if err != nil {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	// Check if it's a duplicate login error
	// 	if strings.Contains(err.Error(), "already exists") {
	// 		w.WriteHeader(http.StatusConflict)
	// 		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	// 	} else {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create user: " + err.Error()})
	// 	}
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// PUT /users/{id}
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing user ID"})
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request format"})
		return
	}

	// Validate required fields
	if req.Login == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Login is required"})
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

	user := repo.GetInstanceByTableName("users")
	if user == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create user instance"})
		return
	}

	user.SetValue("id", id)
	user.SetValue("login", req.Login)
	user.SetValue("fullname", req.Fullname)
	if req.Pwd != "" {
		user.SetValue("pwd", req.Pwd)
	}
	user.SetValue("group_id", req.GroupID)
	u, err := repo.Update(user)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		// Check if it's a duplicate login error
		if strings.Contains(err.Error(), "already exists") {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update user: " + err.Error()})
		}
		return
	}

	// Assign the user to the specified groups
	// First, delete existing associations
	userGroupsInstance := repo.GetInstanceByTableName("users_groups")
	if userGroupsInstance == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create user-groups instance"})
		return
	}
	userGroupsInstance.SetValue("user_id", id)
	existingUserGroups, err := repo.Search(userGroupsInstance, false, false, "")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to get existing user groups: " + err.Error()})
		return
	}
	for _, ug := range existingUserGroups {
		_, err := repo.Delete(ug)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete user-group association: " + err.Error()})
			return
		}
	}
	// Then, add new associations
	for _, gID := range req.GroupIDs {
		dbUserGroup := repo.GetInstanceByTableName("users_groups")
		if dbUserGroup == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create user-group association instance"})
			return
		}
		dbUserGroup.SetValue("user_id", id)
		dbUserGroup.SetValue("group_id", gID)
		_, err := repo.Insert(dbUserGroup)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to assign user to group: " + err.Error()})
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

// DELETE /users/{id}
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing user ID"})
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

	user := repo.GetInstanceByTableName("users")
	if user == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create user instance"})
		return
	}
	user.SetValue("id", id)

	_, err = repo.Delete(user)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete user: " + err.Error()})
		return
	}

	// if err := db.DeleteUser(id); err != nil {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete user: " + err.Error()})
	// 	return
	// }

	w.WriteHeader(http.StatusNoContent)
}

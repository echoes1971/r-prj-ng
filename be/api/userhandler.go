package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"rprj/be/dblayer"

	"github.com/gorilla/mux"
)

// GetAllUsersHandler godoc
//
//	@Summary gets all users
//	@Description Retrieves a list of all users, with optional search and ordering
//	@Tags users
//	@Produce json
//	@Param search query string false "Search term to filter users by login or fullname"
//	@Param order_by query string false "Field to order the results by"
//	@Success 200 {array} map[string]interface{} "List of users"
//	@Failure 401 {object} ErrorResponse "Unauthorized"
//	@Failure 403 {object} ErrorResponse "Forbidden"
//	@Failure 500 {object} ErrorResponse "Internal server error"
//	@Security BearerAuth
//	@Router /users [get]
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
	repo.Verbose = false

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

// GetUserHandler godoc
//
//	@Summary gets a user by ID
//	@Description Retrieves the user specified by the given ID
//	@Tags users
//	@Param id path string true "User ID"
//	@Produce json
//	@Success 200 {object} map[string]interface{} "User data"
//	@Failure 400 {object} ErrorResponse "Invalid request"
//	@Failure 401 {object} ErrorResponse "Unauthorized"
//	@Failure 403 {object} ErrorResponse "Forbidden"
//	@Failure 404 {object} ErrorResponse "User not found"
//	@Failure 500 {object} ErrorResponse "Internal server error"
//	@Security BearerAuth
//	@Router /users/{id} [get]
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
	repo.Verbose = false

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

// CreateUserHandler godoc
//
//	@Summary creates a new user
//	@Description Creates a new user with the provided details
//	@Tags users
//	@Accept json
//	@Produce json
//	@Param request body map[string]interface{} true "Request body containing user details"
//	@Success 201 {object} map[string]interface{} "Created user data"
//	@Failure 400 {object} ErrorResponse "Invalid request"
//	@Failure 401 {object} ErrorResponse "Unauthorized"
//	@Failure 403 {object} ErrorResponse "Forbidden"
//	@Failure 409 {object} ErrorResponse "User already exists"
//	@Failure 500 {object} ErrorResponse "Internal server error"
//	@Security BearerAuth
//	@Router /users [post]
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
	repo.Verbose = false

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

// UpdateUserHandler godoc
//
//	@Summary updates a user by ID
//	@Description Updates the user specified by the given ID
//	@Tags users
//	@Param id path string true "User ID"
//	@Param request body map[string]interface{} true "Request body containing user fields to update"
//	@Produce json
//	@Success 200 {object} map[string]interface{} "Updated user data"
//	@Failure 400 {object} ErrorResponse "Invalid request"
//	@Failure 401 {object} ErrorResponse "Unauthorized"
//	@Failure 403 {object} ErrorResponse "Forbidden"
//	@Failure 404 {object} ErrorResponse "User not found"
//	@Failure 500 {object} ErrorResponse "Internal server error"
//	@Security BearerAuth
//	@Router /users/{id} [put]
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
	repo.Verbose = false

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
	if req.GroupID != "" {
		user.SetValue("group_id", req.GroupID)
	}
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

// DeleteUserHandler godoc
//
//	@Summary deletes a user by ID
//	@Description Deletes the user specified by the given ID
//	@Tags users
//	@Param id path string true "User ID"
//	@Success 204 "No Content"
//	@Failure 400 {object} ErrorResponse "Invalid request"
//	@Failure 401 {object} ErrorResponse "Unauthorized"
//	@Failure 403 {object} ErrorResponse "Forbidden"
//	@Failure 500 {object} ErrorResponse "Internal server error"
//	@Security BearerAuth
//	@Router /users/{id} [delete]
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
	repo.Verbose = false

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

// GetUserPersonHandler godoc
//
//	@Summary gets or creates a Person record linked to the user
//	@Description Retrieves the Person record associated with the specified user ID, or creates one if it doesn't exist
//	@Tags users
//	@Produce json
//	@Param userId path string true "User ID"
//	@Success 200 {object} map[string]interface{} "Person data"
//	@Failure 401 {object} ErrorResponse "Unauthorized"
//	@Failure 403 {object} ErrorResponse "Forbidden"
//	@Failure 404 {object} ErrorResponse "User not found"
//	@Failure 500 {object} ErrorResponse "Internal server error"
//	@Security BearerAuth
//	@Router /users/{userId}/person [get]
func GetUserPersonHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

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

	// Search for existing person with fk_users_id = userId
	person := repo.GetInstanceByTableName("people")
	if person == nil {
		RespondSimpleError(w, ErrInternalServer, "Failed to create person instance", http.StatusInternalServerError)
		return
	}

	if userId == "" {
		userId = claims["user_id"]
		log.Print("UserId not provided, using claim user_id: ", userId)
	}
	person.SetValue("fk_users_id", userId)
	people, err := repo.Search(person, false, false, "")
	if err != nil {
		RespondSimpleError(w, ErrInternalServer, "Failed to search person: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// If person exists, return it
	if len(people) > 0 {
		existingPerson := people[0]

		// Check read permission
		if !repo.CheckReadPermission(existingPerson) {
			RespondSimpleError(w, ErrForbidden, "Permission denied", http.StatusForbidden)
			return
		}

		response := map[string]interface{}{
			"person_id": existingPerson.GetValue("id"),
			"user_id":   userId,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Person doesn't exist, create it
	// Get user details to populate person name
	user := repo.GetInstanceByTableName("users")
	user.SetValue("id", userId)
	users, err := repo.Search(user, false, false, "")
	if err != nil || len(users) == 0 {
		RespondSimpleError(w, ErrUserNotFound, "User not found", http.StatusNotFound)
		return
	}

	currentUser := users[0]
	fullname := currentUser.GetValue("fullname")
	if fullname == nil {
		fullname = currentUser.GetValue("login")
	}

	// Create new person
	newPerson := repo.GetInstanceByTableName("people")
	newPerson.SetValue("name", fullname)
	newPerson.SetValue("fk_users_id", userId)
	newPerson.SetValue("owner", claims["user_id"])
	newPerson.SetValue("group_id", strings.Split(claims["groups"], ",")[0]) // First group
	newPerson.SetValue("permissions", "rwxr-----")                          // Owner and group can read

	createdPerson, err := repo.Insert(newPerson)
	if err != nil {
		RespondSimpleError(w, ErrInternalServer, "Failed to create person: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"person_id": createdPerson.GetValue("id"),
		"user_id":   userId,
		"created":   true,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

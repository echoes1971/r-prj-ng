package api

import (
	"encoding/json"
	"net/http"

	"rprj/be/db"
	"rprj/be/models"

	"github.com/gorilla/mux"
)

// GET /users
func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	searchBy := r.URL.Query().Get("search")
	orderBy := r.URL.Query().Get("order_by")
	users, err := db.GetAllUsers(searchBy, orderBy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

// GET /users/{id}
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.NotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// POST /users
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var u models.DBUser
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.CreateUser(u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

// PUT /users/{id}
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var u models.DBUser
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.UpdateUser(u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(u)
}

// DELETE /users/{id}
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	if err := db.DeleteUser(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

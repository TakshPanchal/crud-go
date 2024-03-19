package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/takshpanchal/crud-go/models"
)

type UserHandler struct {
	InfoLogger *log.Logger
	ErrLogger  *log.Logger
	Model      *models.UserModel
}

// CreateUser is handler for Endpoint to create a new user
func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// only accept POST request
	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	uh.InfoLogger.Printf("User %+v", u)
	id, err := uh.Model.Insert(u)
	if err != nil {
		uh.ErrLogger.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data := map[string]int{"id": id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}

// GetUser is handler for Endpoint to get user by id
func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// only accept GET user

}

func (uh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {}

/*
Endpoint to get user by id
Endpoint to delete a user by id
Endpoint to create a new user
Endpoint to update an existing user by id.
*/

package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

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
	id, err := uh.Model.Insert(&u)
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

func (u *UserHandler) User(w http.ResponseWriter, r *http.Request) {
	uid, err := strconv.Atoi(strings.Split(r.URL.Path, "/user/")[1])
	if err != nil {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		u.GetUser(w, r, uid)
	case http.MethodDelete:
		u.DeleteUser(w, r, uid)
	case http.MethodPatch:
		u.UpdateUserFields(w, r, uid)
	default:
		//TODO: send Method not allowed response
	}
}

// GetUser is handler for Endpoint to get user by id
func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request, uid int) {
	// only accept GET user
	user, err := uh.Model.Get(uid)

	if err != nil {
		uh.ErrLogger.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request, uid int) {
	// only accept DELTE user
	err := h.Model.Delete(uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) UpdateUserFields(w http.ResponseWriter, r *http.Request, uid int) {
	// only accept PATCH user

	u := models.User{ID: uid}
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	h.InfoLogger.Printf("User %+v", u)
	err = h.Model.UpdateFields(uid, &u)

	if err != nil {
		h.ErrLogger.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	h.GetUser(w, r, uid)
}

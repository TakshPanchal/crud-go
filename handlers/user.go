package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/lib/pq"
	"github.com/takshpanchal/crud-go/models"
)

type UserHandler struct {
	iLog  *log.Logger
	eLog  *log.Logger
	model *models.UserModel
}

func NewUserHandlers(infoLogger, errLogger *log.Logger, userModel *models.UserModel) *UserHandler {
	return &UserHandler{infoLogger, errLogger, userModel}
}

// CreateUser is handler for Endpoint to create a new user
func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// only accept POST request
	if r.Method != http.MethodPost {
		methodNotSupported(w, []string{http.MethodPost})
		return
	}

	var u models.User
	err := decodeJSONBody(w, r, u)
	if err != nil {
		var m *malformedRequest
		if errors.As(err, m) {
			clientError(w, m.msg, m.status)

		} else {
			serverError(w, uh.eLog, err)
		}
		return
	}

	id, err := uh.model.Insert(&u)
	if err != nil {
		uh.eLog.Printf("%s \n %T", err.Error(), err)
		// serverError(w, uh.ErrLogger, err)
		var pErr *pq.Error
		if errors.As(err, &pErr) {
			switch pErr.Code.Name() {
			case "unique_violation", "foreign_key_violation":
				clientError(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			}
		}
		return
	}

	data := map[string]int{"id": id}
	sendJSON(w, data)
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
		methodNotSupported(w, []string{http.MethodGet, http.MethodPatch, http.MethodDelete})
	}
}

// GetUser is handler for Endpoint to get user by id
func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request, uid int) {
	user, err := uh.model.Get(uid)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
		} else {
			serverError(w, uh.eLog, err)
		}
		return
	}

	sendJSON(w, user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request, uid int) {
	err := h.model.Delete(uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
		} else {
			serverError(w, h.eLog, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) UpdateUserFields(w http.ResponseWriter, r *http.Request, uid int) {
	u := models.User{ID: uid}
	err := decodeJSONBody(w, r, u)
	if err != nil {
		var m *malformedRequest
		if errors.As(err, m) {
			clientError(w, m.msg, m.status)

		} else {
			serverError(w, h.eLog, err)
		}
		return
	}

	err = h.model.UpdateFields(uid, &u)

	if err != nil {
		serverError(w, h.eLog, err)
		return
	}
	h.GetUser(w, r, uid)
}

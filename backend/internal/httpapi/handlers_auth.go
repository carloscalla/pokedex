package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"pokedex-backend/internal/auth"
)

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authUserDTO struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type authResponse struct {
	Token string      `json:"token"`
	User  authUserDTO `json:"user"`
}

func (h *handlers) signup(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, user, err := h.deps.AuthService.Signup(r.Context(), req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrWeakPassword):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, auth.ErrUsernameTaken):
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "failed to sign up")
		}
		return
	}

	writeJSON(w, http.StatusCreated, authResponse{Token: token, User: authUserDTO{ID: user.ID, Username: user.Username}})
}

func (h *handlers) login(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, user, err := h.deps.AuthService.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to log in")
		return
	}

	writeJSON(w, http.StatusOK, authResponse{Token: token, User: authUserDTO{ID: user.ID, Username: user.Username}})
}

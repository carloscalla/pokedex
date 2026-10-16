package httpapi

import (
	"errors"
	"net/http"
	"strconv"

	"pokedex-backend/internal/auth"
	"pokedex-backend/internal/team"
)

func (h *handlers) getTeam(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())
	resp, err := h.deps.TeamService.List(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch team")
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *handlers) addToTeam(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())
	pokemonID, err := strconv.Atoi(r.PathValue("pokemonId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid pokemon id")
		return
	}

	resp, err := h.deps.TeamService.Add(r.Context(), userID, pokemonID)
	if err != nil {
		switch {
		case errors.Is(err, team.ErrPokemonNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, team.ErrTeamFull):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, team.ErrAlreadyOnTeam):
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "failed to add to team")
		}
		return
	}
	writeJSON(w, http.StatusCreated, resp)
}

func (h *handlers) removeFromTeam(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())
	pokemonID, err := strconv.Atoi(r.PathValue("pokemonId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid pokemon id")
		return
	}

	resp, err := h.deps.TeamService.Remove(r.Context(), userID, pokemonID)
	if err != nil {
		if errors.Is(err, team.ErrNotOnTeam) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to remove from team")
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

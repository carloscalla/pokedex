package httpapi

import (
	"errors"
	"net/http"

	"pokedex-backend/internal/catalog"
)

func (h *handlers) listPokemon(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	resp, err := h.deps.CatalogService.List(r.Context(), catalog.ListQuery{
		Page:     q.Get("page"),
		PageSize: q.Get("pageSize"),
		Search:   q.Get("search"),
		Type:     q.Get("type"),
		Sort:     q.Get("sort"),
	})
	if err != nil {
		switch {
		case errors.Is(err, catalog.ErrInvalidType), errors.Is(err, catalog.ErrInvalidSort):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "failed to list pokemon")
		}
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *handlers) getPokemon(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	detail, err := h.deps.CatalogService.GetDetail(r.Context(), name)
	if err != nil {
		if errors.Is(err, catalog.ErrNotFound) {
			writeError(w, http.StatusNotFound, "pokemon not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to fetch pokemon")
		return
	}
	writeJSON(w, http.StatusOK, detail)
}

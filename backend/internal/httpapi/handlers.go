package httpapi

import (
	"net/http"

	"pokedex-backend/internal/auth"
	"pokedex-backend/internal/catalog"
	"pokedex-backend/internal/team"
)

// Deps are the service-layer dependencies the HTTP layer talks to. Handlers
// never touch the database directly.
type Deps struct {
	CatalogService *catalog.Service
	AuthService    *auth.Service
	TeamService    *team.Service
	JWTService     *auth.JWTService
	AllowedOrigin  string
	AdminToken     string
	// SyncTrigger runs a full catalog re-sync (blocking); the admin handler
	// runs it in its own goroutine so the request can return immediately.
	SyncTrigger func()
}

type handlers struct {
	deps Deps
}

func (h *handlers) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

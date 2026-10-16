package httpapi

import "net/http"

func NewRouter(d Deps) http.Handler {
	mux := http.NewServeMux()
	h := &handlers{deps: d}

	mux.HandleFunc("GET /api/health", h.health)

	mux.HandleFunc("GET /api/pokemon", h.listPokemon)
	mux.HandleFunc("GET /api/pokemon/{name}", h.getPokemon)

	mux.HandleFunc("POST /api/auth/signup", h.signup)
	mux.HandleFunc("POST /api/auth/login", h.login)

	requireAuth := RequireAuth(d.JWTService)
	mux.Handle("GET /api/team", requireAuth(http.HandlerFunc(h.getTeam)))
	mux.Handle("POST /api/team/{pokemonId}", requireAuth(http.HandlerFunc(h.addToTeam)))
	mux.Handle("DELETE /api/team/{pokemonId}", requireAuth(http.HandlerFunc(h.removeFromTeam)))

	requireAdmin := RequireAdminToken(d.AdminToken)
	mux.Handle("POST /api/admin/sync", requireAdmin(http.HandlerFunc(h.triggerSync)))

	return CORS(d.AllowedOrigin)(mux)
}

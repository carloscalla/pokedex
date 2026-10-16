package httpapi

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"pokedex-backend/internal/auth"
)

// CORS allows the configured frontend origin. In day-to-day dev, Vite's
// proxy makes requests same-origin so this mostly sits inert; it matters
// once frontend and backend are deployed to different origins.
func CORS(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAuth collapses every failure mode (missing header, malformed,
// bad signature, expired) into the same 401 response so it never leaks
// which check failed.
func RequireAuth(jwtSvc *auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				writeError(w, http.StatusUnauthorized, "missing or invalid authorization token")
				return
			}
			token := strings.TrimPrefix(header, "Bearer ")

			userID, err := jwtSvc.Verify(token)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "missing or invalid authorization token")
				return
			}

			ctx := auth.ContextWithUserID(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAdminToken guards operator-only routes with a static shared
// secret instead of the JWT bearer scheme, deliberately keeping "operator"
// separate from "logged-in user".
func RequireAdminToken(adminToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			provided := r.Header.Get("X-Admin-Token")
			if adminToken == "" || subtle.ConstantTimeCompare([]byte(provided), []byte(adminToken)) != 1 {
				writeError(w, http.StatusUnauthorized, "missing or invalid admin token")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

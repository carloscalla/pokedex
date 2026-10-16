package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"pokedex-backend/internal/auth"
)

func TestRequireAuth(t *testing.T) {
	jwtSvc := auth.NewJWTService("test-secret", time.Hour)
	validToken, err := jwtSvc.Issue(7)
	if err != nil {
		t.Fatalf("Issue returned error: %v", err)
	}

	var capturedUserID int64
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUserID, _ = auth.UserIDFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	})
	handler := RequireAuth(jwtSvc)(next)

	cases := []struct {
		name       string
		authHeader string
		wantStatus int
	}{
		{"missing header", "", http.StatusUnauthorized},
		{"malformed header", "Bearer garbage", http.StatusUnauthorized},
		{"valid token", "Bearer " + validToken, http.StatusOK},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/team", nil)
			if c.authHeader != "" {
				req.Header.Set("Authorization", c.authHeader)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != c.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, c.wantStatus)
			}
		})
	}

	if capturedUserID != 7 {
		t.Errorf("captured userID = %d, want 7 (from the last, valid-token request)", capturedUserID)
	}
}

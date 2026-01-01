package api

import (
	"net/http"
	"os"
	"strings"
)

// RequireAdminToken vérifie que le token admin est présent dans les headers
func RequireAdminToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		adminToken := os.Getenv("ADMIN_TOKEN")
		if adminToken == "" {
			http.Error(w, "Admin token not configured", http.StatusInternalServerError)
			return
		}

		// Vérifier le token dans le header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Format: "Bearer <token>" ou juste "<token>"
		token := strings.TrimPrefix(authHeader, "Bearer ")
		token = strings.TrimSpace(token)

		if token != adminToken {
			http.Error(w, "Invalid admin token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"strings"

	"upcycle_connect-api/internal/config"
	"upcycle_connect-api/internal/utils"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type contextKey string

const (
	ContextUserID      contextKey = "user_id"
	ContextRoles       contextKey = "roles"
	ContextPermissions contextKey = "permissions"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "header Authorization manquant ou invalide"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "token invalide ou expiré"})
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserID, claims.UserID)
		ctx = context.WithValue(ctx, ContextRoles, claims.Roles)
		ctx = context.WithValue(ctx, ContextPermissions, claims.Permissions)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAdmin vérifie que l'utilisateur a le rôle admin.
// Utilisé comme garde-fou global sur toutes les routes /admin/*.
func RequireAdmin(next http.Handler) http.Handler {
	return RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roles, _ := r.Context().Value(ContextRoles).([]string)

		if slices.Contains(roles, config.RoleAdmin) {
			next.ServeHTTP(w, r)
			return
		}

		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "accès réservé aux administrateurs"})
	}))
}

// RequirePermission vérifie qu'une permission précise est présente dans le token.
// Utilisé pour les routes accessibles à plusieurs rôles selon leurs permissions.
func RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			perms, _ := r.Context().Value(ContextPermissions).([]string)

			if slices.Contains(perms, permission) {
				next.ServeHTTP(w, r)
				return
			}

			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "permission insuffisante"})
		}))
	}
}

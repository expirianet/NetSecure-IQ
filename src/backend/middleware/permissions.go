package middleware

import (
	"database/sql"
	"net/http"
)

// AttachPermissions est un "no-op" middleware ici (perms sont résolues dans /api/auth/me).
// On le garde pour compatibilité avec ton main.go.
func AttachPermissions(_ *sql.DB) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			next(w, r)
		}
	}
}

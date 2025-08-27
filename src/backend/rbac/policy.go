package rbac

import (
	"net/http"

	"netsecure-backend/middleware"
)

// EnforceOperatorScope est un simple wrapper : les handlers appliquent le filtre
// (Admin -> tout ; Operator -> org ; User -> 403 sur endpoints admin/op).
func EnforceOperatorScope(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := middleware.GetRoleName(r)
		if role == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Ici on pourrait imposer des headers/contexte supplémentaires.
		next(w, r)
	}
}

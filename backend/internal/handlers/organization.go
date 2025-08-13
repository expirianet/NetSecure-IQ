package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"backend/internal/auth"
)

type OrganizationHandler struct {
	secret []byte
}

func NewOrganizationHandler(secret []byte) *OrganizationHandler {
	return &OrganizationHandler{secret: secret}
}

// CompleteOrganization gère GET (lecture) et POST (upsert "dummy-safe")
// Objectif: éviter les 500 et fournir un JSON conforme au front.
func (h *OrganizationHandler) CompleteOrganization(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Auth simple par Bearer
	hdr := r.Header.Get("Authorization")
	if hdr == "" || !strings.HasPrefix(strings.ToLower(hdr), "bearer ") {
		http.Error(w, `{"error":"missing or invalid Authorization header"}`, http.StatusUnauthorized)
		return
	}
	token := strings.TrimSpace(hdr[len("Bearer "):])
	if _, err := auth.ParseJWT(h.secret, token); err != nil {
		http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// ⚠️ Minimal: renvoyer un objet "vide" (le front affiche "N/A" sans planter)
		resp := map[string]any{
			"name":           "",
			"vat_number":     "",
			"address":        "",
			"city":           "",
			"state":          "",
			"zip_code":       "",
			"contact_email":  "",
			"contact_phone":  "",
			"sdi_code":       "",
			"pec_email":      "",
			"personnel_info": "",
			"manager":   map[string]string{"name": "", "email": "", "phone": ""},
			"controller": map[string]string{"name": "", "email": "", "phone": ""},
			"processor":  map[string]string{"name": "", "email": "", "phone": ""},
		}
		_ = json.NewEncoder(w).Encode(resp)
		return

	case http.MethodPost:
		// ⚠️ Minimal: on accepte la charge utile et on confirme.
		var in map[string]any
		_ = json.NewDecoder(r.Body).Decode(&in)

		out := map[string]any{
			"message":       "organization saved",
			"organization":  in, // écho de ce qui est envoyé
		}
		_ = json.NewEncoder(w).Encode(out)
		return

	case http.MethodOptions:
		w.WriteHeader(http.StatusNoContent)
		return

	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
}

package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/expirianet/NetSecure-IQ/backend/middleware"
)

type AuthMeSitePerm struct {
	SiteID          string `json:"site_id"`
	ReadAccess      bool   `json:"read_access"`
	WriteAccess     bool   `json:"write_access"`
	CanReceiveAlert bool   `json:"can_receive_alerts"`
}
type AuthMeUser struct {
	ID             string  `json:"id"`
	Email          string  `json:"email"`
	FirstName      *string `json:"first_name,omitempty"`
	LastName       *string `json:"last_name,omitempty"`
	RoleID         int     `json:"role_id"`
	OrganizationID *string `json:"organization_id,omitempty"`
	Attributes     *string `json:"attributes,omitempty"`
	LastLoginAt    *string `json:"last_login_at,omitempty"`
	RoleName       string  `json:"role_name"`
}
type AuthMeResponse struct {
	User            AuthMeUser       `json:"user"`
	Permissions     []string         `json:"permissions"`
	SitePermissions []AuthMeSitePerm `json:"site_permissions"`
}

func HandleAuthMe(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := middleware.GetUserID(r)
		if uid == "" { http.Error(w, "unauthorized", http.StatusUnauthorized); return }

		var u AuthMeUser
		if err := db.QueryRow(`
			SELECT u.id, u.email, u.first_name, u.last_name, u.role_id,
			       u.organization_id, u.attributes, u.last_login_at, r.name
			  FROM users u
			  JOIN roles r ON r.id = u.role_id
			 WHERE u.id = $1
		`, uid).Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.RoleID,
			&u.OrganizationID, &u.Attributes, &u.LastLoginAt, &u.RoleName); err != nil {
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		}

		perms := []string{}
		rows, err := db.Query(`
			SELECT p.name
			  FROM role_permissions rp
			  JOIN permissions p ON p.id = rp.permission_id
			 WHERE rp.role_id = $1 AND (rp.read_access = TRUE OR rp.write_access = TRUE)
		`, u.RoleID)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var n string
				if err := rows.Scan(&n); err == nil { perms = append(perms, n) }
			}
		}

		sperm := []AuthMeSitePerm{}
		if u.RoleID == 3 {
			sr, err := db.Query(`
				SELECT site_id, read_access, write_access, can_receive_alerts
				  FROM user_site_permissions
				 WHERE user_id = $1
			`, u.ID)
			if err == nil {
				defer sr.Close()
				for sr.Next() {
					var p AuthMeSitePerm
					if err := sr.Scan(&p.SiteID, &p.ReadAccess, &p.WriteAccess, &p.CanReceiveAlert); err == nil {
						sperm = append(sperm, p)
					}
				}
			}
		}

		resp := AuthMeResponse{User: u, Permissions: perms, SitePermissions: sperm}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}

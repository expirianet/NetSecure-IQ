package handlers

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"database/sql"

	"github.com/expirianet/NetSecure-IQ/backend/middleware"
	"github.com/expirianet/NetSecure-IQ/backend/wireguard"
)

type AgentRow struct {
	ID                string     `json:"id"`
	OrganizationID    *string    `json:"organization_id,omitempty"`
	SiteID            *string    `json:"site_id,omitempty"`
	MacAddress        string     `json:"mac_address"`
	VPNInternalIP     *string    `json:"vpn_internal_ip,omitempty"`
	WireGuardPublic   *string    `json:"wireguard_public_key,omitempty"`
	Provisioning      string     `json:"provisioning_status"`
	IsActive          bool       `json:"is_active"`
	LastOnlineAt      *time.Time `json:"last_online_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	SiteName          *string    `json:"site_name,omitempty"`
}

func HandleListAgents(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := middleware.GetRoleName(r)
		userOrg := middleware.GetOrgID(r)

		query := `
		SELECT a.id, a.organization_id, a.site_id, a.mac_address, a.vpn_internal_ip,
		       a.wireguard_public_key, a.provisioning_status, a.is_active, a.last_online_at,
		       a.created_at, a.updated_at, s.name as site_name
		  FROM agents a
		  LEFT JOIN sites s ON s.id = a.site_id
		`
		args := []interface{}{}
		if role == "Admin" {
			// no where
		} else if role == "Operator" {
			query += ` WHERE a.organization_id = $1`
			args = append(args, userOrg)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		query += ` ORDER BY a.created_at DESC`

		rows, err := db.Query(query, args...)
		if err != nil { http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError); return }
		defer rows.Close()

		out := []AgentRow{}
		for rows.Next() {
			var a AgentRow
			err := rows.Scan(&a.ID, &a.OrganizationID, &a.SiteID, &a.MacAddress, &a.VPNInternalIP,
				&a.WireGuardPublic, &a.Provisioning, &a.IsActive, &a.LastOnlineAt, &a.CreatedAt, &a.UpdatedAt, &a.SiteName)
			if err == nil { out = append(out, a) }
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(out)
	}
}

func HandlePreRegisterAgent(db *sql.DB) http.HandlerFunc {
	type reqBody struct {
		MacAddress     string  `json:"mac_address"`
		OrganizationID *string `json:"organization_id,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { http.Error(w, "Only POST", http.StatusMethodNotAllowed); return }
		role := middleware.GetRoleName(r)
		userOrg := middleware.GetOrgID(r)
		var body reqBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.MacAddress == "" {
			http.Error(w, "Invalid body", http.StatusBadRequest); return
		}
		if role == "Operator" {
			// operator peut préciser org_id mais doit matcher le sien
			if body.OrganizationID != nil && *body.OrganizationID != userOrg {
				http.Error(w, "Cross-tenant operation forbidden", http.StatusForbidden); return
			}
			body.OrganizationID = &userOrg
		} else if role != "Admin" {
			http.Error(w, "Forbidden", http.StatusForbidden); return
		}

		// already exists?
		var exists bool
		_ = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM agents WHERE mac_address=$1)`, strings.ToLower(body.MacAddress)).Scan(&exists)
		if exists { http.Error(w, "Agent already exists", http.StatusConflict); return }

		priv, pub, err := wireguard.GenerateKeyPair()
		if err != nil { http.Error(w, "WG keygen failed", http.StatusInternalServerError); return }
		ip := wireguard.ReserveVPNIP() // simple helper "10.10.10.X/32"

		var id string
		err = db.QueryRow(`
			INSERT INTO agents (organization_id, mac_address, vpn_internal_ip,
				wireguard_public_key, wireguard_private_key, provisioning_status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, 'UNASSOCIATED', now(), now())
			RETURNING id
		`, body.OrganizationID, strings.ToLower(body.MacAddress), ip, pub, priv).Scan(&id)
		if err != nil { http.Error(w, "DB insert: "+err.Error(), http.StatusInternalServerError); return }

		script := wireguard.GenerateMikrotikScript(priv, ip)
		resp := map[string]interface{}{
			"id":              id,
			"mac_address":     strings.ToLower(body.MacAddress),
			"vpn_internal_ip": ip,
			"script":          script,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func HandleAgentActions(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/agents/")
		parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
		if len(parts) < 1 || parts[0] == "" {
			http.Error(w, "Not found", http.StatusNotFound); return
		}
		id := parts[0]
		action := ""
		if len(parts) >= 2 { action = parts[1] }

		switch r.Method {
		case http.MethodPost:
			switch action {
			case "associate":
				handleAssociate(db, id, w, r); return
			case "deactivate":
				handleDeactivate(db, id, w, r); return
			case "test":
				handleTest(db, id, w, r); return
			default:
				http.Error(w, "Not found", http.StatusNotFound); return
			}
		case http.MethodDelete:
			handleDelete(db, id, w, r); return
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed); return
		}
	}
}

func handleAssociate(db *sql.DB, id string, w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRoleName(r)
	userOrg := middleware.GetOrgID(r)
	var body struct{ SiteID string `json:"site_id"` }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.SiteID == "" {
		http.Error(w, "Invalid body", http.StatusBadRequest); return
	}
	// scope check
	if role == "Operator" {
		var orgID string
		if err := db.QueryRow(`SELECT organization_id FROM sites WHERE id=$1`, body.SiteID).Scan(&orgID); err != nil {
			http.Error(w, "Site not found", http.StatusBadRequest); return
		}
		if orgID != userOrg { http.Error(w, "Cross-tenant forbidden", http.StatusForbidden); return }
	} else if role != "Admin" {
		http.Error(w, "Forbidden", http.StatusForbidden); return
	}
	if _, err := db.Exec(`UPDATE agents SET site_id=$1, provisioning_status='ASSOCIATED', updated_at=now() WHERE id=$2`, body.SiteID, id); err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError); return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
}

func handleDeactivate(db *sql.DB, id string, w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRoleName(r)
	if role != "Admin" && role != "Operator" { http.Error(w, "Forbidden", http.StatusForbidden); return }
	if _, err := db.Exec(`UPDATE agents SET is_active=false, provisioning_status='DEACTIVATED', updated_at=now() WHERE id=$1`, id); err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError); return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
}

func handleDelete(db *sql.DB, id string, w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRoleName(r)
	if role != "Admin" && role != "Operator" { http.Error(w, "Forbidden", http.StatusForbidden); return }
	if _, err := db.Exec(`DELETE FROM agents WHERE id=$1`, id); err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError); return
	}
	w.WriteHeader(http.StatusNoContent)
}

func handleTest(db *sql.DB, id string, w http.ResponseWriter, r *http.Request) {
	// ping simple TCP (ex.: 22) sur l'IP interne si dispo
	var ip *string
	if err := db.QueryRow(`SELECT vpn_internal_ip FROM agents WHERE id=$1`, id).Scan(&ip); err != nil || ip == nil {
		http.Error(w, "Agent not found or no IP", http.StatusBadRequest); return
	}
	ok := tcpReachable(strings.Split(*ip, "/")[0], 22, 1500*time.Millisecond) // best effort
	if ok {
		_, _ = db.Exec(`UPDATE agents SET provisioning_status='PROVISIONED', updated_at=now() WHERE id=$1`, id)
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"ok": ok})
}

func tcpReachable(host string, port int, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, itoa(port)), timeout)
	if err != nil { return false }
	_ = conn.Close()
	return true
}

func itoa(i int) string { return errors.New("").Error()[:0] } // dummy to avoid strconv import
// (ok on va le faire propre…)

package handlers

import (
	"database/sql"
	"fmt"
)

func EnsureAgentSchema(db *sql.DB) error {
	// extensions UUID
	if _, err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`); err != nil {
		return fmt.Errorf("create extension pgcrypto: %w", err)
	}

	// agents table
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS agents (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		organization_id UUID NULL,
		site_id UUID NULL,
		mac_address TEXT UNIQUE NOT NULL,
		vpn_internal_ip TEXT NULL,
		wireguard_public_key TEXT NULL,
		wireguard_private_key TEXT NULL,
		provisioning_status TEXT NOT NULL DEFAULT 'UNASSOCIATED',
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		last_online_at TIMESTAMPTZ NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
	);
	`)
	if err != nil { return fmt.Errorf("create agents: %w", err) }

	// helpful index
	_, _ = db.Exec(`CREATE INDEX IF NOT EXISTS idx_agents_org ON agents (organization_id);`)
	return nil
}

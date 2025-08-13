package repository

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"

	"backend/internal/models"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(dsn string) (*Postgres, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &Postgres{db: db}, nil
}
func (p *Postgres) Close() { _ = p.db.Close() }

// Schéma conforme à ton script SQL
func (p *Postgres) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	const q = `
SELECT id::text, email, password_hash, first_name, last_name, role_id, organization_id::text, is_active
FROM users WHERE email = $1 LIMIT 1`
	var u models.User
	var org sql.NullString
	err := p.db.QueryRowContext(ctx, q, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.RoleID, &org, &u.IsActive,
	)
	if err != nil {
		return nil, err
	}
	if org.Valid {
		tmp := org.String
		u.OrganizationID = &tmp
	}
	return &u, nil
}

func (p *Postgres) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	const q = `
SELECT id::text, email, password_hash, first_name, last_name, role_id, organization_id::text, is_active
FROM users WHERE id = $1 LIMIT 1`
	var u models.User
	var org sql.NullString
	err := p.db.QueryRowContext(ctx, q, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.RoleID, &org, &u.IsActive,
	)
	if err != nil {
		return nil, err
	}
	if org.Valid {
		tmp := org.String
		u.OrganizationID = &tmp
	}
	return &u, nil
}

// Construit la liste des permissions depuis role_permissions + permissions
// (ton schéma: permissions(name) comme 'organizations:read', 'monitoring:write', etc.)
func (p *Postgres) BuildPermissions(ctx context.Context, roleID int) ([]string, error) {
	const q = `
SELECT p.name, rp.read_access, rp.write_access
FROM role_permissions rp
JOIN permissions p ON p.id = rp.permission_id
WHERE rp.role_id = $1`
	rows, err := p.db.QueryContext(ctx, q, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permSet := map[string]struct{}{}
	for rows.Next() {
		var name string
		var read, write bool
		if err := rows.Scan(&name, &read, &write); err != nil {
			return nil, err
		}
		if read || write {
			permSet[name] = struct{}{}
		}
	}
	out := make([]string, 0, len(permSet))
	for k := range permSet {
		out = append(out, k)
	}
	return out, nil
}

package models

type User struct {
	ID              string
	Email           string
	PasswordHash    string
	FirstName       *string
	LastName        *string
	RoleID          int
	OrganizationID  *string
	IsActive        bool
}

type MeDTO struct {
	ID             string   `json:"id"`
	Email          string   `json:"email"`
	RoleID         int      `json:"role_id"`
	OrganizationID *string  `json:"organization_id"`
	Permissions    []string `json:"permissions"`
}

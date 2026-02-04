package domain

import (
	"fmt"
	"time"
)

type UserRole string

const (
	UserRoleAdmin        UserRole = "admin"
	UserRoleFleetManager UserRole = "fleet_manager"
	UserRoleTechnician   UserRole = "technician"
	UserRoleViewer       UserRole = "viewer"
)

// User represents an authenticated entity in the system.
type User struct {
	ID           int64      `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"` // Never expose hash in JSON
	Role         UserRole   `json:"role"`
	IsEnabled    bool       `json:"is_enabled"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// Validate ensures the user state is valid for creation/update.
func (u *User) Validate() error {
	if u.Username == "" {
		return fmt.Errorf("username is required")
	}
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}
	switch u.Role {
	case UserRoleAdmin, UserRoleFleetManager, UserRoleTechnician, UserRoleViewer:
		// valid
	default:
		return fmt.Errorf("invalid role: %s", u.Role)
	}
	return nil
}

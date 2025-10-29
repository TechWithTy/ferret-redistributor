package models

import (
	"fmt"
	"time"
)

// Role represents a user's role in the system
type Role string

// System-wide roles
const (
	RoleSuperAdmin Role = "super_admin" // Full system access
	RoleAdmin      Role = "admin"       // Organization admin
	RoleUser       Role = "user"        // Regular user
	RoleViewer     Role = "viewer"      // Read-only access
)

// Permission represents a specific action a user can perform
type Permission string

// System permissions
const (
	// User management
	PermissionViewUsers   Permission = "view_users"
	PermissionCreateUsers Permission = "create_users"
	PermissionEditUsers   Permission = "edit_users"
	PermissionDeleteUsers Permission = "delete_users"

	// Team management
	PermissionViewTeams   Permission = "view_teams"
	PermissionCreateTeams Permission = "create_teams"
	PermissionEditTeams   Permission = "edit_teams"
	PermissionDeleteTeams Permission = "delete_teams"

	// Content management
	PermissionViewContent   Permission = "view_content"
	PermissionCreateContent Permission = "create_content"
	PermissionEditContent   Permission = "edit_content"
	PermissionDeleteContent Permission = "delete_content"

	// Organization settings
	PermissionViewOrgSettings   Permission = "view_org_settings"
	PermissionEditOrgSettings   Permission = "edit_org_settings"
	PermissionManageBilling     Permission = "manage_billing"
	PermissionManageIntegrations Permission = "manage_integrations"
)

// RolePermissions maps roles to their permissions
var RolePermissions = map[Role][]Permission{
	RoleSuperAdmin: {
		// User management
		PermissionViewUsers, PermissionCreateUsers, PermissionEditUsers, PermissionDeleteUsers,
		// Team management
		PermissionViewTeams, PermissionCreateTeams, PermissionEditTeams, PermissionDeleteTeams,
		// Content management
		PermissionViewContent, PermissionCreateContent, PermissionEditContent, PermissionDeleteContent,
		// Organization settings
		PermissionViewOrgSettings, PermissionEditOrgSettings, PermissionManageBilling, PermissionManageIntegrations,
	},
	RoleAdmin: {
		// User management
		PermissionViewUsers, PermissionCreateUsers, PermissionEditUsers,
		// Team management
		PermissionViewTeams, PermissionCreateTeams, PermissionEditTeams,
		// Content management
		PermissionViewContent, PermissionCreateContent, PermissionEditContent, PermissionDeleteContent,
		// Organization settings (limited)
		PermissionViewOrgSettings, PermissionEditOrgSettings,
	},
	RoleUser: {
		// Basic permissions
		PermissionViewUsers, PermissionViewTeams,
		// Content management (own content)
		PermissionViewContent, PermissionCreateContent, PermissionEditContent,
	},
	RoleViewer: {
		// Read-only access
		PermissionViewUsers, PermissionViewTeams, PermissionViewContent,
	},
}

// HasPermission checks if a role has a specific permission
func (r Role) HasPermission(permission Permission) bool {
	permissions, exists := RolePermissions[r]
	if !exists {
		return false
	}

	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// UserRole represents a user's role within an organization
type UserRole struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	OrgID        string    `json:"org_id"`
	Role         Role      `json:"role"`
	AssignedBy   string    `json:"assigned_by"`
	AssignedAt   time.Time `json:"assigned_at"`
	LastModified time.Time `json:"last_modified"`
}

// TeamRole represents a user's role within a team
type TeamRole struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	TeamID       string    `json:"team_id"`
	Role         Role      `json:"role"`
	AssignedBy   string    `json:"assigned_by"`
	AssignedAt   time.Time `json:"assigned_at"`
	LastModified time.Time `json:"last_modified"`
}

// HasGlobalPermission checks if a user has a global permission
func (u *User) HasGlobalPermission(permission Permission) bool {
	// Super admins have all permissions
	if u.Role == RoleSuperAdmin {
		return true
	}

	// Check role permissions
	return Role(u.Role).HasPermission(permission)
}

// HasOrgPermission checks if a user has a specific permission within an organization
func (u *User) HasOrgPermission(orgID string, permission Permission) bool {
	// Super admins have all permissions
	if u.Role == RoleSuperAdmin {
		return true
	}

	// Check if user has the permission through their role in the organization
	// This would typically involve checking the UserRole for this org
	// For now, we'll just check their global role
	return u.HasGlobalPermission(permission)
}

// HasTeamPermission checks if a user has a specific permission within a team
func (u *User) HasTeamPermission(teamID string, permission Permission) bool {
	// Super admins have all permissions
	if u.Role == RoleSuperAdmin {
		return true
	}

	// Check if user has the permission through their role in the team
	// This would typically involve checking the TeamRole for this team
	// For now, we'll check their global permissions
	return u.HasGlobalPermission(permission)
}

// CanManageUser checks if a user can manage another user
func (u *User) CanManageUser(targetUser *User) bool {
	// Super admins can manage anyone
	if u.Role == RoleSuperAdmin {
		return true
	}

	// Users can't manage users with higher or equal roles
	if u.Role == RoleAdmin && (targetUser.Role == RoleSuperAdmin || targetUser.Role == RoleAdmin) {
		return false
	}

	// Admins can manage non-admin users
	if u.Role == RoleAdmin {
		return true
	}

	// Regular users can only manage themselves
	return u.ID == targetUser.ID
}

// UpdateUserRole updates a user's role with validation
func UpdateUserRole(user *User, newRole Role, updatedBy string) error {
	// Validate the new role
	switch newRole {
	case RoleSuperAdmin, RoleAdmin, RoleUser, RoleViewer:
		// Valid role
	default:
		return fmt.Errorf("invalid role: %s", newRole)
	}

	// Update the user's role
	user.Role = newRole
	// In a real implementation, you would also update timestamps and track who made the change
	user.UpdatedAt = time.Now()

	return nil
}

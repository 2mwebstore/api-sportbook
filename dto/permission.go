package dto

// ── Responses ─────────────────────────────────────────────

type PermissionResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RolePermissionResponse struct {
	Role        string               `json:"role"`
	Permissions []PermissionResponse `json:"permissions"`
}

// ── Requests ──────────────────────────────────────────────

// SetRolePermissionsRequest replaces all permissions for a role.
type SetRolePermissionsRequest struct {
	Permissions []string `json:"permissions" validate:"required,min=1,dive,min=1"`
}

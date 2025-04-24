package dto

// CreateRoleDTO represents the data needed to create a new role
type CreateRoleDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// UpdateRoleDTO represents the data for updating a role
type UpdateRoleDTO struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

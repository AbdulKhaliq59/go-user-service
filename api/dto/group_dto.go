package dto

// CreateGroupDTO represents data needed to create a new group
type CreateGroupDTO struct {
	Name string `json:"name" binding:"required"`
}

// UpdateGroupDTO represents data needed to update an existing group
type UpdateGroupDTO struct {
	Name string `json:"name" binding:"required"`
}

// GroupResponseDTO represents the data structure for group responses
type GroupResponseDTO struct {
	ID    string       `json:"id"`
	Name  string       `json:"name"`
	Roles []RoleSimple `json:"roles,omitempty"`
}

// RoleSimple represents simplified role information for responses
type RoleSimple struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// AssignRolesToGroupDTO represents data needed to assign roles to a group
type AssignRolesToGroupDTO struct {
	GroupID string   `json:"groupId" binding:"required"`
	RoleIDs []string `json:"roleIds" binding:"required"`
}

// UnassignRoleFromGroupDTO represents data needed to unassign roles from a group
type UnassignRoleFromGroupDTO struct {
	GroupID string   `json:"groupId" binding:"required"`
	RoleIDs []string `json:"roleIds" binding:"required"`
}

// GroupsWithRolesDTO represents a map of group names to their assigned roles
type GroupsWithRolesDTO map[string][]RoleSimple

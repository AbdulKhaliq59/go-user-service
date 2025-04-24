package controllers

import (
	"net/http"

	"go-user-service/api/dto"
	"go-user-service/services"
	"go-user-service/utils"

	"github.com/gin-gonic/gin"
)

// GroupController handles group-related HTTP requests
type GroupController struct {
	groupService *services.GroupService
}

// NewGroupController creates a new group controller
func NewGroupController(groupService *services.GroupService) *GroupController {
	return &GroupController{
		groupService: groupService,
	}
}

// @Summary Get all groups
// @Description Get all groups with their roles
// @Tags Groups
// @Accept json
// @Produce json
// @Param pageNumber query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} utils.PaginationResponse
// @Router /api/v1/group [get]
func (c *GroupController) GetGroups(ctx *gin.Context) {
	// Parse pagination parameters
	var query utils.PaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default values if not provided
	if query.PageSize == 0 {
		query.PageSize = 10
	}
	if query.PageNumber == 0 {
		query.PageNumber = 1
	}

	// Get all groups
	groups, err := c.groupService.GetGroups(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch groups: " + err.Error()})
		return
	}

	// Convert to response DTOs
	var groupDTOs []dto.GroupResponseDTO
	for _, group := range groups {
		// Get roles for this group
		roles, err := c.groupService.GetGroupRoles(ctx, *group.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch roles for group: " + err.Error()})
			return
		}

		// Convert roles to simple format
		var roleDTOs []dto.RoleSimple
		for _, role := range roles {
			description := ""
			if role.Description != nil {
				description = *role.Description
			}

			roleDTOs = append(roleDTOs, dto.RoleSimple{
				ID:          *role.ID,
				Name:        *role.Name,
				Description: description,
			})
		}

		// Create group DTO with roles
		groupDTOs = append(groupDTOs, dto.GroupResponseDTO{
			ID:    *group.ID,
			Name:  *group.Name,
			Roles: roleDTOs,
		})
	}

	// Calculate pagination
	total := int64(len(groupDTOs))
	start := (query.PageNumber - 1) * query.PageSize
	end := start + query.PageSize
	if start >= int(total) {
		start = 0
		end = 0
	} else if end > int(total) {
		end = int(total)
	}

	// Return paginated result
	paginatedGroups := groupDTOs
	if start < end {
		paginatedGroups = groupDTOs[start:end]
	} else {
		paginatedGroups = []dto.GroupResponseDTO{}
	}

	response := utils.CreatePaginationResponse(paginatedGroups, total, query)
	ctx.JSON(http.StatusOK, response)
}

// @Summary Get group by ID
// @Description Get a specific group by its ID
// @Tags Groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {object} dto.GroupResponseDTO
// @Router /api/v1/group/{id} [get]
func (c *GroupController) GetGroupByID(ctx *gin.Context) {
	id := ctx.Param("id")

	// Get the group
	group, err := c.groupService.GetGroupById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Group not found: " + err.Error()})
		return
	}

	// Get roles for this group
	roles, err := c.groupService.GetGroupRoles(ctx, *group.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch roles for group: " + err.Error()})
		return
	}

	// Convert roles to simple format
	var roleDTOs []dto.RoleSimple
	for _, role := range roles {
		description := ""
		if role.Description != nil {
			description = *role.Description
		}

		roleDTOs = append(roleDTOs, dto.RoleSimple{
			ID:          *role.ID,
			Name:        *role.Name,
			Description: description,
		})
	}

	// Create response
	response := dto.GroupResponseDTO{
		ID:    *group.ID,
		Name:  *group.Name,
		Roles: roleDTOs,
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Create a new group
// @Description Create a new group
// @Tags Groups
// @Accept json
// @Produce json
// @Param group body dto.CreateGroupDTO true "Group data"
// @Router /api/v1/group [post]
func (c *GroupController) CreateGroup(ctx *gin.Context) {
	var createDto dto.CreateGroupDTO
	if err := ctx.ShouldBindJSON(&createDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.groupService.CreateGroup(ctx, createDto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Group created successfully"})
}

// @Summary Update a group
// @Description Update an existing group
// @Tags Groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param group body dto.UpdateGroupDTO true "Updated group data"
// @Router /api/v1/group/{id} [patch]
func (c *GroupController) UpdateGroup(ctx *gin.Context) {
	id := ctx.Param("id")

	var updateDto dto.UpdateGroupDTO
	if err := ctx.ShouldBindJSON(&updateDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.groupService.UpdateGroup(ctx, id, updateDto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update group: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Group updated successfully"})
}

// @Summary Delete a group
// @Description Delete a group by ID
// @Tags Groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Router /api/v1/group/{id} [delete]
func (c *GroupController) DeleteGroup(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.groupService.DeleteGroup(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete group: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Group deleted successfully"})
}

// @Summary Assign roles to group
// @Description Assign one or more roles to a group
// @Tags Groups
// @Accept json
// @Produce json
// @Param roles body dto.AssignRolesToGroupDTO true "Role assignment data"
// @Router /api/v1/group/assign-roles [post]
func (c *GroupController) AssignRolesToGroup(ctx *gin.Context) {
	var assignDto dto.AssignRolesToGroupDTO
	if err := ctx.ShouldBindJSON(&assignDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.groupService.AssignRolesToGroup(ctx, assignDto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign roles: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Roles assigned successfully"})
}

// @Summary Unassign roles from group
// @Description Remove one or more roles from a group
// @Tags Groups
// @Accept json
// @Produce json
// @Param roles body dto.UnassignRoleFromGroupDTO true "Role unassignment data"
// @Router /api/v1/group/unassign-roles [post]
func (c *GroupController) UnassignRoleFromGroup(ctx *gin.Context) {
	var unassignDto dto.UnassignRoleFromGroupDTO
	if err := ctx.ShouldBindJSON(&unassignDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.groupService.UnassignRoleFromGroup(ctx, unassignDto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unassign roles: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Roles unassigned successfully"})
}

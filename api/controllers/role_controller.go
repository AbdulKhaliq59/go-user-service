package controllers

import (
	"fmt"
	"go-user-service/api/dto"
	"go-user-service/services"
	"go-user-service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoleController handles HTTP requests related to roles
type RoleController struct {
	roleService *services.RoleService
}

// NewRoleController creates a new role controller
func NewRoleController(roleService *services.RoleService) *RoleController {
	return &RoleController{
		roleService: roleService,
	}
}

// CreateRole godoc
// @Summary Create a new role
// @Description Create a new role with the provided information
// @Tags roles
// @Accept json
// @Produce json
// @Param role body dto.CreateRoleDTO true "Role information"
// @Success 201 {object} map[string]interface{} "message: Role created successfully"
// @Failure 400 {object} map[string]interface{} "error: Bad request error message"
// @Failure 500 {object} map[string]interface{} "error: Internal server error message"
// @Router /api/v1/role [post]
// @Security ApiKeyAuth
func (c *RoleController) CreateRole(ctx *gin.Context) {
	var createRoleDto dto.CreateRoleDTO
	if err := ctx.ShouldBindJSON(&createRoleDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.roleService.CreateRole(ctx, createRoleDto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Role created successfully"})
}

// GetRoles godoc
// @Summary Get all roles
// @Description Get a paginated list of all roles
// @Tags roles
// @Accept json
// @Produce json
// @Param page_size query int false "Page size" default(10)
// @Param page_number query int false "Page number" default(1)
// @Success 200 {object} utils.PaginationResponse
// @Failure 400 {object} map[string]interface{} "error: Bad request error message"
// @Failure 500 {object} map[string]interface{} "error: Internal server error message"
// @Router /api/v1/role [get]
func (c *RoleController) GetRoles(ctx *gin.Context) {
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

	roles, err := c.roleService.GetRoles(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error fetching roles",
			"error":   err.Error(),
		})
		return
	}

	// Create a paginated response
	response := utils.CreatePaginationResponse(roles, int64(len(roles)), query)
	ctx.JSON(http.StatusOK, response)
}

// FindById godoc
// @Summary Get a role by ID
// @Description Get detailed information about a specific role by its ID
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Success 200 {object} interface{} "Role object"
// @Failure 404 {object} map[string]interface{} "message: Role not found"
// @Failure 500 {object} map[string]interface{} "error: Internal server error message"
// @Router /api/v1/role/{id} [get]
func (c *RoleController) FindById(ctx *gin.Context) {
	roleId := ctx.Param("id")
	role, err := c.roleService.FindById(ctx, roleId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if role == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Role not found"})
		return
	}

	ctx.JSON(http.StatusOK, role)
}

// UpdateRole godoc
// @Summary Update a role
// @Description Update a role's information by its ID
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Param role body dto.UpdateRoleDTO true "Updated role information"
// @Success 200 {object} map[string]interface{} "message: Role updated successfully"
// @Failure 400 {object} map[string]interface{} "error: Bad request error message"
// @Failure 500 {object} map[string]interface{} "error: Internal server error message"
// @Router /api/v1/role/{id} [patch]
// @Security ApiKeyAuth
func (c *RoleController) UpdateRole(ctx *gin.Context) {
	roleId := ctx.Param("id")
	var updateRoleDto dto.UpdateRoleDTO
	if err := ctx.ShouldBindJSON(&updateRoleDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.roleService.UpdateRole(ctx, roleId, updateRoleDto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error updating role",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Role updated successfully"})
}

// DeleteRole godoc
// @Summary Delete a role
// @Description Delete a role by its ID
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Success 200 {object} map[string]interface{} "message: Role deleted successfully"
// @Failure 500 {object} map[string]interface{} "error: Internal server error message"
// @Router /api/v1/role/{id} [delete]
// @Security ApiKeyAuth
func (c *RoleController) DeleteRole(ctx *gin.Context) {
	roleId := ctx.Param("id")
	fmt.Printf("Deleting role with ID: %s\n", roleId)
	err := c.roleService.DeleteRole(ctx, roleId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error deleting role",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}

// Below are the currently commented-out methods in the original code
// They can be uncommented and implemented as needed

// GetClient godoc
// @Summary Get all clients
// @Description Get a list of all clients
// @Tags clients
// @Accept json
// @Produce json
// @Success 200 {array} interface{} "List of clients"
// @Failure 500 {object} map[string]interface{} "error: Internal server error message"
// @Router /api/v1/role/client [get]
// func (c *RoleController) GetClient(ctx *gin.Context) {
// 	clients, err := c.roleService.GetClient(ctx)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"message": "Error fetching clients",
// 			"error":   err.Error(),
// 		})
// 		return
// 	}
//
// 	ctx.JSON(http.StatusOK, clients)
// }

// GetClientRoles godoc
// @Summary Get roles for a client
// @Description Get all roles associated with a specific client
// @Tags clients
// @Accept json
// @Produce json
// @Param id path string true "Client ID"
// @Success 200 {array} interface{} "List of client roles"
// @Failure 500 {object} map[string]interface{} "error: Internal server error message"
// @Router /api/v1/role/client/{id} [get]
// func (c *RoleController) GetClientRoles(ctx *gin.Context) {
// 	clientId := ctx.Param("id")
// 	roles, err := c.roleService.GetClientRoles(ctx, clientId)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	ctx.JSON(http.StatusOK, roles)
// }

// CreateClientRole godoc
// @Summary Create a client role
// @Description Create a new role for a specific client
// @Tags clients
// @Accept json
// @Produce json
// @Param role body dto.CreateRoleDTO true "Client role information"
// @Success 201 {object} map[string]interface{} "message: Client role created successfully"
// @Failure 400 {object} map[string]interface{} "error: Bad request error message"
// @Failure 500 {object} map[string]interface{} "error: Internal server error message"
// @Router /api/v1/role/client [post]
// @Security ApiKeyAuth
// func (c *RoleController) CreateClientRole(ctx *gin.Context) {
// 	var createRoleDto dto.CreateRoleDTO
// 	if err := ctx.ShouldBindJSON(&createRoleDto); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	err := c.roleService.CreateClientRole(ctx, createRoleDto)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	ctx.JSON(http.StatusCreated, gin.H{"message": "Client role created successfully"})
// }

// api/controllers/user_controller.go
package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-user-service/api/dto"
	"go-user-service/middleware"
	"go-user-service/models"
	"go-user-service/services"
	"go-user-service/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	authService services.AuthService
}

func NewUserController(authService services.AuthService) *UserController {
	return &UserController{
		authService: authService,
	}
}

// Login godoc
// @Summary User login
// @Description Authenticates a user and returns an access token
// @Tags User
// @Accept json
// @Produce json
// @Param login body dto.LoginDto true "Login credentials"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/user/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var loginDto dto.LoginDto
	if err := ctx.ShouldBindJSON(&loginDto); err != nil {
		ctx.JSON(http.StatusBadRequest, &models.Response{
			Timestamp: time.Now(),
			Message:   "Invalid request body",
			Status:    http.StatusBadRequest,
			Data:      nil,
		})
		return
	}

	response, err := c.authService.Login(loginDto.Username, loginDto.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, &models.Response{
			Timestamp: time.Now(),
			Message:   "Incorrect username or password",
			Status:    http.StatusUnauthorized,
			Data:      nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// CreateUser godoc
// @Summary Create new user
// @Description Creates a new user in Keycloak
// @Tags User
// @Accept json
// @Produce json
// @Param user body dto.CreateUserDto true "User details"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 409 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/user [post]
// @Security Bearer
func (c *UserController) CreateUser(ctx *gin.Context) {
	var createUserDto dto.CreateUserDto
	if err := ctx.ShouldBindJSON(&createUserDto); err != nil {
		ctx.JSON(http.StatusBadRequest, &models.Response{
			Timestamp: time.Now(),
			Message:   "Invalid request body",
			Status:    http.StatusBadRequest,
			Data:      nil,
		})
		return
	}

	err := c.authService.CreateNewUser(&createUserDto)
	if err != nil {
		// Handle specific error cases
		if strings.Contains(err.Error(), "username is already taken") ||
			strings.Contains(err.Error(), "email is already taken") {
			ctx.JSON(http.StatusConflict, &models.Response{
				Timestamp: time.Now(),
				Message:   err.Error(),
				Status:    http.StatusConflict,
				Data:      nil,
			})
			return
		}

		// Generic error
		ctx.JSON(http.StatusInternalServerError, &models.Response{
			Timestamp: time.Now(),
			Message:   "Failed to create user",
			Status:    http.StatusInternalServerError,
			Data:      nil,
		})
		return
	}

	ctx.JSON(http.StatusCreated, &models.Response{
		Timestamp: time.Now(),
		Message:   fmt.Sprintf("User %s created successfully", createUserDto.Username),
		Status:    http.StatusCreated,
		Data:      nil,
	})
}

// Logout godoc
// @Summary User logout
// @Description Logs out a user by revoking their refresh token
// @Tags User
// @Accept json
// @Produce json
// @Security Bearer
// @Param logout body dto.LogoutDto true "Refresh token"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/user/logout [post]
func (c *UserController) Logout(ctx *gin.Context) {
	var logoutDto dto.LogoutDto
	if err := ctx.ShouldBindJSON(&logoutDto); err != nil {
		ctx.JSON(http.StatusBadRequest, &models.Response{
			Timestamp: time.Now(),
			Message:   "Invalid request body",
			Status:    http.StatusBadRequest,
			Data:      nil,
		})
		return
	}

	err := c.authService.Logout(logoutDto.RefreshToken.Token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &models.Response{
			Timestamp: time.Now(),
			Message:   "Failed to log out",
			Status:    http.StatusInternalServerError,
			Data:      nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, &models.Response{
		Timestamp: time.Now(),
		Message:   "You are logged out successfully",
		Status:    http.StatusOK,
		Data:      nil,
	})
}

// CreateUser godoc
// @Summary Create new standard user
// @Description Creates a new standard user in Keycloak
// @Tags User
// @Accept json
// @Produce json
// @Param user body dto.UserDto true "User details"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 409 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/user/signup [post]
func (c *UserController) CreateStandardUser(ctx *gin.Context) {
	var createUserDto dto.UserDto
	if err := ctx.ShouldBindJSON(&createUserDto); err != nil {
		ctx.JSON(http.StatusBadRequest, &models.Response{
			Timestamp: time.Now(),
			Message:   "Invalid request body",
			Status:    http.StatusBadRequest,
			Data:      nil,
		})
		return
	}

	err := c.authService.Create(&createUserDto)
	if err != nil {
		// Handle specific error cases
		if strings.Contains(err.Error(), "username is already taken") ||
			strings.Contains(err.Error(), "email is already taken") {
			ctx.JSON(http.StatusConflict, &models.Response{
				Timestamp: time.Now(),
				Message:   err.Error(),
				Status:    http.StatusConflict,
				Data:      nil,
			})
			return
		}

		// Generic error
		ctx.JSON(http.StatusInternalServerError, &models.Response{
			Timestamp: time.Now(),
			Message:   "Failed to create user",
			Status:    http.StatusInternalServerError,
			Data:      nil,
		})
		return
	}

	ctx.JSON(http.StatusCreated, &models.Response{
		Timestamp: time.Now(),
		Message:   fmt.Sprintf("User %s created successfully", createUserDto.Username),
		Status:    http.StatusCreated,
		Data:      nil,
	})
}

// CreatePlayer godoc
// @Summary Create new player user
// @Description Creates a new player user in Keycloak and assigns player role
// @Tags User
// @Accept json
// @Produce json
// @Param user body dto.CreateUserDto true "Player user details"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 409 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/user/signup/player [post]
func (c *UserController) CreatePlayer(ctx *gin.Context) {
	var createUserDto dto.CreateUserDto
	if err := ctx.ShouldBindJSON(&createUserDto); err != nil {
		ctx.JSON(http.StatusBadRequest, &models.Response{
			Timestamp: time.Now(),
			Message:   "Invalid request body",
			Status:    http.StatusBadRequest,
			Data:      nil,
		})
		return
	}

	err := c.authService.CreatePlayer(&createUserDto)
	if err != nil {
		// Handle specific error cases
		if strings.Contains(err.Error(), "username is already taken") ||
			strings.Contains(err.Error(), "email is already taken") {
			ctx.JSON(http.StatusConflict, &models.Response{
				Timestamp: time.Now(),
				Message:   err.Error(),
				Status:    http.StatusConflict,
				Data:      nil,
			})
			return
		}

		// Generic error
		ctx.JSON(http.StatusInternalServerError, &models.Response{
			Timestamp: time.Now(),
			Message:   "Failed to create player: " + err.Error(),
			Status:    http.StatusInternalServerError,
			Data:      nil,
		})
		return
	}

	ctx.JSON(http.StatusCreated, &models.Response{
		Timestamp: time.Now(),
		Message:   fmt.Sprintf("User %s created successfully", createUserDto.Username),
		Status:    http.StatusCreated,
		Data:      nil,
	})
}

// CheckUserAuth godoc
// @Summary Check authentication status
// @Description Returns user information if authenticated
// @Tags User
// @Produce json
// @Security Bearer
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response
// @Router /api/v1/user/check [get]
func (c *UserController) CheckUserAuth(ctx *gin.Context) {
	// Get user from context (set by auth middleware)
	userClaims, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, &models.Response{
			Timestamp: time.Now(),
			Message:   "You are not logged in",
			Status:    http.StatusUnauthorized,
			Data:      nil,
		})
		return
	}

	// Extract user ID from claims
	claims, ok := userClaims.(services.KeycloakClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, &models.Response{
			Timestamp: time.Now(),
			Message:   "Failed to process user information",
			Status:    http.StatusInternalServerError,
			Data:      nil,
		})
		return
	}

	// Get user info using the user ID (sub claim)
	userInfo, err := c.authService.GetUserInfo(claims.Subject)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &models.Response{
			Timestamp: time.Now(),
			Message:   "Failed to retrieve user information",
			Status:    http.StatusInternalServerError,
			Data:      nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, &models.Response{
		Timestamp: time.Now(),
		Message:   "User authenticated",
		Status:    http.StatusOK,
		Data:      userInfo,
	})
}

// AssignUserToGroup godoc
// @Summary Assign user to a group
// @Description Assigns a user to a specified group
// @Tags User
// @Accept json
// @Produce json
// @Security Bearer
// @Param assignment body dto.AssignUserToGroupDto true "User-Group assignment details"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/user/group [post]
func (c *UserController) AssignUserToGroup(ctx *gin.Context) {
	// Check if user has the required role
	if !middleware.HasRole(ctx, "realm:assign_users_to_group") {
		ctx.JSON(http.StatusForbidden, &models.Response{
			Timestamp: time.Now(),
			Message:   "Insufficient permissions",
			Status:    http.StatusForbidden,
			Data:      nil,
		})
		return
	}

	var assignDto dto.AssignUserToGroupDto
	if err := ctx.ShouldBindJSON(&assignDto); err != nil {
		ctx.JSON(http.StatusBadRequest, &models.Response{
			Timestamp: time.Now(),
			Message:   "Invalid request body",
			Status:    http.StatusBadRequest,
			Data:      nil,
		})
		return
	}

	err := c.authService.AssignGroupToUser(&assignDto)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusBadRequest, &models.Response{
				Timestamp: time.Now(),
				Message:   err.Error(),
				Status:    http.StatusBadRequest,
				Data:      nil,
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, &models.Response{
			Timestamp: time.Now(),
			Message:   "Failed to assign user to group: " + err.Error(),
			Status:    http.StatusInternalServerError,
			Data:      nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, &models.Response{
		Timestamp: time.Now(),
		Message:   fmt.Sprintf("Assigned user %s to group %s successfully", assignDto.UserId, assignDto.GroupId),
		Status:    http.StatusOK,
		Data:      nil,
	})
}

// api/controllers/user_controller.go - Add these methods to your UserController

// UnassignUserFromGroup godoc
// @Summary Unassign user from a group
// @Description Removes a user from a specified group
// @Tags User
// @Accept json
// @Produce json
// @Security Bearer
// @Param unassignment body dto.UnassignUserFromGroupDto true "User-Group unassignment details"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/user/unassign-group [post]
func (c *UserController) UnassignUserFromGroup(ctx *gin.Context) {
	// Check if user has the required role
	if !middleware.HasRole(ctx, "realm:unassign_users_to_group") {
		ctx.JSON(http.StatusForbidden, &models.Response{
			Timestamp: time.Now(),
			Message:   "Insufficient permissions",
			Status:    http.StatusForbidden,
			Data:      nil,
		})
		return
	}

	var unassignDto dto.UnassignUserFromGroupDto
	if err := ctx.ShouldBindJSON(&unassignDto); err != nil {
		ctx.JSON(http.StatusBadRequest, &models.Response{
			Timestamp: time.Now(),
			Message:   "Invalid request body",
			Status:    http.StatusBadRequest,
			Data:      nil,
		})
		return
	}

	err := c.authService.UnassignGroupFromUser(&unassignDto)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusBadRequest, &models.Response{
				Timestamp: time.Now(),
				Message:   err.Error(),
				Status:    http.StatusBadRequest,
				Data:      nil,
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, &models.Response{
			Timestamp: time.Now(),
			Message:   "Failed to unassign user from group: " + err.Error(),
			Status:    http.StatusInternalServerError,
			Data:      nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, &models.Response{
		Timestamp: time.Now(),
		Message:   fmt.Sprintf("Unassigned user %s from group %s successfully", unassignDto.UserId, unassignDto.GroupId),
		Status:    http.StatusOK,
		Data:      nil,
	})
}

// UnassignUserFromGroups godoc
// @Summary Unassign user from multiple groups
// @Description Removes a user from multiple specified groups
// @Tags User
// @Accept json
// @Produce json
// @Security Bearer
// @Param unassignment body dto.UnassignUserFromGroupsDto true "User-Groups unassignment details"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/user/unassign-groups [post]
func (c *UserController) UnassignUserFromGroups(ctx *gin.Context) {
	// Check if user has the required role
	if !middleware.HasRole(ctx, "realm:unassign_users_to_group") {
		ctx.JSON(http.StatusForbidden, &models.Response{
			Timestamp: time.Now(),
			Message:   "Insufficient permissions",
			Status:    http.StatusForbidden,
			Data:      nil,
		})
		return
	}

	var unassignDto dto.UnassignUserFromGroupsDto
	if err := ctx.ShouldBindJSON(&unassignDto); err != nil {
		ctx.JSON(http.StatusBadRequest, &models.Response{
			Timestamp: time.Now(),
			Message:   "Invalid request body",
			Status:    http.StatusBadRequest,
			Data:      nil,
		})
		return
	}

	err := c.authService.UnassignGroupsFromUser(&unassignDto)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusBadRequest, &models.Response{
				Timestamp: time.Now(),
				Message:   err.Error(),
				Status:    http.StatusBadRequest,
				Data:      nil,
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, &models.Response{
			Timestamp: time.Now(),
			Message:   "Failed to unassign user from groups: " + err.Error(),
			Status:    http.StatusInternalServerError,
			Data:      nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, &models.Response{
		Timestamp: time.Now(),
		Message:   fmt.Sprintf("Unassigned user %s from groups successfully", unassignDto.UserId),
		Status:    http.StatusOK,
		Data:      nil,
	})
}

// GetAllUsers godoc
// @Summary Get all users with their groups and roles
// @Description Retrieve a paginated list of all users with their associated groups and roles
// @Tags User
// @Accept json
// @Produce json
// @Param pageNumber query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {object} utils.PaginationResponse
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 500 {object} models.Response "Internal Server Error"
// @Router /api/v1/user [get]
// @Security Bearer
func (c *UserController) GetAllUsers(ctx *gin.Context) {
	// Check for required role
	if !middleware.HasRole(ctx, "realm:view_users") {
		ctx.JSON(http.StatusForbidden, &models.Response{
			Timestamp: time.Now(),
			Message:   "Insufficient permissions: realm:view_users role required",
			Status:    http.StatusForbidden,
			Data:      nil,
		})
		return
	}

	// Parse pagination query parameters
	pageNumber, err := strconv.Atoi(ctx.DefaultQuery("pageNumber", "1"))
	if err != nil || pageNumber < 1 {
		pageNumber = 1
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// Get all users with their groups and roles
	users, err := c.authService.GetAllUsersWithGroupsAndRoles()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &models.Response{
			Timestamp: time.Now(),
			Message:   "Failed to fetch users",
			Status:    http.StatusInternalServerError,
			Data:      err.Error(),
		})
		return
	}

	// Create pagination query for the response
	paginationQuery := utils.PaginationQuery{
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}

	// Calculate start and end indices for pagination
	startIndex := (pageNumber - 1) * pageSize
	endIndex := startIndex + pageSize

	if startIndex >= len(users) {
		// Return empty page if start index is beyond the data range
		ctx.JSON(http.StatusOK, utils.CreatePaginationResponse(
			[]map[string]interface{}{},
			int64(len(users)),
			paginationQuery,
		))
		return
	}

	if endIndex > len(users) {
		endIndex = len(users)
	}

	// Apply pagination to the users data
	paginatedUsers := users[startIndex:endIndex]

	// Create and return the pagination response
	response := utils.CreatePaginationResponse(
		paginatedUsers,
		int64(len(users)),
		paginationQuery,
	)

	ctx.JSON(http.StatusOK, response)
}

// UpdateUser godoc
// @Summary Update user information
// @Description Updates an existing user's information
// @Tags User
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Param user body dto.UpdateUserDto true "Updated user details"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/user/{id} [patch]
func (c *UserController) UpdateUser(ctx *gin.Context) {
	// Check if user has the required role
	if !middleware.HasRole(ctx, "realm:update_users") {
		ctx.JSON(http.StatusForbidden, &models.Response{
			Timestamp: time.Now(),
			Message:   "Insufficient permissions",
			Status:    http.StatusForbidden,
			Data:      nil,
		})
		return
	}

	userId := ctx.Param("id")
	if userId == "" {
		ctx.JSON(http.StatusBadRequest, &models.Response{
			Timestamp: time.Now(),
			Message:   "User ID is required",
			Status:    http.StatusBadRequest,
			Data:      nil,
		})
		return
	}

	var updateUserDto dto.UpdateUserDto
	if err := ctx.ShouldBindJSON(&updateUserDto); err != nil {
		ctx.JSON(http.StatusBadRequest, &models.Response{
			Timestamp: time.Now(),
			Message:   "Invalid request body",
			Status:    http.StatusBadRequest,
			Data:      nil,
		})
		return
	}

	err := c.authService.UpdateUser(userId, &updateUserDto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &models.Response{
			Timestamp: time.Now(),
			Message:   "Failed to update user: " + err.Error(),
			Status:    http.StatusInternalServerError,
			Data:      nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, &models.Response{
		Timestamp: time.Now(),
		Message:   "User updated successfully",
		Status:    http.StatusOK,
		Data:      nil,
	})
}

// GetUserById godoc
// @Summary Get user by ID
// @Description Retrieves a user by their ID
// @Tags User
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/user/{id} [get]
func (c *UserController) GetUserById(ctx *gin.Context) {
	userId := ctx.Param("id")
	if userId == "" {
		ctx.JSON(http.StatusBadRequest, &models.Response{
			Timestamp: time.Now(),
			Message:   "User ID is required",
			Status:    http.StatusBadRequest,
			Data:      nil,
		})
		return
	}

	userInfo, err := c.authService.FindUserById(userId)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			ctx.JSON(http.StatusNotFound, &models.Response{
				Timestamp: time.Now(),
				Message:   "User not found",
				Status:    http.StatusNotFound,
				Data:      nil,
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, &models.Response{
			Timestamp: time.Now(),
			Message:   "Failed to retrieve user: " + err.Error(),
			Status:    http.StatusInternalServerError,
			Data:      nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, &models.Response{
		Timestamp: time.Now(),
		Message:   "User found",
		Status:    http.StatusOK,
		Data:      userInfo,
	})
}

// ResetPassword godoc
// @Summary Reset user password
// @Description Resets a user's password
// @Tags User
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Param password body dto.ResetPasswordDto true "New password"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/user/{id}/reset-password [post]
func (c *UserController) ResetPassword(ctx *gin.Context) {
	userId := ctx.Param("id")
	if userId == "" {
		ctx.JSON(http.StatusBadRequest, &models.Response{
			Timestamp: time.Now(),
			Message:   "User ID is required",
			Status:    http.StatusBadRequest,
			Data:      nil,
		})
		return
	}

	var resetPasswordDto dto.ResetPasswordDto
	if err := ctx.ShouldBindJSON(&resetPasswordDto); err != nil {
		ctx.JSON(http.StatusBadRequest, &models.Response{
			Timestamp: time.Now(),
			Message:   "Invalid request body",
			Status:    http.StatusBadRequest,
			Data:      nil,
		})
		return
	}

	err := c.authService.ResetPassword(userId, resetPasswordDto.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &models.Response{
			Timestamp: time.Now(),
			Message:   "Failed to reset password: " + err.Error(),
			Status:    http.StatusInternalServerError,
			Data:      nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, &models.Response{
		Timestamp: time.Now(),
		Message:   "Password updated successfully",
		Status:    http.StatusOK,
		Data:      nil,
	})
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Deletes a user by their ID
// @Tags User
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/user/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	// Check if user has the required role
	if !middleware.HasRole(ctx, "realm:delete_users") {
		ctx.JSON(http.StatusForbidden, &models.Response{
			Timestamp: time.Now(),
			Message:   "Insufficient permissions",
			Status:    http.StatusForbidden,
			Data:      nil,
		})
		return
	}

	userId := ctx.Param("id")
	if userId == "" {
		ctx.JSON(http.StatusBadRequest, &models.Response{
			Timestamp: time.Now(),
			Message:   "User ID is required",
			Status:    http.StatusBadRequest,
			Data:      nil,
		})
		return
	}

	err := c.authService.DeleteUserById(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &models.Response{
			Timestamp: time.Now(),
			Message:   "Failed to delete user: " + err.Error(),
			Status:    http.StatusInternalServerError,
			Data:      nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, &models.Response{
		Timestamp: time.Now(),
		Message:   "User deleted successfully",
		Status:    http.StatusOK,
		Data:      nil,
	})
}

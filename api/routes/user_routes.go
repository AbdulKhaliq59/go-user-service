package routes

import (
	"go-user-service/api/controllers"
	"go-user-service/middleware"
	"go-user-service/services"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.RouterGroup, keycloakService *services.KeycloakService) {
	authService := services.NewAuthService(keycloakService)
	userController := controllers.NewUserController(*authService)

	userRoutes := router.Group("/user")

	// Public routes (no authentication needed)
	{
		userRoutes.POST("/login", userController.Login)
		userRoutes.POST("/signup", userController.CreateStandardUser)
		userRoutes.POST("/signup/player", userController.CreatePlayer)
		userRoutes.GET("/:id", userController.GetUserById)
		userRoutes.POST("/:id/reset-password", userController.ResetPassword)
		userRoutes.POST("/logout", userController.Logout)
	}

	// Protected routes (require authentication)
	authRequired := userRoutes.Group("/")
	authRequired.Use(middleware.AuthMiddleware(keycloakService))
	{
		// Routes requiring ONLY authentication (no specific role)
		authRequired.GET("/check", userController.CheckUserAuth)
		authRequired.PATCH("/:id", userController.UpdateUser)
		// Remove the duplicate routes here

		// Routes requiring authentication + specific roles
		authRequired.POST("", middleware.RoleMiddleware([]string{"realm:create_users"}), userController.CreateUser)
		authRequired.POST("/group", middleware.RoleMiddleware([]string{"realm:assign_users_to_group"}), userController.AssignUserToGroup)
		// authRequired.POST("/groups", middleware.RoleMiddleware([]string{"realm:assign_users_to_group"}), userController.AssignUserToGroups)
		authRequired.POST("/unassign-group", middleware.RoleMiddleware([]string{"realm:unassign_users_to_group"}), userController.UnassignUserFromGroup)
		authRequired.POST("/unassign-groups", middleware.RoleMiddleware([]string{"realm:unassign_users_to_group"}), userController.UnassignUserFromGroups)
		authRequired.GET("", middleware.RoleMiddleware([]string{"realm:view_users"}), userController.GetAllUsers)
		authRequired.DELETE("/:id", middleware.RoleMiddleware([]string{"realm:delete_users"}), userController.DeleteUser)
	}
}

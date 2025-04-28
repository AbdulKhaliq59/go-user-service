// routes/user_routes.go
package routes

import (
	"go-user-service/api/controllers"
	"go-user-service/middleware"
	"go-user-service/services"

	"github.com/gin-gonic/gin"
)

// Update SetupUserRoutes in routes/user_routes.go

func SetupUserRoutes(router *gin.RouterGroup, keycloakService *services.KeycloakService) {
	authService := services.NewAuthService(keycloakService)
	userController := controllers.NewUserController(*authService)

	userRoutes := router.Group("/user")
	{
		// Public routes
		userRoutes.POST("/login", userController.Login)
		userRoutes.POST("", userController.CreateUser)
		userRoutes.POST("/signup", userController.CreateStandardUser)
		userRoutes.POST("/signup/player", userController.CreatePlayer)
		userRoutes.POST("/logout", userController.Logout)
		// Protected routes - require authentication
		authRequired := userRoutes.Group("/")
		authRequired.Use(middleware.AuthMiddleware(keycloakService))
		{
			authRequired.GET("/check", userController.CheckUserAuth)
			authRequired.POST("/group", userController.AssignUserToGroup)
			authRequired.POST("/unassign-group", userController.UnassignUserFromGroup)
			authRequired.POST("/unassign-groups", userController.UnassignUserFromGroups)
			authRequired.GET("", userController.GetAllUsers)
			authRequired.PATCH("/:id", userController.UpdateUser)                  // Update user by ID
			authRequired.DELETE("/:id", userController.DeleteUser)                 // Delete user by ID
			authRequired.POST("/:id/reset-password", userController.ResetPassword) // Reset user password

		}
	}
}

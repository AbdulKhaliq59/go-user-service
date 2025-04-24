package routes

import (
	"go-user-service/api/controllers"
	"go-user-service/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoleRoutes configures all the role routes
func SetupRoleRoutes(router *gin.RouterGroup, db *gorm.DB) {
	// Initialize repository, service, and controller
	roleService := services.NewRoleService()
	roleController := controllers.NewRoleController(roleService)

	roleRoutes := router.Group("/role")
	{
		// Public routes
		roleRoutes.GET("", roleController.GetRoles)
		roleRoutes.GET("/:id", roleController.FindById)

		// Protected routes
		authRoutes := roleRoutes.Group("")
		{
			authRoutes.POST("", roleController.CreateRole)
			authRoutes.PATCH("/:id", roleController.UpdateRole)
			authRoutes.DELETE("/:id", roleController.DeleteRole)
		}

		// Commented out routes from the original code
		// roleRoutes.GET("/client", roleController.GetClient)
		// roleRoutes.GET("/client/:id", roleController.GetClientRoles)
		// roleRoutes.POST("/client", roleController.CreateClientRole)
	}
}

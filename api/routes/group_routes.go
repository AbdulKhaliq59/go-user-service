package routes

import (
	"go-user-service/api/controllers"
	"go-user-service/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupGroupRoutes(router *gin.RouterGroup, db *gorm.DB) {
	groupService := services.NewGroupService()
	groupController := controllers.NewGroupController(groupService)

	groupRoutes := router.Group("/group")
	{
		// Public routes
		groupRoutes.GET("", groupController.GetGroups)
		groupRoutes.GET("/:id", groupController.GetGroupByID)

		// Protected routes
		authRoutes := groupRoutes.Group("")
		{
			authRoutes.POST("", groupController.CreateGroup)
			authRoutes.PATCH("/:id", groupController.UpdateGroup)
			authRoutes.DELETE("/:id", groupController.DeleteGroup)

			authRoutes.POST("/assign-roles", groupController.AssignRolesToGroup)
			authRoutes.POST("/unassign-roles", groupController.UnassignRoleFromGroup)
		}
	}
}

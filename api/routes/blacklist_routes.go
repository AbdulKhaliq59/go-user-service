// api/routes/blacklist_routes.go
package routes

import (
	"go-user-service/api/controllers"
	"go-user-service/middleware"
	"go-user-service/repository"
	"go-user-service/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupBlacklistRoutes(router *gin.RouterGroup, db *gorm.DB) {
	blacklistRepo := repository.NewBlacklistRepository(db)
	blacklistService := services.NewBlacklistService(blacklistRepo)
	blacklistController := controllers.NewBlacklistController(blacklistService)

	blacklistRoutes := router.Group("/blacklist")
	keycloakService, err := services.NewKeycloakService() // Create an instance of KeycloakService
	if err != nil {
		panic("Failed to initialize KeycloakService: " + err.Error())
	}
	blacklistRoutes.Use(middleware.AuthMiddleware(keycloakService)) // Apply authentication middleware
	{
		blacklistRoutes.POST("", blacklistController.CreateBlacklist)
		blacklistRoutes.GET("", blacklistController.GetBlacklists)
		blacklistRoutes.GET("/:id", blacklistController.GetBlacklistByID)
		blacklistRoutes.PUT("/:id", blacklistController.UpdateBlacklist)
		blacklistRoutes.DELETE("/:id", blacklistController.DeleteBlacklist)
	}
}

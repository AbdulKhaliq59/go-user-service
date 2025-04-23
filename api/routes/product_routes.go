// api/routes/product_routes.go
package routes

import (
	"go-user-service/api/controllers"
	"go-user-service/repository"
	"go-user-service/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupProductRoutes(router *gin.RouterGroup, db *gorm.DB) {
	// Initialize repositories
	productRepo := repository.NewProductRepository(db)

	// Initialize services
	productService := services.NewProductService(db, productRepo)

	// Initialize controllers
	productController := controllers.NewProductController(productService, db)

	// Define routes
	productRoutes := router.Group("/product")
	{
		// Public routes
		productRoutes.GET("", productController.GetProducts)
		productRoutes.GET("/ussd", productController.GetProductsForUSSD)

		// File operations
		productRoutes.POST("/upload", productController.UploadFile)
		productRoutes.DELETE("/file/delete/:filename", productController.DeleteFile)

		// Bonus product routes - restructured to avoid conflicts
		productRoutes.GET("/bonus/:bonusProductId/assigned-products", productController.GetProductsByBonus)
		productRoutes.GET("/bonus-products/:productId", productController.GetBonusProducts)
		productRoutes.POST("/assign-bonus/:productId", productController.AssignBonusProducts)
		productRoutes.PATCH("/unassign-bonus/:productId", productController.UnassignBonusProducts)

		// Standard CRUD operations
		productRoutes.POST("", productController.CreateProduct)
		productRoutes.GET("/:id", productController.GetProduct)
		productRoutes.PATCH("/:id", productController.UpdateProduct)
		productRoutes.DELETE("/:id", productController.DeleteProduct)
	}
}

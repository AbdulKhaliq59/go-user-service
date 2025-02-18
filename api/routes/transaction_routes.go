// api/routes/transaction_routes.go
package routes

import (
	"go-user-service/api/controllers"
	"go-user-service/redis"
	"go-user-service/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupTransactionRoutes(router *gin.RouterGroup, db *gorm.DB) {
	redisHelper := redis.NewRedisHelper()

	// Initialize services
	accessKeyService := services.NewAccessKeyService(db, redisHelper)
	transactionService := services.NewTransactionService(db, accessKeyService)

	// Initialize controller
	transactionController := controllers.NewTransactionController(transactionService)

	// Define routes
	transactionRoutes := router.Group("/transactions")
	{
		transactionRoutes.GET("", transactionController.GetAllTransactions)
		transactionRoutes.POST("/initialize", transactionController.Initialize)
		transactionRoutes.POST("/payment/request", transactionController.Payment)
		transactionRoutes.GET("/product-stats", transactionController.GetProductStats)
		transactionRoutes.GET("/player-stats", transactionController.GetPlayerStats)
		transactionRoutes.GET("/:referenceId/status", transactionController.CheckStatus)
		transactionRoutes.GET("/getTokenByPhoneNumber", transactionController.GetTransactionsAndTokensByPhoneNumber)
		transactionRoutes.GET("/my-tokens/:phoneNumber", transactionController.GetMyTokens)
		transactionRoutes.GET("/my-token-stats/:phoneNumber", transactionController.GetTokenStats)
		transactionRoutes.GET("/all-tokens/:phoneNumber", transactionController.GetAllMyTokens)
		transactionRoutes.POST("/resend-token/:token", transactionController.ResendSms)
		transactionRoutes.POST("/regenerate/token/:referenceId", transactionController.RegenerateToken)
		transactionRoutes.GET("/phone_number-hits", transactionController.GetTransactionStats)
		transactionRoutes.GET("/activators/stats", transactionController.GetDailyUserStats)
	}
}

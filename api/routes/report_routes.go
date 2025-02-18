package routes

import (
	"go-user-service/api/controllers"
	"go-user-service/repository"
	"go-user-service/services"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func SetupReportRoutes(router *gin.RouterGroup, db *gorm.DB) {
	// Initialize repository, service, and controller
	reportRepo := repository.NewReportRepository(db)
	reportService := services.NewReportService(reportRepo)
	reportController := controllers.NewReportController(reportService)

	reportRoutes := router.Group("/report")
	{
		reportRoutes.GET("/general", reportController.GeneralReport)
		reportRoutes.GET("/transaction", reportController.TransactionReport)
		reportRoutes.GET("/monthly", reportController.MonthlyReport)
		reportRoutes.GET("/weekly", reportController.WeeklyReport)
		reportRoutes.POST("/expected-transaction", reportController.AddExpectedTransaction)
		reportRoutes.GET("/top-ten", reportController.GetTopTenTransactions)
		reportRoutes.GET("/top-ten/recent", reportController.GetTopTenRecentTransactions)
		reportRoutes.GET("/daily-transactions", reportController.GetDailyTransactionStatistics)
		reportRoutes.GET("/draw", reportController.GetDrawReport)
		reportRoutes.GET("/daily-transactions-report", reportController.GetDailyTransactionReport)
		reportRoutes.GET("/monthly-transaction-report", reportController.GetMonthlyTransactionReport)
		reportRoutes.GET("/hourly-transaction-report", reportController.GetHourlyTransactionReport)

	}
}

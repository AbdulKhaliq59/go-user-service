package controllers

import (
	"go-user-service/models"
	"go-user-service/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	ReportService *services.ReportService
}

func NewReportController(reportService *services.ReportService) *ReportController {
	return &ReportController{ReportService: reportService}
}

// GeneralReport godoc
// @Summary Generate general report
// @Description Generate a general report for the given date range
// @Tags reports
// @Accept json
// @Produce json
// @Param fromDate query string true "From date (YYYY-MM-DD)"
// @Param toDate query string true "To date (YYYY-MM-DD)"
// @Security Bearer
// @Success 200 {object} models.GeneralReportResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /report/general [get]
func (c *ReportController) GeneralReport(ctx *gin.Context) {
	fromDateStr := ctx.Query("fromDate")
	toDateStr := ctx.Query("toDate")

	// Parse dates
	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fromDate format. Expected YYYY-MM-DD"})
		return
	}

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid toDate format. Expected YYYY-MM-DD"})
		return
	}

	// Generate the report
	report, err := c.ReportService.GenerateReport(fromDate, toDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, report)
}

// TransactionReport godoc
// @Summary Generate transaction report
// @Description Generate a transaction report for the given date range, optionally filtered by product ID and status
// @Tags reports
// @Accept json
// @Produce json
// @Param fromDate query string true "From date (YYYY-MM-DD)"
// @Param toDate query string true "To date (YYYY-MM-DD)"
// @Param productId query string false "Product ID"
// @Param status query string false "Transaction status (PENDING, FAILED, SUCCESS)"
// @Security Bearer
// @Success 200 {object} models.TransactionReportResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /report/transaction [get]
func (c *ReportController) TransactionReport(ctx *gin.Context) {
	fromDateStr := ctx.Query("fromDate")
	toDateStr := ctx.Query("toDate")
	productId := ctx.Query("productId")
	status := ctx.Query("status")

	// Parse dates
	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fromDate format. Expected YYYY-MM-DD"})
		return
	}

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid toDate format. Expected YYYY-MM-DD"})
		return
	}

	// Generate the report
	report, err := c.ReportService.TransactionReport(fromDate, toDate, productId, status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, report)
}

// MonthlyReport godoc
// @Summary Generate monthly report
// @Description Generate revenue and transaction statistics grouped by month
// @Tags reports
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} models.MonthlyReportItem
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /report/monthly [get]
func (c *ReportController) MonthlyReport(ctx *gin.Context) {
	// Generate the monthly report
	report, err := c.ReportService.GenerateMonthlyReport()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, report)
}

// WeeklyReport godoc
// @Summary Generate weekly report
// @Description Generate revenue and transaction statistics grouped by week
// @Tags reports
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} models.WeeklyReportItem
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /report/weekly [get]
func (c *ReportController) WeeklyReport(ctx *gin.Context) {
	// Generate the weekly report
	report, err := c.ReportService.GenerateWeeklyReport()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, report)
}

// AddExpectedTransaction godoc
// @Summary Add an expected transaction
// @Description Stores expected revenue and transaction count
// @Tags reports
// @Accept json
// @Produce json
// @Param request body models.ExpectedTransactionRequest true "Expected Transaction Data"
// @Security Bearer
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /report/expected-transaction [post]
func (c *ReportController) AddExpectedTransaction(ctx *gin.Context) {
	var expected models.ExpectedTransactionRequest

	if err := ctx.ShouldBindJSON(&expected); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	transaction := models.ExpectedTransaction{
		ID:           "", // Auto-generated
		Revenue:      expected.Revenue,
		Transactions: expected.Transactions,
	}

	err := c.ReportService.AddExpectedTransaction(transaction)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Expected transaction saved successfully"})
}

// GetTopTenTransactions godoc
// @Summary Get top 10 transactions
// @Description Retrieves the top 10 transactions based on total transaction amount
// @Tags reports
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} models.TopTenTransactionResponse
// @Failure 500 {object} map[string]string
// @Router /report/top-ten [get]
func (c *ReportController) GetTopTenTransactions(ctx *gin.Context) {
	transactions, err := c.ReportService.GetTopTenTransactions()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}

// GetTopTenRecentTransactions godoc
// @Summary Get top 10 recent successful transactions
// @Description Retrieves the top 10 most recent successful transactions
// @Tags reports
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} models.RecentTransaction
// @Failure 500 {object} map[string]string
// @Router /report/top-ten/recent [get]
func (c *ReportController) GetTopTenRecentTransactions(ctx *gin.Context) {
	transactions, err := c.ReportService.GetTopTenRecentTransactions()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}

// GetDailyTransactionStatistics godoc
// @Summary Get daily transaction statistics
// @Description Retrieves transaction statistics for the current day
// @Tags reports
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} models.DailyTransactionStats
// @Failure 500 {object} map[string]string
// @Router /report/daily-transactions [get]
func (c *ReportController) GetDailyTransactionStatistics(ctx *gin.Context) {
	stats, err := c.ReportService.GetDailyTransactionStatistics()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

// Get Daily Transaction Report
// @Tags reports
// @Router /report/daily-transactions-report [get]
func (c *ReportController) GetDailyTransactionReport(ctx *gin.Context) {
	report, err := c.ReportService.GetDailyTransactionReport()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, report)
}

// Get Monthly Transaction Report
// @Tags reports
// @Router /report/monthly-transaction-report [get]
func (c *ReportController) GetMonthlyTransactionReport(ctx *gin.Context) {
	report, err := c.ReportService.GetMonthlyTransactionReport()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, report)
}

// Get Hourly Transaction Report
// @Tags reports
// @Router /report/hourly-transaction-report [get]
func (c *ReportController) GetHourlyTransactionReport(ctx *gin.Context) {
	report, err := c.ReportService.GetHourlyTransactionReport()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, report)
}

// Get Daily Transaction Report godoc
// @Summary Get daily transaction report
// @Tags reports
// @Accept json
// @Produce json
// @Security Bearer
// @Router /report/draw [get]
func (c *ReportController) GetDrawReport(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("pageNumber", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	data, err := c.ReportService.GetDailyTransactionReport()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	paginatedResponse := services.Paginate(data, page, perPage)
	ctx.JSON(http.StatusOK, paginatedResponse)
}

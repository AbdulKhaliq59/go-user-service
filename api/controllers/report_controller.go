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
// @Tags Report
// @Accept json
// @Produce json
// @Param fromDate query string true "From date (YYYY-MM-DD)"
// @Param toDate query string true "To date (YYYY-MM-DD)"
// @Security Bearer
// @Router /api/v1/report/general [get]
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
// @Tags Report
// @Accept json
// @Produce json
// @Param fromDate query string true "From date (YYYY-MM-DD)"
// @Param toDate query string true "To date (YYYY-MM-DD)"
// @Param productId query string false "Product ID"
// @Param status query string false "Transaction status (PENDING, FAILED, SUCCESS)"
// @Security Bearer
// @Router /api/v1/report/transaction [get]
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
// @Tags Report
// @Accept json
// @Produce json
// @Security Bearer
// @Router /api/v1/report/monthly [get]
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
// @Tags Report
// @Accept json
// @Produce json
// @Security Bearer
// @Router /api/v1/report/weekly [get]
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
// @Tags Report
// @Accept json
// @Produce json
// @Param request body models.ExpectedTransactionRequest true "Expected Transaction Data"
// @Security Bearer
// @Router /api/v1/report/expected-transaction [post]
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
// @Tags Report
// @Accept json
// @Produce json
// @Security Bearer
// @Router /api/v1/report/top-ten [get]
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
// @Tags Report
// @Accept json
// @Produce json
// @Security Bearer
// @Router /api/v1/report/top-ten/recent [get]
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
// @Tags Report
// @Accept json
// @Produce json
// @Security Bearer
// @Router /api/v1/report/daily-transactions [get]
func (c *ReportController) GetDailyTransactionStatistics(ctx *gin.Context) {
	stats, err := c.ReportService.GetDailyTransactionStatistics()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

// @Summary Get daily transaction report
// @Description Get detailed daily transaction report for the last 30 days
// @Tags Report
// @Accept json
// @Produce json
// @Router /api/v1/report/daily-transaction-report [get]
func (c *ReportController) GetDailyTransactionReport(ctx *gin.Context) {
	report, err := c.ReportService.GetDailyTransactionReport()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, report)
}

// @Summary Get monthly transaction report
// @Description Get detailed monthly transaction report for the last 12 months
// @Tags Report
// @Accept json
// @Produce json
// @Router /api/v1/report/monthly-transaction-report [get]
func (c *ReportController) GetMonthlyTransactionReport(ctx *gin.Context) {
	report, err := c.ReportService.GetMonthlyTransactionReport()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, report)
}

// @Summary Get hourly transaction report
// @Description Get detailed hourly transaction report for a specific date
// @Tags Report
// @Accept json
// @Produce json
// @Param date query string true "Date in YYYY-MM-DD format"
// @Router /api/v1/report/hourly-transaction-report [get]
func (c *ReportController) GetHourlyTransactionReport(ctx *gin.Context) {
	date := ctx.Query("date")
	if date == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "date parameter is required"})
		return
	}

	report, err := c.ReportService.GetHourlyTransactionReport(date)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, report)
}

// Get Daily Transaction Report godoc
// @Summary Get daily transaction report
// @Tags Report
// @Accept json
// @Produce json
// @Security Bearer
// @Router /api/v1/report/draw [get]
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

// @Summary Get monthly transaction report by product
// @Description Get detailed monthly transaction report filtered by product
// @Tags Report
// @Accept json
// @Produce json
// @Param productId query string false "Product ID (UUID)"
// @Router /api/v1/report/product-monthly-transaction-report [get]
func (c *ReportController) GetProductMonthlyTransactionReport(ctx *gin.Context) {
	productID := ctx.Query("productId")

	report, err := c.ReportService.GetProductMonthlyReport(productID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, report)
}

// @Summary Get daily transaction report by product
// @Description Get detailed daily transaction report filtered by product and date
// @Tags Report
// @Accept json
// @Produce json
// @Param productId query string false "Product ID (UUID)"
// @Param date query string false "Date (YYYY-MM-DD)"
// @Router /api/v1/report/product-daily-transaction-report [get]
func (c *ReportController) GetProductDailyTransactionReport(ctx *gin.Context) {
	productID := ctx.Query("productId")
	date := ctx.Query("date")

	report, err := c.ReportService.GetProductDailyReport(productID, date)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, report)
}

// @Summary Get hourly transaction report by product
// @Description Get detailed hourly transaction report filtered by product and date range
// @Tags Report
// @Accept json
// @Produce json
// @Param productId query string false "Product ID (UUID)"
// @Param date query string false "Date (YYYY-MM-DD)"
// @Param fromDate query string false "From Date (YYYY-MM-DD HH:mm:ss)"
// @Param toDate query string false "To Date (YYYY-MM-DD HH:mm:ss)"
// @Param orderBy query string false "Order by clause (e.g., tot_successful_trx_amount DESC)"
// @Router /api/v1/report/product-hourly-transaction-report [get]
func (c *ReportController) GetProductHourlyTransactionReport(ctx *gin.Context) {
	productID := ctx.Query("productId")
	date := ctx.Query("date")
	fromDate := ctx.Query("fromDate")
	toDate := ctx.Query("toDate")
	orderBy := ctx.Query("orderBy")

	report, err := c.ReportService.GetProductHourlyReport(productID, date, fromDate, toDate, orderBy)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, report)
}

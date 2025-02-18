// api/controllers/transaction_controller.go
package controllers

import (
	"fmt"
	"go-user-service/models"
	"go-user-service/services"
	"go-user-service/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	TransactionService *services.TransactionService
}

func NewTransactionController(transactionService *services.TransactionService) *TransactionController {
	return &TransactionController{
		TransactionService: transactionService,
	}
}

// Initialize godoc
// @Summary Initialize a transaction
// @Description Initializes payment transaction request
// @Tags Transaction
// @Accept json
// @Produce json
// @Param apiKey header string true "API Key"
// @Param transaction body models.CreateTransactionDto true "Transaction data"
// @Success 201
// @Router /api/v1/transactions/initialize [post]
func (c *TransactionController) Initialize(ctx *gin.Context) {
	var createTransactionDto models.CreateTransactionDto
	if err := ctx.ShouldBindJSON(&createTransactionDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apiKey := ctx.GetHeader("apiKey")
	if apiKey == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "API key is required"})
		return
	}

	response, err := c.TransactionService.Initialize(createTransactionDto, apiKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// Payment godoc
// @Summary Process a payment request
// @Description Processes a payment request and sends Kafka event
// @Tags Transaction
// @Accept json
// @Produce json
// @Param payment body models.PaymentDto true "Payment data"
// @Router /api/v1/transactions/payment/request [post]
func (c *TransactionController) Payment(ctx *gin.Context) {
	var paymentDto models.PaymentDto
	if err := ctx.ShouldBindJSON(&paymentDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.TransactionService.Payment(paymentDto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Get all transactions with pagination and filtering
// @Description Retrieve all transactions with optional pagination and filtering
// @Tags Transaction
// @Accept json
// @Produce json
// @Param pageNumber query int false "Page number (default: 1)"
// @Param pageSize query int false "Page size (default: 10)"
// @Param from query string false "Start date (RFC3339 format)"
// @Param to query string false "End date (RFC3339 format)"
// @Param status query string false "Transaction status"
// @Param telecom query string false "Telecom provider"
// @Param productId query string false "Product ID"
// @Param agentCode query string false "Agent code"
// @Param search query string false "Search by phone number"
// @Success 200 {object} utils.PaginationResponse
// @Router /api/v1/transactions [get]
func (c *TransactionController) GetAllTransactions(ctx *gin.Context) {
	var query utils.PaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.TransactionService.FindAll(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Get product statistics
// @Description Retrieve product statistics with optional date range filtering
// @Tags Transaction
// @Accept json
// @Produce json
// @Param from query string false "Start date (RFC3339 format)"
// @Param to query string false "End date (RFC3339 format)"
// @Param range query string false "Date range (daily, weekly, monthly)" Enums(daily, weekly, monthly)
// @Success 200 {array} models.ProductStats
// @Router /api/v1/transactions/product-stats [get]
func (c *TransactionController) GetProductStats(ctx *gin.Context) {
	// Parse query parameters
	fromStr := ctx.Query("from")
	toStr := ctx.Query("to")
	rangeType := ctx.Query("range")

	var from, to *time.Time
	if fromStr != "" {
		parsedFrom, err := time.Parse(time.RFC3339, fromStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'from' date format. Use RFC3339"})
			return
		}
		from = &parsedFrom
	}

	if toStr != "" {
		parsedTo, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'to' date format. Use RFC3339"})
			return
		}
		to = &parsedTo
	}

	// Validate range parameter if provided
	if rangeType != "" && rangeType != "daily" && rangeType != "weekly" && rangeType != "monthly" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid range parameter. Use 'daily', 'weekly', or 'monthly'"})
		return
	}

	// Get product stats
	stats, err := c.TransactionService.GetProductStats(from, to, rangeType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

// GetPlayerStats godoc
// @Summary Get player statistics
// @Description Retrieve player statistics with optional date range filtering
// @Tags Transaction
// @Accept json
// @Produce json
// @Param range query string false "Date range (daily, weekly, monthly)" Enums(daily, weekly, monthly)
// @Success 200 {array} models.PlayerStats
// @Router /api/v1/transactions/player-stats [get]
func (c *TransactionController) GetPlayerStats(ctx *gin.Context) {
	// Parse query parameters
	rangeType := ctx.Query("range")

	// Validate range parameter if provided
	if rangeType != "" && rangeType != "daily" && rangeType != "weekly" && rangeType != "monthly" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid range parameter. Use 'daily', 'weekly', or 'monthly'"})
		return
	}

	// Get player stats
	stats, err := c.TransactionService.GetPlayerStats(rangeType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

// @Summary Check the status of a transaction
// @Description Returns the status and details of a transaction by its reference ID
// @Tags Transaction
// @Accept json
// @Produce json
// @Param apiKey header string true "API Key"
// @Param referenceId path string true "Reference ID of the transaction"
// @Router /api/v1/transactions/{referenceId}/status [get]
func (c *TransactionController) CheckStatus(ctx *gin.Context) {
	referenceID := ctx.Param("referenceId")
	if referenceID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Reference ID is required"})
		return
	}

	apiKey := ctx.GetHeader("apiKey")
	if apiKey == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "API key is required"})
		return
	}

	transaction, err := c.TransactionService.CheckStatus(referenceID, apiKey)
	if err != nil {
		if err.Error() == "invalid API key" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err.Error() == fmt.Sprintf("transaction with reference %s not found", referenceID) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, transaction)
}

// GetTransactionsAndTokensByPhoneNumber godoc
// @Summary Get transactions and tokens by phone number
// @Description Retrieves all transactions and their associated tokens for a given phone number
// @Tags Transaction
// @Accept json
// @Produce json
// @Param phoneNumber query string true "Phone number of the customer" example:"250788205965"
// @Success 200 {array} models.TransactionTokenResponse
// @Failure 400 {object} object{error=string} "Bad Request - Phone number is missing"
// @Failure 500 {object} object{error=string} "Internal Server Error"
// @Router /api/v1/transactions/getTokenByPhoneNumber [get]
func (c *TransactionController) GetTransactionsAndTokensByPhoneNumber(ctx *gin.Context) {
	phoneNumber := ctx.Query("phoneNumber")
	if phoneNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Phone number is required"})
		return
	}

	data, err := c.TransactionService.GetTransactionsAndTokensByPhoneNumber(phoneNumber)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

// @Summary Get all tokens for a phone number
// @Description Retrieves ongoing, won, and expired tokens for a given phone number
// @Tags Transaction
// @Accept json
// @Produce json
// @Param phoneNumber path string true "Phone number" example:"250788205965"
// @Param productId query string false "Product ID" example:"prod123"
// @Router /api/v1/transactions/my-tokens/{phoneNumber} [get]
func (c *TransactionController) GetMyTokens(ctx *gin.Context) {
	phoneNumber := ctx.Param("phoneNumber")
	if phoneNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Phone number is required"})
		return
	}

	productID := ctx.Query("productId")

	response, err := c.TransactionService.GetMyTokens(phoneNumber, productID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Get token statistics for a user
// @Description Retrieves token statistics grouped by product for a given phone number
// @Tags Transaction
// @Accept json
// @Produce json
// @Param phoneNumber query string true "Phone number of the user"
// @Router /api/v1/transactions/my-token-stats/{phoneNumber} [get]
func (c *TransactionController) GetTokenStats(ctx *gin.Context) {
	phoneNumber := ctx.Query("phoneNumber")
	if phoneNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Phone number is required"})
		return
	}

	stats, err := c.TransactionService.GetTokenStats(phoneNumber)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

// @Summary Get all tokens for a phone number
// @Description Retrieves ongoing, won, and expired tokens for a given phone number
// @Tags Transaction
// @Accept json
// @Produce json
// @Param phoneNumber path string true "Phone number" example:"250788205965"
// @Router /api/v1/transactions/all-tokens/{phoneNumber} [get]
func (c *TransactionController) GetAllMyTokens(ctx *gin.Context) {
	phoneNumber := ctx.Param("phoneNumber")
	if phoneNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Phone number is required"})
		return
	}

	productID := ctx.Query("productId")

	// Call service to fetch tokens
	response, err := c.TransactionService.GetMyTokens(phoneNumber, productID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Resend Token SMS
// @Description Resends a token via SMS for the given token ID
// @Tags Transaction
// @Accept json
// @Produce json
// @Param token path string true "Token ID"
// @Router /api/v1/transactions/resend-token/{token} [post]
func (c *TransactionController) ResendSms(ctx *gin.Context) {
	tokenID := ctx.Param("token")
	if tokenID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Token ID is required"})
		return
	}

	// Call the service function to resend SMS
	response, err := c.TransactionService.ResendSms(tokenID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Regenerate Token
// @Description Sends a request to an external API to regenerate a token
// @Tags Transaction
// @Accept json
// @Produce json
// @Param referenceId path string true "Reference ID"
// @Router /api/v1/transactions/regenerate/token/{referenceId} [post]
func (c *TransactionController) RegenerateToken(ctx *gin.Context) {
	referenceID := ctx.Param("referenceId")
	if referenceID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Reference ID is required"})
		return
	}

	// Call the service function
	response, err := c.TransactionService.RegenerateToken(referenceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Get transaction stats by phone number
// @Description Retrieves stats on transactions grouped by phone number
// @Tags Transaction
// @Accept json
// @Produce json
// @Param pageSize query int false "Page size (default: 10)"
// @Param pageNumber query int false "Page number (default: 1)"
// @Router /api/v1/transactions/phone_number-hits [get]
func (c *TransactionController) GetTransactionStats(ctx *gin.Context) {
	// Parse query parameters
	pageSize, err := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pageSize"})
		return
	}

	pageNumber, err := strconv.Atoi(ctx.DefaultQuery("pageNumber", "1"))
	if err != nil || pageNumber <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pageNumber"})
		return
	}

	// Call service function to get transaction stats
	response, err := c.TransactionService.GetTransactionStats(pageSize, pageNumber)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Get daily user transaction stats
// @Description Retrieves user transaction stats grouped by day
// @Tags Transaction
// @Accept json
// @Produce json
// @Param userId query string false "User ID"
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Router /api/v1/transactions/activators/stats [get]
func (c *TransactionController) GetDailyUserStats(ctx *gin.Context) {
	userId := ctx.Query("userId")
	from := ctx.Query("from")
	to := ctx.Query("to")

	stats, err := c.TransactionService.GetDailyUserStats(userId, from, to)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

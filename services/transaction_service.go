package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go-user-service/kafka"
	"go-user-service/models"
	"go-user-service/repository"
	"go-user-service/utils"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/Nerzal/gocloak/v13"
	"gorm.io/gorm"
)

type TransactionService struct {
	db               *gorm.DB
	repository       *repository.TransactionRepository
	accessKeyService *AccessKeyService
	eventHelper      *kafka.EventHelper
}

func NewTransactionService(db *gorm.DB, accessKeyService *AccessKeyService) *TransactionService {
	return &TransactionService{
		db:               db,
		repository:       repository.NewTransactionRepository(db),
		accessKeyService: accessKeyService,
		// eventHelper: kafka.NewEventHelper,
	}
}

// services/transaction_service.go (updated Initialize method)
func (s *TransactionService) Initialize(dto models.CreateTransactionDto, apiKey string) (map[string]interface{}, error) {
	// Validate the access key/API key
	accessKey, err := s.accessKeyService.FindOneAndCache(apiKey)
	if err != nil {
		return nil, err
	}
	if accessKey == nil {
		return nil, errors.New("invalid API key")
	}

	// Verify if Transaction ID is already used
	var existingTransaction models.Transaction
	result := s.db.Where("reference_id = ?", dto.ReferenceID).First(&existingTransaction)
	if result.Error == nil { // Found a transaction with the same reference_id
		return nil, errors.New("transaction ID already used")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	// Create the transaction
	transaction := models.Transaction{
		ReferenceID: dto.ReferenceID,
		PhoneNumber: dto.PhoneNumber,
		Amount:      dto.Amount,
		FinalAmount: dto.Amount, // For now, just use the same amount
		Telco:       (*models.Telco)(&dto.Telco),
		ProductID:   uuid.MustParse(dto.ProductID),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save the transaction
	if result := s.db.Create(&transaction); result.Error != nil {
		return nil, result.Error
	}

	return map[string]interface{}{
		"statusCode": 200,
		"message":    "Transaction request received successfully!",
	}, nil
}

func (s *TransactionService) Payment(paymentDto models.PaymentDto) (map[string]interface{}, error) {
	// Find the product
	var product models.Product
	if err := s.db.Where("product_id = ?", paymentDto.ProductID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	// Prepare event data
	eventData := map[string]interface{}{
		"productId":           product.ProductID,
		"amount":              product.PlayAmount,
		"phoneNumber":         paymentDto.PhoneNumber,
		"collectionAccountId": 2,
		"instance":            "",
		"userId":              paymentDto.AgentCode,
	}

	// Send the event
	err := s.eventHelper.SendEvent("Check transaction status", eventData, os.Getenv("TOPIC"))
	if err != nil {
		return nil, fmt.Errorf("failed to send event: %w", err)
	}

	return map[string]interface{}{
		"message": "Transaction request received successfully!",
	}, nil
}

func (s *TransactionService) FindAll(query utils.PaginationQuery) (*utils.PaginationResponse, error) {
	transactions, total, err := s.repository.FindAll(query)
	if err != nil {
		return nil, err
	}

	response := utils.CreatePaginationResponse(transactions, total, query)
	return &response, nil
}
func (s *TransactionService) GetProductStats(from *time.Time, to *time.Time, rangeType string) ([]models.ProductStats, error) {
	// Build base query with quoted table names and correct column names
	query := s.db.Table("\"transactions\"").
		Select("\"Products\".*, SUM(\"transactions\".final_amount) as totalamount").
		Joins("INNER JOIN \"Products\" ON \"transactions\".product_id = \"Products\".\"productId\"").
		Where("\"transactions\".status = ?", models.TransactionStatusSuccess).
		Group("\"Products\".\"productId\"")

	// Apply date filters based on range or explicit from/to
	if rangeType != "" {
		now := time.Now()
		var start, end time.Time

		switch rangeType {
		case "daily":
			start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			end = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
		case "weekly":
			// Calculate start of week (Sunday)
			daysToSunday := int(now.Weekday())
			start = time.Date(now.Year(), now.Month(), now.Day()-daysToSunday, 0, 0, 0, 0, now.Location())
			end = time.Date(now.Year(), now.Month(), now.Day()-daysToSunday+6, 23, 59, 59, 999999999, now.Location())
		case "monthly":
			start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			end = time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 999999999, now.Location())
		}

		query = query.Where("\"transactions\".created_at BETWEEN ? AND ?", start, end)
	} else if from != nil && to != nil {
		query = query.Where("\"transactions\".created_at BETWEEN ? AND ?", from, to)
	}

	// Execute the query
	var rawResults []map[string]interface{}
	if err := query.Find(&rawResults).Error; err != nil {
		return nil, err
	}

	// Calculate total amount for percentage calculations
	var totalAmount float64
	for _, result := range rawResults {
		if amount, ok := result["totalamount"].(float64); ok {
			totalAmount += amount
		}
	}

	// Process results to calculate percentages and margins
	var productStats []models.ProductStats
	for _, result := range rawResults {
		var productStat models.ProductStats

		// Map product data - with safer type conversions
		product := models.Product{
			ProductID: func() uuid.UUID {
				if id, ok := result["productId"].(string); ok {
					parsedID, _ := uuid.Parse(id)
					return parsedID
				}
				return uuid.Nil
			}(),
			// Handle integer types safely
			ProductIncrementer: safeIntFromInterface(result["productIcrementer"]),
			ProductName:        StringPtr(safeCastString(result["productName"])),
			ProductPicture:     StringPtr(safeCastString(result["productPicture"])),
			Description:        StringPtr(safeCastString(result["description"])),
			IsAvailable:        BoolPtr(safeCastBool(result["isAvailable"])),
			IsCallNeeded:       BoolPtr(safeCastBool(result["isCallNeeded"])),
			ProductCost:        IntPtr(safeIntFromInterface(result["productCost"])),
			DrawPeriod:         IntPtr(safeIntFromInterface(result["drawPeriod"])),
			ProductMargin:      IntPtr(safeIntFromInterface(result["product_margin"])),
			ExpectedAmount:     IntPtr(safeIntFromInterface(result["expected_amount"])),
			NumberOfWinners:    IntPtr(safeIntFromInterface(result["numberOfWinners"])),
			PlayAmount:         IntPtr(safeIntFromInterface(result["playAmount"])),
			CreatedAt:          result["createdAt"].(time.Time),
			UpdatedAt:          result["updatedAt"].(time.Time),
		}

		productStat.Product = product
		productStat.TotalAmount = safeCastFloat64(result["totalamount"])

		// Calculate percentage
		if totalAmount > 0 {
			productStat.Percentage = (productStat.TotalAmount / totalAmount) * 100
			// Round to 2 decimal places
			productStat.Percentage = float64(int(productStat.Percentage*100)) / 100
		}

		// Calculate margin
		// Count tokens for this product within the date range
		var tokenCount int64
		tokenQuery := s.db.Table("\"Tokens\"").Where("product_id = ?", product.ProductID) // Note the quoted table name
		if from != nil && to != nil {
			tokenQuery = tokenQuery.Where("created_at BETWEEN ? AND ?", from, to)
		}
		tokenQuery.Count(&tokenCount)

		if product.ExpectedAmount != nil && *product.ExpectedAmount > 0 && tokenCount > 0 {
			productStat.Margin = ((float64(tokenCount) * float64(*product.PlayAmount)) / float64(*product.ExpectedAmount)) * 100
			// Round to 4 decimal places
			productStat.Margin = float64(int(productStat.Margin*10000)) / 10000
		}

		productStats = append(productStats, productStat)
	}

	// Sort by percentage in descending order
	utils.SortProductStatsByPercentage(productStats)

	// Normalize percentages to ensure they sum to 100
	totalPercentage := 0.0
	for _, stat := range productStats {
		totalPercentage += stat.Percentage
	}

	if totalPercentage != 100 && totalPercentage > 0 {
		for i := range productStats {
			productStats[i].Percentage = (productStats[i].Percentage / totalPercentage) * 100
			// Round to 2 decimal places
			productStats[i].Percentage = float64(int(productStats[i].Percentage*100)) / 100
		}
	}

	return productStats, nil
}

// Helper functions for safe type conversion

func IntPtr(i int) *int {
	return &i
}

func Float64Ptr(f float64) *float64 {
	return &f
}
func BoolPtr(b bool) *bool {
	return &b
}

func safeIntFromInterface(value interface{}) int {
	switch v := value.(type) {
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return 0
	}
}

func safeCastString(value interface{}) string {
	if value == nil {
		return ""
	}
	str, ok := value.(string)
	if !ok {
		return ""
	}
	return str
}

func safeCastBool(value interface{}) bool {
	if value == nil {
		return false
	}
	b, ok := value.(bool)
	if !ok {
		return false
	}
	return b
}

func safeCastFloat64(value interface{}) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	default:
		return 0
	}
}

// Updated method for getting player stats with correct type conversion
func (s *TransactionService) GetPlayerStats(rangeType string) ([]models.PlayerStats, error) {
	// Build base query with quoted table names and correct column names
	query := s.db.Table("\"transactions\"").
		Select("\"Products\".*, COUNT(DISTINCT \"transactions_new\".phone_number) as numberOfPlayers").
		Joins("INNER JOIN \"Products\" ON \"transactions_new\".product_id = \"Products\".\"productId\"").
		Where("\"transactions_new\".status = ?", models.TransactionStatusSuccess).
		Group("\"Products\".\"productId\"")

	// Apply date filters based on range
	if rangeType != "" {
		now := time.Now()
		var start, end time.Time

		switch rangeType {
		case "daily":
			start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			end = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
		case "weekly":
			// Calculate start of week (Sunday)
			daysToSunday := int(now.Weekday())
			start = time.Date(now.Year(), now.Month(), now.Day()-daysToSunday, 0, 0, 0, 0, now.Location())
			end = time.Date(now.Year(), now.Month(), now.Day()-daysToSunday+6, 23, 59, 59, 999999999, now.Location())
		case "monthly":
			start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			end = time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 999999999, now.Location())
		}

		query = query.Where("\"transactions_new\".created_at BETWEEN ? AND ?", start, end)
	}

	// Execute the query
	var rawResults []map[string]interface{}
	if err := query.Find(&rawResults).Error; err != nil {
		return nil, err
	}

	// Process results
	var playerStats []models.PlayerStats
	for _, result := range rawResults {
		var playerStat models.PlayerStats

		// Map product data - with safer type conversions
		product := models.Product{
			ProductID: func() uuid.UUID {
				if id, ok := result["productId"].(string); ok {
					parsedID, _ := uuid.Parse(id)
					return parsedID
				}
				return uuid.Nil
			}(),
			// Handle integer types safely
			ProductIncrementer: safeIntFromInterface(result["productIcrementer"]),
			ProductName:        StringPtr(safeCastString(result["productName"])),
			ProductPicture:     StringPtr(safeCastString(result["productPicture"])),
			Description:        StringPtr(safeCastString(result["description"])),
			IsAvailable:        BoolPtr(safeCastBool(result["isAvailable"])),
			IsCallNeeded:       BoolPtr(safeCastBool(result["isCallNeeded"])),
			ProductCost:        IntPtr(safeIntFromInterface(result["productCost"])),
			DrawPeriod:         IntPtr(safeIntFromInterface(result["drawPeriod"])),
			ProductMargin:      IntPtr(safeIntFromInterface(result["product_margin"])),
			ExpectedAmount:     IntPtr(safeIntFromInterface(result["expected_amount"])),
			NumberOfWinners:    IntPtr(safeIntFromInterface(result["numberOfWinners"])),
			PlayAmount:         IntPtr(safeIntFromInterface(result["playAmount"])),
			CreatedAt:          result["createdAt"].(time.Time),
			UpdatedAt:          result["updatedAt"].(time.Time),
		}

		playerStat.Product = product
		playerStat.NumberOfPlayers = safeIntFromInterface(result["numberOfPlayers"])

		playerStats = append(playerStats, playerStat)
	}

	return playerStats, nil
}

func (s *TransactionService) CheckStatus(referenceID string, apiKey string) (*models.Transaction, error) {
	// Validate the access key/API key
	accessKey, err := s.accessKeyService.FindOneAndCache(apiKey)
	if err != nil {
		return nil, err
	}
	if accessKey == nil {
		return nil, errors.New("invalid API key")
	}

	// Find the transaction by reference ID
	var transaction models.Transaction
	result := s.db.Where("reference_id = ?", referenceID).
		Preload("Product").
		First(&transaction)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("transaction with reference %s not found", referenceID)
		}
		return nil, result.Error
	}

	return &transaction, nil
}

func (s *TransactionService) GetTransactionsAndTokensByPhoneNumber(phoneNumber string) ([]models.TransactionTokenResponse, error) {
	// Query transactions based on phone number
	var transactions []models.Transaction
	if err := s.db.Where("phone_number = ?", phoneNumber).
		Preload("Product").
		Find(&transactions).Error; err != nil {
		return nil, err
	}

	// Exit early if no transactions found
	if len(transactions) == 0 {
		return []models.TransactionTokenResponse{}, nil
	}

	// Extract transaction reference IDs
	referenceIDs := make([]string, len(transactions))
	for i, transaction := range transactions {
		referenceIDs[i] = transaction.ReferenceID
	}

	// Query tokens based on transaction reference IDs with proper table name case
	var tokens []models.Token
	if err := s.db.Table("\"Tokens\""). // Use explicit table name with quotes
						Where("\"referenceId\" IN ?", referenceIDs).
						Find(&tokens).Error; err != nil {
		return nil, err
	}

	// Create a map of reference IDs to tokens for easier lookup
	tokenMap := make(map[string]*models.Token)
	for i := range tokens {
		tokenMap[tokens[i].ReferenceID] = &tokens[i]
	}

	// Combine transactions with their corresponding tokens
	result := make([]models.TransactionTokenResponse, len(transactions))
	for i, transaction := range transactions {
		result[i] = models.TransactionTokenResponse{
			Transaction: &transactions[i],
			Token:       tokenMap[transaction.ReferenceID],
		}
	}

	return result, nil
}

func (s *TransactionService) GetMyTokens(phoneNumber string, productID string) (*models.APIResponse, error) {
	// Query transactions based on phone number and status
	var transactions []models.Transaction
	query := s.db.Where("phone_number = ? AND status = ?", phoneNumber, models.TransactionStatusSuccess)
	if err := query.Find(&transactions).Error; err != nil {
		return nil, err
	}

	// Extract transaction reference IDs
	referenceIDs := make([]string, len(transactions))
	for i, transaction := range transactions {
		referenceIDs[i] = transaction.ReferenceID
	}

	// If no transactions found, return empty response
	if len(referenceIDs) == 0 {
		return &models.APIResponse{
			Status: "success",
			Data: models.TokenResponse{
				OngoingTokens: []*models.Token{},
				WonTokens:     []*models.WonToken{},
				ExpiredTokens: []*models.TokenArchives{},
			},
		}, nil
	}

	// Query tokens based on reference IDs and optional product ID
	var tokens []*models.Token
	tokenQuery := s.db.Where("\"referenceId\" IN ?", referenceIDs).
		Preload("Product")

	if productID != "" {
		tokenQuery = tokenQuery.Where("\"productId\" = ?", productID)
	}

	if err := tokenQuery.Find(&tokens).Error; err != nil {
		return nil, err
	}

	// Query won tokens
	var wonTokens []*models.WonToken
	if err := s.db.Where("reference_id IN ?", referenceIDs).
		Preload("Product").
		Find(&wonTokens).Error; err != nil {
		return nil, err
	}

	// Query archived tokens
	var archivedTokens []*models.TokenArchives
	if err := s.db.Where("reference_id IN ?", referenceIDs).
		Preload("Product").
		Find(&archivedTokens).Error; err != nil {
		return nil, err
	}

	return &models.APIResponse{
		Status: "success",
		Data: models.TokenResponse{
			OngoingTokens: tokens,
			WonTokens:     wonTokens,
			ExpiredTokens: archivedTokens,
		},
	}, nil
}

func (s *TransactionService) GetTokenStats(phoneNumber string) ([]*models.TokenStats, error) {
	// Query transactions based on phone number and success status
	var transactions []models.Transaction
	if err := s.db.Where("phone_number = ? AND status = ?", phoneNumber, models.TransactionStatusSuccess).
		Preload("Product").
		Find(&transactions).Error; err != nil {
		return nil, err
	}

	// Count tokens per product
	productCounts := make(map[string]int)
	for _, transaction := range transactions {
		if transaction.ProductID.String() != "" {
			productCounts[transaction.ProductID.String()]++
		}
	}

	// Get products with counts and latest draws
	var productsWithStats []*models.TokenStats
	for productID, count := range productCounts {
		var product models.Product
		if err := s.db.Where("product_id = ?", productID).First(&product).Error; err != nil {
			continue
		}

		// Get latest draw for the product
		var latestDraw models.Draw
		if err := s.db.Where("product_id = ?", productID).
			Order("created_at DESC").
			First(&latestDraw).Error; err != nil {
			// If no draw found, continue without it
			productsWithStats = append(productsWithStats, &models.TokenStats{
				Product:     &product,
				NbrOfTokens: count,
			})
			continue
		}

		productsWithStats = append(productsWithStats, &models.TokenStats{
			Product:     &product,
			Draw:        &latestDraw,
			NbrOfTokens: count,
		})
	}

	return productsWithStats, nil
}

func (s *TransactionService) ResendSms(tokenID string) (map[string]string, error) {
	// Retrieve the token by token ID
	var playerToken models.Token
	if err := s.db.Where("tokenId = ?", tokenID).Preload("Draw").First(&playerToken).Error; err != nil {
		return nil, fmt.Errorf("token not found with tokenId: %s", tokenID)
	}

	// Retrieve tokenMillion (assuming it might have a different type of token)
	var tokenMillion models.Token
	if err := s.db.Where("tokenId = ?", tokenID).Preload("Draw").First(&tokenMillion).Error; err != nil {
		return nil, fmt.Errorf("token million not found with tokenId: %s", tokenID)
	}

	// Retrieve the transaction linked to this token
	var transaction models.Transaction
	if err := s.db.Where("reference_id = ?", playerToken.ReferenceID).
		Preload("Product").
		First(&transaction).Error; err != nil {
		return nil, fmt.Errorf("transaction not found for reference ID: %s", playerToken.ReferenceID)
	}

	// Compute the draw end dates
	// endDate := playerToken.Draw.EndDate.AddDate(0, 0, 1).Format("2006-01-02")
	// fmt.Printf(endDate)
	endDateMill := tokenMillion.Draw.EndDate.AddDate(0, 0, 1).Format("2006-01-02")

	// Construct the SMS message
	message := fmt.Sprintf("IKUBIRE Lotto ya %s Token ni: %s Tombola izaba: %s 19:00 kuri RTV, Umutingito w'ibihembo!!!",
		transaction.Product.ProductName, playerToken.TokenID, endDateMill)

	// Construct the event data
	eventData := map[string]interface{}{
		"action": "send-message",
		"data": map[string]interface{}{
			"message":     message,
			"phoneNumber": transaction.PhoneNumber,
			"token":       "",
			"userId":      "",
		},
	}

	// Send the event to Kafka
	err := s.eventHelper.SendEvent("resend Token", eventData, os.Getenv("SMS_TOPIC"))
	if err != nil {
		return nil, fmt.Errorf("failed to send event: %w", err)
	}

	return map[string]string{"message": "Message sent successfully"}, nil
}

func (s *TransactionService) RegenerateToken(referenceID string) (map[string]interface{}, error) {
	apiURL := fmt.Sprintf("http://167.86.82.192:8027/api/v1/regenerate-token/%s", referenceID)

	// Create HTTP request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(nil))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to external API: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external API error: %v", responseBody)
	}

	return responseBody, nil
}

func (s *TransactionService) GetTransactionStats(pageSize int, pageNumber int) (map[string]interface{}, error) {
	var results []struct {
		PhoneNumber  string `json:"phone_number"`
		TotalHits    int    `json:"totalHits"`
		TotalSuccess int    `json:"totalSuccess"`
		TotalFailed  int    `json:"totalFailed"`
	}

	// Query: Count transactions per phone number
	query := s.db.Table("\"transactions_new\"").
		Select("\"phone_number\", COUNT(\"phone_number\") as totalHits, "+
			"SUM(CASE WHEN \"status\" = ? THEN 1 ELSE 0 END) as totalSuccess",
			models.TransactionStatusSuccess).
		Group("\"phone_number\"").
		Order("totalHits DESC").
		Limit(pageSize).
		Offset((pageNumber - 1) * pageSize)

	if err := query.Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch transaction stats: %w", err)
	}

	// Compute total unique phone numbers
	var uniquePhoneNumbers int64
	if err := s.db.Table("\"transactions_new\"").
		Select("COUNT(DISTINCT \"phone_number\")").
		Count(&uniquePhoneNumbers).Error; err != nil {
		return nil, fmt.Errorf("failed to count unique phone numbers: %w", err)
	}

	// Compute total failed transactions
	for i := range results {
		results[i].TotalFailed = results[i].TotalHits - results[i].TotalSuccess
	}

	// Pagination calculations
	lastPage := int(math.Ceil(float64(uniquePhoneNumbers) / float64(pageSize)))
	nextPage := pageNumber + 1
	if nextPage > lastPage {
		nextPage = 0
	}
	prevPage := pageNumber - 1
	if prevPage < 1 {
		prevPage = 0
	}

	// Response structure
	response := map[string]interface{}{
		"status":       "success",
		"list":         results,
		"total":        uniquePhoneNumbers,
		"previousPage": prevPage,
		"nextPage":     nextPage,
		"lastPage":     lastPage,
		"currentPage":  pageNumber,
	}

	return response, nil
}

func StringPtr(s string) *string {
	return &s
}

func (s *TransactionService) GetDailyUserStats(userId string, from string, to string) ([]map[string]interface{}, error) {
	var results []struct {
		UserID           string   `json:"userID"`
		Date             string   `json:"date"`
		TotalTransaction int      `json:"total_transaction"`
		TotalAmount      *float64 `json:"total_amount"` // Use pointer to handle NULL values
	}

	// Build base query
	query := s.db.Table("\"transactions_new\"").
		Select("\"user_id\" AS userID, DATE(\"created_at\" + interval '2 hours') AS date, COUNT(\"id\") AS total_transaction, COALESCE(SUM(\"final_amount\"), 0) AS total_amount"). // ✅ Fix: Use COALESCE() to replace NULL with 0
		Where("\"user_id\" IS NOT NULL AND \"user_id\" != ''").
		Where("\"status\" = ?", models.TransactionStatusSuccess).
		Group("\"user_id\", DATE(\"created_at\" + interval '2 hours')")

	// Apply filters only if parameters are provided
	if userId != "" {
		query = query.Where("\"user_id\" = ?", userId)
	}
	if from != "" && to != "" {
		query = query.Where("\"created_at\" BETWEEN ? AND ?", from, to)
	} else if from != "" {
		query = query.Where("\"created_at\" >= ?", from)
	} else if to != "" {
		query = query.Where("\"created_at\" <= ?", to)
	}

	// Execute query and check errors
	if err := query.Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch user transaction stats: %w", err)
	}

	// Handle empty results (no transactions found)
	if len(results) == 0 {
		return []map[string]interface{}{}, nil // ✅ Return an empty array instead of nil to prevent panics
	}

	// Fetch user details from Keycloak
	keycloakService, err := NewKeycloakService()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Keycloak: %w", err)
	}

	// Helper function for safe string conversion
	StringPtr := func(s string) *string { return &s }

	// Enhance results with Keycloak usernames
	var enrichedResults []map[string]interface{}
	for _, result := range results {
		user, err := keycloakService.GetUserDetails(result.UserID)
		if err != nil {
			// ✅ Fix: Use safe pointers to avoid panic
			user = &gocloak.User{
				Username:  StringPtr("Unknown"),
				FirstName: StringPtr(""),
				LastName:  StringPtr(""),
			}
		}

		enrichedResults = append(enrichedResults, map[string]interface{}{
			"userID":            result.UserID,
			"username":          *user.Username,
			"names":             fmt.Sprintf("%s %s", *user.FirstName, *user.LastName),
			"date":              result.Date,
			"total_transaction": result.TotalTransaction,
			"total_amount":      *result.TotalAmount, // ✅ Fix: Dereference pointer safely
		})
	}

	return enrichedResults, nil
}

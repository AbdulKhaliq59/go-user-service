package repository

import (
	"go-user-service/models"
	"go-user-service/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (r *TransactionRepository) FindAll(query utils.PaginationQuery) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	// Start building the query (without JOIN)
	db := r.db.Model(&models.Transaction{})

	// Apply filters
	if query.From != "" {
		fromDate, err := time.Parse(time.RFC3339, query.From)
		if err == nil {
			db = db.Where("created_at >= ?", fromDate)
		}
	}

	if query.To != "" {
		toDate, err := time.Parse(time.RFC3339, query.To)
		if err == nil {
			toDate = toDate.Add(24*time.Hour - time.Second)
			db = db.Where("created_at <= ?", toDate)
		}
	}

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	if query.Telecom != "" {
		db = db.Where("telco = ?", query.Telecom)
	}

	if query.ProductID != "" {
		db = db.Where("product_id = ?", query.ProductID)
	}

	if query.AgentCode != "" {
		db = db.Where("user_id = ?", query.AgentCode)
	}

	if query.Search != "" {
		search := "%" + strings.ToLower(query.Search) + "%"
		db = db.Where("LOWER(phone_number) LIKE ?", search)
	}

	// Count total records
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination defaults
	if query.PageSize == 0 {
		query.PageSize = 10
	}
	if query.PageNumber == 0 {
		query.PageNumber = 1
	}

	offset := (query.PageNumber - 1) * query.PageSize

	// Get paginated transactions (without JOIN)
	err := db.
		Order("created_at DESC").
		Offset(offset).
		Limit(query.PageSize).
		Find(&transactions).Error

	return transactions, total, err
}

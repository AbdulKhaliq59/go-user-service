// models/transaction.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Telco string

const (
	MTN    Telco = "MTN"
	AIRTEL Telco = "AIRTEL"
)

type TransactionStatus string

const (
	StatusPending TransactionStatus = "PENDING"
	StatusFailed  TransactionStatus = "FAILED"
	StatusSuccess TransactionStatus = "SUCCESS"
)

type Transaction struct {
	ID             string            `gorm:"primaryKey;type:uuid" json:"id"`
	ReferenceID    string            `gorm:"column:reference_id;not null;unique" json:"reference_id"`
	PhoneNumber    string            `gorm:"column:phone_number;not null" json:"phone_number"`
	Amount         float64           `gorm:"not null" json:"amount"`
	DiscountCode   string            `json:"discount_code"`
	DiscountType   string            `json:"discount_type"`
	DiscountAmount float64           `json:"discount_amount"`
	FinalAmount    float64           `gorm:"column:final_amount;not null" json:"final_amount"`
	ReferrerID     string            `gorm:"column:referrer_id" json:"referrer_id"`
	Description    string            `json:"description"`
	ProductID      string            `gorm:"column:product_id" json:"product_id"`
	Product        *Product          `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	UserID         string            `gorm:"column:user_id" json:"user_id"`
	Status         TransactionStatus `gorm:"default:PENDING" json:"status"`
	StatusMessage  string            `gorm:"column:status_message" json:"status_message"`
	GwRef          string            `gorm:"column:gw_ref" json:"gw_ref"`
	ChannelRef     string            `gorm:"column:channel_ref" json:"channel_ref"`
	Telco          string            `json:"telco"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	if t.Status == "" {
		t.Status = StatusPending
	}
	return nil
}
func (Transaction) TableName() string {
	return "transactions_new"
}

type CreateTransactionDto struct {
	ReferenceID     string    `json:"reference_id" binding:"required" example:"string"`
	PhoneNumber     string    `json:"phone_number" binding:"required" example:"250788205965"`
	Telco           Telco     `json:"telco" binding:"required" example:"string"`
	Amount          float64   `json:"amount" binding:"required,min=1" example:"0"`
	Token           string    `json:"token" example:"string"`
	TransactionDate time.Time `json:"transaction_date" binding:"required" example:"2025-02-18T08:45:19.803Z"`
	ProductID       string    `json:"productId" binding:"required" example:"string"`
}

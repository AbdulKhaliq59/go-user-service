// models/transaction.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionStatus string
type Telco string

const (
	TransactionStatusPending TransactionStatus = "PENDING"
	TransactionStatusFailed  TransactionStatus = "FAILED"
	TransactionStatusSuccess TransactionStatus = "SUCCESS"

	TelcoTNM    Telco = "TNM"
	TelcoAirtel Telco = "AIRTEL"
)

type Transaction struct {
	ID             uuid.UUID         `gorm:"type:uuid;primary_key;column:id" json:"id"`
	ReferenceID    string            `gorm:"column:reference_id;unique;not null" json:"reference_id"`
	PhoneNumber    string            `gorm:"column:phone_number;not null" json:"phone_number"`
	Amount         float64           `gorm:"column:amount;not null" json:"amount"`
	DiscountCode   *string           `gorm:"column:discount_code" json:"discount_code"`
	DiscountType   *string           `gorm:"column:discount_type" json:"discount_type"`
	DiscountAmount *float64          `gorm:"column:discount_amount" json:"discount_amount"`
	FinalAmount    float64           `gorm:"column:final_amount;not null" json:"final_amount"`
	ReferrerID     *string           `gorm:"column:referrer_id" json:"referrer_id"`
	Description    *string           `gorm:"column:description" json:"description"`
	ProductID      uuid.UUID         `gorm:"type:uuid;column:product_id" json:"product_id"`
	UserID         *string           `gorm:"column:user_id" json:"user_id"`
	Status         TransactionStatus `gorm:"column:status;default:PENDING" json:"status"`
	StatusMessage  *string           `gorm:"column:status_message" json:"status_message"`
	GwRef          *string           `gorm:"column:gw_ref" json:"gw_ref"`
	ChannelRef     *string           `gorm:"column:channel_ref" json:"channel_ref"`
	Telco          *Telco            `gorm:"column:telco" json:"telco"`
	CreatedAt      time.Time         `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time         `gorm:"column:updated_at" json:"updated_at"`

	// Relations
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// TableName specifies the table name for the Transaction model
func (Transaction) TableName() string {
	return "transactions"
}

// BeforeCreate will set a UUID if it hasn't been set
func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return
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

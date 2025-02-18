package models

import "time"

// Token represents the token model in the database
type Token struct {
	TokenID     int64     `gorm:"primaryKey;column:tokenId"`
	ReferenceID string    `gorm:"column:referenceId"`
	ProductID   string    `gorm:"column:productId"`
	DrawID      string    `gorm:"column:drawId"`
	Draw        Draw      `gorm:"foreignKey:DrawID"`
	CreatedAt   time.Time `gorm:"column:createdAt"`
	UpdatedAt   time.Time `gorm:"column:updatedAt"`
	UserID      string    `gorm:"column:userId"`
}

func (Token) TableName() string {
	return "Tokens"
}

type TransactionTokenResponse struct {
	Transaction *Transaction `json:"transaction"`
	Token       *Token       `json:"token,omitempty"`
}

// models/temp_winner.go
package models

import (
	"time"

	"github.com/google/uuid"
)

type TempWinner struct {
	TokenID     int64     `gorm:"type:bigint;primary_key;column:token_id" json:"token_id"`
	ReferenceID string    `gorm:"column:reference_id" json:"reference_id"`
	PhoneNumber *string   `gorm:"column:phone_number" json:"phone_number"`
	ProductID   uuid.UUID `gorm:"type:uuid;column:product_id" json:"product_id"`
	DrawID      uuid.UUID `gorm:"type:uuid;column:draw_id" json:"draw_id"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`

	// Relations
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Draw    *Draw    `gorm:"foreignKey:DrawID" json:"draw,omitempty"`
}

// TableName specifies the table name for the TempWinner model
func (TempWinner) TableName() string {
	return "temp_winner"
}

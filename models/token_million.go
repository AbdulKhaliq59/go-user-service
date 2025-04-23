// models/token_million.go
package models

import (
	"time"

	"github.com/google/uuid"
)

type TokenMillion struct {
	TokenID     int64     `gorm:"type:bigint;primary_key;column:token_id" json:"token_id"`
	ReferenceID string    `gorm:"column:reference_id" json:"reference_id"`
	ProductID   uuid.UUID `gorm:"type:uuid;column:product_id" json:"product_id"`
	DrawID      uuid.UUID `gorm:"type:uuid;column:draw_id" json:"draw_id"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`

	// Relations
	Product   *Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Draw      *Draw      `gorm:"foreignKey:DrawID" json:"draw,omitempty"`
	WonTokens []WonToken `gorm:"foreignKey:TokenMillionID" json:"won_tokens,omitempty"`
}

// TableName specifies the table name for the TokenMillion model
func (TokenMillion) TableName() string {
	return "token_millions"
}

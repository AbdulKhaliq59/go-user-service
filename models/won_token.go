// models/won_token.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WonToken struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;column:id" json:"id"`
	TokenID        string     `gorm:"column:token_id" json:"token_id"`
	ReferenceID    string     `gorm:"column:reference_id" json:"reference_id"`
	ProductID      uuid.UUID  `gorm:"type:uuid;column:product_id" json:"product_id"`
	DrawID         uuid.UUID  `gorm:"type:uuid;column:draw_id" json:"draw_id"`
	TokenMillionID *int64     `gorm:"column:token_million_id" json:"token_million_id"`
	CreatedBy      *string    `gorm:"column:created_by" json:"created_by"`
	Name           *string    `gorm:"column:name" json:"name"`
	PhoneNumber    *string    `gorm:"column:phone_number" json:"phone_number"`
	Age            *int       `gorm:"column:age" json:"age"`
	Gender         *string    `gorm:"column:gender" json:"gender"`
	HeardFrom      *string    `gorm:"column:heard_from" json:"heard_from"`
	IDNumber       *string    `gorm:"column:id_number" json:"id_number"`
	HandedOver     *bool      `gorm:"column:handed_over" json:"handed_over"`
	ConfirmedBy    *string    `gorm:"column:confirmed_by" json:"confirmed_by"`
	ConfirmedAt    *time.Time `gorm:"column:confirmed_at" json:"confirmed_at"`
	IsConfirmed    *bool      `gorm:"column:is_confirmed" json:"is_confirmed"`
	Occupation     *string    `gorm:"column:occupation" json:"occupation"`
	Signature      *string    `gorm:"column:signature" json:"signature"`
	HandedOverDate *time.Time `gorm:"column:handed_over_date" json:"handed_over_date"`
	HandedOverBy   *string    `gorm:"column:handed_over_by" json:"handed_over_by"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"`

	// Relations
	Product      *Product      `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Token        *Token        `gorm:"foreignKey:TokenID" json:"token,omitempty"`
	TokenMillion *TokenMillion `gorm:"foreignKey:TokenMillionID" json:"token_million,omitempty"`
	Draw         *Draw         `gorm:"foreignKey:DrawID" json:"draw,omitempty"`
}

// TableName specifies the table name for the WonToken model
func (WonToken) TableName() string {
	return "wonTokens"
}

// BeforeCreate will set a UUID if it hasn't been set
func (w *WonToken) BeforeCreate(tx *gorm.DB) (err error) {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return
}

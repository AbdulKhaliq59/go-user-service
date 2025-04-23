// models/draw.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Draw struct {
	DrawID    uuid.UUID `gorm:"type:uuid;primary_key;column:drawId" json:"drawId"`
	ProductID uuid.UUID `gorm:"type:uuid;column:productId" json:"productId"`
	StartDate time.Time `gorm:"column:startDate" json:"startDate"`
	EndDate   time.Time `gorm:"column:endDate" json:"endDate"`
	IsPlayed  bool      `gorm:"column:isPlayed" json:"isPlayed"`
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
	// Relations
	Product       *Product       `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Tokens        []Token        `gorm:"foreignKey:DrawID" json:"tokens,omitempty"`
	TempWinners   []TempWinner   `gorm:"foreignKey:DrawID" json:"temp_winners,omitempty"`
	TokenMillions []TokenMillion `gorm:"foreignKey:DrawID" json:"token_millions,omitempty"`
	WonTokens     []WonToken     `gorm:"foreignKey:DrawID" json:"won_tokens,omitempty"`
}

// TableName specifies the table name for the Draw model
func (Draw) TableName() string {
	return "Draws"
}

// BeforeCreate will set a UUID if it hasn't been set
func (d *Draw) BeforeCreate(tx *gorm.DB) (err error) {
	if d.DrawID == uuid.Nil {
		d.DrawID = uuid.New()
	}
	return
}

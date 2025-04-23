package models

import (
	"time"

	"github.com/google/uuid"
)

type BonusProduct struct {
	ProductID      uuid.UUID `gorm:"column:productId;type:uuid;primaryKey"`
	BonusProductID uuid.UUID `gorm:"column:bonusProductId;type:uuid;primaryKey"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName sets the custom table name for the join table
func (BonusProduct) TableName() string {
	return "\"Product_Bonus\""
}

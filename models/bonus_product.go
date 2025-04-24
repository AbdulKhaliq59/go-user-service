package models

import (
	"time"

	"github.com/google/uuid"
)

type ProductBonus struct {
	ProductID      uuid.UUID `gorm:"column:productId;type:uuid;primaryKey"`
	BonusProductID uuid.UUID `gorm:"column:bonusProductId;type:uuid;primaryKey"`
	CreatedAt      time.Time `gorm:"column:createdAt"`
	UpdatedAt      time.Time `gorm:"column:updatedAt"`
}

func (ProductBonus) TableName() string {
	return "Product_Bonus"
}

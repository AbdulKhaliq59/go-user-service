package models

import (
	"time"
)

type Draw struct {
	DrawID    string    `gorm:"primaryKey;column:drawId" json:"drawId"`
	StartDate time.Time `gorm:"column:startDate" json:"startDate"`
	EndDate   time.Time `gorm:"column:endDate" json:"endDate"`
	IsPlayed  bool      `gorm:"column:isPlayed" json:"isPlayed"`
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
	ProductID string    `gorm:"column:productId" json:"productId"`
	Product   *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (Draw) TableName() string {
	return "Draws"
}

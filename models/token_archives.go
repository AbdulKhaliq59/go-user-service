package models

import "time"

type TokenArchives struct {
	ID          int64    `gorm:"primaryKey;column:id" json:"id"`
	ReferenceID string   `gorm:"column:reference_id" json:"reference_id"`
	ProductID   string   `gorm:"column:product_id" json:"product_id"`
	Product     *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	DrawID      string   `gorm:"column:draw_id" json:"draw_id"`
	// Draw        *Draw     `gorm:"foreignKey:DrawID" json:"draw,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (TokenArchives) TableName() string {
	return "tokenArchives"
}

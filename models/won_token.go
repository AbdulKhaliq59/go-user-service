package models

import (
	"time"
)

type WonToken struct {
	ID          string   `gorm:"primaryKey;type:uuid" json:"id"`
	TokenID     string   `gorm:"column:token_id" json:"token_id"`
	ReferenceID string   `gorm:"column:reference_id" json:"reference_id"`
	ProductID   string   `gorm:"column:product_id" json:"product_id"`
	Product     *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Token       *Token   `gorm:"foreignKey:TokenID" json:"token,omitempty"`
	DrawID      string   `gorm:"column:draw_id" json:"draw_id"`
	// Draw           *Draw     `gorm:"foreignKey:DrawID" json:"draw,omitempty"`
	CreatedBy   string `json:"created_by"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	LocationID  string `gorm:"column:location_id" json:"location_id"`
	// Location       *Location `gorm:"foreignKey:LocationID" json:"location,omitempty"`
	HeardFrom      string    `gorm:"column:heard_from" json:"heard_from"`
	IDNumber       string    `gorm:"column:id_number" json:"id_number"`
	HandedOver     bool      `json:"handed_over"`
	ConfirmedBy    string    `json:"confirmed_by"`
	ConfirmedAt    time.Time `json:"confirmed_at"`
	IsConfirmed    bool      `json:"is_confirmed"`
	Occupation     string    `json:"occupation"`
	Signature      string    `json:"signature"`
	HandedOverDate time.Time `json:"handed_over_date"`
	HandedOverBy   string    `json:"handed_over_by"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (WonToken) TableName() string {
	return "wonTokens"
}

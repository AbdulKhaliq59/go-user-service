// models/product.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ProductID          string    `gorm:"primaryKey;column:product_id" json:"id"`
	ProductIncrementer int       `gorm:"column:product_icrementer" json:"productIcrementer"`
	ProductName        string    `gorm:"column:product_name" json:"name"`
	ProductPicture     string    `gorm:"column:product_picture" json:"picture"`
	Description        string    `gorm:"column:description" json:"description"`
	IsAvailable        bool      `gorm:"column:is_available" json:"isAvailable"`
	IsCallNeeded       bool      `gorm:"column:is_call_needed" json:"isCallNeeded"`
	ProductCost        float64   `gorm:"column:product_cost" json:"productCost"`
	DrawPeriod         string    `gorm:"column:draw_period" json:"drawPeriod"`
	ProductMargin      float64   `gorm:"column:product_margin" json:"productMargin"`
	ExpectedAmount     float64   `gorm:"column:expected_amount" json:"expectedAmount"`
	NumberOfWinners    int       `gorm:"column:number_of_winners" json:"numberOfWinners"`
	PlayAmount         float64   `gorm:"column:play_amount" json:"playAmount"`
	CreatedAt          time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ProductID == "" {
		p.ProductID = uuid.New().String()
	}
	return nil
}

func (Product) TableName() string {
	return "\"Products\"" // This will preserve the case
}

type ProductStats struct {
	Product     Product `json:"product"`
	TotalAmount float64 `json:"totalamount"`
	Percentage  float64 `json:"percentage"`
	Margin      float64 `json:"margin"`
}

type PlayerStats struct {
	Product         Product `json:"product"`
	NumberOfPlayers int     `json:"numberOfPlayers"`
}

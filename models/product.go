package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ProductID          uuid.UUID `gorm:"type:uuid;primaryKey;column:productId" json:"productId"`
	ProductIncrementer int       `gorm:"unique;column:productIcrementer" json:"productIcrementer"`
	ProductName        *string   `gorm:"column:productName" json:"productName"`
	ProductPicture     *string   `gorm:"column:productPicture" json:"productPicture"`
	Description        *string   `gorm:"column:description" json:"description"`
	IsAvailable        *bool     `gorm:"column:isAvailable" json:"isAvailable"`
	IsCallNeeded       *bool     `gorm:"column:isCallNeeded" json:"isCallNeeded"`
	AvoidConflict      *bool     `gorm:"column:avoidConflict" json:"avoidConflict"`
	IsBonus            *bool     `gorm:"column:isBonus" json:"isBonus"`
	ProductCost        *int      `gorm:"column:productCost" json:"productCost"`
	DrawPeriod         *int      `gorm:"column:drawPeriod" json:"drawPeriod"`
	NumberOfWinners    *int      `gorm:"column:numberOfWinners" json:"numberOfWinners"`
	PlayAmount         *int      `gorm:"column:playAmount" json:"playAmount"`
	Priority           *int      `gorm:"column:priority" json:"priority"`
	Comment            *string   `gorm:"column:comment" json:"comment"`
	EnglishName        *string   `gorm:"column:englishName" json:"englishName"`
	ProductMargin      *int      `gorm:"column:product_margin" json:"product_margin"`
	ExpectedAmount     *int      `gorm:"column:expected_amount" json:"expected_amount"`
	RequiredDrawDays   []int     `gorm:"type:integer[];column:requiredDrawDays" json:"requiredDrawDays"`
	CreatedAt          time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt          time.Time `gorm:"column:updatedAt" json:"updatedAt"`

	// Self-referencing many-to-many
	BonusProducts     []*Product `gorm:"many2many:Product_Bonus;joinForeignKey:ProductID;joinReferences:BonusProductID" json:"bonusProducts,omitempty"`
	IncludedInBonusOf []*Product `gorm:"many2many:Product_Bonus;joinForeignKey:BonusProductID;joinReferences:ProductID" json:"includedInBonusOf,omitempty"`
}

// TableName overrides default GORM table name
func (Product) TableName() string {
	return "Products"
}

// BeforeCreate will set a UUID if it hasn't been set
func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ProductID == uuid.Nil {
		p.ProductID = uuid.New()
	}
	return
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

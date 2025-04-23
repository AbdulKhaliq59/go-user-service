// api/dto/product_dto.go
package dto

import (
	"time"

	"github.com/google/uuid"
)

// BonusProductDTO represents the request body for assigning or unassigning bonus products
type BonusProductDTO struct {
	BonusProductIds []string `json:"bonusProductIds" binding:"required" example:"['string']"`
}

// ProductUploadResponse represents the response for file upload
type ProductUploadResponse struct {
	FileLink string `json:"fileLink" example:"http://example.com/uploads/image.jpg"`
}

// ProductResponse represents a standard response for product operations
type ProductResponse struct {
	Status  string      `json:"status" example:"success"`
	Message string      `json:"message" example:"Product updated successfully"`
	Data    interface{} `json:"data,omitempty"`
}

type DrawResponse struct {
	DrawID    uuid.UUID `json:"drawId"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	IsPlayed  bool      `json:"isPlayed"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GetProductResponseDto struct {
	ProductID          string        `json:"productId"`
	ProductIncrementer int           `json:"productIcrementer"`
	ProductName        string        `json:"productName"`
	ProductPicture     string        `json:"productPicture"`
	Description        string        `json:"description"`
	IsAvailable        bool          `json:"isAvailable"`
	IsCallNeeded       bool          `json:"isCallNeeded"`
	AvoidConflict      bool          `json:"avoidConflict"`
	IsBonus            bool          `json:"isBonus"`
	ProductCost        int           `json:"productCost"`
	DrawPeriod         int           `json:"drawPeriod"`
	NumberOfWinners    int           `json:"numberOfWinners"`
	PlayAmount         int           `json:"playAmount"`
	Priority           int           `json:"priority"`
	Comment            *string       `json:"comment"`
	EnglishName        *string       `json:"englishName"`
	ProductMargin      int           `json:"product_margin"`
	ExpectedAmount     int           `json:"expected_amount"`
	CreatedAt          string        `json:"createdAt"`
	UpdatedAt          string        `json:"updatedAt"`
	Draw               *DrawResponse `json:"draw"`
}

type CreateProductRequest struct {
	ProductName     string  `json:"productName" binding:"required"`
	Description     string  `json:"description" binding:"required"`
	ProductPicture  string  `json:"productPicture"`
	IsAvailable     bool    `json:"isAvailable" binding:"required"`
	ProductMargin   int     `json:"product_margin"`
	IsCallNeeded    bool    `json:"isCallNeeded" binding:"required"`
	EnglishName     *string `json:"englishName"`
	ProductCost     *int    `json:"productCost" binding:"required"`
	DrawPeriod      *int    `json:"drawPeriod" binding:"required"`
	NumberOfWinners *int    `json:"numberOfWinners" binding:"required"`
	PlayAmount      *int    `json:"playAmount" binding:"required"`
	Priority        *int    `json:"priority"`
	IsBonus         bool    `json:"isBonus" binding:"required"`
}

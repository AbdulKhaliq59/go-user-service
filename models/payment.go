// models/payment.go
package models

type PaymentDto struct {
	ProductID   string `json:"productId" binding:"required" example:"string"`
	PhoneNumber string `json:"phoneNumber" binding:"required" example:"string"`
	AgentCode   string `json:"agentCode" example:"string"` // Optional field
}

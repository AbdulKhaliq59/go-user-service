// api/dto/blacklist_dto.go
package dto

type CreateBlacklistDto struct {
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
}

type UpdateBlacklistDto struct {
	Name        *string `json:"name,omitempty"`
	PhoneNumber *string `json:"phoneNumber,omitempty"`
}

// models/blacklist.go
package models

import (
	"time"
)

type Blacklist struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phoneNumber" gorm:"column:phone_number;unique"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

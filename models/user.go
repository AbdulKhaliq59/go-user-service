// models/user.go
package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"size:255;not null;unique" json:"username"`
	Email     string         `gorm:"size:255;not null;unique" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"` // Password not visible in JSON response
	FirstName string         `gorm:"size:255" json:"first_name"`
	LastName  string         `gorm:"size:255" json:"last_name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

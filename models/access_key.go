package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccessKey struct {
	ID        string          `gorm:"type:uuid;primary_key" json:"id"`
	AccessKey string          `gorm:"unique;not null" json:"access_key"`
	ExpireAt  *time.Time      `json:"expire_at"`
	User      json.RawMessage `gorm:"type:jsonb" json:"user"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func (a *AccessKey) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

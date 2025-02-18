package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"go-user-service/models"
	"go-user-service/redis"
	"time"

	"gorm.io/gorm"
)

type AccessKeyService struct {
	db    *gorm.DB
	redis *redis.RedisHelper
}

func NewAccessKeyService(db *gorm.DB, redis *redis.RedisHelper) *AccessKeyService {
	return &AccessKeyService{
		db:    db,
		redis: redis,
	}
}

func (s *AccessKeyService) GenerateAccessKey(length int) string {
	buffer := make([]byte, length/2)
	rand.Read(buffer)
	return hex.EncodeToString(buffer)[:length]
}

func (s *AccessKeyService) FindOneAndCache(apiKey string) (*models.AccessKey, error) {
	// Check if access key is cached in Redis
	var accessKey models.AccessKey
	err := s.redis.Get("access-key-"+apiKey, &accessKey)
	if err == nil {
		return &accessKey, nil
	}

	// Get access key from DB if not cached
	result := s.db.Where("access_key = ?", apiKey).First(&accessKey)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil without error if not found
		}
		return nil, result.Error
	}

	// Cache the access key
	s.redis.Set("access-key-"+apiKey, accessKey, 10*time.Minute)
	return &accessKey, nil
}

func (s *AccessKeyService) Create(expireAt *time.Time, user interface{}) (*models.AccessKey, error) {
	apiKey := s.GenerateAccessKey(30)
	accessKey := models.AccessKey{
		AccessKey: apiKey,
		ExpireAt:  expireAt,
	}

	// Convert user to JSON and assign
	// This part might need modification depending on how you're handling the JSON field

	result := s.db.Create(&accessKey)
	if result.Error != nil {
		return nil, result.Error
	}

	return &accessKey, nil
}

func (s *AccessKeyService) FindAll(userID string) ([]models.AccessKey, error) {
	var accessKeys []models.AccessKey
	result := s.db.Where("user->>'id' = ?", userID).Find(&accessKeys)
	if result.Error != nil {
		return nil, result.Error
	}
	return accessKeys, nil
}

// repository/blacklist_repository.go
package repository

import (
	"errors"
	"go-user-service/models"
	"go-user-service/utils"

	"gorm.io/gorm"
)

type BlacklistRepository struct {
	DB *gorm.DB
}

func NewBlacklistRepository(db *gorm.DB) *BlacklistRepository {
	return &BlacklistRepository{
		DB: db,
	}
}

func (r *BlacklistRepository) Create(blacklist *models.Blacklist) (*models.Blacklist, error) {
	if err := r.DB.Create(blacklist).Error; err != nil {
		return nil, err
	}
	return blacklist, nil
}

func (r *BlacklistRepository) FindAll(query utils.PaginationQuery) ([]models.Blacklist, int64, error) {
	var blacklists []models.Blacklist
	var total int64

	// Count total records
	r.DB.Model(&models.Blacklist{}).Count(&total)

	// Apply pagination
	offset := utils.GetPaginationOffset(query.PageNumber, query.PageSize)
	if err := r.DB.Offset(offset).Limit(query.PageSize).Find(&blacklists).Error; err != nil {
		return nil, 0, err
	}

	return blacklists, total, nil
}

func (r *BlacklistRepository) FindByID(id uint) (*models.Blacklist, error) {
	var blacklist models.Blacklist
	if err := r.DB.First(&blacklist, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("blacklist entry not found")
		}
		return nil, err
	}
	return &blacklist, nil
}

func (r *BlacklistRepository) FindByPhoneNumber(phoneNumber string) (*models.Blacklist, error) {
	var blacklist models.Blacklist
	if err := r.DB.Where("phone_number = ?", phoneNumber).First(&blacklist).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &blacklist, nil
}

func (r *BlacklistRepository) Update(blacklist *models.Blacklist) (*models.Blacklist, error) {
	if err := r.DB.Save(blacklist).Error; err != nil {
		return nil, err
	}
	return blacklist, nil
}

func (r *BlacklistRepository) Delete(id uint) error {
	if err := r.DB.Delete(&models.Blacklist{}, id).Error; err != nil {
		return err
	}
	return nil
}

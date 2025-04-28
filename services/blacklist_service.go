// services/blacklist_service.go
package services

import (
	"fmt"
	"go-user-service/api/dto"
	"go-user-service/models"
	"go-user-service/repository"
	"go-user-service/utils"
	"time"
)

type BlacklistService struct {
	blacklistRepo *repository.BlacklistRepository
}

func NewBlacklistService(blacklistRepo *repository.BlacklistRepository) *BlacklistService {
	return &BlacklistService{
		blacklistRepo: blacklistRepo,
	}
}

func (s *BlacklistService) CreateBlacklist(createDto *dto.CreateBlacklistDto) (*models.Response, error) {
	// Check if phone number already exists
	existingBlacklist, err := s.blacklistRepo.FindByPhoneNumber(createDto.PhoneNumber)
	if err != nil {
		return nil, err
	}

	if existingBlacklist != nil {
		return nil, fmt.Errorf("phone number %s is already blacklisted", createDto.PhoneNumber)
	}

	// Create new blacklist entry
	blacklist := &models.Blacklist{
		Name:        createDto.Name,
		PhoneNumber: createDto.PhoneNumber,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdBlacklist, err := s.blacklistRepo.Create(blacklist)
	if err != nil {
		return nil, err
	}

	return &models.Response{
		Timestamp: time.Now(),
		Message:   "Blacklist entry created successfully",
		Status:    201,
		Data:      createdBlacklist,
	}, nil
}

func (s *BlacklistService) GetBlacklists(query utils.PaginationQuery) (*models.Response, error) {
	blacklists, total, err := s.blacklistRepo.FindAll(query)
	if err != nil {
		return nil, err
	}

	paginationResponse := utils.CreatePaginationResponse(blacklists, total, query)

	return &models.Response{
		Timestamp: time.Now(),
		Message:   "Blacklist entries retrieved successfully",
		Status:    200,
		Data:      paginationResponse,
	}, nil
}

func (s *BlacklistService) GetBlacklistByID(id uint) (*models.Response, error) {
	blacklist, err := s.blacklistRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &models.Response{
		Timestamp: time.Now(),
		Message:   "Blacklist entry retrieved successfully",
		Status:    200,
		Data:      blacklist,
	}, nil
}

func (s *BlacklistService) UpdateBlacklist(id uint, updateDto *dto.UpdateBlacklistDto) (*models.Response, error) {
	blacklist, err := s.blacklistRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if updateDto.Name != nil {
		blacklist.Name = *updateDto.Name
	}

	if updateDto.PhoneNumber != nil {
		// Check if the new phone number already exists
		if *updateDto.PhoneNumber != blacklist.PhoneNumber {
			existingBlacklist, err := s.blacklistRepo.FindByPhoneNumber(*updateDto.PhoneNumber)
			if err != nil {
				return nil, err
			}

			if existingBlacklist != nil {
				return nil, fmt.Errorf("phone number %s is already blacklisted", *updateDto.PhoneNumber)
			}

			blacklist.PhoneNumber = *updateDto.PhoneNumber
		}
	}

	blacklist.UpdatedAt = time.Now()

	updatedBlacklist, err := s.blacklistRepo.Update(blacklist)
	if err != nil {
		return nil, err
	}

	return &models.Response{
		Timestamp: time.Now(),
		Message:   "Blacklist entry updated successfully",
		Status:    200,
		Data:      updatedBlacklist,
	}, nil
}

func (s *BlacklistService) DeleteBlacklist(id uint) (*models.Response, error) {
	// Check if blacklist entry exists
	_, err := s.blacklistRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if err := s.blacklistRepo.Delete(id); err != nil {
		return nil, err
	}

	return &models.Response{
		Timestamp: time.Now(),
		Message:   fmt.Sprintf("Blacklist entry with id %d has been deleted", id),
		Status:    200,
		Data:      nil,
	}, nil
}

func (s *BlacklistService) IsBlacklisted(phoneNumber string) (bool, error) {
	blacklist, err := s.blacklistRepo.FindByPhoneNumber(phoneNumber)
	if err != nil {
		return false, err
	}

	return blacklist != nil, nil
}

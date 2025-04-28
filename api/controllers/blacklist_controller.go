// api/controllers/blacklist_controller.go
package controllers

import (
	"go-user-service/api/dto"
	"go-user-service/models"
	"go-user-service/services"
	"go-user-service/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type BlacklistController struct {
	blacklistService *services.BlacklistService
}

func NewBlacklistController(blacklistService *services.BlacklistService) *BlacklistController {
	return &BlacklistController{
		blacklistService: blacklistService,
	}
}

// CreateBlacklist godoc
// @Summary Create a new blacklist entry
// @Description Adds a new entry to the blacklist
// @Tags Blacklist
// @Accept json
// @Produce json
// @Param blacklist body dto.CreateBlacklistDto true "Blacklist entry details"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 409 {object} models.Response
// @Failure 500 {object} models.Response
// @Security ApiKeyAuth
// @Router /api/v1/blacklist [post]
func (c *BlacklistController) CreateBlacklist(ctx *gin.Context) {
	var createDto dto.CreateBlacklistDto
	if err := ctx.ShouldBindJSON(&createDto); err != nil {
		ctx.JSON(http.StatusBadRequest, &models.Response{
			Timestamp: time.Now(),
			Message:   "Invalid request body",
			Status:    http.StatusBadRequest,
			Data:      nil,
		})
		return
	}

	response, err := c.blacklistService.CreateBlacklist(&createDto)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "phone number is already blacklisted" {
			status = http.StatusConflict
		}

		ctx.JSON(status, &models.Response{
			Timestamp: time.Now(),
			Message:   err.Error(),
			Status:    status,
			Data:      nil,
		})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

// GetBlacklists godoc
// @Summary Get blacklist entries
// @Description Retrieves all blacklist entries with pagination
// @Tags Blacklist
// @Accept json
// @Produce json
// @Param pageNumber query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Security ApiKeyAuth
// @Router /api/v1/blacklist [get]
func (c *BlacklistController) GetBlacklists(ctx *gin.Context) {
	var query utils.PaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, &models.Response{
			Timestamp: time.Now(),
			Message:   "Invalid query parameters",
			Status:    http.StatusBadRequest,
			Data:      nil,
		})
		return
	}

	// Set default values if not provided
	if query.PageSize == 0 {
		query.PageSize = 10
	}
	if query.PageNumber == 0 {
		query.PageNumber = 1
	}

	response, err := c.blacklistService.GetBlacklists(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &models.Response{
			Timestamp: time.Now(),
			Message:   err.Error(),
			Status:    http.StatusInternalServerError,
			Data:      nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetBlacklistByID godoc
// @Summary Get blacklist entry by ID
// @Description Retrieves a blacklist entry by its ID
// @Tags Blacklist
// @Accept json
// @Produce json
// @Param id path int true "Blacklist ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Security ApiKeyAuth
// @Router /api/v1/blacklist/{id} [get]
func (c *BlacklistController) GetBlacklistByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &models.Response{
			Timestamp: time.Now(),
			Message:   "Invalid ID format",
			Status:    http.StatusBadRequest,
			Data:      nil,
		})
		return
	}

	response, err := c.blacklistService.GetBlacklistByID(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "blacklist entry not found" {
			status = http.StatusNotFound
		}

		ctx.JSON(status, &models.Response{
			Timestamp: time.Now(),
			Message:   err.Error(),
			Status:    status,
			Data:      nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateBlacklist godoc
// @Summary Update blacklist entry
// @Description Updates an existing blacklist entry
// @Tags Blacklist
// @Accept json
// @Produce json
// @Param id path int true "Blacklist ID"
// @Param blacklist body dto.UpdateBlacklistDto true "Updated blacklist details"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 409 {object} models.Response
// @Failure 500 {object} models.Response
// @Security ApiKeyAuth
// @Router /api/v1/blacklist/{id} [put]
func (c *BlacklistController) UpdateBlacklist(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &models.Response{
			Timestamp: time.Now(),
			Message:   "Invalid ID format",
			Status:    http.StatusBadRequest,
			Data:      nil,
		})
		return
	}

	var updateDto dto.UpdateBlacklistDto
	if err := ctx.ShouldBindJSON(&updateDto); err != nil {
		ctx.JSON(http.StatusBadRequest, &models.Response{
			Timestamp: time.Now(),
			Message:   "Invalid request body",
			Status:    http.StatusBadRequest,
			Data:      nil,
		})
		return
	}

	response, err := c.blacklistService.UpdateBlacklist(uint(id), &updateDto)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "blacklist entry not found" {
			status = http.StatusNotFound
		} else if err.Error() == "phone number is already blacklisted" {
			status = http.StatusConflict
		}

		ctx.JSON(status, &models.Response{
			Timestamp: time.Now(),
			Message:   err.Error(),
			Status:    status,
			Data:      nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// DeleteBlacklist godoc
// @Summary Delete blacklist entry
// @Description Deletes a blacklist entry by its ID
// @Tags Blacklist
// @Accept json
// @Produce json
// @Param id path int true "Blacklist ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Security ApiKeyAuth
// @Router /api/v1/blacklist/{id} [delete]
func (c *BlacklistController) DeleteBlacklist(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &models.Response{
			Timestamp: time.Now(),
			Message:   "Invalid ID format",
			Status:    http.StatusBadRequest,
			Data:      nil,
		})
		return
	}

	response, err := c.blacklistService.DeleteBlacklist(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "blacklist entry not found" {
			status = http.StatusNotFound
		}

		ctx.JSON(status, &models.Response{
			Timestamp: time.Now(),
			Message:   err.Error(),
			Status:    status,
			Data:      nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

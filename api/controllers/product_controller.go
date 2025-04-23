package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"go-user-service/api/dto"
	"go-user-service/models"
	"go-user-service/services"
	"go-user-service/utils"
)

type ProductController struct {
	ProductService *services.ProductService
	DB             *gorm.DB
}

func NewProductController(productService *services.ProductService, db *gorm.DB) *ProductController {
	return &ProductController{
		ProductService: productService,
		DB:             db,
	}
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func derefInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the provided details
// @Tags products
// @Accept json
// @Produce json
// @Param product body dto.CreateProductRequest true "Product object"
// @Success 200 {object} models.Product
// @Router /api/v1/product [post]
func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var req dto.CreateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := &models.Product{
		ProductName:     &req.ProductName,
		Description:     &req.Description,
		ProductPicture:  &req.ProductPicture,
		IsAvailable:     &req.IsAvailable,
		IsCallNeeded:    &req.IsCallNeeded,
		ProductMargin:   &req.ProductMargin,
		IsBonus:         &req.IsBonus,
		ProductCost:     req.ProductCost,
		DrawPeriod:      req.DrawPeriod,
		NumberOfWinners: req.NumberOfWinners,
		PlayAmount:      req.PlayAmount,
		Priority:        req.Priority,
		EnglishName:     req.EnglishName,
	}

	createdProduct, err := c.ProductService.Create(product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, createdProduct)
}

// AssignBonusProducts godoc
// @Summary Assign bonus products to a product
// @Description Assign bonus products to a specific product
// @Tags products
// @Accept json
// @Produce json
// @Param productId path string true "Product ID"
// @Param bonusProducts body dto.BonusProductDTO true "Bonus Product IDs"
// @Success 200 {object} models.Product
// @Router /api/v1/product/assign-bonus/{productId} [post]
func (c *ProductController) AssignBonusProducts(ctx *gin.Context) {
	productID := ctx.Param("productId")

	var bonusProductDTO dto.BonusProductDTO
	if err := ctx.ShouldBindJSON(&bonusProductDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := c.ProductService.AssignBonusProducts(productID, bonusProductDTO.BonusProductIds)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// UnassignBonusProducts godoc
// @Summary Unassign bonus products from a product
// @Description Unassign bonus products from a specific product
// @Tags products
// @Accept json
// @Produce json
// @Param productId path string true "Product ID"
// @Param bonusProducts body dto.BonusProductDTO true "Bonus Product IDs"
// @Success 200 {object} models.Product
// @Router /api/v1/product/unassign-bonus/{productId} [patch]
func (c *ProductController) UnassignBonusProducts(ctx *gin.Context) {
	productID := ctx.Param("productId")

	var bonusProductDTO dto.BonusProductDTO
	if err := ctx.ShouldBindJSON(&bonusProductDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := c.ProductService.UnassignBonusProducts(productID, bonusProductDTO.BonusProductIds)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// GetProducts godoc
// @Summary Get all products
// @Description Get a paginated list of products with optional filtering
// @Tags products
// @Accept json
// @Produce json
// @Param pageNumber query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param from query string false "From date (YYYY-MM-DD)"
// @Param to query string false "To date (YYYY-MM-DD)"
// @Param search query string false "Search term"
// @Param isAvailable query bool false "Filter by availability"
// @Success 200 {object} utils.PaginationResponse
// @Router /api/v1/product [get]
func (c *ProductController) GetProducts(ctx *gin.Context) {
	pageNumber, _ := strconv.Atoi(ctx.DefaultQuery("pageNumber", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	from := ctx.Query("from")
	to := ctx.Query("to")
	search := ctx.Query("search")
	isAvailableStr := ctx.Query("isAvailable")

	query := c.DB.Model(&models.Product{}).Order("priority DESC")

	if from != "" && to != "" {
		fromDate, _ := time.Parse("2006-01-02", from)
		toDate, _ := time.Parse("2006-01-02", to)
		query = query.Where("\"createdAt\" BETWEEN ? AND ?", fromDate, toDate)
	} else if from != "" {
		fromDate, _ := time.Parse("2006-01-02", from)
		query = query.Where("\"createdAt\" >= ?", fromDate)
	} else if to != "" {
		toDate, _ := time.Parse("2006-01-02", to)
		query = query.Where("\"createdAt\" <= ?", toDate)
	}

	if search != "" {
		query = query.Where("\"productName\" ILIKE ?", "%"+search+"%")
	}

	if isAvailableStr != "" {
		isAvailable, _ := strconv.ParseBool(isAvailableStr)
		query = query.Where("\"isAvailable\" = ?", isAvailable)
	}

	var total int64
	query.Count(&total)

	offset := (pageNumber - 1) * pageSize
	var products []models.Product
	query.Limit(pageSize).Offset(offset).Find(&products)

	modifiedProducts := make([]dto.GetProductResponseDto, 0, len(products))
	for _, product := range products {
		var latestDraw models.Draw
		c.DB.Where("\"productId\" = ?", product.ProductID).Order("\"createdAt\" DESC").First(&latestDraw)

		p := dto.GetProductResponseDto{
			ProductID:          product.ProductID.String(),
			ProductIncrementer: product.ProductIncrementer,
			ProductName:        derefString(product.ProductName),
			ProductPicture:     derefString(product.ProductPicture),
			Description:        derefString(product.Description),
			IsAvailable:        derefBool(product.IsAvailable),
			IsCallNeeded:       derefBool(product.IsCallNeeded),
			AvoidConflict:      derefBool(product.AvoidConflict),
			IsBonus:            derefBool(product.IsBonus),
			ProductCost:        derefInt(product.ProductCost),
			DrawPeriod:         derefInt(product.DrawPeriod),
			NumberOfWinners:    derefInt(product.NumberOfWinners),
			PlayAmount:         derefInt(product.PlayAmount),
			Priority:           derefInt(product.Priority),
			Comment:            product.Comment,
			EnglishName:        product.EnglishName,
			ProductMargin:      derefInt(product.ProductMargin),
			ExpectedAmount:     derefInt(product.ExpectedAmount),
			CreatedAt:          product.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z"),
			UpdatedAt:          product.UpdatedAt.UTC().Format("2006-01-02T15:04:05.000Z"),
		}

		if latestDraw.DrawID != uuid.Nil {
			p.Draw = &dto.DrawResponse{
				DrawID:    latestDraw.DrawID,
				StartDate: latestDraw.StartDate,
				EndDate:   latestDraw.EndDate,
				IsPlayed:  latestDraw.IsPlayed,
				CreatedAt: latestDraw.CreatedAt,
				UpdatedAt: latestDraw.UpdatedAt,
			}
		}

		modifiedProducts = append(modifiedProducts, p)
	}

	lastPage := int(total) / pageSize
	if int(total)%pageSize > 0 {
		lastPage++
	}

	var nextPage *int
	if pageNumber < lastPage {
		n := pageNumber + 1
		nextPage = &n
	}

	var prevPage *int
	if pageNumber > 1 {
		p := pageNumber - 1
		prevPage = &p
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"list":         modifiedProducts,
		"total":        total,
		"previousPage": prevPage,
		"nextPage":     nextPage,
		"lastPage":     lastPage,
		"currentPage":  pageNumber,
	})
}

// GetProductsForUSSD godoc
// @Summary Get products for USSD
// @Description Get a paginated list of available products for USSD
// @Tags products
// @Accept json
// @Produce json
// @Param pageNumber query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param from query string false "From date (YYYY-MM-DD)"
// @Param to query string false "To date (YYYY-MM-DD)"
// @Param search query string false "Search term"
// @Success 200 {object} utils.PaginationResponse
// @Router /api/v1/product/ussd [get]
func (c *ProductController) GetProductsForUSSD(ctx *gin.Context) {
	// Parse query parameters
	pageNumber, _ := strconv.Atoi(ctx.DefaultQuery("pageNumber", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))
	from := ctx.Query("from")
	to := ctx.Query("to")
	search := ctx.Query("search")

	// Build query
	query := c.DB.Model(&models.Product{}).Where("\"isAvailable\" = ?", true).Order("priority DESC")

	// Apply filters
	if from != "" && to != "" {
		fromDate, _ := time.Parse("2006-01-02", from)
		toDate, _ := time.Parse("2006-01-02", to)
		query = query.Where("createdAt BETWEEN ? AND ?", fromDate, toDate)
	} else if from != "" {
		fromDate, _ := time.Parse("2006-01-02", from)
		query = query.Where("createdAt >= ?", fromDate)
	} else if to != "" {
		toDate, _ := time.Parse("2006-01-02", to)
		query = query.Where("createdAt <= ?", toDate)
	}

	if search != "" {
		query = query.Where("productName ILIKE ?", "%"+search+"%")
	}

	// Count total records
	var total int64
	query.Count(&total)

	// Apply pagination
	offset := (pageNumber - 1) * pageSize
	var products []models.Product
	query.Limit(pageSize).Offset(offset).Find(&products)

	// Create pagination response
	paginationQuery := utils.PaginationQuery{
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}

	response := utils.CreatePaginationResponse(products, total, paginationQuery)
	ctx.JSON(http.StatusOK, response)
}

// GetProduct godoc
// @Summary Get a product by ID
// @Description Get a product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} models.Product
// @Router /api/v1/product/{id} [get]
func (c *ProductController) GetProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	product, err := c.ProductService.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// GetProductsByBonus godoc
// @Summary Get products by bonus product ID
// @Description Get all products that have a specific bonus product
// @Tags products
// @Accept json
// @Produce json
// @Param bonusProductId path string true "Bonus Product ID"
// @Success 200 {array} models.Product
// @Router /api/v1/product/bonus/{bonusProductId}/assigned-products [get]
func (c *ProductController) GetProductsByBonus(ctx *gin.Context) {
	bonusProductID := ctx.Param("bonusProductId")

	products, err := c.ProductService.GetProductsByBonus(bonusProductID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

// GetBonusProducts godoc
// @Summary Get bonus products for a product
// @Description Get all bonus products assigned to a specific product
// @Tags products
// @Accept json
// @Produce json
// @Param productId path string true "Product ID"
// @Router /api/v1/product/bonus-products/{productId} [get]
func (c *ProductController) GetBonusProducts(ctx *gin.Context) {
	productID := ctx.Param("productId")

	product, err := c.ProductService.GetBonusProductsByProductID(productID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update a product with the provided details
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body models.Product true "Product object"
// @Success 200 {object} dto.ProductResponse
// @Router /api/v1/product/{id} [patch]
func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	var product models.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ProductService.Update(id, &product); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.ProductResponse{
		Status:  "success",
		Message: "Product updated successfully",
	})
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Toggle the availability of a product (soft delete)
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.ProductResponse
// @Router /api/v1/product/{id} [delete]
func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.ProductService.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.ProductResponse{
		Status:  "success",
		Message: "Product deleted successfully",
	})
}

// UploadFile godoc
// @Summary Upload a file
// @Description Upload a file and get its link
// @Tags products
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 200 {object} dto.ProductUploadResponse
// @Router /api/v1/product/upload [post]
func (c *ProductController) UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create uploads directory if it doesn't exist
	uploadDir := filepath.Join(".", "uploads")
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	// Generate a unique filename
	filename := fmt.Sprintf("%s-%s", uuid.New().String(), file.Filename)
	filePath := filepath.Join(uploadDir, filename)

	// Save the file
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fileLink := c.ProductService.GetFileLink(filename)
	ctx.JSON(http.StatusOK, dto.ProductUploadResponse{
		FileLink: fileLink,
	})
}

// DeleteFile godoc
// @Summary Delete a file
// @Description Delete a file by its filename
// @Tags products
// @Accept json
// @Produce json
// @Param filename path string true "Filename"
// @Success 200 {object} dto.ProductResponse
// @Router /api/v1/product/file/delete/{filename} [delete]
func (c *ProductController) DeleteFile(ctx *gin.Context) {
	filename := ctx.Param("filename")

	success := c.ProductService.DeleteFile(filename)
	if success {
		ctx.JSON(http.StatusOK, dto.ProductResponse{
			Status:  "success",
			Message: fmt.Sprintf("File %s deleted successfully", filename),
		})
	} else {
		ctx.JSON(http.StatusInternalServerError, dto.ProductResponse{
			Status:  "error",
			Message: fmt.Sprintf("Failed to delete file %s", filename),
		})
	}
}

package services

import (
	"errors"
	"fmt"
	"go-user-service/models"
	"go-user-service/repository"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductService struct {
	DB          *gorm.DB
	ProductRepo *repository.ProductRepository
	UploadPath  string
}

func NewProductService(db *gorm.DB, productRepo *repository.ProductRepository) *ProductService {
	uploadPath := filepath.Join(".", "uploads")
	return &ProductService{
		DB:          db,
		ProductRepo: productRepo,
		UploadPath:  uploadPath,
	}
}

func (s *ProductService) Create(product *models.Product) (*models.Product, error) {
	// Get next incrementer
	lastProduct, err := s.ProductRepo.FindLastProduct()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	nextID := 1
	if lastProduct != nil && lastProduct.ProductIncrementer > 0 {
		nextID = lastProduct.ProductIncrementer + 1
	}
	product.ProductIncrementer = nextID

	// Calculate expected amount
	if product.ProductCost != nil && product.ProductMargin != nil {
		expected := *product.ProductCost + (*product.ProductCost * *product.ProductMargin / 100)
		product.ExpectedAmount = &expected
	}

	if err := s.ProductRepo.Create(product); err != nil {
		return nil, err
	}

	return product, nil
}
func (s *ProductService) FindById(id string) (*models.Product, error) {
	productID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	product, err := s.ProductRepo.FindByID(productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("product not found with id: %s", id)
		}
		return nil, err
	}

	return product, nil
}

func (s *ProductService) Update(id string, updateProduct *models.Product) error {
	productID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid product ID: %w", err)
	}

	// Find the existing product
	existingProduct, err := s.ProductRepo.FindByID(productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("product not found with id: %s", id)
		}
		return err
	}

	// Calculate expected amount if product margin or cost is updated
	if updateProduct.ProductMargin != nil || updateProduct.ProductCost != nil {
		productCost := updateProduct.ProductCost
		if productCost == nil {
			productCost = existingProduct.ProductCost
		}

		productMargin := updateProduct.ProductMargin
		if productMargin == nil {
			productMargin = existingProduct.ProductMargin
		}

		if productCost != nil && productMargin != nil {
			expectedAmount := *productCost * (1 + *productMargin/100)
			updateProduct.ExpectedAmount = &expectedAmount
		}
	}

	// Update the product
	updateProduct.ProductID = productID
	updateProduct.UpdatedAt = time.Now()

	return s.ProductRepo.Update(updateProduct)
}

func (s *ProductService) Delete(id string) error {
	productID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid product ID: %w", err)
	}

	// Find the existing product
	product, err := s.ProductRepo.FindByID(productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("product not found with id: %s", id)
		}
		return err
	}

	// Toggle isAvailable
	isAvailable := false
	if product.IsAvailable != nil {
		isAvailable = !*product.IsAvailable
	}
	product.IsAvailable = &isAvailable

	return s.ProductRepo.Update(product)
}

func (s *ProductService) GetFileLink(filename string) string {
	host := os.Getenv("HOST")
	if host == "" {
		host = "http://localhost:8080"
	}
	return fmt.Sprintf("%s/uploads/%s", host, filename)
}

func (s *ProductService) DeleteFile(filename string) bool {
	filePath := filepath.Join(s.UploadPath, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("File %s does not exist at path %s\n", filename, filePath)
		return false
	}

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		fmt.Printf("Error deleting file %s: %v\n", filename, err)
		return false
	}

	return true
}

func (s *ProductService) AssignBonusProducts(productID string, bonusProductIDs []string) (*models.Product, error) {
	// Parse the product ID
	parsedProductID, err := uuid.Parse(productID)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	// Find the product with its bonus products
	product, err := s.ProductRepo.FindByIDWithBonusProducts(parsedProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("product not found with id: %s", productID)
		}
		return nil, err
	}

	// Parse bonus product IDs
	parsedBonusProductIDs := make([]uuid.UUID, 0, len(bonusProductIDs))
	for _, id := range bonusProductIDs {
		parsedID, err := uuid.Parse(id)
		if err != nil {
			return nil, fmt.Errorf("invalid bonus product ID: %w", err)
		}
		parsedBonusProductIDs = append(parsedBonusProductIDs, parsedID)
	}

	// Find all bonus products
	bonusProducts, err := s.ProductRepo.FindBonusProductsByIDs(parsedBonusProductIDs)
	if err != nil {
		return nil, err
	}

	// Validate that all requested bonus products exist and are marked as bonuses
	if len(bonusProducts) != len(bonusProductIDs) {
		return nil, errors.New("one or more bonus products are invalid or not marked as bonuses")
	}

	// Assign bonus products
	if err := s.ProductRepo.AssignBonusProducts(product, bonusProducts); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) GetProductsByBonus(bonusProductID string) ([]models.Product, error) {
	// Parse the bonus product ID
	parsedBonusProductID, err := uuid.Parse(bonusProductID)
	if err != nil {
		return nil, fmt.Errorf("invalid bonus product ID: %w", err)
	}

	// Find the bonus product
	bonusProduct, err := s.ProductRepo.FindByID(parsedBonusProductID)
	fmt.Print(bonusProduct)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("bonus product not found with id: %s", bonusProductID)
		}
		return nil, err
	}

	// Find products that have this bonus product
	products, err := s.ProductRepo.FindProductsByBonusProductID(parsedBonusProductID)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *ProductService) UnassignBonusProducts(productID string, bonusProductIDs []string) (*models.Product, error) {
	// Parse the product ID
	parsedProductID, err := uuid.Parse(productID)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	// Find the product with its bonus products
	product, err := s.ProductRepo.FindByIDWithBonusProducts(parsedProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("product not found with id: %s", productID)
		}
		return nil, err
	}

	// Parse bonus product IDs to unassign
	parsedBonusProductIDs := make([]uuid.UUID, 0, len(bonusProductIDs))
	for _, id := range bonusProductIDs {
		parsedID, err := uuid.Parse(id)
		if err != nil {
			return nil, fmt.Errorf("invalid bonus product ID: %w", err)
		}
		parsedBonusProductIDs = append(parsedBonusProductIDs, parsedID)
	}

	// Unassign bonus products
	if err := s.ProductRepo.UnassignBonusProducts(product, parsedBonusProductIDs); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) GetBonusProductsByProductID(productID string) (*models.Product, error) {
	// Parse the product ID
	parsedProductID, err := uuid.Parse(productID)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	// Find the product with its bonus products
	product, err := s.ProductRepo.FindByIDWithBonusProducts(parsedProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("product not found with id: %s", productID)
		}
		return nil, err
	}

	return product, nil
}

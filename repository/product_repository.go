package repository

import (
	"go-user-service/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		DB: db,
	}
}

func (r *ProductRepository) Create(product *models.Product) error {
	return r.DB.Create(product).Error
}

func (r *ProductRepository) FindByID(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	if err := r.DB.Where("\"productId\" = ?", id).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) FindByIDWithBonusProducts(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	if err := r.DB.Preload("BonusProducts").
		Where(`"productId" = ?`, id).
		First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) FindLastProduct() (*models.Product, error) {
	var product models.Product
	if err := r.DB.Order("\"productIcrementer\" DESC").First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) Update(product *models.Product) error {
	return r.DB.Model(&models.Product{}).Where("productId = ?", product.ProductID).Updates(product).Error
}

func (r *ProductRepository) FindBonusProductsByIDs(ids []uuid.UUID) ([]*models.Product, error) {
	var products []*models.Product
	if err := r.DB.Where("productId IN ? AND isBonus = ?", ids, true).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) AssignBonusProducts(product *models.Product, bonusProducts []*models.Product) error {
	// Start a transaction
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// Clear existing associations
	if err := tx.Model(product).Association("BonusProducts").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	// Add new associations
	if err := tx.Model(product).Association("BonusProducts").Append(bonusProducts); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *ProductRepository) FindProductsByBonusProductID(bonusProductID uuid.UUID) ([]models.Product, error) {
	var products []models.Product

	// This query finds all products that have the specified bonus product
	err := r.DB.Joins("JOIN product_bonus ON products.productId = product_bonus.productId").
		Where("product_bonus.bonus_product_id = ?", bonusProductID).
		Preload("BonusProducts").
		Find(&products).Error

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) UnassignBonusProducts(product *models.Product, bonusProductIDs []uuid.UUID) error {
	// Start a transaction
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// Remove specific bonus products
	for _, bonusProductID := range bonusProductIDs {
		if err := tx.Exec("DELETE FROM product_bonus WHERE productId = ? AND bonus_product_id = ?",
			product.ProductID, bonusProductID).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

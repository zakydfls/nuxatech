package repository

import (
	"context"
	"nuxatech-nextmedis/config"
	"nuxatech-nextmedis/model"
	"strings"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

// GetAllProducts implements ProductRepository.
func (p *productRepository) GetAllProducts(ctx context.Context, page int, limit int, search string) ([]*model.Product, int64, error) {
	var products []*model.Product
	var total int64

	query := p.db.Model(&model.Product{})

	if search != "" {
		searchQuery := "%" + strings.ToLower(search) + "%"
		query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR LOWER(sku) LIKE ?",
			searchQuery, searchQuery, searchQuery)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (r *productRepository) GetProductForUpdate(ctx context.Context, tx *gorm.DB, id string) (*model.Product, error) {
	var product model.Product
	err := tx.WithContext(ctx).
		Set("gorm:for_update", true).
		First(&product, "id = ?", id).
		Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) UpdateStock(ctx context.Context, tx *gorm.DB, productID string, newStock int) error {
	return tx.WithContext(ctx).
		Model(&model.Product{}).
		Where("id = ?", productID).
		Update("stock", newStock).
		Error
}

// CheckSlug implements ProductRepository.
func (p *productRepository) CheckSlug(ctx context.Context, slug string, sku string) error {
	var product *model.Product
	err := p.db.Where("LOWER(slug) = ? OR LOWER(sku) = ?", slug).Find(&product)
	if err != nil {
		return err.Error
	}
	return nil
}

// CreateProduct implements ProductRepository.
func (p *productRepository) CreateProduct(ctx context.Context, product *model.Product) error {
	return p.db.Create(product).Error
}

// DeleteProduct implements ProductRepository.
func (p *productRepository) DeleteProduct(ctx context.Context, id string) error {
	return p.db.Delete(&model.Product{}, id).Error
}

// GetProduct implements ProductRepository.
func (p *productRepository) GetProduct(ctx context.Context, id string) (*model.Product, error) {
	var product model.Product
	if err := p.db.Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// UpdateProduct implements ProductRepository.
func (p *productRepository) UpdateProduct(ctx context.Context, product *model.Product) error {
	return p.db.Save(product).Error
}

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *model.Product) error
	GetProduct(ctx context.Context, id string) (*model.Product, error)
	UpdateProduct(ctx context.Context, product *model.Product) error
	DeleteProduct(ctx context.Context, id string) error
	GetAllProducts(ctx context.Context, page int, limit int, search string) ([]*model.Product, int64, error)
	CheckSlug(ctx context.Context, slug string, sku string) error
	GetProductForUpdate(ctx context.Context, tx *gorm.DB, id string) (*model.Product, error)
	UpdateStock(ctx context.Context, tx *gorm.DB, productID string, newStock int) error
}

func NewProductRepository() ProductRepository {
	return &productRepository{db: config.GetDB()}
}

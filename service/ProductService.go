package service

import (
	"context"
	"errors"
	"nuxatech-nextmedis/dto/request"
	"nuxatech-nextmedis/dto/response"
	"nuxatech-nextmedis/model"
	"nuxatech-nextmedis/repository"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gosimple/slug"
)

type ProductService interface {
	CreateProduct(ctx context.Context, product *request.CreateProductRequest) (*model.Product, error)
	GetProduct(ctx context.Context, id string) (*model.Product, error)
	UpdateProduct(ctx context.Context, id string, product *request.UpdateProductRequest) (*model.Product, error)
	DeleteProduct(ctx context.Context, id string) error
	GetAllProducts(ctx context.Context, params ProductQueryParams) (*response.PagingResponse, error)
}

type productService struct {
	productRepo repository.ProductRepository
	validate    *validator.Validate
}

type ProductQueryParams struct {
	Page   int
	Limit  int
	Search string
}

// GetAllProducts implements ProductService.
func (p *productService) GetAllProducts(ctx context.Context, params ProductQueryParams) (*response.PagingResponse, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	products, total, err := p.productRepo.GetAllProducts(ctx, params.Page, params.Limit, params.Search)
	if err != nil {
		return nil, err
	}

	return &response.PagingResponse{
		Metadata: response.Metadata{
			TotalCount: int(total),
			Page:       params.Page,
			PerPage:    params.Limit,
		},
		Result: products,
	}, nil
}

// CreateProduct implements ProductService.
func (p *productService) CreateProduct(ctx context.Context, product *request.CreateProductRequest) (*model.Product, error) {
	if err := p.validate.Struct(product); err != nil {
		return nil, err
	}

	slug := slug.Make(product.Name)
	lowerSku := strings.ToLower(product.SKU)
	errSlug := p.productRepo.CheckSlug(ctx, slug, lowerSku)
	if errSlug == nil {
		return nil, errors.New("slug already exist")
	}

	var images model.LocalProductImages
	images = append(images, product.Image...)

	newProduct := &model.Product{
		Name:        product.Name,
		Description: product.Description,
		Image:       images,
		Stock:       product.Stock,
		Price:       product.Price,
		BasePrice:   product.BasePrice,
		SKU:         product.SKU,
		Slug:        slug,
		Weight:      product.Weight,
		Sold:        false,
		CreatedAt:   time.Now().UnixMilli(),
	}

	err := p.productRepo.CreateProduct(ctx, newProduct)
	if err != nil {
		return nil, err
	}

	return newProduct, nil
}

// DeleteProduct implements ProductService.
func (p *productService) DeleteProduct(ctx context.Context, id string) error {
	panic("unimplemented")
}

// GetProduct implements ProductService.
func (p *productService) GetProduct(ctx context.Context, id string) (*model.Product, error) {
	if id == "" {
		return nil, errors.New("missing id params")
	}
	product, err := p.productRepo.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}
	return product, err
}

// UpdateProduct implements ProductService.
func (p *productService) UpdateProduct(ctx context.Context, id string, product *request.UpdateProductRequest) (*model.Product, error) {
	panic("unimplemented")
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{
		productRepo: productRepo,
		validate:    validator.New(),
	}
}

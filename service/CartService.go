package service

import (
	"context"
	"errors"
	"nuxatech-nextmedis/dto/request"
	"nuxatech-nextmedis/dto/response"
	"nuxatech-nextmedis/model"
	"nuxatech-nextmedis/repository"
	"time"

	"github.com/go-playground/validator/v10"
)

type CartService interface {
	AddToCart(ctx context.Context, userID string, req *request.AddToCartRequest) (*response.CartResponse, error)
	GetCart(ctx context.Context, userID string) (*response.CartResponse, error)
	UpdateCartItem(ctx context.Context, userID string, cartItemID string, req *request.UpdateCartItemRequest) (*response.CartResponse, error)
	RemoveFromCart(ctx context.Context, userID string, cartItemID string) error
}

type cartService struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
	validate    *validator.Validate
}

func (s *cartService) AddToCart(ctx context.Context, userID string, req *request.AddToCartRequest) (*response.CartResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	// Get product
	product, err := s.productRepo.GetProduct(ctx, req.ProductID)
	if err != nil {
		return nil, err
	}

	// Check stock
	if product.Stock < req.Quantity {
		return nil, errors.New("insufficient stock")
	}

	// Get or create cart
	cart, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check if product already in cart
	existingItem, err := s.cartRepo.GetCartItem(ctx, cart.ID, req.ProductID)
	if err == nil {
		existingItem.Quantity += req.Quantity
		existingItem.UpdatedAt = time.Now().UnixMilli()
		if err := s.cartRepo.UpdateItem(ctx, existingItem); err != nil {
			return nil, err
		}
	} else {
		cartItem := &model.CartItem{
			CartID:    cart.ID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
			CreatedAt: time.Now().UnixMilli(),
			UpdatedAt: time.Now().UnixMilli(),
		}
		if err := s.cartRepo.AddItem(ctx, cartItem); err != nil {
			return nil, err
		}
	}

	// Get updated cart
	return s.GetCart(ctx, userID)
}

func (s *cartService) GetCart(ctx context.Context, userID string) (*response.CartResponse, error) {
	cart, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Transform to response
	cartItems := make([]response.CartItemResponse, len(cart.Items))
	total := 0
	for i, item := range cart.Items {
		cartItems[i] = response.CartItemResponse{
			ID:       item.ID,
			Product:  toProductResponse(item.Product),
			Quantity: item.Quantity,
		}
		total += item.Product.Price * item.Quantity
	}

	return &response.CartResponse{
		ID:        cart.ID,
		Items:     cartItems,
		Total:     total,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}, nil
}

func (s *cartService) UpdateCartItem(ctx context.Context, userID string, cartItemID string, req *request.UpdateCartItemRequest) (*response.CartResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	// Get cart
	cart, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	var cartItem *model.CartItem
	for _, item := range cart.Items {
		if item.ID == cartItemID {
			cartItem = &item
			break
		}
	}
	if cartItem == nil {
		return nil, errors.New("cart item not found")
	}

	// Check stock
	product, err := s.productRepo.GetProduct(ctx, cartItem.ProductID)
	if err != nil {
		return nil, err
	}
	if product.Stock < req.Quantity {
		return nil, errors.New("insufficient stock")
	}

	cartItem.Quantity = req.Quantity
	cartItem.UpdatedAt = time.Now().UnixMilli()
	if err := s.cartRepo.UpdateItem(ctx, cartItem); err != nil {
		return nil, err
	}

	return s.GetCart(ctx, userID)
}

func (s *cartService) RemoveFromCart(ctx context.Context, userID string, cartItemID string) error {
	// Get cart
	cart, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	// Verify item belongs to user's cart
	var found bool
	for _, item := range cart.Items {
		if item.ID == cartItemID {
			found = true
			break
		}
	}
	if !found {
		return errors.New("cart item not found")
	}

	return s.cartRepo.RemoveItem(ctx, cartItemID)
}

func toProductResponse(product model.Product) response.ProductResponse {
	return response.ProductResponse{
		ID:             product.ID,
		Name:           product.Name,
		Description:    product.Description,
		Image:          product.Image,
		Stock:          product.Stock,
		Price:          product.Price,
		Weight:         product.Weight,
		BasePrice:      product.BasePrice,
		SKU:            product.SKU,
		Slug:           product.Slug,
		UniqueCodeType: product.UniqueCodeType,
		Sold:           product.Sold,
		CreatedAt:      product.CreatedAt,
	}
}

func NewCartService(cartRepo repository.CartRepository, productRepo repository.ProductRepository) CartService {
	return &cartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
		validate:    validator.New(),
	}
}

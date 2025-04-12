package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"nuxatech-nextmedis/dto/request"
	"nuxatech-nextmedis/dto/response"
	"nuxatech-nextmedis/model"
	"nuxatech-nextmedis/repository"
	"sync"
	"time"

	"slices"

	"github.com/go-playground/validator/v10"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userID string, req *request.CreateOrderRequest) (*response.OrderResponse, error)
	GetOrder(ctx context.Context, userID, orderID string) (*response.OrderResponse, error)
	UpdateOrderStatus(ctx context.Context, userID, orderID string, req *request.UpdateOrderStatusRequest) (*response.OrderResponse, error)
	GetUserOrders(ctx context.Context, userID string, params ProductQueryParams) (*response.OrderPagingResponse, error)
}

type orderService struct {
	orderRepo   repository.OrderRepository
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
	accountRepo repository.AccountRepository
	validate    *validator.Validate
	mutex       sync.Mutex
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, userID, orderID string, req *request.UpdateOrderStatusRequest) (*response.OrderResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	order, err := s.orderRepo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	now := time.Now().UnixMilli()
	order.Status = model.OrderStatus(req.Status)
	order.UpdatedAt = now

	if req.Status == "paid" {
		order.PaidAt = &now
	}

	if err := s.orderRepo.UpdateOrder(ctx, order); err != nil {
		return nil, err
	}

	return s.toOrderResponse(order), nil
}

func (s *orderService) GetOrder(ctx context.Context, userID, orderID string) (*response.OrderResponse, error) {
	order, err := s.orderRepo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	return s.toOrderResponse(order), nil
}

func (s *orderService) GetUserOrders(ctx context.Context, userID string, params ProductQueryParams) (*response.OrderPagingResponse, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	orders, total, err := s.orderRepo.GetUserOrders(ctx, userID, params.Page, params.Limit)
	if err != nil {
		return nil, err
	}

	return &response.OrderPagingResponse{
		Metadata: response.Metadata{
			TotalCount: int(total),
			Page:       params.Page,
			PerPage:    params.Limit,
		},
		Result: orders,
	}, nil
}

func (s *orderService) CreateOrder(ctx context.Context, userID string, req *request.CreateOrderRequest) (*response.OrderResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	cart, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	selectedItems := make([]model.CartItem, 0)
	for _, item := range cart.Items {
		if slices.Contains(req.SelectedItems, item.ID) {
			selectedItems = append(selectedItems, item)
		}
	}

	if len(selectedItems) == 0 {
		return nil, errors.New("no valid items selected")
	}

	if len(selectedItems) != len(req.SelectedItems) {
		return nil, errors.New("some selected items were not found in cart")
	}

	tx := s.orderRepo.BeginTx(ctx)
	if tx == nil {
		return nil, errors.New("failed to start transaction")
	}
	defer tx.Rollback()

	var totalAmount int64
	now := time.Now().UnixMilli()
	orderItems := make([]model.OrderItem, len(selectedItems))

	for i, item := range selectedItems {
		product, err := s.productRepo.GetProductForUpdate(ctx, tx, item.ProductID)
		if err != nil {
			return nil, err
		}

		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product: %s (available: %d, requested: %d)",
				product.Name, product.Stock, item.Quantity)
		}

		price := int64(product.Price)
		quantity := item.Quantity
		itemTotal := price * int64(quantity)

		orderItems[i] = model.OrderItem{
			ProductID: item.ProductID,
			Product:   *product,
			Quantity:  quantity,
			Price:     price,
			CreatedAt: now,
		}
		totalAmount += itemTotal

		if err := s.productRepo.UpdateStock(ctx, tx, product.ID, product.Stock-quantity); err != nil {
			return nil, fmt.Errorf("failed to update stock: %w", err)
		}
	}

	order := &model.Order{
		UserID:      userID,
		CartID:      cart.ID,
		Status:      model.OrderStatusPending,
		TotalAmount: totalAmount,
		Items:       orderItems,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.orderRepo.CreateOrder(ctx, tx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	for _, itemID := range req.SelectedItems {
		if err := s.cartRepo.RemoveItem(ctx, itemID); err != nil {
			log.Printf("Failed to remove item %s from cart: %v", itemID, err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return s.toOrderResponse(order), nil
}

func (s *orderService) toOrderResponse(order *model.Order) *response.OrderResponse {
	items := make([]response.OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		items[i] = response.OrderItemResponse{
			ID:       item.ID,
			Product:  toProductResponse(item.Product),
			Quantity: item.Quantity,
			Price:    item.Price,
		}
	}

	return &response.OrderResponse{
		ID:          order.ID,
		Status:      string(order.Status),
		TotalAmount: order.TotalAmount,
		Items:       items,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
		PaidAt:      order.PaidAt,
	}
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	cartRepo repository.CartRepository,
	productRepo repository.ProductRepository,
	accountRepo repository.AccountRepository,
) OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
		accountRepo: accountRepo,
		validate:    validator.New(),
	}
}

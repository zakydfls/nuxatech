package repository

import (
	"context"
	"nuxatech-nextmedis/config"
	"nuxatech-nextmedis/model"

	"gorm.io/gorm"
)

type OrderRepository interface {
	BeginTx(ctx context.Context) *gorm.DB
	CreateOrder(ctx context.Context, tx *gorm.DB, order *model.Order) error
	GetOrder(ctx context.Context, id string) (*model.Order, error)
	UpdateOrder(ctx context.Context, order *model.Order) error
	GetUserOrders(ctx context.Context, userID string, page, limit int) ([]*model.Order, int64, error)
}

type orderRepository struct {
	db *gorm.DB
}

func (r *orderRepository) BeginTx(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Begin()
}

func (r *orderRepository) CreateOrder(ctx context.Context, tx *gorm.DB, order *model.Order) error {
	db := tx
	if tx == nil {
		db = r.db.WithContext(ctx)
	}

	return db.Transaction(func(tx *gorm.DB) error {
		orderOnly := &model.Order{
			ID:          order.ID,
			UserID:      order.UserID,
			CartID:      order.CartID,
			Status:      order.Status,
			TotalAmount: order.TotalAmount,
			CreatedAt:   order.CreatedAt,
			UpdatedAt:   order.UpdatedAt,
		}

		if err := tx.Create(orderOnly).Error; err != nil {
			return err
		}

		items := make([]model.OrderItem, len(order.Items))
		for i, item := range order.Items {
			items[i] = model.OrderItem{
				OrderID:   orderOnly.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Price,
				CreatedAt: item.CreatedAt,
			}
		}

		if len(items) > 0 {
			if err := tx.CreateInBatches(items, len(items)).Error; err != nil {
				return err
			}
		}

		if err := tx.Where("id = ?", orderOnly.ID).
			Preload("Items").
			Preload("Items.Product").
			Find(order).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *orderRepository) GetOrder(ctx context.Context, id string) (*model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Preload("Items.Product").
		Find(&order).Error
	if err != nil {
		return nil, err
	}

	// Check if order was found
	if order.ID == "" {
		return nil, gorm.ErrRecordNotFound
	}

	return &order, nil
}

func (r *orderRepository) UpdateOrder(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}
func (r *orderRepository) GetUserOrders(ctx context.Context, userID string, page, limit int) ([]*model.Order, int64, error) {
	var orders []*model.Order
	var total int64

	query := r.db.Model(&model.Order{}).Where("user_id = ?", userID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data
	offset := (page - 1) * limit
	err := query.Offset(offset).
		Limit(limit).
		Preload("Items.Product").
		Order("created_at DESC").
		Find(&orders).
		Error

	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func NewOrderRepository() OrderRepository {
	return &orderRepository{db: config.GetDB()}
}

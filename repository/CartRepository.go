package repository

import (
	"context"
	"nuxatech-nextmedis/config"
	"nuxatech-nextmedis/model"
	"time"

	"gorm.io/gorm"
)

type CartRepository interface {
	GetCart(ctx context.Context, userID string) (*model.Cart, error)
	AddItem(ctx context.Context, cartItem *model.CartItem) error
	UpdateItem(ctx context.Context, cartItem *model.CartItem) error
	RemoveItem(ctx context.Context, cartItemID string) error
	GetCartItem(ctx context.Context, cartID string, productID string) (*model.CartItem, error)
}

type cartRepository struct {
	db *gorm.DB
}

func (r *cartRepository) GetCart(ctx context.Context, userID string) (*model.Cart, error) {
	var cart model.Cart
	err := r.db.Preload("Items.Product").Where("user_id = ?", userID).First(&cart).Error
	if err == gorm.ErrRecordNotFound {
		cart = model.Cart{
			UserID:    userID,
			CreatedAt: time.Now().UnixMilli(),
			UpdatedAt: time.Now().UnixMilli(),
		}
		if err := r.db.Create(&cart).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) AddItem(ctx context.Context, cartItem *model.CartItem) error {
	return r.db.Create(cartItem).Error
}

func (r *cartRepository) UpdateItem(ctx context.Context, cartItem *model.CartItem) error {
	return r.db.Save(cartItem).Error
}

func (r *cartRepository) RemoveItem(ctx context.Context, cartItemID string) error {
	return r.db.Delete(&model.CartItem{}, cartItemID).Error
}

func (r *cartRepository) GetCartItem(ctx context.Context, cartID string, productID string) (*model.CartItem, error) {
	var cartItem model.CartItem
	err := r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).First(&cartItem).Error
	if err != nil {
		return nil, err
	}
	return &cartItem, nil
}

func NewCartRepository() CartRepository {
	return &cartRepository{db: config.GetDB()}
}

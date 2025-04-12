package model

import "gorm.io/gorm"

type OrderStatus string

const (
	OrderStatusPending  OrderStatus = "pending"
	OrderStatusPaid     OrderStatus = "paid"
	OrderStatusShipped  OrderStatus = "shipped"
	OrderStatusComplete OrderStatus = "complete"
	OrderStatusCanceled OrderStatus = "canceled"
)

type Order struct {
	ID          string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID      string         `gorm:"type:uuid;not null;index" json:"user_id"`
	CartID      string         `gorm:"type:uuid;not null" json:"cart_id"`
	Status      OrderStatus    `gorm:"type:varchar(20);not null" json:"status"`
	TotalAmount int64          `gorm:"type:bigint;not null" json:"total_amount"`
	Items       []OrderItem    `gorm:"foreignKey:OrderID" json:"items"`
	CreatedAt   int64          `gorm:"type:bigint;not null" json:"created_at"`
	UpdatedAt   int64          `gorm:"type:bigint;not null" json:"updated_at"`
	PaidAt      *int64         `gorm:"type:bigint" json:"paid_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

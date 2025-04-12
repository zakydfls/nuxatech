package model

type OrderItem struct {
	ID        string  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	OrderID   string  `gorm:"type:uuid;not null;index" json:"order_id"`
	ProductID string  `gorm:"type:uuid;not null" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  int     `gorm:"not null" json:"quantity"`
	Price     int64   `gorm:"type:bigint;not null" json:"price"`
	CreatedAt int64   `gorm:"type:bigint;not null" json:"created_at"`
}

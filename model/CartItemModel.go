package model

type CartItem struct {
	ID        string  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CartID    string  `gorm:"type:uuid;not null" json:"cart_id"`
	ProductID string  `gorm:"type:uuid;not null" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  int     `gorm:"type:int;not null" json:"quantity"`
	CreatedAt int64   `gorm:"type:bigint;not null" json:"created_at"`
	UpdatedAt int64   `gorm:"type:bigint;not null" json:"updated_at"`
}

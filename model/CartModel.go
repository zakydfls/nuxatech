package model

type Cart struct {
	ID        string     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID    string     `gorm:"type:uuid;not null" json:"user_id"`
	Items     []CartItem `gorm:"foreignKey:CartID" json:"items"`
	CreatedAt int64      `gorm:"type:bigint;not null" json:"created_at"`
	UpdatedAt int64      `gorm:"type:bigint;not null" json:"updated_at"`
}

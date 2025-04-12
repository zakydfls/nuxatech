package model

import "gorm.io/gorm"

type Transaction struct {
	ID          string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	AccountID   string         `gorm:"type:uuid;not null;index" json:"account_id"`
	Amount      int64          `gorm:"type:bigint;not null" json:"amount"`
	Type        string         `gorm:"type:varchar(20);not null" json:"type"`
	Status      string         `gorm:"type:varchar(20);not null" json:"status"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   int64          `gorm:"type:bigint;not null" json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

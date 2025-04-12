package model

import (
	"sync"

	"gorm.io/gorm"
)

type Account struct {
	ID        string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID    string         `gorm:"type:uuid;not null;index" json:"user_id"`
	Balance   int64          `gorm:"type:bigint;not null;default:0" json:"balance"`
	CreatedAt int64          `gorm:"type:bigint;not null" json:"created_at"`
	UpdatedAt int64          `gorm:"type:bigint;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	mu        sync.RWMutex   `gorm:"-" json:"-"` // Mutex for thread safety
}

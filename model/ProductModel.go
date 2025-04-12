package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type LocalProductOption []string
type LocalProductImages []string

type Product struct {
	ID             string             `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" db:"id" json:"id"`
	Name           string             `gorm:"type:varchar(255);not null" db:"name" json:"name"`
	Description    string             `gorm:"type:text" db:"description" json:"description"`
	Image          LocalProductImages `gorm:"type:jsonb" db:"image" json:"image"`
	Stock          int                `gorm:"type:int" db:"stock" json:"stock"`
	Price          int                `gorm:"type:int" db:"price" json:"price"`
	Weight         int                `gorm:"type:int" db:"weight" json:"weight"`
	BasePrice      int                `gorm:"type:int" db:"base_price" json:"base_price"`
	SKU            string             `gorm:"type:varchar(100)" db:"sku" json:"sku"`
	Slug           string             `gorm:"type:varchar(255)" db:"slug" json:"slug"`
	UniqueCodeType string             `gorm:"type:varchar(100)" db:"unique_code_type" json:"unique_code_type"`
	Sold           bool               `gorm:"type:boolean;default:false" db:"sold" json:"sold"`
	CreatedAt      int64              `gorm:"type:bigint;not null" db:"created_at" json:"created_at"`
	DeletedAt      gorm.DeletedAt     `gorm:"" db:"deleted_at" json:"deleted_at"`
	table          string             `gorm:"-"`
}

func (p Product) TableName() string {
	if p.table != "" {
		return p.table
	}
	return "products"
}

func (lpi LocalProductImages) Value() (driver.Value, error) {
	jsonData, err := json.Marshal(lpi)
	if err != nil {
		return nil, err
	}
	return string(jsonData), nil
}

func (lpi *LocalProductImages) Scan(value interface{}) error {
	if value == nil {
		*lpi = nil
		return nil
	}

	byteValue, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan image: value is not []byte")
	}

	return json.Unmarshal(byteValue, lpi)
}

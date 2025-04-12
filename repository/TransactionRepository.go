package repository

import (
	"context"
	"nuxatech-nextmedis/config"
	"nuxatech-nextmedis/model"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *gorm.DB, transaction *model.Transaction) error
	GetByID(ctx context.Context, id string) (*model.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func (r *transactionRepository) Create(ctx context.Context, tx *gorm.DB, transaction *model.Transaction) error {
	db := tx
	if tx == nil {
		db = r.db
	}
	return db.WithContext(ctx).Create(transaction).Error
}

func (r *transactionRepository) GetByID(ctx context.Context, id string) (*model.Transaction, error) {
	var transaction model.Transaction
	if err := r.db.WithContext(ctx).First(&transaction, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func NewTransactionRepository() TransactionRepository {
	return &transactionRepository{db: config.GetDB()}
}

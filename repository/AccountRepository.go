package repository

import (
	"context"
	"errors"
	"nuxatech-nextmedis/config"
	"nuxatech-nextmedis/model"

	"gorm.io/gorm"
)

type AccountRepository interface {
	BeginTx(ctx context.Context) *gorm.DB
	CreateAccount(ctx context.Context, tx *gorm.DB, account *model.Account) error
	GetAccount(ctx context.Context, id string) (*model.Account, error)
	UpdateBalance(ctx context.Context, tx *gorm.DB, id string, newBalance int64) error
	CreateTransaction(ctx context.Context, transaction *model.Transaction) error
	GetAccountByUserID(ctx context.Context, userID string) (*model.Account, error)
}

type accountRepository struct {
	db *gorm.DB
}

func (r *accountRepository) BeginTx(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Begin()
}

func (r *accountRepository) CreateAccount(ctx context.Context, tx *gorm.DB, account *model.Account) error {
	db := tx
	if tx == nil {
		db = r.db
	}
	return db.WithContext(ctx).Create(account).Error
}

func (r *accountRepository) GetAccount(ctx context.Context, id string) (*model.Account, error) {
	var account model.Account
	if err := r.db.WithContext(ctx).First(&account, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) UpdateBalance(ctx context.Context, tx *gorm.DB, id string, newBalance int64) error {
	db := tx
	if tx == nil {
		db = r.db
	}

	result := db.WithContext(ctx).Model(&model.Account{}).
		Where("id = ?", id).
		Update("balance", newBalance)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("account not found")
	}

	return nil
}

func (r *accountRepository) CreateTransaction(ctx context.Context, transaction *model.Transaction) error {
	return r.db.WithContext(ctx).Create(transaction).Error
}

func (r *accountRepository) GetAccountByUserID(ctx context.Context, userID string) (*model.Account, error) {
	var account model.Account
	if err := r.db.WithContext(ctx).First(&account, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func NewAccountRepository() AccountRepository {
	return &accountRepository{db: config.GetDB()}
}

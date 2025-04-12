package repository

import (
	"context"
	"nuxatech-nextmedis/config"
	"nuxatech-nextmedis/model"

	"gorm.io/gorm"
)

type PersonalTokenRepository interface {
	CreateToken(ctx context.Context, token *model.PersonalToken) error
	FindByID(ctx context.Context, id string) (*model.PersonalToken, error)
	FindByToken(ctx context.Context, token string) (*model.PersonalToken, error)
	DeleteToken(ctx context.Context, id string) error
	DeleteAllUserTokens(ctx context.Context, userID string) error
}

type personalTokenRepository struct {
	db *gorm.DB
}

func NewPersonalTokenRepository() PersonalTokenRepository {
	return &personalTokenRepository{
		db: config.GetDB(),
	}
}

func (r *personalTokenRepository) CreateToken(ctx context.Context, token *model.PersonalToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *personalTokenRepository) FindByID(ctx context.Context, id string) (*model.PersonalToken, error) {
	var token model.PersonalToken
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *personalTokenRepository) FindByToken(ctx context.Context, tokenString string) (*model.PersonalToken, error) {
	var token model.PersonalToken
	err := r.db.WithContext(ctx).Where("token = ?", tokenString).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *personalTokenRepository) DeleteToken(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.PersonalToken{}, "id = ?", id).Error
}

func (r *personalTokenRepository) DeleteAllUserTokens(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Delete(&model.PersonalToken{}, "user_id = ?", userID).Error
}

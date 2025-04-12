package repository

import (
	"context"
	"fmt"
	"nuxatech-nextmedis/config"
	"nuxatech-nextmedis/model"
	"time"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindUsersAfterDate(ctx context.Context, date string) ([]*model.User, error)
	FindById(ctx context.Context, id string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id string) error
	CheckUserExists(ctx context.Context, email string, username string) bool
}

func NewUserRepository() UserRepository {
	return &userRepository{
		db: config.GetDB(),
	}
}
func (u *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	return u.db.WithContext(ctx).Create(user).Error
}

// DeleteUser implements UserRepository.
func (u *userRepository) DeleteUser(ctx context.Context, id string) error {
	return u.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

// FindByEmail implements UserRepository.
func (u *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return &user, err
}

// FindById implements UserRepository.
func (u *userRepository) FindById(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := u.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	return &user, err
}

// FindByUsername implements UserRepository.
func (u *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := u.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	return &user, err
}

// FindUsersAfterDate implements UserRepository.
func (u *userRepository) FindUsersAfterDate(ctx context.Context, date string) ([]*model.User, error) {
	parsedDate, err := time.Parse("02-01-2006", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %v", err)
	}
	unixMillis := parsedDate.UnixMilli()

	var users []*model.User
	err = u.db.WithContext(ctx).Where("created_at > ?", unixMillis).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateUser implements UserRepository.
func (u *userRepository) UpdateUser(ctx context.Context, user *model.User) error {
	return u.db.WithContext(ctx).Model(user).Updates(user).Error
}

func (u *userRepository) CheckUserExists(ctx context.Context, email string, username string) bool {
	var user model.User
	err := u.db.WithContext(ctx).Where("email = ? OR username = ?", email, username).First(&user).Error
	if err == nil {
		return true
	}
	return false
}

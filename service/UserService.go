package service

import (
	"context"
	"errors"
	"fmt"
	"nuxatech-nextmedis/dto/request"
	"nuxatech-nextmedis/model"
	"nuxatech-nextmedis/repository"
	"nuxatech-nextmedis/utils"
	"time"

	"github.com/go-playground/validator/v10"
)

type UserService interface {
	CreateUser(ctx context.Context, req request.CreateUserRequest) error
	GetUser(ctx context.Context, id string) (*model.User, error)
	UpdateUser(ctx context.Context, id string, req request.UpdateUserRequest) error
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindUserAfterDate(ctx context.Context, date string) ([]*model.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type userService struct {
	repo     repository.UserRepository
	validate *validator.Validate
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo:     repo,
		validate: validator.New(),
	}
}

// CreateUser implements UserService.
func (u *userService) CreateUser(ctx context.Context, req request.CreateUserRequest) error {
	if err := u.validate.Struct(req); err != nil {
		return err
	}
	checkUser := u.repo.CheckUserExists(ctx, req.Email, req.Username)
	if checkUser {
		return errors.New("user already exists")
	}
	fmt.Println("cek")
	hashedPassword := utils.HashPassword(req.Password)
	user := model.User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  hashedPassword,
		CreatedAt: time.Now().UnixMilli(),
	}
	return u.repo.CreateUser(ctx, &user)
}

// DeleteUser implements UserService.
func (u *userService) DeleteUser(ctx context.Context, id string) error {
	panic("unimplemented")
}

// FindByUsername implements UserService.
func (u *userService) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user *model.User
	user, err := u.repo.FindByUsername(ctx, username)
	return user, err
}

// FindUserAfterDate implements UserService.
func (u *userService) FindUserAfterDate(ctx context.Context, date string) ([]*model.User, error) {
	var user []*model.User
	user, err := u.repo.FindUsersAfterDate(ctx, date)
	return user, err
}

// FindUserByEmail implements UserService.
func (u *userService) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user *model.User
	user, err := u.repo.FindByEmail(ctx, email)
	return user, err
}

// GetUser implements UserService.
func (u *userService) GetUser(ctx context.Context, id string) (*model.User, error) {
	var user *model.User
	user, err := u.repo.FindById(ctx, id)
	return user, err
}

// UpdateUser implements UserService.
func (u *userService) UpdateUser(ctx context.Context, id string, req request.UpdateUserRequest) error {
	panic("unimplemented")
}

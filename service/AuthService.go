package service

import (
	"context"
	"errors"
	"fmt"
	"nuxatech-nextmedis/config"
	"nuxatech-nextmedis/dto/request"
	"nuxatech-nextmedis/dto/response"
	"nuxatech-nextmedis/model"
	"nuxatech-nextmedis/repository"
	"nuxatech-nextmedis/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, req request.RegisterRequest) (*response.RegisterResponse, error)
	Login(ctx context.Context, req request.LoginRequest) (*response.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*response.TokenResponse, error)
	GenerateTokenPair(user *model.User) (*response.TokenResponse, error)
	ValidateAccessToken(tokenString string) (*request.AccessTokenPayload, error)
	ValidateRefreshToken(tokenString string) (*request.RefreshTokenPayload, error)
	Logout(ctx context.Context, tokenString string) error
}

type authService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.PersonalTokenRepository
	validate  *validator.Validate
}

func NewAuthService(
	userRepo repository.UserRepository,
	tokenRepo repository.PersonalTokenRepository,
) AuthService {
	return &authService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		validate:  validator.New(),
	}
}

func (s *authService) Register(ctx context.Context, req request.RegisterRequest) (*response.RegisterResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	exists := s.userRepo.CheckUserExists(ctx, req.Email, req.Username)
	if exists {
		return nil, errors.New("user already exists")
	}

	hashedPassword := utils.HashPassword(req.Password)
	user := &model.User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  hashedPassword,
		CreatedAt: time.Now().UnixMilli(),
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return &response.RegisterResponse{
		User: response.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}, nil
}

func (s *authService) Login(ctx context.Context, req request.LoginRequest) (*response.LoginResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.VerifyPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	tokens, err := s.GenerateTokenPair(user)
	if err != nil {
		return nil, err
	}

	return &response.LoginResponse{
		User: response.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: time.UnixMilli(user.CreatedAt).Format("02-01-2006 15:04:05"),
		},
		Token: *tokens,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*response.TokenResponse, error) {
	payload, err := s.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	token, err := s.tokenRepo.FindByID(ctx, payload.TokenID)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if token.Token != refreshToken {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.userRepo.FindById(ctx, payload.UserID)
	if err != nil {
		return nil, err
	}

	if err := s.tokenRepo.DeleteToken(ctx, payload.TokenID); err != nil {
		return nil, err
	}

	return s.GenerateTokenPair(user)
}

func (s *authService) GenerateTokenPair(user *model.User) (*response.TokenResponse, error) {
	accessPayload := request.AccessTokenPayload{
		UserID:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(time.Duration(config.Envs.AccessTokenTTL) * time.Second),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    accessPayload.UserID,
		"username":   accessPayload.Username,
		"email":      accessPayload.Email,
		"issued_at":  accessPayload.IssuedAt.Unix(),
		"expired_at": accessPayload.ExpiredAt.Unix(),
	})

	accessTokenString, err := accessToken.SignedString([]byte(config.Envs.JwtAccessSecret))
	if err != nil {
		return nil, err
	}

	refreshTokenID := uuid.New().String()
	refreshPayload := request.RefreshTokenPayload{
		TokenID:   refreshTokenID,
		UserID:    user.ID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(time.Duration(config.Envs.RefreshTokenTTL) * time.Second),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token_id":   refreshPayload.TokenID,
		"user_id":    refreshPayload.UserID,
		"issued_at":  refreshPayload.IssuedAt.Unix(),
		"expired_at": refreshPayload.ExpiredAt.Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(config.Envs.JwtRefreshSecret))
	if err != nil {
		return nil, err
	}

	tokenModel := &model.PersonalToken{
		ID:        refreshTokenID,
		Token:     refreshTokenString,
		UserID:    user.ID,
		CreatedAt: time.Now().UnixMilli(),
	}

	if err := s.tokenRepo.CreateToken(context.Background(), tokenModel); err != nil {
		return nil, err
	}

	return &response.TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func (s *authService) ValidateAccessToken(tokenString string) (*request.AccessTokenPayload, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Envs.JwtAccessSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	expiredAt := time.Unix(int64(claims["expired_at"].(float64)), 0)
	if time.Now().After(expiredAt) {
		return nil, errors.New("token has expired")
	}

	return &request.AccessTokenPayload{
		UserID:    claims["user_id"].(string),
		Username:  claims["username"].(string),
		Email:     claims["email"].(string),
		IssuedAt:  time.Unix(int64(claims["issued_at"].(float64)), 0),
		ExpiredAt: expiredAt,
	}, nil
}

func (s *authService) ValidateRefreshToken(tokenString string) (*request.RefreshTokenPayload, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Envs.JwtRefreshSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	expiredAt := time.Unix(int64(claims["expired_at"].(float64)), 0)
	if time.Now().After(expiredAt) {
		return nil, errors.New("token has expired")
	}

	return &request.RefreshTokenPayload{
		TokenID:   claims["token_id"].(string),
		UserID:    claims["user_id"].(string),
		IssuedAt:  time.Unix(int64(claims["issued_at"].(float64)), 0),
		ExpiredAt: expiredAt,
	}, nil
}

func (s *authService) Logout(ctx context.Context, tokenString string) error {
	payload, err := s.ValidateRefreshToken(tokenString)
	if err != nil {
		return err
	}

	_, err = s.tokenRepo.FindByID(ctx, payload.TokenID)
	if err != nil {
		return errors.New("invalid refresh token")
	}

	if err := s.tokenRepo.DeleteToken(ctx, payload.TokenID); err != nil {
		return err
	}

	return nil
}

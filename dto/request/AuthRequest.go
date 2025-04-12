package request

import "time"

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=2,max=100"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AccessTokenPayload struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

type RefreshTokenPayload struct {
	TokenID   string    `json:"token_id"`
	UserID    string    `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

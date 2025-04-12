package request

type CreateAccountRequest struct {
	UserID string `json:"user_id" validate:"required"`
	// Balance int64  `json:"balance" validate:"gte=0"`
}

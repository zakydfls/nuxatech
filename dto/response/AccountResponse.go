package response

type AccountResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Balance   int64  `json:"balance"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

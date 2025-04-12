package response

type CartItemResponse struct {
	ID       string          `json:"id"`
	Product  ProductResponse `json:"product"`
	Quantity int             `json:"quantity"`
}

type CartResponse struct {
	ID        string             `json:"id"`
	Items     []CartItemResponse `json:"items"`
	Total     int                `json:"total"`
	CreatedAt int64              `json:"created_at"`
	UpdatedAt int64              `json:"updated_at"`
}

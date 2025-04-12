package request

type AddToCartRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" validate:"required,min=1"`
}

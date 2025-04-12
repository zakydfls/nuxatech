package request

// CreateOrderRequest represents the request to create a new order
// @Description Order creation request
type CreateOrderRequest struct {
	// Cart ID from which to create the order
	CartID string `json:"cart_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	// Array of cart item IDs to include in the order
	SelectedItems []string `json:"selected_items" validate:"required" example:"['123e4567-e89b-12d3-a456-426614174000']"`
}

// UpdateOrderStatusRequest represents the request to update order status
// @Description Order status update request
type UpdateOrderStatusRequest struct {
	// New status for the order
	// enum: paid,shipped,complete,canceled
	Status string `json:"status" validate:"required,oneof=paid shipped complete canceled" example:"paid"`
}

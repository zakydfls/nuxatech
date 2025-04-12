package response

import "nuxatech-nextmedis/model"

// OrderItemResponse represents a single item in an order
// @Description Order item details
type OrderItemResponse struct {
	ID       string          `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Product  ProductResponse `json:"product"`
	Quantity int             `json:"quantity" example:"2"`
	Price    int64           `json:"price" example:"150000"`
}

// OrderResponse represents the complete order information
// @Description Complete order information
type OrderResponse struct {
	ID          string              `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Status      string              `json:"status" example:"pending"`
	TotalAmount int64               `json:"total_amount" example:"300000"`
	Items       []OrderItemResponse `json:"items"`
	CreatedAt   int64               `json:"created_at" example:"1617183834"`
	UpdatedAt   int64               `json:"updated_at" example:"1617183834"`
	PaidAt      *int64              `json:"paid_at,omitempty" example:"1617183834"`
}

// OrderPagingResponse represents paginated order results
// @Description Paginated order list
type OrderPagingResponse struct {
	Metadata Metadata       `json:"metadata"`
	Result   []*model.Order `json:"result"`
}

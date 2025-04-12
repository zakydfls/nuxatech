package request

type CreateProductRequest struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description"`
	Image       []string `json:"image"`
	Stock       int      `json:"stock" validate:"required"`
	Price       int      `json:"price" validate:"required"`
	Weight      int      `json:"weight" validate:"required"`
	BasePrice   int      `json:"base_price" validate:"required"`
	SKU         string   `json:"sku"`
	// UniqueCodeType string   `json:"unique_code_type"`
}

type UpdateProductRequest struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description"`
	Image       []string `json:"image"`
	Stock       int      `json:"stock" validate:"required"`
	Price       int      `json:"price" validate:"required"`
	Weight      int      `json:"weight" validate:"required"`
	BasePrice   int      `json:"base_price" validate:"required"`
	SKU         string   `json:"sku"`
	// UniqueCodeType string   `json:"unique_code_type"`
}

package response

import "nuxatech-nextmedis/model"

type ProductResponse struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Image          []string `json:"image"`
	Stock          int      `json:"stock"`
	Price          int      `json:"price"`
	Weight         int      `json:"weight"`
	BasePrice      int      `json:"base_price"`
	SKU            string   `json:"sku"`
	Slug           string   `json:"slug"`
	UniqueCodeType string   `json:"unique_code_type"`
	Sold           bool     `json:"sold"`
	CreatedAt      int64    `json:"created_at"`
}

type PagingResponse struct {
	Metadata Metadata         `json:"metadata"`
	Result   []*model.Product `json:"result"`
}

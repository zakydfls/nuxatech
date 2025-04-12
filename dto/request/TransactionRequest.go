package request

type TransactionRequest struct {
	Amount      int64  `json:"amount" validate:"required,min=1"`
	Description string `json:"description"`
}

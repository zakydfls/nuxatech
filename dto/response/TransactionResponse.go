package response

type TransactionResponse struct {
	ID          string `json:"id"`
	AccountID   string `json:"account_id"`
	Amount      int64  `json:"amount"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Description string `json:"description"`
	CreatedAt   int64  `json:"created_at"`
}

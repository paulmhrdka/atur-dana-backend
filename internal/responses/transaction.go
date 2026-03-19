package responses

type TransactionResponse struct {
	ID          uint    `json:"id"`
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	CategoryID  uint    `json:"category_id"`
	Date        string  `json:"date"`
	CreatedAt   string  `json:"created_at"`
}

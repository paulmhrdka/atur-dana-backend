package responses

type BudgetResponse struct {
	Amount    float64 `json:"amount"`
	Category  string  `json:"category"`
	StartDate string  `json:"start_date"`
	EndDate   string  `json:"end_date"`
}

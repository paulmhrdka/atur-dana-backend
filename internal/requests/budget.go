package requests

type CreateBudgetRequest struct {
	Amount     float64 `json:"amount" validate:"required,gt=1"`
	CategoryID uint    `json:"category_id" validate:"required"`
	StartDate  string  `json:"start_date" validate:"required"`
	EndDate    string  `json:"end_date" validate:"required"`
}

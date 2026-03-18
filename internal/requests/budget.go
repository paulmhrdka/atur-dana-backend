package requests

type CreateBudgetRequest struct {
	Amount     float64 `json:"amount" validate:"required,gt=1" example:"1000000"`
	CategoryID uint    `json:"category_id" validate:"required" example:"2"`
	StartDate  string  `json:"start_date" validate:"required" example:"2026-03-01 00:00:00"`
	EndDate    string  `json:"end_date" validate:"required" example:"2026-03-31 23:59:59"`
}

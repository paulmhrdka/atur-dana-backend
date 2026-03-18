package requests

type CreateTransactionRequest struct {
	Type        string  `json:"type" validate:"required,oneof=income expense" example:"income"`
	Amount      float64 `json:"amount" validate:"required,gt=1" example:"50000"`
	Description string  `json:"description" validate:"required" example:"Monthly salary"`
	CategoryID  uint    `json:"category_id" validate:"required" example:"1"`
	Date        string  `json:"date" validate:"required" example:"2026-03-01 00:00:00"`
}

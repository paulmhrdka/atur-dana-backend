package requests

type CreateTransactionRequest struct {
	Type        string  `json:"type" validate:"required,oneof=income expense"`
	Amount      float64 `json:"amount" validate:"required,gt=1"`
	Description string  `json:"description" validate:"required"`
	CategoryID  uint    `json:"category_id" validate:"required"`
	Date        string  `json:"date" validate:"required"`
}

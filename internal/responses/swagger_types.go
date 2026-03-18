package responses

type SwaggerAuthResponse struct {
	Status  int          `json:"status" example:"200"`
	Message string       `json:"message" example:"Login Successfully"`
	Data    AuthResponse `json:"data"`
}

type SwaggerTransactionListResponse struct {
	Status  int                   `json:"status" example:"200"`
	Message string                `json:"message" example:"Success Get Transactions"`
	Data    []TransactionResponse `json:"data"`
}

type SwaggerBudgetListResponse struct {
	Status  int              `json:"status" example:"200"`
	Message string           `json:"message" example:"Success Get Budgets"`
	Data    []BudgetResponse `json:"data"`
}

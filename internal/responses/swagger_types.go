package responses

import "atur-dana/internal/common"

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

type SwaggerTransactionPaginatedResponse struct {
	Status     int                   `json:"status"     example:"200"`
	Message    string                `json:"message"    example:"Success Get Transactions"`
	Data       []TransactionResponse `json:"data"`
	Pagination common.Pagination     `json:"pagination"`
}

type SwaggerCategoryListResponse struct {
	Status  int                `json:"status"  example:"200"`
	Message string             `json:"message" example:"Success Get Categories"`
	Data    []CategoryResponse `json:"data"`
}

type SwaggerTransactionSummaryResponse struct {
	Status  int             `json:"status"  example:"200"`
	Message string          `json:"message" example:"Success Get Transaction Summary"`
	Data    SummaryResponse `json:"data"`
}

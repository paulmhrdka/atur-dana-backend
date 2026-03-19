package handlers

import (
	"atur-dana/internal/common"
	"atur-dana/internal/db"
	"atur-dana/internal/models"
	"atur-dana/internal/requests"
	"atur-dana/internal/responses"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/context"
)

// GetTransactions godoc
// @Summary      List transactions with filtering and pagination
// @Tags         Transactions
// @Produce      json
// @Security     BearerAuth
// @Param        start_date   query  string  false  "Start date (RFC3339, e.g. 2026-01-01T00:00:00Z)"
// @Param        end_date     query  string  false  "End date (RFC3339)"
// @Param        category_id  query  int     false  "Filter by category ID"
// @Param        type         query  string  false  "Filter by type: income or expense"
// @Param        page         query  int     false  "Page number (default 1)"
// @Param        limit        query  int     false  "Items per page (default 20, max 100)"
// @Param        sort_by      query  string  false  "Sort field: date or amount (default date)"
// @Param        sort_order   query  string  false  "Sort direction: asc or desc (default desc)"
// @Success      200  {object}  responses.SwaggerTransactionPaginatedResponse
// @Failure      400  {object}  common.SwaggerBadRequestResponse
// @Failure      500  {object}  common.SwaggerErrorResponse
// @Router       /api/transactions [get]
func GetTransactions(w http.ResponseWriter, r *http.Request) {
	metadata := context.Get(r, "metadata").(jwt.MapClaims)
	userID := uint(metadata["user_id"].(float64))

	q := r.URL.Query()

	// --- parse optional filters ---
	var startDate, endDate *time.Time
	if s := q.Get("start_date"); s != "" {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			common.JSONError(w, http.StatusBadRequest, "invalid start_date format, use RFC3339")
			return
		}
		startDate = &t
	}
	if s := q.Get("end_date"); s != "" {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			common.JSONError(w, http.StatusBadRequest, "invalid end_date format, use RFC3339")
			return
		}
		endDate = &t
	}

	var categoryID *int
	if s := q.Get("category_id"); s != "" {
		id, err := strconv.Atoi(s)
		if err != nil || id < 1 {
			common.JSONError(w, http.StatusBadRequest, "invalid category_id")
			return
		}
		categoryID = &id
	}

	var txType string
	if s := q.Get("type"); s != "" {
		if s != "income" && s != "expense" {
			common.JSONError(w, http.StatusBadRequest, "type must be income or expense")
			return
		}
		txType = s
	}

	// --- pagination ---
	page := 1
	if s := q.Get("page"); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil || v < 1 {
			common.JSONError(w, http.StatusBadRequest, "invalid page")
			return
		}
		page = v
	}
	limit := 20
	if s := q.Get("limit"); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil || v < 1 {
			common.JSONError(w, http.StatusBadRequest, "invalid limit")
			return
		}
		if v > 100 {
			common.JSONError(w, http.StatusBadRequest, "limit must not exceed 100")
			return
		}
		limit = v
	}

	// --- sorting ---
	sortBy := "date"
	if s := q.Get("sort_by"); s != "" {
		if s != "date" && s != "amount" {
			common.JSONError(w, http.StatusBadRequest, "sort_by must be date or amount")
			return
		}
		sortBy = s
	}
	sortOrder := "desc"
	if s := q.Get("sort_order"); s != "" {
		if s != "asc" && s != "desc" {
			common.JSONError(w, http.StatusBadRequest, "sort_order must be asc or desc")
			return
		}
		sortOrder = s
	}

	// --- build base query ---
	baseQuery := db.DB.Model(&models.Transaction{}).Where("user_id = ?", userID)
	if startDate != nil {
		baseQuery = baseQuery.Where("date >= ?", startDate)
	}
	if endDate != nil {
		baseQuery = baseQuery.Where("date <= ?", endDate)
	}
	if categoryID != nil {
		baseQuery = baseQuery.Where("category_id = ?", *categoryID)
	}
	if txType != "" {
		baseQuery = baseQuery.Where("type = ?", txType)
	}

	// count
	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		common.JSONError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// fetch with pagination
	var transactions []models.Transaction
	if err := baseQuery.Preload("Category").
		Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
		Limit(limit).Offset((page - 1) * limit).
		Find(&transactions).Error; err != nil {
		common.JSONError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	data := make([]responses.TransactionResponse, len(transactions))
	for i, v := range transactions {
		data[i] = responses.TransactionResponse{
			ID:          v.ID,
			Type:        v.Type,
			Amount:      v.Amount,
			Description: v.Description,
			Category:    v.Category.Name,
			CategoryID:  v.CategoryID,
			Date:        v.Date.Format("2006-01-02 15:04:05"),
			CreatedAt:   v.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	pages := int(math.Ceil(float64(total) / float64(limit)))
	pagination := common.Pagination{
		Page:  page,
		Limit: limit,
		Total: total,
		Pages: pages,
	}

	common.JSONPaginatedResponse(w, http.StatusOK, data, pagination, "Success Get Transactions")
}

// GetTransactionSummary godoc
// @Summary      Get transaction summary (totals + category breakdown)
// @Tags         Transactions
// @Produce      json
// @Security     BearerAuth
// @Param        start_date  query  string  true  "Start date (RFC3339)"
// @Param        end_date    query  string  true  "End date (RFC3339)"
// @Success      200  {object}  responses.SwaggerTransactionSummaryResponse
// @Failure      400  {object}  common.SwaggerBadRequestResponse
// @Failure      500  {object}  common.SwaggerErrorResponse
// @Router       /api/transactions/summary [get]
func GetTransactionSummary(w http.ResponseWriter, r *http.Request) {
	metadata := context.Get(r, "metadata").(jwt.MapClaims)
	userID := uint(metadata["user_id"].(float64))

	q := r.URL.Query()

	startStr := q.Get("start_date")
	endStr := q.Get("end_date")
	if startStr == "" || endStr == "" {
		common.JSONError(w, http.StatusBadRequest, "start_date and end_date are required")
		return
	}

	startDate, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		common.JSONError(w, http.StatusBadRequest, "invalid start_date format, use RFC3339")
		return
	}
	endDate, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		common.JSONError(w, http.StatusBadRequest, "invalid end_date format, use RFC3339")
		return
	}

	// --- totals by type ---
	type typeTotal struct {
		Type  string
		Total float64
	}
	var typeTotals []typeTotal
	if err := db.DB.Model(&models.Transaction{}).
		Select("type, SUM(amount) as total").
		Where("user_id = ? AND date >= ? AND date <= ?", userID, startDate, endDate).
		Group("type").
		Scan(&typeTotals).Error; err != nil {
		common.JSONError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	var income, expense float64
	for _, t := range typeTotals {
		if t.Type == "income" {
			income = t.Total
		} else if t.Type == "expense" {
			expense = t.Total
		}
	}

	// --- by category ---
	type catRow struct {
		CategoryID   uint
		CategoryName string
		Type         string
		Total        float64
		Count        int64
	}
	var catRows []catRow
	if err := db.DB.Model(&models.Transaction{}).
		Select("transactions.category_id, categories.name as category_name, transactions.type, SUM(transactions.amount) as total, COUNT(*) as count").
		Joins("JOIN categories ON categories.id = transactions.category_id AND categories.deleted_at IS NULL").
		Where("transactions.user_id = ? AND transactions.date >= ? AND transactions.date <= ?", userID, startDate, endDate).
		Group("transactions.category_id, categories.name, transactions.type").
		Scan(&catRows).Error; err != nil {
		common.JSONError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	byCategory := make([]responses.CategoryBreakdown, len(catRows))
	for i, row := range catRows {
		var typeTotal float64
		if row.Type == "income" {
			typeTotal = income
		} else {
			typeTotal = expense
		}
		var pct float64
		if typeTotal > 0 {
			pct = math.Round((row.Total/typeTotal)*100*100) / 100
		}
		byCategory[i] = responses.CategoryBreakdown{
			CategoryID:   row.CategoryID,
			CategoryName: row.CategoryName,
			Type:         row.Type,
			Total:        row.Total,
			Percentage:   pct,
			Count:        row.Count,
		}
	}

	data := responses.SummaryResponse{
		Period: responses.SummaryPeriod{
			StartDate: startDate.Format("2006-01-02 15:04:05"),
			EndDate:   endDate.Format("2006-01-02 15:04:05"),
		},
		Totals: responses.SummaryTotals{
			Income:  income,
			Expense: expense,
			Net:     income - expense,
		},
		ByCategory: byCategory,
	}

	common.JSONResponse(w, http.StatusOK, data, "Success Get Transaction Summary")
}

// CreateTransaction godoc
// @Summary      Create a transaction
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body requests.CreateTransactionRequest true "Transaction payload"
// @Success      201 {object} common.SwaggerCreatedResponse
// @Failure      400 {object} common.SwaggerValidationErrorResponse
// @Failure      401 {object} common.SwaggerErrorResponse
// @Failure      500 {object} common.SwaggerErrorResponse
// @Router       /api/transactions [post]
func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	metadata := context.Get(r, "metadata").(jwt.MapClaims)
	userID := metadata["user_id"].(float64)

	var request requests.CreateTransactionRequest
	json.NewDecoder(r.Body).Decode(&request)

	// Validate the request
	errValidate := common.ValidateRequest(request)
	if errValidate != nil {
		common.JSONValidationError(w, errValidate)
		return
	}

	trxDate, err := time.Parse("2006-01-02 15:04:05", request.Date)
	if err != nil {
		common.JSONError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	transaction := models.Transaction{
		UserID:      uint(userID),
		Type:        request.Type,
		Amount:      request.Amount,
		CategoryID:  request.CategoryID,
		Description: request.Description,
		Date:        trxDate,
	}

	result := db.DB.Create(&transaction)
	if result.Error != nil {
		common.JSONError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	common.JSONResponse(w, http.StatusCreated, nil, "Transaction Created")
}

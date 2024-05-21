package handlers

import (
	"atur-dana/internal/common"
	"atur-dana/internal/db"
	"atur-dana/internal/models"
	"atur-dana/internal/requests"
	"atur-dana/internal/responses"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/context"
	"gorm.io/gorm"
)

func GetBudgets(w http.ResponseWriter, r *http.Request) {
	metadata := context.Get(r, "metadata").(jwt.MapClaims)
	var budgets []models.Budget

	result := db.DB.Where("user_id = ?", metadata["user_id"]).Preload("Category").Find(&budgets)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			common.JSONError(w, http.StatusUnauthorized, "Invalid credentials")
		} else {
			common.JSONError(w, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	data := make([]responses.BudgetResponse, len(budgets))

	for i, v := range budgets {
		data[i] = responses.BudgetResponse{
			Amount:    v.Amount,
			Category:  v.Category.Name,
			StartDate: v.StartDate.Format("2006-01-02 15:04:05"),
			EndDate:   v.EndDate.Format("2006-01-02 15:04:05"),
		}
	}

	common.JSONResponse(w, http.StatusOK, data, "Success Get Budgets")
}

func CreateBudget(w http.ResponseWriter, r *http.Request) {
	metadata := context.Get(r, "metadata").(jwt.MapClaims)
	userID := metadata["user_id"].(float64)

	var request requests.CreateBudgetRequest
	json.NewDecoder(r.Body).Decode(&request)

	// Validate the request
	errValidate := common.ValidateRequest(request)
	if errValidate != nil {
		common.JSONValidationError(w, errValidate)
		return
	}

	startDate, err := time.Parse("2006-01-02 15:04:05", request.StartDate)
	if err != nil {
		common.JSONError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	endDate, err := time.Parse("2006-01-02 15:04:05", request.EndDate)
	if err != nil {
		common.JSONError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	transaction := models.Budget{
		UserID:     uint(userID),
		Amount:     request.Amount,
		CategoryID: request.CategoryID,
		StartDate:  startDate,
		EndDate:    endDate,
	}

	result := db.DB.Create(&transaction)
	if result.Error != nil {
		common.JSONError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	common.JSONResponse(w, http.StatusCreated, nil, "Budget Created")
}

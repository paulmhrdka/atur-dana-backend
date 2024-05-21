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

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	metadata := context.Get(r, "metadata").(jwt.MapClaims)
	var transactions []models.Transaction

	result := db.DB.Where("user_id = ?", metadata["user_id"]).Preload("Category").Find(&transactions)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			common.JSONError(w, http.StatusUnauthorized, "Invalid credentials")
		} else {
			common.JSONError(w, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	data := make([]responses.TransactionResponse, len(transactions))

	for i, v := range transactions {
		data[i] = responses.TransactionResponse{
			Type:        v.Type,
			Amount:      v.Amount,
			Description: v.Description,
			Category:    v.Category.Name,
			Date:        v.Date.Format("2006-01-02 15:04:05"),
		}
	}

	common.JSONResponse(w, http.StatusOK, data, "Success Get Transactions")
}

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

package common

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func ValidateRequest(req interface{}) map[string]string {
	validate = GetValidator()
	err := validate.Struct(req)
	if err == nil {
		return nil
	}

	validationErrors := err.(validator.ValidationErrors)
	errors := make(map[string]string)
	for _, fieldErr := range validationErrors {
		fieldName := fieldErr.Field()
		errorMessage := getCustomErrorMessage(fieldErr)
		errors[fieldName] = errorMessage
	}
	return errors
}

func getCustomErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", fe.Field(), fe.Param())
	case "email":
		return "Invalid email format"
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", fe.Field(), fe.Param())
	case "oneof":
		return fmt.Sprintf("%s value must be one of %s", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("%s is not valid", fe.Field())
	}
}

func JSONValidationError(w http.ResponseWriter, errors map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "error",
		"errors": errors,
	})
}

func JSONResponse(w http.ResponseWriter, status int, data interface{}, message string) {
	response := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func JSONError(w http.ResponseWriter, status int, message string) {
	JSONResponse(w, status, nil, message)
}

type Pagination struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
	Pages int   `json:"pages"`
}

func JSONPaginatedResponse(w http.ResponseWriter, status int, data interface{}, pagination Pagination, message string) {
	response := map[string]interface{}{
		"status":     status,
		"message":    message,
		"data":       data,
		"pagination": pagination,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

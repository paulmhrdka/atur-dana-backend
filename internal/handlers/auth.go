package handlers

import (
	"atur-dana/internal/auth"
	"atur-dana/internal/common"
	"atur-dana/internal/db"
	"atur-dana/internal/models"
	"atur-dana/internal/requests"
	"atur-dana/internal/responses"
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	// Get Payload & Decode
	var request requests.SignupRequest
	json.NewDecoder(r.Body).Decode(&request)

	// Validate the request
	errValidate := common.ValidateRequest(request)
	if errValidate != nil {
		common.JSONValidationError(w, errValidate)
		return
	}

	// Password Hassing Process
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		common.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Create User Process
	user := models.User{
		Username:     request.Username,
		PasswordHash: string(passwordHash),
		Email:        request.Email,
	}

	result := db.DB.Create(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			common.JSONError(w, http.StatusBadRequest, "User already exists")
		} else {
			common.JSONError(w, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	// Generate JWT Token as Response
	token, err := auth.GenerateJWT(request.Username, user.ID)
	if err != nil {
		common.JSONError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	data := responses.AuthResponse{
		Token: token,
		User: responses.User{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}
	common.JSONResponse(w, http.StatusCreated, data, "Sign Up Successfully")
}

func Login(w http.ResponseWriter, r *http.Request) {
	var request requests.LoginRequest
	json.NewDecoder(r.Body).Decode(&request)

	// Validate the request
	errValidate := common.ValidateRequest(request)
	if errValidate != nil {
		common.JSONValidationError(w, errValidate)
		return
	}

	var dbUser models.User
	result := db.DB.Where("email = ?", request.Email).First(&dbUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			common.JSONError(w, http.StatusUnauthorized, "Invalid credentials")
		} else {
			common.JSONError(w, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(request.Password))
	if err != nil {
		common.JSONError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := auth.GenerateJWT(dbUser.Username, dbUser.ID)
	if err != nil {
		common.JSONError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	data := responses.AuthResponse{
		Token: token,
		User: responses.User{
			ID:        dbUser.ID,
			Username:  dbUser.Username,
			Email:     dbUser.Email,
			CreatedAt: dbUser.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: dbUser.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}
	common.JSONResponse(w, http.StatusOK, data, "Login Successfully")
}

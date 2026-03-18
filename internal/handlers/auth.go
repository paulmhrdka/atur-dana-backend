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

// Register godoc
// @Summary      Register a new user
// @Description  Creates a user account and returns a JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body requests.RegisterRequest true "Register payload"
// @Success      201 {object} responses.SwaggerAuthResponse
// @Failure      400 {object} common.SwaggerValidationErrorResponse
// @Failure      500 {object} common.SwaggerErrorResponse
// @Router       /auth/register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	// Get Payload & Decode
	var request requests.RegisterRequest
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
	common.JSONResponse(w, http.StatusCreated, data, "Register Successfully")
}

// Login godoc
// @Summary      Authenticate a user
// @Description  Validates credentials and returns a JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body requests.LoginRequest true "Login payload"
// @Success      200 {object} responses.SwaggerAuthResponse
// @Failure      400 {object} common.SwaggerValidationErrorResponse
// @Failure      401 {object} common.SwaggerErrorResponse
// @Failure      500 {object} common.SwaggerErrorResponse
// @Router       /auth/login [post]
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
			common.JSONError(w, http.StatusUnauthorized, "User not found!")
		} else {
			common.JSONError(w, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(request.Password))
	if err != nil {
		common.JSONError(w, http.StatusUnauthorized, "User not found!")
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

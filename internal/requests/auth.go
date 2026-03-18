package requests

type RegisterRequest struct {
	Username string `json:"username" validate:"required" example:"johndoe"`
	Password string `json:"password" validate:"required" example:"secret123"`
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	Password string `json:"password" validate:"required" example:"secret123"`
}

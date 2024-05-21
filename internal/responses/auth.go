package responses

type User struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"update_at"`
	Username  string `json:"username"`
	Email     string `json:"email"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

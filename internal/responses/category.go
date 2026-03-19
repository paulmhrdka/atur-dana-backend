package responses

type CategoryResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

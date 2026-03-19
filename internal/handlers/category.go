package handlers

import (
	"atur-dana/internal/common"
	"atur-dana/internal/db"
	"atur-dana/internal/models"
	"atur-dana/internal/responses"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/context"
)

// GetCategories godoc
// @Summary      List categories for the authenticated user
// @Tags         Categories
// @Produce      json
// @Security     BearerAuth
// @Param        active_only  query  bool  false  "Filter active categories only (default true)"
// @Success      200  {object}  responses.SwaggerCategoryListResponse
// @Failure      500  {object}  common.SwaggerErrorResponse
// @Router       /api/categories [get]
func GetCategories(w http.ResponseWriter, r *http.Request) {
	metadata := context.Get(r, "metadata").(jwt.MapClaims)
	userID := uint(metadata["user_id"].(float64))

	activeOnly := r.URL.Query().Get("active_only") != "false" // default true

	query := db.DB.Where("user_id = ?", userID)
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	var categories []models.Category
	if err := query.Find(&categories).Error; err != nil {
		common.JSONError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	data := make([]responses.CategoryResponse, len(categories))
	for i, c := range categories {
		data[i] = responses.CategoryResponse{
			ID:        c.ID,
			Name:      c.Name,
			IsActive:  c.IsActive,
			CreatedAt: c.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	common.JSONResponse(w, http.StatusOK, data, "Success Get Categories")
}

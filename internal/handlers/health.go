package handlers

import (
	"atur-dana/internal/common"
	"atur-dana/internal/db"
	"net/http"
	"time"
)

type healthData struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Service   string `json:"service"`
	Database  string `json:"database"`
}

// HealthCheck godoc
// @Summary      Health check
// @Description  Returns the current health status of the service including database connectivity
// @Tags         Health
// @Produce      json
// @Success      200 {object} common.SwaggerSuccessResponse
// @Failure      503 {object} common.SwaggerErrorResponse
// @Router       /health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	dbStatus := "ok"

	sqlDB, err := db.DB.DB()
	if err != nil || sqlDB.Ping() != nil {
		dbStatus = "unreachable"
	}

	data := healthData{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Service:   "atur-dana",
		Database:  dbStatus,
	}

	if dbStatus != "ok" {
		data.Status = "degraded"
		common.JSONResponse(w, http.StatusServiceUnavailable, data, "Service degraded")
		return
	}

	common.JSONResponse(w, http.StatusOK, data, "OK")
}

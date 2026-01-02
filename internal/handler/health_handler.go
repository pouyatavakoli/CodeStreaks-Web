package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pouyatavakoli/CodeStreaks-web/internal/infrastructure/database"
)

type HealthHandler struct {
	db *database.Database
}

func NewHealthHandler(db *database.Database) *HealthHandler {
	return &HealthHandler{db: db}
}

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

// Health godoc
// @Summary Health check
// @Description Check if the service is healthy
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	dbStatus := "up"
	if err := h.db.HealthCheck(); err != nil {
		dbStatus = "down"
	}

	status := "healthy"
	statusCode := http.StatusOK

	if dbStatus == "down" {
		status = "unhealthy"
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, HealthResponse{
		Status:   status,
		Database: dbStatus,
	})
}

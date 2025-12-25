package handler

import (
	"green-light-backend/internal/transport/http/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck returns the health status of the API
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, dto.SuccessResponse(map[string]string{
		"status": "ok",
	}, "Service is healthy"))
}

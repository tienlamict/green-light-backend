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

// HealthCheck godoc
// @Summary Health check
// @Description Check if the API is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} dto.Response
// @Router /healthz [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, dto.SuccessResponse(map[string]string{
		"status": "ok",
	}, "Service is healthy"))
}

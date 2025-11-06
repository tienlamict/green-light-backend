package middleware

import (
	"green-light-backend/internal/transport/http/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RecoveryMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
				)

				c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Internal server error"))
				c.Abort()
			}
		}()
		c.Next()
	}
}

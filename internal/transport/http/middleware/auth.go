package middleware

import (
	"green-light-backend/internal/transport/http/dto"
	"green-light-backend/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	authorizationType   = "Bearer"
	userIDKey           = "user_id"
	userEmailKey        = "user_email"
	userRoleKey         = "user_role"
)

func AuthMiddleware(jwtManager *utils.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authorizationHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Authorization header is required"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != authorizationType {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid authorization header format"))
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid or expired token"))
			c.Abort()
			return
		}

		// Set user info in context
		c.Set(userIDKey, claims.UserID)
		c.Set(userEmailKey, claims.Email)
		c.Set(userRoleKey, claims.Role)

		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get(userRoleKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse("User role not found"))
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Invalid role type"))
			c.Abort()
			return
		}

		// Check if user has required role
		hasRole := false
		for _, role := range roles {
			if roleStr == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, dto.ErrorResponse("Insufficient permissions"))
			c.Abort()
			return
		}

		c.Next()
	}
}

func GetUserID(c *gin.Context) string {
	userID, _ := c.Get(userIDKey)
	if id, ok := userID.(string); ok {
		return id
	}
	return ""
}

func GetUserRole(c *gin.Context) string {
	userRole, _ := c.Get(userRoleKey)
	if role, ok := userRole.(string); ok {
		return role
	}
	return ""
}

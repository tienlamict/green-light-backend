package handler

import (
	"green-light-backend/internal/transport/http/dto"
	"green-light-backend/internal/transport/http/middleware"
	"green-light-backend/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.Response{data=dto.LoginResponse}
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	input := usecase.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.authUseCase.Login(c.Request.Context(), input)
	if err != nil {
		if err == usecase.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to login"))
		return
	}

	response := dto.LoginResponse{
		Token: output.Token,
		User:  dto.ToUserResponse(output.User),
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(response, "Login successful"))
}

// GetUserInfo godoc
// @Summary Get user information
// @Description Get authenticated user's information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.Response{data=dto.UserResponse}
// @Failure 401 {object} dto.Response
// @Router /api/v1/auth/user-info [get]
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("User not authenticated"))
		return
	}

	user, err := h.authUseCase.GetUserInfo(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to get user info"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(dto.ToUserResponse(user), "User info retrieved"))
}

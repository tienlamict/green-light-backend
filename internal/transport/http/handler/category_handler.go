package handler

import (
	"green-light-backend/internal/domain"
	"green-light-backend/internal/transport/http/dto"
	"green-light-backend/internal/usecase"
	"green-light-backend/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CategoryHandler struct {
	categoryUseCase *usecase.CategoryUseCase
}

func NewCategoryHandler(categoryUseCase *usecase.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{
		categoryUseCase: categoryUseCase,
	}
}

// List returns a paginated list of categories with optional filters
func (h *CategoryHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	sort := c.DefaultQuery("sort", "created_at DESC")

	var isActive *bool
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		val, _ := strconv.ParseBool(isActiveStr)
		isActive = &val
	}

	offset := (page - 1) * limit

	filter := domain.CategoryFilter{
		IsActive: isActive,
		Search:   search,
		Limit:    limit,
		Offset:   offset,
		Sort:     sort,
	}

	categories, total, err := h.categoryUseCase.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch categories"))
		return
	}

	c.JSON(http.StatusOK, dto.PaginatedSuccessResponse(
		dto.ToCategoryListResponse(categories),
		total,
		page,
		limit,
	))
}

// Get returns a category by ID or slug
func (h *CategoryHandler) Get(c *gin.Context) {
	idOrSlug := c.Param("id_or_slug")

	// Try to get by ID first
	category, err := h.categoryUseCase.GetByID(c.Request.Context(), idOrSlug)
	if err != nil && err != gorm.ErrRecordNotFound {
		// Try by slug if ID lookup failed
		category, err = h.categoryUseCase.GetBySlug(c.Request.Context(), idOrSlug)
	}

	if err != nil {
		if err == usecase.ErrCategoryNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse("Category not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch category"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(dto.ToCategoryResponse(category), "Category retrieved"))
}

// Create creates a new category (admin/editor only)
func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	input := usecase.CreateCategoryInput{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		IconURL:     req.IconURL,
		IsActive:    req.IsActive,
	}

	category, err := h.categoryUseCase.Create(c.Request.Context(), input)
	if err != nil {
		// Log error for debugging
		logger.Log.Error("Failed to create category",
			zap.String("error", err.Error()),
			zap.String("name", input.Name),
			zap.String("slug", input.Slug),
		)
		
		if err == usecase.ErrCategorySlugExists {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
			return
		}
		
		// Return detailed error message
		errorMsg := err.Error()
		if errorMsg == "" {
			errorMsg = "Failed to create category"
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(errorMsg))
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse(dto.ToCategoryResponse(category), "Category created"))
}

// Update updates an existing category (admin/editor only)
func (h *CategoryHandler) Update(c *gin.Context) {
	categoryID := c.Param("id")

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	input := usecase.UpdateCategoryInput{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		IconURL:     req.IconURL,
		IsActive:    req.IsActive,
	}

	category, err := h.categoryUseCase.Update(c.Request.Context(), categoryID, input)
	if err != nil {
		if err == usecase.ErrCategoryNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(err.Error()))
			return
		}
		if err == usecase.ErrCategorySlugExists {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to update category"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(dto.ToCategoryResponse(category), "Category updated"))
}

// Delete deletes a category (admin only)
func (h *CategoryHandler) Delete(c *gin.Context) {
	categoryID := c.Param("id")

	err := h.categoryUseCase.Delete(c.Request.Context(), categoryID)
	if err != nil {
		if err == usecase.ErrCategoryNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to delete category"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(nil, "Category deleted"))
}

package handler

import (
	"errors"
	"fmt"
	"green-light-backend/internal/domain"
	"green-light-backend/internal/transport/http/dto"
	"green-light-backend/internal/usecase"
	"green-light-backend/pkg/logger"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProductHandler struct {
	productUseCase *usecase.ProductUseCase
	variantUseCase *usecase.ProductVariantUseCase
	imageRepo      domain.ProductImageRepository
}

func NewProductHandler(productUseCase *usecase.ProductUseCase, variantUseCase *usecase.ProductVariantUseCase, imageRepo domain.ProductImageRepository) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
		variantUseCase: variantUseCase,
		imageRepo:      imageRepo,
	}
}

// loadVariantsWithImages loads variants for a product and their images
func (h *ProductHandler) loadVariantsWithImages(c *gin.Context, productID string) []dto.VariantResponse {
	variants, err := h.variantUseCase.GetByProductID(c.Request.Context(), productID)
	if err != nil || len(variants) == 0 {
		return []dto.VariantResponse{}
	}

	// Load images for each variant
	for _, variant := range variants {
		images, err := h.imageRepo.GetByVariantID(c.Request.Context(), variant.VariantID)
		if err == nil && len(images) > 0 {
			variant.Images = images
		}
	}

	return dto.ToVariantListResponse(variants)
}

// List godoc
// @Summary List products
// @Description Get a list of products with optional filters
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param q query string false "Search by name, description, or SKU"
// @Param category query string false "Filter by category ID"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param is_active query bool false "Filter by active status"
// @Param sort query string false "Sort order" default(created_at DESC)
// @Success 200 {object} dto.PaginatedResponse{data=[]dto.ProductResponse}
// @Router /api/v1/products [get]
func (h *ProductHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("q")
	categoryID := c.Query("category")
	sort := c.DefaultQuery("sort", "created_at DESC")

	var minPrice, maxPrice *float64
	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		val, _ := strconv.ParseFloat(minPriceStr, 64)
		minPrice = &val
	}
	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		val, _ := strconv.ParseFloat(maxPriceStr, 64)
		maxPrice = &val
	}

	var isActive *bool
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		val, _ := strconv.ParseBool(isActiveStr)
		isActive = &val
	}

	offset := (page - 1) * limit

	filter := domain.ProductFilter{
		CategoryID: categoryID,
		Search:     search,
		MinPrice:   minPrice,
		MaxPrice:   maxPrice,
		IsActive:   isActive,
		Limit:      limit,
		Offset:     offset,
		Sort:       sort,
	}

	products, total, err := h.productUseCase.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch products"))
		return
	}

	// Convert products to response and add variants with images for each product
	productResponses := dto.ToProductListResponse(products)
	for i, product := range products {
		productResponses[i].Variants = h.loadVariantsWithImages(c, product.ProductID)
	}

	c.JSON(http.StatusOK, dto.PaginatedSuccessResponse(
		productResponses,
		total,
		page,
		limit,
	))
}

// Get godoc
// @Summary Get product
// @Description Get product by ID or slug
// @Tags products
// @Accept json
// @Produce json
// @Param id_or_slug path string true "Product ID or slug"
// @Success 200 {object} dto.Response{data=dto.ProductResponse}
// @Failure 404 {object} dto.Response
// @Router /api/v1/products/{id_or_slug} [get]
func (h *ProductHandler) Get(c *gin.Context) {
	idOrSlug := c.Param("id_or_slug")

	// Try to get by ID first
	product, err := h.productUseCase.GetByID(c.Request.Context(), idOrSlug)
	if err != nil && err != gorm.ErrRecordNotFound {
		// Try by slug if ID lookup failed
		product, err = h.productUseCase.GetBySlug(c.Request.Context(), idOrSlug)
	}

	if err != nil {
		if err == usecase.ErrProductNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse("Product not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch product"))
		return
	}

	// Get variants with images for this product
	response := dto.ToProductResponse(product)
	response.Variants = h.loadVariantsWithImages(c, product.ProductID)

	c.JSON(http.StatusOK, dto.SuccessResponse(response, "Product retrieved"))
}

// Create godoc
// @Summary Create product
// @Description Create a new product (admin/editor only)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateProductRequest true "Product data"
// @Success 201 {object} dto.Response{data=dto.ProductResponse}
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 403 {object} dto.Response
// @Router /api/v1/products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("Validation error when creating product",
			zap.Error(err),
			zap.Any("request_body", req),
		)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Validation error: "+err.Error()))
		return
	}

	// Convert variant requests to usecase inputs
	variants := make([]usecase.CreateProductVariantInput, len(req.Variants))
	for i, v := range req.Variants {
		variants[i] = usecase.CreateProductVariantInput{
			SKU:        v.SKU,
			Name:       v.Name,
			Attributes: v.Attributes,
			Price:      v.Price,
			Stock:      v.Stock,
			IsActive:   v.IsActive,
		}
	}

	// Handle nullable fields
	stock := 0
	if req.Stock != nil {
		stock = *req.Stock
	}
	sku := ""
	if req.SKU != nil && *req.SKU != "" {
		sku = *req.SKU
	}

	input := usecase.CreateProductInput{
		Name:         req.Name,
		Slug:         req.Slug,
		SKU:          sku,
		ShortDesc:    req.ShortDesc,
		Description:  req.Description,
		Stock:        stock,
		ThumbnailURL: req.ThumbnailURL,
		Gallery:      req.Gallery,
		CategoryID:   req.CategoryID,
		IsActive:     req.IsActive,
		Variants:     variants,
	}

	logger.Log.Info("Creating product",
		zap.String("product_name", req.Name),
		zap.String("category_id", req.CategoryID),
		zap.String("slug", req.Slug),
		zap.Int("variants_count", len(variants)),
	)

	product, err := h.productUseCase.Create(c.Request.Context(), input)
	if err != nil {
		logger.Log.Error("Failed to create product",
			zap.Error(err),
			zap.String("product_name", req.Name),
			zap.String("category_id", req.CategoryID),
			zap.String("error_type", fmt.Sprintf("%T", err)),
			zap.String("error_message", err.Error()),
		)
		if errors.Is(err, usecase.ErrProductSlugExists) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
			return
		}
		if errors.Is(err, usecase.ErrCategoryNotFound) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
			return
		}
		// Check for SKU duplicate error (check error message contains "variant SKU already exists")
		errMsg := err.Error()
		if errMsg != "" && (errors.Is(err, usecase.ErrVariantSKUExists) || 
			strings.Contains(errMsg, "variant SKU already exists")) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
			return
		}
		// Return detailed error message for debugging
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to create product: "+err.Error()))
		return
	}

	// Check if product is nil (should not happen, but safety check)
	if product == nil {
		logger.Log.Error("Product is nil after creation")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to create product: product is nil"))
		return
	}

	// Get variants with images for response
	response := dto.ToProductResponse(product)
	response.Variants = h.loadVariantsWithImages(c, product.ProductID)

	c.JSON(http.StatusCreated, dto.SuccessResponse(response, "Product created"))
}

// Update godoc
// @Summary Update product
// @Description Update an existing product (admin/editor only)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param request body dto.UpdateProductRequest true "Product data"
// @Success 200 {object} dto.Response{data=dto.ProductResponse}
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 403 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Router /api/v1/products/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	productID := c.Param("id_or_slug")

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	input := usecase.UpdateProductInput{
		Name:         req.Name,
		Slug:         req.Slug,
		SKU:          req.SKU,
		ShortDesc:    req.ShortDesc,
		Description:  req.Description,
		Stock:        req.Stock,
		ThumbnailURL: req.ThumbnailURL,
		Gallery:      req.Gallery,
		CategoryID:   req.CategoryID,
		IsActive:     req.IsActive,
	}

	product, err := h.productUseCase.Update(c.Request.Context(), productID, input)
	if err != nil {
		if err == usecase.ErrProductNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(err.Error()))
			return
		}
		if err == usecase.ErrProductSlugExists {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
			return
		}
		if err == usecase.ErrCategoryNotFound {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to update product"))
		return
	}

	// Get variants with images for response
	response := dto.ToProductResponse(product)
	response.Variants = h.loadVariantsWithImages(c, product.ProductID)

	c.JSON(http.StatusOK, dto.SuccessResponse(response, "Product updated"))
}

// Delete godoc
// @Summary Delete product
// @Description Delete a product (admin only)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 403 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Router /api/v1/products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	productID := c.Param("id_or_slug")

	err := h.productUseCase.Delete(c.Request.Context(), productID)
	if err != nil {
		if err == usecase.ErrProductNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to delete product"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(nil, "Product deleted"))
}

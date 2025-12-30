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

// loadProductImages loads product-level images (variant_id IS NULL) and converts to ImageResponse
func (h *ProductHandler) loadProductImages(c *gin.Context, productID string) []dto.ImageResponse {
	productImages, err := h.imageRepo.GetByProductID(c.Request.Context(), productID)
	if err != nil || len(productImages) == 0 {
		return []dto.ImageResponse{}
	}

	images := make([]dto.ImageResponse, len(productImages))
	for i, img := range productImages {
		images[i] = dto.ImageResponse{
			ImageID:   img.ImageID,
			URL:       img.URL,
			IsMain:    img.IsMain,
			SortOrder: img.SortOrder,
		}
	}

	return images
}

// List returns a paginated list of products with optional filters
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
		productResponses[i].Images = h.loadProductImages(c, product.ProductID)
	}

	c.JSON(http.StatusOK, dto.PaginatedSuccessResponse(
		productResponses,
		total,
		page,
		limit,
	))
}

// Get returns a product by ID or slug
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
	response.Images = h.loadProductImages(c, product.ProductID)

	c.JSON(http.StatusOK, dto.SuccessResponse(response, "Product retrieved"))
}

// Create creates a new product (admin/editor only)
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
		// Convert images
		images := make([]usecase.CreateProductVariantImageInput, len(v.Images))
		for j, img := range v.Images {
			images[j] = usecase.CreateProductVariantImageInput{
				URL:       img.URL,
				IsMain:    img.IsMain,
				SortOrder: img.SortOrder,
			}
		}
		
		variants[i] = usecase.CreateProductVariantInput{
			SKU:        v.SKU,
			Name:       v.Name,
			Attributes: v.Attributes,
			Price:      v.Price,
			Stock:      v.Stock,
			IsActive:   v.IsActive,
			Images:     images,
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
	response.Images = h.loadProductImages(c, product.ProductID)

	c.JSON(http.StatusCreated, dto.SuccessResponse(response, "Product created"))
}

// Update updates an existing product (admin/editor only)
func (h *ProductHandler) Update(c *gin.Context) {
	productID := c.Param("id_or_slug")

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	// Convert variant requests to usecase inputs
	var variants []usecase.UpdateProductVariantInput
	if len(req.Variants) > 0 {
		variants = make([]usecase.UpdateProductVariantInput, 0, len(req.Variants))
		for _, v := range req.Variants {
			// Variant ID is required when updating via product API
			if v.VariantID == nil || *v.VariantID == "" {
				c.JSON(http.StatusBadRequest, dto.ErrorResponse("variant_id is required for each variant"))
				return
			}
			
			variants = append(variants, usecase.UpdateProductVariantInput{
				VariantID:  *v.VariantID,
				SKU:        v.SKU,
				Name:       v.Name,
				Attributes: v.Attributes,
				Price:      v.Price,
				Stock:      v.Stock,
				IsActive:   v.IsActive,
			})
		}
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
		Variants:     variants,
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
		// Check for variant-related errors
		if strings.Contains(err.Error(), "variant") {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to update product"))
		return
	}

	// Get variants with images for response
	response := dto.ToProductResponse(product)
	response.Variants = h.loadVariantsWithImages(c, product.ProductID)
	response.Images = h.loadProductImages(c, product.ProductID)

	c.JSON(http.StatusOK, dto.SuccessResponse(response, "Product updated"))
}

// Delete deletes a product (admin only)
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

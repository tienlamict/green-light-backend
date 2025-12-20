package handler

import (
	"green-light-backend/internal/domain"
	"green-light-backend/internal/transport/http/dto"
	"green-light-backend/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductHandler struct {
	productUseCase *usecase.ProductUseCase
	variantUseCase *usecase.ProductVariantUseCase
}

func NewProductHandler(productUseCase *usecase.ProductUseCase, variantUseCase *usecase.ProductVariantUseCase) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
		variantUseCase: variantUseCase,
	}
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

	c.JSON(http.StatusOK, dto.PaginatedSuccessResponse(
		dto.ToProductListResponse(products),
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

	// Get variants for this product
	variants, _ := h.variantUseCase.GetByProductID(c.Request.Context(), product.ProductID)

	response := dto.ToProductResponse(product)
	if len(variants) > 0 {
		response.Variants = dto.ToVariantListResponse(variants)
	}

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
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
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

	input := usecase.CreateProductInput{
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

	product, err := h.productUseCase.Create(c.Request.Context(), input)
	if err != nil {
		if err == usecase.ErrProductSlugExists {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
			return
		}
		if err == usecase.ErrCategoryNotFound {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to create product"))
		return
	}

	// Get variants for response
	fetchedVariants, _ := h.variantUseCase.GetByProductID(c.Request.Context(), product.ProductID)

	response := dto.ToProductResponse(product)
	if len(fetchedVariants) > 0 {
		response.Variants = dto.ToVariantListResponse(fetchedVariants)
	}

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

	c.JSON(http.StatusOK, dto.SuccessResponse(dto.ToProductResponse(product), "Product updated"))
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

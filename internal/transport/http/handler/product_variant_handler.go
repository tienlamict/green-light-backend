package handler

import (
	"green-light-backend/internal/transport/http/dto"
	"green-light-backend/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductVariantHandler struct {
	variantUseCase *usecase.ProductVariantUseCase
}

func NewProductVariantHandler(variantUseCase *usecase.ProductVariantUseCase) *ProductVariantHandler {
	return &ProductVariantHandler{
		variantUseCase: variantUseCase,
	}
}

// ListByProduct godoc
// @Summary List product variants
// @Description Get all variants for a specific product
// @Tags product-variants
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Success 200 {object} dto.Response{data=[]dto.VariantResponse}
// @Failure 404 {object} dto.Response
// @Router /api/v1/products/{product_id}/variants [get]
func (h *ProductVariantHandler) ListByProduct(c *gin.Context) {
	productID := c.Param("product_id")

	variants, err := h.variantUseCase.GetByProductID(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch variants"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(dto.ToVariantListResponse(variants), "Variants retrieved"))
}

// Get godoc
// @Summary Get variant
// @Description Get variant by ID
// @Tags product-variants
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Param variant_id path string true "Variant ID"
// @Success 200 {object} dto.Response{data=dto.VariantResponse}
// @Failure 404 {object} dto.Response
// @Router /api/v1/products/{product_id}/variants/{variant_id} [get]
func (h *ProductVariantHandler) Get(c *gin.Context) {
	variantID := c.Param("variant_id")

	variant, err := h.variantUseCase.GetByID(c.Request.Context(), variantID)
	if err != nil {
		if err == usecase.ErrVariantNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse("Variant not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch variant"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(dto.ToVariantResponse(variant), "Variant retrieved"))
}

// GetBySKU godoc
// @Summary Get variant by SKU
// @Description Get variant by SKU
// @Tags product-variants
// @Accept json
// @Produce json
// @Param sku path string true "Variant SKU"
// @Success 200 {object} dto.Response{data=dto.VariantResponse}
// @Failure 404 {object} dto.Response
// @Router /api/v1/variants/sku/{sku} [get]
func (h *ProductVariantHandler) GetBySKU(c *gin.Context) {
	sku := c.Param("sku")

	variant, err := h.variantUseCase.GetBySKU(c.Request.Context(), sku)
	if err != nil {
		if err == usecase.ErrVariantNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse("Variant not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch variant"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(dto.ToVariantResponse(variant), "Variant retrieved"))
}

// Create godoc
// @Summary Create variant
// @Description Create a new product variant (admin/editor only)
// @Tags product-variants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product_id path string true "Product ID"
// @Param request body dto.CreateVariantRequest true "Variant data"
// @Success 201 {object} dto.Response{data=dto.VariantResponse}
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 403 {object} dto.Response
// @Router /api/v1/products/{product_id}/variants [post]
func (h *ProductVariantHandler) Create(c *gin.Context) {
	productID := c.Param("product_id")

	var req dto.CreateVariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	input := usecase.CreateVariantInput{
		ProductID:  productID,
		SKU:        req.SKU,
		Name:       req.Name,
		Attributes: req.Attributes,
		Price:      req.Price,
		Stock:      req.Stock,
		IsActive:   req.IsActive,
	}

	variant, err := h.variantUseCase.Create(c.Request.Context(), input)
	if err != nil {
		if err == usecase.ErrVariantSKUExists {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
			return
		}
		if err == usecase.ErrProductNotFound {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to create variant"))
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse(dto.ToVariantResponse(variant), "Variant created"))
}

// Update godoc
// @Summary Update variant
// @Description Update an existing product variant (admin/editor only)
// @Tags product-variants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product_id path string true "Product ID"
// @Param variant_id path string true "Variant ID"
// @Param request body dto.UpdateVariantRequest true "Variant data"
// @Success 200 {object} dto.Response{data=dto.VariantResponse}
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 403 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Router /api/v1/products/{product_id}/variants/{variant_id} [put]
func (h *ProductVariantHandler) Update(c *gin.Context) {
	variantID := c.Param("variant_id")

	var req dto.UpdateVariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	input := usecase.UpdateVariantInput{
		SKU:        req.SKU,
		Name:       req.Name,
		Attributes: req.Attributes,
		Price:      req.Price,
		Stock:      req.Stock,
		IsActive:   req.IsActive,
	}

	variant, err := h.variantUseCase.Update(c.Request.Context(), variantID, input)
	if err != nil {
		if err == usecase.ErrVariantNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(err.Error()))
			return
		}
		if err == usecase.ErrVariantSKUExists {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to update variant"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(dto.ToVariantResponse(variant), "Variant updated"))
}

// Delete godoc
// @Summary Delete variant
// @Description Delete a product variant (admin only)
// @Tags product-variants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product_id path string true "Product ID"
// @Param variant_id path string true "Variant ID"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 403 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Router /api/v1/products/{product_id}/variants/{variant_id} [delete]
func (h *ProductVariantHandler) Delete(c *gin.Context) {
	variantID := c.Param("variant_id")

	err := h.variantUseCase.Delete(c.Request.Context(), variantID)
	if err != nil {
		if err == usecase.ErrVariantNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to delete variant"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(nil, "Variant deleted"))
}


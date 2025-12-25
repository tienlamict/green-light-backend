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

// ListByProduct returns all variants for a specific product
func (h *ProductVariantHandler) ListByProduct(c *gin.Context) {
	productID := c.Param("id_or_slug")

	variants, err := h.variantUseCase.GetByProductID(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to fetch variants"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(dto.ToVariantListResponse(variants), "Variants retrieved"))
}

// Get returns a variant by ID
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

// GetBySKU returns a variant by SKU
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

// Create creates a new product variant (admin/editor only)
func (h *ProductVariantHandler) Create(c *gin.Context) {
	productID := c.Param("id_or_slug")

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

// Update updates an existing product variant (admin/editor only)
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

// Delete deletes a product variant (admin only)
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

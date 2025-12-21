package handler

import (
	"net/http"

	"green-light-backend/internal/transport/http/dto"
	"green-light-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

// ProductImageHandler handles HTTP requests for product images
type ProductImageHandler struct {
	imageUseCase *usecase.ProductImageUseCase
}

// NewProductImageHandler creates a new product image handler
func NewProductImageHandler(imageUseCase *usecase.ProductImageUseCase) *ProductImageHandler {
	return &ProductImageHandler{
		imageUseCase: imageUseCase,
	}
}

// GeneratePresignedURL generates a presigned URL for direct upload
// POST /api/v1/products/:id_or_slug/images/presign
func (h *ProductImageHandler) GeneratePresignedURL(c *gin.Context) {
	productID := c.Param("id_or_slug")

	var req dto.PresignUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	// Generate presigned URL
	presignedResp, err := h.imageUseCase.GeneratePresignedUploadURL(c.Request.Context(), usecase.GeneratePresignedUploadURLInput{
		ProductID:   productID,
		ContentType: req.ContentType,
		Extension:   req.Extension,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(dto.PresignUploadResponse{
		UploadURL: presignedResp.UploadURL,
		PublicURL: presignedResp.PublicURL,
		ObjectKey: presignedResp.ObjectKey,
		ExpiresAt: presignedResp.ExpiresAt,
	}, "Presigned URL generated successfully"))
}

// ConfirmImageUpload confirms that an image was uploaded successfully
// POST /api/v1/products/:id_or_slug/images
func (h *ProductImageHandler) ConfirmImageUpload(c *gin.Context) {
	productID := c.Param("id_or_slug")

	var req dto.ConfirmImageUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	// Confirm image upload
	image, err := h.imageUseCase.ConfirmImageUpload(c.Request.Context(), usecase.ConfirmImageUploadInput{
		ProductID: productID,
		ObjectKey: req.ObjectKey,
		IsMain:    req.IsMain,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse(dto.ToProductImageResponse(image), "Image uploaded successfully"))
}

// ListProductImages returns all images for a product
// GET /api/v1/products/:id_or_slug/images
func (h *ProductImageHandler) ListProductImages(c *gin.Context) {
	productID := c.Param("id_or_slug")

	images, err := h.imageUseCase.GetProductImages(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(dto.ToProductImageResponseList(images), "Product images retrieved successfully"))
}

// DeleteImage deletes an image
// DELETE /api/v1/products/:id_or_slug/images/:image_id
func (h *ProductImageHandler) DeleteImage(c *gin.Context) {
	imageID := c.Param("image_id")

	if err := h.imageUseCase.DeleteProductImage(c.Request.Context(), imageID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(nil, "Image deleted successfully"))
}

// UpdateImageMetadata updates image metadata (is_main, sort_order)
// PATCH /api/v1/products/:id_or_slug/images/:image_id
func (h *ProductImageHandler) UpdateImageMetadata(c *gin.Context) {
	imageID := c.Param("image_id")

	var req dto.UpdateImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	// Update image metadata
	image, err := h.imageUseCase.UpdateImageMetadata(c.Request.Context(), usecase.UpdateImageMetadataInput{
		ImageID:   imageID,
		IsMain:    req.IsMain,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(dto.ToProductImageResponse(image), "Image metadata updated successfully"))
}


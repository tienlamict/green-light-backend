package dto

import "green-light-backend/internal/domain"

// PresignUploadRequest is the request to get a presigned upload URL
type PresignUploadRequest struct {
	ContentType string `json:"content_type" binding:"required,oneof=image/jpeg image/jpg image/png image/webp"`
	Extension   string `json:"extension" binding:"required,oneof=jpg jpeg png webp"`
}

// PresignUploadResponse is the response containing presigned URL
type PresignUploadResponse struct {
	UploadURL string `json:"upload_url"` // Presigned PUT URL for direct upload
	PublicURL string `json:"public_url"` // Public URL to access the image
	ObjectKey string `json:"object_key"` // MinIO object key
	ExpiresAt string `json:"expires_at"` // Expiration time
}

// ConfirmImageUploadRequest is the request to confirm image upload
type ConfirmImageUploadRequest struct {
	ObjectKey  string  `json:"object_key" binding:"required"`
	VariantID  *string `json:"variant_id,omitempty"` // Optional: if provided, image belongs to variant
	IsMain     bool    `json:"is_main"`
	SortOrder  int     `json:"sort_order"`
}

// ProductImageResponse is the response for a product image
type ProductImageResponse struct {
	ImageID   string `json:"image_id"`
	ProductID string `json:"product_id"`
	URL       string `json:"url"`
	IsMain    bool   `json:"is_main"`
	SortOrder int    `json:"sort_order"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

// ToProductImageResponse converts domain model to response DTO
func ToProductImageResponse(image *domain.ProductImage) ProductImageResponse {
	return ProductImageResponse{
		ImageID:   image.ImageID,
		ProductID: image.ProductID,
		URL:       image.URL,
		IsMain:    image.IsMain,
		SortOrder: image.SortOrder,
		Status:    image.Status,
		CreatedAt: image.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ToProductImageResponseList converts a list of domain models to response DTOs
func ToProductImageResponseList(images []*domain.ProductImage) []ProductImageResponse {
	result := make([]ProductImageResponse, len(images))
	for i, image := range images {
		result[i] = ToProductImageResponse(image)
	}
	return result
}

// UpdateImageRequest is the request to update image metadata
type UpdateImageRequest struct {
	IsMain    *bool `json:"is_main"`
	SortOrder *int  `json:"sort_order"`
}


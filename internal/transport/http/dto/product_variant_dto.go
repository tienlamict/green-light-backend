package dto

import "green-light-backend/internal/domain"

// CreateVariantRequest represents the request to create a product variant
type CreateVariantRequest struct {
	SKU        string            `json:"sku" binding:"required"`
	Name       string            `json:"name" binding:"required"`
	Attributes map[string]string `json:"attributes"`
	Price      float64           `json:"price" binding:"required,gt=0"` // Price is required for variants
	Stock      int               `json:"stock" binding:"min=0"`
	IsActive   bool              `json:"is_active"`
}

// UpdateVariantRequest represents the request to update a product variant
type UpdateVariantRequest struct {
	SKU        *string            `json:"sku"`
	Name       *string            `json:"name"`
	Attributes *map[string]string `json:"attributes"`
	Price      *float64           `json:"price" binding:"omitempty,gt=0"` // Price must be > 0 if provided
	Stock      *int               `json:"stock" binding:"omitempty,min=0"`
	IsActive   *bool              `json:"is_active"`
}

// VariantResponse represents the response for a product variant
type VariantResponse struct {
	VariantID  string            `json:"variant_id"`
	ProductID  string            `json:"product_id"`
	SKU        string            `json:"sku"`
	Name       string            `json:"name"`
	Attributes map[string]string `json:"attributes"`
	Price      float64           `json:"price"` // Price is always present for variants
	Stock      int               `json:"stock"`
	IsActive   bool              `json:"is_active"`
	Images     []ImageResponse   `json:"images,omitempty"` // Variant-specific images
	CreatedAt  string            `json:"created_at"`
	UpdatedAt  string            `json:"updated_at"`
}

// ImageResponse represents an image
type ImageResponse struct {
	ImageID   string `json:"image_id"`
	URL       string `json:"url"`
	IsMain    bool   `json:"is_main"`
	SortOrder int    `json:"sort_order"`
}

// ToVariantResponse converts domain.ProductVariant to VariantResponse
func ToVariantResponse(variant *domain.ProductVariant) VariantResponse {
	response := VariantResponse{
		VariantID:  variant.VariantID,
		ProductID:  variant.ProductID,
		SKU:        variant.SKU,
		Name:       variant.Name,
		Attributes: variant.Attributes,
		Price:      variant.Price,
		Stock:      variant.Stock,
		IsActive:   variant.IsActive,
		CreatedAt:  variant.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  variant.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	
	// Add images if available
	if len(variant.Images) > 0 {
		response.Images = make([]ImageResponse, len(variant.Images))
		for i, img := range variant.Images {
			response.Images[i] = ImageResponse{
				ImageID:   img.ImageID,
				URL:       img.URL,
				IsMain:    img.IsMain,
				SortOrder: img.SortOrder,
			}
		}
	}
	
	return response
}

// ToVariantListResponse converts slice of domain.ProductVariant to slice of VariantResponse
func ToVariantListResponse(variants []*domain.ProductVariant) []VariantResponse {
	responses := make([]VariantResponse, len(variants))
	for i, variant := range variants {
		responses[i] = ToVariantResponse(variant)
	}
	return responses
}


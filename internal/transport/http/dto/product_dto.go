package dto

import "green-light-backend/internal/domain"

type CreateProductRequest struct {
	Name         string                 `json:"name" binding:"required,min=2,max=255"`
	Slug         string                 `json:"slug" binding:"omitempty,max=255"`
	SKU          *string                `json:"sku" binding:"omitempty,max=100"` // Optional base SKU (nullable)
	ShortDesc    string                 `json:"short_desc" binding:"max=500"`
	Description  string                 `json:"description"`
	Stock        *int                   `json:"stock" binding:"omitempty,gte=0"` // Optional, defaults to 0 if null
	ThumbnailURL string                 `json:"thumbnail_url"`
	Gallery      []string               `json:"gallery"`
	CategoryID   string                 `json:"category_id" binding:"required"`
	IsActive     bool                   `json:"is_active"`
	Variants     []CreateVariantRequest `json:"variants" binding:"required,min=1"` // At least 1 variant required
}

type UpdateProductRequest struct {
	Name         *string                 `json:"name" binding:"omitempty,min=2,max=255"`
	Slug         *string                 `json:"slug" binding:"omitempty,max=255"`
	SKU          *string                 `json:"sku" binding:"omitempty,max=100"`
	ShortDesc    *string                 `json:"short_desc" binding:"omitempty,max=500"`
	Description  *string                 `json:"description"`
	Stock        *int                    `json:"stock" binding:"omitempty,gte=0"`
	ThumbnailURL *string                 `json:"thumbnail_url"`
	Gallery      *[]string               `json:"gallery"`
	CategoryID   *string                 `json:"category_id"`
	IsActive     *bool                   `json:"is_active"`
	Variants     []UpdateVariantRequest  `json:"variants"` // Optional: update variants in the same request
}

type ProductResponse struct {
	ProductID    string            `json:"product_id"`
	Name         string            `json:"name"`
	Slug         string            `json:"slug"`
	SKU          string            `json:"sku"`
	ShortDesc    string            `json:"short_desc"`
	Description  string            `json:"description"`
	PriceMin     *float64          `json:"price_min"` // Min price from variants
	PriceMax     *float64          `json:"price_max"` // Max price from variants
	Stock        int               `json:"stock"`
	ThumbnailURL string            `json:"thumbnail_url"`
	Gallery      []string          `json:"gallery"`
	CategoryID   string            `json:"category_id"`
	Category     *CategoryResponse `json:"category,omitempty"`
	IsActive     bool              `json:"is_active"`
	Images       []ImageResponse   `json:"images,omitempty"` // Product-level images
	Variants     []VariantResponse `json:"variants,omitempty"` // Include variants if available
	CreatedAt    string            `json:"created_at"`
	UpdatedAt    string            `json:"updated_at"`
}

func ToProductResponse(product *domain.Product) ProductResponse {
	if product == nil {
		// Return empty response if product is nil
		return ProductResponse{}
	}
	
	response := ProductResponse{
		ProductID:    product.ProductID,
		Name:         product.Name,
		Slug:         product.Slug,
		SKU:          product.SKU,
		ShortDesc:    product.ShortDesc,
		Description:  product.Description,
		PriceMin:     product.PriceMin,
		PriceMax:     product.PriceMax,
		Stock:        product.Stock,
		ThumbnailURL: product.ThumbnailURL,
		Gallery:      []string(product.Gallery),
		CategoryID:   product.CategoryID,
		IsActive:     product.IsActive,
		CreatedAt:    product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if product.Category != nil {
		cat := ToCategoryResponse(product.Category)
		response.Category = &cat
	}

	return response
}

func ToProductListResponse(products []*domain.Product) []ProductResponse {
	response := make([]ProductResponse, len(products))
	for i, prod := range products {
		response[i] = ToProductResponse(prod)
	}
	return response
}

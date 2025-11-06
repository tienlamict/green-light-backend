package dto

import "green-light-backend/internal/domain"

type CreateProductRequest struct {
	Name         string   `json:"name" binding:"required,min=2,max=255"`
	Slug         string   `json:"slug" binding:"omitempty,max=255"`
	SKU          string   `json:"sku" binding:"required,max=100"`
	ShortDesc    string   `json:"short_desc" binding:"max=500"`
	Description  string   `json:"description"`
	Price        float64  `json:"price" binding:"required,gt=0"`
	Stock        int      `json:"stock" binding:"gte=0"`
	ThumbnailURL string   `json:"thumbnail_url"`
	Gallery      []string `json:"gallery"`
	CategoryID   string   `json:"category_id" binding:"required"`
	IsActive     bool     `json:"is_active"`
}

type UpdateProductRequest struct {
	Name         *string   `json:"name" binding:"omitempty,min=2,max=255"`
	Slug         *string   `json:"slug" binding:"omitempty,max=255"`
	SKU          *string   `json:"sku" binding:"omitempty,max=100"`
	ShortDesc    *string   `json:"short_desc" binding:"omitempty,max=500"`
	Description  *string   `json:"description"`
	Price        *float64  `json:"price" binding:"omitempty,gt=0"`
	Stock        *int      `json:"stock" binding:"omitempty,gte=0"`
	ThumbnailURL *string   `json:"thumbnail_url"`
	Gallery      *[]string `json:"gallery"`
	CategoryID   *string   `json:"category_id"`
	IsActive     *bool     `json:"is_active"`
}

type ProductResponse struct {
	ProductID    string            `json:"product_id"`
	Name         string            `json:"name"`
	Slug         string            `json:"slug"`
	SKU          string            `json:"sku"`
	ShortDesc    string            `json:"short_desc"`
	Description  string            `json:"description"`
	Price        float64           `json:"price"`
	Stock        int               `json:"stock"`
	ThumbnailURL string            `json:"thumbnail_url"`
	Gallery      []string          `json:"gallery"`
	CategoryID   string            `json:"category_id"`
	Category     *CategoryResponse `json:"category,omitempty"`
	IsActive     bool              `json:"is_active"`
	CreatedAt    string            `json:"created_at"`
	UpdatedAt    string            `json:"updated_at"`
}

func ToProductResponse(product *domain.Product) ProductResponse {
	response := ProductResponse{
		ProductID:    product.ProductID,
		Name:         product.Name,
		Slug:         product.Slug,
		SKU:          product.SKU,
		ShortDesc:    product.ShortDesc,
		Description:  product.Description,
		Price:        product.Price,
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

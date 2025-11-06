package dto

import "green-light-backend/internal/domain"

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=255"`
	Slug        string `json:"slug" binding:"omitempty,max=255"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=2,max=255"`
	Slug        *string `json:"slug" binding:"omitempty,max=255"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

type CategoryResponse struct {
	CategoryID  string `json:"category_id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func ToCategoryResponse(category *domain.Category) CategoryResponse {
	return CategoryResponse{
		CategoryID:  category.CategoryID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func ToCategoryListResponse(categories []*domain.Category) []CategoryResponse {
	response := make([]CategoryResponse, len(categories))
	for i, cat := range categories {
		response[i] = ToCategoryResponse(cat)
	}
	return response
}

package usecase

import (
	"context"
	"errors"
	"green-light-backend/internal/domain"
	"green-light-backend/pkg/utils"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrCategoryNotFound    = errors.New("category not found")
	ErrCategorySlugExists  = errors.New("category slug already exists")
	ErrInvalidCategoryData = errors.New("invalid category data")
)

type CategoryUseCase struct {
	categoryRepo domain.CategoryRepository
}

func NewCategoryUseCase(categoryRepo domain.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{
		categoryRepo: categoryRepo,
	}
}

type CreateCategoryInput struct {
	Name        string
	Slug        string
	Description string
	IsActive    bool
}

type UpdateCategoryInput struct {
	Name        *string
	Slug        *string
	Description *string
	IsActive    *bool
}

func (uc *CategoryUseCase) Create(ctx context.Context, input CreateCategoryInput) (*domain.Category, error) {
	// Generate slug if not provided
	slug := input.Slug
	if slug == "" {
		slug = generateSlug(input.Name)
	}

	// Check if slug already exists
	_, err := uc.categoryRepo.GetBySlug(ctx, slug)
	if err == nil {
		return nil, ErrCategorySlugExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	category := &domain.Category{
		CategoryID:  utils.GenerateUUIDv7(),
		Name:        input.Name,
		Slug:        slug,
		Description: input.Description,
		IsActive:    input.IsActive,
	}

	if err := uc.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (uc *CategoryUseCase) GetByID(ctx context.Context, categoryID string) (*domain.Category, error) {
	category, err := uc.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return category, nil
}

func (uc *CategoryUseCase) GetBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	category, err := uc.categoryRepo.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return category, nil
}

func (uc *CategoryUseCase) Update(ctx context.Context, categoryID string, input UpdateCategoryInput) (*domain.Category, error) {
	category, err := uc.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	// Update fields if provided
	if input.Name != nil {
		category.Name = *input.Name
	}
	if input.Slug != nil {
		// Check if new slug already exists
		existingCat, err := uc.categoryRepo.GetBySlug(ctx, *input.Slug)
		if err == nil && existingCat.CategoryID != categoryID {
			return nil, ErrCategorySlugExists
		}
		category.Slug = *input.Slug
	}
	if input.Description != nil {
		category.Description = *input.Description
	}
	if input.IsActive != nil {
		category.IsActive = *input.IsActive
	}

	if err := uc.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (uc *CategoryUseCase) Delete(ctx context.Context, categoryID string) error {
	// Check if category exists
	_, err := uc.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCategoryNotFound
		}
		return err
	}

	return uc.categoryRepo.Delete(ctx, categoryID)
}

func (uc *CategoryUseCase) List(ctx context.Context, filter domain.CategoryFilter) ([]*domain.Category, int64, error) {
	return uc.categoryRepo.List(ctx, filter)
}

// generateSlug creates a URL-friendly slug from a string
func generateSlug(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	// Remove special characters (keep alphanumeric and hyphens)
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

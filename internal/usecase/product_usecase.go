package usecase

import (
	"context"
	"errors"
	"green-light-backend/internal/domain"
	"green-light-backend/pkg/utils"

	"gorm.io/gorm"
)

var (
	ErrProductNotFound    = errors.New("product not found")
	ErrProductSlugExists  = errors.New("product slug already exists")
	ErrInvalidProductData = errors.New("invalid product data")
)

type ProductUseCase struct {
	productRepo  domain.ProductRepository
	categoryRepo domain.CategoryRepository
}

func NewProductUseCase(productRepo domain.ProductRepository, categoryRepo domain.CategoryRepository) *ProductUseCase {
	return &ProductUseCase{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

type CreateProductInput struct {
	Name         string
	Slug         string
	SKU          string
	ShortDesc    string
	Description  string
	Price        float64
	Stock        int
	ThumbnailURL string
	Gallery      []string
	CategoryID   string
	IsActive     bool
}

type UpdateProductInput struct {
	Name         *string
	Slug         *string
	SKU          *string
	ShortDesc    *string
	Description  *string
	Price        *float64
	Stock        *int
	ThumbnailURL *string
	Gallery      *[]string
	CategoryID   *string
	IsActive     *bool
}

func (uc *ProductUseCase) Create(ctx context.Context, input CreateProductInput) (*domain.Product, error) {
	// Validate category exists
	_, err := uc.categoryRepo.GetByID(ctx, input.CategoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	// Generate slug if not provided
	slug := input.Slug
	if slug == "" {
		slug = generateSlug(input.Name)
	}

	// Check if slug already exists
	_, err = uc.productRepo.GetBySlug(ctx, slug)
	if err == nil {
		return nil, ErrProductSlugExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	gallery := domain.Gallery(input.Gallery)
	if gallery == nil {
		gallery = domain.Gallery{}
	}

	product := &domain.Product{
		ProductID:    utils.GenerateUUIDv7(),
		Name:         input.Name,
		Slug:         slug,
		SKU:          input.SKU,
		ShortDesc:    input.ShortDesc,
		Description:  input.Description,
		Price:        input.Price,
		Stock:        input.Stock,
		ThumbnailURL: input.ThumbnailURL,
		Gallery:      gallery,
		CategoryID:   input.CategoryID,
		IsActive:     input.IsActive,
	}

	if err := uc.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	// Reload to get category relationship
	return uc.productRepo.GetByID(ctx, product.ProductID)
}

func (uc *ProductUseCase) GetByID(ctx context.Context, productID string) (*domain.Product, error) {
	product, err := uc.productRepo.GetByID(ctx, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return product, nil
}

func (uc *ProductUseCase) GetBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	product, err := uc.productRepo.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return product, nil
}

func (uc *ProductUseCase) Update(ctx context.Context, productID string, input UpdateProductInput) (*domain.Product, error) {
	product, err := uc.productRepo.GetByID(ctx, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	// Validate category if changed
	if input.CategoryID != nil {
		_, err := uc.categoryRepo.GetByID(ctx, *input.CategoryID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrCategoryNotFound
			}
			return nil, err
		}
		product.CategoryID = *input.CategoryID
	}

	// Check slug uniqueness if changed
	if input.Slug != nil {
		existingProd, err := uc.productRepo.GetBySlug(ctx, *input.Slug)
		if err == nil && existingProd.ProductID != productID {
			return nil, ErrProductSlugExists
		}
		product.Slug = *input.Slug
	}

	// Update fields
	if input.Name != nil {
		product.Name = *input.Name
	}
	if input.SKU != nil {
		product.SKU = *input.SKU
	}
	if input.ShortDesc != nil {
		product.ShortDesc = *input.ShortDesc
	}
	if input.Description != nil {
		product.Description = *input.Description
	}
	if input.Price != nil {
		product.Price = *input.Price
	}
	if input.Stock != nil {
		product.Stock = *input.Stock
	}
	if input.ThumbnailURL != nil {
		product.ThumbnailURL = *input.ThumbnailURL
	}
	if input.Gallery != nil {
		product.Gallery = domain.Gallery(*input.Gallery)
	}
	if input.IsActive != nil {
		product.IsActive = *input.IsActive
	}

	if err := uc.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	// Reload to get updated relationships
	return uc.productRepo.GetByID(ctx, productID)
}

func (uc *ProductUseCase) Delete(ctx context.Context, productID string) error {
	// Check if product exists
	_, err := uc.productRepo.GetByID(ctx, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProductNotFound
		}
		return err
	}

	return uc.productRepo.Delete(ctx, productID)
}

func (uc *ProductUseCase) List(ctx context.Context, filter domain.ProductFilter) ([]*domain.Product, int64, error) {
	return uc.productRepo.List(ctx, filter)
}

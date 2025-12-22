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
	ErrProductNotFound    = errors.New("product not found")
	ErrProductSlugExists  = errors.New("product slug already exists")
	ErrInvalidProductData = errors.New("invalid product data")
)

type ProductUseCase struct {
	productRepo  domain.ProductRepository
	categoryRepo domain.CategoryRepository
	variantRepo  domain.ProductVariantRepository
}

func NewProductUseCase(productRepo domain.ProductRepository, categoryRepo domain.CategoryRepository, variantRepo domain.ProductVariantRepository) *ProductUseCase {
	return &ProductUseCase{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		variantRepo:  variantRepo,
	}
}

type CreateProductInput struct {
	Name         string
	Slug         string
	SKU          string
	ShortDesc    string
	Description  string
	Stock        int
	ThumbnailURL string
	Gallery      []string
	CategoryID   string
	IsActive     bool
	Variants     []CreateProductVariantInput // Required: at least 1 variant
}

type CreateProductVariantInput struct {
	SKU        string
	Name       string
	Attributes map[string]string
	Price      float64 // Required for variants
	Stock      int
	IsActive   bool
}

type UpdateProductInput struct {
	Name         *string
	Slug         *string
	SKU          *string
	ShortDesc    *string
	Description  *string
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
		// Slug already exists
		return nil, ErrProductSlugExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Unexpected error
		return nil, err
	}
	// Slug doesn't exist (ErrRecordNotFound) - OK to continue

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
		Stock:        input.Stock,
		ThumbnailURL: input.ThumbnailURL,
		Gallery:      gallery,
		CategoryID:   input.CategoryID,
		IsActive:     input.IsActive,
	}

	if err := uc.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	// Create variants (required - at least 1 variant)
	if len(input.Variants) == 0 {
		_ = uc.productRepo.Delete(ctx, product.ProductID)
		return nil, errors.New("at least one variant is required")
	}

	for _, variantInput := range input.Variants {
		// Check if SKU already exists
		existingVariant, err := uc.variantRepo.GetBySKU(ctx, variantInput.SKU)
		if err == nil && existingVariant != nil {
			// Rollback: delete product if SKU already exists
			_ = uc.productRepo.Delete(ctx, product.ProductID)
			return nil, errors.New("variant SKU already exists: " + variantInput.SKU)
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			// Rollback: delete product on unexpected error
			_ = uc.productRepo.Delete(ctx, product.ProductID)
			return nil, err
		}

		variant := &domain.ProductVariant{
			VariantID:  utils.GenerateUUIDv7(),
			ProductID:  product.ProductID,
			SKU:        variantInput.SKU,
			Name:       variantInput.Name,
			Attributes: variantInput.Attributes,
			Price:      variantInput.Price,
			Stock:      variantInput.Stock,
			IsActive:   variantInput.IsActive,
		}
		if err := uc.variantRepo.Create(ctx, variant); err != nil {
			// Rollback: delete product if variant creation fails
			_ = uc.productRepo.Delete(ctx, product.ProductID)
			// Check if it's a duplicate key error
			errMsg := err.Error()
			if strings.Contains(errMsg, "Duplicate entry") || strings.Contains(errMsg, "duplicate key") || strings.Contains(errMsg, "UNIQUE constraint") {
				return nil, errors.New("variant SKU already exists: " + variantInput.SKU)
			}
			return nil, err
		}
	}

	// Reload product without Preload to avoid relationship issues
	// Category will be loaded manually if needed in handler
	reloadedProduct, err := uc.productRepo.GetByID(ctx, product.ProductID)
	if err != nil {
		return nil, err
	}
	return reloadedProduct, nil
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
	// Note: Price is now managed via variants, not directly on product
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

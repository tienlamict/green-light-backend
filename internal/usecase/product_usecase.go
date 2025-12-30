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
	imageRepo    domain.ProductImageRepository
}

func NewProductUseCase(productRepo domain.ProductRepository, categoryRepo domain.CategoryRepository, variantRepo domain.ProductVariantRepository, imageRepo domain.ProductImageRepository) *ProductUseCase {
	return &ProductUseCase{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		variantRepo:  variantRepo,
		imageRepo:    imageRepo,
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

type CreateProductVariantImageInput struct {
	URL       string
	IsMain    bool
	SortOrder int
}

type CreateProductVariantInput struct {
	SKU        string
	Name       string
	Attributes map[string]string
	Price      float64 // Required for variants
	Stock      int
	IsActive   bool
	Images     []CreateProductVariantImageInput // Optional variant-specific images
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
	Variants     []UpdateProductVariantInput // Optional: update variants in the same request
}

type UpdateProductVariantInput struct {
	VariantID  string // Required to identify which variant to update
	SKU        *string
	Name       *string
	Attributes *map[string]string
	Price      *float64
	Stock      *int
	IsActive   *bool
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

		// Create images for this variant if provided
		if len(variantInput.Images) > 0 {
			for _, imgInput := range variantInput.Images {
				// For external URLs, use the URL as the object_key
				// In a production system, you might want to extract a path or handle this differently
				objectKey := imgInput.URL
				
				image := &domain.ProductImage{
					ImageID:   utils.GenerateUUIDv7(),
					ProductID: product.ProductID,
					VariantID: &variant.VariantID, // Set variant_id for variant-specific images
					URL:       imgInput.URL,
					ObjectKey: objectKey,
					IsMain:    imgInput.IsMain,
					SortOrder: imgInput.SortOrder,
					Status:    "ACTIVE",
				}
				
				if err := uc.imageRepo.Create(ctx, image); err != nil {
					// Log error but don't fail the entire operation
					// You might want to rollback here depending on requirements
					// For now, we'll continue with other images
					continue
				}
			}
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
	if input.Slug != nil && *input.Slug != product.Slug {
		// Only validate if slug is actually different from current slug
		existingProd, err := uc.productRepo.GetBySlug(ctx, *input.Slug)
		if err == nil && existingProd.ProductID != productID {
			return nil, ErrProductSlugExists
		}
		// Only update if slug is different
		if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
			// Unexpected error (not "not found")
			return nil, err
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

	// Update variants if provided
	if len(input.Variants) > 0 {
		for _, variantInput := range input.Variants {
			// Get the variant to update
			variant, err := uc.variantRepo.GetByID(ctx, variantInput.VariantID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errors.New("variant not found: " + variantInput.VariantID)
				}
				return nil, err
			}

			// Verify the variant belongs to this product
			if variant.ProductID != productID {
				return nil, errors.New("variant does not belong to this product: " + variantInput.VariantID)
			}

			// Check SKU uniqueness if changed
			if variantInput.SKU != nil && *variantInput.SKU != variant.SKU {
				existingVariant, err := uc.variantRepo.GetBySKU(ctx, *variantInput.SKU)
				if err == nil && existingVariant.VariantID != variantInput.VariantID {
					return nil, errors.New("variant SKU already exists: " + *variantInput.SKU)
				}
				variant.SKU = *variantInput.SKU
			}

			// Update variant fields
			if variantInput.Name != nil {
				variant.Name = *variantInput.Name
			}
			if variantInput.Attributes != nil {
				variant.Attributes = *variantInput.Attributes
			}
			if variantInput.Price != nil {
				variant.Price = *variantInput.Price
			}
			if variantInput.Stock != nil {
				variant.Stock = *variantInput.Stock
			}
			if variantInput.IsActive != nil {
				variant.IsActive = *variantInput.IsActive
			}

			// Update the variant
			if err := uc.variantRepo.Update(ctx, variant); err != nil {
				return nil, err
			}
		}
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

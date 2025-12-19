package usecase

import (
	"context"
	"errors"
	"green-light-backend/internal/domain"
	"green-light-backend/pkg/utils"

	"gorm.io/gorm"
)

var (
	ErrVariantNotFound   = errors.New("variant not found")
	ErrVariantSKUExists  = errors.New("variant SKU already exists")
	ErrInvalidVariantData = errors.New("invalid variant data")
)

type ProductVariantUseCase struct {
	variantRepo domain.ProductVariantRepository
	productRepo domain.ProductRepository
}

func NewProductVariantUseCase(
	variantRepo domain.ProductVariantRepository,
	productRepo domain.ProductRepository,
) *ProductVariantUseCase {
	return &ProductVariantUseCase{
		variantRepo: variantRepo,
		productRepo: productRepo,
	}
}

type CreateVariantInput struct {
	ProductID  string
	SKU        string
	Name       string
	Attributes map[string]string
	Price      *float64
	Stock      int
	IsActive   bool
}

type UpdateVariantInput struct {
	SKU        *string
	Name       *string
	Attributes *map[string]string
	Price      *float64
	Stock      *int
	IsActive   *bool
}

func (uc *ProductVariantUseCase) Create(ctx context.Context, input CreateVariantInput) (*domain.ProductVariant, error) {
	// Validate product exists
	_, err := uc.productRepo.GetByID(ctx, input.ProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	// Check if SKU already exists
	_, err = uc.variantRepo.GetBySKU(ctx, input.SKU)
	if err == nil {
		return nil, ErrVariantSKUExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	variant := &domain.ProductVariant{
		VariantID:  utils.GenerateUUIDv7(),
		ProductID:  input.ProductID,
		SKU:        input.SKU,
		Name:       input.Name,
		Attributes: input.Attributes,
		Price:      input.Price,
		Stock:      input.Stock,
		IsActive:   input.IsActive,
	}

	if err := uc.variantRepo.Create(ctx, variant); err != nil {
		return nil, err
	}

	return uc.variantRepo.GetByID(ctx, variant.VariantID)
}

func (uc *ProductVariantUseCase) GetByID(ctx context.Context, variantID string) (*domain.ProductVariant, error) {
	variant, err := uc.variantRepo.GetByID(ctx, variantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVariantNotFound
		}
		return nil, err
	}
	return variant, nil
}

func (uc *ProductVariantUseCase) GetBySKU(ctx context.Context, sku string) (*domain.ProductVariant, error) {
	variant, err := uc.variantRepo.GetBySKU(ctx, sku)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVariantNotFound
		}
		return nil, err
	}
	return variant, nil
}

func (uc *ProductVariantUseCase) GetByProductID(ctx context.Context, productID string) ([]*domain.ProductVariant, error) {
	return uc.variantRepo.GetByProductID(ctx, productID)
}

func (uc *ProductVariantUseCase) Update(ctx context.Context, variantID string, input UpdateVariantInput) (*domain.ProductVariant, error) {
	variant, err := uc.variantRepo.GetByID(ctx, variantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVariantNotFound
		}
		return nil, err
	}

	// Check SKU uniqueness if changed
	if input.SKU != nil && *input.SKU != variant.SKU {
		existingVariant, err := uc.variantRepo.GetBySKU(ctx, *input.SKU)
		if err == nil && existingVariant.VariantID != variantID {
			return nil, ErrVariantSKUExists
		}
		variant.SKU = *input.SKU
	}

	// Update fields
	if input.Name != nil {
		variant.Name = *input.Name
	}
	if input.Attributes != nil {
		variant.Attributes = *input.Attributes
	}
	if input.Price != nil {
		variant.Price = input.Price
	}
	if input.Stock != nil {
		variant.Stock = *input.Stock
	}
	if input.IsActive != nil {
		variant.IsActive = *input.IsActive
	}

	if err := uc.variantRepo.Update(ctx, variant); err != nil {
		return nil, err
	}

	return uc.variantRepo.GetByID(ctx, variantID)
}

func (uc *ProductVariantUseCase) Delete(ctx context.Context, variantID string) error {
	// Check if variant exists
	_, err := uc.variantRepo.GetByID(ctx, variantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrVariantNotFound
		}
		return err
	}

	return uc.variantRepo.Delete(ctx, variantID)
}


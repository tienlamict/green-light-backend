package repository

import (
	"context"
	"green-light-backend/internal/domain"

	"gorm.io/gorm"
)

type productVariantRepository struct {
	db *gorm.DB
}

func NewProductVariantRepository(db *gorm.DB) domain.ProductVariantRepository {
	return &productVariantRepository{db: db}
}

func (r *productVariantRepository) Create(ctx context.Context, variant *domain.ProductVariant) error {
	return r.db.WithContext(ctx).Create(variant).Error
}

func (r *productVariantRepository) GetByID(ctx context.Context, variantID string) (*domain.ProductVariant, error) {
	var variant domain.ProductVariant
	err := r.db.WithContext(ctx).
		Preload("Product").
		Where("variant_id = ?", variantID).
		First(&variant).Error
	if err != nil {
		return nil, err
	}
	return &variant, nil
}

func (r *productVariantRepository) GetBySKU(ctx context.Context, sku string) (*domain.ProductVariant, error) {
	var variant domain.ProductVariant
	err := r.db.WithContext(ctx).
		Preload("Product").
		Where("sku = ?", sku).
		First(&variant).Error
	if err != nil {
		return nil, err
	}
	return &variant, nil
}

func (r *productVariantRepository) GetByProductID(ctx context.Context, productID string) ([]*domain.ProductVariant, error) {
	var variants []*domain.ProductVariant
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("created_at ASC").
		Find(&variants).Error
	if err != nil {
		return nil, err
	}
	return variants, nil
}

func (r *productVariantRepository) Update(ctx context.Context, variant *domain.ProductVariant) error {
	return r.db.WithContext(ctx).Save(variant).Error
}

func (r *productVariantRepository) Delete(ctx context.Context, variantID string) error {
	return r.db.WithContext(ctx).Delete(&domain.ProductVariant{}, "variant_id = ?", variantID).Error
}

func (r *productVariantRepository) DeleteByProductID(ctx context.Context, productID string) error {
	return r.db.WithContext(ctx).Delete(&domain.ProductVariant{}, "product_id = ?", productID).Error
}


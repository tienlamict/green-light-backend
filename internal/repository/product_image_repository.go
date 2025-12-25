package repository

import (
	"context"
	"fmt"
	"time"

	"green-light-backend/internal/domain"

	"gorm.io/gorm"
)

type productImageRepository struct {
	db *gorm.DB
}

// NewProductImageRepository creates a new product image repository
func NewProductImageRepository(db *gorm.DB) domain.ProductImageRepository {
	return &productImageRepository{db: db}
}

func (r *productImageRepository) Create(ctx context.Context, image *domain.ProductImage) error {
	return r.db.WithContext(ctx).Create(image).Error
}

func (r *productImageRepository) GetByID(ctx context.Context, imageID string) (*domain.ProductImage, error) {
	var image domain.ProductImage
	err := r.db.WithContext(ctx).Where("image_id = ?", imageID).First(&image).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("image not found")
		}
		return nil, err
	}
	return &image, nil
}

// GetByProductID returns product-level images only (variant_id IS NULL)
func (r *productImageRepository) GetByProductID(ctx context.Context, productID string) ([]*domain.ProductImage, error) {
	var images []*domain.ProductImage
	err := r.db.WithContext(ctx).
		Where("product_id = ? AND variant_id IS NULL AND status = ?", productID, "ACTIVE").
		Order("is_main DESC, sort_order ASC, created_at ASC").
		Find(&images).Error
	if err != nil {
		return nil, err
	}
	return images, nil
}

// GetByVariantID returns images for a specific variant
func (r *productImageRepository) GetByVariantID(ctx context.Context, variantID string) ([]*domain.ProductImage, error) {
	var images []*domain.ProductImage
	err := r.db.WithContext(ctx).
		Where("variant_id = ? AND status = ?", variantID, "ACTIVE").
		Order("is_main DESC, sort_order ASC, created_at ASC").
		Find(&images).Error
	if err != nil {
		return nil, err
	}
	return images, nil
}

// GetAllImagesByProductID returns all images for a product (including all variants)
func (r *productImageRepository) GetAllImagesByProductID(ctx context.Context, productID string) ([]*domain.ProductImage, error) {
	var images []*domain.ProductImage
	err := r.db.WithContext(ctx).
		Where("product_id = ? AND status = ?", productID, "ACTIVE").
		Order("variant_id IS NULL DESC, is_main DESC, sort_order ASC, created_at ASC").
		Find(&images).Error
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (r *productImageRepository) Update(ctx context.Context, image *domain.ProductImage) error {
	return r.db.WithContext(ctx).Save(image).Error
}

func (r *productImageRepository) Delete(ctx context.Context, imageID string) error {
	return r.db.WithContext(ctx).Delete(&domain.ProductImage{}, "image_id = ?", imageID).Error
}

// SetMainImage sets an image as the main image for a product
// It unsets any existing main image first
func (r *productImageRepository) SetMainImage(ctx context.Context, productID, imageID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Unset all main images for this product
		if err := tx.Model(&domain.ProductImage{}).
			Where("product_id = ?", productID).
			Update("is_main", false).Error; err != nil {
			return err
		}

		// Set the new main image
		if err := tx.Model(&domain.ProductImage{}).
			Where("image_id = ?", imageID).
			Update("is_main", true).Error; err != nil {
			return err
		}

		return nil
	})
}

// DeleteOrphanImages deletes images with status=TEMP that are older than the specified duration
// This is useful for cleanup jobs to remove images that were never confirmed after upload
// TODO: Implement a cron job or scheduled task to call this periodically
func (r *productImageRepository) DeleteOrphanImages(ctx context.Context, olderThan time.Duration) error {
	cutoffTime := time.Now().Add(-olderThan)
	return r.db.WithContext(ctx).
		Where("status = ? AND created_at < ?", "TEMP", cutoffTime).
		Delete(&domain.ProductImage{}).Error
}


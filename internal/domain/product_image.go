package domain

import (
	"context"
	"time"
)

// ProductImage represents an image associated with a product or variant
// If VariantID is NULL, the image belongs to the product (shared across variants)
// If VariantID is set, the image belongs to a specific variant
type ProductImage struct {
	ImageID    string    `gorm:"primaryKey;type:varchar(36)" json:"image_id"`
	ProductID  string    `gorm:"type:varchar(36);not null;index" json:"product_id"`
	VariantID  *string   `gorm:"type:varchar(36);index" json:"variant_id,omitempty"` // NULL = product image, NOT NULL = variant image
	URL        string    `gorm:"type:varchar(500);not null" json:"url"`
	ObjectKey  string    `gorm:"type:varchar(500);not null" json:"object_key"` // MinIO object key
	IsMain     bool      `gorm:"default:false;not null" json:"is_main"`
	SortOrder  int       `gorm:"default:0;not null" json:"sort_order"`
	Status     string    `gorm:"type:varchar(20);default:'ACTIVE';not null" json:"status"` // TEMP, ACTIVE
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Product *Product        `gorm:"foreignKey:ProductID;references:ProductID" json:"product,omitempty"`
	Variant *ProductVariant `gorm:"foreignKey:VariantID;references:VariantID" json:"variant,omitempty"`
}

func (ProductImage) TableName() string {
	return "product_images"
}

// ProductImageRepository defines the interface for product image data operations
type ProductImageRepository interface {
	Create(ctx context.Context, image *ProductImage) error
	GetByID(ctx context.Context, imageID string) (*ProductImage, error)
	GetByProductID(ctx context.Context, productID string) ([]*ProductImage, error)                   // Get product-level images only (variant_id IS NULL)
	GetByVariantID(ctx context.Context, variantID string) ([]*ProductImage, error)                   // Get images for a specific variant
	GetAllImagesByProductID(ctx context.Context, productID string) ([]*ProductImage, error)          // Get all images (product + all variants)
	Update(ctx context.Context, image *ProductImage) error
	Delete(ctx context.Context, imageID string) error
	SetMainImage(ctx context.Context, productID, imageID string) error
	DeleteOrphanImages(ctx context.Context, olderThan time.Duration) error
}


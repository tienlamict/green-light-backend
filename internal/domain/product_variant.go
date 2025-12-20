package domain

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// VariantAttributes represents variant attributes (color, size, etc.) as JSON
type VariantAttributes map[string]string

// Scan implements the sql.Scanner interface
func (va *VariantAttributes) Scan(value interface{}) error {
	if value == nil {
		*va = VariantAttributes{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan VariantAttributes")
	}

	return json.Unmarshal(bytes, va)
}

// Value implements the driver.Valuer interface
func (va VariantAttributes) Value() (driver.Value, error) {
	if va == nil {
		return json.Marshal(VariantAttributes{})
	}
	return json.Marshal(va)
}

// ProductVariant represents a product variant (different SKU for same product)
type ProductVariant struct {
	VariantID  string            `gorm:"primaryKey;type:varchar(36)" json:"variant_id"`
	ProductID  string            `gorm:"type:varchar(36);not null;index" json:"product_id"`
	SKU        string            `gorm:"uniqueIndex;type:varchar(100);not null" json:"sku"`
	Name       string            `gorm:"type:varchar(255);not null" json:"name"` // e.g., "Red - Large"
	Attributes VariantAttributes `gorm:"type:json" json:"attributes"`            // {"color": "red", "size": "L"}
	Price      float64           `gorm:"type:decimal(10,2);not null" json:"price"` // Variant price (required)
	Stock      int               `gorm:"type:int;default:0;not null" json:"stock"`
	IsActive   bool              `gorm:"default:true;not null" json:"is_active"`
	CreatedAt  time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time         `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Product *Product `gorm:"foreignKey:ProductID;references:ProductID" json:"product,omitempty"`
}

func (ProductVariant) TableName() string {
	return "product_variants"
}

// ProductVariantRepository defines the interface for variant data operations
type ProductVariantRepository interface {
	Create(ctx context.Context, variant *ProductVariant) error
	GetByID(ctx context.Context, variantID string) (*ProductVariant, error)
	GetBySKU(ctx context.Context, sku string) (*ProductVariant, error)
	GetByProductID(ctx context.Context, productID string) ([]*ProductVariant, error)
	Update(ctx context.Context, variant *ProductVariant) error
	Delete(ctx context.Context, variantID string) error
}


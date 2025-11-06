package domain

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Gallery represents a JSON array of image URLs
type Gallery []string

// Scan implements the sql.Scanner interface
func (g *Gallery) Scan(value interface{}) error {
	if value == nil {
		*g = Gallery{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan Gallery")
	}

	return json.Unmarshal(bytes, g)
}

// Value implements the driver.Valuer interface
func (g Gallery) Value() (driver.Value, error) {
	if g == nil {
		return json.Marshal(Gallery{})
	}
	return json.Marshal(g)
}

// Product represents a product entity
type Product struct {
	ProductID    string    `gorm:"primaryKey;type:varchar(36)" json:"product_id"`
	Name         string    `gorm:"type:varchar(255);not null" json:"name"`
	Slug         string    `gorm:"uniqueIndex;type:varchar(255);not null" json:"slug"`
	SKU          string    `gorm:"type:varchar(100);not null" json:"sku"`
	ShortDesc    string    `gorm:"type:varchar(500)" json:"short_desc"`
	Description  string    `gorm:"type:text" json:"description"`
	Price        float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock        int       `gorm:"type:int;default:0;not null" json:"stock"`
	ThumbnailURL string    `gorm:"type:varchar(500)" json:"thumbnail_url"`
	Gallery      Gallery   `gorm:"type:json" json:"gallery"`
	CategoryID   string    `gorm:"type:varchar(36);not null;index" json:"category_id"`
	IsActive     bool      `gorm:"default:true;not null" json:"is_active"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Category *Category `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"category,omitempty"`
}

func (Product) TableName() string {
	return "products"
}

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, productID string) (*Product, error)
	GetBySlug(ctx context.Context, slug string) (*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, productID string) error
	List(ctx context.Context, filter ProductFilter) ([]*Product, int64, error)
}

// ProductFilter holds filtering options for listing products
type ProductFilter struct {
	CategoryID string
	Search     string
	MinPrice   *float64
	MaxPrice   *float64
	IsActive   *bool
	Limit      int
	Offset     int
	Sort       string
}

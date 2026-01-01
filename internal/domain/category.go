package domain

import (
	"context"
	"time"
)

// Category represents a product category
type Category struct {
	CategoryID  string    `gorm:"primaryKey;type:varchar(36)" json:"category_id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Slug        string    `gorm:"uniqueIndex;type:varchar(255);not null" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`
	IconURL     string    `gorm:"type:varchar(500)" json:"icon_url"` // Icon image URL stored in MinIO
	IsActive    bool      `gorm:"default:true;not null" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships - one category has many products
	// NOTE: Do NOT include constraint in tag - it causes GORM to create wrong FK
	Products []Product `gorm:"-" json:"products,omitempty"`
}

func (Category) TableName() string {
	return "categories"
}

// CategoryRepository defines the interface for category data operations
type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	GetByID(ctx context.Context, categoryID string) (*Category, error)
	GetBySlug(ctx context.Context, slug string) (*Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, categoryID string) error
	List(ctx context.Context, filter CategoryFilter) ([]*Category, int64, error)
}

// CategoryFilter holds filtering options for listing categories
type CategoryFilter struct {
	IsActive *bool
	Search   string
	Limit    int
	Offset   int
	Sort     string
}

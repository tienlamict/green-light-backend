package repository

import (
	"context"
	"green-light-backend/internal/domain"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) GetByID(ctx context.Context, productID string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		First(&product).Error
	if err != nil {
		return nil, err
	}
	// Load Category manually since it has gorm:"-" tag
	if product.CategoryID != "" {
		var category domain.Category
		if err := r.db.WithContext(ctx).Where("category_id = ?", product.CategoryID).First(&category).Error; err == nil {
			product.Category = &category
		}
	}
	return &product, nil
}

func (r *productRepository) GetBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.WithContext(ctx).
		Where("slug = ?", slug).
		First(&product).Error
	if err != nil {
		return nil, err
	}
	// Load Category manually since it has gorm:"-" tag
	if product.CategoryID != "" {
		var category domain.Category
		if err := r.db.WithContext(ctx).Where("category_id = ?", product.CategoryID).First(&category).Error; err == nil {
			product.Category = &category
		}
	}
	return &product, nil
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepository) Delete(ctx context.Context, productID string) error {
	return r.db.WithContext(ctx).Delete(&domain.Product{}, "product_id = ?", productID).Error
}

func (r *productRepository) List(ctx context.Context, filter domain.ProductFilter) ([]*domain.Product, int64, error) {
	var products []*domain.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Product{})

	// Apply filters
	if filter.CategoryID != "" {
		query = query.Where("category_id = ?", filter.CategoryID)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("name LIKE ? OR description LIKE ? OR sku LIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	if filter.MinPrice != nil {
		query = query.Where("price >= ?", *filter.MinPrice)
	}

	if filter.MaxPrice != nil {
		query = query.Where("price <= ?", *filter.MaxPrice)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sortOrder := "created_at DESC"
	if filter.Sort != "" {
		sortOrder = filter.Sort
	}
	query = query.Order(sortOrder)

	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	err := query.Find(&products).Error
	if err != nil {
		return nil, 0, err
	}
	
	// Load Categories manually for all products since it has gorm:"-" tag
	if len(products) > 0 {
		categoryIDs := make([]string, 0, len(products))
		categoryMap := make(map[string]*domain.Category)
		for _, p := range products {
			if p.CategoryID != "" {
				categoryIDs = append(categoryIDs, p.CategoryID)
			}
		}
		if len(categoryIDs) > 0 {
			var categories []domain.Category
			if err := r.db.WithContext(ctx).Where("category_id IN ?", categoryIDs).Find(&categories).Error; err == nil {
				for i := range categories {
					categoryMap[categories[i].CategoryID] = &categories[i]
				}
				// Assign categories to products
				for _, p := range products {
					if cat, ok := categoryMap[p.CategoryID]; ok {
						p.Category = cat
					}
				}
			}
		}
	}
	
	return products, total, nil
}

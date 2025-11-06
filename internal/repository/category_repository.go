package repository

import (
	"context"
	"green-light-backend/internal/domain"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) domain.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *domain.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *categoryRepository) GetByID(ctx context.Context, categoryID string) (*domain.Category, error) {
	var category domain.Category
	err := r.db.WithContext(ctx).Where("category_id = ?", categoryID).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	var category domain.Category
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *domain.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *categoryRepository) Delete(ctx context.Context, categoryID string) error {
	return r.db.WithContext(ctx).Delete(&domain.Category{}, "category_id = ?", categoryID).Error
}

func (r *categoryRepository) List(ctx context.Context, filter domain.CategoryFilter) ([]*domain.Category, int64, error) {
	var categories []*domain.Category
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Category{})

	// Apply filters
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", searchPattern, searchPattern)
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

	err := query.Find(&categories).Error
	return categories, total, err
}

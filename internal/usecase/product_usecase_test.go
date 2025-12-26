package usecase

import (
	"context"
	"green-light-backend/internal/domain"
	"testing"
	"time"

	"gorm.io/gorm"
)

// Mock product repository
type mockProductRepository struct {
	products map[string]*domain.Product
}

func newMockProductRepository() *mockProductRepository {
	return &mockProductRepository{
		products: make(map[string]*domain.Product),
	}
}

func (m *mockProductRepository) Create(ctx context.Context, product *domain.Product) error {
	m.products[product.ProductID] = product
	return nil
}

func (m *mockProductRepository) GetByID(ctx context.Context, productID string) (*domain.Product, error) {
	product, exists := m.products[productID]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}
	return product, nil
}

func (m *mockProductRepository) GetBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	for _, product := range m.products {
		if product.Slug == slug {
			return product, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockProductRepository) Update(ctx context.Context, product *domain.Product) error {
	m.products[product.ProductID] = product
	return nil
}

func (m *mockProductRepository) Delete(ctx context.Context, productID string) error {
	delete(m.products, productID)
	return nil
}

func (m *mockProductRepository) List(ctx context.Context, filter domain.ProductFilter) ([]*domain.Product, int64, error) {
	products := make([]*domain.Product, 0, len(m.products))
	for _, product := range m.products {
		products = append(products, product)
	}
	return products, int64(len(products)), nil
}

// Mock product variant repository
type mockProductVariantRepository struct {
	variants map[string]*domain.ProductVariant
}

func newMockProductVariantRepository() *mockProductVariantRepository {
	return &mockProductVariantRepository{
		variants: make(map[string]*domain.ProductVariant),
	}
}

func (m *mockProductVariantRepository) Create(ctx context.Context, variant *domain.ProductVariant) error {
	m.variants[variant.VariantID] = variant
	return nil
}

func (m *mockProductVariantRepository) GetByID(ctx context.Context, variantID string) (*domain.ProductVariant, error) {
	variant, exists := m.variants[variantID]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}
	return variant, nil
}

func (m *mockProductVariantRepository) GetBySKU(ctx context.Context, sku string) (*domain.ProductVariant, error) {
	for _, variant := range m.variants {
		if variant.SKU == sku {
			return variant, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockProductVariantRepository) GetByProductID(ctx context.Context, productID string) ([]*domain.ProductVariant, error) {
	variants := make([]*domain.ProductVariant, 0)
	for _, variant := range m.variants {
		if variant.ProductID == productID {
			variants = append(variants, variant)
		}
	}
	return variants, nil
}

func (m *mockProductVariantRepository) Update(ctx context.Context, variant *domain.ProductVariant) error {
	m.variants[variant.VariantID] = variant
	return nil
}

func (m *mockProductVariantRepository) Delete(ctx context.Context, variantID string) error {
	delete(m.variants, variantID)
	return nil
}

// Mock category repository
type mockCategoryRepository struct {
	categories map[string]*domain.Category
}

func newMockCategoryRepository() *mockCategoryRepository {
	return &mockCategoryRepository{
		categories: make(map[string]*domain.Category),
	}
}

func (m *mockCategoryRepository) Create(ctx context.Context, category *domain.Category) error {
	m.categories[category.CategoryID] = category
	return nil
}

func (m *mockCategoryRepository) GetByID(ctx context.Context, categoryID string) (*domain.Category, error) {
	category, exists := m.categories[categoryID]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}
	return category, nil
}

func (m *mockCategoryRepository) GetBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	for _, category := range m.categories {
		if category.Slug == slug {
			return category, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockCategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	m.categories[category.CategoryID] = category
	return nil
}

func (m *mockCategoryRepository) Delete(ctx context.Context, categoryID string) error {
	delete(m.categories, categoryID)
	return nil
}

func (m *mockCategoryRepository) List(ctx context.Context, filter domain.CategoryFilter) ([]*domain.Category, int64, error) {
	categories := make([]*domain.Category, 0, len(m.categories))
	for _, category := range m.categories {
		categories = append(categories, category)
	}
	return categories, int64(len(categories)), nil
}

// Mock product image repository
type mockProductImageRepository struct {
	images map[string]*domain.ProductImage
}

func newMockProductImageRepository() *mockProductImageRepository {
	return &mockProductImageRepository{
		images: make(map[string]*domain.ProductImage),
	}
}

func (m *mockProductImageRepository) Create(ctx context.Context, image *domain.ProductImage) error {
	m.images[image.ImageID] = image
	return nil
}

func (m *mockProductImageRepository) GetByID(ctx context.Context, imageID string) (*domain.ProductImage, error) {
	image, exists := m.images[imageID]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}
	return image, nil
}

func (m *mockProductImageRepository) GetByProductID(ctx context.Context, productID string) ([]*domain.ProductImage, error) {
	var result []*domain.ProductImage
	for _, image := range m.images {
		if image.ProductID == productID && image.VariantID == nil {
			result = append(result, image)
		}
	}
	return result, nil
}

func (m *mockProductImageRepository) GetByVariantID(ctx context.Context, variantID string) ([]*domain.ProductImage, error) {
	var result []*domain.ProductImage
	for _, image := range m.images {
		if image.VariantID != nil && *image.VariantID == variantID {
			result = append(result, image)
		}
	}
	return result, nil
}

func (m *mockProductImageRepository) GetAllImagesByProductID(ctx context.Context, productID string) ([]*domain.ProductImage, error) {
	var result []*domain.ProductImage
	for _, image := range m.images {
		if image.ProductID == productID {
			result = append(result, image)
		}
	}
	return result, nil
}

func (m *mockProductImageRepository) Update(ctx context.Context, image *domain.ProductImage) error {
	m.images[image.ImageID] = image
	return nil
}

func (m *mockProductImageRepository) Delete(ctx context.Context, imageID string) error {
	delete(m.images, imageID)
	return nil
}

func (m *mockProductImageRepository) SetMainImage(ctx context.Context, productID, imageID string) error {
	// Unset all other main images for this product
	for _, img := range m.images {
		if img.ProductID == productID && img.ImageID != imageID {
			img.IsMain = false
		}
	}
	// Set this image as main
	if img, exists := m.images[imageID]; exists {
		img.IsMain = true
	}
	return nil
}

func (m *mockProductImageRepository) DeleteOrphanImages(ctx context.Context, olderThan time.Duration) error {
	// Not needed for this test
	return nil
}

func TestProductUseCase_Create(t *testing.T) {
	// Setup
	mockProdRepo := newMockProductRepository()
	mockCatRepo := newMockCategoryRepository()
	mockVariantRepo := newMockProductVariantRepository()
	mockImageRepo := newMockProductImageRepository()
	productUC := NewProductUseCase(mockProdRepo, mockCatRepo, mockVariantRepo, mockImageRepo)

	// Create test category
	testCategory := &domain.Category{
		CategoryID: "test-category-id",
		Name:       "Test Category",
		Slug:       "test-category",
		IsActive:   true,
	}
	mockCatRepo.categories[testCategory.CategoryID] = testCategory

	tests := []struct {
		name    string
		input   CreateProductInput
		wantErr bool
		errType error
	}{
		{
			name: "successful create",
		input: CreateProductInput{
			Name:        "Test Product",
			Slug:        "test-product",
			SKU:         "TEST-001",
			Description: "Test description",
			Stock:       10,
			CategoryID:  "test-category-id",
			IsActive:    true,
			Variants: []CreateProductVariantInput{
				{
					SKU:   "TEST-001-V1",
					Name:  "Default Variant",
					Price: 99.99,
					Stock: 10,
				},
			},
		},
			wantErr: false,
		},
		{
			name: "invalid category",
		input: CreateProductInput{
			Name:        "Test Product 2",
			Slug:        "test-product-2",
			SKU:         "TEST-002",
			Description: "Test description",
			Stock:       10,
			CategoryID:  "nonexistent-category",
			IsActive:    true,
			Variants: []CreateProductVariantInput{
				{
					SKU:   "TEST-002-V1",
					Name:  "Default Variant",
					Price: 99.99,
					Stock: 10,
				},
			},
		},
			wantErr: true,
			errType: ErrCategoryNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := productUC.Create(context.Background(), tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if tt.errType != nil && err != tt.errType {
					t.Errorf("expected error %v, got %v", tt.errType, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if product == nil {
					t.Errorf("expected product, got nil")
				}
				if product != nil && product.Name != tt.input.Name {
					t.Errorf("expected name %s, got %s", tt.input.Name, product.Name)
				}
			}
		})
	}
}

func TestProductUseCase_GetByID(t *testing.T) {
	// Setup
	mockProdRepo := newMockProductRepository()
	mockCatRepo := newMockCategoryRepository()
	mockVariantRepo := newMockProductVariantRepository()
	mockImageRepo := newMockProductImageRepository()
	productUC := NewProductUseCase(mockProdRepo, mockCatRepo, mockVariantRepo, mockImageRepo)

	// Create test product
	priceMin := 99.99
	priceMax := 99.99
	testProduct := &domain.Product{
		ProductID:  "test-product-id",
		Name:       "Test Product",
		Slug:       "test-product",
		SKU:        "TEST-001",
		PriceMin:   &priceMin,
		PriceMax:   &priceMax,
		Stock:      10,
		CategoryID: "test-category-id",
		IsActive:   true,
	}
	mockProdRepo.products[testProduct.ProductID] = testProduct

	tests := []struct {
		name      string
		productID string
		wantErr   bool
	}{
		{
			name:      "existing product",
			productID: "test-product-id",
			wantErr:   false,
		},
		{
			name:      "nonexistent product",
			productID: "nonexistent-id",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := productUC.GetByID(context.Background(), tt.productID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if product == nil {
					t.Errorf("expected product, got nil")
				}
			}
		})
	}
}

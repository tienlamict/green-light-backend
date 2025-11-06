//go:build integration
// +build integration

package repository

import (
	"context"
	"green-light-backend/internal/domain"
	"green-light-backend/pkg/utils"
	"os"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dsn = "root:root@tcp(localhost:3306)/greenlight_test?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Migrate tables
	if err := db.AutoMigrate(&domain.User{}, &domain.Category{}, &domain.Product{}); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func cleanupTestDB(t *testing.T, db *gorm.DB) {
	db.Exec("DELETE FROM products")
	db.Exec("DELETE FROM categories")
	db.Exec("DELETE FROM users")
}

func TestProductRepositoryIntegration_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	categoryRepo := NewCategoryRepository(db)
	productRepo := NewProductRepository(db)

	// Create test category
	category := &domain.Category{
		CategoryID:  utils.GenerateUUIDv7(),
		Name:        "Test Category",
		Slug:        "test-category",
		Description: "Test Description",
		IsActive:    true,
	}
	if err := categoryRepo.Create(context.Background(), category); err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}

	// Create test product
	product := &domain.Product{
		ProductID:    utils.GenerateUUIDv7(),
		Name:         "Test Product",
		Slug:         "test-product",
		SKU:          "TEST-001",
		Description:  "Test Description",
		Price:        99.99,
		Stock:        10,
		ThumbnailURL: "https://example.com/image.jpg",
		Gallery:      domain.Gallery{"https://example.com/1.jpg", "https://example.com/2.jpg"},
		CategoryID:   category.CategoryID,
		IsActive:     true,
	}

	err := productRepo.Create(context.Background(), product)
	if err != nil {
		t.Errorf("Failed to create product: %v", err)
	}

	// Verify product was created
	retrieved, err := productRepo.GetByID(context.Background(), product.ProductID)
	if err != nil {
		t.Errorf("Failed to retrieve product: %v", err)
	}

	if retrieved.Name != product.Name {
		t.Errorf("Expected name %s, got %s", product.Name, retrieved.Name)
	}
}

func TestProductRepositoryIntegration_List(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	categoryRepo := NewCategoryRepository(db)
	productRepo := NewProductRepository(db)

	// Create test category
	category := &domain.Category{
		CategoryID:  utils.GenerateUUIDv7(),
		Name:        "Test Category",
		Slug:        "test-category",
		Description: "Test Description",
		IsActive:    true,
	}
	if err := categoryRepo.Create(context.Background(), category); err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}

	// Create multiple products
	products := []*domain.Product{
		{
			ProductID:  utils.GenerateUUIDv7(),
			Name:       "Product 1",
			Slug:       "product-1",
			SKU:        "P001",
			Price:      100.00,
			Stock:      10,
			CategoryID: category.CategoryID,
			IsActive:   true,
		},
		{
			ProductID:  utils.GenerateUUIDv7(),
			Name:       "Product 2",
			Slug:       "product-2",
			SKU:        "P002",
			Price:      200.00,
			Stock:      20,
			CategoryID: category.CategoryID,
			IsActive:   true,
		},
	}

	for _, product := range products {
		if err := productRepo.Create(context.Background(), product); err != nil {
			t.Fatalf("Failed to create product: %v", err)
		}
	}

	// Test list with filters
	filter := domain.ProductFilter{
		Limit:  10,
		Offset: 0,
		Sort:   "created_at DESC",
	}

	retrieved, total, err := productRepo.List(context.Background(), filter)
	if err != nil {
		t.Errorf("Failed to list products: %v", err)
	}

	if total != 2 {
		t.Errorf("Expected 2 products, got %d", total)
	}

	if len(retrieved) != 2 {
		t.Errorf("Expected 2 products in result, got %d", len(retrieved))
	}
}

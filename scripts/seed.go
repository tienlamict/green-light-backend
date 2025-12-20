package main

import (
	"errors"
	"fmt"
	"green-light-backend/internal/domain"
	"green-light-backend/pkg/config"
	"green-light-backend/pkg/database"
	"green-light-backend/pkg/utils"
	"log"
	"os"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fmt.Printf("Connecting to database: %s@%s:%s/%s\n",
		cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	// Connect to database
	db, err := database.NewConnection(cfg.Database.DSN, logger.Silent)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("✓ Database connection established")
	fmt.Println("")

	fmt.Println("Starting database seeding...")
	fmt.Println("")

	if err := seedDatabase(db); err != nil {
		log.Fatalf("❌ Failed to seed database: %v", err)
	}

	// Seed product variants (example: Product with multiple colors)
	fmt.Println("Seeding product variants...")
	if err := seedProductVariants(db); err != nil {
		fmt.Printf("❌ Failed to seed product variants: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Product variants seeded successfully")

	fmt.Println("")
	fmt.Println("✅ Database seeding completed successfully!")
	fmt.Println("")
	fmt.Println("You can now login with:")
	fmt.Println("  Email: admin@example.com")
	fmt.Println("  Password: admin123")
}

func seedDatabase(db *gorm.DB) error {
	// Seed admin user
	if err := seedUsers(db); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	// Seed categories
	categories, err := seedCategories(db)
	if err != nil {
		return fmt.Errorf("failed to seed categories: %w", err)
	}

	// Seed products
	if err := seedProducts(db, categories); err != nil {
		return fmt.Errorf("failed to seed products: %w", err)
	}

	return nil
}

func seedUsers(db *gorm.DB) error {
	fmt.Println("Seeding users...")

	// Check if admin already exists
	var existingUser domain.User
	err := db.Where("email = ?", "admin@example.com").First(&existingUser).Error
	if err == nil {
		fmt.Printf("  Admin user already exists (ID: %s), skipping...\n", existingUser.UserID)
		return nil
	}

	// If error is not "record not found", it's a real error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Printf("  Warning: Error checking for existing user: %v\n", err)
		// Continue anyway to try creating
	}

	hashedPassword, err := utils.HashPassword("admin123")
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	admin := domain.User{
		UserID:       utils.GenerateUUIDv7(),
		Email:        "admin@example.com",
		PasswordHash: hashedPassword,
		Role:         domain.RoleAdmin,
	}

	if err := db.Create(&admin).Error; err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	fmt.Printf("  ✓ Created admin user: %s (ID: %s)\n", admin.Email, admin.UserID)

	// Verify it was created
	var verifyUser domain.User
	if err := db.Where("email = ?", "admin@example.com").First(&verifyUser).Error; err != nil {
		return fmt.Errorf("failed to verify created user: %w", err)
	}
	fmt.Printf("  ✓ Verified admin user exists in database\n")

	return nil
}

func seedCategories(db *gorm.DB) (map[string]*domain.Category, error) {
	fmt.Println("Seeding categories...")

	categories := []domain.Category{
		{
			CategoryID:  utils.GenerateUUIDv7(),
			Name:        "Electronics",
			Slug:        "electronics",
			Description: "Electronic devices and gadgets",
			IsActive:    true,
		},
		{
			CategoryID:  utils.GenerateUUIDv7(),
			Name:        "Furniture",
			Slug:        "furniture",
			Description: "Home and office furniture",
			IsActive:    true,
		},
	}

	categoryMap := make(map[string]*domain.Category)

	for i := range categories {
		// Check if category already exists
		var existing domain.Category
		if err := db.Where("slug = ?", categories[i].Slug).First(&existing).Error; err == nil {
			fmt.Printf("  Category '%s' already exists, skipping...\n", categories[i].Name)
			categoryMap[categories[i].Slug] = &existing
			continue
		}

		if err := db.Create(&categories[i]).Error; err != nil {
			return nil, err
		}

		categoryMap[categories[i].Slug] = &categories[i]
		fmt.Printf("  Created category: %s\n", categories[i].Name)
	}

	return categoryMap, nil
}

func seedProducts(db *gorm.DB, categories map[string]*domain.Category) error {
	fmt.Println("Seeding products...")

	electronics := categories["electronics"]
	furniture := categories["furniture"]

	products := []domain.Product{
		{
			ProductID:    utils.GenerateUUIDv7(),
			Name:         "Wireless Bluetooth Headphones",
			Slug:         "wireless-bluetooth-headphones",
			SKU:          "ELEC-HP-001",
			ShortDesc:    "Premium wireless headphones with active noise cancellation",
			Description:  "Experience superior sound quality with these premium wireless Bluetooth headphones. Features include active noise cancellation, 30-hour battery life, and comfortable over-ear design.",
			Stock:        50,
			ThumbnailURL: "/uploads/headphones.jpg",
			Gallery:      domain.Gallery{"/uploads/headphones-1.jpg", "/uploads/headphones-2.jpg"},
			CategoryID:   electronics.CategoryID,
			IsActive:     true,
		},
		{
			ProductID:    utils.GenerateUUIDv7(),
			Name:         "4K Smart TV 55 inch",
			Slug:         "4k-smart-tv-55-inch",
			SKU:          "ELEC-TV-002",
			ShortDesc:    "Ultra HD 4K Smart TV with HDR support",
			Description:  "Immerse yourself in stunning 4K resolution with this 55-inch smart TV. Features HDR10+, built-in streaming apps, and voice control compatibility.",
			Stock:        25,
			ThumbnailURL: "/uploads/tv.jpg",
			Gallery:      domain.Gallery{"/uploads/tv-1.jpg", "/uploads/tv-2.jpg", "/uploads/tv-3.jpg"},
			CategoryID:   electronics.CategoryID,
			IsActive:     true,
		},
		{
			ProductID:    utils.GenerateUUIDv7(),
			Name:         "Smartphone 128GB",
			Slug:         "smartphone-128gb",
			SKU:          "ELEC-PHONE-003",
			ShortDesc:    "Latest flagship smartphone with 5G connectivity",
			Description:  "Stay connected with this powerful smartphone featuring 128GB storage, triple camera system, and lightning-fast 5G connectivity. Perfect for photography enthusiasts and power users.",
			Stock:        100,
			ThumbnailURL: "/uploads/phone.jpg",
			Gallery:      domain.Gallery{"/uploads/phone-1.jpg", "/uploads/phone-2.jpg"},
			CategoryID:   electronics.CategoryID,
			IsActive:     true,
		},
		{
			ProductID:    utils.GenerateUUIDv7(),
			Name:         "Modern Office Desk",
			Slug:         "modern-office-desk",
			SKU:          "FURN-DESK-001",
			ShortDesc:    "Spacious office desk with cable management",
			Description:  "Enhance your workspace with this modern office desk. Features include built-in cable management, spacious work surface, and sturdy construction. Perfect for home offices and professional environments.",
			Stock:        30,
			ThumbnailURL: "/uploads/desk.jpg",
			Gallery:      domain.Gallery{"/uploads/desk-1.jpg", "/uploads/desk-2.jpg"},
			CategoryID:   furniture.CategoryID,
			IsActive:     true,
		},
		{
			ProductID:    utils.GenerateUUIDv7(),
			Name:         "Ergonomic Office Chair",
			Slug:         "ergonomic-office-chair",
			SKU:          "FURN-CHAIR-002",
			ShortDesc:    "Comfortable ergonomic chair with lumbar support",
			Description:  "Work in comfort with this ergonomic office chair. Features adjustable lumbar support, breathable mesh back, and smooth-rolling casters. Designed for all-day comfort and productivity.",
			Stock:        45,
			ThumbnailURL: "/uploads/chair.jpg",
			Gallery:      domain.Gallery{"/uploads/chair-1.jpg", "/uploads/chair-2.jpg", "/uploads/chair-3.jpg"},
			CategoryID:   furniture.CategoryID,
			IsActive:     true,
		},
		{
			ProductID:    utils.GenerateUUIDv7(),
			Name:         "Bookshelf Cabinet",
			Slug:         "bookshelf-cabinet",
			SKU:          "FURN-SHELF-003",
			ShortDesc:    "Modern 5-tier bookshelf with storage cabinet",
			Description:  "Organize your books and decor with this stylish bookshelf. Features 5 spacious tiers and a bottom storage cabinet. Constructed from high-quality wood with a contemporary finish.",
			Stock:        20,
			ThumbnailURL: "/uploads/bookshelf.jpg",
			Gallery:      domain.Gallery{"/uploads/bookshelf-1.jpg", "/uploads/bookshelf-2.jpg"},
			CategoryID:   furniture.CategoryID,
			IsActive:     true,
		},
	}

	for i := range products {
		// Check if product already exists
		var existing domain.Product
		if err := db.Where("slug = ?", products[i].Slug).First(&existing).Error; err == nil {
			fmt.Printf("  Product '%s' already exists, skipping...\n", products[i].Name)
			continue
		}

		if err := db.Create(&products[i]).Error; err != nil {
			return err
		}

		fmt.Printf("  Created product: %s\n", products[i].Name)
	}

	return nil
}

func seedProductVariants(db *gorm.DB) error {
	// Find a product to add variants to (e.g., the first product)
	var product domain.Product
	if err := db.First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("⚠️  No products found to add variants to")
			return nil
		}
		return err
	}

	fmt.Printf("Adding variants to product: %s\n", product.Name)

	// Create variants for the first product
	variants := []domain.ProductVariant{
		{
			VariantID:  utils.GenerateUUIDv7(),
			ProductID:  product.ProductID,
			SKU:        "ELEC-HP-001-BLK",
			Name:       "Black",
			Attributes: domain.VariantAttributes{"color": "black"},
			Price:      199.99,
			Stock:      25,
			IsActive:   true,
		},
		{
			VariantID:  utils.GenerateUUIDv7(),
			ProductID:  product.ProductID,
			SKU:        "ELEC-HP-001-WHT",
			Name:       "White",
			Attributes: domain.VariantAttributes{"color": "white"},
			Price:      199.99,
			Stock:      20,
			IsActive:   true,
		},
		{
			VariantID:  utils.GenerateUUIDv7(),
			ProductID:  product.ProductID,
			SKU:        "ELEC-HP-001-RED",
			Name:       "Red",
			Attributes: domain.VariantAttributes{"color": "red"},
			Price:      219.99, // Premium color
			Stock:      15,
			IsActive:   true,
		},
	}

	for _, variant := range variants {
		var existingVariant domain.ProductVariant
		err := db.Where("sku = ?", variant.SKU).First(&existingVariant).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Variant doesn't exist, create it
				if err := db.Create(&variant).Error; err != nil {
					return fmt.Errorf("failed to create variant %s: %w", variant.SKU, err)
				}
				fmt.Printf("  ✓ Created variant: %s (%s)\n", variant.Name, variant.SKU)
			} else {
				return fmt.Errorf("failed to check variant existence: %w", err)
			}
		} else {
			fmt.Printf("  ⚠️  Variant already exists: %s\n", variant.SKU)
		}
	}

	return nil
}


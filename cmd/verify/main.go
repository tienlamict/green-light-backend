package main

import (
	"fmt"
	"green-light-backend/internal/domain"
	"green-light-backend/pkg/config"
	"green-light-backend/pkg/database"
	"log"

	"gorm.io/gorm/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.NewConnection(cfg.Database.DSN, logger.Silent)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Verifying database contents...")
	fmt.Println("")

	// Check users
	var userCount int64
	db.Model(&domain.User{}).Count(&userCount)
	fmt.Printf("Users: %d\n", userCount)
	if userCount > 0 {
		var users []domain.User
		db.Find(&users)
		for _, user := range users {
			fmt.Printf("  - %s (ID: %s, Role: %s)\n", user.Email, user.UserID, user.Role)
		}
	}

	// Check categories
	var categoryCount int64
	db.Model(&domain.Category{}).Count(&categoryCount)
	fmt.Printf("\nCategories: %d\n", categoryCount)
	if categoryCount > 0 {
		var categories []domain.Category
		db.Find(&categories)
		for _, cat := range categories {
			fmt.Printf("  - %s (ID: %s, Slug: %s)\n", cat.Name, cat.CategoryID, cat.Slug)
		}
	}

	// Check products
	var productCount int64
	db.Model(&domain.Product{}).Count(&productCount)
	fmt.Printf("\nProducts: %d\n", productCount)
	if productCount > 0 {
		var products []domain.Product
		db.Limit(5).Find(&products)
		for _, prod := range products {
			priceRange := "N/A"
		if prod.PriceMin != nil && prod.PriceMax != nil {
			if *prod.PriceMin == *prod.PriceMax {
				priceRange = fmt.Sprintf("$%.2f", *prod.PriceMin)
			} else {
				priceRange = fmt.Sprintf("$%.2f - $%.2f", *prod.PriceMin, *prod.PriceMax)
			}
		}
		fmt.Printf("  - %s (ID: %s, Price: %s)\n", prod.Name, prod.ProductID, priceRange)
		}
		if productCount > 5 {
			fmt.Printf("  ... and %d more\n", productCount-5)
		}
	}

	fmt.Println("")
	if userCount == 0 {
		fmt.Println("⚠️  WARNING: No users found! Run seed script:")
		fmt.Println("   docker-compose exec api go run /root/scripts/seed.go")
	} else {
		fmt.Println("✅ Database verification complete!")
	}
}


package main

import (
	"fmt"
	"green-light-backend/pkg/config"
	"green-light-backend/pkg/database"
	"log"
	"os"

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

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	fmt.Println("Starting database migration...")

	// Read SQL migration file
	sqlFile := "migrations/001_create_tables.sql"
	sqlBytes, err := os.ReadFile(sqlFile)
	if err != nil {
		log.Fatalf("Failed to read migration file: %v", err)
	}

	// Execute migration
	_, err = sqlDB.Exec(string(sqlBytes))
	if err != nil {
		log.Fatalf("Failed to execute migration: %v", err)
	}

	fmt.Println("✓ Tables created successfully!")
	fmt.Println("✓ Migration completed!")
}

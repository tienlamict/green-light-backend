package main

import (
	"fmt"
	"green-light-backend/internal/domain"
	httpTransport "green-light-backend/internal/transport/http"
	"green-light-backend/pkg/config"
	"green-light-backend/pkg/database"
	"green-light-backend/pkg/logger"
	"green-light-backend/pkg/utils"
	"log"

	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	if err := logger.InitLogger(cfg.Log.Level); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Log.Info("Starting Green Light Backend API")

	// Connect to database
	logLevel := gormLogger.Silent
	if cfg.Server.ENV == "development" {
		logLevel = gormLogger.Info
	}

	db, err := database.NewConnection(cfg.Database.DSN, logLevel)
	if err != nil {
		logger.Log.Fatal(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	logger.Log.Info("Connected to database successfully")

	// Auto migrate database
	if err := migrateDatabase(db); err != nil {
		logger.Log.Fatal(fmt.Sprintf("Failed to migrate database: %v", err))
	}

	logger.Log.Info("Database migration completed")

	// Setup HTTP router
	router := httpTransport.NewRouter(cfg, db)
	engine := router.Setup()

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	logger.Log.Info(fmt.Sprintf("Server starting on %s", addr))

	if err := engine.Run(addr); err != nil {
		logger.Log.Fatal(fmt.Sprintf("Failed to start server: %v", err))
	}
}

func migrateDatabase(db *gorm.DB) error {
	// NOTE: We DO NOT use GORM AutoMigrate because it creates wrong foreign keys
	// Instead, we rely on SQL migration scripts in migrations/001_create_tables.sql
	
	// Just verify tables exist
	var tableCount int64
	db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name IN ('users', 'categories', 'products')").Scan(&tableCount)
	
	if tableCount < 3 {
		return fmt.Errorf("tables not found - please run SQL migration: migrations/001_create_tables.sql")
	}
	
	return nil
}

func createDefaultAdminUser(db *gorm.DB) error {
	// Check if admin user already exists
	var existingUser domain.User
	result := db.Where("email = ?", "admin@example.com").First(&existingUser)
	if result.Error == nil {
		// Admin user already exists
		return nil
	}
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	// Hash password
	hashedPassword, err := utils.HashPassword("admin123")
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Create admin user
	admin := domain.User{
		UserID:       utils.GenerateUUIDv7(),
		Email:        "admin@example.com",
		PasswordHash: hashedPassword,
		Role:         domain.RoleAdmin,
	}

	if err := db.Create(&admin).Error; err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	logger.Log.Info("Default admin user created: admin@example.com")
	return nil
}

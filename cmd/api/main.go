package main

import (
	"fmt"
	"green-light-backend/internal/domain"
	httpTransport "green-light-backend/internal/transport/http"
	"green-light-backend/pkg/config"
	"green-light-backend/pkg/database"
	"green-light-backend/pkg/logger"
	"log"

	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// @title Green Light Backend API
// @version 1.0
// @description Product showcase backend REST API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@greenlight.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

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
	logger.Log.Info(fmt.Sprintf("Swagger documentation available at http://localhost:%s/api/docs/index.html", cfg.Server.Port))

	if err := engine.Run(addr); err != nil {
		logger.Log.Fatal(fmt.Sprintf("Failed to start server: %v", err))
	}
}

func migrateDatabase(db *gorm.DB) error {
	// Disable foreign key checks temporarily for migration
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")

	// Migrate tables in correct order to avoid foreign key issues
	// 1. Users (no dependencies)
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		db.Exec("SET FOREIGN_KEY_CHECKS = 1") // Re-enable before returning
		return fmt.Errorf("failed to migrate users table: %w", err)
	}

	// 2. Categories (no dependencies)
	if err := db.AutoMigrate(&domain.Category{}); err != nil {
		db.Exec("SET FOREIGN_KEY_CHECKS = 1") // Re-enable before returning
		return fmt.Errorf("failed to migrate categories table: %w", err)
	}

	// 3. Products (depends on Categories)
	if err := db.AutoMigrate(&domain.Product{}); err != nil {
		db.Exec("SET FOREIGN_KEY_CHECKS = 1") // Re-enable before returning
		return fmt.Errorf("failed to migrate products table: %w", err)
	}

	// Re-enable foreign key checks
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")

	return nil
}

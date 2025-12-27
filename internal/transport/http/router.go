package http

import (
	"green-light-backend/internal/repository"
	"green-light-backend/internal/transport/http/handler"
	"green-light-backend/internal/transport/http/middleware"
	"green-light-backend/internal/usecase"
	"green-light-backend/pkg/config"
	"green-light-backend/pkg/logger"
	"green-light-backend/pkg/storage"
	"green-light-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Router struct {
	engine *gin.Engine
	cfg    *config.Config
	db     *gorm.DB
}

func NewRouter(cfg *config.Config, db *gorm.DB) *Router {
	return &Router{
		engine: gin.New(),
		cfg:    cfg,
		db:     db,
	}
}

func (r *Router) Setup() *gin.Engine {
	// Set Gin mode
	if r.cfg.Server.ENV == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Global middleware
	r.engine.Use(middleware.RecoveryMiddleware(logger.Log))
	r.engine.Use(middleware.LoggerMiddleware(logger.Log))
	r.engine.Use(middleware.CORSMiddleware(r.cfg.CORS.Origins))

	// Serve static files (uploads)
	r.engine.Static("/uploads", r.cfg.Upload.Dir)

	// Initialize MinIO client
	minioClient, err := storage.NewMinIOClient(&r.cfg.MinIO)
	if err != nil {
		logger.Log.Fatal("Failed to initialize MinIO client", zap.Error(err))
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(r.db)
	categoryRepo := repository.NewCategoryRepository(r.db)
	productRepo := repository.NewProductRepository(r.db)
	variantRepo := repository.NewProductVariantRepository(r.db)
	imageRepo := repository.NewProductImageRepository(r.db)

	// Initialize use cases
	jwtManager := utils.NewJWTManager(r.cfg.JWT.Secret, r.cfg.JWT.ExpireHours)
	authUseCase := usecase.NewAuthUseCase(userRepo, jwtManager)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepo)
	productUseCase := usecase.NewProductUseCase(productRepo, categoryRepo, variantRepo, imageRepo)
	variantUseCase := usecase.NewProductVariantUseCase(variantRepo, productRepo)
	imageUseCase := usecase.NewProductImageUseCase(imageRepo, productRepo, variantRepo, minioClient)

	// Initialize handlers
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(authUseCase)
	categoryHandler := handler.NewCategoryHandler(categoryUseCase)
	productHandler := handler.NewProductHandler(productUseCase, variantUseCase, imageRepo)
	variantHandler := handler.NewProductVariantHandler(variantUseCase)
	imageHandler := handler.NewProductImageHandler(imageUseCase)
	uploadHandler := handler.NewUploadHandler(r.cfg.Upload.Dir, r.cfg.Upload.MaxFileSize)

	// Health check
	r.engine.GET("/healthz", healthHandler.HealthCheck)

	// API v1 routes
	v1 := r.engine.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.GET("/user-info", middleware.AuthMiddleware(jwtManager), authHandler.GetUserInfo)
		}

		// Public product routes (GET endpoints are public)
		products := v1.Group("/products")
		{
			products.GET("", productHandler.List)

			// Public variant routes (MUST be registered BEFORE /:id_or_slug to avoid route conflict)
			products.GET("/:id_or_slug/variants", variantHandler.ListByProduct)
			products.GET("/:id_or_slug/variants/:variant_id", variantHandler.Get)

			// Public image routes
			products.GET("/:id_or_slug/images", imageHandler.ListProductImages)

			// Product detail route (MUST be after variant routes)
			products.GET("/:id_or_slug", productHandler.Get)
		}

		// Protected product routes (POST/PUT/DELETE require authentication and role)
		productsProtected := v1.Group("/products")
		productsProtected.Use(middleware.AuthMiddleware(jwtManager))
		{
			productsProtected.POST("",
				middleware.RequireRole("admin", "editor"),
				productHandler.Create)
			productsProtected.PUT("/:id_or_slug",
				middleware.RequireRole("admin", "editor"),
				productHandler.Update)
			productsProtected.DELETE("/:id_or_slug",
				middleware.RequireRole("admin"),
				productHandler.Delete)
		}

		// Protected variant routes for POST/PUT/DELETE (MUST be registered BEFORE any /:id_or_slug route)
		productVariantsProtected := v1.Group("/products")
		productVariantsProtected.Use(middleware.AuthMiddleware(jwtManager))
		{
			productVariantsProtected.POST("/:id_or_slug/variants",
				middleware.RequireRole("admin", "editor"),
				variantHandler.Create)
			productVariantsProtected.PUT("/:id_or_slug/variants/:variant_id",
				middleware.RequireRole("admin", "editor"),
				variantHandler.Update)
			productVariantsProtected.DELETE("/:id_or_slug/variants/:variant_id",
				middleware.RequireRole("admin"),
				variantHandler.Delete)

			// Protected image routes for POST/PATCH/DELETE
			productVariantsProtected.POST("/:id_or_slug/images/presign",
				middleware.RequireRole("admin", "editor"),
				imageHandler.GeneratePresignedURL)
			productVariantsProtected.POST("/:id_or_slug/images",
				middleware.RequireRole("admin", "editor"),
				imageHandler.ConfirmImageUpload)
			productVariantsProtected.PATCH("/:id_or_slug/images/:image_id",
				middleware.RequireRole("admin", "editor"),
				imageHandler.UpdateImageMetadata)
			productVariantsProtected.DELETE("/:id_or_slug/images/:image_id",
				middleware.RequireRole("admin"),
				imageHandler.DeleteImage)
		}

		// Public variant lookup by SKU
		v1.GET("/variants/sku/:sku", variantHandler.GetBySKU)

		// Public category routes (GET endpoints are public)
		categories := v1.Group("/categories")
		{
			categories.GET("", categoryHandler.List)
			categories.GET("/:id_or_slug", categoryHandler.Get)
		}

		// Protected category routes (POST/PUT/DELETE require authentication and role)
		categoriesProtected := v1.Group("/categories")
		categoriesProtected.Use(middleware.AuthMiddleware(jwtManager))
		{
			categoriesProtected.POST("",
				middleware.RequireRole("admin", "editor"),
				categoryHandler.Create)
			categoriesProtected.PUT("/:id",
				middleware.RequireRole("admin", "editor"),
				categoryHandler.Update)
			categoriesProtected.DELETE("/:id",
				middleware.RequireRole("admin"),
				categoryHandler.Delete)
		}

		// Upload routes (protected)
		uploads := v1.Group("/uploads")
		uploads.Use(middleware.AuthMiddleware(jwtManager))
		{
			uploads.POST("", uploadHandler.Upload)
		}
	}

	return r.engine
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}

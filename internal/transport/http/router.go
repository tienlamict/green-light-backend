package http

import (
	"green-light-backend/internal/repository"
	"green-light-backend/internal/transport/http/handler"
	"green-light-backend/internal/transport/http/middleware"
	"green-light-backend/internal/usecase"
	"green-light-backend/pkg/config"
	"green-light-backend/pkg/logger"
	"green-light-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	// Initialize repositories
	userRepo := repository.NewUserRepository(r.db)
	categoryRepo := repository.NewCategoryRepository(r.db)
	productRepo := repository.NewProductRepository(r.db)

	// Initialize use cases
	jwtManager := utils.NewJWTManager(r.cfg.JWT.Secret, r.cfg.JWT.ExpireHours)
	authUseCase := usecase.NewAuthUseCase(userRepo, jwtManager)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepo)
	productUseCase := usecase.NewProductUseCase(productRepo, categoryRepo)

	// Initialize handlers
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(authUseCase)
	categoryHandler := handler.NewCategoryHandler(categoryUseCase)
	productHandler := handler.NewProductHandler(productUseCase)
	uploadHandler := handler.NewUploadHandler(r.cfg.Upload.Dir, r.cfg.Upload.MaxFileSize)

	// Health check
	r.engine.GET("/healthz", healthHandler.HealthCheck)

	// Swagger documentation
	r.engine.GET("/api/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := r.engine.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.GET("/user-info", middleware.AuthMiddleware(jwtManager), authHandler.GetUserInfo)
		}

		// Public product routes
		products := v1.Group("/products")
		{
			products.GET("", productHandler.List)
			products.GET("/:id_or_slug", productHandler.Get)
		}

		// Protected product routes
		productsProtected := v1.Group("/products")
		productsProtected.Use(middleware.AuthMiddleware(jwtManager))
		{
			productsProtected.POST("",
				middleware.RequireRole("admin", "editor"),
				productHandler.Create)
			productsProtected.PUT("/:id",
				middleware.RequireRole("admin", "editor"),
				productHandler.Update)
			productsProtected.DELETE("/:id",
				middleware.RequireRole("admin"),
				productHandler.Delete)
		}

		// Public category routes
		categories := v1.Group("/categories")
		{
			categories.GET("", categoryHandler.List)
			categories.GET("/:id_or_slug", categoryHandler.Get)
		}

		// Protected category routes
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

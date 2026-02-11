package main

import (
	"ewallet/config"
	_ "ewallet/docs"
	"ewallet/internal/handlers"
	"ewallet/internal/middleware"
	"ewallet/internal/repository"
	"ewallet/internal/service"
	"ewallet/pkg/utils"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title E-Wallet API
// @version 1.0
// @description RESTful API for simple e-wallet system built with Go and Gin framework
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@ewallet.com

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

	// Run migrations first
	if err := config.RunMigrations(&cfg.Database); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize database
	db, err := config.InitDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize JWT utility
	jwtUtil := utils.NewJWTUtil(cfg.JWT.Secret, cfg.JWT.Expiry)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, walletRepo, jwtUtil, db)
	walletService := service.NewWalletService(walletRepo, transactionRepo, db)
	transactionService := service.NewTransactionService(walletRepo, transactionRepo, userRepo, db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userRepo)
	walletHandler := handlers.NewWalletHandler(walletService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Setup Gin router
	router := gin.Default()

	// Public routes
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(jwtUtil))
		{
			users.GET("/profile", userHandler.GetProfile)
		}

		wallets := api.Group("/wallets")
		wallets.Use(middleware.AuthMiddleware(jwtUtil))
		{
			wallets.GET("/balance", walletHandler.GetBalance)
			wallets.POST("/topup", walletHandler.TopUp)
		}

		transactions := api.Group("/transactions")
		transactions.Use(middleware.AuthMiddleware(jwtUtil))
		{
			transactions.POST("/transfer", transactionHandler.Transfer)
			transactions.GET("/history", transactionHandler.GetHistory)
		}
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

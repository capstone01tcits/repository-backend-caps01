package main

import (
	"fmt"
	"log"

	"go-auth/config"
	"go-auth/internal/handler"
	"go-auth/internal/middleware"
	"go-auth/internal/model"
	"go-auth/internal/repository"
	"go-auth/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load config
	config.Load()

	// Connect database
	db := config.ConnectDB()

	// Auto migrate
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatal("Migration failed:", err)
	}

	// Init dependencies
	userRepo := repository.NewUserRepository(db)
	authSvc := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authSvc)
	aiHandler := handler.NewAIHandler()

	// Init Fiber
	app := fiber.New(fiber.Config{
		AppName: "Go Auth API v1.0",
	})

	// Global Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Routes
	api := app.Group("/api")
	auth := api.Group("/auth")

	// Public routes
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)

	// Protected routes
	auth.Get("/me", middleware.Protected(), authHandler.GetProfile)

	// AI Gateway routes (proxied to Python AI service)
	ai := api.Group("/ai")
	ai.Get("/health", aiHandler.HealthCheck)                     // Public: check AI service status
	ai.All("/*", middleware.Protected(), aiHandler.Proxy)         // Protected: proxy all other AI requests

	// Start server
	addr := fmt.Sprintf(":%s", config.Cfg.AppPort)
	fmt.Printf("✓ Server running on http://localhost%s\n", addr)
	log.Fatal(app.Listen(addr))
}

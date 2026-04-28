package main

import (
	"context"
	"fmt"
	"log"

	"Sevima-AI-Content-Creator/config"
	"Sevima-AI-Content-Creator/internal/handler"
	"Sevima-AI-Content-Creator/internal/middleware"
	"Sevima-AI-Content-Creator/internal/model"
	"Sevima-AI-Content-Creator/internal/queue"
	"Sevima-AI-Content-Creator/internal/repository"
	"Sevima-AI-Content-Creator/internal/service"

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

	// Auto migrate (10 active tables - ContentPillar & ContentTheme removed in April 2026 audit)
	if err := db.AutoMigrate(
		&model.User{},
		&model.Project{},
		&model.BusinessBrief{},
		&model.CreativeBrief{},
		&model.Storyboard{},
		&model.Scene{},
		&model.Video{},
		&model.GenerationJob{},
		&model.VideoVariant{},
		&model.SceneGeneration{},
	); err != nil {
		log.Fatal("Migration failed:", err)
	}

	// Init repositories
	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	briefRepo := repository.NewBriefRepository(db)
	contentRepo := repository.NewContentRepository(db)
	storyboardRepo := repository.NewStoryboardRepository(db)
	jobRepo := repository.NewGenerationJobRepository(db)
	variantRepo := repository.NewVideoVariantRepository(db)
	sceneRepo := repository.NewSceneGenerationRepository(db)

	// Init services
	authSvc := service.NewAuthService(userRepo)
	projectSvc := service.NewProjectService(projectRepo)
	briefSvc := service.NewBriefService(briefRepo, projectRepo, storyboardRepo)
	storyboardSvc := service.NewStoryboardService(storyboardRepo, projectRepo, contentRepo)
	creditSvc := service.NewCreditService(userRepo)
	videoGenSvc := service.NewVideoGenerationService(jobRepo, variantRepo, sceneRepo, creditSvc)

	// Init handlers
	authHandler := handler.NewAuthHandler(authSvc)
	projectHandler := handler.NewProjectHandler(projectSvc)
	briefHandler := handler.NewBriefHandler(briefSvc)
	storyboardHandler := handler.NewStoryboardHandler(storyboardSvc)
	videoHandler := handler.NewVideoHandler(videoGenSvc)
	creditHandler := handler.NewCreditHandler(creditSvc)
	aiHandler := handler.NewAIHandler()

	// Init Fiber
	app := fiber.New(fiber.Config{
		AppName: "AI Video Gen API v1.0",
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

	// ==================== Auth Routes ====================
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)
	auth.Get("/me", middleware.Protected(), authHandler.GetProfile)
	auth.Post("/change-password", middleware.Protected(), authHandler.ChangePassword)
	auth.Delete("/account", middleware.Protected(), authHandler.DeleteAccount)

	// ==================== Project Routes ====================
	projects := api.Group("/projects", middleware.Protected())
	// Initialize project from FE wizard (atomically creates project + briefs)
	projects.Post("/initialize", briefHandler.CreateProjectFromFE)
	// List and get projects
	projects.Get("/", projectHandler.GetProjects)
	projects.Get("/:id", projectHandler.GetProject)

	// ==================== Storyboard Routes ====================
	storyboard := api.Group("/storyboard", middleware.Protected())
	storyboard.Post("/generate", storyboardHandler.GenerateStoryboard)

	// ==================== Video Routes ====================
	videos := api.Group("/videos", middleware.Protected())
	videos.Post("/generate", videoHandler.GenerateVideo)
	videos.Get("/:id", videoHandler.GetVideo)
	videos.Get("/", videoHandler.ListVideos)
	videos.Get("/download/:id", videoHandler.DownloadVideo)
	videos.Post("/:variantId/regenerate", videoHandler.RegenerateVideoVariant)

	// Scene regeneration route
	api.Post("/videos/scene/:sceneId/regenerate", middleware.Protected(), videoHandler.RegenerateScene)

	// ==================== Credit Routes ====================
	credits := api.Group("/credits", middleware.Protected())
	credits.Get("/", creditHandler.GetMyCredits)

	// ==================== Admin Routes ====================
	admin := api.Group("/admin", middleware.Protected(), middleware.RequireRole("admin"))
	admin.Post("/credits", creditHandler.AddCredits)

	// ==================== AI Gateway Routes ====================
	ai := api.Group("/ai")
	ai.Get("/health", aiHandler.HealthCheck)
	ai.All("/*", middleware.Protected(), aiHandler.Proxy)

	// ==================== Job Queue ====================
	// Initialize and start job queue for async video generation
	jobQueue := queue.NewJobQueue(jobRepo, videoGenSvc)
	go func() {
		if err := jobQueue.Start(context.Background(), 3); err != nil {
			log.Printf("Job queue start error: %v", err)
		}
	}()

	// Start server
	addr := fmt.Sprintf(":%s", config.Cfg.AppPort)
	fmt.Printf("✓ Server running on http://localhost%s\n", addr)
	log.Fatal(app.Listen(addr))
}

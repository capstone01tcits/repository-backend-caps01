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
	if err := db.AutoMigrate(
		&model.User{},
		&model.Project{},
		&model.BusinessBrief{},
		&model.CreativeBrief{},
		&model.ContentPillar{},
		&model.ContentTheme{},
		&model.Storyboard{},
		&model.Scene{},
		&model.Video{},
	); err != nil {
		log.Fatal("Migration failed:", err)
	}

	// Init repositories
	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	briefRepo := repository.NewBriefRepository(db)
	contentRepo := repository.NewContentRepository(db)
	storyboardRepo := repository.NewStoryboardRepository(db)
	videoRepo := repository.NewVideoRepository(db)

	// Init services
	authSvc := service.NewAuthService(userRepo)
	projectSvc := service.NewProjectService(projectRepo)
	briefSvc := service.NewBriefService(briefRepo)
	contentSvc := service.NewContentService(contentRepo, projectRepo, userRepo)
	storyboardSvc := service.NewStoryboardService(storyboardRepo, projectRepo, contentRepo, userRepo)
	videoSvc := service.NewVideoService(videoRepo, storyboardRepo, projectRepo, userRepo)
	creditSvc := service.NewCreditService(userRepo)

	// Init handlers
	authHandler := handler.NewAuthHandler(authSvc)
	projectHandler := handler.NewProjectHandler(projectSvc)
	briefHandler := handler.NewBriefHandler(briefSvc)
	contentHandler := handler.NewContentHandler(contentSvc)
	storyboardHandler := handler.NewStoryboardHandler(storyboardSvc)
	videoHandler := handler.NewVideoHandler(videoSvc)
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
	auth.Post("/restore", authHandler.RestoreAccount)
	auth.Get("/me", middleware.Protected(), authHandler.GetProfile)
	auth.Get("/users/:user_id", middleware.Protected(), authHandler.GetUserProfile)
	auth.Post("/change-password", middleware.Protected(), authHandler.ChangePassword)
	auth.Delete("/account", middleware.Protected(), authHandler.DeleteAccount)

	// ==================== Project Routes (Dashboard) ====================
	projects := api.Group("/projects", middleware.Protected())
	projects.Post("/", projectHandler.CreateProject)
	projects.Get("/", projectHandler.GetProjects)
	projects.Get("/:id", projectHandler.GetProject)
	projects.Put("/:id", projectHandler.UpdateProject)
	projects.Delete("/:id", projectHandler.DeleteProject)

	// ==================== Brief Routes ====================
	briefs := api.Group("/briefs", middleware.Protected())
	briefs.Post("/business", briefHandler.CreateBusinessBrief)
	briefs.Get("/business", briefHandler.GetBusinessBriefs)
	briefs.Get("/business/:id", briefHandler.GetBusinessBrief)
	briefs.Put("/business/:id", briefHandler.UpdateBusinessBrief)
	briefs.Delete("/business/:id", briefHandler.DeleteBusinessBrief)
	briefs.Get("/business/:id/creative", briefHandler.GetCreativeBriefsByBusinessBrief)
	briefs.Post("/creative", briefHandler.CreateCreativeBrief)
	briefs.Get("/creative", briefHandler.GetCreativeBriefs)
	briefs.Get("/creative/:id", briefHandler.GetCreativeBrief)
	briefs.Put("/creative/:id", briefHandler.UpdateCreativeBrief)
	briefs.Delete("/creative/:id", briefHandler.DeleteCreativeBrief)

	// ==================== Content Pillar Routes ====================
	projects.Post("/:id/content-pillars/generate", contentHandler.GenerateContentPillars)
	projects.Get("/:id/content-pillars", contentHandler.GetContentPillars)
	contentPillars := api.Group("/content-pillars", middleware.Protected())
	contentPillars.Get("/:id", contentHandler.GetContentPillar)
	contentPillars.Post("/:id/select", contentHandler.SelectContentPillar)
	contentPillars.Get("/:id/themes", contentHandler.GetContentThemes)

	// ==================== Content Theme Routes ====================
	contentThemes := api.Group("/content-themes", middleware.Protected())
	contentThemes.Post("/:id/select", contentHandler.SelectContentTheme)

	// ==================== Storyboard Routes ====================
	projects.Post("/:id/storyboards/generate", storyboardHandler.GenerateStoryboards)
	projects.Get("/:id/storyboards", storyboardHandler.GetStoryboards)
	storyboards := api.Group("/storyboards", middleware.Protected())
	storyboards.Get("/:id", storyboardHandler.GetStoryboard)
	storyboards.Post("/:id/select", storyboardHandler.SelectStoryboard)
	storyboards.Get("/:id/scenes", storyboardHandler.GetScenes)

	// ==================== Video Routes ====================
	videos := api.Group("/videos", middleware.Protected())
	videos.Post("/generate", videoHandler.GenerateVideo)
	videos.Get("/", videoHandler.GetMyVideos)
	videos.Get("/:id", videoHandler.GetVideo)
	videos.Get("/:id/download", videoHandler.DownloadVideo)
	projects.Get("/:id/videos", videoHandler.GetVideosByProject)

	// ==================== Credit Routes ====================
	credits := api.Group("/credits", middleware.Protected())
	credits.Get("/", creditHandler.GetMyCredits)

	// ==================== Admin Routes ====================
	admin := api.Group("/admin", middleware.Protected())
	admin.Post("/credits", creditHandler.AddCredits)

	// ==================== AI Gateway Routes ====================
	ai := api.Group("/ai")
	ai.Get("/health", aiHandler.HealthCheck)
	ai.All("/*", middleware.Protected(), aiHandler.Proxy)

	// Start server
	addr := fmt.Sprintf(":%s", config.Cfg.AppPort)
	fmt.Printf("✓ Server running on http://localhost%s\n", addr)
	log.Fatal(app.Listen(addr))
}

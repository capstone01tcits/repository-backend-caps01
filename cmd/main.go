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

	// Auto migrate
	fmt.Println("Running database migrations...")
	if err := db.AutoMigrate(
		&model.User{},
		&model.Project{},
		&model.BusinessBrief{},
		&model.CreativeBrief{},
		&model.Storyboard{},
		&model.StoryboardSection{},
		&model.Video{},
		&model.GenerationJob{},
		&model.VideoVariant{},
		&model.SceneGeneration{},
	); err != nil {
		log.Fatal("Migration failed:", err)
	}
	fmt.Println("Database migration completed successfully!")

	// Init repositories
	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	briefRepo := repository.NewBriefRepository(db)
	storyboardRepo := repository.NewStoryboardRepository(db)
	jobRepo := repository.NewGenerationJobRepository(db)
	variantRepo := repository.NewVideoVariantRepository(db)
	sceneRepo := repository.NewSceneGenerationRepository(db)

	// Init services
	authSvc := service.NewAuthService(userRepo)
	projectSvc := service.NewProjectService(projectRepo)
	storageSvc := service.NewStorageService()
	storyboardSvc := service.NewStoryboardService(storyboardRepo, projectRepo, briefRepo)
	briefSvc := service.NewBriefService(briefRepo, projectRepo, storyboardSvc, storageSvc)
	creditSvc := service.NewCreditService(userRepo)
	videoGenSvc := service.NewVideoGenerationService(jobRepo, variantRepo, sceneRepo, briefRepo, storyboardRepo, creditSvc, storageSvc)

	// Init handlers
	authHandler := handler.NewAuthHandler(authSvc)
	projectHandler := handler.NewProjectHandler(projectSvc)
	briefHandler := handler.NewBriefHandler(briefSvc)
	storyboardHandler := handler.NewStoryboardHandler(storyboardSvc, videoGenSvc)
	videoHandler := handler.NewVideoHandler(videoGenSvc, storageSvc)
	creditHandler := handler.NewCreditHandler(creditSvc)

	// Init Fiber
	app := fiber.New(fiber.Config{
		AppName: "Sevima AI Video Gen API v1.0",
	})

	// Global Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	// ══════════════════════════════════════════════════════════════════════════
	// HEALTH
	// GET /health
	// ══════════════════════════════════════════════════════════════════════════
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "sevima-ai-video-gen"})
	})

	api := app.Group("/api")

	// ══════════════════════════════════════════════════════════════════════════
	// AUTH
	// POST   /api/auth/register
	// POST   /api/auth/login
	// POST   /api/auth/refresh
	// GET    /api/auth/me              [protected]
	// POST   /api/auth/change-password [protected]
	// DELETE /api/auth/account         [protected]
	// ══════════════════════════════════════════════════════════════════════════
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)
	auth.Get("/me", middleware.Protected(), authHandler.GetProfile)
	auth.Post("/change-password", middleware.Protected(), authHandler.ChangePassword)
	auth.Delete("/account", middleware.Protected(), authHandler.DeleteAccount)

	// ══════════════════════════════════════════════════════════════════════════
	// PROJECTS
	// POST /api/projects/initialize    — wizard: project + briefs + storyboard
	// GET  /api/projects               — list user projects
	// GET  /api/projects/:id           — get single project
	// DELETE /api/projects/:id         — soft delete project
	// POST /api/projects/:id/restore   — restore deleted project
	// ══════════════════════════════════════════════════════════════════════════
	projects := api.Group("/projects", middleware.Protected())
	projects.Post("/initialize", briefHandler.CreateProjectFromFE)
	projects.Get("/", projectHandler.GetProjects)
	projects.Get("/:id", projectHandler.GetProject)
	projects.Delete("/:id", projectHandler.DeleteProject)
	projects.Post("/:id/restore", projectHandler.RestoreProject)

	// ══════════════════════════════════════════════════════════════════════════
	// STORYBOARD
	// POST /api/storyboard/create
	// GET  /api/storyboard/detail/:storyboard_id
	// GET  /api/storyboard/:project_id
	// GET  /api/storyboard/:storyboard_id/sections
	// PUT  /api/storyboard/:storyboard_id
	// DELETE /api/storyboard/:storyboard_id
	// POST /api/storyboard/:storyboard_id/restore
	// ══════════════════════════════════════════════════════════════════════════
	storyboard := api.Group("/storyboard", middleware.Protected())
	storyboard.Post("/create", storyboardHandler.CreateManualStoryboard)
	storyboard.Get("/detail/:storyboard_id", storyboardHandler.GetStoryboard)
	storyboard.Get("/:storyboard_id/sections", storyboardHandler.GetStoryboardSections)
	storyboard.Get("/:project_id", storyboardHandler.GetStoryboardByProject)
	storyboard.Put("/:storyboard_id", storyboardHandler.UpdateStoryboard)
	storyboard.Delete("/:storyboard_id", storyboardHandler.DeleteStoryboard)
	storyboard.Post("/:storyboard_id/restore", storyboardHandler.RestoreStoryboard)

	// ══════════════════════════════════════════════════════════════════════════
	// VIDEO — AI Generation Pipeline
	// POST /api/videos/generate
	// GET  /api/videos/storyboard/:storyboard_id   — find variant IDs after generation
	// GET  /api/videos/download/:id                — get download URL
	// POST /api/videos/:variantId/regenerate        — regenerate a variant
	// POST /api/videos/scene/:sceneId/regenerate   — regenerate a single scene
	// GET  /api/videos/:id                         — get variant status (LAST, wildcard)
	// ══════════════════════════════════════════════════════════════════════════
	videos := api.Group("/videos", middleware.Protected())
	videos.Post("/generate", videoHandler.GenerateVideo)
	videos.Get("/", videoHandler.ListVideos)
	videos.Get("/storyboard/:storyboard_id", videoHandler.GetVideosByStoryboard)
	videos.Get("/download/:id", videoHandler.DownloadVideo)
	videos.Get("/preview/:id", videoHandler.PreviewVideo)
	videos.Post("/scene/:sceneId/regenerate", videoHandler.RegenerateScene)
	videos.Post("/:variantId/regenerate", videoHandler.RegenerateVideoVariant)
	videos.Get("/:id", videoHandler.GetVideo) // wildcard — must be last

	// ══════════════════════════════════════════════════════════════════════════
	// CREDITS
	// GET  /api/credits          — get my balance
	// POST /api/admin/credits    — admin: top-up user credits
	// ══════════════════════════════════════════════════════════════════════════
	credits := api.Group("/credits", middleware.Protected())
	credits.Get("/", creditHandler.GetMyCredits)

	admin := api.Group("/admin", middleware.Protected(), middleware.RequireRole("admin"))
	admin.Post("/credits", creditHandler.AddCredits)

	// ══════════════════════════════════════════════════════════════════════════
	// JOB QUEUE — async video generation worker
	// ══════════════════════════════════════════════════════════════════════════
	jobQueue := queue.NewJobQueue(jobRepo, videoGenSvc)
	go func() {
		if err := jobQueue.Start(context.Background(), 3); err != nil {
			log.Printf("Job queue start error: %v", err)
		}
	}()

	// Start server
	addr := fmt.Sprintf(":%s", config.Cfg.AppPort)
	fmt.Printf("✓ Sevima AI Video Gen API running on http://localhost%s\n", addr)
	log.Fatal(app.Listen(addr))
}

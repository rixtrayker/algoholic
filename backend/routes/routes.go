package routes

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"gorm.io/gorm"

	"github.com/yourusername/algoholic/config"
	"github.com/yourusername/algoholic/handlers"
	"github.com/yourusername/algoholic/middleware"
	"github.com/yourusername/algoholic/services"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	// Initialize services
	authService := services.NewAuthService(db, cfg)
	problemService := services.NewProblemService(db)
	userService := services.NewUserService(db)
	spacedRepService := services.NewSpacedRepetitionService(db)
	questionService := services.NewQuestionService(db, userService, spacedRepService)
	trainingPlanService := services.NewTrainingPlanService(db, questionService, userService)

	// Phase 2: Intelligence services
	embedder := services.NewEmbeddingService(cfg.Ollama.URL, cfg.Ollama.EmbeddingModel)
	vectorService := services.NewVectorService(cfg.ChromaDB.URL, embedder)
	graphService := services.NewGraphService(db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	problemHandler := handlers.NewProblemHandler(problemService)
	questionHandler := handlers.NewQuestionHandler(questionService, userService)
	userHandler := handlers.NewUserHandler(userService, questionService, spacedRepService)
	trainingPlanHandler := handlers.NewTrainingPlanHandler(trainingPlanService)
	listHandler := handlers.NewListHandler(db)
	activityHandler := handlers.NewActivityHandler(db)
	searchHandler := handlers.NewSearchHandler(db, vectorService, graphService)
	topicHandler := handlers.NewTopicHandler(db)

	// Rate limiting helpers
	rateLimitReached := func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": "Too many requests. Please try again later.",
		})
	}

	// Global rate limit: 100 requests/minute per IP
	api := app.Group("/api", limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: rateLimitReached,
	}))

	// Health check with dependency status
	app.Get("/health", func(c *fiber.Ctx) error {
		health := fiber.Map{
			"status":      "healthy",
			"app":         cfg.App.Name,
			"version":     cfg.App.Version,
			"environment": cfg.App.Environment,
			"timestamp":   time.Now().UTC(),
		}

		// Check PostgreSQL
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			health["status"] = "degraded"
			health["database"] = "unreachable"
		} else {
			health["database"] = "ok"
		}

		statusCode := fiber.StatusOK
		if health["status"] == "degraded" {
			statusCode = fiber.StatusServiceUnavailable
		}

		return c.Status(statusCode).JSON(health)
	})

	// Auth routes (public, stricter rate limit: 10/minute per IP)
	auth := api.Group("/auth", limiter.New(limiter.Config{
		Max:          10,
		Expiration:   1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string { return c.IP() },
		LimitReached: rateLimitReached,
	}))
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/forgot-password", authHandler.ForgotPassword)
	auth.Post("/reset-password", authHandler.ResetPassword)
	auth.Post("/refresh", authHandler.RefreshToken)

	// Protected routes
	protected := api.Group("")
	if cfg.Auth.Enabled {
		protected.Use(middleware.AuthMiddleware(authService))
	}

	// Auth routes (protected)
	protected.Get("/auth/me", authHandler.GetMe)
	protected.Post("/auth/change-password", authHandler.ChangePassword)

	// Problem routes
	problems := api.Group("/problems")
	problems.Get("/", problemHandler.GetProblems)
	problems.Get("/search", problemHandler.SearchProblems)
	problems.Get("/:id", problemHandler.GetProblem)
	problems.Get("/slug/:slug", problemHandler.GetProblemBySlug)
	problems.Get("/:id/topics", problemHandler.GetProblemTopics)
	problems.Get("/:id/similar", searchHandler.FindSimilarProblems)

	// Question routes
	questions := api.Group("/questions")
	questions.Get("/", questionHandler.GetQuestions)
	questions.Get("/random", questionHandler.GetRandomQuestion)
	questions.Get("/:id", questionHandler.GetQuestion)
	questions.Get("/:id/hint", questionHandler.GetHint)
	protected.Post("/questions/:id/answer", limiter.New(limiter.Config{
		Max:        20,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			if uid, ok := c.Locals("user_id").(int); ok {
				return fmt.Sprintf("answer:%d", uid)
			}
			return c.IP()
		},
		LimitReached: rateLimitReached,
	}), questionHandler.SubmitAnswer)
	protected.Get("/questions/:id/attempts", questionHandler.GetUserAttempts)

	// Topic routes (public)
	topics := api.Group("/topics")
	topics.Get("/", topicHandler.GetAllTopics)
	topics.Get("/:id", topicHandler.GetTopicByID)
	topics.Get("/:id/prerequisites", searchHandler.GetTopicPrerequisites)

	// Topic performance (protected — uses authenticated user)
	protected.Get("/topics/:topicId/performance", topicHandler.GetTopicPerformance)

	// Problem questions
	api.Get("/problems/:problemId/questions", questionHandler.GetQuestionsByProblem)

	// User routes (all protected)
	users := protected.Group("/users")
	users.Get("/me/stats", userHandler.GetUserStats)
	users.Get("/me/weaknesses", userHandler.GetWeaknesses)
	users.Get("/me/recommendations", userHandler.GetRecommendations)
	users.Get("/me/review-queue", userHandler.GetReviewQueue)
	users.Get("/me/skills", userHandler.GetUserSkills)
	users.Get("/me/skills/:topicId", userHandler.GetUserProgress)
	users.Get("/me/preferences", userHandler.GetPreferences)
	users.Put("/me/preferences", userHandler.UpdatePreferences)
	users.Get("/me/attempts", userHandler.GetRecentAttempts)
	users.Get("/me/due-reviews", userHandler.GetDueReviews)

	// Training plan routes (all protected)
	plans := protected.Group("/training-plans")
	plans.Post("/", trainingPlanHandler.CreateTrainingPlan)
	plans.Get("/", trainingPlanHandler.GetUserPlans)
	plans.Get("/:id", trainingPlanHandler.GetTrainingPlan)
	plans.Get("/:id/next", trainingPlanHandler.GetNextQuestion)
	plans.Get("/:id/items", trainingPlanHandler.GetPlanItems)
	plans.Get("/:id/today", trainingPlanHandler.GetTodaysQuestions)
	plans.Post("/:id/items/:itemId/complete", trainingPlanHandler.CompleteItem)
	plans.Post("/:id/pause", trainingPlanHandler.PausePlan)
	plans.Post("/:id/resume", trainingPlanHandler.ResumePlan)
	plans.Delete("/:id", trainingPlanHandler.DeletePlan)

	// User Lists routes (all protected)
	lists := protected.Group("/lists")
	lists.Get("/", listHandler.GetUserLists)
	lists.Post("/", listHandler.CreateList)
	lists.Get("/:id", listHandler.GetList)
	lists.Put("/:id", listHandler.UpdateList)
	lists.Delete("/:id", listHandler.DeleteList)
	lists.Post("/:id/problems", listHandler.AddProblemToList)
	lists.Delete("/:id/problems/:problemId", listHandler.RemoveProblemFromList)
	lists.Get("/:id/problems", listHandler.GetListProblems)

	// Activity routes (all protected)
	activity := protected.Group("/activity")
	activity.Get("/chart", activityHandler.GetActivityChart)
	activity.Get("/stats", activityHandler.GetActivityStats)
	activity.Get("/history", activityHandler.GetPracticeHistory)
	activity.Post("/record", activityHandler.RecordActivity)

	// Phase 2: Semantic search routes (public)
	search := api.Group("/search")
	search.Get("/problems", searchHandler.SemanticSearchProblems)
	search.Get("/questions", searchHandler.SemanticSearchQuestions)

	// Phase 2: Graph routes (public)
	graph := api.Group("/graph")
	graph.Get("/learning-path", searchHandler.GetLearningPath)

	// Phase 2: Intelligence status (public health-check)
	api.Get("/intelligence/status", searchHandler.IntelligenceStatus)

	// Development-only routes
	if cfg.IsDevelopment() {
		api.Get("/config", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"app":      cfg.App,
				"server":   cfg.Server,
				"database": fiber.Map{"host": cfg.Database.Host, "database": cfg.Database.Database},
				"chromadb": fiber.Map{"url": cfg.ChromaDB.URL},
				"ollama":   fiber.Map{"url": cfg.Ollama.URL},
				"rag":      cfg.RAG,
			})
		})

		// Phase 2: Admin endpoints (development only)
		admin := api.Group("/admin")
		admin.Post("/index", searchHandler.IndexVectors)
		admin.Post("/seed-graph", searchHandler.SeedGraph)
	}
}

package routes

import (
	"github.com/gofiber/fiber/v2"
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
	questionService := services.NewQuestionService(db)
	userService := services.NewUserService(db)
	trainingPlanService := services.NewTrainingPlanService(db, questionService, userService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	problemHandler := handlers.NewProblemHandler(problemService)
	questionHandler := handlers.NewQuestionHandler(questionService, userService)
	userHandler := handlers.NewUserHandler(userService, questionService)
	trainingPlanHandler := handlers.NewTrainingPlanHandler(trainingPlanService)

	// Public routes
	api := app.Group("/api")

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":      "healthy",
			"app":         cfg.App.Name,
			"version":     cfg.App.Version,
			"environment": cfg.App.Environment,
		})
	})

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

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

	// Question routes
	questions := api.Group("/questions")
	questions.Get("/", questionHandler.GetQuestions)
	questions.Get("/random", questionHandler.GetRandomQuestion)
	questions.Get("/:id", questionHandler.GetQuestion)
	protected.Post("/questions/:id/answer", questionHandler.SubmitAnswer)
	protected.Get("/questions/:id/attempts", questionHandler.GetUserAttempts)

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
	}
}

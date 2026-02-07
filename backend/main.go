package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/yourusername/algoholic/backend/config"
	"github.com/yourusername/algoholic/backend/models"
	"github.com/yourusername/algoholic/backend/routes"
)

var (
	DB  *gorm.DB
	cfg *config.Config
)

func initDB() error {
	// Get database configuration
	dbCfg := cfg.Database

	// Set GORM log level based on config
	var gormLogLevel gormlogger.LogLevel
	switch dbCfg.LogLevel {
	case "silent":
		gormLogLevel = gormlogger.Silent
	case "error":
		gormLogLevel = gormlogger.Error
	case "warn":
		gormLogLevel = gormlogger.Warn
	case "info":
		gormLogLevel = gormlogger.Info
	default:
		gormLogLevel = gormlogger.Warn
	}

	// Connect to database
	var err error
	DB, err = gorm.Open(postgres.Open(dbCfg.GetDSN()), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormLogLevel),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL database
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(dbCfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(dbCfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(dbCfg.ConnMaxLifetime) * time.Minute)

	// Auto-migrate (only in development)
	if cfg.IsDevelopment() && dbCfg.AutoMigrate {
		log.Println("Running auto-migration (development mode)")
		if err := models.AutoMigrate(DB); err != nil {
			return fmt.Errorf("failed to auto-migrate: %w", err)
		}
	}

	log.Println("Database connected successfully")
	return nil
}

func main() {
	// Load configuration
	// Priority: env vars > config.yaml > defaults
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	var err error
	cfg, err = config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting %s v%s (%s mode)", cfg.App.Name, cfg.App.Version, cfg.App.Environment)

	// Initialize database
	if err := initDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:           cfg.App.Name,
		EnablePrintRoutes: cfg.IsDevelopment(),
		ReadTimeout:       time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout:      time.Duration(cfg.Server.WriteTimeout) * time.Second,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			// Log error in development
			if cfg.IsDevelopment() {
				log.Printf("Error: %v", err)
			}

			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())

	if cfg.IsDevelopment() {
		app.Use(logger.New(logger.Config{
			Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
		}))
	}

	// CORS configuration from config
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Server.CORS.AllowOrigins[0], // Fiber expects string, not []string
		AllowMethods:     cfg.Server.CORS.AllowMethods[0],
		AllowHeaders:     cfg.Server.CORS.AllowHeaders[0],
		AllowCredentials: cfg.Server.CORS.AllowCredentials,
		MaxAge:           cfg.Server.CORS.MaxAge,
	}))

	// Setup all routes
	routes.SetupRoutes(app, DB, cfg)

	// Start server in goroutine
	go func() {
		addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		log.Printf("Server listening on %s", addr)
		log.Printf("API Documentation: http://localhost:%d/api", cfg.Server.Port)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Close database connections
	if sqlDB, err := DB.DB(); err == nil {
		sqlDB.Close()
	}

	// Shutdown server with timeout
	if err := app.ShutdownWithTimeout(time.Duration(cfg.Server.ShutdownTimeout) * time.Second); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited successfully")
}

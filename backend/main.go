package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/yourusername/algoholic/config"
	appmiddleware "github.com/yourusername/algoholic/middleware"
	"github.com/yourusername/algoholic/models"
	"github.com/yourusername/algoholic/routes"
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
		slog.Info("running auto-migration (development mode)")
		if err := models.AutoMigrate(DB); err != nil {
			return fmt.Errorf("failed to auto-migrate: %w", err)
		}
	}

	slog.Info("database connected successfully")
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
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Configure structured logging
	logLevel := slog.LevelInfo
	if cfg.IsDevelopment() {
		logLevel = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})))

	slog.Info("starting server",
		slog.String("app", cfg.App.Name),
		slog.String("version", cfg.App.Version),
		slog.String("environment", cfg.App.Environment),
	)

	// Initialize database
	if err := initDB(); err != nil {
		slog.Error("failed to initialize database", slog.String("error", err.Error()))
		os.Exit(1)
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

			requestID := appmiddleware.GetRequestID(c)
			slog.Error("request error",
				slog.String("request_id", requestID),
				slog.String("method", c.Method()),
				slog.String("path", c.Path()),
				slog.Int("status", code),
				slog.String("error", err.Error()),
			)

			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())

	// Structured request logging with request IDs (all environments)
	app.Use(appmiddleware.RequestLogger())

	// Additional human-readable logger in development
	if cfg.IsDevelopment() {
		app.Use(logger.New(logger.Config{
			Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
		}))
	}

	// CORS configuration from config
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(cfg.Server.CORS.AllowOrigins, ","),
		AllowMethods:     strings.Join(cfg.Server.CORS.AllowMethods, ","),
		AllowHeaders:     strings.Join(cfg.Server.CORS.AllowHeaders, ","),
		AllowCredentials: cfg.Server.CORS.AllowCredentials,
		MaxAge:           cfg.Server.CORS.MaxAge,
	}))

	// Setup all routes
	routes.SetupRoutes(app, DB, cfg)

	// Start server in goroutine
	go func() {
		addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		slog.Info("server listening", slog.String("address", addr))
		if err := app.Listen(addr); err != nil {
			slog.Error("failed to start server", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")

	// Shutdown server first to drain in-flight requests
	if err := app.ShutdownWithTimeout(time.Duration(cfg.Server.ShutdownTimeout) * time.Second); err != nil {
		slog.Error("server forced to shutdown", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Then close database connections
	if sqlDB, err := DB.DB(); err == nil {
		sqlDB.Close()
	}

	slog.Info("server exited successfully")
}

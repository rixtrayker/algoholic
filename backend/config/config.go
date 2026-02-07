package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

// Config holds all application configuration
type Config struct {
	App      AppConfig      `koanf:"app"`
	Server   ServerConfig   `koanf:"server"`
	Database DatabaseConfig `koanf:"database"`
	Redis    RedisConfig    `koanf:"redis"`
	ChromaDB ChromaDBConfig `koanf:"chromadb"`
	Ollama   OllamaConfig   `koanf:"ollama"`
	RAG      RAGConfig      `koanf:"rag"`
	Auth     AuthConfig     `koanf:"auth"`
	Logging  LoggingConfig  `koanf:"logging"`
}

// AppConfig contains general application settings
type AppConfig struct {
	Name        string `koanf:"name"`
	Version     string `koanf:"version"`
	Environment string `koanf:"environment"` // development, staging, production
	Debug       bool   `koanf:"debug"`
}

// ServerConfig contains HTTP server settings
type ServerConfig struct {
	Host            string `koanf:"host"`
	Port            int    `koanf:"port"`
	ReadTimeout     int    `koanf:"read_timeout"`     // seconds
	WriteTimeout    int    `koanf:"write_timeout"`    // seconds
	ShutdownTimeout int    `koanf:"shutdown_timeout"` // seconds
	CORS            CORSConfig
}

// CORSConfig contains CORS settings
type CORSConfig struct {
	AllowOrigins     []string `koanf:"allow_origins"`
	AllowMethods     []string `koanf:"allow_methods"`
	AllowHeaders     []string `koanf:"allow_headers"`
	AllowCredentials bool     `koanf:"allow_credentials"`
	MaxAge           int      `koanf:"max_age"`
}

// DatabaseConfig contains PostgreSQL settings
type DatabaseConfig struct {
	Host            string `koanf:"host"`
	Port            int    `koanf:"port"`
	User            string `koanf:"user"`
	Password        string `koanf:"password"`
	Database        string `koanf:"database"`
	SSLMode         string `koanf:"sslmode"`
	MaxOpenConns    int    `koanf:"max_open_conns"`
	MaxIdleConns    int    `koanf:"max_idle_conns"`
	ConnMaxLifetime int    `koanf:"conn_max_lifetime"` // minutes
	LogLevel        string `koanf:"log_level"`          // silent, error, warn, info
	AutoMigrate     bool   `koanf:"auto_migrate"`       // development only
}

// RedisConfig contains Redis cache settings
type RedisConfig struct {
	Enabled  bool   `koanf:"enabled"`
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Password string `koanf:"password"`
	DB       int    `koanf:"db"`
	TTL      int    `koanf:"ttl"` // seconds
}

// ChromaDBConfig contains vector database settings
type ChromaDBConfig struct {
	URL        string `koanf:"url"`
	Timeout    int    `koanf:"timeout"` // seconds
	BatchSize  int    `koanf:"batch_size"`
	Collection ChromaCollections
}

// ChromaCollections defines vector database collections
type ChromaCollections struct {
	Problems   string `koanf:"problems"`
	Questions  string `koanf:"questions"`
	Solutions  string `koanf:"solutions"`
	Templates  string `koanf:"templates"`
	Embeddings string `koanf:"embeddings"`
}

// OllamaConfig contains local LLM settings
type OllamaConfig struct {
	URL               string  `koanf:"url"`
	Timeout           int     `koanf:"timeout"` // seconds
	AssessmentModel   string  `koanf:"assessment_model"`
	GenerationModel   string  `koanf:"generation_model"`
	EmbeddingModel    string  `koanf:"embedding_model"`
	Temperature       float64 `koanf:"temperature"`
	AssessmentTemp    float64 `koanf:"assessment_temp"`
	GenerationTemp    float64 `koanf:"generation_temp"`
	MaxTokens         int     `koanf:"max_tokens"`
	ContextWindow     int     `koanf:"context_window"`
}

// RAGConfig contains RAG pipeline settings
type RAGConfig struct {
	Enabled           bool    `koanf:"enabled"`
	TopK              int     `koanf:"top_k"`
	MinSimilarity     float64 `koanf:"min_similarity"`
	MaxContextLength  int     `koanf:"max_context_length"`
	RetrievalStrategy string  `koanf:"retrieval_strategy"` // semantic, hybrid, keyword
	RerankingEnabled  bool    `koanf:"reranking_enabled"`
}

// AuthConfig contains authentication settings
type AuthConfig struct {
	Enabled         bool   `koanf:"enabled"`
	JWTSecret       string `koanf:"jwt_secret"`
	JWTExpiry       int    `koanf:"jwt_expiry"`        // hours
	RefreshExpiry   int    `koanf:"refresh_expiry"`    // hours
	BCryptCost      int    `koanf:"bcrypt_cost"`
	SessionDuration int    `koanf:"session_duration"`  // hours
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	Level      string `koanf:"level"`       // debug, info, warn, error
	Format     string `koanf:"format"`      // json, text
	Output     string `koanf:"output"`      // stdout, file
	FilePath   string `koanf:"file_path"`
	MaxSize    int    `koanf:"max_size"`    // megabytes
	MaxBackups int    `koanf:"max_backups"`
	MaxAge     int    `koanf:"max_age"`     // days
	Compress   bool   `koanf:"compress"`
}

// Global koanf instance
var k = koanf.New(".")

// Load loads configuration from multiple sources
func Load(configPath string) (*Config, error) {
	// Load default configuration
	if err := k.Load(structs.Provider(defaultConfig(), "koanf"), nil); err != nil {
		return nil, fmt.Errorf("error loading default config: %w", err)
	}

	// Load from YAML config file (if exists)
	if configPath != "" {
		if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
			log.Printf("Warning: could not load config file %s: %v", configPath, err)
		} else {
			log.Printf("Loaded configuration from: %s", configPath)
		}
	}

	// Load from environment variables (highest priority)
	// Environment variables should be prefixed with ALGOHOLIC_
	// e.g., ALGOHOLIC_DATABASE_HOST=localhost
	if err := k.Load(env.Provider("ALGOHOLIC_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "ALGOHOLIC_")), "_", ".", -1)
	}), nil); err != nil {
		return nil, fmt.Errorf("error loading environment variables: %w", err)
	}

	// Unmarshal into Config struct
	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Database validation
	if c.Database.Host == "" {
		return fmt.Errorf("database.host is required")
	}
	if c.Database.Database == "" {
		return fmt.Errorf("database.database is required")
	}

	// Server validation
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("server.port must be between 1 and 65535")
	}

	// Auth validation
	if c.Auth.Enabled && c.Auth.JWTSecret == "" {
		return fmt.Errorf("auth.jwt_secret is required when auth is enabled")
	}

	// Environment validation
	validEnvs := map[string]bool{"development": true, "staging": true, "production": true}
	if !validEnvs[c.App.Environment] {
		return fmt.Errorf("app.environment must be one of: development, staging, production")
	}

	return nil
}

// GetDSN returns the PostgreSQL connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode)
}

// GetRedisAddr returns the Redis connection address
func (c *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// IsDevelopment checks if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// IsProduction checks if running in production mode
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// defaultConfig returns default configuration values
func defaultConfig() Config {
	return Config{
		App: AppConfig{
			Name:        "Algoholic API",
			Version:     "1.0.0",
			Environment: "development",
			Debug:       true,
		},
		Server: ServerConfig{
			Host:            "0.0.0.0",
			Port:            4000,
			ReadTimeout:     30,
			WriteTimeout:    30,
			ShutdownTimeout: 10,
			CORS: CORSConfig{
				AllowOrigins:     []string{"*"},
				AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
				AllowHeaders:     []string{"*"},
				AllowCredentials: false,
				MaxAge:           3600,
			},
		},
		Database: DatabaseConfig{
			Host:            "localhost",
			Port:            5432,
			User:            "leetcode",
			Password:        "leetcode123",
			Database:        "leetcode_training",
			SSLMode:         "disable",
			MaxOpenConns:    25,
			MaxIdleConns:    5,
			ConnMaxLifetime: 30,
			LogLevel:        "warn",
			AutoMigrate:     false,
		},
		Redis: RedisConfig{
			Enabled:  false,
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
			TTL:      3600,
		},
		ChromaDB: ChromaDBConfig{
			URL:       "http://localhost:8000",
			Timeout:   30,
			BatchSize: 100,
			Collection: ChromaCollections{
				Problems:   "problems",
				Questions:  "questions",
				Solutions:  "solutions",
				Templates:  "templates",
				Embeddings: "embeddings",
			},
		},
		Ollama: OllamaConfig{
			URL:             "http://localhost:11434",
			Timeout:         120,
			AssessmentModel: "mistral:7b",
			GenerationModel: "codellama:13b",
			EmbeddingModel:  "all-minilm",
			Temperature:     0.7,
			AssessmentTemp:  0.3,
			GenerationTemp:  0.8,
			MaxTokens:       2048,
			ContextWindow:   4096,
		},
		RAG: RAGConfig{
			Enabled:           true,
			TopK:              5,
			MinSimilarity:     0.7,
			MaxContextLength:  4000,
			RetrievalStrategy: "semantic",
			RerankingEnabled:  false,
		},
		Auth: AuthConfig{
			Enabled:         false,
			JWTSecret:       "change-me-in-production",
			JWTExpiry:       24,
			RefreshExpiry:   168,
			BCryptCost:      10,
			SessionDuration: 24,
		},
		Logging: LoggingConfig{
			Level:      "info",
			Format:     "json",
			Output:     "stdout",
			FilePath:   "logs/app.log",
			MaxSize:    100,
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   true,
		},
	}
}

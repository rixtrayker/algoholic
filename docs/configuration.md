# Configuration Guide

Algoholic uses [koanf](https://github.com/knadh/koanf) for flexible, layered configuration management.

## Configuration Priority

Configuration values are loaded in the following order (highest priority last):

1. **Default values** - Built-in defaults in `config/config.go`
2. **YAML file** - `config.yaml` in the backend directory
3. **Environment variables** - Prefixed with `ALGOHOLIC_`

This means environment variables override YAML config, which overrides defaults.

## Quick Start

```bash
# Copy example files
cd backend
cp .env.example .env
cp config.yaml config.local.yaml  # optional: for local overrides

# Edit configuration
vim .env  # or config.local.yaml

# Start server (automatically loads config)
go run main.go
```

## Configuration Sources

### 1. Environment Variables

All environment variables must be prefixed with `ALGOHOLIC_` and use underscores for nested values:

```bash
# Format: ALGOHOLIC_<section>_<key>
ALGOHOLIC_SERVER_PORT=5000
ALGOHOLIC_DATABASE_HOST=db.example.com
ALGOHOLIC_DATABASE_PORT=5432
ALGOHOLIC_OLLAMA_URL=http://ollama:11434
ALGOHOLIC_RAG_TOP_K=10
```

**Example `.env` file:**

```bash
ALGOHOLIC_APP_ENVIRONMENT=production
ALGOHOLIC_SERVER_PORT=4000
ALGOHOLIC_DATABASE_HOST=localhost
ALGOHOLIC_DATABASE_PASSWORD=secure-password-here
ALGOHOLIC_AUTH_JWT_SECRET=your-secret-key-here
```

### 2. YAML Configuration

Create `config.yaml` in the backend directory:

```yaml
app:
  name: "Algoholic API"
  version: "1.0.0"
  environment: "development"
  debug: true

server:
  host: "0.0.0.0"
  port: 4000
  read_timeout: 30
  write_timeout: 30

database:
  host: "localhost"
  port: 5432
  user: "leetcode"
  password: "leetcode123"
  database: "leetcode_training"
  max_open_conns: 25
```

**Tip:** Use `config.local.yaml` for local overrides (git-ignored).

### 3. Custom Config File Path

Override the config file location:

```bash
CONFIG_PATH=/path/to/custom-config.yaml go run main.go
```

## Configuration Sections

### App

General application settings:

```yaml
app:
  name: "Algoholic API"          # Application name
  version: "1.0.0"                # Version string
  environment: "development"      # development | staging | production
  debug: true                     # Enable debug mode
```

### Server

HTTP server configuration:

```yaml
server:
  host: "0.0.0.0"                # Bind address
  port: 4000                      # Port number
  read_timeout: 30                # Request read timeout (seconds)
  write_timeout: 30               # Response write timeout (seconds)
  shutdown_timeout: 10            # Graceful shutdown timeout (seconds)
  cors:
    allow_origins: ["*"]
    allow_methods: ["GET", "POST", "PUT", "DELETE"]
    allow_headers: ["*"]
    allow_credentials: false
    max_age: 3600
```

### Database

PostgreSQL connection settings:

```yaml
database:
  host: "localhost"
  port: 5432
  user: "leetcode"
  password: "leetcode123"
  database: "leetcode_training"
  sslmode: "disable"             # disable | require | verify-ca | verify-full
  max_open_conns: 25             # Maximum open connections
  max_idle_conns: 5              # Maximum idle connections
  conn_max_lifetime: 30          # Connection lifetime (minutes)
  log_level: "warn"              # silent | error | warn | info
  auto_migrate: false            # DANGER: Only use in development
```

**DSN Format:** Automatically constructed from config values.

### Redis

Optional caching layer:

```yaml
redis:
  enabled: false                 # Enable/disable Redis
  host: "localhost"
  port: 6379
  password: ""                   # Leave empty if no auth
  db: 0                          # Redis database number
  ttl: 3600                      # Default TTL (seconds)
```

### ChromaDB

Vector database for semantic search:

```yaml
chromadb:
  url: "http://localhost:8000"
  timeout: 30                    # Request timeout (seconds)
  batch_size: 100                # Batch insert size
  collection:
    problems: "problems"
    questions: "questions"
    solutions: "solutions"
    templates: "templates"
    embeddings: "embeddings"
```

### Ollama

Local LLM settings:

```yaml
ollama:
  url: "http://localhost:11434"
  timeout: 120                   # Request timeout (seconds)
  assessment_model: "mistral:7b" # Model for assessments
  generation_model: "codellama:13b" # Model for generation
  embedding_model: "all-minilm"  # Embedding model
  temperature: 0.7               # Default temperature
  assessment_temp: 0.3           # Lower for consistency
  generation_temp: 0.8           # Higher for creativity
  max_tokens: 2048
  context_window: 4096
```

### RAG

Retrieval-Augmented Generation pipeline:

```yaml
rag:
  enabled: true
  top_k: 5                       # Number of results to retrieve
  min_similarity: 0.7            # Minimum similarity threshold
  max_context_length: 4000       # Max context characters
  retrieval_strategy: "semantic" # semantic | hybrid | keyword
  reranking_enabled: false       # Enable result reranking
```

### Authentication

JWT and session settings:

```yaml
auth:
  enabled: false                 # Enable authentication
  jwt_secret: "change-me-in-production"  # JWT signing secret
  jwt_expiry: 24                 # JWT expiry (hours)
  refresh_expiry: 168            # Refresh token expiry (hours)
  bcrypt_cost: 10                # BCrypt hashing cost
  session_duration: 24           # Session duration (hours)
```

⚠️ **Security:** Always change `jwt_secret` in production!

### Logging

Logging configuration:

```yaml
logging:
  level: "info"                  # debug | info | warn | error
  format: "json"                 # json | text
  output: "stdout"               # stdout | file
  file_path: "logs/app.log"      # Log file path (if output=file)
  max_size: 100                  # Max log file size (MB)
  max_backups: 5                 # Number of log file backups
  max_age: 30                    # Max age of log files (days)
  compress: true                 # Compress old log files
```

## Environment-Specific Configuration

### Development

```yaml
app:
  environment: "development"
  debug: true

database:
  log_level: "info"              # Verbose logging
  auto_migrate: true             # Auto-migrate schema

logging:
  level: "debug"
  format: "text"                 # Human-readable
```

### Production

```yaml
app:
  environment: "production"
  debug: false

database:
  log_level: "error"             # Minimal logging
  auto_migrate: false            # NEVER auto-migrate in prod
  max_open_conns: 100            # Higher connection pool

auth:
  enabled: true
  jwt_secret: "${JWT_SECRET}"    # Use env var

logging:
  level: "warn"
  format: "json"                 # Structured logs
  output: "file"
```

## Best Practices

### Security

1. **Never commit secrets** - Use environment variables for sensitive data
2. **Change default secrets** - Especially `jwt_secret` in production
3. **Use strong passwords** - For database and Redis
4. **Enable SSL in production** - Set `database.sslmode: require`
5. **Restrict CORS** - Don't use `allow_origins: ["*"]` in production

### Performance

1. **Tune connection pools** - Adjust `max_open_conns` based on load
2. **Enable Redis** - Cache frequently accessed data
3. **Adjust timeouts** - Balance between UX and resource usage
4. **Monitor log levels** - Reduce logging overhead in production

### Development

1. **Use local config** - Create `config.local.yaml` (git-ignored)
2. **Enable debug mode** - Get detailed error messages
3. **Use auto-migrate** - Convenient for local schema changes
4. **Verbose logging** - Set `logging.level: debug`

## Configuration Validation

The application validates configuration on startup:

- Required fields must be present
- Port numbers must be valid (1-65535)
- Environment must be valid (development/staging/production)
- JWT secret required when auth is enabled

Validation errors will prevent startup with clear error messages.

## Helper Methods

In code, use configuration helper methods:

```go
// Check environment
if cfg.IsDevelopment() {
    // Development-only code
}

if cfg.IsProduction() {
    // Production-only code
}

// Get DSN string
dsn := cfg.Database.GetDSN()

// Get Redis address
addr := cfg.Redis.GetAddr()
```

## Viewing Current Configuration

In development mode, you can view the current configuration:

```bash
curl http://localhost:4000/api/config
```

Note: This endpoint only returns non-sensitive values and is disabled in production.

## Troubleshooting

### Configuration not loading

1. Check config file path: `CONFIG_PATH` env var or `config.yaml` in cwd
2. Verify YAML syntax: Use a YAML validator
3. Check environment variable format: Must be prefixed with `ALGOHOLIC_`
4. Review startup logs: Configuration loading is logged

### Environment variables not working

```bash
# Correct format
export ALGOHOLIC_SERVER_PORT=5000

# Incorrect formats (won't work)
export SERVER_PORT=5000           # Missing prefix
export ALGOHOLIC_SERVERPORT=5000  # Missing underscore
export algoholic_server_port=5000 # Wrong case
```

### Database connection fails

1. Verify credentials in config/env
2. Check PostgreSQL is running: `docker ps`
3. Test connection manually: `psql -h localhost -U leetcode -d leetcode_training`
4. Review `database.log_level: info` for detailed errors

## Example Configurations

### Docker Compose

```yaml
services:
  backend:
    environment:
      ALGOHOLIC_DATABASE_HOST: postgres
      ALGOHOLIC_CHROMADB_URL: http://chromadb:8000
      ALGOHOLIC_OLLAMA_URL: http://ollama:11434
      ALGOHOLIC_REDIS_HOST: redis
      ALGOHOLIC_REDIS_ENABLED: true
```

### Kubernetes ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: algoholic-config
data:
  config.yaml: |
    app:
      environment: "production"
    server:
      port: 4000
    database:
      host: "postgres-service"
```

### Systemd Service

```ini
[Service]
Environment="ALGOHOLIC_APP_ENVIRONMENT=production"
Environment="ALGOHOLIC_SERVER_PORT=4000"
EnvironmentFile=/etc/algoholic/.env
ExecStart=/usr/local/bin/algoholic
```

## See Also

- [koanf documentation](https://github.com/knadh/koanf)
- [Getting Started Guide](./getting-started.md)
- [Architecture Overview](./architecture.md)

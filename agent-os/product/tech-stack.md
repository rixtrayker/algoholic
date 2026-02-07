# Tech Stack

## Frontend (Phase 2 - Planned)

**Framework**: React 18+ with TypeScript
- Modern hooks-based architecture
- Type safety for better developer experience
- Rich ecosystem of libraries

**State Management**: React Query + Zustand
- React Query for server state (API calls, caching)
- Zustand for client state (UI state, preferences)

**Styling**: Tailwind CSS
- Utility-first approach
- Rapid prototyping
- Consistent design system

**UI Components**: shadcn/ui
- Accessible components built on Radix UI
- Customizable and themeable
- Copy-paste component approach

**Build Tool**: Vite
- Fast HMR (Hot Module Replacement)
- Optimized production builds
- TypeScript out-of-the-box

## Backend (✅ Phase 1 Complete)

**Framework**: Go 1.21+ with Fiber v2
- High performance HTTP server
- Express-like API for Node.js familiarity
- Low memory footprint

**ORM**: GORM
- Type-safe database operations
- Auto-migrations in development
- Support for complex queries

**Authentication**: JWT (golang-jwt/jwt/v5)
- Stateless authentication
- Token-based authorization
- BCrypt password hashing

**Configuration**: Koanf
- Multi-source config (defaults → YAML → env vars)
- Type-safe configuration structs
- Environment-specific settings

## Database

**Primary Database**: PostgreSQL 16+ with Extensions
- Relational data (users, problems, questions, attempts)
- ACID compliance for critical data
- Extensions:
  - **Apache AGE**: Graph queries for topic relationships and learning paths
  - **pgvector**: Native vector search support
  - **pg_trgm**: Trigram-based fuzzy text search

**Vector Database**: ChromaDB (Phase 2)
- Embeddings storage for semantic search
- Similarity-based recommendations
- Collections: problems, questions, solutions, templates

**Cache**: Redis 7+ (Phase 2)
- Session storage
- Query result caching
- Rate limiting
- Real-time features (leaderboards, notifications)

## AI/ML

**LLM Runtime**: Ollama (Phase 2)
- Local LLM execution (no API costs)
- Models:
  - mistral:7b for assessment and hints
  - codellama:13b for code generation
  - all-minilm for embeddings

**Embedding Model**: sentence-transformers
- all-MiniLM-L6-v2 (384 dims, fast)
- bge-large-en-v1.5 (1024 dims, higher quality option)

**RAG Pipeline**:
```
Query → Embed → Vector Search (ChromaDB)
      → Retrieve Context (PostgreSQL)
      → Augment Prompt → LLM (Ollama)
      → Response
```

## DevOps & Infrastructure

**Migrations**: golang-migrate
- Version-controlled SQL migrations
- Up/down migration support
- Hybrid with GORM AutoMigrate for development

**Containerization**: Docker + Docker Compose
- Multi-service orchestration
- Consistent dev/prod environments
- Services: PostgreSQL, ChromaDB, Ollama, Redis, Backend, Frontend

**Version Control**: Git + GitHub
- Feature branch workflow
- Conventional commits
- CI/CD with GitHub Actions (planned)

**Monitoring** (Phase 3):
- Prometheus for metrics
- Grafana for dashboards
- Sentry for error tracking

## Development Tools

**API Documentation**: API reference markdown + development endpoint
- `/api/config` for dev environment inspection
- Comprehensive docs in `docs/api-reference.md`

**Testing**:
- Go: `testify` for assertions, in-memory SQLite for tests
- Frontend (planned): Vitest + React Testing Library

**Linting/Formatting**:
- Go: `gofmt`, `golangci-lint`
- Frontend: ESLint, Prettier

## Architecture Patterns

**Backend**:
- Clean architecture with separation of concerns
- Layers: Models → Services → Handlers → Routes
- Dependency injection via constructors
- Middleware for cross-cutting concerns (auth, logging)

**Frontend** (planned):
- Component-based architecture
- Custom hooks for business logic
- API client abstraction layer
- Responsive design (mobile-first)

**Database**:
- PostgreSQL for transactional data
- Graph database (AGE) for relationships
- Vector database (ChromaDB) for semantic search
- Redis for caching and real-time features

## Key Technology Decisions

### Why Go + Fiber?
- **Performance**: 10x faster than Node.js, low memory usage
- **Concurrency**: Native goroutines for handling multiple requests
- **Type Safety**: Compile-time error checking
- **Deployment**: Single binary, no runtime dependencies

### Why PostgreSQL + AGE?
- **Versatility**: Relational + graph + vector in one database
- **ACID**: Strong consistency guarantees
- **Extensions**: Extensible with custom functionality
- **Cost**: Open source, no licensing fees

### Why Local LLM (Ollama)?
- **Cost**: No API fees, unlimited usage
- **Privacy**: Data stays local
- **Control**: Model selection and fine-tuning
- **Performance**: Low latency for real-time features

### Why React?
- **Ecosystem**: Largest component library ecosystem
- **Talent**: Easier to hire React developers
- **Flexibility**: Can integrate any library
- **Longevity**: Backed by Meta, stable long-term

## Deployment Architecture (Planned)

```
┌─────────────────────────────────────────┐
│   CDN (Frontend - Vercel/Netlify)      │
└───────────────┬─────────────────────────┘
                │
┌───────────────▼─────────────────────────┐
│   Load Balancer (nginx)                 │
└──┬────────────┬────────────┬────────────┘
   │            │            │
┌──▼──────┐ ┌──▼──────┐ ┌──▼──────┐
│Backend  │ │Backend  │ │Backend  │
│Instance │ │Instance │ │Instance │
└──┬──────┘ └──┬──────┘ └──┬──────┘
   │            │            │
┌──▼────────────▼────────────▼──────┐
│   PostgreSQL (Primary + Replica)  │
└───────────────────────────────────┘
         │             │
┌────────▼──────┐ ┌───▼─────────┐
│   ChromaDB    │ │   Redis     │
└───────────────┘ └─────────────┘
```

## Security

- **Authentication**: JWT with httpOnly cookies
- **Password Storage**: BCrypt with cost factor 10
- **SQL Injection**: Prevented by GORM parameterized queries
- **CORS**: Configurable allowed origins
- **Rate Limiting**: Redis-based (Phase 2)
- **Input Validation**: Request validation at handler level
- **HTTPS**: TLS 1.3 in production

## Scalability Considerations

- **Horizontal Scaling**: Stateless backend (multiple instances behind load balancer)
- **Database**: Read replicas for query scaling
- **Caching**: Redis for frequently accessed data
- **CDN**: Static assets served from edge locations
- **Vector Search**: ChromaDB supports distributed mode

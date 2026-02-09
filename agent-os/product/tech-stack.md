# Tech Stack

**✅ Current Status**: Phase 1 Complete - All Core Services Running

## Frontend ✅ **RUNNING**

**Framework**: React 19 with TypeScript
- Modern hooks-based architecture
- Type safety for better developer experience
- Rich ecosystem of libraries

**State Management**: React Query v5 + Zustand
- React Query for server state (API calls, caching) ✅
- Zustand for client state (UI state, preferences) ✅

**Styling**: Tailwind CSS v4
- Utility-first approach ✅
- Rapid prototyping ✅
- Consistent design system ✅

**UI Components**: shadcn/ui
- Accessible components built on Radix UI ✅
- Customizable and themeable ✅
- Copy-paste component approach ✅

**Build Tool**: Vite 7.3.1
- Fast HMR (Hot Module Replacement) ✅
- Optimized production builds ✅
- TypeScript out-of-the-box ✅
- Running on http://localhost:5173 ✅

**Testing**: Vitest + React Testing Library
- 46/46 tests passing ✅

## Backend ✅ **RUNNING**

**Framework**: Go 1.21+ with Fiber v2.52.11
- High performance HTTP server ✅
- Express-like API for Node.js familiarity ✅
- Low memory footprint ✅
- 65 routes registered ✅
- Running on http://localhost:4000 ✅

**ORM**: GORM
- Type-safe database operations ✅
- Auto-migrations in development ✅
- Support for complex queries ✅
- PostgreSQL driver ✅

**Authentication**: JWT (golang-jwt/jwt/v5)
- Stateless authentication ✅
- Token-based authorization ✅
- BCrypt password hashing (cost factor 10) ✅
- Protected endpoints with middleware ✅

**Configuration**: Koanf
- Multi-source config (defaults → YAML → env vars) ✅
- Type-safe configuration structs ✅
- Environment-specific settings ✅

**API**: 22 REST endpoints
- Authentication (3 endpoints) ✅
- Questions (4 endpoints) ✅
- Problems (4 endpoints) ✅
- User Stats (3 endpoints) ✅
- Training Plans (5 endpoints) ✅
- Topics (2 endpoints) ✅
- Health Check (1 endpoint) ✅

## Database

**Primary Database**: PostgreSQL 16 ✅ **RUNNING**
- Relational data (users, problems, questions, attempts) ✅
- ACID compliance for critical data ✅
- Database: `leetcode_training` ✅
- User: `leetcode` ✅
- Port: 5432 ✅
- Connection pool: 10 connections ✅
- Auto-migration via GORM ✅
- Extensions (Planned for Phase 2):
  - **Apache AGE**: Graph queries for topic relationships and learning paths ⏭
  - **pgvector**: Native vector search support ⏭
  - **pg_trgm**: Trigram-based fuzzy text search ⏭

**Vector Database**: ChromaDB ⏭ (Phase 2)
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

**Migrations**: GORM Auto-Migration ✅
- Auto-migration on backend startup ✅
- Development-friendly workflow ✅
- (golang-migrate planned for production) ⏭

**Containerization**: Docker + Docker Compose ⏭ (Planned)
- Currently running natively on macOS ✅
- Multi-service orchestration (planned)
- Consistent dev/prod environments (planned)
- Future services: PostgreSQL, ChromaDB, Ollama, Redis, Backend, Frontend

**Version Control**: Git + GitHub ✅
- Feature branch workflow ✅
- Conventional commits ✅
- CI/CD with GitHub Actions ⏭ (Planned)

## Testing & Quality Assurance ✅

**API Testing**: Postman + Newman
- Collection: 22 endpoints with comprehensive tests ✅
- Newman CLI runner v6.2.1 ✅
- Test scripts: ~85 assertions ✅
- Automatic variable management (tokens, IDs) ✅
- Reports: JSON, JUnit XML, HTML ✅
- Run: `cd postman && ./run-tests.sh` ✅

**Frontend Testing**: Vitest + React Testing Library
- Unit and integration tests ✅
- 46/46 tests passing ✅
- Coverage reports ✅
- Run: `cd frontend && npm test` ✅

**API Documentation**:
- Postman collection as source of truth ✅
- Markdown API reference ✅
- Frontend-API alignment documentation ✅
- See: `postman/README.md`, `docs/api-reference.md` ✅

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

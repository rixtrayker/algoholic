# Algoholic

> A LeetCode-style practice platform with spaced repetition, powered by Go, React, and PostgreSQL.

## ✅ Current Status: Phase 1 Complete - All Systems Running

| Service | Status | URL | Details |
|---------|--------|-----|---------|
| **Backend API** | 🟢 Running | http://localhost:4000 | Go Fiber v2.52.11, 65 routes |
| **Frontend** | 🟢 Running | http://localhost:5173 | React 19 + Vite 7.3.1 |
| **Database** | 🟢 Connected | postgresql://localhost:5432 | PostgreSQL (leetcode_training) |
| **API Tests** | ✅ Passing | - | 22 endpoints via Postman + Newman |

**See**: [RUNNING.md](./RUNNING.md) for detailed service information and management.

---

## Features

### ✅ Implemented (Phase 1)
- 🔐 **JWT Authentication** - Secure user registration and login
- 📝 **Question System** - Random questions with difficulty filters
- 📚 **Problem Library** - Search, filter, and browse problems
- 📊 **User Progress** - Track stats, attempts, and performance
- 🎯 **Training Plans** - Structured learning paths with enrollment
- 🏷️ **Topics** - Categorized learning with performance tracking
- 🧪 **Comprehensive Testing** - 22 API endpoints + 46 frontend tests

### 🔄 In Progress (Phase 5)
- 💻 **Enhanced UI** - Practice interface improvements
- 📈 **Dashboard** - Progress visualization and analytics

### 📋 Planned (Phases 2-4)
- 🧠 **AI Assessment** - LLM-powered code evaluation (Ollama)
- 🔍 **Semantic Search** - Vector-based problem discovery (ChromaDB)
- 🕸️ **Graph Relationships** - Topic dependencies (Apache AGE)
- 🔁 **Spaced Repetition** - SM-2 algorithm for optimal review
- 🎨 **Question Generation** - AI-generated practice problems

---

## Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL 16+
- Newman 6.2.1+ (for API testing)

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd algoholic

# Backend setup
cd backend
go mod download
go run main.go
# → Running on http://localhost:4000

# Frontend setup (in another terminal)
cd frontend
npm install
npm run dev
# → Running on http://localhost:5173

# Database is already configured
# User: leetcode
# Database: leetcode_training
# Port: 5432
```

### Verify Setup

```bash
# Test backend health
curl http://localhost:4000/health

# Run API tests
cd postman
./run-tests.sh

# Run frontend tests
cd frontend
npm test
```

---

## API Documentation

**📦 Complete API Specification**: See [postman/algoholic-api.postman_collection.json](./postman/algoholic-api.postman_collection.json)

**22 Endpoints across 7 categories**:

```
Authentication (3)    Questions (4)         Problems (4)
User Stats (3)        Training Plans (5)    Topics (2)
Health Check (1)
```

**Quick Reference**:
- [API Reference Documentation](./docs/api-reference.md)
- [Postman Collection Guide](./postman/README.md)
- [API-Frontend Alignment Report](./docs/api-frontend-alignment.md)

**Testing**:
```bash
cd postman
./run-tests.sh        # Run all 22 endpoint tests
./list-endpoints.sh   # List all available endpoints
./demo-newman.sh      # See Newman usage examples
```

---

## Technology Stack

| Layer | Technology | Status |
|-------|-----------|--------|
| **Frontend** | React 19 + TypeScript + Vite 7.3.1 | ✅ Running |
| **Styling** | Tailwind CSS v4 + shadcn/ui | ✅ Active |
| **State** | Zustand + React Query v5 | ✅ Active |
| **Backend** | Go 1.21 + Fiber v2.52.11 | ✅ Running |
| **Database** | PostgreSQL 16 + GORM | ✅ Running |
| **Auth** | JWT with BCrypt | ✅ Active |
| **Testing** | Newman 6.2.1 + Vitest | ✅ Passing |
| **Vector DB** | ChromaDB | ⏭ Phase 2 |
| **Graph DB** | Apache AGE | ⏭ Phase 2 |
| **AI/LLM** | Ollama (Mistral, CodeLlama) | ⏭ Phase 2 |

**See**: [agent-os/product/tech-stack.md](./agent-os/product/tech-stack.md) for complete details.

---

## Project Structure

```
algoholic/
├── backend/              # Go Fiber API server
│   ├── config/          # Configuration management (Koanf)
│   ├── database/        # Database connection and setup
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # Auth, logging, CORS
│   ├── models/          # GORM database models
│   ├── routes/          # Route definitions
│   ├── services/        # Business logic
│   └── main.go          # Application entry point
│
├── frontend/            # React application
│   ├── src/
│   │   ├── components/  # Reusable UI components
│   │   ├── pages/       # Page components
│   │   ├── lib/         # API client and utilities
│   │   └── stores/      # Zustand state stores
│   └── package.json
│
├── postman/             # API testing suite
│   ├── algoholic-api.postman_collection.json
│   ├── algoholic-local.postman_environment.json
│   ├── run-tests.sh     # Run all API tests
│   ├── list-endpoints.sh
│   └── README.md
│
├── docs/                # Documentation
│   ├── architecture.md
│   ├── api-reference.md
│   ├── getting-started.md
│   ├── api-frontend-alignment.md
│   └── ...
│
├── db/                  # Database scripts and seed data
└── RUNNING.md          # Service management guide
```

---

## Documentation

### Getting Started
- [Getting Started Guide](./docs/getting-started.md) - Setup and installation
- [RUNNING.md](./RUNNING.md) - Service management and troubleshooting
- [Tech Stack](./agent-os/product/tech-stack.md) - Complete technology details

### Architecture
- [System Architecture](./docs/architecture.md) - Full system design
- [Question Design](./docs/question-design.md) - Question taxonomy and types
- [Topic Reference](./docs/topic-reference.md) - Learning topic structure

### API
- [API Reference](./docs/api-reference.md) - Endpoint documentation
- [Postman Collection](./postman/README.md) - API testing guide
- [API-Frontend Alignment](./docs/api-frontend-alignment.md) - Integration details

---

## Testing

### API Testing (Newman)
```bash
cd postman
./run-tests.sh
```
- 22 endpoints tested
- ~85 test assertions
- Automatic token management
- Reports in JSON, JUnit XML, HTML

### Frontend Testing (Vitest)
```bash
cd frontend
npm test
```
- 46/46 tests passing
- Unit and integration tests
- React Testing Library

---

## Development Workflow

### Making Backend Changes
```bash
cd backend
# Edit Go files
# Restart required (no hot reload)
pkill -f "go run main.go"
go run main.go
```

### Making Frontend Changes
```bash
cd frontend
# Edit files in src/
# Vite hot module reloading active
# Changes appear instantly in browser
```

### Running Both Services
```bash
# Terminal 1 - Backend
cd backend && go run main.go

# Terminal 2 - Frontend
cd frontend && npm run dev

# Terminal 3 - API Tests
cd postman && ./run-tests.sh
```

---

## Roadmap

### Phase 1: Foundation ✅ **COMPLETE**
- [x] PostgreSQL database setup
- [x] Go Fiber API with 22 endpoints
- [x] JWT authentication
- [x] React frontend with Tailwind CSS
- [x] Postman collection with tests
- [x] Frontend-API alignment

### Phase 2: Intelligence ⏭ **NEXT**
- [ ] ChromaDB integration
- [ ] Vector embeddings for semantic search
- [ ] Apache AGE graph database
- [ ] Topic relationship graph
- [ ] Similar problem recommendations

### Phase 3: Training 📋 **PLANNED**
- [ ] Enhanced training plan algorithms
- [ ] Spaced repetition (SM-2)
- [ ] Weakness detection system
- [ ] Progress tracking improvements
- [ ] Review queue management

### Phase 4: AI 📋 **PLANNED**
- [ ] Ollama local LLM integration
- [ ] RAG pipeline for context
- [ ] AI assessment of solutions
- [ ] Question variant generation
- [ ] Personalized hint system

### Phase 5: Frontend 🔄 **IN PROGRESS**
- [ ] Enhanced practice interface
- [ ] Interactive dashboard
- [ ] Code editor integration
- [ ] Progress visualization
- [ ] Mobile responsive design

### Phase 6: Polish 📋 **PLANNED**
- [ ] Performance optimization
- [ ] Difficulty score calibration
- [ ] Docker deployment
- [ ] Production configuration
- [ ] User acceptance testing

---

## Contributing

This is a personal learning project. Contributions, issues, and feature requests are welcome!

---

## License

[MIT](LICENSE)

---

## Acknowledgments

- Inspired by LeetCode and spaced repetition learning systems
- Built with modern web technologies and local AI capabilities

---

**Last Updated**: 2026-02-09
**Current Version**: 1.0.0 (Phase 1 Complete)
**Status**: All services running ✅

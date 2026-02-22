# Algoholic — Getting Started

## Setup Guide & Implementation Checklist

---

## Current Status

**✅ Phase 2 Complete - Intelligence Layer Active**

| Service | URL | Details |
|---------|-----|---------|
| Backend API | http://localhost:4000 | Go Fiber, 44+ endpoints |
| Web (Next.js) | http://localhost:3000 | Next.js 14 + App Router |
| Frontend (Legacy) | http://localhost:5173 | React 19 + Vite (deprecated) |
| Database | postgresql://localhost:5432 | PostgreSQL (leetcode_training) |
| API Tests | - | 30+ endpoints via Postman |

---

## Prerequisites

```bash
go version                # Go 1.21+
node --version            # Node.js 18+
psql --version           # PostgreSQL 16+
newman --version         # Newman 6.2.1+ for API testing
```

---

## Quick Start

### 1. Clone and Setup

```bash
git clone <repository-url>
cd algoholic
```

### 2. Backend Setup

```bash
cd backend
go mod download
go run main.go
# → Running on http://localhost:4000
```

### 3. Web Setup (Next.js 14 - Recommended)

```bash
cd web
npm install
npm run dev
# → Running on http://localhost:3000
```

### 4. Legacy Frontend (Vite - Deprecated)

```bash
cd frontend
npm install
npm run dev
# → Running on http://localhost:5173
```

### 5. Database

PostgreSQL should be running with:
- User: `leetcode`
- Database: `leetcode_training`
- Port: 5432

Tables are auto-migrated via GORM on backend startup.

---

## Verify Setup

```bash
# Test backend health
curl http://localhost:4000/health

# Run API tests
cd postman && ./run-tests.sh

# Build Next.js
cd web && npm run build
```

Expected health response:
```json
{
  "app": "Algoholic API",
  "environment": "development",
  "status": "healthy",
  "version": "1.0.0"
}
```

---

## Project Structure

```
algoholic/
├── backend/              # Go Fiber API server
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # Auth, logging, CORS
│   ├── models/          # GORM database models
│   ├── routes/          # Route definitions
│   └── services/        # Business logic + AI
│
├── web/                 # Next.js 14 application (Primary)
│   ├── src/
│   │   ├── app/         # App Router pages
│   │   ├── components/  # UI components
│   │   ├── lib/         # API client
│   │   └── stores/      # Zustand stores
│   └── public/          # Static assets, logo
│
├── frontend/            # React Vite application (Legacy)
│
├── postman/             # API testing suite
│
├── migrations/          # Database migrations
│
└── docs/                # Documentation
```

---

## API Endpoints

**44+ Endpoints across 9 categories**:

```
Authentication (4)     Questions (6)         Problems (5)
Users (9)             Training Plans (8)    Topics (4)
Lists (7)             Activity (4)          Search (2)
```

See [docs/api-reference.md](./api-reference.md) for complete documentation.

---

## Implementation Checklist

### Phase 1: Foundation ✅ COMPLETE
- [x] PostgreSQL database
- [x] Go Fiber API
- [x] JWT authentication
- [x] React frontend

### Phase 2: Intelligence ✅ COMPLETE
- [x] ChromaDB integration
- [x] Vector embeddings
- [x] Apache AGE graph database
- [x] Semantic search endpoints
- [x] Learning path recommendations

### Phase 3: Frontend Migration ✅ COMPLETE
- [x] Migrate to Next.js 14
- [x] App Router architecture
- [x] All pages with proper API integration
- [x] Custom branding and logo
- [x] PWA manifest

### Phase 4: Enhanced Training 📋 NEXT
- [ ] SM-2 spaced repetition algorithm
- [ ] Adaptive difficulty calibration
- [ ] Code editor integration
- [ ] AI code assessment

---

## Configuration

Backend uses environment variables (prefix: `ALGOHOLIC_`):

```bash
ALGOHOLIC_SERVER_PORT=4000
ALGOHOLIC_DATABASE_HOST=localhost
ALGOHOLIC_CHROMADB_URL=http://localhost:8000
ALGOHOLIC_OLLAMA_URL=http://localhost:11434
```

---

## Useful Commands

```bash
# Backend
cd backend && go run main.go

# Web (Next.js)
cd web && npm run dev

# Legacy Frontend
cd frontend && npm run dev

# API Tests
cd postman && ./run-tests.sh

# Build for production
cd web && npm run build && npm run start
```

---

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Port in use | `lsof -i :4000` to find process |
| DB connection fails | Check PostgreSQL is running |
| Build errors | Delete `node_modules` and reinstall |

---

**Last Updated**: 2026-02-22

# Algoholic

> A LeetCode-style practice platform with spaced repetition, powered by Go, Next.js, and PostgreSQL.

## ✅ Current Status: Phase 2 Complete - Intelligence Layer Active

| Service | Status | URL | Details |
|---------|--------|-----|---------|
| **Backend API** | 🟢 Running | http://localhost:4000 | Go Fiber, 44+ routes |
| **Web (Next.js)** | 🟢 Running | http://localhost:3000 | Next.js 14 + App Router |
| **Frontend (Legacy)** | 🟢 Running | http://localhost:5173 | React 19 + Vite (deprecated) |
| **Database** | 🟢 Connected | postgresql://localhost:5432 | PostgreSQL (leetcode_training) |
| **API Tests** | ✅ Passing | - | 30+ endpoints via Postman |

---

## Features

### ✅ Implemented

**Core Platform**
- 🔐 **JWT Authentication** - Secure user registration and login
- 📝 **Question System** - Random questions with difficulty filters, hints
- 📚 **Problem Library** - Search, filter, browse with semantic search
- 📊 **User Progress** - Stats, attempts, performance tracking
- 🎯 **Training Plans** - Custom study plans with daily goals
- 📋 **User Lists** - Custom problem collections
- 📈 **Activity Tracking** - Commitment chart, practice history

**Intelligence Layer (Phase 2)**
- 🔍 **Semantic Search** - Vector-based problem discovery (ChromaDB)
- 🕸️ **Graph Relationships** - Topic prerequisites & learning paths (Apache AGE)
- 🧠 **Embedding Service** - Local embeddings via Ollama
- 📚 **RAG Pipeline** - Context-aware recommendations
- 🔄 **Review Queue** - Spaced repetition tracking

**Frontend (Next.js 14)**
- 🎨 **Modern UI** - Dark gradient auth, glassmorphism cards
- 📱 **Responsive** - Mobile-friendly navigation
- ⚡ **Optimized** - Static generation, React Query caching
- 🖼️ **Custom Branding** - SVG logo, PWA manifest

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

# Web setup (Next.js - Recommended)
cd web
npm install
npm run dev
# → Running on http://localhost:3000

# OR Legacy frontend (Vite)
cd frontend
npm install
npm run dev
# → Running on http://localhost:5173
```

### Verify Setup

```bash
# Test backend health
curl http://localhost:4000/health

# Run API tests
cd postman && ./run-tests.sh

# Run web build
cd web && npm run build
```

---

## API Documentation

**📦 Complete API Specification**: See [postman/algoholic-api.postman_collection.json](./postman/algoholic-api.postman_collection.json)

**44+ Endpoints across 9 categories**:

```
Authentication (4)     Questions (6)         Problems (5)
Users (9)             Training Plans (8)    Topics (4)
Lists (7)             Activity (4)          Search (2)
```

**Quick Reference**:
- [API Reference Documentation](./docs/api-reference.md)
- [Postman Collection Guide](./postman/README.md)

---

## Technology Stack

| Layer | Technology | Status |
|-------|-----------|--------|
| **Web** | Next.js 14 + TypeScript + App Router | ✅ Active |
| **Frontend (Legacy)** | React 19 + Vite | ⚠️ Deprecated |
| **Styling** | Tailwind CSS v4 | ✅ Active |
| **State** | Zustand + React Query v5 | ✅ Active |
| **Backend** | Go 1.21 + Fiber v2 | ✅ Running |
| **Database** | PostgreSQL 16 + GORM | ✅ Running |
| **Auth** | JWT with BCrypt | ✅ Active |
| **Vector DB** | ChromaDB | ✅ Active |
| **Graph DB** | Apache AGE | ✅ Active |
| **AI/LLM** | Ollama (Mistral) | ✅ Active |
| **Testing** | Newman + Vitest | ✅ Passing |

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
│   └── src/
│
├── postman/             # API testing suite
│
├── migrations/          # Database migrations
│
└── docs/                # Documentation
```

---

## Roadmap

### Phase 1: Foundation ✅ **COMPLETE**
- [x] PostgreSQL database setup
- [x] Go Fiber API
- [x] JWT authentication
- [x] React frontend

### Phase 2: Intelligence ✅ **COMPLETE**
- [x] ChromaDB integration
- [x] Vector embeddings
- [x] Apache AGE graph database
- [x] Semantic search endpoints
- [x] Learning path recommendations

### Phase 3: Frontend ✅ **COMPLETE**
- [x] Migrate to Next.js 14
- [x] App Router architecture
- [x] All pages with proper API integration
- [x] Custom branding and logo
- [x] PWA manifest

### Phase 4: Enhanced Training 📋 **NEXT**
- [ ] SM-2 spaced repetition algorithm
- [ ] Adaptive difficulty calibration
- [ ] Code editor integration
- [ ] AI code assessment
- [ ] Question variant generation

### Phase 5: Polish 📋 **PLANNED**
- [ ] Performance optimization
- [ ] Docker deployment
- [ ] Production configuration
- [ ] Mobile app (React Native)

---

## Contributing

This is a personal learning project. Contributions welcome!

## License

[MIT](LICENSE)

---

**Last Updated**: 2026-02-21
**Current Version**: 2.0.0
**Status**: All services running ✅

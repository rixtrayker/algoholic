# API-Frontend Alignment Report

This document details the alignment between the frontend API clients and the backend API.

## Current Status: ✅ Fully Aligned

All frontend API endpoints match the backend routes (44+ endpoints total).

---

## Frontend Implementations

### 1. Next.js Web Application (Primary) ✅

**Location**: `web/src/lib/api.ts`

Complete API client with:
- All 44+ endpoints implemented
- Proper TypeScript types
- JWT token handling
- Error interceptors

```typescript
// API modules
authAPI          // 4 endpoints
problemsAPI      // 5 endpoints
questionsAPI     // 6 endpoints
userAPI          // 9 endpoints
trainingPlansAPI // 8 endpoints
topicsAPI        // 4 endpoints
listsAPI         // 7 endpoints
activityAPI      // 4 endpoints
searchAPI        // 2 endpoints
graphAPI         // 1 endpoint
intelligenceAPI  // 1 endpoint
```

### 2. Legacy React Frontend ✅

**Location**: `frontend/src/lib/api.ts`

Maintained for backward compatibility with corrected endpoints:
- Training plans use `/training-plans` (fixed from `/plans`)
- Questions use `user_answer` in submit body
- Hint endpoint integrated

---

## Endpoint Mapping

### Authentication (4 endpoints)
| Backend Route | Next.js | Legacy | Status |
|--------------|---------|--------|--------|
| POST /api/auth/register | ✅ | ✅ | Aligned |
| POST /api/auth/login | ✅ | ✅ | Aligned |
| GET /api/auth/me | ✅ | ✅ | Aligned |
| POST /api/auth/change-password | ✅ | - | Aligned |

### Problems (5 endpoints)
| Backend Route | Next.js | Legacy | Status |
|--------------|---------|--------|--------|
| GET /api/problems | ✅ | ✅ | Aligned |
| GET /api/problems/:id | ✅ | ✅ | Aligned |
| GET /api/problems/slug/:slug | ✅ | - | Aligned |
| GET /api/problems/:id/topics | ✅ | ✅ | Aligned |
| GET /api/problems/:id/similar | ✅ | - | Aligned |

### Questions (6 endpoints)
| Backend Route | Next.js | Legacy | Status |
|--------------|---------|--------|--------|
| GET /api/questions | ✅ | ✅ | Aligned |
| GET /api/questions/random | ✅ | ✅ | Aligned |
| GET /api/questions/:id | ✅ | ✅ | Aligned |
| GET /api/questions/:id/hint | ✅ | ✅ | Aligned |
| POST /api/questions/:id/answer | ✅ | ✅ | Aligned |
| GET /api/questions/:id/attempts | ✅ | ✅ | Aligned |

### Users (9 endpoints)
| Backend Route | Next.js | Legacy | Status |
|--------------|---------|--------|--------|
| GET /api/users/me/stats | ✅ | ✅ | Aligned |
| GET /api/users/me/weaknesses | ✅ | ✅ | Aligned |
| GET /api/users/me/recommendations | ✅ | - | Aligned |
| GET /api/users/me/review-queue | ✅ | - | Aligned |
| GET /api/users/me/skills | ✅ | ✅ | Aligned |
| GET /api/users/me/skills/:topicId | ✅ | - | Aligned |
| GET /api/users/me/preferences | ✅ | - | Aligned |
| PUT /api/users/me/preferences | ✅ | - | Aligned |
| GET /api/users/me/attempts | ✅ | ✅ | Aligned |

### Training Plans (8 endpoints)
| Backend Route | Next.js | Legacy | Status |
|--------------|---------|--------|--------|
| POST /api/training-plans | ✅ | ✅ | Aligned |
| GET /api/training-plans | ✅ | ✅ | Aligned |
| GET /api/training-plans/:id | ✅ | ✅ | Aligned |
| GET /api/training-plans/:id/next | ✅ | ✅ | Aligned |
| GET /api/training-plans/:id/items | ✅ | ✅ | Aligned |
| GET /api/training-plans/:id/today | ✅ | ✅ | Aligned |
| POST /api/training-plans/:id/items/:itemId/complete | ✅ | ✅ | Aligned |
| POST /api/training-plans/:id/pause | ✅ | ✅ | Aligned |
| POST /api/training-plans/:id/resume | ✅ | ✅ | Aligned |
| DELETE /api/training-plans/:id | ✅ | ✅ | Aligned |

### Topics (4 endpoints)
| Backend Route | Next.js | Legacy | Status |
|--------------|---------|--------|--------|
| GET /api/topics | ✅ | ✅ | Aligned |
| GET /api/topics/:id | ✅ | - | Aligned |
| GET /api/topics/:id/prerequisites | ✅ | - | Aligned |
| GET /api/topics/:userId/performance/:topicId | ✅ | ✅ | Aligned |

### Lists (7 endpoints)
| Backend Route | Next.js | Legacy | Status |
|--------------|---------|--------|--------|
| GET /api/lists | ✅ | ✅ | Aligned |
| POST /api/lists | ✅ | ✅ | Aligned |
| GET /api/lists/:id | ✅ | ✅ | Aligned |
| PUT /api/lists/:id | ✅ | ✅ | Aligned |
| DELETE /api/lists/:id | ✅ | ✅ | Aligned |
| POST /api/lists/:id/problems | ✅ | ✅ | Aligned |
| DELETE /api/lists/:id/problems/:problemId | ✅ | ✅ | Aligned |

### Activity (4 endpoints)
| Backend Route | Next.js | Legacy | Status |
|--------------|---------|--------|--------|
| GET /api/activity/chart | ✅ | ✅ | Aligned |
| GET /api/activity/stats | ✅ | ✅ | Aligned |
| GET /api/activity/history | ✅ | ✅ | Aligned |
| POST /api/activity/record | ✅ | ✅ | Aligned |

### Search (2 endpoints)
| Backend Route | Next.js | Legacy | Status |
|--------------|---------|--------|--------|
| GET /api/search/problems | ✅ | ✅ | Aligned |
| GET /api/search/questions | ✅ | - | Aligned |

### Graph (1 endpoint)
| Backend Route | Next.js | Legacy | Status |
|--------------|---------|--------|--------|
| GET /api/graph/learning-path | ✅ | - | Aligned |

### Intelligence (1 endpoint)
| Backend Route | Next.js | Legacy | Status |
|--------------|---------|--------|--------|
| GET /api/intelligence/status | ✅ | - | Aligned |

---

## Request/Response Formats

### Submit Answer
```typescript
// Request
POST /api/questions/:id/answer
{
  "user_answer": { "answer": "A" },
  "time_taken_seconds": 45,
  "hints_used": 1,
  "confidence_level": 3,
  "training_plan_id": null
}

// Response
{
  "is_correct": true,
  "correct_answer": { "answer": "A" },
  "explanation": "...",
  "attempt_id": 123,
  "points_earned": 250
}
```

### Training Plans List
```typescript
// Response
{
  "plans": [...],
  "count": 3
}
```

---

## Type Definitions

All types are defined in `web/src/lib/api.ts`:

```typescript
interface User { ... }
interface Problem { ... }
interface Question { ... }
interface UserStats { ... }
interface TrainingPlan { ... }
interface Topic { ... }
interface UserList { ... }
interface WeakTopic { ... }
interface Recommendation { ... }
interface ReviewQueueItem { ... }
```

---

## Authentication Flow

1. User logs in via `/auth/login`
2. JWT token stored in `localStorage`
3. Token added to all requests via axios interceptor
4. 401 responses trigger redirect to login

```typescript
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('auth_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});
```

---

## Testing

### Newman API Tests
```bash
cd postman && ./run-tests.sh
```

### Next.js Build
```bash
cd web && npm run build
```

---

## Files

| File | Purpose |
|------|---------|
| `web/src/lib/api.ts` | Next.js API client (primary) |
| `frontend/src/lib/api.ts` | Legacy API client |
| `backend/routes/routes.go` | Backend route definitions |
| `postman/algoholic-api.postman_collection.json` | API specification |

---

**Last Updated**: 2026-02-22

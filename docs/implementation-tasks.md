# Implementation Plan — Backend Fixes & Missing Features

**Date:** February 24, 2026
**Priority:** Tier 1 (broken core logic) → Tier 2 (security) → Tier 3 (missing features)
**Principle:** Prefer popular, well-maintained Go packages. Keep changes minimal and focused.

---

## Phase 1: Fix Broken Core Logic (5 tasks)

### Task 1.1: Wire SubmitAnswer to UpdateUserProgress, Activity, and Streak

**Problem:** `QuestionService.SubmitAnswer()` records the attempt but never calls `UserService.UpdateUserProgress()`, `UserService.RecordDailyActivity()`, or `UserService.UpdateStreak()`. These functions all exist but are disconnected.

**File:** `backend/services/question_service.go`

**Changes:**

1. Add `userService *UserService` field to `QuestionService` struct
2. Update `NewQuestionService(db, userService)` constructor to accept `UserService`
3. In `SubmitAnswer()`, after line 149 (`UpdateQuestionStats`), add:

```go
// Update user proficiency for the question's topic
if question.ProblemID != nil {
    var pt models.ProblemTopic
    if err := s.db.Where("problem_id = ? AND is_primary = TRUE", *question.ProblemID).First(&pt).Error; err == nil {
        s.userService.UpdateUserProgress(userID, pt.TopicID, isCorrect, req.TimeTaken)
    }
}

// Record daily activity and update streak
s.userService.RecordDailyActivity(userID, req.TimeTaken)
s.userService.UpdateStreak(userID)
s.userService.AddStudyTime(userID, int64(req.TimeTaken))
```

4. Set `response.NewProficiencyLevel` from the updated skill

**Also update:** `backend/routes/routes.go` — pass `userService` when constructing `QuestionService`

**Impact:** Activity chart, streaks, proficiency, recommendations, and weakness detection all start working.

---

### Task 1.2: Implement SM-2 Spaced Repetition Service

**Problem:** The `SpacedRepetitionReview` and `ReviewQueue` models exist but no service implements the SM-2 algorithm. Nothing creates or updates these records.

**New file:** `backend/services/spaced_repetition_service.go`

**Package:** Use the standard SM-2 algorithm (no external package needed — it's a simple formula)

**Spec:**

```go
type SpacedRepetitionService struct {
    db *gorm.DB
}

// ProcessReview — called after each answer submission
// quality: 0-5 (0=complete blackout, 5=perfect recall)
// Maps from: incorrect=1, correct+slow=3, correct+fast=4, correct+confident=5
func (s *SpacedRepetitionService) ProcessReview(userID, questionID int, quality int) error

// SM-2 Algorithm:
// If quality < 3:
//   repetitions = 0, interval = 1
// Else:
//   If repetitions == 0: interval = 1
//   If repetitions == 1: interval = 6
//   Else: interval = round(interval * easinessFactor)
//   repetitions++
//
// EF' = EF + (0.1 - (5-quality) * (0.08 + (5-quality) * 0.02))
// EF = max(1.3, EF')
//
// next_review_at = now + interval days

// GetDueReviews — returns questions due for review
func (s *SpacedRepetitionService) GetDueReviews(userID int, limit int) ([]SpacedRepetitionReview, error)

// QualityFromAttempt — converts attempt data to SM-2 quality rating (0-5)
func QualityFromAttempt(isCorrect bool, timeTaken int, estimatedTime *int, hintsUsed int) int
```

**Wire into:** `QuestionService.SubmitAnswer()` — call `ProcessReview()` after recording the attempt

**Wire into:** `UserHandler.GetReviewQueue()` — call `GetDueReviews()` to populate the review queue endpoint

---

### Task 1.3: Fix Training Plan Adaptive Difficulty SQL

**Problem:** `training_plan_service.go:306-325` — The SQL `WHERE difficulty_score BETWEEN difficulty_score + 5 AND difficulty_score + 15` is self-referencing (column compared to itself). Always evaluates to false.

**File:** `backend/services/training_plan_service.go`

**Fix:** Replace the raw SQL with a proper subquery that gets the current question's difficulty, then finds a replacement:

```go
// Get current average difficulty of incomplete items
var avgDifficulty float64
s.db.Table("training_plan_items tpi").
    Joins("JOIN questions q ON q.question_id = tpi.question_id").
    Where("tpi.plan_id = ? AND tpi.is_completed = FALSE", planID).
    Select("AVG(q.difficulty_score)").
    Scan(&avgDifficulty)

// Calculate target difficulty
var targetMin, targetMax float64
if accuracy > 0.85 {
    targetMin = avgDifficulty + 5
    targetMax = avgDifficulty + 20
} else if accuracy < 0.40 {
    targetMin = math.Max(0, avgDifficulty-20)
    targetMax = math.Max(0, avgDifficulty-5)
}

// Replace questions for incomplete items
var replacementIDs []int
s.db.Table("questions").
    Select("question_id").
    Where("difficulty_score BETWEEN ? AND ?", targetMin, targetMax).
    Where("question_id NOT IN (?)",
        s.db.Table("training_plan_items").Select("question_id").Where("plan_id = ?", planID)).
    Order("RANDOM()").
    Limit(incompletePlanItemCount).
    Pluck("question_id", &replacementIDs)

// Update each incomplete item with a new question
```

---

### Task 1.4: Fix Training Plan Topic Conversion

**Problem:** `training_plan_service.go:59` — `string(rune(topic))` converts int to Unicode char.

**File:** `backend/services/training_plan_service.go`

**Fix:** One-line change:

```go
// Before (broken):
topicsArray[i] = string(rune(topic))

// After (fixed):
topicsArray[i] = strconv.Itoa(topic)
```

**Add import:** `"strconv"` to the file imports

---

### Task 1.5: Fix Code Execution Fallback (Return Error Instead of Silent Pass)

**Problem:** When Judge0 is down, `CheckCode()` falls back to `ValidateCode()` which accepts any code containing `def ` or `class `.

**File:** `backend/services/question_service.go` (lines 241-244)

**Fix:**

```go
// Before:
if err != nil {
    return executor.ValidateCode(code, language)
}

// After:
if err != nil {
    // Code execution service unavailable — don't silently pass
    return false
}
```

And update the `AnswerResponse` to include a `warning` field:

```go
type AnswerResponse struct {
    // ... existing fields ...
    Warning string `json:"warning,omitempty"`
}
```

Set `response.Warning = "Code execution service unavailable. Answer could not be verified."` when Judge0 fails.

---

## Phase 2: Security & Data Integrity (5 tasks)

### Task 2.1: Add Database Transactions to Multi-Step Operations

**Problem:** `SubmitAnswer`, `CreateTrainingPlan`, `AddProblemToList` do multiple DB writes without transactions. Partial failures leave inconsistent data.

**Files:**
- `backend/services/question_service.go` — `SubmitAnswer()`
- `backend/services/training_plan_service.go` — `CreateTrainingPlan()`
- `backend/services/list_service.go` — `AddProblemToList()`, `RemoveProblemFromList()`

**Package:** Built-in GORM transaction support (`db.Transaction()`)

**Spec for SubmitAnswer:**

```go
func (s *QuestionService) SubmitAnswer(userID int, req AnswerRequest) (*AnswerResponse, error) {
    var response *AnswerResponse

    err := s.db.Transaction(func(tx *gorm.DB) error {
        // All DB operations inside use `tx` instead of `s.db`
        // 1. Create attempt
        // 2. Update question stats
        // 3. Update user progress (proficiency)
        // 4. Record daily activity
        // 5. Update streak
        // 6. Process spaced repetition review
        // If any fails, all roll back
        return nil
    })

    return response, err
}
```

**Spec for CreateTrainingPlan:**

```go
func (s *TrainingPlanService) CreateTrainingPlan(userID int, req CreatePlanRequest) (*models.TrainingPlan, error) {
    var plan *models.TrainingPlan

    err := s.db.Transaction(func(tx *gorm.DB) error {
        // 1. Create plan
        // 2. Generate plan items
        // If item generation fails, plan is also rolled back
        return nil
    })

    return plan, err
}
```

---

### Task 2.2: Fix TopicPerformance Authorization

**Problem:** `topic_handlers.go:85` reads `userId` from URL params. Any user can view another user's data.

**File:** `backend/handlers/topic_handlers.go`

**Fix:**

```go
// Before (insecure):
userIDStr := c.Params("userId")
userID, err := strconv.Atoi(userIDStr)

// After (secure):
userID := c.Locals("user_id").(int)
```

**Also:** Move the route from public to protected in `routes.go`.

---

### Task 2.3: Add Rate Limiting Middleware

**Problem:** No rate limiting. Judge0 endpoint and auth endpoints can be hammered.

**New file:** `backend/middleware/rate_limiter.go`

**Package:** `github.com/gofiber/fiber/v2/middleware/limiter` (Fiber's built-in rate limiter — no extra dependency)

**Spec:**

```go
// Global rate limit: 100 requests/minute per IP
app.Use(limiter.New(limiter.Config{
    Max:        100,
    Expiration: 1 * time.Minute,
    KeyGenerator: func(c *fiber.Ctx) string {
        return c.IP()
    },
    LimitReached: func(c *fiber.Ctx) error {
        return c.Status(429).JSON(fiber.Map{
            "error": "Too many requests. Please try again later.",
        })
    },
}))

// Stricter limit for auth endpoints: 10/minute per IP
authGroup.Use(limiter.New(limiter.Config{
    Max:        10,
    Expiration: 1 * time.Minute,
}))

// Stricter limit for code execution: 20/minute per user
questionGroup.Use(limiter.New(limiter.Config{
    Max:        20,
    Expiration: 1 * time.Minute,
    KeyGenerator: func(c *fiber.Ctx) string {
        if uid, ok := c.Locals("user_id").(int); ok {
            return fmt.Sprintf("user:%d", uid)
        }
        return c.IP()
    },
}))
```

---

### Task 2.4: Fix Graceful Shutdown Order

**Problem:** `main.go:163-168` — DB is closed before Fiber is shut down. In-flight requests can hit a closed DB.

**File:** `backend/main.go`

**Fix:** Reverse the order — shutdown Fiber first (drain requests), then close DB:

```go
// Before:
sqlDB, _ := DB.DB()
sqlDB.Close()
app.ShutdownWithTimeout(5 * time.Second)

// After:
app.ShutdownWithTimeout(10 * time.Second) // Wait for in-flight requests
sqlDB, _ := DB.DB()
sqlDB.Close()
```

---

### Task 2.5: Fix CORS Multi-Origin Handling

**Problem:** `main.go:135` takes only `AllowOrigins[0]`. Multiple configured origins are ignored.

**File:** `backend/main.go`

**Fix:**

```go
// Before:
AllowOrigins: cfg.Server.CORS.AllowOrigins[0],

// After:
AllowOrigins: strings.Join(cfg.Server.CORS.AllowOrigins, ","),
```

**Add import:** `"strings"`

---

## Phase 3: Missing Features (7 tasks)

### Task 3.1: Password Reset Flow

**Problem:** No forgot-password / password reset. Users locked out if they forget their password.

**New files:**
- Additions to `backend/services/auth_service.go`
- Additions to `backend/handlers/auth_handlers.go`

**Package:** `github.com/golang-jwt/jwt/v5` (already in use for JWT — reuse for reset tokens)

**New endpoints:**
```
POST /api/auth/forgot-password    — Accepts { "email": "..." }
POST /api/auth/reset-password     — Accepts { "token": "...", "new_password": "..." }
```

**Spec:**

```go
// ForgotPassword — generates a time-limited reset token (1 hour expiry)
func (s *AuthService) ForgotPassword(email string) (string, error) {
    // 1. Find user by email (return success even if not found — prevent enumeration)
    // 2. Generate JWT with short expiry (1 hour), claim: user_id, purpose: "reset"
    // 3. Return token (in real app, send via email; for now return in response)
}

// ResetPassword — validates token and sets new password
func (s *AuthService) ResetPassword(token string, newPassword string) error {
    // 1. Validate token (check expiry, purpose claim)
    // 2. Extract user_id from token
    // 3. Hash new password with bcrypt
    // 4. Update user's password_hash
}
```

**New DB column:** `password_reset_at TIMESTAMP` on `users` table — to invalidate tokens issued before last reset.

---

### Task 3.2: Refresh Token Endpoint

**Problem:** Config has `RefreshExpiry: 7 days` but no endpoint exists. Users must re-login when JWT expires (24h).

**File:** `backend/services/auth_service.go`, `backend/handlers/auth_handlers.go`

**New endpoint:**
```
POST /api/auth/refresh    — Accepts { "refresh_token": "..." }
```

**Spec:**

```go
// Login now returns both access_token (24h) and refresh_token (7d)
type LoginResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int    `json:"expires_in"` // seconds
}

// RefreshToken — accepts a refresh token, returns new access token
func (s *AuthService) RefreshToken(refreshToken string) (*LoginResponse, error) {
    // 1. Validate refresh token (check type claim = "refresh", check expiry)
    // 2. Extract user_id
    // 3. Generate new access token (24h)
    // 4. Optionally rotate refresh token
    // 5. Return new tokens
}
```

**Implementation:** Reuse existing JWT infrastructure. Add a `token_type` claim ("access" vs "refresh") to distinguish token purposes.

---

### Task 3.3: Request Validation Middleware

**Problem:** Every handler does ad-hoc validation with inconsistent error formats.

**Package:** `github.com/go-playground/validator/v10` (most popular Go validation library, 15k+ GitHub stars)

**New file:** `backend/middleware/validator.go`

**Spec:**

```go
import "github.com/go-playground/validator/v10"

var validate = validator.New()

// ValidationError — standardized error response
type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

// ValidateStruct — validates any struct with `validate` tags
func ValidateStruct(s interface{}) []ValidationError {
    var errors []ValidationError
    err := validate.Struct(s)
    if err != nil {
        for _, e := range err.(validator.ValidationErrors) {
            errors = append(errors, ValidationError{
                Field:   e.Field(),
                Message: formatValidationMessage(e),
            })
        }
    }
    return errors
}
```

**Add validation tags to existing request structs:**

```go
type AnswerRequest struct {
    QuestionID int                    `json:"question_id" validate:"required,gt=0"`
    UserAnswer map[string]interface{} `json:"user_answer" validate:"required"`
    TimeTaken  int                    `json:"time_taken_seconds" validate:"required,gte=0"`
    HintsUsed  int                    `json:"hints_used" validate:"gte=0,lte=3"`
}

type CreatePlanRequest struct {
    Name            string `json:"name" validate:"required,min=1,max=200"`
    DurationDays    int    `json:"duration_days" validate:"required,gte=1,lte=365"`
    QuestionsPerDay int    `json:"questions_per_day" validate:"required,gte=1,lte=50"`
    DifficultyMin   float64 `json:"difficulty_min" validate:"gte=0,lte=100"`
    DifficultyMax   float64 `json:"difficulty_max" validate:"gte=0,lte=100,gtfield=DifficultyMin"`
}
```

**Standardized error response:**

```json
{
    "error": "Validation failed",
    "details": [
        {"field": "Name", "message": "Name is required"},
        {"field": "DurationDays", "message": "DurationDays must be at least 1"}
    ]
}
```

---

### Task 3.4: Structured Logging with Request IDs

**Problem:** No production logging. No request tracing.

**Package:** `github.com/rs/zerolog` (popular, fast structured logger for Go, 10k+ GitHub stars)

**New file:** `backend/middleware/logger.go`

**Spec:**

```go
import (
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    "github.com/google/uuid"
)

// RequestLogger middleware — adds request ID and structured logging
func RequestLogger() fiber.Handler {
    return func(c *fiber.Ctx) error {
        requestID := uuid.New().String()
        c.Locals("request_id", requestID)
        c.Set("X-Request-ID", requestID)

        start := time.Now()
        err := c.Next()
        duration := time.Since(start)

        log.Info().
            Str("request_id", requestID).
            Str("method", c.Method()).
            Str("path", c.Path()).
            Int("status", c.Response().StatusCode()).
            Dur("duration", duration).
            Str("ip", c.IP()).
            Msg("request")

        return err
    }
}
```

**Also:** Replace `fmt.Printf` / `log.Printf` calls throughout services with `zerolog` structured logging.

---

### Task 3.5: Health Check for Dependencies

**Problem:** `/health` just returns "healthy" without checking if PostgreSQL, ChromaDB, or Ollama are reachable.

**File:** `backend/main.go` (health endpoint)

**Spec:**

```go
app.Get("/health", func(c *fiber.Ctx) error {
    health := fiber.Map{
        "status":    "healthy",
        "timestamp": time.Now().UTC(),
    }

    // Check PostgreSQL
    sqlDB, err := DB.DB()
    if err != nil || sqlDB.Ping() != nil {
        health["status"] = "degraded"
        health["database"] = "unreachable"
    } else {
        health["database"] = "ok"
    }

    // Check ChromaDB (if configured)
    if cfg.ChromaDB.URL != "" {
        resp, err := http.Get(cfg.ChromaDB.URL + "/api/v1/heartbeat")
        if err != nil || resp.StatusCode != 200 {
            health["chromadb"] = "unreachable"
        } else {
            health["chromadb"] = "ok"
            resp.Body.Close()
        }
    }

    // Check Ollama (if configured)
    if cfg.Ollama.URL != "" {
        resp, err := http.Get(cfg.Ollama.URL + "/api/tags")
        if err != nil || resp.StatusCode != 200 {
            health["ollama"] = "unreachable"
        } else {
            health["ollama"] = "ok"
            resp.Body.Close()
        }
    }

    statusCode := 200
    if health["status"] == "degraded" {
        statusCode = 503
    }

    return c.Status(statusCode).JSON(health)
})
```

---

### Task 3.6: Pagination on All List Endpoints

**Problem:** `GetUserLists`, `GetUserAttempts`, `GetUserPlans`, `GetListProblems`, `GetPlanItems` return unbounded result sets.

**Files:**
- `backend/services/list_service.go`
- `backend/services/training_plan_service.go`
- `backend/services/question_service.go`
- `backend/handlers/` (corresponding handlers)

**Spec:** Add a shared pagination helper:

```go
// backend/utils/pagination.go
type PaginationParams struct {
    Page     int `json:"page"`      // 1-indexed
    PageSize int `json:"page_size"` // default 20, max 100
}

type PaginatedResponse struct {
    Data       interface{} `json:"data"`
    Total      int64       `json:"total"`
    Page       int         `json:"page"`
    PageSize   int         `json:"page_size"`
    TotalPages int         `json:"total_pages"`
}

func ParsePagination(c *fiber.Ctx) PaginationParams {
    page := c.QueryInt("page", 1)
    pageSize := c.QueryInt("page_size", 20)
    if page < 1 { page = 1 }
    if pageSize < 1 { pageSize = 20 }
    if pageSize > 100 { pageSize = 100 }
    return PaginationParams{Page: page, PageSize: pageSize}
}

func (p PaginationParams) Offset() int {
    return (p.Page - 1) * p.PageSize
}
```

Apply to all list endpoints — add `?page=1&page_size=20` query params.

---

### Task 3.7: Email Verification (Optional — can be deferred)

**Problem:** Users register with unverified emails.

**Package:** `github.com/golang-jwt/jwt/v5` (reuse for verification tokens)

**New DB column:** `email_verified_at TIMESTAMP` on `users`, `email_verification_token VARCHAR(500)`

**New endpoints:**
```
POST /api/auth/send-verification    — Resend verification email
POST /api/auth/verify-email         — Accepts { "token": "..." }
```

**Note:** For local development, print the token to console. For production, integrate with an email service (SendGrid, SES, etc.) — that integration is a separate task.

---

## Phase 4: Production Readiness (3 tasks)

### Task 4.1: Soft Deletes with Audit Trail

**Package:** Built-in GORM soft delete support (`gorm.io/gorm` — `gorm.Model` includes `DeletedAt`)

**Spec:** Add `DeletedAt gorm.DeletedAt` field to:
- `TrainingPlan`
- `TrainingPlanItem`
- `UserList`

GORM automatically handles soft deletes when `DeletedAt` field is present — `db.Delete()` sets `deleted_at` instead of removing the row.

---

### Task 4.2: Event System for Answer Submissions

**Package:** Simple Go channels + goroutines (no external dependency needed for v1)

**New file:** `backend/services/event_service.go`

**Spec:**

```go
type EventType string

const (
    EventAnswerSubmitted   EventType = "answer.submitted"
    EventStreakMilestone    EventType = "streak.milestone"
    EventPlanCompleted     EventType = "plan.completed"
    EventProficiencyChange EventType = "proficiency.changed"
)

type Event struct {
    Type    EventType
    UserID  int
    Payload map[string]interface{}
}

type EventService struct {
    subscribers map[EventType][]func(Event)
}

func (es *EventService) Subscribe(eventType EventType, handler func(Event))
func (es *EventService) Publish(event Event)
```

**Wire into:** `SubmitAnswer()` publishes `EventAnswerSubmitted`. Subscribers can handle gamification, notifications, analytics, etc.

---

### Task 4.3: Resolve Schema-Model Mismatches

**Problem:** SQL migration uses `VARCHAR(100)` for `topics.topic_id` but GORM model uses `int`.

**Decision:** Since the GORM model uses `int` and AutoMigrate is used in development, keep `int` as the source of truth. Update the SQL migration to match the model (use `SERIAL PRIMARY KEY` for `topic_id` instead of `VARCHAR`).

**Files:**
- `migrations/000001_core_schema.up.sql` — update topic_id type
- Verify all other field names match between migration and models

---

## Summary — Execution Order

| Order | Task | Risk | Time Est. |
|-------|------|------|-----------|
| 1 | **1.4** Fix topic conversion (`strconv.Itoa`) | Trivial | 5 min |
| 2 | **1.5** Fix code execution fallback | Low | 15 min |
| 3 | **2.5** Fix CORS multi-origin | Trivial | 5 min |
| 4 | **2.4** Fix graceful shutdown order | Low | 10 min |
| 5 | **2.2** Fix TopicPerformance authorization | Low | 15 min |
| 6 | **1.1** Wire SubmitAnswer to proficiency/activity/streak | Medium | 1-2 hrs |
| 7 | **1.3** Fix adaptive difficulty SQL | Medium | 1 hr |
| 8 | **1.2** Implement SM-2 spaced repetition service | Medium | 2-3 hrs |
| 9 | **2.1** Add DB transactions | Medium | 2 hrs |
| 10 | **2.3** Add rate limiting middleware | Low | 30 min |
| 11 | **3.5** Health check for dependencies | Low | 30 min |
| 12 | **3.6** Pagination on all list endpoints | Low | 1 hr |
| 13 | **3.3** Request validation framework | Medium | 2 hrs |
| 14 | **3.4** Structured logging | Medium | 2 hrs |
| 15 | **3.1** Password reset flow | Medium | 2-3 hrs |
| 16 | **3.2** Refresh token endpoint | Medium | 1-2 hrs |
| 17 | **4.1** Soft deletes | Low | 30 min |
| 18 | **4.2** Event system | Medium | 2 hrs |
| 19 | **4.3** Schema-model mismatches | Medium | 1-2 hrs |
| 20 | **3.7** Email verification | Low | 2 hrs |

**Total estimated:** ~20-25 hours of focused development

---

## External Packages to Add

| Package | Purpose | GitHub Stars |
|---------|---------|-------------|
| `github.com/go-playground/validator/v10` | Request validation | 16k+ |
| `github.com/rs/zerolog` | Structured logging | 10k+ |
| `github.com/google/uuid` | Request ID generation | 5k+ |
| `github.com/gofiber/fiber/v2/middleware/limiter` | Rate limiting (already in Fiber) | Built-in |

All other changes use existing dependencies (GORM transactions, Fiber middleware, `golang-jwt`).

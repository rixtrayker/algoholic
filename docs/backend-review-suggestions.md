# Comprehensive Backend Review & Suggestions

**Date:** February 24, 2026
**Scope:** All 31 Go backend files (~7,880 LOC), 2 SQL migrations, 15+ docs/standards files
**Cross-referenced:** mission.md, roadmap.md, architecture.md, gap-analysis.md, implementation-plan.md, question-design.md, skill-gyms.md, all agent-os standards

---

## A. Features Where Business Logic Is Broken (Code Exists But Doesn't Work)

### 1. Spaced Repetition Is Dead (CRITICAL)

The `spaced_repetition_reviews` table, `ReviewQueue` model, and `SpacedRepetitionReview` model all exist. `RecommendationService.GetSpacedRepetitionRecommendations()` queries `user_skills.next_review_at`. But **nothing ever writes to these columns**. No SM-2 calculation exists anywhere in working form.

The mission doc says: *"Spaced Repetition for Long-Term Retention: Scientific scheduling algorithm ensures concepts are reviewed at optimal intervals"* — but the implementation is completely absent.

`UserService.CalculateNextReviewDate()` exists but is a naive if/else (not SM-2), and more importantly, **it is never called** from `SubmitAnswer()`.

**Fix:** Build a `SpacedRepetitionService` implementing SM-2 (or FSRS). Wire it into `QuestionService.SubmitAnswer()`. Update `spaced_repetition_reviews` and sync `user_skills.next_review_at`.

---

### 2. User Skill Proficiency Never Updates (CRITICAL)

`UserSkill.ProficiencyLevel` starts at 0 and **stays at 0 forever**. `SubmitAnswer()` records attempts but never touches `user_skills`. Since the recommendation engine, weakness detection, personalized difficulty, and training plan adaptation ALL depend on `proficiency_level`, they all operate on garbage data.

**Fix:** After each answer, compute a proficiency update for the related topic. ELO-style or Bayesian update based on correctness + question difficulty.

---

### 3. Training Plan Adaptive Difficulty SQL Is Self-Referencing (HIGH)

`training_plan_service.go:306`: The WHERE clause `difficulty_score BETWEEN difficulty_score + 5 AND difficulty_score + 15` compares a column to itself. Always evaluates to false. The feature does nothing.

**Fix:** Reference the current plan item's question difficulty via a subquery or parameter, not a self-referencing column.

---

### 4. Training Plan Topic Conversion Produces Unicode Garbage (HIGH)

`training_plan_service.go:59`: `string(rune(topic))` converts integer topic ID to a Unicode character instead of its string representation (e.g., topic ID 65 becomes "A").

**Fix:** Use `strconv.Itoa(topic)`.

---

### 5. Streaks & Activity Never Update Through Normal Use (HIGH)

`SubmitAnswer()` doesn't call `UpdateStreak()`, `RecordDailyActivity()`, or `RecordActivity()`. The activity chart is blank unless the frontend manually POSTs to `/api/activity/record`. A user who answers 50 questions shows 0 on their activity chart.

**Fix:** Wire `UserService.RecordDailyActivity()` and `UserService.UpdateStreak()` into `QuestionService.SubmitAnswer()`.

---

### 6. Code Execution Fallback Silently Passes Garbage Code (MEDIUM)

When Judge0 is down, `ValidateCode()` checks for `def ` or `class ` in Python. Any syntactically plausible code is marked correct with full points.

**Fix:** Return an explicit "code execution service unavailable" error instead of silently accepting.

---

### 7. Recommendations Exclude Standalone Questions (MEDIUM)

All recommendation queries join through `problem_topics`, so questions with `problem_id = NULL` are invisible. The schema allows nullable `problem_id`, so standalone conceptual questions can never be recommended.

**Fix:** Add a fallback path that queries questions by their `related_concepts` or tags when `problem_id` is NULL.

---

### 8. Streak Logic Race Condition (MEDIUM)

`UpdateStreak()` and `RecordDailyActivity()` aren't called atomically. Streak checks for "yesterday" before today's record exists.

---

## B. Features Missing from ALL Roadmap Phases

Cross-referencing the roadmap (Phases 1–4), implementation plan, and gap analysis — these items are **not planned anywhere**:

### Security & Auth
1. **Password reset / forgot password** — Not mentioned in any roadmap phase
2. **Email verification** — Not planned
3. **Refresh token endpoint** — Config has `RefreshExpiry` but no endpoint; none planned
4. **Rate limiting** — Tech-stack mentions "Redis-based rate limiting" but not in any roadmap phase
5. **Admin role / authorization** — Only appears in Phase 4 "enterprise" with no specifics

### Data Integrity
6. **Database transactions** — `SubmitAnswer`, `CreateTrainingPlan`, `AddProblemToList` all do multi-step operations without transactions
7. **Soft deletes** — All deletes are hard deletes; no audit trail
8. **Concurrency protection on user lists** — JSON column read-modify-write race condition
9. **Request validation framework** — Every handler does ad-hoc validation with inconsistent error formats

### Observability
10. **Structured production logging** — `LoggingConfig` struct exists but is unused
11. **Request IDs / distributed tracing** — Nothing planned
12. **Health check for dependencies** — `/health` returns "healthy" without pinging DB/ChromaDB/Ollama
13. **Metrics endpoint** — No `/metrics` for Prometheus

### Architecture
14. **Event/webhook system** — Answer submissions, streak milestones, plan completions emit no events
15. **Bulk import/export** — No way to load problems from JSON/CSV or export user data

---

## C. Features Described in Docs But Completely Absent from Code

### 1. Anti-Memorization System (CORE MISSION, NOT IMPLEMENTED)

The **mission doc** says: *"Anti-Memorization System: Detects memorization patterns and adapts to ensure genuine mastery"*

The **architecture doc** (Section 6) describes:
- Time pattern analysis (flag suspiciously fast solves)
- Question rotation (never show same question within 30 days)
- Problem variants generation
- Explanation requirement scoring (60% correctness + 40% explanation quality)
- Understanding score (L1–L5 weighted formula)
- `shows_memorization` field on `UserAttempt`
- `memorization_score` field on `Assessment`

**None of this exists in the code.** The `shows_memorization` column is in the model but never written. No memorization detection logic exists.

---

### 2. Assessment System (ARCHITECTURE DESCRIBED, CODE HOLLOW)

The `Assessment` and `WeaknessAnalysis` GORM models exist. The architecture doc describes diagnostic, progress, and mock interview assessments. But there is:
- No `AssessmentService`
- No assessment creation endpoint
- No assessment scoring logic
- No LLM-based evaluation
- No endpoints to start/complete assessments

---

### 3. Skill Gyms (ENTIRE DOC, ZERO CODE)

`docs/skill-gyms.md` describes 8 detailed cognitive training gyms (Pattern Recognition Lab, Root Cause Arena, Assumption Breaker, Decision Speed Track, Decomposition Dojo, etc.) with gamification, rewards, and integration with the main platform. **Zero backend code exists for any of this.**

---

### 4. Multi-Level Weakness Detection (ARCHITECTURE DESCRIBED, NOT IMPLEMENTED)

The architecture doc describes 4 levels:
- Level 1: Statistical (accuracy < 60%)
- Level 2: Pattern-based (mistake type grouping)
- Level 3: Comparative (vs peers)
- Level 4: LLM deep analysis

The current `UserService.GetWeaknesses()` does a SQL query for topics with `proficiency_level < 50`. That's a fraction of Level 1. Levels 2–4 don't exist.

---

### 5. Question Generation / Variant Generation (PLANNED IN DOCS, NOT IN CODE)

The question design doc describes AI-generated variants, follow-ups, and constraint modifications. The architecture mentions a `QuestionGenerationService`. No code exists.

---

### 6. RAG Pipeline (ARCHITECTURE DESCRIBED, HALF-BUILT)

The architecture describes: `Query → Embed → Vector Search → Retrieve Context → Augment Prompt → LLM → Response`. The vector search part exists, but there's no prompt augmentation or LLM-powered response generation. The pipeline stops at returning raw search results.

---

## D. Architectural / Code Issues

### 1. Schema-Model Mismatch (HIGH)

The SQL migration `000001_core_schema.up.sql` defines `topics.topic_id` as `VARCHAR(100)` (string PK like `"dp-01"`), but the GORM model uses `TopicID int`. The standards doc (`struct-tags.md`) shows `TopicID string`. Migration has `statement` for problem body; model has `Description`. Migration has `difficulty_label` enum; model has `OfficialDifficulty` string. These can't both be correct.

### 2. Standards Docs Not Followed in Code

- `enum-types.md` says use `type DifficultyLabel string` with constants. Code uses plain strings.
- `struct-tags.md` says use dual `json` + `gorm` tags with `TableName()`. `Topic` model doesn't match (int vs string PK).
- `seeding-strategy.md` describes two-pass insertion with `OnConflict`. Seed code uses basic `Create()`.

### 3. TopicPerformance Authorization Hole (SECURITY)

`topic_handlers.go:85-86` reads `userId` from URL params instead of JWT. Any user can view any other user's topic performance. Route is public at `routes.go:89`.

### 4. CORS Multi-Origin Bug

`main.go:135` takes `cfg.Server.CORS.AllowOrigins[0]` — only the first origin. Should `strings.Join()` all origins.

### 5. Graceful Shutdown Order Wrong

`main.go:163-168` closes DB before shutting down Fiber. In-flight requests hit a closed database.

### 6. Pagination Missing on Most Endpoints

`GetUserLists`, `GetUserAttempts`, `GetUserPlans`, `GetListProblems`, `GetPlanItems` all return unbounded result sets.

---

## E. Prioritized Recommendations

### Tier 1: Core mission is broken — fix immediately

| # | Issue | Impact |
|---|-------|--------|
| 1 | Implement SM-2 spaced repetition service + wire to SubmitAnswer | Core platform mechanic is dead |
| 2 | Update `user_skills.proficiency_level` after every answer | Recommendations, weakness detection, personalization all broken |
| 3 | Fix adaptive difficulty self-referencing SQL | Training plan feature does nothing |
| 4 | Fix `string(rune(topic))` → `strconv.Itoa(topic)` | Training plan topics are garbage |
| 5 | Wire activity + streak recording into SubmitAnswer | Activity chart and streaks never update |

### Tier 2: Security & data integrity (not in any plan)

| # | Issue |
|---|-------|
| 6 | Wrap multi-step operations in DB transactions |
| 7 | Fix TopicPerformance authorization (use JWT user_id) |
| 8 | Add rate limiting middleware (especially for Judge0) |
| 9 | Fix graceful shutdown order (Fiber first, then DB) |
| 10 | Fix CORS multi-origin handling |

### Tier 3: Missing features not in any roadmap phase

| # | Issue |
|---|-------|
| 11 | Password reset / forgot password flow |
| 12 | Refresh token endpoint |
| 13 | Email verification |
| 14 | Request validation framework with consistent error format |
| 15 | Structured logging with request IDs |
| 16 | Health check that pings all dependencies |
| 17 | Pagination on all list endpoints |

### Tier 4: Core mission features described but not built

| # | Issue | Where described |
|---|-------|-----------------|
| 18 | Anti-memorization detection (time analysis, pattern flagging, question rotation) | mission.md, architecture.md §6 |
| 19 | Multi-level weakness detection (Levels 2–4) | architecture.md §7 |
| 20 | Assessment system (diagnostic, progress, mock interview) | architecture.md §6, roadmap Phase 3 |
| 21 | Full RAG pipeline (prompt augmentation, LLM-powered responses) | architecture.md §4 |
| 22 | Question variant generation | question-design.md, architecture.md §10 |

### Tier 5: Lower priority improvements

| # | Issue |
|---|-------|
| 23 | Resolve schema-model mismatches (TopicID int vs string, etc.) |
| 24 | Align code with standards docs (enum types, struct tags, seeding strategy) |
| 25 | Soft deletes with audit trail |
| 26 | Event/hook system for extensibility |
| 27 | Admin role and authorization |
| 28 | Return explicit error when Judge0 is unavailable |
| 29 | Skill Gyms (described in docs, zero code) |

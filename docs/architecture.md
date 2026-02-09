# Algoholic — System Architecture
## PostgreSQL + Apache AGE + Vector DB + RAG + LLM Assessment

---

## Table of Contents

1. [System Overview](#1-system-overview)
2. [Database Schema (PostgreSQL)](#2-database-schema)
3. [Graph Database (Apache AGE)](#3-graph-database)
4. [Vector Database & RAG](#4-vector-database--rag)
5. [Difficulty Scoring System](#5-difficulty-scoring-system)
6. [Assessment & Anti-Memorization](#6-assessment--anti-memorization)
7. [Weakness Detection](#7-weakness-detection)
8. [Search & Filtering](#8-search--filtering)
9. [Training Plan Builder](#9-training-plan-builder)
10. [LLM Integration](#10-llm-integration)
11. [API Design](#11-api-design)
12. [Local Deployment](#12-local-deployment)
13. [Implementation Phases](#13-implementation-phases)

---

## 1. System Overview

**✅ Current Status**: Phase 1 Complete - All Core Services Running
- Backend API: http://localhost:4000 (Go Fiber v2.52.11, 65 routes)
- Frontend: http://localhost:5173 (React 19 + Vite 7.3.1)
- Database: PostgreSQL (leetcode_training)
- API Tests: Postman Collection + Newman (22 endpoints)

### Technology Stack (Current Implementation)

```
┌─────────────────────────────────────────────────────┐
│     Frontend (React 19 + Tailwind CSS v4) ✅        │
│     Port: 5173 | Vite HMR | 46 Tests Passing        │
└────────────────────────┬────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────┐
│          Backend (Go + Fiber v2.52.11) ✅           │
│     Port: 4000 | GORM Auto-Migration | JWT Auth     │
└──┬──────────┬──────────┬──────────┬─────────────────┘
   │          │          │          │
┌──▼────────┐┌▼────────┐┌▼────────┐┌▼──────────────┐
│PostgreSQL ││ChromaDB/ ││Ollama/  ││Redis          │
│✅ Running ││⏭ Phase 2││⏭ Phase 2││⏭ Phase 2     │
│Port: 5432 ││(Vectors) ││(LLM)   ││(Cache)        │
└───────────┘└──────────┘└─────────┘└────────────────┘
                    │
              ┌─────▼──────┐
              │ RAG Pipeline│
              │  ⏭ Phase 2 │
              └────────────┘
```

### Core Capabilities

| Layer | Technology | Status | Purpose |
|-------|-----------|--------|---------|
| API | Go Fiber v2.52.11 | ✅ Running | 22 REST endpoints with JWT auth |
| Frontend | React 19 + Vite 7.3.1 | ✅ Running | Interactive UI with Tailwind CSS v4 |
| Relational | PostgreSQL 16 | ✅ Running | Problems, questions, users, progress, training plans |
| Testing | Newman 6.2.1 | ✅ Active | API testing with Postman collection |
| Graph | Apache AGE | ⏭ Phase 2 | Problem→Problem, Topic→Topic relationships, learning paths |
| Semantic | ChromaDB/Qdrant | ⏭ Phase 2 | Similarity search, RAG context, recommendations |
| AI | Ollama (local LLM) | ⏭ Phase 2 | Assessment, question generation, weakness analysis |

---

## 2. Database Schema

### Problems

```sql
CREATE TABLE problems (
    problem_id SERIAL PRIMARY KEY,
    leetcode_number INTEGER UNIQUE,
    title VARCHAR(500) NOT NULL,
    slug VARCHAR(200) UNIQUE NOT NULL,
    description TEXT NOT NULL,
    constraints TEXT[],
    examples JSONB NOT NULL,
    hints TEXT[],

    -- Difficulty (custom 0-100 scale, see section 5)
    difficulty_score FLOAT NOT NULL CHECK (difficulty_score >= 0 AND difficulty_score <= 100),
    official_difficulty VARCHAR(20),       -- Easy/Medium/Hard

    -- Categorization
    primary_pattern VARCHAR(100),
    secondary_patterns VARCHAR(100)[],
    source VARCHAR(50),                    -- 'leetcode', 'custom', 'generated'

    -- Solution metadata
    time_complexity VARCHAR(50),
    space_complexity VARCHAR(50),

    -- Stats
    total_attempts INTEGER DEFAULT 0,
    total_solves INTEGER DEFAULT 0,
    average_time_seconds FLOAT,
    acceptance_rate FLOAT,
    companies JSONB,
    tags JSONB,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_problems_difficulty ON problems(difficulty_score);
CREATE INDEX idx_problems_pattern ON problems(primary_pattern);
CREATE INDEX idx_problems_leetcode ON problems(leetcode_number) WHERE leetcode_number IS NOT NULL;
CREATE INDEX idx_problems_tags ON problems USING GIN(tags);

-- Full-text search
ALTER TABLE problems ADD COLUMN search_vector tsvector
GENERATED ALWAYS AS (
    setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
    setweight(to_tsvector('english', coalesce(description, '')), 'B')
) STORED;
CREATE INDEX idx_problems_search ON problems USING GIN(search_vector);
```

### Topics

```sql
CREATE TABLE topics (
    topic_id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    parent_topic_id INTEGER REFERENCES topics(topic_id),
    category VARCHAR(50),                  -- 'data_structure', 'algorithm', 'pattern', 'concept'
    difficulty_level INTEGER,              -- 1-5
    estimated_learning_hours NUMERIC(4,1),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE problem_topics (
    problem_id INTEGER REFERENCES problems(problem_id) ON DELETE CASCADE,
    topic_id INTEGER REFERENCES topics(topic_id) ON DELETE CASCADE,
    relevance_score FLOAT DEFAULT 1.0,
    is_primary BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (problem_id, topic_id)
);
```

### Questions

```sql
CREATE TABLE questions (
    question_id SERIAL PRIMARY KEY,
    problem_id INTEGER REFERENCES problems(problem_id),  -- NULL if standalone

    -- Type (see question-design.md for full taxonomy)
    question_type VARCHAR(50) NOT NULL,    -- 'complexity_analysis', 'ds_selection', 'pattern_recognition', etc.
    question_subtype VARCHAR(50),
    question_format VARCHAR(20) NOT NULL,  -- 'multiple_choice', 'code', 'text', 'ranking'

    -- Content
    question_text TEXT NOT NULL,
    question_data JSONB,                   -- Code snippets, context
    answer_options JSONB,                  -- [{id, text, is_correct}]
    correct_answer JSONB NOT NULL,

    -- Learning
    explanation TEXT NOT NULL,
    wrong_answer_explanations JSONB,
    related_concepts TEXT[],
    common_mistakes TEXT[],

    -- Metadata
    difficulty_score FLOAT NOT NULL,
    estimated_time_seconds INTEGER,

    -- Stats
    total_attempts INTEGER DEFAULT 0,
    correct_attempts INTEGER DEFAULT 0,
    average_time_seconds FLOAT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_questions_type ON questions(question_type, question_subtype);
CREATE INDEX idx_questions_problem ON questions(problem_id);
CREATE INDEX idx_questions_difficulty ON questions(difficulty_score);
```

### Code Templates

```sql
CREATE TABLE code_templates (
    template_id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(200) UNIQUE NOT NULL,
    category VARCHAR(100) NOT NULL,
    description TEXT,
    when_to_use TEXT NOT NULL,
    cpp_template TEXT,
    python_template TEXT,
    complexity_time VARCHAR(100),
    complexity_space VARCHAR(100),
    related_patterns VARCHAR(100)[],
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Users & Progress

```sql
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    preferences JSONB,
    current_streak_days INT DEFAULT 0,
    total_study_time_seconds BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_active_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_attempts (
    attempt_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    question_id INTEGER REFERENCES questions(question_id),
    problem_id INTEGER REFERENCES problems(problem_id),

    user_answer JSONB NOT NULL,
    is_correct BOOLEAN NOT NULL,
    time_taken_seconds INTEGER NOT NULL,
    attempt_number INTEGER DEFAULT 1,
    hints_used INTEGER DEFAULT 0,
    confidence_level INTEGER,              -- 1-5

    -- LLM analysis (populated async)
    detected_patterns TEXT[],
    mistakes_made TEXT[],
    shows_memorization BOOLEAN,

    -- Context
    training_plan_id INTEGER,
    session_id VARCHAR(50),

    attempted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_attempts_user ON user_attempts(user_id, attempted_at DESC);
CREATE INDEX idx_attempts_question ON user_attempts(question_id);

CREATE TABLE user_skills (
    user_id INTEGER REFERENCES users(user_id),
    topic_id INTEGER REFERENCES topics(topic_id),
    proficiency_level FLOAT DEFAULT 0,     -- 0-100
    questions_attempted INT DEFAULT 0,
    questions_correct INT DEFAULT 0,
    improvement_rate FLOAT,
    needs_review BOOLEAN DEFAULT false,
    last_practiced_at TIMESTAMP,
    next_review_at TIMESTAMP,              -- Spaced repetition
    PRIMARY KEY (user_id, topic_id)
);

CREATE INDEX idx_user_skills_weak ON user_skills(user_id, proficiency_level)
    WHERE proficiency_level < 50;
```

### Training Plans

```sql
CREATE TABLE training_plans (
    plan_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    name VARCHAR(200) NOT NULL,
    description TEXT,
    plan_type VARCHAR(50),                 -- 'preset', 'custom', 'ai_generated'
    difficulty_range NUMRANGE,
    target_topics INTEGER[],
    target_patterns VARCHAR(100)[],
    duration_days INTEGER,
    questions_per_day INTEGER DEFAULT 5,
    adaptive_difficulty BOOLEAN DEFAULT TRUE,
    progress_percentage FLOAT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'active',
    start_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE training_plan_items (
    item_id SERIAL PRIMARY KEY,
    plan_id INTEGER REFERENCES training_plans(plan_id) ON DELETE CASCADE,
    question_id INTEGER REFERENCES questions(question_id),
    problem_id INTEGER REFERENCES problems(problem_id),
    sequence_number INTEGER NOT NULL,
    day_number INTEGER,
    scheduled_for DATE,
    item_type VARCHAR(50) NOT NULL,        -- 'question', 'problem', 'review'
    is_completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP,
    UNIQUE(plan_id, sequence_number)
);
```

### Assessments & Weakness Tracking

```sql
CREATE TABLE assessments (
    assessment_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    assessment_type VARCHAR(50),           -- 'diagnostic', 'progress', 'mock_interview'
    topics_covered TEXT[],
    overall_score FLOAT,
    category_scores JSONB,
    strengths TEXT[],
    weaknesses TEXT[],
    recommendations TEXT,
    memorization_score FLOAT,              -- 0-1
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    time_taken_seconds INTEGER
);

CREATE TABLE weakness_analysis (
    analysis_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    weakness_type VARCHAR(100),
    specific_topic INTEGER REFERENCES topics(topic_id),
    severity VARCHAR(20),                  -- 'critical', 'major', 'minor'
    weakness_score FLOAT NOT NULL,
    evidence_question_ids INTEGER[],
    pattern_description TEXT,
    recommended_practice JSONB,
    detected_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    UNIQUE(user_id, weakness_type, specific_topic)
);
```

### LLM Generation Tracking

```sql
CREATE TABLE llm_generations (
    generation_id SERIAL PRIMARY KEY,
    generation_type VARCHAR(50),           -- 'question', 'explanation', 'assessment', 'problem_variant'
    prompt_template VARCHAR(100),
    context_documents JSONB,
    generated_content TEXT,
    model_name VARCHAR(50),
    tokens_used INTEGER,
    quality_score FLOAT,
    is_approved BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## 3. Graph Database

### Why Graph?

Problems and topics form natural graphs: similar problems, follow-ups, prerequisites, learning paths. Graph queries like "find all prerequisites recursively" or "shortest learning path from arrays to DP" are natural in Cypher.

### Setup

```sql
CREATE EXTENSION IF NOT EXISTS age;
LOAD 'age';
SET search_path = ag_catalog, "$user", public;
SELECT create_graph('problem_graph');
```

### Node Types

```cypher
(:Problem {id, title, difficulty_score, primary_pattern})
(:Topic   {id, name, category, difficulty_level})
(:Pattern {name, category})
(:Company {name})
```

### Relationship Types

```cypher
-- Problem relationships
(:Problem)-[:SIMILAR_TO {similarity_score: 0.85, reason: "..."}]->(:Problem)
(:Problem)-[:FOLLOW_UP_OF {difficulty_increase: 15}]->(:Problem)
(:Problem)-[:PREREQUISITE_FOR]->(:Problem)
(:Problem)-[:VARIATION_OF {type: "constraint_change"}]->(:Problem)

-- Topic relationships
(:Problem)-[:HAS_TOPIC {relevance: 0.9, is_primary: true}]->(:Topic)
(:Problem)-[:USES_PATTERN]->(:Pattern)
(:Topic)-[:SUBTOPIC_OF]->(:Topic)
(:Topic)-[:PREREQUISITE_FOR {strength: 0.8}]->(:Topic)
(:Topic)-[:RELATED_TO {relevance: 0.7}]->(:Topic)

-- Company
(:Problem)-[:ASKED_BY {frequency: 5}]->(:Company)
```

### Key Queries

```sql
-- Find similar problems within 2 hops
SELECT * FROM cypher('problem_graph', $$
    MATCH (p:Problem {id: 1})-[r:SIMILAR_TO|SAME_PATTERN_AS*1..2]-(similar:Problem)
    WHERE r.similarity_score > 0.7
    RETURN DISTINCT similar.id, similar.title, similar.difficulty_score
    ORDER BY similar.difficulty_score
$$) AS (id bigint, title text, difficulty numeric);

-- Learning path between topics
SELECT * FROM cypher('problem_graph', $$
    MATCH path = shortestPath(
        (start:Topic {name: 'Array'})-[:PREREQUISITE_FOR*]-(end:Topic {name: 'Dynamic Programming'})
    )
    RETURN [node IN nodes(path) | node.name] AS learning_path
$$) AS (learning_path text[]);

-- Recommend problems based on mastered topics
SELECT * FROM cypher('problem_graph', $$
    MATCH (u:User {id: 123})-[:MASTERED]->(topic:Topic)
    MATCH (topic)<-[:HAS_TOPIC]-(problem:Problem)
    WHERE NOT (u)-[:SOLVED]->(problem)
    RETURN problem.id, problem.title
    ORDER BY problem.difficulty_score
    LIMIT 10
$$) AS (problem_id bigint, title text);

-- Find prerequisite gaps for weak topics
SELECT * FROM cypher('problem_graph', $$
    MATCH (u:User {id: 123})-[:STRUGGLING_WITH]->(weak:Topic)
    MATCH (weak)-[:PREREQUISITE_FOR]->(prereq:Topic)
    WHERE NOT (u)-[:MASTERED]->(prereq)
    RETURN DISTINCT prereq.name, prereq.difficulty_level
    ORDER BY prereq.difficulty_level
$$) AS (topic_name text, difficulty int);
```

---

## 4. Vector Database & RAG

### Collections

| Collection | Content | Use Case |
|-----------|---------|----------|
| `problems` | Problem descriptions + metadata | Semantic search, similar problem finding |
| `questions` | Question text + concepts | Question recommendation |
| `solutions` | Solution explanations + code patterns | Pattern matching, RAG context |
| `templates` | Code template descriptions | Template recommendation |

### Embedding Model

Local: `all-MiniLM-L6-v2` (384 dims, fast) or `bge-large-en-v1.5` (1024 dims, better quality).

### RAG Pipeline

```
User Query → Embed → Vector Search (top-k) → Retrieve from PostgreSQL → Augment Prompt → LLM → Response
```

```go
func RAGQuery(userQuery string, collection string, k int) (string, error) {
    // Query vector DB
    results, err := vectorDB.Query(userQuery, collection, k)
    if err != nil {
        return "", err
    }

    // Build context
    var contextBuilder strings.Builder
    for _, doc := range results.Documents {
        contextBuilder.WriteString(doc)
        contextBuilder.WriteString("\n\n")
    }

    // Generate prompt
    prompt := fmt.Sprintf("Context:\n%s\n\nQuestion: %s\n\nProvide a clear answer based on the context.",
        contextBuilder.String(), userQuery)

    return llm.Generate(prompt)
}
```

### Sync Strategy

On problem/question create/update → generate embedding → upsert to vector DB + store reference in PostgreSQL.

---

## 5. Difficulty Scoring System

### The Magic Unit: 0-100 Scale

Traditional Easy/Medium/Hard is too coarse. We use a multi-dimensional 0-100 score.

### Components

| Component | Weight | What It Measures |
|-----------|--------|-----------------|
| Conceptual Complexity | 25% | Number and depth of concepts needed |
| Algorithm Complexity | 20% | Complexity of the required algorithm |
| Implementation Difficulty | 15% | How hard to code correctly |
| Pattern Recognition | 20% | How obvious/hidden the pattern is |
| Edge Case Density | 10% | Number of edge cases to handle |
| Time Pressure | 10% | Expected solve time |

### Calculation

```go
func CalculateDifficultyScore(problem *Problem) float64 {
    scores := map[string]float64{
        "conceptual":     ScoreConceptual(problem),      // 0-100
        "algorithm":      ScoreAlgorithm(problem),       // 0-100
        "implementation": ScoreImplementation(problem),  // 0-100
        "pattern":        ScorePatternRecognition(problem), // 0-100
        "edge_cases":     ScoreEdgeCases(problem),       // 0-100
        "time_pressure":  ScoreTimePressure(problem),    // 0-100
    }

    weights := []float64{0.25, 0.20, 0.15, 0.20, 0.10, 0.10}
    scoreValues := []float64{
        scores["conceptual"],
        scores["algorithm"],
        scores["implementation"],
        scores["pattern"],
        scores["edge_cases"],
        scores["time_pressure"],
    }

    total := 0.0
    for i, score := range scoreValues {
        total += score * weights[i]
    }

    return total
}
```

### Tiers

```
 0-20  Trivial       Basic application, no tricks
21-35  Easy          Single concept, straightforward
36-50  Medium-Easy   1-2 concepts, some thinking
51-65  Medium        Multiple concepts, optimization needed
66-80  Hard          Advanced algorithms, many edge cases
81-95  Very Hard     Research-level, creative solutions
96-100 Expert        Competition-level, novel
```

### Dynamic Calibration

Difficulty adjusts based on actual user performance (success rate, average time). Capped at ±10% change per recalibration to prevent wild swings.

### Personalized Difficulty

```go
func PersonalizedDifficulty(problemID, userID int) float64 {
    base := GetDifficulty(problemID)
    userProficiency := GetAvgProficiency(userID, GetProblemTopics(problemID))
    adjustment := (50 - userProficiency) * 0.6 // Strong user → easier, weak → harder
    return clamp(base + adjustment, 0, 100)
}

func clamp(value, min, max float64) float64 {
    if value < min {
        return min
    }
    if value > max {
        return max
    }
    return value
}
```

---

## 6. Assessment & Anti-Memorization

### Assessment Types

| Type | Purpose | Duration | Frequency |
|------|---------|----------|-----------|
| Diagnostic | Determine starting level | 45-60 min | Once |
| Progress | Check improvement | 20-30 min | Weekly |
| Mock Interview | Simulate real interview | 45 min/problem | As needed |

### LLM Assessment Prompt

For each submission, the LLM evaluates:
1. Correctness (0-100)
2. Approach quality (0-100)
3. Understanding depth (0-100)
4. Memorization likelihood (0-1)
5. Specific strengths/weaknesses
6. Personalized recommendations

### Anti-Memorization Techniques

**1. Problem Variants** — Generate modified versions (change DS, constraints, direction).

**2. Explanation Requirement** — Score = 60% correctness + 40% explanation quality (LLM-evaluated).

**3. Transfer Questions** — After solving "find duplicates in array", ask "find duplicates in stream".

**4. Time Pattern Analysis** — Flag suspiciously fast solves (< 30% expected time) + no incorrect attempts on hard problems.

**5. Question Rotation** — Never show same question twice within 30 days. If all seen, use least-recently-seen.

**6. Constraint Modification** — "If n changed from 10^5 to 10^9, how would you adapt?"

### Understanding Score

```
Understanding = L1 * 0.15 (correct implementation)
              + L2 * 0.25 (can explain why)
              + L3 * 0.30 (can solve variant)
              + L4 * 0.20 (can adapt to changed constraints)
              + L5 * 0.10 (can compare multiple approaches)
```

---

## 7. Weakness Detection

### Multi-Level Detection

**Level 1 — Statistical:** Accuracy < 60% in a category → weakness.

**Level 2 — Pattern:** Group incorrect attempts by mistake type (off-by-one, wrong DS, missed edge case). Frequency > 30% → persistent weakness.

**Level 3 — Comparative:** Compare user to peers at same level. > 15 points below peer average → weakness.

**Level 4 — LLM Deep Analysis:** Send last 20 incorrect attempts to LLM for root cause analysis (conceptual gap vs careless error vs knowledge gap).

### Severity

```go
func CalculateSeverity(frequency, impact, difficultyToFix float64) string {
    severity := frequency*0.4 + impact*0.4 + difficultyToFix*0.2

    if severity > 0.7 {
        return "critical"
    } else if severity > 0.5 {
        return "major"
    }
    return "minor"
}
```

### Remediation

For each weakness: find prerequisite topics via graph → select practice questions at user's level → schedule in training plan.

---

## 8. Search & Filtering

### Multi-Modal Search

```
User Query → [Keyword (PostgreSQL FTS)] + [Semantic (Vector DB)] + [Graph (AGE)] → Merge & Rank → Results
```

### Filters

- Difficulty range, patterns, topics, companies, tags
- Status: unsolved, attempted, mastered
- Graph: similar_to, follow_up_of, prerequisite_for
- User-specific: weak areas only, due for review

### Ranking

Results from multiple sources get combined scores. Appearing in multiple search modes boosts rank.

---

## 9. Training Plan Builder

### Plan Types

| Type | Description |
|------|-------------|
| Preset | Pre-built curricula (e.g., "DP Bootcamp 30 days") |
| Custom | User picks topics, duration, daily load |
| AI-Generated | LLM creates plan based on user profile + weaknesses |

### Generation Algorithm

1. Assess current level per topic
2. Calculate skill gap to target
3. Prioritize: weaknesses > required topics > prerequisites > frequency
4. Select questions with progressive difficulty curve
5. Inject spaced repetition reviews every 5 questions
6. Schedule across available days

### Adaptive Adjustment

- Accuracy > 85% → increase difficulty by 5 points
- Accuracy < 40% → decrease difficulty by 5 points
- Stuck on category (> 5 attempts, < 30% accuracy) → add prerequisite review
- Topic mastered → remove remaining questions, reallocate time

### Spaced Repetition

SM-2 algorithm: tracks ease factor and interval per question per user. Overdue reviews get priority in daily schedule.

---

## 10. LLM Integration

### Local Models

| Runtime | Models | Best For |
|---------|--------|----------|
| Ollama | Mistral 7B, CodeLlama 13B | Easy setup, good quality |
| llama.cpp | GGUF quantized models | CPU-friendly, fast |

### Use Cases

1. **Question Generation** — Create variants, follow-ups, new problems
2. **Assessment** — Evaluate solutions and explanations
3. **Hint Generation** — Progressive hints (Socratic → direct)
4. **Weakness Analysis** — Deep analysis of mistake patterns
5. **Explanation** — Generate solution explanations with RAG context

### Prompt Architecture

```
Task Router → RAG Context Retriever → Prompt Constructor → LLM Engine → Response Parser
```

All prompts use low temperature (0.3) for assessment consistency, higher (0.7-0.8) for generation creativity.

---

## 11. API Design

**✅ Current Implementation**: 22 endpoints across 7 categories

See the **[Postman Collection](../postman/algoholic-api.postman_collection.json)** for complete specifications with tests.

**Quick Reference**:

```
Authentication (3 endpoints)
  POST   /api/auth/register               Register new user
  POST   /api/auth/login                  Login
  GET    /api/auth/me 🔒                  Get current user

Questions (4 endpoints)
  GET    /api/questions/random            Get random question (with filters)
  POST   /api/questions/{id}/answer 🔒    Submit answer
  GET    /api/questions/{id}/hint 🔒      Request hint

Problems (4 endpoints)
  GET    /api/problems                    Search/filter problems
  GET    /api/problems/{id}               Get problem details
  GET    /api/problems/{id}/questions     Get problem's questions
  GET    /api/problems/search             Search problems

User Stats (3 endpoints)
  GET    /api/users/stats 🔒              User statistics
  GET    /api/users/progress 🔒           Progress history
  GET    /api/users/attempts 🔒           Attempt history

Training Plans (5 endpoints)
  GET    /api/plans                       List all plans
  GET    /api/plans/{id}                  Get plan details
  POST   /api/plans/{id}/enroll 🔒        Enroll in plan
  GET    /api/users/plans 🔒              Get user's enrolled plans
  PUT    /api/plans/{id}/progress 🔒      Update plan progress

Topics (2 endpoints)
  GET    /api/topics                      List all topics
  GET    /api/users/topics/{id}/performance 🔒  Get topic performance

Health Check (1 endpoint)
  GET    /health                          Health check
```

🔒 = Requires JWT authentication

**Testing**: Run `cd postman && ./run-tests.sh` to test all endpoints

**Documentation**: See [docs/api-reference.md](./api-reference.md) and [postman/README.md](../postman/README.md)

---

## 12. Local Deployment

### Docker Compose

```yaml
version: '3.8'

services:
  postgres:
    image: apache/age:PG16_latest
    environment:
      POSTGRES_USER: leetcode
      POSTGRES_PASSWORD: leetcode123
      POSTGRES_DB: leetcode_training
    ports: ["5432:5432"]
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./scripts/init_db.sql:/docker-entrypoint-initdb.d/init.sql
    command: postgres -c shared_preload_libraries=age

  chromadb:
    image: chromadb/chroma:latest
    ports: ["8000:8000"]
    volumes: [chroma_data:/chroma/chroma]
    environment:
      IS_PERSISTENT: "TRUE"
      ANONYMIZED_TELEMETRY: "FALSE"

  ollama:
    image: ollama/ollama:latest
    ports: ["11434:11434"]
    volumes: [ollama_data:/root/.ollama]

  backend:
    build: ./backend
    ports: ["4000:4000"]
    environment:
      DATABASE_URL: postgresql://leetcode:leetcode123@postgres:5432/leetcode_training
      CHROMA_URL: http://chromadb:8000
      OLLAMA_URL: http://ollama:11434
    depends_on: [postgres, chromadb, ollama]

  frontend:
    build: ./frontend
    ports: ["3000:3000"]
    depends_on: [backend]

volumes:
  postgres_data:
  chroma_data:
  ollama_data:
```

### Project Structure

```
backend/
├── main.go
├── config/          # Configuration
├── database/        # DB connection and setup
├── models/          # GORM models
├── services/        # Business logic (search, assessment, training, llm, vector)
├── handlers/        # Fiber route handlers
├── middleware/      # Auth, logging, etc.
├── utils/           # Difficulty calc, graph helpers, embeddings
└── migrations/      # golang-migrate SQL files

frontend/
├── src/
│   ├── pages/       # Dashboard, Practice, Assessment, TrainingPlan, Analytics
│   ├── components/  # QuestionCard, CodeEditor, ProgressChart
│   └── services/    # API client
```

### Startup (Current Method)

**✅ Services are already running**:

```bash
# Backend (Terminal 1)
cd backend
go run main.go
# → http://localhost:4000

# Frontend (Terminal 2)
cd frontend
npm run dev
# → http://localhost:5173

# Database is running as system PostgreSQL service
# postgresql://leetcode:leetcode_password@localhost:5432/leetcode_training
```

**Health Check**:
```bash
curl http://localhost:4000/health
```

**Run API Tests**:
```bash
cd postman
./run-tests.sh
```

**See [RUNNING.md](../RUNNING.md) for complete service management guide.**

---

## 13. Implementation Phases

| Phase | Weeks | Status | Deliverable |
|-------|-------|--------|-------------|
| 1. Foundation | 1-2 | ✅ **COMPLETE** | PostgreSQL schema, 22 API endpoints, JWT auth, Frontend with 46 tests, Postman collection |
| 2. Intelligence | 3-4 | ⏭ **NEXT** | Vector DB + embeddings, graph relationships, semantic search |
| 3. Training | 5-6 | 📋 Planned | Training plans, progress tracking, spaced repetition, weakness detection |
| 4. AI | 7-8 | 📋 Planned | Ollama integration, RAG, assessment analysis, question generation |
| 5. Frontend | 9-10 | 🔄 **IN PROGRESS** | Enhanced UI components, practice interface, dashboard improvements |
| 6. Polish | 11-12 | 📋 Planned | Performance optimization, difficulty calibration, user testing |

**Current Focus**: Phase 1 complete. Ready to begin Phase 2 (Intelligence layer) or continue enhancing Phase 5 (Frontend features).

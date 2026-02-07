# Algoholic — Getting Started
## Setup Guide & Implementation Checklist

---

## Prerequisites

```bash
docker --version          # Docker 20.10+
docker-compose --version  # Docker Compose 2.0+
go version                # Go 1.21+
node --version            # Node.js 18+
```

---

## Quick Start (30 minutes)

### 1. Project Structure

```bash
mkdir -p algoholic/{backend/{api,models,services,utils},frontend/src,scripts,data}
cd algoholic
git init
```

### 2. Docker Compose

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  postgres:
    image: apache/age:PG16_latest
    container_name: algoholic_postgres
    environment:
      POSTGRES_USER: leetcode
      POSTGRES_PASSWORD: leetcode123
      POSTGRES_DB: leetcode_training
    ports: ["5432:5432"]
    volumes:
      - ./data/postgres:/var/lib/postgresql/data
      - ./scripts/init_db.sql:/docker-entrypoint-initdb.d/init.sql
    command: postgres -c shared_preload_libraries=age
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U leetcode"]
      interval: 10s
      timeout: 5s
      retries: 5

  chromadb:
    image: chromadb/chroma:latest
    container_name: algoholic_chromadb
    ports: ["8000:8000"]
    volumes: [./data/chroma:/chroma/chroma]
    environment:
      IS_PERSISTENT: "TRUE"
      ANONYMIZED_TELEMETRY: "FALSE"

  ollama:
    image: ollama/ollama:latest
    container_name: algoholic_ollama
    ports: ["11434:11434"]
    volumes: [./data/ollama:/root/.ollama]

  redis:
    image: redis:7-alpine
    container_name: algoholic_redis
    ports: ["6379:6379"]
    volumes: [./data/redis:/data]
```

### 3. Database Init

Create `scripts/init_db.sql`:

```sql
CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS age;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
LOAD 'age';
SET search_path = ag_catalog, "$user", public;
SELECT create_graph('problem_graph');

-- Core tables (see architecture.md for full schema)
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE problems (
    problem_id SERIAL PRIMARY KEY,
    leetcode_number INTEGER UNIQUE,
    title VARCHAR(500) NOT NULL,
    slug VARCHAR(200) UNIQUE NOT NULL,
    description TEXT NOT NULL,
    constraints TEXT[],
    examples JSONB NOT NULL,
    hints TEXT[],
    difficulty_score FLOAT NOT NULL CHECK (difficulty_score >= 0 AND difficulty_score <= 100),
    primary_pattern VARCHAR(100),
    secondary_patterns VARCHAR(100)[],
    tags JSONB,
    total_attempts INTEGER DEFAULT 0,
    total_solves INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE questions (
    question_id SERIAL PRIMARY KEY,
    problem_id INTEGER REFERENCES problems(problem_id),
    question_type VARCHAR(50) NOT NULL,
    question_text TEXT NOT NULL,
    correct_answer JSONB NOT NULL,
    answer_options JSONB,
    explanation TEXT NOT NULL,
    difficulty_score FLOAT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_attempts (
    attempt_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    question_id INTEGER REFERENCES questions(question_id),
    user_answer JSONB NOT NULL,
    is_correct BOOLEAN NOT NULL,
    time_taken_seconds INTEGER NOT NULL,
    attempted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE training_plans (
    plan_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    name VARCHAR(200) NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    start_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_problems_difficulty ON problems(difficulty_score);
CREATE INDEX idx_questions_type ON questions(question_type);
CREATE INDEX idx_attempts_user ON user_attempts(user_id, attempted_at DESC);

-- Full-text search
ALTER TABLE problems ADD COLUMN search_vector tsvector
GENERATED ALWAYS AS (
    setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
    setweight(to_tsvector('english', coalesce(description, '')), 'B')
) STORED;
CREATE INDEX idx_problems_search ON problems USING GIN(search_vector);

-- Sample data
INSERT INTO users (username, email) VALUES ('demo_user', 'demo@example.com');

INSERT INTO problems (leetcode_number, title, slug, description, constraints, examples, hints, difficulty_score, primary_pattern)
VALUES (1, 'Two Sum', 'two-sum',
    'Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.',
    ARRAY['2 <= nums.length <= 10^4', '-10^9 <= nums[i] <= 10^9', 'Only one valid answer exists'],
    '[{"input": "nums = [2,7,11,15], target = 9", "output": "[0,1]", "explanation": "nums[0] + nums[1] == 9"}]'::jsonb,
    ARRAY['Try using a hash table to store complements'],
    15.5, 'Hash Table');
```

### 4. Backend Setup

Initialize backend:

```bash
cd backend
go mod init github.com/yourusername/algoholic
go get github.com/gofiber/fiber/v2
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/knadh/koanf/v2
go get github.com/knadh/koanf/parsers/yaml
go get github.com/knadh/koanf/providers/env
go get github.com/knadh/koanf/providers/file
go get github.com/knadh/koanf/providers/structs
```

### 5. Configuration

The backend uses [koanf](https://github.com/knadh/koanf) for configuration management with the following priority:

1. **Default values** (lowest priority)
2. **config.yaml file** (medium priority)
3. **Environment variables** (highest priority)

Copy the example files:

```bash
cp .env.example .env
# Edit .env with your values
```

All environment variables must be prefixed with `ALGOHOLIC_`. For example:

```bash
ALGOHOLIC_SERVER_PORT=5000
ALGOHOLIC_DATABASE_HOST=db.example.com
ALGOHOLIC_OLLAMA_URL=http://ollama:11434
```

Configuration file structure (`config.yaml`):

```yaml
app:
  name: "Algoholic API"
  environment: "development"  # development, staging, production
  debug: true

server:
  port: 4000

database:
  host: "localhost"
  port: 5432
  database: "leetcode_training"

# See config.yaml for full configuration options
```

The backend code is structured with koanf configuration. See `backend/main.go` and `backend/config/config.go` for the full implementation.

### 6. Database Migrations

Install golang-migrate:

```bash
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

Create `migrations/000001_init_schema.up.sql`:

```sql
-- Core tables (see docs/architecture.md for full schema)
CREATE TABLE problems (
    problem_id SERIAL PRIMARY KEY,
    leetcode_number INTEGER UNIQUE,
    title VARCHAR(500) NOT NULL,
    description TEXT NOT NULL,
    difficulty_score FLOAT NOT NULL CHECK (difficulty_score >= 0 AND difficulty_score <= 100),
    primary_pattern VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes
CREATE INDEX idx_problems_difficulty ON problems(difficulty_score);
```

### 7. Start Everything

```bash
# Start infrastructure
docker-compose up -d
sleep 30

# Download LLM models
docker exec -it algoholic_ollama ollama pull mistral:7b
docker exec -it algoholic_ollama ollama pull codellama:13b

# Run migrations
migrate -path migrations -database "postgresql://leetcode:leetcode123@localhost:5432/leetcode_training?sslmode=disable" up

# Start backend
cd backend
go run main.go
```

### 8. Verify

```bash
curl http://localhost:4000/health                    # API health
curl http://localhost:4000/api/problems | jq         # List problems
curl http://localhost:8000/api/v1/heartbeat          # ChromaDB
docker exec algoholic_ollama ollama list             # LLM models
docker exec algoholic_postgres psql -U leetcode -d leetcode_training -c "SELECT count(*) FROM problems;"
```

**Access points:**
- API: http://localhost:4000
- ChromaDB: http://localhost:8000
- Ollama: http://localhost:11434

---

## Implementation Checklist

### Phase 1: Foundation (Weeks 1-2)

- [ ] Docker environment running (PostgreSQL, ChromaDB, Ollama, Redis)
- [ ] Full database schema created (see architecture.md)
- [ ] Basic CRUD API endpoints for problems, questions, users
- [ ] Import 50 essential LeetCode problems (10 easy, 25 medium, 15 hard)
- [ ] Generate 200 initial questions across all 10 categories
- [ ] Implement difficulty scoring algorithm
- [ ] Full-text search working

### Phase 2: Intelligence (Weeks 3-4)

- [ ] ChromaDB collections created (problems, solutions, explanations, templates)
- [ ] Embedding generation pipeline (sentence-transformers)
- [ ] Semantic search endpoint
- [ ] Apache AGE graph created with topic/problem nodes
- [ ] Graph relationships populated (similar, follow-up, prerequisite)
- [ ] Graph query helpers (find similar, find prerequisites, learning path)
- [ ] RAG query function working with Ollama

### Phase 3: Training (Weeks 5-6)

- [ ] Training plan generation algorithm
- [ ] Preset plan templates (DP Bootcamp, Interview Prep, etc.)
- [ ] User progress tracking (attempts, accuracy, time)
- [ ] Spaced repetition scheduling (SM-2 algorithm)
- [ ] Adaptive difficulty adjustment
- [ ] Weakness detection (statistical + pattern-based)
- [ ] Review queue endpoint

### Phase 4: AI (Weeks 7-8)

- [ ] LLM assessment prompts (evaluate solutions + explanations)
- [ ] Memorization detection logic
- [ ] Question variant generation
- [ ] Progressive hint system (3 levels)
- [ ] Weakness deep analysis (LLM-powered)
- [ ] Problem generation pipeline with quality validation

### Phase 5: Frontend (Weeks 9-10)

- [ ] React app with routing
- [ ] Problem browser with search/filter
- [ ] Question practice interface
- [ ] Training plan UI (create, view progress, next question)
- [ ] Progress dashboard with charts
- [ ] Weakness report view

### Phase 6: Polish (Weeks 11-12)

- [ ] Load testing and query optimization
- [ ] Difficulty score recalibration from real data
- [ ] Caching for frequent queries
- [ ] Documentation
- [ ] User testing and feedback

---

## Useful Commands

```bash
# Start/stop
docker-compose up -d
docker-compose down

# Logs
docker-compose logs -f backend
docker-compose logs -f postgres

# Database
docker exec -it algoholic_postgres psql -U leetcode -d leetcode_training

# Backup/restore
docker exec algoholic_postgres pg_dump -U leetcode leetcode_training > backup.sql
docker exec -i algoholic_postgres psql -U leetcode leetcode_training < backup.sql

# Clean reset (⚠️ deletes data)
docker-compose down -v
rm -rf data/*

# LLM
docker exec -it algoholic_ollama ollama pull mistral:7b
docker exec -it algoholic_ollama ollama list
```

---

## Troubleshooting

| Issue | Solution |
|-------|----------|
| PostgreSQL won't start | `rm -rf data/postgres/*` then restart |
| Port already in use | `lsof -i :5432` to find process, kill or change port |
| Ollama model download fails | Check disk space (`df -h`), retry |
| Can't connect to DB | `docker exec algoholic_postgres pg_isready -U leetcode` |
| ChromaDB not responding | Check `docker-compose logs chromadb` |

---

## Database Seeding Guide
### n8n Workflow Implementation for PostgreSQL

This section provides detailed implementation instructions for the n8n workflow that seeds the PostgreSQL database using the Master Topic Reference.

### Seeding Sequence

```
1. Parse Topics → Insert Topics (sorted by depth)
2. Extract Prerequisites → Insert Prerequisites
3. Extract Patterns → Insert Patterns (deduplicated)
4. Link Topics-Patterns → Insert topic_patterns
5. Extract Mistakes → Insert common_mistakes
6. Extract Edge Cases → Insert edge_cases
7. Extract Problems → Insert problems
8. Link Problem-Topics → Insert problem_topics
9. Link Problem-Patterns → Insert problem_patterns
10. Validate Integrity
11. Trigger Graph Construction Workflow
```

### Key n8n Node Patterns

**Pattern 1: Parse and Transform**
```javascript
// Read markdown → Parse structure → Transform to DB format
const text = $input.first().binary.data.toString('utf8');
const parsed = parseStructure(text);
return parsed.map(item => ({json: item}));
```

**Pattern 2: Batch Insert with Conflict Handling**
```sql
INSERT INTO table (cols...) VALUES ($1, $2, ...)
ON CONFLICT (unique_col) DO UPDATE SET
    col1 = EXCLUDED.col1,
    updated_at = NOW()
RETURNING id;
```

**Pattern 3: Foreign Key Lookup**
```sql
-- Use subquery for FK lookup
INSERT INTO child_table (parent_id, data) VALUES (
    (SELECT id FROM parent_table WHERE code = $1),
    $2
);
```

**Pattern 4: Validation Query**
```sql
-- Post-insert validation
SELECT COUNT(*) as issues
FROM topics
WHERE parent_topic_id IS NULL AND depth_level > 1;
-- Should return 0
```

### Critical Implementation Notes

**1. Insertion Order Matters**
- Must insert parents before children (sort by depth_level)
- Must insert referenced entities before relationships

**2. Deduplication Strategy**
- Patterns: deduplicate by pattern_name
- Mistakes: deduplicate by topic_code + mistake_name
- Edge Cases: deduplicate by topic_code + case_name
- Problems: deduplicate by leetcode_id

**3. Data Enrichment**
- Auto-generate estimated_practice_hours if missing
- Infer pattern_type from pattern_name
- Infer mistake_category from mistake_name
- Infer edge_case_category from case_name

**4. Error Handling**
- Use ON CONFLICT for idempotency
- Retry on transient DB errors
- Log all insertions with counts
- Validate after each major step

### Validation Checklist

After seeding, verify:

✅ Topics: All inserted, proper hierarchy, no orphans
✅ Prerequisites: No circular dependencies, valid references
✅ Patterns: Deduplicated, all have type
✅ Topic-Patterns: All topics have ≥1 pattern (depth ≥ 2)
✅ Mistakes: Proper categories, all have severity
✅ Edge Cases: All have why_important field
✅ Problems: All have leetcode_id, proper difficulty
✅ Relationships: All FKs valid, no dangling references

### Expected Record Counts

| Table | Expected Count |
|-------|----------------|
| topics | ~150-200 |
| patterns | ~40-60 |
| prerequisites | ~100-150 |
| topic_patterns | ~300-500 |
| common_mistakes | ~200-300 |
| edge_cases | ~100-150 |
| problems | ~50-100 (initial) |
| problem_topics | ~100-200 |

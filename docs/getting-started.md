# Algoholic — Getting Started
## Setup Guide & Implementation Checklist

---

## Prerequisites

```bash
docker --version          # Docker 20.10+
docker-compose --version  # Docker Compose 2.0+
python --version          # Python 3.10+
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

### 4. Backend

Create `backend/requirements.txt`:

```
fastapi==0.104.1
uvicorn[standard]==0.24.0
sqlalchemy==2.0.23
psycopg2-binary==2.9.9
python-dotenv==1.0.0
pydantic==2.5.0
chromadb==0.4.18
sentence-transformers==2.2.2
ollama==0.1.6
```

Create `backend/main.py`:

```python
from fastapi import FastAPI, Depends, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from sqlalchemy import create_engine, Column, Integer, String, Float, Text
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker, Session
from pydantic import BaseModel
from typing import List, Optional
import os

DATABASE_URL = os.getenv("DATABASE_URL", "postgresql://leetcode:leetcode123@localhost:5432/leetcode_training")
engine = create_engine(DATABASE_URL)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)
Base = declarative_base()

app = FastAPI(title="Algoholic API", version="1.0.0")
app.add_middleware(CORSMiddleware, allow_origins=["*"], allow_methods=["*"], allow_headers=["*"])

def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()

class Problem(Base):
    __tablename__ = "problems"
    problem_id = Column(Integer, primary_key=True)
    leetcode_number = Column(Integer, unique=True)
    title = Column(String(500), nullable=False)
    description = Column(Text, nullable=False)
    difficulty_score = Column(Float, nullable=False)
    primary_pattern = Column(String(100))

class ProblemResponse(BaseModel):
    problem_id: int
    title: str
    difficulty_score: float
    primary_pattern: Optional[str]
    class Config:
        from_attributes = True

@app.get("/api/problems", response_model=List[ProblemResponse])
async def get_problems(min_difficulty: float = 0, max_difficulty: float = 100, limit: int = 20, db: Session = Depends(get_db)):
    return db.query(Problem).filter(Problem.difficulty_score.between(min_difficulty, max_difficulty)).limit(limit).all()

@app.get("/api/problems/{problem_id}", response_model=ProblemResponse)
async def get_problem(problem_id: int, db: Session = Depends(get_db)):
    p = db.query(Problem).filter(Problem.problem_id == problem_id).first()
    if not p: raise HTTPException(404, "Problem not found")
    return p

@app.get("/health")
async def health(): return {"status": "healthy"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8080)
```

### 5. Start Everything

```bash
docker-compose up -d
sleep 30

# Download LLM models
docker exec -it algoholic_ollama ollama pull mistral:7b
docker exec -it algoholic_ollama ollama pull codellama:13b

# Start backend
cd backend
pip install -r requirements.txt
python main.py
```

### 6. Verify

```bash
curl http://localhost:8080/health                    # API health
curl http://localhost:8080/api/problems | jq          # List problems
curl http://localhost:8000/api/v1/heartbeat           # ChromaDB
docker exec algoholic_ollama ollama list              # LLM models
docker exec algoholic_postgres psql -U leetcode -d leetcode_training -c "SELECT count(*) FROM problems;"
```

**Access points:**
- API: http://localhost:8080
- API Docs: http://localhost:8080/docs
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

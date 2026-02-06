# LeetCode Training Platform - Quick Start Guide
## Get Running Locally in 30 Minutes

---

## Prerequisites

**Required Software:**
```bash
# Check if you have these installed
docker --version          # Docker 20.10+
docker-compose --version  # Docker Compose 2.0+
python --version         # Python 3.10+
node --version          # Node.js 18+
```

**Install if missing:**
- **Docker Desktop**: https://www.docker.com/products/docker-desktop
- **Python**: https://www.python.org/downloads/
- **Node.js**: https://nodejs.org/

---

## Step 1: Project Setup (5 minutes)

```bash
# Create project directory
mkdir leetcode-training-platform
cd leetcode-training-platform

# Create directory structure
mkdir -p backend/{api,models,services,utils}
mkdir -p frontend/src
mkdir -p data/{postgres,redis,chroma,ollama}
mkdir -p scripts

# Initialize git (optional)
git init
```

---

## Step 2: Docker Compose Configuration (5 minutes)

**Create `docker-compose.yml`:**

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:16
    container_name: leetcode_postgres
    environment:
      POSTGRES_USER: leetcode
      POSTGRES_PASSWORD: leetcode123
      POSTGRES_DB: leetcode_training
    ports:
      - "5432:5432"
    volumes:
      - ./data/postgres:/var/lib/postgresql/data
      - ./scripts/init_db.sql:/docker-entrypoint-initdb.d/init.sql
    command: postgres -c shared_preload_libraries=vector

  chromadb:
    image: chromadb/chroma:latest
    container_name: leetcode_chromadb
    ports:
      - "8000:8000"
    volumes:
      - ./data/chroma:/chroma/chroma
    environment:
      IS_PERSISTENT: "TRUE"

  ollama:
    image: ollama/ollama:latest
    container_name: leetcode_ollama
    ports:
      - "11434:11434"
    volumes:
      - ./data/ollama:/root/.ollama

  redis:
    image: redis:7-alpine
    container_name: leetcode_redis
    ports:
      - "6379:6379"
    volumes:
      - ./data/redis:/data

networks:
  default:
    name: leetcode_network
```

---

## Step 3: Database Initialization (5 minutes)

**Create `scripts/init_db.sql`:**

```sql
-- Enable extensions
CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS btree_gin;

-- Core tables
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE problems (
    problem_id SERIAL PRIMARY KEY,
    leetcode_number INTEGER UNIQUE,
    title VARCHAR(500) NOT NULL,
    statement TEXT NOT NULL,
    constraints TEXT[],
    examples JSONB NOT NULL,
    hints TEXT[],
    difficulty_score FLOAT NOT NULL CHECK (difficulty_score >= 0 AND difficulty_score <= 100),
    primary_pattern VARCHAR(100),
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
    description TEXT,
    start_date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_problems_difficulty ON problems(difficulty_score);
CREATE INDEX idx_questions_type ON questions(question_type);
CREATE INDEX idx_attempts_user ON user_attempts(user_id, attempted_at DESC);

-- Sample data
INSERT INTO users (username, email, password_hash) VALUES
('demo_user', 'demo@example.com', 'dummy_hash_change_in_production');

-- Insert sample problem
INSERT INTO problems (
    leetcode_number, 
    title, 
    statement, 
    constraints, 
    examples,
    hints,
    difficulty_score,
    primary_pattern
) VALUES (
    1,
    'Two Sum',
    'Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.',
    ARRAY['2 <= nums.length <= 10^4', '-10^9 <= nums[i] <= 10^9', 'Only one valid answer exists'],
    '[
        {
            "input": "nums = [2,7,11,15], target = 9",
            "output": "[0,1]",
            "explanation": "Because nums[0] + nums[1] == 9, we return [0, 1]."
        },
        {
            "input": "nums = [3,2,4], target = 6",
            "output": "[1,2]",
            "explanation": "Because nums[1] + nums[2] == 6, we return [1, 2]."
        }
    ]'::jsonb,
    ARRAY['Try using a hash table to store complements'],
    15.5,
    'Hash Table'
);

-- Insert sample question
INSERT INTO questions (
    problem_id,
    question_type,
    question_text,
    correct_answer,
    answer_options,
    explanation,
    difficulty_score
) VALUES (
    1,
    'complexity_analysis',
    'What is the time complexity of the optimal solution for Two Sum using a hash table?',
    '{"value": "O(n)"}'::jsonb,
    '[
        {"id": "A", "text": "O(1)"},
        {"id": "B", "text": "O(n)"},
        {"id": "C", "text": "O(n log n)"},
        {"id": "D", "text": "O(n²)"}
    ]'::jsonb,
    'Using a hash table allows us to look up complements in O(1) time, and we only need one pass through the array, giving us O(n) time complexity.',
    20.0
);
```

---

## Step 4: Backend Setup (10 minutes)

**Create `backend/requirements.txt`:**

```txt
fastapi==0.104.1
uvicorn[standard]==0.24.0
sqlalchemy==2.0.23
psycopg2-binary==2.9.9
python-dotenv==1.0.0
pydantic==2.5.0
chromadb==0.4.18
sentence-transformers==2.2.2
ollama==0.1.6
redis==5.0.1
python-jose[cryptography]==3.3.0
passlib[bcrypt]==1.7.4
python-multipart==0.0.6
```

**Create `backend/main.py`:**

```python
from fastapi import FastAPI, Depends, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from sqlalchemy import create_engine, Column, Integer, String, Float, ARRAY, JSON, Text
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker, Session
from pydantic import BaseModel
from typing import List, Optional
import os

# Database setup
DATABASE_URL = os.getenv("DATABASE_URL", "postgresql://leetcode:leetcode123@localhost:5432/leetcode_training")
engine = create_engine(DATABASE_URL)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)
Base = declarative_base()

# FastAPI app
app = FastAPI(title="LeetCode Training API", version="1.0.0")

# CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Dependency
def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()

# Models
class Problem(Base):
    __tablename__ = "problems"
    problem_id = Column(Integer, primary_key=True)
    leetcode_number = Column(Integer, unique=True)
    title = Column(String(500), nullable=False)
    statement = Column(Text, nullable=False)
    difficulty_score = Column(Float, nullable=False)
    primary_pattern = Column(String(100))

class Question(Base):
    __tablename__ = "questions"
    question_id = Column(Integer, primary_key=True)
    problem_id = Column(Integer)
    question_type = Column(String(50), nullable=False)
    question_text = Column(Text, nullable=False)
    correct_answer = Column(JSON, nullable=False)
    answer_options = Column(JSON)
    explanation = Column(Text, nullable=False)
    difficulty_score = Column(Float, nullable=False)

# Pydantic schemas
class ProblemResponse(BaseModel):
    problem_id: int
    title: str
    statement: str
    difficulty_score: float
    primary_pattern: Optional[str]
    
    class Config:
        from_attributes = True

class QuestionResponse(BaseModel):
    question_id: int
    question_type: str
    question_text: str
    answer_options: Optional[dict]
    difficulty_score: float
    
    class Config:
        from_attributes = True

# Routes
@app.get("/")
async def root():
    return {
        "message": "LeetCode Training API",
        "version": "1.0.0",
        "endpoints": {
            "problems": "/api/problems",
            "questions": "/api/questions",
            "docs": "/docs"
        }
    }

@app.get("/api/problems", response_model=List[ProblemResponse])
async def get_problems(
    skip: int = 0,
    limit: int = 20,
    min_difficulty: float = 0,
    max_difficulty: float = 100,
    db: Session = Depends(get_db)
):
    problems = db.query(Problem).filter(
        Problem.difficulty_score >= min_difficulty,
        Problem.difficulty_score <= max_difficulty
    ).offset(skip).limit(limit).all()
    return problems

@app.get("/api/problems/{problem_id}", response_model=ProblemResponse)
async def get_problem(problem_id: int, db: Session = Depends(get_db)):
    problem = db.query(Problem).filter(Problem.problem_id == problem_id).first()
    if not problem:
        raise HTTPException(status_code=404, detail="Problem not found")
    return problem

@app.get("/api/questions", response_model=List[QuestionResponse])
async def get_questions(
    question_type: Optional[str] = None,
    skip: int = 0,
    limit: int = 20,
    db: Session = Depends(get_db)
):
    query = db.query(Question)
    if question_type:
        query = query.filter(Question.question_type == question_type)
    questions = query.offset(skip).limit(limit).all()
    return questions

@app.get("/api/questions/{question_id}", response_model=QuestionResponse)
async def get_question(question_id: int, db: Session = Depends(get_db)):
    question = db.query(Question).filter(Question.question_id == question_id).first()
    if not question:
        raise HTTPException(status_code=404, detail="Question not found")
    return question

@app.get("/health")
async def health_check():
    return {"status": "healthy"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8080)
```

**Create `backend/Dockerfile`:**

```dockerfile
FROM python:3.10-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8080"]
```

---

## Step 5: Start the System (2 minutes)

```bash
# Start all services
docker-compose up -d

# Wait for services to be ready (30 seconds)
sleep 30

# Check status
docker-compose ps

# Expected output:
# NAME                 STATUS    PORTS
# leetcode_postgres    Up        0.0.0.0:5432->5432/tcp
# leetcode_chromadb    Up        0.0.0.0:8000->8000/tcp
# leetcode_ollama      Up        0.0.0.0:11434->11434/tcp
# leetcode_redis       Up        0.0.0.0:6379->6379/tcp

# View logs
docker-compose logs -f postgres
```

---

## Step 6: Download LLM Models (3 minutes)

```bash
# Download Mistral 7B for general reasoning
docker exec -it leetcode_ollama ollama pull mistral:7b

# Download Code Llama for code understanding
docker exec -it leetcode_ollama ollama pull codellama:13b

# Test LLM
curl http://localhost:11434/api/generate -d '{
  "model": "mistral:7b",
  "prompt": "Explain binary search in one sentence."
}'
```

---

## Step 7: Start Backend (2 minutes)

```bash
# Install dependencies
cd backend
pip install -r requirements.txt

# Run backend
python main.py

# Should see:
# INFO:     Started server process
# INFO:     Uvicorn running on http://0.0.0.0:8080
```

**Test API:**
```bash
# Health check
curl http://localhost:8080/health

# Get problems
curl http://localhost:8080/api/problems

# API documentation
open http://localhost:8080/docs
```

---

## Step 8: Verify Setup

**Database Check:**
```bash
# Connect to PostgreSQL
docker exec -it leetcode_postgres psql -U leetcode -d leetcode_training

# Run query
\dt  # List tables
SELECT COUNT(*) FROM problems;  # Should return 1
SELECT title FROM problems;     # Should show "Two Sum"
\q   # Quit
```

**ChromaDB Check:**
```bash
# Check ChromaDB
curl http://localhost:8000/api/v1/heartbeat
# Should return: {"nanosecond heartbeat": ...}
```

**Ollama Check:**
```bash
# List downloaded models
docker exec -it leetcode_ollama ollama list
# Should show: mistral:7b, codellama:13b
```

---

## Step 9: Load Sample Data (Optional)

**Create `scripts/load_sample_data.py`:**

```python
import psycopg2
import json

# Connect to database
conn = psycopg2.connect(
    host="localhost",
    database="leetcode_training",
    user="leetcode",
    password="leetcode123"
)
cur = conn.cursor()

# Sample problems
problems = [
    {
        "leetcode_number": 15,
        "title": "3Sum",
        "statement": "Given an array nums of n integers, find all unique triplets [nums[i], nums[j], nums[k]] such that i != j, i != k, and j != k, and nums[i] + nums[j] + nums[k] == 0.",
        "constraints": ["3 <= nums.length <= 3000", "-10^5 <= nums[i] <= 10^5"],
        "examples": json.dumps([
            {
                "input": "nums = [-1,0,1,2,-1,-4]",
                "output": "[[-1,-1,2],[-1,0,1]]",
                "explanation": "The distinct triplets are [-1,0,1] and [-1,-1,2]."
            }
        ]),
        "hints": ["Sort the array first", "Use two pointers after fixing one element"],
        "difficulty_score": 52.0,
        "primary_pattern": "Two Pointers"
    },
    {
        "leetcode_number": 121,
        "title": "Best Time to Buy and Sell Stock",
        "statement": "You are given an array prices where prices[i] is the price of a given stock on the ith day. You want to maximize profit by choosing a single day to buy and a single day to sell.",
        "constraints": ["1 <= prices.length <= 10^5", "0 <= prices[i] <= 10^4"],
        "examples": json.dumps([
            {
                "input": "prices = [7,1,5,3,6,4]",
                "output": "5",
                "explanation": "Buy on day 2 (price = 1) and sell on day 5 (price = 6), profit = 6-1 = 5."
            }
        ]),
        "hints": ["Track minimum price seen so far", "Calculate profit at each step"],
        "difficulty_score": 25.0,
        "primary_pattern": "Array"
    }
]

# Insert problems
for p in problems:
    cur.execute("""
        INSERT INTO problems (
            leetcode_number, title, statement, constraints, 
            examples, hints, difficulty_score, primary_pattern
        ) VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
    """, (
        p["leetcode_number"], p["title"], p["statement"], 
        p["constraints"], p["examples"], p["hints"],
        p["difficulty_score"], p["primary_pattern"]
    ))

conn.commit()
cur.close()
conn.close()

print("✅ Sample data loaded!")
```

**Run it:**
```bash
python scripts/load_sample_data.py
```

---

## Testing Everything

**1. Test Database:**
```bash
docker exec -it leetcode_postgres psql -U leetcode -d leetcode_training -c "SELECT COUNT(*) FROM problems;"
# Should return: 3 (or more)
```

**2. Test API:**
```bash
# Get all problems
curl http://localhost:8080/api/problems | jq

# Get specific problem
curl http://localhost:8080/api/problems/1 | jq

# Filter by difficulty
curl "http://localhost:8080/api/problems?min_difficulty=20&max_difficulty=60" | jq
```

**3. Test LLM:**
```python
import ollama

response = ollama.generate(
    model='mistral:7b',
    prompt='What is the time complexity of binary search?'
)

print(response['response'])
```

---

## Common Issues & Solutions

**Issue: PostgreSQL won't start**
```bash
# Check logs
docker-compose logs postgres

# Solution: Remove old data
rm -rf data/postgres/*
docker-compose up -d postgres
```

**Issue: Ollama models not downloading**
```bash
# Check disk space
df -h

# Download manually
docker exec -it leetcode_ollama ollama pull mistral:7b
```

**Issue: Port already in use**
```bash
# Find process using port
lsof -i :5432
# or
netstat -anv | grep 5432

# Kill it or change port in docker-compose.yml
```

**Issue: Cannot connect to database**
```bash
# Test connection
docker exec -it leetcode_postgres pg_isready -U leetcode

# Check credentials
docker exec -it leetcode_postgres env | grep POSTGRES
```

---

## Next Steps

**You now have:**
✅ PostgreSQL database with sample data
✅ ChromaDB vector database
✅ Ollama with LLM models
✅ Redis cache
✅ FastAPI backend with basic endpoints
✅ Sample problems and questions

**To continue development:**

1. **Add Apache AGE** for graph relationships:
```bash
# Use Apache AGE Docker image instead
# Update docker-compose.yml:
  postgres:
    image: apache/age:PG16_latest
```

2. **Implement Vector Embeddings:**
```python
from sentence_transformers import SentenceTransformer
import chromadb

# Initialize
model = SentenceTransformer('all-MiniLM-L6-v2')
chroma_client = chromadb.HttpClient(host='localhost', port=8000)

# Create collection
collection = chroma_client.create_collection(name="problems")

# Add problem embeddings
for problem in problems:
    embedding = model.encode(problem['statement'])
    collection.add(
        documents=[problem['statement']],
        embeddings=[embedding.tolist()],
        metadatas=[{"problem_id": problem['problem_id']}],
        ids=[f"problem_{problem['problem_id']}"]
    )
```

3. **Build Frontend:**
```bash
npx create-react-app frontend
cd frontend
npm install axios react-router-dom recharts
npm start
```

4. **Add Assessment Engine:**
- Implement LLM prompt engineering
- Build weakness detection
- Create personalized recommendations

---

## Quick Reference

**Useful Commands:**

```bash
# Start system
docker-compose up -d

# Stop system
docker-compose down

# Restart a service
docker-compose restart postgres

# View logs
docker-compose logs -f backend

# Clean everything (⚠️ deletes data)
docker-compose down -v
rm -rf data/*

# Backup database
docker exec leetcode_postgres pg_dump -U leetcode leetcode_training > backup.sql

# Restore database
docker exec -i leetcode_postgres psql -U leetcode leetcode_training < backup.sql

# Python dependencies
pip install -r backend/requirements.txt

# Update LLM model
docker exec -it leetcode_ollama ollama pull mistral:7b
```

**Access Points:**
- API: http://localhost:8080
- API Docs: http://localhost:8080/docs
- ChromaDB: http://localhost:8000
- Ollama: http://localhost:11434
- PostgreSQL: localhost:5432
- Redis: localhost:6379

---

You're now ready to start building! 🚀

# LeetCode Training Platform - Technical Architecture
## System Design, Database Schema & Implementation Plan

---

## TABLE OF CONTENTS

1. [System Architecture Overview](#system-architecture)
2. [Database Design (PostgreSQL + Apache AGE)](#database-design)
3. [Graph Database Schema (Apache AGE)](#graph-schema)
4. [Vector Database + RAG Architecture](#vector-rag)
5. [Difficulty Scoring System (Magic Unit)](#difficulty-scoring)
6. [Smart Search & Filtering](#search-filter)
7. [Training Plan Builder](#training-plan)
8. [LLM-Powered Assessment](#llm-assessment)
9. [Problem Generation System](#problem-generation)
10. [Local Deployment Architecture](#local-deployment)
11. [API Design](#api-design)
12. [Implementation Roadmap](#implementation)

---

## 1. SYSTEM ARCHITECTURE OVERVIEW

```
┌─────────────────────────────────────────────────────────────────┐
│                        USER INTERFACE                           │
│                   (Web App / Desktop App)                       │
└────────────────────────┬────────────────────────────────────────┘
                         │
┌────────────────────────┴────────────────────────────────────────┐
│                      API GATEWAY LAYER                          │
│              (FastAPI / Express.js / Flask)                     │
└─┬──────────┬──────────┬──────────┬───────────┬─────────────────┘
  │          │          │          │           │
  │          │          │          │           │
┌─▼──────┐ ┌▼────────┐ ┌▼────────┐ ┌▼────────┐ ┌▼──────────────┐
│Question│ │Training │ │Assess-  │ │Search   │ │Problem        │
│Service │ │Plan     │ │ment     │ │&Filter  │ │Generation     │
│        │ │Builder  │ │Engine   │ │Engine   │ │Service        │
└─┬──────┘ └┬────────┘ └┬────────┘ └┬────────┘ └┬──────────────┘
  │         │           │           │           │
  │         │           │           │           │
┌─▼─────────▼───────────▼───────────▼───────────▼───────────────┐
│                    DATA ACCESS LAYER                           │
└─┬──────────┬──────────────────┬──────────────┬────────────────┘
  │          │                  │              │
┌─▼─────────┐│                 ┌▼─────────────┐│
│PostgreSQL ││                 │Vector DB     ││
│+ Apache   ││                 │(ChromaDB/    ││
│AGE        ││                 │Qdrant/       ││
│(Graph)    ││                 │Weaviate)     ││
└───────────┘│                 └──────────────┘│
             │                                  │
        ┌────▼────────┐                   ┌────▼─────────┐
        │Redis Cache  │                   │LLM Service   │
        │(Session,    │                   │(Ollama/      │
        │Leaderboard) │                   │llama.cpp)    │
        └─────────────┘                   └──────────────┘
```

### **Key Components:**

**1. PostgreSQL + Apache AGE**
- Primary relational database
- Graph extensions for problem relationships
- ACID transactions
- Complex queries

**2. Vector Database (ChromaDB/Qdrant)**
- Embeddings for problems, solutions, explanations
- Semantic search
- RAG context retrieval
- Similar problem finding

**3. LLM Service (Local)**
- Ollama (easy local deployment)
- llama.cpp (performance)
- Code Llama / DeepSeek Coder (code understanding)
- Mistral 7B (general reasoning)

**4. Redis (Optional but Recommended)**
- Session management
- Real-time leaderboards
- Rate limiting
- Cache frequent queries

---

## 2. DATABASE DESIGN (PostgreSQL + Apache AGE)

### **2.1 Core Relational Tables**

#### **Table: `problems`**
```sql
CREATE TABLE problems (
    problem_id SERIAL PRIMARY KEY,
    
    -- Identification
    leetcode_number INTEGER UNIQUE,           -- LeetCode #, null if custom
    problem_slug VARCHAR(200) UNIQUE NOT NULL, -- URL-friendly identifier
    title VARCHAR(500) NOT NULL,
    
    -- Content
    statement TEXT NOT NULL,                   -- Problem description
    constraints TEXT[],                        -- List of constraints
    examples JSONB NOT NULL,                   -- [{input, output, explanation}]
    hints TEXT[],                              -- Progressive hints
    
    -- Metadata
    difficulty_score FLOAT NOT NULL,           -- Our magic unit (0-100)
    official_difficulty VARCHAR(20),           -- LeetCode difficulty (Easy/Medium/Hard)
    
    -- Categorization
    primary_pattern VARCHAR(100),              -- Main algorithmic pattern
    secondary_patterns VARCHAR(100)[],         -- Additional patterns involved
    
    -- Relationships (handled in graph, but denormalized for performance)
    has_follow_ups BOOLEAN DEFAULT FALSE,
    has_prerequisites BOOLEAN DEFAULT FALSE,
    
    -- Tracking
    total_attempts INTEGER DEFAULT 0,
    total_solves INTEGER DEFAULT 0,
    average_time_seconds FLOAT,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Indexes
    CONSTRAINT chk_difficulty_score CHECK (difficulty_score >= 0 AND difficulty_score <= 100)
);

CREATE INDEX idx_problems_difficulty ON problems(difficulty_score);
CREATE INDEX idx_problems_pattern ON problems(primary_pattern);
CREATE INDEX idx_problems_leetcode ON problems(leetcode_number) WHERE leetcode_number IS NOT NULL;
```

**Why This Design:**
- `problem_id` for internal use, `leetcode_number` for external reference
- JSONB for flexible example storage
- Arrays for multi-valued attributes (constraints, hints, patterns)
- Denormalized counters for performance
- Partial index on leetcode_number (saves space for custom problems)

---

#### **Table: `problem_content_versions`**
```sql
CREATE TABLE problem_content_versions (
    version_id SERIAL PRIMARY KEY,
    problem_id INTEGER REFERENCES problems(problem_id) ON DELETE CASCADE,
    
    version_number INTEGER NOT NULL,
    
    -- What changed
    statement TEXT NOT NULL,
    constraints TEXT[],
    examples JSONB NOT NULL,
    
    -- Meta
    change_reason TEXT,
    changed_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(problem_id, version_number)
);
```

**Why This Table:**
- Track problem statement changes
- Allow rollback
- Audit trail
- A/B testing different versions

---

#### **Table: `topics`**
```sql
CREATE TABLE topics (
    topic_id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,        -- "Binary Search", "Dynamic Programming"
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    parent_topic_id INTEGER REFERENCES topics(topic_id), -- For hierarchical topics
    level INTEGER DEFAULT 0,                   -- Depth in hierarchy
    
    -- Meta
    problem_count INTEGER DEFAULT 0,
    average_difficulty FLOAT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_topics_parent ON topics(parent_topic_id);
```

**Example Hierarchy:**
```
Dynamic Programming (level 0)
  ├─ 1D DP (level 1)
  │   ├─ Fibonacci-style (level 2)
  │   └─ House Robber variants (level 2)
  └─ 2D DP (level 1)
      ├─ Grid DP (level 2)
      └─ String DP (level 2)
```

---

#### **Table: `problem_topics`** (Many-to-Many)
```sql
CREATE TABLE problem_topics (
    problem_id INTEGER REFERENCES problems(problem_id) ON DELETE CASCADE,
    topic_id INTEGER REFERENCES topics(topic_id) ON DELETE CASCADE,
    
    relevance_score FLOAT DEFAULT 1.0,         -- How central is this topic? (0-1)
    is_primary BOOLEAN DEFAULT FALSE,          -- Main topic for this problem
    
    PRIMARY KEY (problem_id, topic_id)
);

CREATE INDEX idx_problem_topics_problem ON problem_topics(problem_id);
CREATE INDEX idx_problem_topics_topic ON problem_topics(topic_id);
CREATE INDEX idx_problem_topics_primary ON problem_topics(topic_id) WHERE is_primary = TRUE;
```

---

#### **Table: `questions`**
```sql
CREATE TABLE questions (
    question_id SERIAL PRIMARY KEY,
    problem_id INTEGER REFERENCES problems(problem_id) ON DELETE CASCADE,
    
    -- Question type
    question_type VARCHAR(50) NOT NULL,        -- 'complexity', 'ds_selection', 'pattern', etc.
    question_subtype VARCHAR(50),              -- More specific categorization
    
    -- Content
    question_text TEXT NOT NULL,
    question_data JSONB,                       -- Type-specific data (code snippets, etc.)
    
    -- Answer
    answer_type VARCHAR(20) NOT NULL,          -- 'multiple_choice', 'code', 'text', 'ranking'
    correct_answer JSONB NOT NULL,             -- Flexible answer storage
    answer_options JSONB,                      -- For multiple choice
    
    -- Metadata
    difficulty_score FLOAT NOT NULL,
    estimated_time_seconds INTEGER,            -- Expected solve time
    
    -- Learning
    explanation TEXT NOT NULL,                 -- Why this answer
    wrong_answer_explanations JSONB,           -- Explain each wrong option
    related_concepts TEXT[],
    
    -- Tracking
    total_attempts INTEGER DEFAULT 0,
    correct_attempts INTEGER DEFAULT 0,
    average_time_seconds FLOAT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_questions_problem ON questions(problem_id);
CREATE INDEX idx_questions_type ON questions(question_type, question_subtype);
CREATE INDEX idx_questions_difficulty ON questions(difficulty_score);
```

**Why Flexible Schema:**
- JSONB allows different question structures
- `question_type` determines how to interpret data
- Same table for all question types (simpler queries)
- Easy to add new question types

**Example `question_data` for Complexity Analysis:**
```json
{
  "code_snippet": "for(int i=0; i<n; i++) { for(int j=i; j<n; j++) { ... } }",
  "language": "cpp",
  "context": "Finding pairs with target sum"
}
```

**Example `correct_answer` for Multiple Choice:**
```json
{
  "option_id": "B",
  "value": "O(n²)"
}
```

---

#### **Table: `code_templates`**
```sql
CREATE TABLE code_templates (
    template_id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(200) UNIQUE NOT NULL,
    
    category VARCHAR(100) NOT NULL,            -- 'graph', 'binary_search', 'dp', etc.
    
    description TEXT,
    when_to_use TEXT NOT NULL,
    
    -- Template in multiple languages
    cpp_template TEXT,
    python_template TEXT,
    java_template TEXT,
    
    -- Usage
    complexity_time VARCHAR(100),
    complexity_space VARCHAR(100),
    
    -- Related
    related_patterns VARCHAR(100)[],
    example_problems INTEGER[],                -- Array of problem_ids
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_templates_category ON code_templates(category);
```

---

#### **Table: `user_progress`**
```sql
CREATE TABLE user_progress (
    progress_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
    problem_id INTEGER REFERENCES problems(problem_id) ON DELETE CASCADE,
    
    -- Status
    status VARCHAR(20) NOT NULL,               -- 'attempted', 'solved', 'mastered'
    
    -- Attempts
    total_attempts INTEGER DEFAULT 0,
    successful_attempts INTEGER DEFAULT 0,
    
    -- Performance
    best_time_seconds INTEGER,
    best_complexity_time VARCHAR(100),
    best_complexity_space VARCHAR(100),
    
    -- Tracking
    first_attempt_at TIMESTAMP,
    last_attempt_at TIMESTAMP,
    mastered_at TIMESTAMP,
    
    -- Spaced Repetition
    next_review_at TIMESTAMP,
    review_interval_days INTEGER DEFAULT 1,
    ease_factor FLOAT DEFAULT 2.5,             -- SM-2 algorithm
    
    UNIQUE(user_id, problem_id)
);

CREATE INDEX idx_user_progress_user ON user_progress(user_id);
CREATE INDEX idx_user_progress_status ON user_progress(user_id, status);
CREATE INDEX idx_user_progress_review ON user_progress(user_id, next_review_at) 
    WHERE next_review_at IS NOT NULL;
```

**Spaced Repetition Logic:**
- Implements SuperMemo-2 (SM-2) algorithm
- `ease_factor`: How easy the problem is for this user
- `review_interval_days`: Days until next review
- Adaptive based on performance

---

#### **Table: `user_attempts`**
```sql
CREATE TABLE user_attempts (
    attempt_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
    problem_id INTEGER REFERENCES problems(problem_id) ON DELETE CASCADE,
    question_id INTEGER REFERENCES questions(question_id) ON DELETE SET NULL,
    
    -- Attempt details
    attempt_type VARCHAR(50) NOT NULL,         -- 'practice', 'training_plan', 'contest', 'assessment'
    
    -- Answer
    user_answer JSONB NOT NULL,
    is_correct BOOLEAN NOT NULL,
    
    -- Performance
    time_taken_seconds INTEGER NOT NULL,
    hints_used INTEGER DEFAULT 0,
    
    -- Context
    training_plan_id INTEGER REFERENCES training_plans(plan_id),
    assessment_session_id INTEGER REFERENCES assessment_sessions(session_id),
    
    -- Metadata
    attempted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_attempts_user ON user_attempts(user_id, attempted_at DESC);
CREATE INDEX idx_attempts_problem ON user_attempts(problem_id);
CREATE INDEX idx_attempts_training_plan ON user_attempts(training_plan_id) 
    WHERE training_plan_id IS NOT NULL;
```

---

#### **Table: `training_plans`**
```sql
CREATE TABLE training_plans (
    plan_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
    
    -- Plan details
    name VARCHAR(200) NOT NULL,
    description TEXT,
    goal TEXT,                                 -- "Master Dynamic Programming", etc.
    
    -- Configuration
    target_topics INTEGER[],                   -- Array of topic_ids
    target_patterns VARCHAR(100)[],
    difficulty_range NUMRANGE,                 -- Range of difficulty scores
    
    estimated_duration_days INTEGER,
    questions_per_day INTEGER DEFAULT 5,
    
    -- Scheduling
    start_date DATE NOT NULL,
    end_date DATE,
    
    -- Progress
    total_questions INTEGER NOT NULL,
    completed_questions INTEGER DEFAULT 0,
    
    status VARCHAR(20) DEFAULT 'active',       -- 'active', 'completed', 'paused', 'abandoned'
    
    -- Smart features
    adaptive_difficulty BOOLEAN DEFAULT TRUE,
    include_reviews BOOLEAN DEFAULT TRUE,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_training_plans_user ON training_plans(user_id);
CREATE INDEX idx_training_plans_status ON training_plans(user_id, status);
```

---

#### **Table: `training_plan_questions`**
```sql
CREATE TABLE training_plan_questions (
    plan_question_id SERIAL PRIMARY KEY,
    plan_id INTEGER REFERENCES training_plans(plan_id) ON DELETE CASCADE,
    question_id INTEGER REFERENCES questions(question_id) ON DELETE CASCADE,
    
    -- Order
    sequence_number INTEGER NOT NULL,
    scheduled_for DATE,
    
    -- Status
    status VARCHAR(20) DEFAULT 'pending',      -- 'pending', 'completed', 'skipped'
    
    -- Attempts
    attempts INTEGER DEFAULT 0,
    is_correct BOOLEAN,
    time_taken_seconds INTEGER,
    
    completed_at TIMESTAMP,
    
    UNIQUE(plan_id, sequence_number)
);

CREATE INDEX idx_plan_questions_plan ON training_plan_questions(plan_id, sequence_number);
CREATE INDEX idx_plan_questions_scheduled ON training_plan_questions(plan_id, scheduled_for) 
    WHERE status = 'pending';
```

---

#### **Table: `assessment_sessions`**
```sql
CREATE TABLE assessment_sessions (
    session_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
    
    -- Session config
    assessment_type VARCHAR(50) NOT NULL,      -- 'initial', 'progress_check', 'final_exam'
    topic_focus INTEGER[],                     -- Array of topic_ids, null for comprehensive
    
    -- Timing
    duration_seconds INTEGER NOT NULL,
    started_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP,
    
    -- Results
    total_questions INTEGER NOT NULL,
    correct_answers INTEGER DEFAULT 0,
    
    overall_score FLOAT,                       -- 0-100
    category_scores JSONB,                     -- Score per category
    
    -- Analysis (LLM-powered)
    strengths TEXT[],
    weaknesses TEXT[],
    recommendations TEXT,
    
    status VARCHAR(20) DEFAULT 'in_progress',  -- 'in_progress', 'completed', 'abandoned'
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_assessments_user ON assessment_sessions(user_id, started_at DESC);
```

---

#### **Table: `assessment_questions`**
```sql
CREATE TABLE assessment_questions (
    assessment_question_id SERIAL PRIMARY KEY,
    session_id INTEGER REFERENCES assessment_sessions(session_id) ON DELETE CASCADE,
    question_id INTEGER REFERENCES questions(question_id),
    
    sequence_number INTEGER NOT NULL,
    
    -- Answer
    user_answer JSONB,
    is_correct BOOLEAN,
    time_taken_seconds INTEGER,
    
    -- Analysis
    difficulty_perceived VARCHAR(20),          -- User-reported difficulty
    confidence_level INTEGER,                  -- 1-5 scale
    
    answered_at TIMESTAMP,
    
    UNIQUE(session_id, sequence_number)
);
```

---

#### **Table: `weakness_analysis`**
```sql
CREATE TABLE weakness_analysis (
    analysis_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
    
    -- What's weak
    category VARCHAR(100) NOT NULL,            -- 'complexity_analysis', 'pattern_recognition', etc.
    specific_topic INTEGER REFERENCES topics(topic_id),
    
    -- Severity
    weakness_score FLOAT NOT NULL,             -- 0-1 (higher = more problematic)
    confidence FLOAT NOT NULL,                 -- 0-1 (how sure are we?)
    
    -- Evidence
    failed_questions INTEGER[],
    evidence_text TEXT,
    
    -- Recommendations
    recommended_practice TEXT[],
    recommended_resources TEXT[],
    
    -- Tracking
    identified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP,
    
    UNIQUE(user_id, category, specific_topic)
);

CREATE INDEX idx_weakness_user ON weakness_analysis(user_id, weakness_score DESC) 
    WHERE resolved_at IS NULL;
```

---

#### **Table: `embeddings`**
```sql
CREATE TABLE embeddings (
    embedding_id SERIAL PRIMARY KEY,
    
    -- What is embedded
    entity_type VARCHAR(50) NOT NULL,          -- 'problem', 'question', 'solution', 'explanation'
    entity_id INTEGER NOT NULL,
    
    -- Embedding
    model_name VARCHAR(100) NOT NULL,          -- Which embedding model
    embedding_vector VECTOR(768),              -- Or 1536 for OpenAI, 384 for smaller models
    
    -- Metadata
    text_content TEXT NOT NULL,                -- Original text that was embedded
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(entity_type, entity_id, model_name)
);

-- For similarity search (requires pgvector extension)
CREATE INDEX idx_embeddings_vector ON embeddings 
    USING ivfflat (embedding_vector vector_cosine_ops);

CREATE INDEX idx_embeddings_entity ON embeddings(entity_type, entity_id);
```

**Note:** Requires `pgvector` extension:
```sql
CREATE EXTENSION vector;
```

---

## 3. GRAPH DATABASE SCHEMA (Apache AGE)

### **3.1 Why Graph Database?**

**Problems naturally form a graph:**
- Similar problems
- Follow-up problems
- Prerequisites
- Variations
- Topic relationships

**Graph queries we need:**
- "Find all follow-ups of problem X"
- "What problems should I solve before this?"
- "Find similar problems I haven't tried"
- "Path from beginner to expert in topic Y"
- "Problems that use techniques A AND B"

---

### **3.2 Graph Schema**

**Node Types:**

1. **Problem Node**
```cypher
(:Problem {
    id: integer,
    title: string,
    difficulty_score: float,
    primary_pattern: string
})
```

2. **Topic Node**
```cypher
(:Topic {
    id: integer,
    name: string,
    level: integer
})
```

3. **Pattern Node**
```cypher
(:Pattern {
    name: string,
    category: string
})
```

4. **Company Node**
```cypher
(:Company {
    name: string,
    frequency: integer  // How often they ask questions
})
```

---

**Edge Types:**

1. **SIMILAR_TO**
```cypher
(:Problem)-[:SIMILAR_TO {
    similarity_score: float,      // 0-1
    reason: string               // "Same pattern, different constraints"
}]->(:Problem)
```

2. **FOLLOW_UP_OF**
```cypher
(:Problem)-[:FOLLOW_UP_OF {
    difficulty_increase: float,
    new_constraints: [string]
}]->(:Problem)
```

3. **PREREQUISITE_OF**
```cypher
(:Problem)-[:PREREQUISITE_OF {
    is_hard_requirement: boolean
}]->(:Problem)
```

4. **VARIATION_OF**
```cypher
(:Problem)-[:VARIATION_OF {
    variation_type: string       // "Different DS", "Added constraint", etc.
}]->(:Problem)
```

5. **HAS_TOPIC**
```cypher
(:Problem)-[:HAS_TOPIC {
    relevance: float,
    is_primary: boolean
}]->(:Topic)
```

6. **USES_PATTERN**
```cypher
(:Problem)-[:USES_PATTERN {
    is_primary: boolean
}]->(:Pattern)
```

7. **SUBTOPIC_OF**
```cypher
(:Topic)-[:SUBTOPIC_OF]->(:Topic)
```

8. **ASKED_BY**
```cypher
(:Problem)-[:ASKED_BY {
    frequency: integer,
    last_seen: date
}]->(:Company)
```

---

### **3.3 Apache AGE Setup**

**Install AGE extension:**
```sql
CREATE EXTENSION age;
LOAD 'age';
SET search_path = ag_catalog, "$user", public;
```

**Create graph:**
```sql
SELECT create_graph('problem_graph');
```

**Create nodes:**
```sql
-- Create Problem node
SELECT * FROM cypher('problem_graph', $$
    CREATE (:Problem {
        id: 1,
        title: 'Two Sum',
        difficulty_score: 15.5,
        primary_pattern: 'Hash Table'
    })
$$) as (v agtype);

-- Create Topic node
SELECT * FROM cypher('problem_graph', $$
    CREATE (:Topic {
        id: 1,
        name: 'Array',
        level: 0
    })
$$) as (v agtype);
```

**Create relationships:**
```sql
-- Link problem to topic
SELECT * FROM cypher('problem_graph', $$
    MATCH (p:Problem {id: 1}), (t:Topic {id: 1})
    CREATE (p)-[:HAS_TOPIC {relevance: 0.9, is_primary: true}]->(t)
$$) as (v agtype);

-- Create similar relationship
SELECT * FROM cypher('problem_graph', $$
    MATCH (p1:Problem {id: 1}), (p2:Problem {id: 167})
    CREATE (p1)-[:SIMILAR_TO {
        similarity_score: 0.95,
        reason: 'Same pattern, array is sorted'
    }]->(p2)
$$) as (v agtype);
```

---

### **3.4 Common Graph Queries**

**Find all follow-ups:**
```sql
SELECT * FROM cypher('problem_graph', $$
    MATCH (p:Problem {id: 1})-[:FOLLOW_UP_OF*1..3]->(followup:Problem)
    RETURN followup.title, followup.difficulty_score
    ORDER BY followup.difficulty_score
$$) as (title agtype, difficulty agtype);
```

**Find learning path (prerequisites chain):**
```sql
SELECT * FROM cypher('problem_graph', $$
    MATCH path = (start:Problem {id: 100})<-[:PREREQUISITE_OF*]-(prerequisite:Problem)
    RETURN prerequisite.title, length(path) as distance
    ORDER BY distance DESC
$$) as (title agtype, distance agtype);
```

**Find problems using multiple patterns:**
```sql
SELECT * FROM cypher('problem_graph', $$
    MATCH (p:Problem)-[:USES_PATTERN]->(pat1:Pattern {name: 'Dynamic Programming'})
    MATCH (p)-[:USES_PATTERN]->(pat2:Pattern {name: 'Binary Search'})
    RETURN p.title, p.difficulty_score
$$) as (title agtype, difficulty agtype);
```

**Find similar unsolved problems:**
```sql
-- Assuming we have user_id from application
SELECT * FROM cypher('problem_graph', $$
    MATCH (solved:Problem)<-[:SOLVED]-(u:User {id: $user_id})
    MATCH (solved)-[:SIMILAR_TO]->(similar:Problem)
    WHERE NOT exists((u)-[:SOLVED]->(similar))
    RETURN similar.title, similar.difficulty_score
    LIMIT 10
$$) as (title agtype, difficulty agtype);
```

**Find topic progression path:**
```sql
SELECT * FROM cypher('problem_graph', $$
    MATCH path = (beginner:Topic {name: 'Array'})-[:SUBTOPIC_OF*]->(advanced:Topic)
    RETURN [node in nodes(path) | node.name] as learning_path
$$) as (path agtype);
```

---

## 4. VECTOR DATABASE + RAG ARCHITECTURE

### **4.1 Why Vector DB?**

**Use Cases:**
1. **Semantic search** - Find similar problems by meaning, not keywords
2. **Problem generation** - Generate variations of existing problems
3. **Smart recommendations** - Based on what user struggled with
4. **Context for LLM** - RAG for assessment and explanations

---

### **4.2 Vector DB Choice: ChromaDB**

**Why ChromaDB:**
- Easy local deployment (Python)
- Built-in embedding generation
- Persistent storage
- Good performance for < 1M vectors
- Active development
- Open source

**Alternatives:**
- **Qdrant** - Better for production, Rust-based
- **Weaviate** - More features, heavier
- **pgvector** - Keep everything in PostgreSQL (simpler architecture)

---

### **4.3 What to Embed**

**Collections:**

**1. Problems Collection**
```python
{
    "collection_name": "problems",
    "embedding_function": "all-MiniLM-L6-v2",  # 384 dimensions
    "documents": [
        {
            "id": "problem_1",
            "text": "Problem title + statement + constraints + examples",
            "metadata": {
                "problem_id": 1,
                "difficulty_score": 15.5,
                "patterns": ["Hash Table", "Array"],
                "leetcode_number": 1
            }
        }
    ]
}
```

**2. Solutions Collection**
```python
{
    "collection_name": "solutions",
    "documents": [
        {
            "id": "solution_1_approach_hashmap",
            "text": "Solution approach description + code comments",
            "metadata": {
                "problem_id": 1,
                "approach_name": "Hash Map",
                "time_complexity": "O(n)",
                "space_complexity": "O(n)"
            }
        }
    ]
}
```

**3. Explanations Collection**
```python
{
    "collection_name": "explanations",
    "documents": [
        {
            "id": "explanation_complexity_nested_loops",
            "text": "Detailed explanation of why nested loops are O(n²)",
            "metadata": {
                "category": "complexity_analysis",
                "difficulty": "easy"
            }
        }
    ]
}
```

**4. Code Templates Collection**
```python
{
    "collection_name": "templates",
    "documents": [
        {
            "id": "template_dfs_recursive",
            "text": "DFS template description + when to use + example",
            "metadata": {
                "category": "graph",
                "patterns": ["DFS", "Backtracking"]
            }
        }
    ]
}
```

---

### **4.4 RAG Implementation**

**Architecture:**
```
User Query
    │
    ▼
Query Embedding (using same model)
    │
    ▼
Vector DB Search (top-k similar)
    │
    ▼
Retrieve Documents + Metadata
    │
    ▼
Construct Prompt with Context
    │
    ▼
LLM Generation
    │
    ▼
Post-process & Return
```

**Example RAG Query:**
```python
def rag_query(user_query: str, collection: str, k: int = 5):
    """
    Retrieve relevant context and generate response
    """
    # 1. Get similar documents
    results = chroma_client.query(
        collection_name=collection,
        query_texts=[user_query],
        n_results=k
    )
    
    # 2. Build context
    context = "\n\n".join([
        f"Document {i+1}:\n{doc}"
        for i, doc in enumerate(results['documents'][0])
    ])
    
    # 3. Build prompt
    prompt = f"""You are a coding interview tutor.

Context (relevant information):
{context}

User question: {user_query}

Provide a clear, helpful explanation based on the context above.
"""
    
    # 4. Generate response
    response = llm.generate(prompt)
    
    return response, results['metadatas']
```

---

### **4.5 Embedding Generation Pipeline**

**When to generate embeddings:**
- Problem created/updated → embed problem
- Solution added → embed solution
- Question created → embed question text
- Template added → embed template

**Background job:**
```python
def generate_embeddings_for_problem(problem_id: int):
    """
    Generate all embeddings for a problem
    """
    # Get problem from DB
    problem = db.query(Problem).get(problem_id)
    
    # Combine text
    text = f"""
    Title: {problem.title}
    
    Problem Statement:
    {problem.statement}
    
    Constraints:
    {' '.join(problem.constraints)}
    
    Examples:
    {format_examples(problem.examples)}
    """
    
    # Add to ChromaDB
    chroma_client.add(
        collection_name="problems",
        documents=[text],
        metadatas=[{
            "problem_id": problem.id,
            "difficulty_score": problem.difficulty_score,
            "patterns": problem.secondary_patterns or [],
            "leetcode_number": problem.leetcode_number
        }],
        ids=[f"problem_{problem.id}"]
    )
```

---

## 5. DIFFICULTY SCORING SYSTEM (Magic Unit)

### **5.1 The Magic Unit Formula**

**Goal:** Fair difficulty score (0-100) that accounts for multiple factors

**Components:**

**Base Difficulty (40% weight):**
- Algorithmic complexity
- Number of steps in solution
- Number of edge cases

**Pattern Complexity (25% weight):**
- How many patterns combined
- How obscure the patterns are
- Whether intuition helps

**Implementation Difficulty (20% weight):**
- Code complexity
- Bug-prone areas
- Edge case handling

**Performance from Data (15% weight):**
- Average solve time
- Success rate
- User ratings

---

### **5.2 Detailed Formula**

```python
def calculate_difficulty_score(problem_id: int) -> float:
    """
    Calculate magic unit difficulty score
    Returns: 0-100 float
    """
    
    # 1. BASE DIFFICULTY (0-40 points)
    base_score = calculate_base_difficulty(problem_id)
    
    # 2. PATTERN COMPLEXITY (0-25 points)
    pattern_score = calculate_pattern_complexity(problem_id)
    
    # 3. IMPLEMENTATION DIFFICULTY (0-20 points)
    impl_score = calculate_implementation_difficulty(problem_id)
    
    # 4. EMPIRICAL DATA (0-15 points)
    empirical_score = calculate_empirical_difficulty(problem_id)
    
    # Combine
    total_score = base_score + pattern_score + impl_score + empirical_score
    
    # Normalize to 0-100
    return min(100, max(0, total_score))


def calculate_base_difficulty(problem_id: int) -> float:
    """
    Algorithmic complexity and solution steps
    """
    problem = get_problem(problem_id)
    
    score = 0
    
    # Time complexity required (higher is harder)
    complexity_scores = {
        "O(1)": 1,
        "O(log n)": 3,
        "O(n)": 5,
        "O(n log n)": 10,
        "O(n²)": 15,
        "O(n³)": 25,
        "O(2^n)": 35
    }
    score += complexity_scores.get(problem.optimal_complexity_time, 5)
    
    # Number of solution steps
    num_steps = count_solution_steps(problem_id)
    score += min(15, num_steps * 2)  # Cap at 15
    
    # Edge cases
    num_edge_cases = count_edge_cases(problem_id)
    score += min(10, num_edge_cases * 2)  # Cap at 10
    
    # Number of constraints
    score += min(5, len(problem.constraints))
    
    return min(40, score)


def calculate_pattern_complexity(problem_id: int) -> float:
    """
    How many patterns and how obscure
    """
    problem = get_problem(problem_id)
    
    score = 0
    
    # Number of patterns involved
    num_patterns = 1 + len(problem.secondary_patterns or [])
    
    if num_patterns == 1:
        score += 5  # Single pattern
    elif num_patterns == 2:
        score += 12  # Two patterns combined
    elif num_patterns >= 3:
        score += 20  # Multiple patterns
    
    # Pattern obscurity (some patterns are less intuitive)
    obscurity_scores = {
        "Two Pointers": 2,
        "Sliding Window": 3,
        "Hash Table": 2,
        "Binary Search": 4,
        "DFS": 5,
        "BFS": 5,
        "Dynamic Programming": 12,
        "Backtracking": 10,
        "Union Find": 8,
        "Topological Sort": 9,
        "Monotonic Stack": 7,
        "Trie": 7,
        "Segment Tree": 15,
        "Suffix Array": 18
    }
    
    patterns = [problem.primary_pattern] + (problem.secondary_patterns or [])
    for pattern in patterns:
        score += obscurity_scores.get(pattern, 5)
    
    return min(25, score)


def calculate_implementation_difficulty(problem_id: int) -> float:
    """
    How hard to code correctly
    """
    score = 0
    
    # Lines of code in optimal solution
    loc = get_solution_lines_of_code(problem_id)
    score += min(8, loc / 10)  # Cap at 8
    
    # Pointer manipulation required
    if requires_pointer_manipulation(problem_id):
        score += 5
    
    # Recursion depth/complexity
    if requires_recursion(problem_id):
        score += 4
    
    # Multiple data structures
    num_ds = count_data_structures_needed(problem_id)
    score += min(5, num_ds * 2)
    
    return min(20, score)


def calculate_empirical_difficulty(problem_id: int) -> float:
    """
    Based on actual user performance
    """
    stats = get_problem_stats(problem_id)
    
    if stats.total_attempts < 10:
        return 5  # Not enough data, assume medium
    
    score = 0
    
    # Success rate (lower = harder)
    success_rate = stats.total_solves / stats.total_attempts
    if success_rate < 0.2:
        score += 10  # Very hard
    elif success_rate < 0.4:
        score += 7   # Hard
    elif success_rate < 0.6:
        score += 5   # Medium
    else:
        score += 2   # Easy
    
    # Average time (normalized)
    avg_minutes = stats.average_time_seconds / 60
    if avg_minutes > 45:
        score += 5
    elif avg_minutes > 30:
        score += 3
    elif avg_minutes > 15:
        score += 2
    else:
        score += 1
    
    return min(15, score)
```

---

### **5.3 Difficulty Tiers**

**Mapping scores to tiers:**
```python
def get_difficulty_tier(score: float) -> str:
    """
    Convert score to tier
    """
    if score < 15:
        return "Trivial"
    elif score < 30:
        return "Easy"
    elif score < 50:
        return "Medium"
    elif score < 70:
        return "Hard"
    elif score < 85:
        return "Very Hard"
    else:
        return "Expert"
```

**Visualization:**
```
0────────15───────30──────────50──────────70───────85────100
│ Trivial │  Easy  │   Medium   │   Hard   │ V.Hard│Expert│
```

---

### **5.4 Dynamic Difficulty Adjustment**

**Update difficulty based on data:**
```python
def update_difficulty_score(problem_id: int):
    """
    Recalculate difficulty periodically
    """
    # Only update if enough attempts
    stats = get_problem_stats(problem_id)
    if stats.total_attempts < 10:
        return
    
    # Calculate new score
    new_score = calculate_difficulty_score(problem_id)
    
    # Get current score
    problem = get_problem(problem_id)
    old_score = problem.difficulty_score
    
    # Smooth update (don't change dramatically)
    updated_score = (old_score * 0.7) + (new_score * 0.3)
    
    # Update database
    problem.difficulty_score = updated_score
    problem.updated_at = datetime.now()
    db.commit()
    
    # Log change if significant
    if abs(updated_score - old_score) > 5:
        log_difficulty_change(problem_id, old_score, updated_score)
```

---

## 6. SMART SEARCH & FILTERING

### **6.1 Search Architecture**

**Multi-modal search:**
1. **Keyword search** (PostgreSQL full-text)
2. **Semantic search** (Vector DB)
3. **Graph search** (Apache AGE)
4. **Filtered search** (SQL WHERE clauses)

---

### **6.2 Full-Text Search Setup**

**Add tsvector column:**
```sql
ALTER TABLE problems 
ADD COLUMN search_vector tsvector 
GENERATED ALWAYS AS (
    setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
    setweight(to_tsvector('english', coalesce(statement, '')), 'B') ||
    setweight(to_tsvector('english', array_to_string(constraints, ' ')), 'C')
) STORED;

CREATE INDEX idx_problems_search ON problems USING GIN(search_vector);
```

**Search query:**
```sql
SELECT problem_id, title, difficulty_score,
       ts_rank(search_vector, query) as rank
FROM problems, to_tsquery('english', 'binary & search') query
WHERE search_vector @@ query
ORDER BY rank DESC
LIMIT 20;
```

---

### **6.3 Combined Search API**

```python
def search_problems(
    query: str = None,
    patterns: List[str] = None,
    topics: List[int] = None,
    difficulty_min: float = 0,
    difficulty_max: float = 100,
    companies: List[str] = None,
    similar_to: int = None,
    exclude_solved: bool = False,
    user_id: int = None,
    limit: int = 20
) -> List[Problem]:
    """
    Unified search with multiple filters
    """
    
    results = []
    
    # 1. Text search (if query provided)
    if query:
        text_results = full_text_search(query, limit=limit*2)
        results.extend(text_results)
    
    # 2. Semantic search (if query provided)
    if query:
        semantic_results = semantic_search(query, limit=limit*2)
        results.extend(semantic_results)
    
    # 3. Similar problems (if similar_to provided)
    if similar_to:
        similar_results = find_similar_problems(similar_to, limit=limit*2)
        results.extend(similar_results)
    
    # 4. Apply filters
    filtered = []
    for problem in results:
        # Difficulty filter
        if not (difficulty_min <= problem.difficulty_score <= difficulty_max):
            continue
        
        # Pattern filter
        if patterns and problem.primary_pattern not in patterns:
            continue
        
        # Topic filter
        if topics:
            problem_topics = get_problem_topics(problem.id)
            if not any(t in topics for t in problem_topics):
                continue
        
        # Company filter
        if companies:
            problem_companies = get_problem_companies(problem.id)
            if not any(c in companies for c in problem_companies):
                continue
        
        # Exclude solved
        if exclude_solved and user_id:
            if is_solved_by_user(problem.id, user_id):
                continue
        
        filtered.append(problem)
    
    # 5. Deduplicate and rank
    deduplicated = deduplicate_and_rank(filtered)
    
    return deduplicated[:limit]


def deduplicate_and_rank(problems: List[Problem]) -> List[Problem]:
    """
    Remove duplicates and rank by relevance
    """
    seen = set()
    unique = []
    
    for problem in problems:
        if problem.id not in seen:
            seen.add(problem.id)
            unique.append(problem)
    
    # Rank by multiple factors
    def rank_score(p: Problem) -> float:
        score = 0
        score += p.match_score  # From search
        score += (1 / (p.total_attempts + 1)) * 10  # Prefer less attempted
        score += (p.total_solves / max(p.total_attempts, 1)) * 5  # Prefer solvable
        return score
    
    unique.sort(key=rank_score, reverse=True)
    return unique
```

---

### **6.4 Saved Filters / Smart Lists**

**Predefined searches:**
```sql
CREATE TABLE saved_filters (
    filter_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    
    name VARCHAR(200) NOT NULL,
    description TEXT,
    
    -- Filter criteria (stored as JSON)
    criteria JSONB NOT NULL,
    
    -- Metadata
    is_public BOOLEAN DEFAULT FALSE,
    use_count INTEGER DEFAULT 0,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Example saved filter:**
```json
{
    "name": "Medium DP problems I haven't solved",
    "criteria": {
        "patterns": ["Dynamic Programming"],
        "difficulty_min": 30,
        "difficulty_max": 50,
        "exclude_solved": true
    }
}
```

---

## 7. TRAINING PLAN BUILDER

### **7.1 Plan Generation Algorithm**

```python
def generate_training_plan(
    user_id: int,
    goal: str,
    target_topics: List[int],
    target_patterns: List[str],
    duration_days: int,
    questions_per_day: int,
    difficulty_range: Tuple[float, float],
    include_reviews: bool = True
) -> TrainingPlan:
    """
    Generate personalized training plan
    """
    
    # 1. Analyze user's current level
    user_level = analyze_user_level(user_id, target_topics)
    
    # 2. Get candidate problems
    candidates = get_candidate_problems(
        topics=target_topics,
        patterns=target_patterns,
        difficulty_range=difficulty_range,
        exclude_solved=True,
        user_id=user_id
    )
    
    # 3. Build dependency graph
    dependency_graph = build_dependency_graph(candidates)
    
    # 4. Order by prerequisites (topological sort)
    ordered_problems = topological_sort(dependency_graph)
    
    # 5. Distribute across days
    total_questions = duration_days * questions_per_day
    selected_problems = select_problems(
        ordered_problems,
        count=total_questions,
        user_level=user_level
    )
    
    # 6. Add review sessions
    if include_reviews:
        selected_problems = inject_reviews(selected_problems, user_id)
    
    # 7. Create plan
    plan = create_training_plan_db(
        user_id=user_id,
        name=f"{goal} - {duration_days} Days",
        problems=selected_problems,
        questions_per_day=questions_per_day,
        start_date=datetime.now().date()
    )
    
    return plan


def select_problems(
    candidates: List[Problem],
    count: int,
    user_level: Dict[str, float]
) -> List[Problem]:
    """
    Intelligently select problems for progression
    """
    selected = []
    
    # Start easy, gradually increase
    difficulty_curve = generate_difficulty_curve(count, user_level)
    
    for target_difficulty in difficulty_curve:
        # Find problem closest to target difficulty
        best_match = min(
            candidates,
            key=lambda p: abs(p.difficulty_score - target_difficulty)
        )
        
        selected.append(best_match)
        candidates.remove(best_match)  # Don't repeat
    
    return selected


def generate_difficulty_curve(
    num_problems: int,
    user_level: Dict[str, float]
) -> List[float]:
    """
    Generate progressive difficulty curve
    
    Shape: Start at user level, gradually increase, plateau at end
    """
    start_difficulty = user_level.get('overall', 30)
    
    curve = []
    for i in range(num_problems):
        progress = i / num_problems
        
        if progress < 0.3:
            # Start: Stay at current level
            difficulty = start_difficulty
        elif progress < 0.7:
            # Middle: Linear increase
            difficulty = start_difficulty + (progress - 0.3) * 50
        else:
            # End: Plateau at challenging level
            difficulty = start_difficulty + 20
        
        curve.append(min(100, max(0, difficulty)))
    
    return curve


def inject_reviews(
    problems: List[Problem],
    user_id: int
) -> List[Problem]:
    """
    Add review sessions for spaced repetition
    """
    # Get problems due for review
    due_reviews = get_due_reviews(user_id)
    
    result = []
    review_interval = 5  # Review every 5 problems
    
    for i, problem in enumerate(problems):
        result.append(problem)
        
        # Insert review
        if (i + 1) % review_interval == 0 and due_reviews:
            review_problem = due_reviews.pop(0)
            result.append(review_problem)
    
    return result
```

---

### **7.2 Adaptive Training Plans**

**Adjust based on performance:**
```python
def adapt_training_plan(plan_id: int):
    """
    Adjust remaining questions based on performance
    """
    plan = get_training_plan(plan_id)
    
    # Calculate performance metrics
    recent_performance = calculate_recent_performance(
        plan_id,
        last_n_questions=10
    )
    
    # If performing well, increase difficulty
    if recent_performance.accuracy > 0.8:
        adjust_difficulty(plan_id, increase_by=5)
    
    # If struggling, decrease difficulty
    elif recent_performance.accuracy < 0.4:
        adjust_difficulty(plan_id, decrease_by=5)
    
    # If taking too long, simplify
    if recent_performance.avg_time > 45 * 60:  # 45 minutes
        add_easier_problems(plan_id)
    
    # If blazing through, add challenges
    elif recent_performance.avg_time < 10 * 60:  # 10 minutes
        add_harder_problems(plan_id)
```

---

### **7.3 Plan Templates**

**Pre-built plan templates:**
```sql
CREATE TABLE plan_templates (
    template_id SERIAL PRIMARY KEY,
    
    name VARCHAR(200) NOT NULL,
    description TEXT,
    
    -- Configuration
    default_duration_days INTEGER NOT NULL,
    default_questions_per_day INTEGER NOT NULL,
    
    target_level VARCHAR(50),              -- 'beginner', 'intermediate', 'advanced'
    
    -- Criteria
    topic_requirements JSONB,              -- Which topics to cover
    pattern_requirements JSONB,            -- Which patterns to cover
    difficulty_progression JSONB,          -- How difficulty changes
    
    -- Metadata
    estimated_completion_rate FLOAT,
    average_rating FLOAT,
    use_count INTEGER DEFAULT 0,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Example templates:**
1. "Two Pointers Mastery" - 14 days
2. "Dynamic Programming Bootcamp" - 30 days
3. "Graph Algorithms Sprint" - 21 days
4. "Complete Interview Prep" - 90 days

---

## 8. LLM-POWERED ASSESSMENT

### **8.1 Assessment Architecture**

```
User Answers Questions
        │
        ▼
Store Attempts in DB
        │
        ▼
Aggregate Performance Data
        │
        ▼
RAG: Retrieve Context
(Similar weak areas, solutions, explanations)
        │
        ▼
LLM Analysis
(Structured prompt with data)
        │
        ▼
Generate:
- Weakness identification
- Personalized recommendations
- Next steps
        │
        ▼
Store & Present to User
```

---

### **8.2 Assessment Prompt Engineering**

```python
def generate_assessment_prompt(session_id: int) -> str:
    """
    Build prompt for LLM assessment
    """
    # Get session data
    session = get_assessment_session(session_id)
    questions = get_assessment_questions(session_id)
    user_history = get_user_history(session.user_id, limit=100)
    
    # Build structured data
    performance_data = {
        "overall_score": session.overall_score,
        "category_scores": session.category_scores,
        "questions": []
    }
    
    for q in questions:
        performance_data["questions"].append({
            "category": q.question.question_type,
            "difficulty": q.question.difficulty_score,
            "is_correct": q.is_correct,
            "time_taken": q.time_taken_seconds,
            "user_answer": q.user_answer,
            "correct_answer": q.question.correct_answer
        })
    
    # Build prompt
    prompt = f"""You are an expert coding interview coach analyzing a student's assessment results.

# Assessment Data

## Overall Performance
- Total Score: {session.overall_score}/100
- Questions Answered: {session.total_questions}
- Correct Answers: {session.correct_answers}
- Time Taken: {(session.completed_at - session.started_at).total_seconds() / 60:.1f} minutes

## Category Breakdown
{json.dumps(session.category_scores, indent=2)}

## Detailed Question Performance
{json.dumps(performance_data["questions"], indent=2)}

## Historical Context
The student has previously attempted {len(user_history)} problems with the following patterns:
{summarize_history(user_history)}

# Your Task

Analyze this assessment and provide:

1. **Strengths** (2-3 specific areas where the student performed well)
   - Be specific about what they understand well
   - Reference specific categories or question types

2. **Weaknesses** (2-4 specific areas needing improvement)
   - Identify patterns in mistakes
   - Distinguish between knowledge gaps vs careless errors
   - Note if difficulty level correlates with performance

3. **Root Cause Analysis**
   - Why is the student struggling in weak areas?
   - Is it conceptual understanding, pattern recognition, or implementation?

4. **Personalized Recommendations** (prioritized list of 3-5 action items)
   - Specific topics to study
   - Types of problems to practice
   - Resources or approaches to try
   - Realistic timeline for improvement

5. **Next Steps** (immediate actions)
   - What should they do today?
   - What's the 7-day plan?

# Output Format

Provide your analysis in JSON format:

```json
{{
  "strengths": [
    {{
      "area": "string",
      "evidence": "string",
      "score": 0-10
    }}
  ],
  "weaknesses": [
    {{
      "area": "string",
      "severity": "low|medium|high|critical",
      "evidence": "string",
      "root_cause": "string"
    }}
  ],
  "recommendations": [
    {{
      "priority": 1-5,
      "action": "string",
      "reason": "string",
      "estimated_time": "string"
    }}
  ],
  "next_steps": {{
    "today": ["string"],
    "this_week": ["string"]
  }}
}}
```

Be honest but encouraging. Focus on actionable insights.
"""
    
    return prompt


def analyze_assessment_with_llm(session_id: int) -> Dict:
    """
    Run LLM analysis on assessment
    """
    prompt = generate_assessment_prompt(session_id)
    
    # Call LLM (Ollama)
    response = ollama.generate(
        model="mistral:7b",  # or "codellama:13b"
        prompt=prompt,
        options={
            "temperature": 0.3,  # Lower for more consistent analysis
            "top_p": 0.9
        }
    )
    
    # Parse JSON response
    try:
        analysis = json.loads(extract_json(response["response"]))
    except:
        # Fallback: Use regex to extract JSON
        analysis = extract_json_with_regex(response["response"])
    
    # Store in database
    store_assessment_analysis(session_id, analysis)
    
    return analysis


def extract_json(text: str) -> str:
    """
    Extract JSON from LLM response (may have markdown code blocks)
    """
    # Try to find JSON in markdown code blocks
    match = re.search(r'```json\s*(\{.*?\})\s*```', text, re.DOTALL)
    if match:
        return match.group(1)
    
    # Try to find raw JSON
    match = re.search(r'\{.*\}', text, re.DOTALL)
    if match:
        return match.group(0)
    
    raise ValueError("No JSON found in LLM response")
```

---

### **8.3 Preventing Memorization Bias**

**Problem:** User might memorize answers through repetition

**Solutions:**

**1. Question Variations**
```python
def get_question_for_user(
    question_id: int,
    user_id: int
) -> Dict:
    """
    Generate variation of question if user has seen it before
    """
    # Check if user has seen this before
    attempts = count_user_attempts(user_id, question_id)
    
    if attempts == 0:
        # First time, use original
        return get_question(question_id)
    
    elif attempts < 3:
        # Seen before, use slight variation
        return generate_variation(question_id, variation_level=1)
    
    else:
        # Seen many times, use significant variation
        return generate_variation(question_id, variation_level=2)


def generate_variation(question_id: int, variation_level: int) -> Dict:
    """
    Generate question variation using LLM
    """
    original_question = get_question(question_id)
    
    if variation_level == 1:
        # Change numbers, variable names, but keep structure
        prompt = f"""Generate a variation of this question:

{original_question.question_text}

Requirements:
- Keep the same concept and difficulty
- Change specific numbers, array sizes, or variable names
- Keep the answer choices similar in structure
- Maintain the same learning objective

Return the new question in the same format.
"""
    
    elif variation_level == 2:
        # Change scenario but test same concept
        prompt = f"""Generate a variation of this question:

{original_question.question_text}

Requirements:
- Test the same underlying concept
- Use a different scenario or context
- Keep similar difficulty
- Generate new answer choices
- Maintain the same learning objective

Return the new question in the same format.
"""
    
    # Generate with LLM
    response = ollama.generate(
        model="mistral:7b",
        prompt=prompt,
        options={"temperature": 0.7}
    )
    
    # Parse and return
    return parse_generated_question(response["response"])
```

**2. Time-Weighted Performance**
```python
def calculate_true_mastery(user_id: int, topic_id: int) -> float:
    """
    Calculate mastery accounting for time between attempts
    """
    attempts = get_user_attempts(user_id, topic=topic_id)
    
    mastery_score = 0
    total_weight = 0
    
    for i, attempt in enumerate(attempts):
        # Weight decreases with recency (recent memorization vs old retention)
        days_since = (datetime.now() - attempt.attempted_at).days
        
        if days_since < 1:
            weight = 0.5  # Just tried, might be memorized
        elif days_since < 7:
            weight = 0.8  # Recent memory
        elif days_since < 30:
            weight = 1.0  # Good retention
        else:
            weight = 1.2  # Excellent long-term retention
        
        # Correct = +weight, incorrect = 0
        if attempt.is_correct:
            mastery_score += weight
        
        total_weight += weight
    
    return (mastery_score / total_weight) if total_weight > 0 else 0
```

**3. Conceptual Understanding Tests**
```python
def test_conceptual_understanding(
    user_id: int,
    topic: str
) -> Dict:
    """
    Instead of asking same questions, test concept application
    """
    # Get problems user solved in this topic
    solved_problems = get_solved_problems(user_id, topic=topic)
    
    # Generate new questions that test if they understand WHY
    questions = []
    
    for problem in solved_problems:
        # Ask about trade-offs
        questions.append({
            "type": "tradeoff_analysis",
            "problem": problem,
            "question": f"Why did you choose {problem.solution_approach} over alternative X?"
        })
        
        # Ask about modifications
        questions.append({
            "type": "adaptation",
            "problem": problem,
            "question": f"How would your solution change if constraint X changed to Y?"
        })
        
        # Ask about application
        questions.append({
            "type": "transfer",
            "problem": problem,
            "question": "Identify a different problem where the same technique applies"
        })
    
    return questions
```

---

## 9. PROBLEM GENERATION SYSTEM

### **9.1 LLM-Based Problem Generation**

```python
def generate_problem_variant(
    base_problem_id: int,
    variation_type: str
) -> Dict:
    """
    Generate new problem based on existing one
    
    Variation types:
    - similar: Same pattern, different scenario
    - harder: Added constraints or complexity
    - easier: Simplified version
    - follow_up: Builds on base problem
    """
    
    base_problem = get_problem(base_problem_id)
    
    # Retrieve similar problems for context (RAG)
    similar_problems = semantic_search(
        base_problem.statement,
        collection="problems",
        k=3
    )
    
    # Build prompt
    prompt = build_generation_prompt(
        base_problem,
        variation_type,
        similar_problems
    )
    
    # Generate with LLM
    response = ollama.generate(
        model="codellama:13b",  # Better for code problems
        prompt=prompt,
        options={
            "temperature": 0.8,  # Higher for creativity
            "top_p": 0.95
        }
    )
    
    # Parse response
    generated_problem = parse_generated_problem(response["response"])
    
    # Validate
    if validate_generated_problem(generated_problem):
        return generated_problem
    else:
        # Try again with more specific prompt
        return generate_problem_variant(base_problem_id, variation_type)


def build_generation_prompt(
    base_problem: Problem,
    variation_type: str,
    context_problems: List[Dict]
) -> str:
    """
    Build prompt for problem generation
    """
    
    context = "\n\n".join([
        f"Example {i+1}: {p['document']}"
        for i, p in enumerate(context_problems)
    ])
    
    if variation_type == "similar":
        instruction = """Generate a SIMILAR problem:
- Same algorithmic pattern
- Different real-world scenario
- Similar difficulty
- Different input/output format"""
    
    elif variation_type == "harder":
        instruction = """Generate a HARDER version:
- Add an additional constraint
- Require optimization
- Increase complexity by 15-20 difficulty points
- May combine with another pattern"""
    
    elif variation_type == "easier":
        instruction = """Generate an EASIER version:
- Simplify constraints
- Reduce edge cases
- More direct approach
- Decrease complexity by 15-20 difficulty points"""
    
    elif variation_type == "follow_up":
        instruction = """Generate a FOLLOW-UP problem:
- Assumes base problem is solved
- Builds on that solution
- Tests deeper understanding
- May modify constraints or requirements"""
    
    prompt = f"""You are an expert at creating coding interview problems.

# Base Problem

Title: {base_problem.title}

Statement:
{base_problem.statement}

Constraints:
{chr(10).join(f'- {c}' for c in base_problem.constraints)}

Examples:
{format_examples(base_problem.examples)}

Primary Pattern: {base_problem.primary_pattern}
Difficulty Score: {base_problem.difficulty_score}

# Context (Similar Problems)

{context}

# Task

{instruction}

# Output Format

Provide the new problem in the following JSON format:

```json
{{
  "title": "string",
  "statement": "string (detailed problem description)",
  "constraints": ["string"],
  "examples": [
    {{
      "input": "string",
      "output": "string",
      "explanation": "string"
    }}
  ],
  "hints": ["string"],
  "primary_pattern": "string",
  "secondary_patterns": ["string"],
  "difficulty_score": number,
  "relation_to_base": "string (explain how it relates to base problem)"
}}
```

Ensure the problem:
1. Is clearly stated and unambiguous
2. Has sufficient constraints
3. Includes 2-3 examples with explanations
4. Tests the intended algorithmic concept
5. Has a well-defined optimal solution
"""
    
    return prompt


def validate_generated_problem(problem: Dict) -> bool:
    """
    Validate generated problem meets quality standards
    """
    # Required fields
    required_fields = [
        "title", "statement", "constraints", 
        "examples", "primary_pattern", "difficulty_score"
    ]
    
    for field in required_fields:
        if field not in problem or not problem[field]:
            return False
    
    # Title length
    if len(problem["title"]) < 10 or len(problem["title"]) > 200:
        return False
    
    # Statement length
    if len(problem["statement"]) < 100:
        return False
    
    # Constraints count
    if len(problem["constraints"]) < 2:
        return False
    
    # Examples count
    if len(problem["examples"]) < 2:
        return False
    
    # Example structure
    for example in problem["examples"]:
        if not all(k in example for k in ["input", "output", "explanation"]):
            return False
    
    # Difficulty range
    if not (0 <= problem["difficulty_score"] <= 100):
        return False
    
    return True
```

---

### **9.2 Automatic Graph Relationship Creation**

```python
def create_problem_relationships(
    new_problem_id: int,
    base_problem_id: int,
    relationship_type: str
):
    """
    Automatically create graph edges for generated problems
    """
    # Create in Apache AGE
    if relationship_type == "similar":
        create_similar_edge(new_problem_id, base_problem_id, score=0.9)
    
    elif relationship_type == "follow_up":
        create_follow_up_edge(new_problem_id, base_problem_id)
    
    elif relationship_type == "harder" or relationship_type == "easier":
        create_variation_edge(
            new_problem_id,
            base_problem_id,
            variation_type=relationship_type
        )
    
    # Find and create similar relationships with other problems
    find_and_create_similar_relationships(new_problem_id)


def find_and_create_similar_relationships(problem_id: int):
    """
    Use vector similarity to find and create relationships
    """
    # Get problem embedding
    problem = get_problem(problem_id)
    embedding = get_or_create_embedding(problem)
    
    # Find similar in vector DB
    similar = chroma_client.query(
        collection_name="problems",
        query_embeddings=[embedding],
        n_results=10
    )
    
    # Create edges for highly similar
    for i, similar_id in enumerate(similar['ids'][0]):
        similarity_score = 1 - similar['distances'][0][i]
        
        if similarity_score > 0.85:  # High similarity threshold
            create_similar_edge(
                problem_id,
                int(similar_id.split('_')[1]),
                score=similarity_score
            )
```

---

## 10. LOCAL DEPLOYMENT ARCHITECTURE

### **10.1 Technology Stack**

**Backend:**
- **FastAPI** - REST API (Python)
- **PostgreSQL 16** - Main database
- **Apache AGE** - Graph extension for PostgreSQL
- **pgvector** - Vector similarity in PostgreSQL
- **ChromaDB** - Dedicated vector database
- **Redis** - Caching (optional)
- **Ollama** - Local LLM runtime

**Frontend:**
- **React** - Web UI
- **TailwindCSS** - Styling
- **React Query** - State management
- **Recharts** - Visualizations

**Infrastructure:**
- **Docker Compose** - Container orchestration
- **Nginx** - Reverse proxy (optional)

---

### **10.2 Docker Compose Setup**

```yaml
version: '3.8'

services:
  # PostgreSQL with AGE and pgvector
  postgres:
    image: apache/age:PG16_latest
    container_name: leetcode_postgres
    environment:
      POSTGRES_USER: leetcode
      POSTGRES_PASSWORD: leetcode123
      POSTGRES_DB: leetcode_training
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init_db.sql:/docker-entrypoint-initdb.d/init_db.sql
    command: postgres -c shared_preload_libraries=age -c max_connections=200
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U leetcode"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Redis (optional caching)
  redis:
    image: redis:7-alpine
    container_name: leetcode_redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    command: redis-server --appendonly yes

  # ChromaDB (vector database)
  chromadb:
    image: chromadb/chroma:latest
    container_name: leetcode_chromadb
    ports:
      - "8000:8000"
    volumes:
      - chroma_data:/chroma/chroma
    environment:
      IS_PERSISTENT: "TRUE"
      ANONYMIZED_TELEMETRY: "FALSE"

  # Ollama (LLM runtime)
  ollama:
    image: ollama/ollama:latest
    container_name: leetcode_ollama
    ports:
      - "11434:11434"
    volumes:
      - ollama_data:/root/.ollama
    environment:
      OLLAMA_HOST: 0.0.0.0
    # Optional: GPU support
    # deploy:
    #   resources:
    #     reservations:
    #       devices:
    #         - driver: nvidia
    #           count: 1
    #           capabilities: [gpu]

  # Backend API
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: leetcode_backend
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgresql://leetcode:leetcode123@postgres:5432/leetcode_training
      REDIS_URL: redis://redis:6379
      CHROMA_URL: http://chromadb:8000
      OLLAMA_URL: http://ollama:11434
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_started
      chromadb:
        condition: service_started
      ollama:
        condition: service_started
    volumes:
      - ./backend:/app
    command: uvicorn main:app --host 0.0.0.0 --port 8080 --reload

  # Frontend
  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    container_name: leetcode_frontend
    ports:
      - "3000:3000"
    environment:
      REACT_APP_API_URL: http://localhost:8080
    depends_on:
      - backend
    volumes:
      - ./frontend:/app
      - /app/node_modules
    command: npm start

volumes:
  postgres_data:
  redis_data:
  chroma_data:
  ollama_data:
```

---

### **10.3 Database Initialization Script**

```sql
-- init_db.sql

-- Load extensions
CREATE EXTENSION IF NOT EXISTS age;
CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS pg_trgm;  -- For text search
CREATE EXTENSION IF NOT EXISTS btree_gin;  -- For array indexing

-- Load AGE into search path
LOAD 'age';
SET search_path = ag_catalog, "$user", public;

-- Create graph
SELECT create_graph('problem_graph');

-- Create tables (abbreviated - full schema from section 2)
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE problems (
    problem_id SERIAL PRIMARY KEY,
    leetcode_number INTEGER UNIQUE,
    problem_slug VARCHAR(200) UNIQUE NOT NULL,
    title VARCHAR(500) NOT NULL,
    statement TEXT NOT NULL,
    constraints TEXT[],
    examples JSONB NOT NULL,
    hints TEXT[],
    difficulty_score FLOAT NOT NULL,
    official_difficulty VARCHAR(20),
    primary_pattern VARCHAR(100),
    secondary_patterns VARCHAR(100)[],
    has_follow_ups BOOLEAN DEFAULT FALSE,
    has_prerequisites BOOLEAN DEFAULT FALSE,
    total_attempts INTEGER DEFAULT 0,
    total_solves INTEGER DEFAULT 0,
    average_time_seconds FLOAT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_difficulty_score CHECK (difficulty_score >= 0 AND difficulty_score <= 100)
);

-- Add full-text search
ALTER TABLE problems 
ADD COLUMN search_vector tsvector 
GENERATED ALWAYS AS (
    setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
    setweight(to_tsvector('english', coalesce(statement, '')), 'B') ||
    setweight(to_tsvector('english', array_to_string(constraints, ' ')), 'C')
) STORED;

CREATE INDEX idx_problems_search ON problems USING GIN(search_vector);

-- (Continue with other tables...)

-- Create initial data
INSERT INTO users (username, email) VALUES
('test_user', 'test@example.com');

-- Create sample topics
INSERT INTO topics (name, slug, description) VALUES
('Array', 'array', 'Array manipulation and traversal'),
('Hash Table', 'hash-table', 'Hash table operations'),
('Two Pointers', 'two-pointers', 'Two pointer technique'),
('Binary Search', 'binary-search', 'Binary search algorithm'),
('Dynamic Programming', 'dynamic-programming', 'Dynamic programming patterns');

-- Grant permissions
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO leetcode;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO leetcode;
```

---

### **10.4 Backend Structure**

```
backend/
├── main.py                 # FastAPI app
├── requirements.txt
├── Dockerfile
├── config.py              # Configuration
├── database.py            # DB connection
├── models/                # SQLAlchemy models
│   ├── __init__.py
│   ├── problem.py
│   ├── question.py
│   ├── user.py
│   └── training_plan.py
├── services/              # Business logic
│   ├── __init__.py
│   ├── problem_service.py
│   ├── question_service.py
│   ├── search_service.py
│   ├── training_plan_service.py
│   ├── assessment_service.py
│   ├── llm_service.py
│   └── vector_service.py
├── api/                   # API routes
│   ├── __init__.py
│   ├── problems.py
│   ├── questions.py
│   ├── training_plans.py
│   ├── assessments.py
│   └── search.py
├── utils/                 # Utilities
│   ├── __init__.py
│   ├── difficulty.py      # Difficulty calculation
│   ├── graph.py           # Apache AGE helpers
│   └── embeddings.py      # Vector embeddings
└── tests/
```

---

### **10.5 Startup Script**

```bash
#!/bin/bash
# start.sh

echo "🚀 Starting LeetCode Training Platform..."

# Pull latest images
echo "📥 Pulling Docker images..."
docker-compose pull

# Start services
echo "🐳 Starting Docker containers..."
docker-compose up -d

# Wait for PostgreSQL
echo "⏳ Waiting for PostgreSQL..."
until docker-compose exec -T postgres pg_isready -U leetcode; do
  sleep 1
done

# Download LLM models
echo "🤖 Downloading LLM models (this may take a while)..."
docker-compose exec -T ollama ollama pull mistral:7b
docker-compose exec -T ollama ollama pull codellama:13b

# Run migrations (if using Alembic)
echo "🗄️ Running database migrations..."
docker-compose exec -T backend alembic upgrade head

# Seed database (optional)
echo "🌱 Seeding database..."
docker-compose exec -T backend python seed_database.py

echo "✅ Platform is ready!"
echo "🌐 Frontend: http://localhost:3000"
echo "🔌 Backend API: http://localhost:8080"
echo "📊 API Docs: http://localhost:8080/docs"
```

---

## 11. API DESIGN

### **11.1 API Endpoints Overview**

```
Authentication
- POST /api/auth/register
- POST /api/auth/login
- POST /api/auth/logout
- GET  /api/auth/me

Problems
- GET    /api/problems
- GET    /api/problems/{id}
- GET    /api/problems/{id}/similar
- GET    /api/problems/{id}/follow-ups
- POST   /api/problems/{id}/attempt
- GET    /api/problems/random

Questions
- GET    /api/questions/{id}
- POST   /api/questions/{id}/answer
- GET    /api/questions/daily-challenge

Search
- GET    /api/search/problems
- GET    /api/search/semantic
- POST   /api/search/advanced

Training Plans
- GET    /api/training-plans
- POST   /api/training-plans
- GET    /api/training-plans/{id}
- PUT    /api/training-plans/{id}
- POST   /api/training-plans/{id}/start
- GET    /api/training-plans/{id}/progress
- GET    /api/training-plans/{id}/next-question

Assessments
- POST   /api/assessments/start
- GET    /api/assessments/{id}
- POST   /api/assessments/{id}/answer
- POST   /api/assessments/{id}/complete
- GET    /api/assessments/{id}/analysis

User Progress
- GET    /api/users/me/progress
- GET    /api/users/me/stats
- GET    /api/users/me/weaknesses
- GET    /api/users/me/recommendations
- GET    /api/users/me/review-queue

Analytics
- GET    /api/analytics/difficulty-distribution
- GET    /api/analytics/pattern-coverage
- GET    /api/analytics/performance-trends
```

---

### **11.2 Example API Implementations**

**Search with filters:**
```python
@router.get("/api/search/problems")
async def search_problems(
    query: Optional[str] = None,
    patterns: Optional[List[str]] = Query(None),
    topics: Optional[List[int]] = Query(None),
    difficulty_min: float = 0,
    difficulty_max: float = 100,
    companies: Optional[List[str]] = Query(None),
    similar_to: Optional[int] = None,
    exclude_solved: bool = False,
    limit: int = 20,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """
    Advanced problem search with multiple filters
    """
    results = search_service.search_problems(
        db=db,
        query=query,
        patterns=patterns,
        topics=topics,
        difficulty_min=difficulty_min,
        difficulty_max=difficulty_max,
        companies=companies,
        similar_to=similar_to,
        exclude_solved=exclude_solved,
        user_id=current_user.id if exclude_solved else None,
        limit=limit
    )
    
    return {
        "results": results,
        "count": len(results),
        "filters_applied": {
            "query": query,
            "patterns": patterns,
            "difficulty_range": [difficulty_min, difficulty_max]
        }
    }
```

**Training plan creation:**
```python
@router.post("/api/training-plans")
async def create_training_plan(
    plan_request: TrainingPlanCreate,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """
    Generate personalized training plan
    """
    plan = training_plan_service.generate_plan(
        db=db,
        user_id=current_user.id,
        goal=plan_request.goal,
        target_topics=plan_request.target_topics,
        target_patterns=plan_request.target_patterns,
        duration_days=plan_request.duration_days,
        questions_per_day=plan_request.questions_per_day,
        difficulty_range=(plan_request.difficulty_min, plan_request.difficulty_max),
        include_reviews=plan_request.include_reviews
    )
    
    return plan
```

**Assessment with LLM analysis:**
```python
@router.post("/api/assessments/{session_id}/complete")
async def complete_assessment(
    session_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db),
    background_tasks: BackgroundTasks
):
    """
    Complete assessment and trigger LLM analysis
    """
    # Mark as completed
    session = assessment_service.complete_session(
        db=db,
        session_id=session_id,
        user_id=current_user.id
    )
    
    # Trigger LLM analysis in background
    background_tasks.add_task(
        assessment_service.analyze_with_llm,
        session_id=session_id
    )
    
    return {
        "session_id": session_id,
        "status": "completed",
        "overall_score": session.overall_score,
        "message": "Analysis in progress. Check back in a moment."
    }


@router.get("/api/assessments/{session_id}/analysis")
async def get_assessment_analysis(
    session_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """
    Get LLM analysis results
    """
    session = assessment_service.get_session(db, session_id, current_user.id)
    
    if not session.analysis:
        return {"status": "pending", "message": "Analysis still in progress"}
    
    return {
        "status": "complete",
        "analysis": session.analysis,
        "recommendations": session.recommendations,
        "weaknesses": session.weaknesses
    }
```

---

## 12. IMPLEMENTATION ROADMAP

### **Phase 1: Foundation (Weeks 1-2)**

**Goals:**
- Basic database schema
- Problem and question storage
- Simple API

**Tasks:**
- [ ] Set up PostgreSQL + Apache AGE
- [ ] Create core tables (problems, questions, users)
- [ ] Implement basic CRUD operations
- [ ] Create 50 sample problems
- [ ] Create 200 sample questions
- [ ] Basic search functionality

**Deliverable:** Can store and retrieve problems/questions

---

### **Phase 2: Vector & Graph (Weeks 3-4)**

**Goals:**
- Vector database integration
- Graph relationships
- Semantic search

**Tasks:**
- [ ] Set up ChromaDB
- [ ] Generate embeddings for problems
- [ ] Create graph nodes and edges
- [ ] Implement semantic search
- [ ] Create similar problem finder
- [ ] Build prerequisite chains

**Deliverable:** Smart problem recommendations work

---

### **Phase 3: Training Plans (Weeks 5-6)**

**Goals:**
- Training plan generation
- Progress tracking
- Spaced repetition

**Tasks:**
- [ ] Training plan algorithm
- [ ] User progress tracking
- [ ] Adaptive difficulty
- [ ] Review scheduling
- [ ] Plan templates

**Deliverable:** Users can follow personalized training plans

---

### **Phase 4: LLM Integration (Weeks 7-8)**

**Goals:**
- Local LLM setup
- RAG implementation
- Assessment analysis

**Tasks:**
- [ ] Set up Ollama
- [ ] Download models
- [ ] Implement RAG queries
- [ ] Build assessment prompts
- [ ] LLM weakness detection

**Deliverable:** Automated weakness analysis works

---

### **Phase 5: Frontend (Weeks 9-10)**

**Goals:**
- User interface
- Visualizations
- User experience

**Tasks:**
- [ ] React app setup
- [ ] Problem browser
- [ ] Question interface
- [ ] Training plan UI
- [ ] Progress dashboard

**Deliverable:** Complete working application

---

### **Phase 6: Polish & Deploy (Weeks 11-12)**

**Goals:**
- Bug fixes
- Performance optimization
- Documentation

**Tasks:**
- [ ] Load testing
- [ ] Optimize queries
- [ ] Add caching
- [ ] Write documentation
- [ ] User testing

**Deliverable:** Production-ready local application

---

## SUMMARY

This technical architecture provides:

✅ **Complete database schema** (PostgreSQL + Apache AGE)
✅ **Graph relationships** for problem connections
✅ **Vector database** for semantic search and RAG
✅ **Custom difficulty scoring** (Magic Unit)
✅ **Smart search & filtering** (multi-modal)
✅ **Flexible training plans** (adaptive)
✅ **LLM-powered assessment** (weakness detection)
✅ **Problem generation** (variations & follow-ups)
✅ **Local deployment** (Docker Compose)
✅ **Complete API** design

**Next Steps:**
1. Set up local environment
2. Create database schema
3. Import initial problem set
4. Build core API
5. Test with real usage

The system is designed to run entirely locally, giving you full control and privacy while practicing for interviews.

---

*Technical Architecture Document v1.0*
*Ready for implementation*
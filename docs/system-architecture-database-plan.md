# Interview Training Platform - System Architecture & Database Design
## Complete Technical Blueprint for Local Deployment

---

## DOCUMENT PURPOSE

This is the **complete technical architecture plan** for building a locally-hosted interview training platform with:
- PostgreSQL + Apache AGE (graph database)
- Vector database + RAG for LLM integration
- Intelligent difficulty scoring
- Fair assessment system
- Weakness detection
- Dynamic training plans
- Question generation

**This document defines WHAT to build and HOW to architect it.**

---

# TABLE OF CONTENTS

1. [System Architecture Overview](#system-architecture-overview)
2. [Database Schema Design](#database-schema-design)
3. [Graph Database Structure (Apache AGE)](#graph-database-structure)
4. [Vector Database Integration](#vector-database-integration)
5. [Difficulty Scoring System](#difficulty-scoring-system)
6. [Assessment & Evaluation Framework](#assessment-framework)
7. [Weakness Detection Algorithm](#weakness-detection)
8. [Search & Filter System](#search-filter-system)
9. [Training Plan Builder](#training-plan-builder)
10. [LLM Integration Architecture](#llm-integration)
11. [Question Generation Pipeline](#question-generation)
12. [Local Deployment Setup](#local-deployment)

---

## 1. SYSTEM ARCHITECTURE OVERVIEW

### **Technology Stack**

```
┌─────────────────────────────────────────────────┐
│                  Frontend (Local)                │
│            React/Vue + TailwindCSS               │
└─────────────────┬───────────────────────────────┘
                  │
┌─────────────────▼───────────────────────────────┐
│              Backend API (Local)                 │
│        FastAPI / Node.js + Express               │
└──┬──────────────┬────────────────┬──────────────┘
   │              │                │
   ▼              ▼                ▼
┌──────────┐  ┌──────────┐  ┌─────────────┐
│PostgreSQL│  │ Vector DB│  │   LLM       │
│+ AGE     │  │ (Chroma/ │  │(Ollama/     │
│(Graph)   │  │ Qdrant)  │  │ LlamaCPP)   │
└──────────┘  └──────────┘  └─────────────┘
```

### **Core Components**

**1. Database Layer**
- PostgreSQL: Main relational data
- Apache AGE: Graph relationships (problems, topics, dependencies)
- Vector DB: Embeddings for semantic search + RAG

**2. Backend Services**
- API Server: REST/GraphQL endpoints
- LLM Service: Question generation, assessment
- Analytics Engine: Progress tracking, weakness detection
- Training Plan Engine: Dynamic curriculum generation

**3. Frontend**
- Practice Interface: Question solving
- Dashboard: Progress visualization
- Analytics: Weakness reports
- Training Builder: Custom plans

---

## 2. DATABASE SCHEMA DESIGN (PostgreSQL)

### **Core Tables**

#### **Table: problems**
```sql
CREATE TABLE problems (
    problem_id SERIAL PRIMARY KEY,
    
    -- Identification
    leetcode_number INTEGER UNIQUE,     -- LeetCode problem number
    external_id VARCHAR(50),            -- Other platform IDs
    title VARCHAR(200) NOT NULL,
    slug VARCHAR(200) UNIQUE NOT NULL,
    
    -- Content
    description TEXT NOT NULL,
    constraints TEXT,
    examples JSONB,                     -- [{input, output, explanation}]
    hints JSONB,                        -- [hint1, hint2, ...]
    
    -- Classification
    difficulty_score FLOAT NOT NULL,    -- Our calculated score (0-100)
    official_difficulty VARCHAR(20),    -- Easy/Medium/Hard
    category VARCHAR(50),               -- From our 10 categories
    subcategory VARCHAR(50),            -- Specific type
    
    -- Solution data
    solution_template TEXT,             -- Code template
    solution_pseudocode TEXT,           -- High-level algorithm
    time_complexity VARCHAR(50),
    space_complexity VARCHAR(50),
    
    -- Metadata
    frequency FLOAT DEFAULT 0,          -- How often asked (0-1)
    success_rate FLOAT,                 -- Global success rate
    avg_solve_time INTEGER,             -- Seconds
    companies JSONB,                    -- [company names]
    tags JSONB,                         -- [tag1, tag2, ...]
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    -- Embedding reference
    embedding_id VARCHAR(100)           -- Reference to vector DB
);

CREATE INDEX idx_problems_difficulty ON problems(difficulty_score);
CREATE INDEX idx_problems_category ON problems(category);
CREATE INDEX idx_problems_leetcode ON problems(leetcode_number);
CREATE INDEX idx_problems_tags ON problems USING GIN(tags);
```

---

#### **Table: question_types**
```sql
CREATE TABLE question_types (
    type_id SERIAL PRIMARY KEY,
    type_name VARCHAR(100) NOT NULL,           -- "Complexity Analysis"
    subtype_name VARCHAR(100),                 -- "Code-to-Complexity"
    category VARCHAR(50),                      -- From taxonomy
    
    description TEXT,
    structure JSONB,                           -- Template structure
    learning_objective TEXT,
    
    difficulty_weights JSONB                   -- {easy: 1, medium: 2, hard: 5}
);
```

---

#### **Table: questions**
```sql
CREATE TABLE questions (
    question_id SERIAL PRIMARY KEY,
    
    -- Type classification
    type_id INTEGER REFERENCES question_types(type_id),
    problem_id INTEGER REFERENCES problems(problem_id),  -- NULL if standalone
    
    -- Content
    question_text TEXT NOT NULL,
    question_format VARCHAR(50),               -- 'multiple_choice', 'code', 'pseudocode'
    
    -- For multiple choice
    options JSONB,                             -- [{id: 'A', text: '...', correct: bool}]
    
    -- For open-ended
    rubric JSONB,                              -- Grading criteria
    test_cases JSONB,                          -- For code questions
    
    -- Answer & explanation
    correct_answer TEXT,                       -- 'A' or full answer
    explanation TEXT NOT NULL,
    hints JSONB,
    
    -- Difficulty & scoring
    difficulty_score FLOAT NOT NULL,
    expected_time INTEGER,                     -- Seconds
    points_base INTEGER DEFAULT 10,
    
    -- Learning
    concepts_tested JSONB,                     -- [concept1, concept2]
    prerequisites JSONB,                       -- [question_ids]
    follow_up_questions JSONB,                 -- [question_ids]
    
    -- Metadata
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    -- Stats
    times_attempted INTEGER DEFAULT 0,
    times_correct INTEGER DEFAULT 0,
    avg_time FLOAT,
    
    -- Embedding
    embedding_id VARCHAR(100)
);

CREATE INDEX idx_questions_type ON questions(type_id);
CREATE INDEX idx_questions_problem ON questions(problem_id);
CREATE INDEX idx_questions_difficulty ON questions(difficulty_score);
```

---

#### **Table: topics**
```sql
CREATE TABLE topics (
    topic_id SERIAL PRIMARY KEY,
    topic_name VARCHAR(100) NOT NULL UNIQUE,  -- "Dynamic Programming"
    parent_topic_id INTEGER REFERENCES topics(topic_id),
    
    description TEXT,
    difficulty_level INTEGER,                  -- 1-5
    order_index INTEGER,                       -- Learning order
    
    -- Prerequisites
    prerequisite_topics JSONB,                 -- [topic_ids]
    
    -- Resources
    learning_resources JSONB,                  -- [{type: 'video', url: '...'}]
    
    metadata JSONB
);

CREATE INDEX idx_topics_parent ON topics(parent_topic_id);
```

---

#### **Table: problem_topics** (Many-to-Many)
```sql
CREATE TABLE problem_topics (
    problem_id INTEGER REFERENCES problems(problem_id),
    topic_id INTEGER REFERENCES topics(topic_id),
    
    relevance_score FLOAT DEFAULT 1.0,         -- How relevant is this topic
    is_primary BOOLEAN DEFAULT FALSE,          -- Primary topic or secondary
    
    PRIMARY KEY (problem_id, topic_id)
);

CREATE INDEX idx_pt_problem ON problem_topics(problem_id);
CREATE INDEX idx_pt_topic ON problem_topics(topic_id);
```

---

#### **Table: templates**
```sql
CREATE TABLE templates (
    template_id SERIAL PRIMARY KEY,
    template_name VARCHAR(100) NOT NULL,       -- "BFS Template"
    
    category VARCHAR(50),                      -- "Graph Traversal"
    description TEXT,
    
    code_template TEXT NOT NULL,
    pseudocode TEXT,
    
    when_to_use TEXT,                          -- Explanation
    time_complexity VARCHAR(50),
    space_complexity VARCHAR(50),
    
    variations JSONB,                          -- [{name, code, when}]
    related_problems JSONB,                    -- [problem_ids]
    
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

#### **Table: user_progress**
```sql
CREATE TABLE user_progress (
    progress_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    question_id INTEGER REFERENCES questions(question_id),
    
    -- Attempt data
    attempt_number INTEGER DEFAULT 1,
    is_correct BOOLEAN NOT NULL,
    time_taken INTEGER,                        -- Seconds
    
    -- User's answer
    user_answer TEXT,
    user_code TEXT,                            -- If code question
    
    -- Scoring
    points_earned INTEGER,
    speed_bonus INTEGER,
    
    -- Analysis
    mistakes JSONB,                            -- [{type, description}]
    hints_used JSONB,                          -- [hint_ids]
    
    attempted_at TIMESTAMP DEFAULT NOW(),
    
    -- Spaced repetition
    next_review_at TIMESTAMP,
    ease_factor FLOAT DEFAULT 2.5,
    interval_days INTEGER DEFAULT 1
);

CREATE INDEX idx_progress_user ON user_progress(user_id);
CREATE INDEX idx_progress_question ON user_progress(question_id);
CREATE INDEX idx_progress_review ON user_progress(next_review_at);
```

---

#### **Table: user_skills**
```sql
CREATE TABLE user_skills (
    skill_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    
    skill_name VARCHAR(100) NOT NULL,          -- "Complexity Analysis"
    skill_category VARCHAR(50),                -- Category from taxonomy
    
    -- Proficiency metrics
    level INTEGER DEFAULT 1,                   -- 1-50
    xp INTEGER DEFAULT 0,
    mastery_score FLOAT DEFAULT 0,             -- 0-100
    
    -- Stats
    questions_attempted INTEGER DEFAULT 0,
    questions_correct INTEGER DEFAULT 0,
    avg_accuracy FLOAT,
    avg_time FLOAT,
    
    -- Trends
    improvement_rate FLOAT,                    -- Week over week
    last_practiced_at TIMESTAMP,
    streak_days INTEGER DEFAULT 0,
    
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(user_id, skill_name)
);

CREATE INDEX idx_skills_user ON user_skills(user_id);
CREATE INDEX idx_skills_mastery ON user_skills(mastery_score);
```

---

#### **Table: training_plans**
```sql
CREATE TABLE training_plans (
    plan_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    
    plan_name VARCHAR(200) NOT NULL,
    description TEXT,
    
    -- Structure
    plan_type VARCHAR(50),                     -- 'preset', 'custom', 'ai_generated'
    duration_days INTEGER,
    difficulty_range JSONB,                    -- {min, max}
    
    -- Goals
    target_skills JSONB,                       -- [skill names]
    target_topics JSONB,                       -- [topic_ids]
    
    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    progress_percentage FLOAT DEFAULT 0,
    
    created_at TIMESTAMP DEFAULT NOW(),
    started_at TIMESTAMP,
    completed_at TIMESTAMP
);
```

---

#### **Table: training_plan_items**
```sql
CREATE TABLE training_plan_items (
    item_id SERIAL PRIMARY KEY,
    plan_id INTEGER REFERENCES training_plans(plan_id),
    
    -- Ordering
    day_number INTEGER,
    order_index INTEGER,
    
    -- Content
    item_type VARCHAR(50),                     -- 'question', 'problem', 'review'
    question_id INTEGER REFERENCES questions(question_id),
    problem_id INTEGER REFERENCES problems(problem_id),
    
    -- Requirements
    is_required BOOLEAN DEFAULT TRUE,
    estimated_time INTEGER,
    
    -- Status
    is_completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP,
    
    -- Adaptive
    can_skip BOOLEAN DEFAULT FALSE,
    unlock_conditions JSONB                    -- [{type, requirement}]
);

CREATE INDEX idx_plan_items_plan ON training_plan_items(plan_id);
```

---

#### **Table: assessments**
```sql
CREATE TABLE assessments (
    assessment_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    
    assessment_type VARCHAR(50),               -- 'diagnostic', 'progress', 'final'
    
    -- Questions included
    question_ids JSONB NOT NULL,               -- [question_ids]
    
    -- Scoring
    total_points INTEGER,
    earned_points INTEGER,
    
    -- LLM evaluation
    llm_feedback TEXT,
    strengths JSONB,                           -- [strength descriptions]
    weaknesses JSONB,                          -- [weakness descriptions]
    recommendations JSONB,                     -- [recommendation objects]
    
    -- Bias detection
    memorization_score FLOAT,                  -- 0-1 (1 = likely memorized)
    pattern_recognition_score FLOAT,           -- 0-1
    problem_solving_score FLOAT,               -- 0-1
    
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    time_taken INTEGER                         -- Seconds
);
```

---

#### **Table: weakness_analysis**
```sql
CREATE TABLE weakness_analysis (
    analysis_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    
    weakness_type VARCHAR(100),                -- "Complexity Analysis", "Edge Cases"
    severity VARCHAR(20),                      -- "critical", "major", "minor"
    
    -- Evidence
    evidence_question_ids JSONB,               -- [question_ids where weakness shown]
    pattern_description TEXT,
    
    -- Metrics
    error_frequency FLOAT,                     -- How often this mistake
    impact_score FLOAT,                        -- How much it hurts overall
    
    -- Recommendations
    recommended_practice JSONB,                -- [{type, question_ids}]
    prerequisite_topics JSONB,                 -- [topic_ids to review]
    
    detected_at TIMESTAMP DEFAULT NOW(),
    resolved_at TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE
);

CREATE INDEX idx_weakness_user ON weakness_analysis(user_id, is_active);
```

---

#### **Table: llm_generations**
```sql
CREATE TABLE llm_generations (
    generation_id SERIAL PRIMARY KEY,
    
    generation_type VARCHAR(50),               -- 'question', 'explanation', 'assessment'
    
    -- Input
    prompt_template VARCHAR(100),
    prompt_variables JSONB,
    context_documents JSONB,                   -- RAG context used
    
    -- Output
    generated_content TEXT,
    
    -- Quality
    human_rating INTEGER,                      -- 1-5
    is_approved BOOLEAN DEFAULT FALSE,
    
    -- Metadata
    model_name VARCHAR(50),
    tokens_used INTEGER,
    generation_time FLOAT,                     -- Seconds
    
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

## 3. GRAPH DATABASE STRUCTURE (Apache AGE)

### **Why Graph Database?**

Problems and topics have **natural graph relationships**:
- Similar problems
- Follow-up problems
- Problem → Topic dependencies
- Topic → Topic prerequisites
- Learning paths

### **Graph Schema in AGE**

```cypher
-- Problem nodes
CREATE (p:Problem {
    problem_id: 123,
    title: "Two Sum",
    difficulty_score: 15.5,
    category: "Hash Table"
})

-- Topic nodes
CREATE (t:Topic {
    topic_id: 1,
    name: "Hash Tables",
    difficulty_level: 2
})

-- Template nodes
CREATE (tm:Template {
    template_id: 5,
    name: "Hash Map Pattern"
})
```

### **Relationships**

#### **Problem Relationships**

```cypher
-- Similar problems (bidirectional)
CREATE (p1:Problem)-[:SIMILAR_TO {similarity_score: 0.85}]->(p2:Problem)

-- Follow-up (directional)
CREATE (p1:Problem)-[:FOLLOW_UP {difficulty_increase: 20}]->(p2:Problem)

-- Prerequisite
CREATE (p1:Problem)-[:REQUIRES {importance: 0.9}]->(p2:Problem)

-- Variation
CREATE (p1:Problem)-[:VARIATION_OF {type: "constraint_change"}]->(p2:Problem)

-- Uses template
CREATE (p:Problem)-[:USES_TEMPLATE {applicability: 1.0}]->(tm:Template)
```

#### **Topic Relationships**

```cypher
-- Topic hierarchy
CREATE (parent:Topic)-[:CONTAINS]->(child:Topic)

-- Prerequisites
CREATE (t1:Topic)-[:PREREQUISITE_FOR {strength: 0.8}]->(t2:Topic)

-- Related topics
CREATE (t1:Topic)-[:RELATED_TO {relevance: 0.7}]->(t2:Topic)

-- Problem belongs to topic
CREATE (p:Problem)-[:TAGGED_WITH {relevance: 0.9, is_primary: true}]->(t:Topic)
```

#### **Learning Path Relationships**

```cypher
-- Recommended progression
CREATE (p1:Problem)-[:LEADS_TO {success_rate: 0.75}]->(p2:Problem)

-- Difficulty progression
CREATE (p1:Problem)-[:PROGRESSES_TO {difficulty_step: 5}]->(p2:Problem)
```

### **Graph Queries Examples**

#### **Find Similar Problems**
```cypher
MATCH (p1:Problem {problem_id: 123})-[r:SIMILAR_TO]->(p2:Problem)
WHERE r.similarity_score > 0.7
RETURN p2.problem_id, p2.title, r.similarity_score
ORDER BY r.similarity_score DESC
LIMIT 10
```

#### **Get Learning Path**
```cypher
MATCH path = (start:Problem {problem_id: 1})-[:LEADS_TO*1..5]->(end:Problem)
WHERE end.difficulty_score <= 50
RETURN path
ORDER BY length(path)
LIMIT 1
```

#### **Find All Prerequisites**
```cypher
MATCH (p:Problem {problem_id: 456})-[:REQUIRES*]->(prereq:Problem)
RETURN DISTINCT prereq.problem_id, prereq.title
ORDER BY prereq.difficulty_score
```

#### **Topic Dependency Tree**
```cypher
MATCH path = (root:Topic {name: "Dynamic Programming"})-[:PREREQUISITE_FOR*]->(dependent:Topic)
RETURN path
```

---

## 4. VECTOR DATABASE INTEGRATION

### **Purpose**

Vector DB enables:
- Semantic search ("find problems about sliding windows")
- RAG for LLM (retrieve relevant context)
- Similarity detection (duplicate/similar problems)
- Question generation (find similar problems to base new questions on)

### **Architecture**

```
┌────────────────────────────────────────┐
│         Vector Database                │
│      (Chroma / Qdrant / Weaviate)     │
├────────────────────────────────────────┤
│                                        │
│  Collections:                          │
│  1. problem_descriptions               │
│  2. question_texts                     │
│  3. topic_descriptions                 │
│  4. solution_explanations              │
│                                        │
└────────────────────────────────────────┘
```

### **Embedding Strategy**

**Model Choice (Local):**
- `all-MiniLM-L6-v2` (384 dimensions, fast, good quality)
- `bge-large-en-v1.5` (1024 dimensions, better quality, slower)
- Store in: ChromaDB (simplest for local) or Qdrant (production-ready)

### **Vector Collections**

#### **Collection: problem_embeddings**
```python
{
    "id": "problem_123",
    "embedding": [0.123, -0.456, ...],  # 384 or 1024 dimensions
    "metadata": {
        "problem_id": 123,
        "title": "Two Sum",
        "category": "Hash Table",
        "difficulty_score": 15.5,
        "tags": ["array", "hash_table"],
        "text": "full problem description..."
    }
}
```

#### **Collection: question_embeddings**
```python
{
    "id": "question_5678",
    "embedding": [...],
    "metadata": {
        "question_id": 5678,
        "type": "Complexity Analysis",
        "problem_id": 123,
        "difficulty_score": 20.0,
        "text": "What is the time complexity of..."
    }
}
```

#### **Collection: topic_embeddings**
```python
{
    "id": "topic_42",
    "embedding": [...],
    "metadata": {
        "topic_id": 42,
        "name": "Dynamic Programming",
        "description": "full description...",
        "difficulty_level": 4
    }
}
```

### **Search Operations**

#### **Semantic Search**
```python
# User query: "problems about finding paths in graphs"
query_embedding = embed_text("problems about finding paths in graphs")

results = vector_db.query(
    query_embeddings=[query_embedding],
    n_results=10,
    where={"difficulty_score": {"$lte": 50}}  # Filter by difficulty
)

# Returns: Graph traversal problems (DFS, BFS, Dijkstra, etc.)
```

#### **Similarity Search**
```python
# Find problems similar to problem #123
problem_embedding = get_embedding_for_problem(123)

similar = vector_db.query(
    query_embeddings=[problem_embedding],
    n_results=10,
    where={"problem_id": {"$ne": 123}}  # Exclude self
)
```

#### **RAG Context Retrieval**
```python
# For LLM: Generate explanation for problem
problem_text = get_problem_description(problem_id)
query_embedding = embed_text(problem_text)

# Get relevant context
context_docs = vector_db.query(
    query_embeddings=[query_embedding],
    n_results=5,
    where_document={"$contains": "explanation"}
)

# Use in LLM prompt
prompt = f"""
Context from similar problems:
{context_docs}

Now explain this problem:
{problem_text}
"""
```

---

## 5. DIFFICULTY SCORING SYSTEM

### **The Magic Unit: Difficulty Score (0-100)**

Traditional "Easy/Medium/Hard" is too coarse. We need granular scoring.

### **Scoring Formula Components**

#### **1. Base Complexity Score (0-40 points)**

Based on algorithmic complexity required:

```python
complexity_scores = {
    "O(1)": 5,
    "O(log n)": 10,
    "O(n)": 15,
    "O(n log n)": 20,
    "O(n²)": 25,
    "O(n³)": 30,
    "O(2^n)": 35,
    "O(n!)": 40
}
```

#### **2. Concept Count (0-20 points)**

How many concepts needed:

```python
concepts_needed = len(problem.concepts_tested)
concept_score = min(concepts_needed * 4, 20)
```

#### **3. Implementation Difficulty (0-20 points)**

Factors:
- Edge cases to handle: +2 each (max 10)
- Data structures needed: +3 each (max 10)

```python
edge_cases = count_edge_cases(problem)
data_structures = count_unique_ds(problem)

implementation_score = min(
    edge_cases * 2 + data_structures * 3,
    20
)
```

#### **4. Pattern Recognition (0-10 points)**

How obvious is the pattern?

```python
if problem.pattern_obvious:
    pattern_score = 2
elif problem.pattern_hidden:
    pattern_score = 6
else:  # Novel pattern
    pattern_score = 10
```

#### **5. Historical Performance (0-10 points)**

Based on actual user data:

```python
global_success_rate = problem.success_rate

if global_success_rate > 0.7:
    performance_score = 2
elif global_success_rate > 0.5:
    performance_score = 5
elif global_success_rate > 0.3:
    performance_score = 7
else:
    performance_score = 10
```

### **Final Score Calculation**

```python
def calculate_difficulty_score(problem):
    base_complexity = get_complexity_score(problem.time_complexity)
    concept_count = min(len(problem.concepts_tested) * 4, 20)
    implementation = calculate_implementation_difficulty(problem)
    pattern_recognition = get_pattern_score(problem)
    historical = get_historical_score(problem.success_rate)
    
    raw_score = (
        base_complexity * 0.4 +
        concept_count * 0.2 +
        implementation * 0.2 +
        pattern_recognition * 0.1 +
        historical * 0.1
    )
    
    # Normalize to 0-100
    return round(raw_score, 1)
```

### **Score Ranges**

```
0-15:   Very Easy (warmup)
16-30:  Easy (beginner friendly)
31-45:  Medium-Easy (requires thought)
46-60:  Medium (standard interview)
61-75:  Medium-Hard (challenging)
76-90:  Hard (advanced)
91-100: Expert (rare/research level)
```

### **Dynamic Adjustment**

Score adjusts based on actual performance:

```python
def update_difficulty_score(problem_id):
    recent_attempts = get_recent_attempts(problem_id, days=30)
    
    success_rate = sum(a.is_correct for a in recent_attempts) / len(recent_attempts)
    avg_time = mean(a.time_taken for a in recent_attempts)
    expected_time = problem.expected_time
    
    # Adjust based on performance
    if success_rate > 0.8 and avg_time < expected_time:
        # Too easy, increase by 5%
        new_score = problem.difficulty_score * 1.05
    elif success_rate < 0.3:
        # Too hard, decrease by 5%
        new_score = problem.difficulty_score * 0.95
    else:
        new_score = problem.difficulty_score
    
    return clamp(new_score, 0, 100)
```

---

## 6. ASSESSMENT & EVALUATION FRAMEWORK

### **Goals**

1. **Fair Assessment** - Not biased toward memorization
2. **Detect Weaknesses** - Identify specific gaps
3. **LLM-Powered** - Deep analysis beyond right/wrong
4. **Adaptive** - Adjust to user's level

### **Assessment Types**

#### **Type 1: Diagnostic Assessment**

**Purpose**: Determine starting level and weaknesses

**Structure**:
- 20-30 questions
- Covers all 10 categories
- Adaptive difficulty (adjusts based on performance)
- Time: 45-60 minutes

**Scoring**:
```python
diagnostic_score = {
    "overall": 0-100,
    "per_category": {
        "Complexity Analysis": 0-100,
        "Data Structure Selection": 0-100,
        # ... for all 10
    },
    "strengths": ["Pattern Recognition", "STL Knowledge"],
    "weaknesses": ["Complexity Analysis", "Edge Cases"],
    "recommended_level": "Medium (45-60 difficulty)"
}
```

#### **Type 2: Progress Assessment**

**Purpose**: Check improvement over time

**Structure**:
- 10-15 questions
- Focuses on recently practiced areas
- Time: 20-30 minutes
- Frequency: Weekly

**Metrics**:
```python
progress_metrics = {
    "accuracy_trend": "improving",  # improving/stable/declining
    "speed_trend": "improving",
    "difficulty_handled": 55,  # Can now handle up to 55
    "weak_areas_resolved": ["Edge Cases"],  # Fixed
    "new_weak_areas": ["Greedy vs DP"],  # Emerged
}
```

#### **Type 3: Mock Interview**

**Purpose**: Simulate real interview

**Structure**:
- 2-3 full problems
- 45 minutes per problem
- Cannot pause
- Real-time evaluation

**Evaluation**:
```python
interview_evaluation = {
    "technical_score": 0-100,
    "communication_score": 0-100,  # If recording explanations
    "problem_solving_approach": "structured/unstructured",
    "time_management": "excellent/good/poor",
    "code_quality": 0-100,
    "edge_case_handling": 0-100
}
```

### **Anti-Memorization Measures**

#### **1. Problem Variants**

For each problem, generate variants that test understanding:

```python
# Original: Two Sum with array
# Variant 1: Two Sum with linked list
# Variant 2: Two Sum with BST
# Variant 3: K Sum (generalization)

def generate_variant(original_problem):
    variants = {
        "data_structure_change": change_input_ds(original_problem),
        "constraint_change": modify_constraints(original_problem),
        "generalization": generalize_problem(original_problem),
        "specialization": add_constraints(original_problem)
    }
    return random.choice(variants)
```

#### **2. Question Rotation**

Never show same question twice in short period:

```python
def select_question(user_id, category):
    recently_seen = get_recently_seen_questions(user_id, days=30)
    
    candidates = get_questions_by_category(category)
    candidates = [q for q in candidates if q.id not in recently_seen]
    
    # If all seen, use least recently seen
    if not candidates:
        candidates = get_questions_by_category(category)
        candidates.sort(key=lambda q: last_seen_date(user_id, q.id))
    
    return select_by_difficulty(candidates, user_level)
```

#### **3. Explanation Requirement**

For assessment questions, require explanation:

```python
assessment_question = {
    "question": "...",
    "answer": "multiple_choice_or_code",
    "explanation_required": True,
    "explanation_prompt": "Explain WHY you chose this approach"
}

# LLM evaluates explanation quality
explanation_score = llm_evaluate_explanation(
    question=question,
    user_answer=answer,
    user_explanation=explanation
)

# Score factors in explanation
total_score = (
    0.6 * correctness_score +
    0.4 * explanation_score
)
```

#### **4. Transfer Questions**

Test ability to apply learned patterns to new scenarios:

```python
# User solved: "Find duplicates in array"
# Transfer question: "Find duplicates in stream of integers"

def generate_transfer_question(mastered_problem):
    pattern = extract_pattern(mastered_problem)
    new_context = get_different_context()
    
    transfer_question = apply_pattern_to_context(pattern, new_context)
    return transfer_question
```

### **LLM-Powered Evaluation**

#### **Evaluation Prompt Template**

```python
evaluation_prompt = """
You are an expert coding interview evaluator.

PROBLEM:
{problem_description}

USER'S SOLUTION:
{user_code}

USER'S EXPLANATION:
{user_explanation}

Evaluate on these dimensions:

1. CORRECTNESS (0-100):
   - Does the solution work?
   - Are edge cases handled?
   - Is the logic sound?

2. APPROACH (0-100):
   - Is this the optimal approach?
   - Are there better alternatives?
   - Did they choose the right data structures?

3. UNDERSTANDING (0-100):
   - Does their explanation show deep understanding?
   - Can they explain the complexity?
   - Do they understand WHY this works?

4. MEMORIZATION INDICATOR (0-1):
   - 0 = Clearly understands, not memorized
   - 1 = Likely memorized solution without understanding
   
   Signs of memorization:
   - Can't explain why approach works
   - Can't modify for slight variations
   - Perfect code but poor explanation
   - Fast submission but shallow understanding

Provide your evaluation in JSON format:
{{
    "correctness": 85,
    "approach": 90,
    "understanding": 70,
    "memorization_score": 0.3,
    "strengths": ["chose optimal DS", "clean code"],
    "weaknesses": ["missed edge case X", "could optimize space"],
    "feedback": "detailed constructive feedback...",
    "is_memorized": false
}}
"""

llm_response = llm.generate(evaluation_prompt)
evaluation = json.loads(llm_response)
```

#### **Weakness Pattern Detection**

LLM analyzes multiple attempts to find patterns:

```python
weakness_detection_prompt = """
Analyze this user's last 20 attempts across different problems:

{attempt_history}

Identify patterns in their mistakes:

1. What types of errors repeat?
   - Edge case misses (which types?)
   - Complexity miscalculations
   - Wrong data structure choices
   - Algorithm selection errors

2. What concepts seem weak?
   - List specific concepts
   - Evidence for each

3. What's the root cause?
   - Lack of understanding of X
   - Rushing through Y
   - Not practicing Z enough

Return JSON:
{{
    "weakness_patterns": [
        {{
            "type": "edge_cases",
            "specific": "empty_input_handling",
            "frequency": 0.6,
            "severity": "major",
            "evidence": [question_ids]
        }}
    ],
    "weak_concepts": ["dynamic_programming_state_definition"],
    "root_causes": ["insufficient_practice_with_2D_DP"],
    "recommendations": [
        {{
            "action": "practice",
            "target": "2D_DP_problems",
            "priority": "high",
            "estimated_time": "5_hours"
        }}
    ]
}}
"""
```

---

## 7. WEAKNESS DETECTION ALGORITHM

### **Multi-Level Detection**

#### **Level 1: Statistical Detection**

Simple metrics from database:

```python
def detect_statistical_weaknesses(user_id):
    weaknesses = []
    
    # Get user's performance by category
    for category in ALL_CATEGORIES:
        stats = get_category_stats(user_id, category)
        
        if stats.accuracy < 0.6:
            weaknesses.append({
                "category": category,
                "type": "low_accuracy",
                "severity": "major" if stats.accuracy < 0.4 else "minor",
                "metric": stats.accuracy
            })
        
        if stats.avg_time > expected_time * 1.5:
            weaknesses.append({
                "category": category,
                "type": "slow_solving",
                "severity": "minor",
                "metric": stats.avg_time
            })
    
    return weaknesses
```

#### **Level 2: Pattern Detection**

Analyze mistake patterns:

```python
def detect_mistake_patterns(user_id):
    attempts = get_user_attempts(user_id, limit=50)
    incorrect_attempts = [a for a in attempts if not a.is_correct]
    
    # Group by mistake type
    mistake_groups = defaultdict(list)
    for attempt in incorrect_attempts:
        for mistake in attempt.mistakes:
            mistake_groups[mistake['type']].append(attempt)
    
    # Find frequent patterns
    patterns = []
    for mistake_type, attempts in mistake_groups.items():
        frequency = len(attempts) / len(incorrect_attempts)
        
        if frequency > 0.3:  # Appears in >30% of mistakes
            patterns.append({
                "pattern": mistake_type,
                "frequency": frequency,
                "examples": [a.question_id for a in attempts[:5]]
            })
    
    return patterns
```

#### **Level 3: Comparative Analysis**

Compare to similar users:

```python
def detect_by_comparison(user_id):
    user_level = get_user_difficulty_level(user_id)
    similar_users = find_users_at_level(user_level)
    
    user_performance = get_performance_profile(user_id)
    peer_avg_performance = average_performance(similar_users)
    
    weaknesses = []
    for category in ALL_CATEGORIES:
        user_score = user_performance[category]
        peer_score = peer_avg_performance[category]
        
        # More than 15 points below peers
        if user_score < peer_score - 15:
            weaknesses.append({
                "category": category,
                "type": "below_peers",
                "gap": peer_score - user_score,
                "user_score": user_score,
                "peer_score": peer_score
            })
    
    return weaknesses
```

#### **Level 4: LLM Deep Analysis**

For subtle weaknesses:

```python
def llm_weakness_analysis(user_id):
    # Get user's incorrect attempts with their code/explanations
    attempts = get_detailed_attempts(user_id, incorrect_only=True, limit=20)
    
    analysis_prompt = f"""
    Analyze these 20 incorrect attempts to identify subtle weaknesses:
    
    {format_attempts_for_llm(attempts)}
    
    Look for:
    1. Conceptual misunderstandings (not just implementation bugs)
    2. Blind spots in problem-solving approach
    3. Knowledge gaps that span multiple problems
    4. Thinking patterns that lead to errors
    
    Provide deep analysis of root causes.
    """
    
    llm_response = llm.generate(analysis_prompt)
    return parse_llm_weakness_analysis(llm_response)
```

### **Weakness Severity Calculation**

```python
def calculate_weakness_severity(weakness):
    # Factors
    frequency = weakness.error_frequency        # 0-1
    impact = weakness.impact_on_overall_score   # 0-1
    difficulty_to_fix = estimate_fix_difficulty(weakness)  # 0-1
    
    # Weighted severity
    severity_score = (
        frequency * 0.4 +
        impact * 0.4 +
        difficulty_to_fix * 0.2
    )
    
    if severity_score > 0.7:
        return "critical"
    elif severity_score > 0.5:
        return "major"
    elif severity_score > 0.3:
        return "minor"
    else:
        return "negligible"
```

### **Remediation Recommendations**

```python
def generate_remediation_plan(weakness):
    # Find prerequisite topics
    prerequisites = graph_query("""
        MATCH (weak:Topic {name: $weakness_topic})<-[:PREREQUISITE_FOR*]-(prereq:Topic)
        RETURN prereq
    """, weakness_topic=weakness.topic)
    
    # Find practice questions targeting this weakness
    practice_questions = db.query("""
        SELECT question_id, difficulty_score
        FROM questions
        WHERE $weakness_concept = ANY(concepts_tested)
        AND difficulty_score BETWEEN $min_diff AND $max_diff
        ORDER BY times_correct ASC  -- Less solved = fresher
        LIMIT 20
    """, weakness_concept=weakness.concept,
        min_diff=user_level - 10,
        max_diff=user_level + 10)
    
    return {
        "weakness": weakness.description,
        "severity": weakness.severity,
        "prerequisites": prerequisites,
        "practice_questions": practice_questions,
        "estimated_practice_time": len(practice_questions) * 15,  # minutes
        "review_materials": find_learning_resources(weakness.topic)
    }
```

---

## 8. SEARCH & FILTER SYSTEM

### **Search Architecture**

```
User Query
    ↓
┌───────────────────┐
│ Query Parser      │ → Extract: keywords, filters, intent
└────────┬──────────┘
         ↓
    ┌────────────────────┐
    │ Multi-Search       │
    │ Execution          │
    ├────────────────────┤
    │ 1. SQL Filter      │ → Exact matches (difficulty, category)
    │ 2. Vector Search   │ → Semantic similarity
    │ 3. Graph Traversal │ → Related problems
    └────────┬───────────┘
             ↓
    ┌────────────────┐
    │ Result Ranking │ → Combine & rank results
    └────────┬───────┘
             ↓
    ┌────────────────┐
    │ Return Results │
    └────────────────┘
```

### **Query Types**

#### **1. Keyword Search**

```sql
SELECT * FROM problems
WHERE 
    to_tsvector('english', title || ' ' || description) 
    @@ plainto_tsquery('english', $search_term)
ORDER BY ts_rank(to_tsvector('english', description), plainto_tsquery('english', $search_term)) DESC;
```

#### **2. Filter Search**

```python
filters = {
    "difficulty_range": [30, 60],
    "categories": ["Two Pointers", "Binary Search"],
    "companies": ["Google", "Amazon"],
    "tags": ["array", "sorting"],
    "success_rate_min": 0.5,
    "is_premium": False
}

query = """
SELECT * FROM problems
WHERE 
    difficulty_score BETWEEN :min_diff AND :max_diff
    AND category = ANY(:categories)
    AND companies ?| :companies  -- JSONB contains any
    AND tags ?& :tags           -- JSONB contains all
    AND success_rate >= :success_min
ORDER BY difficulty_score
"""
```

#### **3. Semantic Search**

```python
def semantic_search(query_text, filters=None, limit=10):
    # Get query embedding
    query_embedding = embed_text(query_text)
    
    # Vector search
    results = vector_db.query(
        query_embeddings=[query_embedding],
        n_results=limit * 2,  # Get more for filtering
        where=convert_filters_to_vector_where(filters)
    )
    
    # Enrich with full data from PostgreSQL
    problem_ids = [r['metadata']['problem_id'] for r in results]
    problems = db.query(
        "SELECT * FROM problems WHERE problem_id = ANY(:ids)",
        ids=problem_ids
    )
    
    return problems
```

#### **4. Graph-Based Search**

```cypher
-- Find all follow-ups to a problem
MATCH (p:Problem {problem_id: $start_id})-[:FOLLOW_UP*1..3]->(followup:Problem)
RETURN DISTINCT followup
ORDER BY followup.difficulty_score

-- Find problems requiring same prerequisites
MATCH (p1:Problem {problem_id: $id})-[:REQUIRES]->(prereq:Problem)<-[:REQUIRES]-(p2:Problem)
WHERE p1.problem_id <> p2.problem_id
RETURN p2

-- Find problems in learning path
MATCH path = (start:Problem {problem_id: $start})-[:LEADS_TO*1..5]->(end:Problem)
WHERE end.difficulty_score <= $max_diff
RETURN path
ORDER BY length(path)
```

### **Unified Search Interface**

```python
def unified_search(
    query: str,
    filters: dict = None,
    search_types: list = ["keyword", "semantic", "graph"],
    limit: int = 20
):
    results = []
    
    # Execute each search type
    if "keyword" in search_types:
        keyword_results = keyword_search(query, filters)
        results.extend([(r, "keyword", 1.0) for r in keyword_results])
    
    if "semantic" in search_types:
        semantic_results = semantic_search(query, filters)
        results.extend([(r, "semantic", 0.9) for r in semantic_results])
    
    if "graph" in search_types and filters.get("related_to"):
        graph_results = graph_search(filters["related_to"])
        results.extend([(r, "graph", 0.8) for r in graph_results])
    
    # Merge and rank
    merged = merge_and_rank_results(results)
    
    return merged[:limit]

def merge_and_rank_results(results):
    # Group by problem_id
    problem_scores = defaultdict(list)
    
    for problem, source, confidence in results:
        problem_scores[problem.id].append((problem, source, confidence))
    
    # Calculate combined score
    ranked = []
    for problem_id, entries in problem_scores.items():
        # Multiple sources = higher confidence
        combined_score = sum(conf for _, _, conf in entries) / len(entries)
        combined_score *= (1 + 0.1 * (len(entries) - 1))  # Bonus for multiple sources
        
        ranked.append((entries[0][0], combined_score))
    
    ranked.sort(key=lambda x: x[1], reverse=True)
    return [problem for problem, score in ranked]
```

### **Advanced Filters**

```python
advanced_filters = {
    # Difficulty
    "difficulty_range": [min, max],
    "difficulty_exact": 45,
    
    # Categories
    "categories": ["DP", "Graph"],
    "exclude_categories": ["Math"],
    
    # Problem attributes
    "has_follow_up": True,
    "has_multiple_solutions": True,
    "optimal_uses_template": "BFS",
    
    # Metadata
    "companies": ["Google", "Meta"],
    "frequency": "high",  # high/medium/low
    "tags": ["array", "hash_table"],
    "exclude_tags": ["geometry"],
    
    # User-specific
    "not_attempted": True,
    "attempted_but_incorrect": True,
    "similar_to_mastered": problem_id,
    "recommended_for_weakness": "edge_cases",
    
    # Graph relationships
    "related_to": problem_id,
    "follow_up_to": problem_id,
    "requires_solving": problem_id,
    
    # Performance
    "success_rate_range": [0.4, 0.7],  # Not too easy, not too hard
    "avg_time_range": [10*60, 30*60],  # 10-30 minutes
    
    # Learning
    "uses_pattern": "sliding_window",
    "concepts_tested": ["two_pointers", "sorting"],
    "difficulty_step_from": (problem_id, 5),  # 5 points harder than this
}
```

---

## 9. TRAINING PLAN BUILDER

### **Plan Types**

#### **1. Preset Plans**

Pre-defined curricula:

```python
preset_plans = {
    "beginner_fundamentals": {
        "duration": 30,  # days
        "difficulty_range": [0, 40],
        "categories": [
            "Complexity Analysis",
            "Data Structure Selection",
            "Pattern Recognition"
        ],
        "daily_questions": 5,
        "structure": "linear"
    },
    
    "interview_prep_2_weeks": {
        "duration": 14,
        "difficulty_range": [30, 70],
        "categories": "all",
        "daily_questions": 10,
        "structure": "mixed",
        "includes_mock_interviews": True,
        "mock_frequency": 3  # Every 3 days
    },
    
    "pattern_mastery": {
        "duration": 60,
        "focus": "pattern_recognition",
        "covers_all_patterns": True,
        "structure": "pattern_by_pattern"
    }
}
```

#### **2. Custom Plans**

User-defined:

```python
def create_custom_plan(user_id, preferences):
    plan = {
        "name": preferences.name,
        "duration": preferences.duration_days,
        "goals": preferences.goals,  # ["master_DP", "improve_speed"]
        "difficulty_range": auto_set_difficulty_range(user_id),
        "time_per_day": preferences.time_per_day,  # minutes
        "focus_areas": preferences.focus_categories,
        "avoid_areas": preferences.avoid_categories
    }
    
    # Generate daily schedule
    plan["schedule"] = generate_schedule(plan)
    
    return plan
```

#### **3. AI-Generated Plans**

Personalized by LLM:

```python
def generate_ai_plan(user_id, goal):
    # Get user profile
    user_profile = {
        "current_level": get_user_level(user_id),
        "strengths": get_strengths(user_id),
        "weaknesses": get_weaknesses(user_id),
        "learning_pace": estimate_learning_pace(user_id),
        "available_time": get_available_time(user_id),
        "past_attempts": get_attempt_summary(user_id)
    }
    
    prompt = f"""
    Create a personalized training plan:
    
    User Profile:
    {json.dumps(user_profile, indent=2)}
    
    Goal: {goal}
    
    Generate a plan that:
    1. Addresses weaknesses first
    2. Builds on strengths
    3. Progresses difficulty gradually
    4. Includes spaced repetition
    5. Balances across all categories
    
    Return JSON with daily schedule.
    """
    
    llm_plan = llm.generate(prompt)
    return parse_and_validate_plan(llm_plan)
```

### **Dynamic Plan Adjustment**

Plans adapt based on performance:

```python
def adjust_plan_based_on_performance(plan_id):
    plan = get_plan(plan_id)
    recent_performance = get_recent_performance(plan.user_id, days=7)
    
    adjustments = []
    
    # Too easy?
    if recent_performance.avg_accuracy > 0.85:
        adjustments.append({
            "type": "increase_difficulty",
            "amount": 5,  # points
            "reason": "High accuracy indicates room for challenge"
        })
    
    # Too hard?
    if recent_performance.avg_accuracy < 0.4:
        adjustments.append({
            "type": "decrease_difficulty",
            "amount": 5,
            "reason": "Low accuracy indicates excessive difficulty"
        })
    
    # Stuck on category?
    for category, stats in recent_performance.by_category.items():
        if stats.attempts > 5 and stats.accuracy < 0.3:
            adjustments.append({
                "type": "add_prerequisite_review",
                "category": category,
                "reason": f"Struggling with {category}"
            })
    
    # Apply adjustments
    apply_plan_adjustments(plan_id, adjustments)
```

### **Plan Structure Templates**

#### **Linear Progression**

```python
# Day 1: Easy complexity questions
# Day 2: Easy DS selection
# Day 3: Easy pattern recognition
# ...
# Day 30: Hard hybrid challenges

def generate_linear_plan(duration, categories, difficulty_range):
    questions_per_day = []
    
    difficulty_step = (difficulty_range[1] - difficulty_range[0]) / duration
    
    for day in range(1, duration + 1):
        target_difficulty = difficulty_range[0] + (day * difficulty_step)
        category = categories[day % len(categories)]
        
        questions = select_questions(
            category=category,
            difficulty=target_difficulty,
            count=5
        )
        
        questions_per_day.append({
            "day": day,
            "questions": questions,
            "target_difficulty": target_difficulty
        })
    
    return questions_per_day
```

#### **Spiral Curriculum**

```python
# Cycle through topics but at increasing depth
# Week 1: All topics at easy level
# Week 2: All topics at medium level
# Week 3: All topics at hard level
# Week 4: Review + hybrid

def generate_spiral_plan(duration, categories):
    weeks = duration // 7
    plan = []
    
    difficulty_levels = [
        (0, 30),   # Week 1: Easy
        (30, 60),  # Week 2: Medium
        (60, 80),  # Week 3: Hard
        (40, 70)   # Week 4: Mixed review
    ]
    
    for week in range(weeks):
        diff_range = difficulty_levels[week % len(difficulty_levels)]
        
        for day in range(7):
            category = categories[day % len(categories)]
            questions = select_questions(
                category=category,
                difficulty_range=diff_range,
                count=5
            )
            plan.append({"week": week+1, "day": day+1, "questions": questions})
    
    return plan
```

#### **Mastery-Based Progression**

```python
# Don't move forward until mastery achieved

def generate_mastery_plan(user_id, categories):
    plan = []
    current_category_idx = 0
    
    while current_category_idx < len(categories):
        category = categories[current_category_idx]
        
        # Get questions for this category
        questions = get_category_questions(category, user_level(user_id))
        
        plan.append({
            "category": category,
            "questions": questions,
            "mastery_requirement": 0.8,  # 80% accuracy to move on
            "min_attempts": 10
        })
        
        current_category_idx += 1
    
    return plan
```

### **Spaced Repetition Integration**

```python
def schedule_reviews(user_id):
    # Get questions user has answered correctly
    mastered_questions = db.query("""
        SELECT question_id, MAX(attempted_at) as last_attempt
        FROM user_progress
        WHERE user_id = :user_id AND is_correct = TRUE
        GROUP BY question_id
    """, user_id=user_id)
    
    reviews_needed = []
    
    for q in mastered_questions:
        # SM-2 algorithm for spaced repetition
        days_since = (datetime.now() - q.last_attempt).days
        
        # Get user's ease factor for this question
        ease = get_ease_factor(user_id, q.question_id)
        interval = get_current_interval(user_id, q.question_id)
        
        next_review = q.last_attempt + timedelta(days=interval)
        
        if datetime.now() >= next_review:
            reviews_needed.append({
                "question_id": q.question_id,
                "due_date": next_review,
                "priority": (datetime.now() - next_review).days  # Overdue = higher priority
            })
    
    return sorted(reviews_needed, key=lambda x: x['priority'], reverse=True)
```

---

## 10. LLM INTEGRATION ARCHITECTURE

### **Local LLM Options**

For privacy and offline usage:

```python
llm_options = {
    "ollama": {
        "models": ["llama3.1", "mistral", "codellama"],
        "pros": "Easy setup, good UI",
        "cons": "Requires GPU for good performance"
    },
    
    "llamacpp": {
        "models": ["llama-3.1-8b-gguf", "mistral-7b-gguf"],
        "pros": "CPU-friendly, fast",
        "cons": "CLI-based"
    },
    
    "localai": {
        "models": "Various",
        "pros": "OpenAI-compatible API",
        "cons": "More complex setup"
    }
}
```

### **LLM Service Architecture**

```
┌─────────────────────────────────────┐
│         LLM Service Layer           │
├─────────────────────────────────────┤
│                                     │
│  ┌──────────────────────────────┐  │
│  │   Task Router                │  │
│  │   (Determines which prompt)  │  │
│  └──────────┬───────────────────┘  │
│             ↓                       │
│  ┌──────────────────────────────┐  │
│  │   RAG Context Retriever      │  │
│  │   (Vector DB search)         │  │
│  └──────────┬───────────────────┘  │
│             ↓                       │
│  ┌──────────────────────────────┐  │
│  │   Prompt Constructor         │  │
│  │   (Build final prompt)       │  │
│  └──────────┬───────────────────┘  │
│             ↓                       │
│  ┌──────────────────────────────┐  │
│  │   LLM Engine                 │  │
│  │   (Ollama/LlamaCPP)          │  │
│  └──────────┬───────────────────┘  │
│             ↓                       │
│  ┌──────────────────────────────┐  │
│  │   Response Parser            │  │
│  │   (Extract structured data)  │  │
│  └──────────────────────────────┘  │
│                                     │
└─────────────────────────────────────┘
```

### **Prompt Templates**

#### **Question Generation**

```python
QUESTION_GENERATION_PROMPT = """
You are an expert in creating technical interview questions.

CONTEXT (from similar problems):
{rag_context}

BASE PROBLEM:
{base_problem}

TASK: Generate a new question that:
1. Tests the same algorithmic pattern
2. Uses a different scenario/context
3. Has similar difficulty level ({difficulty_score}/100)
4. Cannot be solved by memorization of the base problem

REQUIREMENTS:
- Clear problem statement
- 2-3 examples
- Constraints clearly stated
- Expected time/space complexity: {expected_complexity}

Return JSON:
{{
    "question_text": "...",
    "examples": [...],
    "constraints": [...],
    "hints": [...],
    "expected_solution_approach": "..."
}}
"""
```

#### **Solution Explanation**

```python
EXPLANATION_PROMPT = """
Explain this solution in a way that helps deep understanding:

PROBLEM:
{problem}

SOLUTION:
{solution_code}

Explain:
1. WHY this approach works (intuition)
2. HOW the algorithm progresses (step-by-step)
3. WHAT makes this optimal (complexity analysis)
4. WHEN to use this pattern (similar problems)

Provide clear, educational explanation suitable for learning.
"""
```

#### **Assessment & Feedback**

```python
ASSESSMENT_PROMPT = """
Evaluate this submission:

PROBLEM: {problem}
USER SOLUTION: {user_solution}
USER EXPLANATION: {user_explanation}

Provide:
1. Correctness score (0-100)
2. Approach quality (0-100)
3. Understanding depth (0-100)
4. Specific strengths (list)
5. Specific weaknesses (list)
6. Improvement suggestions
7. Memorization indicator (0-1)

Be constructive and educational.

Return JSON with structured evaluation.
"""
```

### **RAG Implementation**

```python
def rag_enhanced_prompt(query, task_type):
    # Get relevant context from vector DB
    query_embedding = embed_text(query)
    
    context_results = vector_db.query(
        query_embeddings=[query_embedding],
        n_results=5
    )
    
    # Format context
    context_text = "\n\n".join([
        f"Example {i+1}:\n{r['metadata']['text']}"
        for i, r in enumerate(context_results)
    ])
    
    # Get prompt template for task
    prompt_template = get_prompt_template(task_type)
    
    # Fill in prompt
    final_prompt = prompt_template.format(
        rag_context=context_text,
        query=query,
        **additional_params
    )
    
    return final_prompt
```

---

## 11. QUESTION GENERATION PIPELINE

### **Generation Strategies**

#### **1. Problem Variant Generation**

```python
def generate_problem_variant(base_problem_id):
    base = get_problem(base_problem_id)
    
    # Get similar problems for context
    similar = find_similar_problems(base_problem_id, limit=5)
    context = "\n".join([p.description for p in similar])
    
    # Variation strategies
    strategies = [
        "change_data_structure",
        "modify_constraints",
        "generalize_problem",
        "add_complexity",
        "different_domain"
    ]
    
    strategy = random.choice(strategies)
    
    prompt = f"""
    BASE PROBLEM:
    {base.description}
    
    SIMILAR PROBLEMS FOR CONTEXT:
    {context}
    
    STRATEGY: {strategy}
    
    Generate a NEW problem that applies the same algorithmic pattern but with {strategy}.
    
    Ensure:
    - Same core pattern
    - Similar difficulty ({base.difficulty_score}/100)
    - Different enough that it can't be solved by memory
    - Clear and well-defined
    
    Return complete problem definition.
    """
    
    new_problem = llm.generate(prompt)
    
    # Store with relationship
    new_id = store_problem(new_problem)
    create_graph_relationship(base_problem_id, new_id, "VARIATION_OF")
    
    return new_id
```

#### **2. Follow-Up Generation**

```python
def generate_follow_up(original_problem_id):
    original = get_problem(original_problem_id)
    
    prompt = f"""
    ORIGINAL PROBLEM:
    {original.description}
    
    ORIGINAL DIFFICULTY: {original.difficulty_score}/100
    
    Generate a FOLLOW-UP problem that:
    1. Builds on the original concept
    2. Is 10-20 points harder
    3. Requires understanding the original
    4. Adds new twist or constraint
    
    Common follow-up patterns:
    - Generalize (Two Sum → K Sum)
    - Add dimension (1D → 2D)
    - Add constraint (must be O(1) space)
    - Combine concepts (add graph traversal)
    
    Return complete problem.
    """
    
    follow_up = llm.generate(prompt)
    
    follow_up_id = store_problem(follow_up)
    create_graph_relationship(original_problem_id, follow_up_id, "FOLLOW_UP", 
                             {"difficulty_increase": 15})
    
    return follow_up_id
```

#### **3. Multiple Choice Question Generation**

```python
def generate_mcq_question(problem_id, question_type):
    problem = get_problem(problem_id)
    
    # Get examples of this question type
    examples = get_example_questions(question_type, limit=3)
    
    prompt = f"""
    PROBLEM:
    {problem.description}
    
    QUESTION TYPE: {question_type}
    (e.g., "Complexity Analysis", "Approach Comparison")
    
    EXAMPLES OF THIS TYPE:
    {format_examples(examples)}
    
    Generate a multiple-choice question of this type based on the problem.
    
    Requirements:
    - 4 options (A, B, C, D)
    - All options plausible
    - One clearly best answer
    - Explanation for why each is right/wrong
    
    Return JSON with question, options, correct answer, and explanations.
    """
    
    mcq = llm.generate(prompt)
    return parse_and_store_question(mcq, problem_id, question_type)
```

### **Quality Control Pipeline**

```python
def quality_control_generated_question(question_id):
    question = get_question(question_id)
    
    checks = []
    
    # Check 1: Is it well-formed?
    checks.append(check_well_formed(question))
    
    # Check 2: Is difficulty calibrated correctly?
    estimated_diff = estimate_difficulty(question)
    checks.append(abs(estimated_diff - question.difficulty_score) < 10)
    
    # Check 3: Are options distinct?
    if question.question_format == "multiple_choice":
        checks.append(check_options_distinct(question.options))
    
    # Check 4: Is explanation clear?
    checks.append(check_explanation_quality(question.explanation))
    
    # Check 5: Does it test what it claims?
    checks.append(verify_concepts_tested(question))
    
    quality_score = sum(checks) / len(checks)
    
    if quality_score < 0.8:
        # Needs human review or regeneration
        flag_for_review(question_id, quality_score)
    
    return quality_score
```

---

## 12. LOCAL DEPLOYMENT SETUP

### **System Requirements**

**Minimum:**
- CPU: 4 cores
- RAM: 16 GB
- Storage: 50 GB SSD
- OS: Linux (Ubuntu 22.04+) or macOS

**Recommended:**
- CPU: 8+ cores
- RAM: 32 GB
- Storage: 100 GB SSD
- GPU: NVIDIA GPU with 8+ GB VRAM (for faster LLM)
- OS: Ubuntu 22.04

### **Docker Compose Setup**

```yaml
version: '3.8'

services:
  # PostgreSQL + AGE
  postgres:
    image: apache/age:latest
    environment:
      POSTGRES_DB: interview_platform
      POSTGRES_USER: platform_user
      POSTGRES_PASSWORD: secure_password
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init-scripts:/docker-entrypoint-initdb.d
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U platform_user"]
      interval: 10s
      timeout: 5s
      retries: 5

  # ChromaDB (Vector Database)
  chromadb:
    image: chromadb/chroma:latest
    volumes:
      - chroma_data:/chroma/chroma
    ports:
      - "8000:8000"
    environment:
      ALLOW_RESET: "true"
      ANONYMIZED_TELEMETRY: "false"

  # Ollama (LLM)
  ollama:
    image: ollama/ollama:latest
    volumes:
      - ollama_data:/root/.ollama
    ports:
      - "11434:11434"
    deploy:
      resources:
        reservations:
          devices:
            - driver: nvidia
              count: 1
              capabilities: [gpu]

  # Backend API
  backend:
    build: ./backend
    depends_on:
      - postgres
      - chromadb
      - ollama
    environment:
      DATABASE_URL: postgresql://platform_user:secure_password@postgres:5432/interview_platform
      CHROMA_URL: http://chromadb:8000
      OLLAMA_URL: http://ollama:11434
    ports:
      - "8080:8080"
    volumes:
      - ./backend:/app
      - ./data:/app/data

  # Frontend
  frontend:
    build: ./frontend
    depends_on:
      - backend
    ports:
      - "3000:3000"
    volumes:
      - ./frontend:/app
      - /app/node_modules

volumes:
  postgres_data:
  chroma_data:
  ollama_data:
```

### **Initialization Scripts**

**init-scripts/01-create-extensions.sql:**
```sql
-- Enable AGE extension
CREATE EXTENSION IF NOT EXISTS age;
LOAD 'age';
SET search_path = ag_catalog, "$user", public;

-- Create graph
SELECT create_graph('interview_graph');

-- Enable full-text search
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS btree_gin;
```

**init-scripts/02-create-schema.sql:**
```sql
-- Create all tables (from section 2)
CREATE TABLE problems (...);
CREATE TABLE questions (...);
-- ... etc
```

**init-scripts/03-create-indexes.sql:**
```sql
-- Performance indexes
CREATE INDEX CONCURRENTLY idx_problems_difficulty 
  ON problems(difficulty_score);

CREATE INDEX CONCURRENTLY idx_questions_type 
  ON questions(type_id);

-- Full-text search indexes
CREATE INDEX idx_problems_fts 
  ON problems 
  USING GIN(to_tsvector('english', title || ' ' || description));

-- JSONB indexes
CREATE INDEX idx_problems_tags 
  ON problems 
  USING GIN(tags);
```

### **Backend Structure**

```
backend/
├── main.py                 # FastAPI app entry
├── requirements.txt
├── config.py              # Configuration
├── database/
│   ├── postgres.py        # PostgreSQL connection
│   ├── age.py            # AGE graph queries
│   └── vector.py         # ChromaDB connection
├── services/
│   ├── llm_service.py    # LLM integration
│   ├── assessment.py     # Assessment logic
│   ├── training_plan.py  # Plan generation
│   ├── weakness.py       # Weakness detection
│   └── search.py         # Search implementation
├── api/
│   ├── problems.py       # Problem endpoints
│   ├── questions.py      # Question endpoints
│   ├── progress.py       # Progress tracking
│   └── training.py       # Training plans
└── models/
    ├── schemas.py        # Pydantic models
    └── graph.py          # Graph node/edge definitions
```

### **Frontend Structure**

```
frontend/
├── src/
│   ├── App.jsx
│   ├── pages/
│   │   ├── Dashboard.jsx
│   │   ├── Practice.jsx
│   │   ├── Assessment.jsx
│   │   ├── TrainingPlan.jsx
│   │   └── Analytics.jsx
│   ├── components/
│   │   ├── QuestionCard.jsx
│   │   ├── CodeEditor.jsx
│   │   ├── ProgressChart.jsx
│   │   └── WeaknessReport.jsx
│   ├── services/
│   │   ├── api.js
│   │   └── analytics.js
│   └── utils/
│       ├── scoring.js
│       └── formatting.js
└── package.json
```

### **Setup Commands**

```bash
# 1. Clone repository (hypothetical)
git clone <repo>
cd interview-platform

# 2. Initialize environment
cp .env.example .env
# Edit .env with your settings

# 3. Start services
docker-compose up -d

# 4. Wait for services to be ready
docker-compose ps

# 5. Pull LLM model
docker exec -it interview-platform-ollama-1 ollama pull llama3.1

# 6. Initialize database
docker exec -it interview-platform-postgres-1 psql -U platform_user -d interview_platform -f /docker-entrypoint-initdb.d/01-create-extensions.sql

# 7. Seed initial data
docker exec -it interview-platform-backend-1 python scripts/seed_data.py

# 8. Create vector embeddings
docker exec -it interview-platform-backend-1 python scripts/generate_embeddings.py

# 9. Access application
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
# API Docs: http://localhost:8080/docs
```

### **Data Seeding Strategy**

```python
# scripts/seed_data.py

def seed_database():
    print("Seeding database...")
    
    # 1. Seed topics
    print("Seeding topics...")
    seed_topics(topics_data)
    
    # 2. Seed problems from LeetCode dataset
    print("Seeding problems...")
    seed_problems(leetcode_problems)
    
    # 3. Create graph relationships
    print("Creating graph relationships...")
    create_problem_relationships()
    
    # 4. Seed question types
    print("Seeding question types...")
    seed_question_types(question_types_taxonomy)
    
    # 5. Generate initial questions
    print("Generating questions...")
    for problem in get_all_problems():
        generate_questions_for_problem(problem.id)
    
    # 6. Seed templates
    print("Seeding templates...")
    seed_templates(algorithm_templates)
    
    # 7. Generate embeddings
    print("Generating embeddings...")
    generate_all_embeddings()
    
    print("Database seeded successfully!")
```

### **Backup & Restore**

```bash
# Backup
docker exec interview-platform-postgres-1 pg_dump -U platform_user interview_platform > backup.sql
docker exec interview-platform-chromadb-1 tar -czf - /chroma/chroma > chroma_backup.tar.gz

# Restore
docker exec -i interview-platform-postgres-1 psql -U platform_user interview_platform < backup.sql
docker exec -i interview-platform-chromadb-1 tar -xzf - -C /chroma < chroma_backup.tar.gz
```

### **Monitoring & Logs**

```bash
# View logs
docker-compose logs -f backend
docker-compose logs -f postgres

# Monitor resource usage
docker stats

# Check database health
docker exec interview-platform-postgres-1 pg_isready

# Check LLM status
curl http://localhost:11434/api/tags
```

---

## SUMMARY & NEXT STEPS

### **What We've Defined**

✅ **Complete database schema** (PostgreSQL + AGE graph)  
✅ **Vector database integration** (ChromaDB + RAG)  
✅ **Difficulty scoring formula** (0-100 granular scale)  
✅ **LLM-powered assessment** (fair, anti-memorization)  
✅ **Weakness detection** (multi-level algorithm)  
✅ **Advanced search system** (keyword + semantic + graph)  
✅ **Training plan builder** (preset + custom + AI)  
✅ **Question generation pipeline** (LLM + quality control)  
✅ **Local deployment setup** (Docker Compose stack)

### **This System Enables**

1. **Intelligent Practice** - Questions adapt to user level
2. **Deep Insights** - Know exactly where you're weak
3. **Fair Assessment** - Can't game the system with memorization
4. **Personalized Learning** - Custom plans for your goals
5. **Rich Relationships** - Find similar/follow-up problems easily
6. **Semantic Search** - Find problems by description, not just tags
7. **Continuous Improvement** - LLM generates new questions
8. **Complete Privacy** - Everything runs locally

### **Implementation Priority**

**Phase 1 (Core Infrastructure):**
1. PostgreSQL + AGE setup
2. Basic schema creation
3. Problem seeding
4. Simple search

**Phase 2 (Intelligence Layer):**
1. Vector DB integration
2. Embedding generation
3. LLM integration (Ollama)
4. Difficulty scoring

**Phase 3 (Features):**
1. Assessment system
2. Weakness detection
3. Training plans
4. Question generation

**Phase 4 (Polish):**
1. Frontend UI
2. Analytics dashboard
3. Performance optimization
4. User testing

---

**This is your complete technical blueprint.**  
**Everything you need to build a production-grade interview training platform that runs locally.**

*Ready to start implementation!*
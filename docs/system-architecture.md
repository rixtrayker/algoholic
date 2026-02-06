# Interview Training Platform - System Architecture & Database Design
## PostgreSQL + Apache AGE + Vector DB + RAG + LLM Assessment

---

## TABLE OF CONTENTS

1. [System Overview](#system-overview)
2. [Database Architecture](#database-architecture)
3. [Graph Relationships (Apache AGE)](#graph-relationships)
4. [Vector Database & RAG](#vector-database-rag)
5. [Difficulty Scoring System](#difficulty-scoring-system)
6. [Assessment Framework](#assessment-framework)
7. [Search & Filtering](#search-filtering)
8. [Training Plan Builder](#training-plan-builder)
9. [LLM Integration](#llm-integration)
10. [Local Deployment](#local-deployment)
11. [API Design](#api-design)

---

## 1. SYSTEM OVERVIEW

### **Architecture Components**

```
┌─────────────────────────────────────────────────────────────────┐
│                         USER INTERFACE                          │
│                    (Web App / Local Client)                     │
└────────────────┬────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                      APPLICATION LAYER                          │
│  ┌─────────────┐  ┌──────────────┐  ┌────────────────────┐    │
│  │   FastAPI   │  │   Training   │  │   Assessment      │    │
│  │   Backend   │  │   Planner    │  │   Engine (LLM)    │    │
│  └─────────────┘  └──────────────┘  └────────────────────┘    │
└────────┬──────────────────┬──────────────────┬──────────────────┘
         │                  │                  │
         ▼                  ▼                  ▼
┌────────────────┐  ┌──────────────┐  ┌──────────────────┐
│  PostgreSQL +  │  │  Vector DB   │  │   LLM (Local)    │
│  Apache AGE    │  │  (Chroma/    │  │   Ollama/        │
│  (Graph DB)    │  │   Qdrant)    │  │   LM Studio      │
└────────────────┘  └──────────────┘  └──────────────────┘
         │                  │                  │
         └──────────────────┴──────────────────┘
                            │
                     ┌──────▼───────┐
                     │     RAG      │
                     │   Pipeline   │
                     └──────────────┘
```

### **Core Features**

1. **Relational Data** (PostgreSQL)
   - Problems, questions, solutions
   - User progress, attempts, scores
   - Training plans, assessments

2. **Graph Relationships** (Apache AGE)
   - Problem → Problem (similar, follow-up, prerequisite)
   - Problem → Topic → Tag
   - Topic → Topic (dependency, related)
   - User skill graph

3. **Semantic Search** (Vector DB)
   - Problem similarity
   - Question clustering
   - Solution pattern matching
   - Intelligent recommendations

4. **AI Assessment** (LLM + RAG)
   - Evaluate solutions
   - Detect memorization
   - Identify weaknesses
   - Generate personalized feedback
   - Create new problems

---

## 2. DATABASE ARCHITECTURE

### **PostgreSQL Schema Design**

#### **Core Tables**

```sql
-- ============================================================
-- PROBLEMS & QUESTIONS
-- ============================================================

CREATE TABLE problems (
    problem_id BIGSERIAL PRIMARY KEY,
    
    -- Identification
    leetcode_number INT UNIQUE,           -- Official LeetCode number
    internal_code VARCHAR(50) UNIQUE,     -- Our internal identifier
    title TEXT NOT NULL,
    slug VARCHAR(200) UNIQUE NOT NULL,
    
    -- Content
    description TEXT NOT NULL,
    constraints TEXT[],                    -- Array of constraint strings
    examples JSONB,                        -- Input/output examples
    hints TEXT[],
    
    -- Metadata
    source VARCHAR(50),                    -- 'leetcode', 'custom', 'generated'
    original_url TEXT,
    
    -- Difficulty (our custom scoring - see section 5)
    base_difficulty_score NUMERIC(5,2),   -- 0.00 to 100.00
    community_difficulty_score NUMERIC(5,2),
    time_pressure_score NUMERIC(5,2),
    implementation_complexity_score NUMERIC(5,2),
    
    -- Statistics
    acceptance_rate NUMERIC(5,2),
    submission_count INT DEFAULT 0,
    success_count INT DEFAULT 0,
    average_solve_time_seconds INT,
    
    -- Status
    is_active BOOLEAN DEFAULT true,
    is_premium BOOLEAN DEFAULT false,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_problems_difficulty ON problems(base_difficulty_score);
CREATE INDEX idx_problems_leetcode ON problems(leetcode_number);
CREATE INDEX idx_problems_source ON problems(source);


-- ============================================================
-- PROBLEM FOLLOW-UPS & VARIATIONS
-- ============================================================

CREATE TABLE problem_relationships (
    relationship_id BIGSERIAL PRIMARY KEY,
    
    source_problem_id BIGINT REFERENCES problems(problem_id),
    target_problem_id BIGINT REFERENCES problems(problem_id),
    
    relationship_type VARCHAR(50) NOT NULL,
    -- Types: 'follow_up', 'similar', 'prerequisite', 
    --        'easier_version', 'harder_version', 
    --        'same_pattern', 'alternative_approach'
    
    strength NUMERIC(3,2),                -- 0.00 to 1.00 (similarity strength)
    description TEXT,                      -- Why are they related?
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(source_problem_id, target_problem_id, relationship_type)
);

CREATE INDEX idx_rel_source ON problem_relationships(source_problem_id);
CREATE INDEX idx_rel_target ON problem_relationships(target_problem_id);
CREATE INDEX idx_rel_type ON problem_relationships(relationship_type);


-- ============================================================
-- QUESTIONS (Specific training questions about problems)
-- ============================================================

CREATE TABLE questions (
    question_id BIGSERIAL PRIMARY KEY,
    
    -- Categorization (from taxonomy document)
    category VARCHAR(50) NOT NULL,        -- 'complexity_analysis', 'ds_selection', etc.
    subcategory VARCHAR(50),              -- '1.1', '2.3', etc.
    
    -- Content
    question_text TEXT NOT NULL,
    question_type VARCHAR(50) NOT NULL,   -- 'multiple_choice', 'code', 'pseudocode', 'open_ended'
    
    -- For multiple choice
    options JSONB,                         -- [{id: 'A', text: '...', is_correct: true}, ...]
    
    -- For code questions
    starter_code TEXT,
    test_cases JSONB,                      -- [{input: ..., expected: ..., is_hidden: false}, ...]
    
    -- Metadata
    difficulty_level VARCHAR(20),          -- 'easy', 'medium', 'hard', 'expert'
    estimated_time_seconds INT,
    
    -- Explanation
    explanation TEXT,
    detailed_solution TEXT,
    common_mistakes TEXT[],
    
    -- Related problem (optional - questions can be standalone or problem-specific)
    related_problem_id BIGINT REFERENCES problems(problem_id),
    
    -- Tags for filtering
    tags TEXT[],                           -- ['two-pointers', 'sliding-window', etc.]
    concepts TEXT[],                       -- ['time-complexity', 'space-optimization']
    
    -- Statistics
    attempt_count INT DEFAULT 0,
    correct_count INT DEFAULT 0,
    average_time_seconds INT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_questions_category ON questions(category);
CREATE INDEX idx_questions_difficulty ON questions(difficulty_level);
CREATE INDEX idx_questions_problem ON questions(related_problem_id);
CREATE INDEX idx_questions_tags ON questions USING GIN(tags);
CREATE INDEX idx_questions_concepts ON questions USING GIN(concepts);


-- ============================================================
-- TOPICS & TAGS (Knowledge Graph)
-- ============================================================

CREATE TABLE topics (
    topic_id BIGSERIAL PRIMARY KEY,
    
    name VARCHAR(100) UNIQUE NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    category VARCHAR(50),                  -- 'data_structure', 'algorithm', 'pattern', 'concept'
    
    -- Learning metadata
    difficulty_level VARCHAR(20),
    estimated_learning_hours NUMERIC(4,1),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tags (
    tag_id BIGSERIAL PRIMARY KEY,
    
    name VARCHAR(50) UNIQUE NOT NULL,
    slug VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    color VARCHAR(7),                      -- Hex color for UI
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Many-to-many: Problems to Topics
CREATE TABLE problem_topics (
    problem_id BIGINT REFERENCES problems(problem_id) ON DELETE CASCADE,
    topic_id BIGINT REFERENCES topics(topic_id) ON DELETE CASCADE,
    relevance NUMERIC(3,2) DEFAULT 1.00,   -- How relevant is this topic to this problem
    
    PRIMARY KEY (problem_id, topic_id)
);

-- Many-to-many: Problems to Tags
CREATE TABLE problem_tags (
    problem_id BIGINT REFERENCES problems(problem_id) ON DELETE CASCADE,
    tag_id BIGINT REFERENCES tags(tag_id) ON DELETE CASCADE,
    
    PRIMARY KEY (problem_id, tag_id)
);

-- Many-to-many: Questions to Topics
CREATE TABLE question_topics (
    question_id BIGINT REFERENCES questions(question_id) ON DELETE CASCADE,
    topic_id BIGINT REFERENCES topics(topic_id) ON DELETE CASCADE,
    
    PRIMARY KEY (question_id, topic_id)
);


-- ============================================================
-- TRAINING PLANS
-- ============================================================

CREATE TABLE training_plans (
    plan_id BIGSERIAL PRIMARY KEY,
    
    name VARCHAR(200) NOT NULL,
    description TEXT,
    
    plan_type VARCHAR(50),                 -- 'structured', 'adaptive', 'custom', 'daily'
    difficulty_level VARCHAR(20),
    estimated_duration_days INT,
    
    -- Target goals
    target_topics TEXT[],
    target_skills TEXT[],
    
    -- Configuration
    questions_per_day INT,
    focus_areas JSONB,                     -- {category: weight, ...}
    
    is_template BOOLEAN DEFAULT false,     -- Can be used as template for user plans
    is_active BOOLEAN DEFAULT true,
    
    created_by_user_id BIGINT,             -- NULL if system-generated
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE training_plan_items (
    item_id BIGSERIAL PRIMARY KEY,
    
    plan_id BIGINT REFERENCES training_plans(plan_id) ON DELETE CASCADE,
    
    -- What to practice
    item_type VARCHAR(50) NOT NULL,        -- 'question', 'problem', 'assessment', 'review'
    question_id BIGINT REFERENCES questions(question_id),
    problem_id BIGINT REFERENCES problems(problem_id),
    
    -- Ordering
    sequence_number INT NOT NULL,
    day_number INT,                        -- Which day in the plan
    
    -- Requirements
    is_required BOOLEAN DEFAULT true,
    prerequisite_items INT[],              -- Must complete these items first
    
    -- Metadata
    estimated_time_minutes INT,
    notes TEXT,
    
    UNIQUE(plan_id, sequence_number)
);

CREATE INDEX idx_plan_items_plan ON training_plan_items(plan_id);
CREATE INDEX idx_plan_items_sequence ON training_plan_items(plan_id, sequence_number);


-- ============================================================
-- USER DATA
-- ============================================================

CREATE TABLE users (
    user_id BIGSERIAL PRIMARY KEY,
    
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    
    -- Profile
    display_name VARCHAR(100),
    avatar_url TEXT,
    
    -- Settings
    preferences JSONB,                     -- UI preferences, notification settings
    
    -- Statistics
    total_questions_attempted INT DEFAULT 0,
    total_problems_solved INT DEFAULT 0,
    current_streak_days INT DEFAULT 0,
    longest_streak_days INT DEFAULT 0,
    total_study_time_seconds BIGINT DEFAULT 0,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_active_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- ============================================================
-- USER ATTEMPTS & PROGRESS
-- ============================================================

CREATE TABLE question_attempts (
    attempt_id BIGSERIAL PRIMARY KEY,
    
    user_id BIGINT REFERENCES users(user_id),
    question_id BIGINT REFERENCES questions(question_id),
    
    -- Attempt data
    user_answer JSONB,                     -- Flexible: multiple choice selection, code, text
    is_correct BOOLEAN,
    score NUMERIC(5,2),                    -- 0.00 to 100.00 (partial credit possible)
    
    -- Timing
    time_spent_seconds INT,
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Analysis (LLM-powered)
    detected_patterns TEXT[],              -- What patterns user recognized
    mistakes_made TEXT[],                  -- What mistakes were detected
    shows_memorization BOOLEAN,            -- Did they just memorize the answer?
    confidence_score NUMERIC(3,2),         -- How confident was the answer
    
    -- Context
    attempt_number INT,                    -- nth attempt at this question
    training_plan_id BIGINT REFERENCES training_plans(plan_id),
    session_id VARCHAR(50),                -- Group attempts by session
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_attempts_user ON question_attempts(user_id);
CREATE INDEX idx_attempts_question ON question_attempts(question_id);
CREATE INDEX idx_attempts_user_time ON question_attempts(user_id, submitted_at);


CREATE TABLE problem_attempts (
    attempt_id BIGSERIAL PRIMARY KEY,
    
    user_id BIGINT REFERENCES users(user_id),
    problem_id BIGINT REFERENCES problems(problem_id),
    
    -- Solution submitted
    solution_code TEXT,
    programming_language VARCHAR(50),
    
    -- Results
    is_accepted BOOLEAN,
    test_cases_passed INT,
    test_cases_total INT,
    
    -- Performance
    runtime_ms INT,
    memory_mb NUMERIC(8,2),
    
    -- Analysis
    time_complexity_claimed VARCHAR(50),   -- What user said it is
    space_complexity_claimed VARCHAR(50),
    actual_complexity_detected VARCHAR(50), -- What LLM detected
    
    approach_used VARCHAR(100),            -- Pattern/algorithm used
    
    -- Quality metrics (LLM assessment)
    code_quality_score NUMERIC(5,2),
    readability_score NUMERIC(5,2),
    handles_edge_cases BOOLEAN,
    is_optimal BOOLEAN,
    
    -- Timing
    time_spent_seconds INT,
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    attempt_number INT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_problem_attempts_user ON problem_attempts(user_id);
CREATE INDEX idx_problem_attempts_problem ON problem_attempts(problem_id);


-- ============================================================
-- ASSESSMENTS (LLM-Powered Evaluation)
-- ============================================================

CREATE TABLE assessments (
    assessment_id BIGSERIAL PRIMARY KEY,
    
    user_id BIGINT REFERENCES users(user_id),
    
    assessment_type VARCHAR(50),           -- 'skill', 'weakness', 'comprehensive', 'progress'
    
    -- Scope
    topics_covered TEXT[],
    questions_included BIGINT[],           -- Array of question_ids
    
    -- Results
    overall_score NUMERIC(5,2),
    subscores JSONB,                       -- {category: score, ...}
    
    -- LLM Analysis
    strengths TEXT[],
    weaknesses TEXT[],
    recommendations TEXT[],
    personalized_feedback TEXT,
    
    -- Memorization detection
    memorization_likelihood NUMERIC(3,2),  -- 0.00 to 1.00
    memorization_evidence JSONB,
    
    -- Bias detection
    strong_in_concepts TEXT[],
    weak_in_concepts TEXT[],
    pattern_recognition_score NUMERIC(5,2),
    implementation_skill_score NUMERIC(5,2),
    problem_solving_creativity NUMERIC(5,2),
    
    -- Metadata
    duration_seconds INT,
    completed_at TIMESTAMP,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_assessments_user ON assessments(user_id);
CREATE INDEX idx_assessments_completed ON assessments(completed_at);


-- ============================================================
-- USER SKILL TRACKING (For Weakness Detection)
-- ============================================================

CREATE TABLE user_skills (
    user_id BIGINT REFERENCES users(user_id),
    topic_id BIGINT REFERENCES topics(topic_id),
    
    -- Skill metrics
    proficiency_level NUMERIC(5,2),        -- 0.00 to 100.00
    confidence_level NUMERIC(5,2),
    
    -- Evidence
    questions_attempted INT DEFAULT 0,
    questions_correct INT DEFAULT 0,
    average_time_vs_expected NUMERIC(5,2), -- Ratio: actual/expected time
    
    -- Learning curve
    improvement_rate NUMERIC(5,2),         -- How fast they're improving
    plateau_detected BOOLEAN DEFAULT false,
    needs_review BOOLEAN DEFAULT false,
    
    last_practiced_at TIMESTAMP,
    next_review_at TIMESTAMP,              -- Spaced repetition
    
    PRIMARY KEY (user_id, topic_id)
);

CREATE INDEX idx_user_skills_user ON user_skills(user_id);
CREATE INDEX idx_user_skills_weak ON user_skills(user_id, proficiency_level) 
    WHERE proficiency_level < 50;


-- ============================================================
-- VECTOR EMBEDDINGS (For RAG)
-- ============================================================

CREATE TABLE problem_embeddings (
    problem_id BIGINT PRIMARY KEY REFERENCES problems(problem_id),
    
    -- Embeddings for different aspects
    description_embedding VECTOR(1536),    -- Using pgvector extension
    solution_pattern_embedding VECTOR(768),
    constraint_embedding VECTOR(384),
    
    embedding_model VARCHAR(100),          -- Which model generated it
    embedding_version VARCHAR(20),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_prob_embedding_desc ON problem_embeddings 
    USING ivfflat (description_embedding vector_cosine_ops);


CREATE TABLE question_embeddings (
    question_id BIGINT PRIMARY KEY REFERENCES questions(question_id),
    
    question_embedding VECTOR(1536),
    concept_embedding VECTOR(768),
    
    embedding_model VARCHAR(100),
    embedding_version VARCHAR(20),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_ques_embedding ON question_embeddings 
    USING ivfflat (question_embedding vector_cosine_ops);


-- ============================================================
-- LLM GENERATION HISTORY (Track generated content)
-- ============================================================

CREATE TABLE generated_content (
    content_id BIGSERIAL PRIMARY KEY,
    
    content_type VARCHAR(50),              -- 'problem', 'question', 'solution', 'explanation'
    
    -- What was generated
    generated_text TEXT,
    generated_metadata JSONB,
    
    -- Generation context
    prompt_used TEXT,
    model_used VARCHAR(100),
    temperature NUMERIC(3,2),
    
    -- Source inspiration (for similar problems)
    based_on_problem_id BIGINT REFERENCES problems(problem_id),
    based_on_question_id BIGINT REFERENCES questions(question_id),
    
    -- Quality
    quality_score NUMERIC(5,2),            -- If reviewed
    is_approved BOOLEAN DEFAULT false,
    reviewed_by_user_id BIGINT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_generated_type ON generated_content(content_type);
CREATE INDEX idx_generated_approved ON generated_content(is_approved);
```

---

## 3. GRAPH RELATIONSHIPS (Apache AGE)

### **Why Graph Database?**

1. **Natural Relationships**: Problems, topics, and skills form a graph
2. **Transitive Queries**: "Find all prerequisites recursively"
3. **Pattern Matching**: Cypher queries for complex relationships
4. **Path Finding**: "Learning path from topic A to topic B"

### **Apache AGE Setup**

```sql
-- Enable Apache AGE extension
CREATE EXTENSION IF NOT EXISTS age;
LOAD 'age';
SET search_path = ag_catalog, "$user", public;

-- Create graph
SELECT create_graph('interview_knowledge_graph');
```

### **Graph Schema**

#### **Node Types**

```cypher
-- Problem nodes
CREATE (:Problem {
    id: problem_id,
    title: 'Two Sum',
    leetcode_number: 1,
    difficulty: 45.5,
    topics: ['hash-table', 'array']
})

-- Topic nodes  
CREATE (:Topic {
    id: topic_id,
    name: 'Hash Table',
    category: 'data_structure',
    difficulty: 30.0
})

-- Tag nodes
CREATE (:Tag {
    id: tag_id,
    name: 'two-pointers',
    category: 'pattern'
})

-- Concept nodes (abstract ideas)
CREATE (:Concept {
    id: concept_id,
    name: 'Time-Space Tradeoff',
    category: 'principle'
})

-- User skill nodes
CREATE (:UserSkill {
    user_id: user_id,
    topic_id: topic_id,
    proficiency: 75.0,
    last_practiced: timestamp
})
```

#### **Relationship Types**

```cypher
-- Problem relationships
(:Problem)-[:FOLLOW_UP_OF {difficulty_increase: 15}]->(:Problem)
(:Problem)-[:SIMILAR_TO {similarity: 0.85}]->(:Problem)
(:Problem)-[:PREREQUISITE_FOR]->(:Problem)
(:Problem)-[:EASIER_VERSION_OF]->(:Problem)
(:Problem)-[:SAME_PATTERN_AS {pattern: 'two_pointers'}]->(:Problem)
(:Problem)-[:ALTERNATIVE_APPROACH {approach: 'greedy'}]->(:Problem)

-- Topic relationships
(:Topic)-[:DEPENDS_ON]->(:Topic)
(:Topic)-[:RELATED_TO {strength: 0.7}]->(:Topic)
(:Topic)-[:INCLUDES]->(:Concept)
(:Topic)-[:PREREQUISITE_FOR]->(:Topic)

-- Problem-Topic relationships
(:Problem)-[:USES_TOPIC {relevance: 0.9}]->(:Topic)
(:Problem)-[:TAGGED_WITH]->(:Tag)
(:Problem)-[:DEMONSTRATES]->(:Concept)

-- User learning relationships
(:User)-[:HAS_SKILL]->(:UserSkill)-[:IN_TOPIC]->(:Topic)
(:User)-[:MASTERED]->(:Topic)
(:User)-[:STRUGGLING_WITH]->(:Topic)
(:User)-[:SHOULD_LEARN]->(:Topic)
```

### **Graph Queries Examples**

#### **Find Similar Problems**

```cypher
-- Find problems similar to problem 1 within 2 hops
SELECT * FROM cypher('interview_knowledge_graph', $$
    MATCH (p:Problem {id: 1})-[r:SIMILAR_TO|SAME_PATTERN_AS*1..2]-(similar:Problem)
    WHERE r.similarity > 0.7 OR r IS NULL
    RETURN DISTINCT similar.id, similar.title, similar.difficulty
    ORDER BY similar.difficulty
$$) AS (id bigint, title text, difficulty numeric);
```

#### **Find Learning Path**

```cypher
-- Find shortest learning path from Topic A to Topic B
SELECT * FROM cypher('interview_knowledge_graph', $$
    MATCH path = shortestPath(
        (start:Topic {name: 'Array'})-[:PREREQUISITE_FOR|DEPENDS_ON*]-(end:Topic {name: 'Dynamic Programming'})
    )
    RETURN [node IN nodes(path) | node.name] AS learning_path
$$) AS (learning_path text[]);
```

#### **Find User Weak Prerequisites**

```cypher
-- Find topics user should learn before tackling weak topics
SELECT * FROM cypher('interview_knowledge_graph', $$
    MATCH (user:User {id: 123})-[:STRUGGLING_WITH]->(weak:Topic)
    MATCH (weak)-[:DEPENDS_ON]->(prereq:Topic)
    WHERE NOT (user)-[:MASTERED]->(prereq)
    RETURN DISTINCT prereq.name, prereq.difficulty
    ORDER BY prereq.difficulty
$$) AS (topic_name text, difficulty numeric);
```

#### **Recommend Next Problems**

```cypher
-- Recommend problems based on mastered topics and similar patterns
SELECT * FROM cypher('interview_knowledge_graph', $$
    MATCH (user:User {id: 123})-[:MASTERED]->(topic:Topic)
    MATCH (topic)<-[:USES_TOPIC]-(problem:Problem)
    WHERE NOT (user)-[:SOLVED]->(problem)
    
    // Find similar problems to ones user solved
    OPTIONAL MATCH (user)-[:SOLVED]->(solved:Problem)-[:SIMILAR_TO {similarity: s}]->(problem)
    
    // Calculate recommendation score
    WITH problem, 
         COUNT(DISTINCT topic) AS topic_match_count,
         AVG(s) AS similarity_score,
         problem.difficulty AS difficulty
    
    // Recommend problems slightly harder than user's average
    WHERE difficulty BETWEEN user.avg_solved_difficulty - 5 AND user.avg_solved_difficulty + 10
    
    RETURN problem.id, problem.title, 
           topic_match_count * 10 + similarity_score * 20 AS recommendation_score
    ORDER BY recommendation_score DESC
    LIMIT 10
$$) AS (problem_id bigint, title text, score numeric);
```

#### **Detect Pattern Learning**

```cypher
-- Check if user is learning patterns progressively
SELECT * FROM cypher('interview_knowledge_graph', $$
    MATCH (user:User {id: 123})-[:SOLVED]->(p:Problem)-[:DEMONSTRATES]->(concept:Concept)
    WITH concept, COUNT(p) AS problems_solved, MIN(p.solved_at) AS first_seen
    ORDER BY first_seen
    RETURN concept.name, problems_solved, first_seen
$$) AS (concept text, count bigint, first_date timestamp);
```

### **Hybrid Queries (SQL + Cypher)**

```sql
-- Get problems with their graph relationships
WITH problem_graph AS (
    SELECT * FROM cypher('interview_knowledge_graph', $$
        MATCH (p:Problem {id: $problem_id})-[r:SIMILAR_TO|FOLLOW_UP_OF]-(related:Problem)
        RETURN related.id AS related_id, type(r) AS relationship_type, r.similarity AS similarity
    $$) AS (related_id bigint, relationship_type text, similarity numeric)
)
SELECT 
    p.*,
    pg.relationship_type,
    pg.similarity
FROM problems p
JOIN problem_graph pg ON p.problem_id = pg.related_id
WHERE p.is_active = true
ORDER BY pg.similarity DESC NULLS LAST;
```

---

## 4. VECTOR DATABASE & RAG

### **Why Vector DB + RAG?**

1. **Semantic Search**: Find similar problems by meaning, not just keywords
2. **Question Generation**: LLM generates new questions based on existing patterns
3. **Solution Similarity**: Match user solutions to known patterns
4. **Intelligent Recommendations**: Content-based filtering

### **Architecture**

```
┌───────────────────────────────────────────────────┐
│                 Vector Database                   │
│                  (Chroma/Qdrant)                  │
│                                                   │
│  Collections:                                     │
│  ├─ problems: Problem descriptions + solutions   │
│  ├─ questions: Question text + explanations      │
│  ├─ solutions: Code solutions + patterns         │
│  └─ concepts: Abstract concepts + relationships  │
└───────────────────────────────────────────────────┘
                         │
                         ▼
┌───────────────────────────────────────────────────┐
│                  RAG Pipeline                     │
│                                                   │
│  1. Query → Embedding                            │
│  2. Vector Search → Top K matches                │
│  3. Retrieve context from PostgreSQL             │
│  4. Augment prompt with context                  │
│  5. LLM generates response                       │
└───────────────────────────────────────────────────┘
```

### **Vector Database Schema (Chroma)**

```python
# Collection configurations

# Problems collection
problems_collection = client.create_collection(
    name="problems",
    metadata={"description": "LeetCode problems with solutions"},
    embedding_function=embedding_function
)

# Document structure for each problem
{
    "id": "problem_1",
    "document": "Problem: Two Sum. Given array of integers...",  # Full description
    "metadata": {
        "problem_id": 1,
        "title": "Two Sum",
        "difficulty_score": 45.5,
        "topics": ["hash-table", "array"],
        "patterns": ["hash-lookup", "complement-search"],
        "constraints": ["n <= 10^4", "unique solution"],
        "source": "leetcode",
        "leetcode_number": 1
    }
}

# Questions collection
questions_collection = client.create_collection(
    name="questions",
    metadata={"description": "Training questions"}
)

{
    "id": "question_123",
    "document": "What is the time complexity of...",
    "metadata": {
        "question_id": 123,
        "category": "complexity_analysis",
        "difficulty": "medium",
        "concepts": ["time-complexity", "nested-loops"],
        "related_problems": [1, 15, 18]
    }
}

# Solutions collection (for pattern matching)
solutions_collection = client.create_collection(
    name="solutions",
    metadata={"description": "Solution patterns and code"}
)

{
    "id": "solution_1_hashmap",
    "document": "Use hash map to store complements...",  # Solution explanation
    "metadata": {
        "problem_id": 1,
        "approach": "hash_map",
        "time_complexity": "O(n)",
        "space_complexity": "O(n)",
        "pattern": "hash-lookup",
        "code_snippet": "..."
    }
}
```

### **RAG Query Patterns**

#### **Pattern 1: Similar Problem Search**

```python
def find_similar_problems(problem_description: str, k: int = 5):
    """Find similar problems based on description"""
    
    # Vector search
    results = problems_collection.query(
        query_texts=[problem_description],
        n_results=k,
        include=["metadatas", "documents", "distances"]
    )
    
    # Retrieve full details from PostgreSQL
    similar_problem_ids = [m["problem_id"] for m in results["metadatas"][0]]
    
    query = """
        SELECT p.*, 
               array_agg(DISTINCT t.name) as topics,
               array_agg(DISTINCT tag.name) as tags
        FROM problems p
        LEFT JOIN problem_topics pt ON p.problem_id = pt.problem_id
        LEFT JOIN topics t ON pt.topic_id = t.topic_id
        LEFT JOIN problem_tags ptag ON p.problem_id = ptag.problem_id
        LEFT JOIN tags tag ON ptag.tag_id = tag.tag_id
        WHERE p.problem_id = ANY(%s)
        GROUP BY p.problem_id
    """
    
    return execute_query(query, (similar_problem_ids,))
```

#### **Pattern 2: Generate New Question**

```python
def generate_similar_question(question_id: int):
    """Generate new question similar to existing one"""
    
    # Get original question
    original = get_question(question_id)
    
    # Find similar questions in vector DB
    similar = questions_collection.query(
        query_texts=[original["question_text"]],
        n_results=3
    )
    
    # RAG: Build context
    context = {
        "original": original,
        "similar_questions": similar["documents"][0],
        "concepts": original["concepts"],
        "difficulty": original["difficulty_level"]
    }
    
    # LLM generates new question
    prompt = f"""
    Given these similar questions about {context['concepts']}:
    
    {chr(10).join(context['similar_questions'])}
    
    Generate a NEW question that:
    - Tests the same concepts: {context['concepts']}
    - Has {context['difficulty']} difficulty
    - Is DIFFERENT from the examples
    - Includes 4 multiple choice options
    - Provides clear explanation
    
    Output as JSON.
    """
    
    new_question = llm.generate(prompt)
    return new_question
```

#### **Pattern 3: Solution Pattern Matching**

```python
def match_solution_pattern(user_code: str, problem_id: int):
    """Match user's solution to known patterns"""
    
    # Generate embedding for user code
    user_embedding = embedding_model.embed(user_code)
    
    # Find similar solutions
    similar_solutions = solutions_collection.query(
        query_embeddings=[user_embedding],
        n_results=5,
        where={"problem_id": problem_id}  # Filter by problem
    )
    
    # Analyze which pattern user used
    patterns = [s["pattern"] for s in similar_solutions["metadatas"][0]]
    
    # Return most likely pattern
    from collections import Counter
    pattern_counts = Counter(patterns)
    most_common_pattern = pattern_counts.most_common(1)[0][0]
    
    return {
        "detected_pattern": most_common_pattern,
        "confidence": pattern_counts[most_common_pattern] / len(patterns),
        "similar_solutions": similar_solutions["documents"][0]
    }
```

#### **Pattern 4: Personalized Recommendations**

```python
def recommend_next_question(user_id: int, count: int = 5):
    """Recommend questions based on user's learning history"""
    
    # Get user's weak topics from PostgreSQL
    weak_topics = get_user_weak_topics(user_id)
    
    # Get questions user struggled with
    struggled_questions = get_struggled_questions(user_id)
    
    # Build user profile embedding
    user_profile_text = f"""
    User struggles with: {', '.join(weak_topics)}
    Needs practice in: {', '.join([q['concepts'] for q in struggled_questions])}
    Recent mistakes: {', '.join([q['common_mistakes'] for q in struggled_questions])}
    """
    
    # Vector search for relevant questions
    recommended = questions_collection.query(
        query_texts=[user_profile_text],
        n_results=count * 2,  # Get more than needed
        where={
            "$or": [
                {"concepts": {"$in": weak_topics}},
                {"difficulty": user.target_difficulty}
            ]
        }
    )
    
    # Filter out already attempted
    attempted_ids = get_attempted_question_ids(user_id)
    recommendations = [
        q for q in recommended["metadatas"][0]
        if q["question_id"] not in attempted_ids
    ][:count]
    
    return recommendations
```

### **Embedding Strategy**

```python
# Use different embeddings for different purposes

# For semantic similarity (problems/questions)
semantic_model = "text-embedding-3-small"  # OpenAI
# or
semantic_model = "all-MiniLM-L6-v2"  # Local Sentence-Transformers

# For code similarity
code_model = "code-embedding-ada-002"  # OpenAI
# or  
code_model = "microsoft/codebert-base"  # Local

# Embedding generation
def generate_embeddings(text: str, model_type: str):
    if model_type == "semantic":
        return semantic_embedding_model.encode(text)
    elif model_type == "code":
        return code_embedding_model.encode(text)
```

### **Syncing PostgreSQL ↔ Vector DB**

```python
# Trigger function to sync on insert/update

def sync_problem_to_vector_db(problem_id: int):
    """Sync problem to vector database"""
    
    # Get problem from PostgreSQL
    problem = get_problem(problem_id)
    
    # Create document text
    document = f"""
    Problem: {problem['title']}
    
    Description: {problem['description']}
    
    Constraints: {', '.join(problem['constraints'])}
    
    Topics: {', '.join(problem['topics'])}
    """
    
    # Add to vector DB
    problems_collection.add(
        ids=[f"problem_{problem_id}"],
        documents=[document],
        metadatas=[{
            "problem_id": problem_id,
            "title": problem['title'],
            "difficulty_score": problem['base_difficulty_score'],
            "topics": problem['topics'],
            "leetcode_number": problem['leetcode_number']
        }]
    )
    
    # Update embedding in PostgreSQL for quick access
    embedding = problems_collection.get(ids=[f"problem_{problem_id}"])
    
    update_query = """
        INSERT INTO problem_embeddings (problem_id, description_embedding)
        VALUES (%s, %s)
        ON CONFLICT (problem_id) 
        DO UPDATE SET description_embedding = EXCLUDED.description_embedding
    """
    execute_query(update_query, (problem_id, embedding))
```

---

## 5. DIFFICULTY SCORING SYSTEM

### **The "Magic Unit" - Multi-Dimensional Difficulty**

Traditional difficulty (Easy/Medium/Hard) is too simplistic. We use a 100-point scale with multiple dimensions.

### **Difficulty Components**

```python
class DifficultyScorer:
    """
    Multi-dimensional difficulty scoring system
    
    Total Score: 0-100 points
    """
    
    WEIGHTS = {
        "conceptual_complexity": 0.25,      # 25%
        "implementation_difficulty": 0.20,   # 20%
        "time_pressure": 0.15,               # 15%
        "edge_case_handling": 0.15,          # 15%
        "optimization_required": 0.15,        # 15%
        "debugging_difficulty": 0.10          # 10%
    }
    
    def calculate_base_difficulty(self, problem):
        """Calculate base difficulty score"""
        
        scores = {}
        
        # 1. Conceptual Complexity (0-100)
        scores['conceptual'] = self._score_conceptual(problem)
        
        # 2. Implementation Difficulty (0-100)
        scores['implementation'] = self._score_implementation(problem)
        
        # 3. Time Pressure (0-100)
        scores['time_pressure'] = self._score_time_pressure(problem)
        
        # 4. Edge Cases (0-100)
        scores['edge_cases'] = self._score_edge_cases(problem)
        
        # 5. Optimization Required (0-100)
        scores['optimization'] = self._score_optimization(problem)
        
        # 6. Debugging Difficulty (0-100)
        scores['debugging'] = self._score_debugging(problem)
        
        # Weighted average
        total_score = sum(
            scores[k.replace('_', '')] * v 
            for k, v in self.WEIGHTS.items()
        )
        
        return {
            'total_score': round(total_score, 2),
            'component_scores': scores
        }
    
    def _score_conceptual(self, problem):
        """How many concepts must you know?"""
        
        score = 0
        
        # Number of distinct topics (0-30 points)
        topic_count = len(problem['topics'])
        score += min(30, topic_count * 6)
        
        # Depth of topics (0-40 points)
        topic_difficulties = [get_topic_difficulty(t) for t in problem['topics']]
        score += min(40, sum(topic_difficulties) / len(topic_difficulties))
        
        # Pattern complexity (0-30 points)
        if requires_multiple_patterns(problem):
            score += 30
        elif requires_advanced_pattern(problem):
            score += 20
        else:
            score += 10
        
        return min(100, score)
    
    def _score_implementation(self, problem):
        """How hard to code correctly?"""
        
        score = 0
        
        # Lines of code typically required (0-25 points)
        expected_loc = estimate_lines_of_code(problem)
        score += min(25, expected_loc / 2)
        
        # Number of helper functions (0-20 points)
        helper_count = count_required_helpers(problem)
        score += min(20, helper_count * 5)
        
        # Data structure complexity (0-30 points)
        if requires_custom_ds(problem):
            score += 30
        elif requires_multiple_ds(problem):
            score += 20
        else:
            score += 10
        
        # Pointer manipulation (0-25 points)
        if requires_multiple_pointers(problem):
            score += 25
        elif requires_single_pointer(problem):
            score += 15
        
        return min(100, score)
    
    def _score_time_pressure(self, problem):
        """How quickly must you solve?"""
        
        # Based on community data
        community_avg_time = get_average_solve_time(problem)
        
        if community_avg_time < 600:  # < 10 minutes
            return 20
        elif community_avg_time < 1200:  # < 20 minutes
            return 40
        elif community_avg_time < 1800:  # < 30 minutes
            return 60
        elif community_avg_time < 2400:  # < 40 minutes
            return 80
        else:
            return 100
    
    def _score_edge_cases(self, problem):
        """How many edge cases to handle?"""
        
        score = 0
        
        # Common edge cases (each worth 10 points)
        edge_cases = [
            has_empty_input(problem),
            has_single_element(problem),
            has_all_same_elements(problem),
            has_negative_numbers(problem),
            has_overflow_risk(problem),
            has_cycle_possibility(problem),
            has_null_pointers(problem),
            has_boundary_conditions(problem),
            has_special_characters(problem),
            has_large_constraints(problem)
        ]
        
        score = sum(10 for case in edge_cases if case)
        
        return min(100, score)
    
    def _score_optimization(self, problem):
        """How critical is optimization?"""
        
        score = 0
        
        # Constraints tightness (0-50 points)
        constraints = problem['constraints']
        max_n = extract_max_input_size(constraints)
        
        if max_n >= 10**9:  # Must be O(log n) or O(1)
            score += 50
        elif max_n >= 10**6:  # Must be O(n) or better
            score += 40
        elif max_n >= 10**4:  # Must be O(n log n) or better
            score += 30
        else:
            score += 20
        
        # Gap between naive and optimal (0-50 points)
        naive_complexity = get_naive_complexity(problem)
        optimal_complexity = get_optimal_complexity(problem)
        complexity_gap = calculate_complexity_gap(naive_complexity, optimal_complexity)
        
        if complexity_gap >= 3:  # e.g., O(n^3) to O(n)
            score += 50
        elif complexity_gap >= 2:  # e.g., O(n^2) to O(n log n)
            score += 35
        elif complexity_gap >= 1:  # e.g., O(n log n) to O(n)
            score += 20
        
        return min(100, score)
    
    def _score_debugging(self, problem):
        """How hard to debug when wrong?"""
        
        score = 0
        
        # Off-by-one risk (0-30 points)
        if high_off_by_one_risk(problem):
            score += 30
        
        # Hidden bugs (0-30 points)
        if has_subtle_edge_cases(problem):
            score += 30
        
        # Complex state (0-40 points)
        state_variables = count_state_variables(problem)
        score += min(40, state_variables * 10)
        
        return min(100, score)
```

### **Difficulty Calibration**

```python
def calibrate_difficulty_from_attempts(problem_id: int):
    """
    Adjust difficulty based on actual user performance
    (Collaborative filtering approach)
    """
    
    query = """
        SELECT 
            AVG(is_correct::int) as success_rate,
            AVG(time_spent_seconds) as avg_time,
            COUNT(*) as attempt_count,
            COUNT(DISTINCT user_id) as unique_users,
            AVG(attempt_number) as avg_attempts_to_solve
        FROM question_attempts
        WHERE question_id = %s
    """
    
    stats = execute_query(query, (problem_id,))[0]
    
    # Adjust base difficulty
    base_score = get_base_difficulty(problem_id)
    
    # Success rate adjustment (-20 to +20 points)
    if stats['success_rate'] < 0.2:  # Very hard
        adjustment = +20
    elif stats['success_rate'] < 0.4:
        adjustment = +10
    elif stats['success_rate'] > 0.8:  # Very easy
        adjustment = -20
    elif stats['success_rate'] > 0.6:
        adjustment = -10
    else:
        adjustment = 0
    
    # Time adjustment (-10 to +10 points)
    expected_time = get_expected_time(problem_id)
    time_ratio = stats['avg_time'] / expected_time
    
    if time_ratio > 1.5:
        adjustment += 10
    elif time_ratio > 1.2:
        adjustment += 5
    elif time_ratio < 0.7:
        adjustment -= 10
    elif time_ratio < 0.85:
        adjustment -= 5
    
    # Final calibrated score
    calibrated_score = max(0, min(100, base_score + adjustment))
    
    # Update community difficulty
    update_query = """
        UPDATE problems
        SET community_difficulty_score = %s,
            updated_at = CURRENT_TIMESTAMP
        WHERE problem_id = %s
    """
    execute_query(update_query, (calibrated_score, problem_id))
    
    return calibrated_score
```

### **Personalized Difficulty**

```python
def calculate_personalized_difficulty(problem_id: int, user_id: int):
    """
    Adjust difficulty based on user's specific skills
    """
    
    # Get problem requirements
    problem_topics = get_problem_topics(problem_id)
    
    # Get user proficiency in those topics
    user_skills = get_user_skills(user_id, problem_topics)
    
    # Base difficulty
    base_score = get_base_difficulty(problem_id)
    
    # Calculate skill-adjusted score
    avg_proficiency = sum(s['proficiency_level'] for s in user_skills) / len(user_skills)
    
    # If user is strong in these topics, problem is easier for them
    # Scale: 100 proficiency = -30 points, 0 proficiency = +30 points
    adjustment = (50 - avg_proficiency) * 0.6
    
    personalized_score = max(0, min(100, base_score + adjustment))
    
    return {
        'base_difficulty': base_score,
        'personalized_difficulty': personalized_score,
        'user_proficiency': avg_proficiency,
        'adjustment': adjustment
    }
```

### **Difficulty Ranges**

```
Score Range | Label          | Description
------------|----------------|----------------------------------
0-20        | Trivial        | Basic application, no tricks
21-35       | Easy           | Single concept, straightforward
36-50       | Medium-Easy    | 1-2 concepts, some thinking
51-65       | Medium         | Multiple concepts, optimization
66-75       | Medium-Hard    | Complex patterns, edge cases
76-85       | Hard           | Advanced algorithms, difficult
86-95       | Very Hard      | Research-level, creative
96-100      | Expert         | Competition-level, novel
```

---

## 6. ASSESSMENT FRAMEWORK

### **LLM-Powered Assessment Goals**

1. **Detect Memorization** - Did they memorize or truly understand?
2. **Identify Real Weaknesses** - Beyond just "got it wrong"
3. **Fair Evaluation** - Account for different problem-solving styles
4. **Personalized Feedback** - Specific actionable advice

### **Assessment Types**

```python
class AssessmentEngine:
    """LLM-powered assessment system"""
    
    def assess_solution(self, user_id: int, problem_id: int, solution_code: str):
        """Comprehensive solution assessment"""
        
        # Gather context
        context = self._build_assessment_context(user_id, problem_id, solution_code)
        
        # LLM assessment
        assessment = self._llm_assess(context)
        
        # Store results
        self._store_assessment(user_id, problem_id, assessment)
        
        return assessment
    
    def _build_assessment_context(self, user_id, problem_id, solution_code):
        """Gather all relevant context for LLM"""
        
        return {
            # Problem details
            'problem': get_problem(problem_id),
            'expected_complexity': get_optimal_complexity(problem_id),
            'common_approaches': get_common_approaches(problem_id),
            
            # User history
            'user_history': self._get_user_context(user_id, problem_id),
            
            # Solution analysis
            'code': solution_code,
            'detected_pattern': self._detect_pattern(solution_code),
            'complexity_analysis': self._analyze_complexity(solution_code),
            
            # Similar solutions
            'similar_solutions': self._find_similar_solutions(problem_id, solution_code)
        }
    
    def _get_user_context(self, user_id, problem_id):
        """Get user's relevant history"""
        
        # Previous attempts at this problem
        previous_attempts = get_previous_attempts(user_id, problem_id)
        
        # Similar problems solved
        similar_solved = get_similar_problems_solved(user_id, problem_id)
        
        # User skill levels in relevant topics
        problem_topics = get_problem_topics(problem_id)
        user_skills = get_user_skills(user_id, problem_topics)
        
        # Recent learning activity
        recent_questions = get_recent_questions_attempted(user_id, limit=20)
        
        return {
            'previous_attempts': previous_attempts,
            'similar_solved': similar_solved,
            'skills': user_skills,
            'recent_activity': recent_questions,
            'learning_streak': get_learning_streak(user_id),
            'struggle_patterns': identify_struggle_patterns(user_id)
        }
    
    def _llm_assess(self, context):
        """LLM performs multi-faceted assessment"""
        
        prompt = f"""
        You are an expert programming interview assessor. Analyze this solution comprehensively.
        
        PROBLEM:
        {context['problem']['description']}
        
        Expected: {context['expected_complexity']} time, optimal approach uses {context['common_approaches']}
        
        USER SOLUTION:
        ```
        {context['code']}
        ```
        
        CONTEXT:
        - User has solved {len(context['user_history']['similar_solved'])} similar problems
        - Previous attempts at this problem: {len(context['user_history']['previous_attempts'])}
        - User skill level in required topics: {context['user_history']['skills']}
        - Recent struggle patterns: {context['user_history']['struggle_patterns']}
        
        ASSESS THE SOLUTION ON:
        
        1. CORRECTNESS (0-100)
           - Does it solve the problem?
           - Handles all edge cases?
           - Any bugs?
        
        2. COMPLEXITY (0-100)
           - Actual time complexity?
           - Actual space complexity?
           - Is it optimal?
        
        3. CODE QUALITY (0-100)
           - Readability
           - Variable naming
           - Structure
           - Comments (if needed)
        
        4. MEMORIZATION DETECTION (0-1)
           Evidence of memorization:
           - Identical to known solutions?
           - Solved immediately after viewing similar problem?
           - No reasoning process evident?
           - Pattern match with their history?
           
           Return memorization likelihood: 0.0 (original) to 1.0 (clearly memorized)
        
        5. UNDERSTANDING DEPTH (0-100)
           - Do they understand WHY this works?
           - Can they explain trade-offs?
           - Do they know alternatives?
        
        6. PROBLEM-SOLVING APPROACH (qualitative)
           - Identified the pattern correctly?
           - Started with right approach?
           - Optimized efficiently?
        
        7. SPECIFIC WEAKNESSES DETECTED
           List specific areas where this user needs improvement:
           - Conceptual gaps?
           - Implementation issues?
           - Optimization blindspots?
        
        8. PERSONALIZED FEEDBACK
           Provide 3-5 specific, actionable recommendations for THIS user based on:
           - Their history
           - This attempt
           - Their skill profile
        
        9. NEXT STEPS
           Recommend:
           - Similar problems to solidify learning
           - Prerequisite topics to review
           - Advanced problems when ready
        
        Return as structured JSON.
        """
        
        response = llm.generate(prompt, temperature=0.3)  # Low temp for consistency
        
        return json.loads(response)
    
    def detect_memorization(self, context):
        """Specific memorization detection logic"""
        
        memorization_signals = []
        confidence = 0.0
        
        # Signal 1: Solution extremely similar to known solutions
        similar_solutions = context['similar_solutions']
        if similar_solutions:
            max_similarity = max(s['similarity_score'] for s in similar_solutions)
            if max_similarity > 0.95:
                memorization_signals.append("code_similarity")
                confidence += 0.4
        
        # Signal 2: Solved very quickly relative to difficulty
        problem_difficulty = context['problem']['base_difficulty_score']
        user_proficiency = context['user_history']['skills']['avg_proficiency']
        time_taken = context.get('time_spent_seconds', 0)
        
        expected_time = calculate_expected_time(problem_difficulty, user_proficiency)
        if time_taken < expected_time * 0.3:  # Solved in < 30% of expected time
            memorization_signals.append("suspiciously_fast")
            confidence += 0.3
        
        # Signal 3: Recently viewed similar problem
        recent_activity = context['user_history']['recent_activity']
        similar_problem_ids = [p['problem_id'] for p in context['similar_solved']]
        
        for activity in recent_activity[-5:]:  # Last 5 activities
            if activity['problem_id'] in similar_problem_ids:
                memorization_signals.append("recent_similar_view")
                confidence += 0.2
                break
        
        # Signal 4: No incorrect attempts (suspicious if problem is hard)
        if problem_difficulty > 70 and len(context['user_history']['previous_attempts']) == 0:
            memorization_signals.append("no_trial_and_error")
            confidence += 0.1
        
        return {
            'memorization_likelihood': min(1.0, confidence),
            'signals': memorization_signals,
            'is_likely_memorized': confidence > 0.6
        }
    
    def identify_weaknesses(self, user_id: int, assessment_history: list):
        """Aggregate assessments to identify patterns of weakness"""
        
        # Collect all identified weaknesses
        all_weaknesses = []
        for assessment in assessment_history:
            all_weaknesses.extend(assessment['weaknesses'])
        
        # Cluster and count
        from collections import Counter
        weakness_counts = Counter(all_weaknesses)
        
        # Identify persistent weaknesses (appeared multiple times)
        persistent = [w for w, count in weakness_counts.items() if count >= 3]
        
        # Categorize by type
        categorized = {
            'conceptual': [],      # Don't understand the concept
            'implementation': [],  # Know what to do, struggle to code
            'optimization': [],    # Can solve but not optimally
            'edge_cases': [],      # Miss edge cases
            'debugging': []        # Make bugs frequently
        }
        
        for weakness in persistent:
            category = classify_weakness(weakness)
            categorized[category].append(weakness)
        
        # Generate recommendations
        recommendations = self._generate_weakness_recommendations(categorized)
        
        return {
            'weaknesses': categorized,
            'recommendations': recommendations,
            'priority_areas': list(weakness_counts.most_common(3))
        }
```

### **Anti-Memorization Techniques**

```python
class AntiMemorizationStrategies:
    """Techniques to test understanding beyond memorization"""
    
    def generate_variant_question(self, original_question_id: int):
        """Generate slightly modified version of question"""
        
        original = get_question(original_question_id)
        
        # RAG: Get similar questions for inspiration
        similar = questions_collection.query(
            query_texts=[original['question_text']],
            n_results=3
        )
        
        prompt = f"""
        Original question (user may have memorized):
        {original['question_text']}
        
        Create a VARIANT that:
        1. Tests the SAME concept
        2. Changes surface details (numbers, context, variable names)
        3. Adds a twist that requires understanding
        4. Cannot be solved by pattern matching alone
        
        Example transformations:
        - Change array to linked list
        - Change "find minimum" to "find maximum"
        - Add constraint like "without extra space"
        - Change data type (strings to integers)
        - Reverse the problem (instead of "is valid?", ask "make it valid")
        
        Generate 3 variants with increasing difficulty.
        """
        
        variants = llm.generate(prompt)
        return variants
    
    def test_deeper_understanding(self, problem_id: int, user_id: int):
        """Ask follow-up questions to test true understanding"""
        
        problem = get_problem(problem_id)
        user_solution = get_latest_solution(user_id, problem_id)
        
        follow_up_questions = [
            # Complexity understanding
            {
                "type": "complexity",
                "question": f"Why is your solution {user_solution['complexity']}? Explain the analysis.",
                "expected": "Should reference loop structure, recursive calls, etc."
            },
            
            # Trade-offs
            {
                "type": "trade-off",
                "question": "What is the space-time tradeoff in this problem?",
                "expected": "Should understand optimization options"
            },
            
            # Alternative approaches
            {
                "type": "alternatives",
                "question": "Describe 2 other ways to solve this. What are pros/cons?",
                "expected": "Should know multiple approaches exist"
            },
            
            # Edge cases
            {
                "type": "edge_cases",
                "question": "What edge cases did you consider? How does your solution handle them?",
                "expected": "Should list specific cases and reasoning"
            },
            
            # Modification challenge
            {
                "type": "modification",
                "question": f"If constraint changed from {problem['constraints'][0]} to {generate_harder_constraint()}, how would you modify?",
                "expected": "Should adapt solution, not rebuild from scratch"
            },
            
            # Debugging
            {
                "type": "debugging",
                "question": f"Here's a buggy solution: {generate_buggy_solution(problem_id)}. What's wrong?",
                "expected": "Should identify bug and explain fix"
            }
        ]
        
        return follow_up_questions
    
    def adaptive_difficulty_adjustment(self, user_id: int):
        """If user consistently solves quickly, increase difficulty"""
        
        recent_attempts = get_recent_attempts(user_id, limit=10)
        
        # Calculate suspicion score
        suspicion = 0
        
        for attempt in recent_attempts:
            if attempt['time_spent'] < attempt['expected_time'] * 0.5:
                suspicion += 1
            if attempt['is_correct'] and attempt['attempt_number'] == 1:
                suspicion += 0.5
        
        # If suspicious, inject harder variants
        if suspicion > 5:
            return {
                'adjust_difficulty': True,
                'strategy': 'inject_novel_problems',
                'reason': 'User solving too consistently fast'
            }
        
        return {'adjust_difficulty': False}
```

---

*[Document continues with sections 7-11 covering Search/Filtering, Training Plan Builder, LLM Integration, Local Deployment, and API Design. Would you like me to continue with those sections?]*
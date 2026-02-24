# Master Topic Reference Guide
## Complete DSA Taxonomy for Database Seeding, Graph Construction & RAG Pipeline

---

## Document Purpose

This document serves as the **single source of truth** for:

1. **Database Seeding:** PostgreSQL tables population
2. **Graph Construction:** Apache AGE nodes and edges creation  
3. **Embeddings Generation:** Vector representations for semantic search
4. **RAG Pipeline:** Context retrieval and AI hint generation
5. **n8n Workflows:** Automated content generation and relationship mapping

---

## Table of Contents

1. [Complete Topic Taxonomy](#taxonomy)
2. [Database Schema Reference](#database-schema)
3. [Graph Construction Rules](#graph-construction)
4. [Embeddings Strategy](#embeddings-strategy)
5. [RAG Pipeline Configuration](#rag-pipeline)
6. [n8n Workflow Integration](#n8n-integration)
7. [Topic Details Template](#topic-details)
8. [Implementation Checklist](#implementation-checklist)

---

## 1. COMPLETE TOPIC TAXONOMY {#taxonomy}

### Hierarchical Structure

```
DSA_ROOT
├── 1. FUNDAMENTALS
│   ├── 1.1 Time Complexity
│   ├── 1.2 Space Complexity
│   ├── 1.3 Big O Notation
│   ├── 1.4 Recursion Basics
│   └── 1.5 Iteration vs Recursion
│
├── 2. ARRAYS & STRINGS
│   ├── 2.1 Array Basics
│   │   ├── 2.1.1 Traversal
│   │   ├── 2.1.2 Insertion & Deletion
│   │   └── 2.1.3 Array Manipulation
│   ├── 2.2 Two Pointers
│   │   ├── 2.2.1 Opposite Direction
│   │   ├── 2.2.2 Same Direction
│   │   └── 2.2.3 Fast & Slow Pointers
│   ├── 2.3 Sliding Window
│   │   ├── 2.3.1 Fixed Size Window
│   │   ├── 2.3.2 Variable Size Window
│   │   └── 2.3.3 Shrinkable Window
│   ├── 2.4 Prefix Sum
│   │   ├── 2.4.1 1D Prefix Sum
│   │   ├── 2.4.2 2D Prefix Sum
│   │   └── 2.4.3 Prefix Sum with Hash Map
│   ├── 2.5 Kadane's Algorithm
│   ├── 2.6 Dutch National Flag
│   ├── 2.7 String Manipulation
│   │   ├── 2.7.1 String Matching
│   │   ├── 2.7.2 KMP Algorithm
│   │   └── 2.7.3 Rabin-Karp
│   └── 2.8 Palindromes
│
├── 3. LINKED LISTS
│   ├── 3.1 Singly Linked List
│   ├── 3.2 Doubly Linked List
│   ├── 3.3 Circular Linked List
│   ├── 3.4 Fast & Slow Pointers
│   ├── 3.5 Cycle Detection
│   ├── 3.6 Reversal Techniques
│   └── 3.7 Merge Operations
│
├── 4. STACKS & QUEUES
│   ├── 4.1 Stack Basics
│   ├── 4.2 Monotonic Stack
│   │   ├── 4.2.1 Monotonic Increasing
│   │   ├── 4.2.2 Monotonic Decreasing
│   │   └── 4.2.3 Next Greater/Smaller Element
│   ├── 4.3 Queue Basics
│   ├── 4.4 Deque (Double-ended Queue)
│   ├── 4.5 Monotonic Deque
│   ├── 4.6 Priority Queue (Heap)
│   └── 4.7 Expression Evaluation
│
├── 5. HASH TABLES & SETS
│   ├── 5.1 Hash Map Basics
│   ├── 5.2 Hash Set
│   ├── 5.3 Counting Patterns
│   ├── 5.4 Two Sum Variations
│   └── 5.5 Collision Handling
│
├── 6. TREES
│   ├── 6.1 Binary Trees
│   │   ├── 6.1.1 Tree Traversals
│   │   │   ├── Inorder (DFS)
│   │   │   ├── Preorder (DFS)
│   │   │   ├── Postorder (DFS)
│   │   │   └── Level Order (BFS)
│   │   ├── 6.1.2 Tree Construction
│   │   └── 6.1.3 Tree Properties
│   ├── 6.2 Binary Search Trees (BST)
│   │   ├── 6.2.1 BST Operations
│   │   ├── 6.2.2 BST Validation
│   │   ├── 6.2.3 Inorder Successor
│   │   └── 6.2.4 BST to Sorted Array
│   ├── 6.3 Balanced Trees
│   │   ├── 6.3.1 AVL Trees
│   │   └── 6.3.2 Red-Black Trees
│   ├── 6.4 Tries (Prefix Trees)
│   │   ├── 6.4.1 Trie Construction
│   │   ├── 6.4.2 Word Search
│   │   └── 6.4.3 Autocomplete
│   ├── 6.5 Segment Trees
│   ├── 6.6 Fenwick Trees (Binary Indexed Tree)
│   └── 6.7 Tree DP
│       ├── 6.7.1 Path Sum Problems
│       ├── 6.7.2 Diameter Problems
│       └── 6.7.3 Subtree Problems
│
├── 7. HEAPS
│   ├── 7.1 Min Heap
│   ├── 7.2 Max Heap
│   ├── 7.3 K-way Merge
│   ├── 7.4 Top K Elements
│   ├── 7.5 Median of Stream
│   └── 7.6 Heap Sort
│
├── 8. GRAPHS
│   ├── 8.1 Graph Representations
│   │   ├── 8.1.1 Adjacency Matrix
│   │   ├── 8.1.2 Adjacency List
│   │   └── 8.1.3 Edge List
│   ├── 8.2 Graph Traversal
│   │   ├── 8.2.1 DFS (Depth-First Search)
│   │   │   ├── DFS Recursive
│   │   │   ├── DFS Iterative
│   │   │   └── DFS Applications
│   │   └── 8.2.2 BFS (Breadth-First Search)
│   │       ├── BFS Single Source
│   │       ├── BFS Multi-Source
│   │       └── BFS Applications
│   ├── 8.3 Shortest Path Algorithms
│   │   ├── 8.3.1 Dijkstra's Algorithm
│   │   ├── 8.3.2 Bellman-Ford
│   │   ├── 8.3.3 Floyd-Warshall
│   │   └── 8.3.4 A* Search
│   ├── 8.4 Minimum Spanning Tree
│   │   ├── 8.4.1 Kruskal's Algorithm
│   │   └── 8.4.2 Prim's Algorithm
│   ├── 8.5 Topological Sort
│   │   ├── 8.5.1 Kahn's Algorithm (BFS)
│   │   └── 8.5.2 DFS-based Topological Sort
│   ├── 8.6 Union-Find (Disjoint Set)
│   │   ├── 8.6.1 Basic Union-Find
│   │   ├── 8.6.2 Path Compression
│   │   └── 8.6.3 Union by Rank
│   ├── 8.7 Cycle Detection
│   │   ├── 8.7.1 Undirected Graph Cycle
│   │   └── 8.7.2 Directed Graph Cycle
│   ├── 8.8 Graph Coloring
│   ├── 8.9 Bipartite Graphs
│   └── 8.10 Strongly Connected Components
│       ├── 8.10.1 Kosaraju's Algorithm
│       └── 8.10.2 Tarjan's Algorithm
│
├── 9. DYNAMIC PROGRAMMING
│   ├── 9.1 DP Fundamentals
│   │   ├── 9.1.1 Overlapping Subproblems
│   │   ├── 9.1.2 Optimal Substructure
│   │   ├── 9.1.3 Memoization (Top-Down)
│   │   └── 9.1.4 Tabulation (Bottom-Up)
│   ├── 9.2 1D DP
│   │   ├── 9.2.1 Fibonacci Sequence
│   │   ├── 9.2.2 Climbing Stairs
│   │   ├── 9.2.3 House Robber
│   │   └── 9.2.4 Decode Ways
│   ├── 9.3 2D DP
│   │   ├── 9.3.1 Grid Path Problems
│   │   ├── 9.3.2 Longest Common Subsequence
│   │   ├── 9.3.3 Edit Distance
│   │   └── 9.3.4 Dungeon Game
│   ├── 9.4 Knapsack Problems
│   │   ├── 9.4.1 0/1 Knapsack
│   │   ├── 9.4.2 Unbounded Knapsack
│   │   ├── 9.4.3 Subset Sum
│   │   └── 9.4.4 Partition Equal Subset
│   ├── 9.5 Subsequence DP
│   │   ├── 9.5.1 Longest Increasing Subsequence
│   │   ├── 9.5.2 Longest Common Subsequence
│   │   └── 9.5.3 Longest Palindromic Subsequence
│   ├── 9.6 String DP
│   │   ├── 9.6.1 Longest Palindromic Substring
│   │   ├── 9.6.2 Word Break
│   │   └── 9.6.3 Regular Expression Matching
│   ├── 9.7 State Machine DP
│   ├── 9.8 Interval DP
│   ├── 9.9 Tree DP
│   ├── 9.10 Digit DP
│   ├── 9.11 Bitmask DP
│   └── 9.12 DP Optimizations
│       ├── 9.12.1 Space Optimization
│       ├── 9.12.2 Rolling Array
│       └── 9.12.3 State Compression
│
├── 10. GREEDY ALGORITHMS
│   ├── 10.1 Greedy Fundamentals
│   ├── 10.2 Interval Problems
│   │   ├── 10.2.1 Meeting Rooms
│   │   ├── 10.2.2 Merge Intervals
│   │   └── 10.2.3 Non-overlapping Intervals
│   ├── 10.3 Scheduling Problems
│   ├── 10.4 Huffman Coding
│   └── 10.5 Gas Station
│
├── 11. BACKTRACKING
│   ├── 11.1 Backtracking Fundamentals
│   ├── 11.2 Permutations
│   ├── 11.3 Combinations
│   ├── 11.4 Subsets
│   ├── 11.5 N-Queens
│   ├── 11.6 Sudoku Solver
│   ├── 11.7 Word Search
│   └── 11.8 Partition Problems
│
├── 12. BINARY SEARCH
│   ├── 12.1 Basic Binary Search
│   ├── 12.2 Binary Search on Answer
│   ├── 12.3 Search in Rotated Array
│   ├── 12.4 Finding Boundaries
│   ├── 12.5 Median of Two Sorted Arrays
│   └── 12.6 Capacity to Ship Packages
│
├── 13. BIT MANIPULATION
│   ├── 13.1 Bitwise Operators
│   ├── 13.2 Bit Tricks
│   ├── 13.3 Single Number Problems
│   ├── 13.4 Counting Bits
│   └── 13.5 Bit Masking
│
├── 14. MATH & NUMBER THEORY
│   ├── 14.1 Prime Numbers
│   ├── 14.2 GCD & LCM
│   ├── 14.3 Fast Exponentiation
│   ├── 14.4 Modular Arithmetic
│   ├── 14.5 Combinatorics
│   └── 14.6 Probability
│
└── 15. ADVANCED TOPICS
    ├── 15.1 Advanced Graph Algorithms
    ├── 15.2 Advanced DP Techniques
    ├── 15.3 Computational Geometry
    ├── 15.4 String Algorithms (Advanced)
    └── 15.5 Game Theory
```

---

## 2. DATABASE SCHEMA REFERENCE {#database-schema}

### PostgreSQL Tables Structure

#### 2.1 Core Entity Tables

```sql
-- Topics Table
CREATE TABLE topics (
    topic_id SERIAL PRIMARY KEY,
    topic_code VARCHAR(20) UNIQUE NOT NULL, -- e.g., "2.3.2"
    topic_name VARCHAR(200) NOT NULL,
    parent_topic_id INTEGER REFERENCES topics(topic_id),
    depth_level INTEGER NOT NULL, -- 1=main, 2=sub, 3=sub-sub
    difficulty_range VARCHAR(20), -- "easy-medium", "medium-hard"
    description TEXT,
    key_concepts TEXT[],
    estimated_practice_hours FLOAT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Patterns Table
CREATE TABLE patterns (
    pattern_id SERIAL PRIMARY KEY,
    pattern_name VARCHAR(100) UNIQUE NOT NULL,
    pattern_type VARCHAR(50), -- "technique", "template", "approach"
    description TEXT,
    time_complexity VARCHAR(50),
    space_complexity VARCHAR(50),
    when_to_use TEXT,
    code_template TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Problems Table  
CREATE TABLE problems (
    problem_id SERIAL PRIMARY KEY,
    problem_title VARCHAR(300) NOT NULL,
    problem_slug VARCHAR(300) UNIQUE,
    difficulty INTEGER CHECK (difficulty BETWEEN 1 AND 100),
    leetcode_id INTEGER,
    description TEXT,
    constraints TEXT,
    examples JSONB,
    hints TEXT[],
    solution_approaches JSONB,
    edge_cases TEXT[],
    common_mistakes TEXT[],
    created_at TIMESTAMP DEFAULT NOW()
);

-- Questions Table
CREATE TABLE questions (
    question_id SERIAL PRIMARY KEY,
    question_type VARCHAR(50), -- "multiple_choice", "code_debug", "observation"
    question_text TEXT NOT NULL,
    options JSONB, -- For MC: {"A": "text", "B": "text", ...}
    correct_answer VARCHAR(500),
    explanation TEXT,
    difficulty INTEGER CHECK (difficulty BETWEEN 1 AND 100),
    cognitive_skill VARCHAR(100), -- "observation", "debugging", "optimization"
    estimated_time_seconds INTEGER,
    hints JSONB, -- {"level_1": "hint", "level_2": "hint", ...}
    related_problem_id INTEGER REFERENCES problems(problem_id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Common Mistakes Table
CREATE TABLE common_mistakes (
    mistake_id SERIAL PRIMARY KEY,
    mistake_category VARCHAR(100),
    mistake_name VARCHAR(200),
    description TEXT,
    example_code TEXT,
    correct_code TEXT,
    how_to_avoid TEXT,
    severity VARCHAR(20), -- "critical", "major", "minor"
    created_at TIMESTAMP DEFAULT NOW()
);

-- Edge Cases Table
CREATE TABLE edge_cases (
    edge_case_id SERIAL PRIMARY KEY,
    edge_case_category VARCHAR(100),
    edge_case_name VARCHAR(200),
    description TEXT,
    test_input JSONB,
    expected_behavior TEXT,
    why_important TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Prerequisites Table
CREATE TABLE prerequisites (
    prerequisite_id SERIAL PRIMARY KEY,
    target_topic_id INTEGER REFERENCES topics(topic_id),
    required_topic_id INTEGER REFERENCES topics(topic_id),
    strength VARCHAR(20), -- "strong", "moderate", "weak"
    reason TEXT,
    UNIQUE(target_topic_id, required_topic_id)
);
```

#### 2.2 Relationship Tables (Many-to-Many)

```sql
-- Topic-Pattern Relationships
CREATE TABLE topic_patterns (
    topic_id INTEGER REFERENCES topics(topic_id),
    pattern_id INTEGER REFERENCES patterns(pattern_id),
    relevance_score FLOAT CHECK (relevance_score BETWEEN 0 AND 1),
    PRIMARY KEY (topic_id, pattern_id)
);

-- Problem-Topic Relationships
CREATE TABLE problem_topics (
    problem_id INTEGER REFERENCES problems(problem_id),
    topic_id INTEGER REFERENCES topics(topic_id),
    is_primary BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (problem_id, topic_id)
);

-- Problem-Pattern Relationships
CREATE TABLE problem_patterns (
    problem_id INTEGER REFERENCES problems(problem_id),
    pattern_id INTEGER REFERENCES patterns(pattern_id),
    PRIMARY KEY (problem_id, pattern_id)
);

-- Question-Topic Relationships
CREATE TABLE question_topics (
    question_id INTEGER REFERENCES questions(question_id),
    topic_id INTEGER REFERENCES topics(topic_id),
    PRIMARY KEY (question_id, topic_id)
);

-- Problem-Mistake Relationships
CREATE TABLE problem_mistakes (
    problem_id INTEGER REFERENCES problems(problem_id),
    mistake_id INTEGER REFERENCES common_mistakes(mistake_id),
    frequency VARCHAR(20), -- "very_common", "common", "occasional", "rare"
    PRIMARY KEY (problem_id, mistake_id)
);

-- Problem-EdgeCase Relationships
CREATE TABLE problem_edge_cases (
    problem_id INTEGER REFERENCES problems(problem_id),
    edge_case_id INTEGER REFERENCES edge_cases(edge_case_id),
    must_test BOOLEAN DEFAULT TRUE,
    PRIMARY KEY (problem_id, edge_case_id)
);
```

#### 2.3 User Progress Tables

```sql
-- User Progress
CREATE TABLE user_progress (
    user_id INTEGER,
    topic_id INTEGER REFERENCES topics(topic_id),
    mastery_level FLOAT CHECK (mastery_level BETWEEN 0 AND 1),
    problems_attempted INTEGER DEFAULT 0,
    problems_solved INTEGER DEFAULT 0,
    questions_attempted INTEGER DEFAULT 0,
    questions_correct INTEGER DEFAULT 0,
    last_practiced TIMESTAMP,
    PRIMARY KEY (user_id, topic_id)
);

-- User Weaknesses
CREATE TABLE user_weaknesses (
    user_id INTEGER,
    weakness_category VARCHAR(100),
    weakness_description TEXT,
    occurrences INTEGER DEFAULT 1,
    last_occurred TIMESTAMP,
    recommended_practice JSONB, -- question_ids or problem_ids
    PRIMARY KEY (user_id, weakness_category)
);

-- Question Attempts
CREATE TABLE question_attempts (
    attempt_id SERIAL PRIMARY KEY,
    user_id INTEGER,
    question_id INTEGER REFERENCES questions(question_id),
    user_answer TEXT,
    is_correct BOOLEAN,
    time_taken_seconds INTEGER,
    hints_used JSONB, -- Which hints revealed
    attempt_timestamp TIMESTAMP DEFAULT NOW()
);
```

---

## 3. GRAPH CONSTRUCTION RULES {#graph-construction}

### Apache AGE Graph Schema

#### 3.1 Node Types & Properties

```cypher
// Topic Node
CREATE (t:Topic {
    topic_id: 123,
    topic_code: "2.3.2",
    name: "Variable Size Sliding Window",
    difficulty_range: "medium-hard",
    depth_level: 3,
    description: "...",
    key_concepts: ["window expansion", "shrinking condition"],
    estimated_practice_hours: 8.0
})

// Pattern Node
CREATE (p:Pattern {
    pattern_id: 45,
    name: "Two Pointers - Opposite Direction",
    time_complexity: "O(n)",
    space_complexity: "O(1)",
    template: "int left = 0, right = n-1; ...",
    when_to_use: "When input is sorted"
})

// Problem Node
CREATE (prob:Problem {
    problem_id: 789,
    title: "Container With Most Water",
    leetcode_id: 11,
    difficulty: 45,
    constraints: "...",
    edge_cases: ["all same height", "strictly increasing"],
    common_mistakes: ["forgetting to move both pointers"]
})

// Question Node
CREATE (q:Question {
    question_id: 1024,
    type: "multiple_choice",
    question_text: "What should the base case return?",
    cognitive_skill: "base_case_understanding",
    difficulty: 30
})

// Mistake Node
CREATE (m:Mistake {
    mistake_id: 55,
    category: "integer_overflow",
    name: "Unchecked multiplication overflow",
    severity: "critical"
})

// EdgeCase Node
CREATE (ec:EdgeCase {
    edge_case_id: 88,
    category: "empty_input",
    name: "Empty array",
    why_important: "Causes array[0] crash"
})
```

#### 3.2 Edge Types (Relationships)

```cypher
// Prerequisite relationship
(Topic)-[:PREREQUISITE {strength: "strong", reason: "..."}]->(Topic)

// Topic uses Pattern
(Topic)-[:USES_PATTERN {frequency: "always"}]->(Pattern)

// Problem belongs to Topic
(Problem)-[:BELONGS_TO {is_primary: true}]->(Topic)

// Problem uses Pattern
(Problem)-[:USES_PATTERN]->(Pattern)

// Question tests Topic
(Question)-[:TESTS {cognitive_skill: "observation"}]->(Topic)

// Question related to Problem
(Question)-[:RELATED_TO]->(Problem)

// Problem has common Mistake
(Problem)-[:HAS_MISTAKE {frequency: "very_common"}]->(Mistake)

// Problem requires testing EdgeCase
(Problem)-[:REQUIRES_TEST]->(EdgeCase)

// Topic has common Mistake
(Topic)-[:COMMON_MISTAKE]->(Mistake)

// Pattern fails on EdgeCase
(Pattern)-[:FAILS_ON]->(EdgeCase)

// Learning Path
(Topic)-[:NEXT_RECOMMENDED {difficulty_increase: 10}]->(Topic)

// Similar Problems
(Problem)-[:SIMILAR_TO {similarity_score: 0.85}]->(Problem)
```

#### 3.3 Cypher Query Patterns

```cypher
// Find all prerequisites for a topic (recursive)
MATCH path = (target:Topic {topic_code: "9.4.1"})-[:PREREQUISITE*]->(prereq:Topic)
RETURN prereq.name, length(path) as depth
ORDER BY depth

// Find patterns used in topic
MATCH (t:Topic {name: "Sliding Window"})-[:USES_PATTERN]->(p:Pattern)
RETURN p.name, p.time_complexity

// Find problems testing topic + pattern
MATCH (p:Problem)-[:BELONGS_TO]->(t:Topic {name: "Two Pointers"}),
      (p)-[:USES_PATTERN]->(pat:Pattern {name: "Opposite Direction"})
RETURN p.title, p.difficulty

// Find common mistakes for problem
MATCH (prob:Problem {leetcode_id: 11})-[:HAS_MISTAKE]->(m:Mistake)
RETURN m.name, m.severity, m.how_to_avoid

// Find next problem to practice
MATCH (u:User {user_id: 123})-[progress:MASTERED]->(t:Topic),
      (t)-[:NEXT_RECOMMENDED]->(next_topic:Topic),
      (next_topic)<-[:BELONGS_TO]-(p:Problem)
WHERE progress.mastery_level > 0.7 
  AND NOT EXISTS((u)-[:ATTEMPTED]->(p))
RETURN p.title, next_topic.name, p.difficulty
ORDER BY p.difficulty
LIMIT 5
```

---

## 4. EMBEDDINGS STRATEGY {#embeddings-strategy}

### 4.1 Collections Configuration

```json
{
  "embedding_collections": [
    {
      "collection_name": "topic_descriptions",
      "content_type": "topic",
      "fields_to_embed": ["topic_name", "description", "key_concepts"],
      "embedding_model": "text-embedding-3-small",
      "vector_dimension": 1536,
      "metadata": {
        "topic_id": "integer",
        "topic_code": "string",
        "difficulty_range": "string"
      }
    },
    {
      "collection_name": "problem_statements",
      "content_type": "problem",
      "fields_to_embed": ["problem_title", "description", "constraints"],
      "embedding_model": "text-embedding-3-small",
      "vector_dimension": 1536,
      "metadata": {
        "problem_id": "integer",
        "difficulty": "integer",
        "topics": "array"
      }
    },
    {
      "collection_name": "pattern_templates",
      "content_type": "pattern",
      "fields_to_embed": ["pattern_name", "description", "when_to_use", "code_template"],
      "embedding_model": "text-embedding-3-small",
      "vector_dimension": 1536
    },
    {
      "collection_name": "question_bank",
      "content_type": "question",
      "fields_to_embed": ["question_text", "explanation"],
      "embedding_model": "all-MiniLM-L6-v2",
      "vector_dimension": 384
    },
    {
      "collection_name": "common_mistakes",
      "content_type": "mistake",
      "fields_to_embed": ["mistake_name", "description", "how_to_avoid"],
      "embedding_model": "all-MiniLM-L6-v2",
      "vector_dimension": 384
    }
  ]
}
```

### 4.2 Similarity Search Strategies

```python
# Find similar problems
def find_similar_problems(problem_description, top_k=5):
    results = chroma_db.query(
        collection="problem_statements",
        query_texts=[problem_description],
        n_results=top_k
    )
    return results

# Find relevant hints
def find_relevant_hints(problem_context, user_approach):
    query = f"Problem: {problem_context}\nUser approach: {user_approach}"
    results = chroma_db.query(
        collection="solution_explanations",
        query_texts=[query],
        n_results=3
    )
    return results

# Identify mistakes from code
def identify_mistake(user_buggy_code, problem_id):
    results = chroma_db.query(
        collection="common_mistakes",
        query_texts=[f"Buggy code: {user_buggy_code}"],
        where={"related_problem_id": problem_id},
        n_results=5
    )
    return results
```

---

## 5. RAG PIPELINE CONFIGURATION {#rag-pipeline}

### 5.1 Context Assembly Rules

```json
{
  "rag_context_types": {
    "hint_generation": {
      "description": "Generate personalized hints",
      "context_components": [
        {
          "source": "problem_statements",
          "query": "Current problem description",
          "top_k": 1,
          "weight": 0.3
        },
        {
          "source": "solution_explanations",
          "query": "User's approach + problem",
          "top_k": 2,
          "weight": 0.4
        },
        {
          "source": "common_mistakes",
          "query": "User's buggy code",
          "top_k": 3,
          "weight": 0.2
        },
        {
          "source": "user_progress_db",
          "query": "Past performance",
          "metadata": true,
          "weight": 0.1
        }
      ],
      "max_context_length": 3000
    },
    
    "weakness_analysis": {
      "description": "Analyze weak areas",
      "context_components": [
        {
          "source": "user_attempts_db",
          "query": "Failed attempts (30 days)",
          "weight": 0.4
        },
        {
          "source": "topic_descriptions",
          "query": "Topics with low mastery",
          "top_k": 5,
          "weight": 0.3
        },
        {
          "source": "question_bank",
          "query": "Frequently wrong questions",
          "top_k": 10,
          "weight": 0.3
        }
      ],
      "max_context_length": 4000
    }
  }
}
```

### 5.2 Prompt Templates

```json
{
  "prompt_templates": {
    "hint_level_1_socratic": {
      "name": "Level 1 Hint - Socratic",
      "template": "You are a patient DSA tutor. User solving: {problem_title}\n\nContext:\n{problem_description}\n\nUser's approach:\n{user_approach}\n\nProvide a Level 1 Socratic question that guides without giving away answer.\n\nExamples:\n- 'What happens when capacity is 0?'\n- 'When pointers meet, what does that tell you?'\n\nHint:",
      "max_tokens": 150
    },
    
    "hint_level_2_directional": {
      "name": "Level 2 Hint - Directional",
      "template": "You are a patient DSA tutor. User solving: {problem_title}\n\nContext:\n{problem_description}\n\nUser's approach:\n{user_approach}\n\nCommon mistakes:\n{common_mistakes}\n\nProvide Level 2 hint pointing right direction without full solution.\n\nExamples:\n- 'Consider using hash set to track seen elements'\n- 'Think about what previous cell tells you'\n\nHint:",
      "max_tokens": 200
    },
    
    "hint_level_3_concrete": {
      "name": "Level 3 Hint - Concrete Steps",
      "template": "User solving: {problem_title}\n\nSolution approach:\n{solution_outline}\n\nProvide Level 3 hint with concrete steps but NO code.\n\nFormat:\n1. Initialize X\n2. For each element, do Y\n3. Check condition Z\n4. Return result\n\nHint:",
      "max_tokens": 300
    },
    
    "weakness_analysis": {
      "name": "Weakness Analysis",
      "template": "Analyze user's performance.\n\nFailed attempts:\n{failed_attempts}\n\nWeak topics:\n{weak_topics}\n\nMistake patterns:\n{mistake_patterns}\n\nProvide:\n1. Top 3 specific weaknesses\n2. Root cause for each\n3. Actionable recommendations\n4. Estimated improvement time\n\nAnalysis:",
      "max_tokens": 400
    }
  }
}
```

---

## 6. N8N WORKFLOW INTEGRATION {#n8n-integration}

### 6.1 Workflow Architecture

```
Workflow 1: Database Seeding Pipeline
├── Read Master Reference
├── Parse Topic Taxonomy
├── Insert into PostgreSQL (topics, patterns, problems)
├── Create relationships
└── Trigger Workflow 2

Workflow 2: Graph Construction Pipeline
├── Read from PostgreSQL
├── Create AGE nodes (Topics, Patterns, Problems)
├── Create AGE edges (relationships)
├── Validate graph integrity
└── Trigger Workflow 3

Workflow 3: Embeddings Generation Pipeline
├── Read from PostgreSQL
├── For each collection:
│   ├── Construct embedding text
│   ├── Call embedding API
│   └── Store in ChromaDB
└── Trigger Workflow 4

Workflow 4: Question Generation Pipeline
├── Read problems from PostgreSQL
├── For each problem:
│   ├── Generate questions using Claude
│   ├── Store in PostgreSQL
│   └── Create relationships
└── Log completion

Workflow 5: RAG Hint Generation (Real-time)
├── Receive user request
├── Gather context (PostgreSQL + ChromaDB + AGE)
├── Assemble prompt
├── Call Claude API
└── Return hint
```

### 6.2 Example n8n Node Configurations

```json
{
  "workflow_name": "Embeddings Generation",
  "nodes": [
    {
      "node_type": "PostgreSQL",
      "name": "Fetch Topics",
      "parameters": {
        "operation": "select",
        "query": "SELECT topic_id, topic_name, description, key_concepts FROM topics"
      }
    },
    {
      "node_type": "Code",
      "name": "Construct Embedding Text",
      "parameters": {
        "language": "javascript",
        "code": "const text = `Topic: ${$json.topic_name}\\nDescription: ${$json.description}\\nKey Concepts: ${$json.key_concepts.join(', ')}`;\nreturn {embedding_text: text, metadata: {topic_id: $json.topic_id}};"
      }
    },
    {
      "node_type": "HTTP Request",
      "name": "Generate Embedding",
      "parameters": {
        "method": "POST",
        "url": "https://api.openai.com/v1/embeddings",
        "body": {
          "model": "text-embedding-3-small",
          "input": "={{$json.embedding_text}}"
        }
      }
    },
    {
      "node_type": "HTTP Request",
      "name": "Store in ChromaDB",
      "parameters": {
        "method": "POST",
        "url": "http://chromadb:8000/api/v1/collections/topic_descriptions/add",
        "body": {
          "ids": ["topic_{{$json.metadata.topic_id}}"],
          "embeddings": ["={{$json.data[0].embedding}}"],
          "metadatas": ["={{$json.metadata}}"]
        }
      }
    }
  ]
}
```

---

## 7. TOPIC DETAILS TEMPLATE {#topic-details}

### Standard Template for Each Topic

```json
{
  "topic_code": "X.Y.Z",
  "topic_name": "Topic Name",
  "parent_code": "X.Y",
  "depth_level": 3,
  "difficulty_range": "medium",
  "estimated_practice_hours": 8.0,
  
  "description": "Full description of topic...",
  
  "key_concepts": [
    "Concept 1",
    "Concept 2",
    "Concept 3"
  ],
  
  "prerequisites": [
    {
      "topic_code": "X.Y",
      "topic_name": "Prerequisite Topic",
      "strength": "strong",
      "reason": "Why this is prerequisite"
    }
  ],
  
  "patterns": [
    {
      "pattern_name": "Pattern Name",
      "time_complexity": "O(n)",
      "space_complexity": "O(1)",
      "template": "code template...",
      "when_to_use": "When to apply..."
    }
  ],
  
  "common_mistakes": [
    {
      "mistake_name": "Mistake Name",
      "description": "What the mistake is",
      "severity": "critical|major|minor",
      "example_code": "buggy code...",
      "correct_code": "fixed code...",
      "how_to_avoid": "How to avoid..."
    }
  ],
  
  "edge_cases": [
    {
      "case_name": "Edge Case Name",
      "test_input": "test input",
      "expected_behavior": "what should happen",
      "why_important": "why test this"
    }
  ],
  
  "typical_constraints": {
    "array_length": "1 <= n <= 10^5",
    "element_range": "-10^9 <= arr[i] <= 10^9",
    "time_limit": "O(n) required"
  },
  
  "related_problems": [
    {
      "problem_id": 123,
      "leetcode_id": 11,
      "title": "Problem Title",
      "difficulty": 45,
      "is_classic": true
    }
  ],
  
  "question_types_to_generate": [
    "What happens when...?",
    "Why does X work?",
    "Find the bug in this code"
  ],
  
  "variations": [
    "Variation 1",
    "Variation 2"
  ]
}
```

---

## 8. IMPLEMENTATION CHECKLIST {#implementation-checklist}

### Phase 1: Database Setup
- [ ] Create PostgreSQL database
- [ ] Run schema creation scripts
- [ ] Install Apache AGE extension
- [ ] Set up ChromaDB instance
- [ ] Create database indices

### Phase 2: Topic Data Seeding
- [ ] Parse master reference
- [ ] Extract topic hierarchy
- [ ] Insert topics into PostgreSQL
- [ ] Create prerequisites
- [ ] Validate relationships

### Phase 3: Graph Construction
- [ ] Create AGE graph database
- [ ] Insert Topic nodes
- [ ] Insert Pattern nodes
- [ ] Create PREREQUISITE edges
- [ ] Create USES_PATTERN edges
- [ ] Validate graph integrity

### Phase 4: Embeddings Generation
- [ ] Generate embeddings for topics
- [ ] Generate embeddings for patterns
- [ ] Store in ChromaDB collections
- [ ] Create metadata indices
- [ ] Test similarity searches

### Phase 5: Problem & Question Seeding
- [ ] Import LeetCode problems
- [ ] Map problems to topics
- [ ] Generate questions using Claude
- [ ] Store questions in PostgreSQL
- [ ] Generate question embeddings

### Phase 6: RAG Pipeline Setup
- [ ] Configure context assembly
- [ ] Set up prompt templates
- [ ] Test hint generation
- [ ] Configure caching
- [ ] Monitor API costs

### Phase 7: n8n Workflows
- [ ] Create Database Seeding workflow
- [ ] Create Graph Construction workflow
- [ ] Create Embeddings Generation workflow
- [ ] Create Question Generation workflow
- [ ] Create RAG Hint Generation workflow

### Phase 8: Testing & Validation
- [ ] Test end-to-end hint generation
- [ ] Validate graph queries
- [ ] Test similarity searches
- [ ] Validate question quality
- [ ] Load testing

---

**END OF MASTER REFERENCE GUIDE**

This document serves as the authoritative source for all automated workflows and should be updated whenever the topic taxonomy or system architecture changes.

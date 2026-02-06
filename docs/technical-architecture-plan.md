# Technical Architecture Plan
## LeetCode Training Platform - Database, Scoring, Assessment & LLM Integration

---

## DOCUMENT PURPOSE

This is the **technical blueprint** for building the platform locally with:
- PostgreSQL + Apache AGE (graph database)
- Vector DB + RAG for LLM integration
- Custom difficulty scoring system
- Fair assessment methodology
- Weakness detection
- Training plan builder
- Local deployment architecture

---

## TABLE OF CONTENTS

1. [Database Architecture](#database-architecture)
2. [Difficulty Scoring System](#difficulty-scoring-system)
3. [Assessment Methodology](#assessment-methodology)
4. [Weakness Detection Algorithm](#weakness-detection-algorithm)
5. [Training Plan Builder](#training-plan-builder)
6. [Search & Filter System](#search-filter-system)
7. [LLM Integration Architecture](#llm-integration-architecture)
8. [Question Generation Pipeline](#question-generation-pipeline)
9. [Local Deployment Setup](#local-deployment-setup)
10. [Data Models & Schemas](#data-models-schemas)

---

## 1. DATABASE ARCHITECTURE

### **Hybrid Approach: Relational + Graph + Vector**

**Why This Architecture?**
- **Relational (PostgreSQL)**: Store structured data (users, questions, answers)
- **Graph (Apache AGE)**: Model relationships (similar problems, prerequisites, follow-ups)
- **Vector (pgvector)**: Semantic search, LLM integration, similarity matching

---

### **1.1 PostgreSQL (Relational Layer)**

**Core Tables:**

```
USERS
├── user_id (PK)
├── username
├── email
├── created_at
├── current_skill_level
└── preferences (JSONB)

QUESTIONS
├── question_id (PK)
├── question_type (enum: complexity_analysis, ds_selection, pattern_recognition, etc.)
├── question_subtype
├── difficulty_score (our custom unit - see section 2)
├── leetcode_number (nullable - if from LeetCode)
├── problem_statement (text)
├── constraints (JSONB)
├── test_cases (JSONB)
├── created_at
├── updated_at
└── metadata (JSONB)

QUESTION_OPTIONS
├── option_id (PK)
├── question_id (FK)
├── option_text
├── is_correct
├── explanation
└── distractor_category (why this wrong answer is plausible)

QUESTION_TAGS
├── tag_id (PK)
└── tag_name (e.g., "binary_search", "dynamic_programming", "two_pointers")

QUESTION_TAG_JUNCTION
├── question_id (FK)
├── tag_id (FK)
└── PRIMARY KEY (question_id, tag_id)

TOPICS
├── topic_id (PK)
├── topic_name
├── description
└── parent_topic_id (FK - for hierarchy)

QUESTION_TOPICS
├── question_id (FK)
├── topic_id (FK)
└── relevance_score (0.0 - 1.0)

USER_ATTEMPTS
├── attempt_id (PK)
├── user_id (FK)
├── question_id (FK)
├── selected_answer (for MC) or answer_text (for open-ended)
├── is_correct
├── time_taken_seconds
├── attempt_number (1st try, 2nd try, etc.)
├── timestamp
├── confidence_level (1-5)
└── approach_description (text - how they thought about it)

USER_PROGRESS
├── progress_id (PK)
├── user_id (FK)
├── question_id (FK)
├── mastery_level (0.0 - 1.0)
├── last_attempt_date
├── next_review_date (spaced repetition)
├── attempts_count
└── time_trend (improving or degrading)

TRAINING_PLANS
├── plan_id (PK)
├── user_id (FK)
├── plan_name
├── created_at
├── target_date
├── focus_areas (JSONB - which topics/weaknesses to address)
└── is_active

TRAINING_PLAN_ITEMS
├── item_id (PK)
├── plan_id (FK)
├── question_id (FK)
├── sequence_order
├── is_completed
├── scheduled_date
└── notes

ASSESSMENTS
├── assessment_id (PK)
├── user_id (FK)
├── assessment_type (initial, weekly, final_exam)
├── started_at
├── completed_at
├── overall_score
└── detailed_results (JSONB)

ASSESSMENT_RESPONSES
├── response_id (PK)
├── assessment_id (FK)
├── question_id (FK)
├── user_answer
├── is_correct
├── time_taken
└── llm_feedback (JSONB - from LLM analysis)

WEAKNESSES
├── weakness_id (PK)
├── user_id (FK)
├── skill_category
├── detected_at
├── severity (0.0 - 1.0)
├── evidence (JSONB - which questions exposed this)
└── improvement_recommendations (text)

CODE_TEMPLATES
├── template_id (PK)
├── template_name
├── pattern_category
├── code_snippet (text)
├── language (C++, Python, etc.)
└── usage_notes

USER_TEMPLATE_MASTERY
├── user_id (FK)
├── template_id (FK)
├── mastery_level (0.0 - 1.0)
├── last_practiced
└── PRIMARY KEY (user_id, template_id)
```

**JSONB Fields Explained:**

```jsonb
// QUESTIONS.metadata
{
  "similar_problems": [123, 456, 789],  // LeetCode numbers
  "follow_ups": [234, 567],
  "prerequisites": [12, 45],
  "company_tags": ["Google", "Meta", "Amazon"],
  "frequency": "high",
  "acceptance_rate": 0.42,
  "solution_approaches": ["hash_map", "two_pointers", "binary_search"]
}

// QUESTIONS.constraints
{
  "time_limit_seconds": 45,
  "space_complexity_target": "O(n)",
  "input_constraints": {
    "n": {"min": 1, "max": 100000},
    "array_values": {"min": -10000, "max": 10000}
  }
}

// USER_ATTEMPTS.approach_description
{
  "pattern_identified": "two_pointers",
  "complexity_analysis": "O(n) time, O(1) space",
  "edge_cases_considered": ["empty_array", "single_element"],
  "mistakes_made": ["forgot_to_handle_duplicates"]
}

// ASSESSMENTS.detailed_results
{
  "category_scores": {
    "complexity_analysis": 0.85,
    "pattern_recognition": 0.72,
    "implementation": 0.90
  },
  "strengths": ["fast_implementation", "correct_complexity"],
  "weaknesses": ["edge_case_handling", "pattern_recognition"],
  "time_distribution": {
    "reading": 120,
    "thinking": 300,
    "coding": 600
  }
}
```

---

### **1.2 Apache AGE (Graph Layer)**

**Why Graph Database?**
- Questions have complex relationships (similar, prerequisite, follow-up)
- Topics have hierarchical and cross-cutting relationships
- Skills have dependency graphs
- Path finding (what to learn next)

**Graph Schema:**

```cypher
// Node Types

(:Question {
  question_id: int,
  difficulty_score: float,
  type: string
})

(:Topic {
  topic_id: int,
  name: string,
  level: string  // "fundamental", "intermediate", "advanced"
})

(:Pattern {
  pattern_id: int,
  name: string,
  category: string
})

(:Skill {
  skill_id: int,
  name: string,
  description: string
})

(:User {
  user_id: int,
  username: string
})

// Relationship Types

(:Question)-[:SIMILAR_TO {similarity_score: float}]->(:Question)
(:Question)-[:FOLLOW_UP_OF]->(:Question)
(:Question)-[:REQUIRES]->(:Question)  // Prerequisites
(:Question)-[:BELONGS_TO]->(:Topic)
(:Question)-[:USES_PATTERN]->(:Pattern)
(:Question)-[:TESTS_SKILL]->(:Skill)

(:Topic)-[:SUBTOPIC_OF]->(:Topic)
(:Topic)-[:RELATED_TO {strength: float}]->(:Topic)
(:Topic)-[:REQUIRES]->(:Topic)  // Topic prerequisites

(:Pattern)-[:SIMILAR_TO]->(:Pattern)
(:Pattern)-[:BUILDS_ON]->(:Pattern)

(:Skill)-[:DEPENDS_ON]->(:Skill)
(:Skill)-[:ENABLES]->(:Skill)

(:User)-[:MASTERED {level: float}]->(:Skill)
(:User)-[:WEAK_IN {severity: float}]->(:Skill)
(:User)-[:ATTEMPTED {count: int, avg_score: float}]->(:Question)
```

**Key Graph Queries:**

```cypher
// Find similar problems to a given problem
MATCH (q1:Question {question_id: $qid})-[r:SIMILAR_TO]->(q2:Question)
WHERE r.similarity_score > 0.7
RETURN q2
ORDER BY r.similarity_score DESC

// Find learning path (what to study next)
MATCH path = (start:Topic {topic_id: $current_topic})
             -[:REQUIRES*..3]->
             (end:Topic)
WHERE NOT (:User {user_id: $user_id})-[:MASTERED]->(end)
RETURN path
ORDER BY length(path)

// Find prerequisite chain
MATCH path = (q1:Question {question_id: $qid})
             -[:REQUIRES*]->
             (q2:Question)
WHERE NOT (:User {user_id: $user_id})-[:ATTEMPTED]->(q2)
RETURN path

// Recommend next question based on weaknesses
MATCH (u:User {user_id: $user_id})-[w:WEAK_IN]->(s:Skill)
MATCH (q:Question)-[:TESTS_SKILL]->(s)
WHERE NOT (u)-[:ATTEMPTED]->(q)
RETURN q, w.severity
ORDER BY w.severity DESC, q.difficulty_score ASC
LIMIT 10

// Find follow-up questions
MATCH (q1:Question {question_id: $qid})-[:FOLLOW_UP_OF*1..2]->(q2:Question)
RETURN q2

// Cluster similar problems
MATCH (q:Question)-[s:SIMILAR_TO]-(other:Question)
WHERE s.similarity_score > 0.8
WITH q, collect(other) as similar_questions
RETURN q, similar_questions
```

---

### **1.3 Vector Database (pgvector Extension)**

**Why Vector Embeddings?**
- Semantic similarity (find conceptually similar problems)
- LLM-powered search
- Question generation (find similar patterns)
- Automatic tagging

**Vector Schema:**

```sql
CREATE EXTENSION vector;

CREATE TABLE question_embeddings (
  question_id INT PRIMARY KEY REFERENCES questions(question_id),
  problem_embedding vector(1536),  -- OpenAI ada-002 or similar
  solution_embedding vector(1536),
  tags_embedding vector(1536),
  combined_embedding vector(1536),
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX ON question_embeddings 
USING ivfflat (combined_embedding vector_cosine_ops)
WITH (lists = 100);

CREATE TABLE user_attempt_embeddings (
  attempt_id INT PRIMARY KEY REFERENCES user_attempts(attempt_id),
  approach_embedding vector(1536),  -- Embed user's explanation
  created_at TIMESTAMP DEFAULT NOW()
);
```

**Vector Operations:**

```sql
-- Find semantically similar problems
SELECT q.question_id, q.problem_statement,
       1 - (qe.combined_embedding <=> $query_embedding) as similarity
FROM questions q
JOIN question_embeddings qe ON q.question_id = qe.question_id
WHERE 1 - (qe.combined_embedding <=> $query_embedding) > 0.8
ORDER BY similarity DESC
LIMIT 10;

-- Find problems matching user's approach
SELECT q.question_id,
       1 - (uae.approach_embedding <=> qe.solution_embedding) as match_score
FROM user_attempt_embeddings uae
JOIN question_embeddings qe ON true
JOIN questions q ON q.question_id = qe.question_id
WHERE uae.attempt_id = $attempt_id
ORDER BY match_score DESC;
```

---

## 2. DIFFICULTY SCORING SYSTEM

### **The "Magic Unit" - Skill Score (SS)**

**Philosophy:**
- Not just "easy/medium/hard"
- Multi-dimensional difficulty
- Personalized to user's skill level
- Continuously calibrated

**SS Range: 0.0 - 100.0**

---

### **2.1 Difficulty Dimensions**

**Each question has 7 dimension scores (0-10 each):**

```
1. Conceptual Complexity (CC)
   - How many concepts must be understood?
   - Example: Two Sum = 2 (arrays, hash maps)
   - Example: Median of Two Arrays = 5 (binary search, arrays, median, partitioning, edge cases)

2. Algorithm Complexity (AC)
   - How complex is the algorithm itself?
   - Nested loops = 3
   - Recursion with branching = 6
   - Dynamic programming with optimization = 9

3. Implementation Difficulty (ID)
   - How hard to translate idea to code?
   - Simple loop = 2
   - Pointer manipulation = 5
   - Complex data structure = 8

4. Edge Case Density (ED)
   - How many edge cases exist?
   - 1-2 edge cases = 2
   - 3-5 edge cases = 5
   - 8+ edge cases = 9

5. Pattern Recognition Difficulty (PRD)
   - How hard to recognize the pattern?
   - Obvious pattern (sorted array → binary search) = 2
   - Hidden pattern (requires insight) = 7
   - Novel pattern (never seen before) = 10

6. Time Pressure (TP)
   - Expected time to solve
   - < 10 min = 2
   - 10-20 min = 5
   - 30-45 min = 8
   - > 45 min = 10

7. Optimization Requirement (OR)
   - How much optimization needed?
   - Brute force acceptable = 2
   - Need O(n log n) instead of O(n²) = 5
   - Need O(n) with complex optimization = 8
   - Need O(log n) or better = 10
```

**Base Skill Score Formula:**

```
SS_base = (CC * 0.20) + (AC * 0.20) + (ID * 0.15) + (ED * 0.10) + 
          (PRD * 0.20) + (TP * 0.10) + (OR * 0.05)

Result: 0.0 - 10.0

Then scale to 0-100:
SS = SS_base * 10
```

**Why These Weights?**
- CC & AC & PRD = Most important (60%)
- ID = Important but mechanical (15%)
- ED & TP = Secondary (20%)
- OR = Bonus factor (5%)

---

### **2.2 Personalized Difficulty Adjustment**

**User's Current Skill Level (USL): 0-100**

```
Calculated from:
- Average accuracy across all questions
- Average time vs expected time
- Mastery of different categories
- Improvement rate
```

**Effective Difficulty for User:**

```
ED(question, user) = SS * (1 + adjustment_factor)

adjustment_factor = (SS - USL) / 100

If SS = 70, USL = 50:
  adjustment_factor = +0.20 (feels 20% harder)
  ED = 70 * 1.20 = 84

If SS = 30, USL = 50:
  adjustment_factor = -0.20 (feels 20% easier)
  ED = 30 * 0.80 = 24
```

**Dynamic Calibration:**

After each attempt, recalibrate:

```sql
-- Update question difficulty based on user performance
UPDATE questions
SET difficulty_score = difficulty_score * adjustment
WHERE question_id = $qid;

-- Adjustment based on attempt data:
adjustment = 1 + (
  (actual_accuracy - expected_accuracy) * 0.1 +
  (expected_time - actual_time) / expected_time * 0.1
)

-- Clamp adjustment: 0.9 to 1.1 (max 10% change per update)
```

---

### **2.3 Difficulty Categories (For Display)**

```
Beginner:     SS 0-20    (Green)
Easy:         SS 20-35   (Light Green)
Medium-Easy:  SS 35-50   (Yellow-Green)
Medium:       SS 50-65   (Yellow)
Medium-Hard:  SS 65-80   (Orange)
Hard:         SS 80-90   (Red)
Expert:       SS 90-100  (Dark Red)
```

**Visual Representation:**
```
[●●●●●●○○○○] 60/100 - Medium
 ↓
Conceptual: ●●●●●
Algorithm:  ●●●●●●
Implement:  ●●●●
Edge Cases: ●●●
Pattern:    ●●●●●●●
```

---

## 3. ASSESSMENT METHODOLOGY

### **Goal: Detect Understanding vs Memorization**

**Problem:**
- User might have seen solution before
- Could memorize without understanding
- Need to test TRANSFER of knowledge

---

### **3.1 Anti-Memorization Techniques**

#### **Technique 1: Variant Questions**

When user solves a problem, generate variants:

```
Original: Two Sum (find pair with target sum)

Variants:
1. Two Sum II (sorted array - tests adaptation)
2. Three Sum (tests generalization)
3. Two Sum - Count all pairs (tests different goal)
4. Two Sum - Closest sum (tests optimization variant)
```

**Detection Logic:**
```python
# If user solves original quickly
original_time = 300 seconds
original_correct = True

# Then present variant immediately
variant_time = 400 seconds  # Should be similar if they understood
variant_correct = True

# Calculate understanding score
understanding_score = (
  (1.0 if variant_correct else 0.0) * 0.6 +
  (1.0 - abs(variant_time - original_time) / original_time) * 0.4
)

# understanding_score close to 1.0 = real understanding
# understanding_score close to 0.0 = memorization
```

---

#### **Technique 2: Explain Your Thinking (LLM-Evaluated)**

After solving, ask:
```
"Explain WHY your solution works"
"What pattern did you use?"
"What would happen if constraint X changed?"
"How would you optimize further?"
```

**LLM Prompt for Evaluation:**
```
Evaluate this explanation:
Problem: [problem statement]
User Solution: [code or approach]
User Explanation: [their explanation]

Rate 0-10 on:
1. Correctness of explanation
2. Depth of understanding
3. Awareness of edge cases
4. Complexity analysis accuracy
5. Pattern recognition

Output JSON: {
  "scores": {...},
  "understanding_level": "memorized" | "partial" | "deep",
  "evidence": "...",
  "follow_up_questions": [...]
}
```

---

#### **Technique 3: Constraint Modification**

Change problem constraints, ask how solution changes:

```
Original: Array size n ≤ 10^5, O(n log n) acceptable

Modified: Array size n ≤ 10^9, need O(n) solution

Question: "How would you adapt your solution?"

Expected: Discussion of why O(n log n) won't work, 
          what O(n) approach would be needed
```

**Scoring:**
```python
# User provides explanation of adaptation
adaptation_explanation = user_input()

# LLM evaluates whether they can reason about trade-offs
llm_score = evaluate_reasoning(
  problem_understanding=True/False,
  complexity_awareness=True/False,
  valid_adaptation=True/False
)
```

---

#### **Technique 4: Time Pattern Analysis**

```python
# Memorization indicators:
def detect_memorization(attempt_history):
  indicators = {
    "suspiciously_fast": time_taken < expected_time * 0.3,
    "no_mistakes": all(first_try_correct for q in similar_questions),
    "identical_approach": code_similarity > 0.95 to known_solution,
    "no_debugging_time": time_thinking < 10% of time_coding,
    "perfect_edge_cases": handled_all_edge_cases_on_first_try
  }
  
  memorization_score = sum(indicators.values()) / len(indicators)
  
  return memorization_score > 0.6  # Likely memorized
```

---

### **3.2 Understanding-Based Assessment**

**Multi-Level Testing:**

```
Level 1: Can they solve the exact problem?
  → Tests: Correct implementation

Level 2: Can they explain WHY it works?
  → Tests: Conceptual understanding

Level 3: Can they solve a variant?
  → Tests: Transfer learning

Level 4: Can they adapt to changed constraints?
  → Tests: Deep understanding

Level 5: Can they compare multiple approaches?
  → Tests: Mastery
```

**Final Understanding Score:**

```
Understanding = (
  L1_score * 0.15 +  # Basic implementation
  L2_score * 0.25 +  # Explanation
  L3_score * 0.30 +  # Variant solving
  L4_score * 0.20 +  # Adaptation
  L5_score * 0.10    # Comparison
)
```

---

### **3.3 Fair Scoring System**

**Problem: How to score fairly?**

```
User A: Solves in 10 minutes, memorized solution
User B: Solves in 30 minutes, figured it out

Who should score higher? B!
```

**Multi-Factor Scoring:**

```python
def calculate_fair_score(attempt):
  factors = {
    "correctness": 0.35,      # Did it work?
    "efficiency": 0.15,       # Time complexity correct?
    "understanding": 0.30,    # Can they explain?
    "adaptability": 0.20      # Can they modify?
  }
  
  # Correctness
  correctness_score = 1.0 if attempt.is_correct else 0.0
  
  # Efficiency (time taken vs expected)
  time_ratio = attempt.time_taken / expected_time
  if time_ratio < 0.5:  # Too fast = suspicious
    efficiency_score = 0.7
  elif time_ratio < 1.0:  # Within expected
    efficiency_score = 1.0
  elif time_ratio < 1.5:  # Bit slow
    efficiency_score = 0.8
  else:  # Too slow
    efficiency_score = 0.5
  
  # Understanding (LLM-evaluated)
  understanding_score = llm_evaluate_explanation(attempt.explanation)
  
  # Adaptability (variant performance)
  adaptability_score = solve_variant_and_compare(attempt)
  
  # Weighted sum
  total_score = (
    correctness_score * factors["correctness"] +
    efficiency_score * factors["efficiency"] +
    understanding_score * factors["understanding"] +
    adaptability_score * factors["adaptability"]
  )
  
  return total_score * 100  # Scale to 0-100
```

---

## 4. WEAKNESS DETECTION ALGORITHM

### **4.1 Weakness Categories**

```
1. Pattern Recognition Weaknesses
   - Can't identify two pointers
   - Misses binary search opportunities
   - Doesn't recognize DP problems

2. Implementation Weaknesses
   - Off-by-one errors
   - Pointer manipulation bugs
   - Edge case handling

3. Complexity Analysis Weaknesses
   - Overestimates complexity
   - Underestimates complexity
   - Doesn't recognize amortized complexity

4. Data Structure Selection Weaknesses
   - Uses array when hash map better
   - Uses BST when hash table sufficient
   - Doesn't know when to use heap

5. Algorithm Knowledge Gaps
   - Doesn't know Dijkstra
   - Doesn't understand DP
   - Can't implement DFS correctly

6. Trade-off Understanding Weaknesses
   - Always chooses time over space
   - Doesn't consider iterative vs recursive
   - Misses optimization opportunities

7. Debugging Weaknesses
   - Can't find bugs quickly
   - Doesn't test edge cases
   - Misses logic errors
```

---

### **4.2 Weakness Detection Algorithm**

**Step 1: Track Performance by Category**

```sql
CREATE MATERIALIZED VIEW user_category_performance AS
SELECT 
  u.user_id,
  q.question_type,
  q.question_subtype,
  AVG(CASE WHEN ua.is_correct THEN 1.0 ELSE 0.0 END) as accuracy,
  AVG(ua.time_taken_seconds) as avg_time,
  COUNT(*) as attempts
FROM users u
JOIN user_attempts ua ON u.user_id = ua.user_id
JOIN questions q ON ua.question_id = q.question_id
GROUP BY u.user_id, q.question_type, q.question_subtype;
```

**Step 2: Identify Below-Average Categories**

```python
def detect_weaknesses(user_id):
  # Get user's performance per category
  user_perf = get_category_performance(user_id)
  
  # Get average performance per category (all users)
  avg_perf = get_average_category_performance()
  
  weaknesses = []
  
  for category in all_categories:
    user_score = user_perf[category]['accuracy']
    avg_score = avg_perf[category]['accuracy']
    
    # Weakness if significantly below average
    if user_score < avg_score * 0.7:  # 30% below average
      severity = (avg_score - user_score) / avg_score
      
      weaknesses.append({
        'category': category,
        'severity': severity,
        'user_accuracy': user_score,
        'expected_accuracy': avg_score,
        'attempts': user_perf[category]['attempts']
      })
  
  return sorted(weaknesses, key=lambda w: w['severity'], reverse=True)
```

---

**Step 3: Pattern-Based Detection**

```python
def detect_specific_patterns(user_attempts):
  patterns = {
    "off_by_one_errors": 0,
    "null_pointer_issues": 0,
    "complexity_misunderstanding": 0,
    "wrong_data_structure": 0,
    "missed_edge_cases": 0
  }
  
  for attempt in user_attempts:
    if not attempt.is_correct:
      # Analyze the mistake with LLM
      mistake_analysis = llm_analyze_mistake(
        question=attempt.question,
        user_solution=attempt.answer,
        correct_solution=attempt.question.solution
      )
      
      # Categorize the mistake
      for error_type in mistake_analysis['error_types']:
        if error_type in patterns:
          patterns[error_type] += 1
  
  # Identify most common error patterns
  return {k: v for k, v in patterns.items() if v > 3}
```

---

**Step 4: Time-Based Analysis**

```python
def analyze_time_patterns(user_attempts):
  insights = {}
  
  # Time distribution analysis
  time_data = [a.time_taken for a in user_attempts]
  
  # Check for reading comprehension issues
  # (consistently long reading time)
  reading_times = [a.time_breakdown['reading'] for a in user_attempts]
  if avg(reading_times) > expected_reading_time * 1.5:
    insights['slow_reading'] = {
      'severity': 0.6,
      'recommendation': 'Practice reading problem statements faster'
    }
  
  # Check for planning issues
  # (goes straight to coding without thinking)
  thinking_times = [a.time_breakdown['thinking'] for a in user_attempts]
  if avg(thinking_times) < expected_thinking_time * 0.3:
    insights['insufficient_planning'] = {
      'severity': 0.8,
      'recommendation': 'Spend more time planning before coding'
    }
  
  # Check for implementation speed issues
  coding_times = [a.time_breakdown['coding'] for a in user_attempts]
  if avg(coding_times) > expected_coding_time * 2.0:
    insights['slow_implementation'] = {
      'severity': 0.7,
      'recommendation': 'Practice templates to code faster'
    }
  
  return insights
```

---

**Step 5: LLM-Powered Deep Analysis**

```python
def llm_deep_analysis(user_id, recent_attempts):
  # Prepare context for LLM
  context = {
    'user_profile': get_user_profile(user_id),
    'recent_attempts': [
      {
        'question': a.question.problem_statement,
        'user_solution': a.answer,
        'correct': a.is_correct,
        'time_taken': a.time_taken,
        'user_explanation': a.approach_description
      }
      for a in recent_attempts[-20:]  # Last 20 attempts
    ]
  }
  
  prompt = f"""
  Analyze this user's problem-solving patterns and identify weaknesses:
  
  {json.dumps(context, indent=2)}
  
  Identify:
  1. Recurring mistakes
  2. Knowledge gaps
  3. Pattern recognition issues
  4. Implementation weaknesses
  5. Specific skills to improve
  
  For each weakness, provide:
  - Description
  - Severity (0.0-1.0)
  - Evidence from attempts
  - Specific recommendations
  - Suggested practice problems
  
  Output as JSON.
  """
  
  llm_response = llm.generate(prompt)
  return parse_weaknesses(llm_response)
```

---

### **4.3 Weakness Storage & Tracking**

```sql
INSERT INTO weaknesses (
  user_id,
  skill_category,
  detected_at,
  severity,
  evidence,
  improvement_recommendations
) VALUES (
  $user_id,
  'pattern_recognition:binary_search',
  NOW(),
  0.75,
  jsonb_build_object(
    'failed_questions', ARRAY[123, 456, 789],
    'common_mistakes', ARRAY['missed_sorted_array_cue', 'used_linear_search'],
    'accuracy', 0.25
  ),
  'Practice 10 binary search problems focusing on pattern recognition'
);

-- Track improvement over time
CREATE VIEW weakness_improvement AS
SELECT 
  w.user_id,
  w.skill_category,
  w.detected_at,
  w.severity,
  COALESCE(
    AVG(CASE WHEN ua.is_correct THEN 1.0 ELSE 0.0 END) 
    FILTER (WHERE ua.timestamp > w.detected_at)
  , 0) as post_detection_accuracy
FROM weaknesses w
LEFT JOIN user_attempts ua ON w.user_id = ua.user_id
LEFT JOIN questions q ON ua.question_id = q.question_id
WHERE q.question_type = split_part(w.skill_category, ':', 1)
GROUP BY w.weakness_id, w.user_id, w.skill_category, w.detected_at, w.severity;
```

---

## 5. TRAINING PLAN BUILDER

### **5.1 Plan Structure**

**Training Plan Components:**

```
1. Goal Definition
   - Target interview date
   - Target companies
   - Current skill level
   - Target skill level
   - Focus areas

2. Curriculum
   - Topics to cover
   - Skills to master
   - Templates to memorize
   - Weaknesses to address

3. Schedule
   - Questions per day
   - Review schedule (spaced repetition)
   - Mock interviews
   - Assessments

4. Adaptive Elements
   - Adjust based on progress
   - Add questions for weak areas
   - Remove mastered topics
   - Accelerate if ahead of schedule
```

---

### **5.2 Plan Generation Algorithm**

```python
def generate_training_plan(user_id, goals):
  # Step 1: Assess current state
  current_state = {
    'skill_level': get_user_skill_level(user_id),
    'strengths': get_user_strengths(user_id),
    'weaknesses': get_user_weaknesses(user_id),
    'completed_topics': get_completed_topics(user_id),
    'time_available': goals['hours_per_week']
  }
  
  # Step 2: Define target state
  target_state = {
    'target_level': goals['target_skill_level'],
    'must_know_topics': get_required_topics(goals['target_companies']),
    'target_date': goals['interview_date']
  }
  
  # Step 3: Calculate gap
  gap = calculate_skill_gap(current_state, target_state)
  
  # Step 4: Prioritize topics
  topic_priorities = prioritize_topics(
    weaknesses=current_state['weaknesses'],
    required_topics=target_state['must_know_topics'],
    time_available=current_state['time_available']
  )
  
  # Step 5: Allocate questions
  question_allocation = allocate_questions(
    priorities=topic_priorities,
    total_time=weeks_until_target * goals['hours_per_week'],
    current_skill=current_state['skill_level']
  )
  
  # Step 6: Schedule with spaced repetition
  schedule = create_schedule(
    questions=question_allocation,
    weeks_available=weeks_until_target,
    hours_per_week=goals['hours_per_week']
  )
  
  return {
    'plan_id': create_plan_in_db(user_id, schedule),
    'summary': generate_plan_summary(schedule),
    'milestones': define_milestones(schedule),
    'flexibility': 'high'  # Can adapt based on progress
  }
```

---

### **5.3 Topic Prioritization**

```python
def prioritize_topics(weaknesses, required_topics, time_available):
  priorities = []
  
  for topic in all_topics:
    score = 0
    
    # 1. Is it a weakness? (HIGH priority)
    if topic in [w['category'] for w in weaknesses]:
      weakness_severity = next(w['severity'] for w in weaknesses if w['category'] == topic)
      score += weakness_severity * 10
    
    # 2. Is it required for target? (MEDIUM-HIGH priority)
    if topic in required_topics:
      score += 7
    
    # 3. Is it a prerequisite? (MEDIUM priority)
    dependent_topics = get_topics_depending_on(topic)
    if any(t in required_topics for t in dependent_topics):
      score += 5
    
    # 4. Frequency in real interviews (LOW-MEDIUM priority)
    frequency = get_topic_frequency(topic)
    score += frequency * 3
    
    # 5. Current mastery (INVERSE - study what you don't know)
    mastery = get_user_topic_mastery(user_id, topic)
    score += (1.0 - mastery) * 4
    
    priorities.append({
      'topic': topic,
      'priority_score': score,
      'reasoning': f"Weakness: {weakness_severity if topic in weaknesses else 0}, "
                   f"Required: {topic in required_topics}, "
                   f"Mastery: {mastery}"
    })
  
  return sorted(priorities, key=lambda x: x['priority_score'], reverse=True)
```

---

### **5.4 Question Allocation**

```python
def allocate_questions(priorities, total_time, current_skill):
  allocation = []
  
  # Time budget per topic (minutes)
  time_budget = {}
  remaining_time = total_time * 60  # Convert to minutes
  
  for priority in priorities:
    topic = priority['topic']
    
    # Allocate time based on priority
    if priority['priority_score'] > 8:  # Critical
      time_budget[topic] = remaining_time * 0.25
    elif priority['priority_score'] > 5:  # Important
      time_budget[topic] = remaining_time * 0.15
    else:  # Nice to have
      time_budget[topic] = remaining_time * 0.05
  
  # Normalize so total = remaining_time
  total_allocated = sum(time_budget.values())
  for topic in time_budget:
    time_budget[topic] = (time_budget[topic] / total_allocated) * remaining_time
  
  # Convert time to question count
  for topic, minutes in time_budget.items():
    # Get questions for this topic
    questions = get_questions_for_topic(topic)
    
    # Filter by appropriate difficulty
    appropriate_questions = [
      q for q in questions
      if abs(q.difficulty_score - current_skill) < 20  # Within range
    ]
    
    # Calculate how many questions fit in time budget
    avg_time_per_question = 25  # minutes
    num_questions = int(minutes / avg_time_per_question)
    
    # Select questions (mix of difficulties)
    selected = select_diverse_questions(
      questions=appropriate_questions,
      count=num_questions,
      difficulty_range=(current_skill - 15, current_skill + 15)
    )
    
    allocation.append({
      'topic': topic,
      'questions': selected,
      'estimated_time': len(selected) * avg_time_per_question
    })
  
  return allocation
```

---

### **5.5 Spaced Repetition Scheduling**

```python
def create_schedule(questions, weeks_available, hours_per_week):
  schedule = []
  
  # Spaced repetition intervals (days)
  review_intervals = [1, 3, 7, 14, 30]
  
  # Initial attempt schedule
  days_available = weeks_available * 7
  minutes_per_day = (hours_per_week * 60) / 7
  questions_per_day = int(minutes_per_day / 25)  # 25 min per question
  
  current_day = 0
  
  for topic_allocation in questions:
    for question in topic_allocation['questions']:
      # Schedule initial attempt
      schedule.append({
        'day': current_day,
        'question_id': question.question_id,
        'attempt_type': 'initial',
        'topic': topic_allocation['topic']
      })
      
      # Schedule reviews based on performance
      # (will be dynamically adjusted)
      for interval_days in review_intervals:
        review_day = current_day + interval_days
        if review_day < days_available:
          schedule.append({
            'day': review_day,
            'question_id': question.question_id,
            'attempt_type': 'review',
            'interval': interval_days
          })
      
      # Move to next day if quota filled
      current_day += 1 / questions_per_day
  
  return group_by_day(schedule)
```

---

### **5.6 Adaptive Replanning**

```python
def adapt_plan(plan_id, user_id):
  # Get current progress
  progress = get_plan_progress(plan_id, user_id)
  
  # Check if adjustments needed
  adjustments_needed = []
  
  # 1. Behind schedule?
  if progress['completion_rate'] < progress['expected_completion_rate']:
    adjustments_needed.append({
      'type': 'increase_intensity',
      'reason': 'behind_schedule',
      'action': 'Add 30 minutes per day OR reduce question count'
    })
  
  # 2. Ahead of schedule?
  elif progress['completion_rate'] > progress['expected_completion_rate'] * 1.2:
    adjustments_needed.append({
      'type': 'add_challenges',
      'reason': 'ahead_of_schedule',
      'action': 'Add harder questions'
    })
  
  # 3. New weaknesses detected?
  new_weaknesses = detect_new_weaknesses(user_id, since=plan_created_at)
  if new_weaknesses:
    adjustments_needed.append({
      'type': 'address_weakness',
      'reason': f"New weakness in {new_weaknesses[0]['category']}",
      'action': f"Add 10 questions on {new_weaknesses[0]['category']}"
    })
  
  # 4. Topic mastered?
  mastered_topics = detect_mastered_topics(user_id, plan_id)
  if mastered_topics:
    adjustments_needed.append({
      'type': 'remove_mastered',
      'reason': f"Mastered {mastered_topics[0]}",
      'action': f"Remove remaining {mastered_topics[0]} questions"
    })
  
  # Apply adjustments
  for adjustment in adjustments_needed:
    apply_adjustment(plan_id, adjustment)
  
  return adjustments_needed
```

---

## 6. SEARCH & FILTER SYSTEM

### **6.1 Search Requirements**

**User should be able to search by:**
1. Problem number (LeetCode number)
2. Problem title
3. Keywords in description
4. Pattern/algorithm (semantic)
5. Difficulty range
6. Topics/tags
7. Similar problems
8. Time to solve
9. Company tags

---

### **6.2 Hybrid Search Architecture**

**Combine 3 search methods:**
1. **Full-text search** (PostgreSQL)
2. **Vector similarity** (pgvector)
3. **Graph traversal** (Apache AGE)

```python
def search_questions(query, filters):
  results = {
    'fulltext': [],
    'semantic': [],
    'graph': []
  }
  
  # 1. Full-text search
  if query.strip():
    results['fulltext'] = fulltext_search(query, filters)
  
  # 2. Semantic search (vector similarity)
  if query.strip():
    query_embedding = embed_text(query)
    results['semantic'] = vector_search(query_embedding, filters)
  
  # 3. Graph search (if specific problem given)
  if filters.get('similar_to'):
    results['graph'] = graph_search(filters['similar_to'], filters)
  
  # Merge and rank results
  return merge_and_rank(results, query, filters)
```

---

### **6.3 Full-Text Search**

```sql
-- Add full-text search capabilities
ALTER TABLE questions ADD COLUMN search_vector tsvector
  GENERATED ALWAYS AS (
    setweight(to_tsvector('english', COALESCE(problem_statement, '')), 'A') ||
    setweight(to_tsvector('english', COALESCE(array_to_string(ARRAY(
      SELECT tag_name FROM tags t 
      JOIN question_tag_junction qtj ON t.tag_id = qtj.tag_id 
      WHERE qtj.question_id = questions.question_id
    ), ' ')), 'B'
  ) STORED;

CREATE INDEX questions_search_idx ON questions USING GIN (search_vector);

-- Search query
SELECT 
  q.question_id,
  q.problem_statement,
  ts_rank(q.search_vector, query) AS rank
FROM 
  questions q,
  plainto_tsquery('english', $search_query) query
WHERE 
  q.search_vector @@ query
  AND q.difficulty_score BETWEEN $min_difficulty AND $max_difficulty
ORDER BY rank DESC
LIMIT 50;
```

---

### **6.4 Vector Semantic Search**

```sql
-- Find semantically similar problems
WITH query_embedding AS (
  SELECT $query_vector::vector(1536) as vec
)
SELECT 
  q.question_id,
  q.problem_statement,
  1 - (qe.combined_embedding <=> qv.vec) as similarity_score
FROM 
  questions q
  JOIN question_embeddings qe ON q.question_id = qe.question_id,
  query_embedding qv
WHERE 
  1 - (qe.combined_embedding <=> qv.vec) > 0.7
  AND q.difficulty_score BETWEEN $min_difficulty AND $max_difficulty
ORDER BY similarity_score DESC
LIMIT 50;
```

---

### **6.5 Graph-Based Search**

```cypher
-- Find similar problems via graph
MATCH (q1:Question {question_id: $question_id})
MATCH (q1)-[r:SIMILAR_TO]->(q2:Question)
WHERE r.similarity_score > $min_similarity
AND q2.difficulty_score >= $min_difficulty
AND q2.difficulty_score <= $max_difficulty
RETURN q2
ORDER BY r.similarity_score DESC
LIMIT 50

-- Find problems by pattern path
MATCH (q:Question)-[:USES_PATTERN]->(p:Pattern {name: $pattern_name})
WHERE q.difficulty_score >= $min_difficulty
AND q.difficulty_score <= $max_difficulty
RETURN q
ORDER BY q.difficulty_score
```

---

### **6.6 Filter System**

```python
class QuestionFilters:
  def __init__(self):
    self.difficulty_range = (0, 100)
    self.question_types = []  # ['complexity_analysis', 'pattern_recognition']
    self.tags = []  # ['binary_search', 'dynamic_programming']
    self.companies = []  # ['Google', 'Meta']
    self.time_range = (0, float('inf'))  # Expected time to solve
    self.status = 'all'  # 'all', 'unsolved', 'solved', 'mastered'
    self.leetcode_numbers = []  # Specific problem numbers
    self.similar_to = None  # Question ID for similarity search
    self.pattern = None  # Specific pattern name
    self.topic = None  # Specific topic
    self.exclude_attempted = False
    self.only_weak_areas = False
    self.custom_filter = None  # SQL WHERE clause

def apply_filters(base_query, filters):
  query = base_query
  
  # Difficulty
  query = query.filter(
    Question.difficulty_score.between(
      filters.difficulty_range[0],
      filters.difficulty_range[1]
    )
  )
  
  # Question types
  if filters.question_types:
    query = query.filter(
      Question.question_type.in_(filters.question_types)
    )
  
  # Tags
  if filters.tags:
    query = query.join(QuestionTagJunction).join(Tag).filter(
      Tag.tag_name.in_(filters.tags)
    )
  
  # Status
  if filters.status == 'unsolved':
    query = query.filter(
      ~Question.question_id.in_(
        select(UserAttempt.question_id).where(
          UserAttempt.user_id == current_user_id
        )
      )
    )
  elif filters.status == 'mastered':
    query = query.join(UserProgress).filter(
      UserProgress.mastery_level > 0.8
    )
  
  # Exclude attempted
  if filters.exclude_attempted:
    query = query.filter(
      ~Question.question_id.in_(
        select(UserAttempt.question_id).where(
          UserAttempt.user_id == current_user_id
        )
      )
    )
  
  # Only weak areas
  if filters.only_weak_areas:
    weak_categories = get_user_weaknesses(current_user_id)
    query = query.filter(
      Question.question_type.in_([w['category'] for w in weak_categories])
    )
  
  return query
```

---

## 7. LLM INTEGRATION ARCHITECTURE

### **7.1 LLM Use Cases**

**Where LLMs are used:**
1. **Question Generation**: Create similar problems
2. **Explanation Evaluation**: Judge user's explanations
3. **Hint Generation**: Provide contextual hints
4. **Weakness Analysis**: Deep analysis of mistakes
5. **Solution Comparison**: Compare user solution to optimal
6. **Follow-up Questions**: Generate adaptive follow-ups
7. **Code Review**: Identify bugs and suggest fixes
8. **Pattern Recognition**: Help identify which pattern to use

---

### **7.2 RAG (Retrieval Augmented Generation) Setup**

**Components:**

```
1. Vector Store (pgvector)
   - All questions + solutions + explanations
   - Code templates
   - Pattern descriptions
   - User's past attempts

2. Embedding Model
   - OpenAI ada-002 OR
   - Sentence-Transformers (local)

3. LLM
   - OpenAI GPT-4 OR
   - Claude OR
   - Local: Llama-2-70B, Mixtral

4. Prompt Templates
   - Structured prompts for each use case
```

**RAG Architecture:**

```
User Query
    ↓
[Embed Query]
    ↓
[Vector Search] → Retrieve relevant context
    ↓            (similar problems, templates, past attempts)
[Construct Prompt] ← Add context
    ↓
[LLM Generate]
    ↓
[Parse & Store Response]
    ↓
Display to User
```

---

### **7.3 LLM-Powered Features**

#### **Feature 1: Explanation Evaluator**

```python
def evaluate_explanation(question, user_solution, user_explanation):
  # Retrieve context
  similar_explanations = vector_search(
    embed(user_explanation),
    filter={'question_id': question.question_id}
  )
  
  # Construct prompt
  prompt = f"""
  You are an expert technical interviewer evaluating a candidate's explanation.
  
  Problem:
  {question.problem_statement}
  
  Candidate's Solution:
  {user_solution}
  
  Candidate's Explanation:
  {user_explanation}
  
  Reference Explanations (from top performers):
  {format_references(similar_explanations)}
  
  Evaluate the explanation on:
  1. Correctness (does it accurately describe the solution?)
  2. Completeness (does it cover time/space complexity, edge cases?)
  3. Clarity (is it easy to follow?)
  4. Depth (does it show understanding or just surface knowledge?)
  
  Also identify:
  - What they got right
  - What they missed
  - Whether they truly understand or are parroting
  
  Output as JSON:
  {{
    "correctness_score": 0-10,
    "completeness_score": 0-10,
    "clarity_score": 0-10,
    "depth_score": 0-10,
    "understanding_level": "memorized" | "partial" | "solid" | "deep",
    "strengths": ["...", "..."],
    "gaps": ["...", "..."],
    "follow_up_questions": ["...", "..."],
    "overall_assessment": "..."
  }}
  """
  
  response = llm.generate(prompt, temperature=0.3)
  return json.loads(response)
```

---

#### **Feature 2: Adaptive Hint Generator**

```python
def generate_hint(question, user_context, difficulty_level):
  # difficulty_level: 'gentle', 'moderate', 'direct'
  
  # Retrieve similar problems user solved
  past_successes = vector_search(
    embed(question.problem_statement),
    filter={
      'user_id': user_context['user_id'],
      'is_correct': True
    }
  )
  
  prompt = f"""
  Generate a {difficulty_level} hint for this problem.
  
  Problem:
  {question.problem_statement}
  
  User's Background:
  - Has solved similar problems: {[p.title for p in past_successes]}
  - Current skill level: {user_context['skill_level']}
  - Time spent so far: {user_context['time_spent']} minutes
  - Patterns they know: {user_context['known_patterns']}
  
  Hint Guidelines:
  - {difficulty_level} hint means:
    * gentle: Point to general approach without giving away pattern
    * moderate: Mention the pattern or data structure
    * direct: Outline the algorithm steps
  
  - Build on what they already know (reference their past solutions)
  - Don't give the complete solution
  - Ask a guiding question that leads them to the insight
  
  Output as JSON:
  {{
    "hint_text": "...",
    "hint_type": "approach" | "pattern" | "data_structure" | "edge_case",
    "leads_to": "what realization this hint should trigger"
  }}
  """
  
  response = llm.generate(prompt, temperature=0.7)
  return json.loads(response)
```

---

#### **Feature 3: Similar Question Generator**

```python
def generate_similar_question(original_question, variation_type):
  # variation_type: 'constraint_change', 'inverse', 'generalization', 'specialization'
  
  # Retrieve the pattern and template
  pattern = get_pattern(original_question)
  template = get_template(pattern)
  
  prompt = f"""
  Generate a {variation_type} of this problem.
  
  Original Problem:
  {original_question.problem_statement}
  
  Pattern Used: {pattern.name}
  
  Variation Type: {variation_type}
  
  Guidelines:
  - constraint_change: Modify constraints (array size, value ranges, operations allowed)
  - inverse: Swap input/output (e.g., "find max" → "find min")
  - generalization: Make it harder (e.g., "two sum" → "three sum")
  - specialization: Add constraints (e.g., "sorted array", "positive integers only")
  
  Requirements:
  - Use the same pattern/algorithm
  - Keep difficulty within ±10 of original ({original_question.difficulty_score})
  - Ensure it's testable
  - Provide test cases
  
  Output as JSON:
  {{
    "problem_statement": "...",
    "constraints": {{}},
    "test_cases": [...],
    "expected_pattern": "...",
    "expected_complexity": {{"time": "...", "space": "..."}},
    "how_it_relates": "explanation of how this tests same understanding as original"
  }}
  """
  
  response = llm.generate(prompt, temperature=0.8)
  
  # Store in database
  new_question = parse_and_store_question(response, original_question)
  
  # Create graph relationship
  create_edge(original_question.id, new_question.id, 'SIMILAR_TO', 
              {'generated': True, 'variation_type': variation_type})
  
  return new_question
```

---

### **7.4 Local LLM Setup (For Privacy & Speed)**

**Option 1: Ollama (Easiest)**

```bash
# Install Ollama
curl -fsSL https://ollama.ai/install.sh | sh

# Pull models
ollama pull llama2:70b       # For generation
ollama pull mistral:latest   # Faster, still good

# Run
ollama serve
```

**Python Integration:**

```python
import ollama

def llm_generate(prompt, model="llama2:70b"):
  response = ollama.generate(
    model=model,
    prompt=prompt,
    options={
      "temperature": 0.7,
      "top_p": 0.9,
      "num_predict": 1024
    }
  )
  return response['response']
```

---

**Option 2: llama.cpp (More Control)**

```bash
# Download model
wget https://huggingface.co/TheBloke/Llama-2-70B-GGUF/resolve/main/llama-2-70b.Q4_K_M.gguf

# Run server
./server -m llama-2-70b.Q4_K_M.gguf --host 127.0.0.1 --port 8080
```

**Python Integration:**

```python
import requests

def llm_generate(prompt):
  response = requests.post("http://127.0.0.1:8080/completion", json={
    "prompt": prompt,
    "temperature": 0.7,
    "n_predict": 1024
  })
  return response.json()['content']
```

---

### **7.5 Embedding Model (Local)**

**Using sentence-transformers:**

```python
from sentence_transformers import SentenceTransformer

# Load model once at startup
embedding_model = SentenceTransformer('all-MiniLM-L6-v2')  # Fast, 384 dim
# OR
embedding_model = SentenceTransformer('all-mpnet-base-v2')  # Better quality, 768 dim

def embed_text(text):
  return embedding_model.encode(text, normalize_embeddings=True)

# Store in database
def store_embedding(question_id, text):
  embedding = embed_text(text)
  
  conn.execute("""
    INSERT INTO question_embeddings (question_id, combined_embedding)
    VALUES (%s, %s)
    ON CONFLICT (question_id) DO UPDATE
    SET combined_embedding = EXCLUDED.combined_embedding
  """, (question_id, embedding.tolist()))
```

---

## 8. LOCAL DEPLOYMENT SETUP

### **8.1 System Architecture**

```
┌─────────────────────────────────────────────────────────────┐
│                     Frontend (React/Vue)                     │
│  - Question display                                          │
│  - Code editor                                               │
│  - Progress dashboard                                        │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           │ REST API / GraphQL
                           │
┌──────────────────────────▼──────────────────────────────────┐
│                   Backend (FastAPI/Flask)                    │
│  - API endpoints                                             │
│  - Authentication                                            │
│  - Business logic                                            │
└───────┬──────────┬──────────┬────────────┬──────────────────┘
        │          │          │            │
        │          │          │            │
┌───────▼────┐ ┌──▼─────┐ ┌──▼──────┐ ┌──▼────────────────┐
│ PostgreSQL │ │ Apache │ │ pgvector│ │ LLM Service       │
│ + AGE      │ │  AGE   │ │         │ │ (Ollama/llama.cpp)│
│            │ │ Graph  │ │ Vectors │ │                   │
└────────────┘ └────────┘ └─────────┘ └───────────────────┘
```

---

### **8.2 Docker Compose Setup**

**docker-compose.yml:**

```yaml
version: '3.8'

services:
  # PostgreSQL with AGE extension
  postgres:
    image: apache/age:latest
    container_name: leetcode_db
    environment:
      POSTGRES_USER: leetcode
      POSTGRES_PASSWORD: secure_password
      POSTGRES_DB: leetcode_training
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    restart: unless-stopped
  
  # Backend API
  backend:
    build: ./backend
    container_name: leetcode_backend
    environment:
      DATABASE_URL: postgresql://leetcode:secure_password@postgres:5432/leetcode_training
      LLM_API_URL: http://ollama:11434
    ports:
      - "8000:8000"
    depends_on:
      - postgres
      - ollama
    volumes:
      - ./backend:/app
    restart: unless-stopped
  
  # Ollama LLM service
  ollama:
    image: ollama/ollama:latest
    container_name: leetcode_llm
    ports:
      - "11434:11434"
    volumes:
      - ollama_data:/root/.ollama
    restart: unless-stopped
  
  # Frontend
  frontend:
    build: ./frontend
    container_name: leetcode_frontend
    ports:
      - "3000:3000"
    depends_on:
      - backend
    environment:
      REACT_APP_API_URL: http://localhost:8000
    restart: unless-stopped

volumes:
  postgres_data:
  ollama_data:
```

---

### **8.3 Database Initialization**

**init.sql:**

```sql
-- Enable extensions
CREATE EXTENSION IF NOT EXISTS vector;
LOAD 'age';
SET search_path = ag_catalog, "$user", public;

-- Create graph
SELECT create_graph('leetcode_graph');

-- Create base tables (from section 1)
-- ... (all CREATE TABLE statements)

-- Create indexes
CREATE INDEX idx_questions_difficulty ON questions(difficulty_score);
CREATE INDEX idx_questions_type ON questions(question_type);
CREATE INDEX idx_user_attempts_user ON user_attempts(user_id);
CREATE INDEX idx_user_attempts_question ON user_attempts(question_id);
CREATE INDEX idx_user_attempts_timestamp ON user_attempts(timestamp);

-- Create vector index
CREATE INDEX ON question_embeddings 
USING ivfflat (combined_embedding vector_cosine_ops)
WITH (lists = 100);

-- Create views for analytics
CREATE VIEW user_stats AS
SELECT 
  u.user_id,
  COUNT(DISTINCT ua.question_id) as questions_attempted,
  AVG(CASE WHEN ua.is_correct THEN 1.0 ELSE 0.0 END) as accuracy,
  AVG(ua.time_taken_seconds) as avg_time
FROM users u
LEFT JOIN user_attempts ua ON u.user_id = ua.user_id
GROUP BY u.user_id;
```

---

### **8.4 Backend API Structure**

**backend/main.py:**

```python
from fastapi import FastAPI, Depends, HTTPException
from sqlalchemy.orm import Session
import database as db
import models, schemas
import llm_service, embeddings

app = FastAPI()

# Database dependency
def get_db():
  db_session = db.SessionLocal()
  try:
    yield db_session
  finally:
    db_session.close()

# ===== Questions API =====
@app.get("/questions")
def search_questions(
  query: str = "",
  filters: schemas.QuestionFilters = Depends(),
  db: Session = Depends(get_db)
):
  """Search questions with filters"""
  results = search_service.search_questions(query, filters, db)
  return results

@app.get("/questions/{question_id}")
def get_question(question_id: int, db: Session = Depends(get_db)):
  """Get specific question"""
  question = db.query(models.Question).filter_by(question_id=question_id).first()
  if not question:
    raise HTTPException(status_code=404, detail="Question not found")
  return question

@app.get("/questions/{question_id}/similar")
def get_similar_questions(question_id: int, limit: int = 10, db: Session = Depends(get_db)):
  """Get similar questions using graph + vectors"""
  # Graph-based similarity
  graph_similar = graph_service.find_similar(question_id, limit)
  
  # Vector-based similarity
  question_embedding = embeddings.get_embedding(question_id)
  vector_similar = embeddings.search_similar(question_embedding, limit)
  
  # Merge and rank
  return merge_and_deduplicate(graph_similar, vector_similar)

# ===== Attempts API =====
@app.post("/attempts")
def submit_attempt(
  attempt: schemas.AttemptSubmission,
  current_user: models.User = Depends(get_current_user),
  db: Session = Depends(get_db)
):
  """Submit an attempt"""
  # Store attempt
  db_attempt = models.UserAttempt(**attempt.dict(), user_id=current_user.user_id)
  db.add(db_attempt)
  db.commit()
  
  # Evaluate with LLM if explanation provided
  if attempt.explanation:
    evaluation = llm_service.evaluate_explanation(
      question=db.query(models.Question).get(attempt.question_id),
      user_solution=attempt.answer,
      user_explanation=attempt.explanation
    )
    db_attempt.llm_feedback = evaluation
    db.commit()
  
  # Update progress
  update_user_progress(current_user.user_id, attempt.question_id, db)
  
  # Check for weaknesses
  detect_and_update_weaknesses(current_user.user_id, db)
  
  return {"attempt_id": db_attempt.attempt_id, "evaluation": evaluation}

# ===== Training Plan API =====
@app.post("/training-plans")
def create_training_plan(
  plan_config: schemas.TrainingPlanConfig,
  current_user: models.User = Depends(get_current_user),
  db: Session = Depends(get_db)
):
  """Generate personalized training plan"""
  plan = training_plan_service.generate_plan(current_user.user_id, plan_config, db)
  return plan

@app.get("/training-plans/{plan_id}")
def get_training_plan(
  plan_id: int,
  current_user: models.User = Depends(get_current_user),
  db: Session = Depends(get_db)
):
  """Get training plan details"""
  plan = db.query(models.TrainingPlan).filter_by(
    plan_id=plan_id, 
    user_id=current_user.user_id
  ).first()
  if not plan:
    raise HTTPException(status_code=404)
  return plan

# ===== Assessment API =====
@app.post("/assessments")
def start_assessment(
  assessment_type: str,
  current_user: models.User = Depends(get_current_user),
  db: Session = Depends(get_db)
):
  """Start a new assessment"""
  assessment = assessment_service.create_assessment(
    user_id=current_user.user_id,
    assessment_type=assessment_type,
    db=db
  )
  return assessment

@app.post("/assessments/{assessment_id}/submit")
def submit_assessment(
  assessment_id: int,
  responses: List[schemas.AssessmentResponse],
  current_user: models.User = Depends(get_current_user),
  db: Session = Depends(get_db)
):
  """Submit assessment responses"""
  results = assessment_service.evaluate_assessment(
    assessment_id=assessment_id,
    responses=responses,
    db=db
  )
  return results

# ===== Analytics API =====
@app.get("/users/me/stats")
def get_user_stats(
  current_user: models.User = Depends(get_current_user),
  db: Session = Depends(get_db)
):
  """Get user statistics"""
  stats = analytics_service.get_user_stats(current_user.user_id, db)
  return stats

@app.get("/users/me/weaknesses")
def get_weaknesses(
  current_user: models.User = Depends(get_current_user),
  db: Session = Depends(get_db)
):
  """Get detected weaknesses"""
  weaknesses = db.query(models.Weakness).filter_by(
    user_id=current_user.user_id
  ).order_by(models.Weakness.severity.desc()).all()
  return weaknesses

# ===== LLM-Powered API =====
@app.post("/hint")
def get_hint(
  question_id: int,
  difficulty_level: str = "moderate",
  current_user: models.User = Depends(get_current_user),
  db: Session = Depends(get_db)
):
  """Get AI-generated hint"""
  question = db.query(models.Question).get(question_id)
  user_context = get_user_context(current_user.user_id, db)
  
  hint = llm_service.generate_hint(question, user_context, difficulty_level)
  return hint

@app.post("/questions/{question_id}/generate-variant")
def generate_variant(
  question_id: int,
  variation_type: str,
  db: Session = Depends(get_db)
):
  """Generate similar question using LLM"""
  original = db.query(models.Question).get(question_id)
  variant = llm_service.generate_similar_question(original, variation_type)
  
  # Store variant
  db.add(variant)
  db.commit()
  
  return variant
```

---

### **8.5 Quick Start Guide**

```bash
# 1. Clone repo (or create from scratch)
git clone https://github.com/yourusername/leetcode-training
cd leetcode-training

# 2. Start services
docker-compose up -d

# 3. Wait for services to be healthy
docker-compose ps

# 4. Initialize database
docker-compose exec postgres psql -U leetcode -d leetcode_training -f /init.sql

# 5. Pull LLM models
docker-compose exec ollama ollama pull llama2:70b

# 6. Seed initial data
python backend/scripts/seed_data.py

# 7. Access application
# Frontend: http://localhost:3000
# Backend API: http://localhost:8000
# API Docs: http://localhost:8000/docs
```

---

## 9. NEXT STEPS & ROADMAP

### **Phase 1: Core Setup (Week 1-2)**
- ✓ Set up database schema
- ✓ Implement user authentication
- ✓ Create question CRUD APIs
- ✓ Set up local LLM

### **Phase 2: Question Bank (Week 3-4)**
- ✓ Import initial 200 questions
- ✓ Generate embeddings
- ✓ Build graph relationships
- ✓ Test search functionality

### **Phase 3: Training Features (Week 5-6)**
- ✓ Implement training plan generator
- ✓ Build spaced repetition system
- ✓ Create progress tracking
- ✓ Add weakness detection

### **Phase 4: LLM Integration (Week 7-8)**
- ✓ Implement explanation evaluator
- ✓ Build hint generator
- ✓ Create variant generator
- ✓ Test assessment system

### **Phase 5: Frontend (Week 9-10)**
- ✓ Build question interface
- ✓ Create progress dashboard
- ✓ Implement code editor
- ✓ Add analytics views

### **Phase 6: Testing & Optimization (Week 11-12)**
- ✓ Test all features
- ✓ Optimize query performance
- ✓ Calibrate difficulty scores
- ✓ Gather feedback & iterate

---

## SUMMARY

This technical architecture provides:

1. **Hybrid Database** (Relational + Graph + Vector) for rich relationships
2. **Custom Difficulty Scoring** (7-dimensional, personalized)
3. **Fair Assessment** (anti-memorization, understanding-focused)
4. **Weakness Detection** (multi-method, LLM-powered)
5. **Flexible Training Plans** (adaptive, spaced repetition)
6. **Advanced Search** (fulltext + semantic + graph)
7. **LLM Integration** (RAG, evaluation, generation)
8. **Local Deployment** (Docker-based, privacy-focused)

**All designed to run locally for personal practice and testing.**

Ready to implement! 🚀
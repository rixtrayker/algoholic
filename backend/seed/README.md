# Database Seeding Guide

This directory contains seed data for populating the Algoholic database with initial problems, questions, and topics.

## Quick Start

```bash
cd backend
go run cmd/seed/main.go
```

## What Gets Seeded

### Topics (10 total)
- Arrays
- Hash Table
- Two Pointers
- Sliding Window
- Binary Search
- Dynamic Programming
- Graph Traversal
- Trees
- Stack
- Linked List

### Problems (12 total)

**Easy (5):**
- Two Sum (Hash Table)
- Best Time to Buy and Sell Stock (Array, DP)
- Contains Duplicate (Hash Table)
- Valid Parentheses (Stack)
- Reverse Linked List (Linked List)

**Medium (5):**
- 3Sum (Two Pointers)
- Search in Rotated Sorted Array (Binary Search)
- Longest Substring Without Repeating Characters (Sliding Window)
- Number of Islands (Graph Traversal)
- Maximum Subarray (Dynamic Programming)

**Hard (2):**
- Trapping Rain Water (Two Pointers, DP)
- Merge k Sorted Lists (Linked List, Heap)

### Questions (10 total)

Questions cover multiple types:
- **Complexity Analysis** - Time/space complexity questions
- **Data Structure Selection** - Why use a particular data structure
- **Code Completion** - Implement solutions with test cases
- **Pattern Recognition** - Identify algorithmic patterns
- **Algorithm Explanation** - Explain how algorithms work
- **Concept Explanation** - Fundamental CS concepts
- **Debugging** - Common errors and how to avoid them

## Features

### Test Cases for Code Questions

Code completion questions include test cases that will be validated using the Judge0 code execution engine:

```go
CorrectAnswer: jsonbMap(map[string]interface{}{
    "test_cases": []map[string]string{
        {"input": "[2,7,11,15]\n9", "expected": "[0, 1]"},
        {"input": "[3,2,4]\n6", "expected": "[1, 2]"},
        {"input": "[3,3]\n6", "expected": "[0, 1]"},
    },
})
```

### Multiple Acceptable Answers

Text questions support multiple correct answers for flexible validation:

```go
CorrectAnswer: jsonbMap(map[string]interface{}{
    "answer": []string{
        "hash map allows O(1) lookup",
        "hash table provides constant time access",
        "we can check if complement exists instantly",
    },
})
```

### Wrong Answer Explanations

Multiple choice questions include explanations for each wrong answer:

```go
WrongAnswerExplanations: jsonbMap(map[string]interface{}{
    "a": "O(n^2) would be the brute force approach...",
    "b": "O(n log n) would apply if we sorted...",
    "d": "O(1) is impossible since...",
})
```

## Running the Seeder

### Prerequisites

1. PostgreSQL database running
2. Database configuration in `config.yaml` or environment variables
3. Go modules initialized

### Execution

```bash
# From backend directory
cd backend

# Run the seeder
go run cmd/seed/main.go
```

### Output Example

```
🌱 Starting database seeding...
📋 Running database migrations...
✅ Migrations complete

🏷️  Seeding topics...
   ✅ [1/10] Seeded topic: Arrays
   ✅ [2/10] Seeded topic: Hash Table
   ...
✅ Seeded 10/10 topics

🧩 Seeding problems...
   ✅ [1/12] Seeded problem: Two Sum (Difficulty: 15)
   ✅ [2/12] Seeded problem: Best Time to Buy and Sell Stock (Difficulty: 20)
   ...
✅ Seeded 12/12 problems

❓ Seeding questions...
   ✅ [1/10] Seeded question: What is the time complexity of...
   ...
✅ Seeded 10/10 questions

🔗 Creating problem-topic relationships...
✅ Created 18 problem-topic relationships

============================================================
🎉 Database Seeding Complete!
============================================================
📊 Database Statistics:
   Topics: 10
   Problems: 12
   Questions: 10
   Relationships: 18
============================================================

✅ You can now start practicing!
   Backend: http://localhost:4000
   Frontend: http://localhost:5173
```

## Re-running the Seeder

The seeder is idempotent - it checks if items already exist before inserting:

- **Topics** - Checked by slug
- **Problems** - Checked by slug
- **Questions** - Will create duplicates (run seeder only once)
- **Relationships** - Uses FirstOrCreate to avoid duplicates

## Adding More Seed Data

Edit `seed_data.go` and add items to the respective functions:

- `GetSeedTopics()` - Add more topics
- `GetSeedProblems()` - Add more problems
- `GetSeedQuestions()` - Add more questions

### Example: Adding a New Problem

```go
{
    LeetcodeNumber:     intPtr(70),
    Title:              "Climbing Stairs",
    Slug:               "climbing-stairs",
    Description:        "You are climbing a staircase...",
    Constraints:        models.StringArray{"1 <= n <= 45"},
    Examples:           jsonbFromString(`[...]`),
    Hints:              models.StringArray{"Use dynamic programming"},
    DifficultyScore:    25.0,
    OfficialDifficulty: strPtr("Easy"),
    PrimaryPattern:     strPtr("Dynamic Programming"),
    TimeComplexity:     strPtr("O(n)"),
    SpaceComplexity:    strPtr("O(1)"),
},
```

## Testing

After seeding, test the data:

```bash
# Check problem count
curl http://localhost:4000/api/problems

# Get a random question
curl http://localhost:4000/api/questions/random

# Get specific problem
curl http://localhost:4000/api/problems/two-sum
```

## Troubleshooting

### "Failed to connect to database"
- Check PostgreSQL is running
- Verify database credentials in config.yaml
- Ensure database exists

### "Migration failed"
- Check migrations in `/migrations` directory
- Verify migration files are valid SQL
- Check database user has CREATE TABLE permissions

### "Failed to seed"
- Check for unique constraint violations
- Verify foreign key references exist
- Check JSON formatting in seed data

## Next Steps

After seeding:

1. Test API endpoints with Postman collection
2. Create a test user and practice questions
3. Verify answer validation works correctly
4. Check difficulty scoring algorithms

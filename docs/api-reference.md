## API Reference

**📦 Postman Collection**: For the complete, up-to-date API specification with tests and examples, see:
- Collection: [`postman/algoholic-api.postman_collection.json`](../postman/algoholic-api.postman_collection.json)
- Documentation: [`postman/README.md`](../postman/README.md)
- Alignment Report: [`api-frontend-alignment.md`](./api-frontend-alignment.md)

**✅ Current Status**: 22 endpoints across 7 categories, all tested with Newman

Base URL: `http://localhost:4000/api`

### Authentication

All protected endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <token>
```

**Testing**: Run `cd postman && ./run-tests.sh` to execute all API tests via Newman

---

### Authentication Endpoints

#### POST /auth/register
Register a new user account.

**Request:**
```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**Response:** `201 Created`
```json
{
  "message": "Registration successful",
  "token": "eyJhbGc...",
  "user": {
    "user_id": 1,
    "username": "johndoe",
    "email": "john@example.com"
  }
}
```

#### POST /auth/login
Login to an existing account.

**Request:**
```json
{
  "username": "johndoe",
  "password": "SecurePass123!"
}
```

**Response:** `200 OK`
```json
{
  "message": "Login successful",
  "token": "eyJhbGc...",
  "user": {
    "user_id": 1,
    "username": "johndoe",
    "email": "john@example.com"
  }
}
```

#### GET /auth/me 🔒
Get current user information.

**Response:** `200 OK`
```json
{
  "user": {
    "user_id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "current_streak_days": 7,
    "total_study_time_seconds": 3600
  }
}
```

#### POST /auth/change-password 🔒
Change user password.

**Request:**
```json
{
  "old_password": "OldPass123!",
  "new_password": "NewPass123!"
}
```

---

### Problem Endpoints

#### GET /problems
Get all problems with optional filters.

**Query Parameters:**
- `min_difficulty` (float, default: 0)
- `max_difficulty` (float, default: 100)
- `pattern` (string, optional)
- `limit` (int, default: 20)
- `offset` (int, default: 0)

**Response:** `200 OK`
```json
{
  "problems": [...],
  "total": 100,
  "limit": 20,
  "offset": 0
}
```

#### GET /problems/:id
Get a specific problem by ID.

**Response:** `200 OK`
```json
{
  "problem_id": 1,
  "title": "Two Sum",
  "slug": "two-sum",
  "description": "...",
  "difficulty_score": 15.5,
  "primary_pattern": "Hash Table",
  "examples": {...},
  "hints": [...]
}
```

#### GET /problems/slug/:slug
Get a problem by slug.

#### GET /problems/:id/topics
Get all topics associated with a problem.

**Response:** `200 OK`
```json
{
  "topics": [
    {
      "topic_id": 1,
      "name": "Arrays",
      "slug": "arrays"
    }
  ]
}
```

#### GET /problems/search
Search problems by title or description.

**Query Parameters:**
- `q` (string, required)
- `limit` (int, default: 20)
- `offset` (int, default: 0)

---

### Question Endpoints

#### GET /questions
Get questions with optional filters.

**Query Parameters:**
- `type` (string, optional) - Question type filter
- `min_difficulty` (float, default: 0)
- `max_difficulty` (float, default: 100)
- `limit` (int, default: 20)
- `offset` (int, default: 0)

#### GET /questions/:id
Get a specific question by ID.

**Response:** `200 OK`
```json
{
  "question_id": 1,
  "question_type": "complexity_analysis",
  "question_text": "What is the time complexity?",
  "answer_options": {...},
  "difficulty_score": 25.0
}
```

#### GET /questions/random
Get a random question with optional filters.

**Query Parameters:**
- `type` (string, optional)
- `min_difficulty` (float, default: 0)
- `max_difficulty` (float, default: 100)

#### POST /questions/:id/answer 🔒
Submit an answer to a question.

**Request:**
```json
{
  "user_answer": {
    "answer": "A"
  },
  "time_taken_seconds": 45,
  "hints_used": 1,
  "confidence_level": 3,
  "training_plan_id": null
}
```

**Response:** `200 OK`
```json
{
  "is_correct": true,
  "correct_answer": {"answer": "A"},
  "explanation": "Detailed explanation...",
  "wrong_answer_explanation": "",
  "attempt_id": 123,
  "points_earned": 250
}
```

#### GET /questions/:id/attempts 🔒
Get user's previous attempts for a question.

#### GET /problems/:problemId/questions
Get all questions for a specific problem.

---

### User Endpoints

All user endpoints require authentication 🔒

#### GET /users/me/stats
Get comprehensive user statistics.

**Response:** `200 OK`
```json
{
  "total_attempts": 150,
  "correct_attempts": 120,
  "accuracy_rate": 80.0,
  "total_study_time_seconds": 7200,
  "current_streak_days": 14,
  "problems_attempted": 45,
  "problems_solved": 30,
  "questions_answered": 150,
  "average_difficulty": 52.5,
  "strong_topics": ["Arrays", "Hash Tables"],
  "weak_topics": ["Dynamic Programming"]
}
```

#### GET /users/me/weaknesses
Get user's weak topics.

**Query Parameters:**
- `limit` (int, default: 10)

**Response:** `200 OK`
```json
{
  "weak_topics": [
    {
      "topic_id": 5,
      "name": "Dynamic Programming",
      "proficiency_level": 35.5
    }
  ],
  "count": 1
}
```

#### GET /users/me/recommendations
Get personalized recommendations.

**Response:** `200 OK`
```json
{
  "recommendations": [
    {
      "type": "practice_topic",
      "topic": {...},
      "reason": "Low proficiency - needs practice",
      "priority": "high",
      "action": "Practice questions for this topic"
    }
  ],
  "count": 3
}
```

#### GET /users/me/review-queue
Get topics due for spaced repetition review.

**Response:** `200 OK`
```json
{
  "review_queue": [
    {
      "user_id": 1,
      "topic_id": 3,
      "proficiency_level": 75.0,
      "next_review_at": "2025-02-08T10:00:00Z"
    }
  ],
  "count": 1
}
```

#### GET /users/me/skills
Get all user skills across topics.

#### GET /users/me/skills/:topicId
Get progress for a specific topic.

#### GET /users/me/preferences
Get user preferences.

#### PUT /users/me/preferences
Update user preferences.

**Request:**
```json
{
  "theme": "dark",
  "notifications_enabled": true,
  "preferred_difficulty": "medium"
}
```

#### GET /users/me/attempts
Get recent attempts.

**Query Parameters:**
- `limit` (int, default: 20)

---

### Training Plan Endpoints

All training plan endpoints require authentication 🔒

#### POST /training-plans
Create a new training plan.

**Request:**
```json
{
  "name": "30-Day DP Bootcamp",
  "description": "Master dynamic programming in 30 days",
  "plan_type": "custom",
  "target_topics": [5, 8, 12],
  "target_patterns": ["Dynamic Programming", "Memoization"],
  "duration_days": 30,
  "questions_per_day": 5,
  "difficulty_min": 40.0,
  "difficulty_max": 80.0,
  "adaptive_difficulty": true
}
```

**Response:** `201 Created`
```json
{
  "message": "Training plan created successfully",
  "plan": {
    "plan_id": 1,
    "user_id": 1,
    "name": "30-Day DP Bootcamp",
    "status": "active",
    "progress_percentage": 0.0
  }
}
```

#### GET /training-plans
Get all training plans for the current user.

**Response:** `200 OK`
```json
{
  "plans": [...],
  "count": 3
}
```

#### GET /training-plans/:id
Get a specific training plan.

#### GET /training-plans/:id/next
Get the next question in the training plan.

**Response:** `200 OK`
```json
{
  "question_id": 45,
  "question_text": "...",
  "difficulty_score": 55.0
}
```

#### GET /training-plans/:id/items
Get all items in a training plan.

**Response:** `200 OK`
```json
{
  "items": [
    {
      "item_id": 1,
      "sequence_number": 1,
      "day_number": 1,
      "is_completed": false,
      "question_id": 45
    }
  ],
  "count": 150
}
```

#### GET /training-plans/:id/today
Get today's questions from the plan.

**Response:** `200 OK`
```json
{
  "questions": [...],
  "count": 5
}
```

#### POST /training-plans/:id/items/:itemId/complete
Mark a training plan item as completed.

**Response:** `200 OK`
```json
{
  "message": "Item marked as completed"
}
```

#### POST /training-plans/:id/pause
Pause a training plan.

#### POST /training-plans/:id/resume
Resume a paused training plan.

#### DELETE /training-plans/:id
Delete a training plan.

---

### Health Endpoint

#### GET /health
Check API health status.

**Response:** `200 OK`
```json
{
  "status": "healthy",
  "app": "Algoholic API",
  "version": "1.0.0",
  "environment": "development"
}
```

---

### Error Responses

All endpoints may return error responses in this format:

**4xx Client Errors:**
```json
{
  "error": "Error message describing what went wrong"
}
```

**Common Status Codes:**
- `400 Bad Request` - Invalid request body or parameters
- `401 Unauthorized` - Missing or invalid authentication token
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource already exists (e.g., duplicate username)
- `500 Internal Server Error` - Server error

---

### Rate Limiting

Currently not implemented. Will be added in future versions.

---

### Pagination

Endpoints that return lists support pagination via `limit` and `offset` parameters:
- `limit`: Number of items to return (default: 20, max: 100)
- `offset`: Number of items to skip (default: 0)

**Example:**
```
GET /api/problems?limit=10&offset=20
```

---

### Filtering

Many endpoints support filtering by difficulty:
- `min_difficulty`: Minimum difficulty score (0-100)
- `max_difficulty`: Maximum difficulty score (0-100)

**Example:**
```
GET /api/questions?min_difficulty=50&max_difficulty=75
```

---

## Complete Endpoint List

**Total: 36 endpoints**

### Public Endpoints (3)
- `GET /health`
- `POST /api/auth/register`
- `POST /api/auth/login`

### Protected Endpoints (33)
- Authentication: 2
- Problems: 5
- Questions: 6
- Users: 9
- Training Plans: 11

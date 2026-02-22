# Algoholic Implementation Plan
## Phase 1 Critical Items & Phase 2 Roadmap

**Created:** February 12, 2026
**Status:** Phase 1 Foundation - 100% Complete, Moving to Immediate Fixes

---

## Table of Contents

1. [Current State Summary](#current-state-summary)
2. [Phase 1 Critical Items](#phase-1-critical-items)
3. [Phase 2 Intelligence Layer](#phase-2-intelligence-layer)
4. [Implementation Timeline](#implementation-timeline)
5. [Technical Specifications](#technical-specifications)
6. [Success Criteria](#success-criteria)

---

## Current State Summary

### ✅ Phase 1 Complete (Foundation)
- **Backend API**: 22 endpoints across 7 categories
- **Database**: PostgreSQL with 13 tables, full schema
- **Frontend**: React 19 + Vite, 7 pages, 46 tests passing
- **Testing**: Postman collection with Newman automation
- **Auth**: JWT-based authentication working

### 🔴 Critical Gaps Blocking User Experience
1. **Answer validation TODOs** in `backend/services/question_service.go:208,216`
2. **Empty database** - no problems, questions, or topics to practice
3. **Incomplete service logic** - difficulty scoring, recommendations, streak tracking

---

## Phase 1 Critical Items

### 1. Answer Validation Implementation

**Location:** `backend/services/question_service.go`

#### 1.1 Code Execution Validation (Line 208)

**Current Issue:**
```go
// TODO: Implement actual code execution and validation
// For now, just check if code is provided
code, ok := userAnswer["code"].(string)
return ok && len(code) > 0
```

**Implementation Plan:**

**Option A: Docker-based Sandbox (Recommended)**
```go
func (s *QuestionService) CheckCode(question *models.Question, userAnswer map[string]interface{}) bool {
    code, ok := userAnswer["code"].(string)
    if !ok || len(code) == 0 {
        return false
    }

    language, _ := userAnswer["language"].(string)
    if language == "" {
        language = "python" // default
    }

    // Get test cases from correct_answer
    testCases, ok := question.CorrectAnswer["test_cases"].([]interface{})
    if !ok {
        return false
    }

    // Run code in isolated Docker container
    executor := NewCodeExecutor()
    results, err := executor.RunTests(code, language, testCases)
    if err != nil {
        return false
    }

    // All test cases must pass
    return results.AllPassed
}
```

**Option B: Judge0 API Integration (Quick Win)**
```go
func (s *QuestionService) CheckCode(question *models.Question, userAnswer map[string]interface{}) bool {
    code, _ := userAnswer["code"].(string)
    language, _ := userAnswer["language"].(string)

    // Use Judge0 CE (free, open-source)
    judge0Client := NewJudge0Client("http://localhost:2358")

    testCases := question.CorrectAnswer["test_cases"].([]interface{})

    for _, tc := range testCases {
        testCase := tc.(map[string]interface{})
        input := testCase["input"].(string)
        expected := testCase["expected"].(string)

        result, err := judge0Client.Submit(code, language, input)
        if err != nil || result.Stdout != expected {
            return false
        }
    }

    return true
}
```

**Recommended Approach:** Start with Option B (Judge0) for Phase 1, migrate to Option A in Phase 4.

**Code Execution Service:**
```go
// backend/services/code_executor.go
type CodeExecutor struct {
    judge0URL string
    timeout   int
}

type TestCase struct {
    Input    string
    Expected string
}

type ExecutionResult struct {
    AllPassed   bool
    PassedCount int
    TotalCount  int
    Failures    []FailureDetail
    TimeTaken   int
    MemoryUsed  int
}

func (ce *CodeExecutor) RunTests(code, language string, testCases []interface{}) (*ExecutionResult, error) {
    // Implementation details
}
```

**Estimated Time:** 1-2 days

---

#### 1.2 Text Answer Fuzzy Matching (Line 216)

**Current Issue:**
```go
// TODO: Implement fuzzy matching or keyword-based validation
// For now, simple exact match
userText, ok := userAnswer["answer"].(string)
if !ok {
    return false
}

correctText, ok := question.CorrectAnswer["answer"].(string)
if !ok {
    return false
}

return userText == correctText
```

**Implementation Plan:**

**Fuzzy Matching Algorithm:**
```go
func (s *QuestionService) CheckText(question *models.Question, userAnswer map[string]interface{}) bool {
    userText, ok := userAnswer["answer"].(string)
    if !ok {
        return false
    }

    correctAnswer, ok := question.CorrectAnswer["answer"]
    if !ok {
        return false
    }

    // Support multiple correct answer formats
    switch v := correctAnswer.(type) {
    case string:
        return s.FuzzyMatch(userText, v, 0.85)
    case []interface{}:
        // Multiple acceptable answers
        for _, ans := range v {
            if ansStr, ok := ans.(string); ok {
                if s.FuzzyMatch(userText, ansStr, 0.85) {
                    return true
                }
            }
        }
        return false
    default:
        return false
    }
}

func (s *QuestionService) FuzzyMatch(text1, text2 string, threshold float64) bool {
    // Normalize inputs
    t1 := s.NormalizeText(text1)
    t2 := s.NormalizeText(text2)

    // Strategy 1: Exact match after normalization
    if t1 == t2 {
        return true
    }

    // Strategy 2: Levenshtein distance
    distance := levenshtein.Distance(t1, t2)
    maxLen := max(len(t1), len(t2))
    similarity := 1.0 - float64(distance)/float64(maxLen)

    if similarity >= threshold {
        return true
    }

    // Strategy 3: Keyword matching (for conceptual answers)
    if s.HasRequiredKeywords(t1, t2) {
        return true
    }

    return false
}

func (s *QuestionService) NormalizeText(text string) string {
    // Convert to lowercase
    text = strings.ToLower(text)

    // Remove extra whitespace
    text = strings.TrimSpace(text)
    text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

    // Remove punctuation (except dashes in technical terms)
    text = regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(text, "")

    // Common substitutions
    replacements := map[string]string{
        "o(n)":     "on",
        "o(log n)": "ologn",
        "o(1)":     "o1",
    }

    for old, new := range replacements {
        text = strings.ReplaceAll(text, old, new)
    }

    return text
}

func (s *QuestionService) HasRequiredKeywords(userText, correctText string) bool {
    // Extract key technical terms from correct answer
    keywords := s.ExtractKeywords(correctText)

    // User must mention at least 70% of keywords
    matchCount := 0
    for _, keyword := range keywords {
        if strings.Contains(userText, keyword) {
            matchCount++
        }
    }

    threshold := 0.7
    return float64(matchCount)/float64(len(keywords)) >= threshold
}

func (s *QuestionService) ExtractKeywords(text string) []string {
    // Technical terms, algorithm names, data structures
    technicalTerms := []string{
        "array", "hashmap", "tree", "graph", "dynamic programming",
        "binary search", "two pointers", "sliding window", "dfs", "bfs",
        "stack", "queue", "heap", "trie", "union find",
    }

    keywords := []string{}
    textLower := strings.ToLower(text)

    for _, term := range technicalTerms {
        if strings.Contains(textLower, term) {
            keywords = append(keywords, term)
        }
    }

    return keywords
}
```

**Dependencies:**
```go
import "github.com/agnivade/levenshtein"
```

**Estimated Time:** 1 day

---

### 2. Database Seed Data Creation

**Goal:** Populate database with 50-100 problems and 200-300 questions for MVP

#### 2.1 Seed Data Structure

**Create:** `backend/seed/seed_data.go`

```go
package seed

import (
    "encoding/json"
    "github.com/yourusername/algoholic/models"
)

func GetSeedProblems() []models.Problem {
    return []models.Problem{
        {
            LeetcodeNumber:  intPtr(1),
            Title:           "Two Sum",
            Slug:            "two-sum",
            Description:     "Given an array of integers nums and an integer target...",
            Constraints:     []string{"2 <= nums.length <= 10^4", "-10^9 <= nums[i] <= 10^9"},
            Examples:        jsonb(`[{"input": "nums = [2,7,11,15], target = 9", "output": "[0,1]"}]`),
            Hints:           []string{"Try using a hash map", "For each element, check if target-element exists"},
            DifficultyScore: 15.0,
            OfficialDifficulty: strPtr("Easy"),
            PrimaryPattern:  strPtr("Hash Table"),
            SecondaryPatterns: []string{"Array"},
            TimeComplexity:  strPtr("O(n)"),
            SpaceComplexity: strPtr("O(n)"),
        },
        {
            LeetcodeNumber:  intPtr(15),
            Title:           "3Sum",
            Slug:            "3sum",
            Description:     "Given an integer array nums, return all triplets...",
            DifficultyScore: 45.0,
            PrimaryPattern:  strPtr("Two Pointers"),
            SecondaryPatterns: []string{"Array", "Sorting"},
        },
        // Add 48 more problems...
    }
}

func GetSeedQuestions() []models.Question {
    return []models.Question{
        {
            ProblemID:       intPtr(1), // Two Sum
            QuestionType:    "complexity_analysis",
            QuestionSubtype: strPtr("time_complexity"),
            QuestionFormat:  "multiple_choice",
            QuestionText:    "What is the time complexity of the optimal solution for Two Sum using a hash map?",
            AnswerOptions: jsonb(`[
                {"id": "a", "text": "O(n^2)"},
                {"id": "b", "text": "O(n log n)"},
                {"id": "c", "text": "O(n)"},
                {"id": "d", "text": "O(1)"}
            ]`),
            CorrectAnswer: jsonb(`{"answer": "c"}`),
            Explanation:   "We traverse the array once, and hash map operations are O(1), giving us O(n) overall.",
            WrongAnswerExplanations: jsonb(`{
                "a": "O(n^2) would be the brute force approach with nested loops.",
                "b": "O(n log n) would apply if we sorted the array first.",
                "d": "O(1) is impossible since we must at least read all elements."
            }`),
            RelatedConcepts: []string{"Hash Table", "Time Complexity"},
            DifficultyScore: 20.0,
            EstimatedTimeSeconds: intPtr(60),
        },
        {
            ProblemID:      intPtr(1),
            QuestionType:   "code_completion",
            QuestionFormat: "code",
            QuestionText:   "Complete the Two Sum function in Python",
            QuestionData: jsonb(`{
                "template": "def twoSum(nums: List[int], target: int) -> List[int]:\n    # Your code here\n    pass",
                "language": "python"
            }`),
            CorrectAnswer: jsonb(`{
                "test_cases": [
                    {"input": "[2,7,11,15]\n9", "expected": "[0,1]"},
                    {"input": "[3,2,4]\n6", "expected": "[1,2]"},
                    {"input": "[3,3]\n6", "expected": "[0,1]"}
                ]
            }`),
            Explanation:    "Use a hash map to store seen numbers and their indices.",
            DifficultyScore: 25.0,
        },
        // Add 198 more questions...
    }
}
```

#### 2.2 Seed Script

**Create:** `backend/cmd/seed/main.go`

```go
package main

import (
    "log"
    "github.com/yourusername/algoholic/config"
    "github.com/yourusername/algoholic/database"
    "github.com/yourusername/algoholic/seed"
)

func main() {
    // Load config
    cfg := config.LoadConfig()

    // Connect to database
    db := database.Connect(cfg.Database)

    // Run migrations first
    if err := database.RunMigrations(db); err != nil {
        log.Fatalf("Migration failed: %v", err)
    }

    log.Println("Seeding problems...")
    problems := seed.GetSeedProblems()
    for _, p := range problems {
        if err := db.Create(&p).Error; err != nil {
            log.Printf("Failed to seed problem %s: %v", p.Slug, err)
        }
    }
    log.Printf("Seeded %d problems", len(problems))

    log.Println("Seeding questions...")
    questions := seed.GetSeedQuestions()
    for _, q := range questions {
        if err := db.Create(&q).Error; err != nil {
            log.Printf("Failed to seed question: %v", err)
        }
    }
    log.Printf("Seeded %d questions", len(questions))

    log.Println("Seeding topics...")
    topics := seed.GetSeedTopics()
    for _, t := range topics {
        if err := db.Create(&t).Error; err != nil {
            log.Printf("Failed to seed topic %s: %v", t.Slug, err)
        }
    }
    log.Printf("Seeded %d topics", len(topics))

    log.Println("✅ Database seeding complete!")
}
```

**Run:**
```bash
cd backend
go run cmd/seed/main.go
```

**Estimated Time:** 2-3 days (including creating quality seed data)

---

### 3. Difficulty Scoring Algorithm Implementation

**Location:** `backend/utils/difficulty_scorer.go`

```go
package utils

import (
    "github.com/yourusername/algoholic/models"
)

type DifficultyComponents struct {
    Conceptual     float64 // 0-100
    Algorithm      float64 // 0-100
    Implementation float64 // 0-100
    Pattern        float64 // 0-100
    EdgeCases      float64 // 0-100
    TimePressure   float64 // 0-100
}

func CalculateDifficultyScore(problem *models.Problem) float64 {
    components := DifficultyComponents{
        Conceptual:     ScoreConceptual(problem),
        Algorithm:      ScoreAlgorithm(problem),
        Implementation: ScoreImplementation(problem),
        Pattern:        ScorePatternRecognition(problem),
        EdgeCases:      ScoreEdgeCases(problem),
        TimePressure:   ScoreTimePressure(problem),
    }

    weights := []float64{0.25, 0.20, 0.15, 0.20, 0.10, 0.10}
    scores := []float64{
        components.Conceptual,
        components.Algorithm,
        components.Implementation,
        components.Pattern,
        components.EdgeCases,
        components.TimePressure,
    }

    total := 0.0
    for i, score := range scores {
        total += score * weights[i]
    }

    return clamp(total, 0, 100)
}

func ScoreConceptual(problem *models.Problem) float64 {
    // Count unique concepts needed
    conceptCount := len(problem.SecondaryPatterns) + 1 // primary + secondary

    // More concepts = higher difficulty
    switch {
    case conceptCount == 1:
        return 20.0
    case conceptCount == 2:
        return 40.0
    case conceptCount == 3:
        return 60.0
    case conceptCount >= 4:
        return 80.0
    default:
        return 30.0
    }
}

func ScoreAlgorithm(problem *models.Problem) float64 {
    // Based on time complexity
    if problem.TimeComplexity == nil {
        return 50.0
    }

    complexity := *problem.TimeComplexity

    complexityScores := map[string]float64{
        "O(1)":         10.0,
        "O(log n)":     20.0,
        "O(n)":         30.0,
        "O(n log n)":   50.0,
        "O(n^2)":       60.0,
        "O(n^3)":       80.0,
        "O(2^n)":       90.0,
        "O(n!)":        95.0,
    }

    if score, ok := complexityScores[complexity]; ok {
        return score
    }

    return 50.0
}

func ScoreImplementation(problem *models.Problem) float64 {
    // Based on problem description length and constraints
    descLength := len(problem.Description)
    constraintCount := len(problem.Constraints)

    baseScore := 30.0

    // Longer description = more complex implementation
    if descLength > 500 {
        baseScore += 20.0
    } else if descLength > 300 {
        baseScore += 10.0
    }

    // More constraints = trickier implementation
    baseScore += float64(constraintCount) * 5.0

    return clamp(baseScore, 0, 100)
}

func ScorePatternRecognition(problem *models.Problem) float64 {
    if problem.PrimaryPattern == nil {
        return 60.0 // Unknown pattern = harder to recognize
    }

    pattern := *problem.PrimaryPattern

    // Common patterns are easier to recognize
    patternDifficulty := map[string]float64{
        "Array":              20.0,
        "Hash Table":         25.0,
        "Two Pointers":       35.0,
        "Sliding Window":     40.0,
        "Binary Search":      35.0,
        "Dynamic Programming": 70.0,
        "Graph":              60.0,
        "Backtracking":       75.0,
        "Trie":               65.0,
    }

    if score, ok := patternDifficulty[pattern]; ok {
        return score
    }

    return 50.0
}

func ScoreEdgeCases(problem *models.Problem) float64 {
    // Analyze constraints for edge case complexity
    constraintCount := len(problem.Constraints)

    baseScore := float64(constraintCount) * 10.0

    // Check for specific edge case indicators
    description := strings.ToLower(problem.Description)

    if strings.Contains(description, "negative") {
        baseScore += 10.0
    }
    if strings.Contains(description, "duplicate") {
        baseScore += 10.0
    }
    if strings.Contains(description, "empty") {
        baseScore += 5.0
    }

    return clamp(baseScore, 0, 100)
}

func ScoreTimePressure(problem *models.Problem) float64 {
    // Based on expected solve time
    // For now, use a simple heuristic
    // Will be updated with real user data later

    if problem.OfficialDifficulty == nil {
        return 50.0
    }

    difficulty := *problem.OfficialDifficulty

    timePressure := map[string]float64{
        "Easy":   20.0,
        "Medium": 50.0,
        "Hard":   80.0,
    }

    if score, ok := timePressure[difficulty]; ok {
        return score
    }

    return 50.0
}

func clamp(value, min, max float64) float64 {
    if value < min {
        return min
    }
    if value > max {
        return max
    }
    return value
}

// Personalized difficulty based on user proficiency
func PersonalizedDifficulty(problemID, userID int, db *gorm.DB) float64 {
    // Get base difficulty
    var problem models.Problem
    db.First(&problem, problemID)

    baseDifficulty := problem.DifficultyScore

    // Get user's proficiency in related topics
    var avgProficiency float64
    db.Table("user_skills").
        Select("AVG(proficiency_level)").
        Where("user_id = ? AND topic_id IN (?)",
            userID,
            db.Table("problem_topics").Select("topic_id").Where("problem_id = ?", problemID),
        ).
        Scan(&avgProficiency)

    // Adjust difficulty based on user proficiency
    // Strong user (proficiency > 70) → easier (reduce difficulty)
    // Weak user (proficiency < 30) → harder (increase difficulty)
    adjustment := (50 - avgProficiency) * 0.6

    return clamp(baseDifficulty + adjustment, 0, 100)
}
```

**Estimated Time:** 1 day

---

### 4. Recommendation Logic Implementation

**Location:** `backend/services/recommendation_service.go`

```go
package services

import (
    "gorm.io/gorm"
    "github.com/yourusername/algoholic/models"
)

type RecommendationService struct {
    db *gorm.DB
}

func NewRecommendationService(db *gorm.DB) *RecommendationService {
    return &RecommendationService{db: db}
}

type Recommendation struct {
    QuestionID int     `json:"question_id"`
    ProblemID  *int    `json:"problem_id,omitempty"`
    Reason     string  `json:"reason"`
    Priority   float64 `json:"priority"`
}

func (s *RecommendationService) GetRecommendations(userID int, limit int) ([]Recommendation, error) {
    recommendations := []Recommendation{}

    // Strategy 1: Address weaknesses (highest priority)
    weaknessRecs := s.GetWeaknessBasedRecommendations(userID, limit/3)
    recommendations = append(recommendations, weaknessRecs...)

    // Strategy 2: Progressive difficulty (medium priority)
    progressRecs := s.GetProgressiveRecommendations(userID, limit/3)
    recommendations = append(recommendations, progressRecs...)

    // Strategy 3: Spaced repetition (lower priority)
    reviewRecs := s.GetSpacedRepetitionRecommendations(userID, limit/3)
    recommendations = append(recommendations, reviewRecs...)

    // Sort by priority and return top N
    sort.Slice(recommendations, func(i, j int) bool {
        return recommendations[i].Priority > recommendations[j].Priority
    })

    if len(recommendations) > limit {
        recommendations = recommendations[:limit]
    }

    return recommendations, nil
}

func (s *RecommendationService) GetWeaknessBasedRecommendations(userID int, limit int) []Recommendation {
    // Find user's weak topics
    var weakTopics []struct {
        TopicID          int
        ProficiencyLevel float64
    }

    s.db.Table("user_skills").
        Select("topic_id, proficiency_level").
        Where("user_id = ? AND proficiency_level < 50", userID).
        Order("proficiency_level ASC").
        Limit(limit).
        Scan(&weakTopics)

    recommendations := []Recommendation{}

    for _, wt := range weakTopics {
        // Find questions for this weak topic
        var questions []models.Question
        s.db.Joins("JOIN problem_topics ON problems.problem_id = problem_topics.problem_id").
            Where("problem_topics.topic_id = ? AND questions.difficulty_score < 50", wt.TopicID).
            Limit(2).
            Find(&questions)

        for _, q := range questions {
            recommendations = append(recommendations, Recommendation{
                QuestionID: q.QuestionID,
                ProblemID:  q.ProblemID,
                Reason:     fmt.Sprintf("Practice weak topic (proficiency: %.0f%%)", wt.ProficiencyLevel),
                Priority:   90.0 - wt.ProficiencyLevel, // Lower proficiency = higher priority
            })
        }
    }

    return recommendations
}

func (s *RecommendationService) GetProgressiveRecommendations(userID int, limit int) []Recommendation {
    // Get user's current average difficulty level
    var avgDifficulty float64
    s.db.Table("user_attempts").
        Select("AVG(questions.difficulty_score)").
        Joins("JOIN questions ON user_attempts.question_id = questions.question_id").
        Where("user_attempts.user_id = ? AND user_attempts.is_correct = true", userID).
        Scan(&avgDifficulty)

    // Recommend slightly harder questions
    targetDifficulty := avgDifficulty + 10.0

    var questions []models.Question
    s.db.Where("difficulty_score BETWEEN ? AND ?", targetDifficulty-5, targetDifficulty+5).
        Where("question_id NOT IN (?)",
            s.db.Table("user_attempts").Select("question_id").Where("user_id = ?", userID),
        ).
        Limit(limit).
        Find(&questions)

    recommendations := []Recommendation{}
    for _, q := range questions {
        recommendations = append(recommendations, Recommendation{
            QuestionID: q.QuestionID,
            ProblemID:  q.ProblemID,
            Reason:     "Progressive challenge",
            Priority:   60.0,
        })
    }

    return recommendations
}

func (s *RecommendationService) GetSpacedRepetitionRecommendations(userID int, limit int) []Recommendation {
    // Find questions due for review
    var overdueSkills []models.UserSkill
    s.db.Where("user_id = ? AND next_review_at < NOW() AND needs_review = true", userID).
        Order("next_review_at ASC").
        Limit(limit).
        Find(&overdueSkills)

    recommendations := []Recommendation{}

    for _, skill := range overdueSkills {
        // Find a question for this topic
        var question models.Question
        s.db.Joins("JOIN problem_topics ON problems.problem_id = problem_topics.problem_id").
            Where("problem_topics.topic_id = ?", skill.TopicID).
            Order("RANDOM()").
            First(&question)

        if question.QuestionID > 0 {
            recommendations = append(recommendations, Recommendation{
                QuestionID: question.QuestionID,
                ProblemID:  question.ProblemID,
                Reason:     "Due for review (spaced repetition)",
                Priority:   50.0,
            })
        }
    }

    return recommendations
}
```

**Estimated Time:** 1 day

---

### 5. Streak Tracking Implementation

**Location:** `backend/services/user_service.go`

```go
package services

import (
    "time"
    "gorm.io/gorm"
    "github.com/yourusername/algoholic/models"
)

type UserService struct {
    db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
    return &UserService{db: db}
}

func (s *UserService) UpdateStreak(userID int) error {
    var user models.User
    if err := s.db.First(&user, userID).Error; err != nil {
        return err
    }

    // Get today's date (midnight)
    today := time.Now().Truncate(24 * time.Hour)
    yesterday := today.AddDate(0, 0, -1)

    // Check if user practiced today
    var todayActivity models.DailyActivity
    err := s.db.Where("user_id = ? AND date = ?", userID, today).First(&todayActivity).Error

    if err == gorm.ErrRecordNotFound {
        // Check if user practiced yesterday
        var yesterdayActivity models.DailyActivity
        err = s.db.Where("user_id = ? AND date = ?", userID, yesterday).First(&yesterdayActivity).Error

        if err == gorm.ErrRecordNotFound {
            // Streak broken - reset to 0
            user.CurrentStreakDays = 0
        } else {
            // Continue streak
            user.CurrentStreakDays++
        }
    }

    return s.db.Save(&user).Error
}

func (s *UserService) RecordDailyActivity(userID int) error {
    today := time.Now().Truncate(24 * time.Hour)

    var activity models.DailyActivity
    err := s.db.Where("user_id = ? AND date = ?", userID, today).First(&activity).Error

    if err == gorm.ErrRecordNotFound {
        // Create new activity record
        activity = models.DailyActivity{
            UserID:         userID,
            Date:           today,
            QuestionsCount: 1,
            StudyTime:      0,
        }

        // Get current streak
        var user models.User
        s.db.First(&user, userID)
        activity.Streak = user.CurrentStreakDays + 1

        return s.db.Create(&activity).Error
    }

    // Update existing activity
    activity.QuestionsCount++
    return s.db.Save(&activity).Error
}

func (s *UserService) GetUserStats(userID int) (map[string]interface{}, error) {
    var user models.User
    if err := s.db.First(&user, userID).Error; err != nil {
        return nil, err
    }

    // Total questions attempted
    var totalAttempts int64
    s.db.Model(&models.UserAttempt{}).Where("user_id = ?", userID).Count(&totalAttempts)

    // Total correct
    var correctAttempts int64
    s.db.Model(&models.UserAttempt{}).Where("user_id = ? AND is_correct = true", userID).Count(&correctAttempts)

    // Accuracy
    var accuracy float64
    if totalAttempts > 0 {
        accuracy = float64(correctAttempts) / float64(totalAttempts) * 100
    }

    // Topics practiced
    var topicsPracticed int64
    s.db.Model(&models.UserSkill{}).Where("user_id = ? AND questions_attempted > 0", userID).Count(&topicsPracticed)

    // Recent activity (last 7 days)
    var recentActivities []models.DailyActivity
    s.db.Where("user_id = ? AND date >= ?", userID, time.Now().AddDate(0, 0, -7)).
        Order("date DESC").
        Find(&recentActivities)

    return map[string]interface{}{
        "current_streak":    user.CurrentStreakDays,
        "total_attempts":    totalAttempts,
        "correct_attempts":  correctAttempts,
        "accuracy":          accuracy,
        "topics_practiced":  topicsPracticed,
        "total_study_time":  user.TotalStudyTime,
        "recent_activities": recentActivities,
    }, nil
}
```

**Estimated Time:** 1 day

---

## Phase 2 Intelligence Layer

### Overview

**Goal:** Add semantic search, graph relationships, and vector-based recommendations

**Duration:** 3-4 weeks

**Key Technologies:**
- ChromaDB (vector database)
- Apache AGE (graph extension for PostgreSQL)
- Sentence Transformers (embedding model)

---

### 1. ChromaDB Integration

#### 1.1 Setup

**Install ChromaDB:**
```bash
pip install chromadb
```

**Start ChromaDB Server:**
```bash
chroma run --path ./chroma_data --port 8000
```

**Docker Compose Addition:**
```yaml
chromadb:
  image: chromadb/chroma:latest
  ports:
    - "8000:8000"
  volumes:
    - chroma_data:/chroma/chroma
  environment:
    IS_PERSISTENT: "TRUE"
    ANONYMIZED_TELEMETRY: "FALSE"
```

#### 1.2 Go Client

**Install:**
```bash
go get github.com/amikos-tech/chroma-go
```

**Implementation:** `backend/services/vector_service.go`

```go
package services

import (
    chroma "github.com/amikos-tech/chroma-go"
    "github.com/amikos-tech/chroma-go/types"
)

type VectorService struct {
    client     *chroma.Client
    collections map[string]*chroma.Collection
}

func NewVectorService(chromaURL string) (*VectorService, error) {
    client, err := chroma.NewClient(chromaURL)
    if err != nil {
        return nil, err
    }

    vs := &VectorService{
        client:      client,
        collections: make(map[string]*chroma.Collection),
    }

    // Initialize collections
    if err := vs.InitializeCollections(); err != nil {
        return nil, err
    }

    return vs, nil
}

func (vs *VectorService) InitializeCollections() error {
    collectionNames := []string{"problems", "questions", "solutions", "templates"}

    for _, name := range collectionNames {
        collection, err := vs.client.GetOrCreateCollection(name, nil)
        if err != nil {
            return err
        }
        vs.collections[name] = collection
    }

    return nil
}

// Add problem to vector DB
func (vs *VectorService) AddProblem(problem *models.Problem, embedding []float32) error {
    collection := vs.collections["problems"]

    // Prepare metadata
    metadata := map[string]interface{}{
        "problem_id":       problem.ProblemID,
        "title":            problem.Title,
        "difficulty_score": problem.DifficultyScore,
        "primary_pattern":  problem.PrimaryPattern,
    }

    // Prepare document (text to be embedded)
    document := fmt.Sprintf("%s\n\n%s", problem.Title, problem.Description)

    _, err := collection.Add(
        []string{embedding},
        []map[string]interface{}{metadata},
        []string{document},
        []string{fmt.Sprintf("problem_%d", problem.ProblemID)},
    )

    return err
}

// Semantic search for similar problems
func (vs *VectorService) SearchSimilarProblems(queryText string, limit int) ([]map[string]interface{}, error) {
    collection := vs.collections["problems"]

    // Query with text (ChromaDB will auto-embed)
    results, err := collection.Query(
        []string{queryText},
        limit,
        nil,
        nil,
        nil,
    )

    if err != nil {
        return nil, err
    }

    // Extract and return results
    similarProblems := []map[string]interface{}{}
    for i, metadata := range results.Metadatas[0] {
        similarProblems = append(similarProblems, map[string]interface{}{
            "problem_id": metadata["problem_id"],
            "title":      metadata["title"],
            "score":      results.Distances[0][i],
        })
    }

    return similarProblems, nil
}
```

**Estimated Time:** 2-3 days

---

### 2. Embedding Generation Pipeline

**Location:** `backend/services/embedding_service.go`

```go
package services

import (
    "bytes"
    "encoding/json"
    "net/http"
)

type EmbeddingService struct {
    ollamaURL string
    model     string
}

func NewEmbeddingService(ollamaURL string) *EmbeddingService {
    return &EmbeddingService{
        ollamaURL: ollamaURL,
        model:     "all-minilm", // or "nomic-embed-text"
    }
}

type EmbeddingRequest struct {
    Model  string `json:"model"`
    Prompt string `json:"prompt"`
}

type EmbeddingResponse struct {
    Embedding []float32 `json:"embedding"`
}

func (es *EmbeddingService) GenerateEmbedding(text string) ([]float32, error) {
    reqBody := EmbeddingRequest{
        Model:  es.model,
        Prompt: text,
    }

    jsonData, err := json.Marshal(reqBody)
    if err != nil {
        return nil, err
    }

    resp, err := http.Post(
        fmt.Sprintf("%s/api/embeddings", es.ollamaURL),
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var embResp EmbeddingResponse
    if err := json.NewDecoder(resp.Body).Decode(&embResp); err != nil {
        return nil, err
    }

    return embResp.Embedding, nil
}

// Batch process all problems
func (es *EmbeddingService) EmbedAllProblems(db *gorm.DB, vectorService *VectorService) error {
    var problems []models.Problem
    db.Find(&problems)

    for _, problem := range problems {
        text := fmt.Sprintf("%s\n\n%s", problem.Title, problem.Description)

        embedding, err := es.GenerateEmbedding(text)
        if err != nil {
            log.Printf("Failed to embed problem %d: %v", problem.ProblemID, err)
            continue
        }

        if err := vectorService.AddProblem(&problem, embedding); err != nil {
            log.Printf("Failed to add problem %d to vector DB: %v", problem.ProblemID, err)
        }
    }

    return nil
}
```

**Setup Ollama for Embeddings:**
```bash
ollama pull all-minilm
# or
ollama pull nomic-embed-text
```

**Estimated Time:** 2 days

---

### 3. Apache AGE Graph Database Setup

#### 3.1 PostgreSQL Extension Setup

**Install AGE Extension:**

Already in PostgreSQL image (apache/age:PG16_latest), just need to enable:

```sql
CREATE EXTENSION IF NOT EXISTS age;
LOAD 'age';
SET search_path = ag_catalog, "$user", public;

-- Create graph
SELECT create_graph('problem_graph');
```

#### 3.2 Graph Schema

**Node Types:**
```cypher
(:Problem {id, title, difficulty_score, primary_pattern})
(:Topic {id, name, category, difficulty_level})
(:User {id, username, proficiency_level})
```

**Relationship Types:**
```cypher
(:Problem)-[:SIMILAR_TO {similarity_score: 0.85}]->(:Problem)
(:Problem)-[:FOLLOW_UP_OF]->(:Problem)
(:Problem)-[:HAS_TOPIC {relevance: 0.9}]->(:Topic)
(:Topic)-[:PREREQUISITE_FOR]->(:Topic)
(:User)-[:MASTERED]->(:Topic)
```

#### 3.3 Graph Service Implementation

**Location:** `backend/services/graph_service.go`

```go
package services

import (
    "database/sql"
    "fmt"
)

type GraphService struct {
    db *sql.DB
}

func NewGraphService(db *sql.DB) *GraphService {
    return &GraphService{db: db}
}

func (gs *GraphService) CreateProblemNode(problemID int, title string, difficulty float64) error {
    query := fmt.Sprintf(`
        SELECT * FROM cypher('problem_graph', $$
            MERGE (p:Problem {id: %d})
            SET p.title = '%s', p.difficulty = %f
        $$) AS (v agtype)
    `, problemID, title, difficulty)

    _, err := gs.db.Exec(query)
    return err
}

func (gs *GraphService) CreateSimilarityRelationship(problem1ID, problem2ID int, score float64) error {
    query := fmt.Sprintf(`
        SELECT * FROM cypher('problem_graph', $$
            MATCH (p1:Problem {id: %d}), (p2:Problem {id: %d})
            MERGE (p1)-[r:SIMILAR_TO {score: %f}]->(p2)
        $$) AS (v agtype)
    `, problem1ID, problem2ID, score)

    _, err := gs.db.Exec(query)
    return err
}

func (gs *GraphService) FindSimilarProblems(problemID int, limit int) ([]map[string]interface{}, error) {
    query := fmt.Sprintf(`
        SELECT * FROM cypher('problem_graph', $$
            MATCH (p:Problem {id: %d})-[r:SIMILAR_TO]-(similar:Problem)
            WHERE r.score > 0.7
            RETURN similar.id, similar.title, r.score
            ORDER BY r.score DESC
            LIMIT %d
        $$) AS (id bigint, title text, score numeric)
    `, problemID, limit)

    rows, err := gs.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    results := []map[string]interface{}{}
    for rows.Next() {
        var id int
        var title string
        var score float64

        if err := rows.Scan(&id, &title, &score); err != nil {
            continue
        }

        results = append(results, map[string]interface{}{
            "problem_id": id,
            "title":      title,
            "similarity": score,
        })
    }

    return results, nil
}

func (gs *GraphService) GetLearningPath(startTopic, endTopic string) ([]string, error) {
    query := fmt.Sprintf(`
        SELECT * FROM cypher('problem_graph', $$
            MATCH path = shortestPath(
                (start:Topic {name: '%s'})-[:PREREQUISITE_FOR*]-(end:Topic {name: '%s'})
            )
            RETURN [node IN nodes(path) | node.name]
        $$) AS (learning_path agtype)
    `, startTopic, endTopic)

    var pathJSON string
    err := gs.db.QueryRow(query).Scan(&pathJSON)
    if err != nil {
        return nil, err
    }

    // Parse JSON array
    var path []string
    json.Unmarshal([]byte(pathJSON), &path)

    return path, nil
}
```

**Estimated Time:** 3-4 days

---

### 4. Semantic Search API Endpoint

**Location:** `backend/handlers/search_handler.go`

```go
package handlers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/yourusername/algoholic/services"
)

type SearchHandler struct {
    vectorService *services.VectorService
    graphService  *services.GraphService
}

func NewSearchHandler(vectorService *services.VectorService, graphService *services.GraphService) *SearchHandler {
    return &SearchHandler{
        vectorService: vectorService,
        graphService:  graphService,
    }
}

func (h *SearchHandler) SemanticSearch(c *fiber.Ctx) error {
    query := c.Query("q")
    if query == "" {
        return c.Status(400).JSON(fiber.Map{
            "error": "Query parameter 'q' is required",
        })
    }

    limit := c.QueryInt("limit", 10)

    // Semantic search
    results, err := h.vectorService.SearchSimilarProblems(query, limit)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Search failed",
        })
    }

    return c.JSON(fiber.Map{
        "results": results,
        "count":   len(results),
    })
}

func (h *SearchHandler) FindSimilar(c *fiber.Ctx) error {
    problemID, err := c.ParamsInt("id")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid problem ID",
        })
    }

    // Combine vector similarity + graph relationships
    vectorResults, _ := h.vectorService.SearchSimilarProblems(
        fmt.Sprintf("problem_%d", problemID),
        5,
    )

    graphResults, _ := h.graphService.FindSimilarProblems(problemID, 5)

    // Merge and deduplicate
    allResults := mergeResults(vectorResults, graphResults)

    return c.JSON(fiber.Map{
        "similar_problems": allResults,
    })
}
```

**Estimated Time:** 2 days

---

## Implementation Timeline

### Week 1: Phase 1 Critical Fixes
- **Day 1-2:** Code execution validation (Judge0 integration)
- **Day 3:** Text fuzzy matching implementation
- **Day 4-5:** Create seed data (50 problems, 200 questions)
- **Day 6:** Difficulty scoring algorithm
- **Day 7:** Recommendation logic

### Week 2: Phase 1 Completion
- **Day 1-2:** Streak tracking implementation
- **Day 3:** Testing and bug fixes
- **Day 4-5:** API endpoint testing with Postman
- **Day 6:** Documentation updates
- **Day 7:** Phase 1 review and completion

### Week 3-4: Phase 2 Intelligence (ChromaDB)
- **Week 3 Day 1-2:** ChromaDB setup and integration
- **Week 3 Day 3-4:** Embedding generation pipeline
- **Week 3 Day 5-7:** Vector search implementation
- **Week 4 Day 1-2:** Batch embed existing data
- **Week 4 Day 3-4:** Semantic search API endpoints
- **Week 4 Day 5-7:** Testing and optimization

### Week 5-6: Phase 2 Intelligence (Graph DB)
- **Week 5 Day 1-2:** Apache AGE setup
- **Week 5 Day 3-5:** Graph schema and seeding
- **Week 5 Day 6-7:** Graph service implementation
- **Week 6 Day 1-3:** Graph query API endpoints
- **Week 6 Day 4-5:** Combined vector + graph search
- **Week 6 Day 6-7:** Testing and documentation

---

## Technical Specifications

### Code Execution Sandbox Requirements

**Judge0 CE Setup (Docker):**
```yaml
judge0:
  image: judge0/judge0:latest
  volumes:
    - ./judge0.conf:/judge0.conf:ro
  ports:
    - "2358:2358"
  privileged: true
  environment:
    - REDIS_HOST=judge0-redis
    - POSTGRES_HOST=judge0-postgres
```

**Supported Languages:**
- Python 3
- C++
- Java
- JavaScript (Node.js)
- Go

### Vector Database Specifications

**Embedding Model:** `all-MiniLM-L6-v2`
- Dimensions: 384
- Performance: ~3,000 sentences/sec
- Quality: Good for semantic search

**Collections:**
| Collection | Documents | Update Frequency |
|-----------|-----------|------------------|
| problems  | ~500     | On create/update |
| questions | ~2000    | On create/update |
| solutions | ~1000    | On create       |

### Graph Database Specifications

**Node Count Estimate:**
- Problems: 500 nodes
- Topics: 100 nodes
- Users: 1000+ nodes (grows)

**Relationship Count:**
- SIMILAR_TO: ~2,500 edges
- HAS_TOPIC: ~1,500 edges
- PREREQUISITE_FOR: ~150 edges

**Query Performance Target:**
- Find similar problems: <100ms
- Learning path: <200ms
- Recommendations: <150ms

---

## Success Criteria

### Phase 1 Complete When:
- [ ] All 2 TODOs resolved in question_service.go
- [ ] Database contains 50+ problems, 200+ questions
- [ ] Difficulty scoring produces reasonable values (10-90 range)
- [ ] Recommendations return relevant questions
- [ ] Streak tracking updates correctly on daily practice
- [ ] All 22 API endpoints pass Postman tests

### Phase 2 Complete When:
- [ ] ChromaDB integrated and operational
- [ ] All problems/questions embedded in vector DB
- [ ] Semantic search returns relevant results (>0.7 similarity)
- [ ] Apache AGE graph operational with all relationships
- [ ] Graph queries return correct learning paths
- [ ] Combined search (keyword + semantic + graph) functional
- [ ] API response time <500ms for all searches

---

## Risk Mitigation

### High-Risk Items

1. **Code Execution Security**
   - Mitigation: Use Judge0 with strict timeouts and resource limits
   - Fallback: Disable code execution, use pattern matching only

2. **Vector DB Performance**
   - Mitigation: Index optimization, caching frequent queries
   - Fallback: Use PostgreSQL full-text search only

3. **Graph Query Complexity**
   - Mitigation: Limit traversal depth, add query timeouts
   - Fallback: Precompute common paths, store in cache

### Medium-Risk Items

1. **Embedding Generation Time**
   - Mitigation: Background job processing, batch operations
   - Solution: Pre-generate embeddings during seed

2. **Data Quality**
   - Mitigation: Manual review of seed data
   - Solution: Start with high-quality LeetCode problems

---

## Next Steps

1. **Immediate (Today):**
   - Mark planning todo as complete
   - Start implementing answer validation (code execution)

2. **This Week:**
   - Complete all Phase 1 critical items
   - Test thoroughly with Postman collection

3. **Next Week:**
   - Begin Phase 2 ChromaDB integration
   - Set up embedding pipeline

4. **Following Weeks:**
   - Complete Phase 2 intelligence layer
   - Prepare for Phase 3 (Training Plans)

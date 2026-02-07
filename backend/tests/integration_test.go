package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/yourusername/algoholic/backend/config"
	"github.com/yourusername/algoholic/backend/models"
	"github.com/yourusername/algoholic/backend/routes"
)

var (
	testDB  *gorm.DB
	testApp *fiber.App
	testCfg *config.Config
)

func setupTestApp(t *testing.T) {
	// Create in-memory SQLite database for testing
	var err error
	testDB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations
	if err := models.AutoMigrate(testDB); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Create test configuration
	testCfg = &config.Config{
		App: config.AppConfig{
			Name:        "Algoholic Test",
			Version:     "1.0.0-test",
			Environment: "test",
			Debug:       true,
		},
		Auth: config.AuthConfig{
			Enabled:    false, // Disable auth for easier testing
			JWTSecret:  "test-secret",
			JWTExpiry:  24,
			BCryptCost: 4, // Lower cost for faster tests
		},
	}

	// Create Fiber app
	testApp = fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Setup routes
	routes.SetupRoutes(testApp, testDB, testCfg)
}

func teardownTestApp(t *testing.T) {
	sqlDB, err := testDB.DB()
	if err == nil {
		sqlDB.Close()
	}
}

// Helper function to create a test request
func makeRequest(method, path string, body interface{}) *http.Request {
	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = json.Marshal(body)
	}
	req := httptest.NewRequest(method, path, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	return req
}

// Test Health Endpoint
func TestHealthEndpoint(t *testing.T) {
	setupTestApp(t)
	defer teardownTestApp(t)

	req := makeRequest("GET", "/health", nil)
	resp, err := testApp.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "healthy", result["status"])
}

// Test User Registration
func TestUserRegistration(t *testing.T) {
	setupTestApp(t)
	defer teardownTestApp(t)

	registerReq := map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "testpassword123",
	}

	req := makeRequest("POST", "/api/auth/register", registerReq)
	resp, err := testApp.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Contains(t, result, "token")
	assert.Contains(t, result, "user")
}

// Test User Login
func TestUserLogin(t *testing.T) {
	setupTestApp(t)
	defer teardownTestApp(t)

	// First register a user
	registerReq := map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "testpassword123",
	}
	testApp.Test(makeRequest("POST", "/api/auth/register", registerReq))

	// Then test login
	loginReq := map[string]string{
		"username": "testuser",
		"password": "testpassword123",
	}

	req := makeRequest("POST", "/api/auth/login", loginReq)
	resp, err := testApp.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Contains(t, result, "token")
}

// Test Problem Retrieval
func TestGetProblems(t *testing.T) {
	setupTestApp(t)
	defer teardownTestApp(t)

	// Create a test problem
	testProblem := &models.Problem{
		Title:           "Test Problem",
		Slug:            "test-problem",
		Description:     "This is a test problem",
		DifficultyScore: 50.0,
		Examples:        models.JSONB{"example": "test"},
	}
	testDB.Create(testProblem)

	req := makeRequest("GET", "/api/problems", nil)
	resp, err := testApp.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Contains(t, result, "problems")
	assert.Contains(t, result, "total")
}

// Test Question Answer Submission
func TestQuestionAnswerSubmission(t *testing.T) {
	setupTestApp(t)
	defer teardownTestApp(t)

	// Create a test user
	testUser := &models.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	testDB.Create(testUser)

	// Create a test question
	correctAnswer := models.JSONB{"answer": "A"}
	testQuestion := &models.Question{
		QuestionType:    "multiple_choice",
		QuestionFormat:  "multiple_choice",
		QuestionText:    "What is 2+2?",
		CorrectAnswer:   correctAnswer,
		Explanation:     "2+2 equals 4",
		DifficultyScore: 10.0,
	}
	testDB.Create(testQuestion)

	// Submit an answer
	answerReq := map[string]interface{}{
		"user_answer":        map[string]string{"answer": "A"},
		"time_taken_seconds": 30,
		"hints_used":         0,
	}

	req := makeRequest("POST", "/api/questions/1/answer", answerReq)
	resp, err := testApp.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Contains(t, result, "is_correct")
	assert.Contains(t, result, "explanation")
}

// Test Training Plan Creation
func TestTrainingPlanCreation(t *testing.T) {
	setupTestApp(t)
	defer teardownTestApp(t)

	// Create a test user
	testUser := &models.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	testDB.Create(testUser)

	// Create some test questions
	for i := 1; i <= 5; i++ {
		question := &models.Question{
			QuestionType:    "multiple_choice",
			QuestionFormat:  "multiple_choice",
			QuestionText:    "Test question",
			CorrectAnswer:   models.JSONB{"answer": "A"},
			Explanation:     "Test explanation",
			DifficultyScore: 50.0,
		}
		testDB.Create(question)
	}

	// Create training plan
	planReq := map[string]interface{}{
		"name":                "Test Plan",
		"description":         "A test training plan",
		"plan_type":           "custom",
		"duration_days":       7,
		"questions_per_day":   2,
		"difficulty_min":      0.0,
		"difficulty_max":      100.0,
		"adaptive_difficulty": true,
	}

	req := makeRequest("POST", "/api/training-plans", planReq)
	resp, err := testApp.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Contains(t, result, "plan")
}

// Test User Stats Retrieval
func TestGetUserStats(t *testing.T) {
	setupTestApp(t)
	defer teardownTestApp(t)

	// Create a test user
	testUser := &models.User{
		Username:          "testuser",
		Email:             "test@example.com",
		PasswordHash:      "hashedpassword",
		CurrentStreakDays: 5,
		TotalStudyTime:    3600,
	}
	testDB.Create(testUser)

	req := makeRequest("GET", "/api/users/me/stats", nil)
	resp, err := testApp.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Contains(t, result, "total_attempts")
	assert.Contains(t, result, "accuracy_rate")
	assert.Contains(t, result, "current_streak_days")
}

// Run all tests
func TestMain(m *testing.M) {
	// Setup code before all tests if needed
	m.Run()
	// Teardown code after all tests if needed
}

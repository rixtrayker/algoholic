package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/yourusername/algoholic/models"
	"gorm.io/gorm"
)

// QuestionService handles question-related operations
type QuestionService struct {
	db *gorm.DB
}

// NewQuestionService creates a new question service
func NewQuestionService(db *gorm.DB) *QuestionService {
	return &QuestionService{db: db}
}

// GetQuestions retrieves questions with filters
func (s *QuestionService) GetQuestions(questionType string, minDifficulty, maxDifficulty float64, limit, offset int) ([]models.Question, int64, error) {
	query := s.db.Model(&models.Question{})

	if questionType != "" {
		query = query.Where("question_type = ?", questionType)
	}

	query = query.Where("difficulty_score BETWEEN ? AND ?", minDifficulty, maxDifficulty)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var questions []models.Question
	if err := query.Limit(limit).Offset(offset).Find(&questions).Error; err != nil {
		return nil, 0, err
	}

	return questions, total, nil
}

// GetQuestionByID retrieves a question by ID
func (s *QuestionService) GetQuestionByID(id int) (*models.Question, error) {
	var question models.Question
	if err := s.db.First(&question, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("question not found")
		}
		return nil, err
	}
	return &question, nil
}

// GetQuestionsByProblem retrieves all questions for a problem
func (s *QuestionService) GetQuestionsByProblem(problemID int) ([]models.Question, error) {
	var questions []models.Question
	err := s.db.Where("problem_id = ?", problemID).Find(&questions).Error
	return questions, err
}

// GetRandomQuestion gets a random question with optional filters
func (s *QuestionService) GetRandomQuestion(questionType string, minDifficulty, maxDifficulty float64, excludeIDs []int) (*models.Question, error) {
	query := s.db.Model(&models.Question{})

	if questionType != "" {
		query = query.Where("question_type = ?", questionType)
	}

	query = query.Where("difficulty_score BETWEEN ? AND ?", minDifficulty, maxDifficulty)

	if len(excludeIDs) > 0 {
		query = query.Where("question_id NOT IN ?", excludeIDs)
	}

	var question models.Question
	if err := query.Order("RANDOM()").First(&question).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no questions found matching criteria")
		}
		return nil, err
	}

	return &question, nil
}

// AnswerRequest represents a question answer submission
type AnswerRequest struct {
	QuestionID     int                    `json:"question_id"`
	UserAnswer     map[string]interface{} `json:"user_answer"`
	TimeTaken      int                    `json:"time_taken_seconds"`
	HintsUsed      int                    `json:"hints_used"`
	Confidence     *int                   `json:"confidence_level,omitempty"`
	TrainingPlanID *int                   `json:"training_plan_id,omitempty"`
}

// AnswerResponse represents the result of answering a question
type AnswerResponse struct {
	IsCorrect              bool                   `json:"is_correct"`
	CorrectAnswer          map[string]interface{} `json:"correct_answer"`
	Explanation            string                 `json:"explanation"`
	WrongAnswerExplanation string                 `json:"wrong_answer_explanation,omitempty"`
	AttemptID              int                    `json:"attempt_id"`
	PointsEarned           int                    `json:"points_earned"`
	NewProficiencyLevel    float64                `json:"new_proficiency_level,omitempty"`
}

// SubmitAnswer processes a question answer
func (s *QuestionService) SubmitAnswer(userID int, req AnswerRequest) (*AnswerResponse, error) {
	// Get question
	question, err := s.GetQuestionByID(req.QuestionID)
	if err != nil {
		return nil, err
	}

	// Check if answer is correct
	isCorrect := s.CheckAnswer(question, req.UserAnswer)

	// Create user attempt record
	userAnswerJSON, _ := json.Marshal(req.UserAnswer)
	var userAnswerMap map[string]interface{}
	json.Unmarshal(userAnswerJSON, &userAnswerMap)

	attempt := models.UserAttempt{
		UserID:           userID,
		QuestionID:       &req.QuestionID,
		UserAnswer:       userAnswerMap,
		IsCorrect:        isCorrect,
		TimeTakenSeconds: req.TimeTaken,
		HintsUsed:        req.HintsUsed,
		ConfidenceLevel:  req.Confidence,
		TrainingPlanID:   req.TrainingPlanID,
	}

	// Get attempt number for this user/question
	var attemptCount int64
	s.db.Model(&models.UserAttempt{}).
		Where("user_id = ? AND question_id = ?", userID, req.QuestionID).
		Count(&attemptCount)
	attempt.AttemptNumber = int(attemptCount) + 1

	if err := s.db.Create(&attempt).Error; err != nil {
		return nil, err
	}

	// Update question stats
	s.UpdateQuestionStats(req.QuestionID, isCorrect, req.TimeTaken)

	// Calculate points
	points := s.CalculatePoints(question, isCorrect, req.TimeTaken, req.HintsUsed)

	// Build response
	response := &AnswerResponse{
		IsCorrect:     isCorrect,
		CorrectAnswer: question.CorrectAnswer,
		Explanation:   question.Explanation,
		AttemptID:     attempt.AttemptID,
		PointsEarned:  points,
	}

	// Add wrong answer explanation if applicable
	if !isCorrect && question.WrongAnswerExplanations != nil {
		if userAnswerStr, ok := req.UserAnswer["answer"].(string); ok {
			if explanation, exists := question.WrongAnswerExplanations[userAnswerStr]; exists {
				response.WrongAnswerExplanation = fmt.Sprintf("%v", explanation)
			}
		}
	}

	return response, nil
}

// CheckAnswer validates if the user's answer is correct
func (s *QuestionService) CheckAnswer(question *models.Question, userAnswer map[string]interface{}) bool {
	switch question.QuestionFormat {
	case "multiple_choice":
		return s.CheckMultipleChoice(question, userAnswer)
	case "code":
		return s.CheckCode(question, userAnswer)
	case "text":
		return s.CheckText(question, userAnswer)
	case "ranking":
		return s.CheckRanking(question, userAnswer)
	default:
		return false
	}
}

// CheckMultipleChoice validates multiple choice answers
func (s *QuestionService) CheckMultipleChoice(question *models.Question, userAnswer map[string]interface{}) bool {
	userSelection, ok := userAnswer["answer"].(string)
	if !ok {
		return false
	}

	correctAnswer, ok := question.CorrectAnswer["answer"].(string)
	if !ok {
		return false
	}

	return userSelection == correctAnswer
}

// CheckCode validates code answers with actual execution
func (s *QuestionService) CheckCode(question *models.Question, userAnswer map[string]interface{}) bool {
	code, ok := userAnswer["code"].(string)
	if !ok || len(code) == 0 {
		return false
	}

	language, _ := userAnswer["language"].(string)
	if language == "" {
		language = "python" // default to Python
	}

	// Get test cases from correct_answer
	testCases, ok := question.CorrectAnswer["test_cases"]
	if !ok {
		// Fallback: if no test cases, just validate code structure
		executor := NewCodeExecutor("")
		return executor.ValidateCode(code, language)
	}

	// Convert test cases to proper format
	var testCaseList []interface{}
	switch v := testCases.(type) {
	case []interface{}:
		testCaseList = v
	default:
		// Invalid format, fallback to structure validation
		executor := NewCodeExecutor("")
		return executor.ValidateCode(code, language)
	}

	// Run code execution tests
	executor := NewCodeExecutor("") // Uses default Judge0 URL
	result, err := executor.RunTests(code, language, testCaseList)

	if err != nil {
		// If execution service is unavailable, fallback to validation
		return executor.ValidateCode(code, language)
	}

	return result.AllPassed
}

// CheckText validates text answers with fuzzy matching
func (s *QuestionService) CheckText(question *models.Question, userAnswer map[string]interface{}) bool {
	userText, ok := userAnswer["answer"].(string)
	if !ok {
		return false
	}

	correctAnswer, ok := question.CorrectAnswer["answer"]
	if !ok {
		return false
	}

	validator := NewTextValidator()

	// Support multiple correct answer formats
	switch v := correctAnswer.(type) {
	case string:
		// Single correct answer
		return validator.FuzzyMatch(userText, v)

	case []interface{}:
		// Multiple acceptable answers
		acceptableAnswers := make([]string, 0, len(v))
		for _, ans := range v {
			if ansStr, ok := ans.(string); ok {
				acceptableAnswers = append(acceptableAnswers, ansStr)
			}
		}
		return validator.MatchMultiple(userText, acceptableAnswers)

	default:
		// Fallback to exact match if format is unexpected
		if str, ok := v.(string); ok {
			return validator.FuzzyMatch(userText, str)
		}
		return false
	}
}

// CheckRanking validates ranking answers
func (s *QuestionService) CheckRanking(question *models.Question, userAnswer map[string]interface{}) bool {
	userRanking, ok := userAnswer["ranking"].([]interface{})
	if !ok {
		return false
	}

	correctRanking, ok := question.CorrectAnswer["ranking"].([]interface{})
	if !ok {
		return false
	}

	if len(userRanking) != len(correctRanking) {
		return false
	}

	for i := range userRanking {
		if fmt.Sprint(userRanking[i]) != fmt.Sprint(correctRanking[i]) {
			return false
		}
	}

	return true
}

// CalculatePoints calculates points earned for an answer
func (s *QuestionService) CalculatePoints(question *models.Question, isCorrect bool, timeTaken, hintsUsed int) int {
	if !isCorrect {
		return 0
	}

	// Base points from difficulty
	basePoints := int(question.DifficultyScore * 10)

	// Time bonus (faster = more points, up to 20% bonus)
	timeBonus := 0
	if question.EstimatedTimeSeconds != nil && *question.EstimatedTimeSeconds > 0 {
		ratio := float64(timeTaken) / float64(*question.EstimatedTimeSeconds)
		if ratio < 1.0 {
			timeBonus = int(float64(basePoints) * (1.0 - ratio) * 0.2)
		}
	}

	// Hint penalty (each hint costs 10% of base points)
	hintPenalty := hintsUsed * basePoints / 10

	points := basePoints + timeBonus - hintPenalty
	if points < 0 {
		points = 0
	}

	return points
}

// UpdateQuestionStats updates question statistics
func (s *QuestionService) UpdateQuestionStats(questionID int, isCorrect bool, timeTaken int) error {
	updates := map[string]interface{}{
		"total_attempts": gorm.Expr("total_attempts + 1"),
	}

	if isCorrect {
		updates["correct_attempts"] = gorm.Expr("correct_attempts + 1")
	}

	// Update average time (rolling average)
	var question models.Question
	if err := s.db.First(&question, questionID).Error; err == nil {
		if question.AverageTimeSeconds == nil {
			avgTime := float64(timeTaken)
			updates["average_time_seconds"] = avgTime
		} else {
			// Rolling average: new_avg = (old_avg * count + new_value) / (count + 1)
			newAvg := (*question.AverageTimeSeconds*float64(question.TotalAttempts) + float64(timeTaken)) / float64(question.TotalAttempts+1)
			updates["average_time_seconds"] = newAvg
		}
	}

	return s.db.Model(&models.Question{}).
		Where("question_id = ?", questionID).
		Updates(updates).Error
}

// GetUserAttempts retrieves a user's attempts for a question
func (s *QuestionService) GetUserAttempts(userID, questionID int) ([]models.UserAttempt, error) {
	var attempts []models.UserAttempt
	err := s.db.Where("user_id = ? AND question_id = ?", userID, questionID).
		Order("attempted_at DESC").
		Find(&attempts).Error
	return attempts, err
}

// GetRecentAttempts gets a user's recent attempts
func (s *QuestionService) GetRecentAttempts(userID int, limit int) ([]models.UserAttempt, error) {
	var attempts []models.UserAttempt
	err := s.db.Where("user_id = ?", userID).
		Order("attempted_at DESC").
		Limit(limit).
		Find(&attempts).Error
	return attempts, err
}

// GetHint retrieves a hint for a question
func (s *QuestionService) GetHint(questionID int, hintLevel int) (string, error) {
	var question models.QuestionWithHints
	if err := s.db.Table("questions").Where("question_id = ?", questionID).First(&question).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("question not found")
		}
		return "", err
	}

	var hint *string
	switch hintLevel {
	case 1:
		hint = question.HintLevel1
	case 2:
		hint = question.HintLevel2
	case 3:
		hint = question.HintLevel3
	default:
		return "", errors.New("invalid hint level (must be 1-3)")
	}

	if hint == nil || *hint == "" {
		return "", errors.New("no hint available at this level")
	}

	return *hint, nil
}

// RecordHintUsage records that a user used a hint
func (s *QuestionService) RecordHintUsage(userID, questionID, hintLevel int) error {
	usage := models.QuestionHintUsage{
		UserID:     userID,
		QuestionID: questionID,
		HintLevel:  hintLevel,
	}

	if err := s.db.Create(&usage).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil
		}
		return err
	}

	return nil
}

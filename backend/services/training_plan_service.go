package services

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"github.com/yourusername/algoholic/backend/models"
)

// TrainingPlanService handles training plan operations
type TrainingPlanService struct {
	db              *gorm.DB
	questionService *QuestionService
	userService     *UserService
}

// NewTrainingPlanService creates a new training plan service
func NewTrainingPlanService(db *gorm.DB, questionService *QuestionService, userService *UserService) *TrainingPlanService {
	return &TrainingPlanService{
		db:              db,
		questionService: questionService,
		userService:     userService,
	}
}

// CreatePlanRequest represents a training plan creation request
type CreatePlanRequest struct {
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	PlanType           string   `json:"plan_type"` // preset, custom, ai_generated
	TargetTopics       []int    `json:"target_topics"`
	TargetPatterns     []string `json:"target_patterns"`
	DurationDays       int      `json:"duration_days"`
	QuestionsPerDay    int      `json:"questions_per_day"`
	DifficultyMin      float64  `json:"difficulty_min"`
	DifficultyMax      float64  `json:"difficulty_max"`
	AdaptiveDifficulty bool     `json:"adaptive_difficulty"`
}

// CreateTrainingPlan creates a new training plan
func (s *TrainingPlanService) CreateTrainingPlan(userID int, req CreatePlanRequest) (*models.TrainingPlan, error) {
	// Create training plan
	plan := &models.TrainingPlan{
		UserID:             userID,
		Name:               req.Name,
		Description:        &req.Description,
		PlanType:           &req.PlanType,
		DurationDays:       &req.DurationDays,
		QuestionsPerDay:    req.QuestionsPerDay,
		AdaptiveDifficulty: req.AdaptiveDifficulty,
		Status:             "active",
		StartDate:          time.Now(),
	}

	// Convert target topics and patterns
	if len(req.TargetTopics) > 0 {
		topicsArray := make(models.StringArray, len(req.TargetTopics))
		for i, topic := range req.TargetTopics {
			topicsArray[i] = string(rune(topic))
		}
		plan.TargetTopics = topicsArray
	}

	if len(req.TargetPatterns) > 0 {
		plan.TargetPatterns = models.StringArray(req.TargetPatterns)
	}

	// Save plan
	if err := s.db.Create(plan).Error; err != nil {
		return nil, err
	}

	// Generate plan items
	if err := s.GeneratePlanItems(plan, req.DifficultyMin, req.DifficultyMax); err != nil {
		return nil, err
	}

	return plan, nil
}

// GeneratePlanItems generates questions for a training plan
func (s *TrainingPlanService) GeneratePlanItems(plan *models.TrainingPlan, minDiff, maxDiff float64) error {
	totalQuestions := plan.QuestionsPerDay * *plan.DurationDays
	sequenceNumber := 1

	// Get questions based on target topics/patterns
	questions, _, err := s.questionService.GetQuestions("", minDiff, maxDiff, totalQuestions, 0)
	if err != nil {
		return err
	}

	// Distribute questions across days
	for day := 1; day <= *plan.DurationDays; day++ {
		questionsForDay := plan.QuestionsPerDay
		scheduledDate := plan.StartDate.AddDate(0, 0, day-1)

		for i := 0; i < questionsForDay && len(questions) > 0; i++ {
			question := questions[0]
			questions = questions[1:]

			item := models.TrainingPlanItem{
				PlanID:         plan.PlanID,
				QuestionID:     &question.QuestionID,
				SequenceNumber: sequenceNumber,
				DayNumber:      &day,
				ScheduledFor:   &scheduledDate,
				ItemType:       "question",
				IsCompleted:    false,
			}

			if err := s.db.Create(&item).Error; err != nil {
				return err
			}

			sequenceNumber++
		}
	}

	return nil
}

// GetUserPlans retrieves all training plans for a user
func (s *TrainingPlanService) GetUserPlans(userID int) ([]models.TrainingPlan, error) {
	var plans []models.TrainingPlan
	err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&plans).Error
	return plans, err
}

// GetPlanByID retrieves a training plan by ID
func (s *TrainingPlanService) GetPlanByID(planID, userID int) (*models.TrainingPlan, error) {
	var plan models.TrainingPlan
	err := s.db.Where("plan_id = ? AND user_id = ?", planID, userID).First(&plan).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("training plan not found")
		}
		return nil, err
	}
	return &plan, nil
}

// GetNextQuestion gets the next question in a training plan
func (s *TrainingPlanService) GetNextQuestion(planID, userID int) (*models.Question, error) {
	// Verify plan belongs to user
	plan, err := s.GetPlanByID(planID, userID)
	if err != nil {
		return nil, err
	}

	// Check if plan is active
	if plan.Status != "active" {
		return nil, errors.New("training plan is not active")
	}

	// Get next incomplete item
	var item models.TrainingPlanItem
	err = s.db.Where("plan_id = ? AND is_completed = FALSE", planID).
		Order("sequence_number ASC").
		First(&item).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Mark plan as completed
			s.db.Model(plan).Updates(map[string]interface{}{
				"status":              "completed",
				"progress_percentage": 100.0,
			})
			return nil, errors.New("training plan completed")
		}
		return nil, err
	}

	// Get the question
	if item.QuestionID == nil {
		return nil, errors.New("invalid plan item")
	}

	question, err := s.questionService.GetQuestionByID(*item.QuestionID)
	if err != nil {
		return nil, err
	}

	return question, nil
}

// CompleteItem marks a training plan item as completed
func (s *TrainingPlanService) CompleteItem(planID, userID, itemID int) error {
	// Verify plan belongs to user
	if _, err := s.GetPlanByID(planID, userID); err != nil {
		return err
	}

	// Mark item as completed
	now := time.Now()
	if err := s.db.Model(&models.TrainingPlanItem{}).
		Where("item_id = ? AND plan_id = ?", itemID, planID).
		Updates(map[string]interface{}{
			"is_completed": true,
			"completed_at": now,
		}).Error; err != nil {
		return err
	}

	// Update plan progress
	return s.UpdatePlanProgress(planID)
}

// UpdatePlanProgress updates training plan completion percentage
func (s *TrainingPlanService) UpdatePlanProgress(planID int) error {
	var total, completed int64

	s.db.Model(&models.TrainingPlanItem{}).Where("plan_id = ?", planID).Count(&total)
	s.db.Model(&models.TrainingPlanItem{}).Where("plan_id = ? AND is_completed = TRUE", planID).Count(&completed)

	if total == 0 {
		return nil
	}

	progress := float64(completed) / float64(total) * 100

	return s.db.Model(&models.TrainingPlan{}).
		Where("plan_id = ?", planID).
		Update("progress_percentage", progress).Error
}

// GetPlanItems retrieves all items in a training plan
func (s *TrainingPlanService) GetPlanItems(planID, userID int) ([]models.TrainingPlanItem, error) {
	// Verify plan belongs to user
	if _, err := s.GetPlanByID(planID, userID); err != nil {
		return nil, err
	}

	var items []models.TrainingPlanItem
	err := s.db.Where("plan_id = ?", planID).
		Order("sequence_number ASC").
		Find(&items).Error

	return items, err
}

// GetTodaysQuestions gets questions scheduled for today
func (s *TrainingPlanService) GetTodaysQuestions(planID, userID int) ([]models.Question, error) {
	// Verify plan belongs to user
	if _, err := s.GetPlanByID(planID, userID); err != nil {
		return nil, err
	}

	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var items []models.TrainingPlanItem
	err := s.db.Where("plan_id = ? AND scheduled_for >= ? AND scheduled_for < ? AND is_completed = FALSE",
		planID, startOfDay, endOfDay).
		Order("sequence_number ASC").
		Find(&items).Error

	if err != nil {
		return nil, err
	}

	// Get questions
	var questions []models.Question
	for _, item := range items {
		if item.QuestionID != nil {
			question, err := s.questionService.GetQuestionByID(*item.QuestionID)
			if err == nil {
				questions = append(questions, *question)
			}
		}
	}

	return questions, nil
}

// AdaptPlanDifficulty adjusts plan difficulty based on user performance
func (s *TrainingPlanService) AdaptPlanDifficulty(planID, userID int) error {
	plan, err := s.GetPlanByID(planID, userID)
	if err != nil {
		return err
	}

	if !plan.AdaptiveDifficulty {
		return nil
	}

	// Get user's recent accuracy
	recentAttempts, err := s.questionService.GetRecentAttempts(userID, 20)
	if err != nil || len(recentAttempts) < 5 {
		return nil // Not enough data
	}

	correctCount := 0
	for _, attempt := range recentAttempts {
		if attempt.IsCorrect {
			correctCount++
		}
	}

	accuracy := float64(correctCount) / float64(len(recentAttempts))

	// Adjust difficulty for upcoming questions
	if accuracy > 0.85 {
		// User is doing well, increase difficulty
		s.db.Exec(`
			UPDATE training_plan_items
			SET question_id = (
				SELECT question_id FROM questions
				WHERE difficulty_score BETWEEN difficulty_score + 5 AND difficulty_score + 15
				ORDER BY RANDOM() LIMIT 1
			)
			WHERE plan_id = ? AND is_completed = FALSE
		`, planID)
	} else if accuracy < 0.40 {
		// User is struggling, decrease difficulty
		s.db.Exec(`
			UPDATE training_plan_items
			SET question_id = (
				SELECT question_id FROM questions
				WHERE difficulty_score BETWEEN difficulty_score - 15 AND difficulty_score - 5
				ORDER BY RANDOM() LIMIT 1
			)
			WHERE plan_id = ? AND is_completed = FALSE
		`, planID)
	}

	return nil
}

// PausePlan pauses a training plan
func (s *TrainingPlanService) PausePlan(planID, userID int) error {
	if _, err := s.GetPlanByID(planID, userID); err != nil {
		return err
	}

	return s.db.Model(&models.TrainingPlan{}).
		Where("plan_id = ?", planID).
		Update("status", "paused").Error
}

// ResumePlan resumes a paused training plan
func (s *TrainingPlanService) ResumePlan(planID, userID int) error {
	if _, err := s.GetPlanByID(planID, userID); err != nil {
		return err
	}

	return s.db.Model(&models.TrainingPlan{}).
		Where("plan_id = ?", planID).
		Update("status", "active").Error
}

// DeletePlan deletes a training plan
func (s *TrainingPlanService) DeletePlan(planID, userID int) error {
	if _, err := s.GetPlanByID(planID, userID); err != nil {
		return err
	}

	// Delete plan items first (cascade)
	if err := s.db.Where("plan_id = ?", planID).Delete(&models.TrainingPlanItem{}).Error; err != nil {
		return err
	}

	// Delete plan
	return s.db.Delete(&models.TrainingPlan{}, planID).Error
}

package services

import (
	"errors"
	"math"
	"strconv"
	"time"

	"gorm.io/gorm"
	"github.com/yourusername/algoholic/models"
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
	Name               string   `json:"name" validate:"required,min=1,max=200"`
	Description        string   `json:"description"`
	PlanType           string   `json:"plan_type" validate:"oneof=preset custom ai_generated"` // preset, custom, ai_generated
	TargetTopics       []int    `json:"target_topics"`
	TargetPatterns     []string `json:"target_patterns"`
	DurationDays       int      `json:"duration_days" validate:"required,gte=1,lte=365"`
	QuestionsPerDay    int      `json:"questions_per_day" validate:"required,gte=1,lte=50"`
	DifficultyMin      float64  `json:"difficulty_min" validate:"gte=0,lte=100"`
	DifficultyMax      float64  `json:"difficulty_max" validate:"gte=0,lte=100"`
	AdaptiveDifficulty bool     `json:"adaptive_difficulty"`
}

// CreateTrainingPlan creates a new training plan within a transaction
func (s *TrainingPlanService) CreateTrainingPlan(userID int, req CreatePlanRequest) (*models.TrainingPlan, error) {
	var plan *models.TrainingPlan

	err := s.db.Transaction(func(tx *gorm.DB) error {
		plan = &models.TrainingPlan{
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
				topicsArray[i] = strconv.Itoa(topic)
			}
			plan.TargetTopics = topicsArray
		}

		if len(req.TargetPatterns) > 0 {
			plan.TargetPatterns = models.StringArray(req.TargetPatterns)
		}

		if err := tx.Create(plan).Error; err != nil {
			return err
		}

		// Generate plan items within same transaction
		if err := s.generatePlanItemsTx(tx, plan, req.DifficultyMin, req.DifficultyMax); err != nil {
			return err
		}

		return nil
	})

	return plan, err
}

// generatePlanItemsTx generates questions for a training plan within a transaction
func (s *TrainingPlanService) generatePlanItemsTx(tx *gorm.DB, plan *models.TrainingPlan, minDiff, maxDiff float64) error {
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

			if err := tx.Create(&item).Error; err != nil {
				return err
			}

			sequenceNumber++
		}
	}

	return nil
}

// GetUserPlans retrieves paginated training plans for a user
func (s *TrainingPlanService) GetUserPlans(userID int, limit, offset int) ([]models.TrainingPlan, int64, error) {
	var total int64
	s.db.Model(&models.TrainingPlan{}).Where("user_id = ?", userID).Count(&total)

	var plans []models.TrainingPlan
	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&plans).Error
	return plans, total, err
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

	// Get current average difficulty of incomplete plan items
	var avgDifficulty float64
	s.db.Table("training_plan_items tpi").
		Joins("JOIN questions q ON q.question_id = tpi.question_id").
		Where("tpi.plan_id = ? AND tpi.is_completed = FALSE", planID).
		Select("COALESCE(AVG(q.difficulty_score), 50)").
		Scan(&avgDifficulty)

	// Calculate target difficulty range
	var targetMin, targetMax float64
	if accuracy > 0.85 {
		// User is doing well, increase difficulty
		targetMin = avgDifficulty + 5
		targetMax = math.Min(100, avgDifficulty+20)
	} else if accuracy < 0.40 {
		// User is struggling, decrease difficulty
		targetMin = math.Max(0, avgDifficulty-20)
		targetMax = math.Max(0, avgDifficulty-5)
	} else {
		return nil // Accuracy is in acceptable range, no adjustment needed
	}

	// Get incomplete item IDs
	var incompleteItems []models.TrainingPlanItem
	s.db.Where("plan_id = ? AND is_completed = FALSE", planID).Find(&incompleteItems)

	// Find replacement questions in the target difficulty range
	var replacementIDs []int
	s.db.Table("questions").
		Select("question_id").
		Where("difficulty_score BETWEEN ? AND ?", targetMin, targetMax).
		Where("question_id NOT IN (?)",
			s.db.Table("training_plan_items").Select("question_id").Where("plan_id = ? AND question_id IS NOT NULL", planID)).
		Order("RANDOM()").
		Limit(len(incompleteItems)).
		Pluck("question_id", &replacementIDs)

	// Replace each incomplete item with a new question
	for i, item := range incompleteItems {
		if i < len(replacementIDs) {
			s.db.Model(&models.TrainingPlanItem{}).
				Where("item_id = ?", item.ItemID).
				Update("question_id", replacementIDs[i])
		}
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

// DeletePlan deletes a training plan and its items within a transaction
func (s *TrainingPlanService) DeletePlan(planID, userID int) error {
	if _, err := s.GetPlanByID(planID, userID); err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("plan_id = ?", planID).Delete(&models.TrainingPlanItem{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.TrainingPlan{}, planID).Error
	})
}

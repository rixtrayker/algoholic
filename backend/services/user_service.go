package services

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"github.com/yourusername/algoholic/models"
)

// UserService handles user-related operations
type UserService struct {
	db *gorm.DB
}

// NewUserService creates a new user service
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// UserStats represents user statistics
type UserStats struct {
	TotalAttempts      int     `json:"total_attempts"`
	CorrectAttempts    int     `json:"correct_attempts"`
	AccuracyRate       float64 `json:"accuracy_rate"`
	TotalStudyTime     int64   `json:"total_study_time_seconds"`
	CurrentStreak      int     `json:"current_streak_days"`
	ProblemsAttempted  int     `json:"problems_attempted"`
	ProblemsSolved     int     `json:"problems_solved"`
	QuestionsAnswered  int     `json:"questions_answered"`
	AverageDifficulty  float64 `json:"average_difficulty"`
	StrongTopics       []string `json:"strong_topics"`
	WeakTopics         []string `json:"weak_topics"`
}

// GetUserStats retrieves comprehensive user statistics
func (s *UserService) GetUserStats(userID int) (*UserStats, error) {
	stats := &UserStats{}

	// Get user for streak and study time
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	stats.CurrentStreak = user.CurrentStreakDays
	stats.TotalStudyTime = user.TotalStudyTime

	// Get attempt statistics
	var attempts []models.UserAttempt
	if err := s.db.Where("user_id = ?", userID).Find(&attempts).Error; err != nil {
		return nil, err
	}

	stats.TotalAttempts = len(attempts)
	correctCount := 0
	for _, attempt := range attempts {
		if attempt.IsCorrect {
			correctCount++
		}
	}
	stats.CorrectAttempts = correctCount
	if stats.TotalAttempts > 0 {
		stats.AccuracyRate = float64(correctCount) / float64(stats.TotalAttempts) * 100
	}

	// Get problem statistics
	var problemAttempts []struct {
		ProblemID int
		Solved    bool
	}
	s.db.Table("user_attempts").
		Select("DISTINCT problem_id, MAX(is_correct) as solved").
		Where("user_id = ? AND problem_id IS NOT NULL", userID).
		Group("problem_id").
		Scan(&problemAttempts)

	stats.ProblemsAttempted = len(problemAttempts)
	for _, p := range problemAttempts {
		if p.Solved {
			stats.ProblemsSolved++
		}
	}

	// Get question count
	s.db.Table("user_attempts").
		Select("COUNT(DISTINCT question_id)").
		Where("user_id = ? AND question_id IS NOT NULL", userID).
		Scan(&stats.QuestionsAnswered)

	// Get strong and weak topics
	strongTopics, _ := s.GetStrongTopics(userID, 5)
	weakTopics, _ := s.GetWeakTopics(userID, 5)

	for _, topic := range strongTopics {
		stats.StrongTopics = append(stats.StrongTopics, topic.Name)
	}
	for _, topic := range weakTopics {
		stats.WeakTopics = append(stats.WeakTopics, topic.Name)
	}

	return stats, nil
}

// UpdateUserProgress updates user progress after an attempt
func (s *UserService) UpdateUserProgress(userID, topicID int, isCorrect bool, timeTaken int) error {
	// Get or create user skill for this topic
	var skill models.UserSkill
	result := s.db.Where("user_id = ? AND topic_id = ?", userID, topicID).First(&skill)

	now := time.Now()

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Create new skill
			skill = models.UserSkill{
				UserID:             userID,
				TopicID:            topicID,
				ProficiencyLevel:   0,
				QuestionsAttempted: 0,
				QuestionsCorrect:   0,
				LastPracticedAt:    &now,
			}
		} else {
			return result.Error
		}
	}

	// Update stats
	skill.QuestionsAttempted++
	if isCorrect {
		skill.QuestionsCorrect++
	}
	skill.LastPracticedAt = &now

	// Calculate new proficiency level (simple formula: correct_rate * 100)
	if skill.QuestionsAttempted > 0 {
		oldProficiency := skill.ProficiencyLevel
		skill.ProficiencyLevel = float64(skill.QuestionsCorrect) / float64(skill.QuestionsAttempted) * 100

		// Calculate improvement rate
		if oldProficiency > 0 {
			improvement := (skill.ProficiencyLevel - oldProficiency) / oldProficiency * 100
			skill.ImprovementRate = &improvement
		}
	}

	// Determine if review is needed (proficiency < 70%)
	skill.NeedsReview = skill.ProficiencyLevel < 70.0

	// Set next review date using spaced repetition
	nextReview := s.CalculateNextReviewDate(skill.ProficiencyLevel, isCorrect)
	skill.NextReviewAt = &nextReview

	// Save or update
	if skill.QuestionsAttempted == 1 {
		return s.db.Create(&skill).Error
	}
	return s.db.Save(&skill).Error
}

// CalculateNextReviewDate calculates when to review this topic again
func (s *UserService) CalculateNextReviewDate(proficiency float64, wasCorrect bool) time.Time {
	// Simple spaced repetition algorithm
	var days int

	if !wasCorrect {
		days = 1 // Review tomorrow if incorrect
	} else if proficiency < 50 {
		days = 2
	} else if proficiency < 70 {
		days = 5
	} else if proficiency < 85 {
		days = 10
	} else {
		days = 20
	}

	return time.Now().AddDate(0, 0, days)
}

// GetStrongTopics retrieves user's strongest topics
func (s *UserService) GetStrongTopics(userID int, limit int) ([]models.Topic, error) {
	var topics []models.Topic
	err := s.db.Table("topics").
		Joins("JOIN user_skills ON user_skills.topic_id = topics.topic_id").
		Where("user_skills.user_id = ? AND user_skills.proficiency_level >= 70", userID).
		Order("user_skills.proficiency_level DESC").
		Limit(limit).
		Find(&topics).Error
	return topics, err
}

// GetWeakTopics retrieves user's weakest topics
func (s *UserService) GetWeakTopics(userID int, limit int) ([]models.Topic, error) {
	var topics []models.Topic
	err := s.db.Table("topics").
		Joins("JOIN user_skills ON user_skills.topic_id = topics.topic_id").
		Where("user_skills.user_id = ? AND user_skills.proficiency_level < 50", userID).
		Order("user_skills.proficiency_level ASC").
		Limit(limit).
		Find(&topics).Error
	return topics, err
}

// GetReviewQueue gets topics that need review
func (s *UserService) GetReviewQueue(userID int) ([]models.UserSkill, error) {
	var skills []models.UserSkill
	now := time.Now()
	err := s.db.Where("user_id = ? AND needs_review = TRUE AND next_review_at <= ?", userID, now).
		Order("next_review_at ASC").
		Find(&skills).Error
	return skills, err
}

// UpdateStreak updates the user's practice streak
func (s *UserService) UpdateStreak(userID int) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	now := time.Now()
	lastActive := user.LastActiveAt

	// Check if last active was yesterday
	yesterday := now.AddDate(0, 0, -1)
	if lastActive.Year() == yesterday.Year() &&
		lastActive.Month() == yesterday.Month() &&
		lastActive.Day() == yesterday.Day() {
		// Increment streak
		user.CurrentStreakDays++
	} else if lastActive.Year() != now.Year() ||
		lastActive.Month() != now.Month() ||
		lastActive.Day() != now.Day() {
		// Reset streak if not today or yesterday
		user.CurrentStreakDays = 1
	}

	user.LastActiveAt = now
	return s.db.Save(&user).Error
}

// AddStudyTime adds study time to user's total
func (s *UserService) AddStudyTime(userID int, seconds int64) error {
	return s.db.Model(&models.User{}).
		Where("user_id = ?", userID).
		Update("total_study_time_seconds", gorm.Expr("total_study_time_seconds + ?", seconds)).
		Error
}

// GetUserPreferences retrieves user preferences
func (s *UserService) GetUserPreferences(userID int) (models.JSONB, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return user.Preferences, nil
}

// UpdateUserPreferences updates user preferences
func (s *UserService) UpdateUserPreferences(userID int, preferences models.JSONB) error {
	return s.db.Model(&models.User{}).
		Where("user_id = ?", userID).
		Update("preferences", preferences).Error
}

// GetUserProgress gets detailed progress for a specific topic
func (s *UserService) GetUserProgress(userID, topicID int) (*models.UserSkill, error) {
	var skill models.UserSkill
	err := s.db.Where("user_id = ? AND topic_id = ?", userID, topicID).First(&skill).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no progress found for this topic")
		}
		return nil, err
	}
	return &skill, nil
}

// GetUserSkills retrieves all user skills
func (s *UserService) GetUserSkills(userID int) ([]models.UserSkill, error) {
	var skills []models.UserSkill
	err := s.db.Where("user_id = ?", userID).
		Order("proficiency_level DESC").
		Find(&skills).Error
	return skills, err
}

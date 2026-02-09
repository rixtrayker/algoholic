package services

import (
	"time"

	"github.com/yourusername/algoholic/models"
	"gorm.io/gorm"
)

type ActivityService struct {
	db *gorm.DB
}

func NewActivityService(db *gorm.DB) *ActivityService {
	return &ActivityService{db: db}
}

// ActivityData represents a single day's activity
type ActivityData struct {
	Date           string `json:"date"`
	ProblemsCount  int    `json:"problems_count"`
	QuestionsCount int    `json:"questions_count"`
	StudyTime      int    `json:"study_time_seconds"`
	Streak         int    `json:"streak"`
}

// ActivityStats represents aggregated statistics
type ActivityStats struct {
	TotalDays          int     `json:"total_days"`
	TotalProblems      int     `json:"total_problems"`
	TotalQuestions     int     `json:"total_questions"`
	TotalStudyTime     int     `json:"total_study_time_seconds"`
	CurrentStreak      int     `json:"current_streak"`
	LongestStreak      int     `json:"longest_streak"`
	AveragePerDay      float64 `json:"average_per_day"`
	MostProductiveDay  string  `json:"most_productive_day"`
	MostProductiveDate string  `json:"most_productive_date"`
}

// PracticeHistoryItem represents a detailed history entry
type PracticeHistoryItem struct {
	Date           string `json:"date"`
	ProblemsCount  int    `json:"problems_count"`
	QuestionsCount int    `json:"questions_count"`
	StudyTime      int    `json:"study_time_seconds"`
	TotalAttempts  int    `json:"total_attempts"`
	CorrectAttempts int   `json:"correct_attempts"`
	AccuracyRate   float64 `json:"accuracy_rate"`
}

// GetActivityData returns activity data for the commitment chart
func (s *ActivityService) GetActivityData(userID, days int) ([]ActivityData, error) {
	var activities []models.DailyActivity

	startDate := time.Now().AddDate(0, 0, -days)
	err := s.db.Where("user_id = ? AND date >= ?", userID, startDate).
		Order("date ASC").
		Find(&activities).Error

	if err != nil {
		return nil, err
	}

	result := make([]ActivityData, 0, len(activities))
	for _, activity := range activities {
		result = append(result, ActivityData{
			Date:           activity.Date.Format("2006-01-02"),
			ProblemsCount:  activity.ProblemsCount,
			QuestionsCount: activity.QuestionsCount,
			StudyTime:      activity.StudyTime,
			Streak:         activity.Streak,
		})
	}

	return result, nil
}

// GetActivityStats returns aggregated activity statistics
func (s *ActivityService) GetActivityStats(userID int) (*ActivityStats, error) {
	var activities []models.DailyActivity

	err := s.db.Where("user_id = ?", userID).
		Order("date DESC").
		Find(&activities).Error

	if err != nil {
		return nil, err
	}

	stats := &ActivityStats{}
	stats.TotalDays = len(activities)

	if len(activities) == 0 {
		return stats, nil
	}

	// Calculate totals
	maxActivity := 0
	var mostProductiveDate time.Time
	for _, activity := range activities {
		stats.TotalProblems += activity.ProblemsCount
		stats.TotalQuestions += activity.QuestionsCount
		stats.TotalStudyTime += activity.StudyTime

		total := activity.ProblemsCount + activity.QuestionsCount
		if total > maxActivity {
			maxActivity = total
			mostProductiveDate = activity.Date
		}
	}

	// Calculate current streak (from most recent activity)
	stats.CurrentStreak = activities[0].Streak

	// Calculate longest streak
	for _, activity := range activities {
		if activity.Streak > stats.LongestStreak {
			stats.LongestStreak = activity.Streak
		}
	}

	// Calculate average per day
	totalActivity := stats.TotalProblems + stats.TotalQuestions
	if stats.TotalDays > 0 {
		stats.AveragePerDay = float64(totalActivity) / float64(stats.TotalDays)
	}

	// Set most productive day
	if !mostProductiveDate.IsZero() {
		stats.MostProductiveDay = mostProductiveDate.Weekday().String()
		stats.MostProductiveDate = mostProductiveDate.Format("2006-01-02")
	}

	return stats, nil
}

// GetPracticeHistory returns detailed practice history with attempt data
func (s *ActivityService) GetPracticeHistory(userID, days int) ([]PracticeHistoryItem, error) {
	startDate := time.Now().AddDate(0, 0, -days)

	// Query to join activity data with attempt data
	var results []struct {
		Date            time.Time
		ProblemsCount   int
		QuestionsCount  int
		StudyTime       int
		TotalAttempts   int64
		CorrectAttempts int64
	}

	err := s.db.Table("daily_activities").
		Select(`
			daily_activities.date,
			daily_activities.problems_count,
			daily_activities.questions_count,
			daily_activities.study_time_seconds as study_time,
			COUNT(attempts.attempt_id) as total_attempts,
			SUM(CASE WHEN attempts.is_correct THEN 1 ELSE 0 END) as correct_attempts
		`).
		Joins("LEFT JOIN attempts ON DATE(attempts.attempt_timestamp) = daily_activities.date AND attempts.user_id = daily_activities.user_id").
		Where("daily_activities.user_id = ? AND daily_activities.date >= ?", userID, startDate).
		Group("daily_activities.date, daily_activities.problems_count, daily_activities.questions_count, daily_activities.study_time_seconds").
		Order("daily_activities.date DESC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	history := make([]PracticeHistoryItem, 0, len(results))
	for _, result := range results {
		accuracyRate := 0.0
		if result.TotalAttempts > 0 {
			accuracyRate = float64(result.CorrectAttempts) / float64(result.TotalAttempts) * 100
		}

		history = append(history, PracticeHistoryItem{
			Date:            result.Date.Format("2006-01-02"),
			ProblemsCount:   result.ProblemsCount,
			QuestionsCount:  result.QuestionsCount,
			StudyTime:       result.StudyTime,
			TotalAttempts:   int(result.TotalAttempts),
			CorrectAttempts: int(result.CorrectAttempts),
			AccuracyRate:    accuracyRate,
		})
	}

	return history, nil
}

// RecordActivity records or updates activity for today
func (s *ActivityService) RecordActivity(userID, problemsCount, questionsCount, studyTime int) error {
	today := time.Now().Truncate(24 * time.Hour)

	// Check if activity record exists for today
	var activity models.DailyActivity
	err := s.db.Where("user_id = ? AND date = ?", userID, today).First(&activity).Error

	if err == gorm.ErrRecordNotFound {
		// Create new activity record
		// Calculate streak
		streak := s.calculateStreak(userID, today)

		activity = models.DailyActivity{
			UserID:         userID,
			Date:           today,
			ProblemsCount:  problemsCount,
			QuestionsCount: questionsCount,
			StudyTime:      studyTime,
			Streak:         streak,
		}
		return s.db.Create(&activity).Error
	} else if err != nil {
		return err
	}

	// Update existing record
	updates := map[string]interface{}{
		"problems_count":       activity.ProblemsCount + problemsCount,
		"questions_count":      activity.QuestionsCount + questionsCount,
		"study_time_seconds":   activity.StudyTime + studyTime,
	}

	return s.db.Model(&activity).Updates(updates).Error
}

// calculateStreak calculates the current streak for a user
func (s *ActivityService) calculateStreak(userID int, currentDate time.Time) int {
	var activities []models.DailyActivity

	// Get all activities before current date, ordered by date descending
	yesterday := currentDate.AddDate(0, 0, -1)
	err := s.db.Where("user_id = ? AND date <= ?", userID, yesterday).
		Order("date DESC").
		Find(&activities).Error

	if err != nil || len(activities) == 0 {
		return 1 // First day
	}

	// Check if yesterday had activity
	if activities[0].Date.Equal(yesterday) {
		// Continue streak from yesterday
		return activities[0].Streak + 1
	}

	// Streak broken, start fresh
	return 1
}

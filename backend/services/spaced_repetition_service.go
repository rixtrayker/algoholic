package services

import (
	"errors"
	"math"
	"time"

	"github.com/yourusername/algoholic/models"
	"gorm.io/gorm"
)

// SpacedRepetitionService implements the SM-2 spaced repetition algorithm
type SpacedRepetitionService struct {
	db *gorm.DB
}

// NewSpacedRepetitionService creates a new spaced repetition service
func NewSpacedRepetitionService(db *gorm.DB) *SpacedRepetitionService {
	return &SpacedRepetitionService{db: db}
}

// QualityFromAttempt converts attempt data to an SM-2 quality rating (0-5)
//
// SM-2 quality scale:
//
//	0 = complete blackout
//	1 = incorrect, remembered upon seeing answer
//	2 = incorrect, but answer seemed easy to recall
//	3 = correct, but with serious difficulty
//	4 = correct, after some hesitation
//	5 = perfect, instant recall
func QualityFromAttempt(isCorrect bool, timeTaken int, estimatedTime *int, hintsUsed int) int {
	if !isCorrect {
		if hintsUsed >= 2 {
			return 0 // complete failure with hints
		}
		return 1 // incorrect
	}

	// Correct answer — determine quality based on speed and hints
	if hintsUsed >= 2 {
		return 3 // correct but needed significant help
	}

	if estimatedTime != nil && *estimatedTime > 0 {
		ratio := float64(timeTaken) / float64(*estimatedTime)
		if ratio <= 0.5 {
			if hintsUsed == 0 {
				return 5 // perfect, fast, no hints
			}
			return 4
		}
		if ratio <= 1.0 {
			if hintsUsed == 0 {
				return 4 // good, within time
			}
			return 3
		}
		// Took longer than estimated
		return 3
	}

	// No estimated time available
	if hintsUsed == 0 {
		return 4 // correct, no hints
	}
	return 3 // correct with hints
}

// ProcessReview processes a review using the SM-2 algorithm
// quality: 0-5 rating from QualityFromAttempt
func (s *SpacedRepetitionService) ProcessReview(userID, questionID int, quality int) error {
	if quality < 0 || quality > 5 {
		return errors.New("quality rating must be between 0 and 5")
	}

	var review models.SpacedRepetitionReview
	err := s.db.Where("user_id = ? AND question_id = ?", userID, questionID).First(&review).Error

	now := time.Now()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// First review for this user/question
		review = models.SpacedRepetitionReview{
			UserID:         userID,
			QuestionID:     questionID,
			EasinessFactor: 2.5,
			IntervalDays:   1,
			Repetitions:    0,
			NextReviewAt:   now,
			LastReviewAt:   &now,
			QualityRating:  &quality,
		}
	} else if err != nil {
		return err
	}

	// Apply SM-2 algorithm
	review.LastReviewAt = &now
	review.QualityRating = &quality

	if quality < 3 {
		// Failed recall — reset repetitions, review again soon
		review.Repetitions = 0
		review.IntervalDays = 1
	} else {
		// Successful recall — increase interval
		switch review.Repetitions {
		case 0:
			review.IntervalDays = 1
		case 1:
			review.IntervalDays = 6
		default:
			review.IntervalDays = int(math.Round(float64(review.IntervalDays) * review.EasinessFactor))
		}
		review.Repetitions++
	}

	// Update easiness factor
	// EF' = EF + (0.1 - (5-q) * (0.08 + (5-q) * 0.02))
	ef := review.EasinessFactor + (0.1 - float64(5-quality)*(0.08+float64(5-quality)*0.02))
	review.EasinessFactor = math.Max(1.3, ef)

	// Schedule next review
	review.NextReviewAt = now.AddDate(0, 0, review.IntervalDays)

	// Save or create
	if review.ReviewID == 0 {
		return s.db.Create(&review).Error
	}
	return s.db.Save(&review).Error
}

// GetDueReviews returns questions that are due for review
func (s *SpacedRepetitionService) GetDueReviews(userID int, limit int) ([]models.SpacedRepetitionReview, error) {
	var reviews []models.SpacedRepetitionReview
	err := s.db.Where("user_id = ? AND next_review_at <= ?", userID, time.Now()).
		Order("next_review_at ASC").
		Limit(limit).
		Find(&reviews).Error
	return reviews, err
}

// GetReviewStats returns review statistics for a user
func (s *SpacedRepetitionService) GetReviewStats(userID int) (total int64, due int64, err error) {
	s.db.Model(&models.SpacedRepetitionReview{}).Where("user_id = ?", userID).Count(&total)
	s.db.Model(&models.SpacedRepetitionReview{}).
		Where("user_id = ? AND next_review_at <= ?", userID, time.Now()).
		Count(&due)
	return total, due, nil
}

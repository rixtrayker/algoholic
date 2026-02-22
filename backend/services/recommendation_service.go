package services

import (
	"fmt"
	"sort"
	"time"

	"gorm.io/gorm"
)

// RecommendationService handles question recommendations
type RecommendationService struct {
	db *gorm.DB
}

// NewRecommendationService creates a new recommendation service
func NewRecommendationService(db *gorm.DB) *RecommendationService {
	return &RecommendationService{db: db}
}

// Recommendation represents a recommended question
type Recommendation struct {
	QuestionID   int     `json:"question_id"`
	ProblemID    *int    `json:"problem_id,omitempty"`
	QuestionText string  `json:"question_text"`
	Reason       string  `json:"reason"`
	Priority     float64 `json:"priority"`
	Difficulty   float64 `json:"difficulty"`
}

// GetRecommendations returns personalized question recommendations
func (rs *RecommendationService) GetRecommendations(userID int, limit int) ([]Recommendation, error) {
	recommendations := []Recommendation{}

	// Strategy 1: Address weaknesses (highest priority)
	weaknessRecs := rs.GetWeaknessBasedRecommendations(userID, limit/3)
	recommendations = append(recommendations, weaknessRecs...)

	// Strategy 2: Progressive difficulty (medium priority)
	progressRecs := rs.GetProgressiveRecommendations(userID, limit/3)
	recommendations = append(recommendations, progressRecs...)

	// Strategy 3: Spaced repetition (lower priority)
	reviewRecs := rs.GetSpacedRepetitionRecommendations(userID, limit/3)
	recommendations = append(recommendations, reviewRecs...)

	// Remove duplicates
	recommendations = rs.deduplicateRecommendations(recommendations)

	// Sort by priority (highest first)
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Priority > recommendations[j].Priority
	})

	// Limit results
	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	return recommendations, nil
}

// GetWeaknessBasedRecommendations finds questions for weak topics
func (rs *RecommendationService) GetWeaknessBasedRecommendations(userID int, limit int) []Recommendation {
	// Find user's weak topics (proficiency < 50)
	var weakTopics []struct {
		TopicID          int
		TopicName        string
		ProficiencyLevel float64
	}

	rs.db.Table("user_skills us").
		Select("us.topic_id, t.name as topic_name, us.proficiency_level").
		Joins("JOIN topics t ON us.topic_id = t.topic_id").
		Where("us.user_id = ? AND us.proficiency_level < 50", userID).
		Order("us.proficiency_level ASC").
		Limit(5). // Top 5 weakest topics
		Scan(&weakTopics)

	recommendations := []Recommendation{}

	for _, wt := range weakTopics {
		// Find questions for this weak topic that user hasn't attempted
		var questions []struct {
			QuestionID   int
			QuestionText string
			ProblemID    *int
			Difficulty   float64
		}

		rs.db.Table("questions q").
			Select("q.question_id, q.question_text, q.problem_id, q.difficulty_score").
			Joins("JOIN problem_topics pt ON q.problem_id = pt.problem_id").
			Where("pt.topic_id = ?", wt.TopicID).
			Where("q.question_id NOT IN (?)",
				rs.db.Table("user_attempts").Select("question_id").Where("user_id = ?", userID),
			).
			Where("q.difficulty_score < 60"). // Start with easier questions for weak topics
			Order("q.difficulty_score ASC").
			Limit(2). // 2 questions per weak topic
			Scan(&questions)

		for _, q := range questions {
			recommendations = append(recommendations, Recommendation{
				QuestionID:   q.QuestionID,
				QuestionText: q.QuestionText,
				ProblemID:    q.ProblemID,
				Reason:       fmt.Sprintf("Practice %s (proficiency: %.0f%%)", wt.TopicName, wt.ProficiencyLevel),
				Priority:     90.0 - wt.ProficiencyLevel, // Lower proficiency = higher priority
				Difficulty:   q.Difficulty,
			})
		}

		if len(recommendations) >= limit {
			break
		}
	}

	return recommendations
}

// GetProgressiveRecommendations suggests slightly harder questions
func (rs *RecommendationService) GetProgressiveRecommendations(userID int, limit int) []Recommendation {
	// Get user's recent performance
	var recentPerformance struct {
		AvgDifficulty float64
		SuccessRate   float64
	}

	rs.db.Table("user_attempts ua").
		Select(`
			COALESCE(AVG(q.difficulty_score), 30) as avg_difficulty,
			COALESCE(AVG(CASE WHEN ua.is_correct THEN 1.0 ELSE 0.0 END), 0.5) as success_rate
		`).
		Joins("JOIN questions q ON ua.question_id = q.question_id").
		Where("ua.user_id = ? AND ua.attempted_at > ?", userID, time.Now().AddDate(0, 0, -7)).
		Scan(&recentPerformance)

	// Determine target difficulty based on success rate
	var targetDifficulty float64
	switch {
	case recentPerformance.SuccessRate > 0.80:
		// User is doing well, challenge them more
		targetDifficulty = recentPerformance.AvgDifficulty + 15.0
	case recentPerformance.SuccessRate > 0.60:
		// Gradual increase
		targetDifficulty = recentPerformance.AvgDifficulty + 10.0
	case recentPerformance.SuccessRate < 0.40:
		// User struggling, make it easier
		targetDifficulty = recentPerformance.AvgDifficulty - 10.0
	default:
		// Maintain current level
		targetDifficulty = recentPerformance.AvgDifficulty + 5.0
	}

	// Clamp target difficulty
	targetDifficulty = clampFloat(targetDifficulty, 10, 100)

	// Find questions at target difficulty
	var questions []struct {
		QuestionID   int
		QuestionText string
		ProblemID    *int
		Difficulty   float64
	}

	rs.db.Table("questions q").
		Select("q.question_id, q.question_text, q.problem_id, q.difficulty_score").
		Where("q.difficulty_score BETWEEN ? AND ?", targetDifficulty-8, targetDifficulty+8).
		Where("q.question_id NOT IN (?)",
			rs.db.Table("user_attempts").Select("question_id").Where("user_id = ?", userID),
		).
		Order(fmt.Sprintf("ABS(q.difficulty_score - %f) ASC", targetDifficulty)). // Closest to target
		Limit(limit).
		Scan(&questions)

	recommendations := []Recommendation{}
	for _, q := range questions {
		recommendations = append(recommendations, Recommendation{
			QuestionID:   q.QuestionID,
			QuestionText: q.QuestionText,
			ProblemID:    q.ProblemID,
			Reason:       fmt.Sprintf("Progressive challenge (difficulty: %.0f)", q.Difficulty),
			Priority:     60.0,
			Difficulty:   q.Difficulty,
		})
	}

	return recommendations
}

// GetSpacedRepetitionRecommendations returns questions due for review
func (rs *RecommendationService) GetSpacedRepetitionRecommendations(userID int, limit int) []Recommendation {
	// Find user skills that need review
	var overdueSkills []struct {
		TopicID     int
		TopicName   string
		LastPracticed time.Time
	}

	rs.db.Table("user_skills us").
		Select("us.topic_id, t.name as topic_name, us.last_practiced_at").
		Joins("JOIN topics t ON us.topic_id = t.topic_id").
		Where("us.user_id = ?", userID).
		Where("us.next_review_at < NOW() OR us.needs_review = true").
		Order("us.next_review_at ASC").
		Limit(5).
		Scan(&overdueSkills)

	recommendations := []Recommendation{}

	for _, skill := range overdueSkills {
		// Find questions for this topic
		var questions []struct {
			QuestionID   int
			QuestionText string
			ProblemID    *int
			Difficulty   float64
		}

		rs.db.Table("questions q").
			Select("q.question_id, q.question_text, q.problem_id, q.difficulty_score").
			Joins("JOIN problem_topics pt ON q.problem_id = pt.problem_id").
			Where("pt.topic_id = ?", skill.TopicID).
			Order("RANDOM()"). // Randomize for variety
			Limit(2).
			Scan(&questions)

		for _, q := range questions {
			daysSince := int(time.Since(skill.LastPracticed).Hours() / 24)
			recommendations = append(recommendations, Recommendation{
				QuestionID:   q.QuestionID,
				QuestionText: q.QuestionText,
				ProblemID:    q.ProblemID,
				Reason:       fmt.Sprintf("Review %s (last practiced %d days ago)", skill.TopicName, daysSince),
				Priority:     50.0,
				Difficulty:   q.Difficulty,
			})
		}

		if len(recommendations) >= limit {
			break
		}
	}

	return recommendations
}

// GetTopicRecommendations gets questions for a specific topic
func (rs *RecommendationService) GetTopicRecommendations(userID int, topicID int, limit int) ([]Recommendation, error) {
	// Get user's proficiency in this topic
	var proficiency float64
	rs.db.Table("user_skills").
		Select("COALESCE(proficiency_level, 0)").
		Where("user_id = ? AND topic_id = ?", userID, topicID).
		Scan(&proficiency)

	// Determine appropriate difficulty range
	var minDiff, maxDiff float64
	switch {
	case proficiency < 30:
		minDiff, maxDiff = 10, 40 // Easier questions
	case proficiency < 60:
		minDiff, maxDiff = 30, 65 // Medium questions
	default:
		minDiff, maxDiff = 50, 100 // Harder questions
	}

	// Find questions
	var questions []struct {
		QuestionID   int
		QuestionText string
		ProblemID    *int
		Difficulty   float64
	}

	rs.db.Table("questions q").
		Select("q.question_id, q.question_text, q.problem_id, q.difficulty_score").
		Joins("JOIN problem_topics pt ON q.problem_id = pt.problem_id").
		Where("pt.topic_id = ?", topicID).
		Where("q.difficulty_score BETWEEN ? AND ?", minDiff, maxDiff).
		Where("q.question_id NOT IN (?)",
			rs.db.Table("user_attempts").Select("question_id").Where("user_id = ?", userID),
		).
		Order("q.difficulty_score ASC").
		Limit(limit).
		Scan(&questions)

	recommendations := []Recommendation{}
	for _, q := range questions {
		recommendations = append(recommendations, Recommendation{
			QuestionID:   q.QuestionID,
			QuestionText: q.QuestionText,
			ProblemID:    q.ProblemID,
			Reason:       fmt.Sprintf("Practice this topic (proficiency: %.0f%%)", proficiency),
			Priority:     70.0,
			Difficulty:   q.Difficulty,
		})
	}

	return recommendations, nil
}

// deduplicateRecommendations removes duplicate question IDs
func (rs *RecommendationService) deduplicateRecommendations(recs []Recommendation) []Recommendation {
	seen := make(map[int]bool)
	unique := []Recommendation{}

	for _, rec := range recs {
		if !seen[rec.QuestionID] {
			seen[rec.QuestionID] = true
			unique = append(unique, rec)
		}
	}

	return unique
}

// Helper functions
func clampFloat(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

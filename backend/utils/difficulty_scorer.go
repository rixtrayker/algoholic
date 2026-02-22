package utils

import (
	"strings"

	"github.com/yourusername/algoholic/models"
	"gorm.io/gorm"
)

// DifficultyComponents represents the 6 components of difficulty scoring
type DifficultyComponents struct {
	Conceptual     float64 // 0-100: Number and depth of concepts needed
	Algorithm      float64 // 0-100: Complexity of the required algorithm
	Implementation float64 // 0-100: How hard to code correctly
	Pattern        float64 // 0-100: How obvious/hidden the pattern is
	EdgeCases      float64 // 0-100: Number of edge cases to handle
	TimePressure   float64 // 0-100: Expected solve time
}

// DifficultyScorer handles all difficulty-related calculations
type DifficultyScorer struct {
	db *gorm.DB
}

// NewDifficultyScorer creates a new difficulty scorer
func NewDifficultyScorer(db *gorm.DB) *DifficultyScorer {
	return &DifficultyScorer{db: db}
}

// CalculateDifficultyScore computes the overall difficulty score for a problem
func (ds *DifficultyScorer) CalculateDifficultyScore(problem *models.Problem) float64 {
	components := ds.GetComponents(problem)

	// Weighted components (must sum to 1.0)
	weights := map[string]float64{
		"conceptual":     0.25,
		"algorithm":      0.20,
		"implementation": 0.15,
		"pattern":        0.20,
		"edge_cases":     0.10,
		"time_pressure":  0.10,
	}

	total := 0.0
	total += components.Conceptual * weights["conceptual"]
	total += components.Algorithm * weights["algorithm"]
	total += components.Implementation * weights["implementation"]
	total += components.Pattern * weights["pattern"]
	total += components.EdgeCases * weights["edge_cases"]
	total += components.TimePressure * weights["time_pressure"]

	return clamp(total, 0, 100)
}

// GetComponents calculates all difficulty components
func (ds *DifficultyScorer) GetComponents(problem *models.Problem) DifficultyComponents {
	return DifficultyComponents{
		Conceptual:     ds.ScoreConceptual(problem),
		Algorithm:      ds.ScoreAlgorithm(problem),
		Implementation: ds.ScoreImplementation(problem),
		Pattern:        ds.ScorePatternRecognition(problem),
		EdgeCases:      ds.ScoreEdgeCases(problem),
		TimePressure:   ds.ScoreTimePressure(problem),
	}
}

// ScoreConceptual evaluates conceptual complexity
func (ds *DifficultyScorer) ScoreConceptual(problem *models.Problem) float64 {
	// Count unique concepts needed
	conceptCount := len(problem.SecondaryPatterns)
	if problem.PrimaryPattern != nil {
		conceptCount++
	}

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

// ScoreAlgorithm evaluates algorithm complexity based on time complexity
func (ds *DifficultyScorer) ScoreAlgorithm(problem *models.Problem) float64 {
	if problem.TimeComplexity == nil {
		return 50.0 // Default if unknown
	}

	complexity := *problem.TimeComplexity

	complexityScores := map[string]float64{
		"O(1)":       10.0,
		"O(log n)":   25.0,
		"O(n)":       35.0,
		"O(n log n)": 50.0,
		"O(n^2)":     65.0,
		"O(n^3)":     80.0,
		"O(2^n)":     90.0,
		"O(n!)":      95.0,
	}

	// Normalize complexity string
	complexity = strings.TrimSpace(complexity)

	if score, ok := complexityScores[complexity]; ok {
		return score
	}

	// Try to match partial strings
	complexityLower := strings.ToLower(complexity)
	if strings.Contains(complexityLower, "factorial") || strings.Contains(complexityLower, "n!") {
		return 95.0
	}
	if strings.Contains(complexityLower, "exponential") || strings.Contains(complexityLower, "2^n") {
		return 90.0
	}
	if strings.Contains(complexityLower, "cubic") || strings.Contains(complexityLower, "n^3") {
		return 80.0
	}
	if strings.Contains(complexityLower, "quadratic") || strings.Contains(complexityLower, "n^2") {
		return 65.0
	}
	if strings.Contains(complexityLower, "linearithmic") || strings.Contains(complexityLower, "n log n") {
		return 50.0
	}
	if strings.Contains(complexityLower, "linear") {
		return 35.0
	}
	if strings.Contains(complexityLower, "logarithmic") || strings.Contains(complexityLower, "log n") {
		return 25.0
	}
	if strings.Contains(complexityLower, "constant") {
		return 10.0
	}

	return 50.0 // Default
}

// ScoreImplementation evaluates implementation difficulty
func (ds *DifficultyScorer) ScoreImplementation(problem *models.Problem) float64 {
	baseScore := 30.0

	// Longer description = more complex implementation
	descLength := len(problem.Description)
	if descLength > 800 {
		baseScore += 30.0
	} else if descLength > 500 {
		baseScore += 20.0
	} else if descLength > 300 {
		baseScore += 10.0
	}

	// More constraints = trickier implementation
	constraintCount := len(problem.Constraints)
	baseScore += float64(constraintCount) * 5.0

	// Number of examples indicates complexity
	if problem.Examples != nil {
		exampleCount := len(problem.Examples)
		if exampleCount > 3 {
			baseScore += 10.0
		}
	}

	return clamp(baseScore, 0, 100)
}

// ScorePatternRecognition evaluates how difficult it is to recognize the pattern
func (ds *DifficultyScorer) ScorePatternRecognition(problem *models.Problem) float64 {
	if problem.PrimaryPattern == nil {
		return 60.0 // Unknown pattern = harder to recognize
	}

	pattern := *problem.PrimaryPattern

	// Common patterns are easier to recognize
	patternDifficulty := map[string]float64{
		// Easy to recognize
		"Array":       20.0,
		"String":      20.0,
		"Hash Table":  25.0,
		"Hash Map":    25.0,
		"Stack":       30.0,
		"Queue":       30.0,

		// Medium difficulty
		"Two Pointers":  35.0,
		"Binary Search": 35.0,
		"Linked List":   40.0,
		"Sliding Window": 40.0,
		"Tree":          45.0,
		"Binary Tree":   45.0,

		// Harder to recognize
		"Graph":              60.0,
		"Dynamic Programming": 70.0,
		"Backtracking":       75.0,
		"Trie":               65.0,
		"Heap":               55.0,
		"Union Find":         65.0,
		"Topological Sort":   70.0,
		"Bit Manipulation":   60.0,

		// Advanced
		"Segment Tree":   80.0,
		"Fenwick Tree":   80.0,
		"Suffix Array":   85.0,
		"KMP":            75.0,
		"Manacher":       85.0,
	}

	if score, ok := patternDifficulty[pattern]; ok {
		return score
	}

	// Default for unknown patterns
	return 50.0
}

// ScoreEdgeCases evaluates edge case complexity
func (ds *DifficultyScorer) ScoreEdgeCases(problem *models.Problem) float64 {
	baseScore := float64(len(problem.Constraints)) * 10.0

	// Check for specific edge case indicators in description
	description := strings.ToLower(problem.Description)
	constraints := strings.Join(problem.Constraints, " ")
	constraintsLower := strings.ToLower(constraints)

	edgeCaseKeywords := map[string]float64{
		"negative":   10.0,
		"zero":       5.0,
		"empty":      5.0,
		"duplicate":  10.0,
		"null":       8.0,
		"overflow":   15.0,
		"edge":       8.0,
		"special":    8.0,
		"distinct":   5.0,
		"unique":     5.0,
	}

	for keyword, score := range edgeCaseKeywords {
		if strings.Contains(description, keyword) || strings.Contains(constraintsLower, keyword) {
			baseScore += score
		}
	}

	// Large constraint ranges indicate more edge cases
	if strings.Contains(constraintsLower, "10^9") || strings.Contains(constraintsLower, "10^5") {
		baseScore += 10.0
	}

	return clamp(baseScore, 0, 100)
}

// ScoreTimePressure evaluates time pressure based on difficulty and expected solve time
func (ds *DifficultyScorer) ScoreTimePressure(problem *models.Problem) float64 {
	// Use official difficulty as baseline
	if problem.OfficialDifficulty != nil {
		difficulty := *problem.OfficialDifficulty

		timePressure := map[string]float64{
			"Easy":   20.0,
			"Medium": 50.0,
			"Hard":   80.0,
		}

		if score, ok := timePressure[difficulty]; ok {
			return score
		}
	}

	// Fallback: use pattern complexity
	if problem.PrimaryPattern != nil {
		pattern := *problem.PrimaryPattern

		patternTime := map[string]float64{
			"Array":               25.0,
			"Hash Table":          25.0,
			"Two Pointers":        35.0,
			"Sliding Window":      40.0,
			"Binary Search":       35.0,
			"Dynamic Programming": 75.0,
			"Graph":               65.0,
			"Backtracking":        80.0,
		}

		if score, ok := patternTime[pattern]; ok {
			return score
		}
	}

	return 50.0 // Default
}

// PersonalizedDifficulty calculates user-specific difficulty
func (ds *DifficultyScorer) PersonalizedDifficulty(problemID, userID int) float64 {
	// Get base difficulty
	var problem models.Problem
	if err := ds.db.First(&problem, problemID).Error; err != nil {
		return 50.0 // Default if problem not found
	}

	baseDifficulty := problem.DifficultyScore

	// Get user's proficiency in related topics
	var avgProficiency float64
	err := ds.db.Table("user_skills").
		Select("COALESCE(AVG(proficiency_level), 50)").
		Where("user_id = ? AND topic_id IN (?)",
			userID,
			ds.db.Table("problem_topics").Select("topic_id").Where("problem_id = ?", problemID),
		).
		Scan(&avgProficiency).Error

	if err != nil {
		return baseDifficulty
	}

	// Adjust difficulty based on user proficiency
	// Strong user (proficiency > 70) → easier (reduce difficulty)
	// Weak user (proficiency < 30) → harder (increase difficulty)
	// Neutral user (proficiency ~50) → no adjustment
	adjustment := (50 - avgProficiency) * 0.6

	personalizedScore := baseDifficulty + adjustment

	return clamp(personalizedScore, 0, 100)
}

// RecalibrateDifficulty updates difficulty based on user performance data
func (ds *DifficultyScorer) RecalibrateDifficulty(problemID int) error {
	var problem models.Problem
	if err := ds.db.First(&problem, problemID).Error; err != nil {
		return err
	}

	// Get statistics from user attempts
	var stats struct {
		SuccessRate  float64
		AvgTime      float64
		AttemptCount int64
	}

	ds.db.Table("user_attempts ua").
		Select(`
			COALESCE(AVG(CASE WHEN ua.is_correct THEN 1.0 ELSE 0.0 END), 0) as success_rate,
			COALESCE(AVG(ua.time_taken_seconds), 0) as avg_time,
			COUNT(*) as attempt_count
		`).
		Joins("JOIN questions q ON ua.question_id = q.question_id").
		Where("q.problem_id = ?", problemID).
		Scan(&stats)

	// Need at least 10 attempts for reliable calibration
	if stats.AttemptCount < 10 {
		return nil // Not enough data
	}

	currentDifficulty := problem.DifficultyScore

	// Adjust based on success rate
	var adjustment float64
	switch {
	case stats.SuccessRate > 0.80:
		// Too easy, increase difficulty
		adjustment = 5.0
	case stats.SuccessRate > 0.60:
		// Slightly easy
		adjustment = 2.0
	case stats.SuccessRate < 0.30:
		// Too hard, decrease difficulty
		adjustment = -5.0
	case stats.SuccessRate < 0.45:
		// Slightly hard
		adjustment = -2.0
	default:
		// Just right, no adjustment
		adjustment = 0.0
	}

	// Cap adjustment to prevent wild swings
	newDifficulty := clamp(currentDifficulty+adjustment, currentDifficulty-10, currentDifficulty+10)
	newDifficulty = clamp(newDifficulty, 0, 100)

	// Update problem difficulty
	return ds.db.Model(&problem).Update("difficulty_score", newDifficulty).Error
}

// Helper functions

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// GetDifficultyTier returns a human-readable tier for a difficulty score
func GetDifficultyTier(score float64) string {
	switch {
	case score <= 20:
		return "Trivial"
	case score <= 35:
		return "Easy"
	case score <= 50:
		return "Medium-Easy"
	case score <= 65:
		return "Medium"
	case score <= 80:
		return "Hard"
	case score <= 95:
		return "Very Hard"
	default:
		return "Expert"
	}
}

// GetDifficultyColor returns a color code for UI display
func GetDifficultyColor(score float64) string {
	switch {
	case score <= 35:
		return "#22c55e" // Green
	case score <= 65:
		return "#f59e0b" // Yellow/Orange
	default:
		return "#ef4444" // Red
	}
}

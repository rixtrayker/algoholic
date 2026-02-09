package services

import (
	"errors"

	"gorm.io/gorm"
	"github.com/yourusername/algoholic/models"
)

// ProblemService handles problem-related operations
type ProblemService struct {
	db *gorm.DB
}

// NewProblemService creates a new problem service
func NewProblemService(db *gorm.DB) *ProblemService {
	return &ProblemService{db: db}
}

// GetProblems retrieves problems with filters
func (s *ProblemService) GetProblems(minDifficulty, maxDifficulty float64, pattern string, limit, offset int) ([]models.Problem, int64, error) {
	query := s.db.Model(&models.Problem{})

	// Apply filters
	query = query.Where("difficulty_score BETWEEN ? AND ?", minDifficulty, maxDifficulty)
	if pattern != "" {
		query = query.Where("primary_pattern = ? OR ? = ANY(secondary_patterns)", pattern, pattern)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get problems
	var problems []models.Problem
	if err := query.Limit(limit).Offset(offset).Find(&problems).Error; err != nil {
		return nil, 0, err
	}

	return problems, total, nil
}

// GetProblemByID retrieves a problem by ID
func (s *ProblemService) GetProblemByID(id int) (*models.Problem, error) {
	var problem models.Problem
	if err := s.db.First(&problem, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("problem not found")
		}
		return nil, err
	}
	return &problem, nil
}

// GetProblemBySlug retrieves a problem by slug
func (s *ProblemService) GetProblemBySlug(slug string) (*models.Problem, error) {
	var problem models.Problem
	if err := s.db.Where("slug = ?", slug).First(&problem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("problem not found")
		}
		return nil, err
	}
	return &problem, nil
}

// SearchProblems searches problems by title or description
func (s *ProblemService) SearchProblems(query string, limit, offset int) ([]models.Problem, int64, error) {
	var problems []models.Problem
	var total int64

	searchQuery := s.db.Model(&models.Problem{}).
		Where("title ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%")

	if err := searchQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := searchQuery.Limit(limit).Offset(offset).Find(&problems).Error; err != nil {
		return nil, 0, err
	}

	return problems, total, nil
}

// GetProblemTopics retrieves topics for a problem
func (s *ProblemService) GetProblemTopics(problemID int) ([]models.Topic, error) {
	var topics []models.Topic
	err := s.db.Table("topics").
		Joins("JOIN problem_topics ON problem_topics.topic_id = topics.topic_id").
		Where("problem_topics.problem_id = ?", problemID).
		Order("problem_topics.is_primary DESC, problem_topics.relevance_score DESC").
		Find(&topics).Error

	return topics, err
}

// UpdateProblemStats updates problem statistics after an attempt
func (s *ProblemService) UpdateProblemStats(problemID int, solved bool, timeTaken int) error {
	updates := map[string]interface{}{
		"total_attempts": gorm.Expr("total_attempts + 1"),
	}

	if solved {
		updates["total_solves"] = gorm.Expr("total_solves + 1")
	}

	return s.db.Model(&models.Problem{}).
		Where("problem_id = ?", problemID).
		Updates(updates).Error
}

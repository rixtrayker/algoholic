package services

import (
	"encoding/json"
	"errors"

	"github.com/yourusername/algoholic/models"
	"gorm.io/gorm"
)

type ListService struct {
	db *gorm.DB
}

func NewListService(db *gorm.DB) *ListService {
	return &ListService{db: db}
}

// GetUserLists returns paginated lists for a user
func (s *ListService) GetUserLists(userID int, limit, offset int) ([]models.UserList, int64, error) {
	var total int64
	s.db.Model(&models.UserList{}).Where("user_id = ?", userID).Count(&total)

	var lists []models.UserList
	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&lists).Error
	return lists, total, err
}

// GetList returns a specific list if it belongs to the user
func (s *ListService) GetList(listID, userID int) (*models.UserList, error) {
	var list models.UserList
	err := s.db.Where("list_id = ? AND user_id = ?", listID, userID).
		First(&list).Error
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// CreateList creates a new list for a user
func (s *ListService) CreateList(userID int, name string, description *string, isPublic bool) (*models.UserList, error) {
	emptyProblems := models.JSONB{}
	if err := json.Unmarshal([]byte("[]"), &emptyProblems); err != nil {
		return nil, err
	}

	list := models.UserList{
		UserID:      userID,
		Name:        name,
		Description: description,
		IsPublic:    isPublic,
		ProblemIDs:  emptyProblems,
		TotalItems:  0,
		Completed:   0,
	}

	if err := s.db.Create(&list).Error; err != nil {
		return nil, err
	}

	return &list, nil
}

// UpdateList updates a list's details
func (s *ListService) UpdateList(listID, userID int, name, description *string, isPublic *bool) (*models.UserList, error) {
	list, err := s.GetList(listID, userID)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	if name != nil {
		updates["name"] = *name
	}
	if description != nil {
		updates["description"] = *description
	}
	if isPublic != nil {
		updates["is_public"] = *isPublic
	}

	if err := s.db.Model(list).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Reload the list to get updated values
	if err := s.db.First(list, listID).Error; err != nil {
		return nil, err
	}

	return list, nil
}

// DeleteList deletes a list
func (s *ListService) DeleteList(listID, userID int) error {
	result := s.db.Where("list_id = ? AND user_id = ?", listID, userID).
		Delete(&models.UserList{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// AddProblemToList adds a problem to a list
func (s *ListService) AddProblemToList(listID, userID, problemID int) (*models.UserList, error) {
	list, err := s.GetList(listID, userID)
	if err != nil {
		return nil, err
	}

	// Verify problem exists
	var problem models.Problem
	if err := s.db.First(&problem, problemID).Error; err != nil {
		return nil, errors.New("problem not found")
	}

	// Parse current problem IDs
	var problemIDs []int
	jsonBytes, err := json.Marshal(list.ProblemIDs)
	if err != nil {
		problemIDs = []int{}
	} else if err := json.Unmarshal(jsonBytes, &problemIDs); err != nil {
		problemIDs = []int{}
	}

	// Check if problem already in list
	for _, pid := range problemIDs {
		if pid == problemID {
			return list, nil // Already in list, return as-is
		}
	}

	// Add problem to list
	problemIDs = append(problemIDs, problemID)
	newProblemIDs, err := json.Marshal(problemIDs)
	if err != nil {
		return nil, err
	}

	// Update list
	if err := s.db.Model(list).Updates(map[string]interface{}{
		"problem_ids": newProblemIDs,
		"total_items": len(problemIDs),
	}).Error; err != nil {
		return nil, err
	}

	// Reload list
	if err := s.db.First(list, listID).Error; err != nil {
		return nil, err
	}

	return list, nil
}

// RemoveProblemFromList removes a problem from a list
func (s *ListService) RemoveProblemFromList(listID, userID, problemID int) (*models.UserList, error) {
	list, err := s.GetList(listID, userID)
	if err != nil {
		return nil, err
	}

	// Parse current problem IDs
	var problemIDs []int
	jsonBytes, err := json.Marshal(list.ProblemIDs)
	if err != nil {
		problemIDs = []int{}
	} else if err := json.Unmarshal(jsonBytes, &problemIDs); err != nil {
		problemIDs = []int{}
	}

	// Remove problem from list
	newProblemIDs := []int{}
	for _, pid := range problemIDs {
		if pid != problemID {
			newProblemIDs = append(newProblemIDs, pid)
		}
	}

	// Marshal back to JSON
	updatedProblemIDs, err := json.Marshal(newProblemIDs)
	if err != nil {
		return nil, err
	}

	// Update list
	if err := s.db.Model(list).Updates(map[string]interface{}{
		"problem_ids": updatedProblemIDs,
		"total_items": len(newProblemIDs),
	}).Error; err != nil {
		return nil, err
	}

	// Reload list
	if err := s.db.First(list, listID).Error; err != nil {
		return nil, err
	}

	return list, nil
}

// GetListProblems returns all problems in a list
func (s *ListService) GetListProblems(listID, userID int) ([]models.Problem, error) {
	list, err := s.GetList(listID, userID)
	if err != nil {
		return nil, err
	}

	// Parse problem IDs
	var problemIDs []int
	jsonBytes, err := json.Marshal(list.ProblemIDs)
	if err != nil {
		return []models.Problem{}, nil
	}
	if err := json.Unmarshal(jsonBytes, &problemIDs); err != nil {
		return []models.Problem{}, nil
	}

	if len(problemIDs) == 0 {
		return []models.Problem{}, nil
	}

	// Fetch all problems
	var problems []models.Problem
	err = s.db.Where("problem_id IN ?", problemIDs).Find(&problems).Error
	if err != nil {
		return nil, err
	}

	return problems, nil
}

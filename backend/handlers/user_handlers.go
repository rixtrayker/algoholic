package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/algoholic/middleware"
	"github.com/yourusername/algoholic/models"
	"github.com/yourusername/algoholic/services"
)

type UserHandler struct {
	userService      *services.UserService
	questionService  *services.QuestionService
	spacedRepService *services.SpacedRepetitionService
}

func NewUserHandler(userService *services.UserService, questionService *services.QuestionService, spacedRepService *services.SpacedRepetitionService) *UserHandler {
	return &UserHandler{
		userService:      userService,
		questionService:  questionService,
		spacedRepService: spacedRepService,
	}
}

// GetUserStats retrieves user statistics
func (h *UserHandler) GetUserStats(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	stats, err := h.userService.GetUserStats(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve stats",
		})
	}

	return c.JSON(stats)
}

// GetReviewQueue retrieves topics that need review
func (h *UserHandler) GetReviewQueue(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	reviewQueue, err := h.userService.GetReviewQueue(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve review queue",
		})
	}

	return c.JSON(fiber.Map{
		"review_queue": reviewQueue,
		"count":        len(reviewQueue),
	})
}

// GetWeaknesses retrieves user weaknesses
func (h *UserHandler) GetWeaknesses(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	limit := c.QueryInt("limit", 10)

	weakTopics, err := h.userService.GetWeakTopics(userID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve weaknesses",
		})
	}

	return c.JSON(fiber.Map{
		"weak_topics": weakTopics,
		"count":       len(weakTopics),
	})
}

// GetRecommendations provides personalized recommendations
func (h *UserHandler) GetRecommendations(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Get weak topics
	weakTopics, err := h.userService.GetWeakTopics(userID, 3)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate recommendations",
		})
	}

	recommendations := make([]fiber.Map, 0)

	// Generate recommendations based on weak topics
	for _, topic := range weakTopics {
		recommendations = append(recommendations, fiber.Map{
			"type":        "practice_topic",
			"topic":       topic,
			"reason":      "Low proficiency - needs practice",
			"priority":    "high",
			"action":      "Practice questions for this topic",
		})
	}

	// Get review queue
	reviewQueue, err := h.userService.GetReviewQueue(userID)
	if err == nil && len(reviewQueue) > 0 {
		recommendations = append(recommendations, fiber.Map{
			"type":     "review",
			"count":    len(reviewQueue),
			"reason":   "Topics due for spaced repetition review",
			"priority": "medium",
			"action":   "Review these topics to maintain mastery",
		})
	}

	// Get recent attempts to check for patterns
	recentAttempts, err := h.questionService.GetRecentAttempts(userID, 10)
	if err == nil && len(recentAttempts) > 0 {
		incorrectCount := 0
		for _, attempt := range recentAttempts {
			if !attempt.IsCorrect {
				incorrectCount++
			}
		}

		if incorrectCount > 5 {
			recommendations = append(recommendations, fiber.Map{
				"type":     "difficulty_adjustment",
				"reason":   "High error rate in recent attempts",
				"priority": "high",
				"action":   "Consider reviewing fundamentals or reducing difficulty",
			})
		}
	}

	return c.JSON(fiber.Map{
		"recommendations": recommendations,
		"count":           len(recommendations),
	})
}

// GetUserProgress retrieves progress for a specific topic
func (h *UserHandler) GetUserProgress(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	topicID, err := strconv.Atoi(c.Params("topicId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid topic ID",
		})
	}

	progress, err := h.userService.GetUserProgress(userID, topicID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No progress found for this topic",
		})
	}

	return c.JSON(progress)
}

// GetUserSkills retrieves all user skills
func (h *UserHandler) GetUserSkills(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	skills, err := h.userService.GetUserSkills(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve skills",
		})
	}

	return c.JSON(fiber.Map{
		"skills": skills,
		"count":  len(skills),
	})
}

// GetPreferences retrieves user preferences
func (h *UserHandler) GetPreferences(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	preferences, err := h.userService.GetUserPreferences(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve preferences",
		})
	}

	return c.JSON(fiber.Map{
		"preferences": preferences,
	})
}

// UpdatePreferences updates user preferences
func (h *UserHandler) UpdatePreferences(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var preferences models.JSONB
	if err := c.BodyParser(&preferences); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.userService.UpdateUserPreferences(userID, preferences); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update preferences",
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Preferences updated successfully",
		"preferences": preferences,
	})
}

// GetRecentAttempts retrieves user's recent attempts
func (h *UserHandler) GetRecentAttempts(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	limit := c.QueryInt("limit", 20)

	attempts, err := h.questionService.GetRecentAttempts(userID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve attempts",
		})
	}

	return c.JSON(fiber.Map{
		"attempts": attempts,
		"count":    len(attempts),
	})
}

// GetDueReviews retrieves questions due for spaced repetition review
func (h *UserHandler) GetDueReviews(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	limit := c.QueryInt("limit", 20)

	reviews, err := h.spacedRepService.GetDueReviews(userID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve due reviews",
		})
	}

	total, due, _ := h.spacedRepService.GetReviewStats(userID)

	return c.JSON(fiber.Map{
		"reviews":       reviews,
		"count":         len(reviews),
		"total_tracked": total,
		"total_due":     due,
	})
}

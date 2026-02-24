package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/algoholic/models"
	"gorm.io/gorm"
)

type TopicHandler struct {
	db *gorm.DB
}

func NewTopicHandler(db *gorm.DB) *TopicHandler {
	return &TopicHandler{db: db}
}

func (h *TopicHandler) GetAllTopics(c *fiber.Ctx) error {
	var topics []models.Topic

	category := c.Query("category", "")
	parentID := c.Query("parent_id", "")

	query := h.db.Model(&models.Topic{})

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if parentID != "" {
		if parentID == "null" {
			query = query.Where("parent_topic_id IS NULL")
		} else {
			pid, err := strconv.Atoi(parentID)
			if err == nil {
				query = query.Where("parent_topic_id = ?", pid)
			}
		}
	}

	if err := query.Order("name ASC").Find(&topics).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve topics",
		})
	}

	return c.JSON(fiber.Map{
		"topics": topics,
		"count":  len(topics),
	})
}

func (h *TopicHandler) GetTopicByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid topic ID",
		})
	}

	var topic models.Topic
	if err := h.db.First(&topic, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Topic not found",
		})
	}

	return c.JSON(topic)
}

type TopicPerformance struct {
	TopicID            int     `json:"topic_id"`
	TopicName          string  `json:"topic_name"`
	ProficiencyLevel   float64 `json:"proficiency_level"`
	QuestionsAttempted int     `json:"questions_attempted"`
	QuestionsCorrect   int     `json:"questions_correct"`
	AccuracyRate       float64 `json:"accuracy_rate"`
	NeedsReview        bool    `json:"needs_review"`
	LastPracticedAt    *string `json:"last_practiced_at,omitempty"`
	NextReviewAt       *string `json:"next_review_at,omitempty"`
}

func (h *TopicHandler) GetTopicPerformance(c *fiber.Ctx) error {
	// Use authenticated user ID from JWT instead of URL param
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	topicIDStr := c.Params("topicId")
	topicID, err := strconv.Atoi(topicIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid topic ID",
		})
	}

	var topic models.Topic
	if err := h.db.First(&topic, topicID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Topic not found",
		})
	}

	performance := TopicPerformance{
		TopicID:   topicID,
		TopicName: topic.Name,
	}

	var skill models.UserSkill
	result := h.db.Where("user_id = ? AND topic_id = ?", userID, topicID).First(&skill)

	if result.Error == nil {
		performance.ProficiencyLevel = skill.ProficiencyLevel
		performance.QuestionsAttempted = skill.QuestionsAttempted
		performance.QuestionsCorrect = skill.QuestionsCorrect
		if skill.QuestionsAttempted > 0 {
			performance.AccuracyRate = float64(skill.QuestionsCorrect) / float64(skill.QuestionsAttempted) * 100
		}
		performance.NeedsReview = skill.NeedsReview
		if skill.LastPracticedAt != nil {
			t := skill.LastPracticedAt.Format("2006-01-02T15:04:05Z")
			performance.LastPracticedAt = &t
		}
		if skill.NextReviewAt != nil {
			t := skill.NextReviewAt.Format("2006-01-02T15:04:05Z")
			performance.NextReviewAt = &t
		}
	}

	return c.JSON(performance)
}

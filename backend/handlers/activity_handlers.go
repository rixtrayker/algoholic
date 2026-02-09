package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/algoholic/services"
	"gorm.io/gorm"
)

type ActivityHandler struct {
	activityService *services.ActivityService
}

func NewActivityHandler(db *gorm.DB) *ActivityHandler {
	return &ActivityHandler{
		activityService: services.NewActivityService(db),
	}
}

// GetActivityChart returns daily activity data for the commitment chart
func (h *ActivityHandler) GetActivityChart(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	// Default to 365 days (1 year)
	days := 365
	if daysParam := c.Query("days"); daysParam != "" {
		if d, err := strconv.Atoi(daysParam); err == nil && d > 0 {
			days = d
		}
	}

	activities, err := h.activityService.GetActivityData(userID, days)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch activity data",
		})
	}

	return c.JSON(activities)
}

// GetActivityStats returns aggregated activity statistics
func (h *ActivityHandler) GetActivityStats(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	stats, err := h.activityService.GetActivityStats(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch activity stats",
		})
	}

	return c.JSON(stats)
}

// GetPracticeHistory returns detailed practice history
func (h *ActivityHandler) GetPracticeHistory(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	// Default to 30 days
	days := 30
	if daysParam := c.Query("days"); daysParam != "" {
		if d, err := strconv.Atoi(daysParam); err == nil && d > 0 {
			days = d
		}
	}

	history, err := h.activityService.GetPracticeHistory(userID, days)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch practice history",
		})
	}

	return c.JSON(history)
}

// RecordActivity manually records activity (usually triggered by attempt submission)
func (h *ActivityHandler) RecordActivity(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var req struct {
		ProblemsCount  int `json:"problems_count"`
		QuestionsCount int `json:"questions_count"`
		StudyTime      int `json:"study_time_seconds"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.activityService.RecordActivity(userID, req.ProblemsCount, req.QuestionsCount, req.StudyTime); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to record activity",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

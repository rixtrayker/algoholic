package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/algoholic/middleware"
	"github.com/yourusername/algoholic/services"
	"github.com/yourusername/algoholic/utils"
)

type TrainingPlanHandler struct {
	trainingPlanService *services.TrainingPlanService
}

func NewTrainingPlanHandler(trainingPlanService *services.TrainingPlanService) *TrainingPlanHandler {
	return &TrainingPlanHandler{trainingPlanService: trainingPlanService}
}

// CreateTrainingPlan creates a new training plan
func (h *TrainingPlanHandler) CreateTrainingPlan(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req services.CreatePlanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	plan, err := h.trainingPlanService.CreateTrainingPlan(userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Training plan created successfully",
		"plan":    plan,
	})
}

// GetUserPlans retrieves paginated training plans for the user
func (h *TrainingPlanHandler) GetUserPlans(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	pagination := utils.ParsePagination(c)

	plans, total, err := h.trainingPlanService.GetUserPlans(userID, pagination.PageSize, pagination.Offset())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve training plans",
		})
	}

	return c.JSON(utils.NewPaginatedResponse(plans, total, pagination))
}

// GetTrainingPlan retrieves a specific training plan
func (h *TrainingPlanHandler) GetTrainingPlan(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	planID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid plan ID",
		})
	}

	plan, err := h.trainingPlanService.GetPlanByID(planID, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(plan)
}

// GetNextQuestion gets the next question in the training plan
func (h *TrainingPlanHandler) GetNextQuestion(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	planID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid plan ID",
		})
	}

	question, err := h.trainingPlanService.GetNextQuestion(planID, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(question)
}

// GetPlanItems retrieves all items in a training plan
func (h *TrainingPlanHandler) GetPlanItems(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	planID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid plan ID",
		})
	}

	items, err := h.trainingPlanService.GetPlanItems(planID, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"items": items,
		"count": len(items),
	})
}

// GetTodaysQuestions retrieves today's questions from the plan
func (h *TrainingPlanHandler) GetTodaysQuestions(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	planID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid plan ID",
		})
	}

	questions, err := h.trainingPlanService.GetTodaysQuestions(planID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"questions": questions,
		"count":     len(questions),
	})
}

// CompleteItem marks an item as completed
func (h *TrainingPlanHandler) CompleteItem(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	planID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid plan ID",
		})
	}

	itemID, err := strconv.Atoi(c.Params("itemId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid item ID",
		})
	}

	if err := h.trainingPlanService.CompleteItem(planID, userID, itemID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Item marked as completed",
	})
}

// PausePlan pauses a training plan
func (h *TrainingPlanHandler) PausePlan(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	planID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid plan ID",
		})
	}

	if err := h.trainingPlanService.PausePlan(planID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Training plan paused",
	})
}

// ResumePlan resumes a paused training plan
func (h *TrainingPlanHandler) ResumePlan(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	planID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid plan ID",
		})
	}

	if err := h.trainingPlanService.ResumePlan(planID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Training plan resumed",
	})
}

// DeletePlan deletes a training plan
func (h *TrainingPlanHandler) DeletePlan(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	planID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid plan ID",
		})
	}

	if err := h.trainingPlanService.DeletePlan(planID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Training plan deleted",
	})
}

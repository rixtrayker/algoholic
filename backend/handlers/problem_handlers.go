package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/algoholic/backend/services"
)

type ProblemHandler struct {
	problemService *services.ProblemService
}

func NewProblemHandler(problemService *services.ProblemService) *ProblemHandler {
	return &ProblemHandler{problemService: problemService}
}

// GetProblems retrieves problems with filters
func (h *ProblemHandler) GetProblems(c *fiber.Ctx) error {
	// Parse query parameters
	minDiff := c.QueryFloat("min_difficulty", 0)
	maxDiff := c.QueryFloat("max_difficulty", 100)
	pattern := c.Query("pattern", "")
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	problems, total, err := h.problemService.GetProblems(minDiff, maxDiff, pattern, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve problems",
		})
	}

	return c.JSON(fiber.Map{
		"problems": problems,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	})
}

// GetProblem retrieves a single problem by ID
func (h *ProblemHandler) GetProblem(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid problem ID",
		})
	}

	problem, err := h.problemService.GetProblemByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Problem not found",
		})
	}

	return c.JSON(problem)
}

// GetProblemBySlug retrieves a problem by slug
func (h *ProblemHandler) GetProblemBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")

	problem, err := h.problemService.GetProblemBySlug(slug)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Problem not found",
		})
	}

	return c.JSON(problem)
}

// SearchProblems searches for problems
func (h *ProblemHandler) SearchProblems(c *fiber.Ctx) error {
	query := c.Query("q", "")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Search query is required",
		})
	}

	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	problems, total, err := h.problemService.SearchProblems(query, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search problems",
		})
	}

	return c.JSON(fiber.Map{
		"problems": problems,
		"total":    total,
		"query":    query,
	})
}

// GetProblemTopics retrieves topics for a problem
func (h *ProblemHandler) GetProblemTopics(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid problem ID",
		})
	}

	topics, err := h.problemService.GetProblemTopics(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve topics",
		})
	}

	return c.JSON(fiber.Map{
		"topics": topics,
	})
}

package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/algoholic/services"
	"gorm.io/gorm"
)

type ListHandler struct {
	listService *services.ListService
}

func NewListHandler(db *gorm.DB) *ListHandler {
	return &ListHandler{
		listService: services.NewListService(db),
	}
}

// GetUserLists returns all lists for the authenticated user
func (h *ListHandler) GetUserLists(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	lists, err := h.listService.GetUserLists(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch lists",
		})
	}

	return c.JSON(lists)
}

// GetList returns a specific list by ID
func (h *ListHandler) GetList(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	listID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid list ID",
		})
	}

	list, err := h.listService.GetList(listID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "List not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch list",
		})
	}

	return c.JSON(list)
}

// CreateList creates a new list
func (h *ListHandler) CreateList(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var req struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
		IsPublic    bool    `json:"is_public"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "List name is required",
		})
	}

	list, err := h.listService.CreateList(userID, req.Name, req.Description, req.IsPublic)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create list",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(list)
}

// UpdateList updates an existing list
func (h *ListHandler) UpdateList(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	listID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid list ID",
		})
	}

	var req struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		IsPublic    *bool   `json:"is_public"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	list, err := h.listService.UpdateList(listID, userID, req.Name, req.Description, req.IsPublic)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "List not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update list",
		})
	}

	return c.JSON(list)
}

// DeleteList deletes a list
func (h *ListHandler) DeleteList(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	listID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid list ID",
		})
	}

	if err := h.listService.DeleteList(listID, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "List not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete list",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// AddProblemToList adds a problem to a list
func (h *ListHandler) AddProblemToList(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	listID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid list ID",
		})
	}

	var req struct {
		ProblemID int `json:"problem_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	list, err := h.listService.AddProblemToList(listID, userID, req.ProblemID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "List not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add problem to list",
		})
	}

	return c.JSON(list)
}

// RemoveProblemFromList removes a problem from a list
func (h *ListHandler) RemoveProblemFromList(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	listID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid list ID",
		})
	}

	problemID, err := strconv.Atoi(c.Params("problemId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid problem ID",
		})
	}

	list, err := h.listService.RemoveProblemFromList(listID, userID, problemID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "List not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove problem from list",
		})
	}

	return c.JSON(list)
}

// GetListProblems returns all problems in a list
func (h *ListHandler) GetListProblems(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	listID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid list ID",
		})
	}

	problems, err := h.listService.GetListProblems(listID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "List not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch list problems",
		})
	}

	return c.JSON(problems)
}

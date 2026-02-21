package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/algoholic/middleware"
	"github.com/yourusername/algoholic/services"
)

type QuestionHandler struct {
	questionService *services.QuestionService
	userService     *services.UserService
}

func NewQuestionHandler(questionService *services.QuestionService, userService *services.UserService) *QuestionHandler {
	return &QuestionHandler{
		questionService: questionService,
		userService:     userService,
	}
}

// GetQuestions retrieves questions with filters
func (h *QuestionHandler) GetQuestions(c *fiber.Ctx) error {
	questionType := c.Query("type", "")
	minDiff := c.QueryFloat("min_difficulty", 0)
	maxDiff := c.QueryFloat("max_difficulty", 100)
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	questions, total, err := h.questionService.GetQuestions(questionType, minDiff, maxDiff, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve questions",
		})
	}

	return c.JSON(fiber.Map{
		"questions": questions,
		"total":     total,
		"limit":     limit,
		"offset":    offset,
	})
}

// GetQuestion retrieves a single question by ID
func (h *QuestionHandler) GetQuestion(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid question ID",
		})
	}

	question, err := h.questionService.GetQuestionByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Question not found",
		})
	}

	return c.JSON(question)
}

// GetRandomQuestion gets a random question
func (h *QuestionHandler) GetRandomQuestion(c *fiber.Ctx) error {
	questionType := c.Query("type", "")
	minDiff := c.QueryFloat("min_difficulty", 0)
	maxDiff := c.QueryFloat("max_difficulty", 100)

	question, err := h.questionService.GetRandomQuestion(questionType, minDiff, maxDiff, nil)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No questions found matching criteria",
		})
	}

	return c.JSON(question)
}

// SubmitAnswer handles question answer submission
func (h *QuestionHandler) SubmitAnswer(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid question ID",
		})
	}

	var req services.AnswerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	req.QuestionID = id

	response, err := h.questionService.SubmitAnswer(userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Update user streak
	h.userService.UpdateStreak(userID)

	// Add study time
	h.userService.AddStudyTime(userID, int64(req.TimeTaken))

	return c.JSON(response)
}

// GetQuestionsByProblem retrieves questions for a specific problem
func (h *QuestionHandler) GetQuestionsByProblem(c *fiber.Ctx) error {
	problemID, err := strconv.Atoi(c.Params("problemId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid problem ID",
		})
	}

	questions, err := h.questionService.GetQuestionsByProblem(problemID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve questions",
		})
	}

	return c.JSON(fiber.Map{
		"questions": questions,
		"count":     len(questions),
	})
}

// GetUserAttempts retrieves user's attempts for a question
func (h *QuestionHandler) GetUserAttempts(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	questionID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid question ID",
		})
	}

	attempts, err := h.questionService.GetUserAttempts(userID, questionID)
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

// GetHint retrieves a hint for a question
func (h *QuestionHandler) GetHint(c *fiber.Ctx) error {
	questionID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid question ID",
		})
	}

	hintLevel := c.QueryInt("level", 1)
	if hintLevel < 1 || hintLevel > 3 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Hint level must be between 1 and 3",
		})
	}

	hint, err := h.questionService.GetHint(questionID, hintLevel)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userID, _ := middleware.GetUserID(c)
	if userID > 0 {
		h.questionService.RecordHintUsage(userID, questionID, hintLevel)
	}

	return c.JSON(fiber.Map{
		"hint":        hint,
		"level":       hintLevel,
		"question_id": questionID,
	})
}

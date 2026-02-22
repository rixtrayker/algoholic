package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/yourusername/algoholic/models"
	"github.com/yourusername/algoholic/services"
)

// SearchHandler exposes Phase 2 intelligence endpoints.
type SearchHandler struct {
	db      *gorm.DB
	vector  *services.VectorService
	graph   *services.GraphService
}

// NewSearchHandler wires up the search handler.
func NewSearchHandler(db *gorm.DB, vector *services.VectorService, graph *services.GraphService) *SearchHandler {
	return &SearchHandler{db: db, vector: vector, graph: graph}
}

// ── Semantic search ───────────────────────────────────────────────────────────

// SemanticSearchProblems godoc
// GET /api/search/problems?q=two+sum&limit=10
func (h *SearchHandler) SemanticSearchProblems(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "query parameter 'q' is required"})
	}
	limit := c.QueryInt("limit", 10)
	if limit > 50 {
		limit = 50
	}

	// Try vector search first
	if h.vector != nil && h.vector.IsAvailable() {
		docs, err := h.vector.SearchSimilarProblems(query, limit)
		if err == nil && len(docs) > 0 {
			// Enrich with full problem data
			results := h.enrichProblems(docs)
			return c.JSON(fiber.Map{
				"results": results,
				"count":   len(results),
				"source":  "vector",
			})
		}
	}

	// Fallback: PostgreSQL full-text search
	var problems []models.Problem
	h.db.Where(
		"to_tsvector('english', title || ' ' || description) @@ plainto_tsquery('english', ?)", query,
	).Limit(limit).Find(&problems)

	return c.JSON(fiber.Map{
		"results": problems,
		"count":   len(problems),
		"source":  "fulltext",
	})
}

// SemanticSearchQuestions godoc
// GET /api/search/questions?q=hash+map&limit=10
func (h *SearchHandler) SemanticSearchQuestions(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "query parameter 'q' is required"})
	}
	limit := c.QueryInt("limit", 10)
	if limit > 50 {
		limit = 50
	}

	if h.vector != nil && h.vector.IsAvailable() {
		docs, err := h.vector.SearchSimilarQuestions(query, limit)
		if err == nil && len(docs) > 0 {
			results := h.enrichQuestions(docs)
			return c.JSON(fiber.Map{
				"results": results,
				"count":   len(results),
				"source":  "vector",
			})
		}
	}

	// Fallback: ILIKE search
	var questions []models.Question
	h.db.Where("question_text ILIKE ?", "%"+query+"%").Limit(limit).Find(&questions)

	return c.JSON(fiber.Map{
		"results": questions,
		"count":   len(questions),
		"source":  "keyword",
	})
}

// ── Graph queries ─────────────────────────────────────────────────────────────

// FindSimilarProblems godoc
// GET /api/problems/:id/similar?limit=5
func (h *SearchHandler) FindSimilarProblems(c *fiber.Ctx) error {
	problemID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid problem id"})
	}
	limit := c.QueryInt("limit", 5)

	// Try vector similarity first
	if h.vector != nil && h.vector.IsAvailable() {
		// Build query text from problem
		var problem models.Problem
		if h.db.First(&problem, problemID).Error == nil {
			query := problem.Title + " " + problem.Description
			docs, err := h.vector.SearchSimilarProblems(query, limit+1) // +1 because result includes self
			if err == nil && len(docs) > 0 {
				// Filter out self
				filtered := make([]services.SimilarDoc, 0, len(docs))
				selfID := strconv.Itoa(problemID)
				for _, d := range docs {
					if d.ID != "problem_"+selfID {
						filtered = append(filtered, d)
					}
				}
				if len(filtered) > limit {
					filtered = filtered[:limit]
				}
				return c.JSON(fiber.Map{
					"similar_problems": h.enrichProblems(filtered),
					"source":           "vector",
				})
			}
		}
	}

	// Graph fallback
	results, err := h.graph.FindSimilarProblems(problemID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"similar_problems": results,
		"source":           "graph",
	})
}

// GetLearningPath godoc
// GET /api/graph/learning-path?from=1&to=6
func (h *SearchHandler) GetLearningPath(c *fiber.Ctx) error {
	fromID, err1 := strconv.Atoi(c.Query("from"))
	toID, err2 := strconv.Atoi(c.Query("to"))
	if err1 != nil || err2 != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "from and to query params must be integer topic IDs"})
	}

	path, err := h.graph.GetLearningPath(fromID, toID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"learning_path": path,
		"steps":         len(path),
	})
}

// GetTopicPrerequisites godoc
// GET /api/topics/:id/prerequisites
func (h *SearchHandler) GetTopicPrerequisites(c *fiber.Ctx) error {
	topicID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid topic id"})
	}

	prereqs, err := h.graph.GetTopicPrerequisites(topicID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"topic_id":      topicID,
		"prerequisites": prereqs,
		"count":         len(prereqs),
	})
}

// IndexVectors godoc
// POST /api/admin/index (development only)
// Triggers bulk re-indexing of all problems and questions into ChromaDB.
func (h *SearchHandler) IndexVectors(c *fiber.Ctx) error {
	if h.vector == nil || !h.vector.IsAvailable() {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "ChromaDB is not available",
		})
	}

	if err := h.vector.EnsureCollections(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	pCount, err := h.vector.IndexAllProblems(h.db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	qCount, err := h.vector.IndexAllQuestions(h.db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"indexed_problems":  pCount,
		"indexed_questions": qCount,
		"status":            "ok",
	})
}

// SeedGraph godoc
// POST /api/admin/seed-graph (development only)
func (h *SearchHandler) SeedGraph(c *fiber.Ctx) error {
	if err := h.graph.SeedGraph(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"status":    "ok",
		"age_available": h.graph.IsAvailable(),
	})
}

// IntelligenceStatus godoc
// GET /api/intelligence/status
func (h *SearchHandler) IntelligenceStatus(c *fiber.Ctx) error {
	vectorAvail := h.vector != nil && h.vector.IsAvailable()
	graphAvail := h.graph.IsAvailable()

	return c.JSON(fiber.Map{
		"vector_db": fiber.Map{"available": vectorAvail, "provider": "chromadb"},
		"graph_db":  fiber.Map{"available": graphAvail, "provider": "apache_age"},
		"embedding": fiber.Map{"available": vectorAvail, "fallback": "keyword_hash"},
	})
}

// ── Enrichment helpers ────────────────────────────────────────────────────────

type problemResult struct {
	models.Problem
	SimilarityScore float64 `json:"similarity_score,omitempty"`
}

func (h *SearchHandler) enrichProblems(docs []services.SimilarDoc) []problemResult {
	results := make([]problemResult, 0, len(docs))
	for _, doc := range docs {
		var problemID int
		if v, ok := doc.Metadata["problem_id"].(float64); ok {
			problemID = int(v)
		}
		if problemID == 0 {
			continue
		}
		var problem models.Problem
		if err := h.db.First(&problem, problemID).Error; err != nil {
			continue
		}
		results = append(results, problemResult{Problem: problem, SimilarityScore: doc.Score})
	}
	return results
}

type questionResult struct {
	models.Question
	SimilarityScore float64 `json:"similarity_score,omitempty"`
}

func (h *SearchHandler) enrichQuestions(docs []services.SimilarDoc) []questionResult {
	results := make([]questionResult, 0, len(docs))
	for _, doc := range docs {
		var questionID int
		if v, ok := doc.Metadata["question_id"].(float64); ok {
			questionID = int(v)
		}
		if questionID == 0 {
			continue
		}
		var question models.Question
		if err := h.db.First(&question, questionID).Error; err != nil {
			continue
		}
		results = append(results, questionResult{Question: question, SimilarityScore: doc.Score})
	}
	return results
}

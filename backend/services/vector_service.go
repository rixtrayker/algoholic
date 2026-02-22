package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gorm.io/gorm"
	"github.com/yourusername/algoholic/models"
)

// VectorService manages ChromaDB collections and vector search
type VectorService struct {
	chromaURL   string
	httpClient  *http.Client
	embedder    *EmbeddingService
	collections map[string]string // logical name → chroma collection name
}

// NewVectorService creates a new vector service
func NewVectorService(chromaURL string, embedder *EmbeddingService) *VectorService {
	if chromaURL == "" {
		chromaURL = "http://localhost:8000"
	}
	return &VectorService{
		chromaURL:  chromaURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		embedder:   embedder,
		collections: map[string]string{
			"problems":  "algoholic_problems",
			"questions": "algoholic_questions",
		},
	}
}

// ── ChromaDB HTTP types ──────────────────────────────────────────────────────

type chromaCollection struct {
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type chromaAddRequest struct {
	IDs        []string                 `json:"ids"`
	Embeddings [][]float32              `json:"embeddings"`
	Documents  []string                 `json:"documents"`
	Metadatas  []map[string]interface{} `json:"metadatas"`
}

type chromaQueryRequest struct {
	QueryEmbeddings [][]float32 `json:"query_embeddings"`
	NResults        int         `json:"n_results"`
	Include         []string    `json:"include"`
}

type chromaQueryResponse struct {
	IDs       [][]string                 `json:"ids"`
	Distances [][]float64                `json:"distances"`
	Documents [][]string                 `json:"documents"`
	Metadatas [][]map[string]interface{} `json:"metadatas"`
}

// SimilarDoc represents a document returned from semantic search
type SimilarDoc struct {
	ID         string                 `json:"id"`
	Score      float64                `json:"score"`
	Document   string                 `json:"document"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// ── Collection management ────────────────────────────────────────────────────

// EnsureCollections creates ChromaDB collections if they don't exist
func (vs *VectorService) EnsureCollections() error {
	for _, name := range vs.collections {
		if err := vs.ensureCollection(name); err != nil {
			return fmt.Errorf("ensure collection %s: %w", name, err)
		}
	}
	return nil
}

func (vs *VectorService) ensureCollection(name string) error {
	payload, _ := json.Marshal(chromaCollection{Name: name})
	resp, err := vs.httpClient.Post(
		fmt.Sprintf("%s/api/v1/collections", vs.chromaURL),
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// 200 = created, 409 = already exists — both are fine
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusConflict {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

// ── Upsert helpers ────────────────────────────────────────────────────────────

// UpsertProblem embeds and stores a problem in the vector DB
func (vs *VectorService) UpsertProblem(problem *models.Problem) error {
	text := fmt.Sprintf("%s\n\n%s", problem.Title, problem.Description)

	embedding, err := vs.embedder.GenerateWithFallback(text)
	if err != nil {
		return fmt.Errorf("embed problem %d: %w", problem.ProblemID, err)
	}

	difficulty := ""
	if problem.OfficialDifficulty != nil {
		difficulty = *problem.OfficialDifficulty
	}
	pattern := ""
	if problem.PrimaryPattern != nil {
		pattern = *problem.PrimaryPattern
	}

	metadata := map[string]interface{}{
		"problem_id":       problem.ProblemID,
		"title":            problem.Title,
		"slug":             problem.Slug,
		"difficulty_score": problem.DifficultyScore,
		"difficulty_label": difficulty,
		"primary_pattern":  pattern,
	}

	return vs.upsert(
		vs.collections["problems"],
		fmt.Sprintf("problem_%d", problem.ProblemID),
		embedding,
		text,
		metadata,
	)
}

// UpsertQuestion embeds and stores a question in the vector DB
func (vs *VectorService) UpsertQuestion(question *models.Question) error {
	text := question.QuestionText

	embedding, err := vs.embedder.GenerateWithFallback(text)
	if err != nil {
		return fmt.Errorf("embed question %d: %w", question.QuestionID, err)
	}

	problemID := 0
	if question.ProblemID != nil {
		problemID = *question.ProblemID
	}

	metadata := map[string]interface{}{
		"question_id":      question.QuestionID,
		"question_type":    question.QuestionType,
		"question_format":  question.QuestionFormat,
		"difficulty_score": question.DifficultyScore,
		"problem_id":       problemID,
	}

	return vs.upsert(
		vs.collections["questions"],
		fmt.Sprintf("question_%d", question.QuestionID),
		embedding,
		text,
		metadata,
	)
}

func (vs *VectorService) upsert(collection, id string, embedding []float32, doc string, metadata map[string]interface{}) error {
	req := chromaAddRequest{
		IDs:        []string{id},
		Embeddings: [][]float32{embedding},
		Documents:  []string{doc},
		Metadatas:  []map[string]interface{}{metadata},
	}
	payload, _ := json.Marshal(req)

	// Use upsert endpoint
	url := fmt.Sprintf("%s/api/v1/collections/%s/upsert", vs.chromaURL, collection)
	resp, err := vs.httpClient.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upsert failed (status %d): %s", resp.StatusCode, string(body))
	}
	return nil
}

// ── Search ───────────────────────────────────────────────────────────────────

// SearchSimilarProblems returns problems semantically similar to the query
func (vs *VectorService) SearchSimilarProblems(query string, limit int) ([]SimilarDoc, error) {
	return vs.search(vs.collections["problems"], query, limit)
}

// SearchSimilarQuestions returns questions semantically similar to the query
func (vs *VectorService) SearchSimilarQuestions(query string, limit int) ([]SimilarDoc, error) {
	return vs.search(vs.collections["questions"], query, limit)
}

func (vs *VectorService) search(collection, query string, limit int) ([]SimilarDoc, error) {
	embedding, err := vs.embedder.GenerateWithFallback(query)
	if err != nil {
		return nil, fmt.Errorf("embed query: %w", err)
	}

	req := chromaQueryRequest{
		QueryEmbeddings: [][]float32{embedding},
		NResults:        limit,
		Include:         []string{"distances", "documents", "metadatas"},
	}
	payload, _ := json.Marshal(req)

	url := fmt.Sprintf("%s/api/v1/collections/%s/query", vs.chromaURL, collection)
	resp, err := vs.httpClient.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search failed (status %d): %s", resp.StatusCode, string(body))
	}

	var chromaResp chromaQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&chromaResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if len(chromaResp.IDs) == 0 {
		return []SimilarDoc{}, nil
	}

	docs := make([]SimilarDoc, 0, len(chromaResp.IDs[0]))
	for i, id := range chromaResp.IDs[0] {
		// ChromaDB returns L2 distance; convert to similarity score 0-1
		distance := 0.0
		if i < len(chromaResp.Distances[0]) {
			distance = chromaResp.Distances[0][i]
		}
		score := 1.0 / (1.0 + distance) // bounded 0-1, higher = more similar

		doc := SimilarDoc{
			ID:    id,
			Score: score,
		}
		if i < len(chromaResp.Documents[0]) {
			doc.Document = chromaResp.Documents[0][i]
		}
		if i < len(chromaResp.Metadatas[0]) {
			doc.Metadata = chromaResp.Metadatas[0][i]
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

// ── Bulk indexing ─────────────────────────────────────────────────────────────

// IndexAllProblems embeds all problems in the database into ChromaDB
func (vs *VectorService) IndexAllProblems(db *gorm.DB) (int, error) {
	var problems []models.Problem
	if err := db.Find(&problems).Error; err != nil {
		return 0, err
	}

	count := 0
	for _, p := range problems {
		if err := vs.UpsertProblem(&p); err != nil {
			// Log but continue
			fmt.Printf("Warning: failed to index problem %d: %v\n", p.ProblemID, err)
			continue
		}
		count++
	}
	return count, nil
}

// IndexAllQuestions embeds all questions in the database into ChromaDB
func (vs *VectorService) IndexAllQuestions(db *gorm.DB) (int, error) {
	var questions []models.Question
	if err := db.Find(&questions).Error; err != nil {
		return 0, err
	}

	count := 0
	for _, q := range questions {
		if err := vs.UpsertQuestion(&q); err != nil {
			fmt.Printf("Warning: failed to index question %d: %v\n", q.QuestionID, err)
			continue
		}
		count++
	}
	return count, nil
}

// IsAvailable checks if ChromaDB is reachable
func (vs *VectorService) IsAvailable() bool {
	resp, err := vs.httpClient.Get(fmt.Sprintf("%s/api/v1/heartbeat", vs.chromaURL))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

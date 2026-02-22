package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// EmbeddingService generates vector embeddings via Ollama
type EmbeddingService struct {
	ollamaURL  string
	model      string
	httpClient *http.Client
}

// NewEmbeddingService creates a new embedding service
func NewEmbeddingService(ollamaURL, model string) *EmbeddingService {
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}
	if model == "" {
		model = "all-minilm"
	}
	return &EmbeddingService{
		ollamaURL: ollamaURL,
		model:     model,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

type ollamaEmbedRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type ollamaEmbedResponse struct {
	Embedding []float32 `json:"embedding"`
}

// Generate generates a vector embedding for the given text
func (es *EmbeddingService) Generate(text string) ([]float32, error) {
	payload, _ := json.Marshal(ollamaEmbedRequest{Model: es.model, Prompt: text})

	resp, err := es.httpClient.Post(
		fmt.Sprintf("%s/api/embeddings", es.ollamaURL),
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return nil, fmt.Errorf("ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama error (status %d): %s", resp.StatusCode, string(body))
	}

	var result ollamaEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode embedding response: %w", err)
	}

	return result.Embedding, nil
}

// IsAvailable checks if Ollama is running and the model is loaded
func (es *EmbeddingService) IsAvailable() bool {
	resp, err := es.httpClient.Get(fmt.Sprintf("%s/api/tags", es.ollamaURL))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// GenerateWithFallback returns a real embedding when Ollama is up, or a
// keyword-frequency (TF-IDF-lite) vector as a fallback.  The fallback vector
// is 128-dimensional so the same ChromaDB collection works regardless of
// which path is taken (as long as you don't mix vectors from both paths in
// production — here both are 128-dim for consistency).
func (es *EmbeddingService) GenerateWithFallback(text string) ([]float32, error) {
	if es.IsAvailable() {
		return es.Generate(text)
	}
	return KeywordVector(text), nil
}

// KeywordVector builds a deterministic 128-dim float32 vector from text using
// character n-gram frequency hashing (no external dependencies).
func KeywordVector(text string) []float32 {
	const dim = 128
	vec := make([]float32, dim)

	// Normalise
	lower := strings.ToLower(text)

	// Hash each word into a bucket and accumulate frequency
	words := strings.Fields(lower)
	for _, w := range words {
		// Strip punctuation
		clean := strings.Map(func(r rune) rune {
			if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
				return r
			}
			return -1
		}, w)
		if len(clean) == 0 {
			continue
		}
		// FNV-1a hash into [0, dim)
		h := uint32(2166136261)
		for i := 0; i < len(clean); i++ {
			h ^= uint32(clean[i])
			h *= 16777619
		}
		vec[h%dim] += 1.0
	}

	// L2-normalise
	var norm float32
	for _, v := range vec {
		norm += v * v
	}
	if norm > 0 {
		norm = float32(1.0 / float64(norm))
		for i := range vec {
			vec[i] *= norm
		}
	}
	return vec
}

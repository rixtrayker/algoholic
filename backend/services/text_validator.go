package services

import (
	"regexp"
	"strings"
	"unicode"
)

// TextValidator handles fuzzy text matching for answer validation
type TextValidator struct {
	threshold float64
}

// NewTextValidator creates a new text validator with default threshold
func NewTextValidator() *TextValidator {
	return &TextValidator{
		threshold: 0.85, // 85% similarity required
	}
}

// FuzzyMatch checks if two texts match within the threshold
func (tv *TextValidator) FuzzyMatch(userText, correctText string) bool {
	// Normalize inputs
	t1 := tv.NormalizeText(userText)
	t2 := tv.NormalizeText(correctText)

	// Strategy 1: Exact match after normalization
	if t1 == t2 {
		return true
	}

	// Strategy 2: Levenshtein distance similarity
	similarity := tv.CalculateSimilarity(t1, t2)
	if similarity >= tv.threshold {
		return true
	}

	// Strategy 3: Keyword matching (for conceptual answers)
	if tv.HasRequiredKeywords(t1, t2) {
		return true
	}

	return false
}

// MatchMultiple checks if user text matches any of the acceptable answers
func (tv *TextValidator) MatchMultiple(userText string, acceptableAnswers []string) bool {
	for _, answer := range acceptableAnswers {
		if tv.FuzzyMatch(userText, answer) {
			return true
		}
	}
	return false
}

// NormalizeText prepares text for comparison
func (tv *TextValidator) NormalizeText(text string) string {
	// Convert to lowercase
	text = strings.ToLower(text)

	// Remove extra whitespace
	text = strings.TrimSpace(text)
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	// Remove most punctuation (except dashes and underscores in technical terms)
	text = regexp.MustCompile(`[^\w\s\-()]`).ReplaceAllString(text, "")

	// Common complexity notation normalization
	replacements := map[string]string{
		"o(n)":         "on",
		"o(log n)":     "ologn",
		"o(n log n)":   "onlogn",
		"o(n^2)":       "on2",
		"o(n2)":        "on2",
		"o(1)":         "o1",
		"o(2^n)":       "o2n",
		"log(n)":       "logn",
		"sqrt(n)":      "sqrtn",
		"n^2":          "n2",
		// Common abbreviations
		"binary search": "binarysearch",
		"two pointers":  "twopointers",
		"hash map":      "hashmap",
		"hash table":    "hashtable",
		"linked list":   "linkedlist",
		"binary tree":   "binarytree",
		"bst":           "binarysearchtree",
	}

	for old, new := range replacements {
		text = strings.ReplaceAll(text, old, new)
	}

	return text
}

// CalculateSimilarity computes similarity score using Levenshtein distance
func (tv *TextValidator) CalculateSimilarity(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}

	distance := tv.LevenshteinDistance(s1, s2)
	maxLen := max(len(s1), len(s2))

	if maxLen == 0 {
		return 0.0
	}

	similarity := 1.0 - float64(distance)/float64(maxLen)
	return similarity
}

// LevenshteinDistance calculates the Levenshtein distance between two strings
func (tv *TextValidator) LevenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	// Create a 2D slice for dynamic programming
	d := make([][]int, len(s1)+1)
	for i := range d {
		d[i] = make([]int, len(s2)+1)
	}

	// Initialize first column and row
	for i := 0; i <= len(s1); i++ {
		d[i][0] = i
	}
	for j := 0; j <= len(s2); j++ {
		d[0][j] = j
	}

	// Calculate distances
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			d[i][j] = min(
				d[i-1][j]+1,      // deletion
				d[i][j-1]+1,      // insertion
				d[i-1][j-1]+cost, // substitution
			)
		}
	}

	return d[len(s1)][len(s2)]
}

// HasRequiredKeywords checks if user answer contains key technical terms
func (tv *TextValidator) HasRequiredKeywords(userText, correctText string) bool {
	keywords := tv.ExtractKeywords(correctText)

	if len(keywords) == 0 {
		return false // No keywords to match
	}

	// User must mention at least 70% of keywords
	matchCount := 0
	for _, keyword := range keywords {
		if strings.Contains(userText, keyword) {
			matchCount++
		}
	}

	threshold := 0.7
	return float64(matchCount)/float64(len(keywords)) >= threshold
}

// ExtractKeywords extracts technical terms from text
func (tv *TextValidator) ExtractKeywords(text string) []string {
	text = strings.ToLower(text)

	// Common technical terms, algorithm names, data structures
	technicalTerms := []string{
		// Data structures
		"array", "hashmap", "hashtable", "linkedlist", "stack", "queue",
		"heap", "priorityqueue", "tree", "binarytree", "binarysearchtree",
		"trie", "graph", "set", "map", "deque",

		// Algorithms
		"binarysearch", "dfs", "bfs", "dynamicprogramming", "greedy",
		"backtracking", "divideandconquer", "recursion", "iteration",
		"sorting", "searching", "mergesort", "quicksort",

		// Patterns
		"twopointers", "slidingwindow", "prefixsum", "unionfind",
		"topologicalsort", "dijkstra", "bellmanford", "floydwarshall",
		"knapsack", "kadane", "monotonically", "monotonicstack",

		// Complexity
		"on", "ologn", "onlogn", "on2", "o1", "constant", "linear",
		"logarithmic", "quadratic", "exponential",

		// Concepts
		"memoization", "tabulation", "dp", "optimal", "subproblem",
		"overlapping", "recursion", "base case", "induction",
	}

	keywords := []string{}

	for _, term := range technicalTerms {
		if strings.Contains(text, term) {
			keywords = append(keywords, term)
		}
	}

	return keywords
}

// ValidateComplexityAnswer specifically validates time/space complexity answers
func (tv *TextValidator) ValidateComplexityAnswer(userAnswer, correctAnswer string) bool {
	// Normalize both answers
	user := tv.NormalizeText(userAnswer)
	correct := tv.NormalizeText(correctAnswer)

	// Extract complexity notation
	userComplexity := tv.ExtractComplexity(user)
	correctComplexity := tv.ExtractComplexity(correct)

	return userComplexity == correctComplexity
}

// ExtractComplexity extracts O() notation from text
func (tv *TextValidator) ExtractComplexity(text string) string {
	// Pattern to match O(...)
	re := regexp.MustCompile(`o\s*\(\s*([^)]+)\s*\)`)
	matches := re.FindStringSubmatch(text)

	if len(matches) > 1 {
		complexity := matches[1]
		// Normalize the complexity
		complexity = strings.ReplaceAll(complexity, " ", "")
		return complexity
	}

	// Try to find common complexity terms
	complexities := []string{"on2", "onlogn", "ologn", "on", "o1", "constant", "linear", "logarithmic", "quadratic"}
	for _, c := range complexities {
		if strings.Contains(text, c) {
			return c
		}
	}

	return text
}

// Helper functions
func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// RemovePunctuation removes all punctuation from text
func (tv *TextValidator) RemovePunctuation(text string) string {
	var result strings.Builder
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// PartialMatch checks if key phrases from correct answer appear in user answer
func (tv *TextValidator) PartialMatch(userText, correctText string) bool {
	// Split correct answer into phrases
	correctPhrases := strings.Fields(tv.NormalizeText(correctText))

	if len(correctPhrases) == 0 {
		return false
	}

	userNorm := tv.NormalizeText(userText)

	// Count how many phrases match
	matchCount := 0
	for _, phrase := range correctPhrases {
		if len(phrase) > 2 && strings.Contains(userNorm, phrase) {
			matchCount++
		}
	}

	// At least 60% of phrases should match
	threshold := 0.6
	return float64(matchCount)/float64(len(correctPhrases)) >= threshold
}

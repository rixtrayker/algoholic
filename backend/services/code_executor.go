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

// CodeExecutor handles code execution via Judge0 API
type CodeExecutor struct {
	judge0URL string
	timeout   time.Duration
}

// TestCase represents a single test case
type TestCase struct {
	Input    string `json:"input"`
	Expected string `json:"expected"`
}

// ExecutionResult contains the results of code execution
type ExecutionResult struct {
	AllPassed   bool            `json:"all_passed"`
	PassedCount int             `json:"passed_count"`
	TotalCount  int             `json:"total_count"`
	Failures    []FailureDetail `json:"failures,omitempty"`
	TimeTaken   float64         `json:"time_taken_ms"`
	MemoryUsed  int             `json:"memory_used_kb"`
}

// FailureDetail contains information about a failed test case
type FailureDetail struct {
	TestNumber int    `json:"test_number"`
	Input      string `json:"input"`
	Expected   string `json:"expected"`
	Got        string `json:"got"`
	Error      string `json:"error,omitempty"`
}

// Judge0Submission represents a submission to Judge0
type Judge0Submission struct {
	SourceCode       string  `json:"source_code"`
	LanguageID       int     `json:"language_id"`
	Stdin            string  `json:"stdin,omitempty"`
	ExpectedOutput   string  `json:"expected_output,omitempty"`
	CPUTimeLimit     float64 `json:"cpu_time_limit,omitempty"`
	MemoryLimit      int     `json:"memory_limit,omitempty"`
	WallTimeLimit    float64 `json:"wall_time_limit,omitempty"`
}

// Judge0Response represents Judge0 API response
type Judge0Response struct {
	Token  string `json:"token"`
	Status struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
	} `json:"status"`
	Stdout     string  `json:"stdout"`
	Stderr     string  `json:"stderr"`
	CompileOutput string `json:"compile_output"`
	Time       string  `json:"time"`
	Memory     int     `json:"memory"`
}

// NewCodeExecutor creates a new code executor
func NewCodeExecutor(judge0URL string) *CodeExecutor {
	if judge0URL == "" {
		judge0URL = "http://localhost:2358" // Default Judge0 CE URL
	}

	return &CodeExecutor{
		judge0URL: judge0URL,
		timeout:   30 * time.Second,
	}
}

// RunTests executes code against test cases
func (ce *CodeExecutor) RunTests(code, language string, testCases []interface{}) (*ExecutionResult, error) {
	languageID := ce.getLanguageID(language)
	if languageID == 0 {
		return nil, fmt.Errorf("unsupported language: %s", language)
	}

	result := &ExecutionResult{
		AllPassed:   true,
		PassedCount: 0,
		TotalCount:  len(testCases),
		Failures:    []FailureDetail{},
		TimeTaken:   0,
		MemoryUsed:  0,
	}

	for i, tc := range testCases {
		testCase, ok := tc.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid test case format at index %d", i)
		}

		input, _ := testCase["input"].(string)
		expected, _ := testCase["expected"].(string)

		// Execute code for this test case
		output, execTime, memory, err := ce.executeCode(code, languageID, input, expected)

		result.TimeTaken += execTime
		if memory > result.MemoryUsed {
			result.MemoryUsed = memory
		}

		if err != nil {
			result.AllPassed = false
			result.Failures = append(result.Failures, FailureDetail{
				TestNumber: i + 1,
				Input:      input,
				Expected:   expected,
				Got:        output,
				Error:      err.Error(),
			})
			continue
		}

		// Compare output (trim whitespace)
		if ce.normalizeOutput(output) == ce.normalizeOutput(expected) {
			result.PassedCount++
		} else {
			result.AllPassed = false
			result.Failures = append(result.Failures, FailureDetail{
				TestNumber: i + 1,
				Input:      input,
				Expected:   expected,
				Got:        output,
			})
		}
	}

	return result, nil
}

// executeCode submits code to Judge0 and waits for result
func (ce *CodeExecutor) executeCode(code string, languageID int, stdin, expectedOutput string) (string, float64, int, error) {
	// Create submission
	submission := Judge0Submission{
		SourceCode:     code,
		LanguageID:     languageID,
		Stdin:          stdin,
		ExpectedOutput: expectedOutput,
		CPUTimeLimit:   5.0,  // 5 seconds max
		MemoryLimit:    128000, // 128 MB
		WallTimeLimit:  10.0, // 10 seconds max wall time
	}

	jsonData, err := json.Marshal(submission)
	if err != nil {
		return "", 0, 0, err
	}

	// Submit to Judge0
	url := fmt.Sprintf("%s/submissions?wait=true", ce.judge0URL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", 0, 0, fmt.Errorf("judge0 request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", 0, 0, fmt.Errorf("judge0 error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var judge0Resp Judge0Response
	if err := json.NewDecoder(resp.Body).Decode(&judge0Resp); err != nil {
		return "", 0, 0, fmt.Errorf("failed to parse judge0 response: %w", err)
	}

	// Check execution status
	if judge0Resp.Status.ID != 3 { // 3 = Accepted
		errMsg := judge0Resp.Status.Description
		if judge0Resp.Stderr != "" {
			errMsg = judge0Resp.Stderr
		} else if judge0Resp.CompileOutput != "" {
			errMsg = judge0Resp.CompileOutput
		}
		return "", 0, 0, fmt.Errorf("execution error: %s", errMsg)
	}

	// Parse time (comes as string like "0.001")
	var execTime float64
	if judge0Resp.Time != "" {
		fmt.Sscanf(judge0Resp.Time, "%f", &execTime)
		execTime *= 1000 // Convert to milliseconds
	}

	return judge0Resp.Stdout, execTime, judge0Resp.Memory, nil
}

// getLanguageID maps language names to Judge0 language IDs
func (ce *CodeExecutor) getLanguageID(language string) int {
	languageMap := map[string]int{
		"python":     71, // Python 3.8.1
		"python3":    71,
		"javascript": 63, // JavaScript (Node.js 12.14.0)
		"js":         63,
		"java":       62, // Java (OpenJDK 13.0.1)
		"cpp":        54, // C++ (GCC 9.2.0)
		"c++":        54,
		"c":          50, // C (GCC 9.2.0)
		"go":         60, // Go (1.13.5)
		"rust":       73, // Rust (1.40.0)
		"ruby":       72, // Ruby (2.7.0)
		"php":        68, // PHP (7.4.1)
		"swift":      83, // Swift (5.2.3)
		"kotlin":     78, // Kotlin (1.3.70)
		"typescript": 74, // TypeScript (3.7.4)
		"ts":         74,
	}

	if id, ok := languageMap[language]; ok {
		return id
	}

	return 0 // Unsupported language
}

// normalizeOutput trims and normalizes output for comparison
func (ce *CodeExecutor) normalizeOutput(output string) string {
	// Trim leading/trailing whitespace
	output = strings.TrimSpace(output)

	// Normalize line endings
	output = strings.ReplaceAll(output, "\r\n", "\n")

	return output
}

// ValidateCode provides a simple validation without execution (for fallback)
func (ce *CodeExecutor) ValidateCode(code, language string) bool {
	// Basic validation: code is not empty and has reasonable length
	if len(code) < 10 || len(code) > 50000 {
		return false
	}

	// Language-specific basic checks
	switch language {
	case "python", "python3":
		// Check for basic Python structure
		return bytes.Contains([]byte(code), []byte("def ")) ||
		       bytes.Contains([]byte(code), []byte("class "))
	case "javascript", "js":
		return bytes.Contains([]byte(code), []byte("function ")) ||
		       bytes.Contains([]byte(code), []byte("const ")) ||
		       bytes.Contains([]byte(code), []byte("=>"))
	case "java":
		return bytes.Contains([]byte(code), []byte("class ")) &&
		       bytes.Contains([]byte(code), []byte("public "))
	case "cpp", "c++":
		return bytes.Contains([]byte(code), []byte("#include")) ||
		       bytes.Contains([]byte(code), []byte("int main"))
	default:
		return true // Accept other languages without validation
	}
}

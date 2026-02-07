package data

import "dsa-platform/pkg/models"

// SeedQuestions returns the initial bank of assessment questions.
// Each question references a question_type slug and optionally a problem slug.
func SeedQuestions() []models.Question {
	return []models.Question{

		// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
		// CATEGORY 1: COMPLEXITY ANALYSIS
		// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

		// --- 1.1 Code-to-Complexity ---
		{
			QuestionTypeSlug: "code-to-complexity", Category: "complexity_analysis", Subcategory: "1.1",
			QuestionText: "What is the time complexity of the following code?\n\n```go\nfor i := 0; i < n; i++ {\n    for j := i; j < n; j++ {\n        // O(1) work\n    }\n}\n```",
			QuestionData: models.JSONB{"language": "go", "code_snippet": "for i := 0; i < n; i++ {\n    for j := i; j < n; j++ {\n        // O(1) work\n    }\n}"},
			Format:       models.FmtMultipleChoice,
			CorrectAnswer: models.JSONB{"option_id": "B", "value": "O(n²)"},
			AnswerOptions: []models.AnswerOption{
				{ID: "A", Text: "O(n)", IsCorrect: false},
				{ID: "B", Text: "O(n²)", IsCorrect: true},
				{ID: "C", Text: "O(n log n)", IsCorrect: false},
				{ID: "D", Text: "O(n³)", IsCorrect: false},
			},
			WrongAnswerExplanations: models.JSONB{
				"A": "The inner loop still runs n*(n+1)/2 total iterations, not n.",
				"C": "n log n would require the inner loop to shrink logarithmically, but j goes from i to n.",
				"D": "There are only two nested loops, not three.",
			},
			Explanation:    "The inner loop runs n, n-1, n-2, ... , 1 times. Total = n(n+1)/2 = O(n²). The fact that j starts at i (not 0) halves the constant but doesn't change the asymptotic class.",
			CommonMistakes: []string{"Thinking j=i makes it O(n) because 'it skips half'", "Confusing n(n+1)/2 with O(n)"},
			HintLevel1:     "How many total iterations does the inner loop execute across all values of i?",
			HintLevel2:     "Sum the series: n + (n-1) + (n-2) + ... + 1",
			HintLevel3:     "That's the triangular number n(n+1)/2 which simplifies to O(n²)",
			DifficultyScore: 18, DifficultyLabel: "easy", EstimatedTimeSec: 90,
			RelatedTopicID: "arrays_basics",
			Tags:     []string{"complexity", "nested-loops", "triangular-sum"},
			Concepts: []string{"time-complexity", "summation", "nested-iteration"},
		},

		{
			QuestionTypeSlug: "code-to-complexity", Category: "complexity_analysis", Subcategory: "1.1",
			QuestionText: "What is the time complexity?\n\n```go\nfunc solve(n int) int {\n    if n <= 1 { return 1 }\n    return solve(n-1) + solve(n-1)\n}\n```",
			QuestionData: models.JSONB{"language": "go", "code_snippet": "func solve(n int) int {\n    if n <= 1 { return 1 }\n    return solve(n-1) + solve(n-1)\n}"},
			Format:       models.FmtMultipleChoice,
			CorrectAnswer: models.JSONB{"option_id": "D", "value": "O(2ⁿ)"},
			AnswerOptions: []models.AnswerOption{
				{ID: "A", Text: "O(n)", IsCorrect: false},
				{ID: "B", Text: "O(n²)", IsCorrect: false},
				{ID: "C", Text: "O(n log n)", IsCorrect: false},
				{ID: "D", Text: "O(2ⁿ)", IsCorrect: true},
			},
			Explanation:    "Each call spawns 2 recursive calls. The recursion tree has depth n and each level doubles the calls: 1 + 2 + 4 + ... + 2ⁿ = O(2ⁿ). This is the classic exponential blowup without memoization.",
			CommonMistakes: []string{"Thinking it's O(n) because each call does O(1) work", "Confusing with binary search O(log n)"},
			HintLevel1:     "Draw the recursion tree. How many nodes are at each level?",
			HintLevel2:     "Level 0 has 1 call, level 1 has 2, level 2 has 4...",
			HintLevel3:     "Total nodes = 2⁰ + 2¹ + ... + 2ⁿ = 2ⁿ⁺¹ - 1 = O(2ⁿ)",
			DifficultyScore: 25, DifficultyLabel: "medium", EstimatedTimeSec: 120,
			Tags:     []string{"complexity", "recursion", "exponential"},
			Concepts: []string{"time-complexity", "recursion-tree", "branching-factor"},
		},

		// --- 1.3 Constraint-to-Complexity ---
		{
			QuestionTypeSlug: "constraint-to-complexity", Category: "complexity_analysis", Subcategory: "1.3",
			QuestionText: "A problem has constraint n ≤ 10⁵ and time limit of 1 second. What is the maximum acceptable time complexity for your solution?",
			Format:       models.FmtMultipleChoice,
			CorrectAnswer: models.JSONB{"option_id": "C", "value": "O(n log n)"},
			AnswerOptions: []models.AnswerOption{
				{ID: "A", Text: "O(n³)", IsCorrect: false},
				{ID: "B", Text: "O(n²)", IsCorrect: false},
				{ID: "C", Text: "O(n log n)", IsCorrect: true},
				{ID: "D", Text: "O(n√n)", IsCorrect: false},
			},
			Explanation:    "With n = 10⁵ and ~10⁸ operations per second: O(n²) = 10¹⁰ (too slow), O(n log n) ≈ 1.7×10⁶ (safe), O(n√n) ≈ 3.2×10⁷ (borderline). O(n log n) is the sweet spot for n = 10⁵.",
			CommonMistakes: []string{"Choosing O(n²) — 10¹⁰ operations won't finish in 1s", "Not knowing the ~10⁸ ops/sec rule of thumb"},
			HintLevel1:     "How many operations can a modern CPU do in 1 second?",
			HintLevel2:     "Roughly 10⁸ simple operations per second. Plug in n = 10⁵.",
			HintLevel3:     "n² = 10¹⁰ (too slow), n·log(n) ≈ 1.7M (very fast), n·√n ≈ 31M (borderline OK)",
			DifficultyScore: 22, DifficultyLabel: "medium", EstimatedTimeSec: 90,
			Tags:     []string{"complexity", "constraints", "competitive-programming"},
			Concepts: []string{"constraint-analysis", "operations-per-second"},
		},

		// --- 1.4 Hidden Complexity ---
		{
			QuestionTypeSlug: "hidden-complexity", Category: "complexity_analysis", Subcategory: "1.4",
			QuestionText: "What is the ACTUAL time complexity of this Go code?\n\n```go\nresult := \"\"\nfor i := 0; i < n; i++ {\n    result += string(data[i]) // string concatenation\n}\n```",
			QuestionData: models.JSONB{"language": "go"},
			Format:       models.FmtMultipleChoice,
			CorrectAnswer: models.JSONB{"option_id": "B", "value": "O(n²)"},
			AnswerOptions: []models.AnswerOption{
				{ID: "A", Text: "O(n)", IsCorrect: false},
				{ID: "B", Text: "O(n²)", IsCorrect: true},
				{ID: "C", Text: "O(n log n)", IsCorrect: false},
				{ID: "D", Text: "O(1)", IsCorrect: false},
			},
			Explanation:    "In Go, strings are immutable. Each concatenation creates a new string of increasing length: copy 1 + 2 + 3 + ... + n = O(n²). Use strings.Builder instead for O(n).",
			CommonMistakes: []string{"Assuming string concat is O(1) per operation", "Not knowing Go strings are immutable"},
			HintLevel1:     "Are Go strings mutable or immutable?",
			HintLevel2:     "Each += creates a NEW string, copying all previous characters plus the new one.",
			HintLevel3:     "Total copies: 1 + 2 + 3 + ... + n = n(n+1)/2 = O(n²). Use strings.Builder for O(n).",
			DifficultyScore: 35, DifficultyLabel: "medium", EstimatedTimeSec: 120,
			Tags:     []string{"complexity", "hidden-cost", "golang", "string"},
			Concepts: []string{"hidden-complexity", "immutable-strings", "string-builder"},
		},

		// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
		// CATEGORY 2: DATA STRUCTURE SELECTION
		// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

		{
			QuestionTypeSlug: "requirements-to-ds", Category: "ds_selection", Subcategory: "2.1",
			QuestionText: "You need a data structure that supports:\n• Insert element: O(log n)\n• Delete element: O(log n)\n• Find minimum: O(1)\n• Find maximum: O(1)\n\nWhich data structure best fits?",
			Format:       models.FmtMultipleChoice,
			CorrectAnswer: models.JSONB{"option_id": "C", "value": "Balanced BST (e.g. Red-Black Tree / TreeMap)"},
			AnswerOptions: []models.AnswerOption{
				{ID: "A", Text: "Min-Heap", IsCorrect: false},
				{ID: "B", Text: "Hash Map", IsCorrect: false},
				{ID: "C", Text: "Balanced BST (e.g. Red-Black Tree / TreeMap)", IsCorrect: true},
				{ID: "D", Text: "Sorted Array", IsCorrect: false},
			},
			WrongAnswerExplanations: models.JSONB{
				"A": "Min-heap gives O(1) min but O(n) max, and delete arbitrary element is O(n).",
				"B": "Hash map has O(1) insert/delete but O(n) for min and max.",
				"D": "Sorted array has O(1) min/max but O(n) insert/delete due to shifting.",
			},
			Explanation:    "A balanced BST (like Go's btree or a sorted set) keeps elements ordered. Min = leftmost node O(log n) amortized or cached O(1), max = rightmost. Insert and delete are O(log n). Only a balanced BST satisfies all four requirements simultaneously.",
			HintLevel1:     "Which structures maintain sorted order while allowing efficient insertion?",
			HintLevel2:     "A heap is great for one extreme but not both. What about a sorted structure?",
			HintLevel3:     "A balanced BST has O(log n) insert/delete and the leftmost/rightmost nodes give min/max.",
			DifficultyScore: 28, DifficultyLabel: "medium", EstimatedTimeSec: 120,
			Tags:     []string{"ds-selection", "bst", "heap", "tradeoff"},
			Concepts: []string{"data-structure-selection", "operation-complexity", "balanced-bst"},
		},

		{
			QuestionTypeSlug: "ds-tradeoff", Category: "ds_selection", Subcategory: "2.2",
			QuestionText: "For the 'Sliding Window Maximum' problem (max of each window of size k), why is a deque preferred over a heap?\n\nA) Deque uses less memory\nB) Deque gives O(n) total while heap gives O(n log k)\nC) Heap can't track window boundaries\nD) Both B and C",
			Format:       models.FmtMultipleChoice,
			CorrectAnswer: models.JSONB{"option_id": "D", "value": "Both B and C"},
			AnswerOptions: []models.AnswerOption{
				{ID: "A", Text: "Deque uses less memory", IsCorrect: false},
				{ID: "B", Text: "Deque gives O(n) total while heap gives O(n log k)", IsCorrect: false},
				{ID: "C", Text: "Heap can't easily remove elements leaving the window", IsCorrect: false},
				{ID: "D", Text: "Both B and C", IsCorrect: true},
			},
			Explanation:    "A monotonic deque processes each element at most twice (push+pop) giving O(n) total. A max-heap requires O(log k) per insertion AND has trouble removing stale elements that leave the window (lazy deletion needed). The deque naturally evicts stale elements from the front.",
			DifficultyScore: 38, DifficultyLabel: "medium", EstimatedTimeSec: 150,
			RelatedTopicID: "monotonic_queue",
			Tags:     []string{"ds-selection", "deque", "heap", "monotonic-queue", "sliding-window"},
			Concepts: []string{"monotonic-deque", "amortized-analysis", "lazy-deletion"},
		},

		// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
		// CATEGORY 3: PATTERN RECOGNITION
		// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

		{
			QuestionTypeSlug: "problem-to-pattern", Category: "pattern_recognition", Subcategory: "3.1",
			QuestionText: "Problem: \"Given a string, find the length of the longest substring with at most 2 distinct characters.\"\n\nWhat pattern does this problem use?",
			Format:       models.FmtMultipleChoice,
			CorrectAnswer: models.JSONB{"option_id": "B", "value": "Sliding Window"},
			AnswerOptions: []models.AnswerOption{
				{ID: "A", Text: "Two Pointers (Opposite Direction)", IsCorrect: false},
				{ID: "B", Text: "Sliding Window", IsCorrect: true},
				{ID: "C", Text: "Binary Search", IsCorrect: false},
				{ID: "D", Text: "Dynamic Programming", IsCorrect: false},
			},
			Explanation:    "Key indicators: 'longest substring' + 'at most K' constraint. This is the variable sliding window pattern: expand right to grow window, shrink left when constraint is violated (>2 distinct chars). Track answer as max window size.",
			HintLevel1:     "What does 'longest substring with at most K constraint' usually map to?",
			HintLevel2:     "Think about maintaining a window that satisfies the constraint...",
			HintLevel3:     "Sliding window: expand right, shrink left when >2 distinct chars, track max length.",
			DifficultyScore: 15, DifficultyLabel: "easy", EstimatedTimeSec: 60,
			RelatedTopicID: "sliding_window",
			Tags:     []string{"pattern-recognition", "sliding-window", "substring"},
			Concepts: []string{"pattern-mapping", "keyword-triggers"},
		},

		{
			QuestionTypeSlug: "keyword-to-pattern", Category: "pattern_recognition", Subcategory: "3.3",
			QuestionText: "Match each problem keyword to its most likely algorithmic pattern:\n\n1. \"next greater element\" → ?\n2. \"shortest path in unweighted graph\" → ?\n3. \"all possible combinations\" → ?\n4. \"minimum number of coins\" → ?\n\nWhich mapping is correct?",
			Format:       models.FmtMultipleChoice,
			CorrectAnswer: models.JSONB{"option_id": "A"},
			AnswerOptions: []models.AnswerOption{
				{ID: "A", Text: "1→Monotonic Stack, 2→BFS, 3→Backtracking, 4→DP", IsCorrect: true},
				{ID: "B", Text: "1→Binary Search, 2→DFS, 3→Backtracking, 4→Greedy", IsCorrect: false},
				{ID: "C", Text: "1→Monotonic Stack, 2→Dijkstra, 3→DP, 4→DP", IsCorrect: false},
				{ID: "D", Text: "1→Stack, 2→BFS, 3→DFS, 4→BFS", IsCorrect: false},
			},
			Explanation:    "'Next greater' → monotonic stack (classic usage). 'Shortest path unweighted' → BFS (not Dijkstra which is for weighted). 'All combinations' → backtracking (generate all). 'Minimum coins' → DP (optimization over overlapping subproblems).",
			DifficultyScore: 20, DifficultyLabel: "easy", EstimatedTimeSec: 90,
			Tags:     []string{"pattern-recognition", "keyword-mapping"},
			Concepts: []string{"keyword-triggers", "pattern-classification"},
		},

		{
			QuestionTypeSlug: "problem-to-pattern", Category: "pattern_recognition", Subcategory: "3.1",
			QuestionText: "Problem: \"Given an array, for each element, find the next greater element to its right. If none exists, output -1.\"\n\nWhat pattern?",
			Format:       models.FmtMultipleChoice,
			CorrectAnswer: models.JSONB{"option_id": "C", "value": "Monotonic Stack"},
			AnswerOptions: []models.AnswerOption{
				{ID: "A", Text: "Sliding Window", IsCorrect: false},
				{ID: "B", Text: "Two Pointers", IsCorrect: false},
				{ID: "C", Text: "Monotonic Stack", IsCorrect: true},
				{ID: "D", Text: "Binary Search", IsCorrect: false},
			},
			Explanation:    "\"Next greater element\" is the textbook monotonic stack trigger. Maintain a decreasing stack; when a new element is larger than stack top, it's the 'next greater' for that top element. Process from right to left (or left to right with different logic).",
			RelatedProblemSlug: "largest-rectangle-in-histogram",
			RelatedTopicID: "monotonic_stack",
			DifficultyScore: 18, DifficultyLabel: "easy", EstimatedTimeSec: 60,
			Tags:     []string{"pattern-recognition", "monotonic-stack"},
			Concepts: []string{"next-greater-element", "monotonic-stack-trigger"},
		},

		// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
		// CATEGORY 4: EDGE CASES & OVERFLOW
		// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

		{
			QuestionTypeSlug: "edge-case-identification", Category: "edge_cases", Subcategory: "4.1",
			QuestionText: "For the Two Sum problem, which of these is the trickiest edge case to handle correctly?",
			Format:       models.FmtMultipleChoice,
			CorrectAnswer: models.JSONB{"option_id": "C"},
			AnswerOptions: []models.AnswerOption{
				{ID: "A", Text: "Array of length 2", IsCorrect: false},
				{ID: "B", Text: "All elements are the same", IsCorrect: false},
				{ID: "C", Text: "Target is double of an element that appears once: [3,2,4], target=6", IsCorrect: true},
				{ID: "D", Text: "Negative numbers in the array", IsCorrect: false},
			},
			Explanation:    "When target = 2*x and x appears only once, a naive approach might use the same element twice (index 0 + index 0). The hash map approach handles this naturally because you check the map BEFORE inserting the current element, but it's a common mistake in other implementations.",
			RelatedProblemSlug: "two-sum",
			DifficultyScore: 22, DifficultyLabel: "medium", EstimatedTimeSec: 90,
			Tags:     []string{"edge-case", "two-sum", "duplicate-index"},
			Concepts: []string{"edge-case-identification", "same-element-reuse"},
		},

		{
			QuestionTypeSlug: "overflow-detection", Category: "edge_cases", Subcategory: "4.2",
			QuestionText: "In this binary search code, which line has a potential integer overflow?\n\n```go\nfunc search(nums []int, target int) int {\n    lo, hi := 0, len(nums)-1\n    for lo <= hi {\n        mid := (lo + hi) / 2  // Line A\n        if nums[mid] == target { return mid }\n        if nums[mid] < target { lo = mid + 1 }\n        else { hi = mid - 1 }\n    }\n    return -1\n}\n```",
			QuestionData: models.JSONB{"language": "go"},
			Format:       models.FmtMultipleChoice,
			CorrectAnswer: models.JSONB{"option_id": "A", "value": "Line A: mid := (lo + hi) / 2"},
			AnswerOptions: []models.AnswerOption{
				{ID: "A", Text: "Line A: mid := (lo + hi) / 2 — lo + hi can overflow", IsCorrect: true},
				{ID: "B", Text: "lo = mid + 1 — mid + 1 can overflow", IsCorrect: false},
				{ID: "C", Text: "hi = mid - 1 — mid - 1 can go negative", IsCorrect: false},
				{ID: "D", Text: "No overflow risk in this code", IsCorrect: false},
			},
			Explanation:    "If lo and hi are both close to math.MaxInt, their sum overflows. Fix: mid := lo + (hi - lo) / 2. In Go with int (64-bit), this is rarely hit with array indices, but in 32-bit languages or when binary searching on value ranges (e.g., answer in [0, 2×10⁹]), it's critical.",
			CommonMistakes: []string{"Thinking Go's int is always safe (it is for array indices, but not value ranges)", "Not knowing the lo + (hi-lo)/2 idiom"},
			RelatedTopicID: "binary_search",
			DifficultyScore: 28, DifficultyLabel: "medium", EstimatedTimeSec: 90,
			Tags:     []string{"overflow", "binary-search", "classic-bug"},
			Concepts: []string{"integer-overflow", "binary-search-midpoint"},
		},

		{
			QuestionTypeSlug: "base-case-testing", Category: "edge_cases", Subcategory: "4.3",
			QuestionText: "For the Climbing Stairs DP problem (dp[i] = dp[i-1] + dp[i-2]), what happens if you forget to set dp[0] = 1?\n\n```go\ndp := make([]int, n+1)\n// dp[0] = 1  ← MISSING\ndp[1] = 1\nfor i := 2; i <= n; i++ {\n    dp[i] = dp[i-1] + dp[i-2]\n}\n```",
			Format:       models.FmtMultipleChoice,
			CorrectAnswer: models.JSONB{"option_id": "B"},
			AnswerOptions: []models.AnswerOption{
				{ID: "A", Text: "It still works correctly because dp[0] defaults to 0 in Go", IsCorrect: false},
				{ID: "B", Text: "All values are off — dp[2] becomes 1 instead of 2", IsCorrect: true},
				{ID: "C", Text: "Index out of bounds error", IsCorrect: false},
				{ID: "D", Text: "Infinite loop", IsCorrect: false},
			},
			Explanation:    "dp[0]=0 (Go default), dp[1]=1, so dp[2] = dp[1]+dp[0] = 1+0 = 1 instead of 2. Every subsequent value will be wrong (the Fibonacci sequence starting from 0,1 instead of 1,1). dp[0]=1 represents 'there is 1 way to stand at the ground'.",
			RelatedProblemSlug: "climbing-stairs",
			RelatedTopicID: "dp_1d",
			DifficultyScore: 18, DifficultyLabel: "easy", EstimatedTimeSec: 90,
			Tags:     []string{"base-case", "dp", "off-by-one", "tricky-base-case"},
			Concepts: []string{"dp-base-case", "initialization-error"},
		},

		// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
		// CATEGORY 5: CODE TEMPLATES
		// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

		{
			QuestionTypeSlug: "template-identification", Category: "code_templates", Subcategory: "5.3",
			QuestionText: "For a problem asking 'Find the minimum eating speed such that all bananas are eaten within h hours', which template should you start from?",
			Format:       models.FmtMultipleChoice,
			CorrectAnswer: models.JSONB{"option_id": "B"},
			AnswerOptions: []models.AnswerOption{
				{ID: "A", Text: "Standard binary search on sorted array", IsCorrect: false},
				{ID: "B", Text: "Binary search on answer space with feasibility check", IsCorrect: true},
				{ID: "C", Text: "Sliding window with variable size", IsCorrect: false},
				{ID: "D", Text: "Greedy with sorting", IsCorrect: false},
			},
			Explanation:    "'Minimize X such that condition is met' is the classic binary search on answer template. Define search space [1, max(piles)], for each candidate speed check feasibility in O(n), binary search narrows to the minimum valid speed.",
			RelatedProblemSlug: "koko-eating-bananas",
			RelatedTopicID: "bs_on_answer",
			DifficultyScore: 20, DifficultyLabel: "easy", EstimatedTimeSec: 60,
			Tags:     []string{"template", "binary-search-on-answer"},
			Concepts: []string{"template-selection", "search-on-answer-space"},
		},

		{
			QuestionTypeSlug: "template-completion", Category: "code_templates", Subcategory: "5.1",
			QuestionText: "Complete the BFS template. Fill in the blank:\n\n```go\nqueue := []int{start}\nvisited := map[int]bool{start: true}\n\nfor len(queue) > 0 {\n    node := queue[0]\n    queue = queue[1:]\n    \n    for _, neighbor := range adj[node] {\n        if _____ {\n            visited[neighbor] = true\n            queue = append(queue, neighbor)\n        }\n    }\n}\n```\n\nWhat goes in the blank?",
			QuestionData: models.JSONB{"language": "go"},
			Format:       models.FmtFillBlank,
			CorrectAnswer: models.JSONB{"value": "!visited[neighbor]"},
			AnswerOptions: []models.AnswerOption{
				{ID: "A", Text: "!visited[neighbor]", IsCorrect: true},
				{ID: "B", Text: "visited[neighbor] == false", IsCorrect: false},
				{ID: "C", Text: "neighbor != start", IsCorrect: false},
				{ID: "D", Text: "len(queue) < n", IsCorrect: false},
			},
			Explanation:    "The guard condition prevents revisiting nodes. Without it, BFS would loop forever on cyclic graphs. We mark visited BEFORE enqueueing (not after dequeueing) to prevent the same node from being added multiple times.",
			CommonMistakes: []string{"Marking visited after dequeueing instead of before enqueueing", "Checking visited after push causes duplicate entries in queue"},
			RelatedTopicID: "graph_bfs",
			DifficultyScore: 12, DifficultyLabel: "easy", EstimatedTimeSec: 60,
			Tags:     []string{"template", "bfs", "graph"},
			Concepts: []string{"bfs-template", "visited-check-timing"},
		},

		// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
		// CATEGORY 6: IMPLEMENTATION CORRECTNESS
		// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

		{
			QuestionTypeSlug: "return-value-semantics", Category: "implementation", Subcategory: "6.1",
			QuestionText: "What does this function return for input nums = [1, 1, 2, 2, 3]?\n\n```go\nfunc removeDuplicates(nums []int) int {\n    if len(nums) == 0 { return 0 }\n    slow := 0\n    for fast := 1; fast < len(nums); fast++ {\n        if nums[fast] != nums[slow] {\n            slow++\n            nums[slow] = nums[fast]\n        }\n    }\n    return slow + 1\n}\n```",
			QuestionData: models.JSONB{"language": "go"},
			Format:       models.FmtMultipleChoice,
			CorrectAnswer: models.JSONB{"option_id": "B", "value": "3"},
			AnswerOptions: []models.AnswerOption{
				{ID: "A", Text: "2", IsCorrect: false},
				{ID: "B", Text: "3", IsCorrect: true},
				{ID: "C", Text: "5", IsCorrect: false},
				{ID: "D", Text: "4", IsCorrect: false},
			},
			Explanation:    "Trace: slow=0. fast=1: 1==1 skip. fast=2: 2!=1 → slow=1, nums[1]=2. fast=3: 2==2 skip. fast=4: 3!=2 → slow=2, nums[2]=3. Return slow+1 = 3. Array becomes [1,2,3,2,3] with first 3 elements being unique.",
			RelatedTopicID: "two_ptr_same_dir",
			DifficultyScore: 20, DifficultyLabel: "easy", EstimatedTimeSec: 120,
			Tags:     []string{"implementation", "two-pointers", "dry-run"},
			Concepts: []string{"slow-fast-pointer", "in-place-modification", "return-value"},
		},

		{
			QuestionTypeSlug: "dry-run-trace", Category: "implementation", Subcategory: "6.3",
			QuestionText: "Dry run this Coin Change DP for coins=[1,5,11] and amount=15. What is dp[15]?\n\n```go\ndp := make([]int, amount+1)\nfor i := 1; i <= amount; i++ {\n    dp[i] = amount + 1 // sentinel\n}\nfor i := 1; i <= amount; i++ {\n    for _, c := range coins {\n        if c <= i && dp[i-c]+1 < dp[i] {\n            dp[i] = dp[i-c] + 1\n        }\n    }\n}\n```",
			Format:       models.FmtMultipleChoice,
			CorrectAnswer: models.JSONB{"option_id": "A", "value": "3"},
			AnswerOptions: []models.AnswerOption{
				{ID: "A", Text: "3 (using three 5-coins)", IsCorrect: true},
				{ID: "B", Text: "5 (using fifteen 1-coins... wait that's 15)", IsCorrect: false},
				{ID: "C", Text: "2 (using 11+4... but 4 isn't a coin)", IsCorrect: false},
				{ID: "D", Text: "4 (using 11+1+1+1+1)", IsCorrect: false},
			},
			Explanation:    "dp[15]: check coin 11 → dp[4]+1=4+1=5. Check coin 5 → dp[10]+1=2+1=3. Check coin 1 → dp[14]+1. Best is 3 (5+5+5). Note: the greedy approach of picking 11 first gives 11+1+1+1+1=5 coins — worse than 5+5+5=3 coins. This is why DP beats greedy here.",
			RelatedProblemSlug: "coin-change",
			RelatedTopicID: "dp_knapsack",
			DifficultyScore: 30, DifficultyLabel: "medium", EstimatedTimeSec: 180,
			Tags:     []string{"dry-run", "dp", "coin-change", "greedy-trap"},
			Concepts: []string{"dp-trace", "greedy-vs-dp", "unbounded-knapsack"},
		},

		// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
		// CATEGORY 7: BUG DETECTION
		// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

		{
			QuestionTypeSlug: "find-the-bug", Category: "bug_detection", Subcategory: "7.1",
			QuestionText: "This BFS for 'Number of Islands' has a bug. What is it?\n\n```go\nfunc numIslands(grid [][]byte) int {\n    count := 0\n    for i := range grid {\n        for j := range grid[0] {\n            if grid[i][j] == '1' {\n                count++\n                bfs(grid, i, j)\n            }\n        }\n    }\n    return count\n}\n\nfunc bfs(grid [][]byte, r, c int) {\n    queue := [][2]int{{r, c}}\n    for len(queue) > 0 {\n        cell := queue[0]\n        queue = queue[1:]\n        for _, d := range [][2]int{{0,1},{0,-1},{1,0},{-1,0}} {\n            nr, nc := cell[0]+d[0], cell[1]+d[1]\n            if nr >= 0 && nr < len(grid) && nc >= 0 && nc < len(grid[0]) && grid[nr][nc] == '1' {\n                grid[nr][nc] = '0'\n                queue = append(queue, [2]int{nr, nc})\n            }\n        }\n    }\n}\n```",
			Format:       models.FmtMultipleChoice,
			CorrectAnswer: models.JSONB{"option_id": "B"},
			AnswerOptions: []models.AnswerOption{
				{ID: "A", Text: "Directions array is wrong — should include diagonals", IsCorrect: false},
				{ID: "B", Text: "The starting cell (r,c) is never marked as '0' — it stays '1' and could be re-counted", IsCorrect: true},
				{ID: "C", Text: "Queue should use a struct instead of [2]int", IsCorrect: false},
				{ID: "D", Text: "Should use DFS instead of BFS", IsCorrect: false},
			},
			Explanation:    "The bfs function marks neighbors as '0' but never marks the starting cell itself. This means the same starting cell could be visited again in the outer loop, leading to incorrect island counts. Fix: add grid[r][c] = '0' at the start of bfs.",
			CommonMistakes: []string{"Not marking the root/start node before processing", "Only marking children, forgetting the source"},
			RelatedProblemSlug: "number-of-islands",
			RelatedTopicID: "graph_bfs",
			DifficultyScore: 30, DifficultyLabel: "medium", EstimatedTimeSec: 180,
			Tags:     []string{"bug-detection", "bfs", "grid", "visited-mark"},
			Concepts: []string{"bfs-mark-before-enqueue", "grid-flood-fill-bug"},
		},

		{
			QuestionTypeSlug: "will-it-crash", Category: "bug_detection", Subcategory: "7.3",
			QuestionText: "Will this code crash for input `head = nil` (empty linked list)?\n\n```go\nfunc hasCycle(head *ListNode) bool {\n    slow, fast := head, head\n    for fast != nil && fast.Next != nil {\n        slow = slow.Next\n        fast = fast.Next.Next\n        if slow == fast { return true }\n    }\n    return false\n}\n```",
			Format:       models.FmtMultipleChoice,
			CorrectAnswer: models.JSONB{"option_id": "B"},
			AnswerOptions: []models.AnswerOption{
				{ID: "A", Text: "Yes — nil pointer dereference on slow.Next", IsCorrect: false},
				{ID: "B", Text: "No — the for loop condition 'fast != nil' catches it immediately", IsCorrect: true},
				{ID: "C", Text: "Yes — fast.Next.Next crashes", IsCorrect: false},
				{ID: "D", Text: "Depends on the Go runtime version", IsCorrect: false},
			},
			Explanation:    "When head=nil, both slow and fast are nil. The loop condition `fast != nil` is false, so the loop body never executes. The function safely returns false. The conditions are checked left-to-right with short-circuit evaluation.",
			RelatedTopicID: "two_ptr_cycle",
			DifficultyScore: 15, DifficultyLabel: "easy", EstimatedTimeSec: 60,
			Tags:     []string{"bug-detection", "nil-check", "linked-list", "cycle-detection"},
			Concepts: []string{"nil-safety", "short-circuit-evaluation", "floyd-cycle"},
		},

		// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
		// CATEGORY 9: TRADE-OFF ANALYSIS
		// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

		{
			QuestionTypeSlug: "time-space-tradeoff", Category: "tradeoffs", Subcategory: "9.1",
			QuestionText: "For the House Robber problem, you have two options:\n\nA) dp array of size n: dp[i] = max(dp[i-1], dp[i-2]+nums[i])\nB) Two variables prev1, prev2 (space-optimized)\n\nWhen should you prefer option A over B?",
			Format:       models.FmtMultipleChoice,
			CorrectAnswer: models.JSONB{"option_id": "C"},
			AnswerOptions: []models.AnswerOption{
				{ID: "A", Text: "Always — the array is clearer", IsCorrect: false},
				{ID: "B", Text: "Never — O(1) space is always better", IsCorrect: false},
				{ID: "C", Text: "When you need to reconstruct which houses were robbed (backtrack through dp array)", IsCorrect: true},
				{ID: "D", Text: "When n is very large (>10⁶)", IsCorrect: false},
			},
			Explanation:    "The O(1) space version loses the ability to backtrack and reconstruct the solution. If you only need the maximum value, use two variables. If you need to know WHICH houses to rob (the actual subset), keep the full dp array and trace back from dp[n-1].",
			RelatedProblemSlug: "house-robber",
			RelatedTopicID: "dp_1d",
			DifficultyScore: 30, DifficultyLabel: "medium", EstimatedTimeSec: 120,
			Tags:     []string{"tradeoff", "dp", "space-optimization", "solution-reconstruction"},
			Concepts: []string{"space-time-tradeoff", "dp-backtracking", "solution-reconstruction"},
		},

		{
			QuestionTypeSlug: "approach-comparison", Category: "tradeoffs", Subcategory: "9.2",
			QuestionText: "For Kth Largest Element (n=10⁵, k can vary per query, multiple queries), which approach is best?\n\nA) Sort once, index directly: O(n log n) + O(1) per query\nB) Min-heap of size k: O(n log k) per query\nC) Quickselect: O(n) average per query\nD) Max-heap, pop k times: O(n + k log n) per query",
			Format:       models.FmtMultipleChoice,
			CorrectAnswer: models.JSONB{"option_id": "A"},
			AnswerOptions: []models.AnswerOption{
				{ID: "A", Text: "Sort once, then O(1) per query — best for multiple queries", IsCorrect: true},
				{ID: "B", Text: "Min-heap of size k — best per single query", IsCorrect: false},
				{ID: "C", Text: "Quickselect — best average case for single query", IsCorrect: false},
				{ID: "D", Text: "Max-heap — simplest to implement", IsCorrect: false},
			},
			Explanation:    "With multiple queries and static data, sorting once (O(n log n)) then answering each query in O(1) via indexing is optimal. For a single query, quickselect O(n) average is best. The key insight is amortizing the sort cost across many queries.",
			RelatedProblemSlug: "kth-largest-element",
			DifficultyScore: 35, DifficultyLabel: "medium", EstimatedTimeSec: 180,
			Tags:     []string{"tradeoff", "heap", "sorting", "quickselect", "amortized"},
			Concepts: []string{"query-amortization", "preprocessing-tradeoff"},
		},

		// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
		// CATEGORY 10: HYBRID / MULTI-SKILL
		// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

		{
			QuestionTypeSlug: "optimization-challenge", Category: "hybrid", Subcategory: "10.2",
			QuestionText: "Three Sum problem: Rank these approaches from worst to best time complexity.\n\n1. Three nested loops checking all triplets\n2. Sort + fix one element + two-pointer for remaining two\n3. Fix one element + hash set for complement\n4. Sort + binary search for the third element",
			Format:       models.FmtRanking,
			CorrectAnswer: models.JSONB{
				"ranking": []interface{}{"1", "4", "3", "2"},
				"explanation": "1→O(n³), 4→O(n² log n), 3→O(n²) with more space, 2→O(n²) with O(1) extra space (best overall)",
			},
			AnswerOptions: []models.AnswerOption{
				{ID: "1", Text: "Three nested loops: O(n³)"},
				{ID: "2", Text: "Sort + two pointers: O(n²)"},
				{ID: "3", Text: "Fix one + hash set: O(n²)"},
				{ID: "4", Text: "Sort + binary search: O(n² log n)"},
			},
			Explanation:    "Worst→Best: O(n³) brute force → O(n² log n) sort+binary search → O(n²) hash set (more space) → O(n²) sort+two pointers (optimal: same time, less space). The two-pointer approach is best because it achieves O(n²) with O(1) extra space after sorting.",
			RelatedProblemSlug: "3sum",
			DifficultyScore: 35, DifficultyLabel: "medium", EstimatedTimeSec: 240,
			Tags:     []string{"optimization", "ranking", "3sum", "multiple-approaches"},
			Concepts: []string{"progressive-optimization", "space-time-tradeoff", "approach-ranking"},
		},
	}
}

package seed

import (
	"encoding/json"
	"time"

	"github.com/yourusername/algoholic/models"
)

// Helper functions
func intPtr(i int) *int {
	return &i
}

func strPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}

func jsonbMap(data map[string]interface{}) models.JSONB {
	return models.JSONB(data)
}

func jsonbArray(data []interface{}) models.JSONBArray {
	return models.JSONBArray(data)
}

func jsonbFromString(s string) models.JSONB {
	var data map[string]interface{}
	json.Unmarshal([]byte(s), &data)
	return models.JSONB(data)
}

func jsonbArrayFromString(s string) models.JSONBArray {
	var data []interface{}
	json.Unmarshal([]byte(s), &data)
	return models.JSONBArray(data)
}

// GetSeedTopics returns initial topics for the database
func GetSeedTopics() []models.Topic {
	return []models.Topic{
		{
			Name:                  "Arrays",
			Slug:                  "arrays",
			Description:           strPtr("Fundamental data structure for storing elements in contiguous memory"),
			Category:              strPtr("data_structure"),
			DifficultyLevel:       intPtr(1),
			EstimatedLearningHour: float64Ptr(2.0),
		},
		{
			Name:                  "Hash Table",
			Slug:                  "hash-table",
			Description:           strPtr("Data structure for fast key-value lookups using hashing"),
			Category:              strPtr("data_structure"),
			DifficultyLevel:       intPtr(2),
			EstimatedLearningHour: float64Ptr(3.0),
		},
		{
			Name:                  "Two Pointers",
			Slug:                  "two-pointers",
			Description:           strPtr("Technique using two pointers to traverse data structures"),
			Category:              strPtr("algorithm"),
			DifficultyLevel:       intPtr(2),
			EstimatedLearningHour: float64Ptr(4.0),
		},
		{
			Name:                  "Sliding Window",
			Slug:                  "sliding-window",
			Description:           strPtr("Pattern for solving array/string problems using a moving window"),
			Category:              strPtr("pattern"),
			DifficultyLevel:       intPtr(3),
			EstimatedLearningHour: float64Ptr(5.0),
		},
		{
			Name:                  "Binary Search",
			Slug:                  "binary-search",
			Description:           strPtr("Efficient search algorithm for sorted arrays"),
			Category:              strPtr("algorithm"),
			DifficultyLevel:       intPtr(3),
			EstimatedLearningHour: float64Ptr(4.0),
		},
		{
			Name:                  "Dynamic Programming",
			Slug:                  "dynamic-programming",
			Description:           strPtr("Optimization technique using memoization or tabulation"),
			Category:              strPtr("algorithm"),
			DifficultyLevel:       intPtr(5),
			EstimatedLearningHour: float64Ptr(15.0),
		},
		{
			Name:                  "Graph Traversal",
			Slug:                  "graph-traversal",
			Description:           strPtr("Algorithms for exploring graph structures (DFS, BFS)"),
			Category:              strPtr("algorithm"),
			DifficultyLevel:       intPtr(4),
			EstimatedLearningHour: float64Ptr(8.0),
		},
		{
			Name:                  "Trees",
			Slug:                  "trees",
			Description:           strPtr("Hierarchical data structure with nodes and edges"),
			Category:              strPtr("data_structure"),
			DifficultyLevel:       intPtr(3),
			EstimatedLearningHour: float64Ptr(6.0),
		},
		{
			Name:                  "Stack",
			Slug:                  "stack",
			Description:           strPtr("LIFO (Last In First Out) data structure"),
			Category:              strPtr("data_structure"),
			DifficultyLevel:       intPtr(2),
			EstimatedLearningHour: float64Ptr(2.0),
		},
		{
			Name:                  "Linked List",
			Slug:                  "linked-list",
			Description:           strPtr("Linear data structure with nodes connected by pointers"),
			Category:              strPtr("data_structure"),
			DifficultyLevel:       intPtr(2),
			EstimatedLearningHour: float64Ptr(4.0),
		},
	}
}

// GetSeedProblems returns initial problems for the database
func GetSeedProblems() []models.Problem {
	return []models.Problem{
		// Easy Problems
		{
			LeetcodeNumber:     intPtr(1),
			Title:              "Two Sum",
			Slug:               "two-sum",
			Description:        "Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target. You may assume that each input would have exactly one solution, and you may not use the same element twice. You can return the answer in any order.",
			Constraints:        models.StringArray{"2 <= nums.length <= 10^4", "-10^9 <= nums[i] <= 10^9", "-10^9 <= target <= 10^9", "Only one valid answer exists"},
			Examples:           jsonbArrayFromString(`[{"input": "nums = [2,7,11,15], target = 9", "output": "[0,1]", "explanation": "Because nums[0] + nums[1] == 9, we return [0, 1]"}, {"input": "nums = [3,2,4], target = 6", "output": "[1,2]"}]`),
			Hints:              models.StringArray{"Try using a hash map to store numbers you've seen", "For each element, check if target - element exists in the hash map"},
			DifficultyScore:    15.0,
			OfficialDifficulty: strPtr("Easy"),
			PrimaryPattern:     strPtr("Hash Table"),
			SecondaryPatterns:  models.StringArray{"Array"},
			TimeComplexity:     strPtr("O(n)"),
			SpaceComplexity:    strPtr("O(n)"),
		},
		{
			LeetcodeNumber:     intPtr(121),
			Title:              "Best Time to Buy and Sell Stock",
			Slug:               "best-time-to-buy-and-sell-stock",
			Description:        "You are given an array prices where prices[i] is the price of a given stock on the ith day. You want to maximize your profit by choosing a single day to buy one stock and choosing a different day in the future to sell that stock. Return the maximum profit you can achieve from this transaction. If you cannot achieve any profit, return 0.",
			Constraints:        models.StringArray{"1 <= prices.length <= 10^5", "0 <= prices[i] <= 10^4"},
			Examples:           jsonbArrayFromString(`[{"input": "prices = [7,1,5,3,6,4]", "output": "5", "explanation": "Buy on day 2 (price = 1) and sell on day 5 (price = 6), profit = 6-1 = 5"}, {"input": "prices = [7,6,4,3,1]", "output": "0", "explanation": "No profit possible"}]`),
			Hints:              models.StringArray{"Track the minimum price seen so far", "Calculate profit if selling at current price", "Keep track of maximum profit"},
			DifficultyScore:    20.0,
			OfficialDifficulty: strPtr("Easy"),
			PrimaryPattern:     strPtr("Array"),
			SecondaryPatterns:  models.StringArray{"Dynamic Programming"},
			TimeComplexity:     strPtr("O(n)"),
			SpaceComplexity:    strPtr("O(1)"),
		},
		{
			LeetcodeNumber:     intPtr(217),
			Title:              "Contains Duplicate",
			Slug:               "contains-duplicate",
			Description:        "Given an integer array nums, return true if any value appears at least twice in the array, and return false if every element is distinct.",
			Constraints:        models.StringArray{"1 <= nums.length <= 10^5", "-10^9 <= nums[i] <= 10^9"},
			Examples:           jsonbArrayFromString(`[{"input": "nums = [1,2,3,1]", "output": "true"}, {"input": "nums = [1,2,3,4]", "output": "false"}]`),
			Hints:              models.StringArray{"Use a hash set to track seen numbers", "Return true immediately when a duplicate is found"},
			DifficultyScore:    10.0,
			OfficialDifficulty: strPtr("Easy"),
			PrimaryPattern:     strPtr("Hash Table"),
			SecondaryPatterns:  models.StringArray{"Array"},
			TimeComplexity:     strPtr("O(n)"),
			SpaceComplexity:    strPtr("O(n)"),
		},
		{
			LeetcodeNumber:     intPtr(20),
			Title:              "Valid Parentheses",
			Slug:               "valid-parentheses",
			Description:        "Given a string s containing just the characters '(', ')', '{', '}', '[' and ']', determine if the input string is valid. An input string is valid if: Open brackets must be closed by the same type of brackets, and Open brackets must be closed in the correct order.",
			Constraints:        models.StringArray{"1 <= s.length <= 10^4", "s consists of parentheses only '()[]{}'"},
			Examples:           jsonbArrayFromString(`[{"input": "s = \"()\"", "output": "true"}, {"input": "s = \"()[]{}\"", "output": "true"}, {"input": "s = \"(]\"", "output": "false"}]`),
			Hints:              models.StringArray{"Use a stack to keep track of opening brackets", "When you encounter a closing bracket, check if it matches the top of the stack"},
			DifficultyScore:    18.0,
			OfficialDifficulty: strPtr("Easy"),
			PrimaryPattern:     strPtr("Stack"),
			SecondaryPatterns:  models.StringArray{"String"},
			TimeComplexity:     strPtr("O(n)"),
			SpaceComplexity:    strPtr("O(n)"),
		},
		{
			LeetcodeNumber:     intPtr(206),
			Title:              "Reverse Linked List",
			Slug:               "reverse-linked-list",
			Description:        "Given the head of a singly linked list, reverse the list, and return the reversed list.",
			Constraints:        models.StringArray{"The number of nodes in the list is the range [0, 5000]", "-5000 <= Node.val <= 5000"},
			Examples:           jsonbArrayFromString(`[{"input": "head = [1,2,3,4,5]", "output": "[5,4,3,2,1]"}, {"input": "head = []", "output": "[]"}]`),
			Hints:              models.StringArray{"Use three pointers: prev, current, next", "Iteratively reverse the direction of pointers"},
			DifficultyScore:    22.0,
			OfficialDifficulty: strPtr("Easy"),
			PrimaryPattern:     strPtr("Linked List"),
			SecondaryPatterns:  models.StringArray{},
			TimeComplexity:     strPtr("O(n)"),
			SpaceComplexity:    strPtr("O(1)"),
		},

		// Medium Problems
		{
			LeetcodeNumber:     intPtr(15),
			Title:              "3Sum",
			Slug:               "3sum",
			Description:        "Given an integer array nums, return all the triplets [nums[i], nums[j], nums[k]] such that i != j, i != k, and j != k, and nums[i] + nums[j] + nums[k] == 0. Notice that the solution set must not contain duplicate triplets.",
			Constraints:        models.StringArray{"3 <= nums.length <= 3000", "-10^5 <= nums[i] <= 10^5"},
			Examples:           jsonbArrayFromString(`[{"input": "nums = [-1,0,1,2,-1,-4]", "output": "[[-1,-1,2],[-1,0,1]]"}, {"input": "nums = [0,1,1]", "output": "[]"}]`),
			Hints:              models.StringArray{"Sort the array first", "Use two pointers for each fixed element", "Skip duplicates to avoid duplicate triplets"},
			DifficultyScore:    45.0,
			OfficialDifficulty: strPtr("Medium"),
			PrimaryPattern:     strPtr("Two Pointers"),
			SecondaryPatterns:  models.StringArray{"Array", "Sorting"},
			TimeComplexity:     strPtr("O(n^2)"),
			SpaceComplexity:    strPtr("O(1)"),
		},
		{
			LeetcodeNumber:     intPtr(33),
			Title:              "Search in Rotated Sorted Array",
			Slug:               "search-in-rotated-sorted-array",
			Description:        "There is an integer array nums sorted in ascending order (with distinct values). Prior to being passed to your function, nums is possibly rotated at an unknown pivot index k. Given the array nums after the possible rotation and an integer target, return the index of target if it is in nums, or -1 if it is not in nums. You must write an algorithm with O(log n) runtime complexity.",
			Constraints:        models.StringArray{"1 <= nums.length <= 5000", "-10^4 <= nums[i] <= 10^4", "All values of nums are unique", "nums is an ascending array that is possibly rotated"},
			Examples:           jsonbArrayFromString(`[{"input": "nums = [4,5,6,7,0,1,2], target = 0", "output": "4"}, {"input": "nums = [4,5,6,7,0,1,2], target = 3", "output": "-1"}]`),
			Hints:              models.StringArray{"Use modified binary search", "Determine which half is sorted", "Check if target is in the sorted half"},
			DifficultyScore:    55.0,
			OfficialDifficulty: strPtr("Medium"),
			PrimaryPattern:     strPtr("Binary Search"),
			SecondaryPatterns:  models.StringArray{"Array"},
			TimeComplexity:     strPtr("O(log n)"),
			SpaceComplexity:    strPtr("O(1)"),
		},
		{
			LeetcodeNumber:     intPtr(3),
			Title:              "Longest Substring Without Repeating Characters",
			Slug:               "longest-substring-without-repeating-characters",
			Description:        "Given a string s, find the length of the longest substring without repeating characters.",
			Constraints:        models.StringArray{"0 <= s.length <= 5 * 10^4", "s consists of English letters, digits, symbols and spaces"},
			Examples:           jsonbArrayFromString(`[{"input": "s = \"abcabcbb\"", "output": "3", "explanation": "The answer is abc with length 3"}, {"input": "s = \"bbbbb\"", "output": "1"}, {"input": "s = \"pwwkew\"", "output": "3"}]`),
			Hints:              models.StringArray{"Use sliding window technique", "Use a hash set to track characters in current window", "Expand window when no duplicates, shrink when duplicate found"},
			DifficultyScore:    48.0,
			OfficialDifficulty: strPtr("Medium"),
			PrimaryPattern:     strPtr("Sliding Window"),
			SecondaryPatterns:  models.StringArray{"Hash Table", "String"},
			TimeComplexity:     strPtr("O(n)"),
			SpaceComplexity:    strPtr("O(min(m,n))"),
		},
		{
			LeetcodeNumber:     intPtr(200),
			Title:              "Number of Islands",
			Slug:               "number-of-islands",
			Description:        "Given an m x n 2D binary grid grid which represents a map of '1's (land) and '0's (water), return the number of islands. An island is surrounded by water and is formed by connecting adjacent lands horizontally or vertically.",
			Constraints:        models.StringArray{"m == grid.length", "n == grid[i].length", "1 <= m, n <= 300", "grid[i][j] is '0' or '1'"},
			Examples:           jsonbArrayFromString(`[{"input": "grid = [[\"1\",\"1\",\"1\",\"1\",\"0\"],[\"1\",\"1\",\"0\",\"1\",\"0\"],[\"1\",\"1\",\"0\",\"0\",\"0\"],[\"0\",\"0\",\"0\",\"0\",\"0\"]]", "output": "1"}, {"input": "grid = [[\"1\",\"1\",\"0\",\"0\",\"0\"],[\"1\",\"1\",\"0\",\"0\",\"0\"],[\"0\",\"0\",\"1\",\"0\",\"0\"],[\"0\",\"0\",\"0\",\"1\",\"1\"]]", "output": "3"}]`),
			Hints:              models.StringArray{"Use DFS or BFS to explore each island", "Mark visited cells to avoid recounting", "Count the number of DFS/BFS calls needed"},
			DifficultyScore:    50.0,
			OfficialDifficulty: strPtr("Medium"),
			PrimaryPattern:     strPtr("Graph Traversal"),
			SecondaryPatterns:  models.StringArray{"DFS", "BFS"},
			TimeComplexity:     strPtr("O(m*n)"),
			SpaceComplexity:    strPtr("O(m*n)"),
		},
		{
			LeetcodeNumber:     intPtr(53),
			Title:              "Maximum Subarray",
			Slug:               "maximum-subarray",
			Description:        "Given an integer array nums, find the subarray with the largest sum, and return its sum.",
			Constraints:        models.StringArray{"1 <= nums.length <= 10^5", "-10^4 <= nums[i] <= 10^4"},
			Examples:           jsonbArrayFromString(`[{"input": "nums = [-2,1,-3,4,-1,2,1,-5,4]", "output": "6", "explanation": "The subarray [4,-1,2,1] has the largest sum 6"}, {"input": "nums = [1]", "output": "1"}]`),
			Hints:              models.StringArray{"Use Kadane's algorithm", "Track current sum and maximum sum", "Reset current sum to 0 if it becomes negative"},
			DifficultyScore:    42.0,
			OfficialDifficulty: strPtr("Medium"),
			PrimaryPattern:     strPtr("Dynamic Programming"),
			SecondaryPatterns:  models.StringArray{"Array"},
			TimeComplexity:     strPtr("O(n)"),
			SpaceComplexity:    strPtr("O(1)"),
		},

		// Hard Problems
		{
			LeetcodeNumber:     intPtr(42),
			Title:              "Trapping Rain Water",
			Slug:               "trapping-rain-water",
			Description:        "Given n non-negative integers representing an elevation map where the width of each bar is 1, compute how much water it can trap after raining.",
			Constraints:        models.StringArray{"n == height.length", "1 <= n <= 2 * 10^4", "0 <= height[i] <= 10^5"},
			Examples:           jsonbArrayFromString(`[{"input": "height = [0,1,0,2,1,0,1,3,2,1,2,1]", "output": "6"}, {"input": "height = [4,2,0,3,2,5]", "output": "9"}]`),
			Hints:              models.StringArray{"Water level at each position is min(max_left, max_right)", "Use two pointers approach", "Track left_max and right_max"},
			DifficultyScore:    75.0,
			OfficialDifficulty: strPtr("Hard"),
			PrimaryPattern:     strPtr("Two Pointers"),
			SecondaryPatterns:  models.StringArray{"Array", "Dynamic Programming"},
			TimeComplexity:     strPtr("O(n)"),
			SpaceComplexity:    strPtr("O(1)"),
		},
		{
			LeetcodeNumber:     intPtr(23),
			Title:              "Merge k Sorted Lists",
			Slug:               "merge-k-sorted-lists",
			Description:        "You are given an array of k linked-lists lists, each linked-list is sorted in ascending order. Merge all the linked-lists into one sorted linked-list and return it.",
			Constraints:        models.StringArray{"k == lists.length", "0 <= k <= 10^4", "0 <= lists[i].length <= 500", "-10^4 <= lists[i][j] <= 10^4"},
			Examples:           jsonbArrayFromString(`[{"input": "lists = [[1,4,5],[1,3,4],[2,6]]", "output": "[1,1,2,3,4,4,5,6]"}, {"input": "lists = []", "output": "[]"}]`),
			Hints:              models.StringArray{"Use a min heap (priority queue)", "Extract minimum from all list heads", "Add next element from the same list"},
			DifficultyScore:    78.0,
			OfficialDifficulty: strPtr("Hard"),
			PrimaryPattern:     strPtr("Linked List"),
			SecondaryPatterns:  models.StringArray{"Heap", "Divide and Conquer"},
			TimeComplexity:     strPtr("O(N log k)"),
			SpaceComplexity:    strPtr("O(k)"),
		},
	}
}

// GetSeedQuestions returns initial questions for the database
func GetSeedQuestions() []models.Question {
	now := time.Now()
	return []models.Question{
		// Two Sum Questions
		{
			ProblemID:       intPtr(1),
			QuestionType:    "complexity_analysis",
			QuestionSubtype: strPtr("time_complexity"),
			QuestionFormat:  "multiple_choice",
			QuestionText:    "What is the time complexity of the optimal solution for Two Sum using a hash map?",
			AnswerOptions: jsonbMap(map[string]interface{}{
				"options": []map[string]string{
					{"id": "a", "text": "O(n^2)"},
					{"id": "b", "text": "O(n log n)"},
					{"id": "c", "text": "O(n)"},
					{"id": "d", "text": "O(1)"},
				},
			}),
			CorrectAnswer: jsonbMap(map[string]interface{}{"answer": "c"}),
			Explanation:   "We traverse the array once, and hash map operations (insert and lookup) are O(1) on average, giving us O(n) overall time complexity.",
			WrongAnswerExplanations: jsonbMap(map[string]interface{}{
				"a": "O(n^2) would be the brute force approach with nested loops checking every pair.",
				"b": "O(n log n) would apply if we sorted the array first and used two pointers.",
				"d": "O(1) is impossible since we must at least read all n elements once.",
			}),
			RelatedConcepts:      models.StringArray{"Hash Table", "Time Complexity", "Big O Notation"},
			CommonMistakes:       models.StringArray{"Confusing space complexity with time complexity", "Assuming hash operations are always O(1)"},
			DifficultyScore:      20.0,
			EstimatedTimeSeconds: intPtr(60),
			CreatedAt:            now,
			UpdatedAt:            now,
		},
		{
			ProblemID:      intPtr(1),
			QuestionType:   "data_structure_selection",
			QuestionFormat: "text",
			QuestionText:   "Why is a hash map the optimal data structure for solving Two Sum? Explain in your own words.",
			CorrectAnswer: jsonbMap(map[string]interface{}{
				"answer": []string{
					"hash map allows O(1) lookup",
					"hash table provides constant time access",
					"we can check if complement exists instantly",
					"hashmap gives us fast lookup for the difference",
				},
			}),
			Explanation:          "A hash map is optimal because it provides O(1) average-case lookup time. For each element, we can instantly check if its complement (target - current) exists in the map, avoiding the need for nested loops.",
			RelatedConcepts:      models.StringArray{"Hash Table", "Data Structure Selection", "Algorithm Optimization"},
			CommonMistakes:       models.StringArray{"Using array with linear search", "Sorting unnecessarily"},
			DifficultyScore:      25.0,
			EstimatedTimeSeconds: intPtr(120),
			CreatedAt:            now,
			UpdatedAt:            now,
		},
		{
			ProblemID:      intPtr(1),
			QuestionType:   "code_completion",
			QuestionFormat: "code",
			QuestionText:   "Complete the Two Sum function in Python using a hash map approach.",
			QuestionData: jsonbMap(map[string]interface{}{
				"template": "def twoSum(nums: List[int], target: int) -> List[int]:\n    # Your code here\n    pass",
				"language": "python",
			}),
			CorrectAnswer: jsonbMap(map[string]interface{}{
				"test_cases": []map[string]string{
					{"input": "[2,7,11,15]\n9", "expected": "[0, 1]"},
					{"input": "[3,2,4]\n6", "expected": "[1, 2]"},
					{"input": "[3,3]\n6", "expected": "[0, 1]"},
				},
			}),
			Explanation:          "Use a hash map to store numbers and their indices. For each number, check if target - num exists in the map. If found, return the indices. Otherwise, add the current number to the map.",
			RelatedConcepts:      models.StringArray{"Hash Table", "Array Traversal", "Python Dictionary"},
			CommonMistakes:       models.StringArray{"Returning values instead of indices", "Not handling the same index case", "Using nested loops"},
			DifficultyScore:      30.0,
			EstimatedTimeSeconds: intPtr(300),
			CreatedAt:            now,
			UpdatedAt:            now,
		},

		// Valid Parentheses Questions
		{
			ProblemID:      intPtr(4),
			QuestionType:   "pattern_recognition",
			QuestionFormat: "multiple_choice",
			QuestionText:   "What is the key insight for solving the Valid Parentheses problem?",
			AnswerOptions: jsonbMap(map[string]interface{}{
				"options": []map[string]string{
					{"id": "a", "text": "Count the number of opening and closing brackets"},
					{"id": "b", "text": "Use a stack to match brackets in LIFO order"},
					{"id": "c", "text": "Sort the string first"},
					{"id": "d", "text": "Use recursion to check each pair"},
				},
			}),
			CorrectAnswer: jsonbMap(map[string]interface{}{"answer": "b"}),
			Explanation:   "The stack data structure perfectly models the LIFO (Last In First Out) nature of bracket matching. The most recently opened bracket must be closed first.",
			WrongAnswerExplanations: jsonbMap(map[string]interface{}{
				"a": "Just counting isn't enough - the order matters. '([)]' has equal counts but is invalid.",
				"c": "Sorting would destroy the order information which is crucial for validation.",
				"d": "Recursion is possible but less efficient and more complex than using a stack.",
			}),
			RelatedConcepts:      models.StringArray{"Stack", "LIFO", "Pattern Matching"},
			CommonMistakes:       models.StringArray{"Only checking counts", "Not handling mismatched types"},
			DifficultyScore:      28.0,
			EstimatedTimeSeconds: intPtr(90),
			CreatedAt:            now,
			UpdatedAt:            now,
		},

		// 3Sum Questions
		{
			ProblemID:       intPtr(6),
			QuestionType:    "complexity_analysis",
			QuestionSubtype: strPtr("time_complexity"),
			QuestionFormat:  "text",
			QuestionText:    "What is the time complexity of the two-pointer solution for 3Sum, and why can't we do better than O(n^2)?",
			CorrectAnswer: jsonbMap(map[string]interface{}{
				"answer": []string{
					"O(n^2) - we fix one element and use two pointers for the rest",
					"O(n squared) because we check n elements and for each use two pointers",
					"time complexity is O(n^2) - one loop times two pointer search",
				},
			}),
			Explanation:          "The time complexity is O(n^2). We have an outer loop (O(n)) and for each iteration, we use two pointers to scan the remaining array (O(n)). We can't do better because we need to examine all triplet combinations, which is inherently O(n^2) with optimization.",
			RelatedConcepts:      models.StringArray{"Two Pointers", "Time Complexity", "Algorithm Optimization"},
			CommonMistakes:       models.StringArray{"Thinking it's O(n^3)", "Thinking we can achieve O(n log n)"},
			DifficultyScore:      52.0,
			EstimatedTimeSeconds: intPtr(150),
			CreatedAt:            now,
			UpdatedAt:            now,
		},

		// Sliding Window Questions
		{
			ProblemID:      intPtr(8),
			QuestionType:   "pattern_recognition",
			QuestionFormat: "multiple_choice",
			QuestionText:   "In the Longest Substring Without Repeating Characters problem, when should you shrink the window?",
			AnswerOptions: jsonbMap(map[string]interface{}{
				"options": []map[string]string{
					{"id": "a", "text": "When the window size reaches a certain length"},
					{"id": "b", "text": "When you encounter a duplicate character"},
					{"id": "c", "text": "After processing each character"},
					{"id": "d", "text": "When the window is empty"},
				},
			}),
			CorrectAnswer: jsonbMap(map[string]interface{}{"answer": "b"}),
			Explanation:   "You should shrink the window (move the left pointer) when you encounter a character that's already in the current window. This maintains the 'no repeating characters' invariant.",
			WrongAnswerExplanations: jsonbMap(map[string]interface{}{
				"a": "The window size is variable and depends on when duplicates appear.",
				"c": "You only shrink when necessary (duplicate found), not after every character.",
				"d": "The window should never be empty during processing.",
			}),
			RelatedConcepts:      models.StringArray{"Sliding Window", "Hash Set", "Two Pointers"},
			CommonMistakes:       models.StringArray{"Shrinking too early", "Not updating the result before shrinking"},
			DifficultyScore:      55.0,
			EstimatedTimeSeconds: intPtr(120),
			CreatedAt:            now,
			UpdatedAt:            now,
		},

		// Graph Traversal Questions
		{
			ProblemID:      intPtr(9),
			QuestionType:   "algorithm_explanation",
			QuestionFormat: "text",
			QuestionText:   "Explain the difference between DFS and BFS for the Number of Islands problem. Which would you choose and why?",
			CorrectAnswer: jsonbMap(map[string]interface{}{
				"answer": []string{
					"DFS explores deep into one island before moving on, BFS explores level by level. Either works equally well.",
					"DFS uses recursion or stack, BFS uses queue. Both have same time complexity O(mn).",
					"Both DFS and BFS work fine - DFS is often simpler to implement with recursion",
				},
			}),
			Explanation:          "Both DFS and BFS work equally well for this problem with O(m*n) time complexity. DFS is often preferred because it can be implemented concisely with recursion. BFS would use a queue and requires more code but gives the same result.",
			RelatedConcepts:      models.StringArray{"DFS", "BFS", "Graph Traversal", "Algorithm Choice"},
			CommonMistakes:       models.StringArray{"Thinking one is faster", "Not marking visited cells"},
			DifficultyScore:      58.0,
			EstimatedTimeSeconds: intPtr(180),
			CreatedAt:            now,
			UpdatedAt:            now,
		},

		// Dynamic Programming Questions
		{
			ProblemID:      intPtr(10),
			QuestionType:   "pattern_recognition",
			QuestionFormat: "text",
			QuestionText:   "For Maximum Subarray, explain why Kadane's Algorithm works. What is the key insight?",
			CorrectAnswer: jsonbMap(map[string]interface{}{
				"answer": []string{
					"At each position, decide whether to extend current subarray or start new one",
					"Keep current sum if positive, reset to current element if negative",
					"The key is that negative prefix can only decrease the sum",
					"If current sum becomes negative, starting fresh gives better chance",
				},
			}),
			Explanation:          "Kadane's algorithm works because at each position, we make an optimal local decision: extend the current subarray if it has a positive sum, or start a new subarray at the current position if the previous sum was negative. A negative prefix can only harm our total, so it's always better to start fresh.",
			RelatedConcepts:      models.StringArray{"Dynamic Programming", "Kadane's Algorithm", "Greedy Approach"},
			CommonMistakes:       models.StringArray{"Not handling all negative arrays", "Resetting max instead of current sum"},
			DifficultyScore:      60.0,
			EstimatedTimeSeconds: intPtr(200),
			CreatedAt:            now,
			UpdatedAt:            now,
		},

		// Additional generic questions
		{
			QuestionType:   "concept_explanation",
			QuestionFormat: "text",
			QuestionText:   "Explain the trade-off between time complexity and space complexity. Give an example.",
			CorrectAnswer: jsonbMap(map[string]interface{}{
				"answer": []string{
					"Often we can use extra memory to speed up computation, like hash maps for O(1) lookup",
					"Memoization uses O(n) space to reduce time from exponential to polynomial",
					"Can trade space for time - store results to avoid recomputation",
				},
			}),
			Explanation:          "Time-space tradeoff is a fundamental concept where we use additional memory to reduce computation time. Classic examples include hash tables (O(n) space for O(1) lookup vs O(n) search) and dynamic programming (storing subproblem solutions to avoid recomputation).",
			RelatedConcepts:      models.StringArray{"Time Complexity", "Space Complexity", "Algorithm Design"},
			DifficultyScore:      35.0,
			EstimatedTimeSeconds: intPtr(150),
			CreatedAt:            now,
			UpdatedAt:            now,
		},
		{
			QuestionType:   "debugging",
			QuestionFormat: "text",
			QuestionText:   "What is an off-by-one error? Give an example from array problems.",
			CorrectAnswer: jsonbMap(map[string]interface{}{
				"answer": []string{
					"Off-by-one error is accessing index n instead of n-1, or starting at 1 instead of 0",
					"Using <= when should use < in loop, or vice versa",
					"Common example: for i in range(len(arr)) then accessing arr[i+1]",
				},
			}),
			Explanation:          "Off-by-one errors occur when loop boundaries or array indices are off by exactly one position. Common cases include using <= instead of <, forgetting arrays are 0-indexed, or not accounting for the last element properly.",
			RelatedConcepts:      models.StringArray{"Debugging", "Arrays", "Common Errors"},
			CommonMistakes:       models.StringArray{"Accessing arr[len(arr)]", "Using 1-based indexing in 0-based language"},
			DifficultyScore:      15.0,
			EstimatedTimeSeconds: intPtr(90),
			CreatedAt:            now,
			UpdatedAt:            now,
		},
	}
}

# Master Topic Reference & Workflow Guide
## Complete DSA Topics Taxonomy for Database Seeding & Question Generation

---

## Document Purpose

This document serves as the **single source of truth** for:

1. **Complete Topic Coverage** - Every DSA topic with all subtopics, patterns, and variations
2. **Database Seeding** - Structured data for populating PostgreSQL + Apache AGE
3. **Question Generation** - Templates and patterns for creating questions
4. **Graph Construction** - Relationship definitions between concepts
5. **Embedding Strategy** - What to embed and how for vector search
6. **RAG Pipeline** - Context retrieval and prompt engineering
7. **n8n Workflow** - Automation instructions and data flow

---

## Table of Contents

1. [Topic Taxonomy Structure](#taxonomy-structure)
2. [Complete Topic Coverage](#complete-topics)
3. [Data Schema Specifications](#data-schema)
4. [Graph Relationship Definitions](#graph-relationships)
5. [Embedding Generation Strategy](#embeddings)
6. [RAG Pipeline Configuration](#rag-pipeline)
7. [n8n Workflow Integration](#n8n-workflow)
8. [Question Generation Templates](#question-templates)

---

## 1. TOPIC TAXONOMY STRUCTURE {#taxonomy-structure}

### Hierarchy Levels

```
Level 0: Category (e.g., "Data Structures")
  ├─ Level 1: Topic (e.g., "Arrays")
  │   ├─ Level 2: Subtopic (e.g., "Two Pointers")
  │   │   ├─ Level 3: Pattern (e.g., "Opposite Direction")
  │   │   │   └─ Level 4: Variation (e.g., "Finding Pairs")
```

### Topic Node Schema

```json
{
  "topic_id": "unique_identifier",
  "name": "Topic Name",
  "category": "category_name",
  "level": 0-4,
  "parent_topic_id": "parent_id or null",
  "description": "Detailed description",
  "keywords": ["keyword1", "keyword2"],
  "difficulty_range": [min, max],
  "prerequisites": ["topic_id1", "topic_id2"],
  "related_topics": ["topic_id3", "topic_id4"],
  "common_patterns": [
    {
      "pattern_name": "Pattern Name",
      "description": "What it is",
      "when_to_use": "Conditions",
      "template_code": "code_snippet",
      "time_complexity": "O(...)",
      "space_complexity": "O(...)"
    }
  ],
  "common_pitfalls": [
    {
      "pitfall": "Description of mistake",
      "why_common": "Why students make this",
      "how_to_avoid": "Prevention strategy",
      "example": "Code or scenario"
    }
  ],
  "edge_cases": [
    {
      "case": "Description",
      "example_input": "input",
      "expected_behavior": "what should happen",
      "common_mistake": "what students do wrong"
    }
  ],
  "question_types": [
    "type1", "type2"
  ],
  "practice_problems": [
    {
      "leetcode_number": 1,
      "title": "Problem Title",
      "difficulty": "Easy/Medium/Hard",
      "url": "leetcode url"
    }
  ]
}
```

---

## 2. COMPLETE TOPIC COVERAGE {#complete-topics}

### CATEGORY 1: ARRAYS & STRINGS

#### Topic 1.1: Two Pointers

```json
{
  "topic_id": "two_pointers",
  "name": "Two Pointers",
  "category": "Arrays & Strings",
  "level": 2,
  "parent_topic_id": "arrays",
  "description": "Technique using two pointers moving through array/string simultaneously",
  "keywords": ["two pointers", "opposite direction", "same direction", "sliding window"],
  "difficulty_range": [10, 60],
  "prerequisites": ["arrays_basics", "loops"],
  "related_topics": ["sliding_window", "binary_search"],
  
  "common_patterns": [
    {
      "pattern_name": "Opposite Direction",
      "description": "Pointers start at opposite ends, move towards center",
      "when_to_use": "Finding pairs, palindromes, container problems, sorted array problems",
      "template_code": "int left = 0, right = n-1;\nwhile (left < right) {\n  if (condition) left++;\n  else right--;\n}",
      "time_complexity": "O(n)",
      "space_complexity": "O(1)",
      "examples": ["Two Sum II", "Container With Most Water", "Trapping Rain Water"]
    },
    {
      "pattern_name": "Same Direction (Slow/Fast)",
      "description": "Both pointers move forward, at different speeds",
      "when_to_use": "Remove duplicates, move elements, partition arrays",
      "template_code": "int slow = 0;\nfor (int fast = 0; fast < n; fast++) {\n  if (condition) {\n    arr[slow++] = arr[fast];\n  }\n}",
      "time_complexity": "O(n)",
      "space_complexity": "O(1)",
      "examples": ["Remove Duplicates", "Move Zeroes", "Remove Element"]
    },
    {
      "pattern_name": "Fast and Slow (Cycle Detection)",
      "description": "Fast moves 2 steps, slow moves 1 step",
      "when_to_use": "Linked list cycle detection, finding middle",
      "template_code": "ListNode *slow = head, *fast = head;\nwhile (fast && fast->next) {\n  slow = slow->next;\n  fast = fast->next->next;\n  if (slow == fast) return true;\n}",
      "time_complexity": "O(n)",
      "space_complexity": "O(1)",
      "examples": ["Linked List Cycle", "Find Middle of Linked List", "Happy Number"]
    }
  ],
  
  "common_pitfalls": [
    {
      "pitfall": "Off-by-one errors in pointer initialization",
      "why_common": "Confusion about inclusive/exclusive bounds",
      "how_to_avoid": "Always clarify: does right start at n-1 or n? Test with size 1 array",
      "example": "right = arr.size() causes out-of-bounds; should be arr.size() - 1"
    },
    {
      "pitfall": "Not handling empty arrays",
      "why_common": "Assume array has at least one element",
      "how_to_avoid": "Always check if (arr.empty()) return ...; at start",
      "example": "arr[0] crashes on empty array"
    },
    {
      "pitfall": "Infinite loop when pointers don't progress",
      "why_common": "Condition never changes or pointers stuck",
      "how_to_avoid": "Ensure at least one pointer moves in every iteration",
      "example": "while (left < right) without left++ or right-- inside"
    },
    {
      "pitfall": "Wrong condition for palindrome check",
      "why_common": "Confusion about when to stop (middle)",
      "how_to_avoid": "Use left < right, not left <= right (middle char doesn't need checking)",
      "example": "Checking middle character twice in odd-length string"
    }
  ],
  
  "edge_cases": [
    {
      "case": "Empty array/string",
      "example_input": "[]",
      "expected_behavior": "Return default value (0, false, empty array) or handle gracefully",
      "common_mistake": "Accessing arr[0] without checking size"
    },
    {
      "case": "Single element",
      "example_input": "[5]",
      "expected_behavior": "Pointers might start at same position or loop doesn't run",
      "common_mistake": "Loop condition left < right never true; returns wrong result"
    },
    {
      "case": "Two elements",
      "example_input": "[1, 2]",
      "expected_behavior": "Minimal case for pointer movement",
      "common_mistake": "Pointer update logic breaks (e.g., left = mid causes infinite loop)"
    },
    {
      "case": "All same elements",
      "example_input": "[3, 3, 3, 3]",
      "expected_behavior": "Should handle duplicates correctly",
      "common_mistake": "Logic assumes distinct elements"
    },
    {
      "case": "Already sorted (ascending)",
      "example_input": "[1, 2, 3, 4]",
      "expected_behavior": "Algorithms might have best-case performance",
      "common_mistake": "Not testing if optimization works"
    },
    {
      "case": "Reverse sorted (descending)",
      "example_input": "[4, 3, 2, 1]",
      "expected_behavior": "Worst-case for some algorithms",
      "common_mistake": "Assuming sorted means ascending"
    },
    {
      "case": "Target at boundaries",
      "example_input": "Target is first or last element",
      "expected_behavior": "Should be found without issues",
      "common_mistake": "Boundary conditions in loop exit prematurely"
    },
    {
      "case": "No solution exists",
      "example_input": "Looking for pair that sums to target, none exists",
      "expected_behavior": "Return -1, empty array, or false",
      "common_mistake": "Returning garbage value or not handling explicitly"
    }
  ],
  
  "question_types": [
    "implementation",
    "pattern_recognition",
    "edge_case_identification",
    "complexity_analysis",
    "code_debugging",
    "optimization"
  ],
  
  "practice_problems": [
    {
      "leetcode_number": 167,
      "title": "Two Sum II - Input Array Is Sorted",
      "difficulty": "Medium",
      "url": "https://leetcode.com/problems/two-sum-ii-input-array-is-sorted/",
      "pattern": "Opposite Direction",
      "key_insight": "Use sorted property to eliminate half search space each step"
    },
    {
      "leetcode_number": 15,
      "title": "3Sum",
      "difficulty": "Medium",
      "url": "https://leetcode.com/problems/3sum/",
      "pattern": "Fix one + Two Pointers",
      "key_insight": "Fix first element, use two pointers for remaining two"
    },
    {
      "leetcode_number": 11,
      "title": "Container With Most Water",
      "difficulty": "Medium",
      "url": "https://leetcode.com/problems/container-with-most-water/",
      "pattern": "Opposite Direction",
      "key_insight": "Always move pointer at shorter line"
    },
    {
      "leetcode_number": 26,
      "title": "Remove Duplicates from Sorted Array",
      "difficulty": "Easy",
      "url": "https://leetcode.com/problems/remove-duplicates-from-sorted-array/",
      "pattern": "Same Direction",
      "key_insight": "Slow pointer tracks position for unique elements"
    },
    {
      "leetcode_number": 283,
      "title": "Move Zeroes",
      "difficulty": "Easy",
      "url": "https://leetcode.com/problems/move-zeroes/",
      "pattern": "Same Direction",
      "key_insight": "Swap non-zero elements to front"
    }
  ],
  
  "variations": [
    {
      "name": "Three Pointers",
      "description": "Extension to three pointers for problems like Dutch National Flag",
      "example": "Sort Colors (LeetCode 75)"
    },
    {
      "name": "Two Pointers on Two Arrays",
      "description": "Merge two sorted arrays using pointer in each",
      "example": "Merge Sorted Array (LeetCode 88)"
    }
  ],
  
  "implementation_notes": {
    "python": {
      "syntax": "Common patterns in Python",
      "gotchas": ["list slicing creates copy", "negative indices wrap around"]
    },
    "cpp": {
      "syntax": "Common patterns in C++",
      "gotchas": ["pointer vs iterator", "size() returns unsigned (size_t)"]
    },
    "java": {
      "syntax": "Common patterns in Java",
      "gotchas": ["array length is final field", "ArrayList vs array"]
    }
  }
}
```

#### Topic 1.2: Sliding Window

```json
{
  "topic_id": "sliding_window",
  "name": "Sliding Window",
  "category": "Arrays & Strings",
  "level": 2,
  "parent_topic_id": "arrays",
  "description": "Maintain a window of elements that satisfies certain conditions, expanding and contracting as needed",
  "keywords": ["sliding window", "variable size", "fixed size", "substring", "subarray"],
  "difficulty_range": [20, 70],
  "prerequisites": ["two_pointers", "hash_table"],
  "related_topics": ["two_pointers", "dynamic_programming"],
  
  "common_patterns": [
    {
      "pattern_name": "Fixed Size Window",
      "description": "Window of constant size k slides through array",
      "when_to_use": "Find max/min in all windows of size k",
      "template_code": "// Process first window\nfor (int i = 0; i < k; i++) sum += arr[i];\n// Slide window\nfor (int i = k; i < n; i++) {\n  sum = sum - arr[i-k] + arr[i];\n  // update result\n}",
      "time_complexity": "O(n)",
      "space_complexity": "O(1)",
      "examples": ["Maximum Average Subarray", "Sliding Window Maximum"]
    },
    {
      "pattern_name": "Variable Size Window",
      "description": "Window expands/contracts based on validity condition",
      "when_to_use": "Find longest/shortest subarray satisfying condition",
      "template_code": "int left = 0;\nfor (int right = 0; right < n; right++) {\n  // Add arr[right] to window\n  while (invalid_condition) {\n    // Remove arr[left]\n    left++;\n  }\n  // Update result with window [left, right]\n}",
      "time_complexity": "O(n)",
      "space_complexity": "O(k) where k is unique elements in window",
      "examples": ["Longest Substring Without Repeating", "Minimum Window Substring"]
    }
  ],
  
  "common_pitfalls": [
    {
      "pitfall": "Confusing when to move left pointer",
      "why_common": "Not clear when window becomes invalid",
      "how_to_avoid": "Define 'valid window' explicitly before coding",
      "example": "For 'at most 2 distinct chars', left++ when distinct > 2, not >= 2"
    },
    {
      "pitfall": "Not updating window state correctly",
      "why_common": "Forgetting to remove left element when contracting",
      "how_to_avoid": "Always update data structures when left moves",
      "example": "Removing from hash set when left++: set.erase(s[left++])"
    },
    {
      "pitfall": "Off-by-one in window size calculation",
      "why_common": "Confusion about inclusive/exclusive bounds",
      "how_to_avoid": "Window size = right - left + 1 (both inclusive)",
      "example": "right=5, left=3 → size = 3 (not 2)"
    },
    {
      "pitfall": "Not handling single element windows",
      "why_common": "Edge case where left == right",
      "how_to_avoid": "Test with k=1 or minimal input",
      "example": "Window [5, 5] should have size 1 and contain arr[5]"
    }
  ],
  
  "edge_cases": [
    {
      "case": "Window size k > array length",
      "example_input": "arr = [1,2], k = 3",
      "expected_behavior": "Invalid input or return special value",
      "common_mistake": "Trying to process windows that don't exist"
    },
    {
      "case": "k = 0 or k = n",
      "example_input": "arr = [1,2,3], k = 0 or k = 3",
      "expected_behavior": "Handle as special cases",
      "common_mistake": "Loop logic breaks"
    },
    {
      "case": "All elements same",
      "example_input": "arr = [5,5,5,5], find distinct in window",
      "expected_behavior": "Only 1 distinct element",
      "common_mistake": "Logic assumes variety"
    },
    {
      "case": "Condition never satisfied",
      "example_input": "Find substring with sum >= target, all negative",
      "expected_behavior": "Return -1 or impossible indicator",
      "common_mistake": "Infinite loop or wrong result"
    },
    {
      "case": "Entire array is the answer",
      "example_input": "Find longest valid substring = whole string",
      "expected_behavior": "Return n",
      "common_mistake": "Window never contracts, result not updated"
    }
  ],
  
  "question_types": [
    "implementation",
    "validity_condition_definition",
    "window_state_management",
    "optimization",
    "edge_case_identification"
  ],
  
  "practice_problems": [
    {
      "leetcode_number": 3,
      "title": "Longest Substring Without Repeating Characters",
      "difficulty": "Medium",
      "pattern": "Variable Size Window",
      "key_insight": "Expand until duplicate, contract until valid"
    },
    {
      "leetcode_number": 76,
      "title": "Minimum Window Substring",
      "difficulty": "Hard",
      "pattern": "Variable Size Window",
      "key_insight": "Contract when valid to find minimum"
    },
    {
      "leetcode_number": 239,
      "title": "Sliding Window Maximum",
      "difficulty": "Hard",
      "pattern": "Fixed Size + Deque",
      "key_insight": "Use monotonic deque to track maximum"
    },
    {
      "leetcode_number": 424,
      "title": "Longest Repeating Character Replacement",
      "difficulty": "Medium",
      "pattern": "Variable Size Window",
      "key_insight": "Window size - max_freq <= k"
    }
  ]
}
```

#### Topic 1.3: Prefix Sum

```json
{
  "topic_id": "prefix_sum",
  "name": "Prefix Sum",
  "category": "Arrays & Strings",
  "level": 2,
  "parent_topic_id": "arrays",
  "description": "Precompute cumulative sums to answer range sum queries in O(1)",
  "keywords": ["prefix sum", "cumulative sum", "range query", "subarray sum"],
  "difficulty_range": [15, 50],
  "prerequisites": ["arrays_basics"],
  "related_topics": ["hash_table", "difference_array"],
  
  "common_patterns": [
    {
      "pattern_name": "Basic Prefix Sum",
      "description": "prefix[i] = sum of elements [0...i-1]",
      "when_to_use": "Multiple range sum queries on static array",
      "template_code": "vector<int> prefix(n+1, 0);\nfor (int i = 0; i < n; i++) {\n  prefix[i+1] = prefix[i] + arr[i];\n}\n// Query sum [L, R]: prefix[R+1] - prefix[L]",
      "time_complexity": "O(n) precompute, O(1) query",
      "space_complexity": "O(n)",
      "examples": ["Range Sum Query", "Subarray Sum Equals K"]
    },
    {
      "pattern_name": "2D Prefix Sum",
      "description": "Prefix sum for matrices",
      "when_to_use": "Range sum queries on 2D grid",
      "template_code": "prefix[i][j] = prefix[i-1][j] + prefix[i][j-1] - prefix[i-1][j-1] + matrix[i-1][j-1];\n// Query: prefix[r2][c2] - prefix[r1-1][c2] - prefix[r2][c1-1] + prefix[r1-1][c1-1]",
      "time_complexity": "O(mn) precompute, O(1) query",
      "space_complexity": "O(mn)",
      "examples": ["Range Sum Query 2D"]
    },
    {
      "pattern_name": "Prefix Sum + Hash Table",
      "description": "Find subarrays with specific sum using prefix sum and hash map",
      "when_to_use": "Count subarrays with sum = k, or find if exists",
      "template_code": "unordered_map<int, int> prefixCount;\nprefixCount[0] = 1;\nint sum = 0, count = 0;\nfor (int num : arr) {\n  sum += num;\n  count += prefixCount[sum - k];\n  prefixCount[sum]++;\n}",
      "time_complexity": "O(n)",
      "space_complexity": "O(n)",
      "examples": ["Subarray Sum Equals K", "Continuous Subarray Sum"]
    }
  ],
  
  "common_pitfalls": [
    {
      "pitfall": "Off-by-one in prefix array indexing",
      "why_common": "Confusion about 0-indexed vs 1-indexed prefix array",
      "how_to_avoid": "Use prefix[n+1], prefix[0] = 0, prefix[i] = sum[0...i-1]",
      "example": "To query sum[L, R], use prefix[R+1] - prefix[L], not prefix[R] - prefix[L-1]"
    },
    {
      "pitfall": "Integer overflow in prefix sum",
      "why_common": "Sum can exceed int range even if individual elements fit",
      "how_to_avoid": "Use long long for prefix array",
      "example": "10^5 elements of value 10^5 = 10^10 > INT_MAX"
    },
    {
      "pitfall": "Forgetting to initialize prefix[0] = 0",
      "why_common": "Miss the base case",
      "how_to_avoid": "Always set prefix[0] = 0 explicitly",
      "example": "Without this, queries starting at index 0 are wrong"
    },
    {
      "pitfall": "Modifying array after building prefix",
      "why_common": "Prefix sum is static snapshot",
      "how_to_avoid": "If array changes, rebuild prefix or use different structure (Fenwick Tree)",
      "example": "Update arr[3] = 10 doesn't update prefix[4], prefix[5], etc."
    }
  ],
  
  "edge_cases": [
    {
      "case": "Empty array",
      "example_input": "arr = []",
      "expected_behavior": "prefix = [0]",
      "common_mistake": "Not handling, causing index errors"
    },
    {
      "case": "Single element",
      "example_input": "arr = [5]",
      "expected_behavior": "prefix = [0, 5]",
      "common_mistake": "Off-by-one errors in queries"
    },
    {
      "case": "All negative elements",
      "example_input": "arr = [-5, -2, -8]",
      "expected_behavior": "Prefix sums are decreasing",
      "common_mistake": "Logic assumes positive sums"
    },
    {
      "case": "Sum equals zero",
      "example_input": "arr = [1, -1, 2, -2]",
      "expected_behavior": "prefix[2] = 0, prefix[4] = 0",
      "common_mistake": "Using prefix == 0 as special case incorrectly"
    },
    {
      "case": "Query range [0, 0]",
      "example_input": "Sum of single element at index 0",
      "expected_behavior": "prefix[1] - prefix[0] = arr[0]",
      "common_mistake": "Boundary handling"
    }
  ],
  
  "question_types": [
    "implementation",
    "range_query",
    "subarray_problems",
    "optimization",
    "hash_table_combination"
  ],
  
  "practice_problems": [
    {
      "leetcode_number": 303,
      "title": "Range Sum Query - Immutable",
      "difficulty": "Easy",
      "pattern": "Basic Prefix Sum",
      "key_insight": "Precompute prefix array once"
    },
    {
      "leetcode_number": 560,
      "title": "Subarray Sum Equals K",
      "difficulty": "Medium",
      "pattern": "Prefix Sum + Hash Map",
      "key_insight": "Use prefix sum differences"
    },
    {
      "leetcode_number": 523,
      "title": "Continuous Subarray Sum",
      "difficulty": "Medium",
      "pattern": "Prefix Sum + Modulo + Hash Map",
      "key_insight": "Track remainders of prefix sums"
    },
    {
      "leetcode_number": 304,
      "title": "Range Sum Query 2D - Immutable",
      "difficulty": "Medium",
      "pattern": "2D Prefix Sum",
      "key_insight": "Inclusion-exclusion principle"
    }
  ]
}
```

---

### CATEGORY 2: DYNAMIC PROGRAMMING

#### Topic 2.1: 1D DP (Linear)

```json
{
  "topic_id": "1d_dp",
  "name": "1D Dynamic Programming",
  "category": "Dynamic Programming",
  "level": 2,
  "parent_topic_id": "dynamic_programming",
  "description": "DP problems where state is one-dimensional array",
  "keywords": ["1d dp", "linear dp", "sequence", "fibonacci", "house robber"],
  "difficulty_range": [20, 60],
  "prerequisites": ["recursion", "memoization"],
  "related_topics": ["greedy", "kadanes_algorithm"],
  
  "common_patterns": [
    {
      "pattern_name": "Fibonacci-like",
      "description": "dp[i] depends on dp[i-1] and dp[i-2]",
      "when_to_use": "Two previous states determine current",
      "template_code": "dp[0] = base1, dp[1] = base2;\nfor (int i = 2; i <= n; i++) {\n  dp[i] = dp[i-1] + dp[i-2];\n}",
      "time_complexity": "O(n)",
      "space_complexity": "O(n), optimizable to O(1)",
      "examples": ["Climbing Stairs", "Fibonacci", "Tribonacci"]
    },
    {
      "pattern_name": "Take or Skip",
      "description": "At each position, decide to include or exclude element",
      "when_to_use": "Choosing subset with constraints",
      "template_code": "dp[i] = max(dp[i-1], dp[i-2] + arr[i]);",
      "time_complexity": "O(n)",
      "space_complexity": "O(n), optimizable to O(1)",
      "examples": ["House Robber", "Delete and Earn"]
    },
    {
      "pattern_name": "Maximum Subarray (Kadane's)",
      "description": "dp[i] = max sum ending at i",
      "when_to_use": "Finding best contiguous subarray",
      "template_code": "dp[i] = max(arr[i], dp[i-1] + arr[i]);\nmaxSum = max(maxSum, dp[i]);",
      "time_complexity": "O(n)",
      "space_complexity": "O(1)",
      "examples": ["Maximum Subarray", "Maximum Product Subarray"]
    }
  ],
  
  "common_pitfalls": [
    {
      "pitfall": "Wrong base case initialization",
      "why_common": "Not thinking through what dp[0] and dp[1] should be",
      "how_to_avoid": "Manually trace first few values",
      "example": "Climbing Stairs: dp[0]=1 (1 way to stay), dp[1]=1 (1 step), NOT dp[0]=0"
    },
    {
      "pitfall": "Not handling negative numbers in max subarray",
      "why_common": "Logic assumes positive values",
      "how_to_avoid": "Initialize with arr[0], not 0",
      "example": "All negative array: should return least negative, not 0"
    },
    {
      "pitfall": "Off-by-one in space optimization",
      "why_common": "Confusion when using variables instead of array",
      "how_to_avoid": "Use clear names: prev2, prev1, curr",
      "example": "Swapping prev1 and prev2 at wrong time"
    },
    {
      "pitfall": "Forgetting to update result variable",
      "why_common": "In problems where answer is max(dp[i]), not dp[n]",
      "how_to_avoid": "Check if answer is final state or max of all states",
      "example": "LIS: answer is max(dp[0..n]), not dp[n]"
    }
  ],
  
  "edge_cases": [
    {
      "case": "n = 0 or n = 1",
      "example_input": "Empty or single element",
      "expected_behavior": "Return base case directly",
      "common_mistake": "Array access out of bounds"
    },
    {
      "case": "All elements negative",
      "example_input": "arr = [-5, -2, -8, -1]",
      "expected_behavior": "Return best (least negative)",
      "common_mistake": "Returning 0 or not considering all"
    },
    {
      "case": "Overflow in calculations",
      "example_input": "Large Fibonacci numbers",
      "expected_behavior": "Use long long or modulo",
      "common_mistake": "Integer overflow gives wrong result"
    }
  ],
  
  "question_types": [
    "state_definition",
    "base_case_determination",
    "recurrence_relation",
    "space_optimization",
    "return_value_semantics"
  ],
  
  "practice_problems": [
    {
      "leetcode_number": 70,
      "title": "Climbing Stairs",
      "difficulty": "Easy",
      "pattern": "Fibonacci-like",
      "key_insight": "Ways to reach step i = ways to i-1 + ways to i-2"
    },
    {
      "leetcode_number": 198,
      "title": "House Robber",
      "difficulty": "Medium",
      "pattern": "Take or Skip",
      "key_insight": "Rob current house XOR rob previous house"
    },
    {
      "leetcode_number": 53,
      "title": "Maximum Subarray",
      "difficulty": "Medium",
      "pattern": "Kadane's Algorithm",
      "key_insight": "Either extend previous subarray or start new"
    },
    {
      "leetcode_number": 91,
      "title": "Decode Ways",
      "difficulty": "Medium",
      "pattern": "Fibonacci-like with conditions",
      "key_insight": "Sum ways if 1-digit OR 2-digit is valid"
    }
  ],
  
  "implementation_variations": {
    "bottom_up_array": {
      "description": "Standard DP with array",
      "space": "O(n)",
      "code": "vector<int> dp(n+1);"
    },
    "bottom_up_optimized": {
      "description": "Space optimized with variables",
      "space": "O(1)",
      "code": "int prev2 = base1, prev1 = base2;"
    },
    "top_down_memoization": {
      "description": "Recursive with cache",
      "space": "O(n)",
      "code": "int helper(int i, vector<int>& memo)"
    }
  }
}
```

#### Topic 2.2: 2D DP (Two Sequences)

```json
{
  "topic_id": "2d_dp",
  "name": "2D Dynamic Programming",
  "category": "Dynamic Programming",
  "level": 2,
  "parent_topic_id": "dynamic_programming",
  "description": "DP problems with two-dimensional state, often comparing two sequences",
  "keywords": ["2d dp", "grid", "two sequences", "lcs", "edit distance"],
  "difficulty_range": [30, 75],
  "prerequisites": ["1d_dp", "string_algorithms"],
  "related_topics": ["backtracking", "graph_traversal"],
  
  "common_patterns": [
    {
      "pattern_name": "Two String Comparison",
      "description": "dp[i][j] = answer for s1[0..i-1] and s2[0..j-1]",
      "when_to_use": "Matching, editing, transforming two strings",
      "template_code": "for (int i = 1; i <= m; i++) {\n  for (int j = 1; j <= n; j++) {\n    if (s1[i-1] == s2[j-1]) {\n      dp[i][j] = dp[i-1][j-1] + 1;\n    } else {\n      dp[i][j] = max(dp[i-1][j], dp[i][j-1]);\n    }\n  }\n}",
      "time_complexity": "O(m×n)",
      "space_complexity": "O(m×n), optimizable to O(n)",
      "examples": ["LCS", "Edit Distance", "Distinct Subsequences"]
    },
    {
      "pattern_name": "Grid Path",
      "description": "dp[i][j] = paths/cost to reach cell (i,j)",
      "when_to_use": "Robot movement, path counting, minimum cost path",
      "template_code": "dp[i][j] = min(dp[i-1][j], dp[i][j-1]) + grid[i][j];",
      "time_complexity": "O(m×n)",
      "space_complexity": "O(m×n), optimizable to O(n)",
      "examples": ["Unique Paths", "Minimum Path Sum", "Dungeon Game"]
    },
    {
      "pattern_name": "0/1 Knapsack",
      "description": "dp[i][w] = max value using first i items, capacity w",
      "when_to_use": "Subset selection with capacity constraint",
      "template_code": "for (int i = 1; i <= n; i++) {\n  for (int w = 0; w <= W; w++) {\n    if (weight[i-1] <= w) {\n      dp[i][w] = max(dp[i-1][w], dp[i-1][w-weight[i-1]] + value[i-1]);\n    } else {\n      dp[i][w] = dp[i-1][w];\n    }\n  }\n}",
      "time_complexity": "O(n×W)",
      "space_complexity": "O(n×W), optimizable to O(W)",
      "examples": ["Partition Equal Subset", "Target Sum", "Coin Change"]
    }
  ],
  
  "common_pitfalls": [
    {
      "pitfall": "Off-by-one in 2D array indexing",
      "why_common": "Confusion between 0-indexed strings and 1-indexed DP table",
      "how_to_avoid": "Use dp[m+1][n+1], access string as s[i-1] when at dp[i][j]",
      "example": "LCS: dp[i][j] represents s1[0..i-1] and s2[0..j-1]"
    },
    {
      "pitfall": "Wrong base case for first row/column",
      "why_common": "Not thinking through empty string cases",
      "how_to_avoid": "Manually fill first row and column before loops",
      "example": "Edit Distance: dp[i][0] = i (delete all), dp[0][j] = j (insert all)"
    },
    {
      "pitfall": "Space optimization breaks when need diagonal",
      "why_common": "Overwriting values before using them",
      "how_to_avoid": "Use temp variable to store dp[i-1][j-1]",
      "example": "LCS space optimization needs 'prev' variable"
    },
    {
      "pitfall": "Not considering all transition options",
      "why_common": "Missing one of the recurrence cases",
      "how_to_avoid": "List all possibilities: match, insert, delete, replace",
      "example": "Edit Distance: 3 operations if mismatch, 0 if match"
    }
  ],
  
  "edge_cases": [
    {
      "case": "One or both strings empty",
      "example_input": "s1 = '', s2 = 'abc' or both ''",
      "expected_behavior": "Base cases: LCS = 0, Edit Distance = length of non-empty",
      "common_mistake": "Not initializing first row/column"
    },
    {
      "case": "Strings are identical",
      "example_input": "s1 = 'abc', s2 = 'abc'",
      "expected_behavior": "LCS = 3, Edit Distance = 0",
      "common_mistake": "Logic doesn't optimize for this case"
    },
    {
      "case": "Completely different strings",
      "example_input": "s1 = 'abc', s2 = 'xyz'",
      "expected_behavior": "LCS = 0, Edit Distance = max(m, n)",
      "common_mistake": "Not handling properly"
    },
    {
      "case": "One string is substring of other",
      "example_input": "s1 = 'abc', s2 = 'aabbcc' contains abc",
      "expected_behavior": "LCS = 3",
      "common_mistake": "Logic treats it as no match"
    }
  ],
  
  "question_types": [
    "state_definition_2d",
    "base_case_initialization",
    "recurrence_relation",
    "space_optimization_2d_to_1d",
    "path_reconstruction"
  ],
  
  "practice_problems": [
    {
      "leetcode_number": 1143,
      "title": "Longest Common Subsequence",
      "difficulty": "Medium",
      "pattern": "Two String Comparison",
      "key_insight": "If match, add 1 to diagonal; else take max of top/left"
    },
    {
      "leetcode_number": 72,
      "title": "Edit Distance",
      "difficulty": "Medium",
      "pattern": "Two String Comparison",
      "key_insight": "3 operations: insert, delete, replace"
    },
    {
      "leetcode_number": 62,
      "title": "Unique Paths",
      "difficulty": "Medium",
      "pattern": "Grid Path",
      "key_insight": "paths[i][j] = paths[i-1][j] + paths[i][j-1]"
    },
    {
      "leetcode_number": 64,
      "title": "Minimum Path Sum",
      "difficulty": "Medium",
      "pattern": "Grid Path",
      "key_insight": "Take minimum of top or left, add current cell"
    }
  ],
  
  "space_optimization": {
    "from_2d_to_1d": {
      "condition": "Only need previous row",
      "technique": "Use vector<int> dp(n) instead of dp[m][n]",
      "example": "Unique Paths, LCS"
    },
    "rolling_array": {
      "condition": "Need two rows",
      "technique": "Use dp[2][n], alternate between dp[0] and dp[1]",
      "example": "When need current and previous row"
    },
    "diagonal_storage": {
      "condition": "Need diagonal element",
      "technique": "Store dp[i-1][j-1] in temp before overwriting",
      "example": "LCS space optimized"
    }
  }
}
```

---

### CATEGORY 3: GRAPHS & TREES

#### Topic 3.1: Graph Traversal (BFS)

```json
{
  "topic_id": "graph_bfs",
  "name": "Breadth-First Search (BFS)",
  "category": "Graphs & Trees",
  "level": 2,
  "parent_topic_id": "graph_traversal",
  "description": "Level-order traversal using queue, explores neighbors before going deeper",
  "keywords": ["bfs", "breadth first", "queue", "level order", "shortest path"],
  "difficulty_range": [15, 65],
  "prerequisites": ["queue", "graph_representation"],
  "related_topics": ["dfs", "dijkstra", "topological_sort"],
  
  "common_patterns": [
    {
      "pattern_name": "Standard BFS",
      "description": "Visit all nodes level by level",
      "when_to_use": "Shortest path in unweighted graph, level-order traversal",
      "template_code": "queue<int> q;\nvector<bool> visited(n, false);\nq.push(start);\nvisited[start] = true;\n\nwhile (!q.empty()) {\n  int node = q.front();\n  q.pop();\n  \n  for (int neighbor : adj[node]) {\n    if (!visited[neighbor]) {\n      visited[neighbor] = true;\n      q.push(neighbor);\n    }\n  }\n}",
      "time_complexity": "O(V + E)",
      "space_complexity": "O(V)",
      "examples": ["Shortest Path", "Word Ladder", "Rotting Oranges"]
    },
    {
      "pattern_name": "Level-by-Level BFS",
      "description": "Process each level separately",
      "when_to_use": "Need to track levels/steps explicitly",
      "template_code": "int level = 0;\nwhile (!q.empty()) {\n  int size = q.size();\n  for (int i = 0; i < size; i++) {\n    int node = q.front();\n    q.pop();\n    // Process node\n    for (int neighbor : adj[node]) {\n      if (!visited[neighbor]) {\n        visited[neighbor] = true;\n        q.push(neighbor);\n      }\n    }\n  }\n  level++;\n}",
      "time_complexity": "O(V + E)",
      "space_complexity": "O(V)",
      "examples": ["Binary Tree Level Order", "Minimum Depth"]
    },
    {
      "pattern_name": "Multi-Source BFS",
      "description": "Start BFS from multiple sources simultaneously",
      "when_to_use": "Find distance from ANY of multiple starting points",
      "template_code": "queue<int> q;\nfor (int source : sources) {\n  q.push(source);\n  visited[source] = true;\n}\n// Then standard BFS",
      "time_complexity": "O(V + E)",
      "space_complexity": "O(V)",
      "examples": ["Rotting Oranges", "Walls and Gates", "01 Matrix"]
    },
    {
      "pattern_name": "Bidirectional BFS",
      "description": "BFS from both start and end simultaneously",
      "when_to_use": "Optimize shortest path between two specific nodes",
      "template_code": "queue<int> q1, q2;\nset<int> visited1, visited2;\n// Expand smaller queue each iteration\n// Stop when frontiers meet",
      "time_complexity": "O(V + E), but faster in practice",
      "space_complexity": "O(V)",
      "examples": ["Word Ladder", "Minimum Genetic Mutation"]
    }
  ],
  
  "common_pitfalls": [
    {
      "pitfall": "Not marking visited before adding to queue",
      "why_common": "Think marking after popping is enough",
      "how_to_avoid": "ALWAYS mark visited when adding to queue, not when popping",
      "example": "If mark after pop, same node added multiple times → infinite loop or TLE"
    },
    {
      "pitfall": "Using DFS when need shortest path",
      "why_common": "Confusion about when to use BFS vs DFS",
      "how_to_avoid": "BFS guarantees shortest path in unweighted graphs, DFS doesn't",
      "example": "Finding shortest path in maze: BFS (correct), DFS (wrong - not guaranteed shortest)"
    },
    {
      "pitfall": "Not handling disconnected graphs",
      "why_common": "Assume graph is connected",
      "how_to_avoid": "Loop over all nodes, start BFS from unvisited ones",
      "example": "Counting connected components: need outer loop"
    },
    {
      "pitfall": "Modifying graph during traversal",
      "why_common": "Using grid itself to mark visited",
      "how_to_avoid": "Use separate visited array OR restore graph after",
      "example": "Grid problems: grid[i][j] = '#' to mark visited, restore if needed"
    }
  ],
  
  "edge_cases": [
    {
      "case": "Graph with no edges",
      "example_input": "n nodes, no connections",
      "expected_behavior": "Each node is its own component",
      "common_mistake": "Expecting some traversal"
    },
    {
      "case": "Single node graph",
      "example_input": "n = 1",
      "expected_behavior": "BFS immediately finishes",
      "common_mistake": "Off-by-one errors"
    },
    {
      "case": "Graph with self-loops",
      "example_input": "Node has edge to itself",
      "expected_behavior": "Visited check prevents revisiting",
      "common_mistake": "Infinite loop if not handling"
    },
    {
      "case": "Cyclic graph",
      "example_input": "Graph with cycles",
      "expected_behavior": "Visited array prevents infinite loop",
      "common_mistake": "Not using visited array correctly"
    },
    {
      "case": "Start node unreachable to target",
      "example_input": "Two disconnected components",
      "expected_behavior": "Return -1 or impossible",
      "common_mistake": "Infinite loop or wrong result"
    }
  ],
  
  "question_types": [
    "implementation",
    "shortest_path",
    "level_tracking",
    "multi_source",
    "optimization"
  ],
  
  "practice_problems": [
    {
      "leetcode_number": 102,
      "title": "Binary Tree Level Order Traversal",
      "difficulty": "Medium",
      "pattern": "Level-by-Level BFS",
      "key_insight": "Process each level separately"
    },
    {
      "leetcode_number": 127,
      "title": "Word Ladder",
      "difficulty": "Hard",
      "pattern": "Standard BFS on implicit graph",
      "key_insight": "Each word is node, edge if 1 char different"
    },
    {
      "leetcode_number": 994,
      "title": "Rotting Oranges",
      "difficulty": "Medium",
      "pattern": "Multi-Source BFS",
      "key_insight": "Start from all rotten oranges simultaneously"
    },
    {
      "leetcode_number": 1091,
      "title": "Shortest Path in Binary Matrix",
      "difficulty": "Medium",
      "pattern": "BFS on grid with 8 directions",
      "key_insight": "Can move in 8 directions, not just 4"
    }
  ],
  
  "comparison_with_dfs": {
    "when_use_bfs": [
      "Need shortest path (unweighted)",
      "Need to process level by level",
      "Graph might be very deep (DFS could overflow stack)"
    ],
    "when_use_dfs": [
      "Path finding (any path, not necessarily shortest)",
      "Detecting cycles",
      "Topological sorting",
      "Uses less space for sparse graphs"
    ]
  }
}
```

---

## 3. DATA SCHEMA SPECIFICATIONS {#data-schema}

### Complete Database Schema for n8n Workflow

```json
{
  "database_tables": {
    
    "topics": {
      "description": "Hierarchical topic structure",
      "schema": {
        "topic_id": "VARCHAR(100) PRIMARY KEY",
        "name": "VARCHAR(200) NOT NULL",
        "category": "VARCHAR(100)",
        "level": "INTEGER CHECK (level BETWEEN 0 AND 4)",
        "parent_topic_id": "VARCHAR(100) REFERENCES topics(topic_id)",
        "description": "TEXT",
        "keywords": "TEXT[]",
        "difficulty_range": "NUMRANGE",
        "prerequisites": "TEXT[]",
        "related_topics": "TEXT[]",
        "created_at": "TIMESTAMP DEFAULT CURRENT_TIMESTAMP"
      },
      "indexes": [
        "CREATE INDEX idx_topics_category ON topics(category)",
        "CREATE INDEX idx_topics_parent ON topics(parent_topic_id)",
        "CREATE INDEX idx_topics_keywords ON topics USING GIN(keywords)"
      ]
    },
    
    "patterns": {
      "description": "Reusable solution patterns",
      "schema": {
        "pattern_id": "SERIAL PRIMARY KEY",
        "topic_id": "VARCHAR(100) REFERENCES topics(topic_id)",
        "pattern_name": "VARCHAR(200) NOT NULL",
        "description": "TEXT",
        "when_to_use": "TEXT",
        "template_code_cpp": "TEXT",
        "template_code_python": "TEXT",
        "template_code_java": "TEXT",
        "time_complexity": "VARCHAR(50)",
        "space_complexity": "VARCHAR(50)",
        "examples": "TEXT[]",
        "created_at": "TIMESTAMP DEFAULT CURRENT_TIMESTAMP"
      },
      "indexes": [
        "CREATE INDEX idx_patterns_topic ON patterns(topic_id)",
        "CREATE INDEX idx_patterns_name ON patterns(pattern_name)"
      ]
    },
    
    "pitfalls": {
      "description": "Common mistakes and how to avoid them",
      "schema": {
        "pitfall_id": "SERIAL PRIMARY KEY",
        "topic_id": "VARCHAR(100) REFERENCES topics(topic_id)",
        "pitfall_description": "TEXT NOT NULL",
        "why_common": "TEXT",
        "how_to_avoid": "TEXT",
        "example_code": "TEXT",
        "severity": "VARCHAR(20) CHECK (severity IN ('low', 'medium', 'high', 'critical'))",
        "created_at": "TIMESTAMP DEFAULT CURRENT_TIMESTAMP"
      },
      "indexes": [
        "CREATE INDEX idx_pitfalls_topic ON pitfalls(topic_id)",
        "CREATE INDEX idx_pitfalls_severity ON pitfalls(severity)"
      ]
    },
    
    "edge_cases": {
      "description": "Edge cases to test for each topic",
      "schema": {
        "edge_case_id": "SERIAL PRIMARY KEY",
        "topic_id": "VARCHAR(100) REFERENCES topics(topic_id)",
        "case_description": "TEXT NOT NULL",
        "example_input": "TEXT",
        "expected_behavior": "TEXT",
        "common_mistake": "TEXT",
        "test_category": "VARCHAR(50)",
        "created_at": "TIMESTAMP DEFAULT CURRENT_TIMESTAMP"
      },
      "indexes": [
        "CREATE INDEX idx_edge_cases_topic ON edge_cases(topic_id)",
        "CREATE INDEX idx_edge_cases_category ON edge_cases(test_category)"
      ]
    },
    
    "problems": {
      "description": "LeetCode problems mapped to topics",
      "schema": {
        "problem_id": "SERIAL PRIMARY KEY",
        "leetcode_number": "INTEGER UNIQUE",
        "title": "VARCHAR(500) NOT NULL",
        "statement": "TEXT NOT NULL",
        "constraints": "TEXT[]",
        "examples": "JSONB NOT NULL",
        "hints": "TEXT[]",
        "difficulty_score": "FLOAT CHECK (difficulty_score BETWEEN 0 AND 100)",
        "difficulty_label": "VARCHAR(20)",
        "url": "VARCHAR(500)",
        "created_at": "TIMESTAMP DEFAULT CURRENT_TIMESTAMP",
        "search_vector": "tsvector GENERATED ALWAYS AS (to_tsvector('english', title || ' ' || statement)) STORED"
      },
      "indexes": [
        "CREATE INDEX idx_problems_difficulty ON problems(difficulty_score)",
        "CREATE INDEX idx_problems_leetcode ON problems(leetcode_number)",
        "CREATE INDEX idx_problems_search ON problems USING GIN(search_vector)"
      ]
    },
    
    "problem_topics": {
      "description": "Many-to-many: problems to topics",
      "schema": {
        "problem_id": "INTEGER REFERENCES problems(problem_id)",
        "topic_id": "VARCHAR(100) REFERENCES topics(topic_id)",
        "relevance_score": "FLOAT CHECK (relevance_score BETWEEN 0 AND 1)",
        "is_primary": "BOOLEAN DEFAULT false",
        "pattern_used": "VARCHAR(200)",
        "key_insight": "TEXT",
        "PRIMARY KEY": "(problem_id, topic_id)"
      },
      "indexes": [
        "CREATE INDEX idx_problem_topics_problem ON problem_topics(problem_id)",
        "CREATE INDEX idx_problem_topics_topic ON problem_topics(topic_id)"
      ]
    },
    
    "questions": {
      "description": "Generated questions for each problem",
      "schema": {
        "question_id": "SERIAL PRIMARY KEY",
        "problem_id": "INTEGER REFERENCES problems(problem_id)",
        "topic_id": "VARCHAR(100) REFERENCES topics(topic_id)",
        "question_type": "VARCHAR(50) NOT NULL",
        "question_text": "TEXT NOT NULL",
        "question_data": "JSONB",
        "answer_type": "VARCHAR(50)",
        "correct_answer": "JSONB NOT NULL",
        "answer_options": "JSONB",
        "explanation": "TEXT NOT NULL",
        "difficulty_score": "FLOAT",
        "hint_level_1": "TEXT",
        "hint_level_2": "TEXT",
        "hint_level_3": "TEXT",
        "tags": "TEXT[]",
        "created_at": "TIMESTAMP DEFAULT CURRENT_TIMESTAMP"
      },
      "indexes": [
        "CREATE INDEX idx_questions_problem ON questions(problem_id)",
        "CREATE INDEX idx_questions_topic ON questions(topic_id)",
        "CREATE INDEX idx_questions_type ON questions(question_type)",
        "CREATE INDEX idx_questions_tags ON questions USING GIN(tags)"
      ]
    },
    
    "embeddings": {
      "description": "Vector embeddings for semantic search",
      "schema": {
        "embedding_id": "SERIAL PRIMARY KEY",
        "entity_type": "VARCHAR(50) NOT NULL",
        "entity_id": "INTEGER NOT NULL",
        "embedding_vector": "vector(768)",
        "text_content": "TEXT",
        "model_name": "VARCHAR(100)",
        "created_at": "TIMESTAMP DEFAULT CURRENT_TIMESTAMP"
      },
      "indexes": [
        "CREATE INDEX idx_embeddings_entity ON embeddings(entity_type, entity_id)",
        "CREATE INDEX idx_embeddings_vector ON embeddings USING ivfflat (embedding_vector vector_cosine_ops) WITH (lists = 100)"
      ]
    }
  }
}
```

---

## 4. GRAPH RELATIONSHIP DEFINITIONS {#graph-relationships}

### Apache AGE Graph Schema

```json
{
  "graph_name": "dsa_knowledge_graph",
  
  "node_types": [
    {
      "label": "Topic",
      "properties": {
        "topic_id": "STRING (PRIMARY KEY)",
        "name": "STRING",
        "category": "STRING",
        "level": "INTEGER",
        "difficulty_range": "ARRAY"
      }
    },
    {
      "label": "Problem",
      "properties": {
        "problem_id": "INTEGER (PRIMARY KEY)",
        "leetcode_number": "INTEGER",
        "title": "STRING",
        "difficulty_score": "FLOAT"
      }
    },
    {
      "label": "Pattern",
      "properties": {
        "pattern_id": "INTEGER (PRIMARY KEY)",
        "name": "STRING",
        "time_complexity": "STRING",
        "space_complexity": "STRING"
      }
    },
    {
      "label": "Pitfall",
      "properties": {
        "pitfall_id": "INTEGER (PRIMARY KEY)",
        "description": "STRING",
        "severity": "STRING"
      }
    },
    {
      "label": "Concept",
      "properties": {
        "concept_id": "STRING (PRIMARY KEY)",
        "name": "STRING",
        "description": "STRING"
      }
    }
  ],
  
  "edge_types": [
    {
      "type": "SUBTOPIC_OF",
      "from": "Topic",
      "to": "Topic",
      "description": "Hierarchical topic relationship",
      "properties": {
        "level_difference": "INTEGER"
      },
      "example": "binary_search SUBTOPIC_OF searching"
    },
    {
      "type": "PREREQUISITE_OF",
      "from": "Topic",
      "to": "Topic",
      "description": "Topic A must be learned before Topic B",
      "properties": {
        "is_hard_requirement": "BOOLEAN",
        "strength": "FLOAT (0-1)"
      },
      "example": "arrays PREREQUISITE_OF two_pointers"
    },
    {
      "type": "RELATED_TO",
      "from": "Topic",
      "to": "Topic",
      "description": "Topics are related but not hierarchical",
      "properties": {
        "relationship_type": "STRING",
        "similarity_score": "FLOAT (0-1)"
      },
      "example": "bfs RELATED_TO dfs"
    },
    {
      "type": "HAS_TOPIC",
      "from": "Problem",
      "to": "Topic",
      "description": "Problem belongs to topic",
      "properties": {
        "relevance_score": "FLOAT (0-1)",
        "is_primary": "BOOLEAN"
      },
      "example": "Problem[1] HAS_TOPIC two_pointers"
    },
    {
      "type": "USES_PATTERN",
      "from": "Problem",
      "to": "Pattern",
      "description": "Problem uses solution pattern",
      "properties": {
        "is_optimal_solution": "BOOLEAN"
      },
      "example": "Problem[15] USES_PATTERN opposite_direction_pointers"
    },
    {
      "type": "SIMILAR_TO",
      "from": "Problem",
      "to": "Problem",
      "description": "Problems are similar in approach/structure",
      "properties": {
        "similarity_score": "FLOAT (0-1)",
        "similarity_reason": "STRING",
        "same_pattern": "BOOLEAN"
      },
      "example": "Problem[1] SIMILAR_TO Problem[167]"
    },
    {
      "type": "FOLLOW_UP_OF",
      "from": "Problem",
      "to": "Problem",
      "description": "Problem B is a natural follow-up to Problem A",
      "properties": {
        "difficulty_increase": "FLOAT",
        "new_constraints": "STRING[]",
        "builds_on_concept": "STRING"
      },
      "example": "Problem[15] FOLLOW_UP_OF Problem[1]"
    },
    {
      "type": "VARIATION_OF",
      "from": "Problem",
      "to": "Problem",
      "description": "Problem is a variation (different constraints/input)",
      "properties": {
        "variation_type": "STRING",
        "what_changed": "STRING"
      },
      "example": "Problem[16] VARIATION_OF Problem[15]"
    },
    {
      "type": "HARDER_VERSION",
      "from": "Problem",
      "to": "Problem",
      "description": "Problem A is harder version of Problem B",
      "properties": {
        "difficulty_delta": "FLOAT",
        "added_complexity": "STRING"
      },
      "example": "Problem[18] HARDER_VERSION Problem[15]"
    },
    {
      "type": "HAS_PITFALL",
      "from": "Topic",
      "to": "Pitfall",
      "description": "Topic has common mistake",
      "properties": {
        "frequency": "STRING (rare/common/very_common)"
      },
      "example": "binary_search HAS_PITFALL overflow_in_mid_calculation"
    },
    {
      "type": "REQUIRES_CONCEPT",
      "from": "Topic",
      "to": "Concept",
      "description": "Topic requires understanding of concept",
      "properties": {
        "importance": "STRING (optional/helpful/required)"
      },
      "example": "dp_2d REQUIRES_CONCEPT state_definition"
    },
    {
      "type": "DEMONSTRATES",
      "from": "Problem",
      "to": "Pitfall",
      "description": "Problem commonly triggers this pitfall",
      "properties": {
        "how_it_appears": "STRING"
      },
      "example": "Problem[35] DEMONSTRATES off_by_one_in_binary_search"
    }
  ],
  
  "common_queries": {
    "learning_path": {
      "description": "Find all prerequisites for a topic",
      "cypher": "MATCH path = (target:Topic {topic_id: $topic_id})<-[:PREREQUISITE_OF*]-(prereq:Topic) RETURN path ORDER BY length(path)"
    },
    "similar_problems": {
      "description": "Find problems similar to given problem",
      "cypher": "MATCH (p:Problem {problem_id: $problem_id})-[s:SIMILAR_TO]->(similar:Problem) WHERE s.similarity_score > $threshold RETURN similar ORDER BY s.similarity_score DESC"
    },
    "problems_by_topic": {
      "description": "Get all problems for a topic with their patterns",
      "cypher": "MATCH (t:Topic {topic_id: $topic_id})<-[:HAS_TOPIC]-(p:Problem)-[:USES_PATTERN]->(pat:Pattern) RETURN p, pat"
    },
    "follow_up_chain": {
      "description": "Find follow-up progression from a problem",
      "cypher": "MATCH path = (start:Problem {problem_id: $problem_id})-[:FOLLOW_UP_OF*1..5]->(followup:Problem) RETURN path ORDER BY length(path)"
    },
    "common_pitfalls_for_topic": {
      "description": "Get all pitfalls for a topic ordered by severity",
      "cypher": "MATCH (t:Topic {topic_id: $topic_id})-[:HAS_PITFALL]->(pit:Pitfall) RETURN pit ORDER BY pit.severity DESC"
    },
    "multi_pattern_problems": {
      "description": "Find problems using multiple patterns",
      "cypher": "MATCH (p:Problem)-[:USES_PATTERN]->(pat:Pattern) WITH p, collect(pat.name) as patterns WHERE size(patterns) >= 2 RETURN p, patterns"
    },
    "topic_difficulty_progression": {
      "description": "Order problems within topic by difficulty",
      "cypher": "MATCH (t:Topic {topic_id: $topic_id})<-[:HAS_TOPIC]-(p:Problem) RETURN p ORDER BY p.difficulty_score"
    }
  }
}
```

---

## 5. EMBEDDING GENERATION STRATEGY {#embeddings}

### What to Embed and How

```json
{
  "embedding_strategy": {
    
    "problems": {
      "what_to_embed": "title + statement + constraints + examples",
      "format": "Problem: {title}\n\nDescription: {statement}\n\nConstraints:\n{constraints}\n\nExample:\n{example_1}",
      "model": "text-embedding-3-small OR all-MiniLM-L6-v2",
      "dimensions": 768,
      "chunking": "No chunking - embed entire problem as one",
      "use_case": "Semantic problem search, similar problem discovery",
      "example": {
        "problem_id": 1,
        "text": "Problem: Two Sum\n\nDescription: Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.\n\nConstraints:\n- 2 <= nums.length <= 10^4\n- -10^9 <= nums[i] <= 10^9\n\nExample:\nInput: nums = [2,7,11,15], target = 9\nOutput: [0,1]",
        "embedding_vector": "[0.023, -0.142, ...]"
      }
    },
    
    "topics": {
      "what_to_embed": "name + description + keywords + common patterns",
      "format": "Topic: {name}\n\nDescription: {description}\n\nKey Concepts: {keywords}\n\nCommon Patterns: {pattern_names}",
      "model": "text-embedding-3-small",
      "dimensions": 768,
      "use_case": "Topic recommendation, concept search",
      "example": {
        "topic_id": "two_pointers",
        "text": "Topic: Two Pointers\n\nDescription: Technique using two pointers moving through array/string simultaneously\n\nKey Concepts: opposite direction, same direction, sliding window\n\nCommon Patterns: Opposite Direction, Same Direction (Slow/Fast), Fast and Slow (Cycle Detection)",
        "embedding_vector": "[0.156, 0.089, ...]"
      }
    },
    
    "patterns": {
      "what_to_embed": "pattern_name + description + when_to_use + template_code",
      "format": "Pattern: {pattern_name}\n\nDescription: {description}\n\nWhen to Use: {when_to_use}\n\nTemplate:\n{template_code}",
      "model": "code-embedding-ada-002 OR codebert-base",
      "dimensions": 768,
      "use_case": "Pattern matching, solution template search",
      "example": {
        "pattern_id": 123,
        "text": "Pattern: Opposite Direction\n\nDescription: Pointers start at opposite ends, move towards center\n\nWhen to Use: Finding pairs, palindromes, container problems, sorted array problems\n\nTemplate:\nint left = 0, right = n-1;\nwhile (left < right) {\n  if (condition) left++;\n  else right--;\n}",
        "embedding_vector": "[-0.234, 0.456, ...]"
      }
    },
    
    "questions": {
      "what_to_embed": "question_text + context",
      "format": "Question: {question_text}\n\nContext: {topic_name} - {problem_title}",
      "model": "text-embedding-3-small",
      "dimensions": 768,
      "use_case": "Question similarity, adaptive testing",
      "example": {
        "question_id": 456,
        "text": "Question: What is the time complexity of the optimal solution for Two Sum using a hash table?\n\nContext: Hash Table - Two Sum",
        "embedding_vector": "[0.091, -0.267, ...]"
      }
    },
    
    "explanations": {
      "what_to_embed": "concept explanation + examples",
      "format": "Concept: {concept_name}\n\n{explanation}\n\nExamples:\n{examples}",
      "model": "text-embedding-3-small",
      "dimensions": 768,
      "use_case": "RAG context retrieval for AI hints",
      "example": {
        "concept": "Binary Search Overflow",
        "text": "Concept: Integer Overflow in Binary Search\n\nWhen calculating the middle index in binary search using (left + right) / 2, integer overflow can occur if left and right are both large numbers near INT_MAX. The sum left + right may exceed the maximum value an integer can hold, causing undefined behavior.\n\nExamples:\n- left = 2147483646, right = 2147483647 → left + right overflows\n- Solution: Use left + (right - left) / 2 instead",
        "embedding_vector": "[0.334, -0.123, ...]"
      }
    }
  },
  
  "embedding_sync_strategy": {
    "initial_generation": {
      "description": "Generate all embeddings during database seeding",
      "workflow": [
        "1. Insert problem into problems table",
        "2. Generate embedding from problem text",
        "3. Insert embedding into embeddings table",
        "4. Also insert into ChromaDB collection"
      ]
    },
    "incremental_updates": {
      "description": "Update embeddings when content changes",
      "trigger": "PostgreSQL trigger on UPDATE",
      "workflow": [
        "1. Detect change in problems/topics/patterns table",
        "2. Regenerate embedding for changed entity",
        "3. Update embeddings table",
        "4. Update ChromaDB collection"
      ]
    },
    "batch_processing": {
      "description": "Process embeddings in batches for efficiency",
      "batch_size": 100,
      "rate_limit": "3000 embeddings per minute (API limit)",
      "retry_strategy": "Exponential backoff on failure"
    }
  },
  
  "similarity_search": {
    "query_types": {
      "problem_similarity": {
        "description": "Find similar problems by embedding",
        "threshold": 0.85,
        "max_results": 10,
        "query": "SELECT p.*, 1 - (e1.embedding_vector <=> e2.embedding_vector) as similarity FROM problems p JOIN embeddings e1 ON e1.entity_id = p.problem_id AND e1.entity_type = 'problem' CROSS JOIN embeddings e2 WHERE e2.entity_id = $problem_id AND e2.entity_type = 'problem' AND 1 - (e1.embedding_vector <=> e2.embedding_vector) > $threshold ORDER BY similarity DESC LIMIT $max_results"
      },
      "semantic_search": {
        "description": "Search problems by natural language query",
        "workflow": [
          "1. Generate embedding for user query",
          "2. Search embeddings table with cosine similarity",
          "3. Return top-k most similar problems"
        ],
        "example_query": "Find problems about finding pairs in sorted arrays",
        "expected_results": ["Two Sum II", "3Sum", "4Sum"]
      }
    }
  }
}
```

---

## 6. RAG PIPELINE CONFIGURATION {#rag-pipeline}

### Context Retrieval and Prompt Engineering

```json
{
  "rag_pipeline": {
    
    "components": {
      "vector_store": {
        "type": "ChromaDB",
        "host": "localhost:8000",
        "collections": [
          "problems",
          "topics",
          "patterns",
          "explanations",
          "questions"
        ],
        "distance_metric": "cosine",
        "index_type": "HNSW"
      },
      "embedding_model": {
        "name": "all-MiniLM-L6-v2",
        "dimensions": 384,
        "max_seq_length": 512
      },
      "llm": {
        "type": "Ollama",
        "model": "mistral:7b",
        "host": "localhost:11434",
        "context_window": 8192,
        "temperature": 0.3,
        "top_p": 0.9
      }
    },
    
    "retrieval_strategies": {
      
      "hint_generation": {
        "description": "Retrieve context for generating hints",
        "workflow": [
          "1. User answers question incorrectly",
          "2. Embed question + user's wrong answer",
          "3. Retrieve top-3 similar explanations from 'explanations' collection",
          "4. Retrieve top-2 similar questions from 'questions' collection",
          "5. Retrieve pattern information if pattern-related",
          "6. Build context for LLM"
        ],
        "retrieval_query": {
          "collection": "explanations",
          "query_text": "{question_text} {user_answer} {correct_answer}",
          "n_results": 3,
          "where": {
            "topic_id": "$topic_id"
          }
        },
        "context_assembly": {
          "template": "Question: {question_text}\n\nUser's Answer: {user_answer}\n\nCorrect Answer: {correct_answer}\n\nRelated Concepts:\n{retrieved_explanations}\n\nSimilar Questions:\n{retrieved_questions}",
          "max_context_length": 2000
        }
      },
      
      "problem_recommendation": {
        "description": "Recommend next problem based on current performance",
        "workflow": [
          "1. Analyze user's recent attempts and weak topics",
          "2. Embed weak topic descriptions",
          "3. Retrieve problems from 'problems' collection matching weak topics",
          "4. Filter by difficulty range (user_level ± 10)",
          "5. Exclude already solved problems",
          "6. Rank by relevance and difficulty appropriateness"
        ],
        "retrieval_query": {
          "collection": "problems",
          "query_text": "{weak_topic_description}",
          "n_results": 20,
          "where": {
            "difficulty_score": {
              "$gte": "$user_level - 10",
              "$lte": "$user_level + 10"
            }
          }
        }
      },
      
      "assessment_analysis": {
        "description": "Analyze assessment results using RAG",
        "workflow": [
          "1. Collect all user's answers from assessment",
          "2. Identify patterns in mistakes",
          "3. Embed mistake patterns",
          "4. Retrieve similar pitfalls and explanations",
          "5. Generate personalized feedback"
        ],
        "retrieval_query": {
          "collection": "pitfalls",
          "query_text": "{mistake_pattern_description}",
          "n_results": 5
        },
        "llm_prompt_template": "Based on this student's assessment results:\n\n{assessment_summary}\n\nAnd these related common pitfalls:\n\n{retrieved_pitfalls}\n\nProvide:\n1. Specific weaknesses (2-3)\n2. Root causes\n3. Actionable recommendations\n4. Practice plan\n\nBe specific and constructive."
      },
      
      "solution_explanation": {
        "description": "Generate detailed solution explanation",
        "workflow": [
          "1. Retrieve problem statement",
          "2. Retrieve pattern(s) used in solution",
          "3. Retrieve similar solved problems",
          "4. Generate step-by-step explanation"
        ],
        "retrieval_queries": [
          {
            "collection": "patterns",
            "query_text": "{problem_statement}",
            "n_results": 2
          },
          {
            "collection": "problems",
            "query_text": "{problem_statement}",
            "n_results": 3,
            "where": {
              "has_solution": true
            }
          }
        ],
        "llm_prompt_template": "Problem:\n{problem_statement}\n\nRelevant Patterns:\n{retrieved_patterns}\n\nSimilar Problems:\n{retrieved_similar_problems}\n\nProvide a clear, step-by-step solution explanation:\n1. Key observations\n2. Approach\n3. Algorithm\n4. Complexity analysis\n5. Common mistakes to avoid"
      }
    },
    
    "prompt_templates": {
      
      "hint_level_1": {
        "template": "A student is struggling with this question:\n\n{question_text}\n\nThey answered: {user_answer}\nCorrect answer: {correct_answer}\n\nContext:\n{retrieved_context}\n\nProvide a GENTLE HINT (Level 1):\n- Ask a Socratic question\n- Point to a related concept\n- Don't reveal the answer\n- Be encouraging",
        "max_tokens": 150,
        "temperature": 0.7
      },
      
      "hint_level_2": {
        "template": "A student is struggling with this question (2nd attempt):\n\n{question_text}\n\nThey answered: {user_answer}\nCorrect answer: {correct_answer}\n\nContext:\n{retrieved_context}\n\nProvide a MORE DIRECT HINT (Level 2):\n- Point to specific part of the problem\n- Mention a key insight\n- Still don't give away the answer\n- Be supportive",
        "max_tokens": 200,
        "temperature": 0.6
      },
      
      "hint_level_3": {
        "template": "A student needs help with this question (3rd+ attempt):\n\n{question_text}\n\nThey answered: {user_answer}\nCorrect answer: {correct_answer}\n\nContext:\n{retrieved_context}\n\nProvide a DETAILED HINT (Level 3):\n- Explain the concept needed\n- Walk through reasoning\n- Let them apply it themselves\n- Be patient and clear",
        "max_tokens": 300,
        "temperature": 0.5
      },
      
      "weakness_analysis": {
        "template": "Analyze this student's performance:\n\nTotal Questions: {total_questions}\nCorrect: {correct_count}\nCategories:\n{category_breakdown}\n\nDetailed Mistakes:\n{mistake_details}\n\nRelevant Context:\n{retrieved_pitfalls}\n\nProvide analysis in JSON format:\n{\n  \"weaknesses\": [\n    {\n      \"area\": \"specific weakness\",\n      \"severity\": \"low/medium/high\",\n      \"root_cause\": \"why\",\n      \"evidence\": \"specific examples\"\n    }\n  ],\n  \"recommendations\": [\n    {\n      \"action\": \"specific action\",\n      \"priority\": 1-5,\n      \"estimated_time\": \"time estimate\"\n    }\n  ]\n}",
        "max_tokens": 800,
        "temperature": 0.3,
        "output_format": "json"
      },
      
      "problem_variant_generation": {
        "template": "Given this base problem:\n\n{problem_statement}\n\nAnd these similar problems:\n{retrieved_similar}\n\nGenerate a {variant_type} variant:\n\nVariant Types:\n- SIMILAR: Same pattern, different scenario\n- HARDER: Add constraints, increase difficulty by 15-20\n- EASIER: Simplify, decrease difficulty by 15-20\n- FOLLOW_UP: Build on base problem\n\nReturn JSON:\n{\n  \"title\": \"new problem title\",\n  \"statement\": \"problem statement\",\n  \"constraints\": [\"constraint1\", \"constraint2\"],\n  \"examples\": [{\"input\": \"...\", \"output\": \"...\", \"explanation\": \"...\"}],\n  \"hints\": [\"hint1\", \"hint2\"],\n  \"difficulty_score\": 0-100,\n  \"relation_to_base\": \"how it relates\"\n}",
        "max_tokens": 1000,
        "temperature": 0.7,
        "output_format": "json"
      }
    },
    
    "context_optimization": {
      "max_context_tokens": 3000,
      "truncation_strategy": "Keep most relevant, truncate least relevant",
      "relevance_scoring": {
        "recency": 0.2,
        "similarity_score": 0.5,
        "user_history_match": 0.3
      }
    }
  }
}
```

---

## 7. N8N WORKFLOW INTEGRATION {#n8n-workflow}

### Automated Workflows for Database Seeding and Management

```json
{
  "n8n_workflows": {
    
    "workflow_1_seed_topics": {
      "name": "Seed Topics from Master Reference",
      "description": "Parse this master reference document and populate topics table",
      "trigger": "Manual or Scheduled",
      "nodes": [
        {
          "node_type": "Start",
          "id": "start_node"
        },
        {
          "node_type": "Code",
          "id": "parse_topics",
          "description": "Parse topics from JSON sections of this document",
          "code": "// Extract topics from master reference\nconst topics = JSON.parse($input.item.json.master_reference);\nreturn topics.map(topic => ({\n  json: {\n    topic_id: topic.topic_id,\n    name: topic.name,\n    category: topic.category,\n    level: topic.level,\n    parent_topic_id: topic.parent_topic_id,\n    description: topic.description,\n    keywords: topic.keywords,\n    difficulty_range: topic.difficulty_range,\n    prerequisites: topic.prerequisites,\n    related_topics: topic.related_topics\n  }\n}));"
        },
        {
          "node_type": "PostgreSQL",
          "id": "insert_topics",
          "operation": "executeQuery",
          "query": "INSERT INTO topics (topic_id, name, category, level, parent_topic_id, description, keywords, difficulty_range, prerequisites, related_topics) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (topic_id) DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description",
          "parameters": "topic_id, name, category, level, parent_topic_id, description, keywords, difficulty_range, prerequisites, related_topics"
        },
        {
          "node_type": "Code",
          "id": "log_success",
          "code": "console.log(`Inserted ${$input.all().length} topics`);\nreturn $input.all();"
        }
      ],
      "execution_order": [
        "start_node",
        "parse_topics",
        "insert_topics",
        "log_success"
      ]
    },
    
    "workflow_2_seed_patterns": {
      "name": "Seed Patterns",
      "description": "Extract and insert patterns from topics",
      "trigger": "After workflow_1 completes",
      "nodes": [
        {
          "node_type": "Start",
          "id": "start_patterns"
        },
        {
          "node_type": "PostgreSQL",
          "id": "get_topics",
          "operation": "executeQuery",
          "query": "SELECT * FROM topics"
        },
        {
          "node_type": "Code",
          "id": "extract_patterns",
          "description": "Extract patterns from each topic's common_patterns field",
          "code": "const patterns = [];\n$input.all().forEach(item => {\n  const topic = item.json;\n  if (topic.common_patterns) {\n    topic.common_patterns.forEach(pattern => {\n      patterns.push({\n        json: {\n          topic_id: topic.topic_id,\n          pattern_name: pattern.pattern_name,\n          description: pattern.description,\n          when_to_use: pattern.when_to_use,\n          template_code_cpp: pattern.template_code,\n          time_complexity: pattern.time_complexity,\n          space_complexity: pattern.space_complexity,\n          examples: pattern.examples\n        }\n      });\n    });\n  }\n});\nreturn patterns;"
        },
        {
          "node_type": "PostgreSQL",
          "id": "insert_patterns",
          "operation": "executeQuery",
          "query": "INSERT INTO patterns (topic_id, pattern_name, description, when_to_use, template_code_cpp, time_complexity, space_complexity, examples) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
          "parameters": "topic_id, pattern_name, description, when_to_use, template_code_cpp, time_complexity, space_complexity, examples"
        }
      ]
    },
    
    "workflow_3_seed_pitfalls": {
      "name": "Seed Pitfalls and Edge Cases",
      "description": "Extract pitfalls and edge cases from topics",
      "trigger": "After workflow_2 completes",
      "nodes": [
        {
          "node_type": "Start"
        },
        {
          "node_type": "PostgreSQL",
          "id": "get_topics_with_pitfalls",
          "operation": "executeQuery",
          "query": "SELECT * FROM topics"
        },
        {
          "node_type": "Code",
          "id": "extract_pitfalls",
          "code": "const pitfalls = [];\n$input.all().forEach(item => {\n  const topic = item.json;\n  if (topic.common_pitfalls) {\n    topic.common_pitfalls.forEach(pitfall => {\n      pitfalls.push({\n        json: {\n          topic_id: topic.topic_id,\n          pitfall_description: pitfall.pitfall,\n          why_common: pitfall.why_common,\n          how_to_avoid: pitfall.how_to_avoid,\n          example_code: pitfall.example,\n          severity: determineSeverity(pitfall)\n        }\n      });\n    });\n  }\n});\nreturn pitfalls;\n\nfunction determineSeverity(pitfall) {\n  // Logic to determine severity based on keywords\n  const critical = ['overflow', 'crash', 'infinite loop'];\n  const high = ['wrong result', 'TLE', 'incorrect'];\n  if (critical.some(word => pitfall.pitfall.toLowerCase().includes(word))) return 'critical';\n  if (high.some(word => pitfall.pitfall.toLowerCase().includes(word))) return 'high';\n  return 'medium';\n}"
        },
        {
          "node_type": "PostgreSQL",
          "id": "insert_pitfalls",
          "operation": "executeQuery",
          "query": "INSERT INTO pitfalls (topic_id, pitfall_description, why_common, how_to_avoid, example_code, severity) VALUES ($1, $2, $3, $4, $5, $6)"
        },
        {
          "node_type": "Code",
          "id": "extract_edge_cases",
          "code": "const edgeCases = [];\n$input.all().forEach(item => {\n  const topic = item.json;\n  if (topic.edge_cases) {\n    topic.edge_cases.forEach(edge => {\n      edgeCases.push({\n        json: {\n          topic_id: topic.topic_id,\n          case_description: edge.case,\n          example_input: edge.example_input,\n          expected_behavior: edge.expected_behavior,\n          common_mistake: edge.common_mistake,\n          test_category: categorizeEdgeCase(edge)\n        }\n      });\n    });\n  }\n});\nreturn edgeCases;\n\nfunction categorizeEdgeCase(edge) {\n  const categories = {\n    'empty': ['empty', 'null', 'zero size'],\n    'boundary': ['first', 'last', 'single', 'two'],\n    'extreme': ['maximum', 'minimum', 'overflow'],\n    'special': ['all same', 'sorted', 'reverse']\n  };\n  \n  for (const [category, keywords] of Object.entries(categories)) {\n    if (keywords.some(word => edge.case.toLowerCase().includes(word))) {\n      return category;\n    }\n  }\n  return 'other';\n}"
        },
        {
          "node_type": "PostgreSQL",
          "id": "insert_edge_cases",
          "operation": "executeQuery",
          "query": "INSERT INTO edge_cases (topic_id, case_description, example_input, expected_behavior, common_mistake, test_category) VALUES ($1, $2, $3, $4, $5, $6)"
        }
      ]
    },
    
    "workflow_4_generate_embeddings": {
      "name": "Generate Embeddings for All Content",
      "description": "Create vector embeddings for semantic search",
      "trigger": "After all content is seeded",
      "nodes": [
        {
          "node_type": "Start"
        },
        {
          "node_type": "PostgreSQL",
          "id": "get_problems",
          "operation": "executeQuery",
          "query": "SELECT problem_id, title, statement, constraints, examples FROM problems WHERE NOT EXISTS (SELECT 1 FROM embeddings WHERE entity_type = 'problem' AND entity_id = problem_id)"
        },
        {
          "node_type": "Code",
          "id": "prepare_problem_text",
          "code": "return $input.all().map(item => {\n  const p = item.json;\n  const text = `Problem: ${p.title}\\n\\nDescription: ${p.statement}\\n\\nConstraints:\\n${p.constraints.join('\\n')}\\n\\nExample:\\n${JSON.stringify(p.examples[0])}`;\n  return {\n    json: {\n      entity_id: p.problem_id,\n      entity_type: 'problem',\n      text_content: text\n    }\n  };\n});"
        },
        {
          "node_type": "HTTP Request",
          "id": "generate_embedding",
          "description": "Call embedding model API or local model",
          "method": "POST",
          "url": "http://localhost:8000/embed",
          "body": "{{ json.text_content }}",
          "batch_size": 100,
          "retry_on_fail": true,
          "max_retries": 3
        },
        {
          "node_type": "PostgreSQL",
          "id": "save_embedding",
          "operation": "executeQuery",
          "query": "INSERT INTO embeddings (entity_type, entity_id, embedding_vector, text_content, model_name) VALUES ($1, $2, $3::vector, $4, $5)",
          "parameters": "entity_type, entity_id, embedding_vector, text_content, 'all-MiniLM-L6-v2'"
        },
        {
          "node_type": "HTTP Request",
          "id": "add_to_chromadb",
          "description": "Also add to ChromaDB for faster retrieval",
          "method": "POST",
          "url": "http://localhost:8000/api/v1/collections/problems/add",
          "body": "{\n  \"ids\": [\"problem_{{ json.entity_id }}\"],\n  \"embeddings\": [{{ json.embedding_vector }}],\n  \"metadatas\": [{\"problem_id\": {{ json.entity_id }}, \"entity_type\": \"problem\"}],\n  \"documents\": [\"{{ json.text_content }}\"]\n}"
        },
        {
          "node_type": "Loop",
          "id": "repeat_for_topics",
          "description": "Repeat similar process for topics, patterns, questions"
        }
      ]
    },
    
    "workflow_5_build_graph": {
      "name": "Construct Knowledge Graph",
      "description": "Create nodes and edges in Apache AGE",
      "trigger": "After embeddings are generated",
      "nodes": [
        {
          "node_type": "Start"
        },
        {
          "node_type": "PostgreSQL",
          "id": "create_topic_nodes",
          "operation": "executeQuery",
          "query": "SELECT * FROM cypher('dsa_knowledge_graph', $$ CREATE (t:Topic {topic_id: $topic_id, name: $name, category: $category, level: $level, difficulty_range: $difficulty_range}) $$) as (v agtype)",
          "description": "Create Topic nodes"
        },
        {
          "node_type": "PostgreSQL",
          "id": "create_problem_nodes",
          "operation": "executeQuery",
          "query": "SELECT * FROM cypher('dsa_knowledge_graph', $$ CREATE (p:Problem {problem_id: $problem_id, leetcode_number: $leetcode_number, title: $title, difficulty_score: $difficulty_score}) $$) as (v agtype)",
          "description": "Create Problem nodes"
        },
        {
          "node_type": "Code",
          "id": "determine_relationships",
          "description": "Analyze relationships between entities",
          "code": "// Logic to determine:\n// - PREREQUISITE_OF: Based on prerequisites array\n// - SUBTOPIC_OF: Based on parent_topic_id\n// - SIMILAR_TO: Based on embedding similarity > 0.85\n// - HAS_TOPIC: Based on problem_topics table\n// - USES_PATTERN: Based on patterns used\nreturn relationships;"
        },
        {
          "node_type": "PostgreSQL",
          "id": "create_edges",
          "operation": "executeQuery",
          "query": "SELECT * FROM cypher('dsa_knowledge_graph', $$ MATCH (a:Topic {topic_id: $from_id}), (b:Topic {topic_id: $to_id}) CREATE (a)-[r:PREREQUISITE_OF {is_hard_requirement: $is_hard, strength: $strength}]->(b) $$) as (v agtype)",
          "description": "Create edges with properties"
        }
      ]
    },
    
    "workflow_6_generate_questions": {
      "name": "Generate Questions from Templates",
      "description": "Create questions for each problem using question type templates",
      "trigger": "On-demand or when new problems added",
      "nodes": [
        {
          "node_type": "Start"
        },
        {
          "node_type": "PostgreSQL",
          "id": "get_problems_without_questions",
          "operation": "executeQuery",
          "query": "SELECT p.*, array_agg(pt.topic_id) as topics FROM problems p LEFT JOIN problem_topics pt ON p.problem_id = pt.problem_id WHERE NOT EXISTS (SELECT 1 FROM questions q WHERE q.problem_id = p.problem_id) GROUP BY p.problem_id"
        },
        {
          "node_type": "Loop",
          "id": "for_each_problem",
          "items": "$json"
        },
        {
          "node_type": "Code",
          "id": "get_topic_question_types",
          "description": "For each topic, get applicable question types",
          "code": "const topics = $json.topics;\nconst questionTypes = [];\ntopics.forEach(topicId => {\n  const topicQuestionTypes = getQuestionTypesForTopic(topicId);\n  questionTypes.push(...topicQuestionTypes);\n});\nreturn questionTypes;\n\nfunction getQuestionTypesForTopic(topicId) {\n  // Map from master reference\n  // e.g., two_pointers → ['implementation', 'pattern_recognition', 'edge_case_identification', 'complexity_analysis', 'code_debugging', 'optimization']\n  return topicQuestionTypeMap[topicId] || ['implementation', 'complexity_analysis'];\n}"
        },
        {
          "node_type": "Loop",
          "id": "for_each_question_type"
        },
        {
          "node_type": "HTTP Request",
          "id": "call_llm_for_question",
          "description": "Use LLM to generate question",
          "method": "POST",
          "url": "http://localhost:11434/api/generate",
          "body": "{\n  \"model\": \"mistral:7b\",\n  \"prompt\": \"Generate a {{ json.question_type }} question for this problem:\\n\\nProblem: {{ json.title }}\\n{{ json.statement }}\\n\\nReturn JSON with: question_text, answer_options (if multiple choice), correct_answer, explanation, hint_level_1, hint_level_2, hint_level_3\",\n  \"format\": \"json\"\n}"
        },
        {
          "node_type": "Code",
          "id": "validate_and_parse_question",
          "description": "Validate LLM output and ensure it meets schema",
          "code": "const generated = JSON.parse($json.response);\nif (!generated.question_text || !generated.correct_answer || !generated.explanation) {\n  throw new Error('Invalid question generated');\n}\nreturn {\n  json: {\n    problem_id: $json.problem_id,\n    topic_id: $json.topics[0],\n    question_type: $json.question_type,\n    ...generated\n  }\n};"
        },
        {
          "node_type": "PostgreSQL",
          "id": "insert_question",
          "operation": "executeQuery",
          "query": "INSERT INTO questions (problem_id, topic_id, question_type, question_text, correct_answer, answer_options, explanation, hint_level_1, hint_level_2, hint_level_3) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING question_id"
        },
        {
          "node_type": "HTTP Request",
          "id": "generate_question_embedding",
          "description": "Generate embedding for question"
        },
        {
          "node_type": "PostgreSQL",
          "id": "save_question_embedding",
          "operation": "executeQuery",
          "query": "INSERT INTO embeddings (entity_type, entity_id, embedding_vector, text_content) VALUES ('question', $1, $2::vector, $3)"
        }
      ]
    },
    
    "workflow_7_sync_graph_relationships": {
      "name": "Sync Similarity Relationships",
      "description": "Use embeddings to create SIMILAR_TO edges",
      "trigger": "Daily or when new problems/content added",
      "nodes": [
        {
          "node_type": "Start"
        },
        {
          "node_type": "PostgreSQL",
          "id": "find_similar_problems",
          "operation": "executeQuery",
          "query": "SELECT e1.entity_id as problem1_id, e2.entity_id as problem2_id, 1 - (e1.embedding_vector <=> e2.embedding_vector) as similarity FROM embeddings e1 CROSS JOIN embeddings e2 WHERE e1.entity_type = 'problem' AND e2.entity_type = 'problem' AND e1.entity_id < e2.entity_id AND 1 - (e1.embedding_vector <=> e2.embedding_vector) > 0.85"
        },
        {
          "node_type": "PostgreSQL",
          "id": "create_similar_edges",
          "operation": "executeQuery",
          "query": "SELECT * FROM cypher('dsa_knowledge_graph', $$ MATCH (p1:Problem {problem_id: $problem1_id}), (p2:Problem {problem_id: $problem2_id}) MERGE (p1)-[r:SIMILAR_TO {similarity_score: $similarity, reason: 'embedding_similarity'}]->(p2) $$) as (v agtype)"
        }
      ]
    }
  },
  
  "workflow_orchestration": {
    "execution_sequence": [
      "1. workflow_1_seed_topics",
      "2. workflow_2_seed_patterns",
      "3. workflow_3_seed_pitfalls",
      "4. workflow_4_generate_embeddings (parallel for different entity types)",
      "5. workflow_5_build_graph",
      "6. workflow_6_generate_questions",
      "7. workflow_7_sync_graph_relationships"
    ],
    "error_handling": {
      "on_failure": "Log error, send notification, continue with next item",
      "retry_strategy": "Exponential backoff, max 3 retries",
      "rollback": "If critical error, rollback last transaction"
    },
    "monitoring": {
      "log_level": "INFO",
      "metrics_to_track": [
        "Items processed per workflow",
        "Errors encountered",
        "Execution time",
        "Database write performance",
        "LLM API call success rate"
      ]
    }
  }
}
```

---

## 8. QUESTION GENERATION TEMPLATES {#question-templates}

### Systematic Question Generation for All Topics

```json
{
  "question_generation_framework": {
    
    "question_type_definitions": {
      
      "implementation": {
        "description": "Write code to solve the problem",
        "template": "Implement the following:\n\n{problem_statement}\n\nConstraints: {constraints}\n\nExample:\n{examples}",
        "answer_type": "code",
        "evaluation": "Run test cases",
        "difficulty_multiplier": 1.0
      },
      
      "complexity_analysis": {
        "description": "Analyze time/space complexity",
        "template": "What is the time complexity of the optimal solution for {problem_title}?\n\nA) O(1)\nB) O(log n)\nC) O(n)\nD) O(n log n)\nE) O(n²)\n\nExplain your reasoning.",
        "answer_type": "multiple_choice_with_explanation",
        "correct_answer_generation": "Extract from solution analysis",
        "difficulty_multiplier": 0.6
      },
      
      "pattern_recognition": {
        "description": "Identify which pattern/technique to use",
        "template": "Which technique is most appropriate for solving {problem_title}?\n\nA) {pattern_1}\nB) {pattern_2}\nC) {pattern_3}\nD) {pattern_4}\n\nExplain why.",
        "answer_type": "multiple_choice_with_explanation",
        "options_generation": "Include correct pattern + 3 plausible distractors from related topics",
        "difficulty_multiplier": 0.7
      },
      
      "edge_case_identification": {
        "description": "Identify edge cases to test",
        "template": "For {problem_title}, which of these is a critical edge case to test?\n\nA) {edge_case_1}\nB) {edge_case_2}\nC) {edge_case_3}\nD) All of the above\n\nExplain what happens in each case.",
        "answer_type": "multiple_choice_with_explanation",
        "options_generation": "Pull from topic's edge_cases list",
        "difficulty_multiplier": 0.8
      },
      
      "code_debugging": {
        "description": "Find and fix bug in given code",
        "template": "This code for {problem_title} has a bug:\n\n```{language}\n{buggy_code}\n```\n\nWhat's wrong?\n\nA) {bug_option_1}\nB) {bug_option_2}\nC) {bug_option_3}\nD) {bug_option_4}\n\nHow would you fix it?",
        "answer_type": "multiple_choice_with_code_fix",
        "buggy_code_generation": "Inject common pitfall from topic's pitfalls list",
        "difficulty_multiplier": 0.9
      },
      
      "optimization": {
        "description": "Improve suboptimal solution",
        "template": "This solution for {problem_title} works but is inefficient:\n\n```{language}\n{suboptimal_code}\n```\n\nComplexity: {current_complexity}\n\nHow can you optimize it to {target_complexity}?\n\nA) {optimization_1}\nB) {optimization_2}\nC) {optimization_3}\nD) {optimization_4}",
        "answer_type": "multiple_choice_with_explanation",
        "difficulty_multiplier": 1.1
      },
      
      "state_definition": {
        "description": "Define DP state correctly",
        "template": "For {problem_title}, you're using DP. What does dp[i][j] represent?\n\nA) {definition_1}\nB) {definition_2}\nC) {definition_3}\nD) {definition_4}",
        "answer_type": "multiple_choice",
        "applicable_topics": ["dynamic_programming"],
        "difficulty_multiplier": 0.7
      },
      
      "base_case_determination": {
        "description": "Determine correct base case",
        "template": "In {problem_title}, what should the base case return?\n\nScenario: {base_case_scenario}\n\nA) {return_value_1}\nB) {return_value_2}\nC) {return_value_3}\nD) {return_value_4}\n\nWhy?",
        "answer_type": "multiple_choice_with_explanation",
        "applicable_topics": ["dynamic_programming", "recursion"],
        "difficulty_multiplier": 0.8
      },
      
      "trade_off_analysis": {
        "description": "Compare different approaches",
        "template": "For {problem_title}, compare these approaches:\n\nApproach A: {approach_1_description}\nTime: {approach_1_time}, Space: {approach_1_space}\n\nApproach B: {approach_2_description}\nTime: {approach_2_time}, Space: {approach_2_space}\n\nWhen would you use each approach?",
        "answer_type": "text",
        "difficulty_multiplier": 1.0
      },
      
      "test_case_generation": {
        "description": "Create test case that breaks code",
        "template": "This code looks correct but has a bug:\n\n```{language}\n{buggy_code}\n```\n\nWhat input would cause it to fail?\n\nA) {test_case_1}\nB) {test_case_2}\nC) {test_case_3}\nD) {test_case_4}\n\nWhat would the output be?",
        "answer_type": "multiple_choice_with_explanation",
        "difficulty_multiplier": 1.1
      },
      
      "observation_testing": {
        "description": "Test key insight before coding",
        "template": "Before solving {problem_title}, identify the key observation:\n\n{observation_question}\n\nTrue or False? Explain.",
        "answer_type": "boolean_with_explanation",
        "observations_source": "Extract from problem's key_insight field",
        "difficulty_multiplier": 0.6
      },
      
      "conversion": {
        "description": "Convert between iterative/recursive",
        "template": "Given this {source_type} solution for {problem_title}:\n\n```{language}\n{source_code}\n```\n\nWhich {target_type} version is equivalent?\n\nA) {option_1}\nB) {option_2}\nC) {option_3}\nD) None are equivalent",
        "answer_type": "multiple_choice",
        "difficulty_multiplier": 1.2
      }
    },
    
    "generation_strategy_per_topic": {
      
      "two_pointers": {
        "recommended_question_types": [
          "implementation",
          "pattern_recognition",
          "edge_case_identification",
          "code_debugging",
          "complexity_analysis"
        ],
        "question_distribution": {
          "implementation": 0.25,
          "pattern_recognition": 0.20,
          "edge_case_identification": 0.20,
          "code_debugging": 0.20,
          "complexity_analysis": 0.15
        },
        "specific_templates": {
          "pointer_direction": {
            "question": "For {problem_title}, should pointers move in opposite or same direction?",
            "options": [
              "Opposite direction (towards center)",
              "Same direction (both forward)",
              "One pointer fixed, other moves",
              "Depends on the constraint"
            ]
          },
          "stopping_condition": {
            "question": "In two pointers for {problem_title}, when should the loop stop?",
            "options": [
              "left < right",
              "left <= right",
              "left < n and right >= 0",
              "When condition is met"
            ]
          }
        }
      },
      
      "dynamic_programming": {
        "recommended_question_types": [
          "state_definition",
          "base_case_determination",
          "complexity_analysis",
          "optimization",
          "conversion"
        ],
        "question_distribution": {
          "state_definition": 0.25,
          "base_case_determination": 0.20,
          "complexity_analysis": 0.15,
          "optimization": 0.20,
          "conversion": 0.20
        },
        "specific_templates": {
          "return_value_semantics": {
            "question": "In {problem_title}, if no solution exists, what should dp return?",
            "options": [
              "0",
              "-1",
              "INT_MAX",
              "INT_MIN"
            ],
            "explanation_required": true
          },
          "space_optimization": {
            "question": "This DP solution uses O(n²) space. Can it be optimized to O(n)?",
            "follow_up": "How? Show the optimized code."
          }
        }
      },
      
      "graph_traversal": {
        "recommended_question_types": [
          "pattern_recognition",
          "observation_testing",
          "code_debugging",
          "trade_off_analysis"
        ],
        "specific_templates": {
          "bfs_vs_dfs": {
            "question": "For {problem_title}, should you use BFS or DFS?",
            "options": [
              "BFS - need shortest path",
              "DFS - need any path",
              "Either works",
              "Neither - need Dijkstra"
            ]
          },
          "visited_marking": {
            "question": "In BFS for {problem_title}, when should you mark nodes as visited?",
            "options": [
              "When adding to queue",
              "When popping from queue",
              "Either works",
              "Don't need visited array"
            ]
          }
        }
      }
    },
    
    "question_difficulty_calibration": {
      "formula": "question_difficulty = base_problem_difficulty × question_type_multiplier ± adjustment",
      "adjustments": {
        "add_10_if": [
          "Requires combining multiple concepts",
          "Non-obvious observation required",
          "Multiple valid approaches exist"
        ],
        "subtract_10_if": [
          "Direct application of pattern",
          "Single concept tested",
          "Clear from problem statement"
        ]
      }
    },
    
    "wrong_answer_generation": {
      "description": "Generate plausible wrong answers for multiple choice",
      "strategies": {
        "common_misconception": "Use pitfalls from topic",
        "off_by_one": "Add/subtract 1 from correct answer",
        "complexity_confusion": "O(n) vs O(n log n) vs O(n²)",
        "similar_pattern": "Pattern from related but different topic",
        "partial_solution": "Handles some cases but not all"
      },
      "example": {
        "question": "Time complexity of Binary Search?",
        "correct": "O(log n)",
        "wrong_answers": [
          "O(n) - common if student thinks linear scan",
          "O(1) - misconception about array access",
          "O(n log n) - confusion with sorting"
        ]
      }
    },
    
    "hint_generation_automatic": {
      "level_1": {
        "template": "Think about {key_concept}. What property does it have?",
        "generation": "Extract key_concept from problem's key_insight"
      },
      "level_2": {
        "template": "Consider using {pattern_name}. How would that apply here?",
        "generation": "Extract pattern_name from problem's primary pattern"
      },
      "level_3": {
        "template": "The key insight is: {key_insight}. Now, how would you implement this?",
        "generation": "Use problem's key_insight field directly"
      }
    }
  },
  
  "question_validation": {
    "required_fields": [
      "question_text",
      "correct_answer",
      "explanation",
      "hint_level_1"
    ],
    "checks": {
      "question_text_length": "50-500 characters",
      "explanation_length": "100-1000 characters",
      "answer_options_count": "3-5 for multiple choice",
      "difficulty_score": "0-100",
      "no_ambiguity": "Validate that correct answer is unambiguous"
    },
    "auto_fixes": {
      "too_short_explanation": "Prompt LLM to expand",
      "missing_hints": "Generate from template",
      "unclear_question": "Rephrase for clarity"
    }
  }
}
```

---

## SUMMARY: Complete Workflow

### From Master Reference → Operational System

```
1. PARSE MASTER REFERENCE
   ├─ Extract topics, patterns, pitfalls, edge cases
   └─ Validate JSON structure

2. SEED POSTGRESQL DATABASE
   ├─ Insert topics with hierarchy
   ├─ Insert patterns with code templates
   ├─ Insert pitfalls with severity
   ├─ Insert edge cases with categories
   └─ Insert problems from LeetCode

3. GENERATE EMBEDDINGS
   ├─ For each problem: title + statement + constraints + examples
   ├─ For each topic: name + description + keywords + patterns
   ├─ For each pattern: name + description + when_to_use + template
   └─ Store in both PostgreSQL (embeddings table) and ChromaDB

4. CONSTRUCT KNOWLEDGE GRAPH (Apache AGE)
   ├─ Create Topic nodes
   ├─ Create Problem nodes
   ├─ Create Pattern nodes
   ├─ Create edges: SUBTOPIC_OF, PREREQUISITE_OF, HAS_TOPIC, USES_PATTERN
   └─ Generate similarity edges using embedding distance

5. GENERATE QUESTIONS
   ├─ For each problem:
   │   ├─ Identify applicable question types from topic
   │   ├─ Use LLM with structured prompts to generate questions
   │   ├─ Validate generated questions
   │   └─ Store in questions table
   └─ Generate embeddings for questions

6. SETUP RAG PIPELINE
   ├─ Configure ChromaDB collections
   ├─ Define retrieval strategies
   ├─ Create prompt templates
   └─ Test end-to-end retrieval

7. DEPLOY N8N WORKFLOWS
   ├─ Schedule daily embedding sync
   ├─ Schedule weekly graph relationship updates
   ├─ Setup question generation triggers
   └─ Monitor and log all operations

8. READY FOR USER INTERACTION
   └─ System now serves:
       ├─ Semantic problem search
       ├─ Personalized question selection
       ├─ AI-powered hints
       ├─ Weakness detection
       └─ Adaptive training plans
```

---

**END OF MASTER REFERENCE DOCUMENT**

This document is the complete source of truth for:
✅ Every DSA topic with subtopics, patterns, pitfalls, and edge cases
✅ Complete database schema for PostgreSQL
✅ Graph relationship definitions for Apache AGE
✅ Embedding generation strategy for vector search
✅ RAG pipeline configuration for AI features
✅ N8N workflow definitions for automation
✅ Question generation templates and strategies

Use this document as input to n8n workflows to systematically populate and maintain the entire system.

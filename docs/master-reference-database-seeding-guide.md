# DSA Topics Master Reference & Database Seeding Guide
## Complete Taxonomy for LeetCode Training Platform

**Purpose:** This document serves as:
1. **Complete topic taxonomy** - All DSA topics and subtopics
2. **Database seeding template** - Structure for automated data generation
3. **Graph construction guide** - How to build topic/problem relationships
4. **Embedding generation spec** - What text to embed and store
5. **RAG pipeline seed data** - Context for LLM retrieval
6. **n8n workflow template** - Structured data for automation

---

## Document Structure

Each topic contains:
- **Topic ID** - Unique identifier
- **Topic Name** - Display name
- **Parent Topic** - For hierarchical structure
- **Difficulty Range** - Min-max difficulty scores (0-100)
- **Core Patterns** - Common solution patterns
- **Common Pitfalls** - Mistakes students make
- **Edge Cases** - Special inputs to test
- **Variations** - Problem modifications
- **Prerequisites** - Topics to learn first
- **Related Topics** - Connected concepts
- **Key Observations** - Insights before solving
- **Implementation Notes** - Language-specific issues
- **Question Types** - What questions to generate
- **Graph Relationships** - How to connect in graph DB
- **Embedding Contexts** - Text for vector embeddings
- **RAG Seed Data** - Context for LLM retrieval

---

# PART 1: CORE DATA STRUCTURES

---

## 1. ARRAYS & STRINGS

### 1.1 Basic Arrays

```json
{
  "topic_id": "arrays_basic",
  "topic_name": "Basic Array Operations",
  "parent_topic": "arrays",
  "difficulty_range": [0, 30],
  "subtopics": [
    "array_traversal",
    "array_insertion_deletion",
    "array_searching",
    "subarray_operations"
  ],
  "core_patterns": [
    "Linear scan",
    "Index manipulation",
    "In-place operations",
    "Prefix/suffix arrays"
  ],
  "common_pitfalls": [
    {
      "pitfall": "Off-by-one errors",
      "description": "Accessing arr[n] when size is n",
      "example": "for (int i = 0; i <= n; i++) - should be i < n",
      "fix": "Always use i < n, not i <= n for 0-indexed arrays"
    },
    {
      "pitfall": "Empty array not handled",
      "description": "Accessing arr[0] without checking if empty",
      "example": "int max = arr[0]; // crashes if arr is empty",
      "fix": "if (arr.empty()) return default_value;"
    },
    {
      "pitfall": "Integer overflow in sum",
      "description": "Sum of array elements exceeds int max",
      "example": "int sum = 0; for (int x : arr) sum += x;",
      "fix": "Use long long sum = 0;"
    },
    {
      "pitfall": "Modifying array during iteration",
      "description": "Changing size while iterating",
      "example": "for (int i = 0; i < arr.size(); i++) arr.erase(...)",
      "fix": "Iterate backwards or use separate result array"
    }
  ],
  "edge_cases": [
    {
      "case": "Empty array",
      "input": "[]",
      "expected_behavior": "Return default value or handle gracefully",
      "test_scenarios": ["sum = 0", "max = undefined or INT_MIN", "length = 0"]
    },
    {
      "case": "Single element",
      "input": "[5]",
      "expected_behavior": "Element itself is answer for most queries",
      "test_scenarios": ["max = 5", "min = 5", "sum = 5"]
    },
    {
      "case": "All same elements",
      "input": "[3, 3, 3, 3]",
      "expected_behavior": "Tests deduplication and counting logic",
      "test_scenarios": ["unique count = 1", "most frequent = 3"]
    },
    {
      "case": "Already sorted (ascending)",
      "input": "[1, 2, 3, 4, 5]",
      "expected_behavior": "Optimal case for many algorithms",
      "test_scenarios": ["binary search works", "no swaps needed"]
    },
    {
      "case": "Reverse sorted (descending)",
      "input": "[5, 4, 3, 2, 1]",
      "expected_behavior": "Worst case for some sorts",
      "test_scenarios": ["max swaps for bubble sort", "inversions = n(n-1)/2"]
    },
    {
      "case": "Contains duplicates",
      "input": "[1, 2, 2, 3, 3, 3]",
      "expected_behavior": "Tests handling of repeated elements",
      "test_scenarios": ["set conversion removes dups", "counting frequency"]
    },
    {
      "case": "All negative numbers",
      "input": "[-5, -2, -8, -1]",
      "expected_behavior": "Maximum is least negative",
      "test_scenarios": ["max = -1", "sum is negative"]
    },
    {
      "case": "Mix of positive and negative",
      "input": "[-3, 5, -1, 8, -4]",
      "expected_behavior": "Tests sign handling",
      "test_scenarios": ["max subarray spans positive section"]
    },
    {
      "case": "Contains zero",
      "input": "[1, 0, -1, 2]",
      "expected_behavior": "Zero as neutral element",
      "test_scenarios": ["product becomes 0", "division by zero risk"]
    },
    {
      "case": "Maximum size constraints",
      "input": "arr.size() = 10^5, elements = [-10^9, 10^9]",
      "expected_behavior": "Tests overflow and performance",
      "test_scenarios": ["sum overflows int", "O(n²) too slow"]
    }
  ],
  "variations": [
    {
      "variation": "2D arrays (matrices)",
      "changes": "Two indices [i][j], row/column operations",
      "new_challenges": "Boundary checking in 2 dimensions"
    },
    {
      "variation": "Circular arrays",
      "changes": "Wrap around using modulo",
      "new_challenges": "Index wrapping: (i + 1) % n"
    },
    {
      "variation": "Jagged arrays",
      "changes": "Rows have different lengths",
      "new_challenges": "Can't assume arr[i].size() == arr[j].size()"
    }
  ],
  "prerequisites": [],
  "related_topics": ["strings", "two_pointers", "sliding_window"],
  "key_observations": [
    "Arrays have O(1) random access",
    "Insertion/deletion at arbitrary position is O(n)",
    "Contiguous memory means cache-friendly",
    "Fixed size in most languages (or amortized resize in dynamic arrays)"
  ],
  "implementation_notes": {
    "cpp": [
      "vector<int> is dynamic array, array<int, N> is fixed size",
      "vector.size() returns size_t (unsigned), be careful with subtraction",
      "Accessing out of bounds causes undefined behavior (not exception)"
    ],
    "python": [
      "Lists are dynamic arrays",
      "Negative indices: arr[-1] is last element",
      "Slicing creates copies: arr[1:3]"
    ],
    "java": [
      "Arrays are fixed size: int[] arr = new int[n]",
      "ArrayList<Integer> for dynamic",
      "Arrays.asList() creates fixed-size list"
    ]
  },
  "question_types": [
    {
      "type": "edge_case_identification",
      "template": "What happens when input is: [empty/single/all same/sorted/reverse]?",
      "example": "For 'find maximum', what if array is empty?"
    },
    {
      "type": "overflow_detection",
      "template": "Given constraints [X], will this code overflow?",
      "example": "n=10^5, elements in [-10^9, 10^9], will int sum work?"
    },
    {
      "type": "complexity_analysis",
      "template": "What's the time/space complexity of this approach?",
      "example": "Nested loops over array - O(n²)"
    },
    {
      "type": "boundary_bugs",
      "template": "This code fails on edge case X. Why?",
      "example": "for (i = 0; i <= n; i++) fails when i = n"
    }
  ],
  "sample_problems": [
    {
      "leetcode_id": 1,
      "name": "Two Sum",
      "difficulty": 15,
      "patterns": ["hash_table", "array_traversal"]
    },
    {
      "leetcode_id": 121,
      "name": "Best Time to Buy and Sell Stock",
      "difficulty": 25,
      "patterns": ["single_pass", "tracking_minimum"]
    }
  ],
  "graph_relationships": {
    "PREREQUISITE_OF": ["two_pointers", "sliding_window", "prefix_sum"],
    "RELATED_TO": ["strings", "hash_tables"],
    "USES_PATTERN": ["linear_scan", "in_place_modification"],
    "HAS_VARIATION": ["2d_arrays", "circular_arrays"]
  },
  "embedding_contexts": [
    {
      "context_type": "concept_explanation",
      "text": "Arrays are contiguous memory blocks storing elements of the same type. They provide O(1) random access by index but O(n) insertion/deletion at arbitrary positions. Key operations include traversal, searching, sorting, and subarray manipulation.",
      "metadata": {"topic": "arrays_basic", "type": "definition"}
    },
    {
      "context_type": "common_patterns",
      "text": "Common array patterns: linear scan for single pass operations, two pointers for pair finding, sliding window for subarray problems, prefix sums for range queries, and in-place modification to achieve O(1) space.",
      "metadata": {"topic": "arrays_basic", "type": "patterns"}
    },
    {
      "context_type": "pitfalls",
      "text": "Common array pitfalls: off-by-one errors when accessing indices, not handling empty arrays, integer overflow in sum operations, modifying array during iteration, and forgetting that array.size() returns unsigned type in C++.",
      "metadata": {"topic": "arrays_basic", "type": "mistakes"}
    }
  ],
  "rag_seed_data": {
    "qa_pairs": [
      {
        "question": "How do I handle an empty array in my solution?",
        "answer": "Always check if the array is empty before accessing elements. Use conditions like 'if (arr.empty()) return default_value;' or 'if (arr.size() == 0)'. For problems like finding maximum, decide what to return: INT_MIN, throw exception, or return optional value."
      },
      {
        "question": "Why does my array sum overflow?",
        "answer": "Integer overflow occurs when sum exceeds int max (≈2×10^9). If you have n=10^5 elements each up to 10^9, the sum can reach 10^14. Solution: use 'long long sum = 0;' instead of 'int sum = 0;'."
      },
      {
        "question": "What's the difference between arr[i] and arr.at(i)?",
        "answer": "arr[i] has no bounds checking (undefined behavior if out of range), while arr.at(i) throws std::out_of_range exception. Use at() for safety during development, [] for performance in production after validation."
      }
    ],
    "troubleshooting": [
      {
        "symptom": "Segmentation fault when accessing array",
        "causes": ["Out of bounds access", "Array is null/uninitialized", "Off-by-one error"],
        "solutions": ["Check i < arr.size(), not i <= arr.size()", "Initialize array before use", "Use arr.at(i) to get exception instead of crash"]
      },
      {
        "symptom": "Wrong answer on edge case with single element",
        "causes": ["Assumed array has at least 2 elements", "Loop starts at index 1 without checking"],
        "solutions": ["Add base case for size 1", "Check arr.size() before assumptions"]
      }
    ]
  }
}
```

### 1.2 Two Pointers

```json
{
  "topic_id": "two_pointers",
  "topic_name": "Two Pointers Technique",
  "parent_topic": "arrays",
  "difficulty_range": [20, 60],
  "subtopics": [
    "opposite_direction_pointers",
    "same_direction_pointers",
    "fast_slow_pointers",
    "three_pointers"
  ],
  "core_patterns": [
    "Opposite direction (converging)",
    "Same direction (slow/fast)",
    "Sliding window with two pointers",
    "Partition with two pointers"
  ],
  "common_pitfalls": [
    {
      "pitfall": "Wrong loop condition",
      "description": "Using left < right when should be left <= right",
      "example": "Two Sum II with left < right misses case when both pointers at same valid element",
      "fix": "For pair sum: left < right. For palindrome check: left < right (don't check middle twice)"
    },
    {
      "pitfall": "Not moving pointers",
      "description": "Infinite loop when neither pointer advances",
      "example": "while (left < right) { if (condition) continue; } - pointers never move",
      "fix": "Ensure at least one pointer moves in every iteration"
    },
    {
      "pitfall": "Moving wrong pointer",
      "description": "In pair sum, moving pointer based on wrong comparison",
      "example": "if (sum < target) right-- (should be left++)",
      "fix": "Sum too small → move left (increase). Sum too large → move right (decrease)"
    },
    {
      "pitfall": "Skipping duplicates incorrectly",
      "description": "In 3Sum, not properly skipping duplicate triplets",
      "example": "while (nums[i] == nums[i+1]) i++ can go out of bounds",
      "fix": "while (i < n-1 && nums[i] == nums[i+1]) i++"
    }
  ],
  "edge_cases": [
    {
      "case": "Array with 2 elements",
      "input": "[1, 2], target = 3",
      "expected_behavior": "Minimum valid input for two pointers",
      "test_scenarios": ["left=0, right=1, sum = 3, return immediately"]
    },
    {
      "case": "All duplicates",
      "input": "[2, 2, 2, 2], target = 4",
      "expected_behavior": "Must handle duplicate pairs correctly",
      "test_scenarios": ["Multiple valid pairs", "Skip duplicate solutions in 3Sum/4Sum"]
    },
    {
      "case": "No solution exists",
      "input": "[1, 2, 3], target = 100",
      "expected_behavior": "Pointers meet without finding answer",
      "test_scenarios": ["Return -1 or empty array", "left crosses right"]
    },
    {
      "case": "Answer at boundaries",
      "input": "[1, 5, 7, 9], target = 10 (answer: 1+9)",
      "expected_behavior": "First and last elements",
      "test_scenarios": ["Check initial positions before any moves"]
    }
  ],
  "variations": [
    {
      "variation": "Three pointers (3Sum)",
      "changes": "Fix one pointer, use two pointers for remaining",
      "new_challenges": "Nested loops, O(n²) instead of O(n)"
    },
    {
      "variation": "Fast & slow pointers (linked list cycle)",
      "changes": "Different speeds, not opposite directions",
      "new_challenges": "Detecting cycle when pointers meet"
    },
    {
      "variation": "Multiple arrays merge",
      "changes": "One pointer per array",
      "new_challenges": "Tracking k pointers simultaneously"
    }
  ],
  "prerequisites": ["arrays_basic", "sorting"],
  "related_topics": ["sliding_window", "binary_search", "hash_tables"],
  "key_observations": [
    "Works on sorted arrays for pair sum problems",
    "Opposite direction: when moving one pointer affects validity",
    "Same direction: when fast pointer scouts ahead",
    "Can reduce O(n²) brute force to O(n) in many cases"
  ],
  "question_types": [
    {
      "type": "pointer_movement_logic",
      "template": "In Two Sum II, if sum < target, which pointer should move and why?",
      "example": "Move left++ because we need larger sum"
    },
    {
      "type": "loop_condition",
      "template": "Should loop be 'left < right' or 'left <= right'?",
      "example": "Depends if we want to process middle element or avoid double-counting"
    },
    {
      "type": "duplicate_handling",
      "template": "How to skip duplicates in 3Sum without going out of bounds?",
      "example": "while (i < n-1 && nums[i] == nums[i+1]) i++"
    }
  ],
  "sample_problems": [
    {
      "leetcode_id": 167,
      "name": "Two Sum II - Input Array Is Sorted",
      "difficulty": 25,
      "patterns": ["two_pointers_opposite"]
    },
    {
      "leetcode_id": 15,
      "name": "3Sum",
      "difficulty": 52,
      "patterns": ["two_pointers_opposite", "duplicate_handling"]
    },
    {
      "leetcode_id": 141,
      "name": "Linked List Cycle",
      "difficulty": 30,
      "patterns": ["fast_slow_pointers"]
    }
  ],
  "graph_relationships": {
    "PREREQUISITE_OF": ["3sum", "4sum", "container_with_water"],
    "RELATED_TO": ["sliding_window", "sorting"],
    "VARIATION_OF": ["arrays_basic"],
    "USES_PATTERN": ["opposite_direction", "same_direction"]
  },
  "embedding_contexts": [
    {
      "context_type": "concept_explanation",
      "text": "Two pointers technique uses two array indices moving in specific patterns to solve problems efficiently. Common patterns: opposite direction (converging from ends), same direction (slow/fast for different speeds), and sliding window (maintaining valid range). Often reduces O(n²) brute force to O(n) time.",
      "metadata": {"topic": "two_pointers", "type": "definition"}
    },
    {
      "context_type": "when_to_use",
      "text": "Use two pointers when: finding pairs/triplets with target sum in sorted array, removing duplicates in-place, detecting cycles in linked lists, partitioning arrays, or maintaining a valid subarray range. Key indicator: problem involves pairs or ranges in linear data structure.",
      "metadata": {"topic": "two_pointers", "type": "application"}
    }
  ],
  "rag_seed_data": {
    "qa_pairs": [
      {
        "question": "When should I use two pointers vs hash table for Two Sum?",
        "answer": "Hash table works on unsorted arrays in O(n) time and space. Two pointers requires sorted array but uses O(1) space. Use hash table if: need to preserve original indices, array unsorted and can't modify. Use two pointers if: array is sorted or can be sorted, need O(1) space, or doing 3Sum/4Sum (hash table becomes complex)."
      },
      {
        "question": "In 3Sum, why do we sort first?",
        "answer": "Sorting enables: (1) Two pointers technique for finding pairs efficiently, (2) Easy duplicate skipping by comparing adjacent elements, (3) Early termination when current element is too large. Without sorting, we'd need O(n²) space for hash table or O(n³) brute force."
      }
    ]
  }
}
```

### 1.3 Sliding Window

```json
{
  "topic_id": "sliding_window",
  "topic_name": "Sliding Window Technique",
  "parent_topic": "arrays",
  "difficulty_range": [25, 70],
  "subtopics": [
    "fixed_size_window",
    "variable_size_window",
    "window_with_conditions",
    "minimum_window"
  ],
  "core_patterns": [
    "Fixed window (size k)",
    "Variable window (expand/contract)",
    "Window with constraint (at most k distinct)",
    "Minimum window satisfying condition"
  ],
  "common_pitfalls": [
    {
      "pitfall": "Not contracting window",
      "description": "Only expanding right, never moving left",
      "example": "Finding max length but never shrinking when invalid",
      "fix": "while (window_invalid) { remove left; left++; }"
    },
    {
      "pitfall": "Off-by-one in window size",
      "description": "Calculating window size as right - left (should be right - left + 1)",
      "example": "Window [1,2,3] with left=1, right=3: size = 3-1 = 2 (wrong, should be 3)",
      "fix": "window_size = right - left + 1"
    },
    {
      "pitfall": "Not updating answer at right time",
      "description": "Updating max length only when window is invalid",
      "example": "Update only in while loop, missing valid windows",
      "fix": "Update max length after contracting, when window is valid"
    },
    {
      "pitfall": "Wrong order of operations",
      "description": "Updating answer before adding element to window",
      "example": "maxLen = max(maxLen, right-left+1); window.add(arr[right]);",
      "fix": "Add element first, then check validity, then update answer"
    }
  ],
  "edge_cases": [
    {
      "case": "Window size > array size",
      "input": "arr = [1,2,3], k = 5",
      "expected_behavior": "Entire array is the window",
      "test_scenarios": ["Return array size or entire array"]
    },
    {
      "case": "k = 1 (window size 1)",
      "input": "arr = [5,3,8,2], k = 1",
      "expected_behavior": "Each element is its own window",
      "test_scenarios": ["Max/min of individual elements"]
    },
    {
      "case": "All elements violate condition",
      "input": "arr = [1,1,1,1], find substring with unique chars",
      "expected_behavior": "Window never expands beyond size 1",
      "test_scenarios": ["Return 1 or 0 depending on problem"]
    },
    {
      "case": "Entire array satisfies condition",
      "input": "arr = [1,2,3,4], all distinct, find longest",
      "expected_behavior": "Never contract window",
      "test_scenarios": ["Return array.length"]
    }
  ],
  "variations": [
    {
      "variation": "Fixed-size window",
      "changes": "Window size k is constant",
      "new_challenges": "Easier - just slide by 1, remove left, add right"
    },
    {
      "variation": "Variable window with constraint",
      "changes": "Window grows/shrinks based on condition",
      "new_challenges": "When to expand vs contract"
    },
    {
      "variation": "Multiple windows",
      "changes": "Track multiple non-overlapping windows",
      "new_challenges": "Coordinating multiple left/right pointers"
    }
  ],
  "prerequisites": ["arrays_basic", "two_pointers", "hash_tables"],
  "related_topics": ["two_pointers", "monotonic_queue", "hash_tables"],
  "key_observations": [
    "Window maintains some property (sum, distinct count, frequency)",
    "Expand when property allows, contract when violated",
    "Often uses hash map/set to track window contents",
    "O(n) time because each element added and removed at most once"
  ],
  "question_types": [
    {
      "type": "window_validity",
      "template": "What condition determines when window is valid/invalid?",
      "example": "For 'at most k distinct', invalid when distinct count > k"
    },
    {
      "type": "contraction_logic",
      "template": "When and how should we contract the window?",
      "example": "while (window_invalid) { remove_left_element; left++; }"
    },
    {
      "type": "answer_update_timing",
      "template": "When should we update the answer - before or after adding element?",
      "example": "After adding, check if valid, then update max if valid"
    }
  ],
  "sample_problems": [
    {
      "leetcode_id": 3,
      "name": "Longest Substring Without Repeating Characters",
      "difficulty": 45,
      "patterns": ["variable_window", "hash_set"]
    },
    {
      "leetcode_id": 76,
      "name": "Minimum Window Substring",
      "difficulty": 68,
      "patterns": ["variable_window", "hash_map", "minimum_window"]
    },
    {
      "leetcode_id": 239,
      "name": "Sliding Window Maximum",
      "difficulty": 65,
      "patterns": ["fixed_window", "monotonic_deque"]
    }
  ],
  "graph_relationships": {
    "PREREQUISITE_OF": ["minimum_window_substring", "longest_substring_variations"],
    "RELATED_TO": ["two_pointers", "hash_tables", "monotonic_deque"],
    "VARIATION_OF": ["two_pointers"]
  }
}
```

---

## 2. DYNAMIC PROGRAMMING

### 2.1 Linear DP (1D)

```json
{
  "topic_id": "dp_linear",
  "topic_name": "1D Dynamic Programming",
  "parent_topic": "dynamic_programming",
  "difficulty_range": [30, 75],
  "subtopics": [
    "fibonacci_like",
    "max_subarray",
    "house_robber",
    "decode_ways",
    "climbing_stairs"
  ],
  "core_patterns": [
    "Take or skip current element",
    "Maximum ending at i",
    "Count ways to reach i",
    "Optimal choice at each step"
  ],
  "common_pitfalls": [
    {
      "pitfall": "Wrong base case",
      "description": "dp[0] initialized incorrectly",
      "example": "House Robber: dp[0] = 0 (wrong, should be nums[0])",
      "fix": "dp[0] represents first house robbed, so dp[0] = nums[0]"
    },
    {
      "pitfall": "Wrong recurrence relation",
      "description": "Incorrect transition between states",
      "example": "dp[i] = dp[i-1] + nums[i] (always taking, should be max)",
      "fix": "dp[i] = max(dp[i-1], dp[i-2] + nums[i]) for House Robber"
    },
    {
      "pitfall": "Off-by-one in indexing",
      "description": "Confusing array index with DP state index",
      "example": "dp[i] represents first i elements, but accessing nums[i]",
      "fix": "If dp[i] = first i elements, use nums[i-1] for ith element"
    },
    {
      "pitfall": "Not handling negative numbers",
      "description": "Initializing to 0 when answer can be negative",
      "example": "Max subarray with all negatives returns 0 (should be max negative)",
      "fix": "Initialize to INT_MIN or first element, not 0"
    },
    {
      "pitfall": "Wrong return value for impossible cases",
      "description": "Returning 0 for 'not possible' in minimization problem",
      "example": "Coin change impossible returns 0 (should return -1 or INT_MAX)",
      "fix": "Use sentinel: -1 for impossible, INT_MAX for minimization init"
    }
  ],
  "edge_cases": [
    {
      "case": "Single element array",
      "input": "[5]",
      "expected_behavior": "Element itself is the answer",
      "test_scenarios": ["Max subarray = 5", "House rob = 5", "LIS = 1"]
    },
    {
      "case": "Two elements",
      "input": "[3, 5]",
      "expected_behavior": "Choice between taking both or one",
      "test_scenarios": ["House rob = max(3, 5) = 5", "Max subarray = 8"]
    },
    {
      "case": "All negative",
      "input": "[-5, -2, -8]",
      "expected_behavior": "Return least negative, not 0",
      "test_scenarios": ["Max subarray = -2", "Cannot skip all elements"]
    },
    {
      "case": "Alternating positive/negative",
      "input": "[5, -2, 3, -1, 4]",
      "expected_behavior": "Tests optimal selection strategy",
      "test_scenarios": ["House rob skips adjacent", "Max subarray may skip negatives"]
    }
  ],
  "variations": [
    {
      "variation": "Circular array (House Robber II)",
      "changes": "First and last are adjacent",
      "new_challenges": "Run DP twice: [0..n-2] and [1..n-1], take max"
    },
    {
      "variation": "With k constraints",
      "changes": "Can skip at most k elements",
      "new_challenges": "Additional dimension: dp[i][skips_used]"
    },
    {
      "variation": "With cooldown",
      "changes": "Must wait k steps between selections",
      "new_challenges": "State includes cooldown timer"
    }
  ],
  "prerequisites": ["arrays_basic", "recursion"],
  "related_topics": ["greedy", "prefix_sums", "kadanes_algorithm"],
  "key_observations": [
    "State definition: What does dp[i] represent exactly?",
    "Base case: Smallest subproblem with known answer",
    "Recurrence: How to build dp[i] from smaller solutions",
    "Answer location: Is it dp[n], max(dp), or something else?"
  ],
  "implementation_notes": {
    "space_optimization": {
      "description": "Often can reduce from O(n) to O(1) space",
      "when": "Current state only depends on fixed number of previous states",
      "example": "Fibonacci: only need dp[i-1] and dp[i-2], so use 2 variables",
      "code_pattern": "int prev2 = base1, prev1 = base2; for(...) { int curr = f(prev1, prev2); prev2 = prev1; prev1 = curr; }"
    }
  },
  "question_types": [
    {
      "type": "base_case_identification",
      "template": "What should dp[0] be for this problem?",
      "example": "Coin change: dp[0] = 0 (0 coins for amount 0)"
    },
    {
      "type": "recurrence_derivation",
      "template": "How do we compute dp[i] from previous states?",
      "example": "House Robber: max(skip current, rob current + best two ago)"
    },
    {
      "type": "impossible_state_handling",
      "template": "What value represents 'not possible'?",
      "example": "Minimization: INT_MAX, Maximization: INT_MIN, Counting: 0, Boolean: false"
    },
    {
      "type": "space_optimization",
      "template": "Can we reduce space from O(n) to O(1)?",
      "example": "If dp[i] only uses dp[i-1] and dp[i-2], use 2 variables"
    }
  ],
  "sample_problems": [
    {
      "leetcode_id": 70,
      "name": "Climbing Stairs",
      "difficulty": 20,
      "patterns": ["fibonacci_like"]
    },
    {
      "leetcode_id": 198,
      "name": "House Robber",
      "difficulty": 35,
      "patterns": ["take_or_skip"]
    },
    {
      "leetcode_id": 53,
      "name": "Maximum Subarray",
      "difficulty": 30,
      "patterns": ["kadanes_algorithm", "max_ending_here"]
    },
    {
      "leetcode_id": 91,
      "name": "Decode Ways",
      "difficulty": 45,
      "patterns": ["count_ways", "string_dp"]
    }
  ],
  "graph_relationships": {
    "PREREQUISITE_OF": ["dp_2d", "dp_knapsack", "dp_lis"],
    "RELATED_TO": ["recursion", "memoization"],
    "USES_PATTERN": ["take_or_skip", "fibonacci_sequence"]
  },
  "embedding_contexts": [
    {
      "context_type": "concept_explanation",
      "text": "Linear DP (1D DP) solves problems where state depends on single dimension. Common pattern: dp[i] represents optimal solution for first i elements. Key steps: define state, determine base case, derive recurrence relation, compute bottom-up or top-down with memoization. Often optimizable from O(n) space to O(1) if only recent states needed.",
      "metadata": {"topic": "dp_linear", "type": "definition"}
    },
    {
      "context_type": "common_mistakes",
      "text": "Common 1D DP mistakes: wrong base case initialization (dp[0] = 0 when should be first element), incorrect recurrence relation (always adding instead of taking max), off-by-one errors (confusing dp index with array index), not handling negative numbers (initializing to 0 instead of INT_MIN), wrong sentinel values for impossible states (using 0 in minimization problems instead of INT_MAX).",
      "metadata": {"topic": "dp_linear", "type": "pitfalls"}
    }
  ],
  "rag_seed_data": {
    "qa_pairs": [
      {
        "question": "How do I know if I can optimize DP space to O(1)?",
        "answer": "Check if dp[i] only depends on a fixed number of previous states (usually dp[i-1], dp[i-2], etc.). If yes, you can use rolling variables instead of array. Example: Fibonacci only needs previous 2 values, so use prev1, prev2 variables. If dp[i] depends on all previous values or variable number of states, you need the full array."
      },
      {
        "question": "What should I return for impossible cases in DP?",
        "answer": "Depends on problem type: (1) Minimization: return -1 or INT_MAX, (2) Maximization: return -1 or INT_MIN, (3) Counting ways: return 0, (4) Boolean (is possible): return false. Initialize DP array with sentinel values: INT_MAX for minimization, INT_MIN for maximization, 0 for counting."
      },
      {
        "question": "My DP gives wrong answer for edge cases. How to debug?",
        "answer": "Check these common issues: (1) Base case: Is dp[0] and dp[1] correct? (2) Empty input: Do you handle n=0? (3) Single element: Does dp[1] work correctly? (4) All negatives: Do you handle when all elements are negative? (5) Index confusion: Are you using dp[i] with nums[i] when you should use nums[i-1]?"
      }
    ],
    "troubleshooting": [
      {
        "symptom": "Wrong answer on single element test case",
        "causes": ["Base case dp[0] incorrect", "Loop starts at i=2 but should start at i=1"],
        "solutions": ["Verify dp[0] = first element for most problems", "Check if loop should start at i=1"]
      },
      {
        "symptom": "Negative answer when all elements are negative",
        "causes": ["Initialized dp to 0", "Taking empty subarray as 0"],
        "solutions": ["Initialize dp[0] to first element, not 0", "Problem may require at least one element"]
      }
    ],
    "examples": [
      {
        "problem": "Climbing Stairs",
        "state_definition": "dp[i] = number of ways to reach step i",
        "base_case": "dp[0] = 1 (one way to stay at ground), dp[1] = 1 (one way to reach step 1)",
        "recurrence": "dp[i] = dp[i-1] + dp[i-2] (come from previous step or two steps before)",
        "answer": "dp[n]",
        "space_optimization": "Only need prev1 and prev2, so O(1) space possible"
      },
      {
        "problem": "House Robber",
        "state_definition": "dp[i] = maximum money robbing houses 0 to i",
        "base_case": "dp[0] = nums[0] (rob first house), dp[1] = max(nums[0], nums[1])",
        "recurrence": "dp[i] = max(dp[i-1], dp[i-2] + nums[i]) (skip current or rob current)",
        "answer": "dp[n-1]",
        "common_error": "Using dp[0] = 0 instead of nums[0]"
      }
    ]
  }
}
```

---

## MASTER DATA STRUCTURE FOR DATABASE SEEDING

### Complete Topic Hierarchy

```json
{
  "topics": [
    {
      "category": "Data Structures",
      "topics": [
        {
          "id": "arrays",
          "name": "Arrays & Strings",
          "subtopics": [
            "arrays_basic",
            "two_pointers",
            "sliding_window",
            "prefix_sums",
            "difference_arrays"
          ]
        },
        {
          "id": "linked_lists",
          "name": "Linked Lists",
          "subtopics": [
            "singly_linked_list",
            "doubly_linked_list",
            "circular_linked_list",
            "fast_slow_pointers",
            "linked_list_reversal"
          ]
        },
        {
          "id": "stacks_queues",
          "name": "Stacks & Queues",
          "subtopics": [
            "stack_basic",
            "queue_basic",
            "monotonic_stack",
            "monotonic_queue",
            "deque",
            "priority_queue"
          ]
        },
        {
          "id": "hash_tables",
          "name": "Hash Tables",
          "subtopics": [
            "hash_map",
            "hash_set",
            "frequency_counting",
            "two_sum_pattern"
          ]
        },
        {
          "id": "trees",
          "name": "Trees",
          "subtopics": [
            "binary_tree",
            "binary_search_tree",
            "tree_traversal",
            "tree_construction",
            "lowest_common_ancestor",
            "tree_dp"
          ]
        },
        {
          "id": "heaps",
          "name": "Heaps",
          "subtopics": [
            "min_heap",
            "max_heap",
            "heap_sort",
            "top_k_elements",
            "merge_k_sorted"
          ]
        },
        {
          "id": "graphs",
          "name": "Graphs",
          "subtopics": [
            "graph_representation",
            "bfs",
            "dfs",
            "topological_sort",
            "shortest_path",
            "minimum_spanning_tree",
            "union_find"
          ]
        },
        {
          "id": "tries",
          "name": "Tries",
          "subtopics": [
            "trie_basic",
            "prefix_search",
            "word_search"
          ]
        }
      ]
    },
    {
      "category": "Algorithms",
      "topics": [
        {
          "id": "sorting",
          "name": "Sorting",
          "subtopics": [
            "quicksort",
            "mergesort",
            "heapsort",
            "counting_sort",
            "bucket_sort"
          ]
        },
        {
          "id": "searching",
          "name": "Searching",
          "subtopics": [
            "binary_search",
            "binary_search_on_answer",
            "ternary_search"
          ]
        },
        {
          "id": "dynamic_programming",
          "name": "Dynamic Programming",
          "subtopics": [
            "dp_linear",
            "dp_2d",
            "dp_knapsack",
            "dp_lis",
            "dp_interval",
            "dp_tree",
            "dp_bitmask",
            "dp_digit"
          ]
        },
        {
          "id": "greedy",
          "name": "Greedy",
          "subtopics": [
            "interval_scheduling",
            "activity_selection",
            "huffman_coding"
          ]
        },
        {
          "id": "backtracking",
          "name": "Backtracking",
          "subtopics": [
            "subsets",
            "permutations",
            "combinations",
            "n_queens",
            "sudoku",
            "word_search"
          ]
        },
        {
          "id": "bit_manipulation",
          "name": "Bit Manipulation",
          "subtopics": [
            "bitwise_operations",
            "bit_tricks",
            "counting_bits"
          ]
        }
      ]
    },
    {
      "category": "Advanced Topics",
      "topics": [
        {
          "id": "advanced_data_structures",
          "name": "Advanced Data Structures",
          "subtopics": [
            "segment_tree",
            "fenwick_tree",
            "sparse_table",
            "disjoint_set_union"
          ]
        },
        {
          "id": "string_algorithms",
          "name": "String Algorithms",
          "subtopics": [
            "kmp",
            "rabin_karp",
            "z_algorithm",
            "suffix_array"
          ]
        },
        {
          "id": "math",
          "name": "Mathematics",
          "subtopics": [
            "number_theory",
            "combinatorics",
            "probability",
            "geometry"
          ]
        }
      ]
    }
  ]
}
```

---

## PART 2: GRAPH CONSTRUCTION SPECIFICATION

### Graph Node Types

```json
{
  "node_types": [
    {
      "type": "Topic",
      "properties": {
        "topic_id": "string (unique)",
        "name": "string",
        "description": "text",
        "difficulty_range": "[min, max]",
        "category": "string"
      },
      "indexes": ["topic_id", "name", "category"]
    },
    {
      "type": "Problem",
      "properties": {
        "problem_id": "integer (unique)",
        "leetcode_number": "integer",
        "title": "string",
        "difficulty_score": "float (0-100)",
        "acceptance_rate": "float",
        "statement": "text"
      },
      "indexes": ["problem_id", "leetcode_number", "difficulty_score"]
    },
    {
      "type": "Pattern",
      "properties": {
        "pattern_id": "string (unique)",
        "name": "string",
        "description": "text",
        "complexity": "integer (1-10)"
      },
      "indexes": ["pattern_id", "name"]
    },
    {
      "type": "Question",
      "properties": {
        "question_id": "integer (unique)",
        "question_type": "string",
        "difficulty_score": "float",
        "concept_tested": "string"
      },
      "indexes": ["question_id", "question_type"]
    },
    {
      "type": "Pitfall",
      "properties": {
        "pitfall_id": "string (unique)",
        "description": "text",
        "severity": "integer (1-5)"
      }
    },
    {
      "type": "EdgeCase",
      "properties": {
        "case_id": "string (unique)",
        "description": "text",
        "input_example": "text"
      }
    }
  ]
}
```

### Graph Edge Types

```json
{
  "edge_types": [
    {
      "type": "PREREQUISITE_OF",
      "from": "Topic",
      "to": "Topic",
      "properties": {
        "is_hard_requirement": "boolean",
        "mastery_threshold": "float (0-1)"
      },
      "description": "Topic A must be learned before Topic B"
    },
    {
      "type": "RELATED_TO",
      "from": "Topic",
      "to": "Topic",
      "properties": {
        "strength": "float (0-1)",
        "shared_concepts": "array of strings"
      },
      "description": "Topics share concepts or techniques"
    },
    {
      "type": "HAS_TOPIC",
      "from": "Problem",
      "to": "Topic",
      "properties": {
        "relevance_score": "float (0-1)",
        "is_primary": "boolean"
      },
      "description": "Problem belongs to topic"
    },
    {
      "type": "USES_PATTERN",
      "from": "Problem",
      "to": "Pattern",
      "properties": {
        "importance": "string (primary|secondary|optional)"
      },
      "description": "Problem uses this solution pattern"
    },
    {
      "type": "SIMILAR_TO",
      "from": "Problem",
      "to": "Problem",
      "properties": {
        "similarity_score": "float (0-1)",
        "reason": "string",
        "shared_patterns": "array of strings"
      },
      "description": "Problems are similar in approach or concept"
    },
    {
      "type": "FOLLOW_UP_OF",
      "from": "Problem",
      "to": "Problem",
      "properties": {
        "difficulty_increase": "float",
        "new_constraints": "array of strings"
      },
      "description": "Problem B is harder version of Problem A"
    },
    {
      "type": "VARIATION_OF",
      "from": "Problem",
      "to": "Problem",
      "properties": {
        "variation_type": "string",
        "description": "text"
      },
      "description": "Problem B is variation of Problem A"
    },
    {
      "type": "HAS_PITFALL",
      "from": "Topic",
      "to": "Pitfall",
      "properties": {
        "frequency": "string (common|occasional|rare)"
      },
      "description": "Topic has common pitfall"
    },
    {
      "type": "HAS_EDGE_CASE",
      "from": "Topic",
      "to": "EdgeCase",
      "properties": {
        "importance": "string (critical|important|nice_to_know)"
      },
      "description": "Topic has important edge case"
    },
    {
      "type": "TESTS_CONCEPT",
      "from": "Question",
      "to": "Topic",
      "properties": {
        "concept_area": "string"
      },
      "description": "Question tests understanding of topic"
    },
    {
      "type": "REVEALS_PITFALL",
      "from": "Question",
      "to": "Pitfall",
      "properties": {},
      "description": "Question designed to catch this pitfall"
    }
  ]
}
```

### Graph Query Patterns

```json
{
  "common_queries": [
    {
      "name": "Find learning path",
      "cypher": "MATCH path = (start:Topic {topic_id: $start_id})<-[:PREREQUISITE_OF*]-(end:Topic {topic_id: $end_id}) RETURN path ORDER BY length(path) LIMIT 1",
      "purpose": "Find prerequisite chain from start to end topic"
    },
    {
      "name": "Find similar unsolved problems",
      "cypher": "MATCH (solved:Problem)<-[:SOLVED]-(user:User {user_id: $user_id}), (solved)-[:SIMILAR_TO {similarity_score: $min_sim}]->(similar:Problem) WHERE NOT EXISTS((user)-[:SOLVED]->(similar)) RETURN similar ORDER BY similar.difficulty_score",
      "purpose": "Recommend similar problems user hasn't solved"
    },
    {
      "name": "Find problems for topic mastery",
      "cypher": "MATCH (t:Topic {topic_id: $topic_id})<-[:HAS_TOPIC]-(p:Problem) WHERE p.difficulty_score >= $min_diff AND p.difficulty_score <= $max_diff RETURN p ORDER BY p.difficulty_score",
      "purpose": "Get problems for practicing specific topic"
    },
    {
      "name": "Find multi-pattern problems",
      "cypher": "MATCH (p:Problem)-[:USES_PATTERN]->(pat:Pattern) WITH p, collect(pat.name) as patterns WHERE size(patterns) >= 2 RETURN p, patterns",
      "purpose": "Find problems that combine multiple patterns"
    },
    {
      "name": "Find follow-up chain",
      "cypher": "MATCH path = (start:Problem {problem_id: $problem_id})-[:FOLLOW_UP_OF*1..3]->(followup:Problem) RETURN path",
      "purpose": "Find progression of harder problems"
    },
    {
      "name": "Topic weak areas",
      "cypher": "MATCH (u:User {user_id: $user_id})-[a:ATTEMPTED]->(q:Question)-[:TESTS_CONCEPT]->(t:Topic) WHERE a.is_correct = false WITH t, count(q) as wrong_count ORDER BY wrong_count DESC LIMIT 5 RETURN t, wrong_count",
      "purpose": "Find topics where user makes most mistakes"
    }
  ]
}
```

---

## PART 3: EMBEDDINGS GENERATION SPECIFICATION

### Embedding Types

```json
{
  "embedding_configurations": [
    {
      "embedding_type": "problem_semantic",
      "purpose": "Enable semantic search for similar problems",
      "model": "text-embedding-3-small",
      "dimensions": 1536,
      "content_structure": {
        "components": [
          "problem.title",
          "problem.statement",
          "problem.constraints (formatted as sentences)",
          "problem.examples[].explanation"
        ],
        "template": "Problem: {title}. {statement} Constraints: {constraints}. Examples: {examples}",
        "preprocessing": [
          "Remove code blocks",
          "Convert bullet points to sentences",
          "Expand abbreviations"
        ]
      },
      "metadata": {
        "problem_id": "integer",
        "difficulty_score": "float",
        "primary_pattern": "string",
        "topics": "array of strings"
      },
      "chunking": {
        "strategy": "single_document",
        "max_tokens": 8000
      },
      "update_triggers": [
        "problem_created",
        "problem_statement_updated",
        "examples_modified"
      ]
    },
    {
      "embedding_type": "topic_concept",
      "purpose": "Match user queries to relevant topics",
      "model": "text-embedding-3-small",
      "dimensions": 1536,
      "content_structure": {
        "components": [
          "topic.name",
          "topic.description",
          "topic.key_observations[]",
          "topic.when_to_use",
          "topic.common_patterns[]"
        ],
        "template": "Topic: {name}. {description} Key concepts: {observations}. Use when: {when_to_use}. Patterns: {patterns}",
        "preprocessing": [
          "Expand code examples to prose",
          "Include pattern descriptions"
        ]
      },
      "metadata": {
        "topic_id": "string",
        "difficulty_range": "array",
        "category": "string"
      }
    },
    {
      "embedding_type": "solution_approach",
      "purpose": "Find similar solution strategies",
      "model": "code-embedding-ada-002",
      "dimensions": 1536,
      "content_structure": {
        "components": [
          "solution.approach_description",
          "solution.algorithm_steps[]",
          "solution.time_complexity",
          "solution.space_complexity",
          "solution.pseudocode"
        ],
        "template": "Approach: {description}. Algorithm: {steps}. Complexity: Time O({time}), Space O({space}). {pseudocode}",
        "preprocessing": [
          "Keep code structure",
          "Include complexity analysis"
        ]
      },
      "metadata": {
        "problem_id": "integer",
        "pattern": "string",
        "language": "string"
      }
    },
    {
      "embedding_type": "question_content",
      "purpose": "Match questions to problems and topics",
      "model": "text-embedding-3-small",
      "dimensions": 1536,
      "content_structure": {
        "components": [
          "question.question_text",
          "question.context",
          "question.concept_tested",
          "question.explanation"
        ],
        "template": "{question_text} Context: {context}. Tests: {concept_tested}. {explanation}",
        "preprocessing": [
          "Include code snippets as text",
          "Preserve technical terms"
        ]
      },
      "metadata": {
        "question_id": "integer",
        "question_type": "string",
        "difficulty": "float",
        "topic_id": "string"
      }
    },
    {
      "embedding_type": "pitfall_description",
      "purpose": "Detect when user might encounter specific pitfall",
      "model": "text-embedding-3-small",
      "dimensions": 1536,
      "content_structure": {
        "components": [
          "pitfall.description",
          "pitfall.example",
          "pitfall.symptoms[]",
          "pitfall.fix"
        ],
        "template": "Common mistake: {description}. Example: {example}. Symptoms: {symptoms}. Fix: {fix}",
        "preprocessing": [
          "Include code examples",
          "Highlight key error patterns"
        ]
      },
      "metadata": {
        "pitfall_id": "string",
        "severity": "integer",
        "topics": "array of strings"
      }
    },
    {
      "embedding_type": "edge_case_scenario",
      "purpose": "Retrieve relevant edge cases for testing",
      "model": "text-embedding-3-small",
      "dimensions": 1536,
      "content_structure": {
        "components": [
          "edge_case.name",
          "edge_case.description",
          "edge_case.input_example",
          "edge_case.expected_behavior",
          "edge_case.why_important"
        ],
        "template": "Edge case: {name}. {description} Example input: {input}. Expected: {expected}. Important because: {why_important}",
        "preprocessing": []
      },
      "metadata": {
        "case_id": "string",
        "topics": "array of strings",
        "importance": "string"
      }
    }
  ]
}
```

### Vector Database Collections

```json
{
  "chromadb_collections": [
    {
      "collection_name": "problems",
      "embedding_type": "problem_semantic",
      "distance_metric": "cosine",
      "metadata_indexes": ["difficulty_score", "primary_pattern", "topics"],
      "search_parameters": {
        "default_k": 10,
        "similarity_threshold": 0.75,
        "metadata_filters": {
          "difficulty_range": "optional",
          "topics": "optional",
          "patterns": "optional"
        }
      },
      "update_strategy": "upsert_on_change",
      "sync_trigger": "problem_table_updated"
    },
    {
      "collection_name": "solutions",
      "embedding_type": "solution_approach",
      "distance_metric": "cosine",
      "metadata_indexes": ["pattern", "language"],
      "search_parameters": {
        "default_k": 5,
        "similarity_threshold": 0.80
      }
    },
    {
      "collection_name": "concepts",
      "embedding_type": "topic_concept",
      "distance_metric": "cosine",
      "metadata_indexes": ["topic_id", "category"],
      "search_parameters": {
        "default_k": 8,
        "similarity_threshold": 0.70
      }
    },
    {
      "collection_name": "questions",
      "embedding_type": "question_content",
      "distance_metric": "cosine",
      "metadata_indexes": ["question_type", "difficulty", "topic_id"],
      "search_parameters": {
        "default_k": 15,
        "similarity_threshold": 0.65
      }
    },
    {
      "collection_name": "pitfalls",
      "embedding_type": "pitfall_description",
      "distance_metric": "cosine",
      "metadata_indexes": ["severity", "topics"],
      "search_parameters": {
        "default_k": 5,
        "similarity_threshold": 0.70
      }
    },
    {
      "collection_name": "edge_cases",
      "embedding_type": "edge_case_scenario",
      "distance_metric": "cosine",
      "metadata_indexes": ["topics", "importance"],
      "search_parameters": {
        "default_k": 10,
        "similarity_threshold": 0.65
      }
    }
  ]
}
```

---

## PART 4: RAG PIPELINE SEED DATA

### RAG Context Templates

```json
{
  "rag_contexts": [
    {
      "context_id": "problem_explanation",
      "purpose": "Help user understand problem statement",
      "retrieval_strategy": {
        "collections": ["problems", "concepts"],
        "query_construction": "Embed user question + current problem title",
        "k_documents": 5,
        "rerank": true
      },
      "context_assembly": {
        "sections": [
          {
            "section": "problem_description",
            "source": "current_problem.statement",
            "max_tokens": 500
          },
          {
            "section": "similar_problems",
            "source": "vector_search_results",
            "format": "For each: title, key insight, pattern used",
            "max_tokens": 300
          },
          {
            "section": "relevant_concepts",
            "source": "concepts_collection",
            "format": "Topic name, description, when to use",
            "max_tokens": 200
          }
        ],
        "total_max_tokens": 1000
      },
      "prompt_template": {
        "system": "You are a patient coding interview tutor. Help the user understand the problem without giving away the solution.",
        "user_template": "User asks: {user_query}\n\nCurrent problem: {problem_description}\n\nSimilar problems:\n{similar_problems}\n\nRelevant concepts:\n{relevant_concepts}\n\nProvide a helpful explanation focusing on understanding, not solving."
      }
    },
    {
      "context_id": "hint_generation",
      "purpose": "Generate progressive hints when user is stuck",
      "retrieval_strategy": {
        "collections": ["solutions", "concepts", "pitfalls"],
        "query_construction": "Problem ID + user's current approach (if any) + attempt history",
        "k_documents": 8
      },
      "context_assembly": {
        "sections": [
          {
            "section": "solution_approaches",
            "source": "solutions_collection",
            "filter": "same_problem_id",
            "max_tokens": 400
          },
          {
            "section": "common_pitfalls",
            "source": "pitfalls_collection",
            "filter": "relevant_topics",
            "max_tokens": 200
          },
          {
            "section": "user_history",
            "source": "user_attempts_summary",
            "format": "Previous attempts, time spent, areas of difficulty",
            "max_tokens": 200
          },
          {
            "section": "edge_cases_missed",
            "source": "edge_cases_collection",
            "filter": "not_considered_yet",
            "max_tokens": 200
          }
        ]
      },
      "prompt_template": {
        "system": "Generate a hint at level {hint_level} (1=subtle, 2=directional, 3=concrete). Never give full solution.",
        "user_template": "User struggling with: {problem_title}\nAttempt #{attempt_number}\nTime spent: {time_spent}min\nUser's approach: {user_approach}\n\nSolution approaches:\n{solution_approaches}\n\nCommon pitfalls:\n{common_pitfalls}\n\nUser's weak areas: {user_weaknesses}\n\nGenerate hint at level {hint_level}."
      }
    },
    {
      "context_id": "weakness_analysis",
      "purpose": "Analyze user's mistakes and identify weak concepts",
      "retrieval_strategy": {
        "collections": ["concepts", "pitfalls", "questions"],
        "query_construction": "Aggregate user's incorrect answers + problems struggled with",
        "k_documents": 15
      },
      "context_assembly": {
        "sections": [
          {
            "section": "incorrect_questions",
            "source": "user_attempts_db",
            "filter": "is_correct = false, last_30_days",
            "format": "Question type, concept tested, user's answer, correct answer",
            "max_tokens": 500
          },
          {
            "section": "related_pitfalls",
            "source": "pitfalls_collection",
            "filter": "topics_from_mistakes",
            "max_tokens": 300
          },
          {
            "section": "concept_explanations",
            "source": "concepts_collection",
            "filter": "weak_topics",
            "max_tokens": 400
          },
          {
            "section": "performance_stats",
            "source": "user_stats",
            "format": "Success rate per topic, average time, improvement trend",
            "max_tokens": 200
          }
        ]
      },
      "prompt_template": {
        "system": "You are an expert tutor analyzing a student's performance. Identify patterns in mistakes, root causes, and provide specific recommendations.",
        "user_template": "Student's recent mistakes:\n{incorrect_questions}\n\nRelated common pitfalls:\n{related_pitfalls}\n\nConcept explanations:\n{concept_explanations}\n\nPerformance stats:\n{performance_stats}\n\nProvide:\n1. Pattern in mistakes (be specific)\n2. Root cause (conceptual gap vs careless)\n3. Specific weakness (e.g., 'integer overflow awareness in binary search', not just 'binary search')\n4. Recommended practice (concrete steps)"
      }
    },
    {
      "context_id": "question_generation",
      "purpose": "Generate new questions for existing problems",
      "retrieval_strategy": {
        "collections": ["problems", "questions", "pitfalls", "edge_cases"],
        "query_construction": "Problem ID + desired question type",
        "k_documents": 10
      },
      "context_assembly": {
        "sections": [
          {
            "section": "problem_details",
            "source": "current_problem",
            "format": "Full problem statement, constraints, examples",
            "max_tokens": 600
          },
          {
            "section": "similar_questions",
            "source": "questions_collection",
            "filter": "same_topic, same_question_type",
            "format": "Example questions for this problem type",
            "max_tokens": 400
          },
          {
            "section": "pitfalls_to_test",
            "source": "pitfalls_collection",
            "filter": "relevant_to_problem",
            "max_tokens": 300
          },
          {
            "section": "edge_cases_to_cover",
            "source": "edge_cases_collection",
            "filter": "relevant_to_problem",
            "max_tokens": 200
          }
        ]
      },
      "prompt_template": {
        "system": "Generate a {question_type} question that tests deep understanding, not memorization.",
        "user_template": "Problem:\n{problem_details}\n\nExample questions of this type:\n{similar_questions}\n\nCommon pitfalls to test:\n{pitfalls_to_test}\n\nEdge cases to consider:\n{edge_cases_to_cover}\n\nGenerate a {question_type} question with:\n- Question text\n- 4 answer options (A, B, C, D)\n- Correct answer\n- Explanation for correct answer\n- Why other answers are wrong\n\nFormat as JSON."
      }
    },
    {
      "context_id": "solution_explanation",
      "purpose": "Explain solution approach after user solves or gives up",
      "retrieval_strategy": {
        "collections": ["solutions", "concepts", "problems"],
        "query_construction": "Problem ID + solution pattern",
        "k_documents": 5
      },
      "context_assembly": {
        "sections": [
          {
            "section": "optimal_solution",
            "source": "solutions_collection",
            "filter": "problem_id, primary_approach",
            "max_tokens": 500
          },
          {
            "section": "alternative_approaches",
            "source": "solutions_collection",
            "filter": "problem_id, secondary_approaches",
            "max_tokens": 300
          },
          {
            "section": "key_insights",
            "source": "concepts_collection",
            "filter": "patterns_used",
            "max_tokens": 200
          },
          {
            "section": "similar_problems",
            "source": "problems_collection",
            "filter": "same_pattern",
            "format": "3-5 similar problems to practice",
            "max_tokens": 200
          }
        ]
      },
      "prompt_template": {
        "system": "Explain the solution clearly, emphasizing key insights and patterns.",
        "user_template": "Problem: {problem_title}\n\nOptimal approach:\n{optimal_solution}\n\nAlternative approaches:\n{alternative_approaches}\n\nKey insights:\n{key_insights}\n\nSimilar problems:\n{similar_problems}\n\nProvide a clear explanation that:\n1. Explains the intuition\n2. Walks through the algorithm\n3. Analyzes time/space complexity\n4. Highlights the key pattern\n5. Suggests similar problems"
      }
    },
    {
      "context_id": "code_debugging_help",
      "purpose": "Help user debug their code",
      "retrieval_strategy": {
        "collections": ["pitfalls", "edge_cases", "solutions"],
        "query_construction": "User's code + error message (if any) + problem ID",
        "k_documents": 8
      },
      "context_assembly": {
        "sections": [
          {
            "section": "user_code",
            "source": "user_input",
            "max_tokens": 400
          },
          {
            "section": "common_bugs",
            "source": "pitfalls_collection",
            "filter": "relevant_to_problem",
            "max_tokens": 300
          },
          {
            "section": "edge_cases_to_test",
            "source": "edge_cases_collection",
            "filter": "likely_failures",
            "max_tokens": 200
          },
          {
            "section": "correct_approach",
            "source": "solutions_collection",
            "filter": "problem_id, primary",
            "format": "High-level approach only, not full code",
            "max_tokens": 200
          }
        ]
      },
      "prompt_template": {
        "system": "Help debug code by asking guiding questions and pointing to issues, not giving full solution.",
        "user_template": "User's code:\n{user_code}\n\nError/Issue: {error_description}\n\nCommon bugs for this problem:\n{common_bugs}\n\nEdge cases to test:\n{edge_cases}\n\nCorrect approach (high-level):\n{correct_approach}\n\nProvide debugging help:\n1. Ask clarifying questions\n2. Point to suspicious code sections\n3. Suggest test cases to try\n4. Hint at the issue without revealing fix"
      }
    }
  ]
}
```

---

## PART 5: N8N WORKFLOW TEMPLATES

### Workflow 1: Database Seeding

```json
{
  "workflow_name": "Seed Topics and Problems",
  "trigger": "Manual / Schedule",
  "nodes": [
    {
      "node_id": "1",
      "type": "Data Source",
      "action": "Read JSON",
      "config": {
        "source": "Master Topic Reference Document",
        "parse_structure": "topics_array"
      },
      "output": "array of topic objects"
    },
    {
      "node_id": "2",
      "type": "PostgreSQL",
      "action": "Insert Topics",
      "config": {
        "table": "topics",
        "operation": "upsert",
        "conflict_key": "topic_id",
        "batch_size": 50
      },
      "input_mapping": {
        "topic_id": "{{$json.topic_id}}",
        "name": "{{$json.topic_name}}",
        "parent_topic_id": "{{$json.parent_topic}}",
        "difficulty_min": "{{$json.difficulty_range[0]}}",
        "difficulty_max": "{{$json.difficulty_range[1]}}",
        "description": "{{$json.description}}",
        "key_observations": "{{JSON.stringify($json.key_observations)}}",
        "created_at": "{{$now}}"
      }
    },
    {
      "node_id": "3",
      "type": "Loop",
      "action": "For Each Subtopic",
      "config": {
        "items": "{{$json.subtopics}}"
      }
    },
    {
      "node_id": "4",
      "type": "PostgreSQL",
      "action": "Insert Subtopics",
      "config": {
        "table": "topics",
        "operation": "upsert"
      }
    },
    {
      "node_id": "5",
      "type": "Apache AGE",
      "action": "Create Topic Nodes",
      "config": {
        "cypher": "MERGE (t:Topic {topic_id: $topic_id}) SET t.name = $name, t.description = $description"
      }
    },
    {
      "node_id": "6",
      "type": "Loop",
      "action": "For Each Prerequisite",
      "config": {
        "items": "{{$json.prerequisites}}"
      }
    },
    {
      "node_id": "7",
      "type": "Apache AGE",
      "action": "Create Prerequisite Edges",
      "config": {
        "cypher": "MATCH (from:Topic {topic_id: $prereq_id}), (to:Topic {topic_id: $topic_id}) MERGE (from)-[:PREREQUISITE_OF {is_hard_requirement: $is_hard}]->(to)"
      }
    },
    {
      "node_id": "8",
      "type": "Generate Embeddings",
      "action": "Create Topic Embeddings",
      "config": {
        "model": "text-embedding-3-small",
        "input_template": "{{$json.embedding_contexts[0].text}}"
      }
    },
    {
      "node_id": "9",
      "type": "ChromaDB",
      "action": "Upsert to concepts collection",
      "config": {
        "collection": "concepts",
        "document": "{{$json.embedding_contexts[0].text}}",
        "embedding": "{{$node8.embedding}}",
        "metadata": {
          "topic_id": "{{$json.topic_id}}",
          "category": "{{$json.category}}",
          "difficulty_range": "{{JSON.stringify($json.difficulty_range)}}"
        },
        "id": "topic_{{$json.topic_id}}"
      }
    },
    {
      "node_id": "10",
      "type": "Loop",
      "action": "For Each Pitfall",
      "config": {
        "items": "{{$json.common_pitfalls}}"
      }
    },
    {
      "node_id": "11",
      "type": "PostgreSQL",
      "action": "Insert Pitfalls",
      "config": {
        "table": "pitfalls",
        "operation": "insert"
      },
      "input_mapping": {
        "pitfall_id": "{{$json.topic_id}}_{{$index}}",
        "topic_id": "{{$json.topic_id}}",
        "description": "{{$item.pitfall}}",
        "example": "{{$item.example}}",
        "fix": "{{$item.fix}}",
        "severity": "{{$item.severity || 3}}"
      }
    },
    {
      "node_id": "12",
      "type": "Generate Embeddings",
      "action": "Create Pitfall Embeddings",
      "config": {
        "model": "text-embedding-3-small",
        "input_template": "Common mistake: {{$item.pitfall}}. Example: {{$item.example}}. Fix: {{$item.fix}}"
      }
    },
    {
      "node_id": "13",
      "type": "ChromaDB",
      "action": "Upsert to pitfalls collection",
      "config": {
        "collection": "pitfalls",
        "id": "{{$json.topic_id}}_pitfall_{{$index}}"
      }
    }
  ],
  "flow": "1 -> 2 -> 3 -> 4 -> 5 -> 6 -> 7 (back to 2 for next topic) -> 8 -> 9 -> 10 -> 11 -> 12 -> 13"
}
```

### Workflow 2: Problem Seeding with Graph Construction

```json
{
  "workflow_name": "Seed Problems and Build Graph",
  "trigger": "Manual / API Call",
  "nodes": [
    {
      "node_id": "1",
      "type": "Input",
      "action": "Receive Problem Data",
      "config": {
        "expected_fields": [
          "leetcode_number",
          "title",
          "statement",
          "constraints",
          "examples",
          "hints",
          "patterns",
          "topics"
        ]
      }
    },
    {
      "node_id": "2",
      "type": "Calculate Difficulty",
      "action": "Run Difficulty Algorithm",
      "config": {
        "function": "calculate_difficulty_score",
        "inputs": {
          "time_complexity": "{{$json.time_complexity}}",
          "patterns": "{{$json.patterns}}",
          "loc": "{{$json.solution_code.length}}",
          "success_rate": "{{$json.historical_success_rate || 0.5}}"
        }
      },
      "output": "difficulty_score (0-100)"
    },
    {
      "node_id": "3",
      "type": "PostgreSQL",
      "action": "Insert Problem",
      "config": {
        "table": "problems",
        "operation": "insert",
        "return": "problem_id"
      },
      "input_mapping": {
        "leetcode_number": "{{$json.leetcode_number}}",
        "title": "{{$json.title}}",
        "statement": "{{$json.statement}}",
        "constraints": "{{JSON.stringify($json.constraints)}}",
        "examples": "{{JSON.stringify($json.examples)}}",
        "hints": "{{JSON.stringify($json.hints)}}",
        "difficulty_score": "{{$node2.difficulty_score}}",
        "primary_pattern": "{{$json.patterns[0]}}",
        "created_at": "{{$now}}"
      }
    },
    {
      "node_id": "4",
      "type": "Apache AGE",
      "action": "Create Problem Node",
      "config": {
        "cypher": "CREATE (p:Problem {problem_id: $problem_id, title: $title, difficulty_score: $difficulty_score, leetcode_number: $leetcode_number})"
      }
    },
    {
      "node_id": "5",
      "type": "Loop",
      "action": "For Each Topic",
      "config": {
        "items": "{{$json.topics}}"
      }
    },
    {
      "node_id": "6",
      "type": "PostgreSQL",
      "action": "Insert Problem-Topic Relation",
      "config": {
        "table": "problem_topics",
        "operation": "insert"
      },
      "input_mapping": {
        "problem_id": "{{$node3.problem_id}}",
        "topic_id": "{{$item.topic_id}}",
        "relevance_score": "{{$item.relevance || 0.8}}",
        "is_primary": "{{$index == 0}}"
      }
    },
    {
      "node_id": "7",
      "type": "Apache AGE",
      "action": "Create HAS_TOPIC Edge",
      "config": {
        "cypher": "MATCH (p:Problem {problem_id: $problem_id}), (t:Topic {topic_id: $topic_id}) MERGE (p)-[:HAS_TOPIC {relevance_score: $relevance, is_primary: $is_primary}]->(t)"
      }
    },
    {
      "node_id": "8",
      "type": "Generate Embeddings",
      "action": "Create Problem Embedding",
      "config": {
        "model": "text-embedding-3-small",
        "input_template": "Problem: {{$json.title}}. {{$json.statement}} Constraints: {{$json.constraints.join('. ')}}. Examples: {{$json.examples[0].explanation}}"
      }
    },
    {
      "node_id": "9",
      "type": "PostgreSQL",
      "action": "Store Embedding",
      "config": {
        "table": "embeddings",
        "operation": "insert"
      },
      "input_mapping": {
        "entity_type": "problem",
        "entity_id": "{{$node3.problem_id}}",
        "embedding_vector": "{{$node8.embedding}}",
        "model_name": "text-embedding-3-small",
        "created_at": "{{$now}}"
      }
    },
    {
      "node_id": "10",
      "type": "ChromaDB",
      "action": "Upsert to problems collection",
      "config": {
        "collection": "problems",
        "document": "{{$node8.input_template}}",
        "embedding": "{{$node8.embedding}}",
        "metadata": {
          "problem_id": "{{$node3.problem_id}}",
          "difficulty_score": "{{$node2.difficulty_score}}",
          "primary_pattern": "{{$json.patterns[0]}}",
          "topics": "{{JSON.stringify($json.topics.map(t => t.topic_id))}}"
        },
        "id": "problem_{{$node3.problem_id}}"
      }
    },
    {
      "node_id": "11",
      "type": "ChromaDB",
      "action": "Find Similar Problems",
      "config": {
        "collection": "problems",
        "query_embedding": "{{$node8.embedding}}",
        "n_results": 10,
        "where": {
          "$and": [
            {"problem_id": {"$ne": "{{$node3.problem_id}}"}},
            {"difficulty_score": {"$gte": "{{$node2.difficulty_score - 15}}"}},
            {"difficulty_score": {"$lte": "{{$node2.difficulty_score + 15}}"}}
          ]
        }
      }
    },
    {
      "node_id": "12",
      "type": "Filter",
      "action": "Filter by Similarity Threshold",
      "config": {
        "condition": "{{$item.distance < 0.3}}"
      }
    },
    {
      "node_id": "13",
      "type": "Apache AGE",
      "action": "Create SIMILAR_TO Edges",
      "config": {
        "cypher": "MATCH (p1:Problem {problem_id: $problem_id1}), (p2:Problem {problem_id: $problem_id2}) MERGE (p1)-[:SIMILAR_TO {similarity_score: $similarity, reason: 'semantic_embedding'}]->(p2)"
      },
      "input_mapping": {
        "problem_id1": "{{$node3.problem_id}}",
        "problem_id2": "{{$item.metadata.problem_id}}",
        "similarity": "{{1 - $item.distance}}"
      }
    }
  ]
}
```

### Workflow 3: Question Generation (LLM-Powered)

```json
{
  "workflow_name": "Generate Questions for Problem",
  "trigger": "API Call / Scheduled",
  "nodes": [
    {
      "node_id": "1",
      "type": "Input",
      "action": "Receive Request",
      "config": {
        "fields": {
          "problem_id": "integer",
          "question_types": "array",
          "count_per_type": "integer"
        }
      }
    },
    {
      "node_id": "2",
      "type": "PostgreSQL",
      "action": "Fetch Problem Details",
      "config": {
        "query": "SELECT * FROM problems WHERE problem_id = {{$json.problem_id}}"
      }
    },
    {
      "node_id": "3",
      "type": "PostgreSQL",
      "action": "Fetch Topic Details",
      "config": {
        "query": "SELECT t.* FROM topics t JOIN problem_topics pt ON t.topic_id = pt.topic_id WHERE pt.problem_id = {{$json.problem_id}}"
      }
    },
    {
      "node_id": "4",
      "type": "ChromaDB",
      "action": "Retrieve Similar Questions",
      "config": {
        "collection": "questions",
        "query_texts": "{{$node2.statement}}",
        "n_results": 5,
        "where": {
          "question_type": "{{$item.question_type}}"
        }
      }
    },
    {
      "node_id": "5",
      "type": "ChromaDB",
      "action": "Retrieve Relevant Pitfalls",
      "config": {
        "collection": "pitfalls",
        "query_texts": "{{$node2.statement}}",
        "n_results": 5,
        "where": {
          "topics": {"$contains": "{{$node3[0].topic_id}}"}
        }
      }
    },
    {
      "node_id": "6",
      "type": "ChromaDB",
      "action": "Retrieve Edge Cases",
      "config": {
        "collection": "edge_cases",
        "query_texts": "{{$node2.statement}}",
        "n_results": 5
      }
    },
    {
      "node_id": "7",
      "type": "Assemble RAG Context",
      "action": "Build Prompt Context",
      "config": {
        "template": {
          "problem_details": {
            "title": "{{$node2.title}}",
            "statement": "{{$node2.statement}}",
            "constraints": "{{$node2.constraints}}",
            "examples": "{{$node2.examples}}"
          },
          "similar_questions": "{{$node4.documents.map(d => d.document).join('\\n\\n')}}",
          "pitfalls_to_test": "{{$node5.documents.map(d => d.document).join('\\n')}}",
          "edge_cases": "{{$node6.documents.map(d => d.document).join('\\n')}}"
        }
      }
    },
    {
      "node_id": "8",
      "type": "Loop",
      "action": "For Each Question Type",
      "config": {
        "items": "{{$json.question_types}}"
      }
    },
    {
      "node_id": "9",
      "type": "Ollama LLM",
      "action": "Generate Question",
      "config": {
        "model": "mistral:7b",
        "temperature": 0.7,
        "system_prompt": "You are an expert at creating insightful coding interview questions. Generate questions that test understanding, not memorization. Output valid JSON only.",
        "user_prompt": "Generate a {{$item}} question for this problem:\n\n{{$node7.json}}\n\nOutput JSON with: question_text, options (A,B,C,D), correct_answer, explanation, wrong_answer_explanations",
        "max_tokens": 1000
      }
    },
    {
      "node_id": "10",
      "type": "Parse JSON",
      "action": "Extract Question Data",
      "config": {
        "input": "{{$node9.response}}"
      }
    },
    {
      "node_id": "11",
      "type": "Validate",
      "action": "Check Question Quality",
      "config": {
        "rules": [
          {"field": "question_text", "min_length": 20},
          {"field": "options", "array_length": 4},
          {"field": "correct_answer", "one_of": ["A", "B", "C", "D"]},
          {"field": "explanation", "min_length": 50}
        ]
      }
    },
    {
      "node_id": "12",
      "type": "PostgreSQL",
      "action": "Insert Question",
      "config": {
        "table": "questions",
        "operation": "insert"
      },
      "input_mapping": {
        "problem_id": "{{$json.problem_id}}",
        "question_type": "{{$item}}",
        "question_text": "{{$node10.question_text}}",
        "correct_answer": "{{JSON.stringify({value: $node10.correct_answer})}}",
        "answer_options": "{{JSON.stringify($node10.options)}}",
        "explanation": "{{$node10.explanation}}",
        "wrong_answer_explanations": "{{JSON.stringify($node10.wrong_answer_explanations)}}",
        "difficulty_score": "{{$node2.difficulty_score}}",
        "created_at": "{{$now}}"
      }
    },
    {
      "node_id": "13",
      "type": "Generate Embeddings",
      "action": "Create Question Embedding",
      "config": {
        "model": "text-embedding-3-small",
        "input_template": "{{$node10.question_text}} Context: Problem about {{$node2.title}}. Tests: {{$item}}. {{$node10.explanation}}"
      }
    },
    {
      "node_id": "14",
      "type": "ChromaDB",
      "action": "Upsert to questions collection",
      "config": {
        "collection": "questions",
        "document": "{{$node13.input_template}}",
        "embedding": "{{$node13.embedding}}",
        "metadata": {
          "question_id": "{{$node12.question_id}}",
          "question_type": "{{$item}}",
          "difficulty": "{{$node2.difficulty_score}}",
          "problem_id": "{{$json.problem_id}}"
        },
        "id": "question_{{$node12.question_id}}"
      }
    }
  ]
}
```

---

## PART 6: DATA GENERATION CHECKLIST

### For Each Topic

```yaml
required_data:
  - topic_id: unique_identifier
  - topic_name: display_name
  - parent_topic: hierarchy
  - difficulty_range: [min, max]
  - subtopics: []
  - core_patterns: []
  - common_pitfalls: [] # with example, fix, severity
  - edge_cases: [] # with input, expected, test_scenarios
  - variations: []
  - prerequisites: []
  - related_topics: []
  - key_observations: []
  - question_types: []
  - sample_problems: []

graph_data:
  - Create Topic node
  - Create PREREQUISITE_OF edges
  - Create RELATED_TO edges
  - Create HAS_PITFALL edges
  - Create HAS_EDGE_CASE edges

embeddings:
  - concept_explanation
  - common_patterns
  - pitfalls
  - when_to_use

rag_contexts:
  - qa_pairs: []
  - troubleshooting: []
  - examples: []
```

### For Each Problem

```yaml
required_data:
  - leetcode_number: integer
  - title: string
  - statement: text
  - constraints: []
  - examples: [{input, output, explanation}]
  - hints: []
  - patterns: []
  - topics: []
  - difficulty_score: calculated

graph_data:
  - Create Problem node
  - Create HAS_TOPIC edges (to topics)
  - Create USES_PATTERN edges (to patterns)
  - Create SIMILAR_TO edges (via embedding similarity)
  - Create FOLLOW_UP_OF edges (if applicable)
  - Create VARIATION_OF edges (if applicable)

embeddings:
  - problem_semantic (title + statement + constraints + examples)
  - solution_approach (for each solution)

questions_to_generate:
  - edge_case_identification: 2-3 questions
  - overflow_detection: 1-2 questions
  - complexity_analysis: 2 questions
  - pattern_recognition: 1-2 questions
  - code_debugging: 1-2 questions (with buggy code)
  - observation_testing: 2-3 questions
  - optimization: 1 question

rag_contexts:
  - problem_explanation
  - solution_walkthrough
  - common_mistakes
  - similar_problems
```

### For Each Question

```yaml
required_data:
  - problem_id: foreign_key
  - question_type: string
  - question_text: text
  - correct_answer: json
  - answer_options: json (if multiple choice)
  - explanation: text
  - wrong_answer_explanations: json
  - difficulty_score: float

graph_data:
  - Create Question node
  - Create TESTS_CONCEPT edge (to topic)
  - Create REVEALS_PITFALL edge (if applicable)

embeddings:
  - question_content (question + context + explanation)

hints:
  - level_1: socratic_question
  - level_2: directional_hint
  - level_3: concrete_hint
```

---

## SUMMARY: COMPLETE WORKFLOW

```
1. DATABASE SEEDING
   ├─ Parse master reference document
   ├─ Insert topics → PostgreSQL
   ├─ Create topic nodes → Apache AGE
   ├─ Create prerequisite edges → Apache AGE
   ├─ Generate topic embeddings → OpenAI API
   ├─ Store embeddings → PostgreSQL + ChromaDB
   ├─ Insert pitfalls → PostgreSQL
   └─ Generate pitfall embeddings → ChromaDB

2. PROBLEM INGESTION
   ├─ Calculate difficulty score
   ├─ Insert problem → PostgreSQL
   ├─ Create problem node → Apache AGE
   ├─ Create topic relationships → Apache AGE
   ├─ Generate problem embedding → OpenAI API
   ├─ Store embedding → PostgreSQL + ChromaDB
   ├─ Find similar problems → ChromaDB query
   └─ Create similarity edges → Apache AGE

3. QUESTION GENERATION
   ├─ For each problem:
   │  ├─ Retrieve RAG context (similar questions, pitfalls, edge cases)
   │  ├─ For each question type:
   │  │  ├─ Generate question via LLM (Ollama)
   │  │  ├─ Validate output
   │  │  ├─ Insert question → PostgreSQL
   │  │  ├─ Create question node → Apache AGE
   │  │  ├─ Generate question embedding → OpenAI API
   │  │  └─ Store in ChromaDB
   │  └─ Generate hints (3 levels)

4. GRAPH ENRICHMENT
   ├─ Find similar problems → ChromaDB similarity search
   ├─ Create SIMILAR_TO edges → Apache AGE
   ├─ Identify follow-up problems → difficulty + pattern analysis
   ├─ Create FOLLOW_UP_OF edges → Apache AGE
   ├─ Find problem variations → pattern matching
   └─ Create VARIATION_OF edges → Apache AGE

5. RAG PIPELINE SETUP
   ├─ Create collection for each embedding type
   ├─ Index metadata fields
   ├─ Set similarity thresholds
   └─ Test retrieval quality

6. CONTINUOUS UPDATES
   ├─ Recalculate difficulty scores (weekly)
   ├─ Update embeddings when content changes
   ├─ Rebuild similarity edges (monthly)
   └─ Add new problems/questions as needed
```

---

This master reference is your complete guide for building the database, graph, embeddings, and RAG system. Use it as the source of truth for n8n automation!

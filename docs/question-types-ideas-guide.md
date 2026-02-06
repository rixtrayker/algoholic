# Question Types Ideas Guide
## Creative Approaches to Testing Deep Understanding

---

## Philosophy: Beyond Surface-Level Testing

**Core Principle:** Don't just test if someone can solve a problem. Test if they **understand** how and why the solution works.

**Bad Question:** "What's the time complexity of this algorithm?"
**Good Question:** "If we change constraint X to Y, how does the time complexity change and why?"

**Bad Question:** "Select the correct code."
**Good Question:** "This code gives TLE on test case X. What's the bottleneck and how would you fix it?"

---

## Table of Contents

1. [Dynamic Programming Questions](#dp-questions)
2. [Graph Algorithm Questions](#graph-questions)
3. [Data Structure Selection Questions](#ds-questions)
4. [Debugging & Code Fixing Questions](#debugging-questions)
5. [Edge Case & Overflow Questions](#edge-case-questions)
6. [Observation Testing Questions](#observation-questions)
7. [Conversion Questions (Iterative ↔ Recursive)](#conversion-questions)
8. [Optimization Questions](#optimization-questions)
9. [Trade-off Analysis Questions](#tradeoff-questions)
10. [Pattern Recognition Questions](#pattern-questions)
11. [AI Assistant Integration Points](#ai-integration)

---

## 1. DYNAMIC PROGRAMMING QUESTIONS {#dp-questions}

### 1.1 Base Case Return Values

**Concept:** What should you return in edge cases?

#### Question Type: "What Should We Return?"

**Example 1: Knapsack with Zero Capacity**
```
Problem: 0/1 Knapsack
Question: If capacity is 0 and we still have items to consider, what should the DP function return?

A) INT_MAX (infinity)
B) INT_MIN (negative infinity)
C) 0
D) -1

Correct: C) 0
Reasoning: Zero capacity means we can't take anything, so profit is 0.
```

**Example 2: Rod Cutting with Length 0**
```
Problem: Rod Cutting
Question: What should dp[0] (rod of length 0) be initialized to?

A) 0
B) -1
C) INT_MAX
D) The price of the smallest piece

Correct: A) 0
Reasoning: A rod of length 0 has no value. This is the base case.
```

**Example 3: Longest Common Subsequence - Empty String**
```
Problem: LCS
Question: If we're comparing empty string with any string, what should we return?

A) Length of the other string
B) 0
C) 1
D) -1

Correct: B) 0
Reasoning: LCS of empty string with anything is 0.
```

**Why This Matters:** Many bugs come from wrong base case returns!

**AI Assistant Integration Point:**
- If user gets this wrong, AI can say: "Think about what this state represents. If we have zero capacity/length/items, what's the only logical answer?"
- Track if user consistently struggles with base cases across problems

---

### 1.2 Maximization vs Minimization Return Values

**Concept:** What should impossible states return?

#### Question Type: "Impossible State Return Value"

**Example 1: Coin Change (Minimization)**
```
Problem: Minimum coins to make amount X
Question: If amount is 5 and no combination of coins can make 5, what should dp[5] be?

A) 0
B) -1
C) INT_MAX or amount + 1
D) The number of coins we tried

Correct: C) INT_MAX or amount + 1
Reasoning: We're minimizing, so impossible states should be "infinity" so they're never chosen as the minimum.

Follow-up: What if this was a MAXIMIZATION problem (max coins)? 
Answer: Then impossible states should return INT_MIN or -1.
```

**Example 2: Partition Equal Subset (Boolean DP)**
```
Problem: Can we partition array into two equal sum subsets?
Question: If we can't make target sum S with given elements, what should dp[S] be?

A) 0 or false
B) 1 or true
C) -1
D) INT_MAX

Correct: A) 0 or false
Reasoning: Boolean DP - false means "not possible"
```

**Pattern Recognition:**
- **Minimization problems:** Impossible = INT_MAX
- **Maximization problems:** Impossible = INT_MIN or -1
- **Boolean problems:** Impossible = false
- **Count problems:** Impossible = 0

**AI Assistant Integration:**
- If user picks wrong sentinel value, AI asks: "Are you minimizing or maximizing? What value ensures this state is never chosen?"

---

### 1.3 DP State Definition Questions

#### Question Type: "What Does This State Represent?"

**Example: Knapsack State**
```
Code snippet:
int dp[n+1][W+1];
dp[i][w] = ...

Question: What does dp[i][w] represent?

A) Max value using first i items with exactly weight w
B) Max value using first i items with weight at most w
C) Max value using item i with capacity w
D) Number of ways to achieve weight w with i items

Correct: B) Max value using first i items with weight at most w

Why wrong answers fail:
A) "Exactly w" is wrong - we can use less than w
C) "Item i" is wrong - it's all items up to i, not just item i
D) This would be for a counting problem
```

**AI Hint Pattern:**
If user struggles: "Look at how the state transitions. If dp[i][w] depends on dp[i-1][w], what does that tell you about what i represents?"

---

### 1.4 DP Optimization Space Questions

#### Question Type: "How Can We Reduce Space?"

**Example: Fibonacci**
```
Original: dp[n] array - O(n) space
Question: What's the minimum space we need?

A) O(n) - can't reduce
B) O(1) - only need 2 variables
C) O(log n) - divide and conquer
D) O(√n)

Correct: B) O(1)
Code:
int prev2 = 0, prev1 = 1;
for (int i = 2; i <= n; i++) {
    int curr = prev1 + prev2;
    prev2 = prev1;
    prev1 = curr;
}

Follow-up: Can all DP problems be optimized to O(1) space?
Answer: No. Only if current state depends on constant number of previous states.
```

**AI Hint:** "Look at the recurrence. How many previous values do you actually need to compute the next one?"

---

## 2. GRAPH ALGORITHM QUESTIONS {#graph-questions}

### 2.1 Key Observations Before Solving

**Concept:** Test if user notices critical properties of the graph

#### Question Type: "Graph Property Recognition"

**Example 1: Detecting Bipartite Graph**
```
Given a graph, before writing any code, what observation helps us know if it's bipartite?

A) Count the number of edges
B) Check if we can 2-color it (no adjacent same color)
C) Check if it has cycles
D) Count connected components

Correct: B) 2-coloring
Key Insight: A graph is bipartite ↔ no odd-length cycles
```

**Example 2: Topological Sort**
```
Before implementing topological sort, what MUST be true about the graph?

A) It must be connected
B) It must be a DAG (no cycles)
C) It must have n-1 edges
D) All nodes must have edges

Correct: B) Must be a DAG
If there's a cycle, topological ordering is impossible.

Follow-up: What's a quick way to check for cycles during topological sort?
Answer: If we can't process all nodes (some remain with in-degree > 0), there's a cycle.
```

**Example 3: Shortest Path Algorithm Selection**
```
Graph has negative edge weights. Before coding, what's the problem with using Dijkstra?

A) It will be slower
B) It won't work correctly
C) It needs more space
D) It's harder to implement

Correct: B) Won't work correctly
Dijkstra assumes once a node is processed, its distance is final. Negative weights violate this.

Follow-up: What should we use instead?
Answer: Bellman-Ford
```

**AI Integration:**
- Present graph properties one at a time
- If user misses key observation, hint: "Before jumping to BFS/DFS, what's special about the structure of this graph?"
- Track which observations user consistently misses

---

### 2.2 BFS vs DFS Decision

#### Question Type: "Which Traversal and Why?"

**Example 1:**
```
Problem: Find if path exists between two nodes
Graph: Unweighted

Question: BFS or DFS?

A) BFS - finds shortest path
B) DFS - uses less space
C) Either - both work for connectivity
D) Neither - need Dijkstra

Correct: C) Either works
For simple connectivity, both are O(V+E). However:
- If you also want SHORTEST path → BFS
- If you want to use less space → DFS (recursion or stack)
```

**Example 2:**
```
Problem: Detect cycle in directed graph

Question: Which approach naturally fits this problem?

A) BFS with coloring
B) DFS with recursion stack tracking
C) Dijkstra
D) Union-Find

Correct: B) DFS with recursion stack
Why: Cycles appear as back edges in DFS tree. If we visit a node that's currently in our recursion stack, we found a cycle.
```

**AI Hint:** "Think about what each traversal gives you. BFS = levels/shortest. DFS = going deep/back edges."

---

### 2.3 Graph Reduction Observations

#### Question Type: "Can We Simplify?"

**Example: Tree is a Graph**
```
Problem: Find diameter of tree (longest path between any two nodes)

Observation: Since it's a TREE (not general graph), what's true?

A) It has cycles
B) Between any two nodes, there's exactly ONE path
C) We need to track visited nodes
D) We need Dijkstra

Correct: B) Exactly one path
Implication: We don't need visited array! Trees have no cycles.

Follow-up: How does this change the algorithm?
Answer: Simple DFS, no visited tracking needed.
```

**AI Integration:**
- If user writes complex solution for tree problem, hint: "You're solving this like a general graph. What's special about trees?"

---

## 3. DATA STRUCTURE SELECTION QUESTIONS {#ds-questions}

### 3.1 Why This Data Structure?

#### Question Type: "Justify the Choice"

**Example 1: Hash Table vs Array**
```
Problem: Two Sum

Question: Why use a hash table instead of sorting + two pointers?

A) Hash table is always faster
B) Hash table preserves original indices
C) Sorting requires extra space
D) Two pointers is harder to implement

Correct: B) Preserves original indices
We need to return the indices, not the values. Sorting loses original positions.

Follow-up: What if we only need to return the values (not indices)?
Answer: Then sorting + two pointers is better (no extra space).
```

**Example 2: Stack for Monotonic Problems**
```
Problem: Next Greater Element

Question: Why does a monotonic stack work here?

A) It's always O(n)
B) It naturally maintains elements in sorted order
C) It processes elements only once and removes useless candidates
D) Stacks are faster than queues

Correct: C) Removes useless candidates
Key insight: If we see a larger element, all smaller elements before it can't be the "next greater" for anything, so we pop them.
```

---

### 3.2 Data Structure Constraints

#### Question Type: "What's the Limitation?"

**Example: Priority Queue**
```
You're using a min-heap (priority queue) for Dijkstra's algorithm.

Question: What operation is NOT efficiently supported?

A) Insert new element - O(log n) ✓
B) Get minimum element - O(1) ✓
C) Delete minimum - O(log n) ✓
D) Decrease priority of arbitrary element - O(n) ✗

Correct: D) Decrease priority
This is why Dijkstra's often adds duplicate entries instead of updating.

Follow-up: How do we handle this?
Answer: Mark nodes as processed and skip if already processed when we pop from heap.
```

---

## 4. DEBUGGING & CODE FIXING QUESTIONS {#debugging-questions}

### 4.1 Find the Bug

#### Question Type: "This Code Fails - Why?"

**Example 1: Off-by-One Error**
```cpp
// Binary search
int binarySearch(vector<int>& arr, int target) {
    int left = 0, right = arr.size();  // BUG!
    
    while (left <= right) {
        int mid = left + (right - left) / 2;
        if (arr[mid] == target) return mid;
        if (arr[mid] < target) left = mid + 1;
        else right = mid - 1;
    }
    return -1;
}

Question: This code crashes on some inputs. What's wrong?

A) left should start at 1
B) right should be arr.size() - 1
C) Should use (left + right) / 2
D) Comparison should be >=

Correct: B) right should be arr.size() - 1
Bug: right = arr.size() is out of bounds. When we access arr[mid] with right = size, we access invalid memory.

Test case that fails: arr = [1], target = 1
- right = 1 (out of bounds)
- mid = 0
- arr[0] == 1, returns 0 (works by luck!)
- But if we had arr = [1, 2], target = 3:
- right = 2, mid = 1, arr[1] = 2 < 3, left = 2
- Next iteration: left = 2, right = 2
- mid = 2, arr[2] = OUT OF BOUNDS!
```

**AI Hint:** "Trace through the code with a specific test case. What happens to 'right' when the array has 3 elements?"

---

**Example 2: Integer Overflow**
```cpp
int mid = (left + right) / 2;  // BUG!

Question: What's wrong with this line?

A) Should use right + left
B) left + right might overflow for large values
C) Should be / 2.0
D) Nothing wrong

Correct: B) Overflow risk
If left = 2^30 and right = 2^30, left + right = 2^31 which overflows signed int!

Fix: mid = left + (right - left) / 2
```

---

**Example 3: Infinite Loop**
```cpp
// Find minimum in rotated sorted array
int findMin(vector<int>& nums) {
    int left = 0, right = nums.size() - 1;
    
    while (left < right) {
        int mid = left + (right - left) / 2;
        
        if (nums[mid] > nums[right]) {
            left = mid;  // BUG!
        } else {
            right = mid;
        }
    }
    return nums[left];
}

Question: This code enters an infinite loop. When and why?

A) When array has duplicates
B) When array is already sorted
C) When left and right are adjacent
D) Never, code is correct

Correct: C) When left and right are adjacent
Example: nums = [3, 1], target = 1
- left = 0, right = 1
- mid = 0
- nums[0] = 3 > nums[1] = 1
- left = mid = 0 (no progress!)
- Infinite loop

Fix: left = mid + 1 (not just mid)
```

**AI Hint:** "Simulate the loop with a 2-element array. Does left or right ever change?"

---

### 4.2 Given Code, Find Failing Test Case

#### Question Type: "What Input Breaks This?"

**Example 1:**
```cpp
bool isPalindrome(string s) {
    for (int i = 0; i < s.length() / 2; i++) {
        if (s[i] != s[s.length() - 1 - i]) {
            return false;
        }
    }
    return true;
}

Question: This code looks correct but fails on some input. What input?

A) Empty string ""
B) Single character "a"
C) String with spaces "a b a"
D) None, code is correct

Correct: A) Empty string
Why: s.length() returns unsigned type (size_t). For empty string, length = 0.
s.length() - 1 = -1, but since it's unsigned, it wraps around to MAX_UINT!
s.length() - 1 - i becomes a huge number, causing out-of-bounds access.

Fix: if (s.empty()) return true;
Or: int n = s.length(); ... n - 1 - i
```

---

**Example 2: TLE (Time Limit Exceeded)**
```cpp
// Count number of ways to climb stairs (1 or 2 steps)
int climbStairs(int n) {
    if (n <= 2) return n;
    return climbStairs(n-1) + climbStairs(n-2);
}

Question: This gives TLE for n = 40. Why?

A) Recursion is too slow
B) Exponential time O(2^n) due to repeated subproblems
C) Stack overflow
D) Need iterative solution

Correct: B) Exponential time
The function recalculates climbStairs(20) millions of times!

Test case: n = 40
Expected: Instant
Actual: Takes minutes

Fix: Add memoization or use DP
```

---

### 4.3 Given TLE Code, Fix It

#### Question Type: "Optimize This"

**Example:**
```cpp
// Find if subarray with sum k exists
bool subarraySum(vector<int>& nums, int k) {
    for (int i = 0; i < nums.size(); i++) {
        int sum = 0;
        for (int j = i; j < nums.size(); j++) {
            sum += nums[j];
            if (sum == k) return true;
        }
    }
    return false;
}

Question: This is O(n²). How can we optimize to O(n)?

A) Use binary search
B) Use prefix sum + hash set
C) Sort the array first
D) Can't optimize further

Correct: B) Prefix sum + hash set

Optimized code:
bool subarraySum(vector<int>& nums, int k) {
    unordered_set<int> seen;
    seen.insert(0);
    int sum = 0;
    
    for (int num : nums) {
        sum += num;
        if (seen.count(sum - k)) return true;
        seen.insert(sum);
    }
    return false;
}

Key insight: If prefix_sum[j] - prefix_sum[i] = k, then subarray [i+1...j] has sum k.
```

---

## 5. EDGE CASE & OVERFLOW QUESTIONS {#edge-case-questions}

### 5.1 Integer Overflow Detection

#### Question Type: "Will This Overflow?"

**Example 1:**
```cpp
int sum = 0;
for (int i = 0; i < n; i++) {
    sum += arr[i];  // Possible overflow?
}

Given: arr contains n = 100,000 elements
Each element: -10^9 ≤ arr[i] ≤ 10^9

Question: Can sum overflow?

A) No, int is enough
B) Yes, use long long
C) Only if all elements are positive
D) Depends on the order

Correct: B) Use long long
Worst case: 100,000 * 10^9 = 10^14
int max ≈ 2 * 10^9
long long max ≈ 9 * 10^18 ✓

Fix: long long sum = 0;
```

---

**Example 2:**
```cpp
int multiply(int a, int b) {
    return a * b;
}

Given: -10^5 ≤ a, b ≤ 10^5

Question: Does this overflow?

A) No, multiplication is safe
B) Yes, need long long
C) Only for negative numbers
D) Only if both are maximum

Correct: B) Need long long
10^5 * 10^5 = 10^10 > 2 * 10^9 (int max)

Fix:
long long multiply(int a, int b) {
    return (long long)a * b;
}
```

---

**Example 3: Factorial Overflow**
```cpp
int factorial(int n) {
    int result = 1;
    for (int i = 2; i <= n; i++) {
        result *= i;
    }
    return result;
}

Question: For what value of n does this first overflow?

A) n = 10
B) n = 13
C) n = 20
D) n = 100

Correct: B) n = 13
13! = 6,227,020,800 > 2^31-1 (int max)
12! = 479,001,600 ✓

Even long long overflows at 21!

Fix for counting problems: Use BigInteger or modulo arithmetic
result = (result * i) % MOD;
```

---

### 5.2 Edge Cases to Always Test

#### Question Type: "What Edge Cases?"

**Example: Array Problems**
```
Which edge cases should you ALWAYS test?

✓ Empty array []
✓ Single element [5]
✓ Two elements [1, 2]
✓ All same elements [3, 3, 3, 3]
✓ Already sorted [1, 2, 3, 4]
✓ Reverse sorted [4, 3, 2, 1]
✓ Array with duplicates [1, 2, 2, 3]
✓ Negative numbers [-5, -2, 0, 3]
✓ Maximum/minimum values [INT_MIN, INT_MAX]
```

**Example: String Problems**
```
✓ Empty string ""
✓ Single character "a"
✓ All same characters "aaaa"
✓ Palindrome "racecar"
✓ Special characters "a!b@c"
✓ Spaces "  a  b  "
✓ Unicode characters (if applicable)
```

**Example: Tree Problems**
```
✓ Empty tree (root = nullptr)
✓ Single node tree
✓ Skewed tree (all left or all right)
✓ Balanced tree
✓ Tree with duplicate values (if allowed)
```

---

### 5.3 Null Pointer / Empty Input Handling

#### Question Type: "What if Input is Empty/Null?"

**Example 1:**
```cpp
int findMax(vector<int>& arr) {
    int maxVal = arr[0];  // BUG: What if arr is empty?
    for (int i = 1; i < arr.size(); i++) {
        maxVal = max(maxVal, arr[i]);
    }
    return maxVal;
}

Question: What should we return if arr is empty?

A) 0
B) INT_MIN
C) INT_MAX
D) Throw error or return special value

Correct: D) Throw error or return special value
There's no "maximum" of an empty array. Options:
- Throw exception
- Return INT_MIN (sentinel value)
- Return std::optional<int>

Fixed code:
int findMax(vector<int>& arr) {
    if (arr.empty()) {
        throw std::invalid_argument("Empty array");
        // or: return INT_MIN;
    }
    int maxVal = arr[0];
    ...
}
```

---

**Example 2: Linked List**
```cpp
ListNode* reverse(ListNode* head) {
    ListNode* prev = nullptr;
    ListNode* curr = head;
    
    while (curr != nullptr) {
        ListNode* next = curr->next;
        curr->next = prev;
        prev = curr;
        curr = next;
    }
    
    return prev;
}

Question: Does this handle empty list correctly?

A) Yes, returns nullptr correctly
B) No, crashes on nullptr
C) No, infinite loop
D) No, wrong result

Correct: A) Yes!
If head = nullptr:
- curr = nullptr
- while loop never executes
- returns prev = nullptr ✓

This is a good example of code that naturally handles edge case.
```

---

## 6. OBSERVATION TESTING QUESTIONS {#observation-questions}

### 6.1 Key Observations (Test Before Solving)

**Concept:** Break problem into observations, test each one separately

#### Example: Two Sum Problem

**Observation 1:**
```
Question: Can we have duplicate elements in the array?

Answer from constraints: Yes (not stated otherwise)

Implication: We need to be careful with indices. If target = 6 and array = [3, 3], we return [0, 1], not [0, 0].
```

**Observation 2:**
```
Question: If we use a hash table, what should we store?

A) element → index
B) element → count
C) index → element
D) element → list of indices

Correct: A) element → index (for single pass)
But if we need to handle duplicates carefully: D) element → list of indices

Key insight: Store complement's index, not current element.
```

**Observation 3:**
```
Question: Can we solve without extra space?

A) Yes, sort + two pointers
B) Yes, but time becomes O(n²)
C) No, impossible
D) Only if array is already sorted

Correct: A) Yes, sort + two pointers
Trade-off: O(n log n) time, O(1) space (but loses original indices)
```

**AI Integration:**
- Present observations as multiple choice before showing solution
- If user misses critical observation, hint: "Think about the constraint: 2 ≤ nums.length ≤ 10^4. What does this tell you?"
- Only reveal solution after user demonstrates understanding of observations

---

#### Example: Longest Substring Without Repeating Characters

**Observation 1:**
```
Question: As we expand the window (move right pointer), what causes us to shrink it (move left pointer)?

A) Window becomes too large
B) We see a duplicate character
C) Sum exceeds threshold
D) We've processed enough characters

Correct: B) Duplicate character
Key insight: We maintain a window of unique characters.
```

**Observation 2:**
```
Question: How do we know if current character is a duplicate?

A) Check entire window every time - O(n²)
B) Use a hash set to track characters in window - O(1) lookup
C) Use a frequency array - O(1) lookup
D) Sort the window - O(n log n)

Correct: B or C) Hash set or frequency array
Both are O(1) lookup, but for ASCII/Unicode, array might be better.
```

**Observation 3:**
```
Question: When we find a duplicate, where should we move the left pointer?

A) left++
B) left = right
C) left = last position of duplicate + 1
D) left = 0 (restart)

Correct: C) last position of duplicate + 1
Example: "abcabcbb", when we see second 'a' at index 3, we move left to index 1 (not just left++).
```

---

### 6.2 Graph Observations

#### Example: Detect Cycle in Undirected Graph

**Observation 1:**
```
Question: In an undirected graph, when does a cycle exist during DFS?

A) When we visit any previously visited node
B) When we visit a node that's not the parent
C) When we see a back edge
D) When degree of any node > 2

Correct: B) Visited node that's not the parent
In undirected graph, parent is automatically "visited" but not a cycle.

Example:
  1 --- 2
        |
        3

DFS from 1: 1 → 2, parent of 2 is 1
From 2: 2 → 3, parent of 3 is 2
From 3: Can go back to 2, but 2 is parent, so no cycle (correct!)

Wrong approach: "Any visited node" would give false positive.
```

**Observation 2:**
```
Question: For cycle detection in undirected graph, do we need a separate "visiting" state (3-color DFS)?

A) Yes, always need 3 states
B) No, 2 states (visited/not visited) + parent tracking is enough
C) Depends on graph size
D) Need 4 states

Correct: B) 2 states + parent
3-color is for DIRECTED graphs to detect back edges in recursion stack.
```

---

### 6.3 DP Observations

#### Example: Longest Increasing Subsequence

**Observation 1:**
```
Question: What does dp[i] represent?

A) Length of LIS in entire array
B) Length of LIS ending at index i
C) Length of LIS starting at index i
D) Maximum element in LIS

Correct: B) LIS ending at index i
Key: We must include arr[i] in this subsequence.
```

**Observation 2:**
```
Question: To calculate dp[i], what do we need to look at?

A) All previous elements
B) Only adjacent element dp[i-1]
C) All previous elements j where arr[j] < arr[i]
D) Next element

Correct: C) All j where arr[j] < arr[i]
We can only extend a subsequence ending at j if arr[j] < arr[i].
```

**Observation 3:**
```
Question: Why can't we use a greedy approach?

A) Greedy is too slow
B) Choosing locally optimal element doesn't guarantee global optimum
C) We need to track multiple states
D) Greedy only works on sorted arrays

Correct: B) Local ≠ Global optimum
Example: [1, 100, 2, 3, 4, 5]
Greedy might pick 1 → 100 (length 2)
Optimal: 1 → 2 → 3 → 4 → 5 (length 5)
```

---

## 7. CONVERSION QUESTIONS {#conversion-questions}

### 7.1 Iterative ↔ Recursive DP

#### Question Type: "Convert This Code"

**Example 1: Fibonacci Recursive → Iterative**

**Given (Recursive with Memoization):**
```cpp
int fib(int n, vector<int>& memo) {
    if (n <= 1) return n;
    if (memo[n] != -1) return memo[n];
    
    memo[n] = fib(n-1, memo) + fib(n-2, memo);
    return memo[n];
}
```

**Question: Convert to iterative (bottom-up). Select correct code:**

**A)**
```cpp
int fib(int n) {
    if (n <= 1) return n;
    vector<int> dp(n+1);
    dp[0] = 0, dp[1] = 1;
    
    for (int i = 2; i <= n; i++) {
        dp[i] = dp[i-1] + dp[i-2];
    }
    return dp[n];
}
```

**B)**
```cpp
int fib(int n) {
    int prev = 0, curr = 1;
    for (int i = 2; i <= n; i++) {
        int next = prev + curr;
        prev = curr;
        curr = next;
    }
    return curr;
}
```

**Correct: A) (if we want to maintain same structure with array)**
**Or B) (if we optimize space to O(1))**

Both are valid! Question tests:
- Understanding of state transitions
- Ability to recognize dependencies
- Space optimization awareness

---

**Example 2: LCS Iterative → Recursive**

**Given (Iterative):**
```cpp
int lcs(string s1, string s2) {
    int m = s1.length(), n = s2.length();
    vector<vector<int>> dp(m+1, vector<int>(n+1, 0));
    
    for (int i = 1; i <= m; i++) {
        for (int j = 1; j <= n; j++) {
            if (s1[i-1] == s2[j-1]) {
                dp[i][j] = dp[i-1][j-1] + 1;
            } else {
                dp[i][j] = max(dp[i-1][j], dp[i][j-1]);
            }
        }
    }
    return dp[m][n];
}
```

**Question: Convert to recursive with memoization. Select correct code:**

```cpp
int lcs(string& s1, string& s2, int i, int j, vector<vector<int>>& memo) {
    // Base case
    if (i == 0 || j == 0) return 0;
    
    // Check memo
    if (memo[i][j] != -1) return memo[i][j];
    
    // Recurrence
    if (s1[i-1] == s2[j-1]) {
        memo[i][j] = 1 + lcs(s1, s2, i-1, j-1, memo);
    } else {
        memo[i][j] = max(
            lcs(s1, s2, i-1, j, memo),
            lcs(s1, s2, i, j-1, memo)
        );
    }
    
    return memo[i][j];
}
```

**Key points to recognize:**
- Base case: i == 0 or j == 0 (empty string)
- Recursive calls match loop transitions
- Memoization prevents recomputation

---

### 7.2 DFS ↔ BFS Conversion

#### Question Type: "Rewrite Using Different Traversal"

**Example: Convert BFS to DFS**

**Given (BFS for tree level order):**
```cpp
vector<vector<int>> levelOrder(TreeNode* root) {
    vector<vector<int>> result;
    if (!root) return result;
    
    queue<TreeNode*> q;
    q.push(root);
    
    while (!q.empty()) {
        int size = q.size();
        vector<int> level;
        
        for (int i = 0; i < size; i++) {
            TreeNode* node = q.front();
            q.pop();
            level.push_back(node->val);
            
            if (node->left) q.push(node->left);
            if (node->right) q.push(node->right);
        }
        result.push_back(level);
    }
    return result;
}
```

**Question: Rewrite using DFS. What changes?**

```cpp
void dfs(TreeNode* node, int depth, vector<vector<int>>& result) {
    if (!node) return;
    
    // Ensure result has enough levels
    if (result.size() == depth) {
        result.push_back({});
    }
    
    result[depth].push_back(node->val);
    
    dfs(node->left, depth + 1, result);
    dfs(node->right, depth + 1, result);
}

vector<vector<int>> levelOrder(TreeNode* root) {
    vector<vector<int>> result;
    dfs(root, 0, result);
    return result;
}
```

**Key difference:** DFS needs to track depth parameter. BFS gets it naturally from queue levels.

---

## 8. OPTIMIZATION QUESTIONS {#optimization-questions}

### 8.1 Time Complexity Improvement

#### Question Type: "Can We Do Better?"

**Example 1: Subarray Sum**
```
Brute Force: O(n³)
for i in range(n):
    for j in range(i, n):
        for k in range(i, j):
            sum += arr[k]

Question: First optimization?

A) Use hash table
B) Calculate sum incrementally (don't recalculate)
C) Sort array
D) Use binary search

Correct: B) Calculate sum incrementally
Improved: O(n²)
for i in range(n):
    sum = 0
    for j in range(i, n):
        sum += arr[j]  # Incremental
        
Question: Second optimization (to O(n))?

Answer: Prefix sum + hash table to find sum in one pass
```

---

**Example 2: Contains Duplicate**
```
Approach 1: Nested loops - O(n²)
for i in range(n):
    for j in range(i+1, n):
        if arr[i] == arr[j]: return true

Question: Optimize to O(n log n)?
Answer: Sort first, check adjacent elements

Question: Optimize to O(n)?
Answer: Use hash set

Question: Can we do O(1) space while keeping O(n) time?
Answer: No, trade-off between time and space

Follow-up: What if array has range 1 to n?
Answer: Use counting sort or marking technique (modify array)
```

---

### 8.2 Space Complexity Improvement

#### Question Type: "Reduce Space Without Losing Time"

**Example: Unique Paths**
```
Original: O(m×n) space
int dp[m][n];

Question: Can we reduce space?

Observation: Each dp[i][j] only depends on dp[i-1][j] and dp[i][j-1]

Answer: Yes, O(n) space (only need previous row)
int dp[n];

Further optimization: O(min(m, n))
Process along shorter dimension

Code:
int uniquePaths(int m, int n) {
    vector<int> dp(n, 1);
    
    for (int i = 1; i < m; i++) {
        for (int j = 1; j < n; j++) {
            dp[j] += dp[j-1];
        }
    }
    return dp[n-1];
}
```

---

## 9. TRADE-OFF ANALYSIS QUESTIONS {#tradeoff-questions}

### 9.1 Time vs Space Trade-offs

#### Question Type: "What's the Trade-off?"

**Example 1: Fibonacci**
```
Approach 1: Naive Recursion
Time: O(2^n)
Space: O(n) recursion stack

Approach 2: Memoization (Top-Down DP)
Time: O(n)
Space: O(n) memo + O(n) recursion

Approach 3: Tabulation (Bottom-Up DP)
Time: O(n)
Space: O(n) DP array

Approach 4: Space-Optimized
Time: O(n)
Space: O(1)

Question: Which approach for different scenarios?

A) N = 10, called once
   Answer: Any approach works, even naive

B) N = 40, called once
   Answer: Must use DP (naive is too slow)

C) N = 1000, called millions of times
   Answer: Space-optimized (O(1) space) or memoization with cache

D) N = 10^6, limited memory
   Answer: Space-optimized (O(1) space) only
```

---

**Example 2: Two Sum**
```
Approach 1: Brute Force
Time: O(n²), Space: O(1)

Approach 2: Hash Table
Time: O(n), Space: O(n)

Approach 3: Sort + Two Pointers
Time: O(n log n), Space: O(1)

Question: When to use each?

A) Array size n = 100, memory is very limited
   Answer: Approach 1 or 3 (space is more important than time)

B) Array size n = 10^6, have memory
   Answer: Approach 2 (hash table is fastest)

C) Need to preserve original array
   Answer: Approach 1 or 2 (approach 3 sorts, modifying array)

D) Array is already sorted
   Answer: Approach 3 (two pointers) is best - O(n) time, O(1) space
```

---

## 10. PATTERN RECOGNITION QUESTIONS {#pattern-questions}

### 10.1 Identify the Pattern

#### Question Type: "What Pattern Does This Problem Use?"

**Example 1:**
```
Problem: Given a string, find the length of the longest substring with at most 2 distinct characters.

Question: What pattern is this?

A) Two Pointers
B) Sliding Window
C) Binary Search
D) Dynamic Programming

Correct: B) Sliding Window
Key words: "longest substring" + "at most k constraint"
Window expands when valid, contracts when invalid.
```

**Example 2:**
```
Problem: Merge k sorted arrays

Question: What data structure pattern?

A) Stack
B) Queue
C) Priority Queue (Heap)
D) Hash Table

Correct: C) Priority Queue
Pattern: When you need to repeatedly get minimum across k sources.
```

**Example 3:**
```
Problem: Given an array, for each element, find the next greater element to its right.

Question: What pattern?

A) Sliding Window
B) Two Pointers
C) Monotonic Stack
D) Binary Search

Correct: C) Monotonic Stack
Key phrase: "next greater/smaller element"
```

---

### 10.2 Pattern Variations

#### Question Type: "Same Pattern, Different Problem"

**Concept:** Show that seemingly different problems use the same pattern

**Example: Sliding Window Variations**

**Problem 1:** Longest substring without repeating characters
**Problem 2:** Minimum window substring
**Problem 3:** Maximum sum subarray of size k
**Problem 4:** Longest substring with at most k distinct characters

**Question: What do all these have in common?**

Answer: Sliding window pattern
- Expand window (add right element)
- Check validity
- Contract window if needed (remove left)
- Track optimal result

**Variation point:** What makes window "valid" changes per problem!

---

## 11. AI ASSISTANT INTEGRATION POINTS {#ai-integration}

### 11.1 When to Trigger AI Hints

**Trigger Conditions:**

1. **User answers incorrectly 2+ times on same concept**
   - AI analyzes error pattern
   - Provides targeted hint, not full solution

2. **User takes too long (>2x expected time)**
   - AI: "You're on the right track, but consider..."
   - Nudge toward observation they're missing

3. **User skips too quickly**
   - AI: "Before submitting, have you considered edge case X?"

4. **User shows consistent weakness**
   - AI: "I notice you struggle with base cases. Let's practice..."
   - Redirect to foundation questions

---

### 11.2 AI Hint Levels

**Level 1: Observation Hint (Don't Give Away)**
```
User struggles with: "Find cycle in linked list"

AI Level 1: "What happens when two pointers move at different speeds in a cycle?"

NOT: "Use slow and fast pointers"
```

**Level 2: Direction Hint**
```
AI Level 2: "The slow pointer moves 1 step, fast moves 2 steps. What happens eventually?"

NOT: "They will meet inside the cycle"
```

**Level 3: Concrete Hint**
```
AI Level 3: "If there's a cycle, the fast pointer will eventually lap the slow pointer and they'll meet. Why?"

NOT: Full solution
```

**Level 4: Guided Solution (Last Resort)**
```
AI Level 4: "Let's break this down:
1. What if there's no cycle? Where does fast pointer go?
2. What if there is a cycle? Can fast pointer ever not catch slow?"
```

---

### 11.3 AI Context Data

**What AI should receive:**

```json
{
  "question_id": 123,
  "user_id": 456,
  "attempt_number": 3,
  "time_spent_seconds": 180,
  "user_answer": "B",
  "correct_answer": "C",
  "question_context": {
    "type": "complexity_analysis",
    "difficulty": 35.0,
    "concept": "dynamic_programming"
  },
  "user_history": {
    "similar_questions_attempted": 15,
    "similar_questions_correct": 8,
    "avg_time_vs_expected": 1.4,
    "weak_concepts": ["base_cases", "state_definition"]
  },
  "session_context": {
    "questions_in_session": 5,
    "accuracy_this_session": 0.6,
    "energy_level": "medium"
  }
}
```

**AI Response Format:**
```json
{
  "hint_level": 2,
  "hint_text": "Think about what value you'd return if the capacity was 0. What does 0 capacity mean?",
  "follow_up_question": "If capacity is 0, can we take any items?",
  "weakness_detected": "base_case_return_values",
  "recommended_next": [
    "question_id": 234,
    "reason": "More practice with base cases"
  ]
}
```

---

### 11.4 Offline Hints (No AI)

For users who disable AI or in offline mode:

**Progressive Hint System:**

**Hint 1 (Conceptual):**
"Consider what happens in the base case."

**Hint 2 (Directional):**
"When capacity is 0, no items can be taken."

**Hint 3 (Concrete):**
"Zero capacity means zero value. What number represents zero value?"

**Hint 4 (Near-Solution):**
"The base case should return 0 because we can't add any value with 0 capacity."

**Reveal Answer Button:**
Only after viewing all hints.

---

## QUESTION GENERATION GUIDELINES

### For Each New Problem, Create:

**Minimum Set:**
1. 1-2 Observation questions (test understanding before solving)
2. 1 Edge case question
3. 1 Complexity analysis question
4. 1 Pattern recognition question

**Bonus Questions:**
5. 1 Code debugging question (find bug in solution)
6. 1 Optimization question (improve given solution)
7. 1 Trade-off question (compare approaches)
8. 1 Conversion question (iterative ↔ recursive)

---

### Question Difficulty Progression

**Easy Questions:**
- Multiple choice with obvious wrong answers
- Single concept tested
- Direct questions

**Medium Questions:**
- Multiple concepts combined
- Subtle wrong answers
- Requires deeper thinking

**Hard Questions:**
- Non-obvious connections
- Multiple valid approaches
- Requires synthesis

---

## SUMMARY: QUESTION TYPE CHECKLIST

When creating questions for a problem, include:

- [ ] **Base case questions** - What to return in edge cases
- [ ] **Sentinel value questions** - MAX/MIN/0/-1 for impossible states
- [ ] **State definition questions** - What does dp[i][j] mean?
- [ ] **Overflow questions** - Will int be enough?
- [ ] **Edge case identification** - Empty input, single element, etc.
- [ ] **Bug finding questions** - What's wrong with this code?
- [ ] **Test case generation** - What input breaks this?
- [ ] **Optimization questions** - Can we do better?
- [ ] **Space optimization** - How to reduce space?
- [ ] **Trade-off analysis** - Time vs space, when to use what?
- [ ] **Pattern recognition** - What pattern is this?
- [ ] **Observation testing** - Key insights before solving
- [ ] **Conversion questions** - Iterative ↔ recursive, BFS ↔ DFS
- [ ] **Data structure justification** - Why this DS?
- [ ] **Constraint impact** - How do constraints affect solution?

---

**Remember:** The goal is to test **understanding**, not memorization. If a user can answer these questions correctly, they truly understand the problem and solution, not just memorized the code.

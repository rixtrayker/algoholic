# Question Types & Pedagogical Strategies Guide
## Deep Understanding Through Strategic Question Design

---

## Philosophy: Beyond Memorization

**Core Principle:** Questions should test *understanding* and *reasoning*, not pattern memorization.

**Bad Question:**
> "What is the time complexity of the optimal solution for Two Sum?"

**Good Question:**
> "You solve Two Sum with a hash map and get O(n) time. Your friend says they can do O(1) by sorting first. What's wrong with their reasoning?"

---

## PART 1: DP QUESTION TYPES

### 1.1 Return Value Understanding

**Why It Matters:** Students often memorize "return dp[n]" without understanding what the return value represents.

#### Question Type: Return Value Semantics

**Template:**
```
You're solving [problem] with DP.
In the following scenarios, what should your function return?

A) [Scenario 1] → ?
B) [Scenario 2] → ?
C) [Scenario 3] → ?

Options: -1, 0, 1, INT_MAX, INT_MIN, empty array, null
```

**Example 1: Coin Change**
```
Problem: Find minimum coins to make amount = 11
Available coins: [1, 5, 10]

What should your function return in these cases?

1. It's possible to make 11 (answer exists)
   Answer: 2 (coins: 10 + 1)

2. It's impossible to make 11 (coins: [5, 10] only)
   Answer: -1 (NOT 0, NOT INT_MAX)

3. Amount = 0 (base case)
   Answer: 0 (NOT -1)

Now explain: Why -1 for impossible, not 0?
Why 0 for amount=0, not -1?
```

**Example 2: Longest Increasing Subsequence**
```
What should LIS return for:
1. Empty array → ?
2. All elements decreasing [5,4,3,2,1] → ?
3. All elements equal [5,5,5,5,5] → ?

Common mistakes:
- Returning 0 for case 2 (should be 1 - single element is valid)
- Returning array length for case 3 (should be 1)
```

**Example 3: Maximum Subarray Sum**
```
What should Kadane's algorithm return for:
1. All negative numbers [-5, -2, -8, -1]
   A) 0 (empty subarray allowed)
   B) -1 (least negative)
   C) Depends on problem statement

2. Empty array []
   A) 0
   B) INT_MIN
   C) Throw error

Correct: Case 1 = B (must return best available, even if negative)
         Case 2 = C (undefined behavior - check constraints)
```

#### Question Type: Initialization Values

**Template:**
```
You're filling a DP table for [problem].

Why do we initialize:
dp[0] = X (not Y)?
dp[i] = Z for i > 0 (not W)?

What goes wrong if we use the wrong initialization?
```

**Example: Coin Change**
```
Standard initialization:
dp[0] = 0
dp[1..n] = INT_MAX

Why not:
dp[0] = INT_MAX?  → Base case breaks
dp[1..n] = 0?     → We'd never update (0 is always min)
dp[1..n] = -1?    → Can't do min(-1, anything + 1)

Multiple choice:
If we initialize dp[1..n] = amount + 1 instead of INT_MAX,
what changes?

A) Nothing, both work
B) Breaks - can't detect impossible cases
C) Works but less clear
D) Causes overflow

Answer: C - Works because (amount + 1) is effectively infinity
(can't use more than 'amount' coins of value 1)
```

---

### 1.2 Code Conversion Challenges

#### Question Type: Iterative ↔ Recursive Conversion

**Template:**
```
Given this iterative DP solution:
[CODE]

Which of these recursive versions is equivalent?
A) [Recursive version 1]
B) [Recursive version 2]
C) [Recursive version 3]
D) None are equivalent

Explain why the others are wrong.
```

**Example: Fibonacci**
```
Iterative version:
```cpp
int fib(int n) {
    if (n <= 1) return n;
    int prev2 = 0, prev1 = 1;
    for (int i = 2; i <= n; i++) {
        int curr = prev1 + prev2;
        prev2 = prev1;
        prev1 = curr;
    }
    return prev1;
}
```

Which recursive version is equivalent?

**Option A:**
```cpp
int fib(int n) {
    if (n <= 1) return n;
    return fib(n-1) + fib(n-2);
}
```

**Option B:**
```cpp
int fib(int n, vector<int>& memo) {
    if (n <= 1) return n;
    if (memo[n] != -1) return memo[n];
    return memo[n] = fib(n-1, memo) + fib(n-2, memo);
}
```

**Option C:**
```cpp
int fib(int n, int prev2 = 0, int prev1 = 1, int i = 2) {
    if (i > n) return prev1;
    return fib(n, prev1, prev1 + prev2, i + 1);
}
```

Answer: B (with memoization) is most equivalent in complexity
        A is correct but inefficient
        C is tail-recursive version (same complexity as iterative)

Follow-up: Why does A have exponential time while iterative is O(n)?
```

**Example: 0/1 Knapsack**
```
Given iterative 2D DP:
```cpp
int knapsack(vector<int>& weights, vector<int>& values, int W) {
    int n = weights.size();
    vector<vector<int>> dp(n+1, vector<int>(W+1, 0));
    
    for (int i = 1; i <= n; i++) {
        for (int w = 0; w <= W; w++) {
            if (weights[i-1] <= w) {
                dp[i][w] = max(dp[i-1][w], 
                               dp[i-1][w-weights[i-1]] + values[i-1]);
            } else {
                dp[i][w] = dp[i-1][w];
            }
        }
    }
    return dp[n][W];
}
```

Convert to recursive with memoization. Which is correct?

[Show 3 options with subtle bugs in 2 of them]

Common bugs to include:
- Wrong base case
- Memoization key calculation error
- Off-by-one in array indexing
```

---

### 1.3 Space Optimization

#### Question Type: 2D → 1D DP Conversion

**Template:**
```
This 2D DP can be optimized to 1D.

Original: dp[i][j] depends on dp[i-1][...] only

Which conversion is correct?
[Show original]
[Show 3 attempted 1D versions]

Key question: Do we need to iterate forwards or backwards? Why?
```

**Example: 0/1 Knapsack Space Optimization**
```
Original 2D:
dp[i][w] = max(dp[i-1][w], dp[i-1][w-weight[i]] + value[i])

Option A - Iterate forward:
```cpp
for (int i = 0; i < n; i++) {
    for (int w = 0; w <= W; w++) {
        if (weights[i] <= w)
            dp[w] = max(dp[w], dp[w-weights[i]] + values[i]);
    }
}
```

Option B - Iterate backward:
```cpp
for (int i = 0; i < n; i++) {
    for (int w = W; w >= weights[i]; w--) {
        dp[w] = max(dp[w], dp[w-weights[i]] + values[i]);
    }
}
```

Question: Which is correct and why?

Answer: B (backward)
Why? Forward iteration would use UPDATED values from current row
      (turns 0/1 knapsack into unbounded knapsack)
      Backward ensures we use values from previous row

Follow-up: When would forward iteration be correct?
Answer: Unbounded knapsack!
```

---

## PART 2: INTEGER OVERFLOW QUESTIONS

### 2.1 Overflow Detection

#### Question Type: Will This Overflow?

**Template:**
```
Given constraints: [...]
This code will:

```cpp
int result = a * b + c;
```

A) Never overflow
B) Overflow when [condition]
C) Always overflow
D) Impossible to determine

If it can overflow, how do you fix it?
```

**Example 1: Two Sum Multiplication**
```
Constraints: -10^9 ≤ nums[i] ≤ 10^9

Code:
```cpp
int product = nums[i] * nums[j];
if (product == target) return {i, j};
```

Will this overflow?
A) No, int is enough
B) Yes, need long long
C) Yes, need BigInteger
D) Depends on target value

Answer: B
Explanation: 10^9 × 10^9 = 10^18
             int max ≈ 2×10^9
             long long max ≈ 9×10^18 ✓

Fixed:
```cpp
long long product = (long long)nums[i] * nums[j];
```

Follow-up: What if constraints were -10^10 ≤ nums[i] ≤ 10^10?
Answer: Even long long would overflow! Need to avoid multiplication or use BigInteger.
```

**Example 2: Sum Overflow**
```
Calculating sum of array:
Constraints: n ≤ 10^5, -10^9 ≤ nums[i] ≤ 10^9

```cpp
int sum = 0;
for (int num : nums) {
    sum += num;
}
```

Will this overflow?
A) No
B) Yes - worst case all elements are 10^9
C) Yes - worst case mix of positive/negative
D) Depends on actual values

Answer: B
Worst case: 10^5 × 10^9 = 10^14 > int max (2×10^9)
Need: long long sum = 0;
```

#### Question Type: Overflow Prevention Techniques

**Template:**
```
To compute (a * b) mod m without overflow:

Which technique is correct?
A) int result = (a * b) % m;
B) int result = ((long long)a * b) % m;
C) int result = (a % m) * (b % m) % m;
D) int result = (a % m * b % m) % m;
```

**Example: Modular Arithmetic**
```
Computing n! mod (10^9 + 7):

```cpp
const int MOD = 1e9 + 7;
int factorial(int n) {
    int result = 1;
    for (int i = 2; i <= n; i++) {
        result = (result * i) % MOD;  // Bug?
    }
    return result;
}
```

Is this correct?
A) Yes, perfect
B) Overflows before mod
C) Works for small n, fails for large
D) Logic error in loop

Answer: B (Overflows for n ≥ 13)
Because result * i happens BEFORE mod

Fixed:
```cpp
result = (1LL * result * i) % MOD;
// or
result = ((long long)result * i) % MOD;
```

Multiple choice follow-up:
When n = 100, result = 0. Why?
A) Overflow corrupted the value
B) 100! is divisible by MOD
C) Logic error
D) Both A and B

Answer: B (MOD is prime, but 100! contains it as factor)
```

---

## PART 3: CODE DEBUGGING CHALLENGES

### 3.1 Find the Bug

#### Question Type: Code with TLE/WA/Runtime Error

**Template:**
```
This code gets [TLE/WA/Runtime Error] on test case: [...]

```cpp
[CODE WITH BUG]
```

What's wrong?
A) [Common misconception]
B) [Edge case issue]
C) [Complexity issue]
D) [Off-by-one error]

Fix the code.
```

**Example 1: DFS with Cycle (TLE)**
```
Problem: Count connected components in undirected graph
Test case that causes TLE: [[0,1], [1,2], [2,0]]

Code:
```cpp
int countComponents(int n, vector<vector<int>>& edges) {
    vector<vector<int>> graph(n);
    for (auto& e : edges) {
        graph[e[0]].push_back(e[1]);
        graph[e[1]].push_back(e[0]);
    }
    
    vector<bool> visited(n, false);
    int count = 0;
    
    function<void(int)> dfs = [&](int node) {
        visited[node] = true;
        for (int neighbor : graph[node]) {
            if (!visited[neighbor]) {
                dfs(neighbor);
            }
        }
    };
    
    for (int i = 0; i < n; i++) {
        if (!visited[i]) {
            dfs(i);
            count++;
        }
    }
    
    return count;
}
```

Why TLE?
A) visited check is wrong
B) Infinite recursion in cycle
C) Graph construction is wrong
D) No bug, just slow test case

Answer: A (seems correct!) and B (related)
Wait - the code IS correct! But on edge with cycle, seems to hang?

Actually, this code is CORRECT. The real TLE case would be:
- Missing visited check before recursing
- Or visited check inside loop but AFTER recursive call

Let me show the BUGGY version:
```cpp
function<void(int)> dfs = [&](int node) {
    visited[node] = true;
    for (int neighbor : graph[node]) {
        dfs(neighbor);  // BUG: no check!
        if (!visited[neighbor]) {  // Too late!
            // ...
        }
    }
};
```

Fix: Check BEFORE recursing:
```cpp
if (!visited[neighbor]) {
    dfs(neighbor);
}
```
```

**Example 2: Binary Search (WA)**
```
Problem: Find first position where arr[i] ≥ target
Failing test: arr = [1,2,2,2,3], target = 2
Expected: 1
Got: 3

Code:
```cpp
int lowerBound(vector<int>& arr, int target) {
    int left = 0, right = arr.size() - 1;
    
    while (left < right) {
        int mid = left + (right - left) / 2;
        
        if (arr[mid] < target) {
            left = mid + 1;
        } else {
            right = mid - 1;  // BUG!
        }
    }
    
    return left;
}
```

What's wrong?
A) Should be left <= right
B) Should be right = mid (not mid - 1)
C) Should be mid + 1 when arr[mid] < target
D) Should return right, not left

Answer: B
When arr[mid] >= target, mid might be the answer!
Setting right = mid - 1 skips it.

Fixed:
```cpp
right = mid;  // Keep mid as potential answer
```

Explanation: This is template for finding FIRST occurrence.
Pattern: When condition met, keep mid (right = mid)
         When condition not met, skip mid (left = mid + 1)
```

**Example 3: DP Array Bounds (Runtime Error)**
```
Problem: Climbing Stairs
Failing test: n = 1
Error: Segmentation fault

Code:
```cpp
int climbStairs(int n) {
    vector<int> dp(n);  // BUG!
    dp[0] = 1;
    dp[1] = 2;
    
    for (int i = 2; i < n; i++) {
        dp[i] = dp[i-1] + dp[i-2];
    }
    
    return dp[n-1];
}
```

What's wrong?
A) dp size should be n+1
B) When n=1, dp[1] is out of bounds
C) Loop should be i <= n
D) Should return dp[n], not dp[n-1]

Answer: B
When n=1, dp size is 1 (indices 0 only)
Accessing dp[1] crashes!

Fixed:
```cpp
if (n <= 2) return n;
vector<int> dp(n);
// ... rest of code
```

Or better:
```cpp
vector<int> dp(n + 1);  // Size n+1 to avoid edge cases
```
```

---

### 3.2 Reverse Engineering Test Cases

#### Question Type: Find the Breaking Test Case

**Template:**
```
This code looks correct but has a bug.
Find a test case that makes it fail.

```cpp
[CODE WITH SUBTLE BUG]
```

What test case breaks it?
A) [Edge case 1]
B) [Edge case 2]
C) [Large input]
D) [Special condition]

Why does it fail on your test case?
```

**Example 1: Binary Search**
```
Code:
```cpp
int binarySearch(vector<int>& arr, int target) {
    int left = 0, right = arr.size() - 1;
    
    while (left <= right) {
        int mid = (left + right) / 2;  // Bug here!
        
        if (arr[mid] == target) return mid;
        else if (arr[mid] < target) left = mid + 1;
        else right = mid - 1;
    }
    
    return -1;
}
```

Which test case causes problems?
A) arr = [1, 2, 3], target = 2
B) arr = [INT_MAX - 1, INT_MAX], target = INT_MAX
C) arr = [], target = 5
D) arr = [1], target = 1

Answer: B
Why? (left + right) overflows when both are near INT_MAX

Demo:
left = INT_MAX - 1 = 2147483646
right = INT_MAX = 2147483647
left + right = 4294967293 (overflows to negative!)
mid becomes negative → crash or wrong result

Fixed:
```cpp
int mid = left + (right - left) / 2;
```
```

**Example 2: Two Pointers**
```
Code (Remove duplicates from sorted array):
```cpp
int removeDuplicates(vector<int>& nums) {
    int i = 0;
    for (int j = 1; j < nums.size(); j++) {
        if (nums[j] != nums[i]) {
            nums[i + 1] = nums[j];
            i++;
        }
    }
    return i + 1;
}
```

Find the breaking test case:
A) nums = [1, 1, 2]
B) nums = [1, 2, 3]
C) nums = []
D) nums = [1, 1, 1, 1]

Answer: C
When nums is empty, nums.size() is 0
Loop doesn't run, but we return i + 1 = 1 (wrong! should be 0)

Also: j starts at 1, but if size is 1, loop doesn't run → returns 1 (correct)

Fixed:
```cpp
if (nums.empty()) return 0;
// ... rest of code
```
```

---

## PART 4: GRAPH OBSERVATION QUESTIONS

### 4.1 Key Observations Before Coding

**Philosophy:** Graph problems require key insights ("observations") before coding. Test these separately!

#### Question Type: Observation Testing

**Template:**
```
Problem: [Graph problem]

Before coding, identify key observations:

1. [Observation 1]: True/False + Why?
2. [Observation 2]: True/False + Why?
3. [Observation 3]: True/False + Why?

Which observations are necessary for optimal solution?
```

**Example 1: Course Schedule (Cycle Detection)**
```
Problem: Given prerequisites [[1,0], [2,1]], can you finish all courses?

Test observations:

**Observation 1:** "If there's a cycle, it's impossible"
A) True - can't finish courses in a cycle
B) False - we can skip courses in the cycle
C) Sometimes true
D) Irrelevant to problem

Answer: A

**Observation 2:** "We need to track in-degree for each node"
A) True for Kahn's algorithm
B) True for DFS
C) Not necessary
D) Only for DAG

Answer: A (for Kahn's), C (for DFS cycle detection)

**Observation 3:** "A back edge indicates a cycle"
A) True in undirected graph
B) True in directed graph only
C) True in DAG
D) Never true

Answer: B (directed), A is false (back edges exist in trees)

**Observation 4:** "If we visit a node currently in our DFS path, there's a cycle"
A) True
B) False
C) Only if it's the parent
D) Only in undirected graphs

Answer: A (this is the key insight for DFS cycle detection!)

Which observations do you need for each algorithm?
- Kahn's BFS: Observations 1, 2
- DFS: Observations 1, 4
```

**Example 2: Number of Islands**
```
Problem: Count islands (connected 1s) in a grid.

Test observations:

**Observation 1:** "Each connected component is one island"
A) True
B) False
C) Only if using BFS
D) Only for horizontal/vertical connections

Answer: A (but D is a clarification - diagonal doesn't count in standard problem)

**Observation 2:** "We can modify the input grid to mark visited"
A) Always allowed
B) Never allowed
C) Only if problem permits
D) Required for optimal solution

Answer: C (Check problem - usually allowed for interviews)

**Observation 3:** "DFS and BFS give same answer but different complexity"
A) True - same answer, same complexity
B) False - different answers
C) True - same answer, BFS faster
D) True - same answer, DFS easier to implement

Answer: A/D (same answer, same O(mn) time, DFS usually simpler code)

**Observation 4:** "We should use Union-Find instead of DFS/BFS"
A) Union-Find is better
B) DFS/BFS is better
C) Same complexity, preference
D) Depends on follow-up questions

Answer: D!
If follow-up asks: "What if grid updates dynamically?"
→ Union-Find is much better (O(α(n)) per query)
→ DFS/BFS requires recomputation (O(mn))
```

**Example 3: Shortest Path in Binary Matrix**
```
Problem: Find shortest path from (0,0) to (n-1,n-1) in binary matrix.

Test observations:

**Observation 1:** "Dijkstra's is needed because of weighted edges"
A) True
B) False - edges are unweighted (all cost = 1)
C) True if diagonal moves cost more
D) Only in 3D grids

Answer: B (BFS is sufficient for unweighted!)

**Observation 2:** "We can move in 8 directions"
A) Always true
B) Need to check problem statement
C) Always 4 directions
D) Doesn't matter

Answer: B (CRITICAL - easy to miss!)

**Observation 3:** "If start or end is blocked (1), immediately return -1"
A) True
B) False
C) Only check start
D) Only check end

Answer: A (common bug - forgetting to check!)

**Observation 4:** "We need visited set to avoid revisiting"
A) True - required
B) False - can mark in grid
C) False - BFS guarantees no revisit
D) Only for DFS

Answer: B (can mark grid[i][j] = 1 when visited)
```

---

### 4.2 Multi-Step Observation Chains

#### Question Type: Observation Dependencies

**Template:**
```
To solve [problem], we need these observations IN ORDER:

Step 1: Realize [observation A]
Step 2: This means [observation B]
Step 3: Therefore [observation C]
Step 4: So the approach is [solution]

At which step is the student stuck?
[Show their incorrect reasoning]
```

**Example: Topological Sort**
```
Problem: Course Schedule II - Return course order if possible.

Observation chain:

**Step 1:** Prerequisites form a directed graph
Student: ✓ Got it

**Step 2:** Cycle = impossible to complete
Student: ✓ Got it

**Step 3:** Need topological order (DAG ordering)
Student: ✗ "I'll just use DFS"
→ But DFS gives reverse postorder!

**Step 4:** Use Kahn's algorithm OR reverse DFS postorder
Student: ✗ Stuck

Test question: 
"Your DFS gives [3, 2, 1, 0] but answer should be [0, 1, 2, 3]. Why?"

A) DFS is wrong algorithm
B) Need to reverse the result
C) Need to track entry time, not exit time
D) Need to use BFS instead

Answer: B (DFS postorder is reverse topological order!)

Follow-up: "Why does DFS postorder give reverse?"
Hint: "What does postorder mean?"
→ We record node AFTER visiting all descendants
→ So dependencies (descendants) appear BEFORE node in postorder
→ Reverse this to get topological order
```

---

## PART 5: AI ASSISTANT INTEGRATION

### 5.1 Context for AI Hints

**What to send to AI:**
```json
{
  "question": {
    "id": 123,
    "type": "code_debugging",
    "problem": "Binary Search",
    "question_text": "Find the bug...",
    "user_answer": "Option B",
    "correct_answer": "Option A",
    "is_correct": false
  },
  "user_context": {
    "attempt_number": 3,
    "time_spent": 180,
    "previous_attempts": [
      {"answer": "C", "time": 45},
      {"answer": "D", "time": 60}
    ],
    "recent_performance": {
      "category": "binary_search",
      "success_rate": 0.4,
      "average_time": 120
    },
    "known_strengths": ["two_pointers", "arrays"],
    "known_weaknesses": ["edge_cases", "overflow_handling"]
  },
  "hint_level": 1  // 1 = gentle nudge, 2 = more direct, 3 = detailed
}
```

**AI Prompt Template:**
```
You are a patient coding interview tutor. A student is struggling with this question:

[Question details]

Student answered: {user_answer}
Correct answer: {correct_answer}
This is attempt #{attempt_number}

Student's profile:
- Strengths: {strengths}
- Weaknesses: {weaknesses}
- Success rate in this category: {success_rate}

Provide a hint at level {hint_level}:
Level 1: Socratic question or analogy (don't reveal answer)
Level 2: Point to specific part of problem (narrow down)
Level 3: Explain the concept needed (still let them apply it)

Be encouraging. Don't give the answer directly.
```

**Example AI Response (Level 1):**
```
"I notice you've been working on this for a while. Let me ask you something:

When we do (left + right) / 2, what are the possible values of left and right?

Think about the MAXIMUM values they could have based on the problem constraints."
```

**Example AI Response (Level 2):**
```
"You're on the right track thinking about the algorithm!

The issue is actually in this line:
    int mid = (left + right) / 2;

What happens when left and right are both very large numbers?
Have you learned about integer overflow?"
```

**Example AI Response (Level 3):**
```
"Great persistence! Let me explain the issue:

Integer overflow occurs when a calculation exceeds the maximum value an int can hold (about 2 billion).

In binary search, if we're searching a large array:
- left might be 2000000000
- right might be 2000000000
- left + right = 4000000000 (exceeds int max!)

The solution is to rewrite the calculation to avoid the overflow.
Instead of: (left + right) / 2
Use: left + (right - left) / 2

Can you see why this version doesn't overflow?"
```

---

### 5.2 Offline Hints (Pre-Generated)

**For each question, prepare 3 levels of hints:**

**Example: Binary Search Overflow**

```
Hint 1 (Gentle nudge):
"Consider what happens with very large array indices..."

Hint 2 (More specific):
"Look at the line where mid is calculated. What's the maximum value left + right could have?"

Hint 3 (Nearly explicit):
"When left and right are both near INT_MAX, their sum overflows. Try rewriting (left + right) / 2 without adding first."
```

**Example: DP Return Value**

```
Hint 1:
"Think about what your function returns when no solution exists..."

Hint 2:
"You're returning 0 for impossible cases. But 0 is also a valid answer (when amount is 0). How can you distinguish between them?"

Hint 3:
"Use -1 to indicate impossible. This is different from 0 (valid answer for base case) and different from INT_MAX (initialization value)."
```

---

### 5.3 Progressive Hint Delivery

**Hint Strategy:**

```python
def get_hint_level(attempt_number, time_spent, user_performance):
    """
    Determine hint level based on struggle indicators
    """
    if attempt_number == 1 and time_spent < 120:
        return 0  # No hint yet, let them try
    
    if attempt_number == 2:
        return 1  # Gentle nudge
    
    if attempt_number >= 3 or time_spent > 300:
        return 2  # More direct
    
    if attempt_number >= 5 or time_spent > 600:
        return 3  # Detailed explanation
    
    # Adjust based on category performance
    if user_performance.success_rate < 0.3:
        return min(3, current_level + 1)  # Struggling overall, help more
    
    return min(3, attempt_number)
```

---

### 5.4 Weakness Detection via AI

**After assessment, send to AI:**

```json
{
  "assessment_summary": {
    "total_questions": 20,
    "correct": 12,
    "categories": {
      "complexity_analysis": {"correct": 4, "total": 5},
      "code_debugging": {"correct": 2, "total": 5},
      "pattern_recognition": {"correct": 3, "total": 5},
      "edge_cases": {"correct": 3, "total": 5}
    }
  },
  "detailed_mistakes": [
    {
      "question_id": 42,
      "category": "code_debugging",
      "mistake_type": "missed_overflow",
      "question_summary": "Binary search with large array"
    },
    {
      "question_id": 51,
      "category": "code_debugging",
      "mistake_type": "missed_overflow",
      "question_summary": "Sum of array elements"
    },
    {
      "question_id": 38,
      "category": "edge_cases",
      "mistake_type": "empty_input",
      "question_summary": "Two pointers on empty array"
    }
  ]
}
```

**AI Analysis Prompt:**
```
Analyze this student's assessment results and identify:

1. Pattern in mistakes (are they related?)
2. Root cause (conceptual gap vs. careless mistakes)
3. Specific weakness (be precise)
4. Recommended practice (what to focus on)

Be specific. Don't say "practice more binary search" - say "focus on integer overflow edge cases in binary search".
```

**AI Response Example:**
```
Analysis:

**Pattern Detected:**
The student made overflow-related mistakes in 2 out of 5 debugging questions (40% of debugging errors).

**Root Cause:**
Conceptual gap: The student doesn't automatically think about integer overflow when seeing arithmetic operations with large numbers.

They CAN identify overflow when explicitly asked about it (got complexity questions right), but don't proactively check for it in code.

**Specific Weakness:**
"Integer overflow awareness in implementation"
- Knows the concept
- Forgets to apply it when writing/debugging code
- Needs to build habit of checking constraints first

**Recommended Practice:**
1. Before every problem, write down: "What's the max value of key variables?"
2. Practice 10 problems specifically focused on overflow edge cases
3. Create mental checklist: "Could any arithmetic overflow here?"

**Immediate next steps:**
- Review: "How to detect and prevent integer overflow"
- Practice: Leetcode tagged "overflow" or curated list of overflow-prone problems
- Build habit: Always check constraints before coding
```

---

## PART 6: ADDITIONAL QUESTION TYPES

### 6.1 Complexity Analysis Variations

#### Question Type: Amortized vs. Worst-Case

**Example:**
```
Dynamic array doubling strategy:

```cpp
void push_back(int x) {
    if (size == capacity) {
        capacity *= 2;
        reallocate();  // O(n) operation
    }
    arr[size++] = x;
}
```

What is the time complexity of n push_back operations?

A) O(n²) - each push might reallocate
B) O(n log n) - doubling happens log n times
C) O(n) - amortized O(1) per operation
D) O(n) - but only because n is large

Answer: C
Explain: Why is it O(n) total and not O(n²)?

Deep dive:
- Reallocation happens at sizes: 1, 2, 4, 8, 16, ...
- Total copies: 1 + 2 + 4 + 8 + ... + n/2 = n - 1 ≈ O(n)
- n operations total → O(n)/n = O(1) amortized

Follow-up: What if we doubled + 1 instead of doubling?
capacity = capacity * 2 + 1 vs. capacity = capacity + 1

A) Same amortized complexity
B) Worse - O(n) per operation
C) Better - O(1) worst case
D) Can't determine

Answer: A (still O(1) amortized for *2+1), B for +1 only
```

---

### 6.2 Space Complexity Traps

#### Question Type: Hidden Space Usage

**Example:**
```
Code:

```cpp
void reverseString(string& s) {
    int left = 0, right = s.size() - 1;
    while (left < right) {
        swap(s[left++], s[right--]);
    }
}
```

What is the space complexity?

A) O(1)
B) O(n)
C) O(log n)
D) Depends on string implementation

Answer: D (trick question!)

Explanation:
- Looks like O(1) - no extra data structures
- BUT: C++ strings use COW (Copy-On-Write) in some implementations
- Modifying might create a copy → O(n) space
- In modern C++11+, usually O(1) though

Better question: "What is the INTENDED space complexity?"
Answer: O(1) - we're doing in-place swap

Learning: Always clarify assumptions!
```

---

### 6.3 Trade-off Questions

#### Question Type: Time vs. Space Trade-offs

**Example:**
```
Problem: Check if array has duplicates

Solution A:
```cpp
bool hasDuplicates(vector<int>& arr) {
    for (int i = 0; i < arr.size(); i++) {
        for (int j = i + 1; j < arr.size(); j++) {
            if (arr[i] == arr[j]) return true;
        }
    }
    return false;
}
```

Solution B:
```cpp
bool hasDuplicates(vector<int>& arr) {
    unordered_set<int> seen;
    for (int num : arr) {
        if (seen.count(num)) return true;
        seen.insert(num);
    }
    return false;
}
```

Solution C:
```cpp
bool hasDuplicates(vector<int>& arr) {
    sort(arr.begin(), arr.end());
    for (int i = 1; i < arr.size(); i++) {
        if (arr[i] == arr[i-1]) return true;
    }
    return false;
}
```

Questions:

1. **Which is fastest?**
   A) A  B) B  C) C  D) Depends

   Answer: D (depends on n)
   - Small n: A might be fastest (cache-friendly)
   - Medium n: B is fastest
   - Large n with limited memory: C might be forced

2. **Which uses least memory?**
   A) A  B) B  C) C  D) Same

   Answer: A (O(1) space)
   - A: O(1)
   - B: O(n)
   - C: O(1) if we can modify input, otherwise O(n) for sorting

3. **In an interview, which to use?**
   A) A - simplest
   B) B - optimal time
   C) C - no extra space
   D) Ask about constraints!

   Answer: D then usually B
   
   Why ask constraints?
   - Is array mutable? (affects C)
   - Memory limit? (affects B vs C)
   - Is it sorted? (might change everything)
   - Expected array size? (affects A viability)

4. **Follow-up: What if array is already sorted?**
   Answer: Use modified version of all three, but A/C become O(n) - just check consecutive!

5. **Follow-up: What if we need to do this operation repeatedly?**
   Answer: B becomes even better - build set once, check in O(1) each time
```

---

## PART 7: SYSTEMATIC DEBUGGING APPROACH

### 7.1 Debugging Checklist Questions

**Template:**
```
You have a bug. Apply the debugging checklist:

Step 1: What does the code ACTUALLY do?
Step 2: What SHOULD it do?
Step 3: Where do they differ?
Step 4: Why does it differ there?
Step 5: How to fix?

[Show code with bug]
Answer each step for this student's attempt.
```

**Example: Binary Search Bug**
```
Student's code:
```cpp
int search(vector<int>& arr, int target) {
    int left = 0, right = arr.size() - 1;
    while (left < right) {  // Bug?
        int mid = left + (right - left) / 2;
        if (arr[mid] == target) return mid;
        if (arr[mid] < target) left = mid + 1;
        else right = mid - 1;
    }
    return -1;
}
```

Fails on: arr = [5], target = 5
Returns: -1 (expected: 0)

Debug checklist:

**Step 1: What does it do?**
→ "Loop while left < right"
→ "When left == right, exit"
→ "Never checked arr[left] when left == right"

**Step 2: What should it do?**
→ "Check all possible positions including when left == right"

**Step 3: Where differs?**
→ "When array size = 1, left = right = 0"
→ "Loop condition left < right is false"
→ "Never enters loop, returns -1"

**Step 4: Why?**
→ "Loop condition is too strict"
→ "Should be left <= right to check last position"

**Step 5: Fix**
→ Change to: while (left <= right)

Question: What other test case reveals this bug?
A) arr = [1, 2, 3], target = 1
B) arr = [1, 2, 3], target = 3
C) arr = [], target = 5
D) arr = [1, 2], target = 2

Answer: B (target at last position has similar issue)
```

---

## PART 8: CONCEPT APPLICATION TRANSFER

### 8.1 Recognizing Patterns Across Problems

**Question Type: Same Technique, Different Domain**

**Example:**
```
You learned "two pointers" on arrays.

Which of these problems can use the same technique?

A) Reverse a string
B) Linked list cycle detection
C) Container with most water
D) Merge two sorted arrays
E) All of the above

Answer: E

Now for each, explain:
- How is it "two pointers"?
- What do the pointers represent?
- What's the loop invariant?

**A) Reverse string:**
- Two pointers: left = 0, right = n-1
- Move towards center, swapping
- Invariant: Elements outside [left, right] are reversed

**B) Linked list cycle:**
- Two pointers: slow, fast
- Slow moves 1 step, fast moves 2
- Invariant: If cycle exists, they'll meet

**C) Container with water:**
- Two pointers: left = 0, right = n-1
- Move pointer at shorter line
- Invariant: We've considered all containers with current spread

**D) Merge sorted arrays:**
- Two pointers: i in arr1, j in arr2
- Take smaller element, advance that pointer
- Invariant: Result[0..k] contains k smallest elements

Key insight: "Two pointers" isn't one technique - it's a family!
Different problems, different movement strategies.
```

---

## PART 9: IMPLEMENTATION DETAILS

### 9.1 Language-Specific Gotchas

**Question Type: Python vs. C++ Differences**

**Example:**
```
Python code:
```python
def twoSum(nums, target):
    seen = {}
    for i, num in enumerate(nums):
        if target - num in seen:
            return [seen[target - num], i]
        seen[num] = i
```

Direct C++ translation:
```cpp
vector<int> twoSum(vector<int>& nums, int target) {
    unordered_map<int, int> seen;
    for (int i = 0; i < nums.size(); i++) {
        if (seen.count(target - nums[i])) {
            return {seen[target - nums[i]], i};
        }
        seen[nums[i]] = i;
    }
}
```

What's wrong with C++?

A) Nothing, it's correct
B) Missing return statement
C) map.count() doesn't work like 'in'
D) Should use map.find()

Answer: B (Missing return statement at end for compiler)

Better:
```cpp
return {};  // Or return {-1, -1};
```

Follow-up questions:
1. In Python, what happens if no solution exists?
   → Returns None implicitly

2. In C++, what happens without return statement?
   → Undefined behavior (crashes or garbage)

3. Should we check count() or use find()?
   Answer: count() is clearer, but find() is slightly faster
   ```cpp
   auto it = seen.find(target - nums[i]);
   if (it != seen.end()) {
       return {it->second, i};
   }
   ```

4. Why return {} in C++?
   → Returns empty vector (like Python's None/empty list)
```

---

## SUMMARY: QUESTION TYPE CATEGORIES

### Master List

1. **Understanding-Based**
   - Return value semantics
   - Initialization values
   - Variable naming/meaning
   - Loop invariants

2. **Transformation-Based**
   - Iterative ↔ Recursive
   - 2D DP → 1D DP
   - DFS ↔ BFS
   - Recursion → Iteration

3. **Debugging-Based**
   - Find the bug (TLE/WA/RE)
   - Explain why it's wrong
   - Find breaking test case
   - Fix the code

4. **Analysis-Based**
   - Time complexity (worst/average/amortized)
   - Space complexity (including hidden)
   - Trade-offs (time vs space)
   - When to use which algorithm

5. **Observation-Based**
   - Key insights before coding
   - Observation chains
   - Pattern recognition
   - Problem reduction

6. **Edge Case-Based**
   - Overflow handling
   - Empty input
   - Single element
   - Boundary conditions

7. **Application-Based**
   - Transfer to new domain
   - Recognize same pattern
   - Combine multiple techniques
   - Adapt to constraints

8. **Implementation-Based**
   - Language-specific gotchas
   - API usage
   - Off-by-one errors
   - Index management

---

## AI ASSISTANT GUIDELINES

**When to provide hints:**
- After 2-3 incorrect attempts
- After 5+ minutes on one question
- When user requests help
- When success rate in category < 40%

**How to provide hints:**
- Level 1: Socratic question
- Level 2: Point to relevant concept
- Level 3: Explain concept, let them apply

**What to track:**
- Common mistake patterns
- Time spent per question type
- Categories of weakness
- Improvement over time

**How to personalize:**
- Reference their previous attempts
- Connect to their known strengths
- Acknowledge progress
- Adjust difficulty accordingly

---

This guide should be used to generate 500+ high-quality questions that test *understanding*, not memorization!

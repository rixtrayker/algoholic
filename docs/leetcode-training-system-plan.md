# LeetCode Interview Training System
## Comprehensive Plan for Technical Interview Mastery

---

# TABLE OF CONTENTS

1. [Core Training Modules](#core-training-modules)
2. [Question Type Breakdown](#question-type-breakdown)
3. [Complexity Analysis Training](#complexity-analysis-training)
4. [C++ STL Mastery](#cpp-stl-mastery)
5. [Algorithm Pattern Library](#algorithm-pattern-library)
6. [Data Structure Selection](#data-structure-selection)
7. [Implementation Correctness](#implementation-correctness)
8. [Bug Detection & Fixing](#bug-detection-fixing)
9. [Pseudocode Algorithm Design](#pseudocode-algorithm-design)
10. [Approach Trade-offs](#approach-tradeoffs)
11. [Progressive Learning Path](#progressive-learning-path)

---

## 1. CORE TRAINING MODULES

### **Module Philosophy**

The problem with traditional LeetCode practice:
- **Too much time on syntax** → miss the algorithm
- **Jump to code too fast** → miss the thinking process
- **Don't understand WHY** → can't adapt to new problems
- **Practice randomly** → no systematic skill building

Our approach:
- **Separate algorithm thinking from implementation**
- **Master complexity analysis FIRST**
- **Understand trade-offs BEFORE memorizing solutions**
- **Build pattern recognition through categorization**
- **Use pseudocode to clarify logic before coding**

---

## 2. QUESTION TYPE BREAKDOWN

### 🎯 **Type 1: Complexity Analysis Questions**

**Goal**: Instantly recognize time/space complexity without running code

**Sub-Types**:

#### **A. Given Code, Identify Complexity**
```
Show code snippet with nested loops/recursion
User must identify: O(?) time, O(?) space
```

**Why This Matters**:
- Interview question: "What's the complexity of your solution?"
- Need to know BEFORE implementing
- Helps choose between approaches

**Example Question Format**:
```
Code snippet:
for i in 0 to n:
    for j in i to n:
        if array[i] + array[j] == target:
            return true

What is the time complexity?
A) O(n)
B) O(n log n)
C) O(n²)
D) O(2^n)

What is the space complexity?
A) O(1)
B) O(n)
C) O(n²)
D) O(log n)
```

**Difficulty Progression**:
- **Level 1**: Single loops, simple recursion
- **Level 2**: Nested loops, basic divide & conquer
- **Level 3**: Hidden complexity (amortized, library functions)
- **Level 4**: Complex recursion trees, master theorem
- **Level 5**: Tricky cases (sorted array binary search in loop)

---

#### **B. Complexity Comparison**

**Format**: Given 2-3 solutions, rank by efficiency

**Example**:
```
Problem: Find if array has duplicates

Solution A: Sort then scan for adjacent equals
Solution B: Use hash set
Solution C: Nested loop comparison

Rank these by time complexity (best to worst):
Options show different orders, user picks correct ranking
```

**Why This Matters**:
- Real interviews: "Can you do better than O(n²)?"
- Understanding space/time tradeoffs
- Choosing appropriate solution for constraints

**Variations**:
- Rank by space complexity
- Rank by both (which is overall best?)
- Consider best/worst/average case

---

#### **C. Complexity Under Constraints**

**Format**: Problem + constraints → What complexity is needed?

**Example**:
```
Problem: Search in array
Constraint: Array has 10^9 elements, 10^5 queries

What complexity is required to pass time limit?
A) O(n) per query is fine
B) Need O(log n) per query
C) Need O(1) per query
D) Need preprocessing + O(1) query

Explanation required!
```

**Why This Matters**:
- Real constraints drive solution choice
- Can't just memorize algorithms
- Need to calculate: will this TLE (Time Limit Exceed)?

---

#### **D. Hidden Complexity Detection**

**Format**: Code uses library functions - what's REAL complexity?

**Example**:
```
Code:
for i in 0 to n:
    substring = s.substr(i, k)
    if (set.count(substring)) return true

Appears to be O(n), but what's the ACTUAL complexity?

A) O(n)
B) O(n * k)
C) O(n * k * log n)
D) O(n²)

Why? What's the hidden operation?
```

**Why This Matters**:
- String operations often O(n)
- Set operations might be O(log n)
- Copying data structures costs time
- Real interviews test this knowledge

---

### 🎯 **Type 2: Data Structure Selection**

**Goal**: Choose optimal DS for given requirements

#### **A. Requirements-to-DS Mapping**

**Format**: Given operation requirements, pick best DS

**Example**:
```
Requirements:
- Insert element: Must be O(1) or O(log n)
- Find minimum: Must be O(1)
- Remove minimum: Must be O(log n) or better
- No need for random access

Which data structure is BEST?
A) Array (sorted)
B) Hash set
C) Min heap (priority queue)
D) BST

Why is your choice better than others?
```

**Difficulty Levels**:
- **Level 1**: Single requirement (fast lookup → hash)
- **Level 2**: Two requirements (insert + ordered → BST)
- **Level 3**: Multiple requirements with priorities
- **Level 4**: Tricky combinations (range query + updates)
- **Level 5**: Custom DS needed (combine multiple)

---

#### **B. DS Trade-off Analysis**

**Format**: Compare two DS choices for same problem

**Example**:
```
Problem: Implement LRU Cache

Option A: HashMap + Doubly Linked List
- Get: O(1)
- Put: O(1)
- Space: O(n)

Option B: HashMap + Array
- Get: O(1)
- Put: O(n) worst case
- Space: O(n)

Which is better for LRU cache? Why?
What operations make A better than B?
When might B be acceptable?
```

**Why This Matters**:
- No single "best" DS for everything
- Need to justify choices in interview
- Understanding WHY matters more than memorizing

---

#### **C. STL Container Selection (C++ Specific)**

**Format**: Given problem, pick correct STL container

**Example**:
```
Problem: Need to track last 100 elements seen in stream
Operations: Add new element (remove oldest if >100), check if element exists

Which STL container(s)?
A) vector
B) deque
C) set
D) unordered_set + deque
E) list + unordered_set

Explain your choice and the complexity of key operations.
```

**Key STL Containers to Master**:
- `vector` - when to use, when NOT to
- `deque` - vs vector, front insertion
- `list` - rarely optimal, when it is
- `set` / `multiset` - ordered, O(log n)
- `unordered_set` / `unordered_multiset` - unordered, O(1) average
- `map` / `multimap` - key-value, ordered
- `unordered_map` / `unordered_multimap` - key-value, unordered
- `priority_queue` - heap operations
- `stack` - LIFO
- `queue` - FIFO

---

#### **D. When to Build Custom DS**

**Format**: Identify when STL isn't enough

**Example**:
```
Problem: Need fast insert, delete, and find median

No single STL container solves this optimally.

What combination would you use?
A) Just vector (sort each time)
B) Two priority_queues (max heap + min heap)
C) multiset + iterator to median
D) Custom balanced BST

Analyze time complexity of each operation for your choice.
```

**Why This Matters**:
- Shows advanced thinking
- Real problems often need combinations
- Demonstrates problem-solving, not just memorization

---

### 🎯 **Type 3: Algorithm Pattern Recognition**

**Goal**: Instantly map problem type to algorithm category

#### **A. Problem-to-Pattern Mapping**

**Format**: Read problem, identify which pattern

**Example**:
```
Problem: "Find longest increasing subsequence"

This is a:
A) Two pointers problem
B) Dynamic programming problem
C) Greedy problem
D) Graph problem

Why? What's the key indicator in the problem statement?
```

**Pattern Categories** (from your cheatsheets):
1. Two Pointers (opposite/same direction)
2. Sliding Window (fixed/variable)
3. Binary Search (on array/on answer)
4. DFS/BFS (tree/graph traversal)
5. Dynamic Programming (1D/2D/knapsack/etc)
6. Greedy (scheduling/intervals)
7. Backtracking (generate all/pruning)
8. Union-Find (connectivity)
9. Topological Sort (dependencies)
10. Monotonic Stack/Queue (next greater/sliding window max)

**Difficulty Progression**:
- **Level 1**: Obvious keywords ("shortest path" → BFS)
- **Level 2**: Hidden patterns (need to recognize structure)
- **Level 3**: Multiple applicable patterns (choose best)
- **Level 4**: Hybrid problems (2+ patterns combined)
- **Level 5**: Novel problems (create new pattern)

---

#### **B. Pattern Variant Recognition**

**Format**: Given solved problem, identify variant

**Example**:
```
You know: Two Sum (hash map, O(n))

Which of these uses the SAME pattern?
A) Three Sum
B) Subarray Sum Equals K
C) Longest Substring Without Repeating Characters
D) Two Sum II (sorted array)

For each, explain: same pattern or different? Why?
```

**Why This Matters**:
- One pattern solves 10+ problems
- Interview: "This is similar to X problem"
- Transfer learning between problems

---

#### **C. Anti-Pattern Recognition**

**Format**: Identify why an approach WON'T work

**Example**:
```
Problem: Find all paths in graph from source to destination

Wrong approach: Use Dijkstra's shortest path algorithm

Why won't this work?
A) Dijkstra finds ONE shortest path, not all paths
B) Dijkstra assumes non-negative weights
C) Dijkstra is for weighted graphs only
D) All of the above

What's the CORRECT approach for this problem?
```

**Why This Matters**:
- Avoid wasting time on wrong approaches
- Interviewers test this explicitly
- Shows deep understanding vs surface memorization

---

### 🎯 **Type 4: C++ STL Deep Knowledge**

**Goal**: Master STL operations, complexities, and gotchas

#### **A. STL Function Complexity**

**Format**: Given STL operation, identify complexity

**Example**:
```
vector<int> v = {1,2,3,4,5};
v.insert(v.begin(), 0);  // insert at front

Time complexity?
A) O(1)
B) O(log n)
C) O(n)
D) O(n log n)

Why? What's happening internally?

Alternative: Use deque instead?
```

**Key STL Operations to Know**:

**Vector**:
- `push_back()` - O(1) amortized
- `insert(begin)` - O(n)
- `erase(begin)` - O(n)
- `size()` - O(1)

**Deque**:
- `push_front()` - O(1)
- `push_back()` - O(1)
- `insert(middle)` - O(n)

**Set/Map (Red-Black Tree)**:
- `insert()` - O(log n)
- `find()` - O(log n)
- `erase()` - O(log n)
- `lower_bound()` - O(log n)

**Unordered_Set/Map (Hash Table)**:
- `insert()` - O(1) average, O(n) worst
- `find()` - O(1) average, O(n) worst
- `erase()` - O(1) average

**Priority_Queue**:
- `push()` - O(log n)
- `pop()` - O(log n)
- `top()` - O(1)

---

#### **B. STL Algorithm Usage**

**Format**: Choose correct STL algorithm for task

**Example**:
```
Task: Find first element >= target in sorted vector

Which STL algorithm?
A) find()
B) binary_search()
C) lower_bound()
D) upper_bound()

What's the complexity of your choice?
What does it return?
```

**Critical STL Algorithms**:
- `sort()` - O(n log n)
- `stable_sort()` - O(n log n), preserves order
- `binary_search()` - O(log n), returns bool
- `lower_bound()` - O(log n), returns iterator
- `upper_bound()` - O(log n), returns iterator
- `nth_element()` - O(n), partial sort
- `partition()` - O(n), quicksort helper
- `accumulate()` - O(n), sum/fold
- `reverse()` - O(n)
- `rotate()` - O(n)

---

#### **C. STL Gotchas & Pitfalls**

**Format**: Identify the bug in STL usage

**Example**:
```
vector<int> v = {1,2,3,4,5};
for (auto it = v.begin(); it != v.end(); ++it) {
    if (*it % 2 == 0) {
        v.erase(it);
    }
}

What's wrong with this code?
A) Nothing, works fine
B) Iterator invalidation after erase
C) Should use it--, not ++it
D) Should check if it != end before dereferencing

How to fix it?
```

**Common Gotchas**:
- Iterator invalidation (vector, map, set)
- Reference invalidation (vector reallocation)
- `unordered_map` rehashing
- Comparing iterators from different containers
- `end()` is past-the-end, not last element
- `erase()` returns iterator to next element
- `remove()` doesn't actually remove (use erase-remove idiom)

---

#### **D. STL vs Manual Implementation**

**Format**: When to use STL vs write custom

**Example**:
```
Problem: Implement min heap with decrease-key operation

STL priority_queue doesn't support decrease-key.

What should you do?
A) Use set/map instead (ordered, can erase + reinsert)
B) Use priority_queue with lazy deletion
C) Implement custom heap
D) Depends on constraints

Explain trade-offs of each approach.
```

---

### 🎯 **Type 5: Code Pattern/Macro Library**

**Goal**: Memorize reusable code snippets for speed

#### **A. Essential Code Patterns**

**Format**: Fill in the blank / match pattern to use case

**1. Graph Adjacency List Setup**
```
Pattern name: "Graph Adjacency List"

When to use: Any graph problem

Template:
vector<vector<int>> adj(n);  // for n nodes
// For edge u -> v:
adj[u].push_back(v);

Variations:
- Weighted: vector<vector<pair<int,int>>> adj(n);  // {neighbor, weight}
- Undirected: add both adj[u].push(v) and adj[v].push(u)
```

**2. DFS Template**
```
Pattern name: "Recursive DFS"

When to use: Tree/graph traversal, backtracking

Template:
void dfs(int node, vector<vector<int>>& adj, vector<bool>& visited) {
    visited[node] = true;
    
    for (int neighbor : adj[node]) {
        if (!visited[neighbor]) {
            dfs(neighbor, adj, visited);
        }
    }
}
```

**3. BFS Template**
```
Pattern name: "Level-order BFS"

When to use: Shortest path, level traversal

Template:
queue<int> q;
q.push(start);
visited[start] = true;

while (!q.empty()) {
    int node = q.front();
    q.pop();
    
    for (int neighbor : adj[node]) {
        if (!visited[neighbor]) {
            visited[neighbor] = true;
            q.push(neighbor);
        }
    }
}
```

**4. Binary Search Template**
```
Pattern name: "Binary Search - Lower Bound"

When to use: Find first element >= target

Template:
int left = 0, right = n;
while (left < right) {
    int mid = left + (right - left) / 2;
    if (arr[mid] < target) {
        left = mid + 1;
    } else {
        right = mid;
    }
}
return left;
```

**5. Two Pointers Template**
```
Pattern name: "Two Pointers - Opposite Direction"

When to use: Sorted array, pair sum

Template:
int left = 0, right = n - 1;
while (left < right) {
    if (condition) {
        // found
    } else if (sum < target) {
        left++;
    } else {
        right--;
    }
}
```

**6. Sliding Window Template**
```
Pattern name: "Sliding Window - Variable Size"

When to use: Substring problems, subarray sum

Template:
int left = 0;
for (int right = 0; right < n; right++) {
    // add right to window
    
    while (window_invalid) {
        // remove left from window
        left++;
    }
    
    // update result
}
```

**7. Union-Find Template**
```
Pattern name: "Union-Find with Path Compression"

When to use: Connectivity, cycles, MST

Template:
vector<int> parent(n);
iota(parent.begin(), parent.end(), 0);

int find(int x) {
    if (parent[x] != x) parent[x] = find(parent[x]);
    return parent[x];
}

bool unite(int x, int y) {
    int px = find(x), py = find(y);
    if (px == py) return false;
    parent[py] = px;
    return true;
}
```

**8. Monotonic Stack Template**
```
Pattern name: "Monotonic Decreasing Stack"

When to use: Next greater element

Template:
vector<int> result(n, -1);
stack<int> st;  // stores indices

for (int i = 0; i < n; i++) {
    while (!st.empty() && arr[st.top()] < arr[i]) {
        result[st.top()] = arr[i];
        st.pop();
    }
    st.push(i);
}
```

**9. Trie Template**
```
Pattern name: "Trie (Prefix Tree)"

When to use: Prefix matching, autocomplete

Template:
struct TrieNode {
    TrieNode* children[26] = {};
    bool isEnd = false;
};

void insert(TrieNode* root, string word) {
    TrieNode* node = root;
    for (char c : word) {
        int idx = c - 'a';
        if (!node->children[idx]) {
            node->children[idx] = new TrieNode();
        }
        node = node->children[idx];
    }
    node->isEnd = true;
}
```

**10. Dijkstra Template**
```
Pattern name: "Dijkstra's Shortest Path"

When to use: Weighted graph, single-source shortest path

Template:
vector<long long> dist(n, LLONG_MAX);
priority_queue<pair<long long,int>, 
               vector<pair<long long,int>>, 
               greater<>> pq;

dist[start] = 0;
pq.push({0, start});

while (!pq.empty()) {
    auto [d, node] = pq.top();
    pq.pop();
    
    if (d > dist[node]) continue;
    
    for (auto [neighbor, weight] : adj[node]) {
        if (dist[node] + weight < dist[neighbor]) {
            dist[neighbor] = dist[node] + weight;
            pq.push({dist[neighbor], neighbor});
        }
    }
}
```

---

#### **B. Pattern Selection Quiz**

**Format**: Given problem, pick which template to use

**Example**:
```
Problem: "Find the shortest path in unweighted graph"

Which template should you use?
A) DFS
B) BFS
C) Dijkstra
D) Binary Search

Why is your choice optimal?
What's the time complexity?
```

---

#### **C. Template Customization**

**Format**: Given template + modification needed

**Example**:
```
Template: Standard BFS (finds shortest path)

Modification needed: Track the actual path, not just distance

What changes are required?
A) Add parent array, backtrack from end
B) Store full path in each queue element
C) Use DFS instead
D) Not possible with BFS

Which approach is more space efficient? Why?
```

---

### 🎯 **Type 6: Implementation Correctness**

**Goal**: Identify correct vs buggy implementations

#### **A. Spot the Correct Implementation**

**Format**: 3-4 implementations, only 1 is fully correct

**Example**:
```
Problem: Reverse a linked list

Implementation A:
[Shows code with subtle bug - forgets to update next pointer]

Implementation B:
[Shows correct code]

Implementation C:
[Shows code with off-by-one error]

Which implementation is FULLY correct?
What's wrong with the others?
```

**Why This Matters**:
- Code review skill
- Attention to detail
- Understanding edge cases

**Common Bug Types to Test**:
- Off-by-one errors
- Null pointer handling
- Integer overflow
- Array bounds
- Empty input handling
- Graph disconnected components
- Tree edge cases (null, single node)

---

#### **B. Edge Case Coverage**

**Format**: Given implementation, identify missing edge case

**Example**:
```
Implementation: Binary search in rotated sorted array

Code: [shows implementation]

This code fails on which edge case?
A) Array with duplicates
B) Array not actually rotated (fully sorted)
C) Target not in array
D) Array of size 1

How would you fix it?
```

---

#### **C. Boundary Condition Quiz**

**Format**: Multiple choice on boundary handling

**Example**:
```
When implementing binary search:

The condition "left < right" vs "left <= right" changes:
A) Whether we include the last element in search
B) Whether we return left or right at the end
C) Both A and B
D) Neither - they're equivalent

Explain when to use each.
```

---

### 🎯 **Type 7: Bug Detection & Fixing**

**Goal**: Debug broken code quickly

#### **A. Find the Bug**

**Format**: Working solution with 1-2 bugs inserted

**Example**:
```
Problem: Merge two sorted lists

Code:
ListNode* mergeTwoLists(ListNode* l1, ListNode* l2) {
    ListNode dummy(0);
    ListNode* curr = &dummy;
    
    while (l1 && l2) {
        if (l1->val < l2->val) {
            curr->next = l1;
            l1 = l1->next;
        } else {
            curr->next = l2;
            l2 = l2->next;
        }
    }
    
    // BUG IS HERE - find it
    if (l1) curr = l1;
    if (l2) curr = l2;
    
    return dummy.next;
}

What's wrong? How to fix?
```

**Bug Categories**:
- Logic errors
- Pointer errors
- Loop errors
- Comparison errors
- Edge case misses

---

#### **B. Fix the Bug**

**Format**: Bug location shown, provide fix

**Example**:
```
Bug identified on line 15:
if (l1) curr = l1;  // WRONG

What should it be?
A) if (l1) curr = l1;
B) if (l1) curr->next = l1;
C) while (l1) { curr->next = l1; l1 = l1->next; curr = curr->next; }
D) curr = l1 ? l1 : l2;

Why does your fix work?
```

---

#### **C. Output Prediction with Bug**

**Format**: Given buggy code + input, predict output

**Example**:
```
Buggy code: [shows code with subtle bug]

Input: [1, 2, 3, 4, 5]

What's the actual output?
A) [1, 2, 3, 4, 5]
B) [5, 4, 3, 2, 1]
C) Segmentation fault
D) Infinite loop

Why does the bug cause this behavior?
```

---

### 🎯 **Type 8: Pseudocode Algorithm Design**

**Goal**: Focus on algorithm logic, not syntax

#### **A. Problem to Pseudocode**

**Format**: Given problem, write high-level algorithm

**Example**:
```
Problem: Find median of two sorted arrays

Write pseudocode algorithm (no actual code):

Expected format:
1. Compare sizes of arrays
2. If size1 > size2, swap to ensure size1 <= size2
3. Binary search on smaller array to find partition
4. ...

Then explain:
- Why this approach works
- Time complexity
- Space complexity
```

**Why Pseudocode Matters**:
- Focuses on algorithm, not syntax
- Easier to communicate in interview
- Tests understanding, not memorization
- Faster to iterate on ideas

---

#### **B. Pseudocode to Complexity**

**Format**: Given pseudocode, derive complexity

**Example**:
```
Pseudocode:
1. Sort array A of size n
2. For each element in array A:
   3. Binary search in array B of size m
4. Return results

Time complexity?
A) O(n log n)
B) O(n log m)
C) O(n log n + n log m)
D) O(nm)

Show your reasoning step by step.
```

---

#### **C. Optimize Pseudocode**

**Format**: Given naive pseudocode, improve it

**Example**:
```
Naive pseudocode for "Find duplicates in array":

1. For each element i:
   2. For each element j > i:
      3. If arr[i] == arr[j]:
         4. Add to result

Current complexity: O(n²)

Provide improved pseudocode with better complexity.
What data structure would you use?
```

---

### 🎯 **Type 9: Approach Trade-offs**

**Goal**: Deeply understand when to use which approach

#### **A. Iterative vs Recursive DP**

**Format**: Compare both approaches for same problem

**Example**:
```
Problem: Fibonacci numbers

Recursive DP (Memoization):
```pseudocode
memo = {}
function fib(n):
    if n in memo: return memo[n]
    if n <= 1: return n
    memo[n] = fib(n-1) + fib(n-2)
    return memo[n]
```

Iterative DP (Tabulation):
```pseudocode
dp = array of size n+1
dp[0] = 0, dp[1] = 1
for i from 2 to n:
    dp[i] = dp[i-1] + dp[i-2]
return dp[n]
```

Compare on:
1. Time complexity (both O(n))
2. Space complexity (recursive has call stack)
3. Ease of understanding
4. Iteration order control
5. When would you use each?

Correct answer:
- Recursive: Easier to write, natural for tree problems
- Iterative: Better space, no stack overflow, easier to optimize
- Use recursive when: Problem naturally recursive (trees), top-down easier
- Use iterative when: Bottom-up obvious, space critical, n is large
```

**Key Trade-offs to Understand**:

**Recursive DP (Top-Down)**:
- ✅ More intuitive (matches problem definition)
- ✅ Only computes needed states
- ✅ Easier to write for complex state transitions
- ❌ Call stack overhead (can cause stack overflow)
- ❌ Slower due to function call overhead
- ❌ Harder to optimize space

**Iterative DP (Bottom-Up)**:
- ✅ No stack overflow risk
- ✅ Usually faster (no function calls)
- ✅ Easier to optimize space (rolling array)
- ✅ Clear computation order
- ❌ Must compute ALL states
- ❌ Less intuitive for complex problems
- ❌ Harder to write initially

**When to Choose**:
```
Choose Recursive when:
- Tree/graph problems (natural recursion)
- Complex state transitions
- Don't need all subproblems
- n is small (<10,000)

Choose Iterative when:
- Linear/matrix problems
- Simple state transitions
- Need all subproblems anyway
- n is large (risk of stack overflow)
- Space optimization critical
```

---

#### **B. DFS vs BFS Trade-offs**

**Format**: When to use which traversal

**Example**:
```
Problem: Find if path exists between two nodes

DFS approach:
- Space: O(h) where h = height
- Explores one path completely before backtracking
- Uses stack (recursion or explicit)

BFS approach:
- Space: O(w) where w = max width of tree
- Explores level by level
- Uses queue

For each scenario, choose DFS or BFS and explain:

1. Find shortest path in unweighted graph
   → BFS (level-order guarantees shortest)

2. Find ANY path in large deep tree
   → DFS (less space, finds fast if path is deep)

3. Check if graph is bipartite
   → BFS (but DFS works too)

4. Find connected components
   → Either (both work equally well)

5. Topological sort
   → DFS (natural post-order)

6. Level-order traversal
   → BFS (that's what it does)

7. Detect cycle in directed graph
   → DFS (easier to track current path)

8. Shortest path in unweighted graph
   → BFS (guaranteed shortest)
```

**Decision Matrix**:
```
             | DFS                    | BFS
-------------+------------------------+----------------------
Space        | O(height)              | O(width)
Use Case     | Any path, deep search  | Shortest path, levels
When to Use  | Tree, small width      | Graph, large height
Cycle detect | Easier                 | Possible but harder
Path finding | Any path fast          | Shortest path
Implementation| Recursive natural     | Queue required
```

---

#### **C. Hash Table vs Binary Search Tree**

**Format**: Choose correct DS for requirements

**Example**:
```
Requirement Matrix:

Operation           | Hash Table | BST (set/map)
--------------------+------------+--------------
Insert              | O(1) avg   | O(log n)
Search              | O(1) avg   | O(log n)
Delete              | O(1) avg   | O(log n)
Find min            | O(n)       | O(log n)
Range query         | O(n)       | O(k log n)
Ordered traversal   | O(n log n) | O(n)
Space               | O(n)       | O(n)

Scenarios:

1. Need fast lookup, no ordering needed
   → Hash Table (unordered_set/map)

2. Need to find predecessor/successor
   → BST (set/map with lower_bound)

3. Need to iterate in sorted order frequently
   → BST (in-order traversal is sorted)

4. Need to find all elements in range [a, b]
   → BST (lower_bound to upper_bound)

5. Just checking existence, no other operations
   → Hash Table (fastest)

6. Need to maintain sliding window with min/max
   → Neither - use deque or multiset

Real interview question:
"Implement a data structure that supports:
- Insert in O(1)
- Delete in O(1)  
- GetRandom in O(1)"

Answer: Hash table + array (not BST)
```

---

#### **D. Array vs Linked List**

**Format**: Understand when each shines

**Example**:
```
Operation          | Array      | Linked List
-------------------+------------+-------------
Access by index    | O(1)       | O(n)
Insert at end      | O(1)*      | O(1)
Insert at start    | O(n)       | O(1)
Insert in middle   | O(n)       | O(1)**
Delete at end      | O(1)       | O(1)***
Delete at start    | O(n)       | O(1)
Delete in middle   | O(n)       | O(1)**
Memory overhead    | Low        | High (pointers)
Cache friendly     | Yes        | No
Dynamic size       | Need resize| Natural

* amortized, can be O(n) on resize
** assuming you have pointer to position
*** assuming you have pointer to previous node

When to use Array:
- Random access needed
- Size known or rarely changes
- Sorting, binary search
- Numeric computations

When to use Linked List:
- Frequent insertions/deletions at start
- Size unknown and dynamic
- Iterator invalidation is issue
- Don't need random access
```

---

#### **E. Greedy vs DP**

**Format**: Recognize when greedy works vs needs DP

**Example**:
```
Problem Type 1: Activity Selection (intervals)
Greedy works! Sort by end time, take earliest ending.

Problem Type 2: Coin Change (fewest coins)
Greedy fails with coins [1, 5, 6], target 11
Greedy: 6 + 5 = 2 coins
Optimal: 5 + 5 + 1 = 3 coins (WRONG!)
Actual: 5 + 6 = 2 coins
DP required!

How to tell if Greedy works?
1. Greedy choice property (local optimal → global optimal)
2. Optimal substructure (optimal solution contains optimal subsolutions)
3. Can prove greedy is safe (exchange argument)

Examples where Greedy works:
- Huffman coding
- Dijkstra's algorithm (non-negative weights)
- Kruskal's MST
- Interval scheduling
- Fractional knapsack

Examples needing DP:
- 0/1 Knapsack
- Coin change
- Longest common subsequence
- Matrix chain multiplication

Interview tip: If greedy seems to work, try to find counterexample!
If you can't, and it has optimal substructure, likely correct.
```

---

#### **F. Sorting Algorithm Selection**

**Format**: Which sort for which situation?

**Example**:
```
Algorithm      | Best    | Average  | Worst   | Space | Stable | When to Use
---------------+---------+----------+---------+-------+--------+--------------
QuickSort      | O(n log)| O(n log) | O(n²)   | O(log)| No     | General, cache friendly
MergeSort      | O(n log)| O(n log) | O(n log)| O(n)  | Yes    | Need stable, linked lists
HeapSort       | O(n log)| O(n log) | O(n log)| O(1)  | No     | Limited memory
CountingSort   | O(n+k)  | O(n+k)   | O(n+k)  | O(k)  | Yes    | Small integer range
RadixSort      | O(nk)   | O(nk)    | O(nk)   | O(n+k)| Yes    | Fixed-length strings/ints

Scenarios:

1. Sort 1 million integers, range 0-999
   → Counting Sort (O(n))

2. Sort strings lexicographically, stability matters
   → Merge Sort

3. Sort in-place, average case performance
   → Quick Sort (STL sort uses introsort: quick + heap + insertion)

4. Sort linked list
   → Merge Sort (O(1) space on linked list)

5. Already mostly sorted
   → Insertion Sort (O(n) best case)

6. Need guaranteed O(n log n) worst case
   → Heap Sort or Merge Sort
```

---

### 🎯 **Type 10: Hybrid Challenge Questions**

**Goal**: Combine multiple skills in one question

#### **Format: Multi-Part Analysis**

**Example**:
```
Problem: Design a data structure for time-based key-value store

Requirements:
- set(key, value, timestamp): Store value with timestamp
- get(key, timestamp): Return value at or before timestamp

Part 1 - Data Structure Selection:
What DS would you use? Why?
Options:
A) Hash map of key → value
B) Hash map of key → sorted map of timestamp → value
C) Hash map of key → array of (timestamp, value)
D) Trie of keys + array values

Part 2 - Complexity Analysis:
For your choice, analyze:
- set() complexity: ?
- get() complexity: ?
- Space complexity: ?

Part 3 - Implementation Correctness:
Which is the correct get() implementation?
[Show 3 versions with subtle differences]

Part 4 - Optimization:
If timestamps are always increasing for a key, how can you optimize?

Part 5 - Trade-off Analysis:
Compare your solution to alternative approaches.
When would you use different DS?
```

---

## 10. PROGRESSIVE LEARNING PATH

### **Phase 1: Foundations (Weeks 1-2)**

**Focus**: Complexity analysis + basic patterns

**Daily Practice**:
- 5 complexity analysis questions
- 3 pattern recognition questions
- 2 STL usage questions
- Read 1 algorithm template

**Mastery Goals**:
- Instantly identify O(n), O(n²), O(log n)
- Recognize 5 core patterns
- Know 10 STL containers + operations
- Memorize 5 code templates

---

### **Phase 2: Pattern Mastery (Weeks 3-4)**

**Focus**: Algorithm patterns + pseudocode

**Daily Practice**:
- 3 pseudocode design questions
- 5 pattern application questions
- 2 approach comparison questions
- Implement 2 templates from memory

**Mastery Goals**:
- Write pseudocode for any problem
- Match problem to pattern in <30 seconds
- Explain trade-offs of 3 approaches
- Code 15 templates without reference

---

### **Phase 3: Deep Dive (Weeks 5-6)**

**Focus**: Trade-offs + implementation

**Daily Practice**:
- 3 trade-off analysis questions
- 5 implementation correctness questions
- 2 bug detection questions
- 1 hybrid challenge

**Mastery Goals**:
- Choose optimal approach with justification
- Spot bugs in <2 minutes
- Implement correctly first time
- Handle all edge cases

---

### **Phase 4: Speed & Polish (Weeks 7-8)**

**Focus**: Interview simulation

**Daily Practice**:
- 2 full problems (45 min each)
- 5 rapid-fire complexity questions
- 3 STL gotcha questions
- 1 design question

**Mastery Goals**:
- Complete medium problem in 25 minutes
- Explain solution clearly while coding
- Handle follow-up questions
- No syntax errors

---

## 11. QUESTION DIFFICULTY CALIBRATION

### **Easy (35% of practice)**
- Single pattern application
- Clear complexity
- Standard STL usage
- 1-2 edge cases

**Example**: Two Sum with hash map

---

### **Medium (50% of practice)**
- 2 patterns combined
- Hidden complexity considerations
- Custom DS combinations
- 3-5 edge cases

**Example**: LRU Cache (hash map + doubly linked list)

---

### **Hard (15% of practice)**
- 3+ patterns or novel approach
- Complex complexity analysis
- Advanced DS needed
- Many edge cases or optimization

**Example**: Median of Two Sorted Arrays

---

## 12. INTERVIEW SIMULATION FORMAT

### **Mock Interview Structure**

**Part 1: Problem Statement (2 min)**
- Read problem
- Ask clarifying questions
- Confirm constraints

**Part 2: Approach Discussion (8 min)**
- Identify pattern
- Discuss complexity
- Compare approaches
- Choose optimal

**Part 3: Pseudocode (5 min)**
- Write high-level algorithm
- Get feedback
- Adjust if needed

**Part 4: Implementation (25 min)**
- Code the solution
- Explain while coding
- Handle edge cases

**Part 5: Testing (5 min)**
- Walk through examples
- Identify edge cases
- Fix bugs

**Total: 45 minutes per problem**

---

## 13. COMMON INTERVIEW QUESTIONS ANSWERED

### **"Can you do better?"**

**Preparation**:
- Know complexity of your solution
- Know theoretical lower bound
- Know one better approach if exists

**Example Response**:
"My current solution is O(n²). The theoretical lower bound for this problem is O(n log n) because [reason]. I can achieve this by using [data structure/algorithm] instead."

---

### **"What are the trade-offs?"**

**Preparation**:
- Always think time vs space
- Consider best/average/worst case
- Think about practical factors (cache, constants)

**Example Response**:
"Approach A is O(n) time but O(n) space, while Approach B is O(n log n) time but O(1) space. I'd choose A if memory isn't constrained because the constant factors are better and cache friendliness matters."

---

### **"How would you handle this edge case?"**

**Preparation**:
- Always consider: empty input, single element, duplicates, negatives, overflow
- Think about boundary conditions
- Consider invalid input

**Example Response**:
"For empty input, I'd return early with [default value]. For the overflow case, I'd use long long instead of int here [point to line]."

---

## 14. METACOGNITIVE TRAINING

### **Self-Questioning During Practice**

After each problem:
1. **What pattern was this?** Can I categorize it?
2. **What was the key insight?** The "aha moment"?
3. **What did I miss initially?** How can I spot this next time?
4. **Which part took longest?** Algorithm design or implementation?
5. **What would I do differently?** Better approach exists?

---

### **Mistake Journal**

Track every mistake:
- What type? (Logic, syntax, edge case, complexity)
- Root cause? (Rushed, didn't understand, forgot)
- How to prevent? (Double-check X, remember Y)

**Common patterns in mistakes reveal weak spots**

---

### **Pattern Recognition Training**

Build mental database:
- "Ah, this is like Two Sum but with X twist"
- "This constraint screams binary search"
- "That's clearly a DP state transition"

**Goal: 10 second pattern identification**

---

## 15. FINAL EXAM FORMAT

### **Comprehensive Assessment**

**Section 1: Rapid Fire (20 questions, 15 min)**
- Complexity identification
- STL operations
- Pattern matching
- Quick bug spots

**Section 2: Deep Analysis (5 questions, 25 min)**
- Trade-off comparisons
- Approach justification
- Optimization proposals
- Edge case handling

**Section 3: Implementation (2 problems, 60 min)**
- One medium problem
- One medium-hard problem
- Full interview simulation

**Section 4: Design Question (1 question, 30 min)**
- System design or data structure design
- Multiple follow-ups
- Scalability discussion

**Total: 130 minutes**

**Passing Score**: 80% overall, 70% minimum in each section

---

## SUMMARY

This training system works because:

1. **Separates concerns**: Algorithm vs implementation
2. **Focuses on understanding**: Why, not just how
3. **Builds systematically**: Foundation → application → mastery
4. **Emphasizes patterns**: Transferable knowledge
5. **Simulates reality**: Actual interview conditions
6. **Tracks progress**: Visible improvement metrics
7. **Metacognitive**: Learn how to learn

**The goal isn't to memorize solutions. It's to build intuition.**

When you see a problem, you should INSTANTLY know:
- What pattern is this?
- What's the complexity I need?
- Which DS/algorithm to use?
- What are the edge cases?
- How to implement it correctly?

**That's mastery. That's what gets you the offer.**

---

*Document created for interview preparation*
*Next step: Build question database for each category*
*Target: 1000+ high-quality questions across all types*
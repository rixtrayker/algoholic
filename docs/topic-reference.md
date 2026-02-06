# Algoholic — DSA Topic Reference
## Complete Taxonomy for Database Seeding & Question Generation

---

## Purpose

Single source of truth for:
- Complete DSA topic hierarchy with patterns, pitfalls, and edge cases
- Database seeding structure (PostgreSQL + Apache AGE)
- Embedding generation targets for vector search
- RAG pipeline seed data

---

## Topic Hierarchy

```
Level 0: Category (e.g., "Data Structures")
  └─ Level 1: Topic (e.g., "Arrays")
      └─ Level 2: Subtopic (e.g., "Two Pointers")
          └─ Level 3: Pattern (e.g., "Opposite Direction")
```

---

## Topic Node Schema

Each topic in the database follows this structure:

```json
{
  "topic_id": "unique_identifier",
  "name": "Topic Name",
  "category": "category_name",
  "level": 0-3,
  "parent_topic_id": "parent_id or null",
  "description": "Detailed description",
  "difficulty_range": [min, max],
  "prerequisites": ["topic_id1", "topic_id2"],
  "related_topics": ["topic_id3"],
  "core_patterns": [
    {
      "pattern_name": "Name",
      "when_to_use": "Conditions",
      "template_code": "code_snippet",
      "time_complexity": "O(...)",
      "space_complexity": "O(...)"
    }
  ],
  "common_pitfalls": [
    {
      "pitfall": "Description",
      "why_common": "Why students make this",
      "how_to_avoid": "Fix",
      "example": "Code"
    }
  ],
  "edge_cases": [
    {
      "case": "Description",
      "example_input": "input",
      "common_mistake": "what goes wrong"
    }
  ],
  "practice_problems": [
    {"leetcode_number": 1, "title": "Two Sum", "difficulty": "Easy"}
  ]
}
```

---

## Complete Topic Coverage

### 1. Arrays & Strings

#### 1.1 Basic Arrays
- **Difficulty:** 0-30
- **Patterns:** Linear scan, index manipulation, in-place operations, prefix/suffix arrays
- **Pitfalls:** Off-by-one (`i <= n` vs `i < n`), empty array access, integer overflow in sum, modifying during iteration
- **Edge cases:** Empty, single element, all same, sorted, reverse sorted, negatives, INT_MIN/MAX
- **Problems:** Two Sum (#1), Best Time to Buy/Sell Stock (#121), Contains Duplicate (#217)

#### 1.2 Two Pointers
- **Difficulty:** 10-60
- **Prerequisites:** arrays_basic, loops
- **Patterns:**
  - **Opposite Direction:** Start at ends, move toward center. For pairs, palindromes, container problems.
    ```cpp
    int left = 0, right = n-1;
    while (left < right) { /* move based on condition */ }
    ```
  - **Same Direction (Slow/Fast):** Both forward, different speeds. Remove duplicates, move elements, partition.
    ```cpp
    int slow = 0;
    for (int fast = 0; fast < n; fast++) {
        if (condition) arr[slow++] = arr[fast];
    }
    ```
  - **Fast/Slow (Cycle Detection):** Fast moves 2 steps, slow 1. Linked list cycles, finding middle.
- **Pitfalls:** Off-by-one in pointer init, empty array, infinite loop when pointers don't progress
- **Problems:** Two Sum II (#167), 3Sum (#15), Container With Most Water (#11), Remove Duplicates (#26), Move Zeroes (#283)

#### 1.3 Sliding Window
- **Difficulty:** 20-70
- **Prerequisites:** two_pointers
- **Patterns:**
  - **Fixed Size:** Window of size k, slide across array
  - **Variable Size:** Expand right, shrink left when invalid
- **Pitfalls:** Not shrinking window correctly, wrong window validity check
- **Problems:** Max Sum Subarray of Size K, Longest Substring Without Repeating (#3), Minimum Window Substring (#76)

#### 1.4 Prefix Sum
- **Difficulty:** 15-50
- **Patterns:** Precompute cumulative sums for O(1) range queries
- **Pitfalls:** Off-by-one in range calculation, overflow in cumulative sum
- **Problems:** Range Sum Query (#303), Subarray Sum Equals K (#560)

#### 1.5 Strings
- **Difficulty:** 10-60
- **Patterns:** Character frequency, palindrome check, string matching
- **Pitfalls:** Unsigned size_t subtraction, string copy costs (hidden O(n)), encoding issues
- **Edge cases:** Empty string, single char, all same chars, spaces

### 2. Hash Tables

- **Difficulty:** 10-50
- **Patterns:** Frequency counting, complement lookup, grouping by key
- **When to use:** O(1) lookup needed, no ordering required
- **When NOT to use:** Need ordered iteration, range queries, min/max
- **Pitfalls:** Hash collisions (worst case O(n)), rehashing cost, unordered_map vs map choice
- **Problems:** Two Sum (#1), Group Anagrams (#49), Valid Anagram (#242)

### 3. Linked Lists

- **Difficulty:** 15-60
- **Patterns:** Fast/slow pointers, dummy head, reverse in-place
- **Pitfalls:** Null pointer dereference, losing reference to next node during reversal, not handling single node
- **Edge cases:** Empty list, single node, cycle
- **Problems:** Reverse Linked List (#206), Merge Two Sorted Lists (#21), Linked List Cycle (#141)

### 4. Stacks & Queues

- **Difficulty:** 15-65
- **Patterns:**
  - **Monotonic Stack:** Next greater/smaller element
  - **Stack for Matching:** Parentheses, nested structures
  - **Queue for BFS:** Level-order traversal
- **Pitfalls:** Empty stack pop, not handling unmatched elements
- **Problems:** Valid Parentheses (#20), Daily Temperatures (#739), Next Greater Element (#496)

### 5. Trees

#### 5.1 Binary Trees
- **Difficulty:** 20-60
- **Patterns:** DFS (preorder/inorder/postorder), BFS (level-order), recursive vs iterative
- **Key insight:** Trees have no cycles → no visited array needed
- **Problems:** Max Depth (#104), Invert Tree (#226), Level Order Traversal (#102)

#### 5.2 Binary Search Trees
- **Difficulty:** 25-65
- **Patterns:** Inorder = sorted, search/insert/delete O(log n) average
- **Pitfalls:** Unbalanced BST degrades to O(n)
- **Problems:** Validate BST (#98), Kth Smallest (#230)

#### 5.3 Tries
- **Difficulty:** 40-70
- **Patterns:** Prefix matching, autocomplete, word search
- **Problems:** Implement Trie (#208), Word Search II (#212)

### 6. Heaps / Priority Queues

- **Difficulty:** 30-70
- **Patterns:** Top-K elements, merge K sorted, median stream (two heaps)
- **Limitation:** No efficient decrease-key → add duplicates and skip processed
- **Problems:** Kth Largest (#215), Merge K Sorted Lists (#23), Find Median from Data Stream (#295)

### 7. Graphs

#### 7.1 Graph Traversal
- **Difficulty:** 30-70
- **Patterns:**
  - **DFS:** Recursive or stack-based. Good for: any path, cycle detection, connected components
  - **BFS:** Queue-based. Good for: shortest path (unweighted), level-order
- **Pitfalls:** Forgetting visited set, not handling disconnected components
- **Problems:** Number of Islands (#200), Clone Graph (#133)

#### 7.2 Shortest Path
- **Difficulty:** 40-80
- **Patterns:**
  - **BFS:** Unweighted graphs
  - **Dijkstra:** Non-negative weights, O((V+E) log V)
  - **Bellman-Ford:** Negative weights, O(VE)
- **Key observation:** Dijkstra fails with negative weights
- **Problems:** Network Delay Time (#743), Cheapest Flights (#787)

#### 7.3 Topological Sort
- **Difficulty:** 40-65
- **Patterns:** Kahn's (BFS with in-degree) or DFS postorder (reversed)
- **Prerequisite:** Must be DAG
- **Problems:** Course Schedule (#207), Course Schedule II (#210)

#### 7.4 Union-Find
- **Difficulty:** 35-65
- **Patterns:** Path compression + union by rank → near O(1) per operation
- **Problems:** Number of Connected Components (#323), Redundant Connection (#684)

### 8. Binary Search

- **Difficulty:** 20-70
- **Patterns:**
  - **Standard:** Find exact target
  - **Lower/Upper Bound:** Find first/last occurrence
  - **Search on Answer:** Binary search on result space
- **Pitfalls:** `(left + right) / 2` overflow → use `left + (right - left) / 2`, wrong loop condition (`<` vs `<=`), wrong mid update (`mid` vs `mid+1`)
- **Problems:** Binary Search (#704), Search in Rotated Array (#33), Koko Eating Bananas (#875)

### 9. Dynamic Programming

#### 9.1 1D DP
- **Difficulty:** 30-65
- **Patterns:** Fibonacci-style, house robber, climbing stairs
- **Space optimization:** If depends on last 2 states → O(1) space
- **Problems:** Climbing Stairs (#70), House Robber (#198), Coin Change (#322)

#### 9.2 2D DP
- **Difficulty:** 40-80
- **Patterns:** Grid DP, string DP (LCS, edit distance), knapsack
- **Space optimization:** If depends on previous row only → 1D array
- **Critical:** 0/1 knapsack iterate backwards, unbounded iterate forwards
- **Problems:** Unique Paths (#62), LCS (#1143), Edit Distance (#72), 0/1 Knapsack

#### 9.3 Advanced DP
- **Difficulty:** 60-95
- **Patterns:** Interval DP, bitmask DP, digit DP, DP on trees
- **Problems:** Burst Balloons (#312), Palindrome Partitioning (#131)

### 10. Greedy

- **Difficulty:** 25-75
- **When it works:** Optimal substructure + greedy choice property (provable via exchange argument)
- **When it fails:** Need to consider all options → use DP instead
- **Pitfalls:** Greedy "seems" to work but fails on edge cases → always prove or find counterexample
- **Problems:** Jump Game (#55), Interval Scheduling, Activity Selection

### 11. Backtracking

- **Difficulty:** 35-75
- **Patterns:** Generate all combinations/permutations, constraint satisfaction, pruning
- **Template:** Choose → Explore → Unchoose
- **Problems:** Subsets (#78), Permutations (#46), N-Queens (#51)

### 12. Bit Manipulation

- **Difficulty:** 20-60
- **Patterns:** XOR for finding unique, bit masking, counting bits
- **Problems:** Single Number (#136), Number of 1 Bits (#191)

---

## Graph Relationships for Database

For each topic, define these relationships in Apache AGE:

```cypher
(:Topic)-[:PREREQUISITE_FOR]->(:Topic)    -- Must learn A before B
(:Topic)-[:RELATED_TO]->(:Topic)           -- Connected concepts
(:Topic)-[:SUBTOPIC_OF]->(:Topic)          -- Hierarchy
(:Problem)-[:HAS_TOPIC]->(:Topic)          -- Problem uses topic
(:Problem)-[:USES_PATTERN]->(:Pattern)     -- Problem uses pattern
(:Problem)-[:SIMILAR_TO]->(:Problem)       -- Similar problems
(:Problem)-[:FOLLOW_UP_OF]->(:Problem)     -- Harder version
```

**Example prerequisite chain:**
```
Arrays → Two Pointers → Sliding Window
Arrays → Binary Search
Arrays → Hash Tables
Hash Tables → Graph Traversal (adjacency list)
Recursion → DFS → Backtracking
Recursion → Dynamic Programming
```

---

## Embedding Strategy

For each topic, generate embeddings from:

1. **Concept explanation** — "Arrays are contiguous memory blocks providing O(1) random access..."
2. **Common patterns** — "Common array patterns: linear scan, two pointers, sliding window..."
3. **Pitfalls** — "Common array pitfalls: off-by-one errors, empty array access, integer overflow..."

Store in vector DB with metadata (topic_id, type) for filtered retrieval.

---

## RAG Seed Data

For each topic, pre-generate Q&A pairs:

```json
{
  "question": "How do I handle an empty array?",
  "answer": "Always check if (arr.empty()) before accessing elements. For 'find maximum', return INT_MIN or throw exception."
}
```

And troubleshooting entries:

```json
{
  "symptom": "Segmentation fault when accessing array",
  "causes": ["Out of bounds", "Uninitialized", "Off-by-one"],
  "solutions": ["Check i < arr.size()", "Initialize before use", "Use arr.at(i) for bounds checking"]
}
```

---

## Seeding Workflow

1. Insert topics into PostgreSQL `topics` table (with parent_topic_id for hierarchy)
2. Create topic nodes in Apache AGE graph
3. Create PREREQUISITE_FOR, RELATED_TO, SUBTOPIC_OF edges
4. Insert problems with topic associations
5. Create problem nodes and HAS_TOPIC/USES_PATTERN edges
6. Generate embeddings for all topics and problems
7. Store embeddings in ChromaDB collections
8. Generate initial questions per topic using LLM + templates

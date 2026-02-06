# Algoholic — Question Design Guide
## Taxonomy, Pedagogy & Creative Question Strategies

---

## Philosophy

Test **understanding**, not memorization. If a user can answer these questions correctly, they truly understand the problem — not just memorized the code.

**Bad:** "What's the time complexity of Two Sum?"
**Good:** "Your friend says they can solve Two Sum in O(1) by sorting first. What's wrong with their reasoning?"

---

## Table of Contents

1. [Question Categories (Taxonomy)](#1-question-categories)
2. [Question Design Patterns](#2-question-design-patterns)
3. [DP-Specific Questions](#3-dp-specific-questions)
4. [Graph-Specific Questions](#4-graph-specific-questions)
5. [Debugging & Code Fixing](#5-debugging--code-fixing)
6. [Edge Cases & Overflow](#6-edge-cases--overflow)
7. [Observation Testing](#7-observation-testing)
8. [Conversion Questions](#8-conversion-questions)
9. [AI Hint Integration](#9-ai-hint-integration)
10. [Content Standards](#10-content-standards)

---

## 1. Question Categories

### Category 1: Complexity Analysis

| Subtype | Format | Tests |
|---------|--------|-------|
| 1.1 Code-to-Complexity | Show code → "What's the complexity?" | Reading code, identifying loops/recursion |
| 1.2 Complexity Comparison | Rank 3-4 solutions by efficiency | Relative understanding |
| 1.3 Constraint-to-Complexity | Given constraints → "What complexity is needed?" | Practical constraint interpretation |
| 1.4 Hidden Complexity | Code with STL/library calls → "What's the ACTUAL complexity?" | Knowing hidden costs (string ops, container ops) |
| 1.5 Amortized Analysis | Dynamic array, hash table rehashing | Understanding amortized vs worst-case |

**Difficulty progression:**
- Easy: Single loops, obvious nested loops
- Medium: Recursion with branching, divide & conquer
- Hard: Hidden complexity in library functions, amortized analysis
- Expert: Master theorem, complex recursion trees

### Category 2: Data Structure Selection

| Subtype | Format | Tests |
|---------|--------|-------|
| 2.1 Requirements-to-DS | List operations needed → "Which DS?" | Core DS knowledge |
| 2.2 STL Container Selection | Scenario → "Which STL container(s)?" | C++ practical knowledge |
| 2.3 DS Trade-off Analysis | Compare 2-3 DS options | Nuanced thinking |
| 2.4 Custom DS Design | Problem needing uncommon DS → "What combination?" | Creativity (LRU cache, median stream) |

### Category 3: Pattern Recognition

| Subtype | Format | Tests |
|---------|--------|-------|
| 3.1 Problem-to-Pattern | Problem statement → "Which pattern?" | Core pattern recognition |
| 3.2 Pattern Variant | Known problem + new problem → "Same pattern?" | Transfer learning |
| 3.3 Anti-Pattern Detection | Wrong approach → "Why won't this work?" | Deep understanding |
| 3.4 Multiple Patterns | Problem with multiple valid approaches → compare | Breadth of knowledge |

**Pattern library:** Two Pointers, Binary Search, DFS/BFS, DP, Greedy, Backtracking, Union-Find, Topological Sort, Monotonic Stack/Queue, Trie, Sliding Window.

### Category 4: STL Operations & Complexity

| Subtype | Tests |
|---------|-------|
| 4.1 STL Function Complexity | Know actual complexity of container/algorithm operations |
| 4.2 STL Algorithm Selection | Pick right algorithm (sort vs partial_sort, find vs binary_search) |
| 4.3 STL Gotchas | Iterator invalidation, erase-remove idiom, unsigned size_t |
| 4.4 STL vs Manual | When STL doesn't fit (heap with decrease-key, custom comparators) |

### Category 5: Code Template Mastery

| Subtype | Tests |
|---------|-------|
| 5.1 Template Recognition | Problem → "Which template?" |
| 5.2 Template Customization | Standard template + modification needed |
| 5.3 Template from Memory | Write template without reference |
| 5.4 Speed Drill | 10-second rapid-fire: just name the template |

### Category 6: Implementation Correctness

| Subtype | Tests |
|---------|-------|
| 6.1 Spot Correct Implementation | 3-4 implementations, only one correct |
| 6.2 Edge Case Coverage | "What edge case is missing?" |
| 6.3 Boundary Conditions | `left < right` vs `left <= right`, `i < n` vs `i <= n` |

### Category 7: Bug Detection & Fixing

| Subtype | Tests |
|---------|-------|
| 7.1 Find the Bug | Working solution with 1-2 bugs |
| 7.2 Bug Fixing | Bug shown, multiple fix options |
| 7.3 Output Prediction | Buggy code + input → "What's the output?" |
| 7.4 Multiple Bugs | Code with 3-5 bugs, find all |

### Category 8: Pseudocode Algorithm Design

| Subtype | Tests |
|---------|-------|
| 8.1 Problem to Pseudocode | Write high-level algorithm |
| 8.2 Pseudocode to Complexity | Analyze pseudocode complexity |
| 8.3 Pseudocode Optimization | Improve naive pseudocode |
| 8.4 Pseudocode Verification | Find logic error in pseudocode |

### Category 9: Approach Trade-offs

| Subtype | Tests |
|---------|-------|
| 9.1 Iterative vs Recursive DP | When to use which |
| 9.2 DFS vs BFS | Shortest path, cycle detection, space |
| 9.3 Hash Table vs BST | Ordered vs unordered, range queries |
| 9.4 Greedy vs DP | Prove greedy works or find counterexample |
| 9.5 Space vs Time | When to trade memory for speed |
| 9.6 Sorting Algorithm Selection | HeapSort vs MergeSort vs QuickSort by scenario |

### Category 10: Hybrid Multi-Skill Challenges

| Subtype | Tests |
|---------|-------|
| 10.1 Full Problem Analysis | Pattern → complexity → DS → pseudocode → trade-offs |
| 10.2 Design + Implementation | Design DS, justify choices, implement key methods |
| 10.3 Optimization Challenge | O(n³) → O(n²) → O(n log n) → O(n) progressive |

---

## 2. Question Design Patterns

### Return Value Semantics

Test what to return in edge/impossible cases:

```
Minimization (Coin Change): impossible → INT_MAX or amount+1 (never chosen as min)
Maximization (Knapsack):    impossible → INT_MIN or -1
Boolean (Partition):        impossible → false
Counting:                   impossible → 0
Base case (capacity=0):     → 0 (not -1, not INT_MAX)
```

**Key pattern:** Minimization uses ∞ for impossible, maximization uses -∞.

### State Definition

"What does dp[i][w] represent?"
- Test if user understands the state, not just the recurrence.
- Common wrong answers reveal whether they confuse "exactly w" vs "at most w", or "item i" vs "first i items".

### Initialization Values

"Why initialize dp[1..n] = INT_MAX, not 0?"
- If 0, we'd never update (0 is always min).
- If -1, can't do min(-1, anything + 1).

### Space Optimization

"Can this 2D DP become 1D?"
- Key: does current row only depend on previous row?
- Critical: iterate backwards for 0/1 knapsack (forward = unbounded knapsack).

---

## 3. DP-Specific Questions

### Base Case Return Values

```
Knapsack, capacity=0:     → 0 (can't take anything)
Rod Cutting, length=0:    → 0 (no rod, no value)
LCS, empty string:        → 0 (LCS of empty with anything is 0)
LIS, all decreasing:      → 1 (single element is valid subsequence)
Max Subarray, all negative: → least negative (must pick at least one)
```

### DP Optimization Chain

```
Fibonacci:  O(2^n) → memoization O(n)/O(n) → tabulation O(n)/O(n) → space-optimized O(n)/O(1)
Knapsack:   2D table O(nW)/O(nW) → 1D backwards O(nW)/O(W)
Unique Paths: 2D O(mn)/O(mn) → 1D O(mn)/O(n) → O(mn)/O(min(m,n))
```

---

## 4. Graph-Specific Questions

### Key Observations Before Coding

Test if user notices critical properties:

| Problem | Key Observation |
|---------|----------------|
| Bipartite check | 2-colorable ↔ no odd-length cycles |
| Topological sort | Must be DAG (no cycles) |
| Shortest path, negative weights | Dijkstra won't work → use Bellman-Ford |
| Tree problems | Exactly one path between nodes → no visited array needed |
| Number of Islands | Can modify input grid to mark visited |

### BFS vs DFS Decision

| Need | Use |
|------|-----|
| Shortest path (unweighted) | BFS |
| Any path in deep graph | DFS |
| Level-order traversal | BFS |
| Cycle detection (directed) | DFS with recursion stack |
| Simple connectivity | Either |

### Observation Chains

Test multi-step reasoning:
1. Prerequisites form a directed graph
2. Cycle = impossible to complete
3. Need topological order
4. DFS postorder gives **reverse** topological order → must reverse

---

## 5. Debugging & Code Fixing

### Common Bug Types to Plant

| Bug | Example |
|-----|---------|
| Off-by-one | `right = arr.size()` instead of `arr.size() - 1` |
| Integer overflow | `(left + right) / 2` when both near INT_MAX |
| Infinite loop | `left = mid` instead of `left = mid + 1` |
| Missing base case | No check for empty array before `arr[0]` |
| Wrong loop condition | `left < right` vs `left <= right` |
| Unsigned subtraction | `arr.size() - 1` when size is 0 (wraps to MAX_UINT) |

### Find-the-Breaking-Test-Case Format

Show correct-looking code, ask which input breaks it. Forces mental execution and edge case thinking.

### TLE → Optimization

Show O(n²) code, ask to optimize to O(n). Common patterns:
- Nested loops → hash table
- Recalculating sums → prefix sums
- Naive recursion → memoization

---

## 6. Edge Cases & Overflow

### Always-Test Checklist

**Arrays:** empty, single element, two elements, all same, sorted, reverse sorted, duplicates, negatives, INT_MIN/MAX values.

**Strings:** empty, single char, all same chars, palindrome, spaces, special chars.

**Trees:** null root, single node, skewed (all left/right), balanced, duplicate values.

**Graphs:** disconnected, self-loops, single node, complete graph.

### Overflow Detection

```
n=10^5 elements, each ≤ 10^9 → sum can reach 10^14 → need long long
a,b ≤ 10^5 → a*b can reach 10^10 → need long long
13! overflows int, 21! overflows long long
Fix: (long long)a * b, or 1LL * result * i
Modular: result = (1LL * result * i) % MOD
```

---

## 7. Observation Testing

### Format

Before showing solution, test key insights as separate questions:

**Two Sum:**
1. "Can we have duplicate elements?" → Yes → careful with indices
2. "If using hash table, store element→index or index→element?" → element→index
3. "Can we solve without extra space?" → Sort + two pointers, but loses original indices

**Longest Substring Without Repeating:**
1. "What causes window to shrink?" → Duplicate character
2. "How to check for duplicates?" → Hash set or frequency array, O(1) lookup
3. "When duplicate found, where does left pointer go?" → Last position of duplicate + 1

### Why This Works

Forces understanding of the "why" before the "how". If user misses a key observation, they can't have truly understood the solution.

---

## 8. Conversion Questions

### Iterative ↔ Recursive DP

Show one form, ask to convert. Key things to test:
- Base case mapping
- Memoization key
- Loop direction (forward vs backward)
- Space optimization awareness

### 2D → 1D DP

Critical question: "Do we iterate forwards or backwards?"
- 0/1 Knapsack: backwards (to use previous row values)
- Unbounded Knapsack: forwards (reuse current row values)

### DFS ↔ BFS

Level-order traversal: BFS is natural (queue), DFS needs depth parameter.

---

## 9. AI Hint Integration

### Progressive Hint Levels

| Level | Style | Example (Binary Search Overflow) |
|-------|-------|----------------------------------|
| 1 | Socratic question | "What are the possible values of left and right?" |
| 2 | Point to area | "Look at `(left + right) / 2`. What happens with large values?" |
| 3 | Explain concept | "Integer overflow occurs when left + right exceeds INT_MAX. Rewrite without adding first." |

### Trigger Conditions

- 2+ incorrect attempts on same concept → Level 1
- 5+ minutes on one question → Level 2
- 3+ incorrect + > 5 minutes → Level 3
- Category success rate < 40% → redirect to foundation questions

### Context for AI

Send: question details, user answer, attempt number, time spent, user's strengths/weaknesses, success rate in category. AI responds with hint at appropriate level, never giving the answer directly.

### Offline Fallback

Pre-generate 3 hint levels per question for users without AI access.

---

## 10. Content Standards

### Per Question Requirements

- Clear learning objective
- 4 plausible options (for MC) — no obviously wrong answers
- Detailed explanation for correct AND incorrect options
- Related problems linked
- Difficulty calibrated to 0-100 scale
- Tags assigned (patterns, concepts)

### Difficulty Calibration

| Level | Concepts | Approach | Edge Cases | Time |
|-------|----------|----------|------------|------|
| Easy | Single | Clear, standard template | 1-2 | 5-10 min |
| Medium | 2 combined | Some ambiguity, template adaptation | 3-5 | 15-25 min |
| Hard | 3+ or novel | Non-obvious, custom solution | Many | 30-45 min |
| Expert | Research-level | Multiple steps, optimization critical | All complex | 45+ min |

### Content Targets

| Category | Questions | % of Total |
|----------|-----------|------------|
| Complexity Analysis | 150 | 15% |
| DS Selection | 120 | 12% |
| Pattern Recognition | 200 | 20% |
| STL Operations | 100 | 10% |
| Code Templates | 80 | 8% |
| Implementation | 100 | 10% |
| Bug Detection | 100 | 10% |
| Pseudocode | 80 | 8% |
| Trade-offs | 70 | 7% |
| **Total** | **~1000** | |

### Special Formats

- **Rapid Fire:** 10 questions in 5 minutes, pattern/complexity focus
- **Deep Dive:** One problem, 10 parts, each testing different skill
- **Interview Simulation:** Timed 45 min, full problem + follow-ups
- **Debug Challenge:** Find all bugs in given code
- **Optimization Race:** O(n³) → O(n²) → O(n log n) → O(n) progressive

### Quality Checklist

- [ ] Correct answer verified
- [ ] All distractors plausible
- [ ] Explanation written for each option
- [ ] Related questions linked
- [ ] Difficulty calibrated
- [ ] Tags assigned
- [ ] Edge cases covered in test cases
- [ ] Grammar checked

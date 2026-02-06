# Question Metadata Pattern

Every question includes metadata arrays for organization and learning.

## Fields

```go
Tags                []string  // ["complexity", "nested-loops", "triangular-sum"]
Concepts            []string  // ["time-complexity", "summation", "nested-iteration"]
CommonMistakes      []string  // ["Thinking j=i makes it O(n)", "Confusing n(n+1)/2 with O(n)"]
RelatedProblemSlug  string    // "two-sum"
RelatedTopicID      string    // "arrays_basics"
```

## Rules

- **Tags:** Broad categories for filtering (e.g., "complexity", "dp", "graph")
- **Concepts:** Specific learning topics (e.g., "amortized-analysis", "memoization")
- **CommonMistakes:** Typical wrong approaches or misconceptions
- **RelatedProblemSlug:** Links to the problem this question tests
- **RelatedTopicID:** Links to the topic/pattern this question covers

## Example

```go
{
    Tags: ["complexity", "recursion", "exponential"],
    Concepts: ["time-complexity", "recursion-tree", "branching-factor"],
    CommonMistakes: [
        "Thinking it's O(n) because each call does O(1) work",
        "Confusing with binary search O(log n)",
    ],
    RelatedProblemSlug: "fibonacci-number",
    RelatedTopicID: "recursion_analysis",
}
```

## Why

- Tags enable broad filtering ("show me all DP questions")
- Concepts track specific learning progress ("user weak in memoization")
- CommonMistakes help explain wrong answers and detect patterns
- Related links connect questions to problems and learning paths

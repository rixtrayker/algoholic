# Three-Level Hint System

Every question provides progressive hints from Socratic to direct.

## Pattern

```go
HintLevel1: "How many total iterations does the inner loop execute across all values of i?"
HintLevel2: "Sum the series: n + (n-1) + (n-2) + ... + 1"
HintLevel3: "That's the triangular number n(n+1)/2 which simplifies to O(n²)"
```

## Rules

- **Level 1 (Socratic):** Ask a guiding question, don't reveal the answer
- **Level 2 (Directional):** Point to the key concept or area to examine
- **Level 3 (Explanatory):** Explain the solution approach directly
- Never give the final answer in any hint
- Each level should build on the previous one

## Progressive Disclosure

```
Level 1 → Activates thinking: "What property matters here?"
Level 2 → Narrows focus: "Consider the recurrence relation"
Level 3 → Shows method: "Use master theorem with a=2, b=2"
```

## Why

- Preserves learning: students think before seeing solution
- Progressive disclosure prevents accidentally giving away too much
- Socratic approach builds problem-solving skills

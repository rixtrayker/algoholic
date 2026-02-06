# Answer Option Structure

All multiple-choice questions use structured answer options with explanations for both correct and wrong answers.

## Pattern

```go
AnswerOptions: []AnswerOption{
    {ID: "A", Text: "O(n)", IsCorrect: false},
    {ID: "B", Text: "O(n²)", IsCorrect: true},
    {ID: "C", Text: "O(n log n)", IsCorrect: false},
    {ID: "D", Text: "O(n³)", IsCorrect: false},
}

WrongAnswerExplanations: JSONB{
    "A": "The inner loop still runs n*(n+1)/2 total iterations, not n.",
    "C": "n log n would require the inner loop to shrink logarithmically.",
    "D": "There are only two nested loops, not three.",
}
```

## Rules

- Use IDs "A", "B", "C", "D" for 4-option questions
- All options must be plausible — no obviously wrong answers
- Distractors should reflect common misconceptions
- Provide `WrongAnswerExplanations` for every wrong option
- Explanation field covers the correct answer reasoning

## Complete Example

```go
{
    QuestionText: "What is the time complexity?\n\n```go\nfor i := 0; i < n; i++ {\n    for j := i; j < n; j++ {\n        // O(1) work\n    }\n}\n```",
    AnswerOptions: []AnswerOption{
        {ID: "A", Text: "O(n)", IsCorrect: false},
        {ID: "B", Text: "O(n²)", IsCorrect: true},
        {ID: "C", Text: "O(n log n)", IsCorrect: false},
        {ID: "D", Text: "O(n³)", IsCorrect: false},
    },
    WrongAnswerExplanations: JSONB{
        "A": "The inner loop still runs n*(n+1)/2 total iterations, not n.",
        "C": "n log n would require the inner loop to shrink logarithmically, but j goes from i to n.",
        "D": "There are only two nested loops, not three.",
    },
    Explanation: "The inner loop runs n, n-1, n-2, ..., 1 times. Total = n(n+1)/2 = O(n²).",
}
```

## Anti-Patterns

❌ **Bad: Obvious distractors**
```go
AnswerOptions: []AnswerOption{
    {ID: "A", Text: "O(1)", IsCorrect: false},      // Obviously wrong
    {ID: "B", Text: "O(n²)", IsCorrect: true},
    {ID: "C", Text: "O(n!)", IsCorrect: false},     // Absurd
    {ID: "D", Text: "O(2ⁿ)", IsCorrect: false},    // Unrealistic
}
```

❌ **Bad: Missing wrong answer explanations**
```go
WrongAnswerExplanations: JSONB{
    "A": "Wrong",  // Not helpful
}
// Missing B and D explanations
```

❌ **Bad: No misconception testing**
```go
// All wrong answers are random, don't reflect actual student mistakes
```

## Why

- Wrong answer explanations correct misconceptions at point of failure
- Plausible distractors test understanding, not memorization
- Students learn why approaches fail, not just which one works

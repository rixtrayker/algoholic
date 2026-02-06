# Difficulty Scoring

Questions use a 0-100 numeric score plus a categorical label.

## Fields

```go
DifficultyScore float64         // 0-100 numeric scale
DifficultyLabel DifficultyLabel  // "easy" | "medium" | "hard" | "expert"
```

## Score Guidelines

```
 0-20   Easy       Single concept, standard template
21-40   Medium     2-3 concepts, some ambiguity
41-70   Hard       Multiple concepts, non-obvious approach
71-100  Expert     Novel techniques, research-level
```

## Rules

- Score and label are independent fields
- Score reflects objective difficulty (concepts, steps, edge cases)
- Label is manually assigned based on overall judgment
- Questions with similar scores may have different labels
- Score can be recalibrated based on user performance data

## Example

```go
DifficultyScore: 35  // Numeric: low-medium range
DifficultyLabel: "medium"  // Categorical: judged as medium overall
```

## Why

- Numeric score enables fine-grained filtering and progression
- Categorical label provides intuitive difficulty tiers
- Independence allows calibration without changing labels

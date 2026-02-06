# Enum Type Pattern

Use string-based type aliases with const values for enums.

## Pattern

```go
type DifficultyLabel string

const (
    DiffEasy   DifficultyLabel = "easy"
    DiffMedium DifficultyLabel = "medium"
    DiffHard   DifficultyLabel = "hard"
)
```

## Rules

- Declare enum as `type Name string`
- Define constants with type prefix (e.g., `DiffEasy` not just `Easy`)
- Use lowercase string values matching PostgreSQL enum
- Cast to string when inserting: `string(difficulty)`
- Match DB enum with `::enum_name` in SQL: `$1::difficulty_label`

## Example Usage

```go
type Problem struct {
    DifficultyLabel DifficultyLabel `json:"difficulty_label" db:"difficulty_label"`
}

// Insert
db.Exec(`INSERT INTO problems (..., difficulty_label) VALUES (..., $1::difficulty_label)`,
    string(p.DifficultyLabel))
```

## Why

- Readable in DB queries and JSON: "easy" not `1`
- Type safety in Go prevents invalid values at compile time
- Avoids ALTER TYPE migrations when adding values

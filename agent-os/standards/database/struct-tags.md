# Struct Tag Convention

All model structs use dual `json` and `db` tags for serialization.

## Pattern

```go
type Topic struct {
    TopicID   string `json:"topic_id" db:"topic_id"`
    Name      string `json:"name" db:"name"`
    Slug      string `json:"slug" db:"slug"`

    // Transient fields: json only, no db tag
    TopicLinks []ProblemTopicLink `json:"topic_links,omitempty"`
}
```

## Rules

- Every database field gets both `json` and `db` tags
- Tag values match: use same name for json and db (typically snake_case)
- Transient/computed fields: json tag only, omit db tag
- Optional fields: use `omitempty` in json tag
- Use pointer types for nullable DB columns

## Why

- `json` tags: API serialization to/from JSON
- `db` tags: SQL row scanning with sqlx
- Dual tags enable single struct for both database and API layer

# Struct Tag Convention

All model structs use dual `json` and `gorm` tags for serialization and ORM mapping.

## Pattern

```go
type Topic struct {
    TopicID   string `json:"topic_id" gorm:"column:topic_id;primaryKey"`
    Name      string `json:"name" gorm:"column:name;not null"`
    Slug      string `json:"slug" gorm:"column:slug;uniqueIndex;not null"`

    // Transient fields: json only, gorm:"-" to ignore
    TopicLinks []ProblemTopicLink `json:"topic_links,omitempty" gorm:"-"`
}

// Specify table name explicitly
func (Topic) TableName() string {
    return "topics"
}
```

## Rules

- Every database field gets both `json` and `gorm` tags
- `json` tag: snake_case matching API contract
- `gorm` tag: `column:` directive with snake_case column name
- Transient/computed fields: json tag only, use `gorm:"-"` to exclude from DB
- Optional fields: use `omitempty` in json tag, pointer types for nullable columns
- Primary keys: add `primaryKey` in gorm tag
- Constraints: use gorm tag modifiers (`not null`, `uniqueIndex`, `index`, etc.)
- Always implement `TableName()` method to avoid pluralization issues

## Common GORM Tag Modifiers

- `primaryKey` - Mark as primary key
- `column:name` - Specify column name
- `not null` - NOT NULL constraint
- `unique` - UNIQUE constraint
- `uniqueIndex` or `index` - Create index
- `default:value` - Set default value
- `type:varchar(100)` - Override column type
- `-` - Ignore this field

## Why

- `json` tags: API serialization to/from JSON
- `gorm` tags: GORM ORM mapping and constraints
- Dual tags enable single struct for database, API, and business logic
- `TableName()` prevents GORM from auto-pluralizing table names

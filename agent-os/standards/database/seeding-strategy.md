# Database Seeding Strategy

Use two-pass insertion for entities with self-referential foreign keys.

## Pattern

```go
// Pass 1: Insert without parent references
for _, t := range topics {
    db.Exec(`INSERT INTO topics (...) VALUES (...) ON CONFLICT (topic_id) DO UPDATE SET name=EXCLUDED.name`)
}

// Pass 2: Update parent references
for _, t := range topics {
    if t.ParentTopicID != nil {
        db.Exec(`UPDATE topics SET parent_topic_id = $1 WHERE topic_id = $2`)
    }
}
```

## Rules

- First pass creates all entities with nullable FKs set to NULL
- Second pass updates FK references after all entities exist
- Use `ON CONFLICT` for idempotent re-runs
- Cache slug→ID mappings to avoid repeated lookups
- Seed in dependency order: categories → topics → problems → questions

## Why

Solves circular dependencies: Topic A can reference Topic B as parent, while B hasn't been inserted yet.

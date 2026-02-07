# Database Seeding Strategy

Use two-pass insertion for entities with self-referential foreign keys with GORM and golang-migrate.

## Pattern

### Option 1: GORM Two-Pass (Development/Testing)

```go
// Pass 1: Create all entities with nil parent references
for _, t := range topics {
    topicCopy := t
    topicCopy.ParentTopicID = nil

    result := db.Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "topic_id"}},
        DoUpdates: clause.AssignmentColumns([]string{"name", "slug"}),
    }).Create(&topicCopy)

    if result.Error != nil {
        log.Printf("Error inserting topic %s: %v", t.TopicID, result.Error)
    }
}

// Pass 2: Update parent references
for _, t := range topics {
    if t.ParentTopicID != nil {
        db.Model(&Topic{}).
            Where("topic_id = ?", t.TopicID).
            Update("parent_topic_id", t.ParentTopicID)
    }
}
```

### Option 2: golang-migrate SQL (Production)

Create migration file `migrations/000002_seed_topics.up.sql`:

```sql
-- Pass 1: Insert without parent references
INSERT INTO topics (topic_id, name, slug) VALUES
    ('arrays_basics', 'Array Basics', 'array-basics'),
    ('dp_intro', 'Dynamic Programming Intro', 'dp-intro')
ON CONFLICT (topic_id) DO UPDATE SET
    name = EXCLUDED.name,
    slug = EXCLUDED.slug;

-- Pass 2: Update parent references
UPDATE topics SET parent_topic_id = 'arrays_basics' WHERE topic_id = 'arrays_two_pointer';
UPDATE topics SET parent_topic_id = 'dp_intro' WHERE topic_id = 'dp_knapsack';
```

## Rules

- **Development/Testing:** Use GORM with `OnConflict` clause for idempotent inserts
- **Production/Schema changes:** Use golang-migrate SQL files for version control
- First pass creates all entities with nullable FKs set to NULL
- Second pass updates FK references after all entities exist
- Use `ON CONFLICT` (SQL) or `Clauses(clause.OnConflict{...})` (GORM) for idempotence
- Cache slug→ID mappings to avoid repeated lookups in application code
- Seed in dependency order: categories → topics → problems → questions
- For GORM AutoMigrate: only use in development, never in production

## Hybrid Approach (Recommended)

1. **Schema migrations:** Always use golang-migrate for CREATE/ALTER/DROP
2. **Initial data seeding:** Use golang-migrate SQL files for production
3. **Test data seeding:** Use GORM two-pass in test fixtures
4. **Development convenience:** GORM AutoMigrate in local env only

## Why

- Solves circular dependencies: Topic A can reference Topic B as parent, while B hasn't been inserted yet
- golang-migrate provides version control and rollback for schema changes
- GORM provides type safety and convenience for application-level seeding
- Hybrid approach gives best of both worlds

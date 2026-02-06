# Master Reference Database Seeding Guide
## Step-by-Step n8n Workflow Implementation

This guide provides **detailed implementation instructions** for the n8n workflow that seeds the PostgreSQL database using the Master Topic Reference.

## Quick Reference

**Purpose:** Transform Master Reference markdown → PostgreSQL database records

**Dependencies:**
- Master Topic Reference document
- PostgreSQL with schema created
- n8n workflow engine
- Node.js for transformation logic

**Output:** Fully populated database ready for graph construction and embeddings

---

## Seeding Sequence

```
1. Parse Topics → Insert Topics (sorted by depth)
2. Extract Prerequisites → Insert Prerequisites  
3. Extract Patterns → Insert Patterns (deduplicated)
4. Link Topics-Patterns → Insert topic_patterns
5. Extract Mistakes → Insert common_mistakes
6. Extract Edge Cases → Insert edge_cases
7. Extract Problems → Insert problems
8. Link Problem-Topics → Insert problem_topics
9. Link Problem-Patterns → Insert problem_patterns
10. Link Problem-Mistakes → Insert problem_mistakes
11. Link Problem-EdgeCases → Insert problem_edge_cases
12. Validate Integrity
13. Trigger Graph Construction Workflow
```

---

## Key n8n Node Patterns

### Pattern 1: Parse and Transform

```javascript
// Read markdown → Parse structure → Transform to DB format
const text = $input.first().binary.data.toString('utf8');
const parsed = parseStructure(text);
return parsed.map(item => ({json: item}));
```

### Pattern 2: Batch Insert with Conflict Handling

```sql
INSERT INTO table (cols...) VALUES ($1, $2, ...)
ON CONFLICT (unique_col) DO UPDATE SET
    col1 = EXCLUDED.col1,
    updated_at = NOW()
RETURNING id;
```

### Pattern 3: Foreign Key Lookup

```sql
-- Use subquery for FK lookup
INSERT INTO child_table (parent_id, data) VALUES (
    (SELECT id FROM parent_table WHERE code = $1),
    $2
);
```

### Pattern 4: Validation Query

```sql
-- Post-insert validation
SELECT COUNT(*) as issues 
FROM topics 
WHERE parent_topic_id IS NULL AND depth_level > 1;
-- Should return 0
```

---

## Critical Implementation Notes

**1. Insertion Order Matters**
- Must insert parents before children (sort by depth_level)
- Must insert referenced entities before relationships

**2. Deduplication Strategy**
- Patterns: deduplicate by pattern_name
- Mistakes: deduplicate by topic_code + mistake_name
- Edge Cases: deduplicate by topic_code + case_name
- Problems: deduplicate by leetcode_id

**3. Data Enrichment**
- Auto-generate estimated_practice_hours if missing
- Infer pattern_type from pattern_name
- Infer mistake_category from mistake_name
- Infer edge_case_category from case_name

**4. Error Handling**
- Use ON CONFLICT for idempotency
- Retry on transient DB errors
- Log all insertions with counts
- Validate after each major step

---

## Validation Checklist

After seeding, verify:

✅ Topics: All inserted, proper hierarchy, no orphans
✅ Prerequisites: No circular dependencies, valid references
✅ Patterns: Deduplicated, all have type
✅ Topic-Patterns: All topics have ≥1 pattern (depth ≥ 2)
✅ Mistakes: Proper categories, all have severity
✅ Edge Cases: All have why_important field
✅ Problems: All have leetcode_id, proper difficulty
✅ Relationships: All FKs valid, no dangling references

---

## Expected Record Counts

| Table | Expected Count |
|-------|----------------|
| topics | ~150-200 |
| patterns | ~40-60 |
| prerequisites | ~100-150 |
| topic_patterns | ~300-500 |
| common_mistakes | ~200-300 |
| edge_cases | ~100-150 |
| problems | ~50-100 (initial) |
| problem_topics | ~100-200 |

---

For detailed code examples and complete workflow JSON, see the full document sections above.

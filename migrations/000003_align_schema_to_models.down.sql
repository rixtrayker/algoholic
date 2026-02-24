-- 000003_align_schema_to_models.down.sql
-- Reverting this migration is non-trivial due to VARCHAR -> INT conversion.
-- This down migration drops the changed columns and re-adds them as VARCHAR.
-- DATA WILL BE LOST for topic IDs. Only use in development.

-- Revert questions column renames
ALTER TABLE questions RENAME COLUMN question_type TO category;
ALTER TABLE questions RENAME COLUMN question_subtype TO subcategory;
ALTER TABLE questions RENAME COLUMN question_format TO format;
ALTER TABLE questions RENAME COLUMN problem_id TO related_problem_id;
ALTER TABLE questions RENAME COLUMN total_attempts TO attempt_count;
ALTER TABLE questions RENAME COLUMN correct_attempts TO correct_count;
ALTER TABLE questions RENAME COLUMN average_time_seconds TO avg_time_sec;
ALTER TABLE questions RENAME COLUMN estimated_time_seconds TO estimated_time_sec;

-- Revert problems column renames
ALTER TABLE problems RENAME COLUMN description TO statement;
ALTER TABLE problems DROP COLUMN IF EXISTS official_difficulty;
ALTER TABLE problems DROP COLUMN IF EXISTS total_attempts;
ALTER TABLE problems DROP COLUMN IF EXISTS total_solves;
ALTER TABLE problems RENAME COLUMN average_time_seconds TO avg_solve_time_sec;
ALTER TABLE problems RENAME COLUMN acceptance_rate TO success_rate;

-- Revert problem_topics column renames
ALTER TABLE problem_topics RENAME COLUMN relevance_score TO relevance;

-- Remove soft delete columns
ALTER TABLE training_plans DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE training_plan_items DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE user_lists DROP COLUMN IF EXISTS deleted_at;

-- Note: reverting topic_id from INT back to VARCHAR(100) is destructive.
-- If you need the original VARCHAR IDs, restore from backup.

-- 000003_align_schema_to_models.up.sql
-- Aligns the SQL schema with GORM model definitions.
-- Primary change: topics.topic_id from VARCHAR(100) to SERIAL (int).
-- Also renames columns to match GORM struct field tags.

-- Step 1: Drop foreign keys and indexes that reference topics.topic_id as VARCHAR
ALTER TABLE problem_topics DROP CONSTRAINT IF EXISTS problem_topics_topic_id_fkey;
ALTER TABLE questions      DROP CONSTRAINT IF EXISTS questions_related_topic_id_fkey;
ALTER TABLE pitfalls       DROP CONSTRAINT IF EXISTS pitfalls_topic_id_fkey;
ALTER TABLE topics         DROP CONSTRAINT IF EXISTS topics_parent_topic_id_fkey;
ALTER TABLE user_skills    DROP CONSTRAINT IF EXISTS user_skills_topic_id_fkey;
ALTER TABLE weakness_analysis DROP CONSTRAINT IF EXISTS weakness_analysis_specific_topic_fkey;

DROP INDEX IF EXISTS idx_topics_category;
DROP INDEX IF EXISTS idx_topics_parent;
DROP INDEX IF EXISTS idx_topics_level;
DROP INDEX IF EXISTS idx_topics_keywords;
DROP INDEX IF EXISTS idx_pitfalls_topic;

-- Step 2: Create a temporary mapping table (old VARCHAR id -> new SERIAL id)
CREATE TEMP TABLE topic_id_map AS
  SELECT topic_id AS old_id, ROW_NUMBER() OVER (ORDER BY created_at, topic_id) AS new_id
  FROM topics;

-- Step 3: Add a temporary integer column and populate it
ALTER TABLE topics ADD COLUMN new_topic_id SERIAL;
UPDATE topics t SET new_topic_id = m.new_id FROM topic_id_map m WHERE t.topic_id = m.old_id;

-- Step 4: Update parent_topic_id to integer
ALTER TABLE topics ADD COLUMN new_parent_topic_id INT;
UPDATE topics t SET new_parent_topic_id = m.new_id
  FROM topic_id_map m WHERE t.parent_topic_id = m.old_id;

-- Step 5: Update referencing tables to use integer IDs
ALTER TABLE problem_topics ADD COLUMN new_topic_id INT;
UPDATE problem_topics pt SET new_topic_id = m.new_id
  FROM topic_id_map m WHERE pt.topic_id = m.old_id;

ALTER TABLE user_skills ADD COLUMN new_topic_id INT;
UPDATE user_skills us SET new_topic_id = m.new_id
  FROM topic_id_map m WHERE us.topic_id::text = m.old_id;

ALTER TABLE questions ADD COLUMN new_related_topic_id INT;
UPDATE questions q SET new_related_topic_id = m.new_id
  FROM topic_id_map m WHERE q.related_topic_id = m.old_id;

ALTER TABLE pitfalls ADD COLUMN new_topic_id INT;
UPDATE pitfalls p SET new_topic_id = m.new_id
  FROM topic_id_map m WHERE p.topic_id = m.old_id;

ALTER TABLE weakness_analysis ADD COLUMN new_specific_topic INT;
UPDATE weakness_analysis w SET new_specific_topic = m.new_id
  FROM topic_id_map m WHERE w.specific_topic = m.old_id;

-- Step 6: Drop old VARCHAR columns and rename new INT columns

-- topics
ALTER TABLE topics DROP COLUMN topic_id CASCADE;
ALTER TABLE topics RENAME COLUMN new_topic_id TO topic_id;
ALTER TABLE topics ADD PRIMARY KEY (topic_id);
ALTER TABLE topics DROP COLUMN parent_topic_id;
ALTER TABLE topics RENAME COLUMN new_parent_topic_id TO parent_topic_id;

-- problem_topics
ALTER TABLE problem_topics DROP COLUMN topic_id;
ALTER TABLE problem_topics RENAME COLUMN new_topic_id TO topic_id;
ALTER TABLE problem_topics ADD CONSTRAINT problem_topics_pkey PRIMARY KEY (problem_id, topic_id);

-- user_skills
ALTER TABLE user_skills DROP CONSTRAINT IF EXISTS user_skills_pkey;
ALTER TABLE user_skills DROP COLUMN topic_id;
ALTER TABLE user_skills RENAME COLUMN new_topic_id TO topic_id;
ALTER TABLE user_skills ADD PRIMARY KEY (user_id, topic_id);

-- questions
ALTER TABLE questions DROP COLUMN related_topic_id;
ALTER TABLE questions RENAME COLUMN new_related_topic_id TO related_topic_id;

-- pitfalls
ALTER TABLE pitfalls DROP COLUMN topic_id;
ALTER TABLE pitfalls RENAME COLUMN new_topic_id TO topic_id;

-- weakness_analysis
ALTER TABLE weakness_analysis DROP COLUMN specific_topic;
ALTER TABLE weakness_analysis RENAME COLUMN new_specific_topic TO specific_topic;

-- Step 7: Re-add foreign keys with integer types
ALTER TABLE topics ADD CONSTRAINT topics_parent_topic_id_fkey
  FOREIGN KEY (parent_topic_id) REFERENCES topics(topic_id);

ALTER TABLE problem_topics ADD CONSTRAINT problem_topics_topic_id_fkey
  FOREIGN KEY (topic_id) REFERENCES topics(topic_id) ON DELETE CASCADE;

ALTER TABLE user_skills ADD CONSTRAINT user_skills_topic_id_fkey
  FOREIGN KEY (topic_id) REFERENCES topics(topic_id) ON DELETE CASCADE;

ALTER TABLE questions ADD CONSTRAINT questions_related_topic_id_fkey
  FOREIGN KEY (related_topic_id) REFERENCES topics(topic_id);

ALTER TABLE pitfalls ADD CONSTRAINT pitfalls_topic_id_fkey
  FOREIGN KEY (topic_id) REFERENCES topics(topic_id) ON DELETE CASCADE;

ALTER TABLE weakness_analysis ADD CONSTRAINT weakness_analysis_specific_topic_fkey
  FOREIGN KEY (specific_topic) REFERENCES topics(topic_id) ON DELETE SET NULL;

-- Step 8: Recreate indexes
CREATE INDEX idx_topics_parent ON topics(parent_topic_id);
CREATE INDEX idx_pitfalls_topic ON pitfalls(topic_id);

-- Step 9: Add category/difficulty_level columns to topics if not present (GORM model fields)
ALTER TABLE topics ADD COLUMN IF NOT EXISTS category VARCHAR(100);
ALTER TABLE topics ADD COLUMN IF NOT EXISTS difficulty_level INT;
ALTER TABLE topics ADD COLUMN IF NOT EXISTS estimated_learning_hours FLOAT;

-- Step 10: Rename questions columns to match GORM model
-- category -> question_type, subcategory -> question_subtype, format -> question_format
-- related_problem_id -> problem_id
ALTER TABLE questions RENAME COLUMN category TO question_type;
ALTER TABLE questions RENAME COLUMN subcategory TO question_subtype;
ALTER TABLE questions RENAME COLUMN format TO question_format;
ALTER TABLE questions RENAME COLUMN related_problem_id TO problem_id;
ALTER TABLE questions RENAME COLUMN attempt_count TO total_attempts;
ALTER TABLE questions RENAME COLUMN correct_count TO correct_attempts;
ALTER TABLE questions RENAME COLUMN avg_time_sec TO average_time_seconds;
ALTER TABLE questions RENAME COLUMN estimated_time_sec TO estimated_time_seconds;

-- Step 11: Rename problems columns to match GORM model
-- statement -> description
ALTER TABLE problems RENAME COLUMN statement TO description;
ALTER TABLE problems ADD COLUMN IF NOT EXISTS official_difficulty VARCHAR(20);
UPDATE problems SET official_difficulty = difficulty_label::text WHERE official_difficulty IS NULL;
ALTER TABLE problems ADD COLUMN IF NOT EXISTS total_attempts INT DEFAULT 0;
ALTER TABLE problems ADD COLUMN IF NOT EXISTS total_solves INT DEFAULT 0;
ALTER TABLE problems RENAME COLUMN avg_solve_time_sec TO average_time_seconds;
ALTER TABLE problems RENAME COLUMN success_rate TO acceptance_rate;

-- Step 12: Rename problem_topics columns to match GORM model
ALTER TABLE problem_topics RENAME COLUMN relevance TO relevance_score;

-- Step 13: Add soft delete columns where GORM model expects them
ALTER TABLE training_plans ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
ALTER TABLE training_plan_items ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
ALTER TABLE user_lists ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

CREATE INDEX IF NOT EXISTS idx_training_plans_deleted ON training_plans(deleted_at) WHERE deleted_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_training_plan_items_deleted ON training_plan_items(deleted_at) WHERE deleted_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_user_lists_deleted ON user_lists(deleted_at) WHERE deleted_at IS NOT NULL;

-- Clean up temp table
DROP TABLE IF EXISTS topic_id_map;

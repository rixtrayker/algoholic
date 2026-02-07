-- ============================================================
-- 001_core_schema.sql
-- DSA Problem-Solving Gym - Core Schema
-- PostgreSQL + Apache AGE + pgvector
-- ============================================================

-- Extensions
CREATE EXTENSION IF NOT EXISTS age;
CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS btree_gin;

-- Load AGE
LOAD 'age';
SET search_path = ag_catalog, "$user", public;

-- Create graph for knowledge relationships
SELECT create_graph('dsa_graph');

-- ============================================================
-- ENUMS
-- ============================================================

CREATE TYPE difficulty_label AS ENUM ('easy', 'medium', 'hard', 'expert');
CREATE TYPE question_format AS ENUM ('multiple_choice', 'code', 'ranking', 'open_ended', 'fill_blank', 'debug');
CREATE TYPE topic_level AS ENUM ('category', 'topic', 'subtopic', 'pattern', 'variation');

-- ============================================================
-- CATEGORIES (top-level groupings)
-- ============================================================

CREATE TABLE categories (
    category_id   SERIAL PRIMARY KEY,
    slug          VARCHAR(100) UNIQUE NOT NULL,
    name          VARCHAR(200) NOT NULL,
    description   TEXT,
    icon          VARCHAR(50),        -- emoji or icon name for UI
    color         VARCHAR(7),         -- hex color
    display_order INT DEFAULT 0,
    created_at    TIMESTAMP DEFAULT NOW()
);

-- ============================================================
-- TOPICS (hierarchical knowledge nodes)
-- ============================================================

CREATE TABLE topics (
    topic_id        VARCHAR(100) PRIMARY KEY,  -- e.g. "two_pointers"
    name            VARCHAR(200) NOT NULL,
    slug            VARCHAR(200) UNIQUE NOT NULL,
    category_id     INT REFERENCES categories(category_id),
    parent_topic_id VARCHAR(100) REFERENCES topics(topic_id),
    level           topic_level NOT NULL DEFAULT 'topic',
    description     TEXT,
    keywords        TEXT[],
    difficulty_min  FLOAT DEFAULT 0 CHECK (difficulty_min >= 0 AND difficulty_min <= 100),
    difficulty_max  FLOAT DEFAULT 100 CHECK (difficulty_max >= 0 AND difficulty_max <= 100),
    display_order   INT DEFAULT 0,
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_topics_category ON topics(category_id);
CREATE INDEX idx_topics_parent ON topics(parent_topic_id);
CREATE INDEX idx_topics_level ON topics(level);
CREATE INDEX idx_topics_keywords ON topics USING GIN(keywords);

-- ============================================================
-- TAGS (flat labels for flexible filtering)
-- ============================================================

CREATE TABLE tags (
    tag_id      SERIAL PRIMARY KEY,
    slug        VARCHAR(100) UNIQUE NOT NULL,
    name        VARCHAR(100) NOT NULL,
    tag_group   VARCHAR(50),          -- 'pattern', 'technique', 'company', 'difficulty_trait', 'data_structure'
    description TEXT,
    color       VARCHAR(7),
    created_at  TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_tags_group ON tags(tag_group);

-- ============================================================
-- PROBLEMS (LeetCode-style coding problems)
-- ============================================================

CREATE TABLE problems (
    problem_id          SERIAL PRIMARY KEY,
    leetcode_number     INTEGER UNIQUE,
    slug                VARCHAR(300) UNIQUE NOT NULL,
    title               VARCHAR(500) NOT NULL,
    statement           TEXT NOT NULL,
    constraints         TEXT[],
    examples            JSONB NOT NULL DEFAULT '[]',
    hints               TEXT[],

    -- difficulty
    difficulty_score    FLOAT NOT NULL CHECK (difficulty_score BETWEEN 0 AND 100),
    difficulty_label    difficulty_label NOT NULL,

    -- solution metadata
    primary_pattern     VARCHAR(100),
    secondary_patterns  VARCHAR(100)[],
    time_complexity     VARCHAR(50),
    space_complexity    VARCHAR(50),

    -- stats
    frequency           FLOAT DEFAULT 0,
    success_rate        FLOAT,
    avg_solve_time_sec  INT,
    companies           JSONB DEFAULT '[]',

    -- flags
    has_follow_ups      BOOLEAN DEFAULT FALSE,
    has_prerequisites   BOOLEAN DEFAULT FALSE,

    -- full-text search
    search_vector       tsvector GENERATED ALWAYS AS (
        setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(statement, '')), 'B')
    ) STORED,

    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_problems_difficulty ON problems(difficulty_score);
CREATE INDEX idx_problems_label ON problems(difficulty_label);
CREATE INDEX idx_problems_leetcode ON problems(leetcode_number);
CREATE INDEX idx_problems_pattern ON problems(primary_pattern);
CREATE INDEX idx_problems_search ON problems USING GIN(search_vector);

-- ============================================================
-- JUNCTION: problem <-> topic
-- ============================================================

CREATE TABLE problem_topics (
    problem_id      INT REFERENCES problems(problem_id) ON DELETE CASCADE,
    topic_id        VARCHAR(100) REFERENCES topics(topic_id) ON DELETE CASCADE,
    relevance       FLOAT DEFAULT 1.0 CHECK (relevance BETWEEN 0 AND 1),
    is_primary      BOOLEAN DEFAULT FALSE,
    pattern_used    VARCHAR(200),
    key_insight     TEXT,
    PRIMARY KEY (problem_id, topic_id)
);

CREATE INDEX idx_pt_problem ON problem_topics(problem_id);
CREATE INDEX idx_pt_topic ON problem_topics(topic_id);

-- ============================================================
-- JUNCTION: problem <-> tag
-- ============================================================

CREATE TABLE problem_tags (
    problem_id  INT REFERENCES problems(problem_id) ON DELETE CASCADE,
    tag_id      INT REFERENCES tags(tag_id) ON DELETE CASCADE,
    PRIMARY KEY (problem_id, tag_id)
);

-- ============================================================
-- QUESTION TYPES (taxonomy registry)
-- ============================================================

CREATE TABLE question_types (
    type_id         SERIAL PRIMARY KEY,
    slug            VARCHAR(100) UNIQUE NOT NULL,
    category_code   VARCHAR(10) NOT NULL,         -- '1.1', '2.3', '10.1'
    name            VARCHAR(200) NOT NULL,
    description     TEXT,
    format          question_format NOT NULL,
    parent_category VARCHAR(100),                 -- 'complexity_analysis', 'ds_selection', etc.
    difficulty_default difficulty_label DEFAULT 'medium',
    estimated_time_sec INT,
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMP DEFAULT NOW()
);

-- ============================================================
-- QUESTIONS / CHALLENGES (the core assessment items)
-- ============================================================

CREATE TABLE questions (
    question_id         SERIAL PRIMARY KEY,

    -- type & classification
    question_type_id    INT REFERENCES question_types(type_id),
    category            VARCHAR(50) NOT NULL,
    subcategory         VARCHAR(50),

    -- content
    question_text       TEXT NOT NULL,
    question_data       JSONB,               -- code snippets, extra context
    format              question_format NOT NULL DEFAULT 'multiple_choice',

    -- answers
    correct_answer      JSONB NOT NULL,
    answer_options      JSONB,               -- [{id, text, is_correct}]
    wrong_answer_explanations JSONB,

    -- learning
    explanation         TEXT NOT NULL,
    detailed_solution   TEXT,
    common_mistakes     TEXT[],

    -- hints (3-level progressive)
    hint_level_1        TEXT,                -- socratic question
    hint_level_2        TEXT,                -- directional hint
    hint_level_3        TEXT,                -- concrete hint

    -- difficulty
    difficulty_score    FLOAT CHECK (difficulty_score BETWEEN 0 AND 100),
    difficulty_label    difficulty_label DEFAULT 'medium',
    estimated_time_sec  INT,

    -- relations
    related_problem_id  INT REFERENCES problems(problem_id),
    related_topic_id    VARCHAR(100) REFERENCES topics(topic_id),

    -- tags for filtering
    tags                TEXT[],
    concepts            TEXT[],

    -- stats (updated at runtime)
    attempt_count       INT DEFAULT 0,
    correct_count       INT DEFAULT 0,
    avg_time_sec        INT,

    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_questions_type ON questions(question_type_id);
CREATE INDEX idx_questions_category ON questions(category);
CREATE INDEX idx_questions_difficulty ON questions(difficulty_label);
CREATE INDEX idx_questions_problem ON questions(related_problem_id);
CREATE INDEX idx_questions_topic ON questions(related_topic_id);
CREATE INDEX idx_questions_tags ON questions USING GIN(tags);
CREATE INDEX idx_questions_concepts ON questions USING GIN(concepts);

-- ============================================================
-- JUNCTION: question <-> tag
-- ============================================================

CREATE TABLE question_tags (
    question_id INT REFERENCES questions(question_id) ON DELETE CASCADE,
    tag_id      INT REFERENCES tags(tag_id) ON DELETE CASCADE,
    PRIMARY KEY (question_id, tag_id)
);

-- ============================================================
-- PROBLEM RELATIONSHIPS (follow-ups, variations, prerequisites)
-- ============================================================

CREATE TABLE problem_relationships (
    id                  SERIAL PRIMARY KEY,
    source_problem_id   INT REFERENCES problems(problem_id) ON DELETE CASCADE,
    target_problem_id   INT REFERENCES problems(problem_id) ON DELETE CASCADE,
    relationship_type   VARCHAR(50) NOT NULL, -- 'follow_up', 'variation', 'prerequisite', 'similar'
    description         TEXT,
    created_at          TIMESTAMP DEFAULT NOW(),
    UNIQUE(source_problem_id, target_problem_id, relationship_type)
);

-- ============================================================
-- EMBEDDINGS (vector storage for semantic search)
-- ============================================================

CREATE TABLE embeddings (
    embedding_id    SERIAL PRIMARY KEY,
    entity_type     VARCHAR(50) NOT NULL,    -- 'problem', 'question', 'topic'
    entity_id       VARCHAR(100) NOT NULL,
    embedding       vector(1536),
    text_content    TEXT,
    model_name      VARCHAR(100),
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_embeddings_entity ON embeddings(entity_type, entity_id);

-- ============================================================
-- PITFALLS (common mistakes per topic)
-- ============================================================

CREATE TABLE pitfalls (
    pitfall_id  SERIAL PRIMARY KEY,
    topic_id    VARCHAR(100) REFERENCES topics(topic_id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    example     TEXT,
    fix         TEXT,
    severity    INT DEFAULT 3 CHECK (severity BETWEEN 1 AND 5),
    created_at  TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_pitfalls_topic ON pitfalls(topic_id);

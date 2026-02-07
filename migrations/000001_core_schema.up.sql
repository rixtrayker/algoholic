-- 000001_core_schema.up.sql
-- Core schema for Algoholic DSA training platform

-- Extensions (available on standard PostgreSQL)
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS btree_gin;

-- Enums
CREATE TYPE difficulty_label AS ENUM ('easy', 'medium', 'hard', 'expert');
CREATE TYPE question_format  AS ENUM ('multiple_choice', 'code', 'ranking', 'open_ended', 'fill_blank', 'debug');
CREATE TYPE topic_level      AS ENUM ('category', 'topic', 'subtopic', 'pattern', 'variation');

-- Categories
CREATE TABLE categories (
    category_id   SERIAL PRIMARY KEY,
    slug          VARCHAR(100) UNIQUE NOT NULL,
    name          VARCHAR(200) NOT NULL,
    description   TEXT,
    icon          VARCHAR(50),
    color         VARCHAR(7),
    display_order INT DEFAULT 0,
    created_at    TIMESTAMP DEFAULT NOW()
);

-- Topics (hierarchical)
CREATE TABLE topics (
    topic_id        VARCHAR(100) PRIMARY KEY,
    name            VARCHAR(200) NOT NULL,
    slug            VARCHAR(200) UNIQUE NOT NULL,
    category_id     INT REFERENCES categories(category_id),
    parent_topic_id VARCHAR(100) REFERENCES topics(topic_id),
    level           topic_level NOT NULL DEFAULT 'topic',
    description     TEXT,
    keywords        TEXT[],
    difficulty_min  FLOAT DEFAULT 0  CHECK (difficulty_min >= 0 AND difficulty_min <= 100),
    difficulty_max  FLOAT DEFAULT 100 CHECK (difficulty_max >= 0 AND difficulty_max <= 100),
    display_order   INT DEFAULT 0,
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_topics_category ON topics(category_id);
CREATE INDEX idx_topics_parent   ON topics(parent_topic_id);
CREATE INDEX idx_topics_level    ON topics(level);
CREATE INDEX idx_topics_keywords ON topics USING GIN(keywords);

-- Tags
CREATE TABLE tags (
    tag_id      SERIAL PRIMARY KEY,
    slug        VARCHAR(100) UNIQUE NOT NULL,
    name        VARCHAR(100) NOT NULL,
    tag_group   VARCHAR(50),
    description TEXT,
    color       VARCHAR(7),
    created_at  TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_tags_group ON tags(tag_group);

-- Problems
CREATE TABLE problems (
    problem_id         SERIAL PRIMARY KEY,
    leetcode_number    INTEGER UNIQUE,
    slug               VARCHAR(300) UNIQUE NOT NULL,
    title              VARCHAR(500) NOT NULL,
    statement          TEXT NOT NULL,
    constraints        TEXT[],
    examples           JSONB NOT NULL DEFAULT '[]',
    hints              TEXT[],
    difficulty_score   FLOAT NOT NULL CHECK (difficulty_score BETWEEN 0 AND 100),
    difficulty_label   difficulty_label NOT NULL,
    primary_pattern    VARCHAR(100),
    secondary_patterns VARCHAR(100)[],
    time_complexity    VARCHAR(50),
    space_complexity   VARCHAR(50),
    frequency          FLOAT DEFAULT 0,
    success_rate       FLOAT,
    avg_solve_time_sec INT,
    companies          JSONB DEFAULT '[]',
    has_follow_ups     BOOLEAN DEFAULT FALSE,
    has_prerequisites  BOOLEAN DEFAULT FALSE,
    search_vector      tsvector GENERATED ALWAYS AS (
        setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(statement, '')), 'B')
    ) STORED,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_problems_difficulty ON problems(difficulty_score);
CREATE INDEX idx_problems_label      ON problems(difficulty_label);
CREATE INDEX idx_problems_leetcode   ON problems(leetcode_number);
CREATE INDEX idx_problems_pattern    ON problems(primary_pattern);
CREATE INDEX idx_problems_search     ON problems USING GIN(search_vector);

-- Problem <-> Topic
CREATE TABLE problem_topics (
    problem_id   INT REFERENCES problems(problem_id) ON DELETE CASCADE,
    topic_id     VARCHAR(100) REFERENCES topics(topic_id) ON DELETE CASCADE,
    relevance    FLOAT DEFAULT 1.0 CHECK (relevance BETWEEN 0 AND 1),
    is_primary   BOOLEAN DEFAULT FALSE,
    pattern_used VARCHAR(200),
    key_insight  TEXT,
    PRIMARY KEY (problem_id, topic_id)
);

-- Problem <-> Tag
CREATE TABLE problem_tags (
    problem_id INT REFERENCES problems(problem_id) ON DELETE CASCADE,
    tag_id     INT REFERENCES tags(tag_id) ON DELETE CASCADE,
    PRIMARY KEY (problem_id, tag_id)
);

-- Question Types
CREATE TABLE question_types (
    type_id            SERIAL PRIMARY KEY,
    slug               VARCHAR(100) UNIQUE NOT NULL,
    category_code      VARCHAR(10) NOT NULL,
    name               VARCHAR(200) NOT NULL,
    description        TEXT,
    format             question_format NOT NULL,
    parent_category    VARCHAR(100),
    difficulty_default difficulty_label DEFAULT 'medium',
    estimated_time_sec INT,
    metadata           JSONB DEFAULT '{}',
    created_at         TIMESTAMP DEFAULT NOW()
);

-- Questions
CREATE TABLE questions (
    question_id             SERIAL PRIMARY KEY,
    question_type_id        INT REFERENCES question_types(type_id),
    category                VARCHAR(50) NOT NULL,
    subcategory             VARCHAR(50),
    question_text           TEXT NOT NULL,
    question_data           JSONB,
    format                  question_format NOT NULL DEFAULT 'multiple_choice',
    correct_answer          JSONB NOT NULL,
    answer_options          JSONB,
    wrong_answer_explanations JSONB,
    explanation             TEXT NOT NULL,
    detailed_solution       TEXT,
    common_mistakes         TEXT[],
    hint_level_1            TEXT,
    hint_level_2            TEXT,
    hint_level_3            TEXT,
    difficulty_score        FLOAT CHECK (difficulty_score BETWEEN 0 AND 100),
    difficulty_label        difficulty_label DEFAULT 'medium',
    estimated_time_sec      INT,
    related_problem_id      INT REFERENCES problems(problem_id),
    related_topic_id        VARCHAR(100) REFERENCES topics(topic_id),
    tags                    TEXT[],
    concepts                TEXT[],
    attempt_count           INT DEFAULT 0,
    correct_count           INT DEFAULT 0,
    avg_time_sec            INT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_questions_type       ON questions(question_type_id);
CREATE INDEX idx_questions_category   ON questions(category);
CREATE INDEX idx_questions_difficulty ON questions(difficulty_label);
CREATE INDEX idx_questions_problem    ON questions(related_problem_id);
CREATE INDEX idx_questions_topic      ON questions(related_topic_id);
CREATE INDEX idx_questions_tags       ON questions USING GIN(tags);
CREATE INDEX idx_questions_concepts   ON questions USING GIN(concepts);

-- Question <-> Tag
CREATE TABLE question_tags (
    question_id INT REFERENCES questions(question_id) ON DELETE CASCADE,
    tag_id      INT REFERENCES tags(tag_id) ON DELETE CASCADE,
    PRIMARY KEY (question_id, tag_id)
);

-- Problem Relationships
CREATE TABLE problem_relationships (
    id                SERIAL PRIMARY KEY,
    source_problem_id INT REFERENCES problems(problem_id) ON DELETE CASCADE,
    target_problem_id INT REFERENCES problems(problem_id) ON DELETE CASCADE,
    relationship_type VARCHAR(50) NOT NULL,
    description       TEXT,
    created_at        TIMESTAMP DEFAULT NOW(),
    UNIQUE(source_problem_id, target_problem_id, relationship_type)
);

-- Pitfalls
CREATE TABLE pitfalls (
    pitfall_id SERIAL PRIMARY KEY,
    topic_id   VARCHAR(100) REFERENCES topics(topic_id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    example     TEXT,
    fix         TEXT,
    severity    INT DEFAULT 3 CHECK (severity BETWEEN 1 AND 5),
    created_at  TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_pitfalls_topic ON pitfalls(topic_id);

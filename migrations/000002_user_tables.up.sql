-- 000002_user_tables.up.sql
-- User-related tables for Algoholic DSA training platform

-- Users
CREATE TABLE users (
    user_id                SERIAL PRIMARY KEY,
    username               VARCHAR(100) UNIQUE NOT NULL,
    email                  VARCHAR(255) UNIQUE NOT NULL,
    password_hash          VARCHAR(255) NOT NULL,
    preferences            JSONB DEFAULT '{}',
    current_streak_days    INT DEFAULT 0,
    total_study_time_seconds BIGINT DEFAULT 0,
    created_at             TIMESTAMP DEFAULT NOW(),
    last_active_at         TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);

-- User Attempts (tracks question/problem attempts)
CREATE TABLE user_attempts (
    attempt_id         SERIAL PRIMARY KEY,
    user_id            INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    question_id        INT REFERENCES questions(question_id) ON DELETE SET NULL,
    problem_id         INT REFERENCES problems(problem_id) ON DELETE SET NULL,
    user_answer        JSONB NOT NULL,
    is_correct         BOOLEAN NOT NULL,
    time_taken_seconds INT NOT NULL,
    attempt_number     INT DEFAULT 1,
    hints_used         INT DEFAULT 0,
    confidence_level   INT CHECK (confidence_level BETWEEN 1 AND 5),
    detected_patterns  TEXT[],
    mistakes_made      TEXT[],
    shows_memorization BOOLEAN,
    training_plan_id   INT,
    session_id         VARCHAR(100),
    attempted_at       TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_attempts_user ON user_attempts(user_id);
CREATE INDEX idx_attempts_question ON user_attempts(question_id);
CREATE INDEX idx_attempts_problem ON user_attempts(problem_id);
CREATE INDEX idx_attempts_attempted_at ON user_attempts(attempted_at);

-- User Skills (tracks proficiency per topic)
CREATE TABLE user_skills (
    user_id             INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    topic_id            VARCHAR(100) NOT NULL REFERENCES topics(topic_id) ON DELETE CASCADE,
    proficiency_level   FLOAT DEFAULT 0 CHECK (proficiency_level >= 0 AND proficiency_level <= 100),
    questions_attempted INT DEFAULT 0,
    questions_correct   INT DEFAULT 0,
    improvement_rate    FLOAT,
    needs_review        BOOLEAN DEFAULT FALSE,
    last_practiced_at   TIMESTAMP,
    next_review_at      TIMESTAMP,
    PRIMARY KEY (user_id, topic_id)
);

CREATE INDEX idx_user_skills_needs_review ON user_skills(needs_review) WHERE needs_review = TRUE;
CREATE INDEX idx_user_skills_next_review ON user_skills(next_review_at) WHERE next_review_at IS NOT NULL;

-- Training Plans
CREATE TABLE training_plans (
    plan_id             SERIAL PRIMARY KEY,
    user_id             INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    name                VARCHAR(200) NOT NULL,
    description         TEXT,
    plan_type           VARCHAR(50),
    difficulty_range    VARCHAR(50),
    target_topics       INT[],
    target_patterns     TEXT[],
    duration_days       INT,
    questions_per_day   INT DEFAULT 5,
    adaptive_difficulty BOOLEAN DEFAULT TRUE,
    progress_percentage FLOAT DEFAULT 0,
    status              VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'paused', 'completed', 'cancelled')),
    start_date          TIMESTAMP NOT NULL,
    created_at          TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_training_plans_user ON training_plans(user_id);
CREATE INDEX idx_training_plans_status ON training_plans(status);

-- Training Plan Items
CREATE TABLE training_plan_items (
    item_id         SERIAL PRIMARY KEY,
    plan_id         INT NOT NULL REFERENCES training_plans(plan_id) ON DELETE CASCADE,
    question_id     INT REFERENCES questions(question_id) ON DELETE SET NULL,
    problem_id      INT REFERENCES problems(problem_id) ON DELETE SET NULL,
    sequence_number INT NOT NULL,
    day_number      INT,
    scheduled_for   TIMESTAMP,
    item_type       VARCHAR(20) NOT NULL CHECK (item_type IN ('question', 'problem', 'assessment')),
    is_completed    BOOLEAN DEFAULT FALSE,
    completed_at    TIMESTAMP
);

CREATE INDEX idx_plan_items_plan ON training_plan_items(plan_id);
CREATE INDEX idx_plan_items_scheduled ON training_plan_items(scheduled_for);
CREATE INDEX idx_plan_items_day ON training_plan_items(day_number);

-- Assessments
CREATE TABLE assessments (
    assessment_id      SERIAL PRIMARY KEY,
    user_id            INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    assessment_type    VARCHAR(50),
    topics_covered     TEXT[],
    overall_score      FLOAT CHECK (overall_score >= 0 AND overall_score <= 100),
    category_scores    JSONB DEFAULT '{}',
    strengths          TEXT[],
    weaknesses         TEXT[],
    recommendations    TEXT,
    memorization_score FLOAT CHECK (memorization_score >= 0 AND memorization_score <= 100),
    started_at         TIMESTAMP,
    completed_at       TIMESTAMP,
    time_taken_seconds INT
);

CREATE INDEX idx_assessments_user ON assessments(user_id);
CREATE INDEX idx_assessments_type ON assessments(assessment_type);

-- Weakness Analysis
CREATE TABLE weakness_analysis (
    analysis_id          SERIAL PRIMARY KEY,
    user_id              INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    weakness_type        VARCHAR(100) NOT NULL,
    specific_topic       VARCHAR(100) REFERENCES topics(topic_id) ON DELETE SET NULL,
    severity             VARCHAR(20) NOT NULL CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    weakness_score       FLOAT NOT NULL CHECK (weakness_score >= 0 AND weakness_score <= 100),
    evidence_question_ids INT[],
    pattern_description  TEXT,
    recommended_practice JSONB DEFAULT '{}',
    detected_at          TIMESTAMP DEFAULT NOW(),
    resolved_at          TIMESTAMP,
    is_active            BOOLEAN DEFAULT TRUE
);

CREATE INDEX idx_weakness_user ON weakness_analysis(user_id);
CREATE INDEX idx_weakness_active ON weakness_analysis(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_weakness_type ON weakness_analysis(weakness_type);

-- User Lists (custom problem lists)
CREATE TABLE user_lists (
    list_id     SERIAL PRIMARY KEY,
    user_id     INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    name        VARCHAR(200) NOT NULL,
    description TEXT,
    is_public   BOOLEAN DEFAULT FALSE,
    problem_ids JSONB NOT NULL DEFAULT '[]',
    total_items INT DEFAULT 0,
    completed   INT DEFAULT 0,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_user_lists_user ON user_lists(user_id);
CREATE INDEX idx_user_lists_public ON user_lists(is_public) WHERE is_public = TRUE;

-- Daily Activities (for commitment/streak tracking)
CREATE TABLE daily_activities (
    activity_id     SERIAL PRIMARY KEY,
    user_id         INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    date            DATE NOT NULL,
    problems_count  INT DEFAULT 0,
    questions_count INT DEFAULT 0,
    study_time_seconds INT DEFAULT 0,
    streak          INT DEFAULT 0,
    UNIQUE(user_id, date)
);

CREATE INDEX idx_daily_activities_user ON daily_activities(user_id);
CREATE INDEX idx_daily_activities_date ON daily_activities(date);

-- Spaced Repetition Reviews (SM-2 algorithm support)
CREATE TABLE spaced_repetition_reviews (
    review_id       SERIAL PRIMARY KEY,
    user_id         INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    question_id     INT NOT NULL REFERENCES questions(question_id) ON DELETE CASCADE,
    easiness_factor FLOAT DEFAULT 2.5 CHECK (easiness_factor >= 1.3),
    interval_days   INT DEFAULT 1,
    repetitions     INT DEFAULT 0,
    next_review_at  TIMESTAMP NOT NULL,
    last_review_at  TIMESTAMP,
    quality_rating  INT CHECK (quality_rating BETWEEN 0 AND 5),
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, question_id)
);

CREATE INDEX idx_spaced_rep_user ON spaced_repetition_reviews(user_id);
CREATE INDEX idx_spaced_rep_next_review ON spaced_repetition_reviews(next_review_at);
CREATE INDEX idx_spaced_rep_question ON spaced_repetition_reviews(question_id);

-- Review Queue (questions scheduled for review)
CREATE TABLE review_queue (
    queue_id        SERIAL PRIMARY KEY,
    user_id         INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    question_id     INT NOT NULL REFERENCES questions(question_id) ON DELETE CASCADE,
    scheduled_for   TIMESTAMP NOT NULL,
    priority        INT DEFAULT 0,
    is_overdue      BOOLEAN DEFAULT FALSE,
    added_at        TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, question_id)
);

CREATE INDEX idx_review_queue_user ON review_queue(user_id);
CREATE INDEX idx_review_queue_scheduled ON review_queue(scheduled_for);
CREATE INDEX idx_review_queue_overdue ON review_queue(is_overdue) WHERE is_overdue = TRUE;

-- Code Submissions (for AI assessment - Phase 4)
CREATE TABLE code_submissions (
    submission_id   SERIAL PRIMARY KEY,
    user_id         INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    problem_id      INT NOT NULL REFERENCES problems(problem_id) ON DELETE CASCADE,
    code            TEXT NOT NULL,
    language        VARCHAR(50) NOT NULL,
    status          VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'running', 'passed', 'failed', 'error')),
    test_results    JSONB DEFAULT '{}',
    ai_feedback     JSONB,
    ai_score        FLOAT CHECK (ai_score >= 0 AND ai_score <= 100),
    time_complexity VARCHAR(50),
    space_complexity VARCHAR(50),
    execution_time_ms INT,
    memory_used_kb  INT,
    submitted_at    TIMESTAMP DEFAULT NOW(),
    evaluated_at    TIMESTAMP
);

CREATE INDEX idx_code_submissions_user ON code_submissions(user_id);
CREATE INDEX idx_code_submissions_problem ON code_submissions(problem_id);
CREATE INDEX idx_code_submissions_status ON code_submissions(status);

-- Question Hints Used (track which hints user has seen)
CREATE TABLE question_hint_usage (
    usage_id    SERIAL PRIMARY KEY,
    user_id     INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    question_id INT NOT NULL REFERENCES questions(question_id) ON DELETE CASCADE,
    hint_level  INT NOT NULL CHECK (hint_level BETWEEN 1 AND 3),
    used_at     TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, question_id, hint_level)
);

CREATE INDEX idx_hint_usage_user_question ON question_hint_usage(user_id, question_id);

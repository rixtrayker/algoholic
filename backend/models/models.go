package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// JSONB is a custom type for PostgreSQL JSONB columns
type JSONB map[string]interface{}

// JSONBArray is a custom type for PostgreSQL JSONB columns that store arrays
type JSONBArray []interface{}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal JSONB value")
	}
	return json.Unmarshal(bytes, j)
}

func (j JSONBArray) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONBArray) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal JSONBArray value")
	}
	return json.Unmarshal(bytes, j)
}

// StringArray is a custom type for PostgreSQL text[] columns
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	if len(a) == 0 {
		return "{}", nil
	}
	escaped := make([]string, len(a))
	for i, s := range a {
		escaped[i] = `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
	}
	return "{" + strings.Join(escaped, ",") + "}", nil
}

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return a.parsePostgresArray(string(v))
	case string:
		return a.parsePostgresArray(v)
	default:
		return errors.New("failed to unmarshal StringArray value")
	}
}

func (a *StringArray) parsePostgresArray(s string) error {
	if s == "{}" || s == "" {
		*a = []string{}
		return nil
	}
	s = strings.Trim(s, "{}")
	if s == "" {
		*a = []string{}
		return nil
	}
	var result []string
	var current strings.Builder
	inQuote := false
	escaped := false
	for _, r := range s {
		if escaped {
			current.WriteRune(r)
			escaped = false
			continue
		}
		switch r {
		case '\\':
			escaped = true
		case '"':
			inQuote = !inQuote
		case ',':
			if !inQuote {
				result = append(result, current.String())
				current.Reset()
			} else {
				current.WriteRune(r)
			}
		default:
			current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		result = append(result, current.String())
	}
	*a = result
	return nil
}

// User represents a platform user
type User struct {
	UserID            int       `json:"user_id" gorm:"primaryKey;column:user_id"`
	Username          string    `json:"username" gorm:"column:username;uniqueIndex;not null"`
	Email             string    `json:"email" gorm:"column:email;uniqueIndex;not null"`
	PasswordHash      string    `json:"-" gorm:"column:password_hash;not null"`
	Preferences       JSONB     `json:"preferences,omitempty" gorm:"column:preferences;type:jsonb"`
	CurrentStreakDays int       `json:"current_streak_days" gorm:"column:current_streak_days;default:0"`
	TotalStudyTime    int64     `json:"total_study_time_seconds" gorm:"column:total_study_time_seconds;default:0"`
	CreatedAt         time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	LastActiveAt      time.Time `json:"last_active_at" gorm:"column:last_active_at;autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}

// Problem represents a coding problem
type Problem struct {
	ProblemID          int         `json:"problem_id" gorm:"primaryKey;column:problem_id"`
	LeetcodeNumber     *int        `json:"leetcode_number,omitempty" gorm:"column:leetcode_number;uniqueIndex"`
	Title              string      `json:"title" gorm:"column:title;not null"`
	Slug               string      `json:"slug" gorm:"column:slug;uniqueIndex;not null"`
	Description        string      `json:"description" gorm:"column:description;not null"`
	Constraints        StringArray `json:"constraints,omitempty" gorm:"column:constraints;type:text[]"`
	Examples           JSONBArray  `json:"examples" gorm:"column:examples;type:jsonb;not null"`
	Hints              StringArray `json:"hints,omitempty" gorm:"column:hints;type:text[]"`
	DifficultyScore    float64     `json:"difficulty_score" gorm:"column:difficulty_score;not null;check:difficulty_score >= 0 AND difficulty_score <= 100"`
	OfficialDifficulty *string     `json:"official_difficulty,omitempty" gorm:"column:official_difficulty"`
	PrimaryPattern     *string     `json:"primary_pattern,omitempty" gorm:"column:primary_pattern"`
	SecondaryPatterns  StringArray `json:"secondary_patterns,omitempty" gorm:"column:secondary_patterns;type:text[]"`
	Source             *string     `json:"source,omitempty" gorm:"column:source"`
	TimeComplexity     *string     `json:"time_complexity,omitempty" gorm:"column:time_complexity"`
	SpaceComplexity    *string     `json:"space_complexity,omitempty" gorm:"column:space_complexity"`
	TotalAttempts      int         `json:"total_attempts" gorm:"column:total_attempts;default:0"`
	TotalSolves        int         `json:"total_solves" gorm:"column:total_solves;default:0"`
	AverageTime        *float64    `json:"average_time_seconds,omitempty" gorm:"column:average_time_seconds"`
	AcceptanceRate     *float64    `json:"acceptance_rate,omitempty" gorm:"column:acceptance_rate"`
	Companies          JSONB       `json:"companies,omitempty" gorm:"column:companies;type:jsonb"`
	Tags               JSONB       `json:"tags,omitempty" gorm:"column:tags;type:jsonb"`
	CreatedAt          time.Time   `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          time.Time   `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (Problem) TableName() string {
	return "problems"
}

// Topic represents a learning topic
type Topic struct {
	TopicID               int       `json:"topic_id" gorm:"primaryKey;column:topic_id"`
	Name                  string    `json:"name" gorm:"column:name;uniqueIndex;not null"`
	Slug                  string    `json:"slug" gorm:"column:slug;uniqueIndex;not null"`
	Description           *string   `json:"description,omitempty" gorm:"column:description"`
	ParentTopicID         *int      `json:"parent_topic_id,omitempty" gorm:"column:parent_topic_id"`
	Category              *string   `json:"category,omitempty" gorm:"column:category"`
	DifficultyLevel       *int      `json:"difficulty_level,omitempty" gorm:"column:difficulty_level"`
	EstimatedLearningHour *float64  `json:"estimated_learning_hours,omitempty" gorm:"column:estimated_learning_hours"`
	CreatedAt             time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}

func (Topic) TableName() string {
	return "topics"
}

// ProblemTopic links problems to topics
type ProblemTopic struct {
	ProblemID      int     `json:"problem_id" gorm:"primaryKey;column:problem_id"`
	TopicID        int     `json:"topic_id" gorm:"primaryKey;column:topic_id"`
	RelevanceScore float64 `json:"relevance_score" gorm:"column:relevance_score;default:1.0"`
	IsPrimary      bool    `json:"is_primary" gorm:"column:is_primary;default:false"`
}

func (ProblemTopic) TableName() string {
	return "problem_topics"
}

// Question represents a practice question
type Question struct {
	QuestionID              int         `json:"question_id" gorm:"primaryKey;column:question_id"`
	ProblemID               *int        `json:"problem_id,omitempty" gorm:"column:problem_id"`
	QuestionType            string      `json:"question_type" gorm:"column:question_type;not null"`
	QuestionSubtype         *string     `json:"question_subtype,omitempty" gorm:"column:question_subtype"`
	QuestionFormat          string      `json:"question_format" gorm:"column:question_format;not null"`
	QuestionText            string      `json:"question_text" gorm:"column:question_text;not null"`
	QuestionData            JSONB       `json:"question_data,omitempty" gorm:"column:question_data;type:jsonb"`
	AnswerOptions           JSONB       `json:"answer_options,omitempty" gorm:"column:answer_options;type:jsonb"`
	CorrectAnswer           JSONB       `json:"correct_answer" gorm:"column:correct_answer;type:jsonb;not null"`
	Explanation             string      `json:"explanation" gorm:"column:explanation;not null"`
	WrongAnswerExplanations JSONB       `json:"wrong_answer_explanations,omitempty" gorm:"column:wrong_answer_explanations;type:jsonb"`
	RelatedConcepts         StringArray `json:"related_concepts,omitempty" gorm:"column:related_concepts;type:text[]"`
	CommonMistakes          StringArray `json:"common_mistakes,omitempty" gorm:"column:common_mistakes;type:text[]"`
	DifficultyScore         float64     `json:"difficulty_score" gorm:"column:difficulty_score;not null"`
	EstimatedTimeSeconds    *int        `json:"estimated_time_seconds,omitempty" gorm:"column:estimated_time_seconds"`
	TotalAttempts           int         `json:"total_attempts" gorm:"column:total_attempts;default:0"`
	CorrectAttempts         int         `json:"correct_attempts" gorm:"column:correct_attempts;default:0"`
	AverageTimeSeconds      *float64    `json:"average_time_seconds,omitempty" gorm:"column:average_time_seconds"`
	CreatedAt               time.Time   `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt               time.Time   `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (Question) TableName() string {
	return "questions"
}

// UserAttempt tracks question/problem attempts
type UserAttempt struct {
	AttemptID         int         `json:"attempt_id" gorm:"primaryKey;column:attempt_id"`
	UserID            int         `json:"user_id" gorm:"column:user_id;not null;index:idx_attempts_user"`
	QuestionID        *int        `json:"question_id,omitempty" gorm:"column:question_id;index:idx_attempts_question"`
	ProblemID         *int        `json:"problem_id,omitempty" gorm:"column:problem_id"`
	UserAnswer        JSONB       `json:"user_answer" gorm:"column:user_answer;type:jsonb;not null"`
	IsCorrect         bool        `json:"is_correct" gorm:"column:is_correct;not null"`
	TimeTakenSeconds  int         `json:"time_taken_seconds" gorm:"column:time_taken_seconds;not null"`
	AttemptNumber     int         `json:"attempt_number" gorm:"column:attempt_number;default:1"`
	HintsUsed         int         `json:"hints_used" gorm:"column:hints_used;default:0"`
	ConfidenceLevel   *int        `json:"confidence_level,omitempty" gorm:"column:confidence_level"`
	DetectedPatterns  StringArray `json:"detected_patterns,omitempty" gorm:"column:detected_patterns;type:text[]"`
	MistakesMade      StringArray `json:"mistakes_made,omitempty" gorm:"column:mistakes_made;type:text[]"`
	ShowsMemorization *bool       `json:"shows_memorization,omitempty" gorm:"column:shows_memorization"`
	TrainingPlanID    *int        `json:"training_plan_id,omitempty" gorm:"column:training_plan_id"`
	SessionID         *string     `json:"session_id,omitempty" gorm:"column:session_id"`
	AttemptedAt       time.Time   `json:"attempted_at" gorm:"column:attempted_at;autoCreateTime"`
}

func (UserAttempt) TableName() string {
	return "user_attempts"
}

// UserSkill tracks user proficiency per topic
type UserSkill struct {
	UserID             int        `json:"user_id" gorm:"primaryKey;column:user_id"`
	TopicID            int        `json:"topic_id" gorm:"primaryKey;column:topic_id"`
	ProficiencyLevel   float64    `json:"proficiency_level" gorm:"column:proficiency_level;default:0;check:proficiency_level >= 0 AND proficiency_level <= 100"`
	QuestionsAttempted int        `json:"questions_attempted" gorm:"column:questions_attempted;default:0"`
	QuestionsCorrect   int        `json:"questions_correct" gorm:"column:questions_correct;default:0"`
	ImprovementRate    *float64   `json:"improvement_rate,omitempty" gorm:"column:improvement_rate"`
	NeedsReview        bool       `json:"needs_review" gorm:"column:needs_review;default:false"`
	LastPracticedAt    *time.Time `json:"last_practiced_at,omitempty" gorm:"column:last_practiced_at"`
	NextReviewAt       *time.Time `json:"next_review_at,omitempty" gorm:"column:next_review_at"`
}

func (UserSkill) TableName() string {
	return "user_skills"
}

// TrainingPlan represents a personalized training plan
type TrainingPlan struct {
	PlanID             int            `json:"plan_id" gorm:"primaryKey;column:plan_id"`
	UserID             int            `json:"user_id" gorm:"column:user_id;not null;index"`
	Name               string         `json:"name" gorm:"column:name;not null"`
	Description        *string        `json:"description,omitempty" gorm:"column:description"`
	PlanType           *string        `json:"plan_type,omitempty" gorm:"column:plan_type"`
	DifficultyRange    *string        `json:"difficulty_range,omitempty" gorm:"column:difficulty_range"`
	TargetTopics       StringArray    `json:"target_topics,omitempty" gorm:"column:target_topics;type:integer[]"`
	TargetPatterns     StringArray    `json:"target_patterns,omitempty" gorm:"column:target_patterns;type:text[]"`
	DurationDays       *int           `json:"duration_days,omitempty" gorm:"column:duration_days"`
	QuestionsPerDay    int            `json:"questions_per_day" gorm:"column:questions_per_day;default:5"`
	AdaptiveDifficulty bool           `json:"adaptive_difficulty" gorm:"column:adaptive_difficulty;default:true"`
	ProgressPercentage float64        `json:"progress_percentage" gorm:"column:progress_percentage;default:0"`
	Status             string         `json:"status" gorm:"column:status;default:'active'"`
	StartDate          time.Time      `json:"start_date" gorm:"column:start_date;not null"`
	CreatedAt          time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	DeletedAt          gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`
}

func (TrainingPlan) TableName() string {
	return "training_plans"
}

// TrainingPlanItem represents an item in a training plan
type TrainingPlanItem struct {
	ItemID         int            `json:"item_id" gorm:"primaryKey;column:item_id"`
	PlanID         int            `json:"plan_id" gorm:"column:plan_id;not null;index"`
	QuestionID     *int           `json:"question_id,omitempty" gorm:"column:question_id"`
	ProblemID      *int           `json:"problem_id,omitempty" gorm:"column:problem_id"`
	SequenceNumber int            `json:"sequence_number" gorm:"column:sequence_number;not null"`
	DayNumber      *int           `json:"day_number,omitempty" gorm:"column:day_number"`
	ScheduledFor   *time.Time     `json:"scheduled_for,omitempty" gorm:"column:scheduled_for"`
	ItemType       string         `json:"item_type" gorm:"column:item_type;not null"`
	IsCompleted    bool           `json:"is_completed" gorm:"column:is_completed;default:false"`
	CompletedAt    *time.Time     `json:"completed_at,omitempty" gorm:"column:completed_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`
}

func (TrainingPlanItem) TableName() string {
	return "training_plan_items"
}

// Assessment represents a user assessment
type Assessment struct {
	AssessmentID      int         `json:"assessment_id" gorm:"primaryKey;column:assessment_id"`
	UserID            int         `json:"user_id" gorm:"column:user_id;not null;index"`
	AssessmentType    *string     `json:"assessment_type,omitempty" gorm:"column:assessment_type"`
	TopicsCovered     StringArray `json:"topics_covered,omitempty" gorm:"column:topics_covered;type:text[]"`
	OverallScore      *float64    `json:"overall_score,omitempty" gorm:"column:overall_score"`
	CategoryScores    JSONB       `json:"category_scores,omitempty" gorm:"column:category_scores;type:jsonb"`
	Strengths         StringArray `json:"strengths,omitempty" gorm:"column:strengths;type:text[]"`
	Weaknesses        StringArray `json:"weaknesses,omitempty" gorm:"column:weaknesses;type:text[]"`
	Recommendations   *string     `json:"recommendations,omitempty" gorm:"column:recommendations"`
	MemorizationScore *float64    `json:"memorization_score,omitempty" gorm:"column:memorization_score"`
	StartedAt         *time.Time  `json:"started_at,omitempty" gorm:"column:started_at"`
	CompletedAt       *time.Time  `json:"completed_at,omitempty" gorm:"column:completed_at"`
	TimeTakenSeconds  *int        `json:"time_taken_seconds,omitempty" gorm:"column:time_taken_seconds"`
}

func (Assessment) TableName() string {
	return "assessments"
}

// WeaknessAnalysis tracks detected user weaknesses
type WeaknessAnalysis struct {
	AnalysisID          int         `json:"analysis_id" gorm:"primaryKey;column:analysis_id"`
	UserID              int         `json:"user_id" gorm:"column:user_id;not null;index"`
	WeaknessType        string      `json:"weakness_type" gorm:"column:weakness_type;not null"`
	SpecificTopic       *int        `json:"specific_topic,omitempty" gorm:"column:specific_topic"`
	Severity            string      `json:"severity" gorm:"column:severity;not null"`
	WeaknessScore       float64     `json:"weakness_score" gorm:"column:weakness_score;not null"`
	EvidenceQuestionIDs StringArray `json:"evidence_question_ids,omitempty" gorm:"column:evidence_question_ids;type:integer[]"`
	PatternDescription  *string     `json:"pattern_description,omitempty" gorm:"column:pattern_description"`
	RecommendedPractice JSONB       `json:"recommended_practice,omitempty" gorm:"column:recommended_practice;type:jsonb"`
	DetectedAt          time.Time   `json:"detected_at" gorm:"column:detected_at;autoCreateTime"`
	ResolvedAt          *time.Time  `json:"resolved_at,omitempty" gorm:"column:resolved_at"`
	IsActive            bool        `json:"is_active" gorm:"column:is_active;default:true"`
}

func (WeaknessAnalysis) TableName() string {
	return "weakness_analysis"
}

// UserList represents a custom problem list created by a user
type UserList struct {
	ListID      int            `json:"list_id" gorm:"primaryKey;column:list_id"`
	UserID      int            `json:"user_id" gorm:"column:user_id;not null;index"`
	Name        string         `json:"name" gorm:"column:name;not null"`
	Description *string        `json:"description,omitempty" gorm:"column:description"`
	IsPublic    bool           `json:"is_public" gorm:"column:is_public;default:false"`
	ProblemIDs  JSONB          `json:"problem_ids" gorm:"column:problem_ids;type:jsonb;not null;default:'[]'"`
	TotalItems  int            `json:"total_items" gorm:"column:total_items;default:0"`
	Completed   int            `json:"completed" gorm:"column:completed;default:0"`
	CreatedAt   time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`
}

func (UserList) TableName() string {
	return "user_lists"
}

// DailyActivity tracks user's daily practice activity for commitment chart
type DailyActivity struct {
	ActivityID     int       `json:"activity_id" gorm:"primaryKey;column:activity_id"`
	UserID         int       `json:"user_id" gorm:"column:user_id;not null;uniqueIndex:idx_user_date"`
	Date           time.Time `json:"date" gorm:"column:date;not null;uniqueIndex:idx_user_date;type:date"`
	ProblemsCount  int       `json:"problems_count" gorm:"column:problems_count;default:0"`
	QuestionsCount int       `json:"questions_count" gorm:"column:questions_count;default:0"`
	StudyTime      int       `json:"study_time_seconds" gorm:"column:study_time_seconds;default:0"`
	Streak         int       `json:"streak" gorm:"column:streak;default:0"`
}

func (DailyActivity) TableName() string {
	return "daily_activities"
}

// SpacedRepetitionReview tracks SM-2 algorithm data per question
type SpacedRepetitionReview struct {
	ReviewID       int        `json:"review_id" gorm:"primaryKey;column:review_id"`
	UserID         int        `json:"user_id" gorm:"column:user_id;not null;uniqueIndex:idx_user_question"`
	QuestionID     int        `json:"question_id" gorm:"column:question_id;not null;uniqueIndex:idx_user_question"`
	EasinessFactor float64    `json:"easiness_factor" gorm:"column:easiness_factor;default:2.5"`
	IntervalDays   int        `json:"interval_days" gorm:"column:interval_days;default:1"`
	Repetitions    int        `json:"repetitions" gorm:"column:repetitions;default:0"`
	NextReviewAt   time.Time  `json:"next_review_at" gorm:"column:next_review_at;not null"`
	LastReviewAt   *time.Time `json:"last_review_at,omitempty" gorm:"column:last_review_at"`
	QualityRating  *int       `json:"quality_rating,omitempty" gorm:"column:quality_rating"`
	CreatedAt      time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (SpacedRepetitionReview) TableName() string {
	return "spaced_repetition_reviews"
}

// ReviewQueue tracks questions scheduled for review
type ReviewQueue struct {
	QueueID      int       `json:"queue_id" gorm:"primaryKey;column:queue_id"`
	UserID       int       `json:"user_id" gorm:"column:user_id;not null;uniqueIndex:idx_user_question_queue"`
	QuestionID   int       `json:"question_id" gorm:"column:question_id;not null;uniqueIndex:idx_user_question_queue"`
	ScheduledFor time.Time `json:"scheduled_for" gorm:"column:scheduled_for;not null"`
	Priority     int       `json:"priority" gorm:"column:priority;default:0"`
	IsOverdue    bool      `json:"is_overdue" gorm:"column:is_overdue;default:false"`
	AddedAt      time.Time `json:"added_at" gorm:"column:added_at;autoCreateTime"`
}

func (ReviewQueue) TableName() string {
	return "review_queue"
}

// CodeSubmission tracks user code submissions for AI assessment
type CodeSubmission struct {
	SubmissionID    int        `json:"submission_id" gorm:"primaryKey;column:submission_id"`
	UserID          int        `json:"user_id" gorm:"column:user_id;not null;index"`
	ProblemID       int        `json:"problem_id" gorm:"column:problem_id;not null;index"`
	Code            string     `json:"code" gorm:"column:code;not null;type:text"`
	Language        string     `json:"language" gorm:"column:language;not null"`
	Status          string     `json:"status" gorm:"column:status;default:'pending'"`
	TestResults     JSONB      `json:"test_results,omitempty" gorm:"column:test_results;type:jsonb"`
	AIFeedback      JSONB      `json:"ai_feedback,omitempty" gorm:"column:ai_feedback;type:jsonb"`
	AIScore         *float64   `json:"ai_score,omitempty" gorm:"column:ai_score"`
	TimeComplexity  *string    `json:"time_complexity,omitempty" gorm:"column:time_complexity"`
	SpaceComplexity *string    `json:"space_complexity,omitempty" gorm:"column:space_complexity"`
	ExecutionTimeMs *int       `json:"execution_time_ms,omitempty" gorm:"column:execution_time_ms"`
	MemoryUsedKb    *int       `json:"memory_used_kb,omitempty" gorm:"column:memory_used_kb"`
	SubmittedAt     time.Time  `json:"submitted_at" gorm:"column:submitted_at;autoCreateTime"`
	EvaluatedAt     *time.Time `json:"evaluated_at,omitempty" gorm:"column:evaluated_at"`
}

func (CodeSubmission) TableName() string {
	return "code_submissions"
}

// QuestionHintUsage tracks which hints a user has seen
type QuestionHintUsage struct {
	UsageID    int       `json:"usage_id" gorm:"primaryKey;column:usage_id"`
	UserID     int       `json:"user_id" gorm:"column:user_id;not null;uniqueIndex:idx_user_question_hint"`
	QuestionID int       `json:"question_id" gorm:"column:question_id;not null;uniqueIndex:idx_user_question_hint"`
	HintLevel  int       `json:"hint_level" gorm:"column:hint_level;not null;uniqueIndex:idx_user_question_hint"`
	UsedAt     time.Time `json:"used_at" gorm:"column:used_at;autoCreateTime"`
}

func (QuestionHintUsage) TableName() string {
	return "question_hint_usage"
}

// QuestionWithHints extends Question to include hint fields
type QuestionWithHints struct {
	Question
	HintLevel1 *string `json:"hint_level_1,omitempty" gorm:"column:hint_level_1"`
	HintLevel2 *string `json:"hint_level_2,omitempty" gorm:"column:hint_level_2"`
	HintLevel3 *string `json:"hint_level_3,omitempty" gorm:"column:hint_level_3"`
}

func (QuestionWithHints) TableName() string {
	return "questions"
}

// AutoMigrate runs all model migrations
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Problem{},
		&Topic{},
		&ProblemTopic{},
		&Question{},
		&UserAttempt{},
		&UserSkill{},
		&TrainingPlan{},
		&TrainingPlanItem{},
		&Assessment{},
		&WeaknessAnalysis{},
		&UserList{},
		&DailyActivity{},
		&SpacedRepetitionReview{},
		&ReviewQueue{},
		&CodeSubmission{},
		&QuestionHintUsage{},
	)
}

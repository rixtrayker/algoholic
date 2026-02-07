package models

import (
	"encoding/json"
	"time"
)

// ──────────────────────────────────────────────
// Category
// ──────────────────────────────────────────────

type Category struct {
	CategoryID   int       `json:"category_id" db:"category_id"`
	Slug         string    `json:"slug" db:"slug"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	Icon         string    `json:"icon" db:"icon"`
	Color        string    `json:"color" db:"color"`
	DisplayOrder int       `json:"display_order" db:"display_order"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// ──────────────────────────────────────────────
// Topic
// ──────────────────────────────────────────────

type TopicLevel string

const (
	LevelCategory  TopicLevel = "category"
	LevelTopic     TopicLevel = "topic"
	LevelSubtopic  TopicLevel = "subtopic"
	LevelPattern   TopicLevel = "pattern"
	LevelVariation TopicLevel = "variation"
)

type Topic struct {
	TopicID       string     `json:"topic_id" db:"topic_id"`
	Name          string     `json:"name" db:"name"`
	Slug          string     `json:"slug" db:"slug"`
	CategoryID    int        `json:"category_id" db:"category_id"`
	ParentTopicID *string    `json:"parent_topic_id,omitempty" db:"parent_topic_id"`
	Level         TopicLevel `json:"level" db:"level"`
	Description   string     `json:"description" db:"description"`
	Keywords      []string   `json:"keywords" db:"keywords"`
	DifficultyMin float64    `json:"difficulty_min" db:"difficulty_min"`
	DifficultyMax float64    `json:"difficulty_max" db:"difficulty_max"`
	DisplayOrder  int        `json:"display_order" db:"display_order"`
	Metadata      JSONB      `json:"metadata" db:"metadata"`
}

// ──────────────────────────────────────────────
// Tag
// ──────────────────────────────────────────────

type Tag struct {
	TagID       int    `json:"tag_id" db:"tag_id"`
	Slug        string `json:"slug" db:"slug"`
	Name        string `json:"name" db:"name"`
	TagGroup    string `json:"tag_group" db:"tag_group"`
	Description string `json:"description" db:"description"`
	Color       string `json:"color" db:"color"`
}

// ──────────────────────────────────────────────
// Problem
// ──────────────────────────────────────────────

type DifficultyLabel string

const (
	DiffEasy   DifficultyLabel = "easy"
	DiffMedium DifficultyLabel = "medium"
	DiffHard   DifficultyLabel = "hard"
	DiffExpert DifficultyLabel = "expert"
)

type ProblemExample struct {
	Input       string `json:"input"`
	Output      string `json:"output"`
	Explanation string `json:"explanation,omitempty"`
}

type Problem struct {
	ProblemID         int              `json:"problem_id" db:"problem_id"`
	LeetcodeNumber    *int             `json:"leetcode_number,omitempty" db:"leetcode_number"`
	Slug              string           `json:"slug" db:"slug"`
	Title             string           `json:"title" db:"title"`
	Statement         string           `json:"statement" db:"statement"`
	Constraints       []string         `json:"constraints" db:"constraints"`
	Examples          []ProblemExample `json:"examples" db:"examples"`
	Hints             []string         `json:"hints" db:"hints"`
	DifficultyScore   float64          `json:"difficulty_score" db:"difficulty_score"`
	DifficultyLabel   DifficultyLabel  `json:"difficulty_label" db:"difficulty_label"`
	PrimaryPattern    string           `json:"primary_pattern" db:"primary_pattern"`
	SecondaryPatterns []string         `json:"secondary_patterns" db:"secondary_patterns"`
	TimeComplexity    string           `json:"time_complexity" db:"time_complexity"`
	SpaceComplexity   string           `json:"space_complexity" db:"space_complexity"`
	Frequency         float64          `json:"frequency" db:"frequency"`
	Companies         []string         `json:"companies" db:"companies"`
	HasFollowUps      bool             `json:"has_follow_ups" db:"has_follow_ups"`
	HasPrerequisites  bool             `json:"has_prerequisites" db:"has_prerequisites"`

	// Seed-time relations (not stored in this table directly)
	TopicLinks []ProblemTopicLink `json:"topic_links,omitempty"`
	TagSlugs   []string           `json:"tag_slugs,omitempty"`
}

type ProblemTopicLink struct {
	TopicID     string  `json:"topic_id"`
	Relevance   float64 `json:"relevance"`
	IsPrimary   bool    `json:"is_primary"`
	PatternUsed string  `json:"pattern_used,omitempty"`
	KeyInsight  string  `json:"key_insight,omitempty"`
}

// ──────────────────────────────────────────────
// Question Type (taxonomy entry)
// ──────────────────────────────────────────────

type QuestionFormat string

const (
	FmtMultipleChoice QuestionFormat = "multiple_choice"
	FmtCode           QuestionFormat = "code"
	FmtRanking        QuestionFormat = "ranking"
	FmtOpenEnded      QuestionFormat = "open_ended"
	FmtFillBlank      QuestionFormat = "fill_blank"
	FmtDebug          QuestionFormat = "debug"
)

type QuestionType struct {
	TypeID            int             `json:"type_id" db:"type_id"`
	Slug              string          `json:"slug" db:"slug"`
	CategoryCode      string          `json:"category_code" db:"category_code"`
	Name              string          `json:"name" db:"name"`
	Description       string          `json:"description" db:"description"`
	Format            QuestionFormat  `json:"format" db:"format"`
	ParentCategory    string          `json:"parent_category" db:"parent_category"`
	DifficultyDefault DifficultyLabel `json:"difficulty_default" db:"difficulty_default"`
	EstimatedTimeSec  int             `json:"estimated_time_sec" db:"estimated_time_sec"`
}

// ──────────────────────────────────────────────
// Question / Challenge
// ──────────────────────────────────────────────

type AnswerOption struct {
	ID        string `json:"id"`        // "A", "B", "C", "D"
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

type Question struct {
	QuestionID            int             `json:"question_id" db:"question_id"`
	QuestionTypeSlug      string          `json:"question_type_slug"`           // resolved at insert time
	Category              string          `json:"category" db:"category"`
	Subcategory           string          `json:"subcategory" db:"subcategory"`
	QuestionText          string          `json:"question_text" db:"question_text"`
	QuestionData          JSONB           `json:"question_data,omitempty" db:"question_data"`
	Format                QuestionFormat  `json:"format" db:"format"`
	CorrectAnswer         JSONB           `json:"correct_answer" db:"correct_answer"`
	AnswerOptions         []AnswerOption  `json:"answer_options" db:"answer_options"`
	WrongAnswerExplanations JSONB         `json:"wrong_answer_explanations,omitempty"`
	Explanation           string          `json:"explanation" db:"explanation"`
	DetailedSolution      string          `json:"detailed_solution,omitempty" db:"detailed_solution"`
	CommonMistakes        []string        `json:"common_mistakes,omitempty" db:"common_mistakes"`
	HintLevel1            string          `json:"hint_level_1,omitempty" db:"hint_level_1"`
	HintLevel2            string          `json:"hint_level_2,omitempty" db:"hint_level_2"`
	HintLevel3            string          `json:"hint_level_3,omitempty" db:"hint_level_3"`
	DifficultyScore       float64         `json:"difficulty_score" db:"difficulty_score"`
	DifficultyLabel       DifficultyLabel `json:"difficulty_label" db:"difficulty_label"`
	EstimatedTimeSec      int             `json:"estimated_time_sec" db:"estimated_time_sec"`
	RelatedProblemSlug    string          `json:"related_problem_slug,omitempty"` // resolved at insert
	RelatedTopicID        string          `json:"related_topic_id,omitempty" db:"related_topic_id"`
	Tags                  []string        `json:"tags" db:"tags"`
	Concepts              []string        `json:"concepts" db:"concepts"`
}

// ──────────────────────────────────────────────
// Pitfall
// ──────────────────────────────────────────────

type Pitfall struct {
	TopicID     string `json:"topic_id"`
	Description string `json:"description"`
	Example     string `json:"example"`
	Fix         string `json:"fix"`
	Severity    int    `json:"severity"`
}

// ──────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────

type JSONB map[string]interface{}

func (j JSONB) ToJSON() ([]byte, error) {
	return json.Marshal(j)
}

// SeedBundle is the complete input for one seed run.
type SeedBundle struct {
	Categories    []Category     `json:"categories"`
	Topics        []Topic        `json:"topics"`
	Tags          []Tag          `json:"tags"`
	Problems      []Problem      `json:"problems"`
	QuestionTypes []QuestionType `json:"question_types"`
	Questions     []Question     `json:"questions"`
	Pitfalls      []Pitfall      `json:"pitfalls"`
}

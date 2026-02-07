package seed_test

import (
	"testing"

	"dsa-platform/internal/seed/data"
)

// ─── Categories ────────────────────────────────────────────

func TestCategories_NoDuplicateSlugs(t *testing.T) {
	seen := map[string]bool{}
	for _, c := range data.SeedCategories() {
		if seen[c.Slug] {
			t.Errorf("duplicate category slug: %s", c.Slug)
		}
		seen[c.Slug] = true
	}
}

func TestCategories_AllHaveRequiredFields(t *testing.T) {
	for _, c := range data.SeedCategories() {
		if c.Slug == "" || c.Name == "" || c.Description == "" {
			t.Errorf("category %q missing required fields", c.Slug)
		}
	}
}

// ─── Topics ────────────────────────────────────────────────

func TestTopics_NoDuplicateIDs(t *testing.T) {
	seen := map[string]bool{}
	for _, topic := range data.SeedTopics() {
		if seen[topic.TopicID] {
			t.Errorf("duplicate topic_id: %s", topic.TopicID)
		}
		seen[topic.TopicID] = true
	}
}

func TestTopics_ParentReferencesExist(t *testing.T) {
	ids := map[string]bool{}
	for _, topic := range data.SeedTopics() {
		ids[topic.TopicID] = true
	}
	for _, topic := range data.SeedTopics() {
		if topic.ParentTopicID != nil && !ids[*topic.ParentTopicID] {
			t.Errorf("topic %s references non-existent parent %s", topic.TopicID, *topic.ParentTopicID)
		}
	}
}

func TestTopics_DifficultyRangeValid(t *testing.T) {
	for _, topic := range data.SeedTopics() {
		if topic.DifficultyMin > topic.DifficultyMax {
			t.Errorf("topic %s: min (%f) > max (%f)", topic.TopicID, topic.DifficultyMin, topic.DifficultyMax)
		}
		if topic.DifficultyMin < 0 || topic.DifficultyMax > 100 {
			t.Errorf("topic %s: difficulty out of [0,100] range", topic.TopicID)
		}
	}
}

func TestTopics_CategorySlugReferencesExist(t *testing.T) {
	catSlugs := map[string]bool{}
	for _, c := range data.SeedCategories() {
		catSlugs[c.Slug] = true
	}
	for _, topic := range data.SeedTopics() {
		slug, ok := topic.Metadata["category_slug"].(string)
		if ok && !catSlugs[slug] {
			t.Errorf("topic %s references non-existent category slug %q", topic.TopicID, slug)
		}
	}
}

// ─── Tags ──────────────────────────────────────────────────

func TestTags_NoDuplicateSlugs(t *testing.T) {
	seen := map[string]bool{}
	for _, tag := range data.SeedTags() {
		if seen[tag.Slug] {
			t.Errorf("duplicate tag slug: %s", tag.Slug)
		}
		seen[tag.Slug] = true
	}
}

func TestTags_AllHaveGroup(t *testing.T) {
	validGroups := map[string]bool{
		"pattern": true, "data_structure": true, "difficulty_trait": true,
		"company": true, "technique": true,
	}
	for _, tag := range data.SeedTags() {
		if !validGroups[tag.TagGroup] {
			t.Errorf("tag %s has invalid group %q", tag.Slug, tag.TagGroup)
		}
	}
}

// ─── Problems ──────────────────────────────────────────────

func TestProblems_NoDuplicateSlugs(t *testing.T) {
	seen := map[string]bool{}
	for _, p := range data.SeedProblems() {
		if seen[p.Slug] {
			t.Errorf("duplicate problem slug: %s", p.Slug)
		}
		seen[p.Slug] = true
	}
}

func TestProblems_DifficultyScoreValid(t *testing.T) {
	for _, p := range data.SeedProblems() {
		if p.DifficultyScore < 0 || p.DifficultyScore > 100 {
			t.Errorf("problem %s: difficulty %f out of range", p.Slug, p.DifficultyScore)
		}
	}
}

func TestProblems_TopicLinksReferenceExistingTopics(t *testing.T) {
	topicIDs := map[string]bool{}
	for _, topic := range data.SeedTopics() {
		topicIDs[topic.TopicID] = true
	}
	for _, p := range data.SeedProblems() {
		for _, link := range p.TopicLinks {
			if !topicIDs[link.TopicID] {
				t.Errorf("problem %s links to non-existent topic %s", p.Slug, link.TopicID)
			}
		}
	}
}

func TestProblems_TagSlugsReferenceExistingTags(t *testing.T) {
	tagSlugs := map[string]bool{}
	for _, tag := range data.SeedTags() {
		tagSlugs[tag.Slug] = true
	}
	for _, p := range data.SeedProblems() {
		for _, slug := range p.TagSlugs {
			if !tagSlugs[slug] {
				t.Errorf("problem %s references non-existent tag %q", p.Slug, slug)
			}
		}
	}
}

func TestProblems_AtLeastOnePrimaryTopic(t *testing.T) {
	for _, p := range data.SeedProblems() {
		hasPrimary := false
		for _, link := range p.TopicLinks {
			if link.IsPrimary {
				hasPrimary = true
				break
			}
		}
		if !hasPrimary {
			t.Errorf("problem %s has no primary topic link", p.Slug)
		}
	}
}

// ─── Question Types ────────────────────────────────────────

func TestQuestionTypes_NoDuplicateSlugs(t *testing.T) {
	seen := map[string]bool{}
	for _, qt := range data.SeedQuestionTypes() {
		if seen[qt.Slug] {
			t.Errorf("duplicate question type slug: %s", qt.Slug)
		}
		seen[qt.Slug] = true
	}
}

// ─── Questions ─────────────────────────────────────────────

func TestQuestions_TypeSlugReferencesExist(t *testing.T) {
	qtSlugs := map[string]bool{}
	for _, qt := range data.SeedQuestionTypes() {
		qtSlugs[qt.Slug] = true
	}
	for i, q := range data.SeedQuestions() {
		if !qtSlugs[q.QuestionTypeSlug] {
			t.Errorf("question %d references non-existent type %q", i, q.QuestionTypeSlug)
		}
	}
}

func TestQuestions_ProblemSlugReferencesExist(t *testing.T) {
	problemSlugs := map[string]bool{}
	for _, p := range data.SeedProblems() {
		problemSlugs[p.Slug] = true
	}
	for i, q := range data.SeedQuestions() {
		if q.RelatedProblemSlug != "" && !problemSlugs[q.RelatedProblemSlug] {
			t.Errorf("question %d references non-existent problem %q", i, q.RelatedProblemSlug)
		}
	}
}

func TestQuestions_TopicReferencesExist(t *testing.T) {
	topicIDs := map[string]bool{}
	for _, topic := range data.SeedTopics() {
		topicIDs[topic.TopicID] = true
	}
	for i, q := range data.SeedQuestions() {
		if q.RelatedTopicID != "" && !topicIDs[q.RelatedTopicID] {
			t.Errorf("question %d references non-existent topic %q", i, q.RelatedTopicID)
		}
	}
}

func TestQuestions_AllHaveExplanation(t *testing.T) {
	for i, q := range data.SeedQuestions() {
		if q.Explanation == "" {
			t.Errorf("question %d (%s) missing explanation", i, q.Category)
		}
	}
}

func TestQuestions_DifficultyValid(t *testing.T) {
	for i, q := range data.SeedQuestions() {
		if q.DifficultyScore < 0 || q.DifficultyScore > 100 {
			t.Errorf("question %d: difficulty %f out of range", i, q.DifficultyScore)
		}
	}
}

func TestQuestions_CoverAllCategories(t *testing.T) {
	categories := map[string]int{}
	for _, q := range data.SeedQuestions() {
		categories[q.Category]++
	}

	expected := []string{
		"complexity_analysis", "ds_selection", "pattern_recognition",
		"edge_cases", "code_templates", "implementation",
		"bug_detection", "tradeoffs", "hybrid",
	}
	for _, cat := range expected {
		if categories[cat] == 0 {
			t.Errorf("no questions found for category %q", cat)
		}
	}
}

// ─── Pitfalls ──────────────────────────────────────────────

func TestPitfalls_TopicReferencesExist(t *testing.T) {
	topicIDs := map[string]bool{}
	for _, topic := range data.SeedTopics() {
		topicIDs[topic.TopicID] = true
	}
	for _, p := range data.SeedPitfalls() {
		if !topicIDs[p.TopicID] {
			t.Errorf("pitfall references non-existent topic %q", p.TopicID)
		}
	}
}

func TestPitfalls_SeverityInRange(t *testing.T) {
	for _, p := range data.SeedPitfalls() {
		if p.Severity < 1 || p.Severity > 5 {
			t.Errorf("pitfall for %s has severity %d (must be 1-5)", p.TopicID, p.Severity)
		}
	}
}

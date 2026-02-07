package seed

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"dsa-platform/internal/seed/data"
	"dsa-platform/pkg/models"

	_ "github.com/lib/pq"
)

// Seeder holds the DB connection and runs all seed operations.
type Seeder struct {
	db  *sql.DB
	ctx context.Context

	// caches for resolving slugs → IDs during insertion
	categoryIDBySlug    map[string]int
	problemIDBySlug     map[string]int
	tagIDBySlug         map[string]int
	questionTypeIDBySlug map[string]int
}

// New creates a Seeder connected to the given database.
func New(ctx context.Context, db *sql.DB) *Seeder {
	return &Seeder{
		db:                   db,
		ctx:                  ctx,
		categoryIDBySlug:     make(map[string]int),
		problemIDBySlug:      make(map[string]int),
		tagIDBySlug:          make(map[string]int),
		questionTypeIDBySlug: make(map[string]int),
	}
}

// RunAll executes the full seed pipeline in dependency order.
func (s *Seeder) RunAll() error {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"migrate", s.runMigrations},
		{"categories", s.seedCategories},
		{"topics", s.seedTopics},
		{"tags", s.seedTags},
		{"problems", s.seedProblems},
		{"problem_topics", s.seedProblemTopics},
		{"problem_tags", s.seedProblemTags},
		{"question_types", s.seedQuestionTypes},
		{"questions", s.seedQuestions},
		{"pitfalls", s.seedPitfalls},
		{"graph_nodes", s.seedGraphNodes},
		{"graph_edges", s.seedGraphEdges},
	}

	for _, step := range steps {
		log.Printf("🌱 Seeding: %s ...", step.name)
		if err := step.fn(); err != nil {
			return fmt.Errorf("seed %s failed: %w", step.name, err)
		}
		log.Printf("   ✅ %s done", step.name)
	}

	log.Println("🎉 All seeds completed successfully!")
	return nil
}

// ─── Migrations ────────────────────────────────────────────

func (s *Seeder) runMigrations() error {
	sqlBytes, err := os.ReadFile("internal/seed/migrations/001_core_schema.sql")
	if err != nil {
		return fmt.Errorf("read migration: %w", err)
	}
	_, err = s.db.ExecContext(s.ctx, string(sqlBytes))
	return err
}

// ─── Categories ────────────────────────────────────────────

func (s *Seeder) seedCategories() error {
	for _, c := range data.SeedCategories() {
		var id int
		err := s.db.QueryRowContext(s.ctx,
			`INSERT INTO categories (slug, name, description, icon, color, display_order)
			 VALUES ($1, $2, $3, $4, $5, $6)
			 ON CONFLICT (slug) DO UPDATE SET name=EXCLUDED.name
			 RETURNING category_id`,
			c.Slug, c.Name, c.Description, c.Icon, c.Color, c.DisplayOrder,
		).Scan(&id)
		if err != nil {
			return fmt.Errorf("insert category %s: %w", c.Slug, err)
		}
		s.categoryIDBySlug[c.Slug] = id
	}
	return nil
}

// ─── Topics ────────────────────────────────────────────────

func (s *Seeder) seedTopics() error {
	topics := data.SeedTopics()

	// Two-pass: first insert topics without parents, then update parents.
	// This avoids FK violation on forward references.
	for _, t := range topics {
		catSlug, _ := t.Metadata["category_slug"].(string)
		catID := s.categoryIDBySlug[catSlug]

		_, err := s.db.ExecContext(s.ctx,
			`INSERT INTO topics (topic_id, name, slug, category_id, level, description, keywords, difficulty_min, difficulty_max, display_order, metadata)
			 VALUES ($1, $2, $3, $4, $5::topic_level, $6, $7, $8, $9, $10, $11)
			 ON CONFLICT (topic_id) DO UPDATE SET name=EXCLUDED.name`,
			t.TopicID, t.Name, t.Slug, catID, string(t.Level), t.Description,
			pqArray(t.Keywords), t.DifficultyMin, t.DifficultyMax, t.DisplayOrder,
			jsonMust(t.Metadata),
		)
		if err != nil {
			return fmt.Errorf("insert topic %s: %w", t.TopicID, err)
		}
	}

	// Second pass: set parent references
	for _, t := range topics {
		if t.ParentTopicID != nil {
			_, err := s.db.ExecContext(s.ctx,
				`UPDATE topics SET parent_topic_id = $1 WHERE topic_id = $2`,
				*t.ParentTopicID, t.TopicID,
			)
			if err != nil {
				return fmt.Errorf("set parent for %s: %w", t.TopicID, err)
			}
		}
	}
	return nil
}

// ─── Tags ──────────────────────────────────────────────────

func (s *Seeder) seedTags() error {
	for _, t := range data.SeedTags() {
		var id int
		err := s.db.QueryRowContext(s.ctx,
			`INSERT INTO tags (slug, name, tag_group, description, color)
			 VALUES ($1, $2, $3, $4, $5)
			 ON CONFLICT (slug) DO UPDATE SET name=EXCLUDED.name
			 RETURNING tag_id`,
			t.Slug, t.Name, t.TagGroup, t.Description, t.Color,
		).Scan(&id)
		if err != nil {
			return fmt.Errorf("insert tag %s: %w", t.Slug, err)
		}
		s.tagIDBySlug[t.Slug] = id
	}
	return nil
}

// ─── Problems ──────────────────────────────────────────────

func (s *Seeder) seedProblems() error {
	for _, p := range data.SeedProblems() {
		var id int
		err := s.db.QueryRowContext(s.ctx,
			`INSERT INTO problems (
				leetcode_number, slug, title, statement, constraints, examples, hints,
				difficulty_score, difficulty_label,
				primary_pattern, secondary_patterns,
				time_complexity, space_complexity,
				frequency, companies, has_follow_ups, has_prerequisites
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9::difficulty_label,$10,$11,$12,$13,$14,$15,$16,$17)
			ON CONFLICT (slug) DO UPDATE SET title=EXCLUDED.title
			RETURNING problem_id`,
			p.LeetcodeNumber, p.Slug, p.Title, p.Statement,
			pqArray(p.Constraints), jsonMust(p.Examples), pqArray(p.Hints),
			p.DifficultyScore, string(p.DifficultyLabel),
			p.PrimaryPattern, pqArray(p.SecondaryPatterns),
			p.TimeComplexity, p.SpaceComplexity,
			p.Frequency, jsonMust(p.Companies),
			p.HasFollowUps, p.HasPrerequisites,
		).Scan(&id)
		if err != nil {
			return fmt.Errorf("insert problem %s: %w", p.Slug, err)
		}
		s.problemIDBySlug[p.Slug] = id
	}
	return nil
}

// ─── Problem <-> Topic junction ────────────────────────────

func (s *Seeder) seedProblemTopics() error {
	for _, p := range data.SeedProblems() {
		pid, ok := s.problemIDBySlug[p.Slug]
		if !ok {
			continue
		}
		for _, link := range p.TopicLinks {
			_, err := s.db.ExecContext(s.ctx,
				`INSERT INTO problem_topics (problem_id, topic_id, relevance, is_primary, pattern_used, key_insight)
				 VALUES ($1, $2, $3, $4, $5, $6)
				 ON CONFLICT DO NOTHING`,
				pid, link.TopicID, link.Relevance, link.IsPrimary, link.PatternUsed, link.KeyInsight,
			)
			if err != nil {
				return fmt.Errorf("link problem %s → topic %s: %w", p.Slug, link.TopicID, err)
			}
		}
	}
	return nil
}

// ─── Problem <-> Tag junction ──────────────────────────────

func (s *Seeder) seedProblemTags() error {
	for _, p := range data.SeedProblems() {
		pid, ok := s.problemIDBySlug[p.Slug]
		if !ok {
			continue
		}
		for _, tagSlug := range p.TagSlugs {
			tid, ok := s.tagIDBySlug[tagSlug]
			if !ok {
				log.Printf("   ⚠️  tag slug %q not found, skipping", tagSlug)
				continue
			}
			_, err := s.db.ExecContext(s.ctx,
				`INSERT INTO problem_tags (problem_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
				pid, tid,
			)
			if err != nil {
				return fmt.Errorf("tag problem %s with %s: %w", p.Slug, tagSlug, err)
			}
		}
	}
	return nil
}

// ─── Question Types ────────────────────────────────────────

func (s *Seeder) seedQuestionTypes() error {
	for _, qt := range data.SeedQuestionTypes() {
		var id int
		err := s.db.QueryRowContext(s.ctx,
			`INSERT INTO question_types (slug, category_code, name, description, format, parent_category, difficulty_default, estimated_time_sec)
			 VALUES ($1, $2, $3, $4, $5::question_format, $6, $7::difficulty_label, $8)
			 ON CONFLICT (slug) DO UPDATE SET name=EXCLUDED.name
			 RETURNING type_id`,
			qt.Slug, qt.CategoryCode, qt.Name, qt.Description,
			string(qt.Format), qt.ParentCategory,
			string(qt.DifficultyDefault), qt.EstimatedTimeSec,
		).Scan(&id)
		if err != nil {
			return fmt.Errorf("insert question type %s: %w", qt.Slug, err)
		}
		s.questionTypeIDBySlug[qt.Slug] = id
	}
	return nil
}

// ─── Questions / Challenges ────────────────────────────────

func (s *Seeder) seedQuestions() error {
	for i, q := range data.SeedQuestions() {
		qtID, ok := s.questionTypeIDBySlug[q.QuestionTypeSlug]
		if !ok {
			log.Printf("   ⚠️  question type %q not found for question %d, skipping", q.QuestionTypeSlug, i)
			continue
		}

		// Resolve optional problem reference
		var problemID *int
		if q.RelatedProblemSlug != "" {
			if pid, ok := s.problemIDBySlug[q.RelatedProblemSlug]; ok {
				problemID = &pid
			}
		}

		// Resolve optional topic reference
		var topicID *string
		if q.RelatedTopicID != "" {
			topicID = &q.RelatedTopicID
		}

		_, err := s.db.ExecContext(s.ctx,
			`INSERT INTO questions (
				question_type_id, category, subcategory,
				question_text, question_data, format,
				correct_answer, answer_options, wrong_answer_explanations,
				explanation, detailed_solution, common_mistakes,
				hint_level_1, hint_level_2, hint_level_3,
				difficulty_score, difficulty_label, estimated_time_sec,
				related_problem_id, related_topic_id,
				tags, concepts
			) VALUES (
				$1, $2, $3,
				$4, $5, $6::question_format,
				$7, $8, $9,
				$10, $11, $12,
				$13, $14, $15,
				$16, $17::difficulty_label, $18,
				$19, $20,
				$21, $22
			)`,
			qtID, q.Category, q.Subcategory,
			q.QuestionText, jsonMustNullable(q.QuestionData), string(q.Format),
			jsonMust(q.CorrectAnswer), jsonMust(q.AnswerOptions), jsonMustNullable(q.WrongAnswerExplanations),
			q.Explanation, nilStr(q.DetailedSolution), pqArray(q.CommonMistakes),
			nilStr(q.HintLevel1), nilStr(q.HintLevel2), nilStr(q.HintLevel3),
			q.DifficultyScore, string(q.DifficultyLabel), q.EstimatedTimeSec,
			problemID, topicID,
			pqArray(q.Tags), pqArray(q.Concepts),
		)
		if err != nil {
			return fmt.Errorf("insert question %d (%s): %w", i, q.Category, err)
		}
	}
	return nil
}

// ─── Pitfalls ──────────────────────────────────────────────

func (s *Seeder) seedPitfalls() error {
	for _, p := range data.SeedPitfalls() {
		_, err := s.db.ExecContext(s.ctx,
			`INSERT INTO pitfalls (topic_id, description, example, fix, severity)
			 VALUES ($1, $2, $3, $4, $5)`,
			p.TopicID, p.Description, p.Example, p.Fix, p.Severity,
		)
		if err != nil {
			return fmt.Errorf("insert pitfall for %s: %w", p.TopicID, err)
		}
	}
	return nil
}

// ─── Apache AGE Graph Nodes ────────────────────────────────

func (s *Seeder) seedGraphNodes() error {
	// Create Topic nodes in the graph
	for _, t := range data.SeedTopics() {
		cypher := fmt.Sprintf(
			`SELECT * FROM cypher('dsa_graph', $$
				MERGE (t:Topic {topic_id: '%s'})
				SET t.name = '%s', t.level = '%s', t.difficulty_min = %f, t.difficulty_max = %f
			$$) AS (v agtype)`,
			escAGE(t.TopicID), escAGE(t.Name), escAGE(string(t.Level)),
			t.DifficultyMin, t.DifficultyMax,
		)
		if _, err := s.db.ExecContext(s.ctx, cypher); err != nil {
			// AGE might not be installed; log and continue
			log.Printf("   ⚠️  AGE topic node %s: %v (skipping graph)", t.TopicID, err)
			return nil // don't fail seed if AGE is missing
		}
	}

	// Create Problem nodes
	for _, p := range data.SeedProblems() {
		lc := 0
		if p.LeetcodeNumber != nil {
			lc = *p.LeetcodeNumber
		}
		cypher := fmt.Sprintf(
			`SELECT * FROM cypher('dsa_graph', $$
				MERGE (p:Problem {slug: '%s'})
				SET p.title = '%s', p.leetcode_number = %d, p.difficulty = %f
			$$) AS (v agtype)`,
			escAGE(p.Slug), escAGE(p.Title), lc, p.DifficultyScore,
		)
		if _, err := s.db.ExecContext(s.ctx, cypher); err != nil {
			log.Printf("   ⚠️  AGE problem node %s: %v", p.Slug, err)
			return nil
		}
	}
	return nil
}

// ─── Apache AGE Graph Edges ────────────────────────────────

func (s *Seeder) seedGraphEdges() error {
	// Topic → Topic prerequisite edges
	prereqEdges := [][2]string{
		{"arrays_basics", "two_pointers"},
		{"arrays_basics", "sliding_window"},
		{"arrays_basics", "prefix_sum"},
		{"dp_fundamentals", "dp_1d"},
		{"dp_1d", "dp_2d"},
		{"dp_1d", "dp_knapsack"},
		{"dp_2d", "dp_interval"},
		{"graph_representation", "graph_bfs"},
		{"graph_representation", "graph_dfs"},
		{"graph_dfs", "topological_sort"},
		{"graph_bfs", "shortest_path"},
		{"graph_representation", "union_find"},
		{"binary_tree", "bst"},
		{"binary_tree", "lca"},
		{"stack", "monotonic_stack"},
		{"queue", "monotonic_queue"},
		{"hash_map", "frequency_counting"},
		{"binary_search", "bs_on_answer"},
		{"bit_basics", "bit_tricks"},
	}

	for _, edge := range prereqEdges {
		cypher := fmt.Sprintf(
			`SELECT * FROM cypher('dsa_graph', $$
				MATCH (from:Topic {topic_id: '%s'}), (to:Topic {topic_id: '%s'})
				MERGE (from)-[:PREREQUISITE_OF]->(to)
			$$) AS (v agtype)`,
			escAGE(edge[0]), escAGE(edge[1]),
		)
		if _, err := s.db.ExecContext(s.ctx, cypher); err != nil {
			log.Printf("   ⚠️  AGE edge %s→%s: %v", edge[0], edge[1], err)
			return nil
		}
	}

	// Problem → Topic HAS_TOPIC edges
	for _, p := range data.SeedProblems() {
		for _, link := range p.TopicLinks {
			rel := "HAS_TOPIC"
			if link.IsPrimary {
				rel = "PRIMARY_TOPIC"
			}
			cypher := fmt.Sprintf(
				`SELECT * FROM cypher('dsa_graph', $$
					MATCH (p:Problem {slug: '%s'}), (t:Topic {topic_id: '%s'})
					MERGE (p)-[:%s {relevance: %f}]->(t)
				$$) AS (v agtype)`,
				escAGE(p.Slug), escAGE(link.TopicID), rel, link.Relevance,
			)
			if _, err := s.db.ExecContext(s.ctx, cypher); err != nil {
				log.Printf("   ⚠️  AGE edge %s→%s: %v", p.Slug, link.TopicID, err)
				return nil
			}
		}
	}

	// Related topic edges
	relatedEdges := [][2]string{
		{"two_pointers", "sliding_window"},
		{"two_pointers", "binary_search"},
		{"sliding_window", "prefix_sum"},
		{"graph_bfs", "graph_dfs"},
		{"dp_1d", "greedy"},
		{"union_find", "graph_dfs"},
		{"monotonic_stack", "monotonic_queue"},
		{"binary_search", "sorting_basics"},
		{"priority_queue", "sorting_basics"},
	}

	for _, edge := range relatedEdges {
		cypher := fmt.Sprintf(
			`SELECT * FROM cypher('dsa_graph', $$
				MATCH (a:Topic {topic_id: '%s'}), (b:Topic {topic_id: '%s'})
				MERGE (a)-[:RELATED_TO]->(b)
				MERGE (b)-[:RELATED_TO]->(a)
			$$) AS (v agtype)`,
			escAGE(edge[0]), escAGE(edge[1]),
		)
		if _, err := s.db.ExecContext(s.ctx, cypher); err != nil {
			log.Printf("   ⚠️  AGE related edge: %v", err)
			return nil
		}
	}

	return nil
}

// ─── Helpers ───────────────────────────────────────────────

func jsonMust(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(b)
}

func jsonMustNullable(v models.JSONB) *string {
	if v == nil || len(v) == 0 {
		return nil
	}
	s := jsonMust(v)
	return &s
}

func pqArray(ss []string) string {
	if len(ss) == 0 {
		return "{}"
	}
	escaped := make([]string, len(ss))
	for i, s := range ss {
		escaped[i] = `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
	}
	return "{" + strings.Join(escaped, ",") + "}"
}

func nilStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func escAGE(s string) string {
	return strings.ReplaceAll(s, "'", "\\'")
}

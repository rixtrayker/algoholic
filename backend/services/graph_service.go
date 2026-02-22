package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// GraphService runs Cypher queries via Apache AGE on top of PostgreSQL.
// All methods gracefully return SQL-based results when AGE is not installed.
type GraphService struct {
	db    *gorm.DB
	sqlDB *sql.DB
	graph string
}

const ageGraph = "algoholic_graph"

// NewGraphService creates a new graph service and bootstraps the AGE graph.
func NewGraphService(db *gorm.DB) *GraphService {
	sqlDB, _ := db.DB()
	gs := &GraphService{db: db, sqlDB: sqlDB, graph: ageGraph}
	gs.bootstrap()
	return gs
}

// bootstrap loads the AGE extension and creates the graph. Errors are silent.
func (gs *GraphService) bootstrap() {
	if !gs.IsAvailable() {
		return
	}
	gs.sqlDB.Exec(`LOAD 'age'`)
	gs.sqlDB.Exec(`SET search_path = ag_catalog, "$user", public`)
	gs.sqlDB.Exec(fmt.Sprintf(`SELECT * FROM create_graph('%s')`, gs.graph))
}

// IsAvailable checks whether AGE extension is installed.
func (gs *GraphService) IsAvailable() bool {
	var count int
	err := gs.sqlDB.QueryRow(`SELECT COUNT(*) FROM pg_extension WHERE extname = 'age'`).Scan(&count)
	return err == nil && count > 0
}

// ── Internal helpers ─────────────────────────────────────────────────────────

func escAGE(s string) string {
	return strings.ReplaceAll(s, "'", "\\'")
}

func (gs *GraphService) execCypher(cypher string) error {
	stmt := fmt.Sprintf(`SELECT * FROM cypher('%s', $$ %s $$) AS (v agtype)`, gs.graph, cypher)
	_, err := gs.sqlDB.Exec(stmt)
	return err
}

// ── Node management ───────────────────────────────────────────────────────────

func (gs *GraphService) UpsertProblemNode(id int, title string, difficulty float64, pattern string) error {
	return gs.execCypher(fmt.Sprintf(
		`MERGE (p:Problem {id: %d}) SET p.title = '%s', p.difficulty = %f, p.pattern = '%s'`,
		id, escAGE(title), difficulty, escAGE(pattern),
	))
}

func (gs *GraphService) UpsertTopicNode(id int, name, slug string, level int) error {
	return gs.execCypher(fmt.Sprintf(
		`MERGE (t:Topic {id: %d}) SET t.name = '%s', t.slug = '%s', t.level = %d`,
		id, escAGE(name), escAGE(slug), level,
	))
}

// ── Relationship management ───────────────────────────────────────────────────

func (gs *GraphService) LinkProblemToTopic(problemID, topicID int, relevance float64, isPrimary bool) error {
	rel := "HAS_TOPIC"
	if isPrimary {
		rel = "PRIMARY_TOPIC"
	}
	return gs.execCypher(fmt.Sprintf(
		`MATCH (p:Problem {id: %d}), (t:Topic {id: %d}) MERGE (p)-[r:%s]->(t) SET r.relevance = %f`,
		problemID, topicID, rel, relevance,
	))
}

func (gs *GraphService) LinkTopicPrerequisite(fromID, toID int) error {
	return gs.execCypher(fmt.Sprintf(
		`MATCH (a:Topic {id: %d}), (b:Topic {id: %d}) MERGE (a)-[:PREREQUISITE_OF]->(b)`,
		fromID, toID,
	))
}

func (gs *GraphService) LinkSimilarProblems(id1, id2 int, score float64) error {
	return gs.execCypher(fmt.Sprintf(
		`MATCH (a:Problem {id: %d}), (b:Problem {id: %d}) MERGE (a)-[r:SIMILAR_TO]->(b) SET r.score = %f`,
		id1, id2, score,
	))
}

// ── Queries ───────────────────────────────────────────────────────────────────

type SimilarProblemsResult struct {
	ProblemID  int     `json:"problem_id"`
	Title      string  `json:"title"`
	Difficulty float64 `json:"difficulty"`
	Score      float64 `json:"score"`
}

// FindSimilarProblems returns problems connected by SIMILAR_TO.
// Falls back to SQL (same primary_pattern) when AGE is unavailable.
func (gs *GraphService) FindSimilarProblems(problemID, limit int) ([]SimilarProblemsResult, error) {
	if !gs.IsAvailable() {
		return gs.sqlSimilarProblems(problemID, limit)
	}

	cypher := fmt.Sprintf(
		`MATCH (p:Problem {id: %d})-[r:SIMILAR_TO]-(s:Problem)
		 RETURN s.id AS id, s.title AS title, s.difficulty AS difficulty, r.score AS score
		 ORDER BY r.score DESC LIMIT %d`,
		problemID, limit,
	)
	stmt := fmt.Sprintf(
		`SELECT * FROM cypher('%s', $$ %s $$) AS (id agtype, title agtype, difficulty agtype, score agtype)`,
		gs.graph, cypher,
	)
	rows, err := gs.sqlDB.Query(stmt)
	if err != nil {
		return gs.sqlSimilarProblems(problemID, limit)
	}
	defer rows.Close()

	var results []SimilarProblemsResult
	for rows.Next() {
		var idRaw, titleRaw, diffRaw, scoreRaw string
		if err := rows.Scan(&idRaw, &titleRaw, &diffRaw, &scoreRaw); err != nil {
			continue
		}
		var r SimilarProblemsResult
		json.Unmarshal([]byte(idRaw), &r.ProblemID)
		json.Unmarshal([]byte(titleRaw), &r.Title)
		json.Unmarshal([]byte(diffRaw), &r.Difficulty)
		json.Unmarshal([]byte(scoreRaw), &r.Score)
		results = append(results, r)
	}

	if len(results) == 0 {
		return gs.sqlSimilarProblems(problemID, limit)
	}
	return results, nil
}

func (gs *GraphService) sqlSimilarProblems(problemID, limit int) ([]SimilarProblemsResult, error) {
	var pattern *string
	gs.db.Table("problems").Select("primary_pattern").Where("problem_id = ?", problemID).Scan(&pattern)

	type row struct {
		ProblemID       int
		Title           string
		DifficultyScore float64
	}
	var rows []row
	q := gs.db.Table("problems").
		Select("problem_id, title, difficulty_score").
		Where("problem_id != ?", problemID)
	if pattern != nil {
		q = q.Where("primary_pattern = ?", *pattern)
	}
	q.Order("difficulty_score ASC").Limit(limit).Scan(&rows)

	results := make([]SimilarProblemsResult, len(rows))
	for i, r := range rows {
		results[i] = SimilarProblemsResult{ProblemID: r.ProblemID, Title: r.Title, Difficulty: r.DifficultyScore, Score: 0.7}
	}
	return results, nil
}

// ─────────────────────────────────────────────────────────────────────────────

type LearningPathResult struct {
	TopicID int    `json:"topic_id"`
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	Step    int    `json:"step"`
}

// GetLearningPath returns the shortest prerequisite chain between two topics.
func (gs *GraphService) GetLearningPath(startTopicID, endTopicID int) ([]LearningPathResult, error) {
	if !gs.IsAvailable() {
		return gs.sqlLearningPath(startTopicID, endTopicID)
	}

	cypher := fmt.Sprintf(
		`MATCH path = shortestPath((s:Topic {id: %d})-[:PREREQUISITE_OF*]-(e:Topic {id: %d}))
		 UNWIND nodes(path) AS n RETURN n.id AS id, n.name AS name, n.slug AS slug`,
		startTopicID, endTopicID,
	)
	stmt := fmt.Sprintf(
		`SELECT * FROM cypher('%s', $$ %s $$) AS (id agtype, name agtype, slug agtype)`,
		gs.graph, cypher,
	)
	rows, err := gs.sqlDB.Query(stmt)
	if err != nil {
		return gs.sqlLearningPath(startTopicID, endTopicID)
	}
	defer rows.Close()

	var results []LearningPathResult
	step := 1
	for rows.Next() {
		var idRaw, nameRaw, slugRaw string
		if err := rows.Scan(&idRaw, &nameRaw, &slugRaw); err != nil {
			continue
		}
		var r LearningPathResult
		r.Step = step
		json.Unmarshal([]byte(idRaw), &r.TopicID)
		json.Unmarshal([]byte(nameRaw), &r.Name)
		json.Unmarshal([]byte(slugRaw), &r.Slug)
		results = append(results, r)
		step++
	}

	if len(results) == 0 {
		return gs.sqlLearningPath(startTopicID, endTopicID)
	}
	return results, nil
}

func (gs *GraphService) sqlLearningPath(startTopicID, endTopicID int) ([]LearningPathResult, error) {
	type topicRow struct {
		TopicID int
		Name    string
		Slug    string
	}
	var topics []topicRow
	gs.db.Table("topics").Select("topic_id, name, slug").
		Where("topic_id IN ?", []int{startTopicID, endTopicID}).Scan(&topics)

	results := make([]LearningPathResult, len(topics))
	for i, t := range topics {
		results[i] = LearningPathResult{TopicID: t.TopicID, Name: t.Name, Slug: t.Slug, Step: i + 1}
	}
	return results, nil
}

// ─────────────────────────────────────────────────────────────────────────────

type TopicPrerequisitesResult struct {
	TopicID       int    `json:"topic_id"`
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	DifficultyLvl int    `json:"difficulty_level"`
}

// GetTopicPrerequisites returns prerequisite topics for a given topic.
func (gs *GraphService) GetTopicPrerequisites(topicID int) ([]TopicPrerequisitesResult, error) {
	if !gs.IsAvailable() {
		return gs.sqlPrerequisites(topicID)
	}

	cypher := fmt.Sprintf(
		`MATCH (prereq:Topic)-[:PREREQUISITE_OF]->(t:Topic {id: %d})
		 RETURN prereq.id AS id, prereq.name AS name, prereq.slug AS slug, prereq.level AS level`,
		topicID,
	)
	stmt := fmt.Sprintf(
		`SELECT * FROM cypher('%s', $$ %s $$) AS (id agtype, name agtype, slug agtype, level agtype)`,
		gs.graph, cypher,
	)
	rows, err := gs.sqlDB.Query(stmt)
	if err != nil {
		return gs.sqlPrerequisites(topicID)
	}
	defer rows.Close()

	var results []TopicPrerequisitesResult
	for rows.Next() {
		var idRaw, nameRaw, slugRaw, levelRaw string
		if err := rows.Scan(&idRaw, &nameRaw, &slugRaw, &levelRaw); err != nil {
			continue
		}
		var r TopicPrerequisitesResult
		json.Unmarshal([]byte(idRaw), &r.TopicID)
		json.Unmarshal([]byte(nameRaw), &r.Name)
		json.Unmarshal([]byte(slugRaw), &r.Slug)
		json.Unmarshal([]byte(levelRaw), &r.DifficultyLvl)
		results = append(results, r)
	}
	if len(results) == 0 {
		return gs.sqlPrerequisites(topicID)
	}
	return results, nil
}

func (gs *GraphService) sqlPrerequisites(topicID int) ([]TopicPrerequisitesResult, error) {
	var parentID *int
	gs.db.Table("topics").Select("parent_topic_id").Where("topic_id = ?", topicID).Scan(&parentID)
	if parentID == nil {
		return []TopicPrerequisitesResult{}, nil
	}

	type topicRow struct {
		TopicID         int
		Name            string
		Slug            string
		DifficultyLevel *int
	}
	var parent topicRow
	gs.db.Table("topics").Select("topic_id, name, slug, difficulty_level").
		Where("topic_id = ?", *parentID).Scan(&parent)

	lvl := 0
	if parent.DifficultyLevel != nil {
		lvl = *parent.DifficultyLevel
	}
	return []TopicPrerequisitesResult{{
		TopicID: parent.TopicID, Name: parent.Name, Slug: parent.Slug, DifficultyLvl: lvl,
	}}, nil
}

// ── Bulk seeding ──────────────────────────────────────────────────────────────

// SeedGraph creates graph nodes and edges from the relational DB.
func (gs *GraphService) SeedGraph() error {
	if !gs.IsAvailable() {
		return nil // nothing to do without AGE
	}

	// Topics
	type topicRow struct {
		TopicID         int
		Name            string
		Slug            string
		DifficultyLevel *int
		ParentTopicID   *int
	}
	var topics []topicRow
	gs.db.Table("topics").Select("topic_id, name, slug, difficulty_level, parent_topic_id").Scan(&topics)

	for _, t := range topics {
		lvl := 0
		if t.DifficultyLevel != nil {
			lvl = *t.DifficultyLevel
		}
		gs.UpsertTopicNode(t.TopicID, t.Name, t.Slug, lvl)
	}
	for _, t := range topics {
		if t.ParentTopicID != nil {
			gs.LinkTopicPrerequisite(*t.ParentTopicID, t.TopicID)
		}
	}

	// Problems
	type problemRow struct {
		ProblemID       int
		Title           string
		DifficultyScore float64
		PrimaryPattern  *string
	}
	var problems []problemRow
	gs.db.Table("problems").Select("problem_id, title, difficulty_score, primary_pattern").Scan(&problems)

	for _, p := range problems {
		pat := ""
		if p.PrimaryPattern != nil {
			pat = *p.PrimaryPattern
		}
		gs.UpsertProblemNode(p.ProblemID, p.Title, p.DifficultyScore, pat)
	}

	// Problem → Topic edges
	type ptRow struct {
		ProblemID      int
		TopicID        int
		RelevanceScore float64
		IsPrimary      bool
	}
	var pts []ptRow
	gs.db.Table("problem_topics").Select("problem_id, topic_id, relevance_score, is_primary").Scan(&pts)
	for _, pt := range pts {
		gs.LinkProblemToTopic(pt.ProblemID, pt.TopicID, pt.RelevanceScore, pt.IsPrimary)
	}

	// Auto SIMILAR_TO edges for same-pattern problems
	patGroups := map[string][]int{}
	for _, p := range problems {
		if p.PrimaryPattern != nil {
			patGroups[*p.PrimaryPattern] = append(patGroups[*p.PrimaryPattern], p.ProblemID)
		}
	}
	for _, ids := range patGroups {
		for i := 0; i < len(ids); i++ {
			for j := i + 1; j < len(ids); j++ {
				gs.LinkSimilarProblems(ids[i], ids[j], 0.8)
			}
		}
	}

	return nil
}

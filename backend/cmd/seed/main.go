package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/yourusername/algoholic/config"
	"github.com/yourusername/algoholic/models"
	"github.com/yourusername/algoholic/seed"
)

func main() {
	log.Println("Starting database seeding...")

	// Load configuration
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(cfg.Database.GetDSN()), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run auto-migration to ensure schema is up to date
	log.Println("📋 Running database migrations...")
	if err := models.AutoMigrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("✅ Migrations complete")

	// Seed Topics
	log.Println("\n🏷️  Seeding topics...")
	topics := seed.GetSeedTopics()
	successCount := 0
	for i, topic := range topics {
		// Check if topic already exists
		var existing models.Topic
		result := db.Where("slug = ?", topic.Slug).First(&existing)

		if result.Error == nil {
			log.Printf("   ⏭️  Topic '%s' already exists, skipping", topic.Name)
			continue
		}

		if err := db.Create(&topic).Error; err != nil {
			log.Printf("   ❌ Failed to seed topic '%s': %v", topic.Name, err)
		} else {
			successCount++
			log.Printf("   ✅ [%d/%d] Seeded topic: %s", i+1, len(topics), topic.Name)
		}
	}
	log.Printf("✅ Seeded %d/%d topics\n", successCount, len(topics))

	// Seed Problems
	log.Println("🧩 Seeding problems...")
	problems := seed.GetSeedProblems()
	successCount = 0
	problemIDMap := make(map[int]int) // Map from array index to database ID

	for i, problem := range problems {
		// Check if problem already exists
		var existing models.Problem
		result := db.Where("slug = ?", problem.Slug).First(&existing)

		if result.Error == nil {
			log.Printf("   ⏭️  Problem '%s' already exists, skipping", problem.Title)
			problemIDMap[i+1] = existing.ProblemID
			continue
		}

		if err := db.Create(&problem).Error; err != nil {
			log.Printf("   ❌ Failed to seed problem '%s': %v", problem.Title, err)
		} else {
			successCount++
			problemIDMap[i+1] = problem.ProblemID
			log.Printf("   ✅ [%d/%d] Seeded problem: %s (Difficulty: %.0f)", i+1, len(problems), problem.Title, problem.DifficultyScore)
		}
	}
	log.Printf("✅ Seeded %d/%d problems\n", successCount, len(problems))

	// Seed Questions
	log.Println("❓ Seeding questions...")
	questions := seed.GetSeedQuestions()
	successCount = 0

	for i, question := range questions {
		// Update problem_id reference if it exists
		if question.ProblemID != nil {
			if dbID, ok := problemIDMap[*question.ProblemID]; ok {
				question.ProblemID = &dbID
			}
		}

		if err := db.Create(&question).Error; err != nil {
			log.Printf("   ❌ Failed to seed question %d: %v", i+1, err)
		} else {
			successCount++
			questionType := question.QuestionType
			if question.QuestionSubtype != nil {
				questionType = fmt.Sprintf("%s/%s", question.QuestionType, *question.QuestionSubtype)
			}
			log.Printf("   ✅ [%d/%d] Seeded question: %s (Type: %s)", i+1, len(questions), truncate(question.QuestionText, 50), questionType)
		}
	}
	log.Printf("✅ Seeded %d/%d questions\n", successCount, len(questions))

	// Create Problem-Topic relationships
	log.Println("\n🔗 Creating problem-topic relationships...")
	relationshipCount := 0

	// Map problem patterns to topic IDs
	topicMap := map[string]int{}
	var dbTopics []models.Topic
	db.Find(&dbTopics)
	for _, t := range dbTopics {
		topicMap[t.Slug] = t.TopicID
	}

	// Get all problems
	var dbProblems []models.Problem
	db.Find(&dbProblems)

	for _, problem := range dbProblems {
		if problem.PrimaryPattern != nil {
			patternSlug := slugify(*problem.PrimaryPattern)
			if topicID, ok := topicMap[patternSlug]; ok {
				relationship := models.ProblemTopic{
					ProblemID:      problem.ProblemID,
					TopicID:        topicID,
					RelevanceScore: 1.0,
					IsPrimary:      true,
				}

				if err := db.FirstOrCreate(&relationship, relationship).Error; err != nil {
					log.Printf("   ❌ Failed to create relationship for %s: %v", problem.Title, err)
				} else {
					relationshipCount++
				}
			}
		}

		// Add secondary patterns
		for _, pattern := range problem.SecondaryPatterns {
			patternSlug := slugify(pattern)
			if topicID, ok := topicMap[patternSlug]; ok {
				relationship := models.ProblemTopic{
					ProblemID:      problem.ProblemID,
					TopicID:        topicID,
					RelevanceScore: 0.7,
					IsPrimary:      false,
				}

				if err := db.FirstOrCreate(&relationship, relationship).Error; err != nil {
					log.Printf("   ❌ Failed to create secondary relationship: %v", err)
				} else {
					relationshipCount++
				}
			}
		}
	}
	log.Printf("✅ Created %d problem-topic relationships\n", relationshipCount)

	// Print summary
	log.Println("\n" + strings.Repeat("=", 60))
	log.Println("🎉 Database Seeding Complete!")
	log.Println(strings.Repeat("=", 60))

	var totalProblems, totalQuestions, totalTopics int64
	db.Model(&models.Problem{}).Count(&totalProblems)
	db.Model(&models.Question{}).Count(&totalQuestions)
	db.Model(&models.Topic{}).Count(&totalTopics)

	log.Printf("📊 Database Statistics:")
	log.Printf("   Topics: %d", totalTopics)
	log.Printf("   Problems: %d", totalProblems)
	log.Printf("   Questions: %d", totalQuestions)
	log.Printf("   Relationships: %d", relationshipCount)
	log.Println(strings.Repeat("=", 60))

	log.Println("\n✅ You can now start practicing!")
	log.Println("   Backend: http://localhost:4000")
	log.Println("   Frontend: http://localhost:5173")
}

// Helper functions
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	return s
}

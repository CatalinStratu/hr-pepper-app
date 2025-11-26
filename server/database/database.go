package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// Config holds database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// GetConfigFromEnv reads database configuration from environment variables
func GetConfigFromEnv() Config {
	return Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3306"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", ""),
		Database: getEnv("DB_NAME", "hr_interview"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Initialize sets up the MySQL database connection and creates tables
func Initialize(config Config) error {
	var err error

	// Build DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		config.User, config.Password, config.Host, config.Port, config.Database)

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(10)
	DB.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create tables
	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	// Seed default questions if empty
	if err := seedDefaultQuestions(); err != nil {
		log.Printf("Warning: Could not seed default questions: %v", err)
	}

	log.Println("MySQL database initialized successfully")
	return nil
}

func createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS questions (
			id VARCHAR(36) PRIMARY KEY,
			text TEXT NOT NULL,
			category VARCHAR(100) NOT NULL,
			difficulty VARCHAR(20) DEFAULT 'medium',
			time_limit INT DEFAULT 120,
			question_order INT DEFAULT 0,
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_category (category),
			INDEX idx_is_active (is_active)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		`CREATE TABLE IF NOT EXISTS interviews (
			id VARCHAR(36) PRIMARY KEY,
			candidate_name VARCHAR(255) NOT NULL,
			candidate_email VARCHAR(255),
			position VARCHAR(255) NOT NULL,
			status ENUM('pending', 'in_progress', 'completed', 'cancelled') DEFAULT 'pending',
			started_at TIMESTAMP NULL,
			completed_at TIMESTAMP NULL,
			total_score DECIMAL(5,4),
			notes TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_status (status),
			INDEX idx_created_at (created_at),
			INDEX idx_candidate_email (candidate_email)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		`CREATE TABLE IF NOT EXISTS responses (
			id VARCHAR(36) PRIMARY KEY,
			interview_id VARCHAR(36) NOT NULL,
			question_id VARCHAR(36) NOT NULL,
			response_text TEXT NOT NULL,
			response_audio LONGTEXT,
			duration INT DEFAULT 0,
			sentiment_score DECIMAL(5,4),
			confidence_score DECIMAL(5,4),
			keywords_found JSON,
			analysis TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (interview_id) REFERENCES interviews(id) ON DELETE CASCADE,
			FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
			INDEX idx_interview_id (interview_id),
			INDEX idx_question_id (question_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			return err
		}
	}

	return nil
}

func seedDefaultQuestions() error {
	// Check if questions exist
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM questions").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // Questions already exist
	}

	defaultQuestions := []struct {
		id         string
		text       string
		category   string
		difficulty string
		timeLimit  int
		order      int
	}{
		{"q1", "Tell me about yourself and your background.", "introduction", "easy", 120, 1},
		{"q2", "What interests you about this position?", "motivation", "easy", 90, 2},
		{"q3", "Describe a challenging situation you faced at work and how you handled it.", "behavioral", "medium", 180, 3},
		{"q4", "What are your greatest strengths?", "self-assessment", "easy", 90, 4},
		{"q5", "Where do you see yourself in 5 years?", "goals", "medium", 120, 5},
		{"q6", "Tell me about a time you worked effectively in a team.", "teamwork", "medium", 150, 6},
		{"q7", "How do you handle pressure and tight deadlines?", "stress-management", "medium", 120, 7},
		{"q8", "What is your approach to learning new skills?", "adaptability", "easy", 90, 8},
		{"q9", "Describe a situation where you had to solve a complex problem.", "problem-solving", "hard", 180, 9},
		{"q10", "Do you have any questions for us?", "closing", "easy", 120, 10},
	}

	stmt, err := DB.Prepare(`INSERT INTO questions (id, text, category, difficulty, time_limit, question_order, is_active)
		VALUES (?, ?, ?, ?, ?, ?, TRUE)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, q := range defaultQuestions {
		if _, err := stmt.Exec(q.id, q.text, q.category, q.difficulty, q.timeLimit, q.order); err != nil {
			return err
		}
	}

	log.Println("Default questions seeded successfully")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

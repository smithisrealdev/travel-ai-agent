package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/smithisrealdev/travel-ai-agent/backend/internal/config"
)

// PostgresDB wraps the database connection
type PostgresDB struct {
	DB *sql.DB
}

// NewPostgresDB creates a new PostgreSQL connection
func NewPostgresDB(cfg *config.Config) (*PostgresDB, error) {
	var connStr string
	
	// Use DATABASE_URL if available, otherwise build from components
	if cfg.Database.URL != "" {
		connStr = cfg.Database.URL
	} else {
		connStr = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Database,
		)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")

	// Initialize schema
	if err := initSchema(db); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return &PostgresDB{DB: db}, nil
}

// initSchema creates the necessary database tables
func initSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS travel_searches (
		id SERIAL PRIMARY KEY,
		user_id VARCHAR(255),
		destination VARCHAR(255) NOT NULL,
		start_date DATE,
		end_date DATE,
		budget DECIMAL(10, 2),
		preferences JSONB,
		results JSONB,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_travel_searches_user_id ON travel_searches(user_id);
	CREATE INDEX IF NOT EXISTS idx_travel_searches_destination ON travel_searches(destination);
	CREATE INDEX IF NOT EXISTS idx_travel_searches_created_at ON travel_searches(created_at);

	CREATE TABLE IF NOT EXISTS travel_recommendations (
		id SERIAL PRIMARY KEY,
		search_id INTEGER REFERENCES travel_searches(id) ON DELETE CASCADE,
		recommendation_type VARCHAR(50) NOT NULL,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		price DECIMAL(10, 2),
		rating DECIMAL(3, 2),
		metadata JSONB,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_recommendations_search_id ON travel_recommendations(search_id);
	CREATE INDEX IF NOT EXISTS idx_recommendations_type ON travel_recommendations(recommendation_type);
	`

	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	log.Println("Database schema initialized successfully")
	return nil
}

// Close closes the database connection
func (db *PostgresDB) Close() error {
	return db.DB.Close()
}

// HealthCheck verifies the database connection is healthy
func (db *PostgresDB) HealthCheck() error {
	return db.DB.Ping()
}

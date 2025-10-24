package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Backend  BackendConfig
	Database DatabaseConfig
	Redis    RedisConfig
	OpenAI   OpenAIConfig
	Weather  WeatherConfig
	Flight   FlightConfig
	Hotel    HotelConfig
	JWT      JWTConfig
	Env      EnvironmentConfig
}

// BackendConfig holds backend server configuration
type BackendConfig struct {
	Port string
	Host string
}

// DatabaseConfig holds PostgreSQL database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	URL      string
}

// RedisConfig holds Redis cache configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       string
	URL      string
}

// OpenAIConfig holds OpenAI API configuration
type OpenAIConfig struct {
	APIKey string
	Model  string
}

// WeatherConfig holds Weather API configuration
type WeatherConfig struct {
	APIKey string
	URL    string
}

// FlightConfig holds Flight API configuration
type FlightConfig struct {
	APIKey string
	URL    string
}

// HotelConfig holds Hotel API configuration
type HotelConfig struct {
	APIKey string
	URL    string
}

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	Secret    string
	ExpiresIn string
}

// EnvironmentConfig holds environment-specific settings
type EnvironmentConfig struct {
	Environment string
	Debug       bool
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists (for local development)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{
		Backend: BackendConfig{
			Port: getEnv("BACKEND_PORT", "8080"),
			Host: getEnv("BACKEND_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "travelagent"),
			Password: getEnv("POSTGRES_PASSWORD", ""),
			Database: getEnv("POSTGRES_DB", "travelagent"),
			URL:      getEnv("DATABASE_URL", ""),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnv("REDIS_DB", "0"),
			URL:      getEnv("REDIS_URL", "redis://localhost:6379/0"),
		},
		OpenAI: OpenAIConfig{
			APIKey: getEnv("OPENAI_API_KEY", ""),
			Model:  getEnv("OPENAI_MODEL", "gpt-4-turbo-preview"),
		},
		Weather: WeatherConfig{
			APIKey: getEnv("WEATHER_API_KEY", ""),
			URL:    getEnv("WEATHER_API_URL", "https://api.openweathermap.org/data/2.5"),
		},
		Flight: FlightConfig{
			APIKey: getEnv("FLIGHT_API_KEY", ""),
			URL:    getEnv("FLIGHT_API_URL", "https://api.aviationstack.com/v1"),
		},
		Hotel: HotelConfig{
			APIKey: getEnv("HOTEL_API_KEY", ""),
			URL:    getEnv("HOTEL_API_URL", "https://api.booking.com"),
		},
		JWT: JWTConfig{
			Secret:    getEnv("JWT_SECRET", "default-secret-change-in-production"),
			ExpiresIn: getEnv("JWT_EXPIRES_IN", "24h"),
		},
		Env: EnvironmentConfig{
			Environment: getEnv("ENVIRONMENT", "development"),
			Debug:       getEnv("DEBUG", "false") == "true",
		},
	}

	return config, nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

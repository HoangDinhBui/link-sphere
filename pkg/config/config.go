package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application.
type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string

	// Kafka
	KafkaBrokers string

	// OpenSearch
	OpenSearchURL string

	// JWT
	JWTSecret string
	JWTExpiry time.Duration

	// Server
	ServerPort string
}

// Load reads configuration from environment variables.
func Load() *Config {
	// Load .env file if it exists (for local development)
	_ = godotenv.Load()

	cfg := &Config{
		DBHost:     getEnv("POSTGRES_HOST", "localhost"),
		DBPort:     getEnv("POSTGRES_PORT", "5432"),
		DBUser:     getEnv("POSTGRES_USER", "linksphere"),
		DBPassword: getEnv("POSTGRES_PASSWORD", "linksphere_secret"),
		DBName:     getEnv("POSTGRES_DB", "linksphere"),

		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),

		KafkaBrokers: getEnv("KAFKA_BROKERS", "localhost:9092"),

		OpenSearchURL: getEnv("OPENSEARCH_URL", "http://localhost:9200"),

		JWTSecret: getEnv("JWT_SECRET", "default-secret"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}

	expiry, err := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))
	if err != nil {
		expiry = 24 * time.Hour
	}
	cfg.JWTExpiry = expiry

	return cfg
}

// DSN returns the PostgreSQL connection string.
func (c *Config) DSN() string {
	return "postgres://" + c.DBUser + ":" + c.DBPassword +
		"@" + c.DBHost + ":" + c.DBPort +
		"/" + c.DBName + "?sslmode=disable"
}

// RedisAddr returns the Redis address.
func (c *Config) RedisAddr() string {
	return c.RedisHost + ":" + c.RedisPort
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

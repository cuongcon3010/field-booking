package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Environment string
	FrontendURL string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// JWT
	JWTSecret          string
	JWTExpirationHours int

	// Google OAuth 2.0
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

func LoadConfig() (*Config, error) {
	// Tải file .env nếu có (không báo lỗi nếu không tìm thấy file, hữu ích cho môi trường docker/prod)
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")

	cfg := &Config{
		Port:               getEnv("PORT", "8080"),
		Environment:        getEnv("ENVIRONMENT", "development"),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:3000"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", "postgres"),
		DBName:             getEnv("DB_NAME", "field_booking"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		JWTSecret:          getEnv("JWT_SECRET", "super-secret-jwt-key-field-booking"),
		JWTExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 72),
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/auth/google/callback"),
	}

	return cfg, nil
}

// ValidateOAuth kiểm tra xem cấu hình Google OAuth có đầy đủ hay không
func (c *Config) ValidateOAuth() error {
	if c.GoogleClientID == "" {
		return fmt.Errorf("GOOGLE_CLIENT_ID is required but not set")
	}
	if c.GoogleClientSecret == "" {
		return fmt.Errorf("GOOGLE_CLIENT_SECRET is required but not set")
	}
	if c.GoogleRedirectURL == "" {
		return fmt.Errorf("GOOGLE_REDIRECT_URL is required but not set")
	}
	return nil
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := getEnv(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}

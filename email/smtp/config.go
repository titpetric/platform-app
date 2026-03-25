package smtp

import (
	"os"
	"strconv"
)

// Config holds SMTP configuration
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// ConfigFromEnv creates a Config from environment variables
func ConfigFromEnv() Config {
	return Config{
		Host:     getEnv("PLATFORM_EMAIL_HOST", "localhost"),
		Port:     getEnvInt("PLATFORM_EMAIL_PORT", 1025),
		Username: getEnv("PLATFORM_EMAIL_USERNAME", ""),
		Password: getEnv("PLATFORM_EMAIL_PASSWORD", ""),
		From:     getEnv("PLATFORM_EMAIL_FROM", "noreply@example.com"),
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvInt gets an integer environment variable with a default value
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

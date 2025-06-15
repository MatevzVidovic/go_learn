// internal/config/config.go
// Fixed version with proper MySQL connection string

package config

import (
	"os"
)

// Config holds all our application settings
type Config struct {
	DatabaseURL string // Where to find our database
	MQTTBroker  string // Where to find our MQTT broker
	JWTSecret   string // Secret key for creating secure tokens
	Port        string // What port our web server should listen on
}

// Load reads environment variables and creates a Config struct
func Load() *Config {
	return &Config{
		// Fixed default database URL with parseTime=true parameter
		// This is CRUCIAL for handling MySQL datetime columns properly
		DatabaseURL: getEnv("DATABASE_URL", "storeuser:storepass@tcp(localhost:3306)/onlinestore?parseTime=true"),
		MQTTBroker:  getEnv("MQTT_BROKER", "tcp://localhost:1883"),
		JWTSecret:   getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production"),
		Port:        getEnv("PORT", "8080"),
	}
}

// getEnv is a helper function that gets an environment variable
// If the environment variable doesn't exist, it returns the fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

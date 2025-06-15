// internal/config/config.go
// This file handles reading configuration from environment variables
// Environment variables let us change settings without recompiling code

package config

import (
	"os"
)

// Config holds all our application settings
// In Go, we use structs to group related data together
type Config struct {
	DatabaseURL string // Where to find our database
	MQTTBroker  string // Where to find our MQTT broker
	JWTSecret   string // Secret key for creating secure tokens
	Port        string // What port our web server should listen on
}

// Load reads environment variables and creates a Config struct
// This function returns a Config with default values if env vars aren't set
func Load() *Config {
	return &Config{
		// getEnv is a helper function that gets an env var or returns a default value
		DatabaseURL: getEnv("DATABASE_URL", "storeuser:storepass@tcp(localhost:3306)/onlinestore"),
		MQTTBroker:  getEnv("MQTT_BROKER", "tcp://localhost:1883"),
		JWTSecret:   getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production"),
		Port:        getEnv("PORT", "8080"),
	}
}

// getEnv is a helper function that gets an environment variable
// If the environment variable doesn't exist, it returns the fallback value
// This is a common pattern in Go for handling configuration
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
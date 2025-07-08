package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	CloudflareAPIToken string
	CloudflareAPIKey   string
	CloudflareEmail    string
	Port               string
	AdminUsername      string
	AdminPassword      string
}

func LoadConfig() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{
		CloudflareAPIToken: os.Getenv("CLOUDFLARE_API_TOKEN"),
		CloudflareAPIKey:   os.Getenv("CLOUDFLARE_API_KEY"),
		CloudflareEmail:    os.Getenv("CLOUDFLARE_EMAIL"),
		Port:               getEnvOrDefault("PORT", "8080"),
		AdminUsername:      getEnvOrDefault("ADMIN_USERNAME", "admin"),
		AdminPassword:      getEnvOrDefault("ADMIN_PASSWORD", "password123"),
	}

	// Validate required configuration
	if config.CloudflareAPIToken == "" && (config.CloudflareAPIKey == "" || config.CloudflareEmail == "") {
		log.Fatal("Either CLOUDFLARE_API_TOKEN or both CLOUDFLARE_API_KEY and CLOUDFLARE_EMAIL must be set")
	}

	return config
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

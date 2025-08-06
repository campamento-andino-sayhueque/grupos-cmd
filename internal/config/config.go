package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration.
type Config struct {
	GoogleCloudProjectID string
	PubSubTopicID        string
	FirestoreCollection  string
	Port                 string
}

// New creates a new Config from environment variables.
func New() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		GoogleCloudProjectID: getEnv("GOOGLE_CLOUD_PROJECT_ID", ""),
		PubSubTopicID:        getEnv("PUB_SUB_TOPIC_ID", "grupos"),
		FirestoreCollection:  getEnv("FIRESTORE_COLLECTION", "eventos"),
		Port:                 getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	if fallback == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return fallback
}

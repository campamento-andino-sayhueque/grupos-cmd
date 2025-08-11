
package config

import (
	"github.com/spf13/viper"
)

// Config holds the application configuration.
type Config struct {
	GoogleCloudProjectID string
	PubSubTopicID        string
	FirestoreCollection  string
	Port                 string
}


// New creates a new Config from environment variables using Viper.
func New() *Config {
	viper.SetDefault("PUB_SUB_TOPIC_ID", "grupos")
	viper.SetDefault("FIRESTORE_COLLECTION", "eventos")
	viper.SetDefault("PORT", "8080")

	viper.AutomaticEnv() // Lee variables de entorno automáticamente

	return &Config{
		GoogleCloudProjectID: viper.GetString("GOOGLE_CLOUD_PROJECT_ID"),
		PubSubTopicID:        viper.GetString("PUB_SUB_TOPIC_ID"),
		FirestoreCollection:  viper.GetString("FIRESTORE_COLLECTION"),
		Port:                 viper.GetString("PORT"),
	}
}



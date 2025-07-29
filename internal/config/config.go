package config

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var configInstance *Config

func GetConfig() *Config {
	if configInstance == nil {
		var err error
		configInstance, err = LoadConfig()
		if err != nil {
			panic(fmt.Sprintf("Failed to load config: %v", err))
		}
	}
	return configInstance
}

type Config struct {
	Port          int
	MongoURI      string
	MongoDatabase string
	JWTSecret     string

	// Verification settings
	VerificationCodeLength int
	VerificationCodeExpiry time.Duration

	// Mailchimp configuration
	MailchimpAPIKey string
	MailchimpListID string
	AdminEmail      string

	// EmailJS configuration
	EmailJSServiceID   string
	EmailJSTemplateID  string
	EmailJSUserID      string
	EmailJSAccessToken string
}

func LoadConfig() (*Config, error) {
	// Try to load .env from root directory (for local development only)
	// In Cloud Run, we'll rely on environment variables being set directly
	_ = godotenv.Load() // Ignore errors - will use system env if .env doesn't exist

	config := &Config{
		Port:                   getEnvAsInt("PORT", 8080),
		MongoURI:               getEnv("MONGODB_URI", ""),
		MongoDatabase:          getEnv("MONGODB_DATABASE", ""),
		JWTSecret:              getEnv("JWT_SECRET", ""),
		VerificationCodeLength: 6,
		VerificationCodeExpiry: 15 * time.Minute,
		MailchimpAPIKey:        getEnv("MAILCHIMP_API_KEY", ""),
		MailchimpListID:        getEnv("MAILCHIMP_LIST_ID", ""),
		AdminEmail:             getEnv("ADMIN_EMAIL", ""),
		EmailJSServiceID:       getEnv("EMAILJS_SERVICE_ID", ""),
		EmailJSTemplateID:      getEnv("EMAILJS_TEMPLATE_ID", ""),
		EmailJSUserID:          getEnv("EMAILJS_USER_ID", ""),
		EmailJSAccessToken:     getEnv("EMAILJS_ACCESS_TOKEN", ""),
	}

	// Validate required environment variables
	required := map[string]string{
		"MONGODB_URI":         config.MongoURI,
		"MONGODB_DATABASE":    config.MongoDatabase,
		"JWT_SECRET":          config.JWTSecret,
		"ADMIN_EMAIL":         config.AdminEmail,
		"EMAILJS_SERVICE_ID":  config.EmailJSServiceID,
		"EMAILJS_TEMPLATE_ID": config.EmailJSTemplateID,
		"EMAILJS_USER_ID":     config.EmailJSUserID,
	}

	for key, value := range required {
		if value == "" {
			return nil, fmt.Errorf("required environment variable %s is not set", key)
		}
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	var value int
	_, err := fmt.Sscanf(valueStr, "%d", &value)
	if err != nil {
		return defaultValue
	}
	return value
}

func NewDB(config *Config) (*sql.DB, error) {
	// MongoDB doesn't use sql.DB, this function will be replaced with MongoDB client initialization
	return nil, fmt.Errorf("MongoDB client initialization not implemented yet")
}

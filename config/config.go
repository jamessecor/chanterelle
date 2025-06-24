package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration
var Config struct {
	EmailJSUserID     string
	EmailJSTemplate   string
	EmailJSPublicKey  string
	EmailJSPrivateKey string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() error {
	// Load .env file
	if err := godotenv.Load("/app/.env"); err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	// Load config values from environment variables
	Config.EmailJSUserID = os.Getenv("EMAILJS_USER_ID")
	Config.EmailJSTemplate = os.Getenv("EMAILJS_TEMPLATE")
	Config.EmailJSPublicKey = os.Getenv("EMAILJS_PUBLIC_KEY")
	Config.EmailJSPrivateKey = os.Getenv("EMAILJS_PRIVATE_KEY")
	Config.DBHost = os.Getenv("DB_HOST")
	Config.DBPort = os.Getenv("DB_PORT")
	Config.DBUser = os.Getenv("DB_USER")
	Config.DBPassword = os.Getenv("DB_PASSWORD")
	Config.DBName = os.Getenv("DB_NAME")

	// Validate required config values
	required := []string{
		"EMAILJS_USER_ID",
		"EMAILJS_TEMPLATE",
		"EMAILJS_PRIVATE_KEY",
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
	}

	for _, key := range required {
		value := os.Getenv(key)
		if value == "" {
			return fmt.Errorf("required environment variable %s is not set", key)
		}
	}

	return nil
}

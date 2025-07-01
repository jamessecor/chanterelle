package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration
var Config struct {
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	TwilioSID        string
	TwilioToken      string
	TwilioNumber     string
	TwilioContentSid string
	JWTSecret        string
	AdminPhoneNumber string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() error {
	// Load .env file
	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	// Load config values from environment variables
	Config.DBHost = os.Getenv("DB_HOST")
	Config.DBPort = os.Getenv("DB_PORT")
	Config.DBUser = os.Getenv("DB_USER")
	Config.DBPassword = os.Getenv("DB_PASSWORD")
	Config.DBName = os.Getenv("DB_NAME")
	Config.TwilioSID = os.Getenv("TWILIO_SID")
	Config.TwilioToken = os.Getenv("TWILIO_TOKEN")
	Config.TwilioNumber = os.Getenv("TWILIO_NUMBER")
	Config.TwilioContentSid = os.Getenv("TWILIO_CONTENT_SID")
	Config.JWTSecret = os.Getenv("JWT_SECRET")
	Config.AdminPhoneNumber = os.Getenv("ADMIN_PHONE_NUMBER")

	// Validate required config values
	required := []string{
		"TWILIO_SID",
		"TWILIO_TOKEN",
		"TWILIO_NUMBER",
		"TWILIO_CONTENT_SID",
		"JWT_SECRET",
		"ADMIN_PHONE_NUMBER",
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

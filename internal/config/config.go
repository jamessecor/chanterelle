package config

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"
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

	VerificationCodeLength int
	VerificationCodeExpiry time.Duration

	TwilioAccountSID           string
	TwilioAuthToken            string
	TwilioNumber               string
	TwilioContentSID           string
	AvailableAdminPhoneNumbers []string
	MailchimpAPIKey            string
	MailchimpListID            string
	AdminEmail                 string
}

func LoadConfig() (*Config, error) {
	// Try to load .env from root directory
	// if err := godotenv.Load(".env"); err != nil {
	// 	return nil, fmt.Errorf("error loading .env file: %v", err)
	// }

	config := &Config{
		Port:                       8080,
		MongoURI:                   os.Getenv("MONGODB_URI"),
		MongoDatabase:              os.Getenv("MONGODB_DATABASE"),
		JWTSecret:                  os.Getenv("JWT_SECRET"),
		VerificationCodeLength:     6,
		VerificationCodeExpiry:     15 * time.Minute,
		TwilioAccountSID:           os.Getenv("TWILIO_ACCOUNT_SID"),
		TwilioAuthToken:            os.Getenv("TWILIO_AUTH_TOKEN"),
		TwilioNumber:               os.Getenv("TWILIO_NUMBER"),
		TwilioContentSID:           os.Getenv("TWILIO_CONTENT_SID"),
		AvailableAdminPhoneNumbers: strings.Split(os.Getenv("AVAILABLE_ADMIN_PHONE_NUMBERS"), ","),
		MailchimpAPIKey:            os.Getenv("MAILCHIMP_API_KEY"),
		MailchimpListID:            os.Getenv("MAILCHIMP_LIST_ID"),
		AdminEmail:                 os.Getenv("ADMIN_EMAIL"),
	}

	// Validate required config values
	required := []string{
		"MONGODB_URI",
		"MONGODB_DATABASE",
		"JWT_SECRET",
		"TWILIO_ACCOUNT_SID",
		"TWILIO_AUTH_TOKEN",
		"TWILIO_NUMBER",
		"TWILIO_CONTENT_SID",
		"AVAILABLE_ADMIN_PHONE_NUMBERS",
		"MAILCHIMP_API_KEY",
		"MAILCHIMP_LIST_ID",
		"ADMIN_EMAIL",
	}

	for _, key := range required {
		value := os.Getenv(key)
		if value == "" {
			return nil, fmt.Errorf("required environment variable %s is not set", key)
		}
	}

	return config, nil
}

func NewDB(config *Config) (*sql.DB, error) {
	// MongoDB doesn't use sql.DB, this function will be replaced with MongoDB client initialization
	return nil, fmt.Errorf("MongoDB client initialization not implemented yet")
}

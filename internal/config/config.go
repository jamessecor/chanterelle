package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
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

	// EmailJS configuration
	EmailJSServiceID   string
	EmailJSTemplateID  string
	EmailJSUserID      string
	EmailJSAccessToken string
}

func LoadConfig() (*Config, error) {
	// Try to load .env from root directory
	log.Println("Loading config from .env file...")
	if err := godotenv.Load(".env"); err != nil {
		log.Println("error loading .env file: %v", err)
		log.Println("Warning: we will attempt to use system config...")
	}

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
		EmailJSServiceID:           os.Getenv("EMAILJS_SERVICE_ID"),
		EmailJSTemplateID:          os.Getenv("EMAILJS_TEMPLATE_ID"),
		EmailJSUserID:              os.Getenv("EMAILJS_USER_ID"),
		EmailJSAccessToken:         os.Getenv("EMAILJS_ACCESS_TOKEN"),
	}

	return config, nil
}

func NewDB(config *Config) (*sql.DB, error) {
	// MongoDB doesn't use sql.DB, this function will be replaced with MongoDB client initialization
	return nil, fmt.Errorf("MongoDB client initialization not implemented yet")
}

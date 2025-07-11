package config

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

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
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string

	TwilioAccountSID           string
	TwilioAuthToken            string
	TwilioNumber               string
	TwilioContentSID           string
	AvailableAdminPhoneNumbers []string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	config := &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		JWTSecret:  os.Getenv("JWT_SECRET"),

		TwilioAccountSID:           os.Getenv("TWILIO_ACCOUNT_SID"),
		TwilioAuthToken:            os.Getenv("TWILIO_AUTH_TOKEN"),
		TwilioNumber:               os.Getenv("TWILIO_NUMBER"),
		TwilioContentSID:           os.Getenv("TWILIO_CONTENT_SID"),
		AvailableAdminPhoneNumbers: strings.Split(os.Getenv("AVAILABLE_ADMIN_PHONE_NUMBERS"), ","),
	}

	// Validate required config values
	required := []string{
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"JWT_SECRET",
		"TWILIO_ACCOUNT_SID",
		"TWILIO_AUTH_TOKEN",
		"TWILIO_NUMBER",
		"TWILIO_CONTENT_SID",
		"AVAILABLE_ADMIN_PHONE_NUMBERS",
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
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

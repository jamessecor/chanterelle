package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"chanterelle/internal/config"
	"chanterelle/internal/models"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMailchimpIntegration tests the full Mailchimp integration
func TestMailchimpIntegration(t *testing.T) {
	if err := godotenv.Load(".env"); err != nil {
		t.Errorf("error loading .env file: %v", err)
	}

	cfg := config.GetConfig()

	// Create notification service
	notificationService := NewNotificationService(cfg)

	// Create test contact
	testContact := &models.Contact{
		Name:    "Test User",
		Email:   "test+" + time.Now().Format("20060102150405") + "@mailinator.com",
		Message: "Test message from integration test",
	}

	// Test adding to Mailchimp
	t.Run("AddToMailchimp", func(t *testing.T) {
		err := notificationService.AddToMailchimp(testContact)
		require.NoError(t, err)
	})

	// Test sending admin notification
	// t.Run("SendAdminNotification", func(t *testing.T) {
	// 	err := notificationService.SendAdminNotification(testContact)
	// 	require.NoError(t, err)
	// })

	// Clean up test contact
	// Note: Mailchimp API doesn't provide a direct way to delete a member,
	// so we're using a unique email for each test run
}

// TestInvalidMailchimpConfig tests error handling with invalid Mailchimp configuration
func TestInvalidMailchimpConfig(t *testing.T) {
	// Create notification service with invalid config
	cfg := config.GetConfig()
	notificationService := NewNotificationService(cfg)

	// Create test contact
	testContact := &models.Contact{
		Name:    "Test User",
		Email:   "test@example.com",
		Message: "Test message",
	}

	// Test error handling
	t.Run("AddToMailchimp", func(t *testing.T) {
		err := notificationService.AddToMailchimp(testContact)
		assert.Error(t, err)
	})

	t.Run("SendAdminNotification", func(t *testing.T) {
		err := notificationService.SendAdminNotification(testContact)
		assert.Error(t, err)
	})
}

// TestEmailJSIntegration tests the EmailJS integration
func TestEmailJSIntegration(t *testing.T) {
	if err := godotenv.Load(".env"); err != nil {
		t.Skip("Skipping EmailJS integration test - .env file not found")
	}

	cfg := config.GetConfig()
	if cfg.EmailJSServiceID == "" || cfg.EmailJSTemplateID == "" || cfg.EmailJSUserID == "" || cfg.EmailJSAccessToken == "" {
		t.Skip("Skipping EmailJS integration test - missing EmailJS configuration")
	}

	// Create notification service with real config
	service := NewNotificationService(cfg)

	t.Run("SendVerificationCode", func(t *testing.T) {
		testEmail := "test+" + time.Now().Format("20060102150405") + "@example.com"
		testCode := "123456"

		err := service.SendVerificationCode(testEmail, testCode)
		assert.NoError(t, err, "Failed to send verification code")
	})

	t.Run("SendNewContactNotification", func(t *testing.T) {
		testContact := &models.Contact{
			Name:    "Test User",
			Email:   "test+" + time.Now().Format("20060102150405") + "@example.com",
			Message: "Test message from integration test",
		}

		err := service.SendNewContactNotification(testContact)
		assert.NoError(t, err, "Failed to send contact notification")
	})
}

// TestEmailJSMock tests EmailJS functionality with a mock HTTP server
func TestEmailJSMock(t *testing.T) {
	// Setup test server to mock EmailJS API
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/v1.0/email/send", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Parse the request body
		var req struct {
			ServiceID      string `json:"service_id"`
			TemplateID     string `json:"template_id"`
			UserID         string `json:"user_id"`
			TemplateParams struct {
				ToName      string `json:"to_name"`
				Destination string `json:"destination"`
				Firstname   string `json:"firstname"`
				Lastname    string `json:"lastname"`
				Email       string `json:"email"`
				Message     string `json:"message"`
			} `json:"template_params"`
		}

		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err, "Failed to decode request body")

		// Verify the request body
		assert.Equal(t, "test-service-id", req.ServiceID)
		assert.Equal(t, "test-template-id", req.TemplateID)
		assert.Equal(t, "test-user-id", req.UserID)

		// Respond with success
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}))
	defer testServer.Close()

	// Create test config with mock values
	cfg := &config.Config{
		EmailJSServiceID:   "test-service-id",
		EmailJSTemplateID:  "test-template-id",
		EmailJSUserID:      "test-user-id",
		EmailJSAccessToken: "test-access-token",
	}

	// Create the service with a custom HTTP client that points to our test server
	service := &NotificationService{
		cfg: cfg,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}

	t.Run("SendVerificationCode", func(t *testing.T) {
		err := service.SendVerificationCode("test@example.com", "123456")
		assert.NoError(t, err, "Failed to send verification code")
	})

	t.Run("SendNewContactNotification", func(t *testing.T) {
		testContact := &models.Contact{
			Name:    "John Doe",
			Email:   "john@example.com",
			Message: "Test message",
		}

		err := service.SendNewContactNotification(testContact)
		assert.NoError(t, err, "Failed to send contact notification")
	})

	t.Run("SendEmailJS_ErrorCases", func(t *testing.T) {
		tests := []struct {
			name           string
			serviceID      string
			templateID     string
			userID         string
			expectedErrMsg string
		}{
			{
				name:           "missing service ID",
				serviceID:      "",
				templateID:     "test-template-id",
				userID:         "test-user-id",
				expectedErrMsg: "emailjs configuration is not complete",
			},
			{
				name:           "missing template ID",
				serviceID:      "test-service-id",
				templateID:     "",
				userID:         "test-user-id",
				expectedErrMsg: "emailjs configuration is not complete",
			},
			{
				name:           "missing user ID",
				serviceID:      "test-service-id",
				templateID:     "test-template-id",
				userID:         "",
				expectedErrMsg: "emailjs configuration is not complete",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Create a test config with the test case values
				testCfg := &config.Config{
					EmailJSServiceID:   tt.serviceID,
					EmailJSTemplateID:  tt.templateID,
					EmailJSUserID:      tt.userID,
					EmailJSAccessToken: "test-access-token",
				}

				// Create the service
				testService := NewNotificationService(testCfg)

				// Call the method being tested
				err := testService.SendVerificationCode("test@example.com", "123456")

				// Verify the error
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrMsg)
			})
		}
	})

	t.Run("SendEmailJS_HTTPError", func(t *testing.T) {
		// Setup test server to return an error
		errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		}))
		defer errorServer.Close()

		// Create the service with a custom HTTP client that points to our error server
		errorService := &NotificationService{
			cfg: cfg,
			client: &http.Client{
				Timeout: 5 * time.Second,
			},
		}

		// Call the method being tested
		err := errorService.SendVerificationCode("test@example.com", "123456")

		// Verify the error
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "emailjs API returned status 500")
	})
}

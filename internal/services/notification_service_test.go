package services

import (
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

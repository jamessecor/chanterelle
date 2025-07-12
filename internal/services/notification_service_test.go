package services

import (
	"testing"
	"time"

	"chanterelle/internal/config"
	"chanterelle/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMailchimpIntegration tests the full Mailchimp integration
func TestMailchimpIntegration(t *testing.T) {
	// Skip test if not running in CI or if MAILCHIMP_TEST is not set
	if os.Getenv("MAILCHIMP_TEST") != "true" {
		t.Skip("Skipping Mailchimp integration test. Set MAILCHIMP_TEST=true to run")
	}

	// Load test configuration
	cfg, err := config.LoadConfig()
	require.NoError(t, err)

	// Create notification service
	notificationService := NewNotificationService(cfg)

	// Create test contact
	testContact := &models.Contact{
		Name:    "Test User",
		Email:   "test+" + time.Now().Format("20060102150405") + "@example.com",
		Message: "Test message from integration test",
	}

	// Test adding to Mailchimp
	t.Run("AddToMailchimp", func(t *testing.T) {
		err := notificationService.AddToMailchimp(testContact)
		require.NoError(t, err)
	})

	// Test sending admin notification
	t.Run("SendAdminNotification", func(t *testing.T) {
		err := notificationService.SendAdminNotification(testContact)
		require.NoError(t, err)
	})

	// Clean up test contact
	// Note: Mailchimp API doesn't provide a direct way to delete a member,
	// so we're using a unique email for each test run
}

// TestInvalidMailchimpConfig tests error handling with invalid Mailchimp configuration
func TestInvalidMailchimpConfig(t *testing.T) {
	// Create notification service with invalid config
	cfg := &config.Config{
		MailchimpAPIKey: "invalid-api-key",
		MailchimpListID: "invalid-list-id",
	}
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

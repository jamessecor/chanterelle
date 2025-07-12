package services

import (
	"chanterelle/internal/config"
	"chanterelle/internal/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type NotificationService struct {
	cfg *config.Config
}

func NewNotificationService(cfg *config.Config) *NotificationService {
	return &NotificationService{cfg: cfg}
}

func (s *NotificationService) AddToMailchimp(contact *models.Contact) error {
	// Mailchimp API endpoint
	endpoint := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/lists/%s/members",
		strings.Split(s.cfg.MailchimpAPIKey, "-")[1],
		s.cfg.MailchimpListID)

	// Prepare Mailchimp subscription data
	data := map[string]interface{}{
		"email_address": contact.Email,
		"status":        "subscribed",
		"merge_fields": map[string]interface{}{
			"FNAME": contact.Name,
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal mailchimp data: %v", err)
	}

	// Create request
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s",
		base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("anystring:%s", s.cfg.MailchimpAPIKey)))))

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send mailchimp request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("mailchimp API returned status: %d", resp.StatusCode)
	}

	return nil
}

func (s *NotificationService) SendAdminNotification(contact *models.Contact) error {
	// Mailchimp API endpoint for sending campaign
	endpoint := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/campaigns",
		strings.Split(s.cfg.MailchimpAPIKey, "-")[1])

	// Prepare campaign data
	data := map[string]interface{}{
		"type": "regular",
		"settings": map[string]interface{}{
			"subject_line": fmt.Sprintf("New Contact Form Submission from %s", contact.Name),
			"title":        "New Contact Form Submission",
			"from_name":    "Contact Form",
			"reply_to":     "noreply@example.com",
		},
		"recipients": map[string]interface{}{
			"list_id": s.cfg.MailchimpListID,
		},
		"content": []map[string]interface{}{
			{
				"type": "text/html",
				"content": fmt.Sprintf(`
					<h2>New Contact Form Submission</h2>
					<p>Name: %s</p>
					<p>Email: %s</p>
					<p>Message: %s</p>
				`, contact.Name, contact.Email, contact.Message),
			},
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal campaign data: %v", err)
	}

	// Create request
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("failed to create campaign request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s",
		base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("anystring:%s", s.cfg.MailchimpAPIKey)))))

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send campaign request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("mailchimp campaign API returned status: %d", resp.StatusCode)
	}

	return nil
}

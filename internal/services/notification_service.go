package services

import (
	"chanterelle/internal/config"
	"chanterelle/internal/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type NotificationService struct {
	cfg *config.Config
}

// SendAdminNotification sends an email notification to admin using Mailchimp Transactional API
// Note: Requires setting ADMIN_EMAIL environment variable

func NewNotificationService(cfg *config.Config) *NotificationService {
	return &NotificationService{cfg: cfg}
}

func (s *NotificationService) AddToMailchimp(contact *models.Contact) error {
	// Mailchimp API endpoint
	datacenter := strings.Split(s.cfg.MailchimpAPIKey, "-")[1]
	endpoint := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/lists/%s/members", datacenter, s.cfg.MailchimpListID)

	// Prepare Mailchimp subscription data
	data := map[string]interface{}{
		"email_address": contact.Email,
		"status":        "subscribed",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal mailchimp data: %v", err)
	}

	// Create request
	log.Printf("Request body: %s", string(jsonData))
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("anystring:"+s.cfg.MailchimpAPIKey)))

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send mailchimp request: %v", err)
	}
	defer resp.Body.Close()

	// Read and log the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}
	log.Printf("Mailchimp response: Status=%d, Body=%s", resp.StatusCode, string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("mailchimp API returned status: %d, error: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (s *NotificationService) SendAdminNotification(contact *models.Contact) error {
	// Mailchimp Transactional API endpoint
	endpoint := "https://mandrillapp.com/api/1.0/messages/send.json"

	// Prepare message data
	data := map[string]interface{}{
		"key": s.cfg.MailchimpAPIKey,
		"message": map[string]interface{}{
			"html": fmt.Sprintf(`
			<h2>New Contact Form Submission</h2>
				<p>Name: %s</p>
				<p>Email: %s</p>
				<p>Message: %s</p>
			`, contact.Name, contact.Email, contact.Message),
			"text":       fmt.Sprintf("Name: %s\nEmail: %s\nMessage: %s", contact.Name, contact.Email, contact.Message),
			"subject":    fmt.Sprintf("New Contact Form Submission from %s", contact.Name),
			"from_email": "noreply@example.com",
			"from_name":  "Contact Form",
			"to": []map[string]interface{}{
				{
					"email": s.cfg.AdminEmail,
					"name":  "Admin",
					"type":  "to",
				},
			},
		},
	}
	log.Println(s.cfg.AdminEmail)
	log.Println(data)

	// Convert to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal message data: %v", err)
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

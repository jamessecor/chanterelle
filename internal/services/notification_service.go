package services

import (
	"bytes"
	"chanterelle/internal/config"
	"chanterelle/internal/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type NotificationService struct {
	cfg    *config.Config
	client *http.Client
}

type emailJSParams struct {
	ToName      string `json:"to_name"`
	Destination string `json:"destination"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Email       string `json:"email"`
	Message     string `json:"message"`
}

type emailJSRequest struct {
	ServiceID      string        `json:"service_id"`
	TemplateID     string        `json:"template_id"`
	UserID         string        `json:"user_id"`
	AccessToken    string        `json:"accessToken"`
	TemplateParams emailJSParams `json:"template_params"`
}

// SendAdminNotification sends an email notification to admin using Mailchimp Transactional API
// Note: Requires setting ADMIN_EMAIL environment variable

func NewNotificationService(cfg *config.Config) *NotificationService {
	return &NotificationService{
		cfg:    cfg,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *NotificationService) AddToMailchimp(contact *models.Contact) error {
	// Mailchimp API endpoint
	datacenter := strings.Split(s.cfg.MailchimpAPIKey, "-")[1]
	endpoint := fmt.Sprintf("https://%s.api.mailchimp.com/3.0/lists/%s/members", datacenter, s.cfg.MailchimpListID)

	// Prepare Mailchimp subscription data
	data := map[string]interface{}{
		"email_address": contact.Email,
		"name":          contact.Name,
		"status":        "subscribed",
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

func (s *NotificationService) sendEmailJS(params emailJSParams) error {
	if s.cfg.EmailJSServiceID == "" || s.cfg.EmailJSTemplateID == "" || s.cfg.EmailJSUserID == "" {
		return fmt.Errorf("emailjs configuration is not complete")
	}

	reqBody := emailJSRequest{
		ServiceID:      s.cfg.EmailJSServiceID,
		TemplateID:     s.cfg.EmailJSTemplateID,
		UserID:         s.cfg.EmailJSUserID,
		AccessToken:    s.cfg.EmailJSAccessToken,
		TemplateParams: params,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal emailjs request: %v", err)
	}

	// Send request to EmailJS
	resp, err := s.client.Post(
		"https://api.emailjs.com/api/v1.0/email/send",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to send email via emailjs: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("emailjs API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (s *NotificationService) SendVerificationCode(email, code string) error {
	log.Println("Sending verification code to", email)
	params := emailJSParams{
		ToName:      "Chanterelle member",
		Destination: fmt.Sprintf("Your verification code is: %s", code),
		Firstname:   "",
		Lastname:    "",
		Email:       email,
		Message:     "Please use this code to verify your admin access.",
	}

	return s.sendEmailJS(params)
}

func (s *NotificationService) SendNewContactNotification(contact *models.Contact) error {
	// First name and last name handling (assuming Name is in format "First Last")
	nameParts := strings.Fields(contact.Name)
	firstName := contact.Name
	lastName := ""
	if len(nameParts) > 1 {
		firstName = strings.Join(nameParts[:len(nameParts)-1], " ")
		lastName = nameParts[len(nameParts)-1]
	}

	params := emailJSParams{
		ToName:      "Chanterelle member",
		Destination: "New Contact Form Submission",
		Firstname:   firstName,
		Lastname:    lastName,
		Email:       contact.Email,
		Message:     contact.Message,
	}

	return s.sendEmailJS(params)
}

// Keeping the old Mailchimp implementation for backward compatibility
func (s *NotificationService) SendAdminNotification(contact *models.Contact) error {
	// First try EmailJS if configured
	if s.cfg.EmailJSServiceID != "" && s.cfg.EmailJSTemplateID != "" && s.cfg.EmailJSUserID != "" {
		return s.SendNewContactNotification(contact)
	}

	// Fall back to Mailchimp
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

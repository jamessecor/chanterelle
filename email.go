package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// EmailJSConfig contains EmailJS configuration
var EmailJSConfig = struct {
	UserID     string
	Template   string
	PublicKey  string
	PrivateKey string
}{
	UserID:     "your-emailjs-user-id",
	Template:   "your-template-id",
	PublicKey:  "your-public-key",
	PrivateKey: "your-private-key",
}

// SendMail sends an email using EmailJS
func SendMail(to string, subject string, body string) error {
	// Create the request body
	data := map[string]interface{}{
		"to_email": to,
		"subject":  subject,
		"message":  body,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// Create the request
	req, err := http.NewRequest("POST", "https://api.emailjs.com/api/v1.0/email/send", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", EmailJSConfig.PrivateKey))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("emailjs API error: %s", string(bodyBytes))
	}

	return nil
}

// SendContactEmail sends an email with contact form data
func SendContactEmail(contact Contact) error {
	subject := "New Contact Form Submission"
	body := fmt.Sprintf(`
	Name: %s
	Email: %s
	Message: %s
	`,
		contact.Name,
		contact.Email,
		contact.Message,
	)

	return SendMail("james.secor@gmail.com", subject, body)
}

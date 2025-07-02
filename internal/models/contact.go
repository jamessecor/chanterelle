package models

import "time"

type Contact struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type VerificationCode struct {
	ID          int       `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	Code        string    `json:"code"`
	ExpiresAt   time.Time `json:"expires_at"`
}

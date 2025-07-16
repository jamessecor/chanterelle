package models

import "time"

type Contact struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" validate:"required,min=2,max=100"`
	Email     string    `json:"email" validate:"required,email"`
	Message   string    `json:"message" validate:"max=500"`
	CreatedAt time.Time `json:"created_at"`
}

type VerificationCode struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
}

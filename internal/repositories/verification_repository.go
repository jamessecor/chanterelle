package repositories

import (
	"context"
	"time"
)

type VerificationCode struct {
	ID        string    `bson:"_id,omitempty"`
	Code      string    `bson:"code"`
	Email     string    `bson:"email"`
	CreatedAt time.Time `bson:"created_at"`
	ExpiresAt time.Time `bson:"expires_at"`
}

type VerificationRepository interface {
	CreateVerificationCode(ctx context.Context, email, code string, expiry time.Duration) error
	GetCodeByEmail(ctx context.Context, email string) (*VerificationCode, error)
	DeleteCodeByEmail(ctx context.Context, email string) error
	DeleteExpiredCodes(ctx context.Context) error
}

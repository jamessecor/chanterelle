package services

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"log"
	"strconv"
	"time"

	"chanterelle/internal/config"
	"chanterelle/internal/models"
	"chanterelle/internal/repositories"
)

type VerificationService struct {
	cfg        *config.Config
	repository repositories.VerificationRepository
}

func NewVerificationService(cfg *config.Config, repository repositories.VerificationRepository) *VerificationService {
	return &VerificationService{
		cfg:        cfg,
		repository: repository,
	}
}

// GenerateRandomCode generates a random alphanumeric code of specified length
func (s *VerificationService) GenerateRandomCode() string {
	bytes := make([]byte, s.cfg.VerificationCodeLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return base32.StdEncoding.EncodeToString(bytes)[:s.cfg.VerificationCodeLength]
}

func (s *VerificationService) CreateVerificationCode(ctx context.Context, email string) (string, error) {
	// Generate a 6-digit numeric verification code (100000-999999)
	bytes := make([]byte, 4) // 4 bytes gives us 0-4294967295
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}

	num := int(binary.BigEndian.Uint32(bytes))%900000 + 100000 // Range: 100000-999999
	code := strconv.Itoa(num)

	// Create verification code
	verificationCode := &models.VerificationCode{
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(s.cfg.VerificationCodeExpiry),
	}

	if err := s.repository.CreateVerificationCode(ctx, verificationCode.Email, verificationCode.Code, s.cfg.VerificationCodeExpiry); err != nil {
		return "", err
	}

	return code, nil
}

func (s *VerificationService) GetCodeByEmail(ctx context.Context, email string) (*models.VerificationCode, error) {
	log.Println("Getting verification code for email:", email)
	verificationCode, err := s.repository.GetCodeByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	log.Println("Verification code for email:", email, "is", verificationCode.Code, "expires at", verificationCode.ExpiresAt.Format("2006-01-02 15:04:05"))

	// Check if code is expired
	if verificationCode.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("verification code has expired")
	}

	// Convert repository type to model type
	return &models.VerificationCode{
		Email:     verificationCode.Email,
		Code:      verificationCode.Code,
		ExpiresAt: verificationCode.ExpiresAt,
	}, nil
}

func (s *VerificationService) DeleteCodeByEmail(ctx context.Context, email string) error {
	return s.repository.DeleteCodeByEmail(ctx, email)
}

func (s *VerificationService) DeleteExpiredCodes(ctx context.Context) error {
	return s.repository.DeleteExpiredCodes(ctx)
}

package services

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"chanterelle/internal/config"
	"chanterelle/internal/models"
	"chanterelle/internal/repositories"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type VerificationService struct {
	config           *config.Config
	verificationRepo *repositories.VerificationCodeRepository
}

func NewVerificationService(config *config.Config, verificationRepo *repositories.VerificationCodeRepository) *VerificationService {
	return &VerificationService{
		config:           config,
		verificationRepo: verificationRepo,
	}
}

func (s *VerificationService) GenerateVerificationCode(phoneNumber string) (string, error) {
	code := generateRandomCode()

	verification := &models.VerificationCode{
		PhoneNumber: phoneNumber,
		Code:        code,
		ExpiresAt:   time.Now().Add(5 * time.Minute),
	}

	if err := s.verificationRepo.Create(verification); err != nil {
		return "", errors.Wrap(err, "failed to create verification code")
	}

	if err := s.SendVerificationCode(phoneNumber, code); err != nil {
		return "", errors.Wrap(err, "failed to send verification code")
	}

	return code, nil
}

func (s *VerificationService) VerifyCode(phoneNumber, code string) (string, error) {
	verification, err := s.verificationRepo.GetByCode(phoneNumber, code)
	if err != nil {
		return "", errors.Wrap(err, "failed to verify code")
	}
	if verification == nil {
		return "", errors.New("invalid verification code")
	}

	// Delete the verification code after successful verification
	if err := s.verificationRepo.DeleteByID(verification.ID); err != nil {
		return "", errors.Wrap(err, "failed to delete verification code")
	}

	// Generate JWT token
	claims := jwt.MapClaims{}
	claims["phone_number"] = phoneNumber
	claims["exp"] = time.Now().Add(35 * time.Minute).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", errors.Wrap(err, "failed to generate token")
	}

	return tokenString, nil
}

func (s *VerificationService) SendVerificationCode(phoneNumber, code string) error {
	twilioUrl := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", s.config.TwilioSID)
	data := url.Values{}
	data.Add("To", fmt.Sprintf("whatsapp:%s", phoneNumber))
	data.Add("From", s.config.TwilioNumber)
	data.Add("ContentSid", s.config.TwilioContentSID)
	data.Add("ContentVariables", fmt.Sprintf("{\"1\": \"%s\"}", code))

	req, err := http.NewRequest("POST", twilioUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	req.SetBasicAuth(s.config.TwilioSID, s.config.TwilioToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return errors.Errorf("twilio request failed: %s", string(body))
	}

	return nil
}

func (s *VerificationService) IsValidAdminPhoneNumber(phoneNumber string) bool {
	return slices.Contains(s.config.AvailableAdminPhoneNumbers, phoneNumber)
}

func generateRandomCode() string {
	bytes := make([]byte, 6)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}

	code := fmt.Sprintf("%06d", uint32(bytes[0])%1000000)
	return code
}

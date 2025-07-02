package handlers

import (
	"chanterelle/internal/config"
	"chanterelle/internal/services"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type VerificationHandler struct {
	verificationService *services.VerificationService
	config              *config.Config
}

func NewVerificationHandler(verificationService *services.VerificationService, config *config.Config) *VerificationHandler {
	return &VerificationHandler{verificationService: verificationService, config: config}
}

func (h *VerificationHandler) SendVerification(c *gin.Context) {
	var req struct {
		PhoneNumber string `json:"phoneNumber" validate:"required,e164"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.New().Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !isValidPhoneNumber(req.PhoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number format. Must be in E.164 format (e.g., +18025551234)"})
		return
	}

	if !h.verificationService.IsValidAdminPhoneNumber(req.PhoneNumber) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	_, err := h.verificationService.GenerateVerificationCode(req.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification code sent successfully"})
}

func (h *VerificationHandler) VerifyCode(c *gin.Context) {
	var req struct {
		PhoneNumber string `json:"phoneNumber" validate:"required,e164"`
		Code        string `json:"code" validate:"required,len=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.New().Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !isValidPhoneNumber(req.PhoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number format. Must be in E.164 format (e.g., +18025551234)"})
		return
	}

	token, err := h.verificationService.VerifyCode(req.PhoneNumber, req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid verification code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *VerificationHandler) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format. Expected 'Bearer <token>'"})
			c.Abort()
			return
		}

		tokenString = tokenString[7:] // Remove "Bearer " prefix

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(h.config.JWTSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired or invalid"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Store phone number in context for later use
		phoneNumber := claims["phone_number"].(string)
		c.Set("phone_number", phoneNumber)

		c.Next()
	}
}

func isValidPhoneNumber(phoneNumber string) bool {
	return strings.HasPrefix(phoneNumber, "+") && len(phoneNumber) >= 10
}

package handlers

import (
	"chanterelle/internal/config"
	"chanterelle/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/dgrijalva/jwt-go"
)

type VerificationHandler struct {
	verificationService *services.VerificationService
	config              *config.Config
}

func NewVerificationHandler(verificationService *services.VerificationService, config *config.Config) *VerificationHandler {
	return &VerificationHandler{
		verificationService: verificationService,
		config:              config,
	}
}

func (h *VerificationHandler) SendVerification(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.New().Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Only generate code for admin email
	if req.Email != h.config.AdminEmail {
		// Return success regardless of email
		c.JSON(http.StatusOK, gin.H{
			"message": "If the email was valid, you'll receive a verification code",
		})
		return
	}

	code, err := h.verificationService.CreateVerificationCode(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Store the code in the session for verification
	c.SetCookie("verification_code", code, 300, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "If the email was valid, you'll receive a verification code",
	})
}

func (h *VerificationHandler) VerifyCode(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required,len=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.New().Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Only accept verification for admin email
	if req.Email != h.config.AdminEmail {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid email"})
		return
	}

	// Get the stored code from the cookie
	storedCode, err := c.Cookie("verification_code")
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "No verification code found"})
		return
	}

	// Verify the code
	if req.Code != storedCode {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid verification code"})
		return
	}

	// Clear the cookie after successful verification
	c.SetCookie("verification_code", "", -1, "/", "", false, true)

	// Generate JWT token
	token, err := h.generateToken(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.Writer.Header().Set("X-Verified-Email", req.Email)

	c.JSON(http.StatusOK, gin.H{
		"message": "Verification successful",
		"token":   token,
	})
}

// generateToken creates a new JWT token for the given email
func (h *VerificationHandler) generateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
	})

	tokenString, err := token.SignedString([]byte(h.config.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (h *VerificationHandler) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.GetHeader("X-Verified-Email")
		if email == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Verify that the email is the admin email
		if email != h.config.AdminEmail {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid email"})
			c.Abort()
			return
		}

		// Store email in context for use by other handlers
		c.Set("email", email)
		c.Next()
	}
}

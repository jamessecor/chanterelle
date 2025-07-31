package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"chanterelle/internal/config"
	"chanterelle/internal/models"
	"chanterelle/internal/services"
)

type Handlers struct {
	contactService      *services.ContactService
	notificationService *services.NotificationService
	verificationService *services.VerificationService
	config              *config.Config
}

func NewHandlers(contactService *services.ContactService, notificationService *services.NotificationService, verificationService *services.VerificationService, config *config.Config) *Handlers {
	return &Handlers{
		contactService:      contactService,
		notificationService: notificationService,
		verificationService: verificationService,
		config:              config,
	}
}

func (h *Handlers) CreateContact(c *gin.Context) {
	var contact struct {
		Name    string `json:"name" binding:"required"`
		Email   string `json:"email" binding:"required,email"`
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.contactService.CreateContact(c.Request.Context(), contact.Name, contact.Email, contact.Message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.notificationService.AddToMailchimp(&models.Contact{
		Name:  contact.Name,
		Email: contact.Email,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Contact created successfully"})
}

func (h *Handlers) GetContacts(c *gin.Context) {
	contacts, err := h.contactService.GetContacts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contacts)
}

func (h *Handlers) GetContactByID(c *gin.Context) {
	id := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	contact, err := h.contactService.GetContactByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contact)
}

func (h *Handlers) UpdateContact(c *gin.Context) {
	id := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	var contact struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.contactService.UpdateContact(c.Request.Context(), id, contact.Name, contact.Email, contact.Message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact updated successfully"})
}

func (h *Handlers) DeleteContact(c *gin.Context) {
	id := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	if err := h.contactService.DeleteContact(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact deleted successfully"})
}

func (h *Handlers) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// The token should be in the format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be 'Bearer {token}'"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is what you expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(h.config.JWTSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Add the email from the token to the context
			if email, ok := claims["email"].(string); ok {
				c.Set("email", email)
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (h *Handlers) SendVerification(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Failed to bind JSON:", err)
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

	h.notificationService.SendVerificationCode(req.Email, code)

	c.JSON(http.StatusOK, gin.H{
		"message": "If the email was valid, you'll receive a verification code",
	})
}

func (h *Handlers) VerifyCode(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required,len=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Only accept verification for admin email
	if req.Email != h.config.AdminEmail {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid email"})
		return
	}

	// Verify the code using the repository
	verificationCode, err := h.verificationService.GetCodeByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid or expired verification code"})
		return
	}

	if verificationCode.Code != req.Code {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid verification code"})
		return
	}

	// Delete the verification code after successful verification
	if err := h.verificationService.DeleteCodeByEmail(c.Request.Context(), req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete verification"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": req.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
	})

	tokenString, err := token.SignedString([]byte(h.config.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set the verified email header for subsequent requests
	c.Writer.Header().Set("X-Verified-Email", req.Email)

	c.JSON(http.StatusOK, gin.H{
		"message": "Verification successful",
		"token":   tokenString,
	})
}

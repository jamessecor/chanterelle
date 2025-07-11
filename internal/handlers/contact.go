package handlers

import (
	"chanterelle/internal/models"
	"chanterelle/internal/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ContactHandler struct {
	contactService *services.ContactService
}

func NewContactHandler(contactService *services.ContactService) *ContactHandler {
	return &ContactHandler{contactService: contactService}
}

func (h *ContactHandler) CreateContact(c *gin.Context) {
	var contact models.Contact
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate the contact
	if contact.Name == "" || contact.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and email are required"})
		return
	}

	// Check email format
	if !isValidEmail(contact.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	createdContact, err := h.contactService.Create(&contact)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create contact"})
		return
	}

	c.JSON(http.StatusCreated, createdContact)
}

func isValidEmail(email string) bool {
	// Simple email validation
	return len(email) > 5 && strings.Contains(email, "@") && strings.Contains(email, ".")
}

func (h *ContactHandler) GetContacts(c *gin.Context) {
	contacts, err := h.contactService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contacts"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"contacts": contacts})
}

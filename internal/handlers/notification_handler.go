package handlers

import (
	"chanterelle/internal/config"
	"chanterelle/internal/models"
	"chanterelle/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationService *services.NotificationService
}

func NewNotificationHandler(cfg *config.Config) *NotificationHandler {
	return &NotificationHandler{
		notificationService: services.NewNotificationService(cfg),
	}
}

// NotifyAdmin handles the notification of new contact submissions
func (h *NotificationHandler) NotifyAdmin(c *gin.Context) {
	var contact models.Contact
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Add contact to Mailchimp list
	if err := h.notificationService.AddToMailchimp(&contact); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add contact to Mailchimp"})
		return
	}

	// Send admin notification
	if err := h.notificationService.SendNewContactNotification(&contact); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send admin notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Contact added to Mailchimp and admin notified successfully",
	})
}

// HandleContactSubmission handles the contact form submission and triggers notifications
func (h *NotificationHandler) HandleContactSubmission(c *gin.Context) {
	var contact models.Contact
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Add contact to Mailchimp list
	if err := h.notificationService.AddToMailchimp(&contact); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add contact to Mailchimp"})
		return
	}

	// Send admin notification
	if err := h.notificationService.SendNewContactNotification(&contact); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send admin notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Contact added to Mailchimp and admin notified successfully",
	})
}

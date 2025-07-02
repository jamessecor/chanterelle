package handlers

import (
	"chanterelle/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContactHandler struct {
	contactService *services.ContactService
}

func NewContactHandler(contactService *services.ContactService) *ContactHandler {
	return &ContactHandler{contactService: contactService}
}

func (h *ContactHandler) GetContacts(c *gin.Context) {
	contacts, err := h.contactService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"contacts": contacts})
}

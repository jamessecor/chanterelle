package services

import (
	"chanterelle/internal/models"
	"chanterelle/internal/repositories"
)

type ContactService struct {
	contactRepo *repositories.ContactRepository
}

func NewContactService(contactRepo *repositories.ContactRepository) *ContactService {
	return &ContactService{contactRepo: contactRepo}
}

func (s *ContactService) GetAll() ([]*models.Contact, error) {
	return s.contactRepo.GetAll()
}

func (s *ContactService) Create(contact *models.Contact) (*models.Contact, error) {
	return s.contactRepo.Create(contact)
}

package services

import (
	"context"

	"chanterelle/internal/repositories"
)

type ContactService struct {
	repository repositories.ContactRepository
}

func NewContactService(repository repositories.ContactRepository) *ContactService {
	return &ContactService{
		repository: repository,
	}
}

func (s *ContactService) CreateContact(ctx context.Context, name, email, message string) error {
	return s.repository.CreateContact(ctx, name, email, message)
}

func (s *ContactService) GetContacts(ctx context.Context) ([]repositories.Contact, error) {
	return s.repository.GetContacts(ctx)
}

func (s *ContactService) GetContactByID(ctx context.Context, id string) (repositories.Contact, error) {
	return s.repository.GetContactByID(ctx, id)
}

func (s *ContactService) UpdateContact(ctx context.Context, id string, name, email, message string) error {
	return s.repository.UpdateContact(ctx, id, name, email, message)
}

func (s *ContactService) DeleteContact(ctx context.Context, id string) error {
	return s.repository.DeleteContact(ctx, id)
}

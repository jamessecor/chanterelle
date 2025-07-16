package repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Contact struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string            `bson:"name"`
	Email     string            `bson:"email"`
	Message   string            `bson:"message"`
	CreatedAt time.Time         `bson:"created_at"`
}

type ContactRepository interface {
	CreateContact(ctx context.Context, name, email, message string) error
	GetContacts(ctx context.Context) ([]Contact, error)
	GetContactByID(ctx context.Context, id string) (Contact, error)
	UpdateContact(ctx context.Context, id string, name, email, message string) error
	DeleteContact(ctx context.Context, id string) error
}

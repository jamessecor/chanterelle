package repositories

import (
	"database/sql"
	"github.com/pkg/errors"
	"chanterelle/internal/models"
)

type ContactRepository struct {
	db *sql.DB
}

func NewContactRepository(db *sql.DB) *ContactRepository {
	return &ContactRepository{db: db}
}

func (r *ContactRepository) GetAll() ([]*models.Contact, error) {
	rows, err := r.db.Query(`
		SELECT id, name, email, message, created_at 
		FROM contacts
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query contacts")
	}
	defer rows.Close()

	var contacts []*models.Contact
	for rows.Next() {
		var contact models.Contact
		if err := rows.Scan(
			&contact.ID,
			&contact.Name,
			&contact.Email,
			&contact.Message,
			&contact.CreatedAt,
		); err != nil {
			return nil, errors.Wrap(err, "failed to scan contact")
		}
		contacts = append(contacts, &contact)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to iterate contacts")
	}

	return contacts, nil
}

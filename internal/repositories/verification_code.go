package repositories

import (
	"chanterelle/internal/models"
	"database/sql"

	"github.com/pkg/errors"
)

type VerificationCodeRepository struct {
	db *sql.DB
}

func NewVerificationCodeRepository(db *sql.DB) *VerificationCodeRepository {
	return &VerificationCodeRepository{db: db}
}

func (r *VerificationCodeRepository) Create(code *models.VerificationCode) error {
	_, err := r.db.Exec(`
		INSERT INTO verification_codes (email, code, expires_at)
		VALUES ($1, $2, $3)
	`, code.Email, code.Code, code.ExpiresAt)
	if err != nil {
		return errors.Wrap(err, "failed to create verification code")
	}
	return nil
}

func (r *VerificationCodeRepository) GetByCode(email, code string) (*models.VerificationCode, error) {
	var verificationCode models.VerificationCode

	err := r.db.QueryRow(`
		SELECT id, email, code, expires_at 
		FROM verification_codes 
		WHERE email = $1 AND code = $2 
		AND expires_at > NOW()
	`, email, code).Scan(
		&verificationCode.ID,
		&verificationCode.Email,
		&verificationCode.Code,
		&verificationCode.ExpiresAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to query verification code")
	}

	return &verificationCode, nil
}

func (r *VerificationCodeRepository) DeleteByID(id int) error {
	_, err := r.db.Exec(`
		DELETE FROM verification_codes WHERE id = $1
	`, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete verification code")
	}
	return nil
}

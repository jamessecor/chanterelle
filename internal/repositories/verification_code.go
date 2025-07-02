package repositories

import (
	"database/sql"
	"github.com/pkg/errors"
	"chanterelle/internal/models"
)

type VerificationCodeRepository struct {
	db *sql.DB
}

func NewVerificationCodeRepository(db *sql.DB) *VerificationCodeRepository {
	return &VerificationCodeRepository{db: db}
}

func (r *VerificationCodeRepository) Create(code *models.VerificationCode) error {
	_, err := r.db.Exec(`
		INSERT INTO verification_codes (phone_number, code, expires_at)
		VALUES ($1, $2, $3)
	`, code.PhoneNumber, code.Code, code.ExpiresAt)
	if err != nil {
		return errors.Wrap(err, "failed to create verification code")
	}
	return nil
}

func (r *VerificationCodeRepository) GetByCode(phoneNumber, code string) (*models.VerificationCode, error) {
	var verificationCode models.VerificationCode
	
	err := r.db.QueryRow(`
		SELECT id, phone_number, code, expires_at 
		FROM verification_codes 
		WHERE phone_number = $1 AND code = $2 
		AND expires_at > NOW()
	`, phoneNumber, code).Scan(
		&verificationCode.ID,
		&verificationCode.PhoneNumber,
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

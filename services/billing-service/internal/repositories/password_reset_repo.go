package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"billing-service/internal/models"
)

type PasswordResetRepository struct {
	db *sql.DB
}

func NewPasswordResetRepository(db *sql.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

func (r *PasswordResetRepository) Create(token *models.PasswordResetToken) error {
	err := r.db.QueryRow(`
		INSERT INTO password_reset_tokens (id, user_id, token_hash, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at
	`, token.ID, token.UserID, token.TokenHash, token.ExpiresAt.UTC()).Scan(&token.CreatedAt)
	if err != nil {
		return fmt.Errorf("create password reset token: %w", err)
	}

	token.CreatedAt = token.CreatedAt.UTC()
	return nil
}

func (r *PasswordResetRepository) GetValidByTokenHash(tokenHash string) (*models.PasswordResetToken, error) {
	token := &models.PasswordResetToken{}
	var usedAt sql.NullTime

	err := r.db.QueryRow(`
		SELECT id, user_id, token_hash, expires_at, used_at, created_at
		FROM password_reset_tokens
		WHERE token_hash = $1
		  AND used_at IS NULL
		  AND expires_at > NOW()
	`, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&usedAt,
		&token.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("password reset token not found")
		}
		return nil, fmt.Errorf("get password reset token: %w", err)
	}

	token.ExpiresAt = token.ExpiresAt.UTC()
	token.CreatedAt = token.CreatedAt.UTC()
	if usedAt.Valid {
		value := usedAt.Time.UTC()
		token.UsedAt = &value
	}

	return token, nil
}

func (r *PasswordResetRepository) MarkUsed(id string) error {
	result, err := r.db.Exec(`
		UPDATE password_reset_tokens
		SET used_at = NOW()
		WHERE id = $1
		  AND used_at IS NULL
	`, id)
	if err != nil {
		return fmt.Errorf("mark password reset token used: %w", err)
	}

	return ensureRowsAffected(result, fmt.Sprintf("password reset token not found: %s", id))
}

func (r *PasswordResetRepository) DeleteExpired() error {
	_, err := r.db.Exec(`
		DELETE FROM password_reset_tokens
		WHERE expires_at <= $1
		   OR used_at IS NOT NULL
	`, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("delete expired password reset tokens: %w", err)
	}

	return nil
}

package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Royal17x/flagr/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type tokenRepository struct {
	db *sqlx.DB
}

func NewTokenRepository(db *sqlx.DB) *tokenRepository {
	return &tokenRepository{db: db}
}

func (t *tokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	id := uuid.New()
	query := `
		INSERT INTO refresh_tokens 
		    (id, user_id, token_hash, expires_at, created_at)
			VALUES ($1, $2, $3, $4, NOW())`
	row := t.db.QueryRowContext(ctx, query, id, token.UserID, token.TokenHash, token.ExpiresAt)
	if row.Err() != nil {
		return fmt.Errorf("tokenRepository.Create: %w", row.Err())
	}
	return nil
}

func (t *tokenRepository) GetByTokenHash(ctx context.Context, hash string) (*domain.RefreshToken, error) {
	var token domain.RefreshToken
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at 
		FROM refresh_tokens
		WHERE token_hash = $1`
	err := t.db.GetContext(ctx, &token, query, hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("tokenRepository.GetByTokenHash: %w", err)
	}
	return &token, nil
}

func (t *tokenRepository) DeleteByTokenHash(ctx context.Context, hash string) error {
	query := `DELETE FROM refresh_tokens WHERE token_hash = $1;`
	res, err := t.db.ExecContext(ctx, query, hash)
	if err != nil {
		return fmt.Errorf("tokenRepository.DeleteByTokenHash: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("tokenRepository.DeleteByTokenHash: %w", err)
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (t *tokenRepository) DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1;`
	_, err := t.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("tokenRepository.DeleteAllByUserID: %w", err)
	}
	return nil
}

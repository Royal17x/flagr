package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/Royal17x/flagr/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type sdkKeyRepository struct {
	db *sqlx.DB
}

func NewSDKKeyRepository(db *sqlx.DB) domain.SDKKeyRepository {
	return &sdkKeyRepository{db: db}
}

func (s *sdkKeyRepository) Create(ctx context.Context, key *domain.SDKKey, rawKey string) error {
	id := uuid.New()
	var createdAt time.Time

	query := `
		INSERT INTO sdk_keys
		(id, key_hash, project_id, environment_id, name, created_by, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		RETURNING created_at
	`
	keyHash := sha256.Sum256([]byte(rawKey))
	keyHashStr := hex.EncodeToString(keyHash[:])
	row := s.db.QueryRowContext(ctx, query, id, keyHashStr, key.ProjectID, key.EnvironmentID, key.Name, key.CreatedBy, key.ExpiresAt)
	if err := row.Scan(&createdAt); err != nil {
		return fmt.Errorf("sdkKeyRepository.Create: %w", err)
	}
	return nil
}

func (s *sdkKeyRepository) GetByKeyHash(ctx context.Context, hash string) (*domain.SDKKey, error) {
	var sdkKey domain.SDKKey
	query := `
		SELECT *
		FROM sdk_keys
		WHERE key_hash = $1
	`
	err := s.db.GetContext(ctx, &sdkKey, query, hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("sdkKeyRepository.GetByKeyHash: %w", err)
	}
	return &sdkKey, nil
}

func (s *sdkKeyRepository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]*domain.SDKKey, error) {
	var sdkKeys []*domain.SDKKey
	query := `
		SELECT *
		FROM sdk_keys
		WHERE project_id = $1
	`
	err := s.db.SelectContext(ctx, &sdkKeys, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("sdkKeyRepository.ListByProject: %w", err)
	}
	return sdkKeys, nil
}

func (s *sdkKeyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM sdk_keys WHERE id = $1;`
	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("sdkKeyRepository.Delete: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("sdkKeyRepository.Delete: %w", err)
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

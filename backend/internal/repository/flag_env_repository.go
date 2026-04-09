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

type flagEnvironmentRepository struct {
	db *sqlx.DB
}

func NewFlagEnvironmentRepository(db *sqlx.DB) domain.FlagEnvironmentRepository {
	return &flagEnvironmentRepository{db: db}
}

func (f *flagEnvironmentRepository) Create(ctx context.Context, fe *domain.FlagEnvironment) (uuid.UUID, error) {
	fe.ID = uuid.New()
	query := `
		INSERT INTO flag_environments (id, flag_id, environment_id, enabled, rollout_percentage, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING updated_at`
	row := f.db.QueryRowContext(ctx, query, fe.ID, fe.FlagID, fe.EnvironmentID, fe.Enabled, fe.RolloutPercentage)
	if err := row.Scan(&fe.UpdatedAt); err != nil {
		return uuid.Nil, fmt.Errorf("flagEnvironmentRepository.Create: %w", err)
	}
	return fe.ID, nil
}

func (f *flagEnvironmentRepository) GetByFlagEnvID(ctx context.Context, flagID uuid.UUID, envID uuid.UUID) (*domain.FlagEnvironment, error) {
	var fe domain.FlagEnvironment
	query := `SELECT * FROM flag_environments WHERE flag_id = $1 AND environment_id = $2`
	err := f.db.GetContext(ctx, &fe, query, flagID, envID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("flagEnvironmentRepository.GetByFlagEnvID: %w", err)
	}
	return &fe, nil
}

func (f *flagEnvironmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM flag_environments WHERE id = $1`
	res, err := f.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("flagEnvironmentRepository.Delete: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("flagEnvironmentRepository.Delete: %w", err)
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (f *flagEnvironmentRepository) List(ctx context.Context, flagID uuid.UUID) ([]*domain.FlagEnvironment, error) {
	var fes []*domain.FlagEnvironment
	query := `SELECT * FROM flag_environments WHERE flag_id = $1`
	if err := f.db.SelectContext(ctx, &fes, query, flagID); err != nil {
		return nil, fmt.Errorf("flagEnvironmentRepository.List: %w", err)
	}
	return fes, nil
}

func (f *flagEnvironmentRepository) Upsert(ctx context.Context, fe *domain.FlagEnvironment) error {
	query := `
		INSERT INTO flag_environments (id, flag_id, environment_id, enabled, rollout_percentage, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (flag_id, environment_id)
		DO UPDATE SET
			enabled            = EXCLUDED.enabled,
			rollout_percentage = EXCLUDED.rollout_percentage,
			updated_at         = NOW()
		RETURNING updated_at`
	row := f.db.QueryRowContext(ctx, query, uuid.New(), fe.FlagID, fe.EnvironmentID, fe.Enabled, fe.RolloutPercentage)
	if err := row.Scan(&fe.UpdatedAt); err != nil {
		return fmt.Errorf("flagEnvironmentRepository.Upsert: %w", err)
	}
	return nil
}

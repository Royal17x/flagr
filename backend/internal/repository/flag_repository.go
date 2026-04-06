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

type flagRepository struct {
	db *sqlx.DB
}

func NewFlagRepository(db *sqlx.DB) domain.FlagRepository {
	return &flagRepository{db: db}
}

func (f *flagRepository) Create(ctx context.Context, flag *domain.Flag) (uuid.UUID, error) {
	flag.ID = uuid.New()
	query := `INSERT INTO flags
    (id, project_id, key, name, description, type, created_at, updated_at)
	VALUES($1,$2,$3,$4,$5,$6, NOW(), NOW())
	RETURNING created_at, updated_at`
	row := f.db.QueryRowContext(ctx, query, flag.ID, flag.ProjectID, flag.Key, flag.Name, flag.Description, flag.Type)
	if err := row.Scan(&flag.CreatedAt, &flag.UpdatedAt); err != nil {
		return uuid.Nil, fmt.Errorf("flagRepository.Create: %w", err)
	}
	return flag.ID, nil
}

func (f *flagRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Flag, error) {
	var flag domain.Flag
	query := `SELECT * FROM flags WHERE id = $1;`
	err := f.db.GetContext(ctx, &flag, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("flagRepository.GetByID: %w", err)
	}
	return &flag, nil
}

func (f *flagRepository) GetByKey(ctx context.Context, projectID uuid.UUID, key string) (*domain.Flag, error) {
	query := `SELECT id, project_id, key, name, description, type, created_at, updated_at 
		FROM flags 
		WHERE key = $1
		AND project_id = $2;`
	var flag domain.Flag
	err := f.db.GetContext(ctx, &flag, query, key, projectID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("flagRepository.GetByKey: %w", err)
	}
	return &flag, nil
}

func (f *flagRepository) List(ctx context.Context, projectID uuid.UUID) ([]*domain.Flag, error) {
	var flags []*domain.Flag
	query := `SELECT id, project_id, key, name, description, type, created_at, updated_at
		FROM flags
		WHERE project_id = $1;`
	err := f.db.SelectContext(ctx, &flags, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("flagRepository.List: %w", err)
	}
	return flags, nil
}

func (f *flagRepository) Update(ctx context.Context, flag *domain.Flag) error {
	query := `UPDATE flags
SET key = $1, name = $2, description = $3, type = $4, updated_at = NOW()
WHERE id = $5;`
	res, err := f.db.ExecContext(ctx, query, flag.Key, flag.Name, flag.Description, flag.Type, flag.ID)
	if err != nil {
		return fmt.Errorf("flagRepository.Update: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("flagRepository.Update: %w", err)
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (f *flagRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM flags WHERE id = $1;`
	res, err := f.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("flagRepository.Delete: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("flagRepository.Delete: %w", err)
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

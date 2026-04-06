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

type environmentRepository struct {
	db *sqlx.DB
}

func NewEnvironmentRepository(db *sqlx.DB) domain.EnvironmentRepository {
	return &environmentRepository{db: db}
}

func (e *environmentRepository) Create(ctx context.Context, env *domain.Environment) (uuid.UUID, error) {
	env.ID = uuid.New()
	query := `
		INSERT INTO environments (id, project_id, name, slug, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING created_at`
	row := e.db.QueryRowContext(ctx, query, env.ID, env.ProjectID, env.Name, env.Slug)
	if err := row.Scan(&env.CreatedAt); err != nil {
		return uuid.Nil, fmt.Errorf("environmentRepository.Create: %w", err)
	}
	return env.ID, nil
}

func (e *environmentRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) (*domain.Environment, error) {
	var env domain.Environment
	query := `SELECT * FROM environments WHERE project_id = $1`
	err := e.db.GetContext(ctx, &env, query, projectID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("environmentRepository.GetByProjectID: %w", err)
	}
	return &env, nil
}

func (e *environmentRepository) List(ctx context.Context, projectID uuid.UUID) ([]*domain.Environment, error) {
	var envs []*domain.Environment
	query := `
		SELECT * FROM environments
		WHERE project_id = $1
		ORDER BY created_at DESC
	`
	if err := e.db.SelectContext(ctx, &envs, query, projectID); err != nil {
		return nil, fmt.Errorf("environmentRepository.List: %w", err)
	}
	return envs, nil
}

func (e *environmentRepository) Update(ctx context.Context, env *domain.Environment) error {
	query := `
		UPDATE environments
		SET name = $1, slug = $2
		WHERE id = $3
	`
	res, err := e.db.ExecContext(ctx, query, env.Name, env.Slug, env.ID)
	if err != nil {
		return fmt.Errorf("environmentRepository.Update: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("environmentRepository.Update: %w", err)
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (e *environmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM environments WHERE id = $1`
	res, err := e.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("environmentRepository.Delete: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("environmentRepository.Delete: %w", err)
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

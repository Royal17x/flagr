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

type projectRepository struct {
	db *sqlx.DB
}

func NewProjectRepository(db *sqlx.DB) domain.ProjectRepository {
	return &projectRepository{db: db}
}

func (p *projectRepository) Create(ctx context.Context, project *domain.Project) (uuid.UUID, error) {
	project.ID = uuid.New()
	query := `
		INSERT INTO projects (id, organization_id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING created_at, updated_at
	`
	row := p.db.QueryRowContext(ctx, query, project.ID, project.OrganizationID, project.OrganizationID, project.Name, project.Description)
	if err := row.Scan(&project.CreatedAt, project.UpdatedAt); err != nil {
		return uuid.Nil, fmt.Errorf("projectRepository.Create: %w", err)
	}
	return project.ID, nil
}

func (p *projectRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	var project domain.Project
	query := `
		SELECT * FROM projects WHERE id = $1
	`
	err := p.db.GetContext(ctx, &project, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("projectRepository.GetByID: %w", err)
	}
	return &project, nil
}

func (p *projectRepository) List(ctx context.Context, orgID uuid.UUID) ([]*domain.Project, error) {
	var projects []*domain.Project
	query := `
		SELECT * FROM projects WHERE organization_id = $1 ORDER BY created_at DESC
	`
	if err := p.db.SelectContext(ctx, &projects, query, orgID); err != nil {
		return nil, fmt.Errorf("projectRepository.List: %w", err)
	}
	return projects, nil
}

func (p *projectRepository) Update(ctx context.Context, project *domain.Project) error {
	query := `
		UPDATE projects
		SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3
	`
	res, err := p.db.ExecContext(ctx, query, project.Name, project.Description, project.ID)
	if err != nil {
		return fmt.Errorf("projectRepository.Update: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("projectRepository.Update: %w", err)
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (p *projectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM projects WHERE id = $1
	`
	res, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("projectRepository.Delete: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("projectRepository.Delete: %w", err)

	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

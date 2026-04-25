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

type organizationRepository struct {
	db *sqlx.DB
}

func NewOrganizationRepository(db *sqlx.DB) domain.OrganizationRepository {
	return &organizationRepository{db: db}
}

func (o *organizationRepository) Create(ctx context.Context, org *domain.Organization) (uuid.UUID, error) {
	org.ID = uuid.New()
	query := `INSERT INTO organizations
    (id, name, slug, created_at, updated_at)
	VALUES($1,$2,$3, NOW(), NOW())
	RETURNING created_at, updated_at
	`
	row := o.db.QueryRowContext(ctx, query, org.ID, org.Name, org.Slug)
	if err := row.Scan(&org.CreatedAt, &org.UpdatedAt); err != nil {
		return uuid.Nil, fmt.Errorf("organizationRepository.Create: %w", err)
	}
	return org.ID, nil
}

func (o *organizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
	var org domain.Organization
	query := `SELECT id, name, slug, created_at, updated_at
		FROM organizations
		WHERE id = $1
	`
	err := o.db.GetContext(ctx, &org, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("organizationRepository.GetByID: %w", err)
	}
	return &org, nil
}

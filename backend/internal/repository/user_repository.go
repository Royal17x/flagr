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

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *userRepository {
	return &userRepository{db: db}
}

func (u *userRepository) Create(ctx context.Context, user *domain.User) (uuid.UUID, error) {
	id := uuid.New()
	query := `
		INSERT INTO users 
		    (id, org_id, email, password_hash, created_at, updated_at)
			VALUES($1, $2, $3, $4, NOW(), NOW())
			RETURNING created_at, updated_at`
	row := u.db.QueryRowContext(ctx, query, id, user.OrgID, user.Email, user.PasswordHash)
	if err := row.Scan(&user.CreatedAt, &user.UpdatedAt); err != nil {
		return uuid.Nil, fmt.Errorf("userRepository.Create: %w", err)
	}
	return id, nil
}

func (u *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	query := `
		SELECT id, org_id, email, password_hash, created_at, updated_at 
		FROM users
		WHERE email = $1`
	err := u.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("userRepository.GetByEmail: %w", err)
	}
	return &user, nil
}

func (u *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, org_id, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1`
	var user domain.User
	err := u.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("userRepository.GetByID: %w", err)
	}
	return &user, nil
}

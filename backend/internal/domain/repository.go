package domain

import (
	"context"
	"github.com/google/uuid"
)

type ProjectRepository interface {
	Create(ctx context.Context, project *Project) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Project, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, orgID uuid.UUID) ([]*Project, error)
	Update(ctx context.Context, project *Project) error
}

type EnvironmentRepository interface {
	Create(ctx context.Context, env *Environment) (uuid.UUID, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID) (*Environment, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, projectID uuid.UUID) ([]*Environment, error)
	Update(ctx context.Context, env *Environment) error
}
type FlagEnvironmentRepository interface {
	Create(ctx context.Context, fe *FlagEnvironment) (uuid.UUID, error)
	GetByFlagEnvID(ctx context.Context, flagID uuid.UUID, envID uuid.UUID) (*FlagEnvironment, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, flagID uuid.UUID) ([]*FlagEnvironment, error)
	Upsert(ctx context.Context, fe *FlagEnvironment) error
}

type FlagRepository interface {
	Create(ctx context.Context, flag *Flag) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Flag, error)
	GetByKey(ctx context.Context, projectID uuid.UUID, key string) (*Flag, error)
	List(ctx context.Context, projectID uuid.UUID) ([]*Flag, error)
	Update(ctx context.Context, flag *Flag) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type UserRepository interface {
	Create(ctx context.Context, user *User) (uuid.UUID, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
}
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *RefreshToken) error
	GetByTokenHash(ctx context.Context, hash string) (*RefreshToken, error)
	DeleteByTokenHash(ctx context.Context, hash string) error
	DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error
}

package port

import (
	"context"

	"github.com/Royal17x/flagr/backend/internal/domain"
	"github.com/google/uuid"
)

type FlagServiceInterface interface {
	CreateFlag(ctx context.Context, flag *domain.Flag) (uuid.UUID, error)
	GetFlag(ctx context.Context, id uuid.UUID) (*domain.Flag, error)
	ListFlags(ctx context.Context, projectID uuid.UUID) ([]*domain.Flag, error)
	UpdateFlag(ctx context.Context, flag *domain.Flag) error
	DeleteFlag(ctx context.Context, id uuid.UUID) error
	EvaluateFlag(ctx context.Context, flagKey string, projectID uuid.UUID, environmentID uuid.UUID) (bool, error)
	ToggleFlag(ctx context.Context, flagID uuid.UUID, envID uuid.UUID, enabled bool) error
	GetFlagEnvironment(ctx context.Context, flagID uuid.UUID, envID uuid.UUID) (*domain.FlagEnvironment, error)
}

type AuthServiceInterface interface {
	Register(ctx context.Context, email, password, orgName string) (TokenPair, error)
	Login(ctx context.Context, email, password string) (TokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
	ValidateAccessToken(tokenStr string) (*Claims, error)
}

type ProjectServiceInterface interface {
	CreateProject(ctx context.Context, project *domain.Project) (uuid.UUID, error)
	GetProject(ctx context.Context, id uuid.UUID) (*domain.Project, error)
	ListProjects(ctx context.Context, orgID uuid.UUID) ([]*domain.Project, error)
	UpdateProject(ctx context.Context, project *domain.Project) error
	DeleteProject(ctx context.Context, id uuid.UUID) error
}

type EnvironmentServiceInterface interface {
	CreateEnvironment(ctx context.Context, env *domain.Environment) (uuid.UUID, error)
	ListEnvironments(ctx context.Context, projectID uuid.UUID) ([]*domain.Environment, error)
	DeleteEnvironment(ctx context.Context, id uuid.UUID) error
}

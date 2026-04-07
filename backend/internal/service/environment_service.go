package service

import (
	"context"
	"fmt"

	"github.com/Royal17x/flagr/backend/internal/domain"
	"github.com/google/uuid"
)

type EnvironmentService struct {
	envs     domain.EnvironmentRepository
	projects domain.ProjectRepository
}

func NewEnvironmentService(envs domain.EnvironmentRepository, projects domain.ProjectRepository) *EnvironmentService {
	return &EnvironmentService{envs: envs, projects: projects}
}

func (s *EnvironmentService) CreateEnvironment(ctx context.Context, env *domain.Environment) (uuid.UUID, error) {
	_, err := s.projects.GetByID(ctx, env.ProjectID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("EnvironmentService.CreateEnvironment: project not found: %w", err)
	}
	id, err := s.envs.Create(ctx, env)
	if err != nil {
		return uuid.Nil, fmt.Errorf("EnvironmentService.CreateEnvironment: %w", err)
	}
	return id, nil
}

func (s *EnvironmentService) ListEnvironments(ctx context.Context, projectID uuid.UUID) ([]*domain.Environment, error) {
	envs, err := s.envs.List(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("EnvironmentService.ListEnvironments: %w", err)
	}
	return envs, nil
}

func (s *EnvironmentService) DeleteEnvironment(ctx context.Context, id uuid.UUID) error {
	if err := s.envs.Delete(ctx, id); err != nil {
		return fmt.Errorf("EnvironmentService.DeleteEnvironment: %w", err)
	}
	return nil
}

package service

import (
	"context"
	"fmt"

	"github.com/Royal17x/flagr/backend/internal/domain"
	"github.com/google/uuid"
)

type ProjectService struct {
	projects domain.ProjectRepository
}

func NewProjectService(projects domain.ProjectRepository) *ProjectService {
	return &ProjectService{projects: projects}
}

func (s *ProjectService) CreateProject(ctx context.Context, project *domain.Project) (uuid.UUID, error) {
	id, err := s.projects.Create(ctx, project)
	if err != nil {
		return uuid.Nil, fmt.Errorf("ProjectService.CreateProject: %w", err)
	}
	return id, nil
}

func (s *ProjectService) GetProject(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	project, err := s.projects.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ProjectService.GetProject: %w", err)
	}
	return project, nil
}

func (s *ProjectService) ListProjects(ctx context.Context, orgID uuid.UUID) ([]*domain.Project, error) {
	projects, err := s.projects.List(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("ProjectService.ListProjects: %w", err)
	}
	return projects, nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, project *domain.Project) error {
	if err := s.projects.Update(ctx, project); err != nil {
		return fmt.Errorf("ProjectService.UpdateProject: %w", err)
	}
	return nil
}

func (s *ProjectService) DeleteProject(ctx context.Context, id uuid.UUID) error {
	if err := s.projects.Delete(ctx, id); err != nil {
		return fmt.Errorf("ProjectService.DeleteProject: %w", err)
	}
	return nil
}

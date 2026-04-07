package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Royal17x/flagr/backend/internal/domain"
	"github.com/google/uuid"
)

type FlagService struct {
	flags    domain.FlagRepository
	projects domain.ProjectRepository
	flagEnvs domain.FlagEnvironmentRepository
}

func NewFlagService(flags domain.FlagRepository, projects domain.ProjectRepository, flagEnvs domain.FlagEnvironmentRepository) *FlagService {
	return &FlagService{flags: flags, projects: projects, flagEnvs: flagEnvs}
}

func (f *FlagService) CreateFlag(ctx context.Context, flag *domain.Flag) (uuid.UUID, error) {
	_, err := f.projects.GetByID(ctx, flag.ProjectID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("FlagService.CreateFlag: %w", err)
	}
	_, err = f.flags.GetByKey(ctx, flag.ProjectID, flag.Key)
	if err == nil {
		return uuid.Nil, domain.ErrAlreadyExists
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return uuid.Nil, fmt.Errorf("FlagService.CreateFlag: %w", err)
	}
	return f.flags.Create(ctx, flag)
}

func (f *FlagService) GetFlag(ctx context.Context, id uuid.UUID) (*domain.Flag, error) {
	gotFlag, err := f.flags.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("FlagService.GetFlag: %w", err)
	}
	return gotFlag, nil
}

func (f *FlagService) ListFlags(ctx context.Context, projectID uuid.UUID) ([]*domain.Flag, error) {
	flags, err := f.flags.List(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("FlagService.ListFlags: %w", err)
	}
	return flags, nil
}

func (f *FlagService) UpdateFlag(ctx context.Context, flag *domain.Flag) error {
	_, err := f.flags.GetByID(ctx, flag.ID)
	if err != nil {
		return fmt.Errorf("FlagService.UpdateFlag: %w", err)
	}
	return f.flags.Update(ctx, flag)
}

func (f *FlagService) DeleteFlag(ctx context.Context, id uuid.UUID) error {
	return f.flags.Delete(ctx, id)
}

func (f *FlagService) EvaluateFlag(ctx context.Context, flagKey string, projectID uuid.UUID, environmentID uuid.UUID) (bool, error) {
	gotFlag, err := f.flags.GetByKey(ctx, projectID, flagKey)
	if err != nil {
		return false, fmt.Errorf("FlagService.EvaluateFlag: %w", err)
	}

	fe, err := f.flagEnvs.GetByFlagEnvID(ctx, gotFlag.ID, environmentID)
	if err != nil {
		return false, fmt.Errorf("FlagService.EvaluateFlag: %w", err)
	}
	return fe.Enabled, nil
}

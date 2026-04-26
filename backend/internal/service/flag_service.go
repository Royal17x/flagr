package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Royal17x/flagr/backend/internal/metrics"
	"github.com/Royal17x/flagr/backend/internal/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"log/slog"
	"strconv"
	"time"

	"github.com/Royal17x/flagr/backend/internal/cache"
	"github.com/Royal17x/flagr/backend/internal/domain"
	"github.com/google/uuid"
)

type FlagService struct {
	flags    domain.FlagRepository
	projects domain.ProjectRepository
	flagEnvs domain.FlagEnvironmentRepository
	envs     domain.EnvironmentRepository
	cache    cache.FlagCache
	audit    *AuditService
}

func NewFlagService(
	flags domain.FlagRepository,
	projects domain.ProjectRepository,
	flagEnvs domain.FlagEnvironmentRepository,
	envs domain.EnvironmentRepository,
	cache cache.FlagCache,
	audit *AuditService,
) *FlagService {
	return &FlagService{
		flags:    flags,
		projects: projects,
		flagEnvs: flagEnvs,
		envs:     envs,
		cache:    cache,
		audit:    audit,
	}
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
	id, err := f.flags.Create(ctx, flag)
	if err != nil {
		return uuid.Nil, fmt.Errorf("FlagService.CreateFlag: %w", err)
	}

	actorID, _ := ctx.Value(middleware.UserIDKey).(uuid.UUID)
	f.audit.RecordFlagCreated(ctx, actorID, id, flag.ProjectID)
	envs, err := f.envs.List(ctx, flag.ProjectID)

	if err != nil {
		slog.Warn("failed to auto-create flag environments", "flag_id", id, "error", err)
		return id, nil
	}
	for _, env := range envs {
		if err := f.flagEnvs.Upsert(ctx, &domain.FlagEnvironment{
			FlagID:            id,
			EnvironmentID:     env.ID,
			Enabled:           false,
			RolloutPercentage: 0,
		}); err != nil {
			slog.Warn("failed to auto-create flag environment",
				"flag_id", id,
				"env_id", env.ID,
				"error", err)
		}
	}
	return id, nil
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
	gotFlag, err := f.flags.GetByID(ctx, flag.ID)
	if err != nil {
		return fmt.Errorf("FlagService.UpdateFlag: %w", err)
	}
	if err = f.flags.Update(ctx, flag); err != nil {
		return fmt.Errorf("FlagService.UpdateFlag: %w", err)
	}

	_ = f.cache.InvalidateFlag(ctx, flag.Key, flag.ProjectID)

	actorID, _ := ctx.Value(middleware.UserIDKey).(uuid.UUID)
	f.audit.RecordFlagUpdated(ctx, actorID, gotFlag.ID, flag.ProjectID, *flag)
	return nil
}

func (f *FlagService) DeleteFlag(ctx context.Context, id uuid.UUID) error {
	gotFlag, err := f.flags.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("FlagService.DeleteFlag: %w", err)
	}
	if err = f.flags.Delete(ctx, id); err != nil {
		return fmt.Errorf("FlagService.DeleteFlag: %w", err)
	}
	actorID, _ := ctx.Value(middleware.UserIDKey).(uuid.UUID)
	f.audit.RecordFlagDeleted(ctx, actorID, gotFlag.ID, gotFlag.ProjectID)
	return nil
}

func (f *FlagService) EvaluateFlag(ctx context.Context, flagKey string, projectID uuid.UUID, environmentID uuid.UUID) (bool, error) {
	ctx, span := otel.Tracer("flagr").Start(ctx, "FlagService.EvaluateFlag")
	defer span.End()

	span.SetAttributes(
		attribute.String("flag.key", flagKey),
		attribute.String("project.id", projectID.String()),
		attribute.String("environment.id", environmentID.String()),
	)

	start := time.Now()
	defer func() {
		metrics.FlagEvaluationDuration.Observe(time.Since(start).Seconds())
	}()
	enabled, err := f.cache.GetEvaluation(ctx, flagKey, projectID, environmentID)
	if err == nil {
		span.SetAttributes(attribute.String("cache.result", "hit"))
		metrics.CacheHitsTotal.Inc()
		metrics.FlagEvaluationsTotal.WithLabelValues(flagKey, strconv.FormatBool(enabled), "cache").Inc()
		return enabled, nil
	}
	if !errors.Is(err, cache.ErrCacheMiss) {
		slog.Warn("cache unavailable, falling back to db",
			"error", err,
			"flag", flagKey,
		)
	}
	span.SetAttributes(attribute.String("cache.result", "miss"))
	metrics.CacheMissesTotal.Inc()
	gotFlag, err := f.flags.GetByKey(ctx, projectID, flagKey)
	if err != nil {
		return false, fmt.Errorf("FlagService.EvaluateFlag: %w", err)
	}
	fe, err := f.flagEnvs.GetByFlagEnvID(ctx, gotFlag.ID, environmentID)
	if err != nil {
		return false, fmt.Errorf("FlagService.EvaluateFlag: %w", err)
	}

	_ = f.cache.SetEvaluation(ctx, flagKey, projectID, environmentID, fe.Enabled)
	metrics.FlagEvaluationsTotal.WithLabelValues(flagKey, strconv.FormatBool(enabled), "db").Inc()
	return fe.Enabled, nil
}

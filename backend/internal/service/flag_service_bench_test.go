package service_test

import (
	"context"
	"testing"

	"github.com/Royal17x/flagr/backend/internal/cache"
	"github.com/Royal17x/flagr/backend/internal/domain"
	"github.com/Royal17x/flagr/backend/internal/repository"
	"github.com/Royal17x/flagr/backend/internal/service"
	"github.com/Royal17x/flagr/backend/internal/testhelpers"
	"github.com/google/uuid"
)

func setupBenchData(b *testing.B) (svcWithCache *service.FlagService, svcWithoutCache *service.FlagService, flagKey string, projectID, envID uuid.UUID) {
	b.Helper()
	db := testhelpers.NewTestPostgres(b)
	redisClient := testhelpers.NewTestRedis(b)

	flagRepo := repository.NewFlagRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	flagEnvRepo := repository.NewFlagEnvironmentRepository(db)
	flagCache := cache.NewFlagCache(redisClient)
	nilCache := &testhelpers.NilCache{}

	ctx := context.Background()
	orgID := uuid.New()
	projectID = uuid.New()
	envID = uuid.New()

	db.MustExecContext(ctx, `INSERT INTO organizations (id, name, slug) VALUES ($1, $2, $3)`,
		orgID, "bench", "bench")
	db.MustExecContext(ctx, `INSERT INTO projects (id, organization_id, name, description) VALUES ($1, $2, $3, $4)`,
		projectID, orgID, "bench", "bench")
	db.MustExecContext(ctx, `INSERT INTO environments (id, project_id, name, slug) VALUES ($1, $2, $3, $4)`,
		envID, projectID, "prod", "prod")

	flag := &domain.Flag{ProjectID: projectID, Key: "bench-flag", Name: "Bench", Type: domain.FlagTypeBoolean}
	flagID, _ := flagRepo.Create(ctx, flag)
	flagEnvRepo.Upsert(ctx, &domain.FlagEnvironment{
		FlagID: flagID, EnvironmentID: envID, Enabled: true,
	})

	svcWithCache = service.NewFlagService(flagRepo, projectRepo, flagEnvRepo, flagCache)
	svcWithoutCache = service.NewFlagService(flagRepo, projectRepo, flagEnvRepo, nilCache)

	svcWithCache.EvaluateFlag(ctx, "bench-flag", projectID, envID)

	return svcWithCache, svcWithoutCache, "bench-flag", projectID, envID
}

func BenchmarkEvaluateFlag_WithCache(b *testing.B) {
	svc, _, flagKey, projectID, envID := setupBenchData(b)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		svc.EvaluateFlag(ctx, flagKey, projectID, envID)
	}
}

func BenchmarkEvaluateFlag_WithoutCache(b *testing.B) {
	_, svc, flagKey, projectID, envID := setupBenchData(b)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		svc.EvaluateFlag(ctx, flagKey, projectID, envID)
	}
}

func BenchmarkEvaluateFlag_WithCache_Parallel(b *testing.B) {
	svc, _, flagKey, projectID, envID := setupBenchData(b)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			svc.EvaluateFlag(ctx, flagKey, projectID, envID)
		}
	})
}

func BenchmarkEvaluateFlag_WithoutCache_Parallel(b *testing.B) {
	_, svc, flagKey, projectID, envID := setupBenchData(b)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			svc.EvaluateFlag(ctx, flagKey, projectID, envID)
		}
	})
}

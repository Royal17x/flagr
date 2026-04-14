package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Royal17x/flagr/backend/internal/cache"
	"github.com/Royal17x/flagr/backend/internal/domain"
	"github.com/Royal17x/flagr/backend/internal/service"
	"github.com/Royal17x/flagr/backend/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFlagService_CreateFlag_Success(t *testing.T) {
	// Arrange
	flagRepo := new(mocks.MockFlagRepository)
	projectRepo := new(mocks.MockProjectRepository)
	flagEnvRepo := new(mocks.MockFlagEnvironmentRepository)
	nilCache := &cache.NilCache{}

	projectID := uuid.New()
	expectedID := uuid.New()

	// Expectations
	projectRepo.On("GetByID", mock.Anything, projectID).
		Return(&domain.Project{ID: projectID}, nil)
	flagRepo.On("GetByKey", mock.Anything, projectID, "new-flag").
		Return(nil, domain.ErrNotFound)
	flagRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Flag")).
		Return(expectedID, nil)

	svc := service.NewFlagService(flagRepo, projectRepo, flagEnvRepo, nilCache)

	// Act
	id, err := svc.CreateFlag(context.Background(), &domain.Flag{
		ProjectID: projectID,
		Key:       "new-flag",
		Name:      "New Flag",
		Type:      domain.FlagTypeBoolean,
	})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)
	flagRepo.AssertExpectations(t)
}

func TestFlagService_CreateFlag_AlreadyExists(t *testing.T) {
	// Arrange
	flagRepo := new(mocks.MockFlagRepository)
	projectRepo := new(mocks.MockProjectRepository)
	flagEnvRepo := new(mocks.MockFlagEnvironmentRepository)
	nilCache := &cache.NilCache{}

	projectID := uuid.New()
	flag := &domain.Flag{ProjectID: projectID, Key: "existing-flag"}

	//Expectations
	projectRepo.On("GetByID", mock.Anything, projectID).
		Return(&domain.Project{ID: projectID}, nil)
	flagRepo.On("GetByKey", mock.Anything, projectID, "existing-flag").
		Return(&domain.Flag{}, nil)

	// Act
	svc := service.NewFlagService(flagRepo, projectRepo, flagEnvRepo, nilCache)

	id, err := svc.CreateFlag(context.Background(), flag)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrAlreadyExists, err)
	assert.Equal(t, uuid.Nil, id)
	flagRepo.AssertExpectations(t)
}

func TestFlagService_CreateFlag_ProjectNotFound(t *testing.T) {
	// Arrange
	flagRepo := new(mocks.MockFlagRepository)
	projectRepo := new(mocks.MockProjectRepository)
	flagEnvRepo := new(mocks.MockFlagEnvironmentRepository)
	nilCache := &cache.NilCache{}

	projectID := uuid.New()

	// Expectations
	projectRepo.On("GetByID", mock.Anything, projectID).
		Return(nil, domain.ErrNotFound)

	// Act
	svc := service.NewFlagService(flagRepo, projectRepo, flagEnvRepo, nilCache)
	id, err := svc.CreateFlag(context.Background(), &domain.Flag{ProjectID: projectID})

	// Assert
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
	flagRepo.AssertExpectations(t)
}

func TestFlagService_EvaluateFlag_CacheHit(t *testing.T) {
	// Arrange
	flagRepo := new(mocks.MockFlagRepository)
	projectRepo := new(mocks.MockProjectRepository)
	flagEnvRepo := new(mocks.MockFlagEnvironmentRepository)
	flagCache := new(mocks.MockFlagCache)

	projectID := uuid.New()
	envID := uuid.New()

	// Expecations
	flagCache.On("GetEvaluation", mock.Anything, "existing-flag", projectID, envID).
		Return(true, nil)

	// Act
	svc := service.NewFlagService(flagRepo, projectRepo, flagEnvRepo, flagCache)
	enabled, err := svc.EvaluateFlag(context.Background(), "existing-flag", projectID, envID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, true, enabled)
	flagRepo.AssertNotCalled(t, "GetByKey")
}

func TestFlagService_GetFlag_Success(t *testing.T) {
	// Arrange
	flagRepo := new(mocks.MockFlagRepository)
	projectRepo := new(mocks.MockProjectRepository)
	flagEnvRepo := new(mocks.MockFlagEnvironmentRepository)
	nilCache := &cache.NilCache{}

	id := uuid.New()

	// Expectations
	flagRepo.On("GetByID", mock.Anything, id).
		Return(&domain.Flag{ID: id, Key: "test"}, nil)

	// Act
	svc := service.NewFlagService(flagRepo, projectRepo, flagEnvRepo, nilCache)
	got, err := svc.GetFlag(context.Background(), id)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, id, got.ID)
	flagRepo.AssertExpectations(t)
}

func TestFlagService_GetFlag_NotFound(t *testing.T) {
	// Arrange
	flagRepo := new(mocks.MockFlagRepository)
	projectRepo := new(mocks.MockProjectRepository)
	flagEnvRepo := new(mocks.MockFlagEnvironmentRepository)
	nilCache := &cache.NilCache{}

	// Expectations
	flagRepo.On("GetByID", mock.Anything, mock.Anything).
		Return(nil, domain.ErrNotFound)

	// Act
	svc := service.NewFlagService(flagRepo, projectRepo, flagEnvRepo, nilCache)
	_, err := svc.GetFlag(context.Background(), uuid.New())

	// Assert
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestFlagService_EvaluateFlag_CacheMiss_FallbackToDB(t *testing.T) {
	// Arrange
	flagRepo := new(mocks.MockFlagRepository)
	projectRepo := new(mocks.MockProjectRepository)
	flagEnvRepo := new(mocks.MockFlagEnvironmentRepository)
	flagCache := new(mocks.MockFlagCache)

	projectID := uuid.New()
	envID := uuid.New()
	flagID := uuid.New()

	// Expectations
	flagCache.On("GetEvaluation", mock.Anything, "my-flag", projectID, envID).
		Return(false, cache.ErrCacheMiss)
	flagRepo.On("GetByKey", mock.Anything, projectID, "my-flag").
		Return(&domain.Flag{ID: flagID, Key: "my-flag"}, nil)
	flagEnvRepo.On("GetByFlagEnvID", mock.Anything, flagID, envID).
		Return(&domain.FlagEnvironment{Enabled: true}, nil)
	flagCache.On("SetEvaluation", mock.Anything, "my-flag", projectID, envID, true).
		Return(nil)

	// Act
	svc := service.NewFlagService(flagRepo, projectRepo, flagEnvRepo, flagCache)
	enabled, err := svc.EvaluateFlag(context.Background(), "my-flag", projectID, envID)

	// Assert
	assert.NoError(t, err)
	assert.True(t, enabled)
	flagRepo.AssertExpectations(t)
	flagCache.AssertExpectations(t)
}

func TestFlagService_EvaluateFlag_CacheError_FallbackToDB(t *testing.T) {
	flagRepo := new(mocks.MockFlagRepository)
	projectRepo := new(mocks.MockProjectRepository)
	flagEnvRepo := new(mocks.MockFlagEnvironmentRepository)
	flagCache := new(mocks.MockFlagCache)

	projectID := uuid.New()
	envID := uuid.New()
	flagID := uuid.New()

	// Expectations
	flagCache.On("GetEvaluation", mock.Anything, "my-flag", projectID, envID).
		Return(false, errors.New("redis: connection refused"))
	flagRepo.On("GetByKey", mock.Anything, projectID, "my-flag").
		Return(&domain.Flag{ID: flagID}, nil)
	flagEnvRepo.On("GetByFlagEnvID", mock.Anything, flagID, envID).
		Return(&domain.FlagEnvironment{Enabled: false}, nil)
	flagCache.On("SetEvaluation", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

		// Act
	svc := service.NewFlagService(flagRepo, projectRepo, flagEnvRepo, flagCache)
	enabled, err := svc.EvaluateFlag(context.Background(), "my-flag", projectID, envID)

	// Assert
	assert.NoError(t, err)
	assert.False(t, enabled)
}

package repository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Royal17x/flagr/backend/internal/domain"
	"github.com/Royal17x/flagr/backend/internal/repository"
	"github.com/Royal17x/flagr/backend/internal/testhelpers"
)

func seedProject(t *testing.T, db *sqlx.DB) (orgID, projectID uuid.UUID) {
	t.Helper()
	ctx := context.Background()
	orgID = uuid.New()
	projectID = uuid.New()
	db.MustExecContext(ctx,
		`INSERT INTO organizations (id, name, slug) VALUES ($1, $2, $3)`,
		orgID, "test-org", "test-org-"+orgID.String(),
	)
	db.MustExecContext(ctx,
		`INSERT INTO projects (id, organization_id, name, description) VALUES ($1, $2, $3, $4)`,
		projectID, orgID, "test-project", "desc",
	)
	return orgID, projectID
}

func TestFlagRepository_Create(t *testing.T) {
	// Arrange
	db := testhelpers.NewTestPostgres(t)
	repo := repository.NewFlagRepository(db)
	ctx := context.Background()
	_, projectID := seedProject(t, db)

	flag := &domain.Flag{
		ProjectID:   projectID,
		Key:         "test-flag",
		Name:        "Test Flag",
		Description: "desc",
		Type:        domain.FlagTypeBoolean,
	}

	// Act
	id, err := repo.Create(ctx, flag)

	// Assert
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)
	assert.NotZero(t, flag.CreatedAt)
}

func TestFlagRepository_GetByID(t *testing.T) {
	// Arrange
	db := testhelpers.NewTestPostgres(t)
	repo := repository.NewFlagRepository(db)
	ctx := context.Background()
	_, projectID := seedProject(t, db)

	flag := &domain.Flag{
		ProjectID: projectID, Key: "get-by-id-flag",
		Name: "Get By ID", Type: domain.FlagTypeBoolean,
	}
	id, err := repo.Create(ctx, flag)
	require.NoError(t, err)

	// Act
	got, err := repo.GetByID(ctx, id)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "get-by-id-flag", got.Key)
	assert.Equal(t, "Get By ID", got.Name)
	assert.Equal(t, projectID, got.ProjectID)
}

func TestFlagRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	db := testhelpers.NewTestPostgres(t)
	repo := repository.NewFlagRepository(db)
	ctx := context.Background()

	// Act
	_, err := repo.GetByID(ctx, uuid.New())

	// Assert
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestFlagRepository_Update(t *testing.T) {
	// Arrange
	db := testhelpers.NewTestPostgres(t)
	repo := repository.NewFlagRepository(db)
	ctx := context.Background()
	_, projectID := seedProject(t, db)

	flag := &domain.Flag{
		ProjectID: projectID, Key: "update-flag",
		Name: "Before Update", Type: domain.FlagTypeBoolean,
	}
	id, err := repo.Create(ctx, flag)
	require.NoError(t, err)

	// Act
	flag.ID = id
	flag.Name = "After Update"
	err = repo.Update(ctx, flag)
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, id)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "After Update", got.Name)
}

func TestFlagRepository_Delete(t *testing.T) {
	// Arrange
	db := testhelpers.NewTestPostgres(t)
	repo := repository.NewFlagRepository(db)
	ctx := context.Background()
	_, projectID := seedProject(t, db)

	flag := &domain.Flag{
		ProjectID: projectID, Key: "delete-flag",
		Name: "To Delete", Type: domain.FlagTypeBoolean,
	}
	id, err := repo.Create(ctx, flag)
	require.NoError(t, err)

	// Act
	err = repo.Delete(ctx, id)

	// Assert
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, id)
	assert.ErrorIs(t, err, domain.ErrNotFound)

}

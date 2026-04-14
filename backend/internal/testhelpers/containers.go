package testhelpers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	goredis "github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewTestPostgres(t testing.TB) *sqlx.DB {
	t.Helper()

	ctx := context.Background()

	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("flagr_test"),
		postgres.WithUsername("flagr"),
		postgres.WithPassword("flagr"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Errorf("failed to terminate postgres container: %v", err)
		}
	})

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	driver, err := migratepg.WithInstance(db.DB, &migratepg.Config{})
	if err != nil {
		t.Fatalf("failed to craete migrate driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../migrations",
		"postgres", driver,
	)
	if err != nil {
		t.Fatalf("failed to create migrator: %v", err)
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("failed to run migrations: %v", err)
	}

	return db
}

func NewTestRedis(tb testing.TB) *goredis.Client {
	tb.Helper()
	ctx := context.Background()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "redis:7-alpine",
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor:   wait.ForLog("Ready to accept connections"),
		},
		Started: true,
	})
	if err != nil {
		tb.Fatalf("failed to start redis container: %v", err)
	}

	tb.Cleanup(func() {
		container.Terminate(ctx)
	})

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "6379")

	client := goredis.NewClient(&goredis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port.Port()),
	})

	return client
}

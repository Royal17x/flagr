package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Royal17x/flagr/backend/internal/handler"
	"github.com/Royal17x/flagr/backend/internal/repository"
	"github.com/Royal17x/flagr/backend/internal/service"
	pg "github.com/Royal17x/flagr/backend/pkg/postgres"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// @title           Flagr API
// @version         1.0
// @description     Feature flags as a service — open-source alternative to LaunchDarkly

// @contact.name    Royal17x
// @contact.url     https://github.com/Royal17x/flagr

// @host            localhost:8080
// @BasePath        /api/v1

// @schemes         http https
func main() {
	// postgres
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = "postgres://flagr:flagr@localhost:5433/flagr?sslmode=disable"
	}
	db, err := pg.New(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}

	// repo's
	flagRepo := repository.NewFlagRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	envRepo := repository.NewEnvironmentRepository(db)
	flagEnvRepo := repository.NewFlagEnvironmentRepository(db)

	// services
	flagSvc := service.NewFlagService(flagRepo, projectRepo, flagEnvRepo)
	_ = service.NewProjectService(projectRepo)
	_ = service.NewEnvironmentService(envRepo, projectRepo)

	// handlers
	flagHandler := handler.NewFlagHandler(flagSvc)

	// router
	router := handler.NewRouter(flagHandler)

	// http serv
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	sigChan := make(chan os.Signal, 1)

	// server start on addr
	go func() {
		slog.Info("flagr listening", "addr", "8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

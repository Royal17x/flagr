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

	"github.com/Royal17x/flagr/backend/internal/app"
	cfg "github.com/Royal17x/flagr/backend/internal/config"
	pg "github.com/Royal17x/flagr/backend/pkg/postgres"
	redispkg "github.com/Royal17x/flagr/backend/pkg/redis"
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

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	config := cfg.Load()

	db, err := pg.New(config.Postgres.DSN)
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

	redisClient, err := redispkg.New(config.Redis.Addr, config.Redis.Password)
	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()

	application := app.New(config, db, redisClient)

	srv := &http.Server{
		Addr:         ":" + config.HTTP.Port,
		Handler:      application.Handler,
		ReadTimeout:  config.HTTP.ReadTimeout,
		WriteTimeout: config.HTTP.WriteTimeout,
	}

	go func() {
		slog.Info("flagr listening", "addr", config.HTTP.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

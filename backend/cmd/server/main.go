package main

import (
	"context"
	"errors"
	"github.com/Royal17x/flagr/backend/internal/app"
	cfg "github.com/Royal17x/flagr/backend/internal/config"
	grpcserver "github.com/Royal17x/flagr/backend/internal/grpc"
	"github.com/Royal17x/flagr/backend/pkg/kafka"
	pg "github.com/Royal17x/flagr/backend/pkg/postgres"
	redispkg "github.com/Royal17x/flagr/backend/pkg/redis"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	// config
	config, err := cfg.Load()
	if err != nil {
		log.Fatal(err)
	}
	// postgres
	db, err := pg.New(config.Postgres.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// postgres migrations
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

	// redis
	redisClient, err := redispkg.New(config.Redis.Addr, config.Redis.Password)
	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()

	// kafka topics
	if err := kafka.EnsureTopics(context.Background(), config.Kafka.Broker, config.Kafka.ReplicationFactor); err != nil {
		log.Fatal(err)
	}

	// kafka producer and consumer
	kafkaProducer := kafka.NewProducer(config.Kafka.Broker)
	defer kafkaProducer.Close()

	kafkaConsumer := kafka.NewConsumer(config.Kafka.Broker, "flagr-audit-consumer")
	defer kafkaConsumer.Close()

	// app.go
	application := app.New(config, db, redisClient, kafkaProducer)

	srv := &http.Server{
		Addr:         ":" + config.HTTP.Port,
		Handler:      application.Handler,
		ReadTimeout:  config.HTTP.ReadTimeout,
		WriteTimeout: config.HTTP.WriteTimeout,
	}

	//gRPC server
	go func() {
		slog.Info("grpc listening", "port", config.GRPC.Port)
		if err := grpcserver.StartGRPCServer(application.GRPCServer, config.GRPC.Port); err != nil {
			log.Fatal(err)
		}
	}()

	// http server
	go func() {
		slog.Info("flagr listening", "addr", config.HTTP.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	// kafka consume
	consumerCtx, consumerCancel := context.WithCancel(context.Background())
	defer consumerCancel()

	go kafkaConsumer.Consume(consumerCtx, func(ctx context.Context, msg kafka.AuditMessage) error {
		slog.Info("audit even received",
			"action", msg.Action,
			"resource_id", msg.ResourceID,
			"actor_id", msg.ActorID,
			"occurred_at", msg.OccurredAt,
		)
		return nil
	})

	// graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	application.GRPCServer.GracefulStop()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

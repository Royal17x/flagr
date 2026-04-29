package app

import (
	"net/http"

	"github.com/Royal17x/flagr/backend/pkg/kafka"
	"google.golang.org/grpc"

	"github.com/Royal17x/flagr/backend/internal/cache"
	"github.com/Royal17x/flagr/backend/internal/config"
	grpcserver "github.com/Royal17x/flagr/backend/internal/grpc"
	"github.com/Royal17x/flagr/backend/internal/handler"
	"github.com/Royal17x/flagr/backend/internal/middleware"
	"github.com/Royal17x/flagr/backend/internal/repository"
	"github.com/Royal17x/flagr/backend/internal/service"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type App struct {
	Handler    http.Handler
	GRPCServer *grpc.Server
}

func New(cfg *config.Config, db *sqlx.DB, redisClient *redis.Client, kafkaProducer *kafka.Producer) *App {
	// infrastructure
	flagCache := cache.NewFlagCache(redisClient)

	// repositories
	flagRepo := repository.NewFlagRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	orgRepo := repository.NewOrganizationRepository(db)
	envRepo := repository.NewEnvironmentRepository(db)
	flagEnvRepo := repository.NewFlagEnvironmentRepository(db)
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	sdkKeyRepo := repository.NewSDKKeyRepository(db)

	// services
	auditSvc := service.NewAuditService(kafkaProducer)
	flagSvc := service.NewFlagService(flagRepo, projectRepo, flagEnvRepo, envRepo, flagCache, auditSvc)
	authSvc := service.NewAuthService(
		userRepo,
		tokenRepo,
		orgRepo,
		projectRepo,
		envRepo,
		cfg.Auth.JWTSecret,
		cfg.Auth.AccessTokenDuration,
		cfg.Auth.RefreshTokenDuration,
	)
	projectSvc := service.NewProjectService(projectRepo)
	envSvc := service.NewEnvironmentService(envRepo, projectRepo)

	// handlers & middleware
	flagHandler := handler.NewFlagHandler(flagSvc)
	authHandler := handler.NewAuthHandler(authSvc)
	healthHandler := handler.NewHealthHandler(db, redisClient)
	authMiddleware := middleware.NewAuthMiddleware(authSvc)
	sdkAuthMiddleware := middleware.NewSDKAuthMiddleware(sdkKeyRepo)
	projectHandler := handler.NewProjectHandler(projectSvc)
	envHandler := handler.NewEnvironmentHandler(envSvc)
	sdkKeyHandler := handler.NewSDKKeyHandler(sdkKeyRepo)

	// router
	router := handler.NewRouter(flagHandler,
		authHandler,
		healthHandler,
		authMiddleware,
		sdkAuthMiddleware,
		projectHandler,
		envHandler,
		sdkKeyHandler)

	//gRPC
	grpcSrv, _ := grpcserver.NewGRPCServer(flagSvc, authSvc)

	return &App{
		Handler:    router,
		GRPCServer: grpcSrv,
	}
}

package app

import (
	"net/http"

	"github.com/Royal17x/flagr/backend/internal/cache"
	"github.com/Royal17x/flagr/backend/internal/config"
	"github.com/Royal17x/flagr/backend/internal/handler"
	"github.com/Royal17x/flagr/backend/internal/middleware"
	"github.com/Royal17x/flagr/backend/internal/repository"
	"github.com/Royal17x/flagr/backend/internal/service"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

// App holds the wired application graph and exposes the HTTP handler.
type App struct {
	Handler http.Handler
}

// New wires all layers (repo → service → handler) and returns a ready App.
func New(cfg *config.Config, db *sqlx.DB, redisClient *redis.Client) *App {
	// infrastructure
	flagCache := cache.NewFlagCache(redisClient)

	// repositories
	flagRepo := repository.NewFlagRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	envRepo := repository.NewEnvironmentRepository(db)
	flagEnvRepo := repository.NewFlagEnvironmentRepository(db)
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)

	// services
	flagSvc := service.NewFlagService(flagRepo, projectRepo, flagEnvRepo, flagCache)
	authSvc := service.NewAuthService(
		userRepo,
		tokenRepo,
		cfg.Auth.JWTSecret,
		cfg.Auth.AccessTokenDuration,
		cfg.Auth.RefreshTokenDuration,
	)
	_ = service.NewProjectService(projectRepo)
	_ = service.NewEnvironmentService(envRepo, projectRepo)

	// handlers & middleware
	flagHandler := handler.NewFlagHandler(flagSvc)
	authHandler := handler.NewAuthHandler(authSvc)
	authMiddleware := middleware.NewAuthMiddleware(authSvc)

	// router
	router := handler.NewRouter(flagHandler, authHandler, authMiddleware)

	return &App{Handler: router}
}

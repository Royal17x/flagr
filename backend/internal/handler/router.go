package handler

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"

	_ "github.com/Royal17x/flagr/backend/docs"
	"github.com/Royal17x/flagr/backend/internal/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(flagHandler *FlagHandler, authHandler *AuthHandler, healthHandler *HealthHandler, authMiddleware *middleware.AuthMiddleware, sdkAuthMiddleware *middleware.SDKAuthMiddleware) http.Handler {
	r := chi.NewRouter()

	rl := middleware.NewRateLimiter(60, 120)
	corsConfig := middleware.DefaultCORSConfig()
	authRL := middleware.NewAuthRateLimiter()

	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.CORS(corsConfig))
	r.Use(chiMiddleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(rl.Middleware)
	r.Use(middleware.PrometheusMiddleware)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
	r.Get("/health/live", healthHandler.Live)
	r.Get("/health/ready", healthHandler.Ready)
	r.Handle("/metrics", promhttp.Handler())

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Use(authRL.Middleware)
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.Post("/refresh", authHandler.Refresh)
			r.Post("/logout", authHandler.Logout)

		})
		r.Group(func(r chi.Router) {
			r.Use(sdkAuthMiddleware.Authenticate)
			r.Get("/flags/evaluate", flagHandler.Evaluate)
		})
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)
			r.Route("/flags", func(r chi.Router) {
				r.Post("/", flagHandler.Create)
				r.Get("/", flagHandler.List)
				r.Get("/{id}", flagHandler.GetByID)
				r.Put("/{id}", flagHandler.Update)
				r.Delete("/{id}", flagHandler.Delete)
			})
		})
	})

	return otelhttp.NewHandler(r, "flagr-http")
}

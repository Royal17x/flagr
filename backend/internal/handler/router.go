package handler

import (
	"net/http"

	_ "github.com/Royal17x/flagr/backend/docs"
	"github.com/Royal17x/flagr/backend/internal/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(flagHandler *FlagHandler, authHandler *AuthHandler, authMiddleware *middleware.AuthMiddleware) http.Handler {
	r := chi.NewRouter()

	rl := middleware.NewRateLimiter(60, 120)

	r.Use(chiMiddleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(rl.Middleware)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.Post("/refresh", authHandler.Refresh)
			r.Post("/logout", authHandler.Logout)
		})
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)
			r.Route("/flags", func(r chi.Router) {
				r.Post("/", flagHandler.Create)
				r.Get("/", flagHandler.List)
				r.Get("/evaluate", flagHandler.Evaluate)
				r.Get("/{id}", flagHandler.GetByID)
				r.Put("/{id}", flagHandler.Update)
				r.Delete("/{id}", flagHandler.Delete)
			})
		})
	})

	return r
}

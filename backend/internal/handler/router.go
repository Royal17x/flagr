package handler

import (
	"net/http"

	_ "github.com/Royal17x/flagr/backend/docs"
	"github.com/Royal17x/flagr/backend/internal/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(flagHandler *FlagHandler) http.Handler {
	r := chi.NewRouter()

	rl := middleware.NewRateLimiter(60, 120)

	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RequestID)
	r.Use(rl.Middleware)
	r.Use(middleware.Logger)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/flags", func(r chi.Router) {
			r.Post("/", flagHandler.Create)
			r.Get("/", flagHandler.List)
			r.Get("/evaluate", flagHandler.Evaluate)
			r.Get("/{id}", flagHandler.GetByID)
			r.Put("/{id}", flagHandler.Update)
			r.Delete("/{id}", flagHandler.Delete)
		})
	})

	return r
}

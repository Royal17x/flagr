package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func NewRouter(flagHandler *FlagHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

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

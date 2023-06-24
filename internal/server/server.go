package server

import (
	"net/http"

	"github.com/bobgromozeka/metrics/internal/server/handlers"
	"github.com/bobgromozeka/metrics/internal/server/middlewares"
	"github.com/bobgromozeka/metrics/internal/server/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func new(s storage.Storage) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.StripSlashes)
	r.Use(middlewares.WithLogging)
	r.Post("/update/{type}/{name}/{value}", handlers.Update(s))
	r.Get("/value/{type}/{name}", handlers.Get(s))
	r.Post("/update", handlers.UpdateJSON(s))
	r.Post("/value", handlers.GetJSON(s))
	r.Get("/", handlers.GetAll(s))

	return r
}

func Start(serverAddr string) error {

	s := storage.New()
	server := new(s)

	return http.ListenAndServe(serverAddr, server)
}

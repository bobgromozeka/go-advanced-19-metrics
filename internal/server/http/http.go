package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/bobgromozeka/metrics/internal/server/http/handlers"
	"github.com/bobgromozeka/metrics/internal/server/middlewares"
	"github.com/bobgromozeka/metrics/internal/server/storage"
)

type Config struct {
	Addr          string
	PrivateKey    []byte
	TrustedSubnet string
	HashKey       string
}

func New(s storage.Storage, c Config, privateKey []byte) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.StripSlashes)

	r.Group(
		func(r chi.Router) {
			r.Use(
				middlewares.WithLogging(
					[]string{
						"./http.log",
					},
				),
				middlewares.TrustedSubnet(c.TrustedSubnet),
				middlewares.Gzippify,
				middlewares.Rsa(privateKey),
			)
			r.Post("/update/{type}/{name}/{value}", handlers.Update(s))
			r.Get("/value/{type}/{name}", handlers.Get(s))
			r.Post("/update", handlers.UpdateJSON(s))
			r.Post("/updates", handlers.Updates(s, c.HashKey))
			r.Post("/value", handlers.GetJSON(s))
			r.Get("/", handlers.GetAll(s))
		},
	)
	r.Get("/ping", handlers.Ping)

	return r
}

func Start(ctx context.Context, c Config, s storage.Storage) error {
	router := New(s, c, c.PrivateKey)
	srv := &http.Server{Addr: c.Addr, Handler: router}

	go func() {
		<-ctx.Done()

		hardCtx, hardCancel := context.WithTimeout(context.Background(), time.Second*15)
		defer hardCancel()

		srv.Shutdown(hardCtx)
	}()

	fmt.Printf("Starting http server on addr [%s].......\n", c.Addr)

	return srv.ListenAndServe()
}

package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/bobgromozeka/metrics/internal/server/db"
	"github.com/bobgromozeka/metrics/internal/server/handlers"
	"github.com/bobgromozeka/metrics/internal/server/middlewares"
	"github.com/bobgromozeka/metrics/internal/server/storage"
)

func New(s storage.Storage, config StartupConfig, privateKey []byte) *chi.Mux {
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
				middlewares.TrustedSubnet(config.TrustedSubnet),
				middlewares.Gzippify,
				middlewares.Rsa(privateKey),
			)
			r.Post("/update/{type}/{name}/{value}", handlers.Update(s))
			r.Get("/value/{type}/{name}", handlers.Get(s))
			r.Post("/update", handlers.UpdateJSON(s))
			r.Post("/updates", handlers.Updates(s, config.HashKey))
			r.Post("/value", handlers.GetJSON(s))
			r.Get("/", handlers.GetAll(s))
		},
	)
	r.Get("/ping", handlers.Ping)

	return r
}

func Start(ctx context.Context, startupConfig StartupConfig) error {
	var wg sync.WaitGroup

	var s storage.Storage

	if startupConfig.DatabaseDsn != "" {
		connErr := db.Connect(startupConfig.DatabaseDsn)
		if connErr != nil {
			panic(connErr)
		}

		ddlErr := storage.Bootstrap(db.Connection())
		if ddlErr != nil {
			panic(ddlErr)
		}
		s = storage.NewPG(db.Connection())

		wg.Add(1)

		go func() {
			defer wg.Done()
			<-ctx.Done()

			hardCtx, hardCancel := context.WithTimeout(context.Background(), time.Second*15)
			defer hardCancel()

			db.Connection().Close(hardCtx)
		}()

	} else {
		s = storage.NewMemory()
		s = storage.NewPersistenceStorage(
			s, storage.PersistenceSettings{
				Path:     startupConfig.FileStoragePath,
				Interval: startupConfig.StoreInterval,
				Restore:  startupConfig.Restore,
			},
		)
	}

	privateKey, readErr := os.ReadFile(startupConfig.PrivateKeyPath)
	if readErr != nil {
		return readErr
	}

	router := New(s, startupConfig, privateKey)
	server := &http.Server{Addr: startupConfig.ServerAddr, Handler: router}

	wg.Add(1)

	go func() {
		defer wg.Done()
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln(err)
		}
	}()

	go func() {
		<-ctx.Done()

		hardCtx, hardCancel := context.WithTimeout(context.Background(), time.Second*15)
		defer hardCancel()

		server.Shutdown(hardCtx)
	}()

	wg.Wait()

	return nil
}

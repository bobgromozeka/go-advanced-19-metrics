package server

import (
	"log"
	"net/http"
	"time"

	"github.com/bobgromozeka/metrics/internal/server/handlers"
	"github.com/bobgromozeka/metrics/internal/server/middlewares"
	"github.com/bobgromozeka/metrics/internal/server/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type FileStorageSettings struct {
	Path     string
	Interval uint
	Restore  bool
}

func new(s storage.Storage) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.StripSlashes)

	r.Group(func(r chi.Router) {
		r.Use(
			middlewares.WithLogging,
			middlewares.Gzippify,
		)
		r.Post("/update/{type}/{name}/{value}", handlers.Update(s))
		r.Get("/value/{type}/{name}", handlers.Get(s))
		r.Post("/update", handlers.UpdateJSON(s))
		r.Post("/value", handlers.GetJSON(s))
		r.Get("/", handlers.GetAll(s))
	})

	return r
}

func Start(serverAddr string, fss FileStorageSettings) error {
	s := storage.New()

	if fss.Path != "" {
		if fss.Restore {
			if restoreErr := s.RestoreFrom(fss.Path); restoreErr != nil {
				log.Println("Could not restore data from file: ", restoreErr)
			}
		}

		err := setPersistence(s, fss.Path, fss.Interval)
		if err != nil {
			return err
		}
	}

	server := new(s)

	return http.ListenAndServe(serverAddr, server)
}

func setPersistence(s storage.MemStorage, filePath string, interval uint) error {

	saveFunc := func() {
		if persistenceErr := s.PersistToPath(filePath); persistenceErr != nil {
			log.Println("Error during writing metrics to file: ", persistenceErr)
		}
	}

	if interval == 0 {
		s.Listen(storage.Update, saveFunc)
	} else {
		go func() {
			ticker := time.Tick(time.Second * time.Duration(interval))
			for range ticker {
				saveFunc()
			}
		}()
	}

	return nil
}

package server

import (
	"github.com/bobgromozeka/metrics/internal/server/handlers"
	"github.com/bobgromozeka/metrics/internal/server/storage"
	"net/http"
)

func new(s storage.Storage) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/update/", handlers.UpdateHandler(s))

	return mux
}

func Start() error {
	storage := storage.New()
	server := new(storage)

	return http.ListenAndServe(":8080", server)
}

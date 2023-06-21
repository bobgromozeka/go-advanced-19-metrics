package handlers

import (
	"fmt"
	"net/http"

	"github.com/bobgromozeka/metrics/internal/server/storage"

	"github.com/go-chi/chi/v5"
)

func Get(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricsType := chi.URLParam(r, "type")
		metricsName := chi.URLParam(r, "name")

		m, ok := s.GetMetrics(metricsType, metricsName)

		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Write([]byte(fmt.Sprintf("%v", m)))
		w.WriteHeader(http.StatusOK)
	}
}

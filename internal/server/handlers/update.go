package handlers

import (
	"github.com/bobgromozeka/metrics/internal/metrics"
	"github.com/bobgromozeka/metrics/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func Update(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricsType := chi.URLParam(r, "type")
		metricsName := chi.URLParam(r, "name")
		metricsValue := chi.URLParam(r, "value")

		if !metrics.IsValidType(metricsType) {
			log.Println("Got wrong metrics type in request: ", metricsType)
			http.Error(w, "Wrong metrics type", http.StatusBadRequest)
			return
		}

		if !metrics.IsValidValue(metricsType, metricsValue) {
			log.Println("Got wrong metrics value in request: ", metricsValue)
			http.Error(w, "Wrong metrics value", http.StatusBadRequest)
			return
		}

		s.UpdateMetricsType(metricsType, metricsName, metricsValue)

		w.WriteHeader(http.StatusOK)
	}
}

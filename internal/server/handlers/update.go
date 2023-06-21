package handlers

import (
	"log"
	"net/http"

	"github.com/bobgromozeka/metrics/internal/metrics"
	"github.com/bobgromozeka/metrics/internal/server/storage"

	"github.com/go-chi/chi/v5"
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

		_, err := s.UpdateMetricsType(metricsType, metricsName, metricsValue)

		if err != nil {
			log.Printf("Could not update metrics: [type: %s, name: %s, value: %s]: %s ", metricsType, metricsName, metricsValue, err)
		}

		w.WriteHeader(http.StatusOK)
	}
}

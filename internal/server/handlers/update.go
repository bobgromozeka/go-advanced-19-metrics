package handlers

import (
	"encoding/json"
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

		_, err := s.UpdateMetricsByType(metricsType, metricsName, metricsValue)

		if err != nil {
			log.Printf("Could not update metrics: [type: %s, name: %s, value: %s]: %s ", metricsType, metricsName, metricsValue, err)
		}

		w.WriteHeader(http.StatusOK)
	}
}

func UpdateJSON(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestMetrics metrics.RequestPayload

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&requestMetrics); err != nil {
			http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
			return
		}

		if !metrics.IsValidType(requestMetrics.MType) {
			log.Println("Got wrong metrics type in request: ", requestMetrics.MType)
			http.Error(w, "Wrong metrics type", http.StatusBadRequest)
			return
		}

		if requestMetrics.MType == metrics.CounterType {
			var delta int64
			if requestMetrics.Delta == nil {
				delta = 0
			} else {
				delta = *requestMetrics.Delta
			}
			newValue := s.AddCounter(requestMetrics.ID, delta)
			requestMetrics.Delta = &newValue
		} else {
			var value float64
			if requestMetrics.Value == nil {
				value = 0
			} else {
				value = *requestMetrics.Value
			}
			newValue := s.SetGauge(requestMetrics.ID, value)
			requestMetrics.Value = &newValue
		}

		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		if encodingErr := encoder.Encode(requestMetrics); encodingErr != nil {
			log.Println("Error during encoding update request: ", encodingErr)
			http.Error(w, encodingErr.Error(), http.StatusInternalServerError)
		}
	}
}

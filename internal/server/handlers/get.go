package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bobgromozeka/metrics/internal/metrics"
	"github.com/bobgromozeka/metrics/internal/server/storage"

	"github.com/go-chi/chi/v5"
)

func Get(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricsType := chi.URLParam(r, "type")
		metricsName := chi.URLParam(r, "name")

		m, ok := s.GetMetricsByType(metricsType, metricsName)

		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Write([]byte(fmt.Sprintf("%v", m)))
		w.WriteHeader(http.StatusOK)
	}
}

func GetJSON(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestMetrics metrics.RequestPayload

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&requestMetrics); err != nil {
			http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		}

		if !metrics.IsValidType(requestMetrics.MType) {
			log.Println("Got wrong metrics type in request: ", requestMetrics.MType)
			http.Error(w, "Wrong metrics type", http.StatusBadRequest)
			return
		}

		if requestMetrics.MType == metrics.CounterType {
			val, ok := s.GetCounterMetrics(requestMetrics.ID)
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			requestMetrics.Delta = &val
		} else {
			val, ok := s.GetGaugeMetrics(requestMetrics.ID)
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			requestMetrics.Value = &val
		}

		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		if encodingErr := encoder.Encode(requestMetrics); encodingErr != nil {
			log.Println("Error during encoding update request: ", encodingErr)
			http.Error(w, encodingErr.Error(), http.StatusInternalServerError)
		}
	}
}

package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/bobgromozeka/metrics/internal"
	"github.com/bobgromozeka/metrics/internal/hash"
	"github.com/bobgromozeka/metrics/internal/helpers"
	"github.com/bobgromozeka/metrics/internal/metrics"
	"github.com/bobgromozeka/metrics/internal/server/storage"
)

// Updates Batch metrics update.
func Updates(s storage.Storage, hashKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestMetrics []metrics.RequestPayload

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
			return
		}

		if sum := r.Header.Get(internal.HTTPCheckSumHeader); sum != "" && !hash.IsValidSum(sum, string(body), hashKey) {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if jsonErr := json.Unmarshal(body, &requestMetrics); jsonErr != nil {
			http.Error(w, "Bad request: "+jsonErr.Error(), http.StatusBadRequest)
			return
		}

		metricsMap := metricsArrToMaps(requestMetrics)
		cErr := s.AddCounters(r.Context(), metricsMap.Counter)
		gErr := s.SetGauges(r.Context(), metricsMap.Gauge)

		if gErr != nil || cErr != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		responseBody := []byte(`{"success":true}`)
		helpers.SignResponse(w, responseBody, hashKey, internal.HTTPCheckSumHeader)

		w.WriteHeader(http.StatusOK)
		w.Write(responseBody)
	}
}

func metricsArrToMaps(arr []metrics.RequestPayload) storage.Metrics {
	m := storage.Metrics{
		Gauge:   storage.GaugeMetrics{},
		Counter: storage.CounterMetrics{},
	}

	for _, payload := range arr {
		if !metrics.IsValidType(payload.MType) {
			continue
		}

		if payload.MType == metrics.CounterType {
			var delta int64
			if payload.Delta == nil {
				delta = 0
			} else {
				delta = *payload.Delta
			}

			m.Counter[payload.ID] += delta
		} else if payload.MType == metrics.GaugeType {
			var value float64
			if payload.Value == nil {
				value = 0
			} else {
				value = *payload.Value
			}
			m.Gauge[payload.ID] = value
		}
	}

	return m
}

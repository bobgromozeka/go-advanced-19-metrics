package handlers

import (
	"fmt"
	"net/http"

	"github.com/bobgromozeka/metrics/internal/server/storage"
)

func GetAll(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gaugeMetrics := s.GetAllGaugeMetrics()
		counterMetrics := s.GetAllCounterMetrics()

		response := ""

		for k, v := range gaugeMetrics {
			response += fmt.Sprintf("%s:   %f\r\n", k, v)
		}

		for k, v := range counterMetrics {
			response += fmt.Sprintf("%s:   %d\n", k, v)
		}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}
}

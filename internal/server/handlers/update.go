package handlers

import (
	"github.com/bobgromozeka/metrics/internal/metrics"
	"github.com/bobgromozeka/metrics/internal/server/storage"
	"net/http"
	"strings"
)

const (
	UpdatePartsCount = 4
)

const (
	TypePart = iota + 1
	NamePart
	ValuePart
)

func UpdateHandler(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		urlPath := strings.Trim(r.URL.Path, "/")
		urlParts := strings.Split(urlPath, "/")

		if len(urlParts) < 3 {
			http.Error(w, "Metrics name is not specified", http.StatusNotFound)
			return
		}

		if len(urlParts) != UpdatePartsCount {
			http.Error(w, "Wrong update endpoint signature. Should be /update/type/name/value", http.StatusUnprocessableEntity)
			return
		}

		if !metrics.IsValidType(urlParts[TypePart]) {
			http.Error(w, "Wrong metrics type", http.StatusBadRequest)
			return
		}

		if !metrics.IsValidValue(urlParts[TypePart], urlParts[ValuePart]) {
			http.Error(w, "Wrong metrics value", http.StatusBadRequest)
			return
		}

		s.UpdateMetricsType(urlParts[TypePart], urlParts[NamePart], urlParts[ValuePart])

		w.WriteHeader(http.StatusOK)
	}
}

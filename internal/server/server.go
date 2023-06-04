package server

import (
	"github.com/bobgromozeka/metrics/internal/server/mertics"
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

func new(s *mertics.MemStorage) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) {
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

		if !mertics.IsValidType(urlParts[TypePart]) {
			http.Error(w, "Wrong metrics type", http.StatusBadRequest)
			return
		}

		if !mertics.IsValidValue(urlParts[TypePart], urlParts[ValuePart]) {
			http.Error(w, "Wrong metrics value", http.StatusBadRequest)
			return
		}

		s.UpdateMetricsType(urlParts[TypePart], urlParts[NamePart], urlParts[ValuePart])

		w.WriteHeader(http.StatusOK)
	})

	return mux
}

func Start() error {
	storage := mertics.New()
	server := new(&storage)

	return http.ListenAndServe(":8080", server)
}

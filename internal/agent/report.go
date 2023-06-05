package agent

import (
	"fmt"
	"github.com/bobgromozeka/metrics/internal/metrics"
	"net/http"
	"reflect"
)

func reportToServer(rm runtimeMetrics) {
	serverHost := "http://localhost:8080"

	signatures := makeEndpointsFromStructure(rm)

	for _, signature := range signatures {
		resp, _ := http.Post(serverHost+signature, "text/plain", nil)
		_ = resp.Body.Close()
	}
}

func makeEndpointsFromStructure(rm any) []string {
	v := reflect.ValueOf(rm)
	t := reflect.TypeOf(rm)

	var signatures []string

	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			var metricsValue string
			fieldV := v.Field(i)
			fieldT := t.Field(i)
			switch fieldV.Kind() {
			case reflect.Uint64, reflect.Uint32:
				metricsValue = fmt.Sprintf("%d", fieldV.Interface())
			case reflect.Float64:
				metricsValue = fmt.Sprintf("%f", fieldV.Interface())
			}

			if len(metricsValue) < 1 {
				continue
			}

			metricsType := metrics.Gauge
			if mt, ok := runtimeMetricsTypes[fieldT.Name]; ok {
				metricsType = mt
			}

			signatures = append(signatures, fmt.Sprintf("/update/%s/%s/%s/", metricsType, fieldT.Name, metricsValue))
		}
	}

	return signatures
}

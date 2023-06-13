package agent

import (
	"fmt"
	"github.com/bobgromozeka/metrics/internal/metrics"
	"github.com/go-resty/resty/v2"
	"reflect"
)

func reportToServer(serverAddr string, rm runtimeMetrics) {

	signatures := makeEndpointsFromStructure(rm)

	client := resty.New()
	for _, signature := range signatures {
		_, _ = client.R().Post(serverAddr + signature)
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

			metricsType := metrics.GaugeType
			if mt, ok := runtimeMetricsTypes[fieldT.Name]; ok {
				metricsType = mt
			}

			signatures = append(signatures, fmt.Sprintf("/update/%s/%s/%s", metricsType, fieldT.Name, metricsValue))
		}
	}

	return signatures
}

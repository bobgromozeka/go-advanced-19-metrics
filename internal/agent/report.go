package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/bobgromozeka/metrics/internal"
	"github.com/bobgromozeka/metrics/internal/hash"
	"github.com/bobgromozeka/metrics/internal/helpers"
	"github.com/bobgromozeka/metrics/internal/metrics"
)

func reportToServer(serverAddr string, hashKey string, rm runtimeMetrics) {

	payloads := makeBodiesFromStructure(rm)

	if len(payloads) < 1 {
		return
	}

	//resty client has jitter func to calc wait time between attempts by default (1 + 2^attempt sec)
	client := resty.
		New().
		SetRetryCount(3).
		SetRetryWaitTime(time.Second * 1)

	encodedPayload, err := json.Marshal(payloads)
	if err != nil {
		log.Println("Could not encode request: ", err)
		return
	}

	signature := hash.Sign(hashKey, encodedPayload)
	if signature != "" {
		client.SetHeader(internal.HTTPCheckSumHeader, signature)
	}

	gzippedPayload, gzErr := helpers.Gzip(encodedPayload)
	if gzErr != nil {
		log.Println("Could not gzip request: ", gzErr)
		return
	}

	_, _ = client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(gzippedPayload).
		Post(serverAddr + "/updates")
}

func makeBodiesFromStructure(rm any) []metrics.RequestPayload {
	v := reflect.ValueOf(rm)
	t := reflect.TypeOf(rm)

	var payloads []metrics.RequestPayload

	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			fieldV := v.Field(i)
			fieldT := t.Field(i)
			if fieldV.Kind() == reflect.Slice {
				for j := 0; j < fieldV.Len(); j++ {
					sliceElV := fieldV.Index(j)
					if payload := makeBodyFromStructField(sliceElV, fieldT.Name+strconv.Itoa(j)); payload != nil {
						payloads = append(payloads, *payload)
					}
				}
			} else {
				if payload := makeBodyFromStructField(fieldV, fieldT.Name); payload != nil {
					payloads = append(payloads, *payload)
				}
			}
		}
	}

	return payloads
}

func makeBodyFromStructField(v reflect.Value, name string) *metrics.RequestPayload {
	metricsType := metrics.GaugeType
	if mt, ok := runtimeMetricsTypes[name]; ok {
		metricsType = mt
	}

	rp := metrics.RequestPayload{
		ID:    name,
		MType: metricsType,
	}

	//Shit conversions, but we lose accuracy anyway converting uint64 to float64
	switch metricsType {
	case metrics.GaugeType:
		switch val := v.Interface().(type) {
		case float64:
			rp.Value = &val
		case uint64, uint32:
			strVal := fmt.Sprintf("%d", v.Interface())
			intVal := helpers.StrToInt(strVal)
			fVal := float64(intVal)
			rp.Value = &fVal
		}
	case metrics.CounterType:
		strVal := fmt.Sprintf("%d", v.Interface())
		intVal := helpers.StrToInt(strVal)
		val := int64(intVal)
		rp.Delta = &val
	}

	if rp.Value == nil && rp.Delta == nil {
		return nil
	}

	return &rp
}

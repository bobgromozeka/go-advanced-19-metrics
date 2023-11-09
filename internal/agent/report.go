package agent

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/bobgromozeka/metrics/internal"
	"github.com/bobgromozeka/metrics/internal/hash"
	"github.com/bobgromozeka/metrics/internal/helpers"
	"github.com/bobgromozeka/metrics/internal/metrics"
	proto_interfaces "github.com/bobgromozeka/metrics/internal/proto-interfaces"
	"github.com/bobgromozeka/metrics/internal/utils"
)

func reportToHTTPServer(serverAddr string, hashKey string, publicKey []byte, rm runtimeMetrics) {

	payloads := makeBodiesFromStructure(rm)

	if len(payloads) < 1 {
		return
	}

	//resty client has jitter func to calc wait time between attempts by default (1 + 2^attempt sec)
	client := resty.
		New().
		SetRetryCount(3).
		SetRetryWaitTime(time.Second * 1)
	req := client.R()

	payload, err := json.Marshal(payloads)
	if err != nil {
		log.Println("Could not encode request: ", err)
		return
	}

	signature := hash.Sign(hashKey, payload)
	if signature != "" {
		req.SetHeader(internal.HTTPCheckSumHeader, signature)
	}

	payload, encryptErr := encryptData(payload, publicKey)
	if encryptErr != nil {
		fmt.Printf("Could not encrypt data: %v", encryptErr)
	} else {
		req.SetHeader(internal.RSAEncryptedHeader, "true")
	}

	gzippedPayload, gzErr := helpers.Gzip(payload)
	if gzErr != nil {
		log.Println("Could not gzip request: ", gzErr)
		return
	}

	_, _ = req.
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetHeader(internal.RealIPHeader, utils.GetLocalIPv4().String()).
		SetBody(gzippedPayload).
		Post(serverAddr + "/updates")
}

func reportToGRPCServer(serverAddr string, certPath string, rm runtimeMetrics) {
	payloads := makeBodiesFromStructure(rm)
	bur := metricPayloadsToGRPCPayload(payloads)

	creds, err := credentials.NewClientTLSFromFile(
		certPath, "x.test.example.com",
	) // test certs are valid for this domain
	if err != nil {
		log.Println("Could not create creds for client", err)
	}

	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Println("Error during creating grpc connection", err)
	}

	client := proto_interfaces.NewMetricsClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	_, err = client.BatchUpdate(ctx, bur)
	if err != nil {
		log.Println("Error during sending metrics to grpc server", err)
	}
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

func metricPayloadsToGRPCPayload(rps []metrics.RequestPayload) *proto_interfaces.BatchUpdateRequest {
	bur := &proto_interfaces.BatchUpdateRequest{}

	for _, payload := range rps {
		bur.Data = append(
			bur.Data, &proto_interfaces.Entry{
				Delta: payload.Delta,
				Value: payload.Value,
				ID:    payload.ID,
				MType: payload.MType,
			},
		)
	}

	return bur
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

func encryptData(data []byte, key []byte) ([]byte, error) {
	if len(key) > 0 {
		publicKeyBlock, _ := pem.Decode(key)
		parsedPublicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
		if err != nil {
			return data, err
		}

		res := make([]byte, 0)

		h := sha256.New()
		step := parsedPublicKey.(*rsa.PublicKey).Size() - 2*h.Size() - 2

		for i := 0; i < len(data); i += step {
			end := i + step
			if end > len(data) {
				end = len(data)
			}

			enc, err := rsa.EncryptOAEP(h, rand.Reader, parsedPublicKey.(*rsa.PublicKey), data[i:end], []byte("data"))
			if err != nil {
				return data, err
			}

			res = append(res, enc...)
		}

		return res, nil
	}
	return data, nil
}

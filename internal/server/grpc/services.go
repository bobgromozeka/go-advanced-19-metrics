package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bobgromozeka/metrics/internal/metrics"
	proto_interfaces "github.com/bobgromozeka/metrics/internal/proto-interfaces"
	"github.com/bobgromozeka/metrics/internal/server/storage"
)

type MetricsService struct {
	proto_interfaces.UnimplementedMetricsServer
	stor storage.Storage
}

func NewMetricsService(stor storage.Storage) *MetricsService {
	return &MetricsService{
		stor: stor,
	}
}

func (s *MetricsService) BatchUpdate(ctx context.Context, request *proto_interfaces.BatchUpdateRequest) (*proto_interfaces.Empty, error) {
	metricsMap := entriesToMap(request.Data)
	cErr := s.stor.AddCounters(ctx, metricsMap.Counter)
	gErr := s.stor.SetGauges(ctx, metricsMap.Gauge)

	if cErr != nil || gErr != nil {
		return nil, status.Errorf(codes.Internal, "Could not save metrics")
	}

	return &proto_interfaces.Empty{}, nil
}

func entriesToMap(arr []*proto_interfaces.Entry) storage.Metrics {
	m := storage.Metrics{
		Gauge:   storage.GaugeMetrics{},
		Counter: storage.CounterMetrics{},
	}

	if arr == nil {
		return m
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

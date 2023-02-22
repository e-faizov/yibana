package server

import (
	"context"

	"github.com/e-faizov/yibana/internal"
	"github.com/e-faizov/yibana/internal/interfaces"
	"github.com/e-faizov/yibana/proto"
	"github.com/rs/zerolog/log"
)

type ProtoMetrics struct {
	proto.UnimplementedMetricsServiceServer
	// Store - интерфейс хранилища
	Store interfaces.Store
	// Key - ключ для подсчета и проверки хэша
	Key string
}

func (p *ProtoMetrics) AddMetric(ctx context.Context, data *proto.Metric) (*proto.AddMetricResponse, error) {
	var resp proto.AddMetricResponse

	metric := internal.Convert(data)

	if len(p.Key) != 0 {
		if !internal.CheckHash(p.Key, metric) {
			log.Error().
				Str("hash", metric.Hash).
				Str("id", metric.ID).
				Msg("ProtoMetrics.AddMetric error wrong hash")
			resp.Error = "data error"
			return &resp, nil
		}
	}

	err := p.Store.SetMetrics(ctx, []internal.Metric{metric})
	if err != nil {
		resp.Error = err.Error()
	}

	return &resp, nil
}
func (p *ProtoMetrics) AddMetrics(ctx context.Context, data *proto.Metrics) (*proto.AddMetricResponse, error) {
	var resp proto.AddMetricResponse
	converted := make([]internal.Metric, 0, len(data.Data))

	for _, d := range data.Data {
		converted = append(converted, internal.Convert(d))
	}

	if len(p.Key) != 0 {
		for _, metric := range converted {
			if !internal.CheckHash(p.Key, metric) {
				log.Error().
					Str("hash", metric.Hash).
					Str("id", metric.ID).
					Msg("ProtoMetrics.AddMetrics error wrong hash")
				resp.Error = "data error"
				return &resp, nil
			}
		}
	}

	err := p.Store.SetMetrics(ctx, converted)
	if err != nil {
		resp.Error = err.Error()
	}

	return &resp, nil
}
func (p *ProtoMetrics) GetMetric(ctx context.Context, data *proto.Metric) (*proto.GetMetricResponse, error) {
	var resp proto.GetMetricResponse
	retData := proto.Metric{
		Id:    data.Id,
		MType: data.MType,
	}

	convertedData := internal.Convert(data)

	switch convertedData.MType {
	case internal.GaugeType, internal.CounterType:
		res, ok, err := p.Store.GetMetric(ctx, convertedData)
		if err != nil {
			resp.Error = err.Error()
		} else if !ok {
			resp.Error = "not found"
		} else {
			if len(p.Key) != 0 {
				if res.MType == internal.GaugeType {
					if res.Value != nil {
						tmp := float64(*res.Value)
						retData.Value = &tmp
						retData.Hash = internal.CalcGaugeHash(res.ID, *res.Value, p.Key)

					}
				} else {
					if res.Delta != nil {
						tmp := int64(*res.Delta)
						retData.Delta = &tmp
						retData.Hash = internal.CalcCounterHash(res.ID, *res.Delta, p.Key)

					}
				}
			}

			resp.Data = &retData
		}
	default:
		resp.Error = "invalid type"
	}

	return &resp, nil
}

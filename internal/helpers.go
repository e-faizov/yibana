package internal

import (
	"fmt"
	"runtime"

	"github.com/e-faizov/yibana/proto"
)

// ErrorHelper - функция добавления имени файла и строки к ошибке
func ErrorHelper(err error) error {
	if err == nil {
		return err
	}
	_, filename, line, _ := runtime.Caller(1)
	return fmt.Errorf("[%s:%d] %w", filename, line, err)
}

func Convert(data *proto.Metric) Metric {
	var metric Metric

	metric.ID = data.Id
	metric.Hash = data.Hash
	if proto.Metric_GAUGE == data.MType {
		metric.MType = GaugeType
		val := Gauge(*data.Value)
		metric.Value = &val
	} else {
		metric.MType = CounterType
		val := Counter(*data.Delta)
		metric.Delta = &val
	}
	return metric
}

func UnConvert(data Metric) *proto.Metric {
	var ret proto.Metric
	ret.Id = data.ID
	ret.Hash = data.Hash

	if GaugeType == data.MType {
		ret.MType = proto.Metric_GAUGE
		val := float64(*data.Value)
		ret.Value = &val
	} else {
		ret.MType = proto.Metric_COUNTER
		val := int64(*data.Delta)
		ret.Delta = &val
	}
	return &ret
}

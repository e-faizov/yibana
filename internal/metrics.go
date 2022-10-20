package internal

import (
	"math/rand"
	"runtime"
	"sync"
)

type NamedGauge struct {
	Name  string
	Value Gauge
}

type NamedCounter struct {
	Name  string
	Value Counter
}

type Metric struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *Counter `json:"delta,omitempty"`
	Value *Gauge   `json:"value,omitempty"`
}

func (m *Metric) SetGauge(g Gauge) {
	m.Value = &g
	m.Delta = nil
	m.MType = "gauge"
}

func (m *Metric) SetCounter(c Counter) {
	m.Value = nil
	m.Delta = &c
	m.MType = "counter"
}

type Metrics struct {
	mtx          sync.Mutex
	data         []Metric
	currentCount Counter
}

func (m *Metrics) Update() {

	m.currentCount++
	var tmp []Metric

	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	addGauge := func(nm string, v Gauge) {
		tmpMetric := Metric{
			ID: nm,
		}
		tmpMetric.SetGauge(v)
		tmp = append(tmp, tmpMetric)
	}

	addGauge("alloc", Gauge(rtm.Alloc))
	addGauge("BuckHashSys", Gauge(rtm.BuckHashSys))
	addGauge("Frees", Gauge(rtm.Frees))
	addGauge("GCCPUFraction", Gauge(rtm.GCCPUFraction))

	addGauge("GCSys", Gauge(rtm.GCSys))
	addGauge("HeapAlloc", Gauge(rtm.HeapAlloc))
	addGauge("HeapIdle", Gauge(rtm.HeapIdle))
	addGauge("HeapInuse", Gauge(rtm.HeapInuse))

	addGauge("HeapObjects", Gauge(rtm.HeapObjects))
	addGauge("HeapReleased", Gauge(rtm.HeapReleased))
	addGauge("HeapSys", Gauge(rtm.HeapSys))
	addGauge("LastGC", Gauge(rtm.LastGC))

	addGauge("Lookups", Gauge(rtm.Lookups))
	addGauge("MCacheInuse", Gauge(rtm.MCacheInuse))
	addGauge("MCacheSys", Gauge(rtm.MCacheSys))
	addGauge("MSpanInuse", Gauge(rtm.MSpanInuse))

	addGauge("MSpanSys", Gauge(rtm.MSpanSys))
	addGauge("Mallocs", Gauge(rtm.Mallocs))
	addGauge("NextGC", Gauge(rtm.NextGC))
	addGauge("NumForcedGC", Gauge(rtm.NumForcedGC))

	addGauge("NumGC", Gauge(rtm.NumGC))
	addGauge("OtherSys", Gauge(rtm.OtherSys))
	addGauge("PauseTotalNs", Gauge(rtm.PauseTotalNs))
	addGauge("StackInuse", Gauge(rtm.StackInuse))

	addGauge("StackSys", Gauge(rtm.StackSys))
	addGauge("Sys", Gauge(rtm.Sys))
	addGauge("TotalAlloc", Gauge(rtm.TotalAlloc))
	addGauge("RandomValue", Gauge(rand.Float64()))

	tmpMetric := Metric{
		ID: "PollCount",
	}
	tmpMetric.SetCounter(m.currentCount)
	tmp = append(tmp, tmpMetric)

	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.data = append(m.data, tmp...)
}

func (m *Metrics) Front() (Metric, bool) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if len(m.data) == 0 {
		return Metric{}, false
	}
	return m.data[0], true
}

func (m *Metrics) Pop() {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if len(m.data) != 0 {
		m.data = m.data[1:]
	}
}

package internal

import (
	"fmt"
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
	Name string
	g    *Gauge
	c    *Counter
	t    string
}

func (m *Metric) SetGauge(g Gauge) {
	m.g = &g
	m.c = nil
	m.t = "gauge"
}

func (m *Metric) SetCounter(c Counter) {
	m.g = nil
	m.c = &c
	m.t = "counter"
}

func (m *Metric) GetType() string {
	return m.t
}

func (m *Metric) ToString() string {
	if m.g != nil {
		return fmt.Sprintf("%.3f", *m.g)
	} else if m.c != nil {
		return fmt.Sprintf("%d", *m.c)
	}
	return ""
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
			Name: nm,
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
		Name: "PollCount",
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

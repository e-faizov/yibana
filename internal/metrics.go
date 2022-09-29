package internal

import (
	"math/rand"
	"runtime"
)

type NamedGauge struct {
	Name  string
	Value Gauge
}

type NamedCounter struct {
	Name  string
	Value Counter
}

type Metrics struct {
	Alloc         NamedGauge
	BuckHashSys   NamedGauge
	Frees         NamedGauge
	GCCPUFraction NamedGauge
	GCSys         NamedGauge
	HeapAlloc     NamedGauge
	HeapIdle      NamedGauge
	HeapInuse     NamedGauge
	HeapObjects   NamedGauge
	HeapReleased  NamedGauge
	HeapSys       NamedGauge
	LastGC        NamedGauge
	Lookups       NamedGauge
	MCacheInuse   NamedGauge
	MCacheSys     NamedGauge
	MSpanInuse    NamedGauge
	MSpanSys      NamedGauge
	Mallocs       NamedGauge
	NextGC        NamedGauge
	NumForcedGC   NamedGauge
	NumGC         NamedGauge
	OtherSys      NamedGauge
	PauseTotalNs  NamedGauge
	StackInuse    NamedGauge
	StackSys      NamedGauge
	Sys           NamedGauge
	TotalAlloc    NamedGauge
	PollCount     NamedCounter
	RandomValue   NamedGauge
}

func (m *Metrics) Update() {
	m.PollCount.Value++
	m.RandomValue.Value = Gauge(rand.Float64())

	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	m.Alloc.Value = Gauge(rtm.Alloc)

	m.BuckHashSys.Value = Gauge(rtm.BuckHashSys)
	m.Frees.Value = Gauge(rtm.Frees)
	m.GCCPUFraction.Value = Gauge(rtm.GCCPUFraction)
	m.GCSys.Value = Gauge(rtm.GCSys)
	m.HeapAlloc.Value = Gauge(rtm.HeapAlloc)
	m.HeapIdle.Value = Gauge(rtm.HeapIdle)
	m.HeapInuse.Value = Gauge(rtm.HeapInuse)
	m.HeapObjects.Value = Gauge(rtm.HeapObjects)
	m.HeapReleased.Value = Gauge(rtm.HeapReleased)
	m.HeapSys.Value = Gauge(rtm.HeapSys)
	m.LastGC.Value = Gauge(rtm.LastGC)
	m.Lookups.Value = Gauge(rtm.Lookups)
	m.MCacheInuse.Value = Gauge(rtm.MCacheInuse)
	m.MCacheSys.Value = Gauge(rtm.MCacheSys)
	m.MSpanInuse.Value = Gauge(rtm.MSpanInuse)
	m.MSpanSys.Value = Gauge(rtm.MSpanSys)
	m.Mallocs.Value = Gauge(rtm.Mallocs)
	m.NextGC.Value = Gauge(rtm.NextGC)
	m.NumForcedGC.Value = Gauge(rtm.NumForcedGC)
	m.NumGC.Value = Gauge(rtm.NumGC)
	m.OtherSys.Value = Gauge(rtm.OtherSys)
	m.PauseTotalNs.Value = Gauge(rtm.PauseTotalNs)
	m.StackInuse.Value = Gauge(rtm.StackInuse)
	m.StackSys.Value = Gauge(rtm.StackSys)
	m.Sys.Value = Gauge(rtm.Sys)
	m.TotalAlloc.Value = Gauge(rtm.TotalAlloc)
}

func NewMetrics() Metrics {
	res := initMetrics()
	return res
}

func initMetrics() Metrics {
	var res Metrics

	res.Alloc = NamedGauge{Name: "alloc"}
	res.BuckHashSys = NamedGauge{Name: "BuckHashSys"}
	res.Frees = NamedGauge{Name: "Frees"}
	res.GCCPUFraction = NamedGauge{Name: "GCCPUFraction"}
	res.GCSys = NamedGauge{Name: "GCSys"}
	res.HeapAlloc = NamedGauge{Name: "HeapAlloc"}
	res.HeapIdle = NamedGauge{Name: "HeapIdle"}
	res.HeapInuse = NamedGauge{Name: "HeapInuse"}
	res.HeapObjects = NamedGauge{Name: "HeapObjects"}
	res.HeapReleased = NamedGauge{Name: "HeapReleased"}
	res.HeapSys = NamedGauge{Name: "HeapSys"}
	res.LastGC = NamedGauge{Name: "LastGC"}
	res.Lookups = NamedGauge{Name: "Lookups"}
	res.MCacheInuse = NamedGauge{Name: "MCacheInuse"}
	res.MCacheSys = NamedGauge{Name: "MCacheSys"}
	res.MSpanInuse = NamedGauge{Name: "MSpanInuse"}
	res.MSpanSys = NamedGauge{Name: "MSpanSys"}
	res.Mallocs = NamedGauge{Name: "Mallocs"}
	res.NextGC = NamedGauge{Name: "NextGC"}
	res.NumForcedGC = NamedGauge{Name: "NumForcedGC"}
	res.NumGC = NamedGauge{Name: "NumGC"}
	res.OtherSys = NamedGauge{Name: "OtherSys"}
	res.PauseTotalNs = NamedGauge{Name: "PauseTotalNs"}
	res.StackInuse = NamedGauge{Name: "StackInuse"}
	res.StackSys = NamedGauge{Name: "StackSys"}
	res.Sys = NamedGauge{Name: "Sys"}
	res.TotalAlloc = NamedGauge{Name: "TotalAlloc"}
	res.PollCount = NamedCounter{Name: "PollCount"}
	res.RandomValue = NamedGauge{Name: "RandomValue"}

	return res
}

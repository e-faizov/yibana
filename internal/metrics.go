package internal

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"runtime"
	"sync"

	"github.com/shirou/gopsutil/v3/mem"
)

// CalcGaugeHash - функция подсчета хэша для метрики с типом Gauge
func CalcGaugeHash(id string, g Gauge, k string) string {
	str := fmt.Sprintf("%s:gauge:%f", id, g)
	hash := calcHash(str, k)
	return hash
}

// CalcCounterHash - функция подсчета хэша для метрики с типом Counter
func CalcCounterHash(id string, d Counter, k string) string {
	return calcHash(fmt.Sprintf("%s:counter:%d", id, d), k)
}

func calcHash(s string, k string) string {
	h := hmac.New(sha256.New, []byte(k))
	h.Write([]byte(s))
	hash := h.Sum(nil)
	return fmt.Sprintf("%x", hash)
}

// Metric - структура метрики
type Metric struct {
	// ID - имя метрики
	ID string `json:"id"`
	// MType - тип метрик Counter или Gauge
	MType string `json:"type"`
	// Delta - значения для метрики типа Counter
	Delta *Counter `json:"delta,omitempty"`
	// Value - значения для метрики типа Gauge
	Value *Gauge `json:"value,omitempty"`
	// Hash - хэш метрики
	Hash string `json:"hash,omitempty"`
}

func (m Metric) String() string {
	if m.MType == GaugeType && m.Value != nil {
		return fmt.Sprintf("type: %s, name: %s, value: %f, hash: %s", m.MType, m.ID, *m.Value, m.Hash)
	} else if m.MType == CounterType && m.Delta != nil {
		return fmt.Sprintf("type: %s, name: %s, delta: %d, hash: %s", m.MType, m.ID, *m.Delta, m.Hash)
	}
	return "empty metric id: " + m.ID
}

func (m *Metric) SetGauge(g Gauge) {
	m.Value = &g
	m.Delta = nil
	m.MType = GaugeType
}

func (m *Metric) SetGaugeWithHash(g Gauge, key string) {
	m.SetGauge(g)
	m.Hash = CalcGaugeHash(m.ID, *m.Value, key)
}

func (m *Metric) SetCounter(c Counter) {
	m.Value = nil
	m.Delta = &c
	m.MType = CounterType
}

func (m *Metric) SetCounterWithHash(c Counter, key string) {
	m.SetCounter(c)
	m.Hash = CalcCounterHash(m.ID, *m.Delta, key)
}

// Metrics - структура для сбора метрик
type Metrics struct {
	Key          string
	data         []Metric
	currentCount Counter
	mtx          sync.Mutex
}

// Update - метод сбора новых метрик
func (m *Metrics) Update() error {

	m.currentCount++
	var tmp []Metric

	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	addGauge := func(nm string, v Gauge) {
		tmpMetric := Metric{
			ID: nm,
		}
		tmpMetric.SetGaugeWithHash(v, m.Key)
		tmp = append(tmp, tmpMetric)
	}

	addGauge("Alloc", Gauge(rtm.Alloc))
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
	tmpMetric.SetCounterWithHash(m.currentCount, m.Key)
	tmp = append(tmp, tmpMetric)

	v, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("update error mem.VirtualMemory %w", err)
	}

	addGauge("TotalMemory", Gauge(v.Total))
	addGauge("FreeMemory", Gauge(v.Free))
	addGauge("CPUutilization1", Gauge(v.UsedPercent))

	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.data = append(m.data, tmp...)
	return nil
}

// Batch - метод выделения метрик в отдельный слайс
func (m *Metrics) Batch() []Metric {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if len(m.data) == 0 {
		return []Metric{}
	}
	ret := m.data
	m.data = make([]Metric, 0, 100)
	return ret
}

// Front - метод чтения первой метрики в списке
func (m *Metrics) Front() (Metric, bool) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if len(m.data) == 0 {
		return Metric{}, false
	}
	return m.data[0], true
}

// Pop - метод удаления первой метрики в списке
func (m *Metrics) Pop() {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if len(m.data) != 0 {
		m.data = m.data[1:]
	}
}

package metrics

import (
	"sync"
	"time"
)

type MetricStore struct {
	metrics map[string]SystemMetrics
	mutex   sync.RWMutex
}

var GlobalStore *MetricStore

func InitStore() {
	GlobalStore = &MetricStore{
		metrics: make(map[string]SystemMetrics),
	}
}

func (s *MetricStore) Update(serverID string, m SystemMetrics) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// Ensure LastUpdate is set if not present (though agent should send it)
	if m.LastUpdate.IsZero() {
		m.LastUpdate = time.Now()
	}
	s.metrics[serverID] = m
}

func (s *MetricStore) Get(serverID string) (SystemMetrics, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	m, ok := s.metrics[serverID]
	return m, ok
}

func (s *MetricStore) GetAll() map[string]SystemMetrics {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	// Return a copy to avoid race conditions
	result := make(map[string]SystemMetrics)
	for k, v := range s.metrics {
		result[k] = v
	}
	return result
}

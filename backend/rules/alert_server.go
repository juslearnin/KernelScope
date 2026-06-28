package rules

import (
	"sync"

	"kernelscope/models"
)

type AlertStore struct {
	mu       sync.Mutex
	alerts   []models.Alert
	capacity int
}

func NewAlertStore(capacity int) *AlertStore {
	return &AlertStore{
		alerts:   make([]models.Alert, 0, capacity),
		capacity: capacity,
	}
}

func (s *AlertStore) Add(alert models.Alert) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.alerts = append([]models.Alert{alert}, s.alerts...)

	if len(s.alerts) > s.capacity {
		s.alerts = s.alerts[:s.capacity]
	}
}

func (s *AlertStore) List() []models.Alert {
	s.mu.Lock()
	defer s.mu.Unlock()

	return append([]models.Alert{}, s.alerts...)
}
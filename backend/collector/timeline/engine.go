package timeline

import (
	"sync"

	"kernelscope/models"
)

type TimelineEngine struct {
	mu       sync.Mutex
	events   []models.Event
	capacity int
	next     int
	full     bool
}

func NewTimelineEngine(capacity int) *TimelineEngine {
	return &TimelineEngine{
		events:   make([]models.Event, capacity),
		capacity: capacity,
	}
}

func (t *TimelineEngine) Add(event models.Event) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.events[t.next] = event
	t.next = (t.next + 1) % t.capacity

	if t.next == 0 {
		t.full = true
	}
}

func (t *TimelineEngine) List() []models.Event {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.full {
		return append([]models.Event{}, t.events[:t.next]...)
	}

	result := append([]models.Event{}, t.events[t.next:]...)
	result = append(result, t.events[:t.next]...)

	return result
}
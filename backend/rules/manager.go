package rules

import (
	"fmt"
	"sync"
)

type AlertManager struct {
	mu     sync.Mutex
	active map[string]bool
}

func NewAlertManager() *AlertManager {
	return &AlertManager{
		active: make(map[string]bool),
	}
}

func alertKey(rule string, pid int) string {
	return fmt.Sprintf("%s:%d", rule, pid)
}

func (m *AlertManager) ShouldEmit(rule string, pid int) bool {

	m.mu.Lock()
	defer m.mu.Unlock()

	key := alertKey(rule, pid)

	if m.active[key] {
		return false
	}

	m.active[key] = true
	return true
}
func (m *AlertManager) Resolve(rule string, pid int) {

	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.active, alertKey(rule, pid))

}
func (m *AlertManager) ResolveMissing(currentMatches map[string]bool) []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	var resolved []string

	for key := range m.active {
		if !currentMatches[key] {
			resolved = append(resolved, key)
			delete(m.active, key)
		}
	}

	return resolved
}

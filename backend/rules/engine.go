package rules

import (
	"time"
	"strconv"

	"kernelscope/models"
)

type Rule struct {
	Name      string
	Level     models.AlertLevel
	Condition func(models.Process) bool
	Message   func(models.Process) string
}

var DefaultRules = []Rule{
	{
		Name:  "HIGH_CPU",
		Level: models.Critical,
		Condition: func(p models.Process) bool {
			return p.CPUPercent >= 20.0 // 20%
		},
		Message: func(p models.Process) string {
			return "CPU usage exceeded 80%"
		},
	},
	{
		Name:  "HIGH_RAM",
		Level: models.Warning,
		Condition: func(p models.Process) bool {
			return p.MemoryKB >= 1024*1024 // 1 GB
		},
		Message: func(p models.Process) string {
			return "Memory usage exceeded 1 GB"
		},
	},
	{
		Name:  "MANY_OPEN_FILES",
		Level: models.Warning,
		Condition: func(p models.Process) bool {
			return len(p.OpenFiles) >= 100
		},
		Message: func(p models.Process) string {
			return "Large number of open files"
		},
	},
	{
		Name:  "HIGH_CONNECTIONS",
		Level: models.Warning,
		Condition: func(p models.Process) bool {
			return len(p.Connections) >= 50
		},
		Message: func(p models.Process) string {
			return "Large number of network connections"
		},
	},
}
func CurrentMatches(processes []models.Process) map[string]bool {
	matches := make(map[string]bool)

	for _, process := range processes {
		for _, rule := range DefaultRules {
			if rule.Condition(process) {
				key := rule.Name + ":" + strconv.Itoa(process.PID)
				matches[key] = true
			}
		}
	}

	return matches
}
func Evaluate(processes []models.Process) []models.Alert {
	var alerts []models.Alert

	now := time.Now().UnixMilli()

	for _, process := range processes {
		for _, rule := range DefaultRules {

			if rule.Condition(process) {

				alerts = append(alerts, models.Alert{
					Timestamp: now,
					Level:     rule.Level,
					Rule:      rule.Name,
					PID:       process.PID,
					Process:   process.Name,
					Message:   rule.Message(process),
				})

			}
		}
	}

	return alerts
}
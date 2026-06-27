package timeline

import (
	"strconv"
	"time"
    "strings"
	"kernelscope/models"
)

type ProcessSnapshot struct {
	processes map[int]models.Process
}

func NewProcessSnapshot(processes []models.Process) ProcessSnapshot {
	snapshot := ProcessSnapshot{
		processes: make(map[int]models.Process),
	}

	for _, process := range processes {
		snapshot.processes[process.PID] = process
	}

	return snapshot
}

func DiffProcesses(
	previous ProcessSnapshot,
	current ProcessSnapshot,
) []models.Event {
	var events []models.Event
	now := time.Now().UnixMilli()

	for pid, process := range current.processes {
		if _, existed := previous.processes[pid]; !existed {
			events = append(events, models.Event{
				Timestamp: now,
				Type:      models.ProcessStarted,
				PID:       pid,
				Process:   process.Name,
				Importance: classifyProcessImportance(process),
				Details: map[string]string{
					"ppid":    strconv.Itoa(process.PPID),
					"cmdline": process.Cmdline,
				},
			})
		}
	}

	for pid, process := range previous.processes {
		if _, existsNow := current.processes[pid]; !existsNow {
			events = append(events, models.Event{
				Timestamp: now,
				Type:      models.ProcessExited,
				PID:       pid,
				Process:   process.Name,
				Importance: classifyProcessImportance(process),
				Details: map[string]string{
					"cmdline": process.Cmdline,
				},
			})
		}
	}

	return events
}
func classifyProcessImportance(process models.Process) string {
	name := process.Name
	cmd := process.Cmdline

	lowValue := []string{"curl", "sh", "watch", "grep", "head", "cat"}

	for _, item := range lowValue {
		if name == item {
			return "LOW"
		}
	}

	if strings.Contains(cmd, "KernelScope") || strings.Contains(cmd, "/api/processes") {
		return "LOW"
	}

	highValue := []string{"node", "npm", "python", "go", "chrome", "code", "docker", "mongod"}

	for _, item := range highValue {
		if strings.Contains(strings.ToLower(name), item) ||
			strings.Contains(strings.ToLower(cmd), item) {
			return "HIGH"
		}
	}

	return "NORMAL"
}
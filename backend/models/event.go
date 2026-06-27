package models

type EventType string

const (
	ProcessStarted  EventType = "PROCESS_STARTED"
	ProcessExited   EventType = "PROCESS_EXITED"
	FileOpened      EventType = "FILE_OPENED"
	FileClosed      EventType = "FILE_CLOSED"
	ConnectionOpen  EventType = "CONNECTION_OPENED"
	ConnectionClose EventType = "CONNECTION_CLOSED"
	CPUChanged      EventType = "CPU_CHANGED"
	RAMChanged      EventType = "RAM_CHANGED"
)

type Event struct {
	Timestamp int64             `json:"timestamp"`
	Type      EventType         `json:"type"`
	PID       int               `json:"pid"`
    Importance string            `json:"importance"`
	Process   string            `json:"process"`
	Details   map[string]string `json:"details"`
}
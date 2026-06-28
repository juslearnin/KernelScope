package models

type AlertLevel string

const (
	Info AlertLevel = "INFO"
	Warning AlertLevel = "WARNING"
	Critical AlertLevel = "CRITICAL"
)

type Alert struct {
	Timestamp int64
	Level AlertLevel

	Rule string

	PID int

	Process string

	Message string
}
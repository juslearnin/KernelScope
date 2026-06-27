package models

type FileDescriptor struct {
	FD     int    `json:"fd"`
	Target string `json:"target"`
	Type   string `json:"type"`
}

type Process struct {
	PID        int              `json:"pid"`
	PPID       int              `json:"ppid"`
	Name       string           `json:"name"`
	Cmdline    string           `json:"cmdline"`
	State      string           `json:"state"`
	Threads    int              `json:"threads"`
	MemoryKB   int              `json:"memoryKB"`
	CPUTime    uint64           `json:"cpuTime"`
	CPUPercent float64          `json:"cpuPercent"`
	OpenFiles  []FileDescriptor `json:"openFiles"`
}
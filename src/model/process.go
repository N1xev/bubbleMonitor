package model

import (
	"time"
)

// ProcessSnapshot stores a point-in-time resource snapshot for a process
type ProcessSnapshot struct {
	Timestamp time.Time
	Cpu       float64
	Memory    float64
}

// ProcessInfo holds information about a running process
type ProcessInfo struct {
	Name        string
	Pid         int32
	Cpu         float64
	Memory      float64
	Status      string
	Username    string
	CreateTime  int64
	Cmdline     string
	MemoryBytes uint64
	Nice        int32 // Priority
	Ppid        int32 // Parent PID
}

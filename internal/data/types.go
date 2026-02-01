package data

import (
	"time"
)

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

	// Cached lowercase fields for filtering performance
	NameLower     string
	UsernameLower string
	CmdlineLower  string
}

// ProcessSnapshot stores a point-in-time resource snapshot for a process
type ProcessSnapshot struct {
	Timestamp time.Time
	Cpu       float64
	Memory    float64
}

// DiskPartition holds information about a disk partition
type DiskPartition struct {
	Mountpoint string
	Device     string
	Fstype     string
	Total      uint64
	Used       uint64
	UsedPct    float64
}

// GpuInfo holds information about a GPU
type GpuInfo struct {
	Name        string
	Driver      string
	MemoryTotal string
	MemoryUsed  string
}

// Toast Levels
const (
	ToastInfo    = "info"
	ToastError   = "error"
	ToastWarn    = "warn"
	ToastSuccess = "success"
)

type Toast struct {
	ID        int64
	Message   string
	Level     string
	StartTime time.Time
	Duration  time.Duration
}

type ServiceInfo struct {
	Name        string
	Status      string
	Description string
}

type ConnectionInfo struct {
	LocalAddr  string
	RemoteAddr string
	State      string
	Pid        int32
	Protocol   string
}

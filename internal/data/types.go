package data

import (
	"time"
)

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

	NameLower     string
	UsernameLower string
	CmdlineLower  string
}

type ProcessSnapshot struct {
	Timestamp time.Time
	Cpu       float64
	Memory    float64
}

type DiskPartition struct {
	Mountpoint string
	Device     string
	Fstype     string
	Total      uint64
	Used       uint64
	UsedPct    float64
}

type GpuInfo struct {
	Name        string
	Driver      string
	MemoryTotal string
	MemoryUsed  string
}

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

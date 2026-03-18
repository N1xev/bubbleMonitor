package data

import (
	"time"
)

type ProcessInfo struct {
	NameLower     string
	UsernameLower string
	CmdlineLower  string
	Name          string
	Status        string
	Username      string
	Cmdline       string
	Pid           int32
	CreateTime    int64
	MemoryBytes   uint64
	Cpu           float64
	Memory        float64
	Nice          int32 // Priority
	Ppid          int32 // Parent PID
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
	Slot        string
	Type        string
	Vendor      string
	Temperature string
	FanSpeed    string
	PowerUsage  string
	Utilization string
	ClockSpeed  string
	ErrorMsg    string
}

const (
	ToastInfo    = "info"
	ToastError   = "error"
	ToastWarn    = "warn"
	ToastSuccess = "success"
)

type Toast struct {
	StartTime time.Time
	Message   string
	Level     string
	Duration  time.Duration
	ID        int64
}

type ServiceInfo struct {
	Description string
	Name        string
	Status      string
}

type ConnectionInfo struct {
	LocalAddr  string
	RemoteAddr string
	State      string
	Protocol   string
	Pid        int32
}

type HardwareCapabilities struct {
	HasNvidiaGPU         bool
	HasAmdGPU            bool
	HasIntelGPU          bool
	HasBattery           bool
	HasNetworkInterfaces bool
	HasDiskIO            bool
	HasServices          bool
	HasTempSensors       bool
	HasDocker            bool
	HasKubernetes        bool
}

// ContainerInfo represents container metrics
type ContainerInfo struct {
	Name       string
	ID         string
	Status     string
	State      string
	Image      string
	Created    string
	CPUPercent float64
	MemUsage   uint64
	MemLimit   uint64
	MemPct     float64
	NetRx      uint64
	NetTx      uint64
	Type       string // "docker" or "kubernetes"
}

// K8sPodInfo represents Kubernetes pod metrics
type K8sPodInfo struct {
	Name       string
	Namespace  string
	Status     string
	Ready      string
	Restarts   int
	CPUReq     string
	MemReq     string
	CPULim     string
	MemLim     string
	CPUPercent float64
	MemPct     float64
	Node       string
	Age        string
}

// VmInfo represents VM/hypervisor information
type VmInfo struct {
	Hypervisor     string
	Type           string
	HostName       string
	VirtMemory     uint64
	VirtMemoryUsed uint64
	VirtCPU        int
	VirtCPUUsed    int
}

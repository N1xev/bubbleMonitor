package messages

import (
	"time"

	"github.com/distatus/battery"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"

	// Import shared data types
	"github.com/N1xev/bubbleMonitor/src/data"
)

// Message types for Bubble Tea
type TickMsg time.Time

// CpuMemMsg contains fast-updating metrics (CPU, Memory, Swap)
type CpuMemMsg struct {
	Cpu        float64
	CpuPerCore []float64
	Memory     float64
	Swap       float64
	LoadAvg    *load.AvgStat
	MemInfo    *mem.VirtualMemoryStat
	SwapInfo   *mem.SwapMemoryStat
}

// DiskNetMsg contains slow-updating metrics (Disk, Network)
type DiskNetMsg struct {
	Disk    float64
	NetSent uint64
	NetRecv uint64
}

type MetricsMsg struct {
	Cpu        float64
	CpuPerCore []float64
	Memory     float64
	Disk       float64
	Swap       float64
	NetSent    uint64
	NetRecv    uint64
	LoadAvg    *load.AvgStat
	MemInfo    *mem.VirtualMemoryStat // Cached for render
	SwapInfo   *mem.SwapMemoryStat    // Cached for render
}

type ProcessesMsg []data.ProcessInfo
type HostInfoMsg *host.InfoStat
type DiskInfoMsg []data.DiskPartition // Using data.DiskPartition
type GpuInfoMsg []data.GpuInfo        // Using data.GpuInfo
type DiskIOMsg map[string]disk.IOCountersStat
type TempMsg []host.TemperatureStat
type NetworkInterfacesMsg []net.IOCountersStat
type BatteryMsg []*battery.Battery

// Control Messages
type PriorityChangeMsg struct {
	Pid      int32
	Priority int32
	Err      error
}

type ProcessControlMsg struct {
	Pid    int32
	Action string // "suspend" or "resume"
	Err    error
}

type OpenFilesMsg struct {
	Pid   int32
	Files []process.OpenFilesStat
	Err   error
}

// Toast Messages
type ToastMsg struct {
	Message  string
	Level    string
	Duration time.Duration
}

type ToastTimeoutMsg struct {
	ID int64
}

// KillProcessMsg is sent when a process kill is requested
type KillProcessMsg struct {
	Pid     int32
	Success bool
	Error   string
}

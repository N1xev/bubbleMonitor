package data

import (
	"time"

	"github.com/N1xev/bubbleMonitor/internal/config"
	"github.com/distatus/battery"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// AppState holds all the application state data
type AppState struct {
	Width          int
	Height         int
	Cpu            float64
	CpuPerCore     []float64
	Memory         float64
	Disk           float64
	Swap           float64
	HistoryLength  int
	CpuHistory     *RingBuffer
	MemHistory     *RingBuffer
	NetHistory     *RingBuffer
	SwapHistory    *RingBuffer
	SelectedTab    int
	Processes      []ProcessInfo
	LastNetSent    uint64
	LastNetRecv    uint64
	NetSentRate    float64
	NetRecvRate    float64
	HostInfo       *host.InfoStat
	DiskPartitions []DiskPartition
	SortBy         string
	LoadAvg        *load.AvgStat
	StartTime      time.Time
	Paused         bool
	ShowHelp       bool
	GpuInfo        []GpuInfo
	MemInfo        *mem.VirtualMemoryStat // Cached memory info
	SwapInfo       *mem.SwapMemoryStat    // Cached swap info
	CpuInfoStatic  []cpu.InfoStat         // Static CPU info

	// Temperature
	Sensors     []host.TemperatureStat
	CpuTemp     float64
	HistoryTemp *RingBuffer

	// Network
	NetworkInterfaces     []net.IOCountersStat
	LastNetworkInterfaces map[string]net.IOCountersStat

	// Battery
	Battery []*battery.Battery

	// Disk I/O
	DiskIO        map[string]disk.IOCountersStat
	LastDiskIO    map[string]disk.IOCountersStat
	DiskReadRate  float64
	DiskWriteRate float64
	DiskHORead    *RingBuffer
	DiskHOWrite   *RingBuffer

	// Process navigation and filtering
	SelectedProcess     int
	ProcessCount        int // Lightweight count for Overview
	ProcessScrollOffset int
	ProcessFilter       string
	FilterMode          bool
	ShowKillDialog      bool
	KillTargetPid       int32
	KillTargetName      string

	// Alerts & Configuration
	Config       config.AppConfig
	AlertManager *AlertManager
	ShowSettings bool
	SettingsEdit bool
	SettingsSel  config.MetricType
	SettingsIdx  int

	LastError     string
	LastErrorTime time.Time
	TickCount     uint64 // Throttling counter

	// Toasts
	Toasts      []Toast
	NextToastID int64

	// App-side Suspend State Tracking
	SuspendedState map[int32]bool

	// Open Files Inspector
	ShowOpenFiles         bool
	OpenFilesList         []process.OpenFilesStat
	OpenFilesPid          int32
	OpenFilesScrollOffset int
	OpenFilesView         SimpleViewport

	// Services & Connections
	Services      []ServiceInfo
	Connections   []ConnectionInfo
	SystemLogs    []string
	RemoteUptimes map[string]string

	// Process Tree View
	TreeView       bool
	CollapsedPids  map[int32]bool
	BookmarkedPids map[int32]bool

	// Cached process views (optimization)
	CachedVisibleProcs []ProcessInfo
	CachedTreeIndents  map[int32]int
	ProcessCacheDirty  bool

	// Enhanced Visualization
	ChartType string

	// Customization
	Theme             string
	RefreshRate       int
	BorderType        string
	BorderStyle       string
	BackgroundOpaque  bool
	LastConfigModTime time.Time
	ActiveTabs        []string

	// UI State
	CpuCoreScrollOffset int

	// Analysis & History
	ProcessHistory map[int32]*RingBuffer // History for tracked processes
	HealthScore    int
}

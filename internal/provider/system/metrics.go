package system

import (
	"runtime"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"

	"github.com/N1xev/bubbleMonitor/internal/msg"
)

// TickCmd returns a command that ticks after the specified duration
func TickCmd(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return msg.TickMsg(t)
	})
}

// FastMetricsCmd fetches fast-changing system metrics (CPU, Memory)
func FastMetricsCmd() tea.Cmd {
	return func() tea.Msg {
		cpuPercent, _ := cpu.Percent(0, false)
		cpuVal := 0.0
		if len(cpuPercent) > 0 {
			cpuVal = cpuPercent[0]
		}
		cpuPerCore, _ := cpu.Percent(0, true)
		memInfo, _ := mem.VirtualMemory()
		swapInfo, _ := mem.SwapMemory()
		loadAvg, _ := load.Avg()

		var memUsed, swapUsed float64
		if memInfo != nil {
			memUsed = memInfo.UsedPercent
		}
		if swapInfo != nil {
			swapUsed = swapInfo.UsedPercent
		}

		return msg.CpuMemMsg{
			Cpu:        cpuVal,
			CpuPerCore: cpuPerCore,
			Memory:     memUsed,
			Swap:       swapUsed,
			LoadAvg:    loadAvg,
			MemInfo:    memInfo,
			SwapInfo:   swapInfo,
		}
	}
}

// SlowMetricsCmd fetches slow-changing system metrics (Disk, Network)
func SlowMetricsCmd() tea.Cmd {
	return func() tea.Msg {
		root := "/"
		if runtime.GOOS == "windows" {
			root = "C:"
		}
		diskInfo, _ := disk.Usage(root)

		diskPercent := 0.0
		if diskInfo != nil {
			diskPercent = diskInfo.UsedPercent
		}

		netIO, _ := net.IOCounters(false)
		var netSent, netRecv uint64
		if len(netIO) > 0 {
			netSent = netIO[0].BytesSent
			netRecv = netIO[0].BytesRecv
		}

		return msg.DiskNetMsg{
			Disk:    diskPercent,
			NetSent: netSent,
			NetRecv: netRecv,
		}
	}
}

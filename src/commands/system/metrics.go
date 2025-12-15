package system

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"

	"github.com/N1xev/bubbleMonitor/src/messages"
)

// TickCmd returns a command that ticks after the specified duration
func TickCmd(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return messages.TickMsg(t)
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

		return messages.CpuMemMsg{
			Cpu:        cpuVal,
			CpuPerCore: cpuPerCore,
			Memory:     memInfo.UsedPercent,
			Swap:       swapInfo.UsedPercent,
			LoadAvg:    loadAvg,
			MemInfo:    memInfo,
			SwapInfo:   swapInfo,
		}
	}
}

// SlowMetricsCmd fetches slow-changing system metrics (Disk, Network)
func SlowMetricsCmd() tea.Cmd {
	return func() tea.Msg {
		diskInfo, _ := disk.Usage("/")

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

		return messages.DiskNetMsg{
			Disk:    diskPercent,
			NetSent: netSent,
			NetRecv: netRecv,
		}
	}
}

// Deprecated: Kept for compatibility if anything still calls it, but redirects to Fast
func MetricsCmd() tea.Cmd {
	return func() tea.Msg {
		// Just perform a fast fetch and return as legacy MetricsMsg (incomplete data, but safe)
		// Ideally this should not be called anymore.
		msg := FastMetricsCmd()().(messages.CpuMemMsg)
		return messages.MetricsMsg{
			Cpu:        msg.Cpu,
			CpuPerCore: msg.CpuPerCore,
			Memory:     msg.Memory,
			Swap:       msg.Swap,
			LoadAvg:    msg.LoadAvg,
			MemInfo:    msg.MemInfo,
			SwapInfo:   msg.SwapInfo,
		}
	}
}

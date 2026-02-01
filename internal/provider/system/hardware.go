package system

import (
	"os/exec"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/distatus/battery"
	"github.com/shirou/gopsutil/v3/host"

	"github.com/N1xev/bubbleMonitor/internal/data"
	"github.com/N1xev/bubbleMonitor/internal/msg"
)

// HostInfoCmd fetches host information
func HostInfoCmd() tea.Cmd {
	return func() tea.Msg {
		info, err := host.Info()
		return msg.HostInfoMsg{
			Info: info,
			Err:  err,
		}
	}
}

// GpuInfoCmd fetches GPU information (NVIDIA only)
func GpuInfoCmd() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("nvidia-smi", "--query-gpu=name,driver_version,memory.total,memory.used", "--format=csv,noheader,nounits")
		out, err := cmd.Output()
		var gpuList []data.GpuInfo

		if err != nil {
			// Fallback: lspci
			cmd = exec.Command("lspci")
			out, err = cmd.Output()
			if err == nil {
				lines := strings.Split(string(out), "\n")
				for _, line := range lines {
					if strings.Contains(strings.ToLower(line), "vga") || strings.Contains(strings.ToLower(line), "3d controller") {
						parts := strings.Split(line, ": ")
						if len(parts) > 1 {
							// 00:02.0 VGA compatible controller: Intel Corporation ...
							// We want "Intel Corporation ..."
							// parts[0] is "00", parts[1] is "02.0 VGA..."
							// lspci output format: "Slot Class: Name"
							// We want everything after the second colon if possible, or just the string.
							name := parts[len(parts)-1]
							gpuList = append(gpuList, data.GpuInfo{Name: name, Driver: "Unknown", MemoryTotal: "N/A", MemoryUsed: "N/A"})
						}
					}
				}
				return msg.GpuInfoMsg(gpuList)
			}
			return nil
		}
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				continue
			}
			parts := strings.Split(line, ", ")
			if len(parts) < 4 {
				continue
			}
			gpuList = append(gpuList, data.GpuInfo{
				Name:        parts[0],
				Driver:      parts[1],
				MemoryTotal: parts[2],
				MemoryUsed:  parts[3],
			})
		}
		return msg.GpuInfoMsg(gpuList)
	}
}

// TempCmd fetches temperature sensors
func TempCmd() tea.Cmd {
	return func() tea.Msg {
		temps, err := host.SensorsTemperatures()
		if err != nil {
			return msg.TempMsg{}
		}
		return msg.TempMsg(temps)
	}
}

// BatteryCmd fetches battery information
func BatteryCmd() tea.Cmd {
	return func() tea.Msg {
		batt, err := battery.GetAll()
		if err != nil {
			return msg.BatteryMsg{}
		}
		return msg.BatteryMsg(batt)
	}
}

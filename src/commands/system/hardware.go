package system

import (
	"os/exec"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/distatus/battery"
	"github.com/shirou/gopsutil/v3/host"

	"github.com/N1xev/bubbleMonitor/src/data"
	"github.com/N1xev/bubbleMonitor/src/messages"
)

// HostInfoCmd fetches host information
func HostInfoCmd() tea.Cmd {
	return func() tea.Msg {
		info, _ := host.Info()
		return messages.HostInfoMsg(info)
	}
}

// GpuInfoCmd fetches GPU information (NVIDIA only)
func GpuInfoCmd() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("nvidia-smi", "--query-gpu=name,driver_version,memory.total,memory.used", "--format=csv,noheader,nounits")
		out, err := cmd.Output()
		if err != nil {
			return nil
		}
		lines := strings.Split(string(out), "\n")
		var gpuList []data.GpuInfo
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
		return messages.GpuInfoMsg(gpuList)
	}
}

// TempCmd fetches temperature sensors
func TempCmd() tea.Cmd {
	return func() tea.Msg {
		temps, err := host.SensorsTemperatures()
		if err != nil {
			return messages.TempMsg{}
		}
		return messages.TempMsg(temps)
	}
}

// BatteryCmd fetches battery information
func BatteryCmd() tea.Cmd {
	return func() tea.Msg {
		batt, err := battery.GetAll()
		if err != nil {
			return messages.BatteryMsg{}
		}
		return messages.BatteryMsg(batt)
	}
}

package system

import (
	tea "charm.land/bubbletea/v2"
	"github.com/shirou/gopsutil/v3/disk"

	"github.com/N1xev/bubbleMonitor/internal/data"
	"github.com/N1xev/bubbleMonitor/internal/msg"
)

// DiskInfoCmd fetches disk partition information
func DiskInfoCmd() tea.Cmd {
	return func() tea.Msg {
		partitions, _ := disk.Partitions(false)
		var diskList []data.DiskPartition
		for _, p := range partitions {
			usage, err := disk.Usage(p.Mountpoint)
			if err != nil {
				continue
			}
			diskList = append(diskList, data.DiskPartition{
				Mountpoint: p.Mountpoint,
				Device:     p.Device,
				Fstype:     p.Fstype,
				Total:      usage.Total,
				Used:       usage.Used,
				UsedPct:    usage.UsedPercent,
			})
		}
		return msg.DiskInfoMsg{
			Partitions: diskList,
			Err:        nil,
		}
	}
}

// DiskIOCmd fetches disk I/O statistics
func DiskIOCmd() tea.Cmd {
	return func() tea.Msg {
		ioCounters, err := disk.IOCounters()
		if err != nil {
			return msg.DiskIOMsg{}
		}
		return msg.DiskIOMsg(ioCounters)
	}
}

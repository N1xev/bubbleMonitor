package system

import (
	tea "charm.land/bubbletea/v2"
	"github.com/shirou/gopsutil/v3/disk"

	"github.com/N1xev/bubbleMonitor/src/data"
	"github.com/N1xev/bubbleMonitor/src/messages"
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
		return messages.DiskInfoMsg(diskList)
	}
}

// DiskIOCmd fetches disk I/O statistics
func DiskIOCmd() tea.Cmd {
	return func() tea.Msg {
		ioCounters, err := disk.IOCounters()
		if err != nil {
			return messages.DiskIOMsg{}
		}
		return messages.DiskIOMsg(ioCounters)
	}
}

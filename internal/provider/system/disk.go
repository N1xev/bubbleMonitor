package system

import (
	"os/exec"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/shirou/gopsutil/v3/disk"

	"github.com/N1xev/bubbleMonitor/internal/data"
	"github.com/N1xev/bubbleMonitor/internal/msg"
)

// DiskInfoCmd fetches disk partition information
func DiskInfoCmd() tea.Cmd {
	return func() tea.Msg {
		partitions, err := disk.Partitions(false)
		var diskList []data.DiskPartition

		// First pass: add all partitions with mountpoints
		mountpointSeen := make(map[string]bool)
		for _, p := range partitions {
			if p.Mountpoint != "" {
				mountpointSeen[p.Mountpoint] = true
				usage, usageErr := disk.Usage(p.Mountpoint)
				if usageErr != nil {
					diskList = append(diskList, data.DiskPartition{
						Mountpoint: p.Mountpoint,
						Device:     p.Device,
						Fstype:     p.Fstype,
						Total:      0,
						Used:       0,
						UsedPct:    0,
					})
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
		}

		// Second pass: use lsblk to get unmounted partitions
		cmd := exec.Command("lsblk", "-b", "-o", "NAME,SIZE,TYPE,MOUNTPOINT", "-J")
		out, cmdErr := cmd.Output()
		if cmdErr == nil {
			parseLsblkUnmounted(string(out), &diskList)
		}

		return msg.DiskInfoMsg{
			Partitions: diskList,
			Err:        err,
		}
	}
}

func parseLsblkUnmounted(jsonOutput string, diskList *[]data.DiskPartition) {
	// Simple JSON parsing for lsblk output
	// Format: {"blockdevices":[{"name":"sda","size":...,"type":"disk","mountpoint":null,...},...]}
	lines := strings.Split(jsonOutput, "\n")
	var currentName string
	var currentSize uint64
	var currentType string
	var currentMountpoint string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, `"name":`) {
			currentName = strings.Trim(strings.TrimSpace(strings.Split(line, ":")[1]), `" ,`)
			currentSize = 0
			currentType = ""
		} else if strings.HasPrefix(line, `"size":`) {
			sizeStr := strings.Trim(strings.TrimSpace(strings.Split(line, ":")[1]), `" ,`)
			if sizeStr != "" && sizeStr != "0" {
				if size, err := strconv.ParseUint(sizeStr, 10, 64); err == nil {
					currentSize = size
				}
			}
		} else if strings.HasPrefix(line, `"type":`) {
			currentType = strings.Trim(strings.TrimSpace(strings.Split(line, ":")[1]), `" ,`)
		} else if strings.HasPrefix(line, `"mountpoint":`) {
			mp := strings.TrimSpace(strings.Split(line, ":")[1])
			if mp == "null" || mp == "" {
				currentMountpoint = ""
			} else {
				currentMountpoint = strings.Trim(mp, `"`)
			}

			// When we finish parsing a device entry (next name or end)
			// If it's a partition (type=part) with no mountpoint, add it
			if currentType == "part" && currentMountpoint == "" && currentName != "" && currentSize > 0 {
				devicePath := "/dev/" + currentName
				// Check if already added
				alreadySeen := false
				for _, d := range *diskList {
					if d.Device == devicePath {
						alreadySeen = true
						break
					}
				}
				if !alreadySeen {
					*diskList = append(*diskList, data.DiskPartition{
						Mountpoint: "",
						Device:     devicePath,
						Fstype:     "",
						Total:      currentSize,
						Used:       0,
						UsedPct:    0,
					})
				}
			}
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

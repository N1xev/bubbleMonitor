package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/cobra"

	"github.com/N1xev/bubbleMonitor/internal/util"
)

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show system status overview",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return printStatus(cmd)
		},
	}
}

func printStatus(cmd *cobra.Command) error {
	// CPU
	cpuPcts, _ := cpu.Percent(time.Second, false)
	cpuVal := 0.0
	if len(cpuPcts) > 0 {
		cpuVal = cpuPcts[0]
	}

	// Memory
	memInfo, _ := mem.VirtualMemory()
	memVal := 0.0
	memTotal := uint64(0)
	if memInfo != nil {
		memVal = memInfo.UsedPercent
		memTotal = memInfo.Total
	}

	// Swap
	swapVal := 0.0
	swapInfo, _ := mem.SwapMemory()
	if swapInfo != nil {
		swapVal = swapInfo.UsedPercent
	}

	// Disk
	diskVal := 0.0
	diskUsage, _ := disk.Usage("/")
	if diskUsage != nil {
		diskVal = diskUsage.UsedPercent
	}

	// Temperature
	tempVal := 0.0
	hostInfo, _ := host.Info()
	if temps, err := host.SensorsTemperatures(); err == nil {
		for _, t := range temps {
			if t.Temperature > tempVal {
				tempVal = t.Temperature
			}
		}
	}

	// Uptime
	uptime := "unknown"
	if hostInfo != nil {
		uptime = formatUptime(hostInfo.Uptime)
	}

	barWidth := 12
	out := cmd.OutOrStdout()

	fmt.Fprintf(out, "\n")
	fmt.Fprintf(out, "  CPU     %s  %6s   | Memory  %s  %6s\n",
		renderBar(cpuVal, barWidth), util.FastPercent1(cpuVal),
		renderBar(memVal, barWidth), util.FastPercent1(memVal))
	fmt.Fprintf(out, "  Disk    %s  %6s   | Temp    %s  %6s\n",
		renderBar(diskVal, barWidth), util.FastPercent1(diskVal),
		renderTempBar(tempVal, barWidth), fmt.Sprintf("%.0f°C", tempVal))
	fmt.Fprintf(out, "  Swap    %s  %6s   | Uptime  %s\n",
		renderBar(swapVal, barWidth), util.FastPercent1(swapVal),
		uptime)

	if memTotal > 0 {
		fmt.Fprintf(out, "\n  Memory: %s / %s\n",
			util.FastBytes(memInfo.Used), util.FastBytes(memTotal))
	}
	fmt.Fprintf(out, "\n")

	return nil
}

func renderBar(pct float64, width int) string {
	filled := int(pct / 100.0 * float64(width))
	if filled > width {
		filled = width
	}
	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	return bar
}

func renderTempBar(temp float64, width int) string {
	// Assume 100°C is max
	pct := temp / 100.0 * 100.0
	return renderBar(pct, width)
}

func formatUptime(seconds uint64) string {
	d := seconds / 86400
	h := (seconds % 86400) / 3600
	m := (seconds % 3600) / 60
	if d > 0 {
		return fmt.Sprintf("%dd %dh %dm", d, h, m)
	}
	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}

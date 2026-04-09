package main

import (
	"fmt"
	"os"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/cobra"

	configpkg "github.com/N1xev/bubbleMonitor/internal/config"
)

func newHealthCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "Show system health score",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return printHealth(cmd)
		},
	}
}

func printHealth(cmd *cobra.Command) error {
	cfg, _ := loadConfigWithOverrides()
	thresholds := cfg.Thresholds
	weights := cfg.HealthWeights

	cpuVal := fetchCPUPct()
	memVal := fetchMemPct()
	diskVal := fetchDiskPct()
	tempVal := fetchTempPct()

	score := 100
	out := cmd.OutOrStdout()

	fmt.Fprintf(out, "\n")

	// CPU
	cpuStatus := "OK"
	if cpuVal >= thresholds[configpkg.MetricCPU] {
		cpuStatus = "CRITICAL"
		score -= weights.CpuCritical
	} else if cpuVal >= thresholds[configpkg.MetricCPU]*0.7 {
		cpuStatus = "HIGH"
		score -= weights.CpuHigh
	}
	fmt.Fprintf(out, "  CPU:        %s (%.1f%%)\n", cpuStatus, cpuVal)

	// Memory
	memStatus := "OK"
	if memVal >= thresholds[configpkg.MetricMem] {
		memStatus = "CRITICAL"
		score -= weights.MemCritical
	} else if memVal >= thresholds[configpkg.MetricMem]*0.7 {
		memStatus = "HIGH"
		score -= weights.MemHigh
	}
	fmt.Fprintf(out, "  Memory:     %s (%.1f%%)\n", memStatus, memVal)

	// Disk
	diskStatus := "OK"
	if diskVal >= thresholds[configpkg.MetricDisk] {
		diskStatus = "CRITICAL"
		score -= weights.DiskCritical
	}
	fmt.Fprintf(out, "  Disk:       %s (%.1f%%)\n", diskStatus, diskVal)

	// Temperature
	tempStatus := "OK"
	if tempVal >= thresholds[configpkg.MetricTemp] {
		tempStatus = "CRITICAL"
		score -= weights.TempCritical
	} else if tempVal >= thresholds[configpkg.MetricTemp]*0.8 {
		tempStatus = "HIGH"
		score -= weights.TempHigh
	}
	fmt.Fprintf(out, "  Temp:       %s (%.0f°C)\n", tempStatus, tempVal)

	if score < 0 {
		score = 0
	}

	fmt.Fprintf(out, "\n  Health Score: %d/100\n\n", score)

	if score < 50 {
		os.Exit(2)
	} else if score < 70 {
		os.Exit(1)
	}
	return nil
}

func fetchCPUPct() float64 {
	pcts, _ := cpu.Percent(0, false)
	if len(pcts) == 0 {
		return 0
	}
	return pcts[0]
}

func fetchMemPct() float64 {
	m, _ := mem.VirtualMemory()
	if m == nil {
		return 0
	}
	return m.UsedPercent
}

func fetchDiskPct() float64 {
	d, _ := disk.Usage("/")
	if d == nil {
		return 0
	}
	return d.UsedPercent
}

func fetchTempPct() float64 {
	maxTemp := 0.0
	temps, _ := host.SensorsTemperatures()
	for _, t := range temps {
		if t.Temperature > maxTemp {
			maxTemp = t.Temperature
		}
	}
	return maxTemp
}

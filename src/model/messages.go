package model

import (
	"time"

	"github.com/shirou/gopsutil/v3/load"
)

// Message types for Bubble Tea
type TickMsg time.Time

type MetricsMsg struct {
	Cpu        float64
	CpuPerCore []float64
	Memory     float64
	Disk       float64
	Swap       float64
	NetSent    uint64
	NetRecv    uint64
	LoadAvg    *load.AvgStat
}

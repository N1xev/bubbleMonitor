package data

import (
	"fmt"
	"time"

	"github.com/N1xev/bubbleMonitor/src/config"
)

// Alert represents a triggered alert
type Alert struct {
	Type      config.MetricType
	Value     float64
	Threshold float64
	Message   string
	Timestamp time.Time
}

// AlertManager handles checking and storing alerts
type AlertManager struct {
	ActiveAlerts map[config.MetricType]Alert
}

// NewAlertManager creates a new alert manager
func NewAlertManager() *AlertManager {
	return &AlertManager{
		ActiveAlerts: make(map[config.MetricType]Alert),
	}
}

// CheckAlerts verifies metrics against thresholds and updates active alerts
// Uses AppState to check metrics
func (am *AlertManager) CheckAlerts(s *AppState) {
	// CPU Check
	cpuThreshold := s.Config.Thresholds[config.MetricCPU]
	if cpuThreshold > 0 && s.Cpu > cpuThreshold {
		am.ActiveAlerts[config.MetricCPU] = Alert{
			Type:      config.MetricCPU,
			Value:     s.Cpu,
			Threshold: cpuThreshold,
			Message:   fmt.Sprintf("CPU Usage High: %.1f%% (>%.0f%%)", s.Cpu, cpuThreshold),
			Timestamp: time.Now(),
		}
	} else {
		delete(am.ActiveAlerts, config.MetricCPU)
	}

	// Memory Check
	memThreshold := s.Config.Thresholds[config.MetricMem]
	if memThreshold > 0 && s.Memory > memThreshold {
		am.ActiveAlerts[config.MetricMem] = Alert{
			Type:      config.MetricMem,
			Value:     s.Memory,
			Threshold: memThreshold,
			Message:   fmt.Sprintf("Memory Usage High: %.1f%% (>%.0f%%)", s.Memory, memThreshold),
			Timestamp: time.Now(),
		}
	} else {
		delete(am.ActiveAlerts, config.MetricMem)
	}

	// Disk Check (Overall System Usage)
	diskThreshold := s.Config.Thresholds[config.MetricDisk]
	if diskThreshold > 0 {
		var totalUsed, totalAll uint64
		for _, part := range s.DiskPartitions {
			totalUsed += part.Used
			totalAll += part.Total
		}

		if totalAll > 0 {
			overallPct := (float64(totalUsed) / float64(totalAll)) * 100.0

			if overallPct > diskThreshold {
				am.ActiveAlerts[config.MetricDisk] = Alert{
					Type:      config.MetricDisk,
					Value:     overallPct,
					Threshold: diskThreshold,
					Message:   fmt.Sprintf("Overall Disk High: %.1f%% (>%.0f%%)", overallPct, diskThreshold),
					Timestamp: time.Now(),
				}
			} else {
				delete(am.ActiveAlerts, config.MetricDisk)
			}
		}
	}

	// Temperature Check
	tempThreshold := s.Config.Thresholds[config.MetricTemp]
	if tempThreshold > 0 && s.CpuTemp > tempThreshold {
		am.ActiveAlerts[config.MetricTemp] = Alert{
			Type:      config.MetricTemp,
			Value:     s.CpuTemp,
			Threshold: tempThreshold,
			Message:   fmt.Sprintf("CPU Temp High: %.1f°C (>%.0f°C)", s.CpuTemp, tempThreshold),
			Timestamp: time.Now(),
		}
	} else {
		delete(am.ActiveAlerts, config.MetricTemp)
	}
}

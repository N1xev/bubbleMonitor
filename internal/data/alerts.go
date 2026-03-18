package data

import (
	"fmt"
	"sync"
	"time"

	"github.com/N1xev/bubbleMonitor/internal/config"
)

type Alert struct {
	Timestamp time.Time
	Message   string
	Type      config.MetricType
	Value     float64
	Threshold float64
}

type AlertManager struct {
	ActiveAlerts map[config.MetricType]Alert
	mu           sync.RWMutex
}

func NewAlertManager() *AlertManager {
	return &AlertManager{
		ActiveAlerts: make(map[config.MetricType]Alert),
	}
}

func (am *AlertManager) CheckAlerts(s *AppState) {
	am.mu.Lock()
	defer am.mu.Unlock()

	cpuThreshold := s.Config.Config.Thresholds[config.MetricCPU]
	if cpuThreshold > 0 && s.Metrics.Cpu > cpuThreshold {
		am.ActiveAlerts[config.MetricCPU] = Alert{
			Type:      config.MetricCPU,
			Value:     s.Metrics.Cpu,
			Threshold: cpuThreshold,
			Message:   fmt.Sprintf("CPU Usage High: %.1f%% > (%.0f%%)", s.Metrics.Cpu, cpuThreshold),
			Timestamp: time.Now(),
		}
	} else {
		delete(am.ActiveAlerts, config.MetricCPU)
	}

	memThreshold := s.Config.Config.Thresholds[config.MetricMem]
	if memThreshold > 0 && s.Metrics.Memory > memThreshold {
		am.ActiveAlerts[config.MetricMem] = Alert{
			Type:      config.MetricMem,
			Value:     s.Metrics.Memory,
			Threshold: memThreshold,
			Message:   fmt.Sprintf("Memory Usage High: %.1f%% > (%.0f%%)", s.Metrics.Memory, memThreshold),
			Timestamp: time.Now(),
		}
	} else {
		delete(am.ActiveAlerts, config.MetricMem)
	}

	diskThreshold := s.Config.Config.Thresholds[config.MetricDisk]
	if diskThreshold > 0 {
		var totalUsed, totalAll uint64
		for _, part := range s.Metrics.DiskPartitions {
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
					Message:   fmt.Sprintf("Overall Disk High: %.1f%% > (%.0f%%)", overallPct, diskThreshold),
					Timestamp: time.Now(),
				}
			} else {
				delete(am.ActiveAlerts, config.MetricDisk)
			}
		}
	}

	tempThreshold := s.Config.Config.Thresholds[config.MetricTemp]
	if tempThreshold > 0 && s.Metrics.CpuTemp > tempThreshold {
		am.ActiveAlerts[config.MetricTemp] = Alert{
			Type:      config.MetricTemp,
			Value:     s.Metrics.CpuTemp,
			Threshold: tempThreshold,
			Message:   fmt.Sprintf("CPU Temp High: %.1f°C > (%.0f°C)", s.Metrics.CpuTemp, tempThreshold),
			Timestamp: time.Now(),
		}
	} else {
		delete(am.ActiveAlerts, config.MetricTemp)
	}
}

func (am *AlertManager) GetAlerts() []Alert {
	am.mu.RLock()
	defer am.mu.RUnlock()
	alerts := make([]Alert, 0, len(am.ActiveAlerts))
	for _, alert := range am.ActiveAlerts {
		alerts = append(alerts, alert)
	}
	return alerts
}

func (am *AlertManager) HasAlerts() bool {
	am.mu.RLock()
	defer am.mu.RUnlock()
	return len(am.ActiveAlerts) > 0
}

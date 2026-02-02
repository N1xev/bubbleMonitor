package app

import (
	configpkg "github.com/N1xev/bubbleMonitor/internal/config"
	"github.com/N1xev/bubbleMonitor/internal/data"
)

// UpdateHealthScore calculates the system health score based on critical metrics.
// We start with 100 points and deduct points based on severity of usage (CPU, Memory, Disk, Temp).
func UpdateHealthScore(s *data.AppState) {
	score := data.MaxHealthScore
	weights := s.Config.HealthWeights
	thresholds := s.Config.Thresholds

	// CPU penalty
	if s.Cpu > thresholds[configpkg.MetricCPU] {
		score -= weights.CpuCritical
	} else if s.Cpu > (thresholds[configpkg.MetricCPU] * 0.7) { // 70% of critical is warning
		score -= weights.CpuHigh
	}

	// Memory penalty
	if s.Memory > thresholds[configpkg.MetricMem] {
		score -= weights.MemCritical
	} else if s.Memory > (thresholds[configpkg.MetricMem] * 0.7) {
		score -= weights.MemHigh
	}

	// Disk space penalty
	if s.Disk > thresholds[configpkg.MetricDisk] {
		score -= weights.DiskCritical
	}

	// Temperature penalty
	if s.CpuTemp > thresholds[configpkg.MetricTemp] {
		score -= weights.TempCritical
	} else if s.CpuTemp > (thresholds[configpkg.MetricTemp] * 0.8) {
		score -= weights.TempHigh
	}

	if score < 0 {
		score = 0
	}
	s.HealthScore = score
}

// UpdateProcessHistory maintains a history buffer for the selected process
// and top resource consumers. This allows the charts to show history even
// if you just clicked on a process.
func UpdateProcessHistory(s *data.AppState, alivePids map[int32]bool) {
	if len(s.Processes) == 0 {
		return
	}

	// We want to track history for:
	// 1. The currently selected process (obviously)
	// 2. The top N processes by CPU, so if we switch to them, we have data ready
	targetPids := make(map[int32]bool)

	if s.SelectedProcess >= 0 && s.SelectedProcess < len(s.Processes) {
		pid := s.Processes[s.SelectedProcess].Pid
		targetPids[pid] = true
	}

	// Pre-warm history for top consumers
	limit := data.TopProcessesTrackCount
	if len(s.Processes) < limit {
		limit = len(s.Processes)
	}
	for i := 0; i < limit; i++ {
		targetPids[s.Processes[i].Pid] = true
	}

	for _, p := range s.Processes {
		if targetPids[p.Pid] {
			hist := s.GetOrCreateHistory(p.Pid)
			hist.Push(p.Cpu)
		}
	}

	// Cleanup old data to prevent memory leaks
	s.PruneDeadProcessHistory(alivePids)
}

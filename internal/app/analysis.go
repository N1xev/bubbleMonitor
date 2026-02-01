package app

import (
	"github.com/N1xev/bubbleMonitor/internal/data"
)

func UpdateAnalysis(s *data.AppState, alivePids map[int32]bool) {
	// 1. Calculate Health Score
	score := data.MaxHealthScore
	if s.Cpu > data.HealthThresholdHealthy {
		score -= data.HealthDeductionCPUCritical
	} else if s.Cpu > data.HealthThresholdWarning {
		score -= data.HealthDeductionCPUHigh
	}

	if s.Memory > data.HealthThresholdCritical {
		score -= data.HealthDeductionMemoryCritical
	} else if s.Memory > data.HealthThresholdWarning {
		score -= data.HealthDeductionMemoryHigh
	}

	if s.Disk > data.HealthThresholdCritical {
		score -= data.HealthDeductionDiskCritical
	}

	if s.CpuTemp > data.HealthThresholdHealthy {
		score -= data.HealthDeductionTempCritical
	} else if s.CpuTemp > data.HealthThresholdWarning {
		score -= data.HealthDeductionTempHigh
	}

	if score < 0 {
		score = 0
	}
	s.HealthScore = score

	// 2. Track Process History (Top 5 + Selected)
	// Only if we have processes
	if len(s.Processes) == 0 {
		return
	}

	targetPids := make(map[int32]bool)

	// Track Selected
	if s.SelectedProcess >= 0 && s.SelectedProcess < len(s.Processes) {
		pid := s.Processes[s.SelectedProcess].Pid
		targetPids[pid] = true
	}

	// Track Top 5 (Assuming list is sorted by CPU, which is default)
	// If sorted by Mem, we track top Mem consumers.
	limit := data.TopProcessesTrackCount
	if len(s.Processes) < limit {
		limit = len(s.Processes)
	}
	for i := 0; i < limit; i++ {
		targetPids[s.Processes[i].Pid] = true
	}

	// Update History for targets
	for _, p := range s.Processes {
		if targetPids[p.Pid] {
			hist := s.GetOrCreateHistory(p.Pid)
			hist.Push(p.Cpu)
		}
	}

	// Prune untracked history
	if alivePids == nil {
		alivePids = make(map[int32]bool)
		for _, p := range s.Processes {
			alivePids[p.Pid] = true
		}
	}
	s.PruneDeadProcessHistory(alivePids)
}

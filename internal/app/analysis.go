package app

import (
	"github.com/N1xev/bubbleMonitor/internal/data"
)

func UpdateHealthScore(s *data.AppState) {
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
}

func UpdateProcessHistory(s *data.AppState, alivePids map[int32]bool) {
	if len(s.Processes) == 0 {
		return
	}

	targetPids := make(map[int32]bool)

	if s.SelectedProcess >= 0 && s.SelectedProcess < len(s.Processes) {
		pid := s.Processes[s.SelectedProcess].Pid
		targetPids[pid] = true
	}

	// Track top consumers (CPU by default)
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

	s.PruneDeadProcessHistory(alivePids)
}

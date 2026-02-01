package app

import (
	"github.com/N1xev/bubbleMonitor/internal/data"
)

func UpdateAnalysis(s *data.AppState, alivePids map[int32]bool) {
	// 1. Calculate Health Score
	score := 100
	if s.Cpu > 90 {
		score -= 30
	} else if s.Cpu > 70 {
		score -= 10
	}

	if s.Memory > 95 {
		score -= 30
	} else if s.Memory > 80 {
		score -= 10
	}

	if s.Disk > 95 {
		score -= 20
	}

	if s.CpuTemp > 90 {
		score -= 30
	} else if s.CpuTemp > 80 {
		score -= 10
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
	limit := 5
	if len(s.Processes) < 5 {
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

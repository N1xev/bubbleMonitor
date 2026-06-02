package data

import (
	"sync"
)

type AppState struct {
	Metrics MetricsState
	Process ProcessState
	UI      UIState
	Config  ConfigState
	Remote  RemoteState

	stateMu sync.RWMutex
}

func (s *AppState) SetSuspended(pid int32, suspended bool) {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	if s.Process.SuspendedState == nil {
		s.Process.SuspendedState = make(map[int32]bool)
	}
	if suspended {
		s.Process.SuspendedState[pid] = true
	} else {
		delete(s.Process.SuspendedState, pid)
	}
}

func (s *AppState) IsSuspended(pid int32) bool {
	s.stateMu.RLock()
	defer s.stateMu.RUnlock()
	return s.Process.SuspendedState[pid]
}

func (s *AppState) ToggleCollapsed(pid int32) {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	if s.Process.CollapsedPids == nil {
		s.Process.CollapsedPids = make(map[int32]bool)
	}
	s.Process.CollapsedPids[pid] = !s.Process.CollapsedPids[pid]
}

func (s *AppState) IsCollapsed(pid int32) bool {
	s.stateMu.RLock()
	defer s.stateMu.RUnlock()
	return s.Process.CollapsedPids[pid]
}

func (s *AppState) ToggleBookmark(pid int32) {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	if s.Process.BookmarkedPids == nil {
		s.Process.BookmarkedPids = make(map[int32]bool)
	}
	s.Process.BookmarkedPids[pid] = !s.Process.BookmarkedPids[pid]
}

func (s *AppState) IsBookmarked(pid int32) bool {
	s.stateMu.RLock()
	defer s.stateMu.RUnlock()
	return s.Process.BookmarkedPids[pid]
}

func (s *AppState) GetProcessByPid(pid int32) (ProcessInfo, bool) {
	s.stateMu.RLock()
	defer s.stateMu.RUnlock()
	p, ok := s.Process.ProcessesByPid[pid]
	return p, ok
}

func (s *AppState) SyncProcessesMap() {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	s.Process.ProcessesByPid = make(map[int32]ProcessInfo, len(s.Process.Processes))
	for _, p := range s.Process.Processes {
		s.Process.ProcessesByPid[p.Pid] = p
	}
}

func (s *AppState) ClearProcessMaps() {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	for pid := range s.Process.SuspendedState {
		delete(s.Process.SuspendedState, pid)
	}
	for pid := range s.Process.CollapsedPids {
		delete(s.Process.CollapsedPids, pid)
	}
	for pid := range s.Process.BookmarkedPids {
		delete(s.Process.BookmarkedPids, pid)
	}
}

func (s *AppState) PruneDeadProcessMaps(alivePids map[int32]bool) {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	for pid := range s.Process.SuspendedState {
		if !alivePids[pid] {
			delete(s.Process.SuspendedState, pid)
		}
	}
	for pid := range s.Process.CollapsedPids {
		if !alivePids[pid] {
			delete(s.Process.CollapsedPids, pid)
		}
	}
	for pid := range s.Process.BookmarkedPids {
		if !alivePids[pid] {
			delete(s.Process.BookmarkedPids, pid)
		}
	}
}

func (s *AppState) GetOrCreateHistory(pid int32) *RingBuffer {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	if s.Metrics.ProcessHistory == nil {
		s.Metrics.ProcessHistory = make(map[int32]*RingBuffer)
	}
	hist, ok := s.Metrics.ProcessHistory[pid]
	if !ok {
		hist = NewRingBuffer(s.Config.HistoryLength)
		s.Metrics.ProcessHistory[pid] = hist
	}
	return hist
}

func (s *AppState) GetHistory(pid int32) (*RingBuffer, bool) {
	s.stateMu.RLock()
	defer s.stateMu.RUnlock()
	hist, ok := s.Metrics.ProcessHistory[pid]
	// The returned *RingBuffer shares state with the rest of the app.
	// The caller must not retain the pointer past a tick; concurrent
	// PruneDeadProcessHistory may delete the map entry while the
	// caller still uses the buffer.
	return hist, ok
}

func (s *AppState) PruneDeadProcessHistory(alivePids map[int32]bool) {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	for pid := range s.Metrics.ProcessHistory {
		if !alivePids[pid] {
			delete(s.Metrics.ProcessHistory, pid)
		}
	}
}

package data

import (
	"sync"
	"testing"
)

func TestAppStateConcurrentMapAccess(t *testing.T) {
	state := &AppState{
		SuspendedState: make(map[int32]bool),
		CollapsedPids:  make(map[int32]bool),
		BookmarkedPids: make(map[int32]bool),
		ProcessHistory: make(map[int32]*RingBuffer),
		HistoryLength:  100,
	}

	const numGoroutines = 10
	const iterations = 1000
	var wg sync.WaitGroup

	wg.Add(numGoroutines * 4)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			pid := int32(id % 5)
			for j := 0; j < iterations; j++ {
				state.SetSuspended(pid, j%2 == 0)
			}
		}(i)

		go func(id int) {
			defer wg.Done()
			pid := int32(id % 5)
			for j := 0; j < iterations; j++ {
				state.ToggleCollapsed(pid)
			}
		}(i)

		go func(id int) {
			defer wg.Done()
			pid := int32(id % 5)
			for j := 0; j < iterations; j++ {
				state.ToggleBookmark(pid)
			}
		}(i)

		go func(id int) {
			defer wg.Done()
			pid := int32(id % 5)
			for j := 0; j < iterations; j++ {
				_ = state.GetOrCreateHistory(pid)
			}
		}(i)
	}

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer func() { wg.Done() }()
			pid := int32(id % 5)
			for j := 0; j < iterations; j++ {
				_ = state.IsSuspended(pid)
			}
		}(i)
		wg.Add(1)

		go func(id int) {
			defer func() { wg.Done() }()
			pid := int32(id % 5)
			for j := 0; j < iterations; j++ {
				_ = state.IsCollapsed(pid)
			}
		}(i)
		wg.Add(1)

		go func(id int) {
			defer func() { wg.Done() }()
			pid := int32(id % 5)
			for j := 0; j < iterations; j++ {
				_ = state.IsBookmarked(pid)
			}
		}(i)
		wg.Add(1)

		go func(id int) {
			defer func() { wg.Done() }()
			pid := int32(id % 5)
			for j := 0; j < iterations; j++ {
				_, _ = state.GetHistory(pid)
			}
		}(i)
		wg.Add(1)
	}

	wg.Wait()
}

func TestAppStatePruneDeadProcessMaps(t *testing.T) {
	state := &AppState{
		SuspendedState: make(map[int32]bool),
		CollapsedPids:  make(map[int32]bool),
		BookmarkedPids: make(map[int32]bool),
		ProcessHistory: make(map[int32]*RingBuffer),
		HistoryLength:  100,
	}

	state.SetSuspended(100, true)
	state.SetSuspended(200, true)
	state.SetSuspended(300, true)
	state.ToggleCollapsed(100)
	state.ToggleCollapsed(400)
	state.ToggleBookmark(200)
	state.ToggleBookmark(500)

	alivePids := map[int32]bool{
		100: true,
		200: true,
	}

	state.PruneDeadProcessMaps(alivePids)

	if state.IsSuspended(100) != true {
		t.Error("Expected pid 100 to still be suspended")
	}
	if state.IsSuspended(200) != true {
		t.Error("Expected pid 200 to still be suspended")
	}
	if state.IsSuspended(300) != false {
		t.Error("Expected pid 300 to be removed from suspended")
	}

	if state.IsCollapsed(100) != true {
		t.Error("Expected pid 100 to still be collapsed")
	}
	if state.IsCollapsed(400) != false {
		t.Error("Expected pid 400 to be removed from collapsed")
	}

	if state.IsBookmarked(200) != true {
		t.Error("Expected pid 200 to still be bookmarked")
	}
	if state.IsBookmarked(500) != false {
		t.Error("Expected pid 500 to be removed from bookmarked")
	}
}

func TestAppStatePruneDeadProcessHistory(t *testing.T) {
	state := &AppState{
		ProcessHistory: make(map[int32]*RingBuffer),
		HistoryLength:  100,
	}

	state.GetOrCreateHistory(100)
	state.GetOrCreateHistory(200)
	state.GetOrCreateHistory(300)

	alivePids := map[int32]bool{
		100: true,
		200: true,
	}

	state.PruneDeadProcessHistory(alivePids)

	_, ok := state.GetHistory(100)
	if !ok {
		t.Error("Expected pid 100 history to still exist")
	}

	_, ok = state.GetHistory(200)
	if !ok {
		t.Error("Expected pid 200 history to still exist")
	}

	_, ok = state.GetHistory(300)
	if ok {
		t.Error("Expected pid 300 history to be removed")
	}
}

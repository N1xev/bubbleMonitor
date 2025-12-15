package data

import (
	"sort"
	"strings"
)

// GetVisibleProcesses returns the list of processes that should be currently displayed
// respecting filtering, tree view, and collapsed states.
// It returns the flat list of visible processes and a map of indentation levels.
func (s *AppState) GetVisibleProcesses() ([]ProcessInfo, map[int32]int) {
	if !s.TreeView {
		// Normal view: just filtered list
		return s.GetFilteredProcesses(), make(map[int32]int)
	}

	// Tree View: Build tree respecting collapsed state
	procs := s.GetFilteredProcesses()
	return s.buildProcessTree(procs)
}

// buildProcessTree organizes processes into a tree structure
func (s *AppState) buildProcessTree(procs []ProcessInfo) ([]ProcessInfo, map[int32]int) {
	// Create PID map for quick lookup and index map to preserve order
	procMap := make(map[int32]*ProcessInfo)
	procIdx := make(map[int32]int)
	children := make(map[int32][]int32)
	for i := range procs {
		procMap[procs[i].Pid] = &procs[i]
		procIdx[procs[i].Pid] = i
		children[procs[i].Ppid] = append(children[procs[i].Ppid], procs[i].Pid)
	}

	var flatList []ProcessInfo
	indentMap := make(map[int32]int)

	// Find roots (parent not in valid list)
	var roots []int32
	for _, p := range procs {
		if _, exists := procMap[p.Ppid]; !exists {
			roots = append(roots, p.Pid)
		}
	}

	// Sort roots by their original index to preserve sort order
	sort.Slice(roots, func(i, j int) bool { return procIdx[roots[i]] < procIdx[roots[j]] })

	var build func(pid int32, level int)
	build = func(pid int32, level int) {
		if p, ok := procMap[pid]; ok {
			flatList = append(flatList, *p)
			indentMap[pid] = level

			// Check if collapsed
			if s.CollapsedPids[pid] {
				return
			}

			kids := children[pid]
			// Sort kids by their original index
			sort.Slice(kids, func(i, j int) bool { return procIdx[kids[i]] < procIdx[kids[j]] })

			for _, kid := range kids {
				build(kid, level+1)
			}
		}
	}

	for _, rootPid := range roots {
		build(rootPid, 0)
	}

	return flatList, indentMap
}

// GetFilteredProcesses returns processes matching the current filter
func (s *AppState) GetFilteredProcesses() []ProcessInfo {
	if s.ProcessFilter == "" {
		return s.Processes
	}
	var filtered []ProcessInfo
	filterLower := strings.ToLower(s.ProcessFilter)
	for _, p := range s.Processes {
		if strings.Contains(strings.ToLower(p.Name), filterLower) {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

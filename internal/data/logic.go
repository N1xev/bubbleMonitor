package data

import (
	"sort"
	"strconv"
	"strings"
)

type treeBuilder struct {
	procMap  map[int32]*ProcessInfo
	procIdx  map[int32]int
	children map[int32][]int32
}

func (s *AppState) GetVisibleProcesses() ([]ProcessInfo, map[int32]int) {
	if !s.Process.ProcessCacheDirty && s.Process.CachedVisibleProcs != nil {
		return s.Process.CachedVisibleProcs, s.Process.CachedTreeIndents
	}

	if !s.Process.TreeView {
		filtered := s.GetFilteredProcesses()
		s.Process.CachedVisibleProcs = filtered
		s.Process.CachedTreeIndents = make(map[int32]int)
		s.Process.ProcessCacheDirty = false
		return filtered, s.Process.CachedTreeIndents
	}

	procs := s.GetFilteredProcesses()
	visible, indents := s.buildProcessTree(procs)
	s.Process.CachedVisibleProcs = visible
	s.Process.CachedTreeIndents = indents
	s.Process.ProcessCacheDirty = false
	return visible, indents
}

func (s *AppState) InvalidateProcessCache() {
	s.Process.ProcessCacheDirty = true
}

func (s *AppState) buildProcessTree(procs []ProcessInfo) ([]ProcessInfo, map[int32]int) {
	tb := &treeBuilder{
		procMap:  make(map[int32]*ProcessInfo, 500),
		procIdx:  make(map[int32]int, 500),
		children: make(map[int32][]int32, 500),
	}

	for i := range procs {
		tb.procMap[procs[i].Pid] = &procs[i]
		tb.procIdx[procs[i].Pid] = i
		tb.children[procs[i].Ppid] = append(tb.children[procs[i].Ppid], procs[i].Pid)
	}

	flatList := make([]ProcessInfo, 0, len(procs))
	indentMap := make(map[int32]int)

	roots := make([]int32, 0, 10)
	for _, p := range procs {
		if _, exists := tb.procMap[p.Ppid]; !exists {
			roots = append(roots, p.Pid)
		}
	}

	sort.Slice(roots, func(i, j int) bool { return tb.procIdx[roots[i]] < tb.procIdx[roots[j]] })

	var build func(pid int32, level int)
	build = func(pid int32, level int) {
		if p, ok := tb.procMap[pid]; ok {
			flatList = append(flatList, *p)
			indentMap[pid] = level

			if s.IsCollapsed(pid) {
				return
			}

			kids := tb.children[pid]
			sort.Slice(kids, func(i, j int) bool { return tb.procIdx[kids[i]] < tb.procIdx[kids[j]] })

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

func (s *AppState) GetFilteredProcesses() []ProcessInfo {
	if s.Process.ProcessFilter == "" {
		return s.Process.Processes
	}
	filterLower := s.Process.ProcessFilterLower
	if filterLower == "" {
		filterLower = strings.ToLower(s.Process.ProcessFilter)
	}
	capacity := len(s.Process.Processes) / 2
	if capacity < 4 {
		capacity = 4
	}
	filtered := make([]ProcessInfo, 0, capacity)
	for i := range s.Process.Processes {
		p := &s.Process.Processes[i]
		if strings.Contains(p.NameLower, filterLower) ||
			strings.Contains(p.UsernameLower, filterLower) ||
			strings.Contains(p.CmdlineLower, filterLower) {
			filtered = append(filtered, *p)
		} else if filterLower != "" {
			pidStr := strconv.Itoa(int(p.Pid))
			if strings.Contains(pidStr, filterLower) {
				filtered = append(filtered, *p)
			}
		}
	}
	return filtered
}

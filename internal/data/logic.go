package data

import (
	"fmt"
	"sort"
	"strings"
)

type treeBuilder struct {
	procMap  map[int32]*ProcessInfo
	procIdx  map[int32]int
	children map[int32][]int32
}

func newTreeBuilder() *treeBuilder {
	return &treeBuilder{
		procMap:  make(map[int32]*ProcessInfo, 500),
		procIdx:  make(map[int32]int, 500),
		children: make(map[int32][]int32, 500),
	}
}

func (s *AppState) GetVisibleProcesses() ([]ProcessInfo, map[int32]int) {
	if !s.ProcessCacheDirty && s.CachedVisibleProcs != nil {
		return s.CachedVisibleProcs, s.CachedTreeIndents
	}

	if !s.TreeView {
		filtered := s.GetFilteredProcesses()
		s.CachedVisibleProcs = filtered
		s.CachedTreeIndents = make(map[int32]int)
		s.ProcessCacheDirty = false
		return filtered, s.CachedTreeIndents
	}

	procs := s.GetFilteredProcesses()
	visible, indents := s.buildProcessTree(procs)
	s.CachedVisibleProcs = visible
	s.CachedTreeIndents = indents
	s.ProcessCacheDirty = false
	return visible, indents
}

func (s *AppState) InvalidateProcessCache() {
	s.ProcessCacheDirty = true
}

func (s *AppState) buildProcessTree(procs []ProcessInfo) ([]ProcessInfo, map[int32]int) {
	tb := newTreeBuilder()

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
	if s.ProcessFilter == "" {
		return s.Processes
	}
	var filtered []ProcessInfo
	filterLower := strings.ToLower(s.ProcessFilter)
	for _, p := range s.Processes {
		if strings.Contains(p.NameLower, filterLower) ||
			strings.Contains(p.UsernameLower, filterLower) ||
			strings.Contains(p.CmdlineLower, filterLower) ||
			strings.Contains(fmt.Sprintf("%d", p.Pid), filterLower) {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

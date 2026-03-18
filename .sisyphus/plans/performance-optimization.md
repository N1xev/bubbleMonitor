# Performance Optimization for bubbleMonitor

## TL;DR

> **Quick Summary**: Optimize CPU usage from 2.2% to <0.5% and reduce memory allocations by caching repeated computations, eliminating redundant string operations, and avoiding unnecessary recreations.
> 
> **Deliverables**:
> - Cache CmdlineLower in CachedProcessInfo (avoid repeated ToLower)
> - Cache SortBy uppercase in AppState (avoid ToUpper every render)
> - Cache filterLower in AppState (avoid ToLower on every filter check)
> - Optimize HealthScore calculation to only run when metrics change
> - Fix ZoneManager recreation every render
> 
> **Estimated Effort**: Short (5 tasks)
> **Parallel Execution**: YES - tasks are independent
> **Critical Path**: Task 1 → All others can run in parallel

---

## Context

### Original Request
User wants to reduce CPU from 2.2% to <0.5% and mem from 0.2% to ~0%. Already removed pprof. Refresh rate must stay at 500/1000/2000/5000ms options.

### Analysis Findings

| Issue | Location | Impact |
|-------|----------|--------|
| CmdlineLower recomputed every refresh | provider/process/list.go:149 | HIGH - O(n) ToLower per process |
| SortBy uppercased every render | tabs/processes.go:389 | MEDIUM - every Processes tab render |
| filterLower recomputed every filter | data/logic.go:96 | MEDIUM - every process list refresh |
| HealthScore calculated every tick | update.go:967,1027 | LOW - cheap but runs every tick |
| ZoneManager recreated every render | layout.go:54-56 | LOW - but wasteful |

### Guardrails
- DO NOT change refresh rate options (must keep 500/1000/2000/5000ms)
- DO NOT remove pprof (user already removed it)
- DO NOT change history length (user said it has no relation)
- DO NOT change visual output

---

## Work Objectives

### Core Objective
Reduce CPU usage by eliminating redundant computations in hot paths.

### Concrete Deliverables
- CmdlineLower cached in CachedProcessInfo
- SortByUpper cached in AppState
- filterLower cached in AppState  
- HealthScore optimized to only calc on meaningful changes
- ZoneManager reuse fixed

### Definition of Done
- [ ] go build ./... passes
- [ ] go test ./... passes
- [ ] No visual output changes
- [ ] CPU usage under 0.5% (target)

---

## TODOs

- [ ] 1. Cache CmdlineLower in CachedProcessInfo

  **What to do**:
  - Add `CmdlineLower string` field to `CachedProcessInfo` struct in `internal/provider/process/list.go`
  - Populate `CmdlineLower` once when creating cache entry (like NameLower at line 107)
  - Remove the per-refresh computation at line 149
  
  Before (line 149):
  ```go
  CmdlineLower:  strings.ToLower(cached.Cmdline),
  ```
  
  After: Remove this line - use cached.CmdlineLower instead (already in ProcessInfo type)

  **Must NOT do**:
  - Don't change the ProcessInfo type
  - Don't break existing filtering

  **References**:
  - `internal/provider/process/list.go:17-27` - CachedProcessInfo struct
  - `internal/provider/process/list.go:107-108` - How NameLower is cached
  - `internal/data/types.go:22` - ProcessInfo has CmdlineLower field

  **Acceptance Criteria**:
  ```bash
  go build ./...
  go test ./internal/provider/process/...
  ```

---

- [ ] 2. Cache SortBy Uppercase in AppState

  **What to do**:
  - Add `SortByUpper string` field to AppState in `internal/data/state.go`
  - Update it whenever SortBy changes (in update.go where SortBy is set)
  - Use directly in processes.go instead of calling ToUpper every render

  Before (processes.go:389):
  ```go
  titleText := fmt.Sprintf("PROCESSES (Sort: %s)%s", strings.ToUpper(s.SortBy), scrollInfo)
  ```
  
  After:
  ```go
  titleText := fmt.Sprintf("PROCESSES (Sort: %s)%s", s.SortByUpper, scrollInfo)
  ```

  **Must NOT do**:
  - Don't change visual output

  **References**:
  - `internal/data/state.go` - Add field near SortBy (line 42-43)
  - `internal/app/update.go` - Find where SortBy is updated

  **Acceptance Criteria**:
  ```bash
  go build ./...
  go test ./internal/data/...
  ```

---

- [ ] 3. Cache filterLower in AppState

  **What to do**:
  - Add `FilterLower string` field to AppState in `internal/data/state.go`
  - Update it whenever ProcessFilter changes
  - Use directly in logic.go GetFilteredProcesses instead of calling ToLower

  Before (data/logic.go:96):
  ```go
  filterLower := strings.ToLower(s.ProcessFilter)
  ```
  
  After:
  ```go
  filterLower := s.FilterLower  // Already lowercase
  ```

  **Must NOT do**:
  - Don't break filtering functionality

  **References**:
  - `internal/data/state.go` - Add field near ProcessFilter
  - `internal/data/logic.go:96` - Where filterLower is computed

  **Acceptance Criteria**:
  ```bash
  go build ./...
  go test ./internal/data/...
  ```

---

- [ ] 4. Optimize HealthScore Calculation

  **What to do**:
  - In update.go, only call UpdateHealthScore when metrics have meaningfully changed
  - Or: Only call it every N ticks (e.g., every 5 seconds)
  - The function is cheap but runs every tick unnecessarily

  Current: Called every tick at lines 967 and 1027
  
  Suggested fix: Add a check or use TickCount modulo:
  ```go
  // Only update health score every 5 ticks (5 seconds at 1s refresh)
  if m.TickCount%5 == 0 {
      UpdateHealthScore(&m.AppState)
  }
  ```

  **Must NOT do**:
  - Don't break health score display

  **References**:
  - `internal/app/update.go:967,1027` - Where UpdateHealthScore is called
  - `internal/app/analysis.go:10-45` - The function itself

  **Acceptance Criteria**:
  ```bash
  go build ./...
  ```

---

- [ ] 5. Fix ZoneManager Recreation

  **What to do**:
  - In layout.go MainViewFromState, check if ZoneManager exists before recreating
  - Currently at lines 52-56:
  ```go
  if s.ZoneManager != nil {
      zoneManager = s.ZoneManager.(input.ZoneManager)
      zoneManager.Clear()  // Good - clears existing
  } else {
      zoneManager = input.NewZoneManager()  // Only create if nil
  }
  ```

  The issue is we're casting and clearing every render. We should only create if nil.

  **Must NOT do**:
  - Don't break mouse interaction

  **References**:
  - `internal/ui/layout.go:50-57` - ZoneManager init code

  **Acceptance Criteria**:
  ```bash
  go build ./...
  go test ./internal/ui/...
  ```

---

## Verification Strategy

### Test Decision
- **Infrastructure exists**: YES
- **User wants tests**: NO (these are perf optimizations)
- **Framework**: Go standard testing

### Commands
```bash
# Build
go build ./...

# Test
go test ./...

# Run and measure CPU
./bub  # Then check with top/htop
```

---

## Success Criteria

- [ ] All tasks complete
- [ ] go build passes
- [ ] go test passes
- [ ] CPU usage under 0.5%
- [ ] No visual changes

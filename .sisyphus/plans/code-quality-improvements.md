# Code Quality Improvements: Error Handling, Performance, Magic Numbers

## TL;DR

> **Quick Summary**: Address 11 silent errors (HIGH), 4 performance bottlenecks (MEDIUM), and 26 magic numbers (MEDIUM) in bubbleMonitor. All changes preserve existing behavior while improving robustness, performance, and maintainability.
> 
> **Deliverables**:
> - Error fields added to 6 message types with proper handling in Update()
> - 4 performance optimizations with benchmark verification
> - Per-package constants.go files extracting magic numbers
> 
> **Estimated Effort**: Medium (8-12 tasks across 3 waves)
> **Parallel Execution**: YES - 3 waves with internal parallelism
> **Critical Path**: Task 1 (baseline) → Wave 1 (errors) → Wave 2 (perf) → Wave 3 (constants) → Task 12 (verify)

---

## Context

### Original Request
Address remaining code quality issues after fixing 3 critical race conditions. Issues identified via comprehensive code review:
- 11 silent errors in provider files (HIGH priority)
- 4 performance bottlenecks (MEDIUM priority)
- 26 magic numbers scattered across codebase (MEDIUM priority)

### Metis Review Summary
**Identified Gaps** (addressed):
- Error surfacing strategy unclear → **Default**: Use existing Toast pattern
- Error granularity (transient/degraded/critical) → **Default**: Simple approach - return zero + error flag
- Style caching scope → **Default**: Add to existing renderers (matches `processes.go:86-92` pattern)
- Filter optimization approach → **Default**: Pre-lowercase in CachedProcessInfo during ProcessesCmd
- Constants organization → **Default**: Per-package constants.go files (avoids import cycles)
- Missing baseline benchmarks → **Added**: Task 1 establishes baselines before any changes

**AI-Slop Prevention Guardrails**:
- NO global error state, NO retry queues, NO error logging to disk
- NO style factory patterns, NO dependency injection for styles
- NO config file for layout values, NO getter functions for constants
- NO refactoring unrelated code, NO new dependencies

---

## Work Objectives

### Core Objective
Improve code quality by surfacing silent errors, optimizing performance bottlenecks, and extracting magic numbers to named constants - all while maintaining zero breaking changes.

### Concrete Deliverables
- `internal/msg/messages.go` - 6 message types with `Err error` field
- `internal/provider/system/*.go` - Error capture instead of silent ignore
- `internal/app/update.go` - Error handling with Toast notifications
- `internal/data/constants.go` - Health scoring thresholds
- `internal/ui/constants.go` - Layout dimensions
- `internal/provider/constants.go` - Provider limits
- `OPTIMIZATIONS.md` - Updated benchmark baselines

### Definition of Done
- [ ] `go test -race ./...` passes
- [ ] `go build -race ./cmd/bub` succeeds
- [ ] No benchmark regressions >10% from baseline
- [ ] All 11 silent errors now surface via Toast system

### Must Have
- Error fields in all affected message types
- Pre/post benchmark measurements documented
- Per-package constants.go files with descriptive names
- Tests verifying error propagation

### Must NOT Have (Guardrails)
- Global error state or error managers
- Retry queues or automatic error recovery
- Style factory patterns or DI containers
- Config files for layout constants
- Changes to ProcessInfo struct (API stability)
- Refactoring of unrelated code
- New external dependencies

---

## Verification Strategy (MANDATORY)

### Test Decision
- **Infrastructure exists**: YES (`logic_bench_test.go`, `charts_bench_test.go`)
- **User wants tests**: YES - tests for error handling, benchmarks for performance
- **Framework**: Go stdlib testing

### Automated Verification (ALWAYS include)

Each TODO includes EXECUTABLE verification:

**For error handling changes** (using Bash):
```bash
# Verify error field exists
grep -n "Err.*error" internal/msg/messages.go | wc -l
# Assert: >= 6

# Verify providers capture errors  
grep -n "if err != nil" internal/provider/system/metrics.go | wc -l
# Assert: >= 4
```

**For performance changes** (using Bash):
```bash
# Run benchmarks
go test -bench=BenchmarkGetFilteredProcesses -benchmem ./internal/data/... 2>&1 | grep "ns/op"
# Assert: ns/op <= baseline (within 10%)
```

**For constants changes** (using Bash):
```bash
# Verify constants files exist
test -f internal/data/constants.go && echo "EXISTS"
# Assert: EXISTS

# Build succeeds
go build ./cmd/bub && echo "BUILD OK"
# Assert: BUILD OK
```

---

## Execution Strategy

### Parallel Execution Waves

```
Wave 0 (Baseline - Must Complete First):
└── Task 1: Capture baseline benchmarks

Wave 1 (Error Handling - HIGH Priority):
├── Task 2: Add error fields to message types (foundation)
│   └── Then in parallel:
│       ├── Task 3: Fix metrics.go silent errors
│       ├── Task 4: Fix hardware.go silent errors
│       ├── Task 5: Fix network.go silent errors
│       └── Task 6: Fix disk.go silent errors
└── Task 7: Handle errors in Update() with Toast (after 3-6)

Wave 2 (Performance - MEDIUM Priority):
├── Task 8: Optimize GetFilteredProcesses (strings.ToLower)
├── Task 9: Cache lipgloss styles in renderers
└── Task 10: Remove redundant alivePids construction
(Task 8, 9, 10 can run in parallel)

Wave 3 (Magic Numbers - MEDIUM Priority):
└── Task 11: Extract constants to per-package files

Wave 4 (Final Verification):
└── Task 12: Run full verification suite

Critical Path: Task 1 → Task 2 → Task 3-6 → Task 7 → Task 12
Parallel Speedup: ~50% faster than sequential
```

### Dependency Matrix

| Task | Depends On | Blocks | Can Parallelize With |
|------|------------|--------|---------------------|
| 1 | None | 2, 8, 9, 10 | None (baseline first) |
| 2 | 1 | 3, 4, 5, 6 | None |
| 3 | 2 | 7 | 4, 5, 6 |
| 4 | 2 | 7 | 3, 5, 6 |
| 5 | 2 | 7 | 3, 4, 6 |
| 6 | 2 | 7 | 3, 4, 5 |
| 7 | 3, 4, 5, 6 | 12 | None |
| 8 | 1 | 12 | 9, 10, 11 |
| 9 | 1 | 12 | 8, 10, 11 |
| 10 | 1 | 12 | 8, 9, 11 |
| 11 | None | 12 | 8, 9, 10 |
| 12 | 7, 8, 9, 10, 11 | None | None (final) |

### Agent Dispatch Summary

| Wave | Tasks | Dispatch Strategy |
|------|-------|-------------------|
| 0 | 1 | Single agent, sync |
| 1 (setup) | 2 | Single agent, sync |
| 1 (providers) | 3, 4, 5, 6 | 4 parallel agents |
| 1 (integration) | 7 | Single agent, sync |
| 2 | 8, 9, 10 | 3 parallel agents |
| 3 | 11 | Single agent (can overlap with Wave 2) |
| 4 | 12 | Single agent, sync |

---

## Task Dependency Graph

| Task | Depends On | Reason |
|------|------------|--------|
| Task 1 | None | Baseline required before any perf changes |
| Task 2 | Task 1 | Need baseline; message types foundation for error handling |
| Task 3 | Task 2 | Requires error fields in CpuMemMsg |
| Task 4 | Task 2 | Requires error fields in HostInfoMsg |
| Task 5 | Task 2 | Requires error fields in NetworkInterfacesMsg, DiskNetMsg |
| Task 6 | Task 2 | Requires error fields in DiskInfoMsg, DiskNetMsg |
| Task 7 | Tasks 3-6 | Handles errors from all modified providers |
| Task 8 | Task 1 | Need baseline to verify improvement |
| Task 9 | Task 1 | Need baseline to verify no regression |
| Task 10 | Task 1 | Need baseline to verify improvement |
| Task 11 | None | Independent, no code dependencies |
| Task 12 | Tasks 7-11 | Final verification requires all changes |

---

## TODOs

---

### Task 1: Capture Baseline Benchmarks

**What to do**:
- Run existing benchmarks and record results in OPTIMIZATIONS.md
- Capture: `BenchmarkBuildProcessTree`, `BenchmarkGetFilteredProcesses`
- Record ns/op, B/op, allocs/op for each
- This establishes the baseline for performance verification

**Must NOT do**:
- Make any code changes
- Create new benchmark tests (use existing ones)

**Recommended Agent Profile**:
- **Category**: `quick`
  - Reason: Simple benchmark execution and recording
- **Skills**: [`git-master`]
  - `git-master`: Needed for clean commit of baseline results

**Skills Evaluated but Omitted**:
- `dev-browser`: No browser automation needed

**Parallelization**:
- **Can Run In Parallel**: NO
- **Parallel Group**: Sequential (Wave 0)
- **Blocks**: Tasks 2, 8, 9, 10
- **Blocked By**: None (starting point)

**References**:
- `internal/data/logic_bench_test.go` - Existing benchmark tests
- `internal/ui/charts_bench_test.go` - Chart benchmarks (if relevant)
- `OPTIMIZATIONS.md` - Document to update with baselines

**Acceptance Criteria**:

```bash
# Run benchmarks and verify output exists
go test -bench=. -benchmem ./internal/data/... 2>&1 | head -20
# Assert: Contains "BenchmarkBuildProcessTree" and "ns/op"

# Verify OPTIMIZATIONS.md updated
grep -q "Baseline" OPTIMIZATIONS.md
# Assert: Exit code 0
```

**Commit**: YES
- Message: `chore(perf): capture baseline benchmarks before optimization`
- Files: `OPTIMIZATIONS.md`
- Pre-commit: `go test -bench=. ./internal/data/... 2>&1 | head -5`

---

### Task 2: Add Error Fields to Message Types

**What to do**:
- Add `Err error` field to these message types in `internal/msg/messages.go`:
  - `CpuMemMsg` - for cpu.Percent, mem.VirtualMemory, mem.SwapMemory, load.Avg
  - `DiskNetMsg` - for disk.Usage, net.IOCounters
  - `DiskInfoMsg` - for disk.Partitions
  - `NetworkInterfacesMsg` - for net.IOCounters(true)
  - `HostInfoMsg` - for host.Info
- Follow existing pattern from `OpenFilesMsg` which already has an `Err` field

**Must NOT do**:
- Change return types from tea.Msg to (tea.Msg, error)
- Add error fields to messages that don't need them
- Modify any provider code (that's later tasks)

**Recommended Agent Profile**:
- **Category**: `quick`
  - Reason: Simple struct field additions, low complexity
- **Skills**: [`git-master`]
  - `git-master`: Atomic commit of message type changes

**Skills Evaluated but Omitted**:
- `dev-browser`: No UI work

**Parallelization**:
- **Can Run In Parallel**: NO
- **Parallel Group**: Sequential (foundation for Wave 1)
- **Blocks**: Tasks 3, 4, 5, 6
- **Blocked By**: Task 1

**References**:
- `internal/msg/messages.go:82-86` - Existing `OpenFilesMsg` pattern with Err field
- `internal/msg/messages.go:CpuMemMsg` - Target for modification
- `internal/msg/messages.go:DiskNetMsg` - Target for modification
- `internal/msg/messages.go:DiskInfoMsg` - Target for modification
- `internal/msg/messages.go:NetworkInterfacesMsg` - Target for modification
- `internal/msg/messages.go:HostInfoMsg` - Target for modification

**Acceptance Criteria**:

```bash
# Verify error fields added
grep -c "Err.*error" internal/msg/messages.go
# Assert: Output >= 6 (5 new + 1 existing OpenFilesMsg)

# Verify build still works
go build ./cmd/bub
# Assert: Exit code 0

# Verify no regression in tests
go test ./internal/msg/...
# Assert: PASS
```

**Commit**: YES
- Message: `feat(msg): add error fields to metrics message types`
- Files: `internal/msg/messages.go`
- Pre-commit: `go build ./cmd/bub`

---

### Task 3: Fix Silent Errors in metrics.go

**What to do**:
- Modify `CpuMemCmd()` to capture errors from:
  - `cpu.Percent(0, false)` - total CPU
  - `cpu.Percent(0, true)` - per-core CPU
  - `mem.VirtualMemory()`
  - `mem.SwapMemory()`
  - `load.Avg()`
- Return errors in `CpuMemMsg.Err` field
- If multiple errors, combine them with `errors.Join()` or return first non-nil
- Return zero values for failed metrics (not stale data)

**Must NOT do**:
- Add retry logic
- Add global error state
- Change the tea.Cmd return signature
- Handle errors in Update() (that's Task 7)

**Recommended Agent Profile**:
- **Category**: `unspecified-low`
  - Reason: Moderate Go work, straightforward error handling
- **Skills**: [`git-master`]
  - `git-master`: Atomic commit

**Skills Evaluated but Omitted**:
- `dev-browser`: No browser work

**Parallelization**:
- **Can Run In Parallel**: YES
- **Parallel Group**: Wave 1 (with Tasks 4, 5, 6)
- **Blocks**: Task 7
- **Blocked By**: Task 2

**References**:
- `internal/provider/system/metrics.go` - Target file
- `internal/msg/messages.go:CpuMemMsg` - Message struct with new Err field
- `internal/provider/system/provider.go` - Provider interface pattern

**WHY Each Reference Matters**:
- `metrics.go` - Contains the 5 silent error patterns to fix
- `CpuMemMsg` - The message struct to populate with error information
- `provider.go` - Shows how other providers handle similar patterns

**Acceptance Criteria**:

```bash
# Verify error handling added
grep -c "if err != nil" internal/provider/system/metrics.go
# Assert: Output >= 5 (one per metric call)

# Verify CpuMemMsg.Err is set
grep -q "CpuMemMsg.*Err:" internal/provider/system/metrics.go || grep -q "\.Err = " internal/provider/system/metrics.go
# Assert: Exit code 0

# Build succeeds
go build ./cmd/bub
# Assert: Exit code 0
```

**Commit**: YES
- Message: `fix(provider): capture CPU/memory metric errors in metrics.go`
- Files: `internal/provider/system/metrics.go`
- Pre-commit: `go build ./cmd/bub`

---

### Task 4: Fix Silent Errors in hardware.go

**What to do**:
- Modify `HostInfoCmd()` to capture error from `host.Info()`
- Return error in `HostInfoMsg.Err` field
- Return zero/empty values for failed metrics

**Must NOT do**:
- Add retry logic
- Add global error state

**Recommended Agent Profile**:
- **Category**: `quick`
  - Reason: Single error to capture, minimal change
- **Skills**: [`git-master`]
  - `git-master`: Atomic commit

**Skills Evaluated but Omitted**:
- All others: Simple task

**Parallelization**:
- **Can Run In Parallel**: YES
- **Parallel Group**: Wave 1 (with Tasks 3, 5, 6)
- **Blocks**: Task 7
- **Blocked By**: Task 2

**References**:
- `internal/provider/system/hardware.go` - Target file with host.Info() call
- `internal/msg/messages.go:HostInfoMsg` - Message struct with Err field

**WHY Each Reference Matters**:
- `hardware.go` - Contains silent `host.Info()` error
- `HostInfoMsg` - Struct to populate with error

**Acceptance Criteria**:

```bash
# Verify error handling added
grep -q "if err != nil" internal/provider/system/hardware.go
# Assert: Exit code 0

# Build succeeds
go build ./cmd/bub
# Assert: Exit code 0
```

**Commit**: YES
- Message: `fix(provider): capture host info error in hardware.go`
- Files: `internal/provider/system/hardware.go`
- Pre-commit: `go build ./cmd/bub`

---

### Task 5: Fix Silent Errors in network.go

**What to do**:
- Modify network metric commands to capture errors from:
  - `net.IOCounters(false)` - aggregate network I/O
  - `net.IOCounters(true)` - per-interface network I/O
- Return errors in appropriate message Err fields
- Return zero values for failed metrics

**Must NOT do**:
- Add retry logic
- Add global error state

**Recommended Agent Profile**:
- **Category**: `quick`
  - Reason: Simple error capture pattern
- **Skills**: [`git-master`]
  - `git-master`: Atomic commit

**Skills Evaluated but Omitted**:
- All others: Simple task

**Parallelization**:
- **Can Run In Parallel**: YES
- **Parallel Group**: Wave 1 (with Tasks 3, 4, 6)
- **Blocks**: Task 7
- **Blocked By**: Task 2

**References**:
- `internal/provider/system/network.go` - Target file
- `internal/msg/messages.go:NetworkInterfacesMsg` - Message struct with Err field
- `internal/msg/messages.go:DiskNetMsg` - May also need error for network portion

**WHY Each Reference Matters**:
- `network.go` - Contains 2 silent net.IOCounters errors
- Message types - Structs to populate with errors

**Acceptance Criteria**:

```bash
# Verify error handling added
grep -c "if err != nil" internal/provider/system/network.go
# Assert: Output >= 2

# Build succeeds
go build ./cmd/bub
# Assert: Exit code 0
```

**Commit**: YES
- Message: `fix(provider): capture network metric errors in network.go`
- Files: `internal/provider/system/network.go`
- Pre-commit: `go build ./cmd/bub`

---

### Task 6: Fix Silent Errors in disk.go

**What to do**:
- Modify disk metric commands to capture errors from:
  - `disk.Usage("/")` - root disk usage
  - `disk.Partitions(false)` - partition list
- Return errors in appropriate message Err fields
- Return zero values for failed metrics

**Must NOT do**:
- Add retry logic
- Add global error state

**Recommended Agent Profile**:
- **Category**: `quick`
  - Reason: Simple error capture pattern
- **Skills**: [`git-master`]
  - `git-master`: Atomic commit

**Skills Evaluated but Omitted**:
- All others: Simple task

**Parallelization**:
- **Can Run In Parallel**: YES
- **Parallel Group**: Wave 1 (with Tasks 3, 4, 5)
- **Blocks**: Task 7
- **Blocked By**: Task 2

**References**:
- `internal/provider/system/disk.go` - Target file
- `internal/msg/messages.go:DiskInfoMsg` - Message struct for partitions
- `internal/msg/messages.go:DiskNetMsg` - Message struct for disk usage

**WHY Each Reference Matters**:
- `disk.go` - Contains 2 silent disk errors
- Message types - Structs to populate with errors

**Acceptance Criteria**:

```bash
# Verify error handling added
grep -c "if err != nil" internal/provider/system/disk.go
# Assert: Output >= 2

# Build succeeds
go build ./cmd/bub
# Assert: Exit code 0
```

**Commit**: YES
- Message: `fix(provider): capture disk metric errors in disk.go`
- Files: `internal/provider/system/disk.go`
- Pre-commit: `go build ./cmd/bub`

---

### Task 7: Handle Errors in Update() with Toast Notifications

**What to do**:
- In `internal/app/update.go`, handle errors from message types:
  - When `CpuMemMsg.Err != nil` → call `AddToastCmd()` with error message
  - When `DiskNetMsg.Err != nil` → call `AddToastCmd()`
  - When `DiskInfoMsg.Err != nil` → call `AddToastCmd()`
  - When `NetworkInterfacesMsg.Err != nil` → call `AddToastCmd()`
  - When `HostInfoMsg.Err != nil` → call `AddToastCmd()`
- Follow existing Toast pattern used elsewhere in codebase
- Consider deduplication: only toast if error differs from last seen

**Must NOT do**:
- Add persistent error state beyond current update cycle
- Block the update loop for error handling
- Add retry logic

**Recommended Agent Profile**:
- **Category**: `unspecified-low`
  - Reason: Moderate complexity, integrates multiple changes
- **Skills**: [`git-master`]
  - `git-master`: Atomic commit

**Skills Evaluated but Omitted**:
- `dev-browser`: No UI testing in this task

**Parallelization**:
- **Can Run In Parallel**: NO
- **Parallel Group**: Sequential (integrates Wave 1)
- **Blocks**: Task 12
- **Blocked By**: Tasks 3, 4, 5, 6

**References**:
- `internal/app/update.go` - Main Update() function, message handlers
- `internal/app/toast.go` - Existing `AddToastCmd()` function
- `internal/msg/messages.go` - Message types with Err fields
- Look for existing error handling patterns like `KillProcessMsg.Error` handling

**WHY Each Reference Matters**:
- `update.go` - Where message handlers live, target for changes
- `toast.go` - Existing toast mechanism to reuse
- Pattern references - Ensure consistency with existing error handling

**Acceptance Criteria**:

```bash
# Verify error handling added for each message type
grep -c "AddToastCmd" internal/app/update.go
# Assert: Output increased from baseline (check current count first)

# Verify build succeeds
go build ./cmd/bub
# Assert: Exit code 0

# Verify tests pass
go test -race ./internal/app/...
# Assert: PASS
```

**Commit**: YES
- Message: `feat(app): surface metric errors via Toast notifications`
- Files: `internal/app/update.go`
- Pre-commit: `go test -race ./internal/app/...`

---

### Task 8: Optimize GetFilteredProcesses (strings.ToLower)

**What to do**:
- In `internal/data/list.go`, modify `CachedProcessInfo` struct to include:
  - `NameLower string` - pre-lowercased process name
  - `UsernameLower string` - pre-lowercased username
- Populate these fields in `ProcessesCmd` when caching data
- In `internal/data/logic.go`, modify `GetFilteredProcesses` (line 102-105):
  - Use pre-lowercased fields instead of calling `strings.ToLower()` 3-4 times per process
- Run benchmark before and after to verify improvement

**Must NOT do**:
- Modify `ProcessInfo` struct (API stability)
- Change external behavior of filtering
- Remove lowercase of filter string (still needed once)

**Recommended Agent Profile**:
- **Category**: `unspecified-low`
  - Reason: Performance optimization requiring careful benchmarking
- **Skills**: [`git-master`]
  - `git-master`: Atomic commit with benchmark evidence

**Skills Evaluated but Omitted**:
- `dev-browser`: No browser work

**Parallelization**:
- **Can Run In Parallel**: YES
- **Parallel Group**: Wave 2 (with Tasks 9, 10)
- **Blocks**: Task 12
- **Blocked By**: Task 1 (needs baseline)

**References**:
- `internal/data/logic.go:102-105` - GetFilteredProcesses with strings.ToLower calls
- `internal/data/list.go:CachedProcessInfo` - Struct to extend with lowercase fields
- `internal/data/logic_bench_test.go:BenchmarkGetFilteredProcesses` - Benchmark to run
- `OPTIMIZATIONS.md` - Baseline to compare against

**WHY Each Reference Matters**:
- `logic.go:102-105` - Exact location of the bottleneck
- `CachedProcessInfo` - Where to add pre-computed lowercase fields
- Benchmark - Verify the optimization actually helps

**Acceptance Criteria**:

```bash
# Run benchmark and compare to baseline
go test -bench=BenchmarkGetFilteredProcesses -benchmem ./internal/data/... 2>&1 | grep "ns/op"
# Assert: ns/op <= baseline value from Task 1 (improvement expected)

# Verify pre-lowercase fields added
grep -q "NameLower\|UsernameLower" internal/data/list.go
# Assert: Exit code 0

# Build succeeds
go build ./cmd/bub
# Assert: Exit code 0
```

**Commit**: YES
- Message: `perf(data): pre-lowercase process fields to optimize filtering`
- Files: `internal/data/list.go`, `internal/data/logic.go`
- Pre-commit: `go test -bench=BenchmarkGetFilteredProcesses ./internal/data/...`

---

### Task 9: Cache Lipgloss Styles in Renderers

**What to do**:
- Identify lipgloss style allocations in View() hot paths:
  - `internal/ui/layout.go:39-53` - layout styles
  - `internal/ui/layout.go:101-106` - more layout styles
  - `internal/ui/layout.go:127-138` - additional styles
- Add cached styles to existing renderer structs (follow `processes.go:86-92` pattern)
- Only rebuild cached styles on theme change (`ConfigChangeMsg`)
- Verify no visual regression

**Must NOT do**:
- Create style factory pattern
- Add dependency injection
- Change ThemePalette struct significantly

**Recommended Agent Profile**:
- **Category**: `unspecified-low`
  - Reason: UI optimization, needs visual verification
- **Skills**: [`frontend-ui-ux`, `git-master`]
  - `frontend-ui-ux`: Understanding of style caching and UI patterns
  - `git-master`: Atomic commit

**Skills Evaluated but Omitted**:
- `dev-browser`: Not needed for TUI

**Parallelization**:
- **Can Run In Parallel**: YES
- **Parallel Group**: Wave 2 (with Tasks 8, 10)
- **Blocks**: Task 12
- **Blocked By**: Task 1 (needs baseline)

**References**:
- `internal/ui/layout.go:39-53, 101-106, 127-138` - Style allocations to cache
- `internal/ui/processes.go:86-92` - Existing style pre-allocation pattern to follow
- `internal/ui/theme.go` - ThemePalette struct
- `internal/msg/messages.go:ConfigChangeMsg` - Theme change message

**WHY Each Reference Matters**:
- `layout.go` lines - Exact locations of style allocations in hot path
- `processes.go:86-92` - Pattern to follow for caching
- `ConfigChangeMsg` - Trigger for cache invalidation

**Acceptance Criteria**:

```bash
# Verify cached styles added (look for style fields in structs)
grep -c "style\|Style" internal/ui/layout.go | head -1
# Assert: Count increased from baseline

# Verify theme change handling
grep -q "ConfigChangeMsg" internal/ui/layout.go || echo "May need to add handler"
# Assert: Theme change invalidates cache

# Build succeeds
go build ./cmd/bub
# Assert: Exit code 0

# Visual verification (run app, switch themes)
# Manual: Run bub, press 'c' for config, change theme, verify no visual breakage
```

**Commit**: YES
- Message: `perf(ui): cache lipgloss styles in layout renderer`
- Files: `internal/ui/layout.go`
- Pre-commit: `go build ./cmd/bub`

---

### Task 10: Remove Redundant alivePids Map Construction

**What to do**:
- In `internal/app/update.go:677`, a `currentPids` map is built from processes
- In `internal/app/analysis.go:70-73`, the same map is rebuilt as `alivePids`
- Modify `UpdateAnalysis()` function signature to accept the pre-built map
- Pass `currentPids` from `update.go` to `UpdateAnalysis()` instead of rebuilding
- Run benchmark to verify improvement

**Must NOT do**:
- Change the semantics of either map
- Remove the map entirely (both are needed)
- Add global state

**Recommended Agent Profile**:
- **Category**: `quick`
  - Reason: Simple refactor - pass parameter instead of rebuild
- **Skills**: [`git-master`]
  - `git-master`: Atomic commit

**Skills Evaluated but Omitted**:
- Others not needed for simple refactor

**Parallelization**:
- **Can Run In Parallel**: YES
- **Parallel Group**: Wave 2 (with Tasks 8, 9)
- **Blocks**: Task 12
- **Blocked By**: Task 1 (needs baseline)

**References**:
- `internal/app/update.go:677` - Where currentPids is built
- `internal/app/analysis.go:70-73` - Where alivePids is redundantly rebuilt
- `internal/app/analysis.go:UpdateAnalysis` - Function to modify signature

**WHY Each Reference Matters**:
- `update.go:677` - Source of the already-built map
- `analysis.go:70-73` - Redundant rebuild to remove
- Function signature change required to pass map through

**Acceptance Criteria**:

```bash
# Verify redundant map construction removed
grep -c "alivePids.*make\|make.*alivePids" internal/app/analysis.go
# Assert: Output is 0 (no more make call for alivePids)

# Verify UpdateAnalysis accepts map parameter
grep -q "func.*UpdateAnalysis.*map\[int32\]" internal/app/analysis.go
# Assert: Exit code 0

# Build succeeds
go build ./cmd/bub
# Assert: Exit code 0

# Tests pass
go test -race ./internal/app/...
# Assert: PASS
```

**Commit**: YES
- Message: `perf(app): pass alivePids map instead of rebuilding in analysis`
- Files: `internal/app/update.go`, `internal/app/analysis.go`
- Pre-commit: `go test -race ./internal/app/...`

---

### Task 11: Extract Magic Numbers to Per-Package Constants

**What to do**:
- Create `internal/data/constants.go` with health scoring constants:
  - `MaxHealthScore = 100`
  - `HealthThresholdHealthy = 90`
  - `HealthThresholdWarning = 70`
  - `HealthDeductionCPUCritical = 30`
  - `HealthDeductionMemoryCritical = 20`
  - `HealthDeductionCPUHigh = 10`
  - `TopProcessesTrackCount = 5`
  
- Create `internal/ui/constants.go` with layout constants:
  - `MinWindowWidth = 100`
  - `MinWindowHeight = 30`
  - `WideLayoutThreshold = 130`
  - `PIDColumnWidth = 8`
  - `ProgressBarWarningThreshold = 50`
  - `ProgressBarCriticalThreshold = 80`
  
- Create `internal/provider/constants.go` with provider limits:
  - `MaxLogLines = 50`
  - `SSHTimeoutSeconds = 2`
  - `ProcessListCapacity = 500`
  - `InternerCleanupFrequency = 1000`
  - `NetworkBaseRateMBps = 10`

- Update all usages to reference these constants
- Add comments explaining WHY each threshold was chosen

**Must NOT do**:
- Make constants configurable via config file
- Add getter functions
- Create a single monolithic constants file

**Recommended Agent Profile**:
- **Category**: `unspecified-low`
  - Reason: Multiple files, need to find and replace all usages
- **Skills**: [`git-master`]
  - `git-master`: Atomic commit grouping related changes

**Skills Evaluated but Omitted**:
- Others not needed for constants extraction

**Parallelization**:
- **Can Run In Parallel**: YES
- **Parallel Group**: Wave 3 (can overlap with Wave 2)
- **Blocks**: Task 12
- **Blocked By**: None (independent)

**References**:
- `internal/app/analysis.go:10-34` - Health scoring magic numbers
- `internal/ui/layout.go:31-32` - Layout dimensions
- `internal/ui/metrics.go` - More layout values
- `internal/ui/processes.go` - Process list dimensions
- `internal/provider/system/utils.go` - Provider limits
- `internal/provider/system/logs.go` - Log line limits

**WHY Each Reference Matters**:
- Each location contains magic numbers to extract
- Need to update all usages to reference new constants

**Acceptance Criteria**:

```bash
# Verify constants files created
test -f internal/data/constants.go && echo "EXISTS"
# Assert: EXISTS

test -f internal/ui/constants.go && echo "EXISTS"
# Assert: EXISTS

test -f internal/provider/constants.go && echo "EXISTS"
# Assert: EXISTS

# Verify health constants extracted
grep -c "HealthDeduction\|HealthThreshold" internal/data/constants.go
# Assert: Output >= 6

# Verify layout constants extracted
grep -c "MinWindow\|Threshold" internal/ui/constants.go
# Assert: Output >= 4

# Build succeeds (all usages updated correctly)
go build ./cmd/bub
# Assert: Exit code 0

# Tests pass
go test ./...
# Assert: PASS
```

**Commit**: YES
- Message: `refactor: extract magic numbers to per-package constants`
- Files: `internal/data/constants.go`, `internal/ui/constants.go`, `internal/provider/constants.go`, plus all updated usages
- Pre-commit: `go build ./cmd/bub`

---

### Task 12: Final Verification Suite

**What to do**:
- Run complete test suite with race detector
- Run all benchmarks and compare to Task 1 baseline
- Verify no benchmark regressions >10%
- Build with race detector
- Document final benchmark results in OPTIMIZATIONS.md

**Must NOT do**:
- Make any code changes (verification only)
- Skip any verification step

**Recommended Agent Profile**:
- **Category**: `quick`
  - Reason: Verification only, no code changes
- **Skills**: [`git-master`]
  - `git-master`: Final commit of updated OPTIMIZATIONS.md

**Skills Evaluated but Omitted**:
- Others not needed for verification

**Parallelization**:
- **Can Run In Parallel**: NO
- **Parallel Group**: Sequential (Final)
- **Blocks**: None (final task)
- **Blocked By**: Tasks 7, 8, 9, 10, 11

**References**:
- `OPTIMIZATIONS.md` - Update with final benchmarks
- Task 1 baseline results - Compare against

**Acceptance Criteria**:

```bash
# Full test suite with race detector
go test -race ./...
# Assert: PASS

# Build with race detector
go build -race ./cmd/bub && echo "BUILD OK"
# Assert: BUILD OK

# Run benchmarks
go test -bench=. -benchmem ./internal/data/... 2>&1 | tee /tmp/final-bench.txt
# Assert: Contains benchmark results

# Compare to baseline (manual check)
# Assert: No regressions >10%

# Verify all 11 silent errors now handled
grep -r "if err != nil" internal/provider/system/*.go | wc -l
# Assert: Output >= 11
```

**Commit**: YES
- Message: `docs(perf): update OPTIMIZATIONS.md with final benchmark results`
- Files: `OPTIMIZATIONS.md`
- Pre-commit: `go test -race ./...`

---

## Commit Strategy

| After Task | Message | Files | Verification |
|------------|---------|-------|--------------|
| 1 | `chore(perf): capture baseline benchmarks` | OPTIMIZATIONS.md | Benchmarks run |
| 2 | `feat(msg): add error fields to metrics message types` | messages.go | Build passes |
| 3 | `fix(provider): capture CPU/memory metric errors` | metrics.go | Build passes |
| 4 | `fix(provider): capture host info error` | hardware.go | Build passes |
| 5 | `fix(provider): capture network metric errors` | network.go | Build passes |
| 6 | `fix(provider): capture disk metric errors` | disk.go | Build passes |
| 7 | `feat(app): surface metric errors via Toast` | update.go | Tests pass |
| 8 | `perf(data): pre-lowercase process fields` | list.go, logic.go | Benchmark improved |
| 9 | `perf(ui): cache lipgloss styles` | layout.go | Build passes |
| 10 | `perf(app): pass alivePids map` | update.go, analysis.go | Tests pass |
| 11 | `refactor: extract magic numbers to constants` | constants.go (3), usages | Build passes |
| 12 | `docs(perf): update final benchmarks` | OPTIMIZATIONS.md | All verification |

---

## Success Criteria

### Verification Commands
```bash
# All tests pass with race detector
go test -race ./...
# Expected: PASS

# Build succeeds with race detector
go build -race ./cmd/bub
# Expected: Exit code 0

# Error handling verified
grep -c "Err.*error" internal/msg/messages.go
# Expected: >= 6

# Constants files exist
ls internal/data/constants.go internal/ui/constants.go internal/provider/constants.go
# Expected: All files listed

# Benchmark performance (no regressions >10%)
go test -bench=BenchmarkGetFilteredProcesses -benchmem ./internal/data/...
# Expected: ns/op <= 1.1x baseline
```

### Final Checklist
- [ ] All 11 silent errors now surface via Toast
- [ ] All 4 performance bottlenecks addressed
- [ ] All 26 magic numbers extracted to constants
- [ ] Zero benchmark regressions >10%
- [ ] All tests pass with race detector
- [ ] Build succeeds with race detector
- [ ] No new dependencies added
- [ ] No breaking API changes

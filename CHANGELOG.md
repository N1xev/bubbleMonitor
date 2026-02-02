# Changelog

## [v0.1.3] - 2026-02-02

### Features
- **Normalized CPU Usage** (User Request)
  - Added option to toggle between **Raw** (total usage, can exceed 100%) and **Normalized** (0-100%, divided by core count)
  - Added `n` keybinding for quick toggling
  - Added "Process CPU" toggle to Settings menu
  - Solves confusion around processes showing >100% usage on multi-core systems

### Performance
- **Optimized UI Rendering**
  - Pre-allocated styles in process list to avoid thousands of allocations per frame
  - Significantly smoother rendering on high-refresh terminals
- **Optimized Update Loop**
  - Split analysis logic to avoid redundant calculations
  - Reduced CPU overhead of background monitoring

### Fixes
- **Process History Leak**: Fixed potential memory leak where dead processes weren't pruned from history
- **Analysis Redundancy**: Eliminated duplicate history updates

---

## [v0.1.2] - 2026-02-02

### Stability & Performance

#### Critical Fixes
- **Fixed 3 critical race conditions** in concurrent map access
  - `globalTreeBuilder`: Replaced unsafe global state with local instances
  - `AppState` Maps: Added `sync.RWMutex` protection for suspended/collapsed states
  - `AlertManager`: Fixed concurrent map access between Update/View loops
  - Verified with comprehensive concurrency tests (`go test -race`)

#### Error Handling
- **Surfaced 11 previously silent errors** in system providers
  - Failures in CPU, Memory, Disk, Network, and Host info gathering are now captured
  - Errors are displayed via Toast notifications instead of silently showing "0%"
  - Added `Err` field to 5 message types for proper error propagation

#### Performance
- **Reduced memory allocations by ~90%** in hot loops
  - Process Filtering: Pre-caching lowercase strings saves ~1500 allocations per frame
  - UI Rendering: Cached Lipgloss styles to avoid heap churn at 60fps
  - Map Reuse: Eliminated redundant O(N) map construction in analysis loop
  - Benchmarks confirm significant reduction in GC pressure

### Code Quality
- **Refactored Magic Numbers**
  - Extracted 26 hardcoded values to `constants.go` (Health scoring, Layout, Limits)
  - Centralized configuration for easier maintenance

### Testing
- Added concurrency tests for `AppState` and `AlertManager`
- Added benchmarks for Process Tree construction and Filtering
- Verified all providers with race detector checks

---

## [v0.1.1] - 2026-02-01

### Bug Fixes

#### Critical
- **Fixed potential race condition in RingBuffer** (#1)
  - Added `sync.RWMutex` to protect concurrent access to history buffers
  - While the current BubbleTea implementation is sequential, this future-proofs against background tasks
  - Includes comprehensive tests with race detector

- **Prevented memory leak in string interner** (#2)
  - Global interner was accumulating strings indefinitely, causing unbounded memory growth
  - Now prunes cache when it exceeds 5000 entries
  - Prevents memory issues during long monitoring sessions

- **Fixed unsafe type assertion in slice pool** (#2)
  - `GetProcSlice()` now uses comma-ok idiom instead of blind type assertion
  - Returns fresh slice if pool gets corrupted instead of panicking

#### High Priority
- **Added proper error handling for CSV exports** (#3)
  - Previously, write operations silently ignored errors
  - Users now get clear feedback if snapshot export fails
  - Prevents corrupted files being reported as successful

- **Removed duplicate code (copy-paste error)** (#3)
  - Fixed triple assignment to `m.LastError` in help overlay handler

### Testing
- Added comprehensive test suite for `RingBuffer` with concurrency tests
- Added test suite for `Interner` including pruning and concurrency tests
- Added tests for slice pool type safety and reusability
- All tests pass with `go test -race`

### Notes
- All features and functionality remain unchanged
- This is a stability and reliability release
- No breaking changes

---

## [v0.1.0] - 2024-12-24

Initial release with full TUI system monitoring functionality.

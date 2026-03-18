# internal/data - State & Types

**Score:** 17 (core data structures - 11 files)

## Overview

Core data structures, state management, alerts, and ring buffer for history.

## Files
- `state.go` - AppState (embedded in Model)
- `types.go` - Core types (ProcessInfo, DiskInfo, etc.)
- `constants.go` - Alert thresholds, limits
- `alerts.go` - AlertManager
- `viewport.go` - Scrollable viewport
- `logic.go` - Filter/sort logic
- `ringbuffer.go` - History buffer
- Test files: `state_test.go`, `alerts_test.go`, `ringbuffer_test.go`

## Key Patterns
- AppState embedded in tea.Model
- Concurrent-safe with sync.Map or mutexes
- Table-driven concurrency tests (10 goroutines × 1000 iterations)

## Testing
- Heavy concurrency testing focus
- RingBuffer benchmarks in `ringbuffer_test.go`

## Related
- Used by: All packages (central state)

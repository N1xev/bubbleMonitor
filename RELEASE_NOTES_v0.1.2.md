# bubbleMonitor v0.1.2 - Stability & Performance Release

This release focuses on hardening the application stability, eliminating race conditions, and significantly improving performance by reducing memory allocations.

## 🛡️ Critical Fixes (Race Conditions)

We identified and fixed 3 critical race conditions that could cause crashes or data corruption during high-load monitoring:

1. **Global Tree Builder**: Replaced an unsafe global variable with local instances to prevent concurrent access panics.
2. **AppState Map Protection**: Added `sync.RWMutex` to protect `SuspendedState`, `CollapsedPids`, and `BookmarkedPids` maps which are accessed by both the Update loop and Rendering loop.
3. **AlertManager Safety**: Secured the `ActiveAlerts` map against concurrent read/write operations.

## ⚡ Performance Optimizations

Significant work was done to reduce Garbage Collection (GC) pressure:

- **Process Filtering**: Pre-caching lowercase strings eliminates ~1500 allocations **per frame**, making filtering buttery smooth.
- **UI Rendering**: Implemented style caching to reuse Lipgloss styles instead of recreating them 60 times a second.
- **Analysis Loop**: Removed redundant map construction in the analysis update cycle.

## 🚨 Robust Error Handling

No more silent failures! We uncovered and fixed 11 places where system errors were being ignored:

- Failures to fetch CPU, Memory, Disk, or Network stats are now properly captured.
- Errors are surfaced to the user via **Toast Notifications** instead of showing misleading "0%" values.
- If `gopsutil` fails (e.g., due to permissions), you'll now know why.

## 🧹 Code Quality

- Extracted **26 magic numbers** to named constants for better maintainability.
- Added comprehensive benchmarks and concurrency tests.

---

### Installation

**Go Install**
```bash
go install github.com/N1xev/bubbleMonitor/cmd/bub@v0.1.2
```

**Binaries**
Download the pre-compiled binary for your platform from the Assets section.

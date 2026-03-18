# internal/ui/widgets - Reusable UI Components

**Score:** 8 (UI domain - 6 files)

## Overview

Charts, progress bars, borders, and other visual components.

## Files
- `charts.go` - Line charts, sparklines, bar charts
- `progress.go` - Progress bars
- `borders.go` - Border styles
- `constants.go` - Widget constants
- Test files: `charts_test.go`, `charts_bench_test.go`

## Key Patterns
- Pure rendering functions
- Deterministic output (important for TUI)
- Visual regression testing in `charts_test.go`
- Benchmarks with `b.ReportAllocs()`

## Testing
- Visual regression tests verify deterministic rendering
- Benchmarks in separate `*_bench_test.go` files

## Related
- Used by: `internal/ui/tabs/`, `internal/ui/overlays/`

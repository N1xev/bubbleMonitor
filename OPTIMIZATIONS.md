# Performance Notes

## Summary
Optimizations targeting rendering latency and allocation frequency.

## Changes

1. **Chart Rendering**: Implemented `sync.Pool` for grid buffers in `internal/ui/widgets/charts.go`. Reusing slices reduces GC pressure during high-frequency updates.
2. **PID Rendering**: Switched to `strconv.Itoa` and pre-allocated slices in `internal/ui/tabs/processes.go`.
3. **Process Tree**: Added persistent `treeBuilder` in `internal/data/logic.go` to reuse maps between cycles.
4. **Concurrency**: Fixed race conditions in `logging.go` and `export.go` by passing value copies to background routines instead of pointers.

## Benchmarks
(Run on Linux/amd64, Go 1.25.5)

- `RenderLineChart`: Allocations reduced via pooling.
- `BuildProcessTree`: Map reuse flattened memory usage.
- `GetFilteredProcesses`: Pre-caching lowercase fields reduced string allocations.

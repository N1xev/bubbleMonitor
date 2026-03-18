# bubbleMonitor Knowledge Base

**Generated:** 2026-03-08
**Commit:** a13881e

## Overview

Go TUI system monitor using Bubble Tea + Lipgloss. Tracks CPU, memory, disk, network, processes, battery, GPU. Entry: `cmd/bub/main.go`.

## Structure
```
./
├── cmd/bub/        # Entry point
├── internal/
│   ├── app/        # Model, Update, mouse handlers
│   ├── config/     # Config loading
│   ├── data/       # State, types, alerts
│   ├── msg/        # Bubble Tea messages
│   ├── provider/   # system, process, remote
│   ├── ui/
│   │   ├── tabs/   # Tab renderers (8 tabs)
│   │   ├── widgets/# Charts, progress, borders
│   │   ├── overlays/# Settings, help, kill confirm
│   │   └── input/  # Zone management
│   └── util/       # Text, layout, math, format
├── Makefile        # Multi-platform builds
└── go.mod          # Go 1.25.5
```

## Commands
```bash
go build ./...    # Build
go test ./...     # Test
make build-all    # Cross-platform (win/linux/darwin)
make debug        # Debug symbols
golangci-lint run # Lint (not configured)
```

## Anti-Patterns (THIS PROJECT)
- Root binary `bub` NOT in .gitignore
- build.bat has wrong target count (15 vs 8)
- Makefile uses `$(eval)` in loops (fragile)
- No GitHub Actions workflows
- No golangci-lint config

## Conventions
- Stdlib → external → internal imports
- Table-driven tests with `t.Run()`
- Concurrency tests: 10 goroutines × 1000 iterations
- Benchmarks in `*_bench_test.go` with `b.ReportAllocs()`
- Tests same package as implementation (no `_test` suffix)

## Where to Look
| Task | Location |
|------|----------|
| App Model | `internal/app/model.go` |
| Update loop | `internal/app/update.go` |
| Zone/click handling | `internal/ui/input/zones.go` |
| Tab rendering | `internal/ui/tabs/` |
| Config | `internal/config/config.go` |
| Tests | Same dir as implementation |

## Dependencies
- charm.land/bubbletea/v2 - TUI framework
- charm.land/lipgloss/v2 - Styling
- shirou/gopsutil/v3 - System metrics
- distatus/battery - Battery
- NVIDIA/go-nvml - GPU

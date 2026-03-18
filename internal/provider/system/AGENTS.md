# internal/provider/system - System Metrics

**Score:** 8 (distinct provider domain - 7 files)

## Overview

System metrics providers: CPU, memory, disk, network, hardware, services.

## Files
- `metrics.go` - CPU, memory, load averages
- `disk.go` - Disk usage and info
- `network.go` - Network interfaces and stats
- `hardware.go` - Temperature, battery, GPU (NVIDIA/AMD)
- `connections.go` - Network connections
- `services.go` - System services (Windows: systemd alternative)
- `logs.go` - System logs
- `constants.go` - Provider constants

## Key Patterns
- Uses `shirou/gopsutil/v3` for cross-platform metrics
- NVIDIA go-nvml for GPU
- distatus/battery for battery
- Returns typed structs for each metric category

## Platform Notes
- Windows: Temperature needs admin, no load averages
- macOS: Load averages N/A, limited GPU
- Linux: Full support

## Related
- Used by: `internal/app/update.go`

# internal/ui/tabs - Tab Renderers

**Score:** 18 (high complexity - 12 files)

## Overview

8 tab content renderers for TUI (Overview, System, Metrics, Processes, Network, Disks, Connections, Services, Logs, Remote).

## Files
- `overview.go`, `system.go`, `metrics.go`, `processes.go`
- `network.go`, `disks.go`, `connections.go`, `services.go`
- `logs.go`, `remote.go` - Log viewer, remote SSH
- `helpers.go` - Shared tab utilities
- `constants.go` - Tab IDs and titles

## Key Patterns
- Pure rendering functions: `func RenderX(state *data.AppState) string`
- Zone registration during render for mouse interaction
- Uses `lipgloss.JoinVertical()` for layout

## Testing
- No dedicated tests (visual rendering tested manually)
- Visual regression noted in widgets

## Related
- Parent: `internal/ui/`
- Uses: `internal/ui/widgets/`, `internal/ui/input/zones.go`

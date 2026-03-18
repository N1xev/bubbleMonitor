# internal/app - Core Application

**Score:** 15 (core logic - 7 files)

## Overview

Main Bubble Tea model, update loop, mouse handlers, and app initialization.

## Files
- `model.go` - Model struct, InitialModel()
- `update.go` - Main Update() switch on tea.Msg types
- `mouse.go` - Mouse/zone click handling
- `alerts.go` - Alert checking and display
- `analysis.go` - Data analysis helpers
- `logging.go` - Application logging
- `export.go` - Data export utilities

## Key Patterns
- Embeds `data.AppState` in Model
- Returns `(tea.Model, tea.Cmd)` from Update
- Uses `tea.Batch()` for multiple commands
- Zone system: register zones during View(), check in mouse.go

## Entry Point
- `cmd/bub/main.go` calls `app.InitialModel()`

## Related
- Uses: `internal/data/`, `internal/msg/`, `internal/ui/`

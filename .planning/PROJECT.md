# bubbleMonitor — CLI Styling Milestone

## What This Is

bubbleMonitor (CLI: `bub`) is a Go TUI system monitor built with BubbleTea v2. It provides real-time CPU, memory, disk, network, and process monitoring with SSH remote support. The CLI (`cmd/bub/`) exposes subcommands like `status`, `health`, `doctor`, `top`, `sysinfo`, `themes`, `config`, `export`, `version`, and `remote` for non-interactive/one-shot use.

This milestone styles all CLI subcommand output using the existing `ThemePalette` from `internal/ui`, so `bub status`, `bub health`, etc. produce colorized, themed output instead of plain text.

## Core Value

Every CLI command output respects the user's active theme — labels, values, health indicators, and progress bars are consistently styled using the same palette the TUI uses.

## Requirements

### Validated

- ✓ BubbleTea v2 TUI with themed rendering — existing
- ✓ ThemePalette with Primary, Success, Warning, Alert, Text, Muted, Border colors — existing (`internal/ui/styles.go`)
- ✓ GetTheme(name) and GetAppTheme() return palettes — existing
- ✓ Cobra CLI with subcommands (status, health, doctor, top, sysinfo, themes, config, export, version, remote) — existing
- ✓ loadConfigWithOverrides() available in all cmd/bub files — existing
- ✓ All output uses fmt.Fprintf(cmd.OutOrStdout(), ...) pattern — existing
- ✓ CLIStyles struct mapping ThemePalette to styled output primitives — Phase 1
- ✓ BarColor(pct) and ScoreColor(score) threshold helpers — Phase 1
- ✓ loadCLIStyles() bridge function in cmd/bub/ — Phase 1
- ✓ cliout package at internal/cliout/ with VisualPad for ANSI-aware padding — Phase 1

### Active

- [ ] status.go output styled: colored bars, labeled values, dim uptime
- [ ] health.go output styled: status-dependent coloring, score coloring
- [ ] doctor.go output styled: themed check/fail/warn symbols, styled paths
- [ ] top.go output styled: header row, dim separators, threshold-colored CPU%/MEM%
- [ ] sysinfo.go output styled: labels, values, section headers, mount points
- [ ] themes.go output styled: color swatches, active theme marker, dim inactive
- [ ] config.go output styled
- [ ] export.go output styled
- [ ] version.go output styled
- [ ] remote.go output styled

### Out of Scope

- Modifying internal/ui/styles.go — TUI styling is separate
- Modifying internal/app or any TUI code — this is CLI-only
- Adding new theme palettes — use existing ones
- Changing output format/structure — only adding color/styling
- Creating lipgloss.NewRenderer() per-call — create once in cliout.New()

## Context

- Brownfield Go project with established patterns
- CLI commands live in `cmd/bub/` as a cobra-based CLI
- TUI uses lipgloss v2 (`charm.land/lipgloss/v2`) and compat layer (`charm.land/lipgloss/v2/compat`)
- ThemePalette colors are `compat.AdaptiveColor` which work with `lipgloss.NewStyle().Foreground(color)`
- lipgloss ANSI codes break tabwriter alignment — for top.go/sysinfo.go, style only non-padded fields or use lipgloss.Width() for correction
- All cmd/bub files already call loadConfigWithOverrides() to get config

## Constraints

- **Tech stack**: Go 1.25, lipgloss v2 with compat, cobra CLI
- **No TUI changes**: internal/ui, internal/app are untouched
- **ANSI width**: lipgloss adds invisible ANSI codes; account for this in fixed-width columns (top.go especially)

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| New package `cliout` at internal/cliout/ | Follows internal/ convention, avoids circular deps | ✓ Good |
| CLIStyles as value struct | Stateless after creation, safe to pass around | ✓ Good |
| No renderer needed (lipgloss v2 removed NewRenderer) | Styles are pure value types via lipgloss.NewStyle() | ✓ Good |
| Style before padding for tabwriter | lipgloss ANSI codes break width calculations | — Pending |

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd-transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd-complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state

---
*Last updated: 2026-04-11 after Phase 1 completion*

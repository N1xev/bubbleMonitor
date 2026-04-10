# Project Research Summary

**Project:** bubbleMonitor CLI Styling Layer
**Domain:** Go CLI output styling (non-interactive subcommands using lipgloss v2)
**Researched:** 2026-04-09
**Confidence:** HIGH

## Executive Summary

bubbleMonitor is a Go system monitor that already has a rich TUI built with BubbleTea and 28 carefully designed color themes. This milestone adds themed, colored output to its one-shot CLI subcommands (`bub status`, `bub health`, `bub doctor`, etc.) so that piped and terminal output feel visually consistent with the TUI experience. The domain is well-charted: lipgloss v2 provides exactly the primitives needed -- value-type styles, `compat.AdaptiveColor` bridging, `lipgloss.Width()` for ANSI-aware measurement, and `lipgloss.Fprintf` for automatic color downsampling on non-TTY outputs. No new external dependencies are required; everything is already in go.mod.

The recommended approach is a single new internal package (`internal/cliout`) that converts the existing `ThemePalette` into a `CLIStyles` value struct, plus render helpers for bars, aligned key-value pairs, color swatches, and threshold-colored values. Every command follows the same three-line preamble: load config, resolve palette, create CLIStyles. All styled output goes through `lipgloss.Fprintf` for pipe-safety. The critical architectural constraint is that `cliout` imports only the `ThemePalette` type from `internal/ui` -- never config, never app state -- keeping the styling layer purely presentational.

The key risks are ANSI-aware column alignment (tabwriter and fixed-width `fmt.Sprintf` both break with styled strings) and per-call style allocation in the `top.go` refresh loop. Both are preventable: use `lipgloss.Width()` for manual padding instead of tabwriter, and pre-build all styles once in `cliout.New()` rather than creating them inline. The research is thorough and consistent across all four files, with direct source-code verification of every pattern.

## Key Findings

### Recommended Stack

No new packages needed. The entire milestone uses lipgloss v2 and its bundled sub-packages, which are already in go.mod.

**Core technologies:**
- **lipgloss v2 (v2.0.2):** Style definitions, color rendering, `lipgloss.Width()` for ANSI-aware measurement, `lipgloss.Fprintf` for pipe-safe output -- the single foundation for all CLI styling.
- **lipgloss/v2/compat:** Bridges the existing `ThemePalette` (which stores colors as `compat.AdaptiveColor`) into v2 styles without modifying the TUI layer. AdaptiveColor implements `color.Color`, so it passes directly to `lipgloss.NewStyle().Foreground()`.
- **lipgloss/v2/table:** Optional, for tabular process display in `bub top` and `bub remote`. Handles column sizing with styled cells. Alternative: manual fixed-width columns with `lipgloss.Width()` padding.

**Explicitly NOT used:** `text/tabwriter` (ANSI codes break alignment), `bubbles/progress` (requires BubbleTea model), `glamour` (Markdown renderer, wrong abstraction), `lipgloss.NewRenderer()` (removed in v2).

### Expected Features

**Must have (table stakes):**
- Themed color on labels and values -- wiring the existing 28 theme palettes to CLI output
- Status-dependent coloring (OK/WARN/CRITICAL) for `health` and `status` commands
- Progress bars with threshold colors (green/yellow/red based on percentage)
- Themed check/fail/warn symbols in `doctor` output
- Dim/muted text for secondary info (uptime, paths, mount points)
- Bold for section headers and key values
- NO_COLOR and TERM=dumb respect
- Suppress color when piped (TTY detection)
- Consistent label alignment that accounts for ANSI escape code width
- Export stays machine-readable (no ANSI in JSON/CSV)

**Should have (differentiators):**
- 28+ theme-aware color palettes for CLI output (visual continuity with TUI)
- Color swatches in `bub themes list` (colored block characters, not text names)
- Threshold-colored CPU%/MEM% in the process table
- Large styled health score number

**Defer (v2+):**
- `--color` flag (always/auto/never) -- NO_COLOR + TTY detection covers 95% of cases
- Large styled health score number -- visual polish, not essential

### Architecture Approach

A three-layer architecture with strict dependency rules: config resolution stays in `cmd/bub`, style construction lives in a new `internal/cliout` package, and style application happens in each command file. The cliout package receives only a `ThemePalette` value and returns a `CLIStyles` value struct. It never imports config, app state, or data providers.

**Major components:**
1. **`internal/cliout/` (new)** -- CLIStyles struct, New() constructor, render helpers (BarColor, ScoreColor, RenderBar, KVAligned, RenderSwatch, PadRight). Depends only on `internal/ui` for the ThemePalette type.
2. **`internal/ui/styles.go` (untouched)** -- Existing ThemePalette definition, theme map, GetAppTheme resolution. No changes per project constraints.
3. **`cmd/bub/*.go` (modified in place)** -- Each command adds the three-line preamble (loadConfigWithOverrides + GetAppTheme + cliout.New), then replaces `fmt.Fprintf` with `lipgloss.Fprintf` for styled output.

**Key patterns:** Construct styles once per invocation, never in loops. Write through `lipgloss.Fprintf`, not `fmt.Fprintf`. Use `lipgloss.Width()` for column padding, never tabwriter or `fmt.Sprintf("%-Ns")` on styled strings.

### Critical Pitfalls

1. **ANSI escape codes break tabwriter alignment** -- tabwriter counts invisible escape bytes as column width. Use `lipgloss.Width()` + manual space padding instead. Highest-risk commands: `remote list` (uses tabwriter), `top` (uses `%-30s` format verbs).
2. **Per-call style creation in loops** -- Creating `lipgloss.NewStyle()` inside `top.go`'s refresh loop (every 2s, 20+ processes) causes unnecessary allocations. Pre-build all styles in `cliout.New()`, then pick among `s.OK`/`s.Warn`/`s.Critical` at render time.
3. **Using fmt.Fprintf instead of lipgloss.Fprintf for styled output** -- `fmt.Fprintf` writes raw ANSI bytes regardless of TTY status. Piped output gets escape-code garbage. `NO_COLOR` is ignored. Every styled output line must use `lipgloss.Fprintf`.
4. **compat.AdaptiveColor I/O timing** -- The compat layer queries stdin/stdout globally for light/dark detection. Resolve all colors at CLIStyles construction time, not per-render. Test with `echo "" | bub status` to catch edge cases.
5. **lipgloss tab-to-spaces conversion** -- lipgloss v2 converts tabs to 4 spaces by default. If styled text passes through tabwriter, tabs are already converted before tabwriter sees them, breaking alignment. Either avoid Render() before tabwriter or use `TabWidth(NoTabConversion)`.

## Implications for Roadmap

Based on research, suggested phase structure:

### Phase 1: Foundation -- cliout package
**Rationale:** All other phases depend on the cliout package existing. Must establish the core patterns (value-type styles, lipgloss.Fprintf, width-aware padding) before any command gets styled.
**Delivers:** `internal/cliout/` with CLIStyles struct, New() constructor, and all render helpers (BarColor, ScoreColor, RenderBar, KVAligned, PadRight, RenderSwatch). Unit tests.
**Addresses:** Table stakes foundation, pitfall #2 (per-call styles), pitfall #3 (fmt vs lipgloss writers)
**Avoids:** Pitfalls #1, #2, #3 by baking prevention into the package API from day one

### Phase 2: Simple Commands -- version, config, export
**Rationale:** Simplest commands first to validate the cliout API in real command files with minimal risk. Version is one line. Config is key-value. Export needs conditional (no-styling) path.
**Delivers:** Styled output for `version`, `config`, `export` commands.
**Uses:** CLIStyles.Label, Value, Bold, Dim
**Implements:** The three-line preamble pattern (loadConfig + GetAppTheme + cliout.New) that all subsequent commands will copy

### Phase 3: Monitoring Commands -- status, health, doctor
**Rationale:** The three most-used monitoring commands, where color adds the most value. These exercise threshold coloring, progress bars, and check-mark styles -- the core value proposition.
**Delivers:** Colored progress bars in status, status-dependent coloring in health, themed check/fail symbols in doctor.
**Uses:** BarColor, ScoreColor, RenderBar, CheckOK/CheckFail/CheckWarn styles
**Avoids:** Pitfall #4 (resolve colors once at construction)

### Phase 4: Label-Value Commands -- sysinfo
**Rationale:** sysinfo requires ANSI-aware label alignment (12-char label column), which is a step up in complexity from Phase 2-3 commands. By this point, KVAligned and PadRight are battle-tested.
**Delivers:** Styled sysinfo output with properly aligned labels across disk partitions and system sections.
**Uses:** KVAligned, Header style, section grouping

### Phase 5: Tabulated Commands -- top, remote
**Rationale:** Highest-risk commands due to column alignment with styled strings. top.go uses fixed-width format verbs; remote.go uses tabwriter. Both must replace these with width-aware patterns. Coming last means cliout helpers are proven and any edge cases are already discovered.
**Delivers:** Threshold-colored CPU%/MEM% in process table, styled headers, aligned columns.
**Uses:** PadRight, lipgloss.Width(), threshold styles in loop context
**Avoids:** Pitfalls #1 (tabwriter + ANSI), #5 (tab conversion), #2 (loop allocation)

### Phase 6: Theme Display -- themes list with swatches
**Rationale:** Color swatch rendering is a visual differentiator but isolated from the alignment concerns of other commands. Depends on RenderSwatch helper from Phase 1.
**Delivers:** Colored block swatches in `bub themes list`, active theme marker, dimmed inactive themes.
**Uses:** RenderSwatch, Active style, Background styling

### Phase Ordering Rationale

- Phase 1 is the hard dependency -- nothing else can proceed without cliout existing.
- Phases 2-6 progress from simplest to most complex output patterns. Each phase validates more of the cliout API. If a helper is wrong (e.g., PadRight miscalculates), it is discovered in Phase 2 (cheap fix) before Phase 5 (where alignment bugs are hardest to debug).
- top.go and remote.go are last because they have the trickiest alignment concerns (fixed-width format verbs and tabwriter respectively).
- The ordering matches the dependency graph: each phase builds on patterns proven in the previous phase.

### Research Flags

Phases likely needing deeper research during planning:
- **Phase 5 (top/remote):** Complex integration with existing tabwriter and fixed-width formatting. The lipgloss table package is an alternative to manual padding but needs evaluation for live-refresh performance in top.go. May warrant `/gsd-research-phase` during planning.
- **Phase 6 (themes):** Background-colored block rendering with `compat.AdaptiveColor` needs verification. The swatch rendering pattern (Foreground+Background matching for solid blocks) should work but has not been tested against all 28 themes.

Phases with standard patterns (skip research-phase):
- **Phase 1 (foundation):** Well-documented lipgloss v2 patterns. All code snippets verified against official docs and pkg.go.dev.
- **Phase 2 (simple commands):** Straightforward key-value styling. No alignment complexity.
- **Phase 3 (monitoring):** Standard threshold-coloring pattern. `renderBar` already exists in codebase.
- **Phase 4 (sysinfo):** Standard label-value pattern with `KVAligned` helper.

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Stack | HIGH | lipgloss v2 is the only dependency needed. Already in go.mod. All patterns verified against official docs (pkg.go.dev v2.0.2, GitHub README). No version conflicts. |
| Features | HIGH | Feature list derived from direct codebase analysis of all 10 command files. Every feature maps to existing code that needs styling. Anti-features clearly bounded. |
| Architecture | HIGH | Three-layer separation is standard Go. Import graph verified against actual source files. No circular dependencies possible given the constraints. |
| Pitfalls | HIGH | Top 5 pitfalls all verified against lipgloss v2 source code and project source files. tabwriter+ANSI issue is widely documented. |

**Overall confidence:** HIGH

### Gaps to Address

- **lipgloss table package performance in live-refresh:** The table package may have overhead for top.go's 2-second refresh cycle. If profiling shows issues, fall back to manual fixed-width columns with `lipgloss.Width()`. Decide during Phase 5 planning.
- **Windows ANSI support:** `lipgloss.EnableLegacyWindowsANSI(os.Stdout)` must be called early. The call is a no-op on non-Windows, but Windows testing is needed to confirm it works with Cobra's output writers. Add to Phase 1 as a one-liner in root.go.
- **SSH session edge cases:** `compat.AdaptiveColor` queries stdin for light/dark detection. Over SSH with no stdin (e.g., `bub health` via remote execution), this may default incorrectly. Test with `echo "" | bub status` during Phase 1 validation.

## Sources

### Primary (HIGH confidence)
- lipgloss v2 GitHub README (github.com/charmbracelet/lipgloss) -- official, verified 2026-03-17
- lipgloss v2 pkg.go.dev (pkg.go.dev/charm.land/lipgloss/v2) -- v2.0.2, published 2026-03-11
- lipgloss v2 compat pkg.go.dev -- AdaptiveColor API, color.Color interface
- lipgloss v2 table pkg.go.dev -- table API, StyleFunc, HeaderRow constant
- Project go.mod -- confirmed lipgloss v2.0.0 and bubbletea v2.0.0 in dependencies
- Project source: internal/ui/styles.go -- ThemePalette definition, 28 theme palettes
- Project source: cmd/bub/*.go -- all 10 command files analyzed for output patterns

### Secondary (MEDIUM confidence)
- NO_COLOR spec (no-color.org) -- convention, widely adopted
- WebSearch findings on tabwriter + ANSI alignment -- widely reported issue, consistent across sources
- CLI conventions: GNU coreutils --color, gh CLI, docker CLI, kubectl -- reference for expected behavior

---
*Research completed: 2026-04-09*
*Ready for roadmap: yes*

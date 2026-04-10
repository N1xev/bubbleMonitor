# Feature Landscape

**Domain:** CLI output styling for a Go system monitor (non-interactive/one-shot subcommands)
**Researched:** 2026-04-09
**Confidence:** HIGH (based on codebase analysis + well-established CLI conventions)

## Table Stakes

Features users expect from a polished CLI tool. Missing any of these makes the tool feel unfinished or hostile to scripting.

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| **Themed color on labels and values** | Every modern CLI (`gh`, `docker`, `kubectl` wrappers) colors its output. Plain `fmt.Fprintf` text looks broken by comparison when users have come from the TUI. | Low | Use `ThemePalette.Primary` for labels, `Text` for values. The palette already exists in `internal/ui/styles.go` with 28 themes. |
| **Status-dependent coloring** (OK/HIGH/CRITICAL) | `bub health` and `bub status` output states like "CRITICAL" in plain text -- users expect red for bad, green for good. This is universal in monitoring tools. | Low | Map threshold bands to `Success`/`Warning`/`Alert` palette colors. Health.go already computes the status strings. |
| **Progress bars with threshold colors** | `status.go` already renders block-character bars but in monochrome. Every system monitor colors bars by severity (green/yellow/red). | Low | `renderBar()` already exists -- wrap filled characters in color. Use `BarColor(pct)` helper: <60 green, <80 yellow, >=80 red. |
| **Themed check/fail/warn symbols** | `doctor.go` uses Unicode check/cross/warning symbols with no color. All diagnostic CLIs (`rustup doctor`, `brew doctor`) color these. | Low | Style the existing symbols: green check, red cross, yellow warning using `Success`/`Alert`/`Warning` palette colors. |
| **Dim/muted text for secondary info** | Uptime strings, file paths, mount points, "refreshing every Xs" -- secondary data should recede visually. | Low | Single `Dim` style using `ThemePalette.Muted`. Already have the color. |
| **Bold for section headers and key values** | Health score, process table header, system info labels -- key information should stand out from values. | Low | `Bold` style using `Text` palette color with bold attribute. |
| **Respect NO_COLOR / TERM=dumb** | The [no-color.org](https://no-color.org) convention is adopted by every major CLI tool. Ignoring it breaks scripting and accessibility. | Low | Check `NO_COLOR` env var and `TERM=dumb` at CLI styles initialization. When set, return unstyled output (no ANSI codes). |
| **Suppress color when piped** (TTY detection) | Piping `bub status` to `grep` or `less` should not spray ANSI escape codes. Standard expectation for any CLI. | Low | Use `term.IsTerminal()` from `charm.land/x/term` (already an indirect dependency). When stdout is not a TTY, skip styling. This is the `--color=auto` default. |
| **Consistent label alignment** | `sysinfo.go` and `health.go` use manual padding like `%-12s`. Styled labels break alignment because ANSI codes are invisible width. | Medium | Use `lipgloss.Width()` to measure visible width and pad accordingly. Or use fixed-width label style via `style.Width(12)`. |
| **Export stays machine-readable** | `bub export` outputs JSON/CSV. These must NEVER get ANSI color codes -- they are consumed by scripts and tools. | Low | Export commands write to `cmd.OutOrStdout()` which may be piped; but explicitly: do not apply CLIStyles to export.go output at all. Structural requirement, not a styling choice. |

## Differentiators

Features that set `bub` apart from other system monitors. Not expected by default, but make the tool feel premium.

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| **28+ theme-aware color palettes** | Most CLIs hardcode a few colors. `bub` already has 28 themes with carefully designed palettes. Using the same palette for CLI output as the TUI creates visual continuity. This is a significant differentiator. | Low | Palette data already exists. The work is wiring `ThemePalette` -> `CLIStyles` mapping. |
| **Active theme indicator with color swatches** | `bub themes list` showing actual colored blocks per theme (not just text names) is visually striking and helps users pick a theme without switching. | Medium | Current `renderColorSwatch()` returns a text string. Replace with actual ANSI-colored blocks using the palette's Primary/Success/Warning/Alert colors rendered as `███` blocks. Requires lipgloss rendering in themes.go. |
| **Threshold-colored CPU%/MEM% in process table** | `top.go` showing processes with red-tinted CPU% for high-usage processes provides at-a-glance triage. Most `top`-like tools do NOT color their CLI output. | Medium | Color the percentage value based on threshold: <60 normal, <80 warning, >=80 alert. Must account for ANSI width in column alignment (existing tabwriter pattern). The `top.go` table uses fixed-width `%8s` columns -- styled values will break alignment unless handled. |
| **`--color` flag (always/auto/never)** | Giving explicit control over color output goes beyond `NO_COLOR` compliance. Follows the GNU convention (`grep --color`, `ls --color`). | Low | Add `--color` persistent flag on root command. `auto` = detect TTY (default), `always` = force color, `never` = no color. Resolves before CLIStyles creation. |
| **Health score with large styled number** | `bub health` showing the score as a large, color-coded number (green >=70, yellow >=50, red <50) makes the single most important datum unmissable. | Low-Medium | Render the score with Bold + sized style + threshold color. Single standout element in the output. |

## Anti-Features

Features to explicitly NOT build. These would degrade the tool or create maintenance burden.

| Anti-Feature | Why Avoid | What to Do Instead |
|--------------|-----------|-------------------|
| **Borders and boxes around CLI output** | `lipgloss` makes borders easy, but bordered CLI output looks over-designed for one-shot commands. Borders belong in the TUI, not in `bub status` or `bub health`. Users pipe this output. | Use indentation (2-space prefix, which already exists) and color for visual structure. No box-drawing characters in CLI output. |
| **Animated spinners or progress animations** | These require TTY control and are inappropriate for one-shot commands that complete in <1 second. The TUI already has animations. CLI commands should be instant. | Keep CLI output static. If a command is slow (doctor SSH checks), just print the result when done. |
| **Markdown rendering for CLI help** | `glamour` renders markdown in terminals, but `fang` already handles help formatting. Adding glamour for command output would be scope creep and slow. | Use lipgloss styles for all output. Keep plain text structure. |
| **Table borders (ASCII art tables)** | `tablewriter` and `go-pretty` produce bordered ASCII tables. These look cluttered for system monitor output. The existing whitespace-separated approach (tabwriter for top.go/remote.go, manual alignment for others) is cleaner. | Keep borderless tables. Use color and bold headers to distinguish structure. |
| **New theme palettes or color choices** | The project already has 28 themes. The milestone scope is wiring existing themes to CLI output, not designing new colors. | Use existing `ThemePalette` colors exclusively. |
| **Changing output structure/format** | Only add color/styling to existing output. Do not rearrange, add, or remove lines. The output format is a contract with users. | Style-only changes. Same `fmt.Fprintf` calls, wrapped with styled strings. |
| **Per-command theme overrides** | Each command does not need its own theme setting. The global `--theme` flag and config already handle this. | All commands read the same theme from `loadConfigWithOverrides()`. One theme per invocation. |
| **ANSI color in CSV/JSON export** | Export output is machine-readable. Color would corrupt the data format. | `export.go` is explicitly excluded from styling. |
| **Custom progress bar characters per theme** | Each theme does not need its own block characters. The `█`/`░` characters work universally. | Use same bar characters across all themes. Only the color changes. |

## Feature Dependencies

```
loadCLIStyles() (core helper)
  ├── requires: loadConfigWithOverrides() (existing)
  ├── requires: ui.GetAppTheme() (existing)
  ├── requires: CLIStyles struct definition
  └── produces: CLIStyles used by all styled commands

NO_COLOR / TTY detection (environment awareness)
  └── feeds into: loadCLIStyles() -- returns no-op styles when color disabled

--color flag
  └── feeds into: loadCLIStyles() -- overrides TTY detection

CLIStyles struct
  ├── Label, Value, Bold, Dim, Header (base styles)
  ├── OK, Warn, Critical (status styles)
  ├── CheckOK, CheckFail, CheckWarn (doctor symbols)
  ├── BarColor(pct float64) (progress bar color helper)
  └── ScoreColor(score int) (health score color helper)

status.go styling
  ├── requires: CLIStyles.Label, CLIStyles.BarColor
  └── uses: renderBar() with color (existing function, add color wrapper)

health.go styling
  ├── requires: CLIStyles.OK/Warn/Critical, CLIStyles.ScoreColor
  └── uses: status-dependent coloring on each metric line

doctor.go styling
  ├── requires: CLIStyles.CheckOK/CheckFail/CheckWarn, CLIStyles.Dim
  └── uses: existing check/fail/warn symbols with color

top.go styling
  ├── requires: CLIStyles.Header, CLIStyles.Dim, CLIStyles.BarColor
  ├── requires: ANSI width handling for tabwriter compatibility
  └── depends on: understanding lipgloss.Width() for column alignment

sysinfo.go styling
  ├── requires: CLIStyles.Label, CLIStyles.Value, CLIStyles.Header, CLIStyles.Dim
  └── uses: label-value pairs with aligned styled labels

themes.go styling
  ├── requires: CLIStyles.Active, CLIStyles.Dim
  └── requires: lipgloss rendering for color swatch blocks

version.go styling
  └── requires: CLIStyles.Label, CLIStyles.Dim (minimal)

config.go styling
  └── requires: CLIStyles.Label, CLIStyles.Bold (minimal -- mostly JSON output)

remote.go styling
  ├── requires: CLIStyles.Header, CLIStyles.Dim
  └── uses: tabwriter with styled header (same alignment concern as top.go)

export.go styling
  └── NONE -- explicitly excluded from styling
```

## MVP Recommendation

Prioritize in this order:

1. **CLIStyles struct + loadCLIStyles() helper** -- Foundation. Everything else depends on this. Create the struct with all style fields, the `BarColor()`/`ScoreColor()` helpers, and the `loadCLIStyles()` function that bridges config -> theme -> CLIStyles.

2. **NO_COLOR / TTY detection** -- Must exist before any styling ships. Without this, the tool breaks in pipes and annoys users who set NO_COLOR.

3. **status.go + health.go + doctor.go styling** -- The three monitoring commands. These are the most-used commands and benefit most from color. Status gets colored bars, health gets status-dependent coloring and score coloring, doctor gets themed check/fail/warn symbols.

4. **sysinfo.go + top.go styling** -- System information and process table. These require careful ANSI width handling for alignment, so they come after the simpler commands.

5. **themes.go color swatches** -- Visual differentiator. Replaces the text-only swatch with actual colored blocks.

6. **version.go + config.go + remote.go styling** -- Minimal styling (bold labels, dim secondary info). Low effort, rounds out the polish.

Defer:
- **`--color` flag**: Nice-to-have. `NO_COLOR` + TTY detection covers 95% of use cases. Add the flag in a follow-up.
- **Large styled health score number**: Visual polish that can wait. The colored score text is sufficient for MVP.

## Complexity Assessment

| Command | Styling Complexity | Key Challenge |
|---------|-------------------|---------------|
| status.go | Low | Wrap existing `renderBar()` output in color. Straightforward. |
| health.go | Low | Status strings already computed. Map to colors. |
| doctor.go | Low | Color existing symbols. Dim paths. |
| version.go | Low | One line, minimal styling. |
| config.go | Low | Mostly JSON output. Style the "Set X = Y" confirmation only. |
| sysinfo.go | Medium | Label-value pairs with ANSI width alignment. 12-char label column. |
| top.go | Medium-High | Fixed-width columns + ANSI codes + tabwriter. Must handle `lipgloss.Width()` carefully or style only non-padded fields. Live refresh mode adds complexity. |
| themes.go | Medium | Color swatch rendering. Active theme marker. |
| remote.go | Medium | Same tabwriter + ANSI alignment concern as top.go. |
| export.go | None | No styling. Machine-readable output. |

## Sources

- Codebase analysis: `cmd/bub/*.go`, `internal/ui/styles.go`, `go.mod`
- CLI conventions: GNU coreutils `--color` convention, [no-color.org](https://no-color.org)
- Go libraries: `charm.land/lipgloss/v2`, `charm.land/x/term`, `text/tabwriter`
- Reference tools: `gh` (GitHub CLI), `docker` CLI, `kubectl`, `brew doctor`, `rustup doctor`

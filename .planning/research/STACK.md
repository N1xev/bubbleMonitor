# Technology Stack

**Project:** bubbleMonitor CLI Styling
**Researched:** 2026-04-09
**Scope:** Non-interactive CLI output styling with lipgloss v2 (not TUI)

## Recommended Stack

### Core Styling

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `charm.land/lipgloss/v2` | v2.0.2 | Style definitions, color, rendering | Already in go.mod; pure value-type styles; automatic color downsampling for non-TTY; `lipgloss.Width()` for measuring visible cell width |
| `charm.land/lipgloss/v2/compat` | (bundled) | AdaptiveColor compatibility with v1 theme palette | Already used by `internal/ui/styles.go`; ThemePalette fields are `compat.AdaptiveColor` which implement `color.Color` -- pass directly to `lipgloss.NewStyle().Foreground()` |
| `charm.land/lipgloss/v2/table` | (bundled) | Tabular process display (`bub top`, `bub ps`) | Native lipgloss table with `StyleFunc` per cell; handles column sizing internally; supports colored columns without tabwriter width issues |

### Output Functions

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `lipgloss.Fprint` / `lipgloss.Fprintf` | v2.0.2 | Write styled output to `cmd.OutOrStdout()` | Drop-in replacement for `fmt.Fprintf`; automatically strips ANSI codes when output is not a TTY (piped to file/log) -- zero extra code for pipe-safety |
| `lipgloss.Sprint` / `lipgloss.Sprintf` | v2.0.2 | Build styled strings for composition | For assembling styled fragments before final write |

### NOT Recommended

| Technology | Why Not |
|------------|---------|
| `text/tabwriter` | ANSI escape codes from lipgloss are invisible to tabwriter; it counts escape bytes as visible width, breaking column alignment. Use lipgloss `table` package or manual fixed-width `fmt.Sprintf` with `lipgloss.Width()` instead |
| `lipgloss.NewRenderer()` | Removed in v2. The v1 pattern of creating a renderer per-writer does not exist. In v2, use `lipgloss.HasDarkBackground(os.Stdin, os.Stdout)` + `lipgloss.LightDark()` for standalone adaptive color, or just pass `compat.AdaptiveColor` directly (it already resolves at render time) |
| `github.com/charmbracelet/bubbles/progress` | Designed for BubbleTea TUI programs with a `tea.Model`; not suitable for one-shot CLI output. Build progress bars with styled Unicode block chars (`#`/`-` or `\u2588`/`\u2591`) rendered through `lipgloss.Style.Render()` |
| `github.com/charmbracelet/glamour` | Markdown renderer; heavyweight for this use case. We need styled key-value pairs and bars, not rendered Markdown |
| `fmt.Fprintf` for styled output | Use `lipgloss.Fprint`/`lipgloss.Fprintf` instead. They automatically handle color downsampling when stdout is not a TTY (piped to file, CI, etc.) |

## Key Patterns

### 1. Bridging ThemePalette to CLI Styles

The existing `ThemePalette` stores colors as `compat.AdaptiveColor`. Since `compat.AdaptiveColor` implements `color.Color` (has an `RGBA()` method), it can be passed directly to lipgloss style methods:

```go
import (
    "charm.land/lipgloss/v2"
    "charm.land/lipgloss/v2/compat"
    "github.com/N1xev/bubbleMonitor/internal/ui"
)

// CLIStyles holds pre-built styles for CLI output
type CLIStyles struct {
    Label   lipgloss.Style  // dim, muted color for field names
    Value   lipgloss.Style  // bright, text color for field values
    Bold    lipgloss.Style  // bold for headers
    OK      lipgloss.Style  // green/success for healthy states
    Warn    lipgloss.Style  // yellow/warning for degraded
    Critical lipgloss.Style // red/alert for failures
    Dim     lipgloss.Style  // faint for secondary info (uptime, etc.)
    Header  lipgloss.Style  // bold + primary for section headers
    Active  lipgloss.Style  // bold + primary for active theme marker
    CheckOK lipgloss.Style  // success-colored check mark
    CheckFail lipgloss.Style // alert-colored X mark
    CheckWarn lipgloss.Style // warning-colored warning mark
}

func NewCLIStyles(palette ui.ThemePalette) CLIStyles {
    return CLIStyles{
        Label:    lipgloss.NewStyle().Foreground(palette.Muted),
        Value:    lipgloss.NewStyle().Foreground(palette.Text),
        Bold:     lipgloss.NewStyle().Bold(true).Foreground(palette.Text),
        OK:       lipgloss.NewStyle().Foreground(palette.Success),
        Warn:     lipgloss.NewStyle().Foreground(palette.Warning),
        Critical: lipgloss.NewStyle().Foreground(palette.Alert),
        Dim:      lipgloss.NewStyle().Faint(true).Foreground(palette.Muted),
        Header:   lipgloss.NewStyle().Bold(true).Foreground(palette.Primary),
        Active:   lipgloss.NewStyle().Bold(true).Foreground(palette.Primary),
        CheckOK:  lipgloss.NewStyle().Foreground(palette.Success),
        CheckFail: lipgloss.NewStyle().Foreground(palette.Alert),
        CheckWarn: lipgloss.NewStyle().Foreground(palette.Warning),
    }
}
```

**Confidence: HIGH** -- Verified against official lipgloss v2 README and pkg.go.dev docs. `compat.AdaptiveColor` implements `color.Color` interface, which is exactly what `lipgloss.Style.Foreground()` accepts.

### 2. Threshold-Colored Progress Bars

Build progress bars as styled Unicode block characters. The color changes based on percentage thresholds:

```go
func BarColor(pct float64, palette ui.ThemePalette) compat.AdaptiveColor {
    switch {
    case pct >= 80:
        return palette.Alert
    case pct >= 60:
        return palette.Warning
    default:
        return palette.Success
    }
}

func (s CLIStyles) RenderBar(pct float64, width int, palette ui.ThemePalette) string {
    filled := int(pct / 100.0 * float64(width))
    if filled > width { filled = width }

    barColor := BarColor(pct, palette)
    barStyle := lipgloss.NewStyle().Foreground(barColor)

    filledPart := strings.Repeat("\u2588", filled)   // full block
    emptyPart := strings.Repeat("\u2591", width-filled) // light shade

    return barStyle.Render(filledPart) + s.Dim.Render(emptyPart)
}
```

**Confidence: HIGH** -- Simple string composition with lipgloss styles. The `\u2588`/`\u2591` characters are already used in the codebase (`renderBar` in status.go uses similar approach with plain strings).

### 3. Key-Value Output (sysinfo, status, health)

For labeled output like `  CPU    45.2%`, compose label + value as styled fragments:

```go
func (s CLIStyles) KV(label, value string) string {
    return s.Label.Render(label) + " " + s.Value.Render(value)
}
```

**Do NOT use tabwriter.** Instead, use fixed-width label padding:

```go
func (s CLIStyles) KVAligned(label string, labelWidth int, value string) string {
    padded := fmt.Sprintf("%-*s", labelWidth, label)
    return s.Label.Render(padded) + " " + s.Value.Render(value)
}
```

**Why not tabwriter:** lipgloss embeds ANSI escape codes in the rendered string. tabwriter counts these invisible bytes as part of the column width, causing misalignment. Fixed-width formatting with `fmt.Sprintf("%-*s", ...)` avoids this entirely because the padding is applied to the raw label text before styling.

**Confidence: HIGH** -- This is the documented lipgloss approach. The official README explicitly warns about this by providing `lipgloss.Width()` for measuring visible cell width of styled strings.

### 4. Tabular Data for `bub top` (Process Table)

Use the lipgloss `table` sub-package for the process list. It handles column width calculation internally and supports per-cell styling:

```go
import "charm.land/lipgloss/v2/table"

func renderProcessTable(procs []procEntry, styles CLIStyles, palette ui.ThemePalette) string {
    t := table.New().
        Headers("PID", "NAME", "CPU%", "MEM%", "STATUS").
        Border(lipgloss.HiddenBorder()).  // no visible border for top-like display
        StyleFunc(func(row, col int) lipgloss.Style {
            if row == table.HeaderRow {
                return styles.Bold
            }
            // Color CPU% and MEM% based on threshold
            if col == 2 || col == 3 { // CPU% or MEM% columns
                val := /* parse the percentage from procs[row-1] */
                switch {
                case val >= 80:
                    return lipgloss.NewStyle().Foreground(palette.Alert)
                case val >= 60:
                    return lipgloss.NewStyle().Foreground(palette.Warning)
                default:
                    return lipgloss.NewStyle().Foreground(palette.Text)
                }
            }
            return lipgloss.NewStyle().Foreground(palette.Text)
        })

    for _, p := range procs {
        t.Row(fmt.Sprintf("%d", p.Pid), p.Name,
            util.FastPercent1(p.Cpu),
            util.FastPercent1(p.Mem),
            p.Status)
    }

    return t.Render()
}
```

**Alternative approach (current pattern):** If the table package feels heavyweight for the live-refresh `bub top` command, keep the current `fmt.Fprintf` with fixed-width columns but style only non-padded fields (the numeric values). Style the header row with bold. Apply threshold colors to CPU%/MEM% values. This is simpler and preserves the existing `\033[H\033[2J` screen-clear approach.

**Confidence: HIGH** -- Table package verified in official lipgloss README and pkg.go.dev. `table.HeaderRow` constant and `StyleFunc` pattern confirmed.

### 5. Color Swatches for `bub themes list`

Render color swatches as background-styled Unicode full-block characters:

```go
func renderColorSwatch(p ui.ThemePalette) string {
    colors := []compat.AdaptiveColor{
        p.Primary, p.Secondary, p.Success, p.Warning, p.Alert,
    }
    var b strings.Builder
    for _, c := range colors {
        s := lipgloss.NewStyle().Background(c).Foreground(c).Render(" ")
        b.WriteString(s)
    }
    return b.String()
}
```

Each color renders as a single cell with matching foreground and background, producing a colored square block. Concatenating them creates a visual swatch strip.

**Confidence: HIGH** -- Uses standard lipgloss background styling. The `Render(" ")` produces a single ANSI-colored space character.

### 6. Adaptive Color Handling (CLI Context)

The project's themes all use `compat.AdaptiveColor` with identical Light/Dark values (via `makeColor` and `makeTTYColor`). This means the current themes do NOT actually differentiate between light and dark terminal backgrounds -- each color resolves to the same value regardless.

For CLI output, `compat.AdaptiveColor` resolves automatically when used with `lipgloss.Style.Render()` because it implements `color.Color` and its `RGBA()` method checks the global `compat.HasDarkBackground` variable (initialized from `os.Stdin`/`os.Stdout`).

**Important nuance:** The lipgloss v2 docs note that `compat` is "not recommended for new code" because it uses global I/O state. However, since the project already uses it extensively in `internal/ui/styles.go` and the ThemePalette is locked to this type, we should continue using it in the CLI layer. The alternative (using `lipgloss.HasDarkBackground()` + `lipgloss.LightDark()`) would require changing the ThemePalette type, which is out of scope.

**Confidence: HIGH** -- Verified from compat package docs on pkg.go.dev. The `compat.AdaptiveColor.RGBA()` method checks `compat.HasDarkBackground` (a package-level var) and returns the appropriate color value.

### 7. Pipe-Safe Output (Non-TTY Handling)

Use `lipgloss.Fprint`/`lipgloss.Fprintf` instead of `fmt.Fprintf` for all styled output. These functions automatically strip ANSI codes when the output writer is not a TTY. This means:

- `bub status | cat` -- outputs plain text, no escape codes
- `bub health > report.txt` -- clean text in the file
- `bub doctor` in CI -- readable without color codes

```go
// Instead of:
fmt.Fprintf(cmd.OutOrStdout(), "  CPU: %s\n", styledOutput)

// Use:
lipgloss.Fprint(cmd.OutOrStdout(), "  CPU: ", styledOutput, "\n")
// Or for formatted strings:
lipgloss.Fprintf(cmd.OutOrStdout(), "  CPU: %s\n", styledOutput)
```

**Confidence: HIGH** -- Documented in lipgloss v2 README "Color Downsampling" section and the `Writer` variable docs.

## Installation

No new packages needed. Everything is already in go.mod:

```
charm.land/lipgloss/v2 v2.0.0
```

The `table`, `list`, `tree`, and `compat` sub-packages are bundled within the lipgloss module -- no separate installation required.

## Alternatives Considered

| Category | Recommended | Alternative | Why Not |
|----------|-------------|-------------|---------|
| Tabular display | lipgloss `table` or fixed-width `fmt.Sprintf` | `text/tabwriter` | ANSI escape codes break tabwriter width calculations; tabwriter counts invisible escape bytes as column content |
| Progress bars | Unicode chars + lipgloss styles | `bubbles/progress` | bubbles progress requires BubbleTea model lifecycle; overkill for one-shot CLI output |
| Rich output | lipgloss styles + `lipgloss.Fprint` | `glamour` | Glamour renders Markdown; we need programmatic style control, not Markdown-to-terminal conversion |
| Color resolution | `compat.AdaptiveColor` (existing) | `lipgloss.LightDark()` (new v2 API) | Would require changing ThemePalette type; out of scope. The compat layer works fine for our use case |
| Output writing | `lipgloss.Fprint` / `lipgloss.Fprintf` | `fmt.Fprintf` | `fmt.Fprintf` does not handle color downsampling; ANSI codes leak into piped output |

## Sources

- lipgloss v2 GitHub README: https://github.com/charmbracelet/lipgloss -- official, current (verified 2026-03-17)
- lipgloss v2 pkg.go.dev: https://pkg.go.dev/charm.land/lipgloss/v2 -- v2.0.2, published 2026-03-11
- lipgloss compat pkg.go.dev: https://pkg.go.dev/charm.land/lipgloss/v2/compat -- official API reference
- lipgloss table pkg.go.dev: https://pkg.go.dev/charm.land/lipgloss/v2/table -- official table API
- Project go.mod: confirmed `charm.land/lipgloss/v2 v2.0.0` and `charm.land/bubbletea/v2 v2.0.0` in dependencies
- Project source: `internal/ui/styles.go` -- confirmed ThemePalette uses `compat.AdaptiveColor` fields

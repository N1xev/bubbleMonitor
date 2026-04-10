# Architecture Patterns: CLI Styling Layer (cliout)

**Project:** bubbleMonitor
**Researched:** 2026-04-10
**Scope:** Bridge between ThemePalette and CLI subcommand output

---

## Recommended Architecture

```
                        CONFIG LAYER
                    internal/config
                          |
                    AppConfig.Theme
                    CustomThemeConfig
                          |
                          v
        +-------------------------------------------+
        |  RESOLUTION (lives in cmd/bub, NOT cliout)|
        |                                           |
        |  cfg := loadConfigWithOverrides()         |
        |  palette := ui.GetAppTheme(               |
        |      cfg.Theme, cfg.CustomTheme)          |
        +-------------------------------------------+
                          |
                    ThemePalette
                          |
                          v
               +----------------------+
               |   CLIOUT PACKAGE     |
               |                      |
               |  New(palette) ->     |
               |    CLIStyles         |
               |                      |
               |  Render helpers:     |
               |  - BarColor(pct)    |
               |  - ScoreColor(pct)  |
               |  - RenderBar(...)   |
               |  - KVAligned(...)   |
               +----------------------+
                          |
                   CLIStyles instance
                   + helper functions
                          |
                          v
        +-------------------------------------------+
        |  COMMAND FILES (cmd/bub/*.go)             |
        |                                           |
        |  s := cliout.New(palette)                 |
        |  lipgloss.Fprintf(out, "%s %s\n",         |
        |      s.Label.Render("CPU"),               |
        |      s.Value.Render("45.2%"))             |
        +-------------------------------------------+
                          |
                          v
                     cmd.OutOrStdout()
                     (via lipgloss.Fprintf)
```

### Why This Shape

Three layers, each with a single job:

1. **Config resolution** stays in `cmd/bub` (it already lives there in `loadConfigWithOverrides`). The cliout package never touches config.
2. **Style construction** lives in a new `cliout` package that takes a `ThemePalette` and returns a value struct of pre-built styles plus render helpers. It depends only on `internal/ui` (for the `ThemePalette` type) and `lipgloss`.
3. **Style application** happens in each command file, which calls cliout once to get styles, then uses those styles for all output.

---

## Component Boundaries

### Component 1: `internal/ui/styles.go` (untouched)

| Aspect | Detail |
|--------|--------|
| Responsibility | Define ThemePalette, store theme map, resolve theme by name |
| Imports from | `internal/config` (for `CustomThemeConfig`), `lipgloss`, `lipgloss/compat` |
| Imported by | `internal/ui/layout.go`, `cmd/bub/themes.go` (already), and the new `cliout` package |
| Changes | None. Out of scope per PROJECT.md |

### Component 2: New package `cliout` at `internal/cliout/`

| Aspect | Detail |
|--------|--------|
| Responsibility | Convert ThemePalette into CLI-ready styles and render helpers |
| Imports from | `internal/ui` (for `ThemePalette` type only), `lipgloss/v2`, `lipgloss/v2/compat`, `fmt`, `strings` |
| Imported by | `cmd/bub/*.go` command files |
| Does NOT import | `internal/config`, `internal/app`, `internal/data`, `internal/msg`, `internal/provider` |

**Why `internal/cliout/` and NOT `cmd/bub/cliout/`:**

The PROJECT.md key decisions mention "New package `cliout` in `cmd/bub/`" but this is the wrong location for three reasons:

1. **Import path clarity.** Packages under `cmd/bub/` are conventionally the `main` package. A sub-package there would have the path `github.com/N1xev/bubbleMonitor/cmd/bub/cliout`, which is awkward and suggests it is part of the CLI binary rather than a reusable library.
2. **Separation of concerns.** The cliout package is a styling utility, not command logic. It belongs with other internal libraries (`internal/util`, `internal/config`), not alongside cobra command definitions.
3. **Test isolation.** `internal/cliout/` can have its own `cliout_test.go` that imports only the package under test and lipgloss. Tests under `cmd/bub/` would need the entire `main` package context.

**Alternative considered: `cmd/bub/cliout/`.** Rejected because it mixes library code with the `main` package directory, making the import path longer and the project structure harder to reason about.

### Component 3: `cmd/bub/*.go` command files (modified in place)

| Aspect | Detail |
|--------|--------|
| Responsibility | Load config, resolve theme, create CLIStyles, produce styled output |
| Imports from | `internal/cliout` (new), `internal/ui` (for `GetAppTheme` resolution), `internal/config`, `lipgloss` (for `lipgloss.Fprintf`) |
| Changes | Replace `fmt.Fprintf(out, ...)` with `lipgloss.Fprintf(out, ...)` for styled output; add cliout construction at the top of each command's RunE function |

### Cross-Component Dependency Rules

```
internal/config  --->  (no internal deps)
internal/ui/styles.go  --->  internal/config
internal/cliout  --->  internal/ui (ThemePalette type only)
cmd/bub  --->  internal/cliout, internal/ui, internal/config
```

**Key invariant:** `internal/cliout` never imports `internal/config`. It receives a `ThemePalette` value and nothing else. This keeps the style package agnostic about where themes come from.

---

## Data Flow

### Flow 1: Theme Resolution (per command invocation)

```
1. User runs: bub status --theme nord
2. Cobra parses flags into package vars (cfgTheme = "nord")
3. Command's RunE calls loadConfigWithOverrides()
   - Loads config.json (or defaults)
   - Applies --theme, --refresh-rate, --history-length overrides
   - Returns AppConfig{Theme: "nord", CustomTheme: nil}
4. Command resolves palette:
   palette := ui.GetAppTheme(cfg.Theme, cfg.CustomTheme)
   - ui.GetAppTheme("nord", nil) -> ui.GetTheme("nord") -> themes["nord"]
   - Returns ThemePalette{Primary: makeColor("#82AAFF"), ...}
5. Command creates styles:
   s := cliout.New(palette)
   - cliout.New builds CLIStyles struct with pre-built lipgloss.Style values
   - All color resolution happens once here
6. Command writes output using s.Label, s.Value, etc.
```

This flow already partially exists. `themes.go` already calls `ui.GetTheme(name)` (line 45) and `loadConfigWithOverrides()` (line 34). The new pattern just standardizes it: every command that needs styled output calls `loadConfigWithOverrides()` + `ui.GetAppTheme()` + `cliout.New()`.

### Flow 2: Config-Optional Commands

Some commands do not currently load config (`status.go`, `top.go`, `sysinfo.go`, `version.go`). For these, the pattern is:

```
1. loadConfigWithOverrides() // already handles errors by returning defaults
2. palette := ui.GetAppTheme(cfg.Theme, cfg.CustomTheme)
3. s := cliout.New(palette)
```

`loadConfigWithOverrides()` already falls back to `DefaultConfig()` on error (root.go line 126), so no new error handling is needed. The default theme is "horizon".

### Flow 3: Tabulated Output (top.go, remote.go)

```
1. Resolve styles as above
2. For each row:
   a. Build styled fragments: name := s.Value.Render(proc.Name)
   b. Measure visual width: w := lipgloss.Width(name)
   c. Pad to column width: padded := name + strings.Repeat(" ", colWidth - w)
   d. Write via lipgloss.Fprintf(out, "%s%s\n", styledPID, padded)
```

Do NOT use tabwriter for styled cells. Do NOT use `fmt.Sprintf("%-30s", styled)` (byte-length padding). Use `lipgloss.Width()` + manual space padding.

For `remote.go` which currently uses tabwriter: either (a) replace tabwriter with lipgloss table package, or (b) keep tabwriter but style only non-aligned cells and use `lipgloss.Width()` compensation for styled ones. Option (a) is cleaner.

### Flow 4: Threshold-Based Dynamic Colors

```
1. Pre-built styles: s.OK, s.Warn, s.Critical (from palette.Success, .Warning, .Alert)
2. At render time, pick the right style:
   func thresholdStyle(pct float64, s cliout.CLIStyles) lipgloss.Style {
       switch {
       case pct >= 80: return s.Critical
       case pct >= 60: return s.Warn
       default:        return s.OK
       }
   }
3. No new style creation. Just selection among pre-built instances.
```

This avoids per-call `lipgloss.NewStyle()` in loops. The selection is a simple switch returning an existing value.

---

## cliout Package Internal Structure

```go
package cliout

import (
    "fmt"
    "strings"

    "charm.land/lipgloss/v2"
    "github.com/N1xev/bubbleMonitor/internal/ui"
)

// CLIStyles holds pre-built styles for CLI output.
// Create once per command invocation via New(). Safe to pass by value.
type CLIStyles struct {
    // Structural styles
    Label    lipgloss.Style // Muted color, for field names ("CPU:", "Memory:")
    Value    lipgloss.Style // Text color, for field values ("45.2%", "8.2 GB")
    Header   lipgloss.Style // Bold + Primary, for section headers
    Bold     lipgloss.Style // Bold + Text, for emphasis
    Dim      lipgloss.Style // Faint + Muted, for secondary info (uptime, paths)

    // Status styles
    OK       lipgloss.Style // Success color, for healthy states
    Warn     lipgloss.Style // Warning color, for degraded states
    Critical lipgloss.Style // Alert color, for failures

    // Check mark styles (for doctor output)
    CheckOK   lipgloss.Style // Success-colored
    CheckFail lipgloss.Style // Alert-colored
    CheckWarn lipgloss.Style // Warning-colored

    // Specialized
    Active lipgloss.Style // Bold + Primary, for active theme marker

    // Palette reference (for threshold helpers that need color values)
    palette ui.ThemePalette
}

// New creates a CLIStyles from a ThemePalette.
// This is the only entry point. Call once per command invocation.
func New(palette ui.ThemePalette) CLIStyles {
    return CLIStyles{
        Label:     lipgloss.NewStyle().Foreground(palette.Muted),
        Value:     lipgloss.NewStyle().Foreground(palette.Text),
        Header:    lipgloss.NewStyle().Bold(true).Foreground(palette.Primary),
        Bold:      lipgloss.NewStyle().Bold(true).Foreground(palette.Text),
        Dim:       lipgloss.NewStyle().Faint(true).Foreground(palette.Muted),
        OK:        lipgloss.NewStyle().Foreground(palette.Success),
        Warn:      lipgloss.NewStyle().Foreground(palette.Warning),
        Critical:  lipgloss.NewStyle().Foreground(palette.Alert),
        CheckOK:   lipgloss.NewStyle().Foreground(palette.Success),
        CheckFail: lipgloss.NewStyle().Foreground(palette.Alert),
        CheckWarn: lipgloss.NewStyle().Foreground(palette.Warning),
        Active:    lipgloss.NewStyle().Bold(true).Foreground(palette.Primary),
        palette:   palette,
    }
}
```

### Helper Functions (same package, separate file)

```go
// cliout/helpers.go

// BarColor returns the palette color appropriate for a percentage.
// <60 = Success, <80 = Warning, >=80 = Alert.
func BarColor(pct float64, p ui.ThemePalette) compat.AdaptiveColor { ... }

// ScoreColor returns the palette color for a health score.
// >=70 = Success, >=50 = Warning, <50 = Alert.
func ScoreColor(score int, p ui.ThemePalette) compat.AdaptiveColor { ... }

// RenderBar builds a styled progress bar string.
func (s CLIStyles) RenderBar(pct float64, width int) string { ... }

// KVAligned renders a label-value pair with the label padded to a fixed width.
func (s CLIStyles) KVAligned(label string, labelWidth int, value string) string { ... }

// RenderSwatch builds a color swatch strip from a palette.
func RenderSwatch(p ui.ThemePalette) string { ... }

// PadRight pads a styled string to a target visual width.
func PadRight(styled string, targetWidth int) string { ... }
```

### Package File Layout

```
internal/cliout/
    cliout.go      -- CLIStyles struct, New() constructor
    helpers.go     -- BarColor, ScoreColor, RenderBar, KVAligned, RenderSwatch, PadRight
    cliout_test.go -- unit tests for style rendering and helper functions
```

---

## Patterns to Follow

### Pattern 1: Construct Once, Use Many Times

**What:** Build all `lipgloss.Style` values once in `cliout.New()`. Never call `lipgloss.NewStyle()` inside a render loop or output function.

**When:** Always. This is the primary pattern for the entire milestone.

**Example:**
```go
// In command RunE:
s := cliout.New(palette)

// In a loop:
for _, proc := range procs {
    style := s.Value // pre-built, no allocation
    if proc.Cpu > 80 {
        style = s.Critical // just picking a different pre-built style
    }
    lipgloss.Fprintf(out, "  %s\n", style.Render(fmt.Sprintf("%.1f%%", proc.Cpu)))
}
```

### Pattern 2: Write Through lipgloss.Fprintf, Not fmt.Fprintf

**What:** All styled output goes through `lipgloss.Fprintf` or `lipgloss.Fprint`, which handles color downsampling for non-TTY outputs.

**When:** Every `fmt.Fprintf(cmd.OutOrStdout(), ...)` that might contain a styled string.

**Example:**
```go
// Before:
fmt.Fprintf(out, "  CPU: %s\n", cpuStatus)

// After:
lipgloss.Fprintf(out, "  CPU: %s\n", s.Value.Render(cpuStatus))
```

### Pattern 3: Width-Aware Padding for Columns

**What:** When aligning styled text in columns, measure visual width with `lipgloss.Width()` and pad with literal spaces. Never use `fmt.Sprintf("%-Ns", styled)` or tabwriter for styled strings.

**When:** Any output with fixed-width columns (`top.go`, `sysinfo.go`, `remote.go`).

**Example:**
```go
name := s.Value.Render(proc.Name)
visualWidth := lipgloss.Width(name)
pad := 30 - visualWidth
if pad > 0 {
    name += strings.Repeat(" ", pad)
}
```

### Pattern 4: Config-First Resolution

**What:** Every styled command starts with the same three-line preamble to resolve the theme.

**When:** At the start of each command's RunE function that produces styled output.

**Example:**
```go
RunE: func(cmd *cobra.Command, args []string) error {
    cfg, _ := loadConfigWithOverrides()
    palette := ui.GetAppTheme(cfg.Theme, cfg.CustomTheme)
    s := cliout.New(palette)
    out := cmd.OutOrStdout()
    // ... styled output using s and lipgloss.Fprintf(out, ...)
}
```

---

## Anti-Patterns to Avoid

### Anti-Pattern 1: cliout Importing internal/config

**What:** Passing `AppConfig` or config-related types into the cliout package.

**Why bad:** Couples the styling library to configuration details. If config changes (new fields, new structure), cliout breaks. The cliout package should be able to style output given nothing but a `ThemePalette`.

**Instead:** cliout receives only `ui.ThemePalette`. Config resolution happens in `cmd/bub`.

### Anti-Pattern 2: CLIStyles as a Pointer or Singleton

**What:** Making `New()` return `*CLIStyles` or storing styles in a package-level variable.

**Why bad:** lipgloss v2 `Style` is a value type. `CLIStyles` containing only value types and a `ThemePalette` (also a value type) should itself be a value. No shared mutable state, no synchronization needed, no lifecycle management.

**Instead:** `func New(palette ui.ThemePalette) CLIStyles` -- returns a value. Commands hold it as a local variable.

### Anti-Pattern 3: Inline Style Creation in Loops

**What:** `lipgloss.NewStyle().Foreground(palette.Alert).Render(msg)` inside `for _, proc := range procs`.

**Why bad:** Creates a new Style value per iteration. While Style is a value type (not heap-allocated in the traditional sense), the interface conversion for `color.Color` still occurs each time. In `top.go`'s refresh loop processing 20-100 processes every 2 seconds, this adds up.

**Instead:** Use pre-built styles from `CLIStyles`. Pick among `s.OK`, `s.Warn`, `s.Critical` based on threshold.

### Anti-Pattern 4: fmt.Fprintf for Styled Strings

**What:** `fmt.Fprintf(out, "%s\n", someStyledString)`.

**Why bad:** `fmt.Fprintf` writes raw bytes. It does not check TTY status, `NO_COLOR`, or terminal color profile. Styled strings containing ANSI codes will leak into pipes and files.

**Instead:** `lipgloss.Fprintf(out, "%s\n", someStyledString)`.

### Anti-Pattern 5: Modifying internal/ui/styles.go to Add CLI Concerns

**What:** Adding a `ToCLIStyles()` method on `ThemePalette` or importing lipgloss render helpers into `internal/ui`.

**Why bad:** Violates the project constraint ("no TUI changes"). The TUI and CLI have different rendering lifecycles -- TUI uses BubbleTea's View pipeline, CLI writes directly to stdout. Mixing these concerns creates coupling.

**Instead:** The `internal/cliout` package imports `internal/ui` (for the type), not the other way around.

---

## Scalability Considerations

| Concern | At 1 command | At 10 commands (current) | At 20+ commands |
|---------|-------------|-------------------------|-----------------|
| Style construction overhead | Negligible (one New() call) | Negligible (each command creates its own) | Negligible (value types, ~14 styles per struct) |
| Memory per styled output | ~200 bytes for CLIStyles struct | ~200 bytes per command invocation | Same -- GC collects after command exits |
| Tabulated output rendering | O(n) where n = rows | Same | Consider lipgloss table package for complex tables |
| Theme resolution | O(1) map lookup | Same | Same |

The architecture is intentionally stateless. No caches, no singletons, no shared mutable state. Each command invocation creates its own `CLIStyles`, uses it, and lets the GC clean up when the process exits. This is appropriate for CLI tools where the process lifetime is seconds, not hours.

The one exception is `top.go`'s live-refresh loop, which runs for minutes. But even there, the `CLIStyles` struct is created once before the loop and reused across all refresh cycles within the same invocation.

---

## Build Order (Dependency Sequence)

This section specifies what must be built before what, based on the dependency graph.

```
Phase A: Foundation (no dependencies beyond existing code)
  1. Create internal/cliout/ package directory
  2. Create cliout.go with CLIStyles struct and New() constructor
  3. Create helpers.go with BarColor, ScoreColor, PadRight, KVAligned, RenderBar
  4. Create cliout_test.go with unit tests
  -> Tests pass: CLIStyles can be constructed from any ThemePalette

Phase B: Output Infrastructure (depends on Phase A)
  5. In root.go: add lipgloss.EnableLegacyWindowsANSI(os.Stdout) to initialization
  -> Foundation for all subsequent command styling

Phase C: Simple Commands (depends on Phase A, uses Phase B pattern)
  6. Style version.go (single-line output, simplest case)
  7. Style config.go (key-value output, KVAligned usage)
  8. Style export.go (conditional styling for human output, skip for machine formats)
  -> Validates the New(palette) -> use styles pattern across 3 commands

Phase D: Label-Value Commands (depends on Phase A)
  9. Style status.go (bars, KVAligned, threshold colors)
  10. Style health.go (status-dependent coloring, ScoreColor)
  11. Style sysinfo.go (labels, values, section headers, disk partitions)
  -> Validates RenderBar, threshold helpers, width-aware padding

Phase E: Diagnostic Command (depends on Phase A)
  12. Style doctor.go (CheckOK/CheckFail/CheckWarn, styled paths)
  -> Validates check mark styles, mixed styled/plain output

Phase F: Tabulated Commands (depends on Phase A + PadRight)
  13. Style top.go (header row, threshold-colored CPU%/MEM%, width-aware columns)
  14. Style remote.go (replace tabwriter or use width-compensated styling)
  -> Most complex; validates column alignment with ANSI codes

Phase G: Theme Display (depends on Phase A + RenderSwatch)
  15. Style themes.go (color swatches, active theme marker, dim inactive)
  -> Validates RenderSwatch, background-colored blocks
```

**Ordering rationale:**

- Phase A must come first because all other phases depend on the cliout package existing.
- Phase B is a one-line change but should happen before any styled output is written, for Windows compatibility.
- Phases C through G can technically proceed in any order, but the sequence above progresses from simplest to most complex. Each phase validates more of the cliout API. If a helper is wrong (e.g., PadRight miscalculates), you discover it in Phase C (cheap fix) before hitting Phase F (where alignment bugs are hardest to debug).
- `top.go` and `remote.go` are last because they have the trickiest alignment concerns (fixed-width format verbs and tabwriter respectively). By the time you reach them, the helper functions are battle-tested.

---

## Integration Points

### How cmd/bub/themes.go Already Uses the Pattern

The existing `themes.go` file (line 45) already does `palette := ui.GetTheme(name)` and (line 49) calls `renderColorSwatch(palette)`. The current `renderColorSwatch` is a stub returning `"Primary/Secondary/Success/Warning/Alert"`. The cliout package's `RenderSwatch` replaces this stub with actual colored blocks.

### How cmd/bub/health.go Already Loads Config

The existing `health.go` (line 28) calls `cfg, _ := loadConfigWithOverrides()` and uses `cfg.Thresholds` and `cfg.HealthWeights`. After cliout integration, it adds two lines:

```go
palette := ui.GetAppTheme(cfg.Theme, cfg.CustomTheme)
s := cliout.New(palette)
```

No structural change to the command. Just additive.

### How lipgloss.Fprintf Replaces fmt.Fprintf

The change is mechanical across all command files. Each `fmt.Fprintf(out, ...)` that touches styled text becomes `lipgloss.Fprintf(out, ...)`. The import changes from `"fmt"` to `"charm.land/lipgloss/v2"` (plus keeping `"fmt"` for any non-styled formatting like `fmt.Sprintf`).

---

## Confidence Assessment

| Aspect | Confidence | Reason |
|--------|------------|--------|
| Package location (`internal/cliout`) | HIGH | Standard Go internal package layout; verified against project structure |
| CLIStyles as value struct | HIGH | lipgloss v2 Style is value type; verified from source code |
| No circular imports | HIGH | Traced import graph: cliout -> ui (ThemePalette type only); ui never imports cliout |
| Data flow (config -> palette -> styles) | HIGH | Verified from existing code patterns in themes.go and health.go |
| Width-aware padding pattern | HIGH | Documented lipgloss v2 API; `lipgloss.Width()` exists in source |
| lipgloss.Fprintf for pipe safety | HIGH | Documented in official lipgloss README |
| Build order | MEDIUM | Based on dependency analysis, but actual implementation may surface unforeseen coupling |

---

## Sources

- Project source: `internal/ui/styles.go` -- ThemePalette definition, GetAppTheme(), 32 theme palettes
- Project source: `cmd/bub/root.go` -- loadConfigWithOverrides(), cobra command structure
- Project source: `cmd/bub/themes.go` -- existing ui.GetTheme() usage (proves import path works)
- Project source: `cmd/bub/health.go` -- existing loadConfigWithOverrides() in subcommand
- Project source: `cmd/bub/remote.go` -- tabwriter usage (alignment concern)
- Project source: `cmd/bub/top.go` -- fixed-width format verbs (alignment concern)
- lipgloss v2 source: `/home/samouly/go/pkg/mod/charm.land/lipgloss/v2@v2.0.2/style.go` -- Style is value type, Render is pure string transform
- lipgloss v2 compat: `/home/samouly/go/pkg/mod/charm.land/lipgloss/v2@v2.0.2/compat/color.go` -- AdaptiveColor implements color.Color via RGBA()
- lipgloss v2 set.go: Foreground(c color.Color) accepts any color.Color including compat.AdaptiveColor
- `.planning/research/STACK.md` -- previously verified lipgloss v2 patterns (HIGH confidence)
- `.planning/research/PITFALLS.md` -- previously identified alignment and writer pitfalls (HIGH confidence)

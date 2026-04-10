# Domain Pitfalls: Go CLI Styling with Lipgloss

**Domain:** Go CLI output styling using lipgloss v2 + compat layer, targeting Cobra subcommands
**Researched:** 2026-04-10
**Overall confidence:** HIGH (lipgloss v2 official docs verified, project code reviewed)

---

## Critical Pitfalls

Mistakes that cause rewrites, broken output, or user-visible regressions.

### Pitfall 1: Lipgloss ANSI Escape Codes Break Tabwriter Alignment

**What goes wrong:** `text/tabwriter` counts ANSI escape sequences as visible characters when computing column widths. A string like `\x1b[38;2;61;174;235mprod\x1b[0m` appears as 4 visible characters to a human but tabwriter counts ~25+, causing columns to overflow and misalign.

**Why it happens:** Go's `text/tabwriter` has zero awareness of ANSI escape sequences. It uses `len()` semantics internally for cell width. Lipgloss `Render()` wraps strings in ANSI escape codes for colors, bold, faint, etc. -- all invisible to tabwriter's width calculations.

**Consequences:** Column-aligned output (especially `remote list`, but potentially `top` and `sysinfo` if they adopt tabwriter) becomes jagged and unreadable. The misalignment gets worse with longer styled strings.

**Prevention:**
1. **For tabwriter-based output** (`remote.go` line 49): Style only the *values* that do NOT participate in tab alignment, or style the entire row as a single rendered string before writing to the writer. The safe pattern is: style each cell *before* it goes into tabwriter, then compensate width. Or better: use `lipgloss.Width()` to measure visual width and manually pad with spaces instead of relying on tabwriter for styled cells.
2. **For fixed-width column output** (`top.go` line 84): The `%-8d`, `%-30s` format verbs will break when the string contains ANSI codes. `fmt.Sprintf("%-30s", styledName)` pads based on byte length, not visual width. Use `lipgloss.Width()` to measure the visual width of the styled string, then manually pad to the desired column width using `lipgloss.PlaceHorizontal()` or by adding explicit space padding.
3. **Alternative**: Replace tabwriter with `lipgloss.JoinHorizontal()` for column layout, which is ANSI-aware. This is more work but eliminates the entire class of bugs.

**Detection:** Run `bub remote list` after styling and check that columns still align. Pipe through `cat -v` to see raw ANSI codes and verify byte counts vs visual widths.

**Phase:** Must be addressed in the very first styling phase (the `cliout` package creation), because the width-measurement helper functions need to exist before any styled output is written for tabulated commands.

**Project-specific impact:** `remote.go` uses `tabwriter.NewWriter` directly (line 49). `top.go` uses `fmt.Fprintf` with fixed-width format verbs (line 84-91). Both patterns will break with naive lipgloss styling.

---

### Pitfall 2: Style Objects Created Per-Call Instead of Per-Session

**What goes wrong:** Calling `lipgloss.NewStyle().Foreground(color).Bold(true).Render(text)` inside a loop or on every function invocation causes unnecessary allocations. Each call chain allocates a new Style value (though Style is a value type in v2, so the GC pressure is less than v1's pointer-based approach, it still triggers repeated interface conversion and color profile detection).

**Why it happens:** It feels natural to create styles inline where they are used. The lipgloss API encourages chained method calls.

**Consequences:** In `top.go`'s live-refresh loop (lines 67-111), which runs every 2 seconds and processes 20+ processes, per-call style creation adds measurable GC pressure. For a system monitor that should be lightweight, this is counterproductive.

**Prevention:**
1. Create a `CLIStyles` struct (already planned in PROJECT.md) that holds all pre-built styles as fields.
2. Instantiate `CLIStyles` once in `loadCLIStyles()` at command startup, not per render cycle.
3. Styles are value types in lipgloss v2 -- assigning `styleB := styleA` creates a true copy with zero shared mutation risk. So pre-built styles are safe to use from multiple goroutines.
4. For dynamic styles (e.g., CPU percentage color that changes based on threshold), use the pre-built `OK`, `Warn`, `Critical` styles and pick which one to use at render time. Do NOT create new styles for threshold decisions.

**Detection:** If `go test -bench` or `pprof` shows `lipgloss.NewStyle` or `runtime.convI2I` hot spots in the top command's loop, this pitfall has been hit.

**Phase:** Address in the `CLIStyles` struct creation phase. This is a design decision that must be made before writing any styled output code.

---

### Pitfall 3: Using `fmt.Fprintf(cmd.OutOrStdout(), ...)` Instead of `lipgloss.Fprint`

**What goes wrong:** The project currently uses `fmt.Fprintf(cmd.OutOrStdout(), ...)` everywhere. When lipgloss-styled strings are written via `fmt.Fprint`, the color downsample detection (NO_COLOR, non-TTY, 256-color terminal) does NOT run. The raw ANSI codes are written as-is, meaning:
- Piped output (`bub status > out.txt`) contains raw ANSI escape garbage
- `NO_COLOR=1 bub status` still shows colors
- Terminals that only support ANSI 16 colors receive truecolor escape sequences that may render incorrectly

**Why it happens:** `lipgloss.Render()` produces ANSI strings. These are just strings -- they contain the escape codes regardless of where they are written. Only the lipgloss writer functions (`lipgloss.Fprint`, `lipgloss.Println`, etc.) perform color profile detection and downsample colors for the target output.

**Consequences:** Users who pipe output to files get escape code garbage. Users with `NO_COLOR` set still see colors. Users on limited terminals get broken colors. This violates terminal conventions that CLI tools are expected to follow.

**Prevention:**
1. Replace `fmt.Fprintf(cmd.OutOrStdout(), ...)` with `lipgloss.Fprintf(cmd.OutOrStdout(), ...)` for all styled output.
2. The lipgloss writer functions are drop-in replacements for the `fmt` package: `Fprint`, `Fprintf`, `Fprintln`, `Print`, `Printf`, `Println`, `Sprint`, `Sprintf`, `Sprintln`.
3. For mixed output (some styled, some plain), the lipgloss writer functions handle plain strings fine -- they just pass through.
4. Important: the lipgloss writers use `colorprofile.NewWriter(os.Stdout, os.Environ())` internally, which checks `NO_COLOR`, `COLORTERM`, `TERM`, and whether the output is a TTY.

**Detection:** Run `NO_COLOR=1 bub status` -- if colors still appear, this pitfall is present. Run `bub status | cat -v` -- if ANSI escape codes appear in the pipe, this pitfall is present.

**Phase:** Must be established as the default output pattern from the very first phase. Every command file needs to import and use lipgloss writer functions instead of `fmt`.

**Project-specific note:** The current code uses `fmt.Fprintf` in all 10 command files. Every `fmt.Fprintf(out, ...)` call needs to become `lipgloss.Fprintf(out, ...)` for any output that might contain styled text. Unstyled output (error messages, etc.) can stay with `fmt` for now, but migrating all output to `lipgloss.Fprintf` is safer and more consistent.

---

### Pitfall 4: `compat.AdaptiveColor` vs `lipgloss.LightDark` Confusion

**What goes wrong:** The project's `ThemePalette` stores colors as `compat.AdaptiveColor` (from `charm.land/lipgloss/v2/compat`). This type resolves light/dark by querying stdin/stdout globally. But the CLI commands are standalone (not Bubble Tea), so the light/dark detection happens at the time the color is first used. If multiple commands run concurrently, or if stdin is piped, the detection can fail silently and default to dark mode.

**Why it happens:** `compat.AdaptiveColor` is the v1 compatibility layer. The lipgloss v2 docs explicitly warn: "we don't recommend this for new code as it removes the purity from Lipgloss, computationally speaking, as it removes transparency around when I/O happens, which could cause Lipgloss to compete for resources (like stdin) with other tools."

**Consequences:** Colors may not match the user's terminal background in edge cases (piped input, SSH sessions with no stdin, concurrent tool usage). The `tty` theme (which uses ANSI colors, not AdaptiveColor) is unaffected.

**Prevention:**
1. Since the `ThemePalette` already uses `compat.AdaptiveColor` and we are told NOT to modify `internal/ui/styles.go`, we must work with what we have.
2. When converting `compat.AdaptiveColor` to a `lipgloss.Style`, pass the color directly as `style.Foreground(palette.Primary)` -- the compat layer handles the resolution internally.
3. For the CLI output specifically, detect light/dark background once at startup using `lipgloss.HasDarkBackground(os.Stdin, os.Stdout)` and store the result. This avoids repeated I/O queries.
4. Document that `compat.AdaptiveColor` is an inherited constraint from the TUI and should not be used in new code outside `internal/ui`.

**Detection:** Test with `echo "" | bub status` (stdin is a pipe, not a TTY). If colors look wrong or the command blocks, the adaptive color resolution is failing on non-TTY stdin.

**Phase:** Address during `CLIStyles` initialization -- resolve all adaptive colors at construction time, not per render.

---

## Moderate Pitfalls

### Pitfall 5: Unicode Bar Characters and East Asian Width

**What goes wrong:** The `renderBar()` function in `status.go` uses `█` (U+2588 FULL BLOCK) and `░` (U+2591 LIGHT SHADE). These are typically single-width, but some terminals (especially East Asian locale terminals) may render them differently. If the bar characters are styled with lipgloss, and `lipgloss.Width()` is used for layout, the width measurement should be correct -- but if `len()` or `fmt.Sprintf` width formatting is used, multi-byte UTF-8 characters will break alignment.

**Prevention:** Use `lipgloss.Width()` for any width calculation involving styled or non-ASCII text. The current `renderBar()` builds a string and returns it -- after styling, always measure with `lipgloss.Width()` before padding.

**Phase:** Relevant when styling `status.go` bar output.

---

### Pitfall 6: `top.go` Screen Clear with ANSI Codes in Non-TTY

**What goes wrong:** `top.go` line 81 writes `\033[H\033[2J` (ANSI clear screen) when running in live mode. If the output is not a TTY (e.g., `bub top > file.txt`), these escape codes pollute the output file. This is a pre-existing issue that will become more noticeable once the rest of the output is also ANSI-styled.

**Prevention:** Before writing clear-screen sequences, check if the output is a TTY using `golang.org/x/term.IsTerminal(out.Fd())` or the `lipgloss` internal TTY detection. If not a TTY, skip the clear-screen and just print the snapshot once (or print sequentially without clearing).

**Phase:** Fix when styling `top.go`. This is a pre-existing issue but should be addressed alongside the styling work to avoid compounding the ANSI-in-pipe problem.

---

### Pitfall 7: Styling the Cobra `--help` Output

**What goes wrong:** Cobra generates help text internally. If you attempt to style Cobra's help output by overriding the help template with lipgloss-styled strings, the column alignment in Cobra's flags table will break (same tabwriter issue as Pitfall 1).

**Prevention:** Do NOT style the `--help` output in this milestone. It is out of scope per PROJECT.md and introduces the same tabwriter problems without clear benefit. Leave Cobra's help output as plain text.

**Phase:** Document as an explicit anti-feature for this milestone.

---

### Pitfall 8: Forgetting to Handle `cmd.OutOrStderr()` Cases

**What goes wrong:** Some Cobra commands write to stderr for certain outputs (e.g., error messages, debug info). If the `CLIStyles` package only handles stdout, styled error messages will lose formatting when redirected. Additionally, the `lipgloss.Writer` global variable defaults to stdout -- if you write styled strings to stderr via `fmt.Fprintf(cmd.OutOrStderr(), ...)`, the color profile detection may use the wrong output descriptor.

**Prevention:** When using `lipgloss.Fprintf`, always pass the correct writer. For stderr, use `lipgloss.Fprintf(os.Stderr, ...)`. If using Cobra's `cmd.OutOrStdout()` or `cmd.OutOrStderr()`, pass that writer directly to `lipgloss.Fprintf`.

**Phase:** Review each command file during its styling phase to ensure the correct writer is used.

---

### Pitfall 9: Style Inheritance and Mutation Surprises

**What goes wrong:** Lipgloss v2 styles are value types -- assigning `b := a` creates a copy. But method chaining like `a := baseStyle.Foreground(red)` also creates a new copy; it does NOT mutate `baseStyle`. However, `Inherit()` has subtle behavior: it only copies unset rules from the source. If you expect `Inherit` to overlay all rules, you will get unexpected results where already-set rules on the target are preserved even if you wanted them overridden.

**Prevention:**
1. Use `Inherit()` only for establishing base styles with overrides. Do not use it for dynamic per-render styling.
2. For the `CLIStyles` struct, build each style explicitly from `lipgloss.NewStyle()`. Do not use a chain of inherited styles that might have surprising interaction.
3. Test each style in isolation by rendering it and checking the output.

**Phase:** Relevant during `CLIStyles` struct design.

---

## Minor Pitfalls

### Pitfall 10: Missing `lipgloss.EnableLegacyWindowsANSI()` Call

**What goes wrong:** On older Windows consoles (pre-Windows 10 1511), ANSI escape sequences are not enabled by default. Lipgloss v2 provides `lipgloss.EnableLegacyWindowsANSI(*os.File)` for this purpose, but it must be called explicitly.

**Prevention:** Add `lipgloss.EnableLegacyWindowsANSI(os.Stdout)` early in the CLI initialization (e.g., in `root.go` or `main()`). This is a no-op on non-Windows platforms, so it is safe to call unconditionally. Note that the function signature takes `*os.File`, so pass the actual file, not a wrapped writer.

**Phase:** Address in the root command setup phase.

---

### Pitfall 11: Progress Bar Colors Not Matching Theme Intent

**What goes wrong:** The `BarColor(pct)` helper (planned in PROJECT.md) maps `<60 green, <80 yellow, >=80 red`. But if the theme's Success color is green and Warning is yellow, the bar colors should use `theme.Success` and `theme.Warning`, not hardcoded green/red. Some themes (like `autumn`, `coffee`, `sunset`) have non-standard color mappings where the theme's "Alert" color might not be red.

**Prevention:** Use the theme palette colors directly: `palette.Success` for healthy ranges, `palette.Warning` for caution, `palette.Alert` for critical. Do not hardcode green/yellow/red in the `BarColor` helper.

**Phase:** Address when implementing `BarColor()` and `ScoreColor()` helpers.

---

### Pitfall 12: Lipgloss v2 TabWidth Interference

**What goes wrong:** Lipgloss v2 converts tab characters (`\t`) to 4 spaces at render time by default. If the existing code uses `\t` for alignment (e.g., in `remote.go`'s tabwriter), and the text passes through `lipgloss.Style.Render()` before reaching tabwriter, tabs will already be converted to spaces, potentially doubling the spacing or breaking tabwriter's expected input format.

**Prevention:** For text that goes through tabwriter, either:
1. Do NOT pass it through `lipgloss.Style.Render()` before tabwriter -- style individual cells and handle width manually, OR
2. Use `style.TabWidth(lipgloss.NoTabConversion)` to preserve tabs, then render, then pass to tabwriter.

**Phase:** Critical for `remote.go` styling. Verify with a test case.

---

## Phase-Specific Warnings

| Phase Topic | Likely Pitfall | Mitigation |
|-------------|---------------|------------|
| `cliout` package creation | Pitfall 2 (per-call styles) | Enforce struct-based style cache from day one |
| `cliout` package creation | Pitfall 3 (fmt vs lipgloss writers) | Mandate `lipgloss.Fprintf` in all output |
| `status.go` styling | Pitfall 5 (Unicode bar width) | Use `lipgloss.Width()` for bar measurement |
| `top.go` styling | Pitfall 1 (fixed-width format + ANSI) | Replace `%-30s` with manual width padding |
| `top.go` styling | Pitfall 6 (clear-screen in non-TTY) | Add TTY check before ANSI clear |
| `remote.go` styling | Pitfall 1 (tabwriter + ANSI) | Style cells before tabwriter, compensate width |
| `remote.go` styling | Pitfall 12 (tab conversion) | Preserve tabs or avoid Render before tabwriter |
| `doctor.go` styling | Pitfall 9 (inherit surprises) | Build checkmark/cross styles explicitly |
| `sysinfo.go` styling | Pitfall 1 (fixed-width format + ANSI) | Same as `top.go` -- manual width padding |
| `themes.go` styling | Pitfall 8 (stderr vs stdout) | Verify output writer for theme swatch rendering |
| Root command setup | Pitfall 10 (Windows ANSI) | Add `EnableLegacyWindowsANSI` call |
| All threshold-based styling | Pitfall 11 (hardcoded colors) | Use theme palette, not hardcoded values |
| Adaptive color resolution | Pitfall 4 (compat layer I/O) | Resolve once at CLIStyles construction |

---

## Project-Specific Risk Summary

**Highest risk commands** (most likely to break with naive styling):
1. **`remote list`** -- Uses tabwriter. Any styled cell breaks alignment.
2. **`top`** -- Uses fixed-width format verbs + live refresh loop. Performance and alignment both at risk.
3. **`status`** -- Uses Unicode bar characters. Width measurement must be ANSI-aware.

**Lowest risk commands** (simple label-value output, unlikely to break):
1. **`version`** -- Single-line output, no alignment concerns.
2. **`config`** -- Simple key-value output.
3. **`export`** -- Machine-readable output, may need conditional styling.

---

## Sources

- [Lipgloss v2 official docs](https://pkg.go.dev/charm.land/lipgloss/v2) (pkg.go.dev, published Mar 11, 2026) -- HIGH confidence
- [Lipgloss GitHub README](https://github.com/charmbracelet/lipgloss) -- HIGH confidence
- Project source: `internal/ui/styles.go` (ThemePalette definition, compat.AdaptiveColor usage) -- HIGH confidence
- Project source: `cmd/bub/remote.go` (tabwriter usage, line 49) -- HIGH confidence
- Project source: `cmd/bub/top.go` (fixed-width format verbs, lines 84-91) -- HIGH confidence
- Project source: `cmd/bub/status.go` (Unicode bar rendering, line 98-112) -- HIGH confidence
- [NO_COLOR spec](https://no-color.org/) -- HIGH confidence
- WebSearch findings on tabwriter + ANSI alignment -- MEDIUM confidence (widely reported, consistent across sources)

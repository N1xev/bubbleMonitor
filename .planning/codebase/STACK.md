# Technology Stack

**Analysis Date:** 2026-04-09

## Languages

**Primary:**
- Go 1.25.5 - Entire codebase (module: `github.com/N1xev/bubbleMonitor`)
- CGO enabled conditionally - Required for NVML (NVIDIA) and AMD SMI GPU integrations on Linux only

## Runtime

**Environment:**
- Go 1.25.5 (specified in `go.mod`)
- Statically linked by default (`CGO_ENABLED=0`) except Linux amd64 which uses `CGO_ENABLED=1` for GPU SDKs

**Package Manager:**
- Go modules
- Lockfile: `go.sum` present

## Frameworks

**Core TUI:**
- `charm.land/bubbletea/v2` v2.0.0 - Elm-architecture TUI framework (The Tea architecture)
- `charm.land/bubbles/v2` v2.0.0 - Pre-built TUI components (table, list, etc.)
- `charm.land/lipgloss/v2` v2.0.0 - Terminal styling and layout

**CLI:**
- `github.com/spf13/cobra` v1.10.2 (indirect) - CLI command framework (subcommands, flags, completions)
- `github.com/charmbracelet/fang` v1.0.0 (indirect) - Styled Cobra execution wrapper (custom error rendering)

**Mouse/Zone Tracking:**
- `github.com/lrstanley/bubblezone/v2` v2.0.0 - Mouse zone tracking for Bubble Tea

**Build/Dev:**
- GoReleaser v2 - Cross-platform release builds (`.goreleaser.yml`)
- golangci-lint - Linting (`.golangci.yml` with errcheck, gosimple, govet, ineffassign, staticcheck, unused, gofmt, goimports, revive)

## Key Dependencies

**System Monitoring:**
- `github.com/shirou/gopsutil/v3` v3.24.5 - Cross-platform system metrics (CPU, memory, disk, network, process, host, load, sensors)
- `github.com/NVIDIA/go-nvml` v0.13.0-1 - NVIDIA GPU management library bindings (Linux + CGO only)
- `github.com/hhk7734/amdsmi.go` v0.2.0 - AMD GPU SMI library bindings (Linux + CGO only)
- `github.com/distatus/battery` v0.11.0 - Cross-platform battery status

**TUI Utilities:**
- `github.com/charmbracelet/harmonica` v0.2.0 - Animation/easing library
- `github.com/charmbracelet/x/ansi` v0.11.6 - ANSI escape sequence handling
- `github.com/charmbracelet/x/term` v0.2.2 - Terminal utilities
- `github.com/charmbracelet/x/termios` v0.1.1 - Terminal I/O control
- `github.com/charmbracelet/colorprofile` v0.4.2 - Terminal color profile detection
- `github.com/lucasb-eyer/go-colorful` v1.3.0 - Color manipulation
- `github.com/mattn/go-runewidth` v0.0.20 - East Asian character width
- `github.com/xo/terminfo` v0.0.0-20220910002029-abceb7e1c41e - Terminal capability detection
- `github.com/clipperhouse/displaywidth` v0.11.0 - Display width calculation
- `github.com/clipperhouse/uax29/v2` v2.7.0 - Unicode text segmentation

**Platform-Specific (indirect):**
- `github.com/go-ole/go-ole` v1.2.6 - Windows COM/OLE (for gopsutil on Windows)
- `github.com/power-devops/perfstat` v0.0.0-20210106213030-5aafc221ea8c - AIX performance stats
- `github.com/tklauser/go-sysconf` v0.3.12 - POSIX sysconf
- `github.com/tklauser/numcpus` v0.6.1 - CPU count detection
- `github.com/yusufpapurcu/wmi` v1.2.4 - Windows WMI
- `github.com/shoenig/go-m1cpu` v0.1.6 - Apple M1 CPU detection
- `howett.net/plist` v1.0.0 - macOS plist parsing
- `github.com/lufia/plan9stats` v0.0.0-20211012122336-39d0f177ccd0 - Plan 9 stats
- `golang.org/x/sys` v0.41.0 - Low-level system interaction
- `golang.org/x/text` v0.35.0 - Unicode text processing

**Concurrency:**
- `golang.org/x/sync` v0.20.0 - Extended sync primitives (errgroup, singleflight, etc.)

## Build Configuration

**GoReleaser** (`.goreleaser.yml`):
- Binary name: `bub`
- Entry point: `./cmd/bub`
- Platforms: linux (amd64, arm64), darwin (amd64, arm64), windows (amd64, arm64)
- Linux amd64 built with `CGO_ENABLED=1` (for NVML/AMD SMI)
- All other platforms built with `CGO_ENABLED=0`
- ldflags: `-s -w` (strip debug), `-X main.version={{.Version}}`, `-X main.commit={{.Commit}}`, `-X main.date={{.Date}}`
- Pre-build hook: `go generate ./...`
- Archives: binary format and tar.gz/zip
- GitHub release target: `N1xev/bubbleMonitor`

**Linting** (`.golangci.yml`):
- Enabled: errcheck, gosimple, govet, ineffassign, staticcheck, unused, gofmt, goimports, revive
- govet: `enable-all: true`
- errcheck: `check-type-assertions: true`, `check-blank: false`
- staticcheck: all checks except SA1019 (deprecation warnings suppressed)
- revive rules: blank-imports, context-as-argument, dot-imports, error-return, error-strings, error-naming, receiver-naming, time-naming, errorf
- Issues: `exclude-use-default: false` (all issues reported)
- Timeout: 5 minutes

**Version injection:**
- Variables `buildVersion`, `buildCommit`, `buildDate` in `cmd/bub/root.go` set via ldflags
- Default values: `"dev"`, `"none"`, `"unknown"`

## Build Tags / Conditional Compilation

**Pattern:** `//go:build` constraints for GPU SDK integration

| Build Tag | Files | Purpose |
|-----------|-------|---------|
| `linux && cgo` | `internal/provider/system/detect_nvidia.go` | NVML-based NVIDIA detection |
| `linux && cgo` | `internal/provider/system/hardware_nvidia.go` | NVML NVIDIA GPU metrics |
| `linux && cgo` | `internal/provider/system/hardware_amd.go` | AMD SMI GPU metrics |
| `!linux \|\| !cgo` | `internal/provider/system/detect_nvidia_stub.go` | nvidia-smi CLI fallback (no NVML) |
| `!linux \|\| !cgo` | `internal/provider/system/hardware_nvidia_stub.go` | No-op NVML functions |
| `!linux \|\| !cgo` | `internal/provider/system/hardware_amd_stub.go` | No-op AMD SMI functions |

**Implication:** Full GPU metrics (NVML, AMD SMI) only available when building with CGO on Linux. Other platforms rely on nvidia-smi CLI or sysfs fallbacks.

**Platform branching at runtime** uses `runtime.GOOS` checks in:
- `internal/provider/system/hardware.go` - GPU detection per platform (darwin/windows/linux)
- `internal/provider/system/vm.go` - VM detection (linux vs windows)
- `internal/provider/system/logs.go` - Log reading (linux journalctl only)
- `internal/provider/system/services.go` - Service listing (linux systemctl only)
- `internal/provider/system/metrics.go` - Root disk path (`/` vs `C:`)
- `internal/provider/system/container.go` - Docker socket path (unix vs windows)
- `internal/provider/process/control.go` - Process renice (renice vs powershell)

## Configuration

**Format:** JSON

**Location:** `~/.config/bubble-monitor/config.json` (via `os.UserConfigDir`)

**Overridable via CLI flags:**
- `--config` - Config file path
- `--theme` - Color theme (32 built-in themes)
- `--refresh-rate` - Refresh interval in ms (500, 1000, 2000, 5000)
- `--history-length` - Data points for charts (default: 900)

**Key config options** (`internal/config/config.go`):
- `theme` - dark, light, nord, dracula, gruvbox, solarized, monokai, catppuccin, tokyonight, onedark, ayu, rosepine, everforest, nightowl, palenight, material, synthwave, cobalt2, horizon, oceanic, palefire, github, moonlight, shades, midnight, forest, autumn, cyberpunk, sunset, ocean, coffee, custom
- `chart_type` - Chart rendering style (default: braille)
- `border_type` - normal, rounded
- `border_style` - single, double, dashed
- `remote_hosts` - SSH remote monitoring hosts array
- `thresholds` - Alert thresholds for CPU (90), memory (90), disk (90), temperature (85)
- `health_weights` - Health score calculation weights
- `logging` - Optional file logging (path, enabled)
- `process_cpu_normalized` - Normalize CPU percentages (default: true)
- `background_opaque` - Opaque vs transparent background (default: true)

**Hot reload:** Polls config file every 2 seconds for modification time changes (`internal/config/watcher.go`)

## CLI Subcommands

Defined in `cmd/bub/root.go`, registered via cobra:
- `bub` (default) - Launch TUI monitor
- `bub version` - Show version
- `bub status` - Quick system status
- `bub sysinfo` - Detailed system info
- `bub top` - Top-like process view
- `bub ps` - Process list
- `bub health` - System health score
- `bub export` - Export metrics snapshot
- `bub doctor` - Run diagnostic checks
- `bub themes` - Theme management
- `bub config` - Config management
- `bub remote` (add/list/remove) - SSH remote host management

## Platform Requirements

**Development:**
- Go 1.25.5+
- CGO toolchain (gcc/clang) required for GPU SDK integration on Linux
- NVIDIA driver + NVML library for NVIDIA GPU monitoring
- AMD ROCm SMI for AMD GPU monitoring

**Production:**
- Static binary for most platforms (CGO_ENABLED=0)
- Dynamic binary for linux/amd64 (CGO_ENABLED=1 for GPU SDKs)
- Optional system tools for enhanced monitoring:
  - `nvidia-smi` - NVIDIA GPU info
  - `rocm-smi` - AMD GPU detection
  - `systemctl` - Linux service listing
  - `journalctl` - System log reading
  - `lsblk` - Block device enumeration
  - `docker` - Container monitoring
  - `kubectl` - Kubernetes pod monitoring
  - `ssh` - Remote host monitoring
  - `renice` - Process priority adjustment (Unix)
  - `system_profiler` - macOS GPU info
  - `wmic` - Windows GPU info
  - `dmidecode` - VM hypervisor detection
  - `powershell` - Windows process priority / VM detection
  - `upower` - Battery detection (Linux)

---

*Stack analysis: 2026-04-09*

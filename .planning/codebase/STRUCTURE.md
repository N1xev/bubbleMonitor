# Codebase Structure

**Analysis Date:** 2026-04-09

## Directory Layout

```
bubbleMonitor/
├── cmd/                        # CLI entry points
│   └── bub/                    # Main CLI binary ("bub")
├── dist/                       # Build output (GoReleaser)
├── internal/                   # Private application packages
│   ├── app/                    # BubbleTea model, update loop, handlers
│   │   ├── handlers/           # Input handling (keyboard, mouse, toast, config)
│   │   └── testutil/           # Test model factories and mock providers
│   ├── config/                 # Configuration loading, saving, watching
│   ├── data/                   # State types, business logic, ring buffer
│   ├── msg/                    # BubbleTea message type definitions
│   ├── provider/               # Provider interfaces and adapters
│   │   ├── process/            # Process listing and control (gopsutil)
│   │   ├── remote/             # Remote host monitoring (SSH)
│   │   └── system/             # System metrics (gopsutil, NVML, sysfs)
│   ├── ui/                     # Terminal UI rendering
│   │   ├── input/              # Mouse zone management
│   │   ├── overlays/           # Modal overlays (help, settings, kill, files)
│   │   ├── tabs/               # Per-tab content renderers
│   │   └── widgets/            # Reusable UI primitives (charts, progress, borders)
│   └── util/                   # Shared utilities (format, math, text, layout)
├── .github/                    # GitHub config
│   └── workflows/              # CI workflows
├── .planning/                  # Planning documents
│   └── codebase/               # Codebase analysis documents
├── go.mod                      # Go module definition
├── go.sum                      # Dependency checksums
├── .goreleaser.yml             # GoReleaser build config
└── .golangci.yml               # Linter config
```

## Directory Purposes

### `cmd/bub/`
- Purpose: CLI binary entry point and Cobra command definitions
- Contains: `main.go` (program launch), `root.go` (Cobra root + TUI launch), plus one file per subcommand
- Key files:
  - `main.go` -- `main()` function, creates and runs `tea.Program`
  - `root.go` -- Cobra root command, `launchTUI()`, `loadConfigWithOverrides()`, CLI flags
  - `version.go` -- `bub version` subcommand
  - `status.go` -- `bub status` subcommand (quick system overview)
  - `sysinfo.go` -- `bub sysinfo` subcommand (detailed system info)
  - `top.go` -- `bub top` subcommand (top-like process view)
  - `health.go` -- `bub health` subcommand (health score)
  - `export.go` -- `bub export` subcommand (snapshot export)
  - `doctor.go` -- `bub doctor` subcommand (diagnostic checks)
  - `themes.go` -- `bub themes` subcommand (theme management)
  - `config.go` -- `bub config` subcommand (config management)
  - `remote.go` -- `bub remote` subcommand (SSH remote host management)

### `internal/app/`
- Purpose: Core BubbleTea application model, update dispatcher, and handlers
- Contains: The `Model` struct, `Init()`, `Update()`, `View()` methods, and extracted handler logic
- Key files:
  - `model.go` -- `Model` struct (embeds `AppState`), `InitialModel()`, `InitialModelWithConfig()`, `Init()`, `View()`, `RenderCache`, `GetBorder()`
  - `update.go` -- `Update()` method: 422-line type-switch dispatcher for all message types
  - `analysis.go` -- `UpdateHealthScore()` (weighted penalty system), `UpdateProcessHistory()` (per-process ring buffers)
  - `export.go` -- `SaveSnapshotCmd()` (JSON + CSV snapshot export)
  - `logging.go` -- `LogMetricsCmd()` (optional file-based metrics logging)

### `internal/app/handlers/`
- Purpose: Extracted input handling logic from the update loop
- Contains: Functions that take `*data.AppState` and return `tea.Cmd`
- Key files:
  - `keyboard.go` -- `HandleKey()` (939 lines), `KillProcessCmd()`, `handleSettingsChange()`, `reorderTab()`
  - `mouse.go` -- `HandleMouse()`, `handleLeftClick()`, `handleRightClick()`, `handleScrollUp()`, `handleScrollDown()`
  - `toast.go` -- `HandleToast()`, `HandleToastTimeout()`, `AddToastCmd()`
  - `config.go` -- Settings change helper logic

### `internal/app/testutil/`
- Purpose: Test infrastructure -- model factories and mock providers
- Contains: Functions to create test models with pre-populated state
- Key files:
  - `model.go` -- `TestModel` struct, `NewTestModel()`, `NewModelWithProcesses()`, `NewModelWithMetrics()`, `RenderCache`
  - `mock_system.go` -- `MockSystemProvider` with overridable function fields, mock data generators (`MockCpuMem`, `MockProcesses`, `MockDisk`, `MockNetwork`, `MockTemp`, `MockBattery`, `MockHostInfo`, `MockGpuInfo`, `MockCpuInfo`)
  - `mock_process.go` -- `MockProcessProvider` with overridable function fields
  - `mock_remote.go` -- `MockRemoteProvider` with overridable function fields

### `internal/config/`
- Purpose: JSON configuration loading, saving, validation, and hot-reload
- Key files:
  - `config.go` -- `AppConfig` struct (19 fields), `DefaultConfig()`, `LoadConfig()`, `LoadConfigFromPath()`, `SaveConfig()`, `GetConfigPath()`, `ResolvePath()`
  - `watcher.go` -- `WatchConfig()` (poll-based mtime check), `ConfigChangeMsg`, `ConfigWatchTickMsg`, `scheduleNextCheck()`
  - `options.go` -- `GetThemeNames()` (32 themes), `GetBorderStyles()`, `GetBorderTypes()`, `GetRefreshRates()`, `DefaultCustomTheme()`

### `internal/data/`
- Purpose: Central state definition, data structures, and business logic
- Key files:
  - `state.go` -- `AppState` struct (root state), mutex-protected methods (`SetSuspended`, `IsSuspended`, `ToggleCollapsed`, `IsCollapsed`, `ToggleBookmark`, `IsBookmarked`, `GetProcessByPid`, `SyncProcessesMap`, `ClearProcessMaps`, `PruneDeadProcessMaps`, `GetOrCreateHistory`, `GetHistory`, `PruneDeadProcessHistory`)
  - `groups.go` -- `MetricsState` (history buffers, current metrics, hardware flags, GPU info, alerts), `ProcessState` (process list, tree view, filtering, selection, scroll, kill dialog, open files, services, connections, logs, context menu), `UIState` (window dimensions, tab state, interactive flags, toasts, scroll offsets, mouse position, zones, chart display), `ConfigState` (display config, sorting, full config reference), `RemoteState` (remote host metrics map), `RemoteHostMetrics`, `RemoteProcessInfo`
  - `types.go` -- `ProcessInfo` (12 fields + lowercase caches), `ProcessSnapshot`, `DiskPartition`, `GpuInfo` (13 fields), `Toast` (6 fields), `ServiceInfo`, `ConnectionInfo`, `HardwareCapabilities` (10 bool flags), `ContainerInfo`, `K8sPodInfo`, `VmInfo`
  - `ringbuffer.go` -- `RingBuffer` struct with `Push`, `Get`, `Len`, `Max`, `Avg`, `Resize`, `Accessor` interface
  - `logic.go` -- `GetVisibleProcesses()` (cached filter + tree build), `InvalidateProcessCache()`, `GetFilteredProcesses()` (name/username/cmdline/pid substring match), `buildProcessTree()` (recursive tree with collapse support)
  - `alerts.go` -- `AlertManager` struct, `NewAlertManager()`, `CheckAlerts()`, `GetAlerts()`, `HasAlerts()`
  - `constants.go` -- `ReservedContentRows=19`, settings indices (`ThresholdCount=4`, `DisplayCount=6`, `TabCount=9`, `AppearanceCount=5`, `TotalSettingsCount=24`), `AllAvailableTabs`, health score constants, overlay dimension constants
  - `viewport.go` -- `SimpleViewport` wrapper around `bubbles/v2/viewport.Model`

### `internal/msg/`
- Purpose: All BubbleTea message type definitions
- Key files:
  - `messages.go` -- ~25 message types (TickMsg, CpuMemMsg, DiskNetMsg, MetricsMsg, ProcessesMsg, ProcessCountMsg, HostInfoMsg, DiskInfoMsg, GpuInfoMsg, DiskIOMsg, TempMsg, NetworkInterfacesMsg, BatteryMsg, ServicesMsg, ConnectionsMsg, SysLogMsg, RemoteMsg, PriorityChangeMsg, ProcessControlMsg, OpenFilesMsg, OpenFilesRequestMsg, ProcessCmdlineMsg, ProcessUsernameMsg, ToastMsg, ToastTimeoutMsg, QuitMsg, ExportSnapshotMsg, ForceRefreshMsg, KillProcessMsg, ContainerInfoMsg, VmInfoMsg)

### `internal/provider/`
- Purpose: Provider interfaces and adapter structs (thin delegation layer)
- Key files:
  - `interfaces.go` -- `SystemProvider` (17 methods), `ProcessProvider` (9 methods), `RemoteProvider` (1 method)
  - `adapters.go` -- `SystemAdapter`, `ProcessAdapter`, `RemoteAdapter` structs wrapping package-level functions; compile-time interface checks at bottom
  - `constants.go` -- Shared provider constants

### `internal/provider/system/`
- Purpose: Concrete system metrics collection using gopsutil, NVML, sysfs
- Key files:
  - `metrics.go` -- `TickCmd()`, `FastMetricsCmd()` (CPU, memory, swap, load), `SlowMetricsCmd()` (disk usage, network aggregate)
  - `hardware.go` -- `DetectHardware()` (returns `HardwareCapabilities`), GPU fetch for macOS/Windows/Linux, PCI device name database
  - `hardware_nvidia.go` -- NVML-based NVIDIA GPU metrics (build tag: `linux && cgo`)
  - `hardware_nvidia_stub.go` -- No-op stubs when NVML unavailable (build tag: `!linux || !cgo`)
  - `hardware_amd.go` -- AMD SMI GPU metrics (build tag: `linux && cgo`)
  - `hardware_amd_stub.go` -- No-op stubs when AMD SMI unavailable
  - `detect.go` -- Hardware detection helpers (`detectAmd`, `detectBattery`, etc.)
  - `detect_nvidia.go` -- NVML-based NVIDIA detection (build tag: `linux && cgo`)
  - `detect_nvidia_stub.go` -- Fallback nvidia-smi detection (build tag: `!linux || !cgo`)
  - `disk.go` -- `DiskInfoCmd()`, `DiskIOCmd()`, `parseLsblkUnmounted()`
  - `network.go` -- `NetworkInterfacesCmd()`
  - `connections.go` -- `ConnectionsCmd()`
  - `services.go` -- `ServicesCmd()` (systemctl wrapper, Linux-only)
  - `logs.go` -- `SystemLogsCmd()` (journalctl wrapper, Linux-only)
  - `container.go` -- `ContainerInfoMsg` producer (Docker/Kubernetes detection)
  - `vm.go` -- `VmInfoMsg` producer (VM/hypervisor detection via /proc, DMI, dmidecode)
  - `constants.go` -- System-specific constants (SSH timeout defaults)

### `internal/provider/process/`
- Purpose: Process listing, caching, sorting, and control operations
- Key files:
  - `list.go` -- `ProcessesCmd()`, `ProcessCountCmd()`, `PidsOnlyCmd()`, process cache management, sort logic
  - `control.go` -- `ReniceProcessCmdSafe()`, `SuspendProcessCmd()`, `ResumeProcessCmd()`, `KillProcessCmd()`
  - `files.go` -- `FetchOpenFilesCmd()`, `FetchProcessCmdlineCmd()`, `FetchProcessUsernameCmd()`
  - `utils.go` -- String interner (LRU, 5000 cap), `sync.Pool` for slice reuse (`GetProcSlice`/`PutProcSlice`), `Intern()` function
  - `utils_test.go` -- Tests for process utilities
  - `constants.go` -- Process-specific constants

### `internal/provider/remote/`
- Purpose: SSH-based remote host monitoring
- Key files:
  - `ssh.go` -- `CheckRemoteCmd()`, `buildSSHCmd()`, remote output parsers (`parseRemoteOutput`, `parseUptime`, `parseLoadAvg`, `parseMeminfo`, `parseDisk`, `parseNet`, `parseProcesses`)
  - `constants.go` -- SSH timeout defaults

### `internal/ui/`
- Purpose: Terminal UI rendering from AppState (pure functions, no state mutation)
- Key files:
  - `layout.go` -- `MainViewFromState()` (1070 lines), `RenderFromAppState()`, `ViewModel` interface, `ThemePalette` resolution, header/footer/tab/overlay assembly, compositor layer management, zone registration
  - `styles.go` -- `ThemePalette` struct, 32 theme definitions (dark, light, nord, dracula, gruvbox, solarized, monokai, catppuccin, tokyonight, onedark, ayu, rosepine, everforest, nightowl, palenight, material, synthwave, cobalt2, horizon, oceanic, palefire, github, moonlight, shades, midnight, forest, autumn, cyberpunk, sunset, ocean, coffee, tty, custom), `GetTheme()`, `GetThemeFromCustom()`, `GetAppTheme()`
  - `constants.go` -- `MinWindowWidth=80`, `MinWindowHeight=24`, `WideLayoutThreshold=130`, re-exported overlay dimensions

### `internal/ui/tabs/`
- Purpose: Per-tab content rendering functions
- Key files:
  - `metrics.go` -- `RenderMetrics()` -- Charts for CPU, memory, swap, disk IO, network, temperature; per-core CPU grid
  - `processes.go` -- `RenderProcesses()` -- Scrollable process list with tree/flat view, bookmarks, inline selection
  - `system.go` -- `RenderSystem()` -- Multi-block responsive layout (host info, CPU, memory, load, GPU, battery, temperatures) with per-block scrolling
  - `disks.go` -- `RenderDisks()` -- Disk partition table with usage bars
  - `network.go` -- `RenderNetwork()` -- Network interface statistics table
  - `services.go` -- `RenderServices()` -- Systemd service list (scrollable)
  - `connections.go` -- `RenderConnections()` -- Network connections table (scrollable)
  - `logs.go` -- `RenderLogs()` -- System log viewer (scrollable)
  - `remote.go` -- `RenderRemote()` -- Remote host metric cards
  - `overview.go` -- `RenderOverview()` -- Overview tab (if present)
  - `helpers.go` -- Shared tab rendering utilities
  - `constants.go` -- Tab-specific constants

### `internal/ui/overlays/`
- Purpose: Modal overlay rendering (centered dialogs on top of content)
- Key files:
  - `help.go` -- `RenderHelp()` -- Keybinding reference overlay
  - `settings.go` -- `RenderSettingsOverlay()` -- 3-column settings panel (thresholds, display, tabs, appearance)
  - `kill.go` -- `RenderKillDialog()` -- Kill confirmation dialog
  - `files.go` -- `RenderOpenFilesOverlay()` -- Open files viewer for a process
  - `samlab.go` -- `RenderSamLab()` -- About/credits overlay

### `internal/ui/widgets/`
- Purpose: Reusable UI primitives (charts, progress bars, borders)
- Key files:
  - `charts.go` -- Chart rendering (braille, line, bar types) using `RingBuffer.Accessor` interface
  - `charts_test.go` -- Chart rendering tests
  - `charts_bench_test.go` -- Chart rendering benchmarks
  - `progress.go` -- Progress bar widget
  - `borders.go` -- `GetBorder()` function (single, double, dashed, rounded)
  - `constants.go` -- Widget-specific constants

### `internal/ui/input/`
- Purpose: Mouse zone management for interactive elements
- Key files:
  - `types.go` -- `Zone` struct, `ZoneType` enum, `ZoneManager` interface, `zoneManager` implementation
  - `zones.go` -- `NewZoneManager()`, `Register()`, `FindZoneAt()`, `GetHoveredZone()`, `IsHovered()`, `UpdateMousePos()`, `Clear()`, `GetZones()`
  - `zones_test.go` -- Zone manager tests

### `internal/util/`
- Purpose: Shared leaf-package utilities (no internal dependencies)
- Key files:
  - `format.go` -- Number/string formatting helpers
  - `layout.go` -- Layout calculation helpers
  - `math.go` -- Math utility functions
  - `text.go` -- Text manipulation helpers
  - `fast.go` -- Performance-optimized utility functions

### `dist/`
- Purpose: Build output directory (GoReleaser)
- Contains: Compiled binaries for cross-platform distribution
- Generated: Yes (by GoReleaser)
- Committed: No (should be gitignored)

## Package Dependency Graph

```
cmd/bub
  └── internal/app
        ├── internal/config
        ├── internal/data
        │     ├── internal/config
        │     ├── internal/ui/input  (for Zone types)
        │     └── gopsutil/v3 (metric types)
        ├── internal/msg
        │     ├── internal/data
        │     └── gopsutil/v3 (raw types)
        ├── internal/provider
        │     ├── internal/provider/system
        │     │     ├── internal/msg
        │     │     ├── internal/data
        │     │     └── gopsutil/v3, NVML, AMD-SMI, battery
        │     ├── internal/provider/process
        │     │     ├── internal/msg
        │     │     ├── internal/data
        │     │     └── gopsutil/v3/process
        │     └── internal/provider/remote
        │           ├── internal/config
        │           ├── internal/data
        │           └── internal/msg
        ├── internal/ui
        │     ├── internal/data
        │     ├── internal/config
        │     ├── internal/msg
        │     ├── internal/ui/input
        │     ├── internal/ui/tabs
        │     ├── internal/ui/overlays
        │     ├── internal/ui/widgets
        │     └── lipgloss/v2, bubbletea/v2
        └── internal/util (leaf, no deps)
```

**Import cycle avoidance:** `internal/data` imports `internal/ui/input` for zone types, while `internal/ui` imports `internal/data` for state types. The `UIState.ZoneManager` field is typed as `interface{}` rather than `*input.ZoneManager` to break what would otherwise be a circular import (`data` -> `ui/input` -> `data`). This is a known concern (see CONCERNS.md C-022).

## Key File Locations

### Entry Points
- `cmd/bub/main.go` -- Standalone entry point (direct `tea.Program` launch)
- `cmd/bub/root.go` -- Cobra CLI entry point (`launchTUI()` at line 102, `Execute()` at line 70)

### BubbleTea Model
- `internal/app/model.go` -- `Model` struct definition, `Init()`, `View()`, factory functions
- `internal/app/update.go` -- `Update()` dispatcher (all message handling)

### State Definitions
- `internal/data/state.go` -- `AppState` root struct with mutex-protected methods
- `internal/data/groups.go` -- `MetricsState`, `ProcessState`, `UIState`, `ConfigState`, `RemoteState`
- `internal/data/types.go` -- `ProcessInfo`, `GpuInfo`, `Toast`, `HardwareCapabilities`, etc.

### Provider Contracts
- `internal/provider/interfaces.go` -- `SystemProvider`, `ProcessProvider`, `RemoteProvider` interfaces
- `internal/provider/adapters.go` -- Adapter structs with compile-time interface checks

### Configuration
- `internal/config/config.go` -- `AppConfig` struct, load/save/defaults
- `internal/config/watcher.go` -- Hot-reload file watcher
- `internal/config/options.go` -- Theme/border/rate enumerations

### Rendering
- `internal/ui/layout.go` -- Main render function (1070 lines)
- `internal/ui/styles.go` -- 32 theme palettes

### Business Logic
- `internal/data/logic.go` -- Process filtering, tree building, cache invalidation
- `internal/data/alerts.go` -- Threshold-based alert system
- `internal/data/ringbuffer.go` -- Circular buffer for metric history
- `internal/app/analysis.go` -- Health scoring, process history tracking

## Test File Locations

### Unit Tests
| File | What It Tests |
|------|--------------|
| `internal/app/update_integration_test.go` | Update loop: all message types, key sequences, error paths, startup flow, kill flow, tab switching, settings, render cache (1520 lines, comprehensive) |
| `internal/data/state_test.go` | AppState mutex-protected methods (concurrent access, process maps) |
| `internal/data/alerts_test.go` | AlertManager concurrent access and threshold checking |
| `internal/data/ringbuffer_test.go` | RingBuffer Push/Get/Max/Avg/Resize correctness |
| `internal/data/logic_bench_test.go` | Benchmarks for buildProcessTree (500 procs) and GetFilteredProcesses (500 procs) |
| `internal/provider/process/utils_test.go` | String interner and slice pool |
| `internal/ui/widgets/charts_test.go` | Chart rendering correctness |
| `internal/ui/widgets/charts_bench_test.go` | Chart rendering performance |
| `internal/ui/input/zones_test.go` | Zone manager registration, lookup, hover |

### Test Infrastructure
| File | Purpose |
|------|---------|
| `internal/app/testutil/model.go` | `TestModel` factory, `NewTestModel()`, `NewModelWithProcesses()`, `NewModelWithMetrics()` |
| `internal/app/testutil/mock_system.go` | `MockSystemProvider` with overridable function fields, mock data generators |
| `internal/app/testutil/mock_process.go` | `MockProcessProvider` with overridable function fields |
| `internal/app/testutil/mock_remote.go` | `MockRemoteProvider` with overridable function fields |

### Test Patterns
- Mock providers use function fields: each method checks if `XxxFunc` is set; if so, calls it; otherwise returns default mock data
- Integration tests in `update_integration_test.go` create a `Model` via `createTestModel()`, send messages via `m.Update(msg)`, then validate state changes
- Table-driven tests throughout: each test case has `name`, `msg`/`keyPress`, and `validate` function
- Benchmarks use `b.ResetTimer()` with pre-populated data (500 processes for tree/filter benchmarks)

### Untested Areas
- `internal/provider/system/` -- 15+ files with no test files (gopsutil wrappers, GPU detection, sysfs parsing)
- `internal/provider/remote/ssh.go` -- No tests for SSH output parsing
- `internal/config/` -- No tests for config loading/saving
- `internal/ui/tabs/` -- No tests for tab rendering
- `internal/ui/overlays/` -- No tests for overlay rendering
- `internal/app/handlers/mouse.go` -- No dedicated tests (mouse handling tested indirectly)

## Naming Conventions

### Files
- Go packages map 1:1 to directories (standard Go convention)
- Test files: `*_test.go` suffix, co-located with source
- Benchmark files: `*_bench_test.go` suffix
- Build-tag variants: `*_stub.go` for no-op fallbacks (e.g., `hardware_nvidia_stub.go`)
- Platform files: named by feature (`hardware_nvidia.go`, `detect_nvidia.go`)

### Directories
- `internal/` -- Private packages (not importable outside this module)
- `internal/app/handlers/` -- Extracted handler logic (functions, not methods)
- `internal/app/testutil/` -- Test utilities (mock providers, model factories)
- `internal/provider/{system,process,remote}/` -- Concrete implementations per domain
- `internal/ui/{tabs,overlays,widgets,input}/` -- UI decomposition by rendering concern

## Where to Add New Code

### New Tab
1. Create renderer: `internal/ui/tabs/<tabname>.go` implementing `Render<TabName>(s *data.AppState, ...) string`
2. Add tab name to `AllAvailableTabs` in `internal/data/constants.go`
3. Add tab name to default `Tabs` in `internal/config/config.go` `DefaultConfig()`
4. Add `case "<TabName>":` in `internal/ui/layout.go` (the switch at line 614) to call the new renderer
5. Add `case "<TabName>":` in `internal/app/update.go` `TickMsg` handler for tab-specific provider commands
6. Add scroll handling in `internal/app/handlers/keyboard.go` (if scrollable)
7. Add scroll handling in `internal/app/handlers/mouse.go` (if scrollable)

### New Message Type
1. Define message struct in `internal/msg/messages.go`
2. Add `case messages.<NewMsg>:` handler in `internal/app/update.go` `Update()` method
3. Add provider command that produces the message (in appropriate `internal/provider/*/` package)
4. Dispatch the command from `Init()` or `TickMsg` handler as needed

### New Provider Method
1. Add method to the appropriate interface in `internal/provider/interfaces.go`
2. Add delegation method to the adapter struct in `internal/provider/adapters.go`
3. Implement the concrete function in the appropriate `internal/provider/*/` package
4. Update mock in `internal/app/testutil/mock_<provider>.go` (manually)

### New Theme
1. Add `ThemePalette` entry to the `themes` map in `internal/ui/styles.go`
2. Add theme name to the list in `internal/config/options.go` `GetThemeNames()`

### New Overlay
1. Create renderer: `internal/ui/overlays/<name>.go` implementing `Render<Name>(...) string`
2. Add layer assembly in `internal/ui/layout.go` `MainViewFromState()` with Z-order
3. Add toggle flag to `UIState` in `internal/data/groups.go`
4. Add keyboard handler in `internal/app/handlers/keyboard.go`

### New Setting
1. Add field to `AppConfig` in `internal/config/config.go`
2. Add field to `ConfigState` in `internal/data/groups.go`
3. Add field to the corresponding state sub-struct if needed
4. Update `TotalSettingsCount` / related constants in `internal/data/constants.go`
5. Add handler case in `handleSettingsChange()` in `internal/app/handlers/keyboard.go`
6. Add rendering in `internal/ui/overlays/settings.go`

### New Utility Function
- Add to the appropriate file in `internal/util/` (format, math, text, layout, or fast)
- No internal imports allowed (leaf package)

## Special Directories

### `dist/`
- Purpose: GoReleaser build output (compiled binaries)
- Generated: Yes (by `goreleaser release`)
- Committed: No (should be in `.gitignore`)

### `.planning/codebase/`
- Purpose: Codebase analysis documents consumed by `/gsd-plan-phase` and `/gsd-execute-phase`
- Generated: Yes (by `/gsd-map-codebase`)
- Committed: Yes (tracked in git for reference)

### `internal/app/testutil/`
- Purpose: Test-only code (mock providers, test model factories)
- Import path: `github.com/N1xev/bubbleMonitor/internal/app/testutil`
- Note: Uses `package testutil` (separate from `package app`) to avoid import cycles

---

*Structure analysis: 2026-04-09*

# Architecture

**Analysis Date:** 2026-04-09

## Pattern Overview

**Overall:** BubbleTea (Elm Architecture) -- Model/Update/View

The application follows the Elm Architecture implemented via BubbleTea v2 (`charm.land/bubbletea/v2`). The core loop: messages (`tea.Msg`) flow into `Update()`, which mutates state and returns commands (`tea.Cmd`). Commands execute asynchronously, produce new messages, and the cycle repeats.

**Key Characteristics:**
- Unidirectional data flow: providers emit messages, Update mutates state, View reads state
- All provider operations return `tea.Cmd` (functions that produce `tea.Msg`), never mutate state directly
- State is centralized in a single `AppState` struct, embedded in the top-level `Model`
- Rendering is decoupled from state: `ui.RenderFromAppState()` takes `*AppState` and produces a view
- Render throttling via `RenderCache` (33ms / ~30 FPS) with forced invalidation on key/mouse events
- Config hot-reload via filesystem modification time polling (2-second interval)
- Zone-based mouse interaction system with click handlers and hover detection
- Hardware capability detection gates which provider commands run (avoids wasted work on headless/bare systems)
- Provider abstraction via interfaces + adapters enables test mocking (compile-time checks at `internal/provider/adapters.go:200-204`)

## Layers

### Provider Layer (Data Sources)
- Purpose: Fetch system metrics, process data, and remote host data from the OS and network
- Location: `internal/provider/`
- Contains: Interface definitions (`interfaces.go`), adapter implementations (`adapters.go`), and concrete provider packages (`system/`, `process/`, `remote/`)
- Depends on: `gopsutil/v3`, `battery`, `NVIDIA/go-nvml`, `hhk7734/amdsmi.go`, SSH (`os/exec`)
- Used by: `internal/app` (Model)

### Data Layer (State + Types + Business Logic)
- Purpose: Central state definition, data structures, business logic (filtering, tree building, health scoring, alerts)
- Location: `internal/data/`
- Contains:
  - `state.go` -- `AppState` root struct with mutex-protected methods for process maps
  - `groups.go` -- Sub-state struct definitions (`MetricsState`, `ProcessState`, `UIState`, `ConfigState`, `RemoteState`)
  - `types.go` -- Value types (`ProcessInfo`, `GpuInfo`, `Toast`, `ServiceInfo`, `ConnectionInfo`, `HardwareCapabilities`, container/VM types)
  - `alerts.go` -- `AlertManager` with threshold-based alert checking
  - `ringbuffer.go` -- Thread-safe circular buffer with cached Max/Avg
  - `logic.go` -- Process filtering, tree building, visible process caching
  - `constants.go` -- Layout constants, health thresholds, overlay dimensions
  - `viewport.go` -- `SimpleViewport` wrapper around `bubbles/v2/viewport`
- Depends on: `internal/config` (for threshold types), `gopsutil/v3` (for raw metric types), `internal/ui/input` (for zone types)
- Used by: `internal/app`, `internal/ui`, `internal/msg`, `internal/provider/process`

### Message Layer
- Purpose: Define all BubbleTea message types that flow between providers and the update loop
- Location: `internal/msg/messages.go`
- Contains: ~25 message types organized as:
  - **Tick/Metrics:** `TickMsg`, `CpuMemMsg`, `DiskNetMsg`, `MetricsMsg`
  - **System Info:** `HostInfoMsg`, `DiskInfoMsg`, `GpuInfoMsg`, `DiskIOMsg`, `TempMsg`, `NetworkInterfacesMsg`, `BatteryMsg`, `ServicesMsg`, `ConnectionsMsg`, `SysLogMsg`
  - **Process:** `ProcessesMsg`, `ProcessCountMsg`, `OpenFilesMsg`, `OpenFilesRequestMsg`, `ProcessCmdlineMsg`, `ProcessUsernameMsg`
  - **Process Control:** `KillProcessMsg`, `PriorityChangeMsg`, `ProcessControlMsg`
  - **Remote:** `RemoteMsg`
  - **UI/Control:** `ToastMsg`, `ToastTimeoutMsg`, `QuitMsg`, `ExportSnapshotMsg`, `ForceRefreshMsg`
  - **Container/VM:** `ContainerInfoMsg`, `VmInfoMsg`
- Depends on: `internal/data` (for shared data types like `ProcessInfo`, `GpuInfo`), `gopsutil/v3` (raw types), `battery`
- Used by: All packages that produce or consume messages

### Application Layer (Model/Update/View)
- Purpose: Orchestrate the BubbleTea update loop, bind providers to the model, manage configuration
- Location: `internal/app/`
- Contains:
  - `model.go` -- `Model` struct (embeds `AppState`), `Init()`, `View()`, `InitialModel()`, `InitialModelWithConfig()`, `RenderCache`
  - `update.go` -- `Update()` method: 422-line type-switch dispatcher for all message types
  - `handlers/keyboard.go` -- `HandleKey()` function for all keyboard input (939 lines)
  - `handlers/mouse.go` -- `HandleMouse()` for click, right-click, scroll, and motion events
  - `handlers/toast.go` -- Toast creation and timeout handling
  - `handlers/config.go` -- Settings change logic
  - `analysis.go` -- `UpdateHealthScore()` and `UpdateProcessHistory()`
  - `export.go` -- `SaveSnapshotCmd()` for JSON/CSV export
  - `logging.go` -- `LogMetricsCmd()` for optional file logging
  - `testutil/` -- Test model factories and mock providers
- Depends on: `internal/provider`, `internal/data`, `internal/msg`, `internal/config`, `internal/ui`
- Used by: `cmd/bub` (entry point)

### UI Layer (Rendering)
- Purpose: Transform `AppState` into terminal output using lipgloss styles and compositor
- Location: `internal/ui/`
- Contains:
  - `layout.go` -- Main render function `MainViewFromState()`: builds header, tabs, content, footer, overlays, and assembles compositor layers
  - `styles.go` -- Theme palette definitions (32 themes) and `GetAppTheme()` resolver
  - `constants.go` -- Window size minimums, responsive thresholds
  - `tabs/` -- Per-tab render functions (`metrics.go`, `processes.go`, `system.go`, `disks.go`, `network.go`, `services.go`, `connections.go`, `logs.go`, `remote.go`, `overview.go`, `helpers.go`)
  - `overlays/` -- Overlay renderers (`help.go`, `settings.go`, `kill.go`, `files.go`, `samlab.go`)
  - `widgets/` -- Reusable UI primitives (`charts.go`, `progress.go`, `borders.go`, `constants.go`)
  - `input/` -- Zone management system (`zones.go`, `types.go`)
- Depends on: `internal/data`, `internal/config`, `internal/msg`, `lipgloss/v2`, `bubbletea/v2`
- Used by: `internal/app` (`Model.View()`)

### Configuration Layer
- Purpose: Load, save, validate, and hot-reload JSON configuration
- Location: `internal/config/`
- Contains:
  - `config.go` -- `AppConfig` struct, `LoadConfig()`, `SaveConfig()`, `DefaultConfig()`, `LoadConfigFromPath()`, `ResolvePath()`
  - `watcher.go` -- `WatchConfig()` poll-based file change detection, `ConfigChangeMsg`, `ConfigWatchTickMsg`
  - `options.go` -- `GetThemeNames()`, `GetBorderStyles()`, `GetBorderTypes()`, `GetRefreshRates()`, `DefaultCustomTheme()`
- Depends on: `bubbletea/v2` (for `tea.Cmd` in watcher)
- Used by: `internal/app`, `internal/data`, `internal/ui`, `cmd/bub`

### Utility Layer
- Purpose: Shared formatting, layout, math, and text helpers
- Location: `internal/util/`
- Contains: `format.go`, `layout.go`, `math.go`, `text.go`, `fast.go`
- Depends on: Nothing internal (leaf package)
- Used by: UI layer, data layer, providers

## Data Flow

### Primary Tick Cycle

```
1. TickMsg fires (interval from config.RefreshRate, default 1000ms)
2. Update() increments TickCount
3. AlertManager.CheckAlerts() evaluates thresholds against current metrics
4. Provider.System.TickCmd() schedules next tick
5. Provider.System.FastMetricsCmd() dispatched (CPU, memory, swap)
6. Tab-specific commands dispatched based on active tab:
   - "Network"     -> NetworkInterfacesCmd
   - "Disks"       -> DiskIOCmd
   - "System"      -> TempCmd + HostInfoCmd + GpuInfoCmd
   - "Metrics"     -> DiskIOCmd + NetworkInterfacesCmd + TempCmd
   - "Processes"   -> ProcessesCmd + lazy cmdline/username fetch
   - "Connections" -> ConnectionsCmd
7. Every other tick (TickCount%2==0):
   - SlowMetricsCmd() (disk usage, network aggregate)
   - HostInfoCmd()
   - BatteryCmd() (if HasBattery)
   - ServicesCmd() (if tab=="Services" && HasServices)
   - SystemLogsCmd() (if tab=="Logs")
   - CheckRemoteCmd() (if tab=="Remote") for each configured host
8. Other tabs get ProcessCountCmd() (lightweight) every other tick
```

### Metrics Collection Flow

```
Update() dispatches tea.Cmd (e.g., FastMetricsCmd)
  -> system.FastMetricsCmd() runs in goroutine:
     - cpu.Percent(0, false) -> aggregate CPU%
     - cpu.Percent(0, true) -> per-core CPU%
     - mem.VirtualMemory() -> RAM stats
     - mem.SwapMemory() -> swap stats
     - load.Avg() -> load averages
     - Returns CpuMemMsg{Cpu, CpuPerCore, Memory, Swap, LoadAvg, MemInfo, SwapInfo, Err}
  -> Update() receives CpuMemMsg:
     - Mutates m.Metrics.Cpu, .Memory, .Swap, .CpuPerCore, .LoadAvg, .MemInfo, .SwapInfo
     - Pushes into ring buffers: CpuHistory, MemHistory, SwapHistory
     - Calls UpdateHealthScore()
```

### Process Management Flow

```
Update() dispatches Process.ProcessesCmd(sortBy, sortDirection)
  -> process.ProcessesCmd() runs in goroutine:
     - Prevents concurrent execution via atomic.Bool (isFetching)
     - Fetches PIDs via gopsutil process.Pids()
     - Checks processCache map for static data (Name, CreateTime, Nice, Ppid)
     - Creates CachedProcessInfo for new processes
     - Fetches dynamic data per process: CPUPercent, MemoryPercent, Status, MemoryInfo(RSS)
     - Uses string interner for process names (LRU, 5000 cap)
     - Uses sync.Pool for process slice reuse
     - Sorts by requested field (cpu/mem/pid/name) and direction
     - Prunes dead processes from processCache
     - Returns ProcessesMsg([]ProcessInfo)
  -> Update() receives ProcessesMsg:
     - Recycles previous slice via PutProcSlice
     - Replaces m.Process.Processes
     - Syncs ProcessesByPid map
     - Prunes dead process maps (SuspendedState, CollapsedPids, BookmarkedPids)
     - Invalidates process render cache (ProcessCacheDirty = true)
     - Updates health score and process history
     - Clamps SelectedProcess to valid range
```

### View Rendering Flow

```
BubbleTea calls Model.View()
  -> Render cache check: if !Force and < 33ms since last render, return cached content
  -> Otherwise: ui.RenderFromAppState(&m.AppState)
  -> MainViewFromState(s, getBorder, getColors):
     1. Create/clear ZoneManager
     2. Resolve theme palette from s.Config.Theme
     3. Rebuild style cache if theme changed (box, title, dim, warn, header, key, val, activeTab, tab)
     4. Minimum dimension check (80x24) -- shows "WINDOW TOO SMALL" if below
     5. Build top bar:
        - Header: "BUBBLE MONITOR" + clock + alert indicator (blinking if active)
        - Tab buttons: registered as click zones (ZoneTypeTab), hover-highlighted
     6. Build footer:
        - Context-sensitive buttons (help, settings, filter, kill, sort, quit)
        - Registered as click zones (ZoneTypeButton)
        - SamLab link as click zone (ZoneTypeLink)
     7. Dispatch content rendering by active tab:
        - "Metrics"    -> tabs.RenderMetrics() -- charts for CPU, mem, disk, net, temp
        - "Processes"  -> tabs.RenderProcesses() -- scrollable process list with tree/flat view
        - "System"     -> tabs.RenderSystem() -- multi-block layout with per-block scrolling
        - "Disks"      -> tabs.RenderDisks() -- partition table
        - "Network"    -> tabs.RenderNetwork() -- interface stats
        - "Services"   -> tabs.RenderServices() -- systemd service list
        - "Connections"-> tabs.RenderConnections() -- network connections table
        - "Logs"       -> tabs.RenderLogs() -- journalctl output
        - "Remote"     -> tabs.RenderRemote() -- remote host cards
     8. Assemble layers via lipgloss.NewCompositor:
        Z=0:  Base layer (header + content + footer)
        Z=2:  Toast overlays
        Z=3:  Kill dialog overlay
        Z=4:  Help overlay / Open files overlay
        Z=5:  Settings overlay / SamLab overlay
        Z=10: Context menu (right-click on process)
     9. Return tea.View with AltScreen=true, MouseMode=MouseModeAllMotion
```

### Config Hot-Reload Flow

```
Init() starts WatchConfig(lastModTime)
  -> WatchConfig returns tea.Cmd (polls file mtime)
  -> If config mtime changed:
     -> Returns ConfigChangeMsg{NewModTime}
     -> Update() receives ConfigChangeMsg:
        - Reloads config from disk via config.LoadConfig()
        - Diffs against current state via reflect.DeepEqual
        - If unchanged, reschedules WatchConfig only
        - If changed:
          - Applies all config fields to state
          - Resizes ring buffers if HistoryLength changed
          - Invalidates process cache if TreeView changed
          - Updates active tabs, theme, refresh rate, borders, etc.
          - Shows "Config Reloaded" toast
        - Reschedules WatchConfig
  -> If no change:
     -> Returns ConfigWatchTickMsg
     -> Update() reschedules WatchConfig
```

Polling interval: 2 seconds (see `internal/config/watcher.go:35`)

### State Management

- **Single source of truth:** `AppState` struct in `internal/data/state.go` holds all application state
- **Embedded in Model:** `Model` embeds `AppState` directly, so handler functions can take `*data.AppState` without importing the app package
- **Concurrency model:** BubbleTea calls `Update()` sequentially on the main goroutine. Provider `Cmd` functions run in background goroutines. The `AppState.stateMu` (sync.RWMutex) protects process maps that could theoretically be accessed concurrently, though the current architecture ensures sequential access via the BubbleTea loop. Methods using the mutex: `SetSuspended`, `IsSuspended`, `ToggleCollapsed`, `IsCollapsed`, `ToggleBookmark`, `IsBookmarked`, `GetProcessByPid`, `PruneDeadProcessMaps`, `PruneDeadProcessHistory`, `GetOrCreateHistory`, `GetHistory`
- **Process caching:** `CachedVisibleProcs` + `ProcessCacheDirty` flag for lazy recomputation of filtered/visible process list. `ProcessesByPid` map synced via `SyncProcessesMap()`
- **Ring buffers:** Pre-allocated fixed-size arrays with cursor wrapping. Cached Max and Sum with dirty flags for O(1) amortized access. Thread-safe via per-buffer `sync.RWMutex`

## Key Abstractions

### Provider Interfaces (`internal/provider/interfaces.go`)

**SystemProvider** (17 methods):
- Methods return `tea.Cmd` (async) or `bool` (hardware detection)
- Data methods: `TickCmd`, `FastMetricsCmd`, `SlowMetricsCmd`, `HostInfoCmd`, `GpuInfoCmd`, `TempCmd`, `BatteryCmd`, `DiskInfoCmd`, `DiskIOCmd`, `NetworkInterfacesCmd`, `ConnectionsCmd`, `ServicesCmd`, `SystemLogsCmd`
- Capability methods: `HasNvidiaGPU`, `HasAmdGPU`, `HasBattery`, `HasNetworkInterfaces`, `HasDiskIO`, `HasServices`, `HasTempSensors`
- Lifecycle: `DetectHardware`
- Adapter: `SystemAdapter` in `internal/provider/adapters.go` (wraps package-level functions)
- Implementation: `internal/provider/system/` (concrete package functions using gopsutil, NVML, sysfs)

**ProcessProvider** (9 methods):
- Data methods: `ProcessesCmd(sortBy, sortDirection)`, `ProcessCountCmd`, `PidsOnlyCmd`
- Control methods: `ReniceProcessCmdSafe(pid, delta)`, `SuspendProcessCmd(pid)`, `ResumeProcessCmd(pid)`
- Lazy-load methods: `FetchOpenFilesCmd(pid)`, `FetchProcessCmdlineCmd(pid)`, `FetchProcessUsernameCmd(pid)`
- Adapter: `ProcessAdapter` in `internal/provider/adapters.go`
- Implementation: `internal/provider/process/list.go`, `control.go`, `files.go`

**RemoteProvider** (1 method):
- `CheckRemoteCmd(host RemoteHostConfig)` -- SSH-based remote metrics collection
- Adapter: `RemoteAdapter` in `internal/provider/adapters.go`
- Implementation: `internal/provider/remote/ssh.go`

### RingBuffer (`internal/data/ringbuffer.go`)
- Fixed-size circular buffer for metric history (used for all charts)
- Thread-safe via `sync.RWMutex`
- Cached Max and Sum with dirty flags for O(1) amortized access
- `Resize(newSize)` preserves newest data
- `Accessor` interface: `Len()`, `Get(i)`, `Max()` -- used by chart rendering

### AlertManager (`internal/data/alerts.go`)
- Threshold-based alert system for CPU, memory, disk, temperature
- `map[MetricType]Alert` with thread-safe access via `sync.RWMutex`
- `CheckAlerts(state)` evaluates current metrics against configured thresholds
- Used by UI to display blinking warning indicator in header

### ZoneManager (`internal/ui/input/types.go`, `zones.go`)
- Mouse interaction regions for clickable elements
- Interface: `ZoneManager{ Register, FindZoneAt, GetHoveredZone, IsHovered, UpdateMousePos, Clear, GetZones }`
- Zone types: `ZoneTypeTab`, `ZoneTypeButton`, `ZoneTypeListItem`, `ZoneTypeLink`, `ZoneTypeMenuItem`
- Zones registered during View() rendering with absolute positions
- Mouse events resolved against zones in `handlers/mouse.go`
- Stored as `interface{}` in `UIState.ZoneManager` to avoid circular import (known concern -- see CONCERNS.md C-022)

## Entry Points

### TUI Mode (Primary)
- Location: `cmd/bub/main.go:14-25` (standalone), `cmd/bub/root.go:102-113` (via Cobra)
- Flow (standalone): `main()` -> `app.InitialModel()` -> `tea.NewProgram()` -> `p.Run()`
- Flow (Cobra): `main()` -> `Execute()` -> `newRootCmd()` -> default RunE -> `launchTUI()` -> `loadConfigWithOverrides()` -> `app.InitialModelWithConfig(cfg)` -> `tea.NewProgram()` -> `p.Run()`
- Config loading: CLI flags (`--config`, `--theme`, `--refresh-rate`, `--history-length`) override JSON config values

### Subcommand Mode
- Location: `cmd/bub/root.go:55-65`
- Subcommands: version, status, sysinfo, top, ps, health, export, doctor, themes, config, remote
- Each subcommand defined in a separate file in `cmd/bub/`

## Error Handling

**Strategy:** Toast-based user notification; errors never crash the application.

**Patterns:**
- Provider commands return error field in message structs (e.g., `CpuMemMsg.Err`, `DiskNetMsg.Err`, `HostInfoMsg.Err`)
- `Update()` checks `.Err != nil` and creates a toast: `handlers.AddToastCmd("error description", data.ToastError)`
- Toasts have levels: `ToastInfo`, `ToastError`, `ToastWarn`, `ToastSuccess` (defined in `internal/data/types.go:56-60`)
- Toasts auto-expire via `ToastTimeoutMsg` after configurable `Duration`
- Toasts render as layered overlays in the compositor (Z=2)
- Config load failure falls back to `DefaultConfig()` silently
- No error causes the application to exit; all errors are non-fatal and surfaced to the user

## View Rendering Pipeline

**Render Throttle:** `Model.View()` (`internal/app/model.go:96`) caches output and skips re-render if <33ms since last render (~30 FPS cap). Key/mouse events set `renderCache.Force = true` to bypass throttle. Style cache rebuilt on theme change.

**Layer Assembly:** Uses `lipgloss.NewCompositor` with Z-ordered layers:
- Z=0: Base (header + tab bar + content + footer)
- Z=2: Toast notifications (bottom-right positioned)
- Z=3: Kill confirmation dialog (centered)
- Z=4: Help overlay / Open files overlay (centered)
- Z=5: Settings overlay / SamLab overlay (centered)
- Z=10: Context menu (positioned at mouse coordinates)

## Cross-Cutting Concerns

**Logging:** Optional file-based metrics logging via `internal/app/logging.go`. Controlled by `config.LoggingConfig.Enabled` and `.Path`. Writes one line per snapshot in format: `timestamp | CPU: X% | Mem: X% | Disk: X% | Net: X MB/s | Procs: N`.

**Validation:** Config validation happens at load time in `config.LoadConfigFromPath()` -- missing/zero fields get default values. No explicit range validation on numeric values.

**Authentication:** Remote monitoring uses SSH key-based auth. Config stores `key_path` and optional `port`/`timeout`. SSH commands use `BatchMode=yes` (no interactive password prompts).

**Process Control:** Kill (with confirmation dialog), suspend, resume, renice operations. All operations return error messages via toasts on failure.

**Health Scoring:** Weighted penalty system starting at 100. CPU, memory, disk, temperature each contribute configurable penalties via `HealthWeights`. Evaluated by `UpdateHealthScore()` in `internal/app/analysis.go`.

**Export:** Snapshot export to JSON + CSV files triggered by pressing 'e'. Writes to `~/bubble_snapshot.json` and `~/bubble_snapshot.csv` via `SaveSnapshotCmd()` in `internal/app/export.go`.

**Process Memory Optimization:**
- String interner (`internal/provider/process/utils.go`) -- LRU cache with 5000 entry cap, evicts 500 at a time
- `sync.Pool` for `[]ProcessInfo` slice reuse (`GetProcSlice`/`PutProcSlice`)
- `atomic.Bool` (`isFetching`) prevents concurrent process list fetches

---

*Architecture analysis: 2026-04-09*

# Coding Conventions

**Analysis Date:** 2026-04-09

## Naming Patterns

### Files

- **Lowercase, no underscores for general files:** `model.go`, `update.go`, `layout.go`, `styles.go`
- **Build-tag suffix files use underscores:** `hardware_nvidia.go`, `hardware_nvidia_stub.go`, `hardware_amd.go`, `hardware_amd_stub.go`, `detect_nvidia.go`, `detect_nvidia_stub.go`
- **Test files use `_test.go` suffix:** `ringbuffer_test.go`, `alerts_test.go`
- **Benchmark test files use `_bench_test.go` suffix:** `logic_bench_test.go`, `charts_bench_test.go`
- **Mock/test utility files prefixed `mock_`:** `mock_system.go`, `mock_process.go`, `mock_remote.go`

### Packages

- **Single-word, lowercase:** `data`, `config`, `util`, `msg`, `input`
- **Sub-packages under provider:** `provider/system`, `provider/process`, `provider/remote`
- **Sub-packages under ui:** `ui/widgets`, `ui/tabs`, `ui/overlays`, `ui/input`
- **CLI commands in `cmd/bub`:** package `main`
- **Internal packages use `internal/` prefix** to prevent external imports

### Types (Structs)

- **PascalCase:** `AppState`, `ProcessInfo`, `MetricsState`, `RingBuffer`, `Toast`
- **State groupings use `*State` suffix:** `MetricsState`, `ProcessState`, `UIState`, `ConfigState`, `RemoteState`
- **Message types use `*Msg` suffix:** `CpuMemMsg`, `DiskNetMsg`, `ProcessesMsg`, `TickMsg`, `KillProcessMsg`
- **Adapter types use `*Adapter` suffix:** `SystemAdapter`, `ProcessAdapter`, `RemoteAdapter`
- **Provider interfaces use `*Provider` suffix:** `SystemProvider`, `ProcessProvider`, `RemoteProvider`
- **Info types use `*Info` suffix:** `GpuInfo`, `ServiceInfo`, `ConnectionInfo`, `ContainerInfo`
- **Config types use `*Config` suffix:** `AppConfig`, `LoggingConfig`, `RemoteHostConfig`, `CustomThemeConfig`
- **Private types use camelCase:** `treeBuilder`, `internerEntry`, `renderCache`, `zoneManager`

### Functions and Methods

- **PascalCase for exported:** `NewRingBuffer`, `InitialModel`, `LoadConfig`, `DefaultConfig`
- **Constructor pattern uses `New*` prefix:** `NewRingBuffer(size int)`, `NewTestModel()`, `NewZoneManager()`, `NewSimpleViewport()`
- **Constructor with config uses `InitialModelWithConfig(cfg)`:** `internal/app/model.go`
- **Boolean getters use `Is*` prefix:** `IsSuspended(pid)`, `IsCollapsed(pid)`, `IsBookmarked(pid)`, `IsHovered(zoneID)`
- **Action methods use verb prefix:** `SetSuspended`, `ToggleCollapsed`, `ToggleBookmark`, `InvalidateProcessCache`
- **Command constructors use `*Cmd` suffix:** `TickCmd(d)`, `FastMetricsCmd()`, `ProcessesCmd()`, `KillProcessCmd(pid)`

### Variables and Constants

- **Exported constants are PascalCase:** `MaxHealthScore`, `HealthThresholdHealthy`, `ProtoTCP`, `ProtoUDP`
- **Private constants are camelCase:** `maxInternerSize` (in `internal/provider/process/utils.go`)
- **Package-level vars for singletons:** `var detectionDone atomic.Bool`, `var themes = map[string]ThemePalette{...}`
- **Constants grouped in `const` blocks by category:** See `internal/data/constants.go`, `internal/provider/constants.go`

### Type Aliases for Messages

Some message types are defined as type aliases rather than structs:
```go
type TickMsg time.Time                          // internal/msg/messages.go
type ProcessesMsg []data.ProcessInfo            // internal/msg/messages.go
type ProcessCountMsg int                        // internal/msg/messages.go
type DiskIOMsg map[string]disk.IOCountersStat   // internal/msg/messages.go
```

Complex messages with error fields use structs:
```go
type CpuMemMsg struct {
    LoadAvg    *load.AvgStat
    MemInfo    *mem.VirtualMemoryStat
    Err        error
    Cpu        float64
    Memory     float64
}
```

## Code Style

### Formatting

- **Tool:** `gofmt` and `goimports` (enforced via `.golangci.yml`)
- **All code is `gofmt`-formatted**

### Linting

- **Tool:** golangci-lint (configured in `.golangci.yml`)
- **Enabled linters:** `errcheck`, `gosimple`, `govet`, `ineffassign`, `staticcheck`, `unused`, `gofmt`, `goimports`, `revive`
- **errcheck:** checks type assertions (`check-type-assertions: true`), does not check blank assignments (`check-blank: false`)
- **govet:** `enable-all: true`
- **revive rules:** blank-imports, context-as-argument, dot-imports, error-return, error-strings, error-naming, receiver-naming, time-naming, errorf

### Import Organization

Imports follow this order (separated by blank lines):

1. **Standard library:** `"fmt"`, `"os"`, `"sync"`, `"time"`
2. **Third-party / framework:** `tea "charm.land/bubbletea/v2"`, `"github.com/shirou/gopsutil/v3/mem"`
3. **Internal project packages:** `"github.com/N1xev/bubbleMonitor/internal/data"`, `"github.com/N1xev/bubbleMonitor/internal/config"`

Example from `internal/app/model.go`:
```go
import (
    "time"

    tea "charm.land/bubbletea/v2"
    "charm.land/lipgloss/v2"
    configpkg "github.com/N1xev/bubbleMonitor/internal/config"
    "github.com/N1xev/bubbleMonitor/internal/data"
    "github.com/N1xev/bubbleMonitor/internal/provider"
    "github.com/N1xev/bubbleMonitor/internal/provider/system"
    "github.com/N1xev/bubbleMonitor/internal/ui"
    "github.com/shirou/gopsutil/v3/cpu"
)
```

**Import aliases used:**
- `tea "charm.land/bubbletea/v2"` -- universal alias for Bubble Tea
- `messages "github.com/N1xev/bubbleMonitor/internal/msg"` -- alias to avoid collision with `msg` parameter names
- `configpkg "github.com/N1xev/bubbleMonitor/internal/config"` -- alias to avoid collision with `cfg` variable names

## Error Handling

### Patterns

**Errors returned as struct fields in messages:**
```go
type CpuMemMsg struct {
    Err    error
    Cpu    float64
    Memory float64
    // ...
}
```

**Error handling in Update uses `msg.Err` check:**
```go
case messages.CpuMemMsg:
    if msg.Err != nil {
        m.UI.LastError = msg.Err.Error()
        m.UI.LastErrorTime = time.Now()
        return m, nil
    }
    // process successful data...
```
Reference: `internal/app/update.go`

**Error messages displayed as toasts:**
```go
return m, handlers.AddToastCmd(fmt.Sprintf("Priority Error: %v", msg.Err), data.ToastError)
```

**Error wrapping with `fmt.Errorf` and `%w`:**
```go
return fmt.Errorf("failed to load config: %w", err)
```
Reference: `cmd/bub/root.go`

**Panic recovery in hardware detection:**
```go
defer func() {
    if r := recover(); r != nil {
        log.Printf("NVIDIA GPU detection panicked: %v", r)
    }
}()
```
Reference: `internal/provider/system/detect.go`, `internal/provider/system/hardware_nvidia.go`

**Explicit error discarding with `_ =`:**
```go
defer func() { _ = file.Close() }()
```
Reference: `internal/config/config.go`

## Build Tags

Platform-specific and CGO-dependent code uses Go build tags:

```go
//go:build linux && cgo
```

Files using build tags:
- `internal/provider/system/hardware_nvidia.go` -- `//go:build linux && cgo`
- `internal/provider/system/hardware_nvidia_stub.go` -- `//go:build !linux || !cgo`
- `internal/provider/system/hardware_amd.go` -- `//go:build linux && cgo`
- `internal/provider/system/hardware_amd_stub.go` -- `//go:build !linux || !cgo`
- `internal/provider/system/detect_nvidia.go` -- `//go:build linux && cgo`
- `internal/provider/system/detect_nvidia_stub.go` -- `//go:build !linux || !cgo`

**Convention:** Pair each implementation file with a `_stub.go` that provides the no-op fallback. The stub file lives in the same package with the same function signatures.

## Concurrency Patterns

### Mutex-Protected State

`AppState` uses `sync.RWMutex` for concurrent map access:
```go
type AppState struct {
    Metrics MetricsState
    Process ProcessState
    // ...
    stateMu sync.RWMutex
}
```

- **Write operations use `Lock()`:** `SetSuspended`, `ToggleCollapsed`, `ToggleBookmark`
- **Read operations use `RLock()`:** `IsSuspended`, `IsCollapsed`, `IsBookmarked`, `GetHistory`
- Reference: `internal/data/state.go`

### Atomic Values for Hardware Detection

```go
var (
    detectionDone    atomic.Bool
    nvidiaDetected   atomic.Bool
    amdDetected      atomic.Bool
    // ...
)
```
Reference: `internal/provider/system/detect.go`

### RingBuffer Thread Safety

`RingBuffer` uses `sync.RWMutex` for all operations:
```go
type RingBuffer struct {
    data []float64
    mu   sync.RWMutex
    // ...
}
```
Reference: `internal/data/ringbuffer.go`

## Constants Organization

### Per-Package Constants

Constants are grouped in dedicated `constants.go` files within each package:

- `internal/data/constants.go` -- Layout constants, health scoring, settings overlay dimensions
- `internal/provider/constants.go` -- Provider limits, timeouts, capacities
- `internal/provider/system/constants.go` -- Protocol constants, max log lines
- `internal/provider/process/constants.go` -- Process list capacity
- `internal/provider/remote/constants.go` -- SSH timeout
- `internal/ui/constants.go` -- Layout dimensions, re-exports from data package
- `internal/config/config.go` -- MetricType enum values (`MetricCPU`, `MetricMem`, `MetricDisk`, `MetricTemp`)

### Constants are grouped in `const` blocks by category:

```go
const (
    ThresholdCount     = 4
    DisplayCount       = 6
    TabCount           = 9
    AppearanceCount    = 5
    TotalSettingsCount = ThresholdCount + DisplayCount + TabCount + AppearanceCount
)
```

## Interface Definition Patterns

### Provider Interfaces

Defined in `internal/provider/interfaces.go` with the following conventions:

- **Methods that return `tea.Cmd`:** All provider methods return `tea.Cmd` for async execution in the Bubble Tea runtime
- **Capability detection methods return `bool`:** `HasNvidiaGPU()`, `HasBattery()`, etc.
- **Comment per method:** Each interface method has a `//` doc comment

### Compile-Time Interface Checks

```go
var (
    _ SystemProvider  = (*SystemAdapter)(nil)
    _ ProcessProvider = (*ProcessAdapter)(nil)
    _ RemoteProvider  = (*RemoteAdapter)(nil)
)
```
Reference: `internal/provider/adapters.go`

### Adapter Pattern

Each provider interface has an `*Adapter` struct in `internal/provider/adapters.go` that delegates to package-level functions:
```go
type SystemAdapter struct{}
func (a *SystemAdapter) TickCmd(d time.Duration) tea.Cmd {
    return system.TickCmd(d)
}
```

## Bubble Tea Patterns

### Model Composition

The `Model` struct embeds `data.AppState` directly:
```go
type Model struct {
    data.AppState
    renderCache *RenderCache
    Provider    struct {
        System  provider.SystemProvider
        Process provider.ProcessProvider
        Remote  provider.RemoteProvider
    }
}
```
Reference: `internal/app/model.go`

### State Grouping

State is organized into domain-specific sub-structs:
- `MetricsState` -- system metrics and history
- `ProcessState` -- process list, sort, filter, tree view state
- `UIState` -- window size, tabs, dialogs, scroll offsets
- `ConfigState` -- theme, refresh rate, thresholds
- `RemoteState` -- remote host metrics

### Update Method Pattern

The `Update` method uses a type switch on `tea.Msg`:
```go
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        // handle
    case messages.CpuMemMsg:
        // handle
    // ...
    }
    return m, nil
}
```

### Handler Delegation

Key handling and mouse handling are delegated to separate handler packages:
- `internal/app/handlers/keyboard.go` -- `HandleKey(...)`
- `internal/app/handlers/mouse.go` -- `HandleMouse(...)`
- `internal/app/handlers/toast.go` -- `HandleToast(...)`, `AddToastCmd(...)`

### Command Batching

Multiple commands are batched using `tea.Batch`:
```go
return tea.Batch(
    m.Provider.Process.ProcessesCmd(m.Process.SortBy, m.Process.SortDirection),
    handlers.AddToastCmd("Priority Changed", data.ToastSuccess),
)
```

## Configuration Patterns

### Config Loading

Chain: `LoadConfig()` -> `GetConfigPath()` -> `LoadConfigFromPath(path)` with fallback to `DefaultConfig()`:
```go
func LoadConfig() (AppConfig, error) {
    path, err := GetConfigPath()
    if err != nil {
        return DefaultConfig(), err
    }
    return LoadConfigFromPath(path)
}
```
Reference: `internal/config/config.go`

### Config Defaults

`DefaultConfig()` returns a fully-populated `AppConfig`. Loading always merges with defaults for missing fields.

### Config Persistence

JSON format with `encoding/json`, stored at `~/.config/bubble-monitor/config.json`.

## Comment/Documentation Style

### Exported Types and Functions

Every exported type, function, and method has a `//` doc comment:
```go
// SystemProvider defines the interface for system metrics providers.
// Each method returns a tea.Cmd that produces a message when executed.
type SystemProvider interface { ... }

// TickCmd returns a command that sends tick messages at the specified duration.
func (a *SystemAdapter) TickCmd(d time.Duration) tea.Cmd { ... }
```

### Inline Comments

Used sparingly for non-obvious logic:
```go
// Re-export from data package for convenience
KillDialogDefaultWidth = data.KillDialogDefaultWidth
```

### Benchmark Target Comments

Benchmarks include baseline and target metrics in comments:
```go
// BenchmarkRenderSparkline tests memory allocations for sparkline rendering.
// Baseline: 565 allocs/op
// Target after optimization: < 113 allocs/op (80% reduction)
```
Reference: `internal/ui/widgets/charts_bench_test.go`

## Module Design

### No Barrel/Init Files

No `init()` functions are used. All initialization is explicit via `New*()` constructors.

### No Exported Package-Level Variables for Mutable State

Mutable state lives in structs. Package-level vars are used only for:
- Theme maps (`themes` in `internal/ui/styles.go`)
- Hardware detection flags using `atomic.Bool` (immutable after first detection)

### Package Dependencies Flow

```
cmd/bub -> internal/app -> internal/data, internal/provider, internal/ui, internal/config
                     internal/app/handlers -> internal/data, internal/msg
                     internal/provider/system -> internal/data, internal/msg
                     internal/provider/process -> internal/data, internal/msg
                     internal/ui -> internal/data, internal/config, internal/ui/input
```

---

*Convention analysis: 2026-04-09*

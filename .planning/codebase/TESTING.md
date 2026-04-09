# Testing Patterns

**Analysis Date:** 2026-04-09

## Test Framework

**Runner:**
- Go standard `testing` package
- No third-party assertion libraries (no testify, no gomega)
- All assertions use `t.Errorf`, `t.Error`, `t.Fatal`, `t.Fatalf`

**Run Commands:**
```bash
go test ./...                                    # Run all tests
go test -v ./internal/data/...                   # Verbose output for data package
go test -v -race ./...                           # Run with race detector (CI uses this)
go test -v ./internal/app/... -run Integration   # Run integration test group
go test -bench=. ./internal/data/...             # Run benchmarks
go test -bench=. -benchmem ./internal/ui/widgets/...  # Benchmarks with memory profile
```

## Test File Inventory

### Unit Tests

| File | Package | Tests | Focus |
|------|---------|-------|-------|
| `internal/data/ringbuffer_test.go` | `data` | 4 tests + 2 benchmarks | RingBuffer operations, concurrency, wrap-around |
| `internal/data/alerts_test.go` | `data` | 2 tests | AlertManager concurrent access, alert triggering |
| `internal/data/state_test.go` | `data` | 3 tests | AppState concurrent map access, process map pruning |
| `internal/data/logic_bench_test.go` | `data` | 2 benchmarks | Process tree building, filtering performance |
| `internal/provider/process/utils_test.go` | `process` | 5 tests + 1 benchmark | Interner, slice pool, concurrency |
| `internal/ui/widgets/charts_test.go` | `widgets` | 1 test (3 subtests) | Chart visual regression, determinism |
| `internal/ui/widgets/charts_bench_test.go` | `widgets` | 5 benchmarks | Chart rendering allocation profiling |
| `internal/ui/input/zones_test.go` | `input` | 13 tests | ZoneManager register, find, hover, clear, z-order |

### Integration Tests

| File | Package | Tests | Focus |
|------|---------|-------|-------|
| `internal/app/update_integration_test.go` | `app` | 8 test groups, ~50 subtests | Full Update() message handling, key sequences, flows |

### Benchmark Tests

| File | Package | Benchmarks | Focus |
|------|---------|------------|-------|
| `internal/data/ringbuffer_test.go` | `data` | `BenchmarkRingBufferPush`, `BenchmarkRingBufferConcurrentAccess` | Push throughput, concurrent read/write |
| `internal/data/logic_bench_test.go` | `data` | `BenchmarkBuildProcessTree`, `BenchmarkGetFilteredProcesses` | Process tree building, filtering with 500 processes |
| `internal/provider/process/utils_test.go` | `process` | `BenchmarkIntern` | String interning throughput |
| `internal/ui/widgets/charts_bench_test.go` | `widgets` | `BenchmarkRenderSparkline`, `BenchmarkRenderLineChart`, `BenchmarkRenderBrailleChart`, `BenchmarkRenderLineChartLarge`, `BenchmarkRenderAllCharts` | Rendering allocation tracking |

### Test Utilities

| File | Package | Purpose |
|------|---------|---------|
| `internal/app/testutil/model.go` | `testutil` | TestModel factory, NewTestModel(), NewModelWithProcesses(), NewModelWithMetrics() |
| `internal/app/testutil/mock_system.go` | `testutil` | MockSystemProvider, standalone mock data functions |
| `internal/app/testutil/mock_process.go` | `testutil` | MockProcessProvider with configurable function fields |
| `internal/app/testutil/mock_remote.go` | `testutil` | MockRemoteProvider with configurable function fields |

## Test File Organization

**Location:** Tests are co-located with source files in the same package.

**Naming:**
- Unit tests: `*_test.go` (e.g., `ringbuffer_test.go`, `alerts_test.go`)
- Benchmark tests: `*_bench_test.go` (e.g., `logic_bench_test.go`, `charts_bench_test.go`)
- Integration tests: `*_integration_test.go` (e.g., `update_integration_test.go`)
- Mock helpers: `mock_*.go` in `internal/app/testutil/` (not `_test.go` suffix -- accessible from test files)

**Structure:**
```
internal/
  app/
    update.go
    update_integration_test.go    # Integration tests for Update()
    testutil/
      model.go                    # TestModel factories (shared, non-_test.go)
      mock_system.go              # MockSystemProvider (shared, non-_test.go)
      mock_process.go             # MockProcessProvider (shared, non-_test.go)
      mock_remote.go              # MockRemoteProvider (shared, non-_test.go)
  data/
    ringbuffer.go
    ringbuffer_test.go            # Unit + benchmark tests
    alerts.go
    alerts_test.go                # Unit tests
    state.go
    state_test.go                 # Unit tests
    logic.go
    logic_bench_test.go           # Benchmark-only file
  provider/process/
    utils.go
    utils_test.go                 # Unit + benchmark tests
  ui/
    widgets/
      charts.go
      charts_test.go              # Visual regression tests
      charts_bench_test.go        # Benchmark tests
    input/
      zones.go
      zones_test.go               # Unit tests
```

## Test Structure

### Test Function Naming

Tests use descriptive `Test*` names that state what is being tested:

```go
func TestRingBufferBasicOperations(t *testing.T) { ... }
func TestRingBufferWrapAround(t *testing.T) { ... }
func TestRingBufferEmptyState(t *testing.T) { ... }
func TestRingBufferConcurrency(t *testing.T) { ... }
func TestAppStateConcurrentMapAccess(t *testing.T) { ... }
func TestAppStatePruneDeadProcessMaps(t *testing.T) { ... }
func TestAlertManagerConcurrentAccess(t *testing.T) { ... }
```

Benchmarks follow `Benchmark*` naming:
```go
func BenchmarkRingBufferPush(b *testing.B) { ... }
func BenchmarkBuildProcessTree(b *testing.B) { ... }
func BenchmarkRenderLineChart(b *testing.B) { ... }
```

### Table-Driven Tests

The primary testing pattern is table-driven with subtests:

```go
func TestMessageTypes(t *testing.T) {
    tests := []struct {
        name     string
        msg      tea.Msg
        validate func(*testing.T, *Model)
    }{
        {
            name: "TickMsg",
            msg:  messages.TickMsg(time.Now()),
            validate: func(t *testing.T, m *Model) {
                if m.UI.TickCount != 1 {
                    t.Errorf("expected TickCount 1, got %d", m.UI.TickCount)
                }
            },
        },
        // ...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            m := createTestModel()
            newModel, _ := m.Update(tt.msg)
            updatedModel := newModel.(*Model)
            tt.validate(t, updatedModel)
        })
    }
}
```
Reference: `internal/app/update_integration_test.go`

### Assertion Style

All assertions use the standard `testing` package directly:

```go
// Equality check
if m.Metrics.Cpu != 25.5 {
    t.Errorf("expected Cpu 25.5, got %f", m.Metrics.Cpu)
}

// Boolean assertion
if !m.Process.ProcessesLoaded {
    t.Error("expected ProcessesLoaded to be true")
}

// Fatal assertion (stops test immediately)
if rb == nil || *rb == nil {
    t.Fatal("GetProcSlice returned nil or nil slice")
}
```

**Pattern:** Use `t.Errorf` for non-fatal failures (test continues), `t.Fatal`/`t.Fatalf` when the test cannot proceed.

### Test Grouping

The integration test file uses a master test function that calls sub-groups:

```go
func TestIntegration(t *testing.T) {
    t.Run("MessageTypes", func(t *testing.T) {
        TestMessageTypes(t)
    })
    t.Run("KeySequences", func(t *testing.T) {
        TestKeySequences(t)
    })
    t.Run("ErrorPaths", func(t *testing.T) {
        TestErrorPaths(t)
    })
    // ...
}
```
Reference: `internal/app/update_integration_test.go`

## Mock/Test Utility Infrastructure

### Location

`internal/app/testutil/` -- a dedicated sub-package containing shared test utilities. Files are NOT named `_test.go` so they are importable by test files in other packages.

### TestModel Factory

`internal/app/testutil/model.go` provides pre-configured test model builders:

```go
// Creates a minimal valid TestModel with default config
func NewTestModel() *TestModel { ... }

// Creates a TestModel pre-populated with 5 sample processes
func NewModelWithProcesses() *TestModel { ... }

// Creates a TestModel pre-populated with metrics and history data
func NewModelWithMetrics() *TestModel { ... }
```

**Note:** The integration tests in `update_integration_test.go` define their own `createTestModel()` helper instead of using `testutil.NewTestModel()`. This is because the integration tests need the full `Model` (which embeds providers) rather than `TestModel` (which only embeds `AppState`).

### Mock Provider Pattern

Mocks use configurable function fields -- set a `*Func` field to override behavior, leave nil for default:

```go
type MockSystemProvider struct {
    CpuFunc     func() (float64, []float64, error)
    MemoryFunc  func() (*mem.VirtualMemoryStat, *mem.SwapMemoryStat, error)
    // ...
}

func (m *MockSystemProvider) Cpu() (float64, []float64, error) {
    if m.CpuFunc != nil {
        return m.CpuFunc()
    }
    cpuVal, _, _ := MockCpuMem()
    return cpuVal, MockCpuPerCore(), nil
}
```
Reference: `internal/app/testutil/mock_system.go`

### Standalone Mock Data Functions

The testutil package also provides standalone data factory functions:

```go
func MockCpuMem() (float64, *mem.VirtualMemoryStat, *mem.SwapMemoryStat)
func MockCpuPerCore() []float64
func MockProcesses() []data.ProcessInfo
func MockDisk() []data.DiskPartition
func MockNetwork() ([]net.IOCountersStat, uint64, uint64, error)
func MockTemp() []host.TemperatureStat
func MockBattery() []*battery.Battery
func MockHostInfo() *host.InfoStat
func MockGpuInfo() []data.GpuInfo
func MockCpuInfo() []cpu.InfoStat
```
Reference: `internal/app/testutil/mock_system.go`

## Concurrency Testing

Concurrency tests are a distinct pattern in this codebase, testing thread safety of shared data structures:

```go
func TestRingBufferConcurrency(t *testing.T) {
    rb := NewRingBuffer(100)
    var wg sync.WaitGroup
    numGoroutines := 10
    iterations := 1000

    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func(val float64) {
            defer wg.Done()
            for j := 0; j < iterations; j++ {
                rb.Push(val)
                _ = rb.Max()
                _ = rb.Avg()
            }
        }(float64(i))
    }

    wg.Wait()
    // assertions...
}
```

**Packages with concurrency tests:**
- `internal/data` -- RingBuffer, AppState maps, AlertManager
- `internal/provider/process` -- Interner

## Benchmark Tests

### Naming Convention

Benchmarks are named `Benchmark*` and placed in either:
- The same file as unit tests (e.g., `ringbuffer_test.go` contains `BenchmarkRingBufferPush`)
- Separate `*_bench_test.go` files (e.g., `logic_bench_test.go`, `charts_bench_test.go`)

### Benchmark Patterns

```go
func BenchmarkRenderSparkline(b *testing.B) {
    b.ReportAllocs()                          // Track memory allocations
    ring := data.NewRingBuffer(100)
    for i := 0; i < 100; i++ {
        ring.Push(float64(i % 100))
    }

    b.ResetTimer()                            // Exclude setup from timing
    for i := 0; i < b.N; i++ {
        _ = RenderSparkline(ring, 80, 1, c1, c2, 100.0, "default")
    }
}
```

```go
func BenchmarkRingBufferConcurrentAccess(b *testing.B) {
    rb := NewRingBuffer(1000)
    for i := 0; i < 100; i++ {
        rb.Push(float64(i))
    }

    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            rb.Push(42.0)
            _ = rb.Max()
            _ = rb.Avg()
        }
    })
}
```

### Benchmark Targets

Benchmarks include target metrics as comments:

```go
// Baseline: 565 allocs/op
// Target after optimization: < 113 allocs/op (80% reduction)
```

```go
// Baseline: 5,632 allocs/op
// Target after optimization: < 1,126 allocs/op (80% reduction)
```
Reference: `internal/ui/widgets/charts_bench_test.go`

## CI Test Configuration

**File:** `.github/workflows/ci.yml`

### Pipeline

```
lint (ubuntu-latest, Go 1.25, golangci-lint)
  -> test (matrix: ubuntu/macos/windows x Go 1.22/1.23/1.24/1.25, -race flag)
  -> release (on tags only, GoReleaser)
```

### Test Job Details

- **Matrix:** 3 OS (ubuntu-latest, macos-latest, windows-latest) x 4 Go versions (1.22, 1.23, 1.24, 1.25)
- **Race detector enabled:** `go test -v -race ./...`
- **Build verification:** `go build -v ./...` runs after tests
- **Fail-fast disabled:** `fail-fast: false` -- all matrix combinations run regardless of individual failures

### Lint Job Details

- **Tool:** golangci-lint-action@v6
- **Timeout:** 5 minutes
- **Config:** `.golangci.yml`

## Coverage Assessment

### Well-Tested Packages (HIGH)

| Package | Test Files | Coverage Focus |
|---------|-----------|----------------|
| `internal/data` | 5 test files (3 unit + 2 bench) | RingBuffer all operations, AppState concurrent map access, AlertManager, process tree building, filtering |
| `internal/app` | 1 integration test (1519 lines) | Full Update() message handling, key sequences, error paths, startup flow, kill flow, tab switch flow, settings flow, render cache |
| `internal/ui/widgets` | 2 test files (1 unit + 1 bench) | Chart visual regression (determinism), all chart renderers benchmarked with alloc tracking |
| `internal/ui/input` | 1 test file (290 lines) | ZoneManager register/find/hover/clear/z-order/metadata/edge cases |

### Moderately Tested Packages (MEDIUM)

| Package | Test Files | Coverage Focus | Gaps |
|---------|-----------|----------------|------|
| `internal/provider/process` | 1 test file | Interner, slice pool, concurrency | `list.go`, `control.go`, `files.go` not tested |

### Untested Packages (LOW/NONE)

| Package | Test Files | Notes |
|---------|-----------|-------|
| `internal/provider/system` | None | Hardware detection, metrics collection, network, disk -- all untested |
| `internal/provider/remote` | None | SSH connectivity untested |
| `internal/config` | None | Config loading, saving, defaults untested |
| `internal/util` | None | Math, format, layout, text helpers untested |
| `internal/ui` (root) | None | `layout.go` (1069 lines), `styles.go` (455 lines) untested |
| `internal/ui/tabs` | None | All tab rendering untested |
| `internal/ui/overlays` | None | Settings overlay untested |
| `internal/app/handlers` | None | Keyboard (939 lines), mouse (228 lines) handler logic untested directly |
| `internal/msg` | None | Message type definitions (no logic to test) |
| `cmd/bub` | None | CLI commands untested |

## Testing Guidance for New Code

### When Adding a New Message Type

1. Add the type to `internal/msg/messages.go`
2. Add a handler case in `internal/app/update.go`
3. Add a table-driven test case in `internal/app/update_integration_test.go` following the `TestMessageTypes` pattern:
   ```go
   {
       name: "MyNewMsg with valid data",
       msg: messages.MyNewMsg{ /* data */ },
       validate: func(t *testing.T, m *Model) {
           // assertions
       },
   },
   {
       name: "MyNewMsg with error",
       msg: messages.MyNewMsg{Err: errors.New("mock error")},
       validate: func(t *testing.T, m *Model) {
           // verify error handling
       },
   },
   ```

### When Adding a New Data Structure

1. Create the type in `internal/data/`
2. Add a `Test*Concurrency` test if the type uses maps or shared state
3. Add a benchmark in a `*_bench_test.go` file if the type is performance-sensitive
4. Use `b.ReportAllocs()` and document baseline/target metrics in comments

### When Adding a New Provider Method

1. Add the method to the interface in `internal/provider/interfaces.go`
2. Add the adapter method in `internal/provider/adapters.go`
3. Add the implementation in the appropriate provider package
4. Create a mock method in `internal/app/testutil/mock_*.go` with configurable function field

### When Adding a New UI Component

1. Add chart/widget tests in `internal/ui/widgets/` following the visual regression pattern
2. Use determinism checks: render twice, assert equal output
3. Add a benchmark with `b.ReportAllocs()` for allocation tracking

---

*Testing analysis: 2026-04-09*

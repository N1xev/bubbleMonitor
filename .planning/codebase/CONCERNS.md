# Codebase Concerns

**Analysis Date:** 2026-04-09

## Tech Debt

### C-001: Giant Update Function (God Function)

- **Severity:** HIGH
- **Category:** Maintainability
- **Location:** `internal/app/update.go` (lines 20-442, 422 lines in single function)
- **Description:** The `Model.Update()` method is a 422-line type-switch monolith handling every message type in the application. Every new feature, tab, or message type requires modification to this single function. The function mixes business logic (health scoring, process history), UI logic (tab selection, dialog state), and data transformation (rate calculations, parsing).
- **Impact:** Adding new features is error-prone. Changes to one message handler can accidentally break unrelated handlers. Testing individual message handlers in isolation requires constructing the full Model. Code review is difficult due to sheer size.
- **Recommendation:** Split into per-message handler functions following the pattern already used in `internal/app/handlers/` for keyboard and mouse. Create a dispatcher map or handler interface. Example structure: each message type gets its own `handleXMsg(m *Model, msg XMsg) (tea.Model, tea.Cmd)` function.

### C-002: Duplicated Keyboard Handler Logic (939 lines)

- **Severity:** HIGH
- **Category:** Maintainability
- **Location:** `internal/app/handlers/keyboard.go` (lines 1-939)
- **Description:** The `HandleKey` function is 939 lines with massive duplication in scroll handling. The scroll-up/scroll-down/page-up/page-down/home/end logic for Services, Connections, Logs, and System tabs is nearly identical code blocks repeated 6+ times. Each tab has the same pattern: calculate `rows`, `maxScroll`, and update offset.
- **Impact:** Bug fixes must be applied in multiple places. Adding a new scrollable tab requires copy-pasting ~40 lines. Inconsistencies easily creep in.
- **Recommendation:** Extract a generic scroll handler struct or function:
  ```go
  type Scrollable struct { Offset, MaxOffset, ItemCount, VisibleRows int }
  func (s *Scrollable) Up(n int)    { s.Offset = max(s.Offset - n, 0) }
  func (s *Scrollable) Down(n int)  { s.Offset = min(s.Offset + n, s.MaxOffset) }
  ```
  Replace all duplicated scroll blocks with calls to this abstraction.

### C-003: Hard-Coded PCI ID Database (GPU Detection)

- **Severity:** MEDIUM
- **Category:** Maintainability
- **Location:** `internal/provider/system/hardware.go` (lines 472-595)
- **Description:** Two large hard-coded maps (`pciNames` ~40 entries, `pciVRAM` ~20 entries) map PCI vendor:device IDs to GPU names and VRAM amounts. This list is incomplete and will become stale as new GPUs release. The VRAM values are approximations that may not match actual card variants (e.g., GTX 1060 comes in 3GB and 6GB variants).
- **Impact:** New GPUs show as "Unknown GPU". Incorrect VRAM is displayed. Each new GPU requires a code change and release.
- **Recommendation:** Consider using the Linux `pci.ids` database file (usually at `/usr/share/hwdata/pci.ids`) for name lookups, with the hard-coded map as fallback only. For VRAM, prefer runtime detection via sysfs (`mem_info_vram_total`) and only fall back to static data.

### C-004: Ad-Hoc JSON Parsing in lsblk Handler

- **Severity:** MEDIUM
- **Category:** Maintainability
- **Location:** `internal/provider/system/disk.go` (lines 63-120)
- **Description:** The `parseLsblkUnmounted` function parses JSON output from `lsblk -J` using string splitting on newlines and colon separators. It manually tracks `currentName`, `currentSize`, `currentType`, `currentMountpoint` across lines. This is fragile and will break if `lsblk` output format changes (different indentation, ordering, or whitespace).
- **Impact:** Unmounted partitions may not display correctly. Silent parsing failures produce no error; data is simply missing.
- **Recommendation:** Use `encoding/json` with a proper struct to unmarshal `lsblk -J` output. The output format is documented and stable. Define:
  ```go
  type LsblkOutput struct {
      BlockDevices []struct {
          Name       string `json:"name"`
          Size       uint64 `json:"size,string"`
          Type       string `json:"type"`
          Mountpoint *string `json:"mountpoint"`
      } `json:"blockdevices"`
  }
  ```

### C-005: Kill Dialog Code Duplication

- **Severity:** LOW
- **Category:** Maintainability
- **Location:** `internal/app/handlers/keyboard.go` (lines 52-86)
- **Description:** The "y" key handler (lines 58-66) and "enter" key handler (lines 72-84) contain identical code to dismiss the kill dialog and execute the kill. The cleanup logic (setting `ShowKillDialog`, `KillTargetPid`, `KillTargetName`, `KillDialogSel` to zero values) is repeated 4 times across the kill dialog block.
- **Impact:** If cleanup logic changes (e.g., adding a new field to reset), it must be updated in 4 places.
- **Recommendation:** Extract a `dismissKillDialog(m *data.AppState)` helper and a `confirmKill(m *data.AppState) tea.Cmd` helper.

## Security Concerns

### C-006: SSH Remote Execution Without Credential Validation

- **Severity:** HIGH
- **Category:** Security
- **Location:** `internal/provider/remote/ssh.go` (lines 22-48, `buildSSHCmd`)
- **Description:** The remote monitoring feature executes arbitrary shell commands on remote hosts via SSH. While the commands are currently hard-coded in the source (line 44), the `host.Host` value comes directly from user config (`config.RemoteHostConfig.Host`) and is passed to `exec.Command("ssh", ..., host.Host, script)`. The `KeyPath` is also user-provided and passed as `-i` argument without validation. There is no sanitization or validation of the host string -- a malicious config could inject additional SSH arguments.
- **Impact:** A crafted `Host` value starting with `-` could be interpreted as an SSH flag rather than a hostname. A crafted `KeyPath` could point to unexpected files. The config file is stored as plain JSON with no integrity checks.
- **Recommendation:** Validate that `host.Host` does not start with `-`. Validate that `host.KeyPath` (when provided) references an existing file with appropriate permissions (0600). Consider using `--` to separate SSH options from the destination: `exec.Command("ssh", args..., "--", host.Host, script)`.

### C-007: Config File Written with World-Readable Permissions

- **Severity:** MEDIUM
- **Category:** Security
- **Location:** `internal/config/config.go` (line 198, `SaveConfig`)
- **Description:** `SaveConfig` uses `os.Create(path)` which creates files with the process umask (typically 0644). The config file contains `RemoteHostConfig` entries with SSH key paths, hostnames, and potentially sensitive connection details. The parent directory is created with `0755` (line 114).
- **Impact:** Other users on the system can read the configuration, learning about remote hosts and SSH key paths.
- **Recommendation:** Create the config file with `0600` permissions using `os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)`. Set the config directory to `0700`.

### C-008: Snapshot Export Files with Broad Permissions

- **Severity:** LOW
- **Category:** Security
- **Location:** `internal/app/export.go` (lines 36, 51)
- **Description:** Snapshot JSON and CSV files are created with `0644` permissions. While the data (CPU, memory, disk usage) is not highly sensitive, the home directory file path is predictable (`~/bubble_snapshot.json`, `~/bubble_snapshot.csv`).
- **Impact:** Any user on the system can read system metrics snapshots.
- **Recommendation:** Use `0600` for the output files, or write to the XDG config/data directory instead of home.

### C-009: PowerShell Command Injection Vector in Priority Change

- **Severity:** MEDIUM
- **Category:** Security
- **Location:** `internal/provider/process/control.go` (lines 66-68)
- **Description:** The Windows priority change constructs a PowerShell command string using `fmt.Sprintf` with `pid` (int32) and `newClass` (string from a fixed list). While `pid` is a number and `newClass` comes from a hardcoded slice, the pattern of building PowerShell commands via string formatting is a potential injection risk if the code is later modified to accept user input.
- **Impact:** Currently low risk because both values are controlled. Future modifications accepting user-provided class names would be vulnerable.
- **Recommendation:** Add a comment noting the security assumption. Validate `newClass` against the `priorities` slice before use.

## Performance Concerns

### C-010: cpu.prof Profiling Artifact Committed to Git

- **Severity:** MEDIUM
- **Category:** Performance / Process
- **Location:** `/home/samouly/Projects/Golang/bubbleMonitor/cpu.prof` (in git staging area)
- **Description:** A CPU profiling artifact (`cpu.prof`, gzip compressed, ~4.5KB) is staged for commit. The `.gitignore` file does not exclude `*.prof` files. Additionally, `cmd/bub/main.go` imports `_ "net/http/pprof"` (line 8), registering pprof HTTP handlers unconditionally, though the server startup is commented out (lines 15-18).
- **Impact:** Profiling data in git bloats the repository. The pprof import adds the pprof HTTP handler to the default serve mux, which could be accidentally activated. Binary size increase from the import.
- **Recommendation:** Add `*.prof` to `.gitignore`. Unstage `cpu.prof`. Guard the pprof import with a build tag (e.g., `//go:build pprof`) or remove it entirely.

### C-011: Unbounded Process History Map Growth

- **Severity:** MEDIUM
- **Category:** Performance
- **Location:** `internal/data/state.go` (lines 108-115, `GetOrCreateHistory`); `internal/app/analysis.go` (lines 50-83, `UpdateProcessHistory`)
- **Description:** `Metrics.ProcessHistory` is a `map[int32]*RingBuffer` that grows as processes are tracked. While `PruneDeadProcessHistory` is called during each tick, it only removes entries for PIDs not in the current process list. On a system with high process churn (e.g., containers, CI runners), hundreds of short-lived processes could create RingBuffers before they are cleaned up. Each RingBuffer allocates `history_length` float64 values (default 900 = 7.2KB per buffer).
- **Impact:** Memory usage can spike on systems with high process churn. With default settings, 1000 dead-but-not-yet-pruned processes would consume ~7.2MB. The pruning only happens when `ProcessesMsg` is received, not on a timer.
- **Recommendation:** Add a cap on the total number of history entries (e.g., 100). When the cap is exceeded, evict the oldest entries. Consider reducing the default history length for per-process tracking.

### C-012: CPUPercent Requires Persistent Process Objects

- **Severity:** MEDIUM
- **Category:** Performance
- **Location:** `internal/provider/process/list.go` (lines 116-118)
- **Description:** `cached.Proc.CPUPercent()` relies on gopsutil's internal state within the `*process.Process` object. This requires keeping references to these objects across poll cycles (via `processCache`). Each `*process.Process` holds file descriptors and internal state. The `isFetching` atomic flag prevents concurrent execution, but if a single fetch is slow (e.g., a zombie process), all subsequent fetches are blocked.
- **Impact:** A single slow/hanging process query blocks all process list updates. Memory is held for cached `*process.Process` objects for all running processes.
- **Recommendation:** Add a timeout per-process query. Consider a maximum time budget for the entire `ProcessesCmd` and return partial results if exceeded.

### C-013: Net Percent Calculation Uses Hard-Coded 10 MB/s Reference

- **Severity:** LOW
- **Category:** Performance / Correctness
- **Location:** `internal/app/update.go` (lines 307-309)
- **Description:** The network history percentage is calculated as `(totalNetRate / 10) * 100`, where 10 is a hard-coded reference of 10 MB/s as "100%". This means on a 10 Gbps connection the graph will always show near-max values, while on a slow connection it will barely register.
- **Description (code):**
  ```go
  netPercent := (totalNetRate / 10) * 100
  if netPercent > 100 { netPercent = 100 }
  ```
- **Impact:** Network history graph is not meaningful on networks significantly faster or slower than 10 MB/s.
- **Recommendation:** Make the network reference speed configurable, or auto-scale based on observed peak values (using the RingBuffer's `Max()` method).

## Reliability Concerns

### C-014: SyncProcessesMap Missing Mutex Protection

- **Severity:** HIGH
- **Category:** Reliability
- **Location:** `internal/data/state.go` (lines 67-72, `SyncProcessesMap`)
- **Description:** `SyncProcessesMap` rebuilds `ProcessesByPid` from the `Processes` slice without acquiring the `stateMu` mutex. It is called from `update.go` line 341 inside the `ProcessesMsg` handler. Other methods like `GetProcessByPid` (line 60) do acquire `stateMu.RLock()`. While Bubble Tea's single-threaded Update loop means this specific call is safe in practice, the inconsistency is dangerous if the code is refactored to use concurrent access or if `SyncProcessesMap` is called from elsewhere.
- **Impact:** Data race if `SyncProcessesMap` is ever called from a goroutine. Inconsistent API contract where some methods require mutex and others do not.
- **Recommendation:** Either add mutex protection to `SyncProcessesMap` for API consistency, or document clearly that it must only be called from the Bubble Tea update loop.

### C-015: AppState Mutex Not Used in Most Update Handler Access

- **Severity:** MEDIUM
- **Category:** Reliability
- **Location:** `internal/app/update.go` (entire file); `internal/data/state.go`
- **Description:** The `AppState` struct has a `sync.RWMutex` (`stateMu`) that is used by `SetSuspended`, `IsSuspended`, `ToggleCollapsed`, `IsCollapsed`, `ToggleBookmark`, `IsBookmarked`, `GetProcessByPid`, `PruneDeadProcessMaps`, and `PruneDeadProcessHistory`. However, the vast majority of state mutations in `update.go` (lines 272-441) access fields directly without any mutex: `m.Metrics.Cpu`, `m.Process.Processes`, `m.UI.SelectedTab`, etc. The mutex is only used for the map-type fields.
- **Impact:** If any provider command ever accesses state concurrently (e.g., remote SSH checks running in background goroutines), there will be data races. Bubble Tea's model guarantees sequential message processing, so this is safe today but fragile.
- **Recommendation:** Document the concurrency contract clearly: all mutations happen in the Bubble Tea Update loop, the mutex is only needed for methods that could be called from goroutines. Consider renaming `stateMu` to `goroutineMu` to clarify intent.

### C-016: Remote SSH Commands Have No Cancellation/Timeout Enforcement

- **Severity:** MEDIUM
- **Category:** Reliability
- **Location:** `internal/provider/remote/ssh.go` (lines 50-66, `CheckRemoteCmd`)
- **Description:** While `ConnectTimeout` is set on the SSH command (line 27), the `cmd.Output()` call (line 55) has no context or deadline. If the remote host accepts the connection but hangs during command execution (e.g., a NFS-mounted `/proc/meminfo`), the goroutine will block indefinitely. Multiple configured remote hosts are checked every other tick (update.go line 263), so a single hanging host could accumulate blocked goroutines.
- **Impact:** Goroutine leak when remote hosts are unresponsive. Over time, this consumes memory and goroutine IDs.
- **Recommendation:** Use `exec.CommandContext` with a context deadline that encompasses both connection and execution time. The configured timeout should apply to the entire operation:
  ```go
  ctx, cancel := context.WithTimeout(context.Background(), timeout)
  defer cancel()
  cmd := exec.CommandContext(ctx, "ssh", args...)
  ```

### C-017: Disk IO Rate Can Show Negative After Counter Wrap

- **Severity:** LOW
- **Category:** Reliability
- **Location:** `internal/app/update.go` (lines 370-394)
- **Description:** The code checks `if totalRead >= lastTotalRead` (line 383) before calculating the read rate, but `DiskLastIO` is set to `msg` (line 394) which means `LastDiskIO` is updated to the current snapshot. If counters wrap or are reset (e.g., module reload), the rate drops to 0 rather than handling the discontinuity. Additionally, line 393 sets `m.Metrics.DiskIO = msg` and then line 394 immediately sets `m.Metrics.LastDiskIO = msg`, making `DiskIO` and `LastDiskIO` point to the same map, which is correct for this tick but means the "last" values are always identical to current.
- **Impact:** After a counter reset, disk IO rate shows 0 for one tick then resumes correctly. The shared reference is not a bug because maps are replaced, not mutated, but it is confusing.
- **Recommendation:** Add a comment explaining the counter-wrap behavior. No code change strictly needed but the pattern is fragile.

### C-018: Config File Race Between Watcher and Save

- **Severity:** MEDIUM
- **Category:** Reliability
- **Location:** `internal/config/watcher.go` (lines 16-32); `internal/config/config.go` (line 192, `SaveConfig`)
- **Description:** The config watcher polls the file's modification time every 2 seconds. Meanwhile, settings changes in the TUI save the config file immediately (keyboard.go lines 166, 223, 232, 254). This creates a race: the user changes settings in the TUI, which saves the file, which the watcher then detects as an "external change" and reloads. The `reflect.DeepEqual` check (update.go line 140) mitigates false reloads, but if the user is actively editing settings, the watcher could reload a partially-written JSON file.
- **Impact:** Config could be corrupted or reverted during rapid settings changes. The JSON encoder writes incrementally, so a read mid-write could see truncated JSON.
- **Recommendation:** Write config atomically using a temp file + rename pattern:
  ```go
  tmp := path + ".tmp"
  os.WriteFile(tmp, data, 0600)
  os.Rename(tmp, path)
  ```
  Alternatively, set a "self-write" flag that the watcher checks to skip its own saves.

## Platform-Specific Concerns

### C-019: CGO_ENABLED=0 Conflicts with NVML/AMD GPU Detection

- **Severity:** HIGH
- **Category:** Platform
- **Location:** `.goreleaser.yml` (line 14: `CGO_ENABLED=0`); `internal/provider/system/hardware_nvidia.go` (build tag: `//go:build linux && cgo`)
- **Description:** The GoReleaser config sets `CGO_ENABLED=0` globally for all builds, then overrides only `linux/amd64` to `CGO_ENABLED=1`. This means: (1) Linux ARM64 builds have no NVML/AMD GPU support. (2) The `amdsmi.go` and `go-nvml` libraries require CGO, so Linux ARM64 will fail to link if those packages are imported. (3) The stub files (`hardware_nvidia_stub.go`, `hardware_amd_stub.go`) correctly handle the `!cgo` case, so the build succeeds, but GPU detection silently returns nothing on ARM64 Linux.
- **Impact:** Users on ARM64 Linux (e.g., ARM servers, Raspberry Pi) get no GPU monitoring even with NVIDIA/AMD GPUs present. macOS and Windows users never get NVML support.
- **Recommendation:** Document platform limitations clearly. If ARM64 GPU support is needed, add `CGO_ENABLED=1` for `linux/arm64` in GoReleaser overrides.

### C-020: Multiple exec.Command Calls with OS-Specific Paths

- **Severity:** MEDIUM
- **Category:** Platform
- **Location:**
  - `internal/provider/system/services.go` (line 19: `systemctl`, Linux-only)
  - `internal/provider/system/logs.go` (line 19: `journalctl`, Linux-only)
  - `internal/provider/system/vm.go` (line 160: `dmidecode`, Linux-only; line 195: `powershell`, Windows-only)
  - `internal/provider/system/hardware.go` (line 230: `nvidia-smi`, Linux/Windows)
  - `internal/provider/system/container.go` (lines 43, 59: `docker`, `kubectl`)
  - `internal/provider/system/disk.go` (line 50: `lsblk`, Linux-only)
- **Description:** Multiple external commands are invoked without verifying the OS beforehand or using build tags. Some have runtime checks (`runtime.GOOS == "linux"`), others do not. `disk.go` calls `lsblk` on all platforms but it only exists on Linux. `container.go` calls `docker info` on Windows but uses Unix socket on Linux.
- **Impact:** On non-Linux platforms, `lsblk` invocation fails silently (error is ignored, line 52: `if cmdErr == nil`). On macOS, `dmidecode` is not available. The `fetchDarwinGpus` and `fetchWindowsGpus` functions handle their platforms, but `fetchLinuxGpus` is called on all platforms via `runtime.GOOS` checks that may not cover all edge cases.
- **Recommendation:** Add explicit `runtime.GOOS` checks before every `exec.Command` call. For `disk.go`, guard the `lsblk` call with `if runtime.GOOS == "linux"`.

### C-021: No macOS/Windows Remote Monitoring Support

- **Severity:** LOW
- **Category:** Platform
- **Location:** `internal/provider/remote/ssh.go` (lines 44)
- **Description:** The remote monitoring script (line 44) runs Linux-specific commands: `/proc/loadavg`, `/proc/meminfo`, `ps aux`, `df -B1`. If the remote host is macOS or Windows, the output parsing will fail silently and most metrics will be zero or empty.
- **Impact:** Remote monitoring is effectively Linux-only for monitored hosts, even though the TUI itself supports macOS and Windows.
- **Recommendation:** Document that remote hosts must be Linux. Optionally, add OS detection and use appropriate commands for macOS (`sysctl`, `vm_stat`, `ps aux`) and Windows (`wmic`, `tasklist`).

## Code Quality Concerns

### C-022: UIState.ZoneManager Stored as interface{} to Avoid Circular Import

- **Severity:** MEDIUM
- **Category:** Maintainability
- **Location:** `internal/data/groups.go` (line 203: `ZoneManager interface{}`)
- **Description:** The `UIState.ZoneManager` field is typed as `interface{}` with a comment saying `*input.ZoneManager - stored as interface to avoid circular import`. This means any access requires a type assertion (seen in `internal/ui/layout.go` line 53: `if zm, ok := s.UI.ZoneManager.(input.ZoneManager); ok`). If the type assertion fails, the zone manager is silently replaced with a new one.
- **Impact:** Loss of type safety. Runtime panics if the assertion is done incorrectly elsewhere. The circular import indicates the dependency graph needs restructuring.
- **Recommendation:** Break the circular dependency by defining a `ZoneManagerInterface` in the `data` package (or a shared interface package) that `input.ZoneManager` satisfies, then use that interface instead of `interface{}`.

### C-023: Inline String Interner with O(n) Eviction

- **Severity:** LOW
- **Category:** Performance
- **Location:** `internal/provider/process/utils.go` (lines 78-111, `evictOldest`)
- **Description:** When the string interner cache exceeds 5000 entries, it performs an O(n) partial sort to find the 500 oldest entries. With `n=5000`, this does ~2.5M comparisons per eviction. The interner is global and shared across all process fetch operations.
- **Impact:** Minor performance spike when the cache first fills. After that, eviction happens rarely (only when many new process names appear). On systems with stable process sets, this is unlikely to trigger.
- **Recommendation:** Consider using a container/heap for O(log n) eviction. Or simply clear the entire cache when full since process names are short strings and the allocation cost is minimal.

### C-024: Direct process.Package Usage in Handlers Bypasses Interface

- **Severity:** LOW
- **Category:** Maintainability
- **Location:** `internal/app/handlers/keyboard.go` (lines 15, 280, 309, 551, 552, 559, 573, 583, 591)
- **Description:** `keyboard.go` imports the `process` package directly and calls `process.ProcessesCmd`, `process.ReniceProcessCmdSafe`, `process.FetchOpenFilesCmd`, `process.SuspendProcessCmd`, `process.ResumeProcessCmd` directly, bypassing the `ProcessProvider` interface defined in `internal/provider/interfaces.go`. This undermines the adapter pattern used by `Model.Provider.Process`.
- **Impact:** The test infrastructure (`testutil/mock_process.go`) provides mock implementations of the `ProcessProvider` interface, but the keyboard handler bypasses it. Tests that inject a mock provider cannot intercept these calls.
- **Recommendation:** Pass the `ProcessProvider` interface to `HandleKey` (or the specific command functions it needs) rather than calling the process package directly.

### C-025: No Graceful Shutdown for Background Operations

- **Severity:** MEDIUM
- **Category:** Reliability
- **Location:** `cmd/bub/main.go` (lines 14-25)
- **Description:** The main function creates a `tea.Program` and runs it. When the user quits, `tea.Quit` is returned from the Update handler (update.go line 131), but there is no explicit cleanup of: (1) running SSH connections to remote hosts, (2) the `processCache` map holding `*process.Process` objects, (3) the string interner's global state, (4) open log files if logging is enabled. The `processCache` is a package-level variable that persists across test runs.
- **Impact:** Minor resource leak on exit (OS reclaims everything). More problematic for tests: the global `processCache`, `globalInterner`, and `containerChecked` state persist between test runs, potentially causing test pollution.
- **Recommendation:** Add a `Shutdown()` method to relevant providers that clears caches. For tests, add `ResetForTesting()` functions or use `t.Cleanup()` to reset global state.

### C-026: Test File Staged with Potentially Stale Mock Code

- **Severity:** LOW
- **Category:** Maintainability
- **Location:** `internal/app/testutil/mock_system.go`, `internal/app/testutil/mock_process.go`, `internal/app/testutil/mock_remote.go`
- **Description:** The testutil directory contains mock implementations of the provider interfaces. These mocks must be manually kept in sync with the interfaces in `internal/provider/interfaces.go`. If an interface method is added or changed, the mock may not compile or may silently implement the wrong behavior.
- **Recommendation:** Consider using a mock generation tool (e.g., `mockgen` or `gomock`) to auto-generate mocks from interfaces. Add a build tag or CI check that verifies mocks match interfaces.

### C-027: Health Threshold Constants Inconsistently Defined

- **Severity:** LOW
- **Category:** Maintainability
- **Location:** `internal/data/constants.go` (lines 26-46); `internal/app/analysis.go` (lines 10-13)
- **Description:** Health scoring thresholds and deductions are defined in two places. `constants.go` defines `HealthThresholdHealthy=90`, `HealthThresholdWarning=70`, `HealthDeductionCPUCritical=30`, etc. But `analysis.go` uses hardcoded multipliers like `0.7` (line 18: `thresholds[configpkg.MetricCPU] * 0.7`) and `0.8` (line 38: `thresholds[configpkg.MetricTemp] * 0.8`) instead of the defined constants. The constants in `constants.go` (`HealthThresholdWarning`, `HealthThresholdHealthy`) appear to be unused.
- **Impact:** Confusion about which thresholds are actually used. If someone changes `HealthThresholdWarning` expecting it to affect behavior, nothing will change. The magic numbers 0.7 and 0.8 are not documented.
- **Recommendation:** Use the defined constants in `analysis.go` or remove the unused constants. Replace `0.7` and `0.8` with named constants like `WarningThresholdFraction = 0.7`.

## Missing Critical Features

### C-028: No Input Validation on Config Values

- **Severity:** MEDIUM
- **Category:** Reliability
- **Location:** `internal/config/config.go` (lines 120-181, `LoadConfigFromPath`)
- **Description:** When loading config, the code checks for empty/zero values and replaces them with defaults. However, there is no validation of out-of-range values: a user could set `RefreshRate: 10` (10ms, extremely fast), `HistoryLength: 1` (useless), `Thresholds: {CPU: -5}` (negative threshold), or `Port: 99999` (invalid port for SSH). These would be accepted silently.
- **Impact:** Invalid config values cause poor behavior (excessive CPU usage from fast polling, broken charts from tiny history, broken alerts from negative thresholds) with no error message.
- **Recommendation:** Add a `Validate()` method to `AppConfig` that checks all values are within acceptable ranges and returns errors for invalid configurations. Call it after loading.

### C-029: No Unit Tests for Provider Package

- **Severity:** MEDIUM
- **Category:** Maintainability / Quality
- **Location:** `internal/provider/system/*.go`, `internal/provider/remote/ssh.go`
- **Description:** The `internal/provider/system/` directory contains 15+ Go files with no test files. The remote SSH provider (`internal/provider/remote/ssh.go`) also has no tests. The only provider tests are in `internal/provider/process/utils_test.go`. The parsing functions (`parseRemoteOutput`, `parseMeminfo`, `parseDisk`, `parseProcesses`, `parseLsblkUnmounted`, etc.) are good candidates for unit testing since they operate on string input.
- **Impact:** Parsing regressions in remote output, lsblk output, or sysfs data are not caught by CI. The CI pipeline (`ci.yml`) runs `go test -v -race ./...` but most provider code has no test coverage.
- **Recommendation:** Add table-driven tests for all parsing functions. Start with `parseRemoteOutput` and `parseLsblkUnmounted` since they parse external command output.

---

*Concerns audit: 2026-04-09*

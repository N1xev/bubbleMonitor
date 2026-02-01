---
name: Complete Feature Enhancement Plan
overview: Comprehensive plan to add all 16 feature categories to Bubble Monitor, organized into logical implementation phases covering monitoring enhancements, process management, visualization, customization, data export, and advanced features.
todos: []
---

# Complete Feature Enhancement Plan for Bubble Monitor

This plan implements all 16 feature categories across multiple phases, enhancing the TUI system monitor with advanced monitoring, process management, visualization, and customization capabilities.

## Architecture Overview

The implementation will extend the existing Bubble Tea architecture:

```
Model (State)
├── Extended Metrics (temperature, disk I/O, network interfaces, battery)
├── Configurable Settings (history length, refresh rate, themes, thresholds)
├── Alert System (threshold monitoring, notifications)
├── Process Management (priority, suspend/resume, tree view, open files)
├── Export/Logging (CSV, JSON, file logging)
└── Remote Monitoring (SSH connections, multi-host)

Commands (Data Fetching)
├── TemperatureCmd (CPU, GPU, disk)
├── DiskIOCmd (read/write rates per partition)
├── NetworkInterfacesCmd (per-interface stats)
├── BatteryCmd (laptop battery status)
├── ServicesCmd (system services/daemons)
├── ConnectionsCmd (network connections)
└── RemoteMetricsCmd (SSH-based remote monitoring)

Views (Rendering)
├── Enhanced Metrics Tab (extended history, new metrics)
├── Alerts Tab (threshold monitoring)
├── Services Tab (system services)
├── Network Tab (connections, interfaces)
└── Settings Tab (configuration UI)
```

## Phase 1: Extended History & Additional Metrics

### 1.1 Configurable History Length

**Files:** `src/app.go`, `src/update.go`, `src/render.go`

- Add `HistoryLength` field to Model (default 60, options: 60, 300, 900, 3600 seconds)
- Add `HistoryLengthSelector` state for UI
- Modify history append logic to respect `HistoryLength` instead of hardcoded `MaxHistoryLen`
- Add keyboard shortcut to cycle history lengths (e.g., `H` key)
- Update sparkline rendering to show appropriate time labels

### 1.2 Disk I/O Metrics

**Files:** `src/app.go`, `src/commands.go`, `src/render.go`, `src/update.go`

- Add `DiskIO` struct with `ReadRate`, `WriteRate`, `ReadIOPS`, `WriteIOPS` per partition
- Add `DiskIOHistory` map[string][]DiskIOData to track per-partition history
- Create `DiskIOCmd()` using `gopsutil/v3/disk` to fetch I/O counters
- Add disk I/O history tracking in `Update()` for `MetricsMsg`
- Render disk I/O charts in Metrics tab alongside existing charts
- Show per-partition I/O rates in Disks tab

### 1.3 Temperature Monitoring

**Files:** `src/app.go`, `src/commands.go`, `src/render.go`

- Add `TemperatureInfo` struct: `CPUTemp`, `GPUTemp`, `DiskTemps` map
- Create platform-specific temperature commands:
  - Linux: Read from `/sys/class/thermal/thermal_zone*/temp`
  - Windows: Use WMI queries or third-party libraries
  - macOS: Use `powermetrics` or `sysctl`
- Add temperature display in Overview and System tabs
- Show temperature warnings when thresholds exceeded
- Add temperature history tracking

### 1.4 Network Interface Details

**Files:** `src/app.go`, `src/commands.go`, `src/render.go`

- Add `NetworkInterface` struct: `Name`, `BytesSent`, `BytesRecv`, `PacketsSent`, `PacketsRecv`, `Errors`, `Dropped`
- Add `NetworkInterfaces` slice to Model
- Create `NetworkInterfacesCmd()` using `gopsutil/v3/net`
- Add new "NETWORK" tab showing per-interface statistics
- Display interface selection and detailed stats per interface
- Show interface history charts

### 1.5 Battery Status

**Files:** `src/app.go`, `src/commands.go`, `src/render.go`

- Add `BatteryInfo` struct: `Percent`, `Status`, `TimeRemaining`, `PowerConsumption`
- Add `BatteryInfo` field to Model (nullable)
- Create `BatteryCmd()` using `gopsutil/v3/battery` (if available) or platform-specific methods
- Display battery info in Overview tab (when available)
- Show battery status in System tab

## Phase 2: Alerts & Thresholds

### 2.1 Alert System

**Files:** `src/app.go`, `src/update.go`, `src/render.go`, `src/alerts.go` (new)

- Create `src/alerts.go` with alert management
- Add `Alert` struct: `Type`, `Metric`, `Threshold`, `CurrentValue`, `Timestamp`, `Acknowledged`
- Add `Alerts` slice and `AlertThresholds` map to Model
- Add `AlertSettings` struct with configurable thresholds per metric
- Implement threshold checking in `Update()` when processing `MetricsMsg`
- Create alert rendering overlay (similar to help dialog)
- Add alert acknowledgment (press `a` to acknowledge)
- Visual indicators in tabs when alerts are active
- Optional audio alerts (beep on threshold breach)

### 2.2 Threshold Configuration

**Files:** `src/app.go`, `src/render.go`, `src/update.go`

- Add Settings tab (tab 6) for configuration
- Allow setting thresholds for: CPU, Memory, Disk, Temperature, Network
- Store thresholds in Model state
- Persist thresholds to config file (JSON) in user home directory
- Load thresholds on startup

## Phase 3: Enhanced Process Management

### 3.1 Process Priority Management

**Files:** `src/app.go`, `src/update.go`, `src/commands.go`, `src/render.go`

- Add `ChangePriorityCmd(pid, priority)` using `gopsutil/v3/process`
- Add priority change dialog (similar to kill dialog)
- Show current priority in process details panel
- Add keyboard shortcuts: `+` to increase priority, `-` to decrease
- Display priority in process list (nice value)

### 3.2 Suspend/Resume Processes

**Files:** `src/app.go`, `src/update.go`, `src/commands.go`

- Add `SuspendProcessCmd()` and `ResumeProcessCmd()` using process.Suspend()/Resume()
- Add suspend/resume actions in process details
- Keyboard shortcuts: `z` to suspend, `Z` to resume
- Visual indicator for suspended processes in list

### 3.3 Process Tree View

**Files:** `src/app.go`, `src/commands.go`, `src/render.go`, `src/process_tree.go` (new)

- Create `src/process_tree.go` for tree building logic
- Add `ProcessTreeNode` struct with parent/children relationships
- Build process tree from parent PIDs
- Add tree view toggle in Processes tab (`t` key)
- Render hierarchical process tree with indentation
- Allow navigation and selection in tree view

### 3.4 Open Files (lsof)

**Files:** `src/app.go`, `src/commands.go`, `src/render.go`

- Add `OpenFilesCmd(pid)` to fetch open files using `process.OpenFiles()`
- Add `OpenFiles` field to `ProcessInfo`
- Display open files in process details panel (scrollable)
- Show file count and list of open files
- Filter by file type (regular files, sockets, pipes)

### 3.5 Process Grouping & Search

**Files:** `src/app.go`, `src/update.go`, `src/render.go`

- Add `GroupBy` field to Model: "none", "user", "parent", "command"
- Implement grouping logic in `getFilteredProcesses()`
- Add group headers in process list when grouped
- Enhanced search: search by PID, name, command line, user
- Add search mode toggle (`/` for search vs `f` for filter)
- Bookmark processes: add `BookmarkedProcesses` map[int32]bool
- Keyboard shortcut `b` to bookmark/unbookmark selected process
- Show bookmarked processes indicator

## Phase 4: Enhanced Visualization

### 4.1 Advanced Charts

**Files:** `src/render.go`, `src/helpers.go`, `src/charts.go` (new)

- Create `src/charts.go` for advanced chart rendering
- Implement line chart renderer (more detailed than sparklines)
- Add chart type toggle: sparkline vs line chart
- Per-core CPU history charts (one chart per core or combined)
- Disk I/O charts per partition
- Network charts per interface
- Chart zoom functionality (time range selection)

### 4.2 Per-Process Resource History

**Files:** `src/app.go`, `src/update.go`, `src/render.go`

- Add `ProcessHistory` map[int32][]ProcessSnapshot
- Track CPU and memory usage over time for selected process
- Display process history chart in process details panel
- Limit history to top N processes to manage memory

## Phase 5: Customization

### 5.1 Theme System

**Files:** `src/colors.go`, `src/app.go`, `src/render.go`, `src/themes.go` (new)

- Create `src/themes.go` with multiple theme definitions
- Add `Theme` struct and `CurrentTheme` field to Model
- Implement themes: "default", "dark", "light", "high-contrast", "monochrome"
- Theme switching via Settings tab or keyboard shortcut (`T` key)
- Persist theme preference to config file

### 5.2 Configurable Refresh Rate

**Files:** `src/app.go`, `src/update.go`, `src/commands.go`

- Add `RefreshInterval` field to Model (default 1 second)
- Make `TickCmd()` respect `RefreshInterval`
- Allow setting refresh rate in Settings tab (0.5s, 1s, 2s, 5s)
- Update tick interval dynamically


## Phase 6: Data Export & Logging

### 6.1 Export to CSV/JSON

**Files:** `src/app.go`, `src/update.go`, `src/export.go` (new)

- Create `src/export.go` for export functionality should export all tabs informations in that!
- Add export dialog (similar to kill dialog)
- Implement `ExportMetricsCSV()` and `ExportMetricsJSON()`
- Export current metrics snapshot or historical data
- Keyboard shortcut `e` to open export dialog with CSV and JSON buttons they should be accesed by arrow keys
- Allow selecting which metrics to export or export all in the dialog with tab to select/unselect

### 6.2 Metrics Logging

**Files:** `src/app.go`, `src/update.go`, `src/logging.go` (new)

- Create `src/logging.go` for file logging
- Add `LoggingEnabled` and `LogFilePath` to Model
- Implement background logging to file (append mode)
- Log metrics at configured interval
- Rotate log files when size limit reached
- Toggle logging via Settings tab

### 6.3 Snapshots

**Files:** `src/app.go`, `src/update.go`, `src/export.go`

- Add `TakeSnapshot()` function to capture current state
- Save snapshot to JSON file with timestamp
- Keyboard shortcut `S` (capital) to take snapshot
- Display snapshot count in footer

## Phase 7: Remote Monitoring

### 7.1 SSH Remote Monitoring

**Files:** `src/app.go`, `src/commands.go`, `src/remote.go` (new), `src/render.go`

- Create `src/remote.go` for SSH connection management
- Add `RemoteHost` struct: `Name`, `Host`, `User`, `Port`, `KeyPath`
- Add `RemoteHosts` slice and `ActiveRemoteHost` to Model
- Implement SSH connection using `golang.org/x/crypto/ssh`
- Create `RemoteMetricsCmd(host)` to fetch metrics via SSH
- Add Remote tab for managing connections
- Display remote host metrics in existing tabs (with host indicator)
- Support multiple remote hosts (switch with `Ctrl+H`)

### 7.2 Multi-Host Comparison

**Files:** `src/app.go`, `src/render.go`

- Add comparison view showing metrics from multiple hosts side-by-side
- Compare tab showing CPU, Memory, etc. across hosts
- Color-code hosts for easy identification

## Phase 8: Additional System Information

### 8.1 System Services

**Files:** `src/app.go`, `src/commands.go`, `src/render.go`

- Add `ServiceInfo` struct: `Name`, `Status`, `PID`, `Description`
- Create `ServicesCmd()` using platform-specific methods:
  - Linux: `systemctl list-units`
  - Windows: `Get-Service` (WMI)
  - macOS: `launchctl list`
- Add Services tab showing running services
- Allow filtering and searching services
- Show service status (running, stopped, failed)

### 8.2 Network Connections

**Files:** `src/app.go`, `src/commands.go`, `src/render.go`

- Add `ConnectionInfo` struct: `LocalAddr`, `RemoteAddr`, `State`, `PID`, `ProcessName`
- Create `ConnectionsCmd()` using `gopsutil/v3/net`
- Add Connections view in Network tab
- Filter by state, process, address
- Show connection count per process

### 8.3 System Logs Viewer

**Files:** `src/app.go`, `src/commands.go`, `src/render.go`, `src/logs.go` (new)

- Create `src/logs.go` for log viewing
- Add `LogEntry` struct: `Timestamp`, `Level`, `Message`, `Source`
- Platform-specific log reading:
  - Linux: `/var/log/syslog`, `journalctl`
  - Windows: Event Log
  - macOS: `log show`
- Add Logs tab with scrollable log viewer
- Filter by log level, search logs
- Auto-scroll to latest logs

## Phase 9: GPU Enhancements

### 9.1 AMD GPU Support

**Files:** `src/commands.go`, `src/app.go`

- Extend `GpuInfoCmd()` to detect AMD GPUs
- Use `rocm-smi` command for AMD GPU info
- Parse AMD GPU metrics (similar to NVIDIA)

### 9.2 Intel GPU Support

**Files:** `src/commands.go`, `src/app.go`

- Add Intel GPU detection using `intel_gpu_top` or similar
- Parse Intel GPU metrics
- Display in System tab alongside NVIDIA/AMD

### 9.3 Enhanced GPU Metrics

**Files:** `src/app.go`, `src/commands.go`, `src/render.go`

- Add GPU utilization percentage, temperature, power consumption
- GPU history tracking
- Per-GPU charts in Metrics tab

## Phase 10: User Experience Enhancements

### 10.1 Enhanced Keyboard Shortcuts

**Files:** `src/update.go`, `src/view.go`

- Add comprehensive keyboard shortcut system
- Document all shortcuts in help overlay
- Add `Ctrl+K` for quick kill (skip dialog with confirmation flag)
- Custom keybindings support (load from config)

### 10.2 Mouse Support

**Files:** `src/update.go`, `src/view.go`

- Add mouse event handling in `Update()`
- Support clicking tabs to switch
- Click processes to select
- Scroll with mouse wheel
- Click and drag for column resizing

### 10.3 Layout Improvements

**Files:** `src/render.go`, `src/app.go`

- Add split view mode (multiple metrics side-by-side)
- Resizable panels (when in split mode)
- Full-screen mode for specific tabs (`F11` or `f` key)
- Panel layout persistence

## Phase 11: Performance Optimizations

### 11.1 Configurable Update Intervals

**Files:** `src/app.go`, `src/update.go`, `src/commands.go`

- Different update intervals for different metrics (fast for CPU, slower for disk)
- Lazy loading: only update visible tabs
- Background updates with debouncing

### 11.2 Process List Optimization

**Files:** `src/app.go`, `src/update.go`, `src/render.go`

- Limit process list to top N by default (configurable)
- Virtual scrolling for large process lists
- Process list caching with smart invalidation

## Phase 12: Resource Analysis

### 12.1 Top Resource Consumers

**Files:** `src/app.go`, `src/render.go`, `src/analysis.go` (new)

- Create `src/analysis.go` for analysis functions
- Identify top CPU, memory, disk I/O consumers
- Display "Top Consumers" panel in Overview tab
- Historical top consumers tracking

### 12.2 Resource Usage Trends

**Files:** `src/app.go`, `src/render.go`, `src/analysis.go`

- Calculate trends (increasing/decreasing) for metrics
- Show trend indicators (arrows, colors)
- Predict future usage based on trends

### 12.3 Process Comparison

**Files:** `src/app.go`, `src/render.go`

- Select multiple processes for comparison
- Side-by-side comparison view
- Compare CPU, memory, and other metrics

## Phase 13: System Health

### 13.1 Health Score

**Files:** `src/app.go`, `src/update.go`, `src/render.go`, `src/health.go` (new)

- Create `src/health.go` for health calculation
- Calculate composite health score (0-100) based on:
  - CPU usage (weighted)
  - Memory usage
  - Disk usage
  - Temperature
  - Process errors
- Display health score in Overview tab
- Color-code health status

### 13.2 Optimization Recommendations

**Files:** `src/app.go`, `src/render.go`, `src/health.go`

- Analyze system state and provide recommendations
- Examples: "High CPU usage detected, consider closing applications"
- Display recommendations in System tab or dedicated panel

### 13.3 Historical Health Trends

**Files:** `src/app.go`, `src/render.go`

- Track health score over time
- Display health history chart
- Alert when health score drops significantly

## Implementation Order

1. **Phase 1** (Extended History & Metrics) - Foundation for other features
2. **Phase 2** (Alerts) - Important for monitoring
3. **Phase 5** (Customization) - User experience improvements
4. **Phase 3** (Process Management) - Core functionality enhancement
5. **Phase 4** (Visualization) - Better data presentation
6. **Phase 6** (Export/Logging) - Data persistence
7. **Phase 8** (System Info) - Additional information
8. **Phase 9** (GPU) - Hardware support
9. **Phase 10** (UX) - Polish and usability
10. **Phase 11** (Performance) - Optimization
11. **Phase 12** (Analysis) - Advanced features
12. **Phase 13** (Health) - System insights
13. **Phase 7** (Remote) - Advanced feature (can be done in parallel)

## Configuration File Structure

Create `~/.bubblem/config.json`:

```json
{
  "theme": "default",
  "refreshInterval": 1,
  "historyLength": 60,
  "alertThresholds": {
    "cpu": 80,
    "memory": 85,
    "disk": 90,
    "temperature": 75
  },
  "columnWidths": {
    "pid": 8,
    "name": 30,
    "cpu": 8,
    "memory": 8
  },
  "remoteHosts": [],
  "logging": {
    "enabled": false,
    "path": "~/.bubblem/metrics.log"
  }
}
```

## Dependencies to Add

- `golang.org/x/crypto/ssh` - For remote monitoring
- Consider adding JSON config library if not using standard library
- Platform-specific libraries for temperature monitoring (may need CGO)

## Testing Considerations

- Test each feature on Linux, Windows, and macOS
- Handle cases where platform-specific features aren't available gracefully
- Test with large numbers of processes
- Test remote monitoring with various SSH configurations
- Verify config file persistence works correctly
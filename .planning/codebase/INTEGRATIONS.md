# External Integrations

**Analysis Date:** 2026-04-09

## OS-Level System Integrations

### Linux /proc Filesystem

- `/proc/cpuinfo` - CPU flags for hypervisor/VM detection (`internal/provider/system/vm.go` line 82)
- `/proc/loadavg` - Load averages (via SSH remote monitoring, `internal/provider/remote/ssh.go` line 44)
- `/proc/meminfo` - Memory info (via SSH remote monitoring, `internal/provider/remote/ssh.go` line 44)
- `/proc/net/dev` - Network device stats (via SSH remote monitoring, `internal/provider/remote/ssh.go` line 44)
- `/proc/1/cgroup` - Container detection (Docker/LXC) (`internal/provider/system/vm.go` line 267)

### Linux /sys Filesystem

- `/sys/class/drm/card*/device/` - GPU detection via DRM sysfs enumeration
  - `device/uevent` - PCI ID and driver info
  - `device/product_name` - GPU product name
  - `device/product_version` - GPU product version
  - `device/chip_info` - GPU chip identification
  - `device/mem_info_vram_total` - VRAM size
  - `device/meminfo` - GPU memory info
  - `device/hwmon/hwmon*/temp1_input` - Intel GPU temperature
  - `device/thermal_zone/temp` - Intel GPU temperature (fallback)
  - `device/gt/gt0/punit/gpu_freq_mhz` - Intel GPU frequency
  - `device/freq0/freq` - Intel GPU frequency (fallback)
- `/sys/class/dmi/id/product_name` - DMI product name for VM detection
- `/sys/class/dmi/id/sys_vendor` - DMI vendor for VM detection
- `/sys/class/dmi/id/board_vendor` - DMI board vendor for VM detection
- `/sys/class/dmi/id/bios_vendor` - DMI BIOS vendor for VM detection
- `/sys/devices/system/cpu/possible` - CPU count (VM/container)
- `/sys/fs/cgroup/cpuset/cpuset.cpus` - CPUset configuration (VM/container)
- `/sys/fs/cgroup/cpuset/cpuset.effective_cpus` - Effective CPUs (VM/container)

Files: `internal/provider/system/hardware.go`, `internal/provider/system/vm.go`

### File Existence Checks

- `/.dockerenv` - Docker container detection (`internal/provider/system/vm.go` line 275)
- `/run/.containerenv` - Container environment detection (`internal/provider/system/vm.go` line 278)
- `/usr/bin/lxc-checkconfig` - LXC detection (`internal/provider/system/vm.go` line 283)
- `/sys/fs/cgroup/cpuset/lxc` - LXC detection (`internal/provider/system/vm.go` line 283)
- `/var/run/docker.sock` - Docker daemon socket (`internal/provider/system/container.go` line 47)
- `/run/docker.sock` - Docker daemon socket fallback (`internal/provider/system/container.go` line 48)

### gopsutil Integration

**Package:** `github.com/shirou/gopsutil/v3` - Cross-platform system metrics

| Subsystem | Package | Functions Used | Files |
|-----------|---------|----------------|-------|
| CPU | `gopsutil/v3/cpu` | `Percent(false)`, `Percent(true)` | `internal/provider/system/metrics.go` |
| Memory | `gopsutil/v3/mem` | `VirtualMemory()`, `SwapMemory()` | `internal/provider/system/metrics.go`, `vm.go` |
| Disk | `gopsutil/v3/disk` | `Usage()`, `Partitions()`, `IOCounters()` | `internal/provider/system/metrics.go`, `disk.go` |
| Network | `gopsutil/v3/net` | `IOCounters(false)`, `IOCounters(true)`, `Connections("all")` | `internal/provider/system/metrics.go`, `network.go`, `connections.go` |
| Host | `gopsutil/v3/host` | `Info()`, `SensorsTemperatures()` | `internal/provider/system/hardware.go`, `vm.go` |
| Load | `gopsutil/v3/load` | `Avg()` | `internal/provider/system/metrics.go` |
| Process | `gopsutil/v3/process` | `Pids()`, `NewProcess()`, `Name()`, `CPUPercent()`, `MemoryPercent()`, `MemoryInfo()`, `Status()`, `Cmdline()`, `Username()`, `OpenFiles()`, `CreateTime()`, `Nice()`, `Ppid()`, `Suspend()`, `Resume()` | `internal/provider/process/list.go`, `files.go`, `control.go` |

## Hardware Detection

### NVIDIA GPU

**Detection strategy (layered):**

1. **NVML library** (CGO, Linux only) - `github.com/NVIDIA/go-nvml/pkg/nvml`
   - `nvml.Init()` -> `nvml.DeviceGetCount()` -> `nvml.DeviceGetHandleByIndex()`
   - Used for both detection and full metrics collection
   - Files: `internal/provider/system/detect_nvidia.go` (build: `linux && cgo`), `internal/provider/system/hardware_nvidia.go` (build: `linux && cgo`)
2. **nvidia-smi CLI** (all platforms, fallback) - `exec.Command("nvidia-smi", ...)`
   - Detection: `nvidia-smi --query-gpu=index --format=csv,noheader`
   - Metrics: `nvidia-smi --query-gpu=index,name,pci.bus_id,memory.total,memory.used,utilization.gpu,temperature.gpu,power.draw,clocks.sm --format=csv,noheader,nounits`
   - Files: `internal/provider/system/detect_nvidia_stub.go` (build: `!linux || !cgo`), `internal/provider/system/hardware.go` (`fetchNvidiaSmiGpus`)

**Metrics collected:** GPU name, driver, UUID/slot, memory total/used (MB), GPU utilization %, temperature (C), power draw (W), clock speed (MHz)

### AMD GPU

**Detection strategy:**

1. **AMD SMI library** (CGO, Linux only) - `github.com/hhk7734/amdsmi.go`
   - `amdsmi.New()` -> `Init(INIT_AMD_GPUS)` -> `Sockets()` -> `Processors()` -> `GPUMemoryTotal()`, `GPUMemoryUsage()`, `GPUMetricsInfo()`
   - File: `internal/provider/system/hardware_amd.go` (build: `linux && cgo`)
2. **rocm-smi CLI** (Linux only) - detection only via `exec.LookPath("rocm-smi")`
   - File: `internal/provider/system/detect.go` (`detectAmd`)
3. **sysfs fallback** - `/sys/class/drm/` enumeration with vendor ID `1002`
   - File: `internal/provider/system/hardware.go` (`fetchSysfsGpus`)

**Metrics collected:** VRAM total/used, GFX activity %, hotspot temperature, socket power (W), GFX clock (MHz)

### Intel GPU

**Detection:** sysfs `/sys/class/drm/` enumeration with PCI vendor ID `8086`
- File: `internal/provider/system/hardware.go` (`fetchSysfsGpus`, `readIntelTemperature`, `readIntelFrequency`)
- Driver match: `i915`
- Temperature: `device/hwmon/hwmon*/temp1_input` or `device/thermal_zone/temp`
- Frequency: `device/gt/gt0/punit/gpu_freq_mhz` or `device/freq0/freq`

### Apple GPU (macOS)

**Detection:** `system_profiler SPDisplaysDataType`
- File: `internal/provider/system/hardware.go` (`fetchDarwinGpus`)
- Parses chipset model and VRAM from output

### Windows GPU

**Detection:** `wmic path win32_VideoController get name /format:list`
- File: `internal/provider/system/hardware.go` (`fetchWindowsGpus`)

### Battery

**Library:** `github.com/distatus/battery` - `battery.GetAll()`
- File: `internal/provider/system/hardware.go` (`BatteryCmd`)
- Detection: `upower` binary on Linux, always available on macOS/Windows

### PCI Device Name Database

- Hardcoded map of PCI vendor:device IDs to GPU names in `internal/provider/system/hardware.go`
- Functions: `getPciDeviceName()` (line 472), `getVramFromPciDb()` (line 527)
- Covers: NVIDIA (10de:*), AMD (1002:*), Intel (8086:*) devices

## SSH Remote Monitoring

**Implementation:** `internal/provider/remote/ssh.go`

**Mechanism:**
- Executes a remote shell script via system `ssh` binary (`os/exec`)
- Uses `BatchMode=yes` for non-interactive authentication (key-based only)
- Configurable per host: port, key path, timeout (default: 2s)

**Connection options built:**
```
ssh -o ConnectTimeout=N -o BatchMode=yes [-i key_path] [-o Port=N] user@host <script>
```

**Remote script collects:**
```bash
uptime                          # System uptime
cat /proc/loadavg               # Load averages
nproc                           # CPU count
cat /proc/meminfo | head -20    # Memory info
df -B1 --output=size,used,pcent /  # Disk usage
awk '{if(NR>2)print $1,$2,$3,$10}' /proc/net/dev | head -5  # Network stats
ps aux --no-headers | sort -k3 -rn | head -15  # Top processes
```

**Configuration:** `internal/config/config.go` - `RemoteHostConfig` struct
- Managed via CLI: `bub remote add`, `bub remote remove`, `bub remote list`
- Stored in JSON config at `~/.config/bubble-monitor/config.json`

## Container & Orchestration Detection

### Docker

**Detection:**
- Linux: Unix socket probe at `/var/run/docker.sock` or `/run/docker.sock`
- Windows: `docker info` command execution
- File: `internal/provider/system/container.go` lines 41-56

**Data collection:**
- `docker ps -a --format "{{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}\t{{.State}}\t{{.CreatedAt}}"` - Container listing
- `docker stats --no-stream --format "{{json .}}"` - Per-container CPU/memory/network stats
- File: `internal/provider/system/container.go` lines 126-256

### Kubernetes

**Detection:** `kubectl cluster-info` execution
- File: `internal/provider/system/container.go` lines 58-62

**Data collection:**
- `kubectl get pods -o json` - Full pod listing with status, resources, node info, container statuses, restart counts
- File: `internal/provider/system/container.go` lines 304-421

## Network Monitoring

**Implementation:** `internal/provider/system/network.go`, `connections.go`

**Data sources (gopsutil):**
- `net.IOCounters(true)` - Per-interface I/O stats (bytes sent/received, packets, errors, drop)
- `net.IOCounters(false)` - Aggregate network I/O for charts
- `net.Connections("all")` - All network connections (TCP/UDP, local/remote addresses, state, PID)

**Remote:** `/proc/net/dev` parsed via SSH (`internal/provider/remote/ssh.go`)

## System Logs

**Implementation:** `internal/provider/system/logs.go`

**Mechanism:** Executes `journalctl -n 1000 --no-pager -o short-iso`
- Linux only
- Falls back to error message if `journalctl` unavailable or no permissions
- Max log lines constant: `MaxLogLines = 1000` (`internal/provider/system/constants.go`)

**Other platforms:** Returns "Logs not implemented for {OS}" message

## System Services

**Implementation:** `internal/provider/system/services.go`

**Mechanism:** Executes `systemctl list-units --type=service --no-legend --all`
- Linux only (requires `systemctl`)
- Parses: unit name, sub-state (running/dead/exited), description

**Windows:** Service detection flagged as available (uses gopsutil internally)

**Other platforms:** No service listing

## Process Management

**Implementation:** `internal/provider/process/`

**Data sources (gopsutil):**
- `process.Pids()` - PID enumeration
- `process.NewProcess(pid)` - Per-process handle
- `proc.Name()`, `proc.CPUPercent()`, `proc.MemoryPercent()`, `proc.MemoryInfo()`
- `proc.Status()`, `proc.Cmdline()`, `proc.Username()`, `proc.OpenFiles()`
- `proc.CreateTime()`, `proc.Nice()`, `proc.Ppid()`
- `proc.Suspend()`, `proc.Resume()` - Process control

**External commands for priority:**
- Unix: `renice -n <value> -p <pid>` (`internal/provider/process/control.go` line 90)
- Windows: `powershell -NoProfile -Command "Get-Process -Id <pid> | foreach { $_.PriorityClass = '<class>' }"` (`internal/provider/process/control.go` line 68)
  - Priority classes: Idle, BelowNormal, Normal, AboveNormal, High, RealTime

**Performance optimizations:**
- Process cache with `sync.RWMutex` to avoid re-fetching static data
- String interning for process names/status (LRU, 5000 entry cap)
- `sync.Pool` for process slice reuse (capacity 500)
- Concurrent fetch guard via `atomic.Bool`
- Files: `internal/provider/process/list.go`, `internal/provider/process/utils.go`

## Virtual Machine Detection

**Implementation:** `internal/provider/system/vm.go`

**Detection strategy (layered, Linux):**
1. CPU flags from `/proc/cpuinfo` (hyperv, vmware, kvm, qemu, xen, parallels, bhyve, openvz, virtio)
2. DMI product/vendor from `/sys/class/dmi/id/` files (product_name, sys_vendor, board_vendor, bios_vendor)
3. `dmidecode -s system-product-name` or `dmidecode -s system-manufacturer` CLI

**Detection strategy (Windows):**
- PowerShell: `(Get-WmiObject Win32_ComputerSystem).Manufacturer`
- PowerShell: `Get-WindowsOptionalFeature -FeatureName Microsoft-Hyper-V -Online`

**Container type detection (Linux):**
- `/proc/1/cgroup` for docker/container strings
- `/.dockerenv` file existence
- `/run/.containerenv` file existence
- `/usr/bin/lxc-checkconfig` or `/sys/fs/cgroup/cpuset/lxc` for LXC

**Supported hypervisors:** VMware, VirtualBox, KVM, Hyper-V, QEMU, Xen, Parallels, Bhyve, OpenVZ, LXC

## External Commands Summary

| Command | Platform | Purpose | File |
|---------|----------|---------|------|
| `nvidia-smi --query-gpu=... --format=csv` | All | NVIDIA GPU detection and metrics | `internal/provider/system/hardware.go`, `detect_nvidia_stub.go` |
| `rocm-smi` (LookPath) | Linux | AMD GPU detection | `internal/provider/system/detect.go` |
| `system_profiler SPDisplaysDataType` | macOS | GPU discovery | `internal/provider/system/hardware.go` |
| `wmic path win32_VideoController` | Windows | GPU discovery | `internal/provider/system/hardware.go` |
| `lsblk -b -o NAME,SIZE,TYPE,MOUNTPOINT -J` | Linux | Unmounted disk partitions | `internal/provider/system/disk.go` |
| `journalctl -n 1000 --no-pager -o short-iso` | Linux | System logs | `internal/provider/system/logs.go` |
| `systemctl list-units --type=service --no-legend --all` | Linux | Service listing | `internal/provider/system/services.go` |
| `ssh -o BatchMode=yes ...` | All | Remote host monitoring | `internal/provider/remote/ssh.go` |
| `renice -n <val> -p <pid>` | Unix | Process priority change | `internal/provider/process/control.go` |
| `powershell -NoProfile -Command ...` | Windows | Process priority, VM detection | `internal/provider/process/control.go`, `vm.go` |
| `docker ps -a --format ...` | All | Container listing | `internal/provider/system/container.go` |
| `docker stats --no-stream --format json` | All | Container stats | `internal/provider/system/container.go` |
| `docker info` | Windows | Docker availability | `internal/provider/system/container.go` |
| `kubectl get pods -o json` | All | Kubernetes pod listing | `internal/provider/system/container.go` |
| `kubectl cluster-info` | All | K8s availability check | `internal/provider/system/container.go` |
| `dmidecode -s system-product-name` | Linux | VM hypervisor detection | `internal/provider/system/vm.go` |
| `upower` (LookPath) | Linux | Battery detection | `internal/provider/system/detect.go` |
| `systemctl` (LookPath) | Linux | Service capability detection | `internal/provider/system/detect.go` |

## Configuration File Handling

**Format:** JSON

**Location:** `~/.config/bubble-monitor/config.json` (via `os.UserConfigDir`)

**Implementation:** `internal/config/config.go`

**Load process:**
1. Resolve config path (default or `--config` flag)
2. Read and decode JSON via `json.NewDecoder`
3. Merge with defaults for any missing fields (thresholds, tabs, theme, refresh rate, etc.)
4. Allow CLI flag overrides for theme, refresh-rate, history-length

**Save:** Pretty-printed JSON via `json.Encoder` with indentation, triggered by `bub remote add/remove`, `bub config`

**Hot reload:** Polls file modification time every 2 seconds (`internal/config/watcher.go`)
- Sends `ConfigChangeMsg` to TUI on change

## Authentication & Identity

**No authentication framework.** The application:
- Uses system-level process information (requires appropriate OS permissions)
- SSH remote monitoring relies on existing SSH key-based auth (no password auth)
- Configured via `RemoteHostConfig.KeyPath` for specific SSH keys

## Monitoring & Observability

**Error Tracking:** None (no external error tracking service)

**Logs:**
- Optional file logging configured via `LoggingConfig.Path` and `LoggingConfig.Enabled`
- File: `internal/app/logging.go`
- GPU/AMD detection panics are recovered and logged via `log.Printf`

**Profiling:**
- `net/http/pprof` imported in `cmd/bub/main.go` (currently commented out but import retained)
- `cpu.prof` file exists in repo root

## CI/CD & Deployment

**Hosting:** GitHub releases via GoReleaser
- Repository: `github.com/N1xev/bubbleMonitor`

**CI Pipeline:** Not detected (no `.github/workflows/` or equivalent)

**Release Process:**
- GoReleaser builds binaries for: linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64, windows/arm64
- Archives: binary format, tar.gz, zip
- Changelog: auto-generated, excludes docs/test/chore/refactor commits

## Environment Configuration

**Required env vars:** None (all configuration via JSON config file and CLI flags)

**Secrets location:** No secrets management. SSH keys referenced by path in config.

## Webhooks & Callbacks

**Incoming:** None

**Outgoing:** None

---

*Integration audit: 2026-04-09*

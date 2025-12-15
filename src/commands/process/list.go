package process

import (
	"sort"
	"strings"
	"sync"

	tea "charm.land/bubbletea/v2"
	"github.com/shirou/gopsutil/v3/process"

	"github.com/N1xev/bubbleMonitor/src/data"
	"github.com/N1xev/bubbleMonitor/src/messages"
)

// CachedProcessInfo stores the actual process object and its static data
type CachedProcessInfo struct {
	Proc       *process.Process // Persistent object for accurate CPU deltas
	Name       string
	Username   string
	Cmdline    string
	CreateTime int64
	Nice       int32
	Ppid       int32
}

var (
	// processCache stores static info by PID
	processCache = make(map[int32]CachedProcessInfo)
	cacheMutex   sync.RWMutex
)

// ProcessesCmd fetches running processes and sorts them
func ProcessesCmd(sortBy string) tea.Cmd {
	return func() tea.Msg {
		// Use Pids() instead of Processes() -> Cheaper, returns only []int32
		pids, err := process.Pids()
		if err != nil {
			return messages.ProcessesMsg{}
		}

		// Pre-allocate to avoid re-sizing (Memory Optimization)
		procList := make([]data.ProcessInfo, 0, len(pids))

		// Map for quick lookup of current PIDs to clean up cache
		currentPids := make(map[int32]bool)

		cacheMutex.Lock()
		defer cacheMutex.Unlock()

		for _, pid := range pids {
			currentPids[pid] = true

			// Try to get from cache first
			cached, exists := processCache[pid]

			if !exists {
				// Create NEW process object only once
				newProc, err := process.NewProcess(pid)
				if err != nil {
					continue // Process might have died between Pids() and NewProcess()
				}

				// Fetch static data (Expensive calls on Windows)
				name, _ := newProc.Name()
				username, _ := newProc.Username()
				createTime, _ := newProc.CreateTime()
				cmdline, _ := newProc.Cmdline()
				nice, _ := newProc.Nice()
				ppid, _ := newProc.Ppid()

				cached = CachedProcessInfo{
					Proc:       newProc,
					Name:       name,
					Username:   username,
					Cmdline:    cmdline,
					CreateTime: createTime,
					Nice:       nice,
					Ppid:       ppid,
				}
				processCache[pid] = cached
			}

			// Always fetch dynamic data (CPU, Memory, Status) using the PERSISTENT object
			// This allows gopsutil to calculate true CPU usage over time intervals
			cpuPercent, _ := cached.Proc.CPUPercent()
			memPercent, _ := cached.Proc.MemoryPercent()
			status, _ := cached.Proc.Status()
			memInfo, _ := cached.Proc.MemoryInfo()

			var memBytes uint64
			if memInfo != nil {
				memBytes = memInfo.RSS
			}

			// Get a readable status
			statusStr := strings.Join(status, ",")
			if statusStr == "" {
				statusStr = "running"
			}

			procList = append(procList, data.ProcessInfo{
				Name:        cached.Name,
				Pid:         pid,
				Cpu:         cpuPercent,
				Memory:      float64(memPercent),
				Status:      statusStr,
				Username:    cached.Username,
				CreateTime:  cached.CreateTime,
				Cmdline:     cached.Cmdline,
				MemoryBytes: memBytes,
				Nice:        cached.Nice,
				Ppid:        cached.Ppid,
			})
		}

		// Clean up cache: remove PIDs that are no longer running
		for pid := range processCache {
			if !currentPids[pid] {
				delete(processCache, pid)
			}
		}

		// Sort in background thread
		sort.Slice(procList, func(i, j int) bool {
			switch sortBy {
			case "cpu":
				return procList[i].Cpu > procList[j].Cpu
			case "memory":
				return procList[i].Memory > procList[j].Memory
			case "pid":
				return procList[i].Pid > procList[j].Pid
			default:
				return false
			}
		})

		return messages.ProcessesMsg(procList)
	}
}

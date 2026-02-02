package app

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/shirou/gopsutil/v3/net"

	"github.com/N1xev/bubbleMonitor/internal/config"
	"github.com/N1xev/bubbleMonitor/internal/data"
	messages "github.com/N1xev/bubbleMonitor/internal/msg"
	"github.com/N1xev/bubbleMonitor/internal/provider/process"
	"github.com/N1xev/bubbleMonitor/internal/provider/remote"
	"github.com/N1xev/bubbleMonitor/internal/provider/system"
)

func KillProcessCmd(pid int32) tea.Cmd {
	return func() tea.Msg {
		proc, err := os.FindProcess(int(pid))
		if err != nil {
			return messages.KillProcessMsg{Pid: pid, Success: false, Error: err.Error()}
		}
		err = proc.Kill()
		if err != nil {
			return messages.KillProcessMsg{Pid: pid, Success: false, Error: err.Error()}
		}
		return messages.KillProcessMsg{Pid: pid, Success: true}
	}
}

// Update handles all messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case messages.KillProcessMsg:
		m.ShowKillDialog = false
		m.KillTargetPid = 0
		m.KillTargetName = ""
		return m, process.ProcessesCmd(m.SortBy)

	case messages.PriorityChangeMsg:
		if msg.Err != nil {
			return m, AddToastCmd(fmt.Sprintf("Priority Error: %v", msg.Err), data.ToastError)
		}
		return m, tea.Batch(process.ProcessesCmd(m.SortBy), AddToastCmd("Priority Changed", data.ToastSuccess))

	case messages.ProcessControlMsg:
		if msg.Err != nil {
			return m, AddToastCmd(fmt.Sprintf("%s Failed: %v", strings.Title(msg.Action), msg.Err), data.ToastError)
		}
		switch msg.Action {
		case "suspend":
			m.SetSuspended(msg.Pid, true)
		case "resume":
			m.SetSuspended(msg.Pid, false)
		}
		return m, tea.Batch(process.ProcessesCmd(m.SortBy), AddToastCmd(fmt.Sprintf("Process %sd", strings.Title(msg.Action)), data.ToastSuccess))

	case messages.ToastMsg:
		id := m.NextToastID
		m.NextToastID++
		t := data.Toast{
			ID:        id,
			Message:   msg.Message,
			Level:     msg.Level,
			StartTime: time.Now(),
			Duration:  msg.Duration,
		}
		m.Toasts = append(m.Toasts, t)
		return m, TickToastCmd(id, msg.Duration)

	case messages.ToastTimeoutMsg:
		for i := 0; i < len(m.Toasts); i++ {
			if m.Toasts[i].ID == msg.ID {
				m.Toasts = append(m.Toasts[:i], m.Toasts[i+1:]...)
				break
			}
		}

	case messages.OpenFilesMsg:
		if msg.Err != nil {
			m.ShowOpenFiles = false
			return m, AddToastCmd(fmt.Sprintf("Open Files Error: %v", msg.Err), data.ToastError)
		}
		m.OpenFilesList = msg.Files
		m.OpenFilesPid = msg.Pid

		var lines []string
		for _, f := range msg.Files {
			lines = append(lines, f.Path)
		}
		if len(lines) == 0 {
			lines = []string{"No open files found or access denied."}
		}
		content := strings.Join(lines, "\n")
		m.OpenFilesView.SetContent(content)
		m.OpenFilesView.GotoTop()

		return m, nil

	case config.ConfigWatchTickMsg:
		return m, config.WatchConfig(m.LastConfigModTime)

	case config.ConfigChangeMsg:
		m.LastConfigModTime = msg.NewModTime
		newConfig, err := config.LoadConfig()
		if err == nil {
			// Prevent reload loop if config effectively hasn't changed (e.g. self-save)
			if reflect.DeepEqual(m.Config, newConfig) {
				return m, config.WatchConfig(m.LastConfigModTime)
			}

			m.Config = newConfig
			m.HistoryLength = newConfig.HistoryLength
			m.ChartType = newConfig.ChartType
			m.SortBy = newConfig.SortBy
			oldTreeView := m.TreeView
			m.TreeView = newConfig.ViewType == "tree"
			if oldTreeView != m.TreeView {
				m.InvalidateProcessCache()
			}
			m.Theme = newConfig.Theme
			m.RefreshRate = newConfig.RefreshRate
			m.BorderType = newConfig.BorderType
			m.BorderStyle = newConfig.BorderStyle
			m.BackgroundOpaque = newConfig.BackgroundOpaque
			return m, tea.Batch(config.WatchConfig(m.LastConfigModTime), AddToastCmd("Config Reloaded", data.ToastSuccess))
		}
		return m, tea.Batch(config.WatchConfig(m.LastConfigModTime), AddToastCmd(fmt.Sprintf("Config Error: %v", err), data.ToastError))

	case tea.KeyMsg:
		// Handle kill dialog first
		if m.ShowKillDialog {
			switch msg.String() {
			case "y", "enter":
				pid := m.KillTargetPid
				m.ShowKillDialog = false
				return m, KillProcessCmd(pid)
			case "n", "esc":
				m.ShowKillDialog = false
				m.KillTargetPid = 0
				m.KillTargetName = ""
			}
			return m, nil
		}

		// Handle help overlay
		if m.ShowHelp {
			if msg.String() == "?" || msg.String() == "esc" {
				m.ShowHelp = false
			}
			return m, nil
		}

		// Handle Open Files overlay
		if m.ShowOpenFiles {
			switch msg.String() {
			case "o", "esc":
				m.ShowOpenFiles = false
				m.OpenFilesView.GotoTop()
			default:

				switch msg.String() {
				case "j", "down":
					m.OpenFilesView.LineDown(1)
				case "k", "up":
					m.OpenFilesView.LineUp(1)
				case "pgdown", "ctrl+d":
					m.OpenFilesView.HalfViewDown()
				case "pgup", "ctrl+u":
					m.OpenFilesView.HalfViewUp()
				case "home":
					m.OpenFilesView.GotoTop()
				case "end":
					m.OpenFilesView.GotoBottom()
				}
				return m, nil
			}
			return m, nil
		}

		// Handle filter mode
		if m.FilterMode {
			switch msg.String() {
			case "esc":
				m.FilterMode = false
			case "backspace":
				if len(m.ProcessFilter) > 0 {
					m.ProcessFilter = m.ProcessFilter[:len(m.ProcessFilter)-1]
					m.InvalidateProcessCache()
				}
			case "enter":
				m.FilterMode = false
			default:
				if len(msg.String()) == 1 {
					m.ProcessFilter += msg.String()
					m.InvalidateProcessCache()
				}
			}
			m.SelectedProcess = 0
			m.ProcessScrollOffset = 0
			return m, nil
		}

		// Settings overlay key handling
		if m.ShowSettings {
			// Total settings: 19
			// 4 Thresholds (0-3)
			// 4 Display (4-7)
			// 6 Tabs (8-13)
			// 5 Appearance (14-18)
			totalSettings := 19

			switch msg.String() {
			case "esc", ".":
				m.ShowSettings = false
				config.SaveConfig(m.Config)
				return m, nil
			case "up", "k":
				m.SettingsIdx = (m.SettingsIdx - 1 + totalSettings) % totalSettings
				if m.SettingsIdx < 4 {
					metrics := []config.MetricType{config.MetricCPU, config.MetricMem, config.MetricDisk, config.MetricTemp}
					m.SettingsSel = metrics[m.SettingsIdx]
				}
			case "down", "j":
				m.SettingsIdx = (m.SettingsIdx + 1) % totalSettings
				if m.SettingsIdx < 4 {
					metrics := []config.MetricType{config.MetricCPU, config.MetricMem, config.MetricDisk, config.MetricTemp}
					m.SettingsSel = metrics[m.SettingsIdx]
				}
			case "+", "=", "right", "l":
				if m.SettingsIdx < 4 {
					curr := m.Config.Thresholds[m.SettingsSel]
					if curr < 100 {
						m.Config.Thresholds[m.SettingsSel] = curr + 1
					}
				} else {
					m.handleSettingsChange(1)
					return m, AddToastCmd("Setting Changed", data.ToastSuccess)
				}
			case "-", "_", "left", "h":
				if m.SettingsIdx < 4 {
					curr := m.Config.Thresholds[m.SettingsSel]
					if curr > 0 {
						m.Config.Thresholds[m.SettingsSel] = curr - 1
					}
				} else {
					m.handleSettingsChange(-1)
					return m, AddToastCmd("Setting Changed", data.ToastSuccess)
				}
			}
			return m, nil
		}

		// Global Keybindings
		switch msg.String() {
		case "q", "ctrl+c":
			config.SaveConfig(m.Config)
			return m, tea.Quit
		case "e":
			return m, SaveSnapshotCmd(m.Cpu, m.Memory, m.Disk, m.ProcessCount)
		case ".":
			m.ShowSettings = !m.ShowSettings
			if !m.ShowSettings {
				config.SaveConfig(m.Config)
			}
			return m, nil
		case "space":
			if m.TreeView && m.SelectedTab == 2 {
				procs, _ := m.GetVisibleProcesses()
				if m.SelectedProcess >= 0 && m.SelectedProcess < len(procs) {
					proc := procs[m.SelectedProcess]
					m.ToggleCollapsed(proc.Pid)
					m.InvalidateProcessCache()
				}
			}
			return m, nil
		case "n":
			m.ProcessCpuNormalized = !m.ProcessCpuNormalized
			m.Config.ProcessCpuNormalized = m.ProcessCpuNormalized
			config.SaveConfig(m.Config)
			status := "Raw (>100%)"
			if m.ProcessCpuNormalized {
				status = "Normalized (0-100%)"
			}
			return m, AddToastCmd("CPU Mode: "+status, data.ToastInfo)
		}

		currentTab := "Overview"
		if m.SelectedTab < len(m.ActiveTabs) {
			currentTab = m.ActiveTabs[m.SelectedTab]
		}

		switch msg.String() {
		case "tab", "right", "l":
			if len(m.ActiveTabs) > 0 {
				m.SelectedTab = (m.SelectedTab + 1) % len(m.ActiveTabs)
			}
			if len(m.ActiveTabs) > 0 && currentTab != "Processes" && m.ActiveTabs[m.SelectedTab] != "Processes" {
				m.SelectedProcess = 0
				m.ProcessScrollOffset = 0
			}
		case "shift+tab", "left":
			if len(m.ActiveTabs) > 0 {
				m.SelectedTab = (m.SelectedTab - 1 + len(m.ActiveTabs)) % len(m.ActiveTabs)
			}
			if len(m.ActiveTabs) > 0 && currentTab != "Processes" && m.ActiveTabs[m.SelectedTab] != "Processes" {
				m.SelectedProcess = 0
				m.ProcessScrollOffset = 0
			}
		case "H":
			// Cycle history length
			switch m.HistoryLength {
			case 60:
				m.HistoryLength = 300
			case 300:
				m.HistoryLength = 900
			case 900:
				m.HistoryLength = 3600
			default:
				m.HistoryLength = 60
			}
			m.Config.HistoryLength = m.HistoryLength
			// Resize Ring Buffers
			m.CpuHistory = data.NewRingBuffer(m.HistoryLength)
			m.MemHistory = data.NewRingBuffer(m.HistoryLength)
			m.NetHistory = data.NewRingBuffer(m.HistoryLength)
			m.SwapHistory = data.NewRingBuffer(m.HistoryLength)
			m.HistoryTemp = data.NewRingBuffer(m.HistoryLength)
			m.DiskHORead = data.NewRingBuffer(m.HistoryLength)
			m.DiskHOWrite = data.NewRingBuffer(m.HistoryLength)
		case "C":
			// Cycle chart type (Metrics tab)
			switch m.ChartType {
			case "line":
				m.ChartType = "bar"
			case "bar":
				m.ChartType = "braille"
			default:
				m.ChartType = "line"
			}
			// Update config
			m.Config.ChartType = m.ChartType
		case "pgup":
			if currentTab == "Processes" {
				m.SelectedProcess -= 10
				if m.SelectedProcess < 0 {
					m.SelectedProcess = 0
				}
				if m.SelectedProcess < m.ProcessScrollOffset {
					m.ProcessScrollOffset = m.SelectedProcess
				}
			} else if currentTab == "Metrics" {
				m.CpuCoreScrollOffset -= 4
				if m.CpuCoreScrollOffset < 0 {
					m.CpuCoreScrollOffset = 0
				}
			}
		case "pgdown":
			if currentTab == "Processes" {
				visibleProcs, _ := m.GetVisibleProcesses()
				filteredLen := len(visibleProcs)
				if filteredLen > 0 {
					m.SelectedProcess += 10
					if m.SelectedProcess >= filteredLen {
						m.SelectedProcess = filteredLen - 1
					}
					// Ensure scroll follows
					rows := m.getVisibleProcessRows()
					if m.SelectedProcess >= m.ProcessScrollOffset+rows {
						m.ProcessScrollOffset = m.SelectedProcess - rows + 1
					}
				}
			} else if currentTab == "Metrics" {
				m.CpuCoreScrollOffset += 4
			}
		case "1":
			m.SelectedTab = 0
		case "2":
			m.SelectedTab = 1
		case "3":
			m.SelectedTab = 2
		case "4":
			m.SelectedTab = 3
		case "5":
			m.SelectedTab = 4
		case "S":
			switch m.SortBy {
			case "cpu":
				m.SortBy = "mem"
			case "mem":
				m.SortBy = "pid"
			default:
				m.SortBy = "cpu"
			}
			m.Config.SortBy = m.SortBy
			m.InvalidateProcessCache()

		case "b":
			if currentTab == "Processes" {
				procs, _ := m.GetVisibleProcesses()
				if m.SelectedProcess >= 0 && m.SelectedProcess < len(procs) {
					proc := procs[m.SelectedProcess]
					m.ToggleBookmark(proc.Pid)
					m.InvalidateProcessCache()
				}
			}
		case "T":
			if currentTab == "Processes" {
				m.TreeView = !m.TreeView
				viewName := "normal"
				if m.TreeView {
					viewName = "tree"
				}
				m.Config.ViewType = viewName
				m.InvalidateProcessCache()
			}
		case "+", "=":
			if currentTab == "Processes" {
				filtered := m.GetFilteredProcesses()
				if m.SelectedProcess < len(filtered) {
					proc := filtered[m.SelectedProcess]
					// Decrease delta (increase priority: -1 on Unix, Step Up on Windows logic)
					return m, process.ReniceProcessCmdSafe(proc.Pid, -1)
				}
			}
		case "-", "_":
			if currentTab == "Processes" {
				filtered := m.GetFilteredProcesses()
				if m.SelectedProcess < len(filtered) {
					proc := filtered[m.SelectedProcess]
					// Increase delta (decrease priority: +1 on Unix, Step Down on Windows logic)
					return m, process.ReniceProcessCmdSafe(proc.Pid, 1)
				}
			}
		case "o":
			if currentTab == "Processes" {
				// Toggle Open Files Inspector
				if m.ShowOpenFiles {
					m.ShowOpenFiles = false
				} else {
					filtered := m.GetFilteredProcesses()
					if m.SelectedProcess < len(filtered) {
						proc := filtered[m.SelectedProcess]
						m.ShowOpenFiles = true
						m.OpenFilesList = nil // Clear previous
						m.OpenFilesPid = proc.Pid
						return m, process.FetchOpenFilesCmd(proc.Pid)
					}
				}
			}

		case "z":
			if currentTab == "Processes" {
				filtered := m.GetFilteredProcesses()
				if m.SelectedProcess < len(filtered) {
					proc := filtered[m.SelectedProcess]
					return m, process.SuspendProcessCmd(proc.Pid)
				}
			}
		case "x":
			if currentTab == "Processes" {
				filtered := m.GetFilteredProcesses()
				if m.SelectedProcess < len(filtered) {
					proc := filtered[m.SelectedProcess]
					return m, process.ResumeProcessCmd(proc.Pid)
				}
			}
		case "p":
			m.Paused = !m.Paused
		case "r":
			return m, tea.Batch(
				system.FastMetricsCmd(),
				process.ProcessesCmd(m.SortBy),
				system.HostInfoCmd(),
				system.DiskInfoCmd(),
				system.GpuInfoCmd(),
			)
		case "?":
			m.ShowHelp = true
			m.LastError = ""

			// Process navigation (only on Processes tab)
		case "j", "down":
			if currentTab == "Processes" {
				// Use visible processes (tree aware)
				visibleProcs, _ := m.GetVisibleProcesses()
				filteredLen := len(visibleProcs)

				if filteredLen > 0 && m.SelectedProcess < filteredLen-1 {
					m.SelectedProcess++
					visibleRows := m.getVisibleProcessRows()
					if m.SelectedProcess >= m.ProcessScrollOffset+visibleRows {
						m.ProcessScrollOffset = m.SelectedProcess - visibleRows + 1
					}
				}
			}
		case "k", "up":
			if currentTab == "Processes" {
				if m.SelectedProcess > 0 {
					m.SelectedProcess--
					if m.SelectedProcess < m.ProcessScrollOffset {
						m.ProcessScrollOffset = m.SelectedProcess
					}
				}
			}
		case "K":
			if currentTab == "Processes" && len(m.Processes) > 0 {
				filtered := m.GetFilteredProcesses()
				if m.SelectedProcess < len(filtered) {
					proc := filtered[m.SelectedProcess]
					m.ShowKillDialog = true
					m.KillTargetPid = proc.Pid
					m.KillTargetName = proc.Name
				}
			}
		case "f":
			if currentTab == "Processes" {
				m.FilterMode = true
			}
		case "c":
			if currentTab == "Processes" {
				m.ProcessFilter = ""
				m.SelectedProcess = 0
				m.ProcessScrollOffset = 0
				m.InvalidateProcessCache()
			}
		case "g":
			if currentTab == "Processes" {
				m.SelectedProcess = 0
				m.ProcessScrollOffset = 0
			}
		case "G":
			if currentTab == "Processes" {
				filteredLen := m.getFilteredProcessCount()
				if filteredLen > 0 {
					m.SelectedProcess = filteredLen - 1
					visibleRows := m.getVisibleProcessRows()
					if m.SelectedProcess >= visibleRows {
						m.ProcessScrollOffset = m.SelectedProcess - visibleRows + 1
					}
				}
			}
		}

	case messages.TickMsg:
		if m.AlertManager != nil {
			m.AlertManager.CheckAlerts(&m.AppState)
		}

		m.TickCount++
		cmds := []tea.Cmd{
			system.TickCmd(time.Duration(m.RefreshRate) * time.Millisecond),
		}

		cmds = append(cmds, system.FastMetricsCmd())

		currentTab := "Overview"
		if m.SelectedTab < len(m.ActiveTabs) {
			currentTab = m.ActiveTabs[m.SelectedTab]
		}

		// Update Process List, Network, Disk IO every 2nd tick (2s)
		if m.TickCount%2 == 0 {
			cmds = append(cmds,
				system.DiskIOCmd(),
				system.NetworkInterfacesCmd(),
				system.TempCmd(),
			)

			// Optimize CPU usage: Only fetch full process list if visible
			if currentTab == "Processes" {
				cmds = append(cmds, process.ProcessesCmd(m.SortBy))
			} else {
				cmds = append(cmds, process.ProcessCountCmd())
			}

			if currentTab == "Connections" {
				cmds = append(cmds, system.ConnectionsCmd())
			}
		}

		// Update Slow Metrics (Disk Usage/Net Totals) every 5th tick (5s)
		if m.TickCount%5 == 0 {
			cmds = append(cmds,
				system.SlowMetricsCmd(),
				system.BatteryCmd(),
				system.HostInfoCmd(),
				system.DiskInfoCmd(),
				LogMetricsCmd(m.Cpu, m.Memory, m.Disk, m.NetSentRate+m.NetRecvRate, m.ProcessCount, m.Config.Logging.Enabled, m.Config.Logging.Path),
			)

			if currentTab == "Services" {
				cmds = append(cmds, system.ServicesCmd())
			}
			if currentTab == "Logs" {
				cmds = append(cmds, system.SystemLogsCmd())
			}
			if currentTab == "Remote" {
				for _, h := range m.Config.RemoteHosts {
					cmds = append(cmds, remote.CheckRemoteCmd(h.Host))
				}
			}
		}

		return m, tea.Batch(cmds...)

	case messages.CpuMemMsg:
		if msg.Err != nil {
			return m, AddToastCmd("CPU/Memory metrics error: "+msg.Err.Error(), data.ToastError)
		}
		m.Cpu = msg.Cpu
		m.CpuPerCore = msg.CpuPerCore
		m.Memory = msg.Memory
		m.Swap = msg.Swap
		m.LoadAvg = msg.LoadAvg
		m.MemInfo = msg.MemInfo
		m.SwapInfo = msg.SwapInfo

		m.CpuHistory.Push(m.Cpu)
		m.MemHistory.Push(m.Memory)
		m.SwapHistory.Push(m.Swap)
		UpdateHealthScore(&m.AppState)

	case messages.DiskNetMsg:
		if msg.Err != nil {
			return m, AddToastCmd("Disk/Network metrics error: "+msg.Err.Error(), data.ToastError)
		}
		if m.LastNetSent > 0 && m.LastNetRecv > 0 {
			m.NetSentRate = float64(msg.NetSent-m.LastNetSent) / 1024 / 1024
			m.NetRecvRate = float64(msg.NetRecv-m.LastNetRecv) / 1024 / 1024
		}
		// Fallback for huge jumps (first run)
		if m.NetSentRate < 0 {
			m.NetSentRate = 0
		}
		if m.NetRecvRate < 0 {
			m.NetRecvRate = 0
		}

		m.LastNetSent = msg.NetSent
		m.LastNetRecv = msg.NetRecv
		m.Disk = msg.Disk

		totalNetRate := m.NetSentRate + m.NetRecvRate
		netPercent := (totalNetRate / 10) * 100 // Arbitrary 10MB/s cap for graph scaling logic
		if netPercent > 100 {
			netPercent = 100
		}

		m.NetHistory.Push(netPercent)

	case messages.ProcessCountMsg:
		m.ProcessCount = int(msg)

	case messages.ProcessesMsg:
		allProcesses := msg
		m.ProcessCount = len(allProcesses)

		// Prune dead processes from AppState maps to prevent memory leaks
		currentPids := make(map[int32]bool, len(allProcesses))
		for _, p := range allProcesses {
			currentPids[p.Pid] = true
		}

		m.PruneDeadProcessMaps(currentPids)

		maxProcesses := len(allProcesses)
		if len(allProcesses) > maxProcesses {
			allProcesses = allProcesses[:maxProcesses]
		}

		// Recycle old slice buffer
		if cap(m.Processes) > 0 {
			oldSlice := m.Processes
			process.PutProcSlice(&oldSlice)
		}

		m.Processes = allProcesses
		m.InvalidateProcessCache()
		UpdateHealthScore(&m.AppState)
		UpdateProcessHistory(&m.AppState, currentPids)

		filteredLen := m.getFilteredProcessCount()
		if m.SelectedProcess >= filteredLen {
			m.SelectedProcess = filteredLen - 1
			if m.SelectedProcess < 0 {
				m.SelectedProcess = 0
			}
		}

	case messages.HostInfoMsg:
		if msg.Err != nil {
			return m, AddToastCmd("Host info error: "+msg.Err.Error(), data.ToastError)
		}
		m.HostInfo = msg.Info
	case messages.DiskInfoMsg:
		if msg.Err != nil {
			return m, AddToastCmd("Disk partitions error: "+msg.Err.Error(), data.ToastError)
		}
		m.DiskPartitions = msg.Partitions
	case messages.GpuInfoMsg:
		m.GpuInfo = msg
	case messages.DiskIOMsg:
		if len(m.LastDiskIO) > 0 {
			var totalRead, totalWrite uint64
			var lastTotalRead, lastTotalWrite uint64

			for k, v := range msg {
				totalRead += v.ReadBytes
				totalWrite += v.WriteBytes
				if last, ok := m.LastDiskIO[k]; ok {
					lastTotalRead += last.ReadBytes
					lastTotalWrite += last.WriteBytes
				}
			}

			if totalRead >= lastTotalRead {
				m.DiskReadRate = float64(totalRead-lastTotalRead) / 1024 / 1024
			}
			if totalWrite >= lastTotalWrite {
				m.DiskWriteRate = float64(totalWrite-lastTotalWrite) / 1024 / 1024
			}

			m.DiskHORead.Push(m.DiskReadRate)
			m.DiskHOWrite.Push(m.DiskWriteRate)
		}
		m.DiskIO = msg
		m.LastDiskIO = msg
	case messages.TempMsg:
		m.Sensors = msg
		// Calculate CPU temp (average of coretemp or k10temp)
		var totalTemp float64
		var count int
		for _, t := range msg {
			key := strings.ToLower(t.SensorKey)
			if strings.Contains(key, "core") || strings.Contains(key, "cpu") || strings.Contains(key, "package") {
				if t.Temperature > 0 {
					totalTemp += t.Temperature
					count++
				}
			}
		}
		if count > 0 {
			m.CpuTemp = totalTemp / float64(count)
		} else if len(msg) > 0 {
			// Fallback: take first sensor
			m.CpuTemp = msg[0].Temperature
		}

		m.HistoryTemp.Push(m.CpuTemp)

	case messages.NetworkInterfacesMsg:
		if msg.Err != nil {
			return m, AddToastCmd("Network interfaces error: "+msg.Err.Error(), data.ToastError)
		}
		m.NetworkInterfaces = msg.Interfaces
		if m.LastNetworkInterfaces == nil {
			m.LastNetworkInterfaces = make(map[string]net.IOCountersStat)
		}

		for _, nic := range msg.Interfaces {
			m.LastNetworkInterfaces[nic.Name] = nic
		}

	case messages.BatteryMsg:
		m.Battery = msg
	case messages.ServicesMsg:
		m.Services = msg
	case messages.ConnectionsMsg:
		m.Connections = msg
	case messages.SysLogMsg:
		m.SystemLogs = msg
	case messages.RemoteMsg:
		m.RemoteUptimes[msg.Host] = msg.Output
	}
	return m, nil
}

// AddToastCmd creates a command to show a toast
func AddToastCmd(msg string, level string) tea.Cmd {
	return func() tea.Msg {
		return messages.ToastMsg{Message: msg, Level: level, Duration: 3 * time.Second}
	}
}

// TickToastCmd waits for duration then sends timeout
func TickToastCmd(id int64, duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return messages.ToastTimeoutMsg{ID: id}
	})
}

// getFilteredProcessCount returns the count of filtered processes
func (m Model) getFilteredProcessCount() int {
	return len(m.GetFilteredProcesses())
}

// getVisibleProcessRows returns how many process rows can be displayed
func (m Model) getVisibleProcessRows() int {
	rows := m.Height - 19
	if rows < 3 {
		rows = 3
	}
	return rows
}

// Helper to sort tabs according to standard order
func sortActiveTabs(active []string) []string {
	standard := []string{"Overview", "Metrics", "Processes", "Disks", "Network", "System"}
	var sorted []string
	for _, std := range standard {
		for _, act := range active {
			if act == std {
				sorted = append(sorted, act)
			}
		}
	}
	return sorted
}

// handleSettingsChange handles non-threshold settings updates
func (m *Model) handleSettingsChange(dir int) {
	switch m.SettingsIdx {
	case 4: // Chart Type
		types := []string{"line", "bar", "braille"}
		for i, t := range types {
			if t == m.ChartType {
				nextIdx := (i + dir + len(types)) % len(types)
				m.ChartType = types[nextIdx]
				break
			}
		}
		m.Config.ChartType = m.ChartType

	case 5:
		m.TreeView = !m.TreeView
		viewName := "normal"
		if m.TreeView {
			viewName = "tree"
		}
		m.Config.ViewType = viewName
		m.InvalidateProcessCache()

	case 6:
		opts := []string{"cpu", "mem", "pid"}
		for i, o := range opts {
			if o == m.SortBy {
				nextIdx := (i + dir + len(opts)) % len(opts)
				m.SortBy = opts[nextIdx]
				break
			}
		}
		m.Config.SortBy = m.SortBy
		m.InvalidateProcessCache()

	case 7: // History Length
		lens := []int{60, 300, 900, 3600}
		for i, l := range lens {
			if l == m.HistoryLength {
				nextIdx := (i + dir + len(lens)) % len(lens)
				m.HistoryLength = lens[nextIdx]
				break
			}
		}
		m.Config.HistoryLength = m.HistoryLength

	case 8: // Process CPU Normalized
		m.ProcessCpuNormalized = !m.ProcessCpuNormalized
		m.Config.ProcessCpuNormalized = m.ProcessCpuNormalized

	case 9, 10, 11, 12, 13, 14: // Tabs
		allTabs := []string{"Overview", "Metrics", "Processes", "Disks", "Network", "System"}
		tabIdx := m.SettingsIdx - 9
		if tabIdx >= 0 && tabIdx < len(allTabs) {
			targetTab := allTabs[tabIdx]

			idxInActive := -1
			for i, t := range m.ActiveTabs {
				if t == targetTab {
					idxInActive = i
					break
				}
			}

			if idxInActive >= 0 {
				m.ActiveTabs = append(m.ActiveTabs[:idxInActive], m.ActiveTabs[idxInActive+1:]...)
			} else {
				m.ActiveTabs = append(m.ActiveTabs, targetTab)
				m.ActiveTabs = sortActiveTabs(m.ActiveTabs)
			}
			m.Config.Tabs = m.ActiveTabs
		}

	case 15: // Theme
		themes := config.GetThemeNames()
		for i, t := range themes {
			if t == m.Theme {
				nextIdx := (i + dir + len(themes)) % len(themes)
				m.Theme = themes[nextIdx]
				break
			}
		}
		m.Config.Theme = m.Theme
		if m.Theme == "custom" && m.Config.CustomTheme == nil {
			m.Config.CustomTheme = config.DefaultCustomTheme()
		}

	case 16: // Refresh Rate
		rates := config.GetRefreshRates()
		for i, r := range rates {
			if r == m.RefreshRate {
				nextIdx := (i + dir + len(rates)) % len(rates)
				m.RefreshRate = rates[nextIdx]
				break
			}
		}
		m.Config.RefreshRate = m.RefreshRate

	case 17: // Border Type
		types := config.GetBorderTypes()
		for i, t := range types {
			if t == m.BorderType {
				nextIdx := (i + dir + len(types)) % len(types)
				m.BorderType = types[nextIdx]
				break
			}
		}
		m.Config.BorderType = m.BorderType

	case 18: // Border Style
		styles := config.GetBorderStyles()
		for i, s := range styles {
			if s == m.BorderStyle {
				nextIdx := (i + dir + len(styles)) % len(styles)
				m.BorderStyle = styles[nextIdx]
				break
			}
		}
		m.Config.BorderStyle = m.BorderStyle

	case 19: // Background
		m.BackgroundOpaque = !m.BackgroundOpaque
		m.Config.BackgroundOpaque = m.BackgroundOpaque
	}
}

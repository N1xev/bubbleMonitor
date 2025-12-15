package model

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/N1xev/bubbleMonitor/src/commands/process"
	"github.com/N1xev/bubbleMonitor/src/commands/system"
	configpkg "github.com/N1xev/bubbleMonitor/src/config"
	"github.com/N1xev/bubbleMonitor/src/data"
	"github.com/N1xev/bubbleMonitor/src/ui"
	"github.com/shirou/gopsutil/v3/cpu"
)

type Model struct {
	data.AppState
}

// Init initializes the model and returns start commands
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		system.TickCmd(time.Duration(m.RefreshRate)*time.Millisecond),
		system.FastMetricsCmd(),
		system.SlowMetricsCmd(),
		process.ProcessesCmd(m.SortBy),
		system.HostInfoCmd(),
		system.DiskInfoCmd(),
		system.DiskIOCmd(),
		system.TempCmd(),
		system.NetworkInterfacesCmd(),
		system.BatteryCmd(),
		system.GpuInfoCmd(),
		configpkg.WatchConfig(m.LastConfigModTime),
	)
}

// View renders the UI
func (m Model) View() tea.View {
	return ui.RenderFromAppState(&m.AppState)
}

// GetBorder returns the lipgloss border for the current style
func (m Model) GetBorder() lipgloss.Border {
	switch m.BorderStyle {
	case "double":
		return lipgloss.DoubleBorder()
	case "dashed":
		return lipgloss.Border{
			Top:         "-",
			Bottom:      "-",
			Left:        "|",
			Right:       "|",
			TopLeft:     "+",
			TopRight:    "+",
			BottomLeft:  "+",
			BottomRight: "+",
		}
	default: // single
		if m.BorderType == "rounded" {
			return lipgloss.RoundedBorder()
		}
		return lipgloss.NormalBorder()
	}
}

// InitialModel creates a new Model with default values
func InitialModel() Model {
	// Load Configuration
	cfg, err := configpkg.LoadConfig()
	if err != nil {
		// Fallback to default if error
		cfg = configpkg.DefaultConfig()
	}

	// Initialize Alert Manager
	am := data.NewAlertManager()

	// Fetch static CPU info once at startup (doesn't change)
	cpuInfo, _ := cpu.Info()

	return Model{
		AppState: data.AppState{
			SelectedTab:       0,
			Config:            cfg,
			AlertManager:      am,
			SettingsSel:       configpkg.MetricCPU,
			HistoryLength:     cfg.HistoryLength,
			CpuHistory:        data.NewRingBuffer(cfg.HistoryLength),
			MemHistory:        data.NewRingBuffer(cfg.HistoryLength),
			NetHistory:        data.NewRingBuffer(cfg.HistoryLength),
			SwapHistory:       data.NewRingBuffer(cfg.HistoryLength),
			DiskHORead:        data.NewRingBuffer(cfg.HistoryLength),
			DiskHOWrite:       data.NewRingBuffer(cfg.HistoryLength),
			HistoryTemp:       data.NewRingBuffer(cfg.HistoryLength),
			Processes:         []data.ProcessInfo{},
			StartTime:         time.Now(),
			Toasts:            []data.Toast{},
			SuspendedState:    make(map[int32]bool),
			ChartType:         cfg.ChartType,
			TreeView:          cfg.ViewType == "tree",
			CollapsedPids:     make(map[int32]bool),
			SortBy:            cfg.SortBy,
			CpuInfoStatic:     cpuInfo, // Static CPU info fetched once
			Theme:             cfg.Theme,
			RefreshRate:       cfg.RefreshRate,
			BorderType:        cfg.BorderType,
			BorderStyle:       cfg.BorderStyle,
			BackgroundOpaque:  cfg.BackgroundOpaque,
			LastConfigModTime: time.Now(),
			ActiveTabs:        cfg.Tabs,
			OpenFilesView:     data.NewSimpleViewport(0, 0),
		},
	}
}

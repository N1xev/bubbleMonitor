package app

import (
	"fmt"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
	configpkg "github.com/N1xev/bubbleMonitor/internal/config"
	"github.com/N1xev/bubbleMonitor/internal/data"
	"github.com/N1xev/bubbleMonitor/internal/msg"
)

func LogMetricsCmd(cpu, memory, disk, netRate float64, procCount int, enabled bool, logPath string) tea.Cmd {
	return func() tea.Msg {
		if !enabled {
			return nil
		}

		path, err := configpkg.ResolvePath(logPath, "bubble_metrics.log")
		if err != nil {
			return msg.ToastMsg{Message: fmt.Sprintf("Path Error: %v", err), Level: data.ToastError, Duration: 3 * time.Second}
		}

		file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return msg.ToastMsg{Message: fmt.Sprintf("Failed to open log file %s: %v", path, err), Level: data.ToastError, Duration: 3 * time.Second}
		}
		defer file.Close()

		timestamp := time.Now().Format(time.RFC3339)
		line := fmt.Sprintf("%s | CPU: %.2f%% | Mem: %.2f%% | Disk: %.2f%% | Net: %.2f MB/s | Procs: %d\n",
			timestamp, cpu, memory, disk, netRate, procCount)

		if _, err := file.WriteString(line); err != nil {
			return msg.ToastMsg{Message: fmt.Sprintf("Log Write Error: %v", err), Level: data.ToastError, Duration: 3 * time.Second}
		}

		return nil
	}
}

package app

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/N1xev/bubbleMonitor/internal/data"
	"github.com/N1xev/bubbleMonitor/internal/msg"
)

func LogMetricsCmd(cpu, memory, disk, netRate float64, procCount int, enabled bool, logPath string) tea.Cmd {
	return func() tea.Msg {
		if !enabled {
			return nil
		}

		path := logPath
		if path == "" {
			home, _ := os.UserHomeDir()
			path = filepath.Join(home, "bubble_metrics.log")
		} else if !filepath.IsAbs(path) {
			home, _ := os.UserHomeDir()
			path = filepath.Join(home, path)
		}

		file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return msg.ToastMsg{Message: "Log Error: " + err.Error(), Level: data.ToastError, Duration: 3 * time.Second}
		}
		defer file.Close()

		timestamp := time.Now().Format(time.RFC3339)
		line := fmt.Sprintf("%s | CPU: %.2f%% | Mem: %.2f%% | Disk: %.2f%% | Net: %.2f MB/s | Procs: %d\n",
			timestamp, cpu, memory, disk, netRate, procCount)

		if _, err := file.WriteString(line); err != nil {
			return msg.ToastMsg{Message: "Log Write Error", Level: data.ToastError, Duration: 3 * time.Second}
		}

		return nil
	}
}

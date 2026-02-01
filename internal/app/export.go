package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/N1xev/bubbleMonitor/internal/data"
	"github.com/N1xev/bubbleMonitor/internal/msg"
)

func SaveSnapshotCmd(cpu, memory, disk float64, processCount int) tea.Cmd {
	return func() tea.Msg {
		snapshot := struct {
			Timestamp time.Time
			Cpu       float64
			Memory    float64
			Disk      float64
			Processes int
		}{
			Timestamp: time.Now(),
			Cpu:       cpu,
			Memory:    memory,
			Disk:      disk,
			Processes: processCount,
		}

		home, err := os.UserHomeDir()
		if err != nil {
			return msg.ToastMsg{Message: "Failed to get home dir", Level: data.ToastError, Duration: 3 * time.Second}
		}

		path := filepath.Join(home, "bubble_snapshot.json")
		file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return msg.ToastMsg{Message: "Failed to open file", Level: data.ToastError, Duration: 3 * time.Second}
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		if err := encoder.Encode(snapshot); err != nil {
			return msg.ToastMsg{Message: "Failed to write snapshot", Level: data.ToastError, Duration: 3 * time.Second}
		}

		csvPath := filepath.Join(home, "bubble_snapshot.csv")
		csvFile, err := os.OpenFile(csvPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			defer csvFile.Close()
			stat, _ := csvFile.Stat()
			if stat != nil && stat.Size() == 0 {
				if _, err := csvFile.WriteString("Timestamp,CPU,Memory,Disk,Processes\n"); err != nil {
					return msg.ToastMsg{Message: "Failed to write CSV header: " + err.Error(), Level: data.ToastError, Duration: 3 * time.Second}
				}
			}
			line := fmt.Sprintf("%s,%.2f,%.2f,%.2f,%d\n", snapshot.Timestamp.Format(time.RFC3339), snapshot.Cpu, snapshot.Memory, snapshot.Disk, snapshot.Processes)
			if _, err := csvFile.WriteString(line); err != nil {
				return msg.ToastMsg{Message: "Failed to write CSV data: " + err.Error(), Level: data.ToastError, Duration: 3 * time.Second}
			}
		}

		return msg.ToastMsg{Message: "Snapshot saved to " + path + " & .csv", Level: data.ToastSuccess, Duration: 3 * time.Second}
	}
}

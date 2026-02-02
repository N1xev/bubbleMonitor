package config

import (
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
)

type ConfigChangeMsg struct {
	NewModTime time.Time
}

type ConfigWatchTickMsg struct{}

func WatchConfig(lastModTime time.Time) tea.Cmd {
	return func() tea.Msg {
		path, err := GetConfigPath()
		if err != nil {
			return scheduleNextCheck()
		}
		info, err := os.Stat(path)
		if err != nil {
			return scheduleNextCheck()
		}
		currentModTime := info.ModTime()
		if currentModTime.After(lastModTime) {
			return ConfigChangeMsg{NewModTime: currentModTime}
		}
		return scheduleNextCheck()
	}
}

func scheduleNextCheck() tea.Msg {
	return tea.Tick(2*time.Second, func(time.Time) tea.Msg {
		return ConfigWatchTickMsg{}
	})()
}

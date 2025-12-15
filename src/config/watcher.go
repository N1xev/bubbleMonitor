package config

import (
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
)

// ConfigChangeMsg is sent when the configuration file changes
type ConfigChangeMsg struct {
	NewModTime time.Time
}

// WatchConfig watches the config file for changes (polling)
func WatchConfig(lastModTime time.Time) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(2 * time.Second) // Poll every 2 seconds
		path, err := GetConfigPath()
		if err != nil {
			return WatchConfig(lastModTime)()
		}
		info, err := os.Stat(path)
		if err != nil {
			return WatchConfig(lastModTime)()
		}
		currentModTime := info.ModTime()
		if currentModTime.After(lastModTime) {
			return ConfigChangeMsg{NewModTime: currentModTime}
		}
		return WatchConfig(lastModTime)()
	}
}

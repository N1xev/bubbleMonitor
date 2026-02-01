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

// ConfigWatchTickMsg is sent to trigger the next config check
type ConfigWatchTickMsg struct{}

// WatchConfig watches the config file for changes (polling)
// Returns a message immediately and schedules the next check
func WatchConfig(lastModTime time.Time) tea.Cmd {
	return func() tea.Msg {
		path, err := GetConfigPath()
		if err != nil {
			// On error, schedule next check
			return tea.Tick(2*time.Second, func(time.Time) tea.Msg {
				return ConfigWatchTickMsg{}
			})()
		}
		info, err := os.Stat(path)
		if err != nil {
			// On error, schedule next check
			return tea.Tick(2*time.Second, func(time.Time) tea.Msg {
				return ConfigWatchTickMsg{}
			})()
		}
		currentModTime := info.ModTime()
		if currentModTime.After(lastModTime) {
			return ConfigChangeMsg{NewModTime: currentModTime}
		}
		// No change, schedule next check
		return tea.Tick(2*time.Second, func(time.Time) tea.Msg {
			return ConfigWatchTickMsg{}
		})()
	}
}

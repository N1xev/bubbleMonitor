package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// MetricType defines the type of system metric
type MetricType string

const (
	MetricCPU  MetricType = "CPU"
	MetricMem  MetricType = "Memory"
	MetricDisk MetricType = "Disk"
	MetricTemp MetricType = "Temperature"
)

// AppConfig holds persistent configuration
type AppConfig struct {
	Thresholds       map[MetricType]float64 `json:"thresholds"`
	HistoryLength    int                    `json:"history_length"`
	ChartType        string                 `json:"chart_type"`
	ViewType         string                 `json:"view_type"`         // "normal" or "tree"
	SortBy           string                 `json:"sort_by"`           // "cpu", "mem", "pid"
	Theme            string                 `json:"theme"`             // dark, light, nord, dracula, custom, etc
	RefreshRate      int                    `json:"refresh_rate"`      // milliseconds: 500, 1000, 2000, 5000
	BorderType       string                 `json:"border_type"`       // normal, rounded
	BorderStyle      string                 `json:"border_style"`      // single, double, dashed
	BackgroundOpaque bool                   `json:"background_opaque"` // true = opaque, false = transparent
	Tabs             []string               `json:"tabs,omitempty"`
	CustomTheme      *CustomThemeConfig     `json:"custom_theme,omitempty"`
}

// CustomThemeConfig holds user-configurable theme colors
type CustomThemeConfig struct {
	Primary    string `json:"primary"`
	Secondary  string `json:"secondary"`
	Success    string `json:"success"`
	Warning    string `json:"warning"`
	Alert      string `json:"alert"`
	Text       string `json:"text"`
	Muted      string `json:"muted"`
	Border     string `json:"border"`
	Background string `json:"background"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() AppConfig {
	return AppConfig{
		HistoryLength:    60,
		ChartType:        "sparkline",
		ViewType:         "normal",
		SortBy:           "cpu",
		Theme:            "dark",
		RefreshRate:      1000,
		BorderType:       "rounded",
		BorderStyle:      "single",
		BackgroundOpaque: false,
		Tabs:             []string{"Overview", "Metrics", "Processes", "Disks", "Network", "System"},
		Thresholds: map[MetricType]float64{
			MetricCPU:  90.0,
			MetricMem:  90.0,
			MetricDisk: 90.0,
			MetricTemp: 85.0,
		},
	}
}

// GetConfigPath resolves the configuration file path
func GetConfigPath() (string, error) {
	config, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(config, "bubble-monitor")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.json"), nil
}

// LoadConfig loads the configuration from file
func LoadConfig() (AppConfig, error) {
	path, err := GetConfigPath()
	if err != nil {
		return DefaultConfig(), err
	}

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return DefaultConfig(), err
	}
	defer file.Close()

	var config AppConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return DefaultConfig(), err
	}

	// Ensure defaults for missing keys if any
	defaults := DefaultConfig()
	if config.Thresholds == nil {
		config.Thresholds = defaults.Thresholds
	}
	// Verify history length logic if strictly needed, but 0 is valid? No, default 60.
	if config.HistoryLength == 0 {
		config.HistoryLength = defaults.HistoryLength
	}
	if config.ChartType == "" {
		config.ChartType = defaults.ChartType
	}
	if config.ViewType == "" {
		config.ViewType = defaults.ViewType
	}
	if config.SortBy == "" {
		config.SortBy = defaults.SortBy
	}
	// If config file existed but Tabs missing/empty, populate defaults
	if len(config.Tabs) == 0 {
		config.Tabs = defaults.Tabs
		// We could save here, but maybe better to just use defaults in memory
		// User requested: "added to config by default". So we SHOULD save.
		// However, LoadConfig returns config. Saving here might be circular or side-effect heavy?
		// Better to rely on app saving on exit or change.
		// But to ensure user sees it, we must ensure it's in memory.
	}
	if config.Theme == "" {
		config.Theme = defaults.Theme
	}
	if config.RefreshRate == 0 {
		config.RefreshRate = defaults.RefreshRate
	}
	if config.BorderType == "" {
		config.BorderType = defaults.BorderType
	}
	if config.BorderStyle == "" {
		config.BorderStyle = defaults.BorderStyle
	}

	return config, nil
}

// SaveConfig writes the configuration to file
func SaveConfig(config AppConfig) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}

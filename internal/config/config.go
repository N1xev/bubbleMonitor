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

type AppConfig struct {
	Thresholds           map[MetricType]float64 `json:"thresholds"`
	HistoryLength        int                    `json:"history_length"`
	ChartType            string                 `json:"chart_type"`
	ViewType             string                 `json:"view_type"`         // "normal" or "tree"
	SortBy               string                 `json:"sort_by"`           // "cpu", "mem", "pid"
	Theme                string                 `json:"theme"`             // dark, light, nord, dracula, custom, etc
	RefreshRate          int                    `json:"refresh_rate"`      // milliseconds: 500, 1000, 2000, 5000
	BorderType           string                 `json:"border_type"`       // normal, rounded
	BorderStyle          string                 `json:"border_style"`      // single, double, dashed
	BackgroundOpaque     bool                   `json:"background_opaque"` // true = opaque, false = transparent
	ProcessCpuNormalized bool                   `json:"process_cpu_normalized"`
	Tabs                 []string               `json:"tabs,omitempty"`
	Logging              LoggingConfig          `json:"logging"`
	RemoteHosts          []RemoteHostConfig     `json:"remote_hosts"`
	HealthWeights        HealthWeights          `json:"health_weights"`
	CustomTheme          *CustomThemeConfig     `json:"custom_theme,omitempty"`
}

type HealthWeights struct {
	CpuCritical  int `json:"cpu_critical"`
	CpuHigh      int `json:"cpu_high"`
	MemCritical  int `json:"mem_critical"`
	MemHigh      int `json:"mem_high"`
	DiskCritical int `json:"disk_critical"`
	TempCritical int `json:"temp_critical"`
	TempHigh     int `json:"temp_high"`
}

type RemoteHostConfig struct {
	Name string `json:"name"`
	Host string `json:"host"` // user@hostname:port
}

type LoggingConfig struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`
}

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

func DefaultConfig() AppConfig {
	return AppConfig{
		HistoryLength:        60,
		ChartType:            "braille",
		ViewType:             "normal",
		SortBy:               "cpu",
		Theme:                "dark",
		RefreshRate:          1000,
		BorderType:           "rounded",
		BorderStyle:          "dashed",
		BackgroundOpaque:     true,
		ProcessCpuNormalized: true,
		Tabs:                 []string{"Overview", "Metrics", "Processes", "Disks", "Network", "System", "Services", "Connections", "Logs", "Remote"},
		Thresholds: map[MetricType]float64{
			MetricCPU:  90.0,
			MetricMem:  90.0,
			MetricDisk: 90.0,
			MetricTemp: 85.0,
		},
		HealthWeights: HealthWeights{
			CpuCritical:  30,
			CpuHigh:      10,
			MemCritical:  30,
			MemHigh:      10,
			DiskCritical: 20,
			TempCritical: 30,
			TempHigh:     10,
		},
	}
}

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

	// Ensure defaults for HealthWeights if missing (zero value check)
	if config.HealthWeights.CpuCritical == 0 {
		config.HealthWeights = defaults.HealthWeights
	}

	return config, nil
}

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

func ResolvePath(path string, defaultName string) (string, error) {
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, defaultName), nil
	}
	if !filepath.IsAbs(path) {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, path), nil
	}
	return path, nil
}

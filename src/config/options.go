package config

// GetThemeNames returns available theme names
func GetThemeNames() []string {
	return []string{
		"dark",
		"light",
		"nord",
		"dracula",
		"gruvbox",
		"solarized",
		"monokai",
		"catppuccin",
		"tokyonight",
		"onedark",
		"ayu",
		"tty",
		"rosepine",
		"everforest",
		"nightowl",
		"palenight",
		"material",
		"synthwave",
		"cobalt2",
		"horizon",
		"oceanic",
		"palefire",
		"github",
		"moonlight",
		"shades",
		"midnight",
		"forest",
		"autumn",
		"cyberpunk",
		"sunset",
		"ocean",
		"coffee",
		"custom",
	}
}

// GetBorderStyles returns available border styles
func GetBorderStyles() []string {
	return []string{"single", "double", "dashed"}
}

// GetBorderTypes returns available border types
func GetBorderTypes() []string {
	return []string{"normal", "rounded"}
}

// GetRefreshRates returns available refresh rates in ms
func GetRefreshRates() []int {
	return []int{500, 1000, 2000, 5000}
}

// DefaultCustomTheme returns a default custom theme template
func DefaultCustomTheme() *CustomThemeConfig {
	return &CustomThemeConfig{
		Primary:    "#7D56F4",
		Secondary:  "#EE6FF8",
		Success:    "#A1E3AD",
		Warning:    "#F5A962",
		Alert:      "#F25D94",
		Text:       "#F0F0F0",
		Muted:      "#A0A0A0",
		Border:     "#4A4A4A",
		Background: "#1C1C1C",
	}
}

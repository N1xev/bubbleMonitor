package widgets

import (
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
)

// GetColorForValue returns a color based on the value threshold
func GetColorForValue(val float64, su, w, a compat.AdaptiveColor) compat.AdaptiveColor {
	if val < 50 {
		return su
	} else if val < 80 {
		return w
	}
	return a
}

// Pre-allocate the empty style as it's constant
var emptyStyle = lipgloss.NewStyle().Foreground(compat.AdaptiveColor{Light: lipgloss.Color("#E5E7EB"), Dark: lipgloss.Color("#374151")})

// RenderProgressBar creates a colored progress bar
func RenderProgressBar(val float64, width int, su, w, a compat.AdaptiveColor) string {
	if val < 0 {
		val = 0
	}
	if val > 100 {
		val = 100
	}
	filled := int(val / 100 * float64(width))
	if filled > width {
		filled = width
	}
	empty := width - filled
	color := GetColorForValue(val, su, w, a)

	// Create filled style on demand (color varies), but reuse emptyStyle
	filledStyle := lipgloss.NewStyle().Foreground(color)

	return filledStyle.Render(strings.Repeat("█", filled)) + emptyStyle.Render(strings.Repeat("░", empty))
}

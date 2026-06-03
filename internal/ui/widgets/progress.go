package widgets

import (
	"image/color"

	"charm.land/bubbles/v2/progress"
)

// Pre-configured progress model with exact same characters as original implementation
// Using '█' (full block) for filled and '░' (light shade) for empty
var progressBar = progress.New(
	progress.WithFillCharacters('█', '░'),
	progress.WithoutPercentage(),
)

// GetColorForValue returns a color based on the value threshold
func GetColorForValue(val float64, su, w, a color.Color) color.Color {
	if val < ProgressBarWarningThreshold {
		return su
	} else if val < ProgressBarCriticalThreshold {
		return w
	}
	return a
}

// RenderProgressBar creates a colored progress bar using Bubbles progress.Model
// Maintains exact same visual output as original implementation
func RenderProgressBar(val float64, width int, su, w, a color.Color) string {
	if val < 0 {
		val = 0
	}
	if val > 100 {
		val = 100
	}

	// Set width dynamically
	progressBar.SetWidth(width)

	// Set colors based on value thresholds
	c := GetColorForValue(val, su, w, a)
	progressBar.FullColor = c
	progressBar.EmptyColor = c

	// ViewAs renders at the given percentage without needing Update cycle
	// Dividing by 100 because ViewAs expects 0.0-1.0 range
	return progressBar.ViewAs(val / 100)
}

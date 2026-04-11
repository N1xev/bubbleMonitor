package cliout

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// BarColor returns a lipgloss.Style colored by the percentage threshold:
// <60% -> Success (green), <80% -> Warning (yellow), >=80% -> Alert (red).
func (s CLIStyles) BarColor(pct float64) lipgloss.Style {
	switch {
	case pct >= 80:
		return lipgloss.NewStyle().Foreground(s.palette.Alert)
	case pct >= 60:
		return lipgloss.NewStyle().Foreground(s.palette.Warning)
	default:
		return lipgloss.NewStyle().Foreground(s.palette.Success)
	}
}

// ScoreColor returns a lipgloss.Style colored by the health score threshold:
// >=70 -> Success (green), >=50 -> Warning (yellow), <50 -> Alert (red).
func (s CLIStyles) ScoreColor(score int) lipgloss.Style {
	return lipgloss.NewStyle() // stub: will be filled in Task 2
}

// RenderBar produces a colored progress bar string.
func (s CLIStyles) RenderBar(pct float64, width int) string {
	return strings.Repeat("\u2591", width) // stub: will be filled in Task 2
}

// VisualPad pads a styled string to targetWidth, accounting for ANSI escape codes.
func VisualPad(styled string, targetWidth int) string {
	return styled // stub: will be filled in Task 2
}

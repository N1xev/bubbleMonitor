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
	switch {
	case score >= 70:
		return lipgloss.NewStyle().Foreground(s.palette.Success)
	case score >= 50:
		return lipgloss.NewStyle().Foreground(s.palette.Warning)
	default:
		return lipgloss.NewStyle().Foreground(s.palette.Alert)
	}
}

// RenderBar produces a colored progress bar string.
// Filled blocks are colored by BarColor threshold; empty blocks use the Dim style.
func (s CLIStyles) RenderBar(pct float64, width int) string {
	filled := int(pct / 100.0 * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	barStyle := s.BarColor(pct)
	filledPart := barStyle.Render(strings.Repeat("\u2588", filled))
	emptyPart := s.Dim.Render(strings.Repeat("\u2591", width-filled))
	return filledPart + emptyPart
}

// VisualPad pads a styled string to targetWidth, accounting for ANSI escape codes.
// Uses lipgloss.Width() to measure the visible cell width, then appends spaces
// to reach the target. If the string is already at or beyond targetWidth, it is
// returned unchanged.
func VisualPad(styled string, targetWidth int) string {
	visibleWidth := lipgloss.Width(styled)
	if visibleWidth >= targetWidth {
		return styled
	}
	return styled + strings.Repeat(" ", targetWidth-visibleWidth)
}

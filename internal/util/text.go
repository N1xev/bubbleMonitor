package util

import (
	"strings"

	"charm.land/lipgloss/v2"
)

var spaceCache = make([]string, 128)

func init() {
	for i := range spaceCache {
		spaceCache[i] = strings.Repeat(" ", i)
	}
}

func getSpaces(n int) string {
	if n < len(spaceCache) {
		return spaceCache[n]
	}
	return strings.Repeat(" ", n)
}

// PadRight pads string s with spaces on the right to visual width w
func PadRight(s string, w int) string {
	width := lipgloss.Width(s)
	if width >= w {
		return s
	}
	return s + getSpaces(w-width)
}

// PadLeft pads string s with spaces on the left to visual width w
func PadLeft(s string, w int) string {
	width := lipgloss.Width(s)
	if width >= w {
		return s
	}
	return getSpaces(w-width) + s
}

// Truncate truncates string s to visual width w
func Truncate(s string, w int) string {
	return lipgloss.NewStyle().MaxWidth(w).Render(s)
}

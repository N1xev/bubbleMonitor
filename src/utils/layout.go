package utils

import (
	"charm.land/lipgloss/v2"
)

// CalculateColumnWidths distributes the total width among n columns
func CalculateColumnWidths(totalWidth, n int) []int {
	if n <= 0 {
		return []int{}
	}
	baseWidth := totalWidth / n
	remainder := totalWidth % n
	widths := make([]int, n)
	for i := 0; i < n; i++ {
		widths[i] = baseWidth
		if i < remainder {
			widths[i]++
		}
	}
	return widths
}

// FullWidthBg fills a string to the given width with spaces using lipgloss
func FullWidthBg(s string, width int) string {
	return lipgloss.NewStyle().Width(width).Render(s)
}

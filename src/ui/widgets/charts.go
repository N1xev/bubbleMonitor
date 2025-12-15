package widgets

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"

	"github.com/N1xev/bubbleMonitor/src/data"
)

// Helper to get max value from Accessor
func maxFromAccessor(a data.Accessor) float64 {
	if a.Len() == 0 {
		return 0
	}
	max := 0.0
	for i := 0; i < a.Len(); i++ {
		v := a.Get(i)
		if v > max {
			max = v
		}
	}
	return max
}

// RenderSparkline creates a sparkline chart from data
func RenderSparkline(data data.Accessor, width, height int, c1, c2 compat.AdaptiveColor) string {
	if data.Len() == 0 {
		return "No data"
	}
	chars := []string{" ", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	maxV := maxFromAccessor(data)
	if maxV == 0 {
		maxV = 1
	}
	var result strings.Builder
	startIdx := 0
	if data.Len() > width {
		startIdx = data.Len() - width
	}
	for i := startIdx; i < data.Len(); i++ {
		val := data.Get(i)
		normVal := val / maxV
		chIdx := int(normVal * float64(len(chars)-1))
		if chIdx >= len(chars) {
			chIdx = len(chars) - 1
		}
		if chIdx < 0 {
			chIdx = 0
		}
		color := c1
		if val > 70 {
			color = c2
		}
		result.WriteString(lipgloss.NewStyle().Foreground(color).Render(chars[chIdx]))
	}

	for i := data.Len(); i < width; i++ {
		result.WriteString(" ")
	}
	return result.String()
}

// RenderLineChart creates a multi-line chart using block characters
func RenderLineChart(data data.Accessor, width, height int, c1, c2 compat.AdaptiveColor) string {
	if data.Len() == 0 || height < 1 {
		return "No data"
	}
	maxV := maxFromAccessor(data)
	if maxV == 0 {
		maxV = 1
	}

	grid := make([][]string, height)
	for r := 0; r < height; r++ {
		grid[r] = make([]string, width)
		for c := 0; c < width; c++ {
			grid[r][c] = " "
		}
	}

	startIdx := 0
	if data.Len() > width {
		startIdx = data.Len() - width
	}

	for col := 0; col < width && (startIdx+col) < data.Len(); col++ {
		val := data.Get(startIdx + col)
		normalized := val / maxV
		filledRows := int(normalized * float64(height))

		for row := height - 1; row >= height-filledRows && row >= 0; row-- {
			grid[row][col] = "█"
		}
	}

	// Convert grid to string
	var lines []string
	for r := 0; r < height; r++ {
		line := ""
		for c := 0; c < width; c++ {
			char := grid[r][c]
			color := c1
			if char == "█" {
				color = c2
			}
			line += lipgloss.NewStyle().Foreground(color).Render(char)
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// RenderBarChart creates a horizontal bar chart
func RenderBarChart(data data.Accessor, width, height int, c1, c2 compat.AdaptiveColor) string {
	if data.Len() == 0 {
		return "No data"
	}
	maxV := maxFromAccessor(data)
	if maxV == 0 {
		maxV = 1
	}

	// Show latest N values that fit in height
	startIdx := 0
	if data.Len() > height {
		startIdx = data.Len() - height
	}

	var lines []string
	for i := startIdx; i < data.Len(); i++ {
		val := data.Get(i)
		barLen := int((val / maxV) * float64(width-8))
		if barLen < 0 {
			barLen = 0
		}

		color := c1
		if val > 70 {
			color = c2
		}

		label := fmt.Sprintf("%5.1f%% ", val)
		bar := lipgloss.NewStyle().Foreground(color).Render(strings.Repeat("█", barLen))
		lines = append(lines, label+bar)
	}

	return strings.Join(lines, "\n")
}

// RenderBrailleChart creates a high-resolution chart using Unicode Braille patterns
func RenderBrailleChart(data data.Accessor, width, height int, c1, c2 compat.AdaptiveColor) string {
	if data.Len() == 0 || height < 1 {
		return "No data"
	}
	maxV := maxFromAccessor(data)
	if maxV == 0 {
		maxV = 1
	}

	// Braille has 4 rows per character
	dotsPerCol := height * 4

	sampleWidth := width * 2
	startIdx := 0
	if data.Len() > sampleWidth {
		startIdx = data.Len() - sampleWidth
	}

	// Create dot matrix: dotsPerCol rows x sampleWidth cols
	dots := make([][]bool, dotsPerCol)
	for r := 0; r < dotsPerCol; r++ {
		dots[r] = make([]bool, sampleWidth)
	}

	// Fill dots based on data
	for col := 0; col < sampleWidth && (startIdx+col) < data.Len(); col++ {
		val := data.Get(startIdx + col)
		normalized := val / maxV
		filledDots := int(normalized * float64(dotsPerCol))

		for row := dotsPerCol - 1; row >= dotsPerCol-filledDots && row >= 0; row-- {
			dots[row][col] = true
		}
	}

	var lines []string
	for charRow := 0; charRow < height; charRow++ {
		var line strings.Builder
		for charCol := 0; charCol < width; charCol++ {
			dotRow := charRow * 4
			dotCol := charCol * 2

			var braille rune = 0x2800
			if dotCol < sampleWidth {
				if dotRow < dotsPerCol && dots[dotRow][dotCol] {
					braille += 1
				}
				if dotRow+1 < dotsPerCol && dots[dotRow+1][dotCol] {
					braille += 2
				}
				if dotRow+2 < dotsPerCol && dots[dotRow+2][dotCol] {
					braille += 4
				}
				if dotRow+3 < dotsPerCol && dots[dotRow+3][dotCol] {
					braille += 64
				}
			}
			if dotCol+1 < sampleWidth {
				if dotRow < dotsPerCol && dots[dotRow][dotCol+1] {
					braille += 8
				}
				if dotRow+1 < dotsPerCol && dots[dotRow+1][dotCol+1] {
					braille += 16
				}
				if dotRow+2 < dotsPerCol && dots[dotRow+2][dotCol+1] {
					braille += 32
				}
				if dotRow+3 < dotsPerCol && dots[dotRow+3][dotCol+1] {
					braille += 128
				}
			}
			line.WriteRune(braille)
		}
		color := c1
		if charRow < height/2 {
			color = c2
		}
		lines = append(lines, lipgloss.NewStyle().Foreground(color).Render(line.String()))
	}

	return strings.Join(lines, "\n")
}

// RenderTTYChart creates a simple chart using block characters for TTY compatibility
func RenderTTYChart(data data.Accessor, width, height int, c1, c2 compat.AdaptiveColor) string {
	if data.Len() == 0 || height < 1 {
		return "No data"
	}
	maxV := maxFromAccessor(data)
	if maxV == 0 {
		maxV = 1
	}

	startIdx := 0
	if data.Len() > width {
		startIdx = data.Len() - width
	}

	grid := make([][]rune, height)
	for r := 0; r < height; r++ {
		grid[r] = make([]rune, width)
		for c := 0; c < width; c++ {
			grid[r][c] = ' '
		}
	}

	for col := 0; col < width && (startIdx+col) < data.Len(); col++ {
		val := data.Get(startIdx + col)
		normalized := val / maxV
		filledRows := int(normalized * float64(height))

		for row := height - 1; row >= height-filledRows && row >= 0; row-- {
			grid[row][col] = '█'
		}
	}

	var lines []string
	for r := 0; r < height; r++ {
		var line strings.Builder
		for c := 0; c < width; c++ {
			char := grid[r][c]
			color := c1
			if r < height/3 {
				color = c2
			}
			line.WriteString(lipgloss.NewStyle().Foreground(color).Render(string(char)))
		}
		lines = append(lines, line.String())
	}

	return strings.Join(lines, "\n")
}

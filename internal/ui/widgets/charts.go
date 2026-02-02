package widgets

import (
	"fmt"
	"strings"
	"sync"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"

	"github.com/N1xev/bubbleMonitor/internal/data"
)

var (
	stringGridPool = sync.Pool{
		New: func() interface{} {
			grid := make([][]string, 0, InitialGridCapacity)
			return &grid
		},
	}
	boolGridPool = sync.Pool{
		New: func() interface{} {
			grid := make([][]bool, 0, 20)
			return &grid
		},
	}
)

func RenderSparkline(data data.Accessor, width, height int, c1, c2 compat.AdaptiveColor, fixedMax float64) string {
	if data.Len() == 0 {
		return "No data"
	}
	chars := []string{" ", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	maxV := fixedMax
	if maxV <= 0 {
		maxV = data.Max()
	}
	if maxV == 0 {
		maxV = 1
	}
	var result strings.Builder
	result.Grow(width * 4) // Pre-allocate for UTF-8 chars
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

func RenderLineChart(data data.Accessor, width, height int, c1, c2 compat.AdaptiveColor, fixedMax float64) string {
	if data.Len() == 0 || height < 1 {
		return "No data"
	}
	maxV := fixedMax
	if maxV <= 0 {
		maxV = data.Max()
	}
	if maxV == 0 {
		maxV = 1
	}

	gridPtr := stringGridPool.Get().(*[][]string)
	grid := *gridPtr

	if cap(grid) < height {
		grid = make([][]string, height)
	} else {
		grid = grid[:height]
	}

	for r := 0; r < height; r++ {
		if cap(grid[r]) < width {
			grid[r] = make([]string, width)
		} else {
			grid[r] = grid[r][:width]
		}
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

	lines := make([]string, 0, height)
	for r := 0; r < height; r++ {
		var line strings.Builder
		line.Grow(width * 10) // Estimate for styled chars
		for c := 0; c < width; c++ {
			char := grid[r][c]
			color := c1
			if char == "█" {
				color = c2
			}
			line.WriteString(lipgloss.NewStyle().Foreground(color).Render(char))
		}
		lines = append(lines, line.String())
	}

	*gridPtr = grid
	stringGridPool.Put(gridPtr)

	return strings.Join(lines, "\n")
}

func RenderBarChart(data data.Accessor, width, height int, c1, c2 compat.AdaptiveColor, fixedMax float64) string {
	if data.Len() == 0 {
		return "No data"
	}
	maxV := fixedMax
	if maxV <= 0 {
		maxV = data.Max()
	}
	if maxV == 0 {
		maxV = 1
	}

	startIdx := 0
	count := data.Len() - startIdx
	if data.Len() > height {
		startIdx = data.Len() - height
		count = height
	}

	lines := make([]string, 0, count)
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

func RenderBrailleChart(data data.Accessor, width, height int, c1, c2 compat.AdaptiveColor, fixedMax float64) string {
	if data.Len() == 0 || height < 1 {
		return "No data"
	}
	maxV := fixedMax
	if maxV <= 0 {
		maxV = data.Max()
	}
	if maxV == 0 {
		maxV = 1
	}

	dotsPerCol := height * 4

	sampleWidth := width * 2
	startIdx := 0
	if data.Len() > sampleWidth {
		startIdx = data.Len() - sampleWidth
	}

	gridPtr := boolGridPool.Get().(*[][]bool)
	dots := *gridPtr

	if cap(dots) < dotsPerCol {
		dots = make([][]bool, dotsPerCol)
	} else {
		dots = dots[:dotsPerCol]
	}

	for r := 0; r < dotsPerCol; r++ {
		if cap(dots[r]) < sampleWidth {
			dots[r] = make([]bool, sampleWidth)
		} else {
			dots[r] = dots[r][:sampleWidth]
			for c := range dots[r] {
				dots[r][c] = false
			}
		}
	}

	for col := 0; col < sampleWidth && (startIdx+col) < data.Len(); col++ {
		val := data.Get(startIdx + col)
		normalized := val / maxV
		filledDots := int(normalized * float64(dotsPerCol))

		for row := dotsPerCol - 1; row >= dotsPerCol-filledDots && row >= 0; row-- {
			dots[row][col] = true
		}
	}

	lines := make([]string, 0, height)
	for charRow := 0; charRow < height; charRow++ {
		var line strings.Builder
		line.Grow(width)
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

	*gridPtr = dots
	boolGridPool.Put(gridPtr)

	return strings.Join(lines, "\n")
}

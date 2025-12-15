package widgets

import (
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
)

// GetBorder returns the lipgloss border for the current style
func GetBorder(borderStyle, borderType string) lipgloss.Border {
	switch borderStyle {
	case "double":
		return lipgloss.DoubleBorder()
	case "dashed":
		return lipgloss.Border{
			Top:         "-",
			Bottom:      "-",
			Left:        "|",
			Right:       "|",
			TopLeft:     "+",
			TopRight:    "+",
			BottomLeft:  "+",
			BottomRight: "+",
		}
	default:
		if borderType == "rounded" {
			return lipgloss.RoundedBorder()
		}
		return lipgloss.NormalBorder()
	}
}

// RenderTopBorderWithBg renders a top border with title
func RenderTopBorderWithBg(title string, width int, border lipgloss.Border, borderColor, titleColor compat.AdaptiveColor) string {
	borderStyle := lipgloss.NewStyle().Foreground(borderColor)
	titleStyle := lipgloss.NewStyle().Foreground(titleColor).Bold(true)

	leftPart := borderStyle.Render(border.TopLeft+border.Top+" ") +
		titleStyle.Render(title) +
		borderStyle.Render(" ")

	rightPart := borderStyle.Render(border.TopRight)

	leftWidth := lipgloss.Width(leftPart)
	rightWidth := lipgloss.Width(rightPart)

	remainingWidth := width - leftWidth - rightWidth
	if remainingWidth < 0 {
		remainingWidth = 0
	}

	hLine := border.Top
	hLineWidth := lipgloss.Width(hLine)
	if hLineWidth == 0 {
		hLineWidth = 1
	}

	repeatCount := remainingWidth / hLineWidth
	middlePart := borderStyle.Render(strings.Repeat(hLine, repeatCount))

	return leftPart + middlePart + rightPart
}

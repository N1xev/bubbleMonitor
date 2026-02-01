package tabs

import (
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"

	"github.com/N1xev/bubbleMonitor/internal/data"
	"github.com/N1xev/bubbleMonitor/internal/ui/widgets"
)

func RenderLogs(s *data.AppState, container lipgloss.Style, su, w, a, t, mu, p, b compat.AdaptiveColor, availHeight int) string {
	boxWidth := s.Width
	border := widgets.GetBorder(s.BorderStyle, s.BorderType)

	contentWidth := boxWidth - 4

	visibleRows := availHeight - 2
	if visibleRows < 1 {
		visibleRows = 1
	}

	var logLines []string
	logStyle := lipgloss.NewStyle().Foreground(t)

	// Show last N logs that fit
	count := len(s.SystemLogs)
	start := 0
	if count > visibleRows {
		start = count - visibleRows
	}

	for i := start; i < count; i++ {
		line := s.SystemLogs[i]
		// Basic highlighting
		if strings.Contains(strings.ToLower(line), "error") || strings.Contains(strings.ToLower(line), "fail") {
			logLines = append(logLines, lipgloss.NewStyle().Foreground(a).Width(contentWidth).Render(line))
		} else if strings.Contains(strings.ToLower(line), "warn") {
			logLines = append(logLines, lipgloss.NewStyle().Foreground(w).Width(contentWidth).Render(line))
		} else {
			logLines = append(logLines, logStyle.Width(contentWidth).Render(line))
		}
	}

	content := strings.Join(logLines, "\n")

	c := container.Width(boxWidth).Height(visibleRows).BorderTop(false)
	body := c.Render(content)
	topBorder := widgets.RenderTopBorderWithBg("SYSTEM LOGS (Last 50)", boxWidth, border, b, p)

	return lipgloss.JoinVertical(lipgloss.Left, topBorder, body)
}

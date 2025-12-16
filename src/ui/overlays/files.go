package overlays

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"

	"github.com/N1xev/bubbleMonitor/src/data"
	"github.com/N1xev/bubbleMonitor/src/ui/widgets"
)

// RenderOpenFilesOverlay renders the open files list
func RenderOpenFilesOverlay(s *data.AppState, width, height int, b, p, t, mu, bg compat.AdaptiveColor) string {
	boxWidth := 80
	if boxWidth > width-4 {
		boxWidth = width - 4
	}
	boxHeight := 20
	if boxHeight > height-4 {
		boxHeight = height - 4
	}

	border := widgets.GetBorder(s.BorderStyle, s.BorderType)
	title := fmt.Sprintf("OPEN FILES (PID %d)", s.OpenFilesPid)

	vpWidth := boxWidth - 10
	vpHeight := boxHeight - 4
	if vpWidth < 10 {
		vpWidth = 10
	}
	if vpHeight < 5 {
		vpHeight = 5
	}

	s.OpenFilesView.Width = vpWidth
	s.OpenFilesView.Height = vpHeight

	container := lipgloss.NewStyle().
		Border(border).
		BorderForeground(b).
		Padding(1, 2).
		Width(boxWidth - 6).
		Height(boxHeight).
		BorderTop(false)

	hint := lipgloss.NewStyle().Foreground(mu).Italic(true).Render("↑↓/jk to scroll • O or ESC to close")

	// Render Viewport
	renderedView := s.OpenFilesView.View()

	contentWithHint := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Foreground(t).Render(renderedView),
		"",
		hint,
	)

	body := container.Render(contentWithHint)
	actualWidth := lipgloss.Width(body)

	topBorder := widgets.RenderTopBorderWithBg(title, actualWidth, border, b, p)

	return lipgloss.JoinVertical(lipgloss.Left, topBorder, body)
}

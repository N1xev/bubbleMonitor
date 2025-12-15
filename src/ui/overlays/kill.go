package overlays

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
	"github.com/N1xev/bubbleMonitor/src/data"
	"github.com/N1xev/bubbleMonitor/src/ui/widgets"
)

// RenderKillDialog renders the kill confirmation dialog using AppState
func RenderKillDialog(s *data.AppState, b, p, danger, t, mu compat.AdaptiveColor) string {
	boxWidth := 50
	if boxWidth > s.Width-4 {
		boxWidth = s.Width - 4
	}

	warningStyle := lipgloss.NewStyle().Foreground(danger).Bold(true)
	labelStyle := lipgloss.NewStyle().Foreground(mu)
	valueStyle := lipgloss.NewStyle().Foreground(t).Bold(true)
	keyStyle := lipgloss.NewStyle().Foreground(p).Bold(true)

	content := lipgloss.JoinVertical(lipgloss.Center,
		warningStyle.Render("⚠ KILL PROCESS?"),
		"",
		labelStyle.Render("PID: ")+valueStyle.Render(fmt.Sprintf("%d", s.KillTargetPid)),
		labelStyle.Render("Name: ")+valueStyle.Render(s.KillTargetName),
		"",
		lipgloss.JoinHorizontal(lipgloss.Center,
			keyStyle.Render("[Y]")+" "+labelStyle.Render("Kill"),
			"   ",
			keyStyle.Render("[N]")+" "+labelStyle.Render("Cancel"),
		),
	)

	border := widgets.GetBorder(s.BorderStyle, s.BorderType)

	container := lipgloss.NewStyle().
		Border(border).
		BorderForeground(danger).
		Padding(1, 3).
		Width(boxWidth - 6).
		BorderTop(false)

	body := container.Render(content)
	actualWidth := lipgloss.Width(body)
	topBorder := widgets.RenderTopBorderWithBg("CONFIRM KILL", actualWidth, border, danger, p)

	return lipgloss.JoinVertical(lipgloss.Left, topBorder, body)
}

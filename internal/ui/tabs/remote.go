package tabs

import (
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"

	"github.com/N1xev/bubbleMonitor/internal/data"
	"github.com/N1xev/bubbleMonitor/internal/ui/widgets"
)

func RenderRemote(s *data.AppState, container lipgloss.Style, su, w, a, t, mu, p, b compat.AdaptiveColor, availHeight int) string {
	boxWidth := s.Width
	border := widgets.GetBorder(s.BorderStyle, s.BorderType)

	contentWidth := boxWidth - 4

	visibleRows := availHeight - 2
	if visibleRows < 1 {
		visibleRows = 1
	}

	var rows []string
	cellStyle := lipgloss.NewStyle().Foreground(t)

	hdrStyle := lipgloss.NewStyle().Bold(true).Underline(true)
	header := hdrStyle.Width(20).Render("HOST") + " " + hdrStyle.Width(contentWidth-21).Render("STATUS")

	for _, hostCfg := range s.Config.RemoteHosts {
		status := s.RemoteUptimes[hostCfg.Host]
		if status == "" {
			status = "Waiting..."
		}

		name := hostCfg.Name
		if name == "" {
			name = hostCfg.Host
		}

		row := lipgloss.JoinHorizontal(lipgloss.Top,
			cellStyle.Width(20).Render(name),
			" ",
			cellStyle.Width(contentWidth-21).Render(status),
		)
		rows = append(rows, row)
	}

	if len(s.Config.RemoteHosts) == 0 {
		rows = append(rows, "No remote hosts configured in config.json")
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Width(contentWidth).Render(header),
		"",
		strings.Join(rows, "\n"),
	)

	c := container.Width(boxWidth).Height(visibleRows).BorderTop(false)
	body := c.Render(content)
	topBorder := widgets.RenderTopBorderWithBg("REMOTE HOSTS", boxWidth, border, b, p)

	return lipgloss.JoinVertical(lipgloss.Left, topBorder, body)
}

package tabs

import (
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"

	"github.com/N1xev/bubbleMonitor/internal/data"
	"github.com/N1xev/bubbleMonitor/internal/ui/widgets"
)

func RenderServices(s *data.AppState, container lipgloss.Style, su, w, a, t, mu, p, b compat.AdaptiveColor, availHeight int) string {
	boxWidth := s.Width
	border := widgets.GetBorder(s.BorderStyle, s.BorderType)

	contentWidth := boxWidth - 4
	nameWidth := 40
	statusWidth := 15
	descWidth := contentWidth - nameWidth - statusWidth - 4
	if descWidth < 20 {
		descWidth = 20
	}

	sp := func(str string) string { return str }

	hdrStyle := lipgloss.NewStyle().Bold(true).Underline(true)
	headerRow := hdrStyle.Width(nameWidth).Render("UNIT") + sp(" ") +
		hdrStyle.Width(statusWidth).Render("STATUS") + sp(" ") +
		hdrStyle.Width(descWidth).Render("DESCRIPTION")

	visibleRows := availHeight - 2
	if visibleRows < 1 {
		visibleRows = 1
	}

	var rows []string
	cellStyle := lipgloss.NewStyle().Foreground(t)
	runningStyle := lipgloss.NewStyle().Foreground(su)
	stoppedStyle := lipgloss.NewStyle().Foreground(mu)
	failedStyle := lipgloss.NewStyle().Foreground(a)

	for i, svc := range s.Services {
		if i >= visibleRows {
			break
		}

		stStyle := cellStyle
		if svc.Status == "running" || svc.Status == "active" {
			stStyle = runningStyle
		} else if svc.Status == "failed" {
			stStyle = failedStyle
		} else {
			stStyle = stoppedStyle
		}

		nameCell := cellStyle.Width(nameWidth).Render(svc.Name)
		statusCell := stStyle.Width(statusWidth).Render(svc.Status)
		descCell := cellStyle.Width(descWidth).Render(svc.Description)

		row := lipgloss.JoinHorizontal(lipgloss.Top, nameCell, " ", statusCell, " ", descCell)
		rows = append(rows, row)
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Width(contentWidth).Render(headerRow),
		"",
		strings.Join(rows, "\n"),
	)

	c := container.Width(boxWidth).Height(visibleRows).BorderTop(false)
	body := c.Render(content)
	topBorder := widgets.RenderTopBorderWithBg("SYSTEM SERVICES", boxWidth, border, b, p)

	return lipgloss.JoinVertical(lipgloss.Left, topBorder, body)
}

package tabs

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"

	"github.com/N1xev/bubbleMonitor/internal/data"
	"github.com/N1xev/bubbleMonitor/internal/ui/widgets"
)

func RenderConnections(s *data.AppState, container lipgloss.Style, su, w, a, t, mu, p, b compat.AdaptiveColor, availHeight int) string {
	boxWidth := s.Width
	border := widgets.GetBorder(s.BorderStyle, s.BorderType)

	contentWidth := boxWidth - 4

	protoWidth := 6
	localWidth := 25
	remoteWidth := 25
	stateWidth := 15
	pidWidth := 8

	hdrStyle := lipgloss.NewStyle().Bold(true).Underline(true)
	sp := func(str string) string { return str }

	headerRow := hdrStyle.Width(protoWidth).Render("PROTO") + sp(" ") +
		hdrStyle.Width(localWidth).Render("LOCAL ADDR") + sp(" ") +
		hdrStyle.Width(remoteWidth).Render("REMOTE ADDR") + sp(" ") +
		hdrStyle.Width(stateWidth).Render("STATE") + sp(" ") +
		hdrStyle.Width(pidWidth).Render("PID")

	visibleRows := availHeight - 2
	if visibleRows < 1 {
		visibleRows = 1
	}

	var rows []string
	cellStyle := lipgloss.NewStyle().Foreground(t)

	for i, conn := range s.Connections {
		if i >= visibleRows {
			break
		}

		protoCell := cellStyle.Width(protoWidth).Render(conn.Protocol)
		localCell := cellStyle.Width(localWidth).Render(conn.LocalAddr)
		remoteCell := cellStyle.Width(remoteWidth).Render(conn.RemoteAddr)
		stateCell := cellStyle.Width(stateWidth).Render(conn.State)
		pidCell := cellStyle.Width(pidWidth).Render(fmt.Sprintf("%d", conn.Pid))

		row := lipgloss.JoinHorizontal(lipgloss.Top, protoCell, " ", localCell, " ", remoteCell, " ", stateCell, " ", pidCell)
		rows = append(rows, row)
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Width(contentWidth).Render(headerRow),
		"",
		strings.Join(rows, "\n"),
	)

	c := container.Width(boxWidth).Height(visibleRows).BorderTop(false)
	body := c.Render(content)
	topBorder := widgets.RenderTopBorderWithBg("NETWORK CONNECTIONS", boxWidth, border, b, p)

	return lipgloss.JoinVertical(lipgloss.Left, topBorder, body)
}

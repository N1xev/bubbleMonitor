package tabs

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"

	"github.com/N1xev/bubbleMonitor/src/data"
	"github.com/N1xev/bubbleMonitor/src/ui/widgets"
	"github.com/N1xev/bubbleMonitor/src/utils"
)

// RenderNetwork renders the network interfaces tab
func RenderNetwork(s *data.AppState, container, titleStyle, labelStyle, valueStyle lipgloss.Style, t, mu, p, b, bg compat.AdaptiveColor, availHeight int) string {
	if len(s.NetworkInterfaces) == 0 {
		return "Loading network interfaces..."
	}

	width := s.Width
	cols := 1
	if width >= 80 {
		cols = 2
	}
	if width >= 100 {
		cols = 3
	}

	colWidths := utils.CalculateColumnWidths(width, cols)
	border := widgets.GetBorder(s.BorderStyle, s.BorderType)

	fwLine := func(str string, w int) string {
		return utils.FullWidthBg(str, w)
	}

	var netBlocks []string
	for i, nic := range s.NetworkInterfaces {
		cW := colWidths[i%cols] - 4

		// Check for inactivity
		if nic.BytesRecv == 0 && nic.BytesSent == 0 {
			inactiveStyle := lipgloss.NewStyle().Foreground(mu).Align(lipgloss.Center).Width(cW)
			msg := inactiveStyle.Render("(Inactive / Down)")

			c := container.Width(colWidths[i%cols]).Height(8).BorderTop(false)
			body := c.Render(lipgloss.NewStyle().MarginTop(3).Render(msg))

			topBorder := widgets.RenderTopBorderWithBg(nic.Name, colWidths[i%cols], border, mu, p)
			netBlocks = append(netBlocks, lipgloss.JoinVertical(lipgloss.Left, topBorder, body))
			continue
		}

		var rxRate, txRate float64
		if last, ok := s.LastNetworkInterfaces[nic.Name]; ok {
			if nic.BytesRecv >= last.BytesRecv {
				rxRate = float64(nic.BytesRecv-last.BytesRecv) / 1024 / 1024
			}
			if nic.BytesSent >= last.BytesSent {
				txRate = float64(nic.BytesSent-last.BytesSent) / 1024 / 1024
			}
		}

		stats := lipgloss.JoinVertical(lipgloss.Left,
			fwLine(lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Rx Total: ")+valueStyle.Render(utils.FormatBytes(nic.BytesRecv))), cW),
			fwLine(lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Tx Total: ")+valueStyle.Render(utils.FormatBytes(nic.BytesSent))), cW),
			fwLine(lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Rx Rate:  ")+valueStyle.Render(fmt.Sprintf("%.2f MB/s", rxRate))), cW),
			fwLine(lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Tx Rate:  ")+valueStyle.Render(fmt.Sprintf("%.2f MB/s", txRate))), cW),
			fwLine(lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Errors:   ")+valueStyle.Render(fmt.Sprintf("In:%d Out:%d", nic.Errin, nic.Errout))), cW),
			fwLine(lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Dropped:  ")+valueStyle.Render(fmt.Sprintf("In:%d Out:%d", nic.Dropin, nic.Dropout))), cW),
		)

		c := container.Width(colWidths[i%cols]).Height(8).BorderTop(false)
		body := c.Render(lipgloss.NewStyle().MarginTop(1).Render(stats))
		topBorder := widgets.RenderTopBorderWithBg(nic.Name, colWidths[i%cols], border, b, p)

		netBlocks = append(netBlocks, lipgloss.JoinVertical(lipgloss.Left, topBorder, body))
	}

	var rows []string
	for i := 0; i < len(netBlocks); i += cols {
		end := i + cols
		if end > len(netBlocks) {
			end = len(netBlocks)
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, netBlocks[i:end]...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

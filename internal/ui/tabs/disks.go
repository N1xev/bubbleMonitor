package tabs

import (
	"image/color"

	"charm.land/lipgloss/v2"

	"github.com/N1xev/bubbleMonitor/internal/data"
	"github.com/N1xev/bubbleMonitor/internal/ui/widgets"
	"github.com/N1xev/bubbleMonitor/internal/util"
)

// RenderDisks renders the disk partitions tab
func RenderDisks(s *data.AppState, container lipgloss.Style, su, w, a, t, mu, p, b color.Color, availHeight int) string {
	if len(s.Metrics.DiskPartitions) == 0 {
		return "Loading disk information..."
	}

	width := s.UI.Width
	boxWidth := width
	contentWidth := boxWidth - 4
	barWidth := boxWidth - 6
	border := widgets.GetBorder(s.Config.BorderStyle, s.Config.BorderType)

	mountStyle := lipgloss.NewStyle().Bold(true).Foreground(t)
	infoStyle := lipgloss.NewStyle().Foreground(mu)

	s.UI.ContentBuilder.Reset()
	for i, d := range s.Metrics.DiskPartitions {
		bar := widgets.RenderProgressBar(d.UsedPct, barWidth, su, w, a)
		info := "Used: " + util.FormatBytes(d.Used) + " / " + util.FormatBytes(d.Total)

		s.UI.ContentBuilder.WriteString(mountStyle.Render(util.PadRight(d.Mountpoint, contentWidth)))
		s.UI.ContentBuilder.WriteString("\n")
		s.UI.ContentBuilder.WriteString(bar)
		s.UI.ContentBuilder.WriteString("\n")
		s.UI.ContentBuilder.WriteString(infoStyle.Render(util.PadRight(info, contentWidth)))

		if i < len(s.Metrics.DiskPartitions)-1 {
			s.UI.ContentBuilder.WriteByte('\n')
		}
	}

	content := s.UI.ContentBuilder.String()

	titleText := "DISK PARTITIONS (Total R: " + util.FastMbPerSec(s.Metrics.DiskReadRate) + " W: " + util.FastMbPerSec(s.Metrics.DiskWriteRate) + ")"

	contentHeight := max(availHeight-2, 0)

	c := container.Width(boxWidth).Height(contentHeight).BorderTop(false)
	body := c.Render(content)
	topBorder := widgets.RenderTopBorderWithBg(titleText, boxWidth, border, b, p)

	s.UI.ContentBuilder.Reset()
	s.UI.ContentBuilder.WriteString(topBorder)
	s.UI.ContentBuilder.WriteByte('\n')
	s.UI.ContentBuilder.WriteString(body)

	return s.UI.ContentBuilder.String()
}

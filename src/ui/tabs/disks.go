package tabs

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"

	"github.com/N1xev/bubbleMonitor/src/data"
	"github.com/N1xev/bubbleMonitor/src/ui/widgets"
	"github.com/N1xev/bubbleMonitor/src/utils"
)

// RenderDisks renders the disk partitions tab
func RenderDisks(s *data.AppState, container lipgloss.Style, su, w, a, t, mu, p, b compat.AdaptiveColor, availHeight int) string {
	if len(s.DiskPartitions) == 0 {
		return "Loading disk information..."
	}

	width := s.Width
	boxWidth := width
	contentWidth := boxWidth - 4
	border := widgets.GetBorder(s.BorderStyle, s.BorderType)

	fwLine := func(str string) string {
		return utils.FullWidthBg(str, contentWidth)
	}

	mountStyle := lipgloss.NewStyle().Bold(true).Foreground(t)
	infoStyle := lipgloss.NewStyle().Foreground(mu)

	var diskBlocks []string
	for _, d := range s.DiskPartitions {
		bar := widgets.RenderProgressBar(d.UsedPct, contentWidth, su, w, a)

		info := fmt.Sprintf("Used: %s / %s", utils.FormatBytes(d.Used), utils.FormatBytes(d.Total))

		block := lipgloss.JoinVertical(lipgloss.Left,
			fwLine(mountStyle.Render(d.Mountpoint)),
			fwLine(bar),
			fwLine(infoStyle.Render(info)),
		)
		diskBlocks = append(diskBlocks, block)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, diskBlocks...)

	titleText := fmt.Sprintf("DISK PARTITIONS (Total R: %.2f MB/s W: %.2f MB/s)", s.DiskReadRate, s.DiskWriteRate)

	contentHeight := availHeight - 2
	if contentHeight < 0 {
		contentHeight = 0
	}

	c := container.Width(boxWidth).Height(contentHeight).BorderTop(false)
	body := c.Render(content)
	topBorder := widgets.RenderTopBorderWithBg(titleText, boxWidth, border, b, p)

	return lipgloss.JoinVertical(lipgloss.Left, topBorder, body)
}

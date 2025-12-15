package tabs

import (
	"fmt"
	"time"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"

	"github.com/N1xev/bubbleMonitor/src/data"
	"github.com/N1xev/bubbleMonitor/src/ui/widgets"
	"github.com/N1xev/bubbleMonitor/src/utils"
)

// RenderOverview renders the overview tab
func RenderOverview(s *data.AppState, container, titleStyle, labelStyle, valueStyle lipgloss.Style, su, w, a, t, mu, bg, p, b compat.AdaptiveColor, availHeight int) string {
	width := s.Width
	cols := 1
	if width >= 80 {
		cols = 2
	}
	if width >= 120 {
		cols = 3
	}

	colWidths := utils.CalculateColumnWidths(width, cols)
	border := widgets.GetBorder(s.BorderStyle, s.BorderType)

	getContentWidth := func(i int) int {
		return colWidths[i%cols] - 4
	}

	fwLine := func(str string, idx int) string {
		return utils.FullWidthBg(str, getContentWidth(idx))
	}

	sp := func(str string) string { return str }

	// CPU (Index 0)
	idx := 0
	cw := getContentWidth(idx)
	lAS := "N/A"
	if s.LoadAvg != nil {
		lAS = fmt.Sprintf("%.2f, %.2f, %.2f", s.LoadAvg.Load1, s.LoadAvg.Load5, s.LoadAvg.Load15)
	}
	cpuBar := widgets.RenderProgressBar(s.Cpu, cw, su, w, a)
	cpuVal := valueStyle.Foreground(widgets.GetColorForValue(s.Cpu, su, w, a)).Render(fmt.Sprintf("%.1f%%", s.Cpu))
	cpuTemp := "N/A"
	if s.CpuTemp > 0 {
		cpuTemp = fmt.Sprintf("%.1f°C", s.CpuTemp)
	} else {
		// Hint for Windows users
		cpuTemp = "N/A (Admin?)"
	}
	cpuBlock := lipgloss.JoinVertical(lipgloss.Left,
		fwLine(cpuVal, idx),
		fwLine(cpuBar, idx),
		fwLine(labelStyle.Render("Load:")+sp(" ")+labelStyle.Render(lAS), idx),
		fwLine(labelStyle.Render("Temp:")+sp(" ")+labelStyle.Render(cpuTemp), idx),
	)

	// Memory (Index 1) - Use cached MemInfo
	idx = 1
	cw = getContentWidth(idx)
	memBar := widgets.RenderProgressBar(s.Memory, cw, su, w, a)
	memVal := valueStyle.Foreground(widgets.GetColorForValue(s.Memory, su, w, a)).Render(fmt.Sprintf("%.1f%%", s.Memory))

	memUsed := uint64(0)
	memTotal := uint64(0)
	if s.MemInfo != nil {
		memUsed = s.MemInfo.Used
		memTotal = s.MemInfo.Total
	}
	memBlock := lipgloss.JoinVertical(lipgloss.Left,
		fwLine(memVal, idx),
		fwLine(memBar, idx),
		fwLine(labelStyle.Render(utils.FormatBytes(memUsed))+sp(" / ")+labelStyle.Render(utils.FormatBytes(memTotal)), idx),
	)

	// Disk (Index 2) - Use cached disk data from partitions
	idx = 2
	cw = getContentWidth(idx)
	diskBar := widgets.RenderProgressBar(s.Disk, cw, su, w, a)
	diskVal := valueStyle.Foreground(widgets.GetColorForValue(s.Disk, su, w, a)).Render(fmt.Sprintf("%.1f%%", s.Disk))

	diskUsed := uint64(0)
	diskTotal := uint64(0)
	if len(s.DiskPartitions) > 0 {
		diskUsed = s.DiskPartitions[0].Used
		diskTotal = s.DiskPartitions[0].Total
	}
	diskBlock := lipgloss.JoinVertical(lipgloss.Left,
		fwLine(diskVal, idx),
		fwLine(diskBar, idx),
		fwLine(labelStyle.Render(utils.FormatBytes(diskUsed))+sp(" / ")+labelStyle.Render(utils.FormatBytes(diskTotal)), idx),
	)

	// Network (Index 3)
	idx = 3
	cw = getContentWidth(idx)
	nP := CalcNetPercent(s)
	netBar := widgets.RenderProgressBar(nP, cw, su, w, a)
	netVal := valueStyle.Foreground(widgets.GetColorForValue(nP, su, w, a)).Render(fmt.Sprintf("%.1f%%", nP))

	recvMb := s.NetRecvRate * 8
	sentMb := s.NetSentRate * 8

	netBlock := lipgloss.JoinVertical(lipgloss.Left,
		fwLine(netVal, idx),
		fwLine(netBar, idx),
		fwLine(labelStyle.Render("↓")+sp(" ")+labelStyle.Render(fmt.Sprintf("%.2f Mb/s", recvMb))+sp(" ")+labelStyle.Render("↑")+sp(" ")+labelStyle.Render(fmt.Sprintf("%.2f Mb/s", sentMb)), idx),
	)

	// Quick Stats (Index 4)
	idx = 4

	uptimeStr := "Loading..."
	if s.HostInfo != nil {
		uptimeStr = utils.FormatDuration(time.Duration(s.HostInfo.Uptime) * time.Second)
	}

	loadStr := "Loading..."
	if s.LoadAvg != nil {
		if s.LoadAvg.Load1 == 0 && s.LoadAvg.Load5 == 0 && s.LoadAvg.Load15 == 0 {
			loadStr = "N/A (Windows)"
		} else {
			loadStr = fmt.Sprintf("%.2f, %.2f, %.2f", s.LoadAvg.Load1, s.LoadAvg.Load5, s.LoadAvg.Load15)
		}
	}

	batteryStr := "N/A"
	if len(s.Battery) > 0 {
		batt := s.Battery[0]
		pct := batt.Current / batt.Full * 100
		batteryStr = fmt.Sprintf("%.0f%% (%s)", pct, batt.State.String())
	}

	quickStatsBlock := lipgloss.JoinVertical(lipgloss.Left,
		fwLine(lipgloss.JoinHorizontal(lipgloss.Top,
			labelStyle.Render("Uptime:"+sp("      ")),
			valueStyle.Render(uptimeStr),
		), idx),
		fwLine(lipgloss.JoinHorizontal(lipgloss.Top,
			labelStyle.Render("Processes:"+sp("   ")),
			valueStyle.Render(fmt.Sprintf("%d", len(s.Processes))),
		), idx),
		fwLine(lipgloss.JoinHorizontal(lipgloss.Top,
			labelStyle.Render("Load Avg:"+sp("    ")),
			valueStyle.Render(loadStr),
		), idx),
		fwLine(lipgloss.JoinHorizontal(lipgloss.Top,
			labelStyle.Render("Battery:"+sp("     ")),
			valueStyle.Render(batteryStr),
		), idx),
	)

	// Swap (Index 5) - Use cached SwapInfo
	swapBlock := ""
	idx = 5
	if s.Swap > 0 && s.SwapInfo != nil {
		cw = getContentWidth(idx)
		swapBar := widgets.RenderProgressBar(s.Swap, cw, su, w, a)
		swapVal := valueStyle.Foreground(widgets.GetColorForValue(s.Swap, su, w, a)).Render(fmt.Sprintf("%.1f%%", s.Swap))
		swapBlock = lipgloss.JoinVertical(lipgloss.Left,
			fwLine(swapVal, idx),
			fwLine(swapBar, idx),
			fwLine(labelStyle.Render(utils.FormatBytes(s.SwapInfo.Used))+sp(" / ")+labelStyle.Render(utils.FormatBytes(s.SwapInfo.Total)), idx),
		)
	}

	blocks := []string{cpuBlock, memBlock, diskBlock, netBlock, quickStatsBlock}
	if swapBlock != "" {
		blocks = append(blocks, swapBlock)
	}

	numRows := (len(blocks) + cols - 1) / cols
	rowHeight := availHeight / numRows
	contentHeight := rowHeight - 2
	if contentHeight < 1 {
		contentHeight = 1
	}

	titles := []string{"CPU USAGE", "MEMORY", "DISK USAGE", "NETWORK", "QUICK STATS"}
	if swapBlock != "" {
		titles = append(titles, "SWAP")
	}

	var renderedBlocks []string
	for i, block := range blocks {
		title := titles[i]
		bw := colWidths[i%cols]

		blockWithMargin := lipgloss.NewStyle().MarginTop(1).Render(block)
		c := container.Width(bw).Height(contentHeight).BorderTop(false).Padding(0, 1)
		body := c.Render(blockWithMargin)
		topBorder := widgets.RenderTopBorderWithBg(title, bw, border, b, p)

		renderedBlocks = append(renderedBlocks, lipgloss.JoinVertical(lipgloss.Left, topBorder, body))
	}

	for i := len(renderedBlocks); i < numRows*cols; i++ {
		bw := colWidths[i%cols]
		emptyBlock := lipgloss.Place(bw, rowHeight, lipgloss.Center, lipgloss.Center, "")
		renderedBlocks = append(renderedBlocks, emptyBlock)
	}

	var rows []string
	for i := 0; i < len(renderedBlocks); i += cols {
		end := i + cols
		if end > len(renderedBlocks) {
			end = len(renderedBlocks)
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, renderedBlocks[i:end]...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

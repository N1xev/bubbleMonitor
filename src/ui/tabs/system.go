package tabs

import (
	"fmt"
	"runtime"
	"time"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"

	"github.com/N1xev/bubbleMonitor/src/data"
	"github.com/N1xev/bubbleMonitor/src/ui/widgets"
	"github.com/N1xev/bubbleMonitor/src/utils"
)

// RenderSystem renders the system information tab
func RenderSystem(s *data.AppState, container, titleStyle, labelStyle, valueStyle lipgloss.Style, t, mu, p, b, bg compat.AdaptiveColor, availHeight int) string {
	if s.HostInfo == nil {
		return "Loading system information..."
	}

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

	// Host Info (Index 0)
	idx := 0
	hostInfo := lipgloss.JoinVertical(lipgloss.Left,
		fwLine(lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Hostname:"+sp(" ")), valueStyle.Render(s.HostInfo.Hostname)), idx),
		fwLine(lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("OS:"+sp("       ")), valueStyle.Render(s.HostInfo.OS)), idx),
		fwLine(lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Platform:"+sp(" ")), valueStyle.Render(s.HostInfo.Platform)), idx),
		fwLine(lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Kernel:"+sp("   ")), valueStyle.Render(s.HostInfo.KernelVersion)), idx),
	)

	// CPU Info (Index 1) - Use cached CpuInfoStatic
	idx = 1
	cpuModel := "N/A"
	if len(s.CpuInfoStatic) > 0 {
		cpuModel = s.CpuInfoStatic[0].ModelName
	}
	cpuInfo := lipgloss.JoinVertical(lipgloss.Left,
		fwLine(lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Model:")+sp("    "), valueStyle.Render(cpuModel)), idx),
		fwLine(lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Cores:")+sp("    "), valueStyle.Render(fmt.Sprintf("%d physical,", runtime.NumCPU()))+sp(" ")+valueStyle.Render(fmt.Sprintf("%d logical", runtime.GOMAXPROCS(0)))), idx),
	)

	// Uptime & Load (Index 2)
	idx = 2

	l1, l5, l15 := "...", "...", "..."
	if s.LoadAvg != nil {
		if s.LoadAvg.Load1 == 0 && s.LoadAvg.Load5 == 0 && s.LoadAvg.Load15 == 0 {
			l1, l5, l15 = "N/A", "N/A", "N/A"
		} else {
			l1 = fmt.Sprintf("%.2f", s.LoadAvg.Load1)
			l5 = fmt.Sprintf("%.2f", s.LoadAvg.Load5)
			l15 = fmt.Sprintf("%.2f", s.LoadAvg.Load15)
		}
	}

	uptimeLoad := lipgloss.JoinVertical(lipgloss.Left,
		fwLine(lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Uptime:"+sp("    ")), valueStyle.Render(utils.FormatDuration(time.Duration(s.HostInfo.Uptime)*time.Second))), idx),
		fwLine(lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Load 1m:"+sp("   ")), valueStyle.Render(l1)), idx),
		fwLine(lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Load 5m:"+sp("   ")), valueStyle.Render(l5)), idx),
		fwLine(lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Load 15m:"+sp("  ")), valueStyle.Render(l15)), idx),
	)

	// Uptime (Index 3)
	idx = 3
	bootTime := time.Unix(int64(s.HostInfo.BootTime), 0)
	uptime := lipgloss.JoinVertical(lipgloss.Left,
		fwLine(lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Boot:"+sp("   ")), valueStyle.Render(bootTime.Format("2006-01-02 15:04:05"))), idx),
		fwLine(lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("System:"+sp(" ")), valueStyle.Render(utils.FormatDuration(time.Since(bootTime)))), idx),
		fwLine(lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Monitor:"), valueStyle.Render(utils.FormatDuration(time.Since(s.StartTime)))), idx),
	)

	// Memory Info (Index 4) - Use cached MemInfo
	idx = 4
	memTotal := uint64(0)
	memAvailable := uint64(0)
	memUsed := uint64(0)
	memUsedPct := 0.0
	if s.MemInfo != nil {
		memTotal = s.MemInfo.Total
		memAvailable = s.MemInfo.Available
		memUsed = s.MemInfo.Used
		memUsedPct = s.MemInfo.UsedPercent
	}
	memInfo := lipgloss.JoinVertical(lipgloss.Left,
		fwLine(lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Total:"+sp("     ")), valueStyle.Render(utils.FormatBytes(memTotal))), idx),
		fwLine(lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Available:"+sp(" ")), valueStyle.Render(utils.FormatBytes(memAvailable))), idx),
		fwLine(lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Used:"+sp("      ")), valueStyle.Render(fmt.Sprintf("%s"+sp(" ")+"(%.2f%%)", utils.FormatBytes(memUsed), memUsedPct))), idx),
	)

	// GPU Info (Index 5)
	idx = 5
	var gpuInfo string
	if len(s.GpuInfo) > 0 {
		gpu := s.GpuInfo[0]
		gpuInfo = lipgloss.JoinVertical(lipgloss.Left,
			fwLine(lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Name:")+sp("   "), valueStyle.Render(gpu.Name)), idx),
			fwLine(lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Driver:")+sp(" "), valueStyle.Render(gpu.Driver)), idx),
			fwLine(lipgloss.JoinHorizontal(lipgloss.Top, labelStyle.Render("Memory:")+sp(" "), valueStyle.Render(fmt.Sprintf("%sMB", gpu.MemoryUsed))+sp(" / ")+valueStyle.Render(fmt.Sprintf("%sMB", gpu.MemoryTotal))), idx),
		)
	}

	blocks := []string{hostInfo, cpuInfo, uptimeLoad, uptime, memInfo}
	if gpuInfo != "" {
		blocks = append(blocks, gpuInfo)
	}

	numRows := (len(blocks) + cols - 1) / cols
	rowHeight := availHeight / numRows
	contentHeight := rowHeight - 2
	if contentHeight < 1 {
		contentHeight = 1
	}

	titles := []string{"HOST", "CPU", "SYSTEM LOAD", "TIMINGS", "MEMORY"}
	if gpuInfo != "" {
		titles = append(titles, "GPU")
	}

	var renderedBlocks []string
	for i, block := range blocks {
		title := titles[i]
		bw := colWidths[i%cols]

		c := container.Width(bw).Height(contentHeight).BorderTop(false)
		body := c.Render(block)
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

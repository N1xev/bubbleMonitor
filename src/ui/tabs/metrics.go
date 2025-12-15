package tabs

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"

	"github.com/N1xev/bubbleMonitor/src/data"
	"github.com/N1xev/bubbleMonitor/src/ui/widgets"
	"github.com/N1xev/bubbleMonitor/src/utils"
)

// SumAccessor sums two accessors at each index
type SumAccessor struct {
	A, B data.Accessor
}

func (s *SumAccessor) Len() int {
	if s.A == nil {
		return 0
	}
	return s.A.Len()
}

func (s *SumAccessor) Get(i int) float64 {
	val := 0.0
	if s.A != nil {
		val += s.A.Get(i)
	}
	if s.B != nil {
		val += s.B.Get(i)
	}
	return val
}

// RenderMetrics renders the metrics/charts tab
func RenderMetrics(app *data.AppState, container lipgloss.Style, su, w, a, s, t, mu, p, b compat.AdaptiveColor, availHeight int) string {
	width := app.Width

	// 1. Calculate Layout Columns
	// Charts
	chartCols := 1
	if width >= 100 {
		chartCols = 2
	}

	// Cores
	coreCols := 1
	if width >= 60 {
		coreCols = 2
	}
	if width >= 90 {
		coreCols = 3
	}
	if width >= 120 {
		coreCols = 4
	}
	if width >= 150 {
		coreCols = 5
	}

	// 2. Calculate Required Core Height
	numCores := len(app.CpuPerCore)
	if numCores == 0 {
		numCores = 1
	}

	numCoreRows := (numCores + coreCols - 1) / coreCols
	coreSectionHeight := numCoreRows + 2

	// 3. Calculate Available Chart Height
	availChartSpace := availHeight - coreSectionHeight
	minChartSpace := 6
	if availChartSpace < minChartSpace {
		availChartSpace = minChartSpace
	}

	// 4. Distribute Space to Charts
	numChartRows := (4 + chartCols - 1) / chartCols
	chartBlockHeight := availChartSpace / numChartRows
	if chartBlockHeight < 5 {
		chartBlockHeight = 5
	}

	// 5. Render Charts
	textStyle := lipgloss.NewStyle().Foreground(t)
	chartsTitles := []string{
		fmt.Sprintf("CPU HISTORY (Window: %ds)", app.HistoryLength),
		fmt.Sprintf("MEMORY HISTORY (Window: %ds)", app.HistoryLength),
		fmt.Sprintf("NETWORK ACTIVITY (Window: %ds)", app.HistoryLength),
		fmt.Sprintf("DISK I/O HISTORY (Window: %ds)", app.HistoryLength),
	}
	var renderedCharts []string

	chartWidths := utils.CalculateColumnWidths(width, chartCols)

	// Helper to render chart based on ChartType
	renderChart := func(data data.Accessor, chartW, chartH int, c1, c2 compat.AdaptiveColor) string {
		switch app.ChartType {
		case "line":
			return widgets.RenderLineChart(data, chartW, chartH, c1, c2)
		case "bar":
			return widgets.RenderBarChart(data, chartW, chartH, c1, c2)
		case "braille":
			return widgets.RenderBrailleChart(data, chartW, chartH, c1, c2)
		case "tty":
			return widgets.RenderTTYChart(data, chartW, chartH, c1, c2)
		default:
			return widgets.RenderSparkline(data, chartW, chartH, c1, c2)
		}
	}

	for i := 0; i < 4; i++ {
		var boxW int
		if chartCols == 1 {
			boxW = width
		} else {
			boxW = chartWidths[i%chartCols]
		}

		sparklineH := chartBlockHeight - 3
		if sparklineH < 1 {
			sparklineH = 1
		}

		contentW := boxW - 4
		chartW := contentW - 6
		if chartW < 5 {
			chartW = 5
		}

		var innerBlock string

		switch i {
		case 0: // CPU
			ch := renderChart(app.CpuHistory, chartW, sparklineH, p, w)
			stats := fmt.Sprintf("Cur: %.1f%% Avg: %.1f%% Peak: %.1f%%", app.Cpu, app.CpuHistory.Avg(), app.CpuHistory.Max())
			innerBlock = lipgloss.JoinVertical(lipgloss.Left, ch, textStyle.Render(stats))
		case 1: // Mem
			ch := renderChart(app.MemHistory, chartW, sparklineH, s, w)
			stats := fmt.Sprintf("Cur: %.1f%% Avg: %.1f%% Peak: %.1f%%", app.Memory, app.MemHistory.Avg(), app.MemHistory.Max())
			innerBlock = lipgloss.JoinVertical(lipgloss.Left, ch, textStyle.Render(stats))
		case 2: // Net
			ch := renderChart(app.NetHistory, chartW, sparklineH, su, w)
			stats := fmt.Sprintf("Peak: %.1f%% Recv: %.2f MB/s Sent: %.2f MB/s", app.NetHistory.Max(), app.NetRecvRate, app.NetSentRate)
			innerBlock = lipgloss.JoinVertical(lipgloss.Left, ch, textStyle.Render(stats))
		case 3: // Disk I/O
			totalIO := &SumAccessor{A: app.DiskHORead, B: app.DiskHOWrite}
			ch := renderChart(totalIO, chartW, sparklineH, mu, w)
			stats := fmt.Sprintf("Read: %.2f MB/s Write: %.2f MB/s", app.DiskReadRate, app.DiskWriteRate)
			innerBlock = lipgloss.JoinVertical(lipgloss.Left, ch, textStyle.Render(stats))
		}

		c := container.Width(boxW).Height(chartBlockHeight - 1).BorderTop(false)
		body := c.Render(innerBlock)

		topBorder := widgets.RenderTopBorderWithBg(chartsTitles[i], boxW, widgets.GetBorder(app.BorderStyle, app.BorderType), b, p)
		renderedCharts = append(renderedCharts, lipgloss.JoinVertical(lipgloss.Left, topBorder, body))
	}

	// Assemble Top Section (Charts)
	var rows []string
	for i := 0; i < len(renderedCharts); i += chartCols {
		end := i + chartCols
		if end > len(renderedCharts) {
			end = len(renderedCharts)
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, renderedCharts[i:end]...))
	}
	topSection := lipgloss.JoinVertical(lipgloss.Left, rows...)

	// 6. Render Cores
	coreColWidths := utils.CalculateColumnWidths(width, coreCols)
	var coreBlocks []string
	textStyle = lipgloss.NewStyle().Foreground(t)

	for i, usage := range app.CpuPerCore {
		cW := coreColWidths[i%coreCols] - 4
		if cW < 10 {
			cW = 10
		}
		barW := cW - 16
		if barW < 5 {
			barW = 5
		}

		bar := widgets.RenderProgressBar(usage, barW, su, w, a)
		line := lipgloss.JoinHorizontal(lipgloss.Left,
			textStyle.Width(16).Render(fmt.Sprintf("Core %-2d: %5.1f%% ", i, usage)),
			bar,
		)
		coreBlocks = append(coreBlocks, line)
	}

	var coreRows []string
	for i := 0; i < len(coreBlocks); i += coreCols {
		end := i + coreCols
		if end > len(coreBlocks) {
			end = len(coreBlocks)
		}
		rowItems := coreBlocks[i:end]

		var rowStr string
		for j, item := range rowItems {
			w := coreColWidths[(i+j)%coreCols]
			rowStr = lipgloss.JoinHorizontal(lipgloss.Top, rowStr, lipgloss.NewStyle().Width(w).Render(item))
		}
		coreRows = append(coreRows, rowStr)
	}
	coresC := strings.Join(coreRows, "\n")

	// Assemble Bottom Section (Cores)
	coresBoxWidth := width
	topBorder := widgets.RenderTopBorderWithBg("CPU PER CORE", coresBoxWidth, widgets.GetBorder(app.BorderStyle, app.BorderType), b, p)

	coreBodyHeight := coreSectionHeight - 1
	c := container.Width(coresBoxWidth).Height(coreBodyHeight).BorderTop(false)
	body := c.Render(coresC)

	bottomSection := lipgloss.JoinVertical(lipgloss.Left, topBorder, body)

	return lipgloss.JoinVertical(lipgloss.Left, topSection, bottomSection)
}

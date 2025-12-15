package tabs

import (
	"fmt"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"

	"github.com/N1xev/bubbleMonitor/src/data"
	"github.com/N1xev/bubbleMonitor/src/ui/widgets"
	"github.com/N1xev/bubbleMonitor/src/utils"
)

// RenderProcesses renders the processes tab
func RenderProcesses(s *data.AppState, visibleProcs []data.ProcessInfo, treeIndents map[int32]int, container lipgloss.Style, su, w, a, t, mu, p, b compat.AdaptiveColor, availHeight int) string {
	boxWidth := s.Width
	border := widgets.GetBorder(s.BorderStyle, s.BorderType)

	// Split available height: processes list (70%) + details panel (30%)
	detailsHeight := 7
	listHeight := availHeight - detailsHeight - 1
	if listHeight < 10 {
		listHeight = 10
		detailsHeight = availHeight - listHeight - 1
	}

	contentWidth := boxWidth - 4

	pidWidth := 8
	statusWidth := 12
	cpuWidth := 8
	memWidth := 8

	nameWidth := contentWidth - pidWidth - statusWidth - cpuWidth - memWidth - 4
	if nameWidth < 20 {
		nameWidth = 20
	}

	sp := func(str string) string { return str }

	sI, mSI, pSI := "", "", ""
	switch s.SortBy {
	case "cpu":
		sI = " ▼"
	case "mem":
		mSI = " ▼"
	case "pid":
		pSI = " ▼"
	}

	hdrStyle := lipgloss.NewStyle().Bold(true).Underline(true)
	headerRow := hdrStyle.Width(pidWidth).Render("PID"+pSI) + sp(" ") +
		hdrStyle.Width(nameWidth).Render("NAME") + sp(" ") +
		hdrStyle.Width(statusWidth).Render("STATUS") + sp(" ") +
		hdrStyle.Width(cpuWidth).Align(lipgloss.Right).Render("CPU"+sI) + sp(" ") +
		hdrStyle.Width(memWidth).Align(lipgloss.Right).Render("MEM"+mSI)

	filtered := visibleProcs

	visibleRows := listHeight - 4
	if visibleRows < 1 {
		visibleRows = 1
	}

	startIdx := s.ProcessScrollOffset
	endIdx := startIdx + visibleRows
	if endIdx > len(filtered) {
		endIdx = len(filtered)
	}
	if startIdx >= len(filtered) {
		startIdx = 0
		endIdx = visibleRows
		if endIdx > len(filtered) {
			endIdx = len(filtered)
		}
	}

	selColor := compat.AdaptiveColor{Light: lipgloss.Color("#E0E7FF"), Dark: lipgloss.Color("#3730A3")}
	selectedStyle := lipgloss.NewStyle().Background(selColor)
	cellStyle := lipgloss.NewStyle()

	// Pre-allocate styles to avoid garbage creation in the loop
	styleLow := cellStyle.Foreground(su)
	styleMed := cellStyle.Foreground(w)
	styleHigh := cellStyle.Foreground(a)

	styleLowSel := styleLow.Background(selColor)
	styleMedSel := styleMed.Background(selColor)
	styleHighSel := styleHigh.Background(selColor)

	// Helper to pick style
	getStyle := func(val float64, selected bool) lipgloss.Style {
		if selected {
			if val < 50 {
				return styleLowSel
			}
			if val < 80 {
				return styleMedSel
			}
			return styleHighSel
		}
		if val < 50 {
			return styleLow
		}
		if val < 80 {
			return styleMed
		}
		return styleHigh
	}

	var rows []string
	var selectedProc *data.ProcessInfo

	for i := startIdx; i < endIdx; i++ {
		proc := filtered[i]
		isSelected := i == s.SelectedProcess

		if isSelected {
			selectedProc = &proc
		}

		currCellStyle := cellStyle
		if isSelected {
			currCellStyle = currCellStyle.Background(selColor)
		}

		// Use pre-allocated styles
		cpuStyle := getStyle(proc.Cpu, isSelected)
		memStyle := getStyle(proc.Memory, isSelected)

		name := proc.Name

		if s.TreeView {
			if level, ok := treeIndents[proc.Pid]; ok {
				prefix := strings.Repeat("  ", level)

				indicator := ""
				if s.CollapsedPids[proc.Pid] {
					indicator = "▶ "
				} else {
					if level > 0 {
						indicator = "└─ "
					}
				}
				name = prefix + indicator + name
			}
		}

		if len(name) > nameWidth {
			name = name[:nameWidth-3] + "..."
		}

		status := proc.Status
		if status == "" {
			status = "running"
		}

		if s.SuspendedState[proc.Pid] {
			status = "SUSPENDED"
		}

		if len(status) > statusWidth {
			status = status[:statusWidth-3] + "..."
		}

		// Apply widths via style
		pidCell := currCellStyle.Width(pidWidth).Render(fmt.Sprintf("%d", proc.Pid))
		nameCell := currCellStyle.Width(nameWidth).Render(name)

		var statusStr string
		if s.SuspendedState[proc.Pid] {
			warnStyle := currCellStyle.Foreground(lipgloss.Color("#F59E0B"))
			if isSelected {
				warnStyle = warnStyle.Background(selColor)
			}
			statusStr = warnStyle.Width(statusWidth).Render(status)
		} else {
			statusStr = currCellStyle.Width(statusWidth).Render(status)
		}

		cpuStr := fmt.Sprintf("%.1f%%", proc.Cpu)
		memStr := fmt.Sprintf("%.1f%%", proc.Memory)

		cpuCell := cpuStyle.Width(cpuWidth).Align(lipgloss.Right).Render(cpuStr)
		memCell := memStyle.Width(memWidth).Align(lipgloss.Right).Render(memStr)

		space := " "
		if isSelected {
			space = selectedStyle.Render(" ")
		}

		// Compose row
		rowContent := pidCell + space + nameCell + space + statusStr + space + cpuCell + space + memCell

		row := lipgloss.NewStyle().Width(contentWidth).Render(rowContent)

		if isSelected {
			row = selectedStyle.Width(contentWidth).Render(rowContent)
		}
		rows = append(rows, row)
	}

	scrollInfo := ""
	if len(filtered) > visibleRows {
		scrollInfo = fmt.Sprintf(" [%d-%d of %d]", startIdx+1, endIdx, len(filtered))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lipgloss.NewStyle().Width(contentWidth).Render(headerRow), "", strings.Join(rows, "\n"))

	titleText := fmt.Sprintf("PROCESSES (Sort: %s)%s", strings.ToUpper(s.SortBy), scrollInfo)

	listContentHeight := listHeight - 2
	if listContentHeight < 0 {
		listContentHeight = 0
	}

	c := container.Width(boxWidth).Height(listContentHeight).BorderTop(false)
	body := c.Render(content)
	topBorder := widgets.RenderTopBorderWithBg(titleText, boxWidth, border, b, p)

	listBlock := lipgloss.JoinVertical(lipgloss.Left, topBorder, body)

	detailsBlock := renderProcessDetails(s, selectedProc, container, boxWidth, detailsHeight-2, t, mu, p, b, su, w, a)

	var filterIndicator string
	if s.FilterMode {
		filterIndicator = lipgloss.NewStyle().
			Foreground(p).
			Bold(true).
			MarginLeft(2).
			Render(fmt.Sprintf(" Filter: %s█", s.ProcessFilter))
	} else if s.ProcessFilter != "" {
		filterIndicator = lipgloss.NewStyle().
			Foreground(mu).
			MarginLeft(2).
			Render(fmt.Sprintf(" Filter: %s (press 'c' to clear, 'f' to edit)", s.ProcessFilter))
	}

	result := lipgloss.JoinVertical(lipgloss.Left, listBlock, detailsBlock)
	if filterIndicator != "" {
		result = lipgloss.JoinVertical(lipgloss.Left, filterIndicator, result)
	}

	return result
}

// renderProcessDetails renders the details panel for the selected process
func renderProcessDetails(s *data.AppState, proc *data.ProcessInfo, container lipgloss.Style, boxWidth, contentHeight int, t, mu, p, b, su, w, a compat.AdaptiveColor) string {
	border := widgets.GetBorder(s.BorderStyle, s.BorderType)

	if proc == nil {
		c := container.Width(boxWidth).Height(contentHeight).BorderTop(false)
		body := c.Render(lipgloss.NewStyle().Foreground(mu).Render("No process selected - use j/k or ↑↓ to navigate"))
		topBorder := widgets.RenderTopBorderWithBg("PROCESS DETAILS", boxWidth, border, b, p)
		return lipgloss.JoinVertical(lipgloss.Left, topBorder, body)
	}

	labelStyle := lipgloss.NewStyle().Foreground(mu)
	valueStyle := lipgloss.NewStyle().Foreground(t).Bold(true)
	contentWidth := boxWidth - 4

	col1Width := contentWidth / 3
	col2Width := contentWidth / 3
	col3Width := contentWidth - col1Width - col2Width

	statusColor := su
	status := proc.Status
	if status == "" {
		status = "running"
	}

	if s.SuspendedState[proc.Pid] {
		status = "SUSPENDED"
		statusColor = compat.AdaptiveColor{Light: lipgloss.Color("#F59E0B"), Dark: lipgloss.Color("#F59E0B")}
	}

	memStr := utils.FormatBytes(proc.MemoryBytes)

	username := proc.Username
	if len(username) > 15 {
		username = username[:12] + "..."
	}
	if username == "" {
		username = "N/A"
	}

	cpuColor := widgets.GetColorForValue(proc.Cpu, su, w, a)
	memColor := widgets.GetColorForValue(proc.Memory, su, w, a)

	leftCol := lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.Render("PID: ")+valueStyle.Render(fmt.Sprintf("%d", proc.Pid)),
		labelStyle.Render("Name: ")+valueStyle.Render(proc.Name),
		labelStyle.Render("Status: ")+lipgloss.NewStyle().Foreground(statusColor).Bold(true).Render(status),
	)

	midCol := lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.Render("User: ")+valueStyle.Render(username),
		labelStyle.Render("Nice: ")+valueStyle.Render(fmt.Sprintf("%d", proc.Nice)),
		labelStyle.Render("PPID: ")+valueStyle.Render(fmt.Sprintf("%d", proc.Ppid)),
	)

	rightCol := lipgloss.JoinVertical(lipgloss.Left,
		labelStyle.Render("CPU: ")+lipgloss.NewStyle().Foreground(cpuColor).Bold(true).Render(fmt.Sprintf("%.1f%%", proc.Cpu)),
		labelStyle.Render("Memory: ")+lipgloss.NewStyle().Foreground(memColor).Bold(true).Render(fmt.Sprintf("%.1f%% (%s)", proc.Memory, memStr)),
		labelStyle.Render("Started: ")+valueStyle.Render(time.Unix(proc.CreateTime/1000, 0).Format("15:04:05")),
	)

	details := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(col1Width).Render(leftCol),
		lipgloss.NewStyle().Width(col2Width).Render(midCol),
		lipgloss.NewStyle().Width(col3Width).Render(rightCol),
	)

	c := container.Width(boxWidth).Height(contentHeight).BorderTop(false)
	body := c.Render(details)
	topBorder := widgets.RenderTopBorderWithBg("PROCESS DETAILS", boxWidth, border, b, p)

	return lipgloss.JoinVertical(lipgloss.Left, topBorder, body)
}

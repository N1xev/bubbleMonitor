package overlays

import (
	"fmt"
	"slices"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"

	"github.com/N1xev/bubbleMonitor/internal/config"
	"github.com/N1xev/bubbleMonitor/internal/data"
	"github.com/N1xev/bubbleMonitor/internal/ui/widgets"
)

// RenderSettingsOverlay renders the settings configuration modal
func RenderSettingsOverlay(s *data.AppState, width, height int, b, p, t, mu, bg compat.AdaptiveColor) string {
	boxWidth := min(data.SettingsDefaultWidth, width-4)

	boxHeight := min(data.SettingsDefaultHeight, height-4)
	border := widgets.GetBorder(s.Config.BorderStyle, s.Config.BorderType)

	itemStyle := lipgloss.NewStyle().Foreground(t)
	selectedStyle := lipgloss.NewStyle().Foreground(p).Bold(true).Border(border, false, false, false, true).BorderForeground(p).PaddingLeft(1)
	headerStyle := lipgloss.NewStyle().Foreground(p).Bold(true).MarginBottom(1)

	var col1 []string
	col1 = append(col1, headerStyle.Render("THRESHOLDS & DISPLAY"))

	thresholdItems := []struct {
		label  string
		metric config.MetricType
		idx    int
	}{
		{"CPU Alert:", config.MetricCPU, 0},
		{"Mem Alert:", config.MetricMem, 1},
		{"Disk Alert:", config.MetricDisk, 2},
		{"Temp Alert:", config.MetricTemp, 3},
	}

	for _, item := range thresholdItems {
		val := fmt.Sprintf("%.0f%%", s.Config.Config.Thresholds[item.metric])
		if item.metric == config.MetricTemp {
			val = fmt.Sprintf("%.0f°C", s.Config.Config.Thresholds[item.metric])
		}
		line := fmt.Sprintf("%-15s %s", item.label, val)
		if s.UI.SettingsIdx == item.idx {
			col1 = append(col1, selectedStyle.Render(line))
		} else {
			col1 = append(col1, itemStyle.Render("  "+line))
		}
	}

	viewName := "normal"
	if s.Process.TreeView {
		viewName = "tree"
	}

	displayItems := []struct {
		label string
		value string
		idx   int
	}{
		{"Chart Type:", s.UI.ChartType, 4},
		{"View Type:", viewName, 5},
		{"Sort By:", s.Process.SortBy, 6},
		{"History Length:", fmt.Sprintf("%ds", s.UI.HistoryLength), 7},
		{"Process CPU:", "Raw", 8},
		{"Sort Direction:", s.Process.SortDirection, 9},
	}
	if s.Config.ProcessCpuNormalized {
		displayItems[4].value = "Normalized"
	}

	for _, item := range displayItems {
		line := fmt.Sprintf("%-15s %s", item.label, item.value)
		if s.UI.SettingsIdx == item.idx {
			col1 = append(col1, selectedStyle.Render(line))
		} else {
			col1 = append(col1, itemStyle.Render("  "+line))
		}
	}

	var col2 []string
	col2 = append(col2, headerStyle.Render("TABS"))

	allTabs := []string{"Metrics", "Processes", "Disks", "Network", "System", "Services", "Connections", "Logs", "Remote"}
	currentTabIdxBase := 10

	for i, tabName := range allTabs {
		idx := currentTabIdxBase + i
		isEnabled := slices.Contains(s.UI.ActiveTabs, tabName)

		status := "[ ]"
		if isEnabled {
			status = "[x]"
		}

		line := fmt.Sprintf("%-15s %s", tabName, status)
		if s.UI.SettingsIdx == idx {
			col2 = append(col2, selectedStyle.Render(line))
		} else {
			col2 = append(col2, itemStyle.Render("  "+line))
		}
	}

	var col3 []string
	col3 = append(col3, headerStyle.Render("APPEARANCE"))

	appearanceIdxBase := currentTabIdxBase + len(allTabs)

	bgLabel := "transparent"
	if s.Config.BackgroundOpaque {
		bgLabel = "opaque"
	}

	appItems := []struct {
		label string
		value string
		idx   int
	}{
		{"Theme:", s.Config.Theme, appearanceIdxBase},
		{"Refresh Rate:", fmt.Sprintf("%dms", s.Config.RefreshRate), appearanceIdxBase + 1},
		{"Border Type:", s.Config.BorderType, appearanceIdxBase + 2},
		{"Border Style:", s.Config.BorderStyle, appearanceIdxBase + 3},
		{"Background:", bgLabel, appearanceIdxBase + 4},
	}

	for _, item := range appItems {
		line := fmt.Sprintf("%-15s %s", item.label, item.value)
		if s.UI.SettingsIdx == item.idx {
			col3 = append(col3, selectedStyle.Render(line))
		} else {
			col3 = append(col3, itemStyle.Render("  "+line))
		}
	}

	contentWidth := boxWidth - 6
	colWidth := contentWidth / 3
	col1Content := lipgloss.JoinVertical(lipgloss.Left, col1...)
	col2Content := lipgloss.JoinVertical(lipgloss.Left, col2...)
	col3Content := lipgloss.JoinVertical(lipgloss.Left, col3...)

	contentBlock := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(colWidth).Render(col1Content),
		lipgloss.NewStyle().Width(colWidth).Render(col2Content),
		lipgloss.NewStyle().Width(colWidth).Render(col3Content),
	)

	hint := lipgloss.NewStyle().Foreground(mu).Align(lipgloss.Center).Width(boxWidth - 6).MarginTop(1).Render("↑/↓ select • ←/→ change/toggle • . to close")

	container := lipgloss.NewStyle().
		Border(border).
		BorderForeground(b).
		Padding(1, 2).
		Width(boxWidth - 6).
		Height(boxHeight).
		BorderTop(false)

	body := container.Render(lipgloss.JoinVertical(lipgloss.Left, contentBlock, hint))
	actualWidth := lipgloss.Width(body)

	topBorder := widgets.RenderTopBorderWithBg("CONFIGURATION", actualWidth, border, b, p)

	return lipgloss.JoinVertical(lipgloss.Left, topBorder, body)
}

package overlays

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"

	"github.com/N1xev/bubbleMonitor/src/config"
	"github.com/N1xev/bubbleMonitor/src/data"
	"github.com/N1xev/bubbleMonitor/src/ui/widgets"
)

// RenderSettingsOverlay renders the settings configuration modal
func RenderSettingsOverlay(s *data.AppState, width, height int, b, p, t, mu, bg compat.AdaptiveColor) string {
	boxWidth := 90
	if boxWidth > width-4 {
		boxWidth = width - 4
	}

	boxHeight := 22
	border := widgets.GetBorder(s.BorderStyle, s.BorderType)

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
		val := fmt.Sprintf("%.0f%%", s.Config.Thresholds[item.metric])
		if item.metric == config.MetricTemp {
			val = fmt.Sprintf("%.0f°C", s.Config.Thresholds[item.metric])
		}
		line := fmt.Sprintf("%-15s %s", item.label, val)
		if s.SettingsIdx == item.idx {
			col1 = append(col1, selectedStyle.Render(line))
		} else {
			col1 = append(col1, itemStyle.Render("  "+line))
		}
	}

	col1 = append(col1, "")

	viewName := "normal"
	if s.TreeView {
		viewName = "tree"
	}

	displayItems := []struct {
		label string
		value string
		idx   int
	}{
		{"Chart Type:", s.ChartType, 4},
		{"View Type:", viewName, 5},
		{"Sort By:", s.SortBy, 6},
		{"History Length:", fmt.Sprintf("%ds", s.HistoryLength), 7},
	}

	for _, item := range displayItems {
		line := fmt.Sprintf("%-15s %s", item.label, item.value)
		if s.SettingsIdx == item.idx {
			col1 = append(col1, selectedStyle.Render(line))
		} else {
			col1 = append(col1, itemStyle.Render("  "+line))
		}
	}

	var col2 []string
	col2 = append(col2, headerStyle.Render("TABS & APPEARANCE"))

	allTabs := []string{"Overview", "Metrics", "Processes", "Disks", "Network", "System"}
	currentTabIdxBase := 8

	for i, tabName := range allTabs {
		idx := currentTabIdxBase + i
		isEnabled := false
		for _, tab := range s.ActiveTabs {
			if tab == tabName {
				isEnabled = true
				break
			}
		}

		status := "[ ]"
		if isEnabled {
			status = "[x]"
		}

		line := fmt.Sprintf("%-15s %s", tabName, status)
		if s.SettingsIdx == idx {
			col2 = append(col2, selectedStyle.Render(line))
		} else {
			col2 = append(col2, itemStyle.Render("  "+line))
		}
	}

	col2 = append(col2, "")

	appearanceIdxBase := currentTabIdxBase + len(allTabs)

	bgLabel := "transparent"
	if s.BackgroundOpaque {
		bgLabel = "opaque"
	}

	appItems := []struct {
		label string
		value string
		idx   int
	}{
		{"Theme:", s.Theme, appearanceIdxBase},
		{"Refresh Rate:", fmt.Sprintf("%dms", s.RefreshRate), appearanceIdxBase + 1},
		{"Border Type:", s.BorderType, appearanceIdxBase + 2},
		{"Border Style:", s.BorderStyle, appearanceIdxBase + 3},
		{"Background:", bgLabel, appearanceIdxBase + 4},
	}

	for _, item := range appItems {
		line := fmt.Sprintf("%-15s %s", item.label, item.value)
		if s.SettingsIdx == item.idx {
			col2 = append(col2, selectedStyle.Render(line))
		} else {
			col2 = append(col2, itemStyle.Render("  "+line))
		}
	}

	col1Content := lipgloss.JoinVertical(lipgloss.Left, col1...)
	col2Content := lipgloss.JoinVertical(lipgloss.Left, col2...)

	contentBlock := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(boxWidth/2).Render(col1Content),
		lipgloss.NewStyle().Width(boxWidth/2).Render(col2Content),
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

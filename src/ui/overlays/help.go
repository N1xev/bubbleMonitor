package overlays

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"

	"github.com/N1xev/bubbleMonitor/src/data"
	"github.com/N1xev/bubbleMonitor/src/ui/widgets"
)

// RenderHelp renders the help overlay
func RenderHelp(s *data.AppState, b, p, bg compat.AdaptiveColor) string {
	sec := lipgloss.NewStyle().Foreground(compat.AdaptiveColor{Light: lipgloss.Color("#1E40AF"), Dark: lipgloss.Color("#3B82F6")}).Bold(true).MarginTop(1)
	key := lipgloss.NewStyle().Foreground(compat.AdaptiveColor{Light: lipgloss.Color("#047857"), Dark: lipgloss.Color("#10B981")}).Bold(true)
	desc := lipgloss.NewStyle().Foreground(compat.AdaptiveColor{Light: lipgloss.Color("#6B7280"), Dark: lipgloss.Color("#9CA3AF")})
	spacer := lipgloss.NewStyle()
	border := widgets.GetBorder(s.BorderStyle, s.BorderType)

	useTwoColumns := s.Width >= 100
	isCompact := s.Width < 60

	var boxWidth int
	if useTwoColumns {
		boxWidth = int(float64(s.Width) * 0.7)
		if boxWidth > 100 {
			boxWidth = 100
		}
	} else {
		boxWidth = int(float64(s.Width) * 0.8)
		if boxWidth > 70 {
			boxWidth = 70
		}
	}
	if boxWidth < 40 {
		boxWidth = 40
	}
	if boxWidth > s.Width-4 {
		boxWidth = s.Width - 4
	}

	sp := spacer.Render

	var helpContent string

	if useTwoColumns {
		colWidth := (boxWidth - 12) / 2
		if colWidth < 25 {
			colWidth = 25
		}

		leftCol := lipgloss.JoinVertical(lipgloss.Left,
			sec.Width(colWidth).Render("NAVIGATION"),
			spacer.Width(colWidth).Render(key.Render("Tab / → / L")+sp("  ")+desc.Render("Next tab")),
			spacer.Width(colWidth).Render(key.Render("Shift+Tab / ← / h")+sp("    ")+desc.Render("Previous tab")),
			spacer.Width(colWidth).Render(key.Render("1-6")+sp("  ")+desc.Render("Jump to tab")),
			spacer.Width(colWidth).Render(""),
			sec.Width(colWidth).Render("CONTROLS"),
			spacer.Width(colWidth).Render(key.Render("P")+sp(" ")+desc.Render("Pause/Resume")),
			spacer.Width(colWidth).Render(key.Render("R")+sp(" ")+desc.Render("Refresh")),
			spacer.Width(colWidth).Render(key.Render("S")+sp(" ")+desc.Render("Sort processes")),
			spacer.Width(colWidth).Render(key.Render("?")+sp(" ")+desc.Render("Toggle help")),
			spacer.Width(colWidth).Render(key.Render("Q")+sp(" ")+desc.Render("Quit")),
			spacer.Width(colWidth).Render(key.Render(".")+sp(" ")+desc.Render("Settings")),
			spacer.Width(colWidth).Render(key.Render("H")+sp(" ")+desc.Render("History len")),
			spacer.Width(colWidth).Render(key.Render("C")+sp(" ")+desc.Render("Chart type")),
		)

		rightCol := lipgloss.JoinVertical(lipgloss.Left,
			sec.Width(colWidth).Render("PROCESSES TAB"),
			spacer.Width(colWidth).Render(key.Render("j / ↓")+sp(" ")+desc.Render("Move down")),
			spacer.Width(colWidth).Render(key.Render("k / ↑")+sp(" ")+desc.Render("Move up")),
			spacer.Width(colWidth).Render(key.Render("g / G")+sp(" ")+desc.Render("Top / bottom")),
			spacer.Width(colWidth).Render(key.Render("f")+sp("     ")+desc.Render("Filter")),
			spacer.Width(colWidth).Render(key.Render("c")+sp("     ")+desc.Render("Clear filter")),
			spacer.Width(colWidth).Render(key.Render("z / x")+sp(" ")+desc.Render("Suspend/Resume")),
			spacer.Width(colWidth).Render(key.Render("K")+sp("     ")+desc.Render("Kill process")),
			spacer.Width(colWidth).Render(key.Render("o")+sp("     ")+desc.Render("Open files")),
			spacer.Width(colWidth).Render(key.Render("T")+sp("     ")+desc.Render("Tree view")),
			spacer.Width(colWidth).Render(key.Render("Space")+sp(" ")+desc.Render("Collapse/Exp")),
			spacer.Width(colWidth).Render(key.Render("+ / -")+sp(" ")+desc.Render("Nice +/-")),
			spacer.Width(colWidth).Render(""),
			spacer.Width(colWidth).Render(""),
			lipgloss.NewStyle().Foreground(compat.AdaptiveColor{Light: lipgloss.Color("#6B7280"), Dark: lipgloss.Color("#9CA3AF")}).Italic(true).Width(colWidth).Render("Press ? or ESC to close"),
		)

		helpContent = lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().Width(colWidth).Render(leftCol),
			lipgloss.NewStyle().Width(4).Render(""),
			lipgloss.NewStyle().Width(colWidth).Render(rightCol),
		)
	} else if isCompact {
		contentWidth := boxWidth - 8
		if contentWidth < 20 {
			contentWidth = 20
		}
		helpContent = lipgloss.JoinVertical(lipgloss.Left,
			sec.Width(contentWidth).Render("NAVIGATION"),
			spacer.Width(contentWidth).Render(key.Render("Tab/→/L")+" "+desc.Render("Next")),
			spacer.Width(contentWidth).Render(key.Render("S-Tab/←")+" "+desc.Render("Prev")),
			spacer.Width(contentWidth).Render(key.Render("1-5")+" "+desc.Render("Jump to tab")),
			spacer.Width(contentWidth).Render(""),
			sec.Width(contentWidth).Render("CONTROLS"),
			spacer.Width(contentWidth).Render(key.Render("P")+" "+desc.Render("Pause")),
			spacer.Width(contentWidth).Render(key.Render("R")+" "+desc.Render("Refresh")),
			spacer.Width(contentWidth).Render(key.Render("S")+" "+desc.Render("Sort")),
			spacer.Width(contentWidth).Render(key.Render("?")+" "+desc.Render("Help")),
			spacer.Width(contentWidth).Render(key.Render("Q")+" "+desc.Render("Quit")),
			spacer.Width(contentWidth).Render(key.Render("H/C")+" "+desc.Render("History/Charts")),
			spacer.Width(contentWidth).Render(""),
			lipgloss.NewStyle().Foreground(compat.AdaptiveColor{Light: lipgloss.Color("#6B7280"), Dark: lipgloss.Color("#9CA3AF")}).Italic(true).Width(contentWidth).Render("? or ESC to close"),
		)
	} else {
		contentWidth := boxWidth - 8
		if contentWidth < 20 {
			contentWidth = 20
		}
		helpContent = lipgloss.JoinVertical(lipgloss.Left,
			sec.Width(contentWidth).Render("NAVIGATION"),
			spacer.Width(contentWidth).Render(key.Render("Tab / → / L")+sp("       ")+desc.Render("Next tab")),
			spacer.Width(contentWidth).Render(key.Render("Shift+Tab / ←")+sp("     ")+desc.Render("Previous tab")),
			spacer.Width(contentWidth).Render(key.Render("1-5")+sp("               ")+desc.Render("Jump to specific tab")),
			spacer.Width(contentWidth).Render(""),
			sec.Width(contentWidth).Render("CONTROLS"),
			spacer.Width(contentWidth).Render(key.Render("P")+sp("   ")+desc.Render("Pause/Resume monitoring")),
			spacer.Width(contentWidth).Render(key.Render("R")+sp("   ")+desc.Render("Refresh all data")),
			spacer.Width(contentWidth).Render(key.Render("S")+sp("   ")+desc.Render("Sort processes (CPU/Memory/PID)")),
			spacer.Width(contentWidth).Render(key.Render("?")+sp("   ")+desc.Render("Show/Hide this help")),
			spacer.Width(contentWidth).Render(key.Render("Q")+sp("   ")+desc.Render("Quit application")),
			spacer.Width(contentWidth).Render(key.Render("H")+sp("   ")+desc.Render("Cycle history length (1m/5m/15m/1h)")),
			spacer.Width(contentWidth).Render(key.Render("C")+sp("   ")+desc.Render("Cycle chart type (Spark/Line/Bar)")),
			spacer.Width(contentWidth).Render(""),
			sec.Width(contentWidth).Render("PROCESSES TAB"),
			spacer.Width(contentWidth).Render(key.Render("j / ↓")+sp("   ")+desc.Render("Move down")),
			spacer.Width(contentWidth).Render(key.Render("k / ↑")+sp("   ")+desc.Render("Move up")),
			spacer.Width(contentWidth).Render(key.Render("g / G")+sp("   ")+desc.Render("Go to top / bottom")),
			spacer.Width(contentWidth).Render(key.Render("f")+sp("       ")+desc.Render("Filter processes")),
			spacer.Width(contentWidth).Render(key.Render("c")+sp("       ")+desc.Render("Clear filter")),
			spacer.Width(contentWidth).Render(key.Render("K")+sp("       ")+desc.Render("Kill selected process")),
			spacer.Width(contentWidth).Render(key.Render("o")+sp("       ")+desc.Render("Open files")),
			spacer.Width(contentWidth).Render(key.Render("T")+sp("       ")+desc.Render("Toggle tree view")),
			spacer.Width(contentWidth).Render(key.Render("Space")+sp("   ")+desc.Render("Collapse/Expand tree node")),
			spacer.Width(contentWidth).Render(key.Render("+ / -")+sp("   ")+desc.Render("Increase/Decrease priority")),
			spacer.Width(contentWidth).Render(""),
			lipgloss.NewStyle().Foreground(compat.AdaptiveColor{Light: lipgloss.Color("#6B7280"), Dark: lipgloss.Color("#9CA3AF")}).Italic(true).Width(contentWidth).Render("Press ? or ESC to close"),
		)
	}

	var boxHeight int
	if useTwoColumns {
		boxHeight = 18
	} else if isCompact {
		boxHeight = 18
	} else {
		boxHeight = 32
	}
	maxHeight := int(float64(s.Height) * 0.8)
	if boxHeight > maxHeight {
		boxHeight = maxHeight
	}
	if boxHeight < 10 {
		boxHeight = 10
	}

	container := lipgloss.NewStyle().
		Border(border).
		BorderForeground(b).
		Padding(1, 2).
		Width(boxWidth - 6).
		Height(boxHeight).
		BorderTop(false)

	body := container.Render(helpContent)
	actualWidth := lipgloss.Width(body)

	topBorder := widgets.RenderTopBorderWithBg("HELP", actualWidth, border, b, p)

	helpBox := lipgloss.JoinVertical(lipgloss.Left, topBorder, body)

	return helpBox
}

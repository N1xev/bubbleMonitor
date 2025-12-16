package ui

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/N1xev/bubbleMonitor/src/config"
	"github.com/N1xev/bubbleMonitor/src/data"
	"github.com/N1xev/bubbleMonitor/src/ui/overlays"
	"github.com/N1xev/bubbleMonitor/src/ui/tabs"
	"github.com/N1xev/bubbleMonitor/src/ui/widgets"
)

// ViewModel is an interface for rendering the UI
// This avoids import cycles by not importing model package
type ViewModel interface {
	GetBorder() lipgloss.Border
	GetColors() ThemePalette
	GetVisibleProcesses() ([]data.ProcessInfo, map[int32]int)
	// Direct field accessors from AppState (embedded in Model)
	GetAppState() *data.AppState
}

// MainViewFromState renders the entire application UI from AppState
func MainViewFromState(s *data.AppState, getBorder func() lipgloss.Border, getColors func() ThemePalette) tea.View {
	// Minimum dimensions check
	const minWidth = 100
	const minHeight = 30
	if s.Width < minWidth || s.Height < minHeight {
		theme := getColors()
		p := theme.Primary
		t := theme.Text
		bg := theme.Background

		boxStyle := lipgloss.NewStyle().
			Border(getBorder()).
			BorderForeground(p).
			Padding(1, 2).
			Margin(1).
			Foreground(t).
			Align(lipgloss.Center)

		titleStyle := lipgloss.NewStyle().
			Foreground(p).
			Bold(true).
			MarginBottom(1)

		dimStyle := lipgloss.NewStyle().
			Foreground(theme.Muted)

		msg := fmt.Sprintf(
			"%s\n\n%s\n%s",
			titleStyle.Render("WINDOW TOO SMALL"),
			fmt.Sprintf("Current: %dx%d", s.Width, s.Height),
			dimStyle.Render(fmt.Sprintf("Minimum: %dx%d", minWidth, minHeight)),
		)

		v := tea.NewView(lipgloss.Place(
			s.Width, s.Height,
			lipgloss.Center, lipgloss.Center,
			boxStyle.Render(msg),
		))

		v.AltScreen = true

		if s.BackgroundOpaque {
			v.BackgroundColor = bg
		} else {
			v.BackgroundColor = lipgloss.NoColor{}
		}

		return v
	}

	// Theme colors from current theme
	theme := getColors()
	t := theme.Text
	mu := theme.Muted
	bg := theme.Background
	p := theme.Primary
	b := theme.Border
	// su := theme.Success
	// w := theme.Warning
	a := theme.Alert

	// Border style
	border := getBorder()

	// Render Header with top margin
	headerText := " BUBBLE MONITOR"

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(p).
		MarginTop(1).
		Render(headerText) + lipgloss.NewStyle().Foreground(mu).Render("  /////  ") +
		lipgloss.NewStyle().Foreground(mu).Render(time.Now().Format("15:04:05"))

	// Create Alert String
	var alertStr string
	if s.AlertManager != nil && len(s.AlertManager.ActiveAlerts) > 0 {
		warnStyle := lipgloss.NewStyle().Foreground(a).Bold(true).Blink(true)
		rawText := "  ⚠️  ALERT: "

		// Append first alert message
		for _, alert := range s.AlertManager.ActiveAlerts {
			rawText += alert.Message
			break
		}
		alertStr = warnStyle.Render(rawText)
	}

	// Render Tabs with top margin
	var tabBlocks []string
	for i, titleRaw := range s.ActiveTabs {
		title := strings.ToUpper(titleRaw)
		if s.SelectedTab == i {
			tabBlocks = append(tabBlocks, lipgloss.NewStyle().
				Bold(true).
				Foreground(p).
				Border(border, false, false, true, false).
				BorderForeground(p).
				Render(" "+title+" "))
		} else {
			tabBlocks = append(tabBlocks, lipgloss.NewStyle().
				Foreground(mu).
				MarginBottom(1).
				Render(" "+title+" "))
		}
	}
	tabRow := lipgloss.NewStyle().MarginTop(1).Render(lipgloss.JoinHorizontal(lipgloss.Bottom, tabBlocks...))

	// TopBar Assembly
	var topBar string
	if s.Width >= 130 && alertStr != "" {
		alertBlock := lipgloss.NewStyle().MarginTop(1).Render(alertStr)
		topBar = lipgloss.JoinHorizontal(lipgloss.Top, header, "    ", tabRow, alertBlock)
	} else {
		topBar = lipgloss.JoinHorizontal(lipgloss.Top, header, "    ", tabRow)
	}

	// Footer - context-aware
	var footerText string
	if s.SelectedTab == 2 {
		if s.FilterMode {
			footerText = "Type to filter • ESC/Return to apply filter"
		} else {
			footerText = "Press ? for Help • f to Filter • K to Kill • S to Sort"
		}
	} else {
		footerText = "Press ? for Help • q to Quit"
	}

	// Footer Assembly
	var footer string
	if s.Width < 130 && alertStr != "" {
		footerLeft := lipgloss.NewStyle().Foreground(mu).Render(footerText)
		footerContent :=  lipgloss.JoinHorizontal(lipgloss.Bottom, footerLeft, lipgloss.NewStyle().Foreground(mu).Render("  /////  "), alertStr)
		footer = lipgloss.NewStyle().MarginBottom(1).Render(footerContent)
	} else {
		footer = lipgloss.NewStyle().
			Foreground(mu).
			MarginBottom(1).
			Render(footerText)
	}

	// Calculate Content Area Height
	topGap := 1
	topPad := lipgloss.NewStyle().Height(topGap).Render("")

	topBarH := lipgloss.Height(topBar)
	footerH := lipgloss.Height(footer)

	reservedHeight := topBarH + topGap + footerH
	availHeight := s.Height - reservedHeight
	if availHeight < 5 {
		availHeight = 5
	}
	// Determine active tab
	activeTab := ""
	if s.SelectedTab >= 0 && s.SelectedTab < len(s.ActiveTabs) {
		activeTab = s.ActiveTabs[s.SelectedTab]
	}

	// Styles for tabs
	titleStyle := lipgloss.NewStyle().Foreground(p).Bold(true).MarginBottom(1)
	labelStyle := lipgloss.NewStyle().Foreground(mu)
	valueStyle := lipgloss.NewStyle().Foreground(t).Bold(true)
	su := theme.Success
	w := theme.Warning
	sColor := theme.Secondary

	var content string
	container := lipgloss.NewStyle().
		Border(border).
		BorderForeground(b).
		Padding(0, 1)

	switch activeTab {
	case "Overview":
		content = tabs.RenderOverview(s, container, titleStyle, labelStyle, valueStyle, su, w, a, t, mu, bg, p, b, availHeight)
	case "Metrics":
		content = tabs.RenderMetrics(s, container, su, w, a, sColor, t, mu, p, b, availHeight)
	case "Processes":
		// Get visible processes
		visibleProcs, treeIndents := s.GetVisibleProcesses()
		content = tabs.RenderProcesses(s, visibleProcs, treeIndents, container, su, w, a, t, mu, p, b, availHeight)
	case "Disks":
		content = tabs.RenderDisks(s, container, su, w, a, t, mu, p, b, availHeight)
	case "Network":
		content = tabs.RenderNetwork(s, container, titleStyle, labelStyle, valueStyle, t, mu, p, b, bg, availHeight)
	case "System":
		content = tabs.RenderSystem(s, container, titleStyle, labelStyle, valueStyle, t, mu, p, b, bg, availHeight)
	default:
		content = lipgloss.NewStyle().Foreground(mu).Render("Tab not found: " + activeTab)
	}

	ui := lipgloss.JoinVertical(lipgloss.Left,
		topBar,
		topPad,
		content,
	)

	// Determine how much space left for footer to stick to bottom
	currH := lipgloss.Height(ui)

	rem := s.Height - currH - footerH
	if rem < 0 {
		rem = 0
	}

	ui = lipgloss.JoinVertical(lipgloss.Left,
		ui,
		strings.Repeat("\n", rem),
		footer,
	)

	// Create base layer
	baseLayer := lipgloss.NewLayer(ui).Width(s.Width).Height(s.Height)

	var layers []*lipgloss.Layer
	layers = append(layers, baseLayer)

	// Toast Overlay
	if len(s.Toasts) > 0 {
		var toastBlocks []string
		for _, toast := range s.Toasts {
			color := "#10B981" // Green
			icon := "✔ "
			if toast.Level == "error" {
				color = "#EF4444"
				icon = "✖ "
			}
			if toast.Level == "warn" {
				color = "#F59E0B"
				icon = "⚠"
			}

			block := lipgloss.NewStyle().
				Border(border).
				BorderForeground(lipgloss.Color(color)).
				Padding(0, 1).
				Render(lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(icon + " " + toast.Message))
			toastBlocks = append(toastBlocks, block)
		}
		toastStack := lipgloss.JoinVertical(lipgloss.Right, toastBlocks...)
		toastWidth := lipgloss.Width(toastStack)
		toastHeight := lipgloss.Height(toastStack)

		toastX := s.Width - toastWidth - 2
		toastY := s.Height - footerH - toastHeight - 1

		if toastX < 0 {
			toastX = 0
		}
		if toastY < 0 {
			toastY = 0
		}

		toastLayer := lipgloss.NewLayer(toastStack).X(toastX).Y(toastY).Z(2)
		layers = append(layers, toastLayer)
	}

	if s.ShowKillDialog {
		killDialog := overlays.RenderKillDialog(s, b, p, a, t, mu)
		dialogWidth := lipgloss.Width(killDialog)
		dialogHeight := lipgloss.Height(killDialog)

		dialogX := (s.Width - dialogWidth) / 2
		dialogY := (s.Height - dialogHeight) / 2
		if dialogX < 0 {
			dialogX = 0
		}
		if dialogY < 0 {
			dialogY = 0
		}

		dialogLayer := lipgloss.NewLayer(killDialog).X(dialogX).Y(dialogY).Z(3)
		layers = append(layers, dialogLayer)
	}

	if s.ShowHelp {
		helpBox := overlays.RenderHelp(s, b, p, bg)
		hWidth := lipgloss.Width(helpBox)
		hHeight := lipgloss.Height(helpBox)
		hX := (s.Width - hWidth) / 2
		hY := (s.Height - hHeight) / 2
		if hX < 0 {
			hX = 0
		}
		if hY < 0 {
			hY = 0
		}
		layers = append(layers, lipgloss.NewLayer(helpBox).X(hX).Y(hY).Z(4))
	}

	if s.ShowSettings {
		settingsBox := overlays.RenderSettingsOverlay(s, s.Width, s.Height, b, p, t, mu, bg)
		sWidth := lipgloss.Width(settingsBox)
		sHeight := lipgloss.Height(settingsBox)
		sX := (s.Width - sWidth) / 2
		sY := (s.Height - sHeight) / 2
		if sX < 0 {
			sX = 0
		}
		if sY < 0 {
			sY = 0
		}
		layers = append(layers, lipgloss.NewLayer(settingsBox).X(sX).Y(sY).Z(5))
	}

	if s.ShowOpenFiles {
		filesBox := overlays.RenderOpenFilesOverlay(s, s.Width, s.Height, b, p, t, mu, bg)
		fWidth := lipgloss.Width(filesBox)
		fHeight := lipgloss.Height(filesBox)
		fX := (s.Width - fWidth) / 2
		fY := (s.Height - fHeight) / 2
		if fX < 0 {
			fX = 0
		}
		if fY < 0 {
			fY = 0
		}
		layers = append(layers, lipgloss.NewLayer(filesBox).X(fX).Y(fY).Z(4))
	}

	// Create view from layers
	canvas := lipgloss.NewCanvas(layers...)
	v := tea.NewView(canvas)
	v.AltScreen = true

	// Apply background based on transparency setting
	if s.BackgroundOpaque {
		v.BackgroundColor = bg
	} else {
		v.BackgroundColor = lipgloss.NoColor{}
	}
	return v
}

// Helper function for rendering - gets theme from state
func getThemeFromState(s *data.AppState) ThemePalette {
	return GetAppTheme(s.Theme, s.Config.CustomTheme)
}

// Helper function for rendering - gets border from state
func getBorderFromState(s *data.AppState) lipgloss.Border {
	return widgets.GetBorder(s.BorderStyle, s.BorderType)
}

// RenderFromAppState is a convenience wrapper for main.go to use
func RenderFromAppState(s *data.AppState) tea.View {
	return MainViewFromState(s, func() lipgloss.Border {
		return getBorderFromState(s)
	}, func() ThemePalette {
		return getThemeFromState(s)
	})
}

// CustomThemeConfig re-export for backward compatibility
type CustomThemeConfig = config.CustomThemeConfig

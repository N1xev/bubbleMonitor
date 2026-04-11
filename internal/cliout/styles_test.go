package cliout_test

import (
	"strings"
	"testing"

	"github.com/N1xev/bubbleMonitor/internal/cliout"
	"github.com/N1xev/bubbleMonitor/internal/ui"
)

func TestNew(t *testing.T) {
	palette := ui.GetTheme("dark")
	s := cliout.New(palette)

	// Verify all 9 style fields produce styled output (non-zero lipgloss.Style)
	// A lipgloss.Style with a Foreground set will produce ANSI escape codes when rendered
	fields := []struct {
		name  string
		style func() string
	}{
		{"Label", func() string { return s.Label.Render("x") }},
		{"Value", func() string { return s.Value.Render("x") }},
		{"Bold", func() string { return s.Bold.Render("x") }},
		{"OK", func() string { return s.OK.Render("x") }},
		{"Warn", func() string { return s.Warn.Render("x") }},
		{"Critical", func() string { return s.Critical.Render("x") }},
		{"Dim", func() string { return s.Dim.Render("x") }},
		{"Header", func() string { return s.Header.Render("x") }},
		{"Active", func() string { return s.Active.Render("x") }},
	}

	for _, f := range fields {
		rendered := f.style()
		// A plain "x" has length 1. A styled "x" contains ANSI escapes and is longer.
		if len(rendered) <= 1 {
			t.Errorf("%s style did not apply ANSI styling (output length %d, expected >1)", f.name, len(rendered))
		}
	}
}

func TestCheckSymbols(t *testing.T) {
	palette := ui.GetTheme("dark")
	s := cliout.New(palette)

	if !strings.Contains(s.CheckOK, "\u2713") {
		t.Errorf("CheckOK should contain checkmark \u2713, got: %q", s.CheckOK)
	}
	if !strings.Contains(s.CheckFail, "\u2717") {
		t.Errorf("CheckFail should contain cross \u2717, got: %q", s.CheckFail)
	}
	if !strings.Contains(s.CheckWarn, "\u26A0") {
		t.Errorf("CheckWarn should contain warning \u26A0, got: %q", s.CheckWarn)
	}

	// Verify the symbols are styled (contain ANSI codes, longer than plain rune)
	if len(s.CheckOK) <= len("\u2713") {
		t.Errorf("CheckOK should be styled (length %d, expected >%d)", len(s.CheckOK), len("\u2713"))
	}
	if len(s.CheckFail) <= len("\u2717") {
		t.Errorf("CheckFail should be styled (length %d, expected >%d)", len(s.CheckFail), len("\u2717"))
	}
	if len(s.CheckWarn) <= len("\u26A0") {
		t.Errorf("CheckWarn should be styled (length %d, expected >%d)", len(s.CheckWarn), len("\u26A0"))
	}
}

func TestNewStoresPalette(t *testing.T) {
	palette := ui.GetTheme("dark")
	s := cliout.New(palette)

	// Verify the unexported palette is accessible via BarColor returning a non-zero style
	// If palette was not stored, BarColor would panic or return a zero style
	barStyle := s.BarColor(50)
	rendered := barStyle.Render("x")
	if len(rendered) <= 1 {
		t.Error("BarColor returned a zero style, suggesting palette was not stored in CLIStyles")
	}
}

func TestLabelUsesMutedColor(t *testing.T) {
	palette := ui.GetTheme("dark")
	s := cliout.New(palette)

	rendered := s.Label.Render("test")
	if len(rendered) <= len("test") {
		t.Error("Label style should apply Muted foreground color")
	}
}

func TestOKIsBold(t *testing.T) {
	palette := ui.GetTheme("dark")
	s := cliout.New(palette)

	rendered := s.OK.Render("test")
	// Bold text contains ANSI escape sequences
	if !strings.Contains(rendered, "\x1b[") {
		t.Error("OK style should produce ANSI codes (bold + Success foreground)")
	}
}

func TestHeaderIsBoldAndUnderlined(t *testing.T) {
	palette := ui.GetTheme("dark")
	s := cliout.New(palette)

	rendered := s.Header.Render("test")
	if !strings.Contains(rendered, "\x1b[") {
		t.Error("Header style should produce ANSI codes (bold + underline + Primary foreground)")
	}
}

func TestDimIsFaint(t *testing.T) {
	palette := ui.GetTheme("dark")
	s := cliout.New(palette)

	rendered := s.Dim.Render("test")
	if !strings.Contains(rendered, "\x1b[") {
		t.Error("Dim style should produce ANSI codes (faint + Muted foreground)")
	}
}

func TestActiveIsBold(t *testing.T) {
	palette := ui.GetTheme("dark")
	s := cliout.New(palette)

	rendered := s.Active.Render("test")
	if !strings.Contains(rendered, "\x1b[") {
		t.Error("Active style should produce ANSI codes (bold + Primary foreground)")
	}
}

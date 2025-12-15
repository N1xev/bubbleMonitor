package data

import (
	"strings"
)

// SimpleViewport is a lightweight scrollable text viewer
type SimpleViewport struct {
	Width   int
	Height  int
	YOffset int
	Content []string // Split lines
}

func NewSimpleViewport(width, height int) SimpleViewport {
	return SimpleViewport{
		Width:  width,
		Height: height,
	}
}

func (v *SimpleViewport) SetContent(s string) {
	v.Content = strings.Split(s, "\n")
	v.YOffset = 0
}

func (v *SimpleViewport) View() string {
	if len(v.Content) == 0 {
		return ""
	}
	start := v.YOffset
	end := start + v.Height
	if start >= len(v.Content) {
		start = len(v.Content) - 1
		if start < 0 {
			start = 0
		}
	}
	if end > len(v.Content) {
		end = len(v.Content)
	}

	visible := v.Content[start:end]

	return strings.Join(visible, "\n")
}

func (v *SimpleViewport) LineDown(n int) {
	maxOffset := len(v.Content) - v.Height
	if maxOffset < 0 {
		maxOffset = 0
	}
	v.YOffset += n
	if v.YOffset > maxOffset {
		v.YOffset = maxOffset
	}
}

func (v *SimpleViewport) LineUp(n int) {
	v.YOffset -= n
	if v.YOffset < 0 {
		v.YOffset = 0
	}
}

func (v *SimpleViewport) GotoTop() {
	v.YOffset = 0
}

func (v *SimpleViewport) GotoBottom() {
	maxOffset := len(v.Content) - v.Height
	if maxOffset < 0 {
		maxOffset = 0
	}
	v.YOffset = maxOffset
}

func (v *SimpleViewport) HalfViewDown() {
	v.LineDown(v.Height / 2)
}

func (v *SimpleViewport) HalfViewUp() {
	v.LineUp(v.Height / 2)
}

package data

import (
	"charm.land/bubbles/v2/viewport"
)

// SimpleViewport wraps Bubbles viewport.Model to provide the same API
type SimpleViewport struct {
	viewport viewport.Model
}

func NewSimpleViewport(width, height int) SimpleViewport {
	vp := viewport.New(
		viewport.WithWidth(width),
		viewport.WithHeight(height),
	)
	return SimpleViewport{viewport: vp}
}

func (v *SimpleViewport) SetContent(s string) {
	v.viewport.SetContent(s)
}

func (v *SimpleViewport) View() string {
	return v.viewport.View()
}

func (v *SimpleViewport) LineDown(n int) {
	v.viewport.ScrollDown(n)
}

func (v *SimpleViewport) LineUp(n int) {
	v.viewport.ScrollUp(n)
}

func (v *SimpleViewport) GotoTop() {
	v.viewport.GotoTop()
}

func (v *SimpleViewport) GotoBottom() {
	v.viewport.GotoBottom()
}

func (v *SimpleViewport) HalfViewDown() {
	v.viewport.HalfPageDown()
}

func (v *SimpleViewport) HalfViewUp() {
	v.viewport.HalfPageUp()
}

func (v *SimpleViewport) SetWidth(w int) {
	v.viewport.SetWidth(w)
}

func (v *SimpleViewport) SetHeight(h int) {
	v.viewport.SetHeight(h)
}

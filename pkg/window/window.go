package window

import (
	"github.com/AllenDang/giu"
	"github.com/gucio321/HellSpawner/pkg/app/state"
)

// Window represents project's window
type Window struct {
	*giu.WindowWidget
	Visible bool
}

// New creates new window
func New(title string, x, y float32) *Window {
	return (&Window{
		WindowWidget: giu.Window(title),
	}).Pos(x, y)
}

// State returns window's state
func (t *Window) State() state.WindowState {
	x, y := t.CurrentPosition()
	w, h := t.CurrentSize()

	return state.WindowState{
		Visible: t.Visible,
		PosX:    x,
		PosY:    y,
		Width:   w,
		Height:  h,
	}
}

// ToggleVisibility toggles visibility
func (t *Window) ToggleVisibility() {
	t.Visible = !t.Visible
}

// Show turns visibility to true
func (t *Window) Show() {
	t.Visible = true
}

// Build builds window
func (t *Window) Build() {
	// noop
}

// RegisterKeyboardShortcuts sets a local shortcuts list
func (t *Window) RegisterKeyboardShortcuts(s ...giu.WindowShortcut) {
	t.WindowWidget.RegisterKeyboardShortcuts(s...)
}

// KeyboardShortcuts returns a list of local keyboard shortcuts
func (t *Window) KeyboardShortcuts() []giu.WindowShortcut {
	return []giu.WindowShortcut{}
}

// IsVisible returns true if window is visible
func (t *Window) IsVisible() bool {
	return t.Visible
}

// SetVisible sets window's visibility
func (t *Window) SetVisible(visible bool) {
	t.Visible = visible
}

// Cleanup hides window
func (t *Window) Cleanup() {
	t.Visible = false
}

func (t *Window) Pos(x, y float32) *Window {
	t.WindowWidget.Pos(x, y)
	return t
}

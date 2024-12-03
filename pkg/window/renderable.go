package window

import "github.com/AllenDang/giu"

// Renderable represents top-level renderable objects (window of one of types: editor, toolwindow, dialog)
type Renderable interface {
	Build()
	Cleanup()
	// KeyboardShortcuts returns a list of keyboard shortcuts
	KeyboardShortcuts() []giu.WindowShortcut
	IsVisible() bool
	// RegisterKeyboardShortcuts wraps giu.RegisterKeyboardShortcuts
	RegisterKeyboardShortcuts(...giu.WindowShortcut)
	GetLayout() giu.Widget
}

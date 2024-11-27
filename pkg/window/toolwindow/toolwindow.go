package toolwindow

import (
	"github.com/gucio321/HellSpawner/pkg/common/state"
	"github.com/gucio321/HellSpawner/pkg/window"
)

// ToolWindow represents a tool window
type ToolWindow struct {
	*window.Window
	Type state.ToolWindowType
}

// New creates a new tool window
func New(title string, toolWindowType state.ToolWindowType, x, y float32) *ToolWindow {
	return &ToolWindow{
		Window: window.New(title, x, y),
		Type:   toolWindowType,
	}
}

// State returns state of tool window
func (t *ToolWindow) State() state.ToolWindowState {
	return state.ToolWindowState{
		WindowState: t.Window.State(),
		Type:        t.Type,
	}
}

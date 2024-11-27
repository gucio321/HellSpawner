package toolwindow

import (
	"github.com/gucio321/HellSpawner/pkg/common/hsstate"
	"github.com/gucio321/HellSpawner/pkg/window"
)

// ToolWindow represents a tool window
type ToolWindow struct {
	*window.Window
	Type hsstate.ToolWindowType
}

// New creates a new tool window
func New(title string, toolWindowType hsstate.ToolWindowType, x, y float32) *ToolWindow {
	return &ToolWindow{
		Window: window.New(title, x, y),
		Type:   toolWindowType,
	}
}

// State returns state of tool window
func (t *ToolWindow) State() hsstate.ToolWindowState {
	return hsstate.ToolWindowState{
		WindowState: t.Window.State(),
		Type:        t.Type,
	}
}

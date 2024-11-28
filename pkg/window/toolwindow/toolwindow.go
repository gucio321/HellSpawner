package toolwindow

import (
	"github.com/AllenDang/giu"
	"github.com/gucio321/HellSpawner/pkg/app/state"
	"github.com/gucio321/HellSpawner/pkg/window"
)

// ToolWindow represents tool windows
type ToolWindow interface {
	window.Renderable

	HasFocus() (hasFocus bool)
	Show()
	SetVisible(bool)
	BringToFront()
	State() state.ToolWindowState
	Pos(x, y float32) *window.Window
	Size(float32, float32) *giu.WindowWidget
	CurrentSize() (float32, float32)
}

// ToolWindowBase represents a base for a tool window
// Window implementations should embed this struct
type ToolWindowBase struct {
	*window.Window
	Type state.ToolWindowType
}

// New creates a new tool window
func New(title string, toolWindowType state.ToolWindowType, x, y float32) *ToolWindowBase {
	return &ToolWindowBase{
		Window: window.New(title, x, y),
		Type:   toolWindowType,
	}
}

// State returns state of tool window
func (t *ToolWindowBase) State() state.ToolWindowState {
	return state.ToolWindowState{
		WindowState: t.Window.State(),
		Type:        t.Type,
	}
}

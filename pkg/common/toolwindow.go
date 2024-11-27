package common

import (
	"github.com/AllenDang/giu"

	"github.com/gucio321/HellSpawner/pkg/common/hsstate"
)

// ToolWindow represents tool windows
type ToolWindow interface {
	Renderable

	HasFocus() (hasFocus bool)
	Show()
	SetVisible(bool)
	BringToFront()
	State() hsstate.ToolWindowState
	Pos(x, y float32) *giu.WindowWidget
	Size(float32, float32) *giu.WindowWidget
	CurrentSize() (float32, float32)
}
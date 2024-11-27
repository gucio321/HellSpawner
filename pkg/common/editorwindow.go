package common

import (
	"github.com/AllenDang/giu"

	"github.com/gucio321/HellSpawner/pkg/common/state"
)

// EditorWindow represents editor window
type EditorWindow interface {
	Renderable
	MainMenuUpdater

	// HasFocus returns true if editor is focused
	HasFocus() (hasFocus bool)
	// GetWindowTitle controls what the window title for this editor appears as
	GetWindowTitle() string
	// Show sets Visible to true
	Show()
	// SetVisible can be used to set Visible to false if the editor should be closed
	SetVisible(bool)
	// GetID returns a unique identifier for this editor window
	GetID() string
	// BringToFront brings this editor to the front of the application, giving it focus
	BringToFront()
	// State returns the current state of this editor, in a JSON-serializable struct
	State() state.EditorState
	// Save writes any changes made in the editor to the file that is open in the editor.
	Save()

	Size(float32, float32) *giu.WindowWidget
}

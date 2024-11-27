package editor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/AllenDang/giu"

	"github.com/gucio321/HellSpawner/pkg/common/hsproject"

	"github.com/gucio321/HellSpawner/pkg/common"
	"github.com/gucio321/HellSpawner/pkg/common/state"
	"github.com/gucio321/HellSpawner/pkg/window"
)

// EditorWindow represents editor window
type Editor interface {
	window.Renderable
	window.MainMenuUpdater

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

// EditorBase represents an editor
type EditorBase struct {
	*window.Window
	Path    *common.PathEntry
	Project *hsproject.Project
}

// New creates a new editor
func New(path *common.PathEntry, x, y float32, project *hsproject.Project) *EditorBase {
	return &EditorBase{
		Window:  window.New(generateWindowTitle(path), x, y),
		Path:    path,
		Project: project,
	}
}

// State returns editors state
func (e *EditorBase) State() state.EditorState {
	path, err := json.Marshal(e.Path)
	if err != nil {
		log.Print("failed to marshal editor path to JSON: ", err)
	}

	result := state.EditorState{
		WindowState: e.Window.State(),
		Path:        path,
		Encoded:     e.EncodeState(),
	}

	return result
}

// GetWindowTitle returns window title
func (e *EditorBase) GetWindowTitle() string {
	return generateWindowTitle(e.Path)
}

// GetID returns editors ID
func (e *EditorBase) GetID() string {
	return e.Path.GetUniqueID()
}

// Save saves an editor
func (e *EditorBase) Save(editor Saveable) {
	if e.Path.Source != common.PathEntrySourceProject {
		// saving to MPQ not yet supported
		return
	}

	saveData := editor.GenerateSaveData()
	if saveData == nil {
		return
	}

	existingFileData, err := e.Path.GetFileBytes()
	if err != nil {
		fmt.Println("failed to read file before saving: ", err)
		return
	}

	if bytes.Equal(saveData, existingFileData) {
		// nothing to save
		return
	}

	err = e.Path.WriteFile(saveData)
	if err != nil {
		fmt.Println("failed to save file: ", err)
		return
	}
}

// HasChanges returns true if editor has changed data
func (e *EditorBase) HasChanges(editor Saveable) bool {
	if e.Path.Source != common.PathEntrySourceProject {
		// saving to MPQ not yet supported
		return false
	}

	newData := editor.GenerateSaveData()
	if newData != nil {
		oldData, err := e.Path.GetFileBytes()
		if err == nil {
			return !bytes.Equal(oldData, newData)
		}
	}

	// err on the side of caution; if any errors occurred, just say nothing has changed so no changes get saved
	return false
}

// Cleanup cides an editor
func (e *EditorBase) Cleanup() {
	e.Window.Cleanup()
}

func generateWindowTitle(path *common.PathEntry) string {
	return path.Name + "##" + path.GetUniqueID()
}

// EncodeState returns widget's state (unique for each editor type) in byte slice format
func (e *EditorBase) EncodeState() []byte {
	id := giu.ID(fmt.Sprintf("widget_%s", e.Path.GetUniqueID()))

	if s := giu.Context.GetState(id); s != nil {
		data, err := json.Marshal(s)
		if err != nil {
			log.Printf("error encoding state of editor at path %v: %v", e.Path, err)
		}

		return data
	}

	return nil
}

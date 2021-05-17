// Package hscofeditor contains cof editor's data
package hscharactersaveeditor

import (
	"fmt"

	"github.com/OpenDiablo2/dialog"
	g "github.com/ianling/giu"

	"github.com/gucio321/d2d2s"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hswidget/charactersavewidget"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hsconfig"
	"github.com/OpenDiablo2/HellSpawner/hsinput"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

// static check, to ensure, if cof editor implemented editoWindow
var _ hscommon.EditorWindow = &CharacterSaveEditor{}

// CharacterSaveEditor represents a cof editor
type CharacterSaveEditor struct {
	*hseditor.Editor
	char  *d2d2s.D2S
	state []byte
}

// Create creates a new cof editor
func Create(config *hsconfig.Config,
	tl hscommon.TextureLoader,
	pathEntry *hscommon.PathEntry,
	state []byte,
	data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	char, err := d2d2s.Unmarshal(*data)
	if err != nil {
		return nil, fmt.Errorf("error loading cof file: %w", err)
	}

	result := &CharacterSaveEditor{
		Editor: hseditor.New(pathEntry, x, y, project),
		char:   char,
		state:  state,
	}

	return result, nil
}

// Build builds a cof editor
func (e *CharacterSaveEditor) Build() {
	uid := e.Path.GetUniqueID()
	widget := charactersavewidget.Create(e.state, uid, e.char)

	e.IsOpen(&e.Visible)
	e.Flags(g.WindowFlagsAlwaysAutoResize)
	e.Layout(g.Layout{widget})
}

// UpdateMainMenuLayout updates a main menu layout, to it contains COFViewer's settings
func (e *CharacterSaveEditor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("Character Save Editor").Layout(g.Layout{
		g.MenuItem("Save\t\t\t\tCtrl+Shift+S").OnClick(e.Save),
		g.Separator(),
		g.MenuItem("Add to project").OnClick(func() {}),
		g.MenuItem("Remove from project").OnClick(func() {}),
		g.Separator(),
		g.MenuItem("Import from file...").OnClick(func() {}),
		g.MenuItem("Export to file...").OnClick(func() {}),
		g.Separator(),
		g.MenuItem("Close").OnClick(func() {
			e.Cleanup()
		}),
	})

	*l = append(*l, m)
}

// RegisterKeyboardShortcuts adds a local shortcuts for this editor
func (e *CharacterSaveEditor) RegisterKeyboardShortcuts(inputManager *hsinput.InputManager) {
	// Ctrl+Shift+S saves file
	inputManager.RegisterShortcut(func() {
		e.Save()
	}, g.KeyS, g.ModShift+g.ModControl, false)
}

// GenerateSaveData generates data to be saved
func (e *CharacterSaveEditor) GenerateSaveData() []byte {
	data, _ := e.Path.GetFileBytes()

	return data
}

// Save saves an editor
func (e *CharacterSaveEditor) Save() {
	e.Editor.Save(e)
}

// Cleanup hides an editor
func (e *CharacterSaveEditor) Cleanup() {
	const strPrompt = "There are unsaved changes to %s, save before closing this editor?"

	if e.HasChanges(e) {
		if shouldSave := dialog.Message(strPrompt, e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.Editor.Cleanup()
}

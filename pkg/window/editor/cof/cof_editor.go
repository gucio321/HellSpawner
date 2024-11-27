// Package cof contains cof editor's data
package cof

import (
	"fmt"
	"github.com/gucio321/HellSpawner/pkg/app/config"

	g "github.com/AllenDang/giu"
	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2cof"

	"github.com/gucio321/HellSpawner/pkg/common/hsproject"

	"github.com/gucio321/HellSpawner/pkg/common"
	"github.com/gucio321/HellSpawner/pkg/window/editor"

	"github.com/gucio321/HellSpawner/pkg/widgets/cofwidget"
)

// static check, to ensure, if cof editor implemented editoWindow
var _ editor.Editor = &Editor{}

// Editor represents a cof editor
type Editor struct {
	*editor.EditorBase
	cof           *d2cof.COF
	textureLoader common.TextureLoader
	state         []byte
}

// Create creates a new cof editor
func Create(_ *config.Config,
	tl common.TextureLoader,
	pathEntry *common.PathEntry,
	state []byte,
	data *[]byte, x, y float32, project *hsproject.Project,
) (editor.Editor, error) {
	cof, err := d2cof.Unmarshal(*data)
	if err != nil {
		return nil, fmt.Errorf("error loading cof file: %w", err)
	}

	result := &Editor{
		EditorBase:    editor.New(pathEntry, x, y, project),
		cof:           cof,
		textureLoader: tl,
		state:         state,
	}

	return result, nil
}

// Build builds a cof editor
func (e *Editor) Build() {
	uid := e.Path.GetUniqueID()
	cofWidget := cofwidget.Create(e.state, e.textureLoader, uid, e.cof)

	e.IsOpen(&e.Visible)
	e.Flags(g.WindowFlagsAlwaysAutoResize)
	e.Layout(g.Layout{cofWidget})
}

// UpdateMainMenuLayout updates a main menu layout, to it contains COFViewer's settings
func (e *Editor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("COF Editor").Layout(g.Layout{
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

// GenerateSaveData generates data to be saved
func (e *Editor) GenerateSaveData() []byte {
	data := e.cof.Marshal()

	return data
}

// Save saves an editor
func (e *Editor) Save() {
	e.EditorBase.Save(e)
}

// Cleanup hides an editor
func (e *Editor) Cleanup() {
	const strPrompt = "There are unsaved changes to %s, save before closing this editor?"

	if e.HasChanges(e) {
		if shouldSave := dialog.Message(strPrompt, e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.EditorBase.Cleanup()
}

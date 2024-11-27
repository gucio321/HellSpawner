// Package ds1 contains ds1 editor's data
package ds1

import (
	"fmt"

	g "github.com/AllenDang/giu"
	"github.com/OpenDiablo2/dialog"

	"github.com/gucio321/HellSpawner/pkg/common/hsproject"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2ds1"

	"github.com/gucio321/HellSpawner/pkg/assets"
	"github.com/gucio321/HellSpawner/pkg/common"
	"github.com/gucio321/HellSpawner/pkg/config"
	"github.com/gucio321/HellSpawner/pkg/widgets/ds1widget"
	"github.com/gucio321/HellSpawner/pkg/window/editor"
)

// static check if Editor implemented editor.Editor
var _ editor.Editor = &Editor{}

// Editor represents ds1 editor
type Editor struct {
	*editor.EditorBase
	ds1                 *d2ds1.DS1
	deleteButtonTexture *g.Texture
	textureLoader       common.TextureLoader
	state               []byte
}

// Create creates a new ds1 editor
func Create(_ *config.Config,
	tl common.TextureLoader,
	pathEntry *common.PathEntry,
	state []byte,
	data *[]byte, x, y float32, project *hsproject.Project,
) (editor.Editor, error) {
	ds1, err := d2ds1.Unmarshal(*data)
	if err != nil {
		return nil, fmt.Errorf("error loading DS1 file: %w", err)
	}

	result := &Editor{
		EditorBase:    editor.New(pathEntry, x, y, project),
		ds1:           ds1,
		textureLoader: tl,
		state:         state,
	}

	result.Path = pathEntry

	tl.CreateTextureFromFile(assets.DeleteIcon, func(texture *g.Texture) {
		result.deleteButtonTexture = texture
	})

	return result, nil
}

// Build builds an editor
func (e *Editor) Build() {
	e.IsOpen(&e.Visible).
		Flags(g.WindowFlagsAlwaysAutoResize).
		Layout(g.Layout{
			ds1widget.Create(e.textureLoader, e.Path.GetUniqueID(), e.ds1, e.deleteButtonTexture, e.state),
		})
}

// UpdateMainMenuLayout updates main menu layout to it contains editors options
func (e *Editor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("DS1 Editor").Layout(g.Layout{
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
	data := e.ds1.Marshal()

	return data
}

// Save saves editors data
func (e *Editor) Save() {
	e.EditorBase.Save(e)
}

// Cleanup hides editor
func (e *Editor) Cleanup() {
	if e.HasChanges(e) {
		if shouldSave := dialog.Message("There are unsaved changes to %s, save before closing this editor?",
			e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.EditorBase.Cleanup()
}

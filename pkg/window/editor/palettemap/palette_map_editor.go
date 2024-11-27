// Package palettemap contains palette map editor's data
package palettemap

import (
	"fmt"
	"github.com/gucio321/HellSpawner/pkg/app/config"

	g "github.com/AllenDang/giu"

	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2pl2"

	"github.com/gucio321/HellSpawner/pkg/common"
	"github.com/gucio321/HellSpawner/pkg/common/hsproject"
	"github.com/gucio321/HellSpawner/pkg/widgets/palettemapwidget"
	"github.com/gucio321/HellSpawner/pkg/window/editor"
)

// static check, to ensure, if palette map editor implemented editoWindow
var _ editor.Editor = &Editor{}

// Editor represents a palette map editor
type Editor struct {
	*editor.EditorBase
	pl2           *d2pl2.PL2
	textureLoader common.TextureLoader
	state         []byte
}

// Create creates a new palette map editor
func Create(_ *config.Config,
	textureLoader common.TextureLoader,
	pathEntry *common.PathEntry,
	state []byte,
	data *[]byte, x, y float32, project *hsproject.Project,
) (editor.Editor, error) {
	pl2, err := d2pl2.Load(*data)
	if err != nil {
		return nil, fmt.Errorf("error loading PL2 file: %w", err)
	}

	result := &Editor{
		EditorBase:    editor.New(pathEntry, x, y, project),
		pl2:           pl2,
		textureLoader: textureLoader,
		state:         state,
	}

	result.Path = pathEntry

	return result, nil
}

// Build builds an editor
func (e *Editor) Build() {
	e.IsOpen(&e.Visible).
		Flags(g.WindowFlagsAlwaysAutoResize).
		Layout(g.Layout{
			palettemapwidget.Create(e.textureLoader, e.Path.GetUniqueID(), e.pl2, e.state),
		})
}

// UpdateMainMenuLayout updates a main menu layout to it contains editors options
func (e *Editor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("Palette Map Editor").Layout(g.Layout{
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

// GenerateSaveData creates data to be saved
func (e *Editor) GenerateSaveData() []byte {
	data := e.pl2.Marshal()

	return data
}

// Save saves an editor
func (e *Editor) Save() {
	e.EditorBase.Save(e)
}

// Cleanup hides an editor
func (e *Editor) Cleanup() {
	if e.HasChanges(e) {
		if shouldSave := dialog.Message("There are unsaved changes to %s, save before closing this editor?",
			e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.EditorBase.Cleanup()
}

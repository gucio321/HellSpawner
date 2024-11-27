// Package palette contains palette editor's data
package palette

import (
	"fmt"
	"github.com/gucio321/HellSpawner/pkg/app/config"

	"github.com/OpenDiablo2/dialog"

	g "github.com/AllenDang/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dat"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"

	"github.com/gucio321/HellSpawner/pkg/common"
	"github.com/gucio321/HellSpawner/pkg/common/hsproject"
	"github.com/gucio321/HellSpawner/pkg/widgets/palettegrideditorwidget"
	"github.com/gucio321/HellSpawner/pkg/widgets/palettegridwidget"
	"github.com/gucio321/HellSpawner/pkg/window/editor"
)

// static check, to ensure, if palette editor implemented editoWindow
var _ editor.Editor = &Editor{}

// Editor represents a palette editor
type Editor struct {
	*editor.EditorBase
	palette       d2interface.Palette
	textureLoader common.TextureLoader
	state         []byte
}

// Create creates a new palette editor
func Create(
	_ *config.Config,
	tl common.TextureLoader,
	pathEntry *common.PathEntry,
	state []byte,
	data *[]byte, x, y float32, project *hsproject.Project,
) (editor.Editor, error) {
	palette, err := d2dat.Load(*data)
	if err != nil {
		return nil, fmt.Errorf("error loading dat palette: %w", err)
	}

	result := &Editor{
		EditorBase:    editor.New(pathEntry, x, y, project),
		palette:       palette,
		textureLoader: tl,
		state:         state,
	}

	return result, nil
}

// Build builds a palette editor
func (e *Editor) Build() {
	const colorsPerPalette = 256

	col := make([]palettegridwidget.PaletteColor, colorsPerPalette)
	for n, i := range e.palette.GetColors() {
		col[n] = palettegridwidget.PaletteColor(i)
	}

	e.IsOpen(&e.Visible).Flags(g.WindowFlagsAlwaysAutoResize).Layout(g.Layout{
		palettegrideditorwidget.Create(e.state, e.textureLoader, e.GetID(), &col),
	})
}

// UpdateMainMenuLayout updates a main menu layout to it contain palette editor's options
func (e *Editor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("Palette Editor").Layout(g.Layout{
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
	palette, ok := e.palette.(*d2dat.DATPalette)
	if ok {
		data := palette.Marshal()

		return data
	}

	return make([]byte, 0)
}

// Save saves editor
func (e *Editor) Save() {
	e.EditorBase.Save(e)
}

// Cleanup hides palette editor
func (e *Editor) Cleanup() {
	if e.HasChanges(e) {
		if shouldSave := dialog.Message("There are unsaved changes to %s, save before closing this editor?",
			e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.EditorBase.Cleanup()
}

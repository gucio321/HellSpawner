// Package dc6 represents a dc6 editor window
package dc6

import (
	"fmt"
	"github.com/gucio321/HellSpawner/pkg/app/config"

	g "github.com/AllenDang/giu"
	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dc6"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"

	"github.com/gucio321/HellSpawner/pkg/common"
	"github.com/gucio321/HellSpawner/pkg/common/hsproject"
	"github.com/gucio321/HellSpawner/pkg/widgets/dc6widget"
	"github.com/gucio321/HellSpawner/pkg/widgets/selectpalettewidget"
	"github.com/gucio321/HellSpawner/pkg/window/editor"
)

// static check, to ensure, if dc6 editor implemented editoWindow
var _ editor.Editor = &Editor{}

// Editor represents a dc6 editor
type Editor struct {
	*editor.EditorBase
	dc6                 *d2dc6.DC6
	textureLoader       common.TextureLoader
	config              *config.Config
	selectPalette       bool
	palette             *[256]d2interface.Color
	selectPaletteWidget g.Widget
	state               []byte
}

// Create creates a new dc6 editor
func Create(cfg *config.Config,
	textureLoader common.TextureLoader,
	pathEntry *common.PathEntry,
	state []byte,
	data *[]byte, x, y float32, project *hsproject.Project,
) (editor.Editor, error) {
	dc6, err := d2dc6.Load(*data)
	if err != nil {
		return nil, fmt.Errorf("error loading DC6 animation: %w", err)
	}

	result := &Editor{
		EditorBase:    editor.New(pathEntry, x, y, project),
		dc6:           dc6,
		textureLoader: textureLoader,
		selectPalette: false,
		config:        cfg,
		state:         state,
	}

	return result, nil
}

// Build builds a new dc6 editor
func (e *Editor) Build() {
	e.IsOpen(&e.Visible)
	e.Flags(g.WindowFlagsAlwaysAutoResize)

	if !e.selectPalette {
		e.Layout(g.Layout{
			dc6widget.Create(e.state, e.palette, e.textureLoader, e.Path.GetUniqueID(), e.dc6),
		})

		return
	}

	if e.selectPaletteWidget == nil {
		e.selectPaletteWidget = selectpalettewidget.NewSelectPaletteWidget(
			e.Path.GetUniqueID()+"selectPalette",
			e.Project,
			e.config,
			func(palette *[256]d2interface.Color) {
				e.palette = palette
			},
			func() {
				e.selectPalette = false
			},
		)
	}

	e.Layout(e.selectPaletteWidget)
}

// UpdateMainMenuLayout updates main menu to it contain DC6's editor menu
func (e *Editor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("DC6 Editor").Layout(g.Layout{
		g.MenuItem("Change Palette").OnClick(func() {
			e.selectPalette = true
		}),
		g.Separator(),
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

// GenerateSaveData generates save data
func (e *Editor) GenerateSaveData() []byte {
	data := e.dc6.Marshal()

	return data
}

// Save saves editor's data
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

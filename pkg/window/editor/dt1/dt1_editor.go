// Package dt1 contains dt1 editor's data
package dt1

import (
	"fmt"

	"github.com/gucio321/HellSpawner/pkg/app/config"

	g "github.com/AllenDang/giu"

	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dt1"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"

	"github.com/gucio321/HellSpawner/pkg/common"
	"github.com/gucio321/HellSpawner/pkg/common/hsproject"
	"github.com/gucio321/HellSpawner/pkg/widgets/dt1widget"
	"github.com/gucio321/HellSpawner/pkg/widgets/selectpalettewidget"
	"github.com/gucio321/HellSpawner/pkg/window/editor"
)

// static check, to ensure, if dt1 editor implemented editoWindow
var _ editor.Editor = &Editor{}

// Editor represents a dt1 editor
type Editor struct {
	*editor.EditorBase
	dt1                 *d2dt1.DT1
	config              *config.Config
	selectPalette       bool
	palette             *[256]d2interface.Color
	selectPaletteWidget g.Widget
	state               []byte
}

// Create creates new dt1 editor
func Create(cfg *config.Config,
	pathEntry *common.PathEntry,
	state []byte,
	data *[]byte, x, y float32, project *hsproject.Project,
) (editor.Editor, error) {
	dt1, err := d2dt1.LoadDT1(*data)
	if err != nil {
		return nil, fmt.Errorf("error loading dt1 file: %w", err)
	}

	result := &Editor{
		EditorBase:    editor.New(pathEntry, x, y, project),
		dt1:           dt1,
		config:        cfg,
		selectPalette: false,
		state:         state,
	}

	return result, nil
}

// Build prepares the editor for rendering, but does not actually render it
func (e *Editor) Build() {
	e.IsOpen(&e.Visible).
		Flags(g.WindowFlagsAlwaysAutoResize).
		Layout(e.GetLayout())
}

func (e *Editor) GetLayout() g.Widget {
	if !e.selectPalette {
		return g.Layout{
			dt1widget.Create(
				e.state, e.palette,
				e.Path.GetUniqueID(),
				e.dt1,
			),
		}
	}

	// create mpq explorer if doesn't exist for now
	if e.selectPaletteWidget == nil {
		e.selectPaletteWidget = selectpalettewidget.NewSelectPaletteWidget(
			e.Path.GetUniqueID(),
			e.Project,
			e.config,
			func(colors *[256]d2interface.Color) {
				e.palette = colors
			},
			func() {
				e.selectPalette = false
			},
		)
	}

	return g.Layout{e.selectPaletteWidget}
}

// UpdateMainMenuLayout updates main menu layout to it contains editors options
func (e *Editor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("DT1 Editor").Layout(g.Layout{
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

// KeyboardShortcuts register a new keyboard shortcut
func (e *Editor) KeyboardShortcuts() []g.WindowShortcut {
	// https://github.com/gucio321/HellSpawner/issues/329
	return []g.WindowShortcut{
		/*
			// right arrow goes to the next tile group
			giu.WindowShortcut{
				Callback: func() {
					e.dt1Viewer.SetTileGroup(e.dt1Viewer.TileGroup() + 1)
				},
				Key:      g.KeyRight,
				Modifier: g.ModNone,
			},

			// left arrow goes to the previous tile group
			giu.WindowShortcut{
				Callback: func() {
					e.dt1Viewer.SetTileGroup(e.dt1Viewer.TileGroup() - 1)
				},
				Key:      g.KeyLeft,
				Modifier: g.ModNone,
			},
		*/
	}
}

// GenerateSaveData generates data to be saved
func (e *Editor) GenerateSaveData() []byte {
	data := e.dt1.Marshal()

	return data
}

// Save saves editor
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

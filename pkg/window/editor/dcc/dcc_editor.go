// Package dcc contains dcc editor's data
package dcc

import (
	"fmt"

	"github.com/gucio321/HellSpawner/pkg/app/config"

	g "github.com/AllenDang/giu"

	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dcc"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"

	"github.com/gucio321/HellSpawner/pkg/common"
	"github.com/gucio321/HellSpawner/pkg/common/hsproject"
	"github.com/gucio321/HellSpawner/pkg/widgets/dccwidget"
	"github.com/gucio321/HellSpawner/pkg/widgets/selectpalettewidget"
	"github.com/gucio321/HellSpawner/pkg/window/editor"
)

// static check, to ensure, if dc6 editor implemented editoWindow
var _ editor.Editor = &Editor{}

// Editor represents a new dcc editor
type Editor struct {
	*editor.EditorBase
	dcc                 *d2dcc.DCC
	config              *config.Config
	selectPalette       bool
	palette             *[256]d2interface.Color
	selectPaletteWidget g.Widget
	state               []byte
}

// Create creates a new dcc editor
func Create(cfg *config.Config,
	pathEntry *common.PathEntry,
	state []byte,
	data *[]byte, x, y float32, project *hsproject.Project,
) (editor.Editor, error) {
	dcc, err := d2dcc.Load(*data)
	if err != nil {
		return nil, fmt.Errorf("error loading dcc animation: %w", err)
	}

	result := &Editor{
		EditorBase:    editor.New(pathEntry, x, y, project),
		dcc:           dcc,
		config:        cfg,
		selectPalette: false,
		state:         state,
	}

	return result, nil
}

// Build builds a dcc editor
func (e *Editor) Build() {
	e.IsOpen(&e.Visible).
		Flags(g.WindowFlagsAlwaysAutoResize).
		Layout(e.GetLayout())
}

func (e *Editor) GetLayout() g.Widget {
	if !e.selectPalette {
		return g.Layout{
			dccwidget.Create(e.state, e.palette, e.Path.GetUniqueID(), e.dcc),
		}
	}

	if e.selectPaletteWidget == nil {
		e.selectPaletteWidget = selectpalettewidget.NewSelectPaletteWidget(
			"##"+e.Path.GetUniqueID()+"SelectPaletteWidget",
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

// UpdateMainMenuLayout updates main menu to it contain editor's options
func (e *Editor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("DCC Editor").Layout(g.Layout{
		g.MenuItem("Change Palette").OnClick(func() {
			e.selectPalette = true
		}),
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

// GenerateSaveData generates data to save
func (e *Editor) GenerateSaveData() []byte {
	// https://github.com/gucio321/HellSpawner/issues/181
	data, _ := e.Path.GetFileBytes()

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

// Package hsdcceditor contains dcc editor's data
package hsdcceditor

import (
	"fmt"

	g "github.com/AllenDang/giu"

	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dcc"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"

	"github.com/gucio321/HellSpawner/pkg/common"
	"github.com/gucio321/HellSpawner/pkg/common/hsproject"
	"github.com/gucio321/HellSpawner/pkg/config"
	"github.com/gucio321/HellSpawner/pkg/widgets/dccwidget"
	"github.com/gucio321/HellSpawner/pkg/widgets/selectpalettewidget"
	"github.com/gucio321/HellSpawner/pkg/window/hseditor"
)

// static check, to ensure, if dc6 editor implemented editoWindow
var _ common.EditorWindow = &DCCEditor{}

// DCCEditor represents a new dcc editor
type DCCEditor struct {
	*hseditor.Editor
	dcc                 *d2dcc.DCC
	config              *config.Config
	selectPalette       bool
	palette             *[256]d2interface.Color
	selectPaletteWidget g.Widget
	state               []byte
	textureLoader       common.TextureLoader
}

// Create creates a new dcc editor
func Create(config *config.Config,
	tl common.TextureLoader,
	pathEntry *common.PathEntry,
	state []byte,
	data *[]byte, x, y float32, project *hsproject.Project) (common.EditorWindow, error) {
	dcc, err := d2dcc.Load(*data)
	if err != nil {
		return nil, fmt.Errorf("error loading dcc animation: %w", err)
	}

	result := &DCCEditor{
		Editor:        hseditor.New(pathEntry, x, y, project),
		dcc:           dcc,
		config:        config,
		selectPalette: false,
		state:         state,
		textureLoader: tl,
	}

	return result, nil
}

// Build builds a dcc editor
func (e *DCCEditor) Build() {
	e.IsOpen(&e.Visible)
	e.Flags(g.WindowFlagsAlwaysAutoResize)

	if !e.selectPalette {
		e.Layout(g.Layout{
			dccwidget.Create(e.textureLoader, e.state, e.palette, e.Path.GetUniqueID(), e.dcc),
		})

		return
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

	e.Layout(g.Layout{e.selectPaletteWidget})
}

// UpdateMainMenuLayout updates main menu to it contain editor's options
func (e *DCCEditor) UpdateMainMenuLayout(l *g.Layout) {
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
func (e *DCCEditor) GenerateSaveData() []byte {
	// https://github.com/gucio321/HellSpawner/issues/181
	data, _ := e.Path.GetFileBytes()

	return data
}

// Save saves editor
func (e *DCCEditor) Save() {
	e.Editor.Save(e)
}

// Cleanup hides editor
func (e *DCCEditor) Cleanup() {
	if e.HasChanges(e) {
		if shouldSave := dialog.Message("There are unsaved changes to %s, save before closing this editor?",
			e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.Editor.Cleanup()
}

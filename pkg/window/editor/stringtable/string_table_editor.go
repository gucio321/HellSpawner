// Package stringtable contains string tables editor's data
package stringtable

import (
	"fmt"
	"github.com/gucio321/HellSpawner/pkg/app/config"

	g "github.com/AllenDang/giu"

	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2tbl"

	"github.com/gucio321/HellSpawner/pkg/common"
	"github.com/gucio321/HellSpawner/pkg/common/hsproject"
	"github.com/gucio321/HellSpawner/pkg/widgets/stringtablewidget"
	"github.com/gucio321/HellSpawner/pkg/window/editor"
)

const (
	mainWindowW, mainWindowH = 600, 500
)

// static check, to ensure, if string table editor implemented editoWindow
var _ editor.Editor = &Editor{}

// Editor represents a string table editor
type Editor struct {
	*editor.EditorBase
	dict  d2tbl.TextDictionary
	state []byte
}

// Create creates a new string table editor
func Create(_ *config.Config,
	_ common.TextureLoader,
	pathEntry *common.PathEntry,
	state []byte,
	data *[]byte, x, y float32, project *hsproject.Project,
) (editor.Editor, error) {
	dict, err := d2tbl.LoadTextDictionary(*data)
	if err != nil {
		return nil, fmt.Errorf("error loading string table: %w", err)
	}

	result := &Editor{
		EditorBase: editor.New(pathEntry, x, y, project),
		dict:       dict,
		state:      state,
	}

	if w, h := result.CurrentSize(); w == 0 || h == 0 {
		result.Size(mainWindowW, mainWindowH)
	}

	result.Path = pathEntry

	return result, nil
}

// Build builds an editor
func (e *Editor) Build() {
	l := stringtablewidget.Create(e.state, e.Path.GetUniqueID(), e.dict)

	e.IsOpen(&e.Visible).
		Flags(g.WindowFlagsHorizontalScrollbar).
		Layout(g.Layout{l})
}

// UpdateMainMenuLayout updates main menu layout to it contain editors options
func (e *Editor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("String Table Editor").Layout(g.Layout{
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
	data := e.dict.Marshal()

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

// Package fonttable represents fontTableEditor's window
package fonttable

import (
	"fmt"
	"github.com/gucio321/HellSpawner/pkg/app/config"

	"github.com/OpenDiablo2/dialog"

	g "github.com/AllenDang/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2font"

	"github.com/gucio321/HellSpawner/pkg/common"
	"github.com/gucio321/HellSpawner/pkg/common/hsproject"
	"github.com/gucio321/HellSpawner/pkg/widgets/fonttablewidget"
	"github.com/gucio321/HellSpawner/pkg/window/editor"
)

const (
	mainWindowW, mainWindowH = 550, 400
)

// static check, to ensure, if font table editor implemented editoWindow
var _ editor.Editor = &Editor{}

// Editor represents font table editor
type Editor struct {
	*editor.EditorBase
	fontTable *d2font.Font
	state     []byte
}

// Create creates a new font table editor
func Create(_ *config.Config,
	pathEntry *common.PathEntry,
	state []byte,
	data *[]byte, x, y float32, project *hsproject.Project,
) (editor.Editor, error) {
	table, err := d2font.Load(*data)
	if err != nil {
		return nil, fmt.Errorf("error font table: %w", err)
	}

	result := &Editor{
		EditorBase: editor.New(pathEntry, x, y, project),
		fontTable:  table,
		state:      state,
	}

	if w, h := result.CurrentSize(); w == 0 || h == 0 {
		result.Size(mainWindowW, mainWindowH)
	}

	return result, nil
}

// Build builds a font table editor's window
func (e *Editor) Build() {
	e.IsOpen(&e.Visible).Flags(g.WindowFlagsHorizontalScrollbar).
		Layout(g.Layout{
			fonttablewidget.Create(e.state, e.Path.GetUniqueID(), e.fontTable),
		})
}

// UpdateMainMenuLayout updates mainMenu layout's to it contain Editor's options
func (e *Editor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("Font Table Editor").Layout(g.Layout{
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
	data := e.fontTable.Marshal()

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

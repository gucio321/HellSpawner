// Package hsfonttableeditor represents fontTableEditor's window
package hsfonttableeditor

import (
	"fmt"

	"github.com/OpenDiablo2/dialog"

	g "github.com/AllenDang/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2font"

	"github.com/gucio321/HellSpawner/pkg/common"
	"github.com/gucio321/HellSpawner/pkg/common/hsproject"
	"github.com/gucio321/HellSpawner/pkg/config"
	"github.com/gucio321/HellSpawner/pkg/widgets/fonttablewidget"
	"github.com/gucio321/HellSpawner/pkg/window/hseditor"
)

const (
	mainWindowW, mainWindowH = 550, 400
)

// static check, to ensure, if font table editor implemented editoWindow
var _ common.EditorWindow = &FontTableEditor{}

// FontTableEditor represents font table editor
type FontTableEditor struct {
	*hseditor.Editor
	fontTable     *d2font.Font
	state         []byte
	textureLoader common.TextureLoader
}

// Create creates a new font table editor
func Create(_ *config.Config,
	tl common.TextureLoader,
	pathEntry *common.PathEntry,
	state []byte,
	data *[]byte, x, y float32, project *hsproject.Project) (common.EditorWindow, error) {
	table, err := d2font.Load(*data)
	if err != nil {
		return nil, fmt.Errorf("error font table: %w", err)
	}

	result := &FontTableEditor{
		Editor:        hseditor.New(pathEntry, x, y, project),
		fontTable:     table,
		state:         state,
		textureLoader: tl,
	}

	if w, h := result.CurrentSize(); w == 0 || h == 0 {
		result.Size(mainWindowW, mainWindowH)
	}

	return result, nil
}

// Build builds a font table editor's window
func (e *FontTableEditor) Build() {
	e.IsOpen(&e.Visible).Flags(g.WindowFlagsHorizontalScrollbar).
		Layout(g.Layout{
			fonttablewidget.Create(e.state, e.textureLoader, e.Path.GetUniqueID(), e.fontTable),
		})
}

// UpdateMainMenuLayout updates mainMenu layout's to it contain FontTableEditor's options
func (e *FontTableEditor) UpdateMainMenuLayout(l *g.Layout) {
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
func (e *FontTableEditor) GenerateSaveData() []byte {
	data := e.fontTable.Marshal()

	return data
}

// Save saves an editor
func (e *FontTableEditor) Save() {
	e.Editor.Save(e)
}

// Cleanup hides an editor
func (e *FontTableEditor) Cleanup() {
	if e.HasChanges(e) {
		if shouldSave := dialog.Message("There are unsaved changes to %s, save before closing this editor?",
			e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.Editor.Cleanup()
}
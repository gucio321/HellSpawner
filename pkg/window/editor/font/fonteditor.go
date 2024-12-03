// Package font contains font editor's data
package font

import (
	"fmt"

	"github.com/gucio321/HellSpawner/pkg/app/config"

	g "github.com/AllenDang/giu"
	"github.com/OpenDiablo2/dialog"

	"github.com/gucio321/HellSpawner/pkg/common/hsfiletypes/hsfont"
	"github.com/gucio321/HellSpawner/pkg/common/hsproject"

	"github.com/gucio321/HellSpawner/pkg/common"
	"github.com/gucio321/HellSpawner/pkg/window/editor"
)

const (
	mainWindowW, mainWindowH = 400, 300
	pathSize                 = 245
	browseW, browseH         = 30, 0
)

// static check, to ensure, if font editor implemented editoWindow
var _ editor.Editor = &Editor{}

// Editor represents a font editor
type Editor struct {
	*editor.EditorBase
	*hsfont.Font
}

// Create creates a new font editor
func Create(_ *config.Config,
	pathEntry *common.PathEntry,
	_ []byte,
	data *[]byte, x, y float32, project *hsproject.Project,
) (editor.Editor, error) {
	font, err := hsfont.LoadFromJSON(*data)
	if err != nil {
		return nil, fmt.Errorf("error loading JSON font: %w", err)
	}

	result := &Editor{
		EditorBase: editor.New(pathEntry, x, y, project),
		Font:       font,
	}

	if w, h := result.CurrentSize(); w == 0 || h == 0 {
		result.Size(mainWindowW, mainWindowH)
	}

	return result, nil
}

// Build builds an editor
func (e *Editor) Build() {
	e.IsOpen(&e.Visible).
		Layout(e.GetLayout())
}

func (e *Editor) GetLayout() g.Widget {
	return g.Layout{
		g.Label("DC6 Path"),
		g.Row(
			g.InputText(&e.SpriteFile).Size(pathSize).Flags(g.InputTextFlagsReadOnly),
			g.Button("...##EditorDC6Browse").Size(browseW, browseH).OnClick(e.onBrowseDC6PathClicked),
		),
		g.Separator(),
		g.Label("TBL Path"),
		g.Row(
			g.InputText(&e.TableFile).Size(pathSize).Flags(g.InputTextFlagsReadOnly),
			g.Button("...##EditorTBLBrowse").Size(browseW, browseH).OnClick(e.onBrowseTBLPathClicked),
		),
		g.Separator(),
		g.Label("PL2 Path"),
		g.Row(
			g.InputText(&e.PaletteFile).Size(pathSize).Flags(g.InputTextFlagsReadOnly),
			g.Button("...##EditorPL2Browse").Size(browseW, browseH).OnClick(e.onBrowsePL2PathClicked),
		),
	}
}

func (e *Editor) onBrowseDC6PathClicked() {
	path := dialog.File().SetStartDir(e.Project.GetProjectFileContentPath())
	path.Filter("DC6 File", "dc6", "DC6")

	filePath, err := path.Load()

	if err != nil || filePath == "" {
		return
	}

	e.SpriteFile = filePath
}

func (e *Editor) onBrowseTBLPathClicked() {
	path := dialog.File().SetStartDir(e.Project.GetProjectFileContentPath())
	path.Filter("TBL File", "tbl", "TBL")

	filePath, err := path.Load()

	if err != nil || filePath == "" {
		return
	}

	e.TableFile = filePath
}

func (e *Editor) onBrowsePL2PathClicked() {
	path := dialog.File().SetStartDir(e.Project.GetProjectFileContentPath())
	path.Filter("PL2 File", "pl2", "PL2")

	filePath, err := path.Load()

	if err != nil || filePath == "" {
		return
	}

	e.PaletteFile = filePath
}

// UpdateMainMenuLayout updates main menu layout to it contains editors options
func (e *Editor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("Font Editor").Layout(g.Layout{
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
	data, err := e.JSON()
	if err != nil {
		fmt.Println("failed to marshal font to JSON:, ", err)
		return nil
	}

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

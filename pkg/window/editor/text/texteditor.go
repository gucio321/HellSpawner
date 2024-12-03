// Package text contains text editor's data
package text

import (
	"log"
	"strings"

	"github.com/gucio321/HellSpawner/pkg/app/config"

	g "github.com/AllenDang/giu"

	"github.com/OpenDiablo2/dialog"

	"github.com/gucio321/HellSpawner/pkg/common"
	"github.com/gucio321/HellSpawner/pkg/common/hsproject"
	"github.com/gucio321/HellSpawner/pkg/window/editor"
)

const (
	mainWindowW, mainWindowH = 400, 300
	tableViewModW            = 80
	maxTableColumns          = 64
)

// static check, to ensure, if text editor implemented editoWindow
var _ editor.Editor = &Editor{}

// Editor represents a text editor
type Editor struct {
	*editor.EditorBase

	text      string
	tableView bool
	tableRows []*g.TableRowWidget
	columns   int
}

// Create creates a new text editor
func Create(_ *config.Config,
	pathEntry *common.PathEntry,
	_ []byte,
	data *[]byte, x, y float32, project *hsproject.Project,
) (editor.Editor, error) {
	result := &Editor{
		EditorBase: editor.New(pathEntry, x, y, project),
		text:       string(*data),
	}

	if w, h := result.CurrentSize(); w == 0 || h == 0 {
		result.Size(mainWindowW, mainWindowH)
	}

	lines := strings.Split(result.text, "\n")
	firstLine := lines[0]
	result.tableView = strings.Count(firstLine, "\t") > 0

	if !result.tableView {
		return result, nil
	}

	result.tableRows = make([]*g.TableRowWidget, len(lines))

	columns := strings.Split(firstLine, "\t")

	result.columns = len(columns)
	if result.columns > maxTableColumns {
		result.columns = maxTableColumns
		columns = columns[:maxTableColumns]

		log.Print("Waring: Table is too wide (more than 64 columns)! Only first 64 columns will be displayed" +
			"See: https://github.com/ocornut/imgui/issues/3572")
	}

	columnWidgets := make([]g.Widget, result.columns)

	for idx := range columns {
		columnWidgets[idx] = g.Label(columns[idx])
	}

	result.tableRows[0] = g.TableRow(columnWidgets...)

	for lineIdx := range lines[1:] {
		columns := strings.Split(lines[lineIdx+1], "\t")
		columnWidgets := make([]g.Widget, len(columns))

		for idx := range columns {
			columnWidgets[idx] = g.Label(columns[idx])
		}

		result.tableRows[lineIdx+1] = g.TableRow(columnWidgets...)
	}

	return result, nil
}

// Build builds an editor
func (e *Editor) Build() {
	e.IsOpen(&e.Visible).
		Flags(g.WindowFlagsHorizontalScrollbar).
		Layout(e.GetLayout())
}

func (e *Editor) GetLayout() g.Widget {
	if !e.tableView {
		return g.InputTextMultiline(&e.text).
			Flags(g.InputTextFlagsAllowTabInput)
	}

	return g.Child().Border(false).Size(float32(e.columns*tableViewModW), 0).Layout(
		g.Table().FastMode(true).Freeze(0, 1).Rows(e.tableRows...),
	)
}

// UpdateMainMenuLayout updates mainMenu layout to it contains editor's options
func (e *Editor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("Text Editor").Layout(g.Layout{
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
	data := []byte(e.text)

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

// Package animdata contains D2 editor's data
package animdata

import (
	"fmt"

	"github.com/gucio321/HellSpawner/pkg/app/config"

	g "github.com/AllenDang/giu"
	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2animdata"

	"github.com/gucio321/HellSpawner/pkg/common/hsproject"

	"github.com/gucio321/HellSpawner/pkg/common"
	"github.com/gucio321/HellSpawner/pkg/widgets/animdatawidget"
	"github.com/gucio321/HellSpawner/pkg/window/editor"
)

// static check, to ensure, if D2 editor implemented editoWindow
var _ editor.Editor = &AnimationDataEditor{}

// AnimationDataEditor represents a cof editor
type AnimationDataEditor struct {
	*editor.EditorBase
	d2    *d2animdata.AnimationData
	state []byte
}

// Create creates a new cof editor
func Create(_ *config.Config,
	pathEntry *common.PathEntry,
	state []byte,
	data *[]byte, x, y float32, project *hsproject.Project,
) (editor.Editor, error) {
	d2, err := d2animdata.Load(*data)
	if err != nil {
		return nil, fmt.Errorf("error loading animation data file: %w", err)
	}

	result := &AnimationDataEditor{
		EditorBase: editor.New(pathEntry, x, y, project),
		d2:         d2,
		state:      state,
	}

	return result, nil
}

// Build builds a D2 editor
func (e *AnimationDataEditor) Build() {
	e.IsOpen(&e.Visible).
		Flags(g.WindowFlagsAlwaysAutoResize).
		Layout(e.GetLayout())
}

func (e *AnimationDataEditor) GetLayout() g.Widget {
	uid := e.Path.GetUniqueID()
	return animdatawidget.Create(e.state, uid, e.d2)
}

// UpdateMainMenuLayout updates a main menu layout, to it contains anim data viewer's settings
func (e *AnimationDataEditor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("Animation Data Editor").Layout(g.Layout{
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
func (e *AnimationDataEditor) GenerateSaveData() []byte {
	data := e.d2.Marshal()

	return data
}

// Save saves an editor
func (e *AnimationDataEditor) Save() {
	e.EditorBase.Save(e)
}

// Cleanup hides an editor
func (e *AnimationDataEditor) Cleanup() {
	const strPrompt = "There are unsaved changes to %s, save before closing this editor?"

	if e.HasChanges(e) {
		if shouldSave := dialog.Message(strPrompt, e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.EditorBase.Cleanup()
}

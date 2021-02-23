// Package hscofeditor contains cof editor's data
package hsanimdataeditor

import (
	"github.com/OpenDiablo2/dialog"
	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2data"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswidget/animdatawidget"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

// AnimDataEditor represents a cof editor
type AnimDataEditor struct {
	*hseditor.Editor
	animData d2data.AnimationData
}

// Create creates a new cof editor
func Create(_ *hscommon.TextureLoader,
	pathEntry *hscommon.PathEntry,
	data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	animData, err := d2data.LoadAnimationData(*data)
	if err != nil {
		return nil, err
	}

	result := &AnimDataEditor{
		Editor:   hseditor.New(pathEntry, x, y, project),
		animData: animData,
	}

	return result, nil
}

// Build builds a cof editor
func (e *AnimDataEditor) Build() {
	e.IsOpen(&e.Visible).Flags(g.WindowFlagsAlwaysAutoResize).Layout(g.Layout{
		animdatawidget.AnimDataViewer(e.Path.GetUniqueID(), e.animData),
	})
}

// UpdateMainMenuLayout updates a main menu layout, to it contains COFViewer's settings
func (e *AnimDataEditor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("Animation Data Editor").Layout(g.Layout{
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
func (e *AnimDataEditor) GenerateSaveData() []byte {
	// https://github.com/OpenDiablo2/HellSpawner/issues/181
	data, _ := e.Path.GetFileBytes()

	return data
}

// Save saves an editor
func (e *AnimDataEditor) Save() {
	e.Editor.Save(e)
}

// Cleanup hides an editor
func (e *AnimDataEditor) Cleanup() {
	if e.HasChanges(e) {
		if shouldSave := dialog.Message("There are unsaved changes to %s, save before closing this editor?",
			e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.Editor.Cleanup()
}

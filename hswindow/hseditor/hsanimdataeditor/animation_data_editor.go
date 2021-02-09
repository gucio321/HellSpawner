package hsanimdataeditor

import (
	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2animdata"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hswidget"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

type AnimationDataEditor struct {
	*hseditor.Editor
	animData *d2animdata.AnimationData
}

func Create(_ *hscommon.TextureLoader,
	pathEntry *hscommon.PathEntry,
	data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	animData, err := d2animdata.Load(*data)
	if err != nil {
		return nil, err
	}

	result := &AnimationDataEditor{
		Editor:   hseditor.New(pathEntry, x, y, project),
		animData: animData,
	}

	return result, nil
}

func (e *AnimationDataEditor) Build() {
	e.IsOpen(&e.Visible).Flags(g.WindowFlagsAlwaysAutoResize).Layout(g.Layout{
		hswidget.AnimDataViewer(e.Path.GetUniqueID(), e.animData),
	})
}

func (e *AnimationDataEditor) UpdateMainMenuLayout(_ *g.Layout) {
}

// GenerateSaveData generates data to be saved
func (e *AnimationDataEditor) GenerateSaveData() []byte {
	// https://github.com/OpenDiablo2/HellSpawner/issues/181
	data, _ := e.Path.GetFileBytes()

	return data
}

// Save saves an editor
func (e *AnimationDataEditor) Save() {
	e.Editor.Save(e)
}

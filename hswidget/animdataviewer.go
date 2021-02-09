package hswidget

import (
	"fmt"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2animdata"
)

// AnimDataViewerState represents state ov animation data viewer
type AnimDataViewerState struct {
}

// Dispose clears viewer's state
func (s *AnimDataViewerState) Dispose() {
	// noop
}

type AnimDataViewerWidget struct {
	id       string
	animData *d2animdata.AnimationData
}

func AnimDataViewer(id string, animData *d2animdata.AnimationData) *AnimDataViewerWidget {
	result := &AnimDataViewerWidget{
		id:       id,
		animData: animData,
	}

	return result
}

func (p *AnimDataViewerWidget) Build() {
	stateID := fmt.Sprintf("AnimDataViewerWidget_%s", p.id)
	s := giu.Context.GetState(stateID)

	if s == nil {
		giu.Context.SetState(stateID, &AnimDataViewerState{})

		return
	}

	//state := s.(*AnimDataViewerState)

	records := p.animData.GetRecordNames()[0]

	giu.TabBar("AnimDataViewerTabs").Layout(giu.Layout{
		giu.TabItem("Records").Layout(giu.Layout{
			giu.Label(fmt.Sprintf("%v", p.animData.GetRecord(records).FPS())),
		}),
	}).Build()
}

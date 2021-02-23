package animdatawidget

import (
	"fmt"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2data"
)

// AnimDataViewerState represents state ov animation data viewer
type AnimDataViewerState struct {
	record int32
}

// Dispose clears viewer's state
func (s *AnimDataViewerState) Dispose() {
	// noop
}

type AnimDataViewerWidget struct {
	id       string
	animData d2data.AnimationData
}

func AnimDataViewer(id string, animData d2data.AnimationData) *AnimDataViewerWidget {
	result := &AnimDataViewerWidget{
		id:       id,
		animData: animData,
	}

	return result
}

func (p *AnimDataViewerWidget) Build() {
	state := p.getState()

	var recordsList []string = make([]string, 0)

	for i, _ := range p.animData {
		recordsList = append(recordsList, i)
	}

	fmt.Println(recordsList[0])

	giu.Layout{
		giu.Combo("##"+p.id+"recordsList", recordsList[state.record], recordsList, &state.record),
	}.Build()
}

func (p *AnimDataViewerWidget) getStateID() string {
	return fmt.Sprintf("AnimDataViewerWidget_%s", p.id)
}

func (p *AnimDataViewerWidget) getState() *AnimDataViewerState {
	var state *AnimDataViewerState

	s := giu.Context.GetState(p.getStateID())

	if s != nil {
		state = s.(*AnimDataViewerState)
	} else {
		p.initState()
		state = p.getState()
	}

	return state
}

func (p *AnimDataViewerWidget) initState() {
	state := &AnimDataViewerState{}

	p.setState(state)
}

func (p *AnimDataViewerWidget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}

package charactersavewidget

import (
	"fmt"

	"github.com/ianling/giu"
)

type widgetState struct {
	difficultyStatus int32
	questsDifficulty int32
	questsAct        int32
	questsIdx        int32
}

func (s *widgetState) Dispose() {
	// noop
}

func (p *widget) getStateID() string {
	// return fmt.Sprintf("widget_%s", p.id)
	return fmt.Sprintf("charSaveWidget_%s", p.id)
}

func (p *widget) getState() *widgetState {
	var state *widgetState

	s := giu.Context.GetState(p.getStateID())

	if s != nil {
		state = s.(*widgetState)
	} else {
		p.initState()
		state = p.getState()
	}

	return state
}

func (p *widget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}

func (p *widget) initState() {
	state := &widgetState{}

	p.setState(state)
}

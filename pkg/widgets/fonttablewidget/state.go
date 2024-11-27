package fonttablewidget

import (
	"fmt"
	"github.com/gucio321/HellSpawner/pkg/app/assets"

	"github.com/AllenDang/giu"
)

type widgetMode int32

const (
	modeViewer widgetMode = iota
	modeEditRune
	modeAddItem
)

type widgetState struct {
	Mode                widgetMode
	EditRuneState       editRuneState
	AddItemState        addItemState
	deleteButtonTexture *giu.Texture
}

// Dispose cleans state
func (s *widgetState) Dispose() {
	s.EditRuneState.Dispose()
	s.AddItemState.Dispose()
}

type editRuneState struct {
	EditedRune int32
	RuneBefore rune
}

// Dispose disposes a rune state
func (e *editRuneState) Dispose() {
	e.EditedRune = rune(0)
	e.RuneBefore = rune(0)
}

type addItemState struct {
	NewRune,
	Width,
	Height int32
}

func (s *addItemState) Dispose() {
	s.NewRune = rune(0)
	s.Height = 0
	s.Width = 0
}

func (p *widget) getStateID() giu.ID {
	return giu.ID(fmt.Sprintf("widget_%s", p.id))
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

func (p *widget) initState() {
	state := &widgetState{
		Mode: modeViewer,
	}

	p.textureLoader.CreateTextureFromFile(assets.DeleteIcon, func(texture *giu.Texture) {
		state.deleteButtonTexture = texture
	})

	p.setState(state)
}

func (p *widget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}

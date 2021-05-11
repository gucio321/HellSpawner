package dccwidget

import (
	"fmt"
	"log"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dcc"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswidget/animationwidget"
)

const (
	maxAlpha = uint8(255)
)

type widget struct {
	id            string
	dcc           *d2dcc.DCC
	palette       *[256]d2interface.Color
	textureLoader hscommon.TextureLoader
}

// Create creates a new dcc widget
func Create(tl hscommon.TextureLoader, state []byte, palette *[256]d2interface.Color, id string, dcc *d2dcc.DCC) giu.Widget {
	result := &widget{
		id:            id,
		dcc:           dcc,
		palette:       palette,
		textureLoader: tl,
	}

	if giu.Context.GetState(result.getStateID()) == nil && state != nil {
		s := result.getState()
		s.Decode(state)
		result.setState(s)
	}

	return result
}

// Build build a widget
func (p *widget) Build() {
	viewerState := p.getState()

	err := giu.Context.GetRenderer().SetTextureMagFilter(giu.TextureFilterNearest)
	if err != nil {
		log.Print(err)
	}

	giu.Layout{
		giu.Line(
			giu.Label(fmt.Sprintf("Signature: %v", p.dcc.Signature)),
			giu.Label(fmt.Sprintf("Version: %v", p.dcc.Version)),
		),
		giu.Line(
			giu.Label(fmt.Sprintf("Directions: %v", p.dcc.NumberOfDirections)),
			giu.Label(fmt.Sprintf("Frames per Direction: %v", p.dcc.FramesPerDirection)),
		),
		animationwidget.Create(p.id+"widget", viewerState.images, p.dcc.FramesPerDirection, p.dcc.NumberOfDirections, p.textureLoader),
	}.Build()
}

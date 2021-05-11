package dc6widget

import (
	"fmt"
	"log"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dc6"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswidget/animationwidget"
)

const (
	inputIntW        = 30
	buttonW, buttonH = 200, 30
)

const (
	maxAlpha = uint8(255)
)

// widget represents dc6viewer's widget
type widget struct {
	id            string
	dc6           *d2dc6.DC6
	textureLoader hscommon.TextureLoader
	palette       *[256]d2interface.Color
}

// Create creates new widget
func Create(state []byte, palette *[256]d2interface.Color, textureLoader hscommon.TextureLoader, id string, dc6 *d2dc6.DC6) giu.Widget {
	result := &widget{
		id:            id,
		dc6:           dc6,
		textureLoader: textureLoader,
		palette:       palette,
	}

	if giu.Context.GetState(result.getStateID()) == nil && state != nil {
		s := result.getState()
		s.Decode(state)

		if s.mode == dc6WidgetTiledView {
			result.createImage(s)
		}

		result.setState(s)
	}

	return result
}

// TiledView switches widget into tiled view mode
func (p *widget) TiledView() {
	state := p.getState()
	state.mode = dc6WidgetTiledView
}

// Build builds a widget
func (p *widget) Build() {
	state := p.getState()

	switch state.mode {
	case dc6WidgetViewer:
		p.makeViewerLayout().Build()
	case dc6WidgetTiledView:
		p.makeTiledViewLayout(state).Build()
	}
}

func (p *widget) makeViewerLayout() giu.Layout {
	viewerState := p.getState()

	err := giu.Context.GetRenderer().SetTextureMagFilter(giu.TextureFilterNearest)
	if err != nil {
		log.Print(err)
	}

	return giu.Layout{
		giu.Label(fmt.Sprintf(
			"Version: %v\t Flags: %b\t Encoding: %v\t",
			p.dc6.Version,
			int64(p.dc6.Flags),
			p.dc6.Encoding,
		)),
		giu.Label(fmt.Sprintf("Directions: %v\tFrames per Direction: %v", p.dc6.Directions, p.dc6.FramesPerDirection)),
		giu.Separator(),
		animationwidget.Create(p.id+"widget", viewerState.rgb, int(p.dc6.FramesPerDirection), int(p.dc6.Directions), p.textureLoader),
	}
}

func (p *widget) makeTiledViewLayout(state *widgetState) giu.Layout {
	return giu.Layout{
		giu.Line(
			giu.Label("Tiled view:"),
			giu.InputInt("Width##"+p.id+"tiledWidth", &state.width).Size(inputIntW).OnChange(func() {
				p.recalculateTiledViewHeight(state)
			}),
			giu.InputInt("Height##"+p.id+"tiledHeight", &state.height).Size(inputIntW).OnChange(func() {
				p.recalculateTiledViewWidth(state)
			}),
		),
		giu.Image(state.tiled).Size(float32(state.imgw), float32(state.imgh)),
		giu.Button("Back##"+p.id+"tiledBack").Size(buttonW, buttonH).OnClick(func() {
			state.mode = dc6WidgetViewer
		}),
	}
}

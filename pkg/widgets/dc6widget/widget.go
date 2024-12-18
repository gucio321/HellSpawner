package dc6widget

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/giu"
	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dc6"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"

	"github.com/gucio321/HellSpawner/pkg/common/hsutil"
	"github.com/gucio321/HellSpawner/pkg/widgets"
)

const (
	comboW              = 125
	inputIntW           = 30
	playPauseButtonSize = 15
	buttonW, buttonH    = 200, 30
)

const (
	maxAlpha = uint8(255)
)

// widget represents dc6viewer's widget
type widget struct {
	id      string
	dc6     *d2dc6.DC6
	palette *[256]d2interface.Color
}

// Create creates new widget
func Create(state []byte, palette *[256]d2interface.Color, id string, dc6 *d2dc6.DC6) giu.Widget {
	result := &widget{
		id:      id,
		dc6:     dc6,
		palette: palette,
	}

	if giu.Context.GetState(result.getStateID()) == nil && state != nil {
		s := result.getState()

		if err := json.Unmarshal(state, s); err != nil {
			log.Printf("error decoding dc6 widget state: %v", err)
		}

		s.ticker.Reset(time.Second * time.Duration(s.TickTime) / miliseconds)

		if s.Mode == dc6WidgetTiledView {
			result.createImage(s)
		}

		result.setState(s)
	}

	return result
}

// Build builds a widget
func (p *widget) Build() {
	state := p.getState()

	switch state.Mode {
	case dc6WidgetViewer:
		p.makeViewerLayout().Build()
	case dc6WidgetTiledView:
		p.makeTiledViewLayout(state).Build()
	}
}

func (p *widget) makeViewerLayout() giu.Layout {
	viewerState := p.getState()

	//nolint:gosec // we need to cast this here
	imageScale := uint32(viewerState.Controls.Scale)
	curFrameIndex := int(viewerState.Controls.Frame) + (int(viewerState.Controls.Direction) * int(p.dc6.FramesPerDirection))
	dirIdx := int(viewerState.Controls.Direction)

	textureIdx := dirIdx*int(p.dc6.FramesPerDirection) + int(viewerState.Controls.Frame)

	if imageScale < 1 {
		imageScale = 1
	}

	// TODO: doesn't work on latest giu
	/*
		err := giu.Context.GetRenderer().SetTextureMagFilter(giu.TextureFilterNearest)
		if err != nil {
			log.Print(err)
		}
	*/

	w := float32(p.dc6.Frames[curFrameIndex].Width * imageScale)
	h := float32(p.dc6.Frames[curFrameIndex].Height * imageScale)

	var widget *giu.ImageWidget
	if viewerState.textures == nil || len(viewerState.textures) <= int(viewerState.Controls.Frame) ||
		viewerState.textures[curFrameIndex] == nil {
		widget = giu.Image(nil).Size(w, h)
	} else {
		widget = giu.Image(viewerState.textures[textureIdx]).Size(w, h)
	}

	return giu.Layout{
		giu.Label(fmt.Sprintf(
			"Version: %v\t Flags: %b\t Encoding: %v\t",
			p.dc6.Version,
			int64(p.dc6.Flags),
			p.dc6.Encoding,
		)),
		giu.Label(fmt.Sprintf("Directions: %v\tFrames per Direction: %v", p.dc6.Directions, p.dc6.FramesPerDirection)),
		giu.Custom(func() {
			imgui.BeginGroup()
			if p.dc6.Directions > 1 {
				//nolint:gosec // imgui magic here.
				imgui.SliderInt("Direction", &viewerState.Controls.Direction, 0, int32(p.dc6.Directions-1))
			}

			if p.dc6.FramesPerDirection > 1 {
				//nolint:gosec // imgui magic here.
				imgui.SliderInt("Frames", &viewerState.Controls.Frame, 0, int32(p.dc6.FramesPerDirection-1))
			}

			const minScale, maxScale = 1, 8

			imgui.SliderInt("Scale", &viewerState.Controls.Scale, minScale, maxScale)

			imgui.EndGroup()
		}),
		giu.Separator(),
		p.makePlayerLayout(viewerState),
		giu.Separator(),
		widget,
		giu.Separator(),
		giu.Button("Tiled View##"+p.id+"tiledViewButton").Size(buttonW, buttonH).OnClick(func() {
			viewerState.Mode = dc6WidgetTiledView
			p.createImage(viewerState)
		}),
	}
}

func (p *widget) makePlayerLayout(state *widgetState) giu.Layout {
	playModeList := make([]string, 0)
	for i := playModeForward; i <= playModePingPong; i++ {
		playModeList = append(playModeList, i.String())
	}

	pm := int32(state.PlayMode)

	return giu.Layout{
		giu.Row(
			giu.Checkbox("Loop##"+p.id+"PlayRepeat", &state.Repeat),
			giu.Combo("##"+p.id+"PlayModeList", playModeList[state.PlayMode], playModeList, &pm).OnChange(func() {
				state.PlayMode = animationPlayMode(pm)
			}).Size(comboW),
			giu.InputInt(&state.TickTime).Label("Tick time").Size(inputIntW).OnChange(func() {
				state.ticker.Reset(time.Second * time.Duration(state.TickTime) / miliseconds)
			}),
			widgets.PlayPauseButton(&state.IsPlaying).
				Size(playPauseButtonSize, playPauseButtonSize),
			giu.Button("Export GIF##"+p.id+"exportGif").OnClick(func() {
				err := p.exportGif(state)
				if err != nil {
					dialog.Message("%v", err).Error()
				}
			}),
			giu.Button("Export Frames (PNG)##"+p.id+"exportpng").OnClick(func() {
				err := p.exportPng(state)
				if err != nil {
					dialog.Message("%v", err).Error()
				}
			}),
		),
	}
}

func (p *widget) makeTiledViewLayout(state *widgetState) giu.Layout {
	return giu.Layout{
		giu.Row(
			giu.Label("Tiled view:"),
			giu.InputInt(&state.Width).Label("Width").Size(inputIntW).OnChange(func() {
				p.recalculateTiledViewHeight(state)
			}),
			giu.InputInt(&state.Height).Label("Height").Size(inputIntW).OnChange(func() {
				p.recalculateTiledViewWidth(state)
			}),
		),
		giu.Image(state.tiled).Size(float32(state.Imgw), float32(state.Imgh)),
		giu.Button("Back##"+p.id+"tiledBack").Size(buttonW, buttonH).OnClick(func() {
			state.Mode = dc6WidgetViewer
		}),
	}
}

func (p *widget) exportGif(state *widgetState) error {
	//nolint:gosec // we need to cast this here
	fpd := int32(p.dc6.FramesPerDirection)
	firstFrame := state.Controls.Direction * fpd
	images := state.rgb[firstFrame : firstFrame+fpd]

	err := hsutil.ExportToGif(images, state.TickTime)
	if err != nil {
		return fmt.Errorf("error creating gif file: %w", err)
	}

	return nil
}

func (p *widget) exportPng(state *widgetState) error {
	images := state.rgb

	err := hsutil.ExportToPng(images)
	if err != nil {
		return fmt.Errorf("error creating png file: %w", err)
	}

	return nil
}

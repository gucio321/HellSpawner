package animationwidget

import (
	"fmt"
	"image"
	"log"
	"time"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsutil"
	"github.com/OpenDiablo2/HellSpawner/hswidget"
	"github.com/OpenDiablo2/dialog"
	"github.com/ianling/giu"
	"github.com/ianling/imgui-go"
)

const (
	inputIntW           = 30
	playPauseButtonSize = 15
	comboW              = 125
	miliseconds         = 1000
	imageW, imageH      = 32, 32
	defaultTickTime     = 100
)

type widget struct {
	images        []*image.RGBA
	fpd           int
	id            string
	numDirs       int
	textureLoader hscommon.TextureLoader
}

func Create(id string, images []*image.RGBA, fpd, numDirs int, tl hscommon.TextureLoader) giu.Widget {
	result := &widget{
		id:            id,
		images:        images,
		fpd:           fpd,
		numDirs:       numDirs,
		textureLoader: tl,
	}

	return result
}

func (p *widget) Build() {
	state := p.getState()

	imageScale := uint32(state.controls.scale)
	dirIdx := int(state.controls.direction)
	frameIdx := state.controls.frame

	textureIdx := dirIdx*p.fpd + int(frameIdx)

	if imageScale < 1 {
		imageScale = 1
	}

	err := giu.Context.GetRenderer().SetTextureMagFilter(giu.TextureFilterNearest)
	if err != nil {
		log.Print(err)
	}

	var widget *giu.ImageWidget
	if state.textures == nil || len(state.textures) <= int(frameIdx) || state.textures[frameIdx] == nil {
		widget = giu.Image(nil).Size(imageW, imageH)
	} else {
		/*bw := p.dcc.Directions[dirIdx].Box.Width
		bh := p.dcc.Directions[dirIdx].Box.Height
		w := float32(uint32(bw) * imageScale)
		h := float32(uint32(bh) * imageScale)*/
		b := p.images[textureIdx].Bounds()
		w, h := b.Dx()*int(imageScale), b.Dy()*int(imageScale)
		widget = giu.Image(state.textures[textureIdx]).Size(float32(w), float32(h))
	}

	giu.Layout{
		giu.Custom(func() {
			imgui.BeginGroup()
			if p.numDirs > 1 {
				imgui.SliderInt("Direction", &state.controls.direction, 0, int32(p.numDirs-1))
			}

			if p.fpd > 1 {
				imgui.SliderInt("Frames", &state.controls.frame, 0, int32(p.fpd-1))
			}

			imgui.SliderInt("Scale", &state.controls.scale, 1, 8)

			imgui.EndGroup()
		}),
		giu.Separator(),
		p.makePlayerLayout(state),
		giu.Separator(),
		widget,
	}.Build()
}

func (p *widget) makePlayerLayout(state *widgetState) giu.Layout {
	playModeList := make([]string, 0)
	for i := playModeForward; i <= playModePingPong; i++ {
		playModeList = append(playModeList, i.String())
	}

	pm := int32(state.playMode)

	return giu.Layout{
		giu.Line(
			giu.Checkbox("Loop##"+p.id+"PlayRepeat", &state.repeat),
			giu.Combo("##"+p.id+"PlayModeList", playModeList[state.playMode], playModeList, &pm).OnChange(func() {
				state.playMode = animationPlayMode(pm)
			}).Size(comboW),
			giu.InputInt("Tick time##"+p.id+"PlayTickTime", &state.tickTime).Size(inputIntW).OnChange(func() {
				state.ticker.Reset(time.Second * time.Duration(state.tickTime) / miliseconds)
			}),
			hswidget.PlayPauseButton("##"+p.id+"PlayPauseAnimation", &state.isPlaying, p.textureLoader).
				Size(playPauseButtonSize, playPauseButtonSize),
			giu.Button("Export GIF##"+p.id+"exportGif").OnClick(func() {
				err := p.exportGif(state)
				if err != nil {
					dialog.Message(err.Error()).Error()
				}
			}),
		),
	}
}

func (p *widget) exportGif(state *widgetState) error {
	firstFrame := state.controls.direction * int32(p.fpd)
	images := p.images[firstFrame : firstFrame+int32(p.fpd)]

	err := hsutil.ExportToGif(images, state.tickTime)
	if err != nil {
		return fmt.Errorf("error creating gif file: %w", err)
	}

	return nil
}

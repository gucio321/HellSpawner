package dccwidget

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/AllenDang/giu"

	"github.com/gucio321/HellSpawner/pkg/common/hsutil"
)

const miliseconds = 1000

type animationPlayMode byte

const (
	playModeForward animationPlayMode = iota
	playModeBackword
	playModePingPong
)

func (a animationPlayMode) String() string {
	s := map[animationPlayMode]string{
		playModeForward:  "Forwards",
		playModeBackword: "Backwords",
		playModePingPong: "Ping-Pong",
	}

	k, ok := s[a]
	if !ok {
		return "Unknown"
	}

	return k
}

const defaultTickTime = 100

type widgetState struct {
	Controls struct {
		Direction int32
		Frame     int32
		Scale     int32
	}

	IsPlaying bool
	Repeat    bool
	TickTime  int32
	PlayMode  animationPlayMode

	// cache - will not be saved
	images   []*image.RGBA
	textures []*giu.Texture

	isForward bool // determines a direction of animation
	ticker    *time.Ticker
}

// Dispose cleans viewers state
func (s *widgetState) Dispose() {
	s.textures = nil
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
	// Prevent multiple invocation to LoadImage.
	state := &widgetState{
		IsPlaying: false,
		Repeat:    false,
		TickTime:  defaultTickTime,
		PlayMode:  playModeForward,
	}

	state.ticker = time.NewTicker(time.Second * time.Duration(state.TickTime) / miliseconds)

	p.setState(state)

	go p.runPlayer(state)

	totalFrames := p.dcc.NumberOfDirections * p.dcc.FramesPerDirection
	state.images = make([]*image.RGBA, totalFrames)

	for dirIdx := range p.dcc.Directions {
		fw := p.dcc.Directions[dirIdx].Box.Width
		fh := p.dcc.Directions[dirIdx].Box.Height

		for frameIdx := range p.dcc.Directions[dirIdx].Frames {
			absoluteFrameIdx := (dirIdx * p.dcc.FramesPerDirection) + frameIdx

			frame := p.dcc.Directions[dirIdx].Frames[frameIdx]
			pixels := frame.PixelData

			state.images[absoluteFrameIdx] = image.NewRGBA(image.Rect(0, 0, fw, fh))

			for y := 0; y < fh; y++ {
				for x := 0; x < fw; x++ {
					idx := x + (y * fw)
					if idx >= len(pixels) {
						continue
					}

					val := pixels[idx]

					RGBAColor := p.makeImagePixel(val)
					state.images[absoluteFrameIdx].Set(x, y, RGBAColor)
				}
			}
		}
	}

	go func() {
		textures := make([]*giu.Texture, totalFrames)

		for frameIndex := 0; frameIndex < totalFrames; frameIndex++ {
			frameIndex := frameIndex
			giu.EnqueueNewTextureFromRgba(state.images[frameIndex], func(t *giu.Texture) {
				textures[frameIndex] = t
			})
		}

		s := p.getState()
		s.textures = textures
		p.setState(s)
	}()
}

func (p *widget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}

func (p *widget) makeImagePixel(val byte) color.RGBA {
	alpha := maxAlpha

	if val == 0 {
		alpha = 0
	}

	var r, g, b uint8

	if p.palette != nil {
		col := p.palette[val]
		r, g, b = col.R(), col.G(), col.B()
	} else {
		r, g, b = val, val, val
	}

	RGBAColor := color.RGBA{
		R: r,
		G: g,
		B: b,
		A: alpha,
	}

	return RGBAColor
}

func (p *widget) runPlayer(state *widgetState) {
	for range state.ticker.C {
		if !state.IsPlaying {
			continue
		}

		//nolint:mnd // this comes from od2 and we have no control over that. On the other hand we need to pass it to giu which supports int32 only.
		numFrames := int32(p.dcc.FramesPerDirection - 1)
		isLastFrame := state.Controls.Frame == numFrames

		// update play direction
		switch state.PlayMode {
		case playModeForward:
			state.isForward = true
		case playModeBackword:
			state.isForward = false
		case playModePingPong:
			if isLastFrame || state.Controls.Frame == 0 {
				state.isForward = !state.isForward
			}
		}

		// now update the frame number
		if state.isForward {
			state.Controls.Frame++
		} else {
			state.Controls.Frame--
		}

		//nolint:gomnd // this comes from od2 and needs to be passed to giu.
		state.Controls.Frame = hsutil.Wrap(state.Controls.Frame, int32(p.dcc.FramesPerDirection))

		// next, check for stopping/repeat
		isStoppingFrame := (state.Controls.Frame == 0) || (state.Controls.Frame == numFrames)

		if isStoppingFrame && !state.Repeat {
			state.IsPlaying = false
		}
	}
}

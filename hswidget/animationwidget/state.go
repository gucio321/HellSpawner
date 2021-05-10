package animationwidget

import (
	"fmt"
	"time"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsutil"
	"github.com/ianling/giu"
)

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

type widgetState struct {
	controls struct {
		direction int32
		frame     int32
		scale     int32
	}

	isPlaying bool
	repeat    bool
	tickTime  int32
	playMode  animationPlayMode

	// cache
	textures  []*giu.Texture
	isForward bool // determines a direction of animation
	ticker    *time.Ticker
}

func (s *widgetState) Dispose() {
	// noop
}

func (p *widget) getStateID() string {
	return fmt.Sprintf("widget_%s", p.id)
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
	state := &widgetState{
		isPlaying: false,
		repeat:    false,
		tickTime:  defaultTickTime,
		playMode:  playModeForward,
	}

	state.ticker = time.NewTicker(time.Second * time.Duration(state.tickTime) / miliseconds)

	go p.runPlayer(state)

	totalFrames := p.numDirs * p.fpd
	go func() {
		textures := make([]*giu.Texture, totalFrames)

		for frameIndex := 0; frameIndex < totalFrames; frameIndex++ {
			frameIndex := frameIndex
			p.textureLoader.CreateTextureFromARGB(p.images[frameIndex], func(t *giu.Texture) {
				textures[frameIndex] = t
			})
		}

		s := p.getState()
		s.textures = textures
		p.setState(s)
	}()

	p.setState(state)
}

func (p *widget) runPlayer(state *widgetState) {
	for range state.ticker.C {
		if !state.isPlaying {
			continue
		}

		numFrames := int32(p.fpd - 1)
		isLastFrame := state.controls.frame == numFrames

		// update play direction
		switch state.playMode {
		case playModeForward:
			state.isForward = true
		case playModeBackword:
			state.isForward = false
		case playModePingPong:
			if isLastFrame || state.controls.frame == 0 {
				state.isForward = !state.isForward
			}
		}

		// now update the frame number
		if state.isForward {
			state.controls.frame++
		} else {
			state.controls.frame--
		}

		state.controls.frame = int32(hsutil.Wrap(int(state.controls.frame), p.fpd))

		// next, check for stopping/repeat
		isStoppingFrame := (state.controls.frame == 0) || (state.controls.frame == numFrames)

		if isStoppingFrame && !state.repeat {
			state.isPlaying = false
		}
	}
}

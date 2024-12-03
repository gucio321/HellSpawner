package widgets

import (
	"fmt"
	"math"

	"github.com/gucio321/HellSpawner/pkg/app/assets"

	"github.com/AllenDang/giu"

	"github.com/gucio321/HellSpawner/pkg/common"
)

// MakeImageButton is a hack for giu.ImageButton that creates image button
// as a giu.child
func MakeImageButton(w, h int, t *giu.Texture, fn func()) giu.Widget {
	// the image button
	btnW, btnH := float32(w), float32(h)
	return giu.ImageButton(t).Size(btnW, btnH).OnClick(fn)
}

type playPauseButtonState struct {
	playTexture,
	pauseTexture *giu.Texture
}

func (s *playPauseButtonState) Dispose() {
	s.playTexture = nil
	s.pauseTexture = nil
}

// PlayPauseButtonWidget represents a play/pause button
type PlayPauseButtonWidget struct {
	id giu.ID

	onChange,
	onPauseClicked,
	onPlayClicked func()

	width,
	height float32

	isPlaying *bool
}

// PlayPauseButton creates a play/pause button
func PlayPauseButton(isPlaying *bool) *PlayPauseButtonWidget {
	return &PlayPauseButtonWidget{
		id:        giu.GenAutoID("##PlayPauseButton"),
		isPlaying: isPlaying,
	}
}

// ID is for compatibility with giu widgets. It allows to manually set widget's id.
func (p *PlayPauseButtonWidget) ID(id giu.ID) *PlayPauseButtonWidget {
	p.id = id
	return p
}

// Size sets button's size
func (p *PlayPauseButtonWidget) Size(w, h float32) *PlayPauseButtonWidget {
	p.width, p.height = w, h
	return p
}

// OnPlayClicked sets onPlayClicked callback (called when the user clicks on play button)
func (p *PlayPauseButtonWidget) OnPlayClicked(cb func()) *PlayPauseButtonWidget {
	p.onPlayClicked = cb
	return p
}

// OnPauseClicked sets onPauseClicked callback (called when the user clicks on pause button)
func (p *PlayPauseButtonWidget) OnPauseClicked(cb func()) *PlayPauseButtonWidget {
	p.onPauseClicked = cb
	return p
}

// OnChange sets onChange callback (called the user click on any button)
func (p *PlayPauseButtonWidget) OnChange(cb func()) *PlayPauseButtonWidget {
	p.onChange = cb
	return p
}

// Build build a widget
func (p *PlayPauseButtonWidget) Build() {
	stateID := giu.ID(fmt.Sprintf("%s_state", p.id))
	state := giu.Context.GetState(stateID)

	var widget giu.Widget

	if state == nil {
		widget = giu.Image(nil).Size(p.width, p.height)

		state := &playPauseButtonState{}

		common.LoadTexture(assets.PlayButtonIcon, func(t *giu.Texture) {
			state.playTexture = t
		})

		common.LoadTexture(assets.PauseButtonIcon, func(t *giu.Texture) {
			state.pauseTexture = t
		})

		giu.Context.SetState(stateID, state)

		widget.Build()

		return
	}

	imgState := state.(*playPauseButtonState)

	w, h := int(p.width), int(p.height)

	var texture *giu.Texture

	var callback func() // callback

	setIsPlaying := func(b bool) {
		*p.isPlaying = b

		if cb := p.onChange; cb != nil {
			cb()
		}

		if cb := p.onPlayClicked; cb != nil {
			cb()
		}
	}

	if !*p.isPlaying {
		texture = imgState.playTexture
		callback = func() { setIsPlaying(true) }
	} else {
		texture = imgState.pauseTexture
		callback = func() { setIsPlaying(false) }
	}

	widget = MakeImageButton(w, h, texture, callback)

	widget.Build()
}

// SetByteToInt sets byte given to intager
// if intager > max possible byte size, sets to 255
func SetByteToInt(input int32, output *byte) {
	if input > int32(math.MaxUint8) {
		*output = math.MaxUint8

		return
	}

	*output = byte(input)
}

// MakeInputInt creates input intager using POINTER given
// additionally, for byte checks, if value smaller than 255
func MakeInputInt(width int32, output interface{}, optionalCB func()) *giu.InputIntWidget {
	var input int32
	switch o := output.(type) {
	case *byte:
		input = int32(*o)
	case *int:
		input = int32(*o) //nolint:gosec // lets suppose this is safe
	default:
		panic(fmt.Sprintf("MakeInputInt: invalid value type %T given", o))
	}

	return giu.InputInt(&input).Size(float32(width)).OnChange(func() {
		switch o := output.(type) {
		case *byte:
			SetByteToInt(input, o)
		case *int:
			*o = int(input)
		}

		if optionalCB != nil {
			optionalCB()
		}
	})
}

// MakeCheckboxFromByte creates a checkbox using a byte as input/output
func MakeCheckboxFromByte(id giu.ID, value *byte) *giu.CheckboxWidget {
	v := *value > 0

	return giu.Checkbox("", &v).ID(id).OnChange(func() {
		if v {
			*value = 1
		} else {
			*value = 0
		}
	})
}

// OnDoubleClick detects if item is double-clicked
// this can be used as an alternative to OnClick methos of some widgets
// e.g.:
//
//	giu.Layout{
//		giu.Button("double click me"),
//		OnDoubleClick(func() { fmt.Println("I was double-clicked") }),
//	}
func OnDoubleClick(cb func()) giu.Widget {
	return giu.Custom(func() {
		if giu.IsItemHovered() && giu.IsMouseDoubleClicked(giu.MouseButtonLeft) {
			go cb()
		}
	})
}

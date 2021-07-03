package palettegrideditorwidget

import (
	"github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsutil"
	"github.com/OpenDiablo2/HellSpawner/hswidget/palettegridwidget"
)

const (
	actionButtonW, actionButtonH = 250, 30
)

// PaletteGridEditorWidget represents a palette grid editor
type PaletteGridEditorWidget struct {
	id            string
	colors        *[]palettegridwidget.PaletteColor
	textureLoader hscommon.TextureLoader
	onChange      func()
}

// Create creates a new palette grid editor widget
func Create(state []byte,
	textureLoader hscommon.TextureLoader,
	id string,
	colors *[]palettegridwidget.PaletteColor) *PaletteGridEditorWidget {
	result := &PaletteGridEditorWidget{
		id:            id,
		colors:        colors,
		textureLoader: textureLoader,
		onChange:      nil,
	}

	if giu.Context.GetState(result.getStateID()) == nil && state != nil {
		s := result.getState()
		s.Decode(state)
		result.setState(s)
	}

	return result
}

// OnChange sets on change callback
// this callback is ran, when editor's slider or field gets change
func (p *PaletteGridEditorWidget) OnChange(onChange func()) *PaletteGridEditorWidget {
	p.onChange = onChange
	return p
}

// Build Builds a widget
func (p *PaletteGridEditorWidget) Build() {
	state := p.getState()

	colors := make([]palettegridwidget.PaletteColor, len(*p.colors))
	for n := range *(p.colors) {
		colors[n] = (*p.colors)[n]
	}

	grid := palettegridwidget.Create(p.textureLoader, p.id, &colors).OnClick(func(idx int) {
		color := hsutil.Color((*p.colors)[idx].RGBA())
		state.rgba = color
		state.idx = idx

		state.mode = widgetModeEdit
	})

	switch state.mode {
	case widgetModeGrid:
		grid.Build()
	case widgetModeEdit:
		p.buildEditor()
	}
}

func (p *PaletteGridEditorWidget) buildEditor() {
	state := p.getState()

	isOpen := state.mode == widgetModeEdit
	onChange := func() {
		p.changeColor(state)

		if p.onChange != nil {
			p.onChange()
		}
	}

	onClick := func() {
		onChange()

		state.mode = widgetModeGrid
	}

	const (
		popupTitle = "Edit color"
	)

	colorEditID := "##edit color"
	buttonID := "OK##" + p.id + "editColorOK"

	colorEditor := giu.PopupModal(popupTitle).IsOpen(&isOpen).Layout(
		giu.ColorEdit(colorEditID, &state.rgba).Flags(giu.ColorEditFlagsNoAlpha),
		giu.Separator(),
		giu.Button(buttonID).Size(actionButtonW, actionButtonH).OnClick(onClick),
	)

	// handle clicking on "X" button of popup
	closeButtonHandler := giu.Custom(func() {
		if !isOpen {
			onChange()
			state.mode = widgetModeGrid
		}
	})

	giu.Layout{colorEditor, closeButtonHandler}.Build()
}

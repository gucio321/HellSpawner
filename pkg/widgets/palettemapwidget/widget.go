package palettemapwidget

import (
	"encoding/json"
	"log"

	"github.com/AllenDang/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2pl2"

	"github.com/gucio321/HellSpawner/pkg/common/hsutil"
	"github.com/gucio321/HellSpawner/pkg/widgets/palettegrideditorwidget"
	"github.com/gucio321/HellSpawner/pkg/widgets/palettegridwidget"
)

const (
	comboW             = 280
	layoutW, layoutH   = 475, 300
	actionButtonW      = layoutW
	numColorsInPalette = 256
)

type widget struct {
	id  string
	pl2 *d2pl2.PL2
}

// Create creates a new palette map viewer's widget
func Create(id string, pl2 *d2pl2.PL2, state []byte) giu.Widget {
	result := &widget{
		id:  id,
		pl2: pl2,
	}

	if giu.Context.GetState(result.getStateID()) == nil && state != nil {
		s := result.getState()
		if err := json.Unmarshal(state, s); err != nil {
			log.Printf("error decoding palette map widget state: %v", err)
		}
	}

	return result
}

// Build builds a new widget
func (p *widget) Build() {
	state := p.getState()

	switch state.Mode {
	case widgetModeView:
		p.buildViewer(state)
	case widgetModeEditTransform:
		p.buildEditor(state)
	}
}

func (p *widget) buildViewer(state *widgetState) {
	// TODO: this is disabled in giu since cimgui-go migration
	/*
		err := giu.Context.GetRenderer().SetTextureMagFilter(giu.TextureFilterNearest)
		if err != nil {
			log.Print(err)
		}
	*/

	baseColors := make([]palettegridwidget.PaletteColor, numColorsInPalette)

	for n := range baseColors {
		baseColors[n] = palettegridwidget.PaletteColor(&p.pl2.BasePalette.Colors[n])
	}

	left := giu.Layout{
		giu.Label("Base Palette"),
		palettegrideditorwidget.Create(nil, p.id+"basePalette", &baseColors).OnChange(func() {
			state.textures = make(map[string]giu.Widget)
		}),
	}

	selections := getPaletteTransformString()
	right := giu.Layout{
		giu.Label("Palette Map"),
		giu.Layout{
			giu.Combo("", selections[state.Selection], selections, &state.Selection).Size(comboW),
			p.getTransformViewLayout(state.Selection),
		},
	}

	w1, h1 := float32(layoutW), float32(layoutH)
	w2, h2 := float32(layoutW), float32(layoutH)

	// nolint:gomnd // special case for alpha blend
	if state.Selection == 3 {
		h2 += 32
	}

	layout := giu.Layout{
		giu.Child().Size(w1, h1).Layout(left),
		giu.Child().Size(w2, h2).Layout(right),
	}

	layout.Build()
}

func (p *widget) getTransformViewLayout(transformIdx int32) giu.Layout {
	buildLayout := []func() giu.Layout{
		func() giu.Layout {
			return p.transformMulti("LightLevelVariations", p.pl2.LightLevelVariations[:])
		},
		func() giu.Layout {
			return p.transformMulti("InvColorVariations", p.pl2.InvColorVariations[:])
		},
		func() giu.Layout {
			return p.transformSingle("SelectedUintShift", &p.pl2.SelectedUintShift.Indices)
		},
		func() giu.Layout {
			return p.transformMultiGroup("AlphaBlend", p.pl2.AlphaBlend[:]...)
		},
		func() giu.Layout {
			return p.transformMulti("AdditiveBlend", p.pl2.AdditiveBlend[:])
		},
		func() giu.Layout {
			return p.transformMulti("MultiplicativeBlend", p.pl2.MultiplicativeBlend[:])
		},
		func() giu.Layout {
			return p.transformMulti("HueVariations", p.pl2.HueVariations[:])
		},
		func() giu.Layout {
			return p.transformSingle("RedTones", &p.pl2.RedTones.Indices)
		},
		func() giu.Layout {
			return p.transformSingle("GreenTones", &p.pl2.GreenTones.Indices)
		},
		func() giu.Layout {
			return p.transformSingle("BlueTones", &p.pl2.BlueTones.Indices)
		},
		func() giu.Layout {
			return p.transformMulti("UnknownVariations", p.pl2.UnknownVariations[:])
		},
		func() giu.Layout {
			return p.transformMulti("MaxComponentBlend", p.pl2.MaxComponentBlend[:])
		},
		func() giu.Layout {
			return p.transformSingle("DarkendColorShift", &p.pl2.DarkendColorShift.Indices)
		},
		func() giu.Layout {
			return p.textColors("TextColors", p.pl2.TextColors[:])
		},
		func() giu.Layout {
			return p.transformMulti("TextColorShifts", p.pl2.TextColorShifts[:])
		},
	}

	return buildLayout[transformIdx]()
}

func (p *widget) buildEditor(state *widgetState) {
	var grid giu.Widget

	indices := p.getPaletteIndices(state)

	colors := make([]palettegridwidget.PaletteColor, len(p.pl2.BasePalette.Colors))

	for n := range colors {
		colors[n] = palettegridwidget.PaletteColor(&p.pl2.BasePalette.Colors[n])
	}

	grid = palettegridwidget.Create(p.id+"transformEdit", &colors).OnClick(func(idx int) {
		// this is save, because idx is always less than 256
		indices[state.Idx] = byte(idx)

		// reset textures list
		state.textures = make(map[string]giu.Widget)

		state.Mode = widgetModeView
	})
	labelColor := hsutil.Color(p.pl2.BasePalette.Colors[indices[state.Idx]].RGBA())
	giu.Layout{
		giu.Style().SetColor(giu.StyleColorText, labelColor).To(
			giu.Label("Select color from base palette"),
		),
		grid,
		giu.Separator(),
		// if height > 0, then pushItemHeight
		giu.Button("Cancel##"+p.id+"cancelEditorButton").Size(actionButtonW, 0).OnClick(func() {
			state.Mode = widgetModeView
		}),
	}.Build()
}

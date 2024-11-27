package palettegridwidget

import (
	"image"

	"github.com/AllenDang/giu"
)

const (
	gridWidth  = 16
	gridHeight = 16
	cellSize   = 12
)

// PaletteGridWidget represents a palette grid
type PaletteGridWidget struct {
	id      string
	colors  *[]PaletteColor
	onClick func(idx int)
}

// Create creates a new palette grid widget
func Create(id string, colors *[]PaletteColor) *PaletteGridWidget {
	result := &PaletteGridWidget{
		id:      id,
		colors:  colors,
		onClick: nil,
	}

	return result
}

// OnClick sets onClick callback
func (p *PaletteGridWidget) OnClick(onClick func(idx int)) *PaletteGridWidget {
	p.onClick = onClick
	return p
}

// UpdateImage updates a palette image.
// should be called when palete colors gets changed
func (p *PaletteGridWidget) UpdateImage() {
	p.rebuildImage()
}

// Build build a new widget
func (p *PaletteGridWidget) Build() {
	state := p.getState()

	// cache variable for a base position of image
	var imgBase image.Point

	giu.Layout{
		// just save base cursor position
		giu.Custom(func() {
			imgBase = giu.GetCursorScreenPos()
		}),
		giu.Image(state.rgba).
			Size(gridWidth*cellSize, gridHeight*cellSize),
		// event detector - detects clicking in a cell
		giu.Custom(func() {
			mousePos := giu.GetMousePos()

			// x, y - cursor position on an image
			x := mousePos.X - imgBase.X
			y := mousePos.Y - imgBase.Y

			// cellX, cellY - cell cords
			cellX, cellY := x/cellSize, y/cellSize

			// check if cell cords are out of bounds
			if cellX < 0 || cellY < 0 || cellX >= gridWidth || cellY >= gridHeight {
				return
			}

			idx := cellY*gridHeight + cellX

			if giu.IsWindowFocused(giu.FocusedFlags(giu.FocusedFlagsNone)) && giu.IsMouseClicked(giu.MouseButtonLeft) {
				p.onClick(idx)
				p.rebuildImage()
			}
		}),
	}.Build()
}

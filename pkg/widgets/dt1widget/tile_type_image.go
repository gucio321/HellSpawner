package dt1widget

import (
	"github.com/AllenDang/giu"

	"github.com/gucio321/HellSpawner/pkg/widgets/dt1widget/tiletypeimage"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
)

func drawTileTypeImage(t d2enum.TileType) giu.Widget {
	return giu.Custom(func() {
		canvas := giu.GetCanvas()
		pos := giu.GetCursorScreenPos()
		b := tiletypeimage.TileTypeImage(canvas, pos)
		lookup := map[d2enum.TileType]func(){
			d2enum.TileFloor:                      func() { b.Floor() },
			d2enum.TileLeftWall:                   func() { b.Floor().WestWall(true) },
			d2enum.TileRightWall:                  func() { b.Floor().NorthWall(true) },
			d2enum.TileRightPartOfNorthCornerWall: func() { b.Floor().WestWall(false).NorthWall(true) },
			d2enum.TileLeftPartOfNorthCornerWall:  func() { b.Floor().WestWall(true).NorthWall(false) },
			d2enum.TileLeftEndWall:                func() { b.Floor().EastWall() },
			d2enum.TileRightEndWall:               func() { b.Floor().SouthWall() },
			d2enum.TileSouthCornerWall:            func() { b.Floor().Corner() },
			d2enum.TileLeftWallWithDoor:           func() { b.Floor().WestDoor() },
			d2enum.TileRightWallWithDoor:          func() { b.Floor().NorthDoor() },
		}

		if creator, ok := lookup[t]; ok {
			creator()
		}
	})
}
